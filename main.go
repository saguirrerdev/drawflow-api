package main

import (
	"context"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	node "github.com/zebek95/draflow-api/controllers"
	"github.com/zebek95/draflow-api/database"
	"github.com/zebek95/draflow-api/error_handler"
	"github.com/zebek95/draflow-api/models"
)

func init() {
	os.Setenv("ENVIRONMENT", "development")
	os.Setenv("API_PORT", "3333")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "9080")

	database.Init()
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://*"},
		AllowedMethods: []string{"GET", "POST", "DELETE"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("root."))
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	})

	r.Route("/nodes", func(r chi.Router) {
		r.Get("/", node.ListNodes)
		r.Post("/", node.CreateNode)

		r.Route("/{nodeID}", func(r chi.Router) {
			r.Use(NodeCtx)
			r.Get("/", node.GetNode)
			r.Get("/code", node.GetCode)
			r.Put("/", node.UpdateNode)
			r.Delete("/", node.DeleteNode)
		})
	})

	http.ListenAndServe(":3333", r)
}

func NodeCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var node *models.Node
		var err error

		if nodeID := chi.URLParam(r, "nodeID"); nodeID != "" {
			node, err = models.DbGetNode(nodeID)
		} else {
			render.Render(w, r, error_handler.ErrNotFound)
			return
		}
		if err != nil {
			render.Render(w, r, error_handler.ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "node", node)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
