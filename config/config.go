package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RedditClientID     string
	RedditClientSecret string
	UserAgent          string

	// token info
	InitialToken string
	ExpiresIn    int

	SubReddit string
	Port      int
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	redditClientID := os.Getenv("REDDIT_CLIENT_ID")
	redditClientSecret := os.Getenv("REDDIT_CLIENT_SECRET")
	userAgent := os.Getenv("USER_AGENT")

	if redditClientID == "" || redditClientSecret == "" {
		return nil, errors.New("environment variables REDDIT_CLIENT_ID and REDDIT_CLIENT_SECRET must be set")
	}

	return &Config{
		RedditClientID:     redditClientID,
		RedditClientSecret: redditClientSecret,
		UserAgent:          userAgent,
	}, nil

}
