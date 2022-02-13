package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/whiskey-back/internal/types"
	"github.com/whiskey-back/pkg/infrastruct"
)

func (h *Handlers) UploadCSV(w http.ResponseWriter, r *http.Request) {

	_, header, err := r.FormFile("file")
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	if err = h.services.Verification.UploadCSV(header); err != nil {
		apiErrorEncode(w, err)
		return
	}
}

func (h *Handlers) GetAllNotVerifiPhoto(w http.ResponseWriter, r *http.Request) {

	claims, err := h.tokenManager.GetClaimsByRequest(r, h.JWTkey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	files, err := h.services.Profile.GetAllNotVerifiPhoto(claims.UserID)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, files)
}

type createVerdict struct {
	Date        string       `json:"date"`
	CheckAmount string       `json:"check_amount"`
	FN          string       `json:"fn"`
	FD          string       `json:"fd"`
	FP          string       `json:"fp"`
	Check       types.Status `json:"check"`
	UserID      string       `json:"user_id"`
	PhotoID     string       `json:"photo_id"`
}

func (h *Handlers) VerdictPhoto(w http.ResponseWriter, r *http.Request) {

	claims, err := h.tokenManager.GetClaimsByRequest(r, h.JWTkey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	var req createVerdict
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	err = h.services.Profile.VerifiPhoto(&types.Cheque{
		Date:        req.Date,
		CheckAmount: req.CheckAmount,
		FN:          req.FN,
		FD:          req.FD,
		FP:          req.FP,
		Check:       req.Check,
		UserID:      req.UserID,
	}, req.PhotoID, claims.UserID)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}
}

func (h *Handlers) Statistics(w http.ResponseWriter, r *http.Request) {

	claims, err := h.tokenManager.GetClaimsByRequest(r, h.JWTkey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	statistics, err := h.services.Verification.Statistics(claims.UserID)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, statistics)
}

func (h *Handlers) RecentCheques(w http.ResponseWriter, r *http.Request) {

	claims, err := h.tokenManager.GetClaimsByRequest(r, h.JWTkey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	cheques, err := h.services.Verification.RecentCheques(claims.UserID)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, cheques)
}

func (h *Handlers) GetAllUsers(w http.ResponseWriter, r *http.Request) {

	claims, err := h.tokenManager.GetClaimsByRequest(r, h.JWTkey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	users, err := h.services.Verification.GetAllUsers(claims.UserID)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, users)
}

func (h *Handlers) GetUser(w http.ResponseWriter, r *http.Request) {

	claims, err := h.tokenManager.GetClaimsByRequest(r, h.JWTkey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	id := mux.Vars(r)["id"]

	user, err := h.services.Verification.GetUser(claims.UserID, id)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, user)
}

func (h *Handlers) GetAllCheques(w http.ResponseWriter, r *http.Request) {

	claims, err := h.tokenManager.GetClaimsByRequest(r, h.JWTkey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	us, err := h.services.Verification.GetAllCheques(claims.UserID)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, us)
}

func (h *Handlers) GetCountUsersGift(w http.ResponseWriter, r *http.Request) {

	claims, err := h.tokenManager.GetClaimsByRequest(r, h.JWTkey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	countUsGi, err := h.services.Verification.GetCountUsersGift(claims.UserID)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, countUsGi)
}

func (h *Handlers) GetAllRequestGift(w http.ResponseWriter, r *http.Request) {

	claims, err := h.tokenManager.GetClaimsByRequest(r, h.JWTkey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	allRequestGift, err := h.services.Verification.GetAllRequestGift(claims.UserID)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, allRequestGift)
}
