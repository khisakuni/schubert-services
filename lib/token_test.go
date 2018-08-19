package token

import "testing"

func TestToken(t *testing.T) {
	email := "kohei@example.com"
	tok, err := NewToken(email)
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
}
