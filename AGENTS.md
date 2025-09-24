# Agent Instructions for Go HTMX Web App

## Build/Test Commands
- **Build**: `go build`
- **Test all**: `go test ./...`
- **Test single**: `go test -run TestName ./package`
- **Lint**: `go vet ./...`
- **Format Go**: `gofmt -w .`
- **Format Templ**: `templ fmt`
- **Generate Templ**: `templ generate`

## Code Style Guidelines
- **Formatting**: Use `gofmt` for Go files, `templ fmt` for .templ files
- **Imports**: Group standard library, third-party, then local imports
- **Naming**: PascalCase for exported identifiers, camelCase for unexported
- **Error Handling**: Return errors, use `errors.Is()` for checking
- **Structs**: Align fields for readability, use meaningful names
- **HTTP Handlers**: Follow standard Go patterns with ResponseWriter and Request
- **Templating**: Use context.Background() for rendering, keep templates clean
- **Sessions**: Use secure random IDs, HTTP-only cookies, 24h expiration