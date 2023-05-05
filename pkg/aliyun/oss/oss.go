package oss

import "github.com/aliyun/aliyun-oss-go-sdk/oss"

func NewRoot(endpoint, accessKey, secretKey, bucket string) (*Storage, error) {
	c, e := oss.New(endpoint, accessKey, secretKey)
	if e != nil {
		return nil, e
	}
	b, e := c.Bucket(bucket)
	if e != nil {
		return nil, e
	}
	return &Storage{
		c: c,
		b: b,
	}, nil
}

type Storage struct {
	path string
	c    *oss.Client
	b    *oss.Bucket
}
