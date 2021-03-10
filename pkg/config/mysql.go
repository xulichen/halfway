package config

type MySqlConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       string
}

func NewMySqlConfig(s string) *MySqlConfig {
	m := ConfigMap[s]
	if m != nil {
		return &MySqlConfig{
			Host:     m["host"].(string),
			Port:     m["port"].(string),
			User:     m["user"].(string),
			Password: m["password"].(string),
		}
	}
	return nil
}
