package routes

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/VincentZhao12/secret-hitler/backend/internal/envs"
	"github.com/VincentZhao12/secret-hitler/backend/internal/game"
	"github.com/VincentZhao12/secret-hitler/backend/internal/handlers"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func SetupRouter(m *game.Manager, env envs.Env) http.Handler {
	r := chi.NewRouter()

	if env == envs.Development {
		fmt.Println("CORS enabled for development")
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:8080"},
			AllowedMethods:   []string{"GET", "POST"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Upgrade", "Connection"},
			AllowCredentials: true,
			ExposedHeaders:   []string{"Upgrade", "Connection"},
		}))
	}

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api/v1", func(api chi.Router) {
		api.Post("/games/create", handlers.CreateGame(m))
		api.Post("/games/join", handlers.JoinGame(m))
		api.Post("/games/bots/add", handlers.AddBot(m))
		api.Post("/games/bots/remove", handlers.RemoveBot(m))
		api.Get("/play", handlers.Play(m))
	})

	// Serve static files from web/dist
	workDir, _ := filepath.Abs("./static")
	filesDir := http.Dir(workDir)
	FileServer(r, "/", filesDir)

	return r
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))

		// Check if the file exists
		filePath := strings.TrimPrefix(r.URL.Path, pathPrefix)
		if filePath == "" {
			filePath = "/"
		}

		// Try to open the file
		file, err := root.Open(filePath)
		if err != nil {
			// File doesn't exist, serve index.html for SPA routing
			indexFile, err := root.Open("/index.html")
			if err != nil {
				http.NotFound(w, r)
				return
			}
			defer indexFile.Close()

			w.Header().Set("Content-Type", "text/html")
			http.ServeContent(w, r, "index.html", time.Time{}, indexFile)
			return
		}
		defer file.Close()

		// File exists, serve it normally
		fs.ServeHTTP(w, r)
	})
}
