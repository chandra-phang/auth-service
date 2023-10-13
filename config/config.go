package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

var config *Config

type Config struct {
	OauthConfig *oauth2.Config
}

func InitConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config = &Config{
		OauthConfig: &oauth2.Config{
			ClientID:     os.Getenv("CLIENT_ID"),
			ClientSecret: os.Getenv("CLIENT_SECRET"),
			RedirectURL:  os.Getenv("REDIRECT_URI"),
			Endpoint: oauth2.Endpoint{
				AuthURL:  os.Getenv("AUTH_URL"),
				TokenURL: os.Getenv("TOKEN_URL"),
			},
			Scopes: []string{"openid", "profile", "email"},
		},
	}
}

func GetConfig() *Config {
	return config
}
