package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/CodebyTecs/wishlist-service/internal/domain"
)

type Service struct {
	secret []byte
	ttl    time.Duration
}

func NewService(secret string, ttl time.Duration) *Service {
	return &Service{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

func (s *Service) GenerateAccessToken(userID string) (string, error) {
	now := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *Service) ParseAccessToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrUnauthorized
		}
		return s.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !token.Valid {
		return "", domain.ErrUnauthorized
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || claims.Subject == "" {
		return "", domain.ErrUnauthorized
	}
	return claims.Subject, nil
}
