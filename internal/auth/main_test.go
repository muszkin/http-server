package auth

import (
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestJWT(t *testing.T) {
	userId := uuid.New()
	const tokenSecret = "secret"
	duration, _ := time.ParseDuration("5s")
	token, err := MakeJWT(userId, tokenSecret, duration)
	if err != nil {
		t.Fatal(err)
	}
	userIdFromToken, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatal(err)
	}
	if userIdFromToken != userId {
		t.Fatal("user id from token does not match")
		return
	}
}

func TestJWTExpire(t *testing.T) {
	userId := uuid.New()
	const tokenSecret = "secret"
	duration, _ := time.ParseDuration("1s")
	token, err := MakeJWT(userId, tokenSecret, duration)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Duration(2) * time.Second)
	_, err = ValidateJWT(token, tokenSecret)
	if err == nil {
		t.Fatal("token should expire")
	}
}
