package models

import (
	"github.com/ncuhome/GeniusAuthoritarian/pkg/drivers/mysql"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/drivers/redis"
)

type Config struct {
	Mysql     mysql.Config `yaml:"mysql"`
	Redis     redis.Config `yaml:"redis"`
	Feishu    Feishu       `yaml:"feishu"`
	DingTalk  DingTalk     `yaml:"dingTalk"`
	Jwt       Jwt          `yaml:"jwt"`
	SshDev    SshDev       `yaml:"sshDev"`
	Sms       Sms          `yaml:"sms"`
	TraceMode bool         `yaml:"traceMode" config:"omitempty"`
}

type Feishu struct {
	ClientID string `yaml:"clientID"`
	Secret   string `yaml:"secret"`
}

type DingTalk struct {
	ClientID string `yaml:"clientID"`
	Secret   string `yaml:"secret"`
}

type Jwt struct {
	SignKey string `yaml:"signKey"`
}

type SshDev struct {
	Token string `yaml:"token" config:"omitempty"`
}

type Sms struct {
	SpCode    string `yaml:"spCode"`
	LoginName string `yaml:"loginName"`
	Password  string `yaml:"password"`
}
