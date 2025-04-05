package util

import (
	"fmt"
	"net/url"
	"testing"
	// "github.com/stretchr/testify/assert"
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
	shortPath := DecToHex(int(hashVal))
	fmt.Println(parsedUrl.Host + "/" + shortPath)
}
