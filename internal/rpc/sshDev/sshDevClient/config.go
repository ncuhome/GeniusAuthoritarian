package sshDevClient

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Token string `yaml:"token"`
	Addr  string `yaml:"addr"`
}

func ReadConfig() Config {
	f, err := os.Open("config.yaml")
	if err != nil {
		log.Fatalln("读取配置文件失败:", err)
	}
	defer f.Close()

	var conf Config
	err = yaml.NewDecoder(f).Decode(&conf)
	if err != nil {
		log.Fatalln("解析配置文件失败:", err)
	}
	return conf
}
