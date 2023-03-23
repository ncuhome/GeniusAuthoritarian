package util

import (
	"github.com/Mmx233/tool"
	"time"
)

var Http *tool.Http

func init() {
	defaultTimeout := time.Second * 30

	Http = tool.NewHttpTool(tool.GenHttpClient(&tool.HttpClientOptions{
		Transport: tool.GenHttpTransport(&tool.HttpTransportOptions{
			Timeout: defaultTimeout,
		}),
		Timeout: defaultTimeout,
	}))
}
