package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var sampleJwtKey = "sampleJwtKey"

type Claims struct {
	UID uint
	jwt.MapClaims
}

func ReleaseToken(userID string) (string, error) {
	id := []byte(userID)
	var uid uint = 0
	for _, b := range id {
		uid = uid + uint(b)
	}
	expirationTime := time.Now().Add(30 * time.Hour * 24)
	mapClaims := jwt.MapClaims{}
	mapClaims["exp"] = expirationTime.Unix()
	mapClaims["iat"] = time.Now().Unix()
	mapClaims["iss"] = "Mou1ght"
	mapClaims["sub"] = "user token"
	mapClaims["uid"] = userID
	claims := &Claims{
		UID:       uid,
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
