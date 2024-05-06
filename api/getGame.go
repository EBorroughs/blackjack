package api

import (
	"encoding/json"
	"net/http"

	"blackjack/game"
	"blackjack/middleware"
)

func (s Server) GetGame(w http.ResponseWriter, r *http.Request) {
	sessionID := middleware.GetSessionID(r.Context())
	gameState, err := game.Get(s.storage, sessionID)
	if err != nil {
		respondInternalServerError(w)
		return
	}

	if gameState == nil {
		respondError(w, "no existing game was found for this session", http.StatusNotFound)
		return
	}

	resp := gameStateToResponse(gameState)
	respBytes, err := json.Marshal(resp)
	if err != nil {
		respondInternalServerError(w)
		return
	}

	respondSuccess(w, respBytes, http.StatusOK)
}
