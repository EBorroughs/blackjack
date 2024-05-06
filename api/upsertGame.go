package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"blackjack/game"
	"blackjack/middleware"
)

func (s Server) UpsertGame(w http.ResponseWriter, r *http.Request) {
	sessionID := middleware.GetSessionID(r.Context())
	var body upsertGameRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		respondError(w, fmt.Sprintf("unable to read request body: %+v", err), http.StatusBadRequest)
		return
	}

	var gameState *game.State
	var respBytes []byte
	switch body.Action {
	case upsertCreate:
		gameState, err = game.New(s.storage, sessionID)
		if err != nil {
			respondInternalServerError(w)
			return
		}

		resp := gameStateToResponse(gameState)
		respBytes, err = json.Marshal(resp)
		if err != nil {
			respondInternalServerError(w)
			return
		}

	case upsertHit:
		gameState, err = game.Hit(s.storage, sessionID)
		if err != nil {
			if gameState != nil {
				respondError(w, err.Error(), http.StatusBadRequest)
				return
			}

			respondInternalServerError(w)
			return
		}

		if gameState == nil {
			respondError(w, "no existing game was found for this session", http.StatusNotFound)
			return
		}

		resp := gameStateToResponse(gameState)
		respBytes, err = json.Marshal(resp)
		if err != nil {
			respondInternalServerError(w)
			return
		}

	case upsertStand:
		gameState, err = game.Stand(s.storage, sessionID)
		if err != nil {
			if gameState != nil {
				respondError(w, err.Error(), http.StatusBadRequest)
				return
			}

			respondInternalServerError(w)
			return
		}

		if gameState == nil {
			respondError(w, "no existing game was found for this session", http.StatusNotFound)
			return
		}

		resp := gameStateToResponse(gameState)
		respBytes, err = json.Marshal(resp)
		if err != nil {
			respondInternalServerError(w)
			return
		}

	default:
		respondError(w, fmt.Sprintf("invalid action value '%s'", body.Action), http.StatusBadRequest)
		return
	}

	respondSuccess(w, respBytes, http.StatusOK)
}
