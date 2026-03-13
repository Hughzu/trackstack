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
	TursoCaloriesURL                string
	TursoCaloriesURLHTTP            string
	TursoCaloriesURLWS              string
	TursoCaloriesToken              string
	TursoExpensesURL                string
	TursoExpensesURLHTTP            string
	TursoExpensesURLWS              string
	TursoExpensesToken              string
	TursoHeatURL                    string
	TursoHeatURLHTTP                string
	TursoHeatURLWS                  string
	TursoHeatToken                  string
	TursoUsersURL                   string
	TursoUsersURLHTTP               string
	TursoUsersURLWS                 string
	TursoUsersToken                 string
	OriginVerifyHeader              string
	OriginVerifyValue               string
	CORSAllowedOrigin               string
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
	cfg.TursoCaloriesURL = getEnv("TURSO_CALORIES_URL", "")
	cfg.TursoCaloriesURLHTTP = getEnv("TURSO_CALORIES_URL_HTTP", "")
	cfg.TursoCaloriesURLWS = getEnv("TURSO_CALORIES_URL_WS", "")
	cfg.TursoCaloriesToken = getEnv("TURSO_CALORIES_TOKEN", "")
	cfg.TursoExpensesURL = getEnv("TURSO_EXPENSES_URL", "")
	cfg.TursoExpensesURLHTTP = getEnv("TURSO_EXPENSES_URL_HTTP", "")
	cfg.TursoExpensesURLWS = getEnv("TURSO_EXPENSES_URL_WS", "")
	cfg.TursoExpensesToken = getEnv("TURSO_EXPENSES_TOKEN", "")
	cfg.TursoHeatURL = getEnv("TURSO_HEAT_URL", "")
	cfg.TursoHeatURLHTTP = getEnv("TURSO_HEAT_URL_HTTP", "")
	cfg.TursoHeatURLWS = getEnv("TURSO_HEAT_URL_WS", "")
	cfg.TursoHeatToken = getEnv("TURSO_HEAT_TOKEN", "")
	cfg.TursoUsersURL = getEnv("TURSO_USERS_URL", "")
	cfg.TursoUsersURLHTTP = getEnv("TURSO_USERS_URL_HTTP", "")
	cfg.TursoUsersURLWS = getEnv("TURSO_USERS_URL_WS", "")
	cfg.TursoUsersToken = getEnv("TURSO_USERS_TOKEN", "")
	cfg.OriginVerifyHeader = getEnv("ORIGIN_VERIFY_HEADER", "")
	cfg.OriginVerifyValue = getEnv("ORIGIN_VERIFY_VALUE", "")
	cfg.CORSAllowedOrigin = getEnv("CORS_ALLOWED_ORIGIN", "")
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
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (cfg Config) Validate() error {
	mode := strings.ToUpper(strings.TrimSpace(cfg.DBConnectionMode))
	switch mode {
	case "", "HTTP":
		if err := requireModeURLs("HTTP", [][2]string{
			{"TURSO_CALORIES_URL_HTTP", cfg.TursoCaloriesURLHTTP},
			{"TURSO_EXPENSES_URL_HTTP", cfg.TursoExpensesURLHTTP},
			{"TURSO_HEAT_URL_HTTP", cfg.TursoHeatURLHTTP},
			{"TURSO_USERS_URL_HTTP", cfg.TursoUsersURLHTTP},
		}); err != nil {
			return err
		}
	case "WS":
		if err := requireModeURLs("WS", [][2]string{
			{"TURSO_CALORIES_URL_WS", cfg.TursoCaloriesURLWS},
			{"TURSO_EXPENSES_URL_WS", cfg.TursoExpensesURLWS},
			{"TURSO_HEAT_URL_WS", cfg.TursoHeatURLWS},
			{"TURSO_USERS_URL_WS", cfg.TursoUsersURLWS},
		}); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid DB_CONNECTION_MODE %q: expected HTTP or WS", cfg.DBConnectionMode)
	}

	return nil
}

func requireModeURLs(mode string, values [][2]string) error {
	for _, item := range values {
		if strings.TrimSpace(item[1]) == "" {
			return fmt.Errorf("%s requires %s to be set", mode, item[0])
		}
	}
	return nil
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
