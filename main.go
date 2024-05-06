package main

import (
	"net/http"
	"time"

	"blackjack/api/routes"
	"blackjack/config"
	"blackjack/middleware"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	store := sessions.NewCookieStore([]byte("very-secure-key"))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   config.SessionMaxAgeSeconds,
		HttpOnly: true,
	}

	r := chi.NewRouter()
	r.Use(middleware.Session(store))
	r.Use(chiMiddleware.Timeout(time.Duration(config.RequestTimeoutSeconds) * time.Second))

	r.Get("/game", routes.GetGame)
	r.Delete("/game", routes.DeleteGame)
	r.Post("/game", routes.UpsertGame)

	http.ListenAndServe(":8080", r)
}
