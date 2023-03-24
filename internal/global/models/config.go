package models

import (
	"github.com/ncuhome/GeniusAuthoritarian/pkg/drivers/mysql"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/drivers/redis"
)

type Config struct {
	Mysql  mysql.Config
	Redis  redis.Config
	Ldap   Ldap
	Feishu Feishu
	Jwt    Jwt
}

type Ldap struct {
	Addr     string // example: ldap://ldap.example.com:389
	AdminCN  string
	AdminPWD string
}

type Feishu struct {
	ClientID string
	Secret   string
}

type Jwt struct {
	SignKey string
}
