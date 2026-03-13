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

type TursoSelection struct {
	Mode   string
	URL    string
	Scheme string
	Host   string
}

func BuildTursoDSN(cfg TursoConfig) (string, error) {
	selection, err := ResolveTursoSelection(cfg)
	if err != nil {
		return "", err
	}

	if cfg.Token == "" {
		return selection.URL, nil
	}

	if strings.Contains(selection.URL, "authToken=") {
		return selection.URL, nil
	}

	sep := "?"
	if strings.Contains(selection.URL, "?") {
		sep = "&"
	}

	return selection.URL + sep + "authToken=" + url.QueryEscape(cfg.Token), nil
}

func ResolveTursoSelection(cfg TursoConfig) (TursoSelection, error) {
	mode := strings.ToUpper(strings.TrimSpace(cfg.Mode))
	var selected string
	switch mode {
	case "", "HTTP":
		selected = strings.TrimSpace(cfg.URLHTTP)
		if selected == "" {
			return TursoSelection{}, errors.New("missing turso http url")
		}
	case "WS":
		selected = strings.TrimSpace(cfg.URLWS)
		if selected == "" {
			return TursoSelection{}, errors.New("missing turso ws url")
		}
	default:
		return TursoSelection{}, errors.New("invalid turso mode")
	}

	parsed, err := url.Parse(selected)
	if err != nil {
		return TursoSelection{}, err
	}

	return TursoSelection{
		Mode:   normalizeMode(mode),
		URL:    selected,
		Scheme: parsed.Scheme,
		Host:   parsed.Host,
	}, nil
}

func normalizeMode(mode string) string {
	if mode == "" {
		return "HTTP"
	}
	return mode
}
