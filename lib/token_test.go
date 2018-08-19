package token

import (
	"testing"
	"time"
)

func TestToken(t *testing.T) {
	// Check for email
	email := "kohei@example.com"
	tok, err := NewToken(email, time.Hour*1)
	if err != nil {
		t.Error(err)
	}
	result, err := ParseToken(tok)
	if err != nil {
		t.Error(err)
	}
	if result != email {
		t.Errorf("Expected %s, got %s", email, result)
	}

	// Check for expiration
	expiresIn := time.Second * 1
	tok, err = NewToken(email, expiresIn)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Second * 2)
	_, err = ParseToken(tok)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
