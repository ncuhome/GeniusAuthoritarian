package oss

import "github.com/aliyun/aliyun-oss-go-sdk/oss"

func New(endpoint, accessKey, secretKey, bucket string) (*Client, error) {
	c, e := oss.New(endpoint, accessKey, secretKey)
	if e != nil {
		return nil, e
	}
	b, e := c.Bucket(bucket)
	if e != nil {
		return nil, e
	}
	return &Client{
		c: c,
		b: b,
	}, nil
}

type Client struct {
	path string
	c    *oss.Client
	b    *oss.Bucket
}
