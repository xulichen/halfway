package http

import "fmt"

type Config struct {
	// baseConfig ServiceName
	Name string
	Addr string `yaml:"http.address"`
	// Host Port build Addr
	Host string `yaml:"http.host"`
	Port string `yaml:"http.port"`
}

func (c *Config) Address() string {
	if c.Addr != "" {
		return c.Addr
	}
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
