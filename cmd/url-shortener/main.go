package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sillkiw/url-shorten/internal/config"
	"github.com/sillkiw/url-shorten/internal/http/handlers/delete"
	"github.com/sillkiw/url-shorten/internal/http/handlers/redirect"
	"github.com/sillkiw/url-shorten/internal/http/handlers/url/save"
	mvLogger "github.com/sillkiw/url-shorten/internal/http/middleware/logger"
	"github.com/sillkiw/url-shorten/internal/lib/validation"
	"github.com/sillkiw/url-shorten/internal/storage/postgres"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := postgres.New(cfg.DB.URL)
	if err != nil {
		log.Error("failed to init storage", slog.Any("error", err))
		os.Exit(1)
	}
	defer func() {
		if err := storage.Close(); err != nil {
			log.Error("failed to close storage", slog.Any("error", err))
		}
	}()

	validator := validation.New(
		cfg.Valid.MyHost,
		cfg.Valid.MaxURLLen,
		cfg.Valid.MinAliasLen,
		cfg.Valid.MaxAliasLen,
		cfg.Valid.DefaultAliasLen,
	)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(mvLogger.New(log))
	router.Use(middleware.URLFormat)

	router.Route("/admin", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTP.Admin: cfg.HTTP.Password,
		}))

		r.Delete("/{alias}", delete.New(log, storage, validator))
	})

	router.Post("/url", save.New(log, storage, validator))
	router.Get("/{alias}", redirect.New(log, storage, validator))

	log.Info("starting server", slog.String("adress", cfg.HTTP.Addr))

	srv := &http.Server{
		Addr:         cfg.HTTP.Addr,
		Handler:      router,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}
	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}

// func FileServer(r chi.Router, path string, root http.FileSystem) error {
// 	if path == "" || path[len(path)-1] != '/' {
// 		return fmt.Errorf("FileServer: path must end with /")
// 	}

// 	fs := http.StripPrefix(path, http.FileServer(root))

// 	r.Handle(path+"*", fs)
// 	return nil
// }
