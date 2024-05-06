package routes

import (
	"net/http"

	"blackjack/api"
	"blackjack/game"
	"blackjack/middleware"
)

func DeleteGame(w http.ResponseWriter, r *http.Request) {
	sessionID := middleware.GetSessionID(r.Context())
	game.Delete(sessionID)
	api.RespondSuccess(w, nil, http.StatusNoContent)
}
