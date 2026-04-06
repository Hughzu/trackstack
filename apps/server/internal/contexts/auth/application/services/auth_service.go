package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/auth/application/ports"
	"github.com/Hughzu/trackstack/apps/server/internal/contexts/auth/domain"
	"github.com/google/uuid"
)

type AuthService struct {
	userProvider       ports.UserProvider
	tokenIssuer        ports.TokenIssuer
	sessionRepo        ports.SessionRepository
	refreshTokens      ports.RefreshTokenManager
	refreshTTL         time.Duration
	refreshAbsoluteTTL time.Duration
	now                func() time.Time
}

var _ ports.AuthUseCase = (*AuthService)(nil)

func NewAuthService(
	userProvider ports.UserProvider,
	tokenIssuer ports.TokenIssuer,
	sessionRepo ports.SessionRepository,
	refreshTokens ports.RefreshTokenManager,
	refreshTTL time.Duration,
	refreshAbsoluteTTL time.Duration,
) *AuthService {
	return &AuthService{
		userProvider:       userProvider,
		tokenIssuer:        tokenIssuer,
		sessionRepo:        sessionRepo,
		refreshTokens:      refreshTokens,
		refreshTTL:         refreshTTL,
		refreshAbsoluteTTL: refreshAbsoluteTTL,
		now:                func() time.Time { return time.Now().UTC() },
	}
}

func (s *AuthService) Login(ctx context.Context, command ports.LoginCommand) (ports.AuthTokens, error) {
	email := strings.ToLower(strings.TrimSpace(command.Email))
	if email == "" || command.Password == "" {
		return ports.AuthTokens{}, domain.ErrInvalidInput
	}

	userID, err := s.userProvider.VerifyCredentials(ctx, email, command.Password)
	if err != nil {
		return ports.AuthTokens{}, domain.ErrUnauthorized
	}

	tokens, err := s.createSessionTokens(ctx, userID, metadataFromRequest(command.UserAgent, command.IP))
	if err != nil {
		return ports.AuthTokens{}, err
	}

	asyncCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second)
	go func() {
		defer cancel()
		_ = s.userProvider.UpdateLastLogin(asyncCtx, userID, s.now().Format(time.RFC3339))
	}()

	return tokens, nil
}

func (s *AuthService) Refresh(ctx context.Context, command ports.RefreshCommand) (ports.AuthTokens, error) {
	refreshToken := strings.TrimSpace(command.RefreshToken)
	if refreshToken == "" {
		return ports.AuthTokens{}, domain.ErrUnauthorized
	}

	metadata := metadataFromRequest(command.UserAgent, command.IP)
	tokenHash := s.refreshTokens.HashToken(refreshToken)
	session, err := s.sessionRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, domain.ErrSessionNotFound) {
			return ports.AuthTokens{}, domain.ErrUnauthorized
		}
		return ports.AuthTokens{}, err
	}

	now := s.now()
	if session.IsRotated() {
		_ = s.sessionRepo.RevokeFamily(ctx, sessionRootID(session), now)
		return ports.AuthTokens{}, domain.ErrUnauthorized
	}
	if session.IsRevoked() || session.IsExpired(now) {
		_ = s.sessionRepo.RevokeSession(ctx, session.ID, now)
		return ports.AuthTokens{}, domain.ErrUnauthorized
	}

	newRefreshToken, err := s.refreshTokens.GenerateToken()
	if err != nil {
		return ports.AuthTokens{}, fmt.Errorf("generate refresh token: %w", err)
	}

	newSessionID := uuid.NewString()
	rotatedSession, err := s.sessionRepo.RotateSession(ctx, ports.RotateSessionInput{
		SessionID:         session.ID,
		TokenHash:         tokenHash,
		RotatedAt:         now,
		NewSessionID:      newSessionID,
		NewTokenHash:      s.refreshTokens.HashToken(newRefreshToken),
		NewExpiresAt:      now.Add(s.refreshTTL),
		NewLastSeenAt:     now,
		AbsoluteExpiresAt: session.AbsoluteExpiresAt,
		UserAgentHash:     metadata.userAgentHash,
		IPPrefix:          metadata.ipPrefix,
	})
	if err != nil {
		if errors.Is(err, domain.ErrSessionNotFound) {
			return ports.AuthTokens{}, domain.ErrUnauthorized
		}
		return ports.AuthTokens{}, err
	}

	accessToken, err := s.tokenIssuer.IssueToken(ports.IssueTokenInput{UserID: rotatedSession.UserID, SessionID: rotatedSession.ID})
	if err != nil {
		return ports.AuthTokens{}, err
	}

	return ports.AuthTokens{
		AccessToken:      accessToken.Value,
		TokenType:        "Bearer",
		ExpiresAt:        accessToken.ExpiresAt,
		UserID:           rotatedSession.UserID,
		SessionID:        rotatedSession.ID,
		RefreshToken:     newRefreshToken,
		RefreshExpiresAt: rotatedSession.ExpiresAt,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, command ports.LogoutCommand) error {
	refreshToken := strings.TrimSpace(command.RefreshToken)
	if refreshToken == "" {
		return nil
	}

	session, err := s.sessionRepo.FindByTokenHash(ctx, s.refreshTokens.HashToken(refreshToken))
	if err != nil {
		if errors.Is(err, domain.ErrSessionNotFound) {
			return nil
		}
		return err
	}

	return s.sessionRepo.RevokeSession(ctx, session.ID, s.now())
}

