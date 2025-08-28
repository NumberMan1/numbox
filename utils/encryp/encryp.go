package encryp

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/NumberMan1/numbox/utils/env"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var (
	DefaultSecret = env.GetEnv("JWT_SECRET_KEY", "MC-ROGUE-SECRET").String()
)

func GeneratorCustomSecret() (string, error) {
	// 生成 32 字节的随机数据
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// 将字节数组转换为Base64编码的字符串
	randomSecret := base64.URLEncoding.EncodeToString(randomBytes)

	// 返回生成的随机字符串
	return randomSecret, nil
}

func GeneratorDefaultAccessKey(claims jwt.Claims) (string, error) {
	return GeneratorAccessKey(DefaultSecret, claims)
}

func GeneratorAccessKey(secret string, claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseDefaultAccessToken[T any](tokenString string, val *T) error {
	return ParseAccessToken(tokenString, DefaultSecret, val)
}

func ParseAccessToken[T any](tokenString string, secret string, val *T) error {
	claims := &valueClaims[T]{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 确保令牌使用了正确的签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return err
	}
	if err = token.Claims.Valid(); err != nil {
		return err
	}
	*val = claims.Value
	return nil
}

func NewValueClaims[T any](val T, expireDuration time.Duration) jwt.Claims {
	claims := &valueClaims[T]{
		Value: val,
	}
	claims.ExpiresAt = time.Now().Add(expireDuration).UnixMilli()
	return claims
}

type valueClaims[T any] struct {
	Value T
	jwt.StandardClaims
}

func (claims *valueClaims[T]) Valid() error {
	if !claims.VerifyExpiresAt(time.Now().UnixMilli(), true) {
		return errors.New("token is expired")
	}
	return nil
}
