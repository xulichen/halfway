package rpc

import (
	"fmt"
	"github.com/xulichen/halfway/pkg/discovery/consul"
	"go.elastic.co/apm/module/apmgrpc"

	"google.golang.org/grpc"
)

// DialWithConsul 客户端通过consul链接服务端
func DialWithConsul(cr consul.Resource, service string) *grpc.ClientConn {
	done := cr.ClaimServices([]string{service})
	if !done {
		panic("claim service failed")
	}
	conn, err := cr.Dial(service)
	if err != nil {
		panic(err)
	}
	return conn
}

// Dial 客户端直连服务端... 需配合linkerd或istio做负载均衡
func Dial(serviceAddr string, options ...grpc.DialOption) (*grpc.ClientConn, error) {
	options = append(options, grpc.WithInsecure(), grpc.WithUnaryInterceptor(apmgrpc.NewUnaryClientInterceptor()))

	conn, err := grpc.Dial(
		fmt.Sprintf("%s://%s", "tcp", serviceAddr),
		options...,
	)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
