package api

import (
	"encoding/json"
	"github.com/whiskey-back/internal/config"
	"github.com/whiskey-back/internal/service"
	"github.com/whiskey-back/internal/types"
	"github.com/whiskey-back/pkg/infrastruct"
	"github.com/whiskey-back/pkg/jwt"
	"github.com/whiskey-back/pkg/logger"
	"net/http"
)

type Handlers struct {
	services     *service.Service
	tokenManager *jwt.Manager
	JWTkey       string
	dataStart    string
}

func NewHandlers(cfg *config.Config, services *service.Service, tokenManager *jwt.Manager) *Handlers {
	return &Handlers{
		services:     services,
		tokenManager: tokenManager,
		JWTkey:       cfg.JWTKey,
		dataStart:    cfg.DateStart,
	}
}

func (h *Handlers) Ping(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("pong"))
}

func apiErrorEncode(w http.ResponseWriter, err error) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if customError, ok := err.(*infrastruct.CustomError); ok {
		w.WriteHeader(customError.Code)
	}

	result := struct {
		Err string `json:"error"`
	}{
		Err: err.Error(),
	}

	if err = json.NewEncoder(w).Encode(result); err != nil {
		logger.LogError(err)
	}
}

func apiResponseEncoder(w http.ResponseWriter, res interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.LogError(err)
	}
}

func (h *Handlers) Gifts(w http.ResponseWriter, r *http.Request) {

	gifts, err := h.services.Common.GetGifts()
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, gifts)
}

func (h *Handlers) Support(w http.ResponseWriter, r *http.Request) {

	req := new(types.Support)
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	if err = h.services.Common.SupportSend(req); err != nil {
		apiErrorEncode(w, err)
		return
	}
}
