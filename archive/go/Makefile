.PHONY: templ css build run clean

# Generate Go code from Templ files
templ:
	templ generate

# Build Tailwind CSS
css:
	./bin/tailwindcss -i static/css/input.css -o static/css/output.css --minify

# Build the binary
build: templ css
	go build -o bin/trackstack ./cmd/monolith

# Run the server
run: build
	./bin/trackstack

# Clean generated files
clean:
	rm -rf bin/trackstack
	rm -f static/css/output.css
	find . -name "*_templ.go" -delete

# Development mode with watch (requires air)
dev:
	templ generate --watch &
	./bin/tailwindcss -i static/css/input.css -o static/css/output.css --watch
