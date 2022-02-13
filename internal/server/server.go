package server

import (
	"context"
	"github.com/rs/cors"
	"github.com/whiskey-back/internal/config"
	"github.com/whiskey-back/internal/delivery/api"
	"github.com/whiskey-back/pkg/logger"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, handlers *api.Handlers) *Server {

	router := NewRouter(handlers)

	corsOpts := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
		AllowedHeaders: []string{
			"*",
		},
	})

	return &Server{
		httpServer: &http.Server{
			Addr:    cfg.HTTP.Port,
			Handler: corsOpts.Handler(router),
		},
	}
}

func (s *Server) Run() error {
	logger.LogInfo("Restart server")
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	logger.LogInfo("Server shutdown")
	return s.httpServer.Shutdown(ctx)
}
