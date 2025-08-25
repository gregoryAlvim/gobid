package api

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gregoryAlvim/gobid/internal/utils"
)

func (api *Api) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !api.Sessions.Exists(r.Context(), "AuthenticateUserId") {
			utils.EncodeJson(w, r, http.StatusUnauthorized, map[string]any{"message": "must be logged in"})
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (api *Api) HandleGetCSRFToken(w http.ResponseWriter, r *http.Request) {
	token := csrf.Token(r)
	utils.EncodeJson(w, r, http.StatusOK, map[string]any{"csrf_token": token})
}
