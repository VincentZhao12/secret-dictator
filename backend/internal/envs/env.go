package envs

import (
	"fmt"
	"os"
)

const DefaultOpenRouterBaseURL = "https://openrouter.ai/api/v1"

type Env string

const (
	Development Env = "development"
	Production  Env = "production"
)

func GetEnv() Env {
	fmt.Println("ENV:", os.Getenv("ENV"))
	if os.Getenv("ENV") == "production" {
		return Production
	}
	return Development
}

func OpenRouterAPIKey() string {
	return os.Getenv("OPENROUTER_API_KEY")
}

func OpenRouterBaseURL() string {
	if url := os.Getenv("OPENROUTER_BASE_URL"); url != "" {
		return url
	}
	return DefaultOpenRouterBaseURL
}
