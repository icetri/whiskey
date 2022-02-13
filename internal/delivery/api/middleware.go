package api

import (
	"net/http"
)

func (h *Handlers) CheckUsers(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		claims, err := h.tokenManager.GetClaimsByRequest(r, h.JWTkey)
		if err != nil {
			apiErrorEncode(w, err)
			return
		}

		if err = h.services.Common.CheckUsers(claims.UserID); err != nil {
			apiErrorEncode(w, err)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
