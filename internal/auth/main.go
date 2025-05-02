package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

type ChirpyClaims struct {
	jwt.RegisteredClaims
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret []byte, expiresIn time.Duration) (string, error) {
	expire := &jwt.NumericDate{Time: time.Now().Add(expiresIn)}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		Subject:   fmt.Sprintf("%s", userID),
		ExpiresAt: expire,
		IssuedAt:  &jwt.NumericDate{Time: time.Now()},
	})
	signedString, err := token.SignedString(tokenSecret)
	if err != nil {
		return "", err
	}
	return signedString, nil
}

func ValidateJWT(tokenString string, tokenSecret []byte) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ChirpyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return tokenSecret, nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	claims, ok := token.Claims.(*ChirpyClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid token claims")
	}
	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, err
	}
	return userId, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	header := headers.Get("Authorization")
	if len(header) == 0 {
		return "", errors.New("no bearer token")
	}
	token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
	return token, nil
}

func MakeRefreshToken() (string, error) {
	var randomString [32]byte
	_, err := rand.Read(randomString[:])
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(randomString[:]), nil
}

func GetApiKey(headers http.Header) (string, error) {
	header := headers.Get("Authorization")
	if len(header) == 0 {
		return "", errors.New("no api key")
	}
	apiKey := strings.TrimSpace(strings.TrimPrefix(header, "ApiKey "))
	return apiKey, nil
}
