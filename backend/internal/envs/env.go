package envs

import (
	"fmt"
	"os"
)

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
