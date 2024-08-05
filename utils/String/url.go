package String

import (
	URL "net/url"
)

func EncodeURI(url string) string {
	u, _ := URL.Parse(url)
	u.RawQuery = u.Query().Encode()
	return u.String()
}
