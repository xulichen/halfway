package config

type RedisConfig struct {
	Host     string
	Password string
	Port     string
}

//func NewRedisConfig(s string) *RedisConfig {
//	m := ConfigMap[s]
//	if m != nil {
//		return &RedisConfig{
//			Host:     m["host"].(string),
//			Port:     m["port"].(string),
//			Password: m["password"].(string),
//		}
//	}
//	return nil
//}
