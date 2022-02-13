package api

import (
	"encoding/json"
	"github.com/whiskey-back/internal/types"
	"github.com/whiskey-back/pkg/infrastruct"
	"net/http"
)

type createPreUser struct {
	Phone string `json:"phone"`
}

func (h *Handlers) PreRegistration(w http.ResponseWriter, r *http.Request) {
	var req createPreUser

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	err = h.services.Authorization.PreRegistrationUser(&types.CheckUser{
		Phone: req.Phone,
	})
	if err != nil {
		apiErrorEncode(w, err)
		return
	}
}

type createUser struct {
	FirstName string `json:"first_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Code      string `json:"code"`
}

func (h *Handlers) Registration(w http.ResponseWriter, r *http.Request) {
	var req createUser

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	token, err := h.services.Authorization.RegistrationUser(&types.User{
		FirstName: req.FirstName,
		Phone:     req.Phone,
		Email:     req.Email,
		Code:      req.Code,
	})
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, token)
}

func (h *Handlers) PreAuth(w http.ResponseWriter, r *http.Request) {
	var req createPreUser

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	err = h.services.Authorization.PreAuthorizationUser(&types.User{
		Phone: req.Phone,
	})
	if err != nil {
		apiErrorEncode(w, err)
		return
	}
}

type code struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

func (h *Handlers) Auth(w http.ResponseWriter, r *http.Request) {
	var req code

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	token, err := h.services.Authorization.AuthorizationUser(&types.User{
		Phone: req.Phone,
		Code:  req.Code,
	})
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, token)
}
