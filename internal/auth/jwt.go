package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

type JwtManager struct {
	signKey        []byte
	accessTokenTTL time.Duration
}

func NewJwtManager(signKey string, accessTokenTTL time.Duration) *JwtManager {
	return &JwtManager{
		signKey:        []byte(signKey),
		accessTokenTTL: accessTokenTTL,
	}
}

func (m *JwtManager) CreateToken(userID, role string) (string, error) {
	claims := TokenClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(m.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.signKey)
}

func (m *JwtManager) ParseToken(raw string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(raw, &TokenClaims{}, func(t *jwt.Token) (any, error) {
		return m.signKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

func (m *JwtManager) AccessTokenTTL() time.Duration {
	return m.accessTokenTTL
}
