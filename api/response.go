package api

import (
	"net/http"

	"blackjack/game"
)

type GameStateResponse struct {
	DealerHand []string `json:"dealerHand"`
	PlayerHand []string `json:"playerHand"`
	Wins       int      `json:"wins"`
	Losses     int      `json:"losses"`
	Ties       int      `json:"ties"`
}

func gameStateToResponse(gameState *game.State) GameStateResponse {
	dealerHand, playerHand := gameState.ToCardFormat()
	resp := GameStateResponse{
		DealerHand: dealerHand,
		PlayerHand: playerHand,
		Wins:       gameState.Wins(),
		Losses:     gameState.Losses(),
		Ties:       gameState.Ties(),
	}

	return resp
}

func respondSuccess(w http.ResponseWriter, body []byte, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if len(body) > 0 {
		w.Write(body)
	}
}

func respondError(w http.ResponseWriter, message string, status int) {
	http.Error(w, message, status)
}

func respondInternalServerError(w http.ResponseWriter) {
	respondError(w, "internal server error", http.StatusInternalServerError)
}
