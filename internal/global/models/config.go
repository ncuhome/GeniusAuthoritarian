package models

import (
	"github.com/ncuhome/GeniusAuthoritarian/pkg/drivers/mysql"
	"github.com/ncuhome/GeniusAuthoritarian/pkg/drivers/redis"
)

type Config struct {
	Mysql mysql.Config
	Redis redis.Config
}
