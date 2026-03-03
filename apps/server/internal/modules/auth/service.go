package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"
)

var ErrInvalidInput = errors.New("invalid input")

type Service struct {
	store SessionStore
	cfg   Config
}

func NewService(store SessionStore, cfg Config) *Service {
	return &Service{store: store, cfg: cfg}
}

func (s *Service) CreateSession(ctx context.Context, req CreateSessionRequest) (string, Session, error) {
	if strings.TrimSpace(req.UserID) == "" {
		return "", Session{}, ErrInvalidInput
	}

	rawToken, err := newToken()
	if err != nil {
		return "", Session{}, err
	}

	now := time.Now().UTC()
	absoluteExpiresAt := now.Add(time.Duration(s.cfg.SessionAbsoluteSeconds) * time.Second)
	expiresAt := resolveIdleExpiry(now, absoluteExpiresAt, s.cfg.SessionIdleSeconds)

	var userAgentHash *string
	if req.Context.UserAgent != nil {
		hashed := hashToken(*req.Context.UserAgent)
		userAgentHash = &hashed
	}

	session := Session{
		ID:                hashToken(rawToken),
		UserID:            req.UserID,
		CreatedAt:         now.Format(time.RFC3339),
		ExpiresAt:         expiresAt.Format(time.RFC3339),
		RotatedAt:         now.Format(time.RFC3339),
		LastSeenAt:        now.Format(time.RFC3339),
		AbsoluteExpiresAt: absoluteExpiresAt.Format(time.RFC3339),
		ParentID:          nil,
		RevokedAt:         nil,
		UserAgentHash:     userAgentHash,
		IPPrefix:          req.Context.IPPrefix,
	}

	if err := s.store.InsertSession(ctx, session); err != nil {
		return "", Session{}, err
	}

	return rawToken, session, nil
}

func (s *Service) RevokeSessionByRawToken(ctx context.Context, req RevokeSessionRequest) error {
	rawToken := strings.TrimSpace(req.RawToken)
	if rawToken == "" {
		return nil
	}

	return s.store.RevokeSession(ctx, hashToken(rawToken), time.Now().UTC().Format(time.RFC3339))
}

func hashToken(value string) string {
	hash := sha256.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}

func newToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func resolveIdleExpiry(now time.Time, absoluteExpiresAt time.Time, idleSeconds int) time.Time {
	idleExpiry := now.Add(time.Duration(idleSeconds) * time.Second)
	if idleExpiry.Before(absoluteExpiresAt) {
		return idleExpiry
	}
	return absoluteExpiresAt
}
