# Project: GoContextPacker
- Language: Go 1.22+
- Style: Idiomatic Go (Effective Go). Use `camelCase` for internal, `PascalCase` for exported.
- Error Handling: Return errors, don't panic. Wrap errors with context (fmt.Errorf("%w")).
- Architecture: CLI driven. Logic in `internal/`, entry in `cmd/`.

# Commands
- Build: `go build -o bin/gopack ./cmd/gopack`
- Run: `go run ./cmd/gopack`
- Test: `go test ./...`
- Lint: `golangci-lint run` (if available)

# Workflow rules
- Always run `go fmt ./...` after modifying code.
- If adding external dependencies, run `go mod tidy`.
