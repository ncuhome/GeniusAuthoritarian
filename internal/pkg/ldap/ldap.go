package ldap

import (
	"github.com/go-ldap/ldap/v3"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
	log "github.com/sirupsen/logrus"
)

var Conn *ldap.Conn

func init() {
	var e error
	Conn, e = ldap.DialURL(global.Config.Ldap.Addr)
	if e != nil {
		log.Fatalln("连接 ldap 服务失败:", e)
	}
}
