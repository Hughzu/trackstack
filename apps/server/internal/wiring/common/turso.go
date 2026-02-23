package common

import (
	"errors"
	"net/url"
	"strings"
)

type TursoConfig struct {
	Mode    string
	URL     string
	URLHTTP string
	URLWS   string
	Token   string
}

func BuildTursoDSN(cfg TursoConfig) (string, error) {
	urlValue := selectURL(cfg)
	if urlValue == "" {
		return "", errors.New("missing turso url")
	}

	if cfg.Token == "" {
		return urlValue, nil
	}

	if strings.Contains(urlValue, "authToken=") {
		return urlValue, nil
	}

	sep := "?"
	if strings.Contains(urlValue, "?") {
		sep = "&"
	}

	return urlValue + sep + "authToken=" + url.QueryEscape(cfg.Token), nil
}

func selectURL(cfg TursoConfig) string {
	mode := strings.ToUpper(strings.TrimSpace(cfg.Mode))
	if mode == "HTTP" && cfg.URLHTTP != "" {
		return cfg.URLHTTP
	}
	if mode == "WS" && cfg.URLWS != "" {
		return cfg.URLWS
	}
	return cfg.URL
}
