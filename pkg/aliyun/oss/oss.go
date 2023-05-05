package oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"path"
	"time"
)

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

func (a *Storage) NewDir(childDir string) *Storage {
	var dir = *a
	dir.path = path.Join(dir.path, childDir)
	return &dir
}

func (a *Storage) Put(name string, reader io.Reader, options ...oss.Option) error {
	return a.b.PutObject(path.Join(a.path, name), reader, options...)
}

func (a *Storage) Del(name string, options ...oss.Option) error {
	return a.b.DeleteObject(path.Join(a.path, name), options...)
}

func (a *Storage) SignGetUrl(name string, exprIn time.Duration, options ...oss.Option) (string, error) {
	return a.b.SignURL(path.Join(a.path, name), oss.HTTPGet, int64(exprIn.Seconds()), options...)
}
