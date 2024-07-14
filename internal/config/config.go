package config

import (
	"time"

	"github.com/spf13/viper"
)

type MongoDBConfig struct {
	URI    string
	DBName string
}

type OktaConfig struct {
	ClientID     string
	ClientSecret string
	Issuer       string
	Nonce        string
}

type JWTConfig struct {
	SecretKey    string
	ExpiredAfter time.Duration
}

type Config struct {
	MongoDB       *MongoDBConfig
	OktaConfig    *OktaConfig
	JWTConfig     *JWTConfig
	SessionConfig string
}

func NewConfig() *Config {
	return &Config{
		MongoDB: &MongoDBConfig{
			URI:    viper.GetString("MONGODB_URI"),
			DBName: viper.GetString("MONGODB_DB"),
		},
		OktaConfig: &OktaConfig{
			ClientID:     viper.GetString("OKTA_CLIENT_ID"),
			ClientSecret: viper.GetString("OKTA_CLIENT_SECRET"),
			Issuer:       viper.GetString("OKTA_ISSUER"),
			Nonce:        viper.GetString("OKTA_NONCE"),
		},
		JWTConfig: &JWTConfig{
			SecretKey:    viper.GetString("JWT_SECRET_KEY"),
			ExpiredAfter: viper.GetDuration("JWT_EXPIRED_AFTER"),
		},
		SessionConfig: viper.GetString("SESSION_SECRET"),
	}
}
