package dingTalk

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/dingtalk/contact_1_0"
	"github.com/alibabacloud-go/dingtalk/oauth2_1_0"
)

func New() (*Client, error) {
	var (
		protocol = "https"
		region   = "central"
	)
	config := &openapi.Config{
		Protocol: &protocol,
		RegionId: &region,
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
	}, nil
}

type Client struct {
	ContactClient *contact_1_0.Client
	OauthClient   *oauth2_1_0.Client
}
