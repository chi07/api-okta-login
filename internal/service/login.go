package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	loginHistoryRepo LoginHistoryRepo
}

func NewOktaService(log *zerolog.Logger, oktaConfig *config.OktaConfig, jwtConfig *config.JWTConfig, userRepo UserRepo, loginHistoryRepo LoginHistoryRepo) *OktaService {
	return &OktaService{
		log:              log,
		oktaConfig:       oktaConfig,
		jwtConfig:        jwtConfig,
		userRepo:         userRepo,
		loginHistoryRepo: loginHistoryRepo,
	}
}

func (s *OktaService) Login(ctx context.Context, oktaToken string) (map[string]interface{}, error) {
	// @TODO: verify token and do something with the token
	claims, ok, err := s.VerifyToken(ctx, oktaToken)
	if err != nil {
		s.log.Error().Err(err).Msg("s.VerifyToken() got err: " + err.Error())
		return nil, err
	}

	if !ok {
		return nil, fmt.Errorf("token is invalid")
	}

	// Test if the token is valid
	username, ok := claims.Claims["sub"].(string)
	if !ok {
		s.log.Error().Msg("sub claim not found")
		return nil, fmt.Errorf("sub claim not found")
	}

	email, ok := claims.Claims["email"].(string)
	if !ok {
		s.log.Error().Msg("email claim not found")
		return nil, fmt.Errorf("email claim not found")
	}

	name, ok := claims.Claims["name"].(string)
	if !ok {
		s.log.Error().Msg("name claim not found")
		return nil, fmt.Errorf("name claim not found")
	}

	// @TODO: check if the user is already registered in the database
	user, err := s.CheckUser(ctx, username, email, name)
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
		"email":  user.Email,
		"exp":    time.Now().Add(time.Second * s.jwtConfig.ExpiredAfter).Unix(),
		"iat":    time.Now().Unix(),
	})

	token, err := jwtToken.SignedString([]byte(s.jwtConfig.SecretKey))
	if err != nil {
		s.log.Error().Err(err).Msg("jwtToken.SignedString() got err: " + err.Error())
		return nil, err
	}

	return map[string]interface{}{
		"accessToken": token,
	}, nil
}

func (s *OktaService) VerifyToken(ctx context.Context, requestToken string) (*jwtverifier.Jwt, bool, error) {
	toValidate := map[string]string{}
	toValidate["nonce"] = s.oktaConfig.Nonce
	toValidate["aud"] = s.oktaConfig.ClientID

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
