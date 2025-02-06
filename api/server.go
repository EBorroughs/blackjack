package api

import (
	"net/http"
	"time"

	"blackjack/middleware"
	"blackjack/storage"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
)

type Server struct {
	address string
	port    string
	router  *chi.Mux
	storage storage.GameStateStorage
}

func NewServer(store *sessions.CookieStore, storage storage.GameStateStorage, address, port string, requestTimeoutSeconds int) *Server {
	s := &Server{address: address, port: port, storage: storage}
	r := chi.NewRouter()
	r.Use(middleware.Session(store))
	r.Use(chiMiddleware.Timeout(time.Duration(requestTimeoutSeconds) * time.Second))

	r.Get("/game", s.GetGame)
	r.Delete("/game", s.DeleteGame)
	r.Post("/game", s.UpsertGame)
	s.router = r
	return s
}

func (s Server) Start() {
	address := s.address
	if s.port != "" {
		address += ":" + s.port
	}
	http.ListenAndServe(address, s.router)
}
