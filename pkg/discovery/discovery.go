package discovery

// ServiceConfig 需要注册的服务配置
type ServiceConfig struct {
	IP   string
	Port int
	Tag  []string
	Name string
	ID   string
	//HealthyCheckURL   健康检查的 URL
	HealthyCheckURL string
	IsRPC           bool
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
