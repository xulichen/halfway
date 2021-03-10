package test

import (
	"fmt"
	"github.com/xulichen/halfway/pkg/config"
	"testing"
)

func TestInitConfig(t *testing.T) {
	type args struct {
		configDir string
	}
	tests := []struct {
		name string
		args args
	}{
		{args: args{
			configDir: "",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.InitConfig(tt.args.configDir)
			fmt.Println(config.ConfigMap)
		})
	}
}
