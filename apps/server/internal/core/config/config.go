package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Env                             string
	Port                            string
	LogLevel                        string
	DBConnectionMode                string
	TursoHeatURL                    string
	TursoHeatURLHTTP                string
	TursoHeatURLWS                  string
	TursoHeatToken                  string
	HardcodedUserID                 string
	OriginVerifyHeader              string
	OriginVerifyValue               string
	AuthCookieName                  string
	AuthCookieSecure                bool
	AuthCookieSameSite              string
	AuthSessionIdleSeconds          int
	AuthSessionAbsoluteSeconds      int
	AuthSessionRotateAfterSeconds   int
	AuthSessionRotationGraceSeconds int
	AuthSessionTouchSeconds         int
}

func Load() (Config, error) {
	var cfg Config
	cfg.Env = getEnv("APP_ENV", "local")
	cfg.Port = getEnv("PORT", "8080")
	cfg.LogLevel = getEnv("LOG_LEVEL", "info")
	cfg.DBConnectionMode = getEnv("DB_CONNECTION_MODE", "HTTP")
	cfg.TursoHeatURL = getEnv("TURSO_HEAT_URL", "")
	cfg.TursoHeatURLHTTP = getEnv("TURSO_HEAT_URL_HTTP", "")
	cfg.TursoHeatURLWS = getEnv("TURSO_HEAT_URL_WS", "")
	cfg.TursoHeatToken = getEnv("TURSO_HEAT_TOKEN", "")
	cfg.HardcodedUserID = getEnv("HARD_CODED_USER_ID", "")
	cfg.OriginVerifyHeader = getEnv("ORIGIN_VERIFY_HEADER", "")
	cfg.OriginVerifyValue = getEnv("ORIGIN_VERIFY_VALUE", "")
	cfg.AuthCookieName = getEnv("AUTH_COOKIE_NAME", "trackstack_session")
	cfg.AuthCookieSecure = getEnvBool("AUTH_COOKIE_SECURE", false)
	cfg.AuthCookieSameSite = getEnv("AUTH_COOKIE_SAMESITE", "lax")

	var err error
	cfg.AuthSessionIdleSeconds, err = getEnvInt("AUTH_SESSION_IDLE_SECONDS", 1800)
	if err != nil {
		return Config{}, err
	}
	cfg.AuthSessionAbsoluteSeconds, err = getEnvInt("AUTH_SESSION_ABSOLUTE_SECONDS", 86400)
	if err != nil {
		return Config{}, err
	}
	cfg.AuthSessionRotateAfterSeconds, err = getEnvInt("AUTH_SESSION_ROTATE_AFTER_SECONDS", 1800)
	if err != nil {
		return Config{}, err
	}
	cfg.AuthSessionRotationGraceSeconds, err = getEnvInt("AUTH_SESSION_ROTATION_GRACE_SECONDS", 300)
	if err != nil {
		return Config{}, err
	}
	cfg.AuthSessionTouchSeconds, err = getEnvInt("AUTH_SESSION_TOUCH_SECONDS", 300)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func getEnvInt(key string, fallback int) (int, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	return parsed, nil
}
