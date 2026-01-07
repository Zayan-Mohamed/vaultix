# Development Guide

Guide for developers who want to contribute to Vaultix or understand its internals.

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Git
- Basic understanding of Go and cryptography

### Clone and Build

```bash
git clone https://github.com/Zayan-Mohamed/vaultix.git
cd vaultix
go build -o vaultix
```

### Run Tests

```bash
go test ./...
```

### Install Locally

```bash
go install
```

## Project Structure

```
vaultix/
├── internal/
│   ├── crypto/       # Cryptographic operations
│   │   └── crypto.go
│   ├── storage/      # File system operations
│   │   └── storage.go
│   ├── vault/        # Business logic
│   │   └── vault.go
│   └── cli/          # Command-line interface
│       └── cli.go
├── docs/             # MkDocs documentation
├── .github/          # GitHub workflows and config
├── main.go           # Entry point
├── go.mod            # Go module definition
├── go.sum            # Dependency checksums
├── install.sh        # Linux/macOS installer
├── install.ps1       # Windows installer
├── mkdocs.yml        # MkDocs configuration
├── LICENSE           # MIT License
└── README.md         # Project readme
```

## Architecture

### Layer Overview

```
┌─────────────────────────────────────┐
│            CLI Layer                │
│  (User interaction, arg parsing)    │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│          Vault Layer                │
│  (Business logic, orchestration)    │
└──────────────┬──────────────────────┘
               │
       ┌───────┴────────┐
       │                │
┌──────▼──────┐  ┌──────▼──────┐
│   Crypto    │  │   Storage   │
│  (Encrypt)  │  │   (I/O)     │
└─────────────┘  └─────────────┘
```

### Package Responsibilities

#### `internal/crypto/`

Handles all cryptographic operations:

- Key derivation (Argon2id)
- Encryption/decryption (AES-256-GCM)
- Random number generation
- Salt generation

**Key Functions**:

- `GenerateSalt()` - Create random salt
- `DeriveKey(password, salt)` - Generate encryption key
- `Encrypt(plaintext, key)` - Encrypt data
- `Decrypt(ciphertext, key)` - Decrypt data

#### `internal/storage/`

Manages filesystem operations:

- Vault structure creation
- File reading/writing
- Object storage
- Metadata serialization
- Secure file deletion

**Key Functions**:

- `InitializeVault(path)` - Create vault directory structure
- `WriteObject(path, id, data)` - Write encrypted object
- `ReadObject(path, id)` - Read encrypted object
- `SecureDelete(path)` - Overwrite and delete file

#### `internal/vault/`

Business logic layer:

- Coordinates crypto + storage
- Manages metadata
- Implements vault operations (init, add, extract, etc.)
- Error handling

**Key Functions**:

- `Initialize(password)` - Create new vault
- `AddFile(password, path)` - Encrypt and add file
- `ExtractFile(password, name, dest)` - Decrypt file
- `ListFiles(password)` - Get file list

#### `internal/cli/`

Command-line interface:

- Argument parsing
- Command routing
- User prompts
- Password input
- Output formatting

**Key Functions**:

- `Init(args)` - Handle `init` command
- `Add(args)` - Handle `add` command
- `Extract(args)` - Handle `extract` command
- `readPassword()` - Secure password input

## Code Style

### Naming Conventions

- **Packages**: lowercase, single word (`crypto`, `vault`)
- **Files**: lowercase, descriptive (`crypto.go`, `storage.go`)
- **Functions**: CamelCase, exported start with capital (`Initialize`, `AddFile`)
- **Variables**: camelCase (`vaultPath`, `objectID`)
- **Constants**: CamelCase or UPPER_CASE (`saltFileName`, `AES_KEY_SIZE`)

### Error Handling

Always wrap errors with context:

```go
// Good
if err != nil {
    return fmt.Errorf("failed to read file: %w", err)
}

// Bad
if err != nil {
    return err
}
```

### Comments

- Comment all exported functions
- Explain "why", not "what"
- Use godoc format

```go
// DeriveKey generates an encryption key from a password using Argon2id.
// The salt must be 32 bytes long and randomly generated.
// Returns a 32-byte key suitable for AES-256 encryption.
func DeriveKey(password string, salt []byte) ([]byte, error) {
    // ...
}
```

## Adding a New Command

1. **Add CLI handler** in `internal/cli/cli.go`:

```go
func MyCommand(args []string) error {
    // Parse arguments
    // Read password
    // Call vault method
    // Handle errors
    // Print success message
    return nil
}
```

2. **Add vault method** in `internal/vault/vault.go`:

```go
func (v *Vault) MyOperation(password string, arg string) error {
    // Get salt and derive key
    // Read metadata
    // Perform operation
    // Update metadata
    // Return result
    return nil
}
```

3. **Register command** in `main.go`:

```go
switch command {
case "mycommand":
    err = cli.MyCommand(args)
// ...
}
```

4. **Update help text** in `cli.PrintUsage()`:

```go
fmt.Println("  vaultix mycommand <arg>    Description of command")
```

5. **Add documentation** in `docs/commands.md`

6. **Write tests**

## Testing

### Unit Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/crypto

# Verbose output
go test -v ./...
```

### Integration Tests

Create test vaults in `/tmp`:

```bash
#!/bin/bash
cd /tmp
mkdir test_vault
cd test_vault
echo "test" > file.txt

# Test init
expect << 'EOF'
spawn vaultix init
expect "Enter password:"
send "test123\r"
expect "Confirm password:"
send "test123\r"
expect eof
EOF

# Test list
expect << 'EOF'
spawn vaultix list
expect "Enter vault password:"
send "test123\r"
expect eof
EOF
```

### Manual Testing

```bash
# Build
go build

# Test in temp directory
cd /tmp/test_vault
~/vaultix/vaultix init
~/vaultix/vaultix list
```

## Dependencies

Vaultix has minimal external dependencies:

```go
require (
    golang.org/x/crypto v0.x.x  // Argon2id
    golang.org/x/term v0.x.x    // Password input
)
```

### Updating Dependencies

```bash
go get -u ./...
go mod tidy
```

## Release Process

1. Update version in code (if versioned)
2. Update CHANGELOG.md
3. Run tests: `go test ./...`
4. Build for all platforms
5. Create git tag: `git tag v1.0.0`
6. Push tag: `git push origin v1.0.0`
7. Create GitHub release with binaries

## CI/CD

GitHub Actions workflows (`.github/workflows/`):

- **build.yml**: Build and test on push
- **release.yml**: Build binaries on tag
- **docs.yml**: Deploy MkDocs to GitHub Pages

## Documentation

Documentation uses MkDocs with Material theme.

### Local Preview

```bash
# Install MkDocs
pip install mkdocs-material

# Serve locally
mkdocs serve

# Open http://localhost:8000
```

### Deploy Documentation

```bash
mkdocs gh-deploy
```

## Debugging

### Enable Debug Output

Add debug prints (remove before committing):

```go
fmt.Fprintf(os.Stderr, "DEBUG: key length = %d\n", len(key))
```

### Common Issues

**Build errors**:

```bash
go mod tidy
go clean -cache
```

**Import errors**:

```bash
go get golang.org/x/crypto
go mod tidy
```

## Security Considerations

When contributing:

- ✓ Never add custom crypto
- ✓ Use standard library when possible
- ✓ Validate all inputs
- ✓ Clear sensitive data from memory
- ✓ Test error paths
- ✗ Don't log passwords or keys
- ✗ Don't store sensitive data

## Questions?

- Open an issue on GitHub
- Check existing issues
- Read the documentation
- Review the code

Happy coding! 🚀
