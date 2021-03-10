package config

import (
	"gopkg.in/yaml.v3"
	"github.com/xulichen/halfway/pkg/consts"
	"io/ioutil"
	"os"
	"strings"
)

var ConfigMap = make(map[string]map[string]interface{}, 0)

func InitConfig(configDir string) {
	if configDir == "" {
		configDir = consts.DefaultBaseConfigDir
	}
	InitConfigMap(configDir)
}

func InitConfigMap(baseDir string) {
	files, _ := ioutil.ReadDir(baseDir)
	kv := make(map[string]interface{}, 0)
	for _, file := range files {
		if file.IsDir() {
			if strings.HasPrefix(file.Name(), "..") {
				continue
			}
			InitConfigMap(baseDir + "/" + file.Name())
		} else {
			f, _ := ioutil.ReadFile(baseDir + "/" + file.Name())
			value := strings.ReplaceAll(string(f), "\n", "")
			kv[strings.ReplaceAll(strings.ReplaceAll(file.Name(), "\n", ""), "_", "-")] = value
		}
	}
	keys := strings.Split(baseDir, "/")
	ConfigMap[strings.ReplaceAll(keys[len(keys)-1], "_", "-")] = kv
}

func BuildConfigFile(m map[string]map[string]interface{}, file string) {
	bytes, err := yaml.Marshal(m)
	if err != nil {
		panic(err)
	}
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.WriteString(string(bytes) + "\n")
	if err != nil {
		panic(err)
	}
}