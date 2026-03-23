package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Env      string
	Port     string
	LogLevel string

	// DB
	TursoCaloriesURLHTTP string
	TursoCaloriesToken   string
	TursoExpensesURLHTTP string
	TursoExpensesToken   string
	TursoHeatURLHTTP     string
	TursoHeatToken       string
	TursoUsersURLHTTP    string
	TursoUsersToken      string

	DBMaxOpenConns           int
	DBMaxIdleConns           int
	DBConnMaxLifetimeSeconds int
	DBConnMaxIdleTimeSeconds int

	JWTSecret        string
	AuthCookieName   string
	AuthCookieSecure bool
}

func Load() (Config, error) {
	var err error
	cfg := Config{
		Env:      getEnv("APP_ENV", "local"),
		Port:     getEnv("PORT", "8080"),
		LogLevel: getEnv("LOG_LEVEL", "info"),

		// DB
		TursoCaloriesURLHTTP: getEnv("TURSO_CALORIES_URL_HTTP", ""),
		TursoCaloriesToken:   getEnv("TURSO_CALORIES_TOKEN", ""),
		TursoExpensesURLHTTP: getEnv("TURSO_EXPENSES_URL_HTTP", ""),
		TursoExpensesToken:   getEnv("TURSO_EXPENSES_TOKEN", ""),
		TursoHeatURLHTTP:     getEnv("TURSO_HEAT_URL_HTTP", ""),
		TursoHeatToken:       getEnv("TURSO_HEAT_TOKEN", ""),
		TursoUsersURLHTTP:    getEnv("TURSO_USERS_URL_HTTP", ""),
		TursoUsersToken:      getEnv("TURSO_USERS_TOKEN", ""),

		JWTSecret:        getEnv("JWT_SECRET", ""),
		AuthCookieName:   getEnv("AUTH_COOKIE_NAME", "trackstack_session"),
		AuthCookieSecure: getEnvBool("AUTH_COOKIE_SECURE", true),
	}

	cfg.DBMaxOpenConns, err = getEnvInt("DB_MAX_OPEN_CONNS", 10)
	if err != nil {
		return Config{}, err
	}

	cfg.DBMaxIdleConns, err = getEnvInt("DB_MAX_IDLE_CONNS", 5)
	if err != nil {
		return Config{}, err
	}

	cfg.DBConnMaxLifetimeSeconds, err = getEnvInt("DB_CONN_MAX_LIFETIME_SECONDS", 300)
	if err != nil {
		return Config{}, err
	}

	cfg.DBConnMaxIdleTimeSeconds, err = getEnvInt("DB_CONN_MAX_IDLE_TIME_SECONDS", 60)
	if err != nil {
		return Config{}, err
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (cfg Config) Validate() error {
	if strings.TrimSpace(cfg.Port) == "" {
		return fmt.Errorf("Port must not be empty.")
	}

	if strings.TrimSpace(cfg.TursoCaloriesURLHTTP) == "" {
		return fmt.Errorf("TursoCaloriesURLHTTP must not be empty.")
	}

	if strings.TrimSpace(cfg.TursoCaloriesToken) == "" {
		return fmt.Errorf("TursoCaloriesToken must not be empty.")
	}

	if strings.TrimSpace(cfg.TursoExpensesURLHTTP) == "" {
		return fmt.Errorf("TursoExpensesURLHTTP must not be empty.")
	}

	if strings.TrimSpace(cfg.TursoExpensesToken) == "" {
		return fmt.Errorf("TursoExpensesToken must not be empty.")
	}

	if strings.TrimSpace(cfg.TursoHeatURLHTTP) == "" {
		return fmt.Errorf("TursoHeatURLHTTP must not be empty.")
	}

	if strings.TrimSpace(cfg.TursoHeatToken) == "" {
		return fmt.Errorf("TursoHeatToken must not be empty.")
	}

	if strings.TrimSpace(cfg.TursoUsersURLHTTP) == "" {
		return fmt.Errorf("TursoUsersURLHTTP must not be empty.")
	}

	if strings.TrimSpace(cfg.TursoUsersToken) == "" {
		return fmt.Errorf("TursoUsersToken must not be empty.")
	}

	if strings.TrimSpace(cfg.JWTSecret) == "" {
		return fmt.Errorf("JWT_SECRET must not be empty.")
	}
	if strings.TrimSpace(cfg.AuthCookieName) == "" {
		return fmt.Errorf("AUTH_COOKIE_NAME must not be empty.")
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
