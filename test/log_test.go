package test

import (
	"fmt"
	"github.com/xulichen/halfway/pkg/utils"
	"testing"
)

func TestSetLogConfig(t *testing.T) {
	type args struct {
		logConfig utils.LogConfig
	}
	tests := []struct {
		name string
		args args
	}{
		{args: args{
			logConfig: utils.LogConfig{},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			utils.SetLogConfig(tt.args.logConfig)
			fmt.Println(utils.GetLogger())
		})
	}
}
