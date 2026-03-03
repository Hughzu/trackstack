package auth

import "time"

type ClientContext struct {
	UserAgent *string
	IPPrefix  *string
}

type Session struct {
	ID                string
	UserID            string
	CreatedAt         string
	ExpiresAt         string
	RotatedAt         string
	LastSeenAt        string
	AbsoluteExpiresAt string
	ParentID          *string
	RevokedAt         *string
	UserAgentHash     *string
	IPPrefix          *string
}

type Config struct {
	SessionIdleSeconds          int
	SessionAbsoluteSeconds      int
	SessionRotateAfterSeconds   int
	SessionRotationGraceSeconds int
	SessionTouchSeconds         int
}

type CreateSessionRequest struct {
	UserID  string
	Context ClientContext
}

type RevokeSessionRequest struct {
	RawToken string
}

type AuthenticateRequest struct {
	RawToken string
	Context  ClientContext
}

type AuthenticateResponse struct {
	UserID         string
	SessionID      string
	ReplacementRaw *string
	CookieExpires  *time.Time
}
