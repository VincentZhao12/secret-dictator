package main

import (
	"fmt"
	"net/http"

	"github.com/VincentZhao12/secret-hitler/backend/internal/game"
	"github.com/VincentZhao12/secret-hitler/backend/internal/routes"
)

func main() {
	m := game.NewManager()
	r := routes.SetupRouter(m)
	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", r)
}
