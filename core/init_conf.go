package core

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

var ConfPath = "settings.yaml"

type System struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
}
type Config struct {
	System System
}

func ReadConf() {
	byteData, err := os.ReadFile(ConfPath)
	if err != nil {
		panic(err)
	}
	var config Config
	err = yaml.Unmarshal(byteData, &config)
	if err != nil {
		panic(fmt.Sprintf("yaml格式错误 %v", err))
	}
	fmt.Printf("%#v\n", config)
}
