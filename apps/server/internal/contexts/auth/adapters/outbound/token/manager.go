package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/auth/application/ports"
)

type Manager struct{}

var _ ports.RefreshTokenManager = (*Manager)(nil)

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) GenerateToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func (m *Manager) HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
