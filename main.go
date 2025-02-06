package main

import (
	"blackjack/api"
	"blackjack/config"
	"blackjack/storage"

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

	storage, err := storage.NewGameStateStorage(config.StorageBackend)
	if err != nil {
		panic(err)
	}

	server := api.NewServer(store, storage, config.ServerAddress, config.Port, config.RequestTimeoutSeconds)
	server.Start()
}
