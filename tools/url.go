package tools

import (
	"net/url"
)

func AddTokenToUrlQuery(link, token string) (string, error) {
	callbackUrl, e := url.Parse(link)
	if e != nil {
		return "", e
	}
	callbackQuery := callbackUrl.Query()
	callbackQuery.Set("token", token)
	callbackUrl.RawQuery = callbackQuery.Encode()
	return callbackUrl.String(), nil
}
