package jwt

import (
	"time"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/auth/application/ports"
	"github.com/Hughzu/trackstack/apps/server/internal/contexts/auth/domain"
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

type Issuer struct {
	secret string
	ttl    time.Duration
}

var _ ports.TokenIssuer = (*Issuer)(nil)

func NewIssuer(secret string, ttl time.Duration) *Issuer {
	return &Issuer{secret: secret, ttl: ttl}
}

func (i *Issuer) IssueToken(input ports.IssueTokenInput) (ports.IssuedToken, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(i.ttl)

	claims := domain.SessionClaims{
		UserID:    input.UserID,
		SessionID: input.SessionID,
		TokenUse:  domain.TokenUseAccess,
		RegisteredClaims: jwtv5.RegisteredClaims{
			Subject:   input.UserID,
			ExpiresAt: jwtv5.NewNumericDate(expiresAt),
			IssuedAt:  jwtv5.NewNumericDate(now),
			NotBefore: jwtv5.NewNumericDate(now),
		},
	}

	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(i.secret))
	if err != nil {
		return ports.IssuedToken{}, err
	}

	return ports.IssuedToken{
		Value:     signedToken,
		ExpiresAt: expiresAt,
	}, nil
}
