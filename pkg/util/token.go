package util

import (
	"errors"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang/glog"
)

var jwtKey = []byte("#goshop2023023290^")

func JwtGenerateToken(user_id uint32) (token string, err error) {
	expirationTime := time.Now().Add(30 * time.Minute).Unix()
	claims := jwt.MapClaims{
		"sub": strconv.Itoa(int(user_id)),
		"iss": "go-shop",	// 签发者
		"exp": expirationTime,		// 过期时间
		"nbf": time.Now().Unix(),		// 生效时间
		"iat": time.Now().Unix(),		// 签发时间
	}
	then := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = then.SignedString(jwtKey)
	return
}

func secret() jwt.Keyfunc {
	return func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	}
}

func JwtExtractTokenUserId(tokenString string) (uint32, error) {
	glog.Errorln(tokenString)
	token, err := jwt.Parse(tokenString, secret())
	if err != nil {
		return 0, err
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("jwt token parse error.")
		return 0, err
	}
	if !token.Valid {
		err = errors.New("jwt token is invalid.")
		return 0, err
	}
	id, err := strconv.Atoi(claim["sub"].(string))
	return uint32(id), err
}

