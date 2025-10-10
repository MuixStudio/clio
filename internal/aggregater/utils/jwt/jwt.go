package jwt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	accessTokenSecret  = []byte("access_token")
	refreshTokenSecret = []byte("refresh_token")
	accessTokenExp     = time.Hour * 24     // 1 day
	refreshTokenExp    = time.Hour * 24 * 7 // 7 days

	signingMethod = jwt.SigningMethodHS256
)

func ParseAccessToken(tokenStr string) (jwt.MapClaims, error) {
	claims := &jwt.MapClaims{}

	// 解析 & 验证
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		// 校验签名用的密钥
		return accessTokenSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// 检查 token 是否有效
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return *claims, nil
}

func ParseRefreshToken(tokenStr string) (jwt.MapClaims, error) {
	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return refreshTokenSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return *claims, nil
}

func GenerateAccessToken(userID uint64) (string, error) {
	accessTokenClaims := jwt.MapClaims{
		"sub":     strconv.Itoa(int(userID)),
		"exp":     time.Now().Add(accessTokenExp).Unix(),
		"iat":     time.Now().Unix(),
		"nbf":     time.Now().Unix(),
		"iss":     "clio",
		"aud":     "webcall",
		"user_id": userID,
	}
	accessToken := jwt.NewWithClaims(signingMethod, accessTokenClaims)

	accessTokenStr, err := accessToken.SignedString(accessTokenSecret)
	if err != nil {
		return "", err
	}
	return accessTokenStr, nil
}

func GenerateRefreshToken(userID uint64) (string, error) {
	refreshTokenClaims := jwt.MapClaims{
		"sub":     strconv.Itoa(int(userID)),
		"exp":     time.Now().Add(refreshTokenExp).Unix(),
		"iat":     time.Now().Unix(),
		"nbf":     time.Now().Unix(),
		"iss":     "clio",
		"aud":     "webcall",
		"user_id": userID,
	}
	refreshToken := jwt.NewWithClaims(signingMethod, refreshTokenClaims)
	refreshTokenStr, err := refreshToken.SignedString(refreshTokenSecret)
	if err != nil {
		return "", err
	}
	return refreshTokenStr, nil
}
