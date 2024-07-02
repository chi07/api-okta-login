package main

import (
	"context"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/chi07/api-okta-login/internal/config"
	"github.com/chi07/api-okta-login/internal/http/handler"
	"github.com/chi07/api-okta-login/internal/repository"
	"github.com/chi07/api-okta-login/internal/service"
	"github.com/go-playground/validator/v10"
)

func main() {
	viper.AutomaticEnv()
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "It's work")
	})

	cnf := config.NewConfig()
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	var validate *validator.Validate
	validate = validator.New()

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cnf.MongoDB.URI))
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to MongoDB")
	}

	db := client.Database(cnf.MongoDB.DBName)
	userRepo := repository.NewUser(db)
	loginHistoryRepo := repository.NewLoginHistory(db)
	oktaService := service.NewOktaService(&logger, cnf.OktaConfig, cnf.JWTConfig, userRepo, loginHistoryRepo)
	loginHandler := handler.NewLoginHandler(cnf, logger, validate, oktaService)

	e.POST("/login", loginHandler.LoginWithOkta)

	e.Logger.Fatal(e.Start(":1323"))
}
