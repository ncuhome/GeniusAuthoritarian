package ldap

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
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

func (c Client) DC() string {
	return "dc=ncuos,dc=com"
}

func (c Client) UserIdentify(cn string) string {
	return "cn=" + cn + "," + c.DC()
}

// Bind 向该连接绑定账户，相当于登录
func (c Client) Bind(cn, password string) error {
	return c.Conn.Bind(c.UserIdentify(cn), password)
}

func (c Client) Sudo() error {
	return c.Bind(global.Config.Ldap.AdminCN, global.Config.Ldap.AdminPWD)
}

func (c Client) PasswordModify(cn, oldPassword, newPassword string) error {
	_, e := c.Conn.PasswordModify(ldap.NewPasswordModifyRequest(c.UserIdentify(cn), oldPassword, newPassword))
	return e
}

func (c Client) Create(cn, sn, password string, ou []string) error {
	req := ldap.NewAddRequest(c.UserIdentify(cn), nil)
	req.Attribute("objectClass", []string{"top", "person", "organizationalPerson", "inetOrgPerson"})
	req.Attribute("cn", []string{cn})
	req.Attribute("sn", []string{sn})
	if ou != nil {
		req.Attribute("ou", ou)
	}
	encryptPwd, e := c.GenerateSSHA(password)
	if e != nil {
		return e
	}
	req.Attribute("userPassword", []string{encryptPwd})
	return c.Conn.Add(req)
}

func (c Client) GenerateSSHA(password string) (string, error) {
	salt := make([]byte, 4)
	_, e := rand.Read(salt)
	if e != nil {
		return "", e
	}

	hash := sha1.New()
	hash.Write([]byte(password))
	hash.Write(salt)

	digest := hash.Sum(nil)
	digest = append(digest, salt...)

	return "{SSHA}" + base64.StdEncoding.EncodeToString(digest), nil
}

func (c Client) Del(cn string) error {
	return c.Conn.Del(ldap.NewDelRequest(c.UserIdentify(cn), nil))
}

func (c Client) FindSingle(cn, ou string) (*ldap.SearchResult, error) {
	return c.Conn.Search(ldap.NewSearchRequest(
		c.UserIdentify(cn),
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		1, 0, false,
		fmt.Sprintf("(&(ou=%s))", ou),
		[]string{"sn", "ou"},
		nil,
	))
}
