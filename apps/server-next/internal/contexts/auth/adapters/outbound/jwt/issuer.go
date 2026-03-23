package jwt

import (
	"time"

	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/auth/application/ports"
	"github.com/Hughzu/trackstack/apps/server-next/internal/contexts/auth/domain"
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

type Issuer struct {
	secret string
}

var _ ports.TokenIssuer = (*Issuer)(nil)

func NewIssuer(secret string) *Issuer {
	return &Issuer{secret: secret}
}

func (i *Issuer) IssueToken(userID string) (string, error) {
	now := time.Now().UTC()
	// TTL: 30 days
	expiresAt := now.Add(30 * 24 * time.Hour)

	claims := domain.SessionClaims{
		UserID: userID,
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: jwtv5.NewNumericDate(expiresAt),
			IssuedAt:  jwtv5.NewNumericDate(now),
			NotBefore: jwtv5.NewNumericDate(now),
		},
	}

	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	return token.SignedString([]byte(i.secret))
}
