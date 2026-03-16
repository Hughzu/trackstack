package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"log/slog"
	"strings"
	"time"
)

var ErrInvalidInput = errors.New("invalid input")
var ErrUnauthorized = errors.New("unauthorized")

type Service struct {
	store SessionStore
	cfg   Config
}

func NewService(store SessionStore, cfg Config) *Service {
	return &Service{store: store, cfg: cfg}
}

func (s *Service) CreateSession(ctx context.Context, req CreateSessionRequest) (string, Session, error) {
	start := time.Now()
	defer logAuthTiming(ctx, "create_session.total", start, nil)

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

	insertStart := time.Now()
	if err := s.store.InsertSession(ctx, session); err != nil {
		logAuthTiming(ctx, "create_session.insert", insertStart, err)
		return "", Session{}, err
	}
	logAuthTiming(ctx, "create_session.insert", insertStart, nil)

	return rawToken, session, nil
}

func (s *Service) RevokeSessionByRawToken(ctx context.Context, req RevokeSessionRequest) error {
	rawToken := strings.TrimSpace(req.RawToken)
	if rawToken == "" {
		return nil
	}

	return s.store.RevokeSession(ctx, hashToken(rawToken), time.Now().UTC().Format(time.RFC3339))
}

func (s *Service) Authenticate(ctx context.Context, req AuthenticateRequest) (AuthenticateResponse, error) {
	start := time.Now()
	defer logAuthTiming(ctx, "authenticate.total", start, nil)

	rawToken := strings.TrimSpace(req.RawToken)
	if rawToken == "" {
		return AuthenticateResponse{}, ErrUnauthorized
	}

	findStart := time.Now()
	session, err := s.store.FindSessionByID(ctx, hashToken(rawToken))
	logAuthTiming(ctx, "authenticate.find_session", findStart, err)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return AuthenticateResponse{}, ErrUnauthorized
		}
		return AuthenticateResponse{}, err
	}

	now := time.Now().UTC()
	evaluation := evaluateSession(session, now, s.cfg)
	if !evaluation.Valid {
		return AuthenticateResponse{}, ErrUnauthorized
	}

	response := AuthenticateResponse{
		UserID:    session.UserID,
		SessionID: session.ID,
	}

	if evaluation.NeedsRotation {
		replacementRaw, replacement, err := s.rotateSession(ctx, session, req.Context, now)
		if err != nil {
			return AuthenticateResponse{}, err
		}
		response.SessionID = replacement.ID
		response.ReplacementRaw = &replacementRaw
		expiresAt, err := time.Parse(time.RFC3339, replacement.ExpiresAt)
		if err == nil {
			response.CookieExpires = &expiresAt
		}
		return response, nil
	}

	if evaluation.NeedsTouch {
		absoluteExpiresAt, err := time.Parse(time.RFC3339, session.AbsoluteExpiresAt)
		if err == nil {
			expiresAt := resolveIdleExpiry(now, absoluteExpiresAt, s.cfg.SessionIdleSeconds)
			touchStart := time.Now()
			err = s.store.TouchSession(ctx, session.ID, now.Format(time.RFC3339), expiresAt.Format(time.RFC3339))
			logAuthTiming(ctx, "authenticate.touch_session", touchStart, err)
			if err == nil {
				response.CookieExpires = &expiresAt
			}
		}
	}

	return response, nil
}

type sessionEvaluation struct {
	Valid         bool
	NeedsRotation bool
	NeedsTouch    bool
}

func evaluateSession(session Session, now time.Time, cfg Config) sessionEvaluation {
	expiresAt, err := time.Parse(time.RFC3339, session.ExpiresAt)
	if err != nil || !expiresAt.After(now) {
		return sessionEvaluation{}
	}

	absoluteExpiresAt, err := time.Parse(time.RFC3339, session.AbsoluteExpiresAt)
	if err != nil || !absoluteExpiresAt.After(now) {
		return sessionEvaluation{}
	}

	rotatedAt, err := time.Parse(time.RFC3339, session.RotatedAt)
	if err != nil {
		rotatedAt = now
	}
	lastSeenAt, err := time.Parse(time.RFC3339, session.LastSeenAt)
	if err != nil {
		lastSeenAt = now
	}

	needsRotation := session.RevokedAt != nil || now.Sub(rotatedAt) > time.Duration(cfg.SessionRotateAfterSeconds)*time.Second
	needsTouch := now.Sub(lastSeenAt) > time.Duration(cfg.SessionTouchSeconds)*time.Second

	return sessionEvaluation{
		Valid:         true,
		NeedsRotation: needsRotation,
		NeedsTouch:    needsTouch,
	}
}

func (s *Service) rotateSession(ctx context.Context, session Session, context ClientContext, now time.Time) (string, Session, error) {
	replacementRaw, err := newToken()
	if err != nil {
		return "", Session{}, err
	}

	absoluteExpiresAt, err := time.Parse(time.RFC3339, session.AbsoluteExpiresAt)
	if err != nil {
		return "", Session{}, err
	}
	expiresAt := resolveIdleExpiry(now, absoluteExpiresAt, s.cfg.SessionIdleSeconds)

	var userAgentHash *string
	if context.UserAgent != nil {
		hashed := hashToken(*context.UserAgent)
		userAgentHash = &hashed
	} else {
		userAgentHash = session.UserAgentHash
	}

	replacement := Session{
		ID:                hashToken(replacementRaw),
		UserID:            session.UserID,
		CreatedAt:         now.Format(time.RFC3339),
		ExpiresAt:         expiresAt.Format(time.RFC3339),
		RotatedAt:         now.Format(time.RFC3339),
		LastSeenAt:        now.Format(time.RFC3339),
		AbsoluteExpiresAt: session.AbsoluteExpiresAt,
		ParentID:          &session.ID,
		RevokedAt:         nil,
		UserAgentHash:     userAgentHash,
		IPPrefix:          context.IPPrefix,
	}

	if replacement.IPPrefix == nil {
		replacement.IPPrefix = session.IPPrefix
	}

	insertStart := time.Now()
	if err := s.store.InsertSession(ctx, replacement); err != nil {
		logAuthTiming(ctx, "authenticate.rotate.insert_replacement", insertStart, err)
		return "", Session{}, err
	}
	logAuthTiming(ctx, "authenticate.rotate.insert_replacement", insertStart, nil)

	graceExpiry := now.Add(time.Duration(s.cfg.SessionRotationGraceSeconds) * time.Second)
	if graceExpiry.After(absoluteExpiresAt) {
		graceExpiry = absoluteExpiresAt
	}

	rotateStart := time.Now()
	if err := s.store.RotateOutSession(
		ctx,
		session.ID,
		now.Format(time.RFC3339),
		graceExpiry.Format(time.RFC3339),
		now.Format(time.RFC3339),
	); err != nil {
		logAuthTiming(ctx, "authenticate.rotate.rotate_out", rotateStart, err)
		return "", Session{}, err
	}
	logAuthTiming(ctx, "authenticate.rotate.rotate_out", rotateStart, nil)

	return replacementRaw, replacement, nil
}

func logAuthTiming(ctx context.Context, step string, start time.Time, err error) {
	attrs := []any{"step", step, "duration", time.Since(start)}
	if err != nil {
		attrs = append(attrs, "error", err)
	}
	slog.DebugContext(ctx, "auth timing", attrs...)
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
