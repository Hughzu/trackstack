package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/23St/trackstack/internal/common/db"
)

// IPWhitelistMiddleware validates client IP against allowed mappings
func IPWhitelistMiddleware(database *db.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// DEV_MODE bypass: allow all IPs in development
			if os.Getenv("DEV_MODE") == "true" {
				slog.Debug("DEV_MODE enabled, bypassing IP whitelist")
				next.ServeHTTP(w, r)
				return
			}

			// Always allow health checks
			if r.URL.Path == "/health" {
				next.ServeHTTP(w, r)
				return
			}

			// Extract real client IP
			clientIP := extractClientIP(r)
			slog.Debug("checking IP whitelist", "ip", clientIP, "path", r.URL.Path)

			// Look up user by IP
			user, err := database.GetUserByIP(ctx, clientIP)
			if err != nil {
				slog.Warn("unauthorized IP access attempt",
					"ip", clientIP,
					"path", r.URL.Path,
					"user_agent", r.UserAgent(),
				)
				http.Error(w, "Forbidden: Unauthorized device", http.StatusForbidden)
				return
			}

			// Update last seen
			if err := database.UpdateUserLastSeen(ctx, user.ID); err != nil {
				slog.Error("failed to update last seen", "error", err, "user_id", user.ID)
			}

			// Check for existing session cookie
			var sessionID string
			cookie, err := r.Cookie("trackstack_session")
			if err == nil && cookie.Value != "" {
				// Validate session
				session, err := database.GetSessionByID(ctx, cookie.Value)
				if err == nil && session.UserID == user.ID {
					sessionID = session.ID
				} else {
					slog.Debug("invalid or expired session, creating new one",
						"user_id", user.ID,
						"old_session", cookie.Value,
					)
				}
			}

			// Create new session if needed
			if sessionID == "" {
				session := db.NewSession(user.ID)
				if err := database.CreateSession(ctx, session); err != nil {
					slog.Error("failed to create session", "error", err, "user_id", user.ID)
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
				}
				sessionID = session.ID
				http.SetCookie(w, createSessionCookie(session.ID))
				slog.Info("created new session",
					"user_id", user.ID,
					"session_id", session.ID,
					"ip", clientIP,
				)
			}

			// Inject user into context
			ctx = context.WithValue(ctx, UserContextKey, user)

			// Continue to next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// extractClientIP gets the real client IP, handling reverse proxies
func extractClientIP(r *http.Request) string {
	var ip string

	// Check X-Forwarded-For header (used by reverse proxies like Caddy, Nginx, AWS LB)
	// Format: "client, proxy1, proxy2"
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// Get the first IP (the actual client)
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			ip = strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-Ip header (some proxies use this)
	if ip == "" {
		xri := r.Header.Get("X-Real-Ip")
		if xri != "" {
			ip = strings.TrimSpace(xri)
		}
	}

	// Fallback to RemoteAddr (direct connection)
	// Format: "ip:port" or "[ipv6]:port", we need to strip the port
	if ip == "" {
		remoteAddr := r.RemoteAddr
		ip = stripPort(remoteAddr)
	}

	return normalizeIP(ip)
}

// stripPort removes the port from an address
// Handles both "ip:port" and "[ipv6]:port" formats
func stripPort(addr string) string {
	// Check if it's IPv6 with brackets: [::1]:8080
	if strings.HasPrefix(addr, "[") {
		// Find closing bracket
		if idx := strings.Index(addr, "]"); idx != -1 {
			// Return content inside brackets
			return addr[1:idx]
		}
	}

	// IPv4 or plain format: 127.0.0.1:8080
	if idx := strings.LastIndex(addr, ":"); idx != -1 {
		return addr[:idx]
	}

	return addr
}

// normalizeIP converts common IP representations to a standard format
func normalizeIP(ip string) string {
	ip = strings.TrimSpace(ip)

	// Handle IPv6 localhost variants
	if ip == "::1" || ip == "[::1]" || ip == "0:0:0:0:0:0:0:1" {
		return "127.0.0.1"
	}

	// Remove brackets from IPv6 addresses
	ip = strings.TrimPrefix(ip, "[")
	ip = strings.TrimSuffix(ip, "]")

	return ip
}
