package api

import (
	"encoding/json"
	"github.com/whiskey-back/internal/types"
	"github.com/whiskey-back/pkg/infrastruct"
	"net/http"
	"time"
)

func (h *Handlers) Profile(w http.ResponseWriter, r *http.Request) {

	claims, err := h.tokenManager.GetClaimsByRequest(r, h.JWTkey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	profile, err := h.services.Profile.GetProfile(claims.UserID)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, profile)
}

func (h *Handlers) UploadCheck(w http.ResponseWriter, r *http.Request) {

	dateStart, err := time.Parse("2006-01-02T15:04:05", h.dataStart)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	if time.Now().Before(dateStart) {
		w.WriteHeader(http.StatusLocked)
		return
	}

	claims, err := h.tokenManager.GetClaimsByRequest(r, h.JWTkey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	err = r.ParseMultipartForm(209715200)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	files := r.MultipartForm.File["file"]
	for _, file := range files {
		if err = h.services.Profile.UploadCheck(file, claims.UserID); err != nil {
			apiErrorEncode(w, err)
			return
		}
	}
}

type handWrite struct {
	Date        string `json:"date"`
	CheckAmount string `json:"check_amount"`
	FN          string `json:"fn"`
	FD          string `json:"fd"`
	FP          string `json:"fp"`
}

func (h *Handlers) HandWriteCheck(w http.ResponseWriter, r *http.Request) {

	dateStart, err := time.Parse("2006-01-02T15:04:05", h.dataStart)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	if time.Now().Before(dateStart) {
		w.WriteHeader(http.StatusLocked)
		return
	}

	claims, err := h.tokenManager.GetClaimsByRequest(r, h.JWTkey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	var req handWrite
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	check, err := h.services.Profile.HandWriteCheck(&types.Cheque{
		Date:        req.Date,
		CheckAmount: req.CheckAmount,
		FN:          req.FN,
		FD:          req.FD,
		FP:          req.FP,
		UserID:      claims.UserID,
	})
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, check)
}

func (h *Handlers) Prize(w http.ResponseWriter, r *http.Request) {

	dateStart, err := time.Parse("2006-01-02T15:04:05", h.dataStart)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	if time.Now().Before(dateStart) {
		w.WriteHeader(http.StatusLocked)
		return
	}

	claims, err := h.tokenManager.GetClaimsByRequest(r, h.JWTkey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	req := new(types.Prize)
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	products, err := h.services.Profile.PrizeLogic(req, claims.UserID)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, products)
}
