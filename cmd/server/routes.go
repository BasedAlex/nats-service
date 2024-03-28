package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/basedalex/nats-service/internal/cache"
	"github.com/basedalex/nats-service/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func serve(cache *cache.Cache, cfg *config.Config) {
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.Port),
		Handler: routes(cache),
	}

	fmt.Printf("starting web server on port %d\n", cfg.Port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func routes(cache *cache.Cache) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)

	mux.Get("/", Ping)

	mux.Get("/orders/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		data, ok := cache.Get(id)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		newData, err := json.Marshal(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return 
		}
		w.WriteHeader(http.StatusOK)
		w.Write(newData)
	})

	return mux
}


func Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello")
}