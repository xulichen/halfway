package config

import (
	"flag"
	"fmt"
	"github.com/xulichen/halfway/pkg/consts"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// InitViper 初始化viper
func InitViper(configFilePath string) {
	if _, err := os.Stat(configFilePath); !os.IsNotExist(err) {
		configMap := NewConfigMap(consts.DefaultBaseConfigDir)
		BuildConfigFile(configMap, configFilePath)
	}
	configFile := flag.String("conf", configFilePath, "path of config file")
	viper.SetConfigFile(*configFile)
	err := viper.ReadInConfig()
	if err != nil {
		errStr := fmt.Sprintf("viper read config is failed, err is %v configFile is %v ", err, configFile)
		panic(errStr)
	}
}

// NewConfigMap 基于K8S文件加载解析子文件, 构造map结构体
func NewConfigMap(baseDir string) map[string]map[string]interface{} {
	configMap := make(map[string]map[string]interface{}, 0)
	var r func(string2 string)
	r = func(baseDir string) {
		files, _ := ioutil.ReadDir(baseDir)
		kv := make(map[string]interface{}, 0)
		for _, file := range files {
			if file.IsDir() {
				if strings.HasPrefix(file.Name(), "..") {
					continue
				}
				r(baseDir + "/" + file.Name())
			} else {
				f, _ := ioutil.ReadFile(baseDir + "/" + file.Name())
				value := strings.ReplaceAll(string(f), "\n", "")
				kv[strings.ReplaceAll(strings.ReplaceAll(file.Name(), "\n", ""), "_", "-")] = value
			}
		}
		keys := strings.Split(baseDir, "/")
		configMap[strings.ReplaceAll(keys[len(keys)-1], "_", "-")] = kv
	}
	r(baseDir)
	return configMap
}

// BuildConfigFile 生成yaml文件
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
