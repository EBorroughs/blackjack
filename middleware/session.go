package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

type userID string

var idCtxKey userID = "userID"

func Session(store *sessions.CookieStore) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			session, err := store.Get(r, "session-key")
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to get session: %+v", err), http.StatusInternalServerError)
			}

			if session.IsNew {
				session.Values["uuid"] = uuid.New().String()
			}

			session.Save(r, w)
			ctx := context.WithValue(r.Context(), idCtxKey, session.Values["uuid"].(string))
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func GetSessionID(ctx context.Context) string {
	return ctx.Value(idCtxKey).(string)
}
