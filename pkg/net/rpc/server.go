package rpc

import (
	"github.com/xulichen/halfway/pkg/discovery"
	"net"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/xulichen/halfway/pkg/log"
	"go.elastic.co/apm/module/apmgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	*grpc.Server
	*Config
	d discovery.Discovery
}

// NewServer return *Server
func NewServer(cf Config, opt ...grpc.ServerOption) *Server {
	s := &Server{
		Config: &cf,
	}
	opt = append(opt, grpcmiddleware.WithUnaryServerChain(
		// apm 集成了recover
		apmgrpc.NewUnaryServerInterceptor(),
		grpczap.UnaryServerInterceptor(log.GetLogger().GetZapLog()),
	))
	opt = append(opt, grpcmiddleware.WithStreamServerChain(
		grpczap.StreamServerInterceptor(log.GetLogger().GetZapLog())),
	)
	keepParam := grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     s.Config.MaxConnectionIdle,
		MaxConnectionAgeGrace: s.Config.MaxConnectionAgeGrace,
		Time:                  s.Config.KeepAliveInterval,
		Timeout:               s.Config.KeepAliveTimeout,
		MaxConnectionAge:      s.Config.MaxConnectionAge,
	})
	opt = append(opt, keepParam)
	grpcServer := grpc.NewServer(opt...)
	s.Server = grpcServer

	// 通用健康检查注册
	hsrv := health.NewServer()
	hsrv.SetServingStatus(cf.Name, healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(grpcServer, hsrv)

	return s
}

// 服务发现 注册服务
func (s Server) Discovery(d discovery.Discovery) error {
	if d == nil {
		return nil
	}
	s.d = d
	return s.d.RegisterService()
}

// 服务发现 注销服务
func (s Server) FadeAway() error {
	if s.d == nil {
		return nil
	}
	return s.d.DeregisterService()
}

//
func (s Server) Start() error {
	lis, err := net.Listen(s.Config.NetWork, s.Address())
	if err != nil {
		return err
	}
	if s.Debug {
		reflection.Register(s.Server)
	}
	go func() {
		if err := s.Server.Serve(lis); err != nil {
			panic(err)
		}
	}()
	return nil
}

// Stop 直接退出
func (s Server) Stop() {
	_ = s.FadeAway()
	s.Server.Stop()
}

// GracefulStop 优雅退出
func (s Server) GracefulStop() {
	_ = s.FadeAway()
	s.Server.GracefulStop()
}
