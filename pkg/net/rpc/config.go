package rpc

import (
	"fmt"
	"time"
)

type Config struct {
	// baseConfig ServiceName
	Name string
	// default tcp
	NetWork string `yaml:"network"`
	Addr    string `yaml:"address"`
	// Host Port build Addr
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	// Timeout is context timeout for per rpc call.
	Timeout time.Duration `yaml:"grpc.timeout"`
	// MaxConnectionIdle is a duration for the amount of time after which an
	// idle connection would be closed by sending a GoAway. Idleness duration
	// is defined since the most recent time the number of outstanding RPCs
	// became zero or the connection establishment.
	MaxConnectionIdle time.Duration `yaml:"grpc.maxConnectionIdle"`
	// MaxConnectionAge is a duration for the maximum amount of time a
	// connection may exist before it will be closed by sending a GoAway. A
	// random jitter of +/-10% will be added to MaxConnectionAge to spread out
	// connection storms.
	MaxConnectionAge time.Duration `yaml:"grpc.maxConnectionAge"`
	// MaxConnectionAgeGrace is an additive period after MaxConnectionAge after
	// which the connection will be forcibly closed.
	MaxConnectionAgeGrace time.Duration `yaml:"grpc.maxConnectionAgeGrace"`
	// After a duration of this time if the server doesn't see any activity it
	// pings the client to see if the transport is still alive.
	// If set below 1s, a minimum value of 1s will be used instead.
	KeepAliveInterval time.Duration `yaml:"grpc.keepaliveInterval"`
	// After having pinged for keepalive check, the server waits for a duration
	// of Timeout and if no activity is seen even after that the connection is
	// closed. default 20s
	KeepAliveTimeout time.Duration `yaml:"grpc.keepaliveTimeout"`
	// Debug mode
	Debug bool
	// DisableAPM 关闭apm
	DisableAPM bool
}

func (c *Config) Address() string {
	if c.Addr != "" {
		return c.Addr
	}
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
