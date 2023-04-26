package dingTalkApi

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/dingtalk/contact_1_0"
	"github.com/alibabacloud-go/dingtalk/oauth2_1_0"
	"github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

func New(c Config) (*Client, error) {
	config := &openapi.Config{
		Protocol: tea.String("https"),
		RegionId: tea.String("central"),
	}

	oc, e := oauth2_1_0.NewClient(config)
	if e != nil {
		return nil, e
	}

	cc, e := contact_1_0.NewClient(config)
	if e != nil {
		return nil, e
	}

	return &Client{
		ContactClient: cc,
		OauthClient:   oc,
		Config:        c,
	}, nil
}

type Client struct {
	ContactClient *contact_1_0.Client
	OauthClient   *oauth2_1_0.Client
	Config
}

func (c Client) GetUserInfo(accessToken string) (*contact_1_0.GetUserResponse, error) {
	r, e := c.ContactClient.GetUserWithOptions(tea.String("me"), &contact_1_0.GetUserHeaders{
		XAcsDingtalkAccessToken: &accessToken,
	}, &service.RuntimeOptions{})
	return r, e
}

func (c Client) GetUserToken(authCode string) (*oauth2_1_0.GetUserTokenResponse, error) {
	return c.OauthClient.GetUserToken(&oauth2_1_0.GetUserTokenRequest{
		ClientId:     &c.ClientID,
		ClientSecret: &c.Secret,
		Code:         &authCode,
		GrantType:    tea.String("authorization_code"),
	})
}
