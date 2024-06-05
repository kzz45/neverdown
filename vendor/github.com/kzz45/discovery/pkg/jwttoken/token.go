package jwttoken

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	TokenKey                 = "Token"
	TokenExpiredTimeInSecond = 3600
	TokenExpiration          = "TOKEN_EXPIRATION"
	SecretSalt               = "NEVERDOWN_SECRET"
	Issuer                   = "Neverdown-Issue"
)

var (
	secretKey       interface{}
	tokenExpiration int64
)

type Claims struct {
	UserId   int    `json:"userId"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"isAdmin"`
	jwt.StandardClaims
}

func Generate(username string, isAdmin bool) (string, error) {
	now := time.Now()
	claims := Claims{
		UserId:   0,
		Username: username,
		IsAdmin:  isAdmin,
		StandardClaims: jwt.StandardClaims{
			Audience:  "",
			ExpiresAt: now.Add(time.Second * time.Duration(GetTokenExpirationFromEnv())).Unix(),
			Id:        "",
			IssuedAt:  now.Unix(),
			Issuer:    Issuer,
			NotBefore: now.Unix(),
			Subject:   "",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString(GetSecretKeyFromEnv())
}

func Parse(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return GetSecretKeyFromEnv(), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("failed vaild tokenClaims:%#v", tokenClaims)
}

func GetSecretKeyFromEnv() interface{} {
	if secretKey != nil {
		return secretKey
	}
	secretKey = []byte(SecretSalt)
	key := os.Getenv(SecretSalt)
	if key == "" {
		return secretKey
	}
	secretKey = []byte(key)
	return secretKey
}

func GetTokenExpirationFromEnv() int64 {
	if tokenExpiration != 0 {
		return tokenExpiration
	}
	tokenExpiration = TokenExpiredTimeInSecond
	key := os.Getenv(TokenExpiration)
	if key == "" {
		return tokenExpiration
	}
	t, err := strconv.Atoi(key)
	if err != nil {
		return tokenExpiration
	}
	tokenExpiration = int64(t)
	return tokenExpiration
}
