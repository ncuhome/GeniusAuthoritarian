package ldap

import (
	"github.com/go-ldap/ldap/v3"
	"github.com/ncuhome/GeniusAuthoritarian/internal/global"
)

type Client struct {
	Conn *ldap.Conn
}

func New() (Client, error) {
	conn, e := ldap.DialURL(global.Config.Ldap.Addr)
	return Client{
		Conn: conn,
	}, e
}

func (c Client) Close() {
	c.Conn.Close()
}

func (c Client) UserIdentify(username string) string {
	return "cn=" + username
}

// Bind 向该连接绑定账户，相当于登录
func (c Client) Bind(username, password string) error {
	return c.Conn.Bind(c.UserIdentify(username), password)
}

func (c Client) PasswordModify(username, oldPassword, newPassword string) error {
	_, e := c.Conn.PasswordModify(ldap.NewPasswordModifyRequest(c.UserIdentify(username), oldPassword, newPassword))
	return e
}
