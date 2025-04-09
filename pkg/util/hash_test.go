package util

import (
	"net/url"
	"testing"
)

func TestShortUrl(t *testing.T) {
	rawURL := "https://chat.deepseek.com/a/chat/s/0c107bad-d30c7sf9gh-4b4f1934h-87df-1d4ea3fd6a95"
	parsedUrl, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	path := parsedUrl.Path
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}
	hashVal := HashMurmur32(path)
	// shortPath := base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(int(hashVal))))
	shortPath := Base62(int(hashVal))
	t.Log(parsedUrl.Host + "/" + shortPath)
}
