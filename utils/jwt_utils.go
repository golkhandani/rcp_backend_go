package utils

import (
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golkhandani/shopWise/configs"
)

var _SecretKey = []byte(configs.Env.JWTSecret)
var JWTSigningKey = jwtware.SigningKey{Key: _SecretKey}
var JWTContextKey = "token"
var JWTUserIDKey = "uid"
var LocalUserKey = "user"

func CreateAccessToken(uid string) (string, int64, error) {
	day := time.Hour * 24
	exp := time.Now().Add(day).Unix()
	claims := jwt.MapClaims{
		JWTUserIDKey: uid,
		"exp":        exp,
	}
	// Create token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	signed, err := accessToken.SignedString(_SecretKey)
	return signed, exp, err
}

func CreateRefreshToken(uid string) (string, int64, error) {
	days := time.Hour * 24 * 30
	exp := time.Now().Add(days).Unix()
	claims := jwt.MapClaims{
		JWTUserIDKey: uid,
		"exp":        exp,
	}
	// Create token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	signed, err := refreshToken.SignedString(_SecretKey)
	return signed, exp, err
}

func DecodeToken(tokenString string) (jwt.MapClaims, int64, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(_SecretKey), nil
	})

	begin := time.Date(1970, 1, 1, 0, 0, 0, 0, time.Now().Location()).Unix()
	if err != nil {
		return nil, begin, err
	}

	dt, err := claims.GetExpirationTime()
	if err != nil {
		return nil, begin, err
	}
	return claims, dt.Unix(), nil
}
