package config

import (
	"reflect"
	"testing"
)

func TestNewConfigMap(t *testing.T) {
	type args struct {
		baseDir string
	}
	tests := []struct {
		name string
		args args
		want map[string]map[string]interface{}
	}{
		{args: args{baseDir: "/credential"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConfigMap(tt.args.baseDir); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfigMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildConfigFile(t *testing.T) {
	type args struct {
		m    map[string]map[string]interface{}
		file string
	}
	tests := []struct {
		name string
		args args
	}{
		{args: args{
			m:    NewConfigMap("/credential"),
			file: "./config.yaml",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BuildConfigFile(tt.args.m, tt.args.file)
		})
	}
}
