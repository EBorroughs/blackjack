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

func GameStateToResponse(gameState *game.State) GameStateResponse {
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

func RespondSuccess(w http.ResponseWriter, body []byte, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if len(body) > 0 {
		w.Write(body)
	}
}

func RespondError(w http.ResponseWriter, message string, status int) {
	http.Error(w, message, status)
}

func RespondInternalServerError(w http.ResponseWriter) {
	RespondError(w, "internal server error", http.StatusInternalServerError)
}
