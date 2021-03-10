package config

import "github.com/spf13/viper"

type BaseConfig struct {
	ENV         string
	HttpPort    string
	RpcPort     string
	ServiceName string
}

func NewBaseConfig() *BaseConfig {
	return &BaseConfig{
		ENV:      viper.GetString("golang-base.env"),
		HttpPort: viper.GetString("golang-base.http-port"),
		RpcPort:  viper.GetString("golang-base.rpc-port"),
	}
}
