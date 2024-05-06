package api

import (
	"net/http"

	"blackjack/game"
	"blackjack/middleware"
)

func (s Server) DeleteGame(w http.ResponseWriter, r *http.Request) {
	sessionID := middleware.GetSessionID(r.Context())
	game.Delete(s.storage, sessionID)
	respondSuccess(w, nil, http.StatusNoContent)
}
