package util

import (
	"Mou1ght/internal/config"
	"crypto/rand"
	"encoding/hex"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	sampleJwtKey  []byte
	visitorJwtKey []byte
	keyOnce       sync.Once
)

const (
	defaultJwtKey        = "sampleJwtKey"
	defaultVisitorJwtKey = "sampleVisitorJwtKey"
	envJwtKey            = "JWT_KEY"
	envVisitorJwtKey     = "VISITOR_JWT_KEY"
)

func initKeys() {
	keyOnce.Do(func() {
		sampleJwtKey = resolveKey(envJwtKey, defaultJwtKey, true)
		visitorJwtKey = resolveKey(envVisitorJwtKey, defaultVisitorJwtKey, false)
	})
}

func resolveKey(envKey string, fallback string, isPrimary bool) []byte {
	cfg := config.GetConfig().GetSecurity()
	if isPrimary && cfg.JwtKey != "" {
		return []byte(cfg.JwtKey)
	}
	if !isPrimary && cfg.VisitorJwtKey != "" {
		return []byte(cfg.VisitorJwtKey)
	}
	if v := os.Getenv(envKey); v != "" {
		return []byte(v)
	}
	return []byte(fallback)
}

func GetSampleJwtKey() []byte {
	initKeys()
	return sampleJwtKey
}

func GetVisitorJwtKey() []byte {
	initKeys()
	return visitorJwtKey
}

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
	tokenString, err := token.SignedString(GetSampleJwtKey())
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return GetSampleJwtKey(), err
	})
	return token, claims, err
}

func ClearToken(tokenString string) error {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return GetSampleJwtKey(), err
	})
	token.Valid = false
	return err
}

type VisitorClaims struct {
	IP string `json:"ip"`
	UA string `json:"ua"`
	jwt.RegisteredClaims
}

func ReleaseVisitorToken(ip string, ua string) (string, error) {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	now := time.Now()
	claims := &VisitorClaims{
		IP: ip,
		UA: ua,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Mou1ght",
			Subject:   "visitor token",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(365 * 24 * time.Hour)),
			ID:        hex.EncodeToString(b),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(GetVisitorJwtKey())
}

func ParseVisitorToken(tokenString string) (*jwt.Token, *VisitorClaims, error) {
	claims := &VisitorClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return GetVisitorJwtKey(), err
	})
	return token, claims, err
}
