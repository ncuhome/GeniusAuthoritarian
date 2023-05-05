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
	Aliyun    Aliyun       `yaml:"aliyun"`
	Jwt       Jwt          `yaml:"jwt"`
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

type Aliyun struct {
	AccessKey string `yaml:"accessKey"`
	SecretKey string `yaml:"secretKey"`
	Oss       `yaml:"oss"`
}

type Oss struct {
	Endpoint string `yaml:"endpoint"`
	Bucket   string `yaml:"bucket"`
}
