package config

type BaseConfig struct {
	ENV      string
	HttpPort string
	RpcPort  string
	Service  string
}

func NewBaseConfig(s string) *BaseConfig {
	m := ConfigMap[s]
	if m != nil {
		return &BaseConfig{
			ENV:      m["env"].(string),
			HttpPort: m["http-port"].(string),
			RpcPort:  m["rpc-port"].(string),
		}
	}
	return nil
}
