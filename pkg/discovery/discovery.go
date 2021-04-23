package discovery

// ServiceConfig 需要注册的服务配置
type ServiceConfig struct {
	IP              string
	Port            int
	Tags            []string
	Name            string
	ID              string
	HealthyCheckURL string //HealthyCheckURL   健康检查的 URL
	IsRPC           bool
	Meta            map[string]string
}

// ServerConfig Discovery服务配置
type ServerConfig struct {
	Address string
	Port    int
	Token   string
}

type Discovery interface {
	RegisterService() error
	DeregisterService() error
}
