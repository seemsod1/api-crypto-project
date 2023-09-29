package main

import (
	"api-crypto-project/internal/config"
	"api-crypto-project/internal/http-server/handlers/mail/rate"
	"api-crypto-project/internal/http-server/handlers/mail/sendEmails"
	"api-crypto-project/internal/http-server/handlers/mail/subscribe"
	"api-crypto-project/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func LaunchProject() {
	cfg := config.MustLoad()
	storage := InitStorage(cfg)
	router := InitRouter(storage)
	InitServer(cfg, router)

}

func InitStorage(cfg *config.Config) *sqlite.Storage {
	storage, err := sqlite.NewDB(cfg.StoragePath)
	if err != nil {
		log.Fatalf("failed to init storage")
	}

	return storage
}
func InitServer(cfg *config.Config, router *chi.Mux) {
	log.Print("starting server")
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Print("failed to start server")
	}

	log.Print("server stopped")
}
func InitRouter(storage *sqlite.Storage) *chi.Mux {
	router := chi.NewRouter()

	InitMiddleware(router)
	InitPostMethods(router, storage)

	return router
}
func InitMiddleware(router *chi.Mux) {
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
}
func InitPostMethods(router *chi.Mux, storage *sqlite.Storage) {
	router.Post("/subscribe", subscribe.New(storage))
	router.Post("/sendEmails", sendEmails.New(storage))
	router.Post("/rate", rate.New())
}