func (s *AuthService) createSessionTokens(ctx context.Context, userID string, metadata requestMetadata) (ports.AuthTokens, error) {
	now := s.now()
	refreshToken, err := s.refreshTokens.GenerateToken()
	if err != nil {
		return ports.AuthTokens{}, fmt.Errorf("generate refresh token: %w", err)
	}

	session, err := s.sessionRepo.CreateSession(ctx, ports.CreateSessionInput{
		ID:                uuid.NewString(),
		UserID:            userID,
		TokenHash:         s.refreshTokens.HashToken(refreshToken),
		CreatedAt:         now,
		ExpiresAt:         now.Add(s.refreshTTL),
		AbsoluteExpiresAt: now.Add(s.refreshAbsoluteTTL),
		UserAgentHash:     metadata.userAgentHash,
		IPPrefix:          metadata.ipPrefix,
	})
	if err != nil {
		return ports.AuthTokens{}, err
	}

	accessToken, err := s.tokenIssuer.IssueToken(ports.IssueTokenInput{UserID: userID, SessionID: session.ID})
	if err != nil {
		return ports.AuthTokens{}, err
	}

	return ports.AuthTokens{
		AccessToken:      accessToken.Value,
		TokenType:        "Bearer",
		ExpiresAt:        accessToken.ExpiresAt,
		UserID:           userID,
		SessionID:        session.ID,
		RefreshToken:     refreshToken,
		RefreshExpiresAt: session.ExpiresAt,
	}, nil
}

type requestMetadata struct {
	userAgentHash *string
	ipPrefix      *string
}

func metadataFromRequest(userAgent string, ip string) requestMetadata {
	result := requestMetadata{}

	if trimmed := strings.TrimSpace(userAgent); trimmed != "" {
		hash := sha256Hex(trimmed)
		result.userAgentHash = &hash
	}

	if prefix := normalizeIPPrefix(ip); prefix != "" {
		result.ipPrefix = &prefix
	}

	return result
}

func normalizeIPPrefix(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}

	ip := net.ParseIP(trimmed)
	if ip == nil {
		return ""
	}

	if v4 := ip.To4(); v4 != nil {
		return fmt.Sprintf("%d.%d.%d", v4[0], v4[1], v4[2])
	}

	masked := ip.Mask(net.CIDRMask(64, 128))
	return masked.String()
}

func sha256Hex(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func sessionRootID(session domain.RefreshSession) string {
	if session.ParentID != nil && *session.ParentID != "" {
		return *session.ParentID
	}
	return session.ID
}
