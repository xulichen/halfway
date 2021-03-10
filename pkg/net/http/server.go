package http

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/xulichen/halfway/pkg/net/http/middleware"
	"go.elastic.co/apm/module/apmechov4"
)

type Server struct {
	*Config
	*echo.Echo
}

func NewServer(cf *Config) *Server {
	server := echo.New()
	middleware.InitValidate(server)
	server.Use(middleware.Recover(), middleware.Logger(), apmechov4.Middleware())
	return &Server{
		Config: cf,
		Echo:   server,
	}
}

func (s Server) Start() {
	s.Logger.Fatal(s.Echo.Start(s.Address()))
}

// Stop 直接退出
func (s *Server) Stop() error {
	return s.Echo.Close()
}

// GracefulStop 优雅退出
func (s *Server) GracefulStop(ctx context.Context) error {
	return s.Echo.Shutdown(ctx)
}
