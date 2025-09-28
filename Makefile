.PHONY: build run clean test tidy

# Build the application
build:
	go build -o line-oa-backend main.go

# Run the application
run:
	go run main.go

# Clean build artifacts
clean:
	rm -f line-oa-backend

# Run tests
test:
	go test ./...

# Tidy dependencies
tidy:
	go mod tidy

# Install dependencies
deps:
	go mod download

# Run with hot reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	air

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Run all checks
check: fmt vet test

# Build for production
build-prod:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o line-oa-backend main.go
