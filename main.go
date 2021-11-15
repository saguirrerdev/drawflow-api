package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	node "github.com/zebek95/draflow-api/controllers"
	"github.com/zebek95/draflow-api/database"
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
			r.Delete("/", node.DeleteNode)
		})
	})

	workDir, _ := os.Getwd()
	fileServer := http.FileServer(http.Dir(filepath.Join(workDir, "app")))
	r.Handle("/*", http.StripPrefix("/", fileServer))

	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("API_PORT")), r)
}

func NodeCtx(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var node *models.Node

		if nodeID := chi.URLParam(r, "nodeID"); nodeID != "" {
			node = models.DbGetNode(nodeID)
		}

		ctx := context.WithValue(r.Context(), "node", node)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
