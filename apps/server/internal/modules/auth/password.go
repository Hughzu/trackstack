package auth

import (
	"crypto/subtle"
	"encoding/base64"
	"strconv"
	"strings"

	"golang.org/x/crypto/scrypt"
)

const (
	defaultScryptN = 16384
	defaultScryptR = 8
	defaultScryptP = 1
)

func VerifyPassword(password string, storedHash string) bool {
	if !strings.HasPrefix(storedHash, "$scrypt$") {
		return false
	}

	parts := strings.Split(storedHash, "$")
	if len(parts) < 5 {
		return false
	}

	params := parseParams(parts[2])
	saltEncoded := parts[3]
	hashEncoded := parts[4]
	if saltEncoded == "" || hashEncoded == "" {
		return false
	}

	salt, err := base64.StdEncoding.DecodeString(saltEncoded)
	if err != nil {
		return false
	}
	expected, err := base64.StdEncoding.DecodeString(hashEncoded)
	if err != nil {
		return false
	}

	derived, err := scrypt.Key([]byte(password), salt, params.N, params.R, params.P, len(expected))
	if err != nil {
		return false
	}

	return subtle.ConstantTimeCompare(derived, expected) == 1
}

type scryptParams struct {
	N int
	R int
	P int
}

func parseParams(segment string) scryptParams {
	params := scryptParams{N: defaultScryptN, R: defaultScryptR, P: defaultScryptP}
	for _, part := range strings.Split(segment, ",") {
		pieces := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(pieces) != 2 {
			continue
		}
		value, err := strconv.Atoi(strings.TrimSpace(pieces[1]))
		if err != nil {
			continue
		}
		switch pieces[0] {
		case "N":
			params.N = value
		case "r":
			params.R = value
		case "p":
			params.P = value
		}
	}
	return params
}
