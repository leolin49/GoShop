package util

import (
	"errors"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	issuer        = "go-shop"
	accessJwtKey  = []byte("Iv80ZxU7I8YhfJIIFpol")
	refreshJwtKey = []byte("FVusF9efEyEEjjBH3Yk")
)

func jwtGenerateToken(user_id uint32, sec int64, jwtKey []byte) (token string, err error) {
	expirationTime := time.Now().Add(time.Duration(sec) * time.Second).Unix()
	claims := jwt.MapClaims{
		"sub": strconv.Itoa(int(user_id)),
		"iss": issuer,            // 签发者
		"exp": expirationTime,    // 过期时间
		"nbf": time.Now().Unix(), // 生效时间
		"iat": time.Now().Unix(), // 签发时间
	}
	then := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = then.SignedString(jwtKey)
	return
}

func JwtDoubleToken(user_id uint32, accessSec, refreshSec int64) (accessToken, refreshToken string, err error) {
	accessToken, err = jwtGenerateToken(user_id, accessSec, accessJwtKey)
	if err != nil {
		return "", "", err
	}
	refreshToken, err = jwtGenerateToken(user_id, refreshSec, refreshJwtKey)
	if err != nil {
		return "", "", err
	}
	return
}

func JwtExtractAccessTokenUserId(tokenString string) (uint32, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return accessJwtKey, nil
	})
	if err != nil {
		return 0, err
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("jwt access token parse error.")
		return 0, err
	}
	if !token.Valid {
		err = errors.New("jwt access token is invalid.")
		return 0, err
	}
	id, err := strconv.Atoi(claim["sub"].(string))
	return uint32(id), err
}

func JwtExtractRefreshTokenUserId(tokenString string) (uint32, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return refreshJwtKey, nil
	})
	if err != nil {
		return 0, err
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("jwt refresh token parse error.")
		return 0, err
	}
	if !token.Valid {
		err = errors.New("jwt refresh token is invalid.")
		return 0, err
	}
	id, err := strconv.Atoi(claim["sub"].(string))
	return uint32(id), err
}
