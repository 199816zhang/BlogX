package core

import (
	"blogx_server/flags"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type System struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
}
type Config struct {
	System System
}

func ReadConf() {
	byteData, err := os.ReadFile(flags.FlagOptions.File)
	if err != nil {
		panic(err)
	}
	var config Config
	err = yaml.Unmarshal(byteData, &config)
	if err != nil {
		panic(fmt.Sprintf("yaml格式错误 %v", err))
	}
	fmt.Printf("读取配置文件%s成功", flags.FlagOptions.File)
}
