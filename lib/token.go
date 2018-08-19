package token

import (
	"errors"
	"time"

	"github.com/caarlos0/env"
	"github.com/dgrijalva/jwt-go"
)

const (
	Issuer = "schubert"
)

type config struct {
	Secret           string `env:"SECRET" envDefault:""`
	TokenExpireHours int    `env:"TOKEN_EXPIREHOURES" envDefault:1`
}

type claim struct {
	Email string
	jwt.StandardClaims
}

func NewToken(email string, expireIn time.Duration) (string, error) {
	c := claim{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration(expireIn),
			Issuer:    Issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	s, err := secret()
	if err != nil {
		return "", err
	}
	return token.SignedString([]byte(s))
}

func ParseToken(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &claim{}, func(token *jwt.Token) (interface{}, error) {
		if jwt.SigningMethodHS256 != token.Method {
			return nil, errors.New("Invalid signing method")
		}
		s, err := secret()
		if err != nil {
			return nil, err
		}
		return []byte(s), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*claim); ok && token.Valid {
		return claims.Email, nil
	}

	return "", errors.New("Invalid token")
}

func expiration(expiresIn time.Duration) int64 {
	return time.Now().Add(expiresIn).Unix()
}

func secret() (string, error) {
	c := config{}
	err := env.Parse(&c)
	if err != nil {
		return "", err
	}
	return c.Secret, nil
}
