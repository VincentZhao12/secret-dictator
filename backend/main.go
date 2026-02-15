package main

import (
	"fmt"
	"net/http"

	"github.com/VincentZhao12/secret-hitler/backend/internal/bot"
	"github.com/VincentZhao12/secret-hitler/backend/internal/envs"
	"github.com/VincentZhao12/secret-hitler/backend/internal/game"
	"github.com/VincentZhao12/secret-hitler/backend/internal/routes"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func main() {
	env := envs.GetEnv()
	aiClient := openai.NewClient(
		option.WithAPIKey(envs.OpenRouterAPIKey()),
		option.WithBaseURL(envs.OpenRouterBaseURL()),
	)
	botProvider := bot.NewOpenRouterProvider(&aiClient, "openai/gpt-oss-120b:free")
	m := game.NewManager(botProvider)
	r := routes.SetupRouter(m, env)
	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", r)
}
