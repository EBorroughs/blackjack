package routes

import (
	"encoding/json"
	"net/http"

	"blackjack/api"
	"blackjack/game"
	"blackjack/middleware"
)

func GetGame(w http.ResponseWriter, r *http.Request) {
	sessionID := middleware.GetSessionID(r.Context())
	gameState, err := game.Get(sessionID)
	if err != nil {
		api.RespondInternalServerError(w)
		return
	}

	if gameState == nil {
		api.RespondError(w, "no existing game was found for this session", http.StatusNotFound)
		return
	}

	resp := api.GameStateToResponse(gameState)
	respBytes, err := json.Marshal(resp)
	if err != nil {
		api.RespondInternalServerError(w)
		return
	}

	api.RespondSuccess(w, respBytes, http.StatusOK)
}
