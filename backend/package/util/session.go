package util

import (
	"Mou1ght-Server/config"
	"Mou1ght-Server/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var sampleJwtKey = config.Conf.JwtKey

type Claims struct {
	UID uint
	jwt.MapClaims
}

func ReleaseToken(user *model.User) (string, error) {
	expirationTime := time.Now().Add(30 * time.Hour * 24)
	mapClaims := jwt.MapClaims{}
	mapClaims["exp"] = expirationTime.Unix()
	mapClaims["iat"] = time.Now().Unix()
	mapClaims["iss"] = "Mou1ght"
	mapClaims["sub"] = "user token"
	claims := &Claims{
		UID:       user.ID,
		MapClaims: mapClaims,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(sampleJwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	// 复杂的捏
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return sampleJwtKey, err
	})
	return token, claims, err
}

func ClearToken(tokenString string) error {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return sampleJwtKey, err
	})
	token.Valid = false
	return err
}
