.PHONY: backend-guard

backend-guard:
	cd apps/server && go run ./cmd/archguard && go run ./cmd/lintguard && go test ./...
