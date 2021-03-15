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

func NewServer(cf *Config, middlewareFunc ...echo.MiddlewareFunc) *Server {
	defaultMiddlewareSlice := []echo.MiddlewareFunc{middleware.Recover(), middleware.Logger(), apmechov4.Middleware()}
	newMiddleware := append(defaultMiddlewareSlice, middlewareFunc...)
	server := echo.New()
	middleware.InitValidate(server)
	server.Use(newMiddleware...)
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
func (s *Server) GracefulStop() error {
	return s.Echo.Shutdown(context.Background())
}
