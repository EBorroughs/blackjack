package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"blackjack/api"
	"blackjack/game"
	"blackjack/middleware"
)

func UpsertGame(w http.ResponseWriter, r *http.Request) {
	sessionID := middleware.GetSessionID(r.Context())
	var body api.UpsertGameRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		api.RespondError(w, fmt.Sprintf("unable to read request body: %+v", err), http.StatusBadRequest)
		return
	}

	var gameState *game.State
	var respBytes []byte
	switch body.Action {
	case api.UpsertCreate:
		gameState, err = game.New(sessionID)
		if err != nil {
			api.RespondInternalServerError(w)
			return
		}

		resp := api.GameStateToResponse(gameState)
		respBytes, err = json.Marshal(resp)
		if err != nil {
			api.RespondInternalServerError(w)
			return
		}

	case api.UpsertHit:
		gameState, err = game.Hit(sessionID)
		if err != nil {
			if gameState != nil {
				api.RespondError(w, err.Error(), http.StatusBadRequest)
				return
			}

			api.RespondInternalServerError(w)
			return
		}

		if gameState == nil {
			api.RespondError(w, "no existing game was found for this session", http.StatusNotFound)
			return
		}

		resp := api.GameStateToResponse(gameState)
		respBytes, err = json.Marshal(resp)
		if err != nil {
			api.RespondInternalServerError(w)
			return
		}

	case api.UpsertStand:
		gameState, err = game.Stand(sessionID)
		if err != nil {
			if gameState != nil {
				api.RespondError(w, err.Error(), http.StatusBadRequest)
				return
			}

			api.RespondInternalServerError(w)
			return
		}

		if gameState == nil {
			api.RespondError(w, "no existing game was found for this session", http.StatusNotFound)
			return
		}

		resp := api.GameStateToResponse(gameState)
		respBytes, err = json.Marshal(resp)
		if err != nil {
			api.RespondInternalServerError(w)
			return
		}

	default:
		api.RespondError(w, fmt.Sprintf("invalid action value '%s'", body.Action), http.StatusBadRequest)
		return
	}

	api.RespondSuccess(w, respBytes, http.StatusOK)
}
