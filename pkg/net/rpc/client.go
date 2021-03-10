package rpc

import (
	"github.com/xulichen/halfway/pkg/discovery/consul"

	"google.golang.org/grpc"
)

func Dial(cr consul.Resource, service string) *grpc.ClientConn {
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
