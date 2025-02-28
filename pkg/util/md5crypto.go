package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

const DefaultSalt = "$1$WNQC2GY5$Yxw4ghncpwfqEN/zRz3LX0"

func DefaultSalfFunc() string {
	return DefaultSalt
}

func MD5(text string) string {
	data := text
	hash := md5.New()
	hash.Write([]byte(data))
	hashedBytes := hash.Sum(nil)
	return hex.EncodeToString(hashedBytes)
}

func MD5WithSalt(text, salt string) string {
	return MD5(text + salt)
}

func MD5WithSaltFun(text string, saltFunc func() string) string {
	if saltFunc == nil {
		saltFunc = DefaultSalfFunc
	}
	return MD5WithSalt(text, fmt.Sprintf("%s|%d", saltFunc(), time.Now().UnixNano()))
}
