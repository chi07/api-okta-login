package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/okta/okta-jwt-verifier-golang/v2"
	"github.com/rs/zerolog"

	"github.com/chi07/api-okta-login/internal/config"
	"github.com/chi07/api-okta-login/internal/model"
)

type OktaService struct {
	log              *zerolog.Logger
	oktaConfig       *config.OktaConfig
	jwtConfig        *config.JWTConfig
	userRepo         UserRepo
	sessionStore     *sessions.CookieStore
	loginHistoryRepo LoginHistoryRepo
	client           *resty.Client
}

func NewOktaService(log *zerolog.Logger, oktaConfig *config.OktaConfig, jwtConfig *config.JWTConfig, userRepo UserRepo, loginHistoryRepo LoginHistoryRepo, sessionStore *sessions.CookieStore, client *resty.Client) *OktaService {
	return &OktaService{
		log:              log,
		oktaConfig:       oktaConfig,
		jwtConfig:        jwtConfig,
		userRepo:         userRepo,
		loginHistoryRepo: loginHistoryRepo,
		sessionStore:     sessionStore,
		client:           client,
	}
}

func (s *OktaService) Login(ctx echo.Context, oktaToken string) (map[string]interface{}, error) {
	// @TODO: verify token and do something with the token
	claims, ok, err := s.VerifyToken(ctx.Request().Context(), oktaToken)
	if err != nil {
		s.log.Error().Err(err).Msg("s.VerifyToken() got err: " + err.Error())
		return nil, err
	}

	if !ok {
		return nil, fmt.Errorf("token is invalid")
	}

	// Test if the token is valid
	userOktaID, ok := claims.Claims["uid"].(string)
	email, ok := claims.Claims["sub"].(string)
	if !ok {
		s.log.Error().Msg("email claim not found")
		return nil, fmt.Errorf("email claim not found")
	}

	// @TODO get userinfo from Okta
	var oktaUserInfo model.OktaUserInfo
	oktaURL := fmt.Sprintf("%s/v1/userinfo", s.oktaConfig.Issuer)
	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+oktaToken).
		SetHeader("Accept", "application/json").
		SetResult(&oktaUserInfo).
		Get(oktaURL)

	if err != nil {
		s.log.Error().Err(err).Msg("client.R().Get() got err: " + err.Error())
		return nil, err
	}

	if resp.StatusCode() != 200 {
		s.log.Error().Msg("failed to get userinfo from Okta")
		return nil, fmt.Errorf("failed to get userinfo from Okta")
	}

	// @TODO: check if the user is already registered in the database
	user, err := s.CheckUser(ctx.Request().Context(), userOktaID, oktaUserInfo.Email, oktaUserInfo.Name)
	if err != nil {
		s.log.Error().Err(err).Msg("s.CheckUser() got err: " + err.Error())
		return nil, err
	}

	// Create a new login history
	go func() {
		_, err = s.loginHistoryRepo.Create(context.Background(), &model.LoginHistory{
			OktaID:    user.OktaID,
			Email:     user.Email,
			Status:    model.StatusSuccess,
			CreatedAt: time.Now(),
		})
		if err != nil {
			s.log.Error().Err(err).Msg("s.loginHistoryRepo.Create() got err: " + err.Error())
		}
	}()

	// @TODO: generate a new token for the user
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":    user.OktaID,
		"oktaID": user.OktaID,
		"email":  email,
		"exp":    time.Now().Add(time.Second * s.jwtConfig.ExpiredAfter).Unix(),
		"iat":    time.Now().Unix(),
	})

	token, err := jwtToken.SignedString([]byte(s.jwtConfig.SecretKey))
	if err != nil {
		s.log.Error().Err(err).Msg("jwtToken.SignedString() got err: " + err.Error())
		return nil, err
	}

	// @TODO: save the token to the session
	session, err := s.sessionStore.Get(ctx.Request(), "currentUser")
	if err != nil {
		s.log.Error().Err(err).Msg("s.sessionStore.Get() got err: " + err.Error())
		return nil, err
	}

	session.Values["accessToken"] = token
	if err = session.Save(ctx.Request(), ctx.Response()); err != nil {
		s.log.Error().Err(err).Msg("session.Save() got err: " + err.Error())
		return nil, err
	}

	return map[string]interface{}{
		"accessToken": token,
	}, nil
}

func (s *OktaService) VerifyToken(ctx context.Context, requestToken string) (*jwtverifier.Jwt, bool, error) {
	toValidate := map[string]string{}
	toValidate["nonce"] = s.oktaConfig.Nonce
	toValidate["aud"] = "api://default" // Get this from Okta application

	jwtVerifierSetup := jwtverifier.JwtVerifier{
		Issuer:           s.oktaConfig.Issuer,
		ClaimsToValidate: toValidate,
	}

	verifier, err := jwtVerifierSetup.New()
	if err != nil {
		s.log.Error().Err(err).Msg("jwtVerifierSetup.New() got err: " + err.Error())
		return nil, false, nil

	}

	token, err := verifier.VerifyAccessToken(requestToken)
	if err != nil {
		s.log.Error().Err(err).Msg("verifier.VerifyAccessToken() got err: " + err.Error())
		return nil, false, nil
	}

	return token, true, nil
}

func (s *OktaService) CheckUser(ctx context.Context, oktaID, email, name string) (*model.User, error) {
	user, err := s.userRepo.GetByOktaID(ctx, oktaID)
	if err != nil {
		s.log.Error().Err(err).Msg("s.userRepo.GetByOktaID() got err: " + err.Error())
		return nil, err
	}

	// if user is not found, create a new user
	if user == nil {
		now := time.Now()
		user = &model.User{
			Name:      name,
			OktaID:    oktaID,
			Email:     email,
			Status:    "active",
			CreatedAt: now,
			UpdatedAt: now,
		}

		user, err = s.userRepo.Create(ctx, user)
		if err != nil {
			s.log.Error().Err(err).Msg("s.userRepo.Create() got err: " + err.Error())
			return nil, err
		}
	}

	return user, nil
}

func (s *OktaService) Logout(ctx echo.Context) error {
	// @TODO: remove the token from the session
	session, err := s.sessionStore.Get(ctx.Request(), "currentUser")
	if err != nil {
		s.log.Error().Err(err).Msg("s.sessionStore.Get() got err: " + err.Error())
		return err
	}

	// set accessToken to empty
	session.Values["accessToken"] = ""
	if err = session.Save(ctx.Request(), ctx.Response()); err != nil {
		s.log.Error().Err(err).Msg("session.Save() got err: " + err.Error())
		return err
	}

	return nil
}
