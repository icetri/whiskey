package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/whiskey-back/internal/delivery/api"
)

func NewRouter(h *api.Handlers) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	router.Methods(http.MethodGet).Path("/ping").HandlerFunc(h.Ping)
	router.Methods(http.MethodGet).Path("/gifts").HandlerFunc(h.Gifts)
	router.Methods(http.MethodPost).Path("/support").HandlerFunc(h.Support)

	auth := router.PathPrefix("/auth").Subrouter()
	auth.Methods(http.MethodPost).Path("/pre-sign-up").HandlerFunc(h.PreRegistration)
	auth.Methods(http.MethodPost).Path("/sign-up").HandlerFunc(h.Registration)
	auth.Methods(http.MethodPost).Path("/pre-sign-in").HandlerFunc(h.PreAuth)
	auth.Methods(http.MethodPost).Path("/sign-in").HandlerFunc(h.Auth)

	profile := router.PathPrefix("/profile").Subrouter()
	profile.Use(h.CheckUsers)

	profile.Methods(http.MethodGet).Path("").HandlerFunc(h.Profile)
	profile.Methods(http.MethodPost).Path("/upload").HandlerFunc(h.UploadCheck)
	profile.Methods(http.MethodPost).Path("/write").HandlerFunc(h.HandWriteCheck)
	profile.Methods(http.MethodPost).Path("/prize").HandlerFunc(h.Prize)

	admin := router.PathPrefix("/admin").Subrouter()
	admin.Methods(http.MethodPost).Path("/secret").HandlerFunc(h.UploadCSV)
	admin.Methods(http.MethodGet).Path("/cheques").HandlerFunc(h.GetAllNotVerifiPhoto)
	admin.Methods(http.MethodPost).Path("/verdict").HandlerFunc(h.VerdictPhoto)
	admin.Methods(http.MethodGet).Path("/statistics").HandlerFunc(h.Statistics)
	admin.Methods(http.MethodGet).Path("/recentCheques").HandlerFunc(h.RecentCheques)
	admin.Methods(http.MethodGet).Path("/users").HandlerFunc(h.GetAllUsers)
	admin.Methods(http.MethodGet).Path("/users/{id}").HandlerFunc(h.GetUser)
	admin.Methods(http.MethodGet).Path("/allcheques").HandlerFunc(h.GetAllCheques)
	admin.Methods(http.MethodGet).Path("/statistics/users").HandlerFunc(h.GetCountUsersGift)
	admin.Methods(http.MethodGet).Path("/request-gifts").HandlerFunc(h.GetAllRequestGift)

	return router
}
