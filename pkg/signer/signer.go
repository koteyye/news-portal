package signer

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/koteyye/news-portal/pkg/models"
)

// ErrTokenExpired возвращается, если токен просрочен.
var ErrTokenExpired = errors.New("token is expired")

const (
	defaultTTL =  time.Hour * 12 // ttl базовое время жизни токена
)

type claims struct {
	jwt.RegisteredClaims
	Profile *models.Profile
}

// Signer интерфейс подписанта полезной нагрузки
type Signer interface {
	// Sign подписывает полезную нагрузку и возвращает токен
	Sign(payload *models.Profile) (token string, err error)

	// Parse парсит токен и возвращает полезную нагрузку
	Parse(token string) (payload *models.Profile, err error)
}

type jwtSigner struct {
	secret []byte
	ttl time.Duration
}

func New(secret []byte) Signer {
	return &jwtSigner{
		secret: secret,
		ttl: defaultTTL,
	}
}

func (s *jwtSigner) Sign(payload *models.Profile) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS384, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(defaultTTL)),
		},
		Profile: payload,
	})

	tokenString, err := jwtToken.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("ошибка при получении токена: %w", err)
	}

	return tokenString, nil
}

func (s *jwtSigner) Parse(token string) ( *models.Profile, error) {
	claims, err := s.parseToken(token)
	if err != nil {
		return nil, fmt.Errorf("token parsing: %w", err)
	}
	return claims.Profile, nil
}

func (s *jwtSigner) parseToken(token string) (claims, error) {
	var c claims

	_, err := jwt.ParseWithClaims(token, &c, s.secretKey)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			err = ErrTokenExpired
		}
		return claims{}, err
	}
	return c, nil
}

func (s *jwtSigner) secretKey(t *jwt.Token) (any, error) {
	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
	}
	return s.secret, nil
}