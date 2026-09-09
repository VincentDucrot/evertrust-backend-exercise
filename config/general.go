package config

import (
	"log"
	"os"
	"time"
)

type GeneralConfig struct {
	Port          string
	TokenDuration time.Duration // in minutes
}

const defaultPort = "8080"
const defaultTokenDuration = 10 * time.Minute

func LoadGeneralConfig() GeneralConfig {
	port := os.Getenv("PORT")
	tokenDurationString := os.Getenv("TOKEN_DURATION")
	var tokenDuration time.Duration
	var err error

	if port == "" {
		port = defaultPort
	}
	if tokenDurationString == "" {
		tokenDuration = defaultTokenDuration
	} else {
		tokenDuration, err = time.ParseDuration(tokenDurationString + "m")
		if err != nil {
			log.Println("error parsing token duration. Default value will be used.")
			tokenDuration = defaultTokenDuration
		}
	}

	return GeneralConfig{
		Port:          port,
		TokenDuration: tokenDuration,
	}
}
