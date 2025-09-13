package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/RomanKovalev007/mai_news/internal/config"
	"github.com/RomanKovalev007/mai_news/internal/handlers"
	"github.com/RomanKovalev007/mai_news/internal/storage/sqlstore"
)
const (
	envLocal = "local"
	envDev = "dev"
	envProd = "prod"
)
func setupLogger(env string) *slog.Logger{
	var log *slog.Logger
	switch env{
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log

}


func main(){
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	storage, err := sqlstore.New(cfg.StoragePath)
	if err != nil{
		log.Error("failed to start storage", err)
		os.Exit(1)
	}


	log.Info("starting mai_news", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	r := http.NewServeMux()

	r.HandleFunc("GET /posts/", handlers.GetAllPostsHandler(storage, log))
	r.HandleFunc("POST /posts/", handlers.CreatePostHandler(storage, log))
	r.HandleFunc("GET /posts/{id}/",handlers.GetPostHandler(storage, log))
	r.HandleFunc("PATCH /posts/{id}/",handlers.PatchPostHandler(storage, log))
	r.HandleFunc("DELETE /posts/{id}/",handlers.DeletePostHandler(storage, log))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      r,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	log.Info("server started: http/localhost:8000/" )
	if err := srv.ListenAndServe(); err != nil{
		log.Error("failed to start server", err)
	}
	
}