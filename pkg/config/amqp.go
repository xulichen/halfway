package config

type AmqpConfig struct {
	AccessKeyId     string
	AccessKeySecret string
	InstanceId      string
	Address         string
}

func NewAmqpConfig(s string) *AmqpConfig {
	m := ConfigMap[s]
	if m != nil {
		return &AmqpConfig{
			AccessKeyId:     m["access-key-id"].(string),
			AccessKeySecret: m["access-key-secret"].(string),
			InstanceId:      m["instance-id"].(string),
			Address:         m["address-id"].(string),
		}
	}
	return nil
}
