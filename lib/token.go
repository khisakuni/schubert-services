package token

import (
	"errors"
	"time"

	"github.com/caarlos0/env"
	"github.com/dgrijalva/jwt-go"
)

const (
	Issuer        = "schubert"
	ExpireInHours = 1
)

type config struct {
	Secret string `env:"SECRET" envDefault:""`
}

type claim struct {
	Email string
	jwt.StandardClaims
}

func NewToken(email string) (string, error) {
	c := claim{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration(),
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

func expiration() int64 {
	return time.Now().Add(time.Hour * ExpireInHours).Unix()
}

func secret() (string, error) {
	c := config{}
	err := env.Parse(&c)
	if err != nil {
		return "", err
	}
	return c.Secret, nil
}
