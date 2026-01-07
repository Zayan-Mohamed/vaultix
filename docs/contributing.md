# Contributing to Vaultix

Thank you for your interest in contributing to vaultix! This document provides guidelines and instructions for contributing.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Documentation](#documentation)
- [Pull Request Process](#pull-request-process)
- [Release Process](#release-process)

## Code of Conduct

### Our Pledge

We pledge to make participation in our project a harassment-free experience for everyone, regardless of age, body size, disability, ethnicity, gender identity and expression, level of experience, nationality, personal appearance, race, religion, or sexual identity and orientation.

### Our Standards

✓ **Examples of behavior that contributes to a positive environment:**

- Using welcoming and inclusive language
- Being respectful of differing viewpoints
- Gracefully accepting constructive criticism
- Focusing on what is best for the community
- Showing empathy towards other community members

✗ **Examples of unacceptable behavior:**

- Trolling, insulting/derogatory comments, and personal attacks
- Public or private harassment
- Publishing others' private information without permission
- Other conduct which could reasonably be considered inappropriate

## Getting Started

### Prerequisites

- Go 1.21 or later
- Git
- Basic understanding of cryptography (helpful but not required)

### Development Tools

Recommended:

- [golangci-lint](https://golangci-lint.run/) for linting
- [gopls](https://github.com/golang/tools/tree/master/gopls) for IDE support
- [delve](https://github.com/go-delve/delve) for debugging

## Development Setup

### 1. Fork the Repository

```bash
# Click "Fork" on GitHub, then:
git clone https://github.com/YOUR_USERNAME/vaultix.git
cd vaultix
```

### 2. Add Upstream Remote

```bash
git remote add upstream https://github.com/zayan-mohamed/vaultix.git
git fetch upstream
```

### 3. Create a Branch

```bash
git checkout -b feature/your-feature-name
```

### 4. Build and Test

```bash
# Build
go build -o vaultix

# Run tests
go test ./...

# Run with race detector
go test -race ./...
```

## How to Contribute

### Reporting Bugs

Before creating a bug report:

1. **Search existing issues** to see if the problem has already been reported
2. **Check if the problem is reproducible** in the latest version
3. **Collect information** about the bug

**Good bug report includes:**

- Clear, descriptive title
- Steps to reproduce
- Expected behavior
- Actual behavior
- Environment (OS, Go version, vaultix version)
- Error messages or logs
- Screenshots (if applicable)

**Example bug report:**

```text
**Title:** vaultix crashes when adding files larger than 2GB

**Description:**
When trying to add a file larger than 2GB, vaultix crashes with "out of memory" error.

**Steps to Reproduce:**

1. Create a 3GB test file: `fallocate -l 3G bigfile.bin`
2. Initialize vault: `vaultix init`
3. Add file: `vaultix add bigfile.bin`

**Expected Behavior:**
File should be encrypted successfully.

**Actual Behavior:**
Program crashes with error:
```

fatal error: runtime: out of memory

```

**Environment:**
- OS: Ubuntu 22.04
- Go version: 1.21.5
- vaultix version: v0.1.0
- RAM: 8GB
```

### Suggesting Features

Before suggesting a feature:

1. **Check if it aligns** with project goals (see the project philosophy in the main documentation)
2. **Search existing issues** for similar proposals
3. **Consider the scope** - does this belong in vaultix?

**Good feature request includes:**

- Clear description of the problem it solves
- Proposed solution
- Alternatives considered
- Why this belongs in vaultix (not a plugin/wrapper)

**Out of scope features:**

- GUI interfaces
- Cloud synchronization
- Password manager UI
- Network features
- Operating system integration beyond basic CLI

### Writing Code

#### Before You Start

1. **Discuss major changes** by opening an issue first
2. **Check existing PRs** to avoid duplicate work
3. **Review the architecture** in [architecture.md](architecture.md)

#### While Coding

✓ **DO:**

- Follow existing code style
- Write clear commit messages
- Add tests for new functionality
- Update documentation
- Keep changes focused and atomic

✗ **DON'T:**

- Mix multiple features in one PR
- Break existing functionality
- Add dependencies unnecessarily
- Ignore test failures
- Commit sensitive data or test vaults

## Coding Standards

### Go Style

Follow [Effective Go](https://golang.org/doc/effective_go.html) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).

**Key points:**

```go
// Good: Clear, concise names
func encryptFile(path string, key []byte) error

// Bad: Unclear abbreviations
func encF(p string, k []byte) error

// Good: Early returns
func process() error {
    if err := validate(); err != nil {
        return err
    }
    return execute()
}

// Bad: Deep nesting
func process() error {
    if validate() == nil {
        if execute() == nil {
            return nil
        }
    }
    return errors.New("failed")
}

// Good: Error context
return fmt.Errorf("failed to encrypt %s: %w", filename, err)

// Bad: Unclear errors
return err
```

### Code Organization

```
vaultix/
├── cmd/               # CLI commands (init, add, etc.)
├── crypto/            # Cryptographic operations
├── storage/           # File system operations
├── vault/             # Vault management
└── main.go
```

**Rules:**

- One package per directory
- Clear package responsibilities
- Minimal inter-package dependencies
- No circular dependencies

### Security Standards

**Critical rules:**

1. **No crypto invention**: Use standard library only
2. **No password logging**: Never log sensitive data
3. **Clear error messages**: Don't leak vault contents in errors
4. **Secure cleanup**: Wipe sensitive data from memory
5. **Input validation**: Sanitize all user inputs

**Example:**

```go
// Good: Use standard crypto
import "crypto/aes"
import "crypto/cipher"

cipher, err := aes.NewCipher(key)
gcm, err := cipher.NewGCM(cipher)

// Bad: Custom crypto
func myEncryption(data []byte) []byte {
    // DON'T DO THIS
}

// Good: Clear sensitive data
defer func() {
    for i := range password {
        password[i] = 0
    }
}()

// Good: Don't leak info in errors
return fmt.Errorf("decryption failed")

// Bad: Leaks information
return fmt.Errorf("decryption failed: invalid AES key for file %s", filename)
```

### Comments

```go
// Good: Explain why, not what
// Hash password using Argon2id to derive encryption key.
// We use Argon2id over bcrypt because it's memory-hard,
// protecting against GPU-based attacks.
key := argon2.IDKey(...)

// Bad: Obvious comment
// Create a new Argon2 key
key := argon2.IDKey(...)

// Good: Package documentation
// Package crypto provides cryptographic operations for vaultix.
// All encryption uses AES-256-GCM with keys derived from Argon2id.
package crypto

// Good: Exported function documentation
// EncryptFile encrypts the file at path using the provided key.
// Returns the encrypted data and metadata required for decryption.
//
// The key must be 32 bytes (AES-256).
// Encryption uses AES-256-GCM with a randomly generated nonce.
func EncryptFile(path string, key []byte) ([]byte, error)
```

## Testing

### Unit Tests

**Every package must have tests.**

```go
// Good test structure
func TestEncryptFile(t *testing.T) {
    tests := []struct {
        name    string
        path    string
        key     []byte
        wantErr bool
    }{
        {
            name: "valid file and key",
            path: "testdata/sample.txt",
            key:  make([]byte, 32),
        },
        {
            name:    "invalid key length",
            path:    "testdata/sample.txt",
            key:     make([]byte, 16),
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := EncryptFile(tt.path, tt.key)
            if (err != nil) != tt.wantErr {
                t.Errorf("EncryptFile() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Integration Tests

Test real workflows:

```bash
#!/bin/bash
# test_integration.sh

# Setup
TEMP=$(mktemp -d)
cd "$TEMP"

# Test init
echo "test" > file.txt
echo "password123" | ../vaultix init
if [ $? -ne 0 ]; then
    echo "FAIL: init"
    exit 1
fi

# Test list
echo "password123" | ../vaultix list | grep "file.txt"
if [ $? -ne 0 ]; then
    echo "FAIL: list"
    exit 1
fi

# Cleanup
cd -
rm -rf "$TEMP"

echo "PASS: integration tests"
```

### Running Tests

```bash
# All tests
go test ./...

# With coverage
go test -cover ./...

# Race detector
go test -race ./...

# Specific package
go test ./crypto

# Verbose
go test -v ./...

# Integration tests
./test_integration.sh
```

### Test Coverage

Aim for:

- **90%+** coverage for crypto package
- **80%+** coverage for vault package
- **70%+** coverage overall

Check coverage:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Documentation

### Code Documentation

- All exported functions must have doc comments
- Package-level documentation required
- Complex algorithms need explanation

### User Documentation

When adding features, update:

- [README.md](../README.md) - If it affects quick start
- [commands.md](commands.md) - New commands or flags
- [usage.md](usage.md) - Usage patterns
- [examples.md](examples.md) - Real-world examples

### Writing Style

- **Clear and concise**: No jargon unless necessary
- **Examples**: Show, don't just tell
- **Honest**: Don't oversell security claims
- **Accurate**: Test all examples

## Pull Request Process

### 1. Before Submitting

Checklist:

- [ ] Code follows style guidelines
- [ ] Tests added and passing
- [ ] Documentation updated
- [ ] Commit messages are clear
- [ ] Branch is up-to-date with main

```bash
# Update your branch
git fetch upstream
git rebase upstream/main
```

### 2. PR Description

**Good PR description:**

```markdown
## Summary

Add support for extracting multiple files at once.

## Motivation

Users often need to extract several related files. Currently this requires multiple commands.

## Changes

- Modified `extract` command to accept multiple filenames
- Added fuzzy matching for each filename
- Updated tests and documentation

## Testing

- Added unit tests for multi-file extraction
- Tested with 1, 5, and 100 files
- Verified fuzzy matching works for all files

## Related Issues

Fixes #42
```

### 3. Review Process

**What reviewers look for:**

- Code quality and style
- Test coverage
- Security implications
- Documentation completeness
- Performance impact

**Expected turnaround:**

- Initial review within 1-3 days
- Discussion and revisions as needed
- Merge when approved by maintainer

### 4. After Merge

- Your code will be included in the next release
- Update your fork
- Close related issues (if not auto-closed)

## Release Process

(For maintainers)

### Version Numbering

Follow [Semantic Versioning](https://semver.org/):

- **Major** (v1.0.0): Breaking changes
- **Minor** (v0.1.0): New features, backwards compatible
- **Patch** (v0.0.1): Bug fixes

### Release Checklist

1. **Update CHANGELOG.md**
2. **Tag release**:
   ```bash
   git tag -a v0.2.0 -m "Release v0.2.0"
   git push origin v0.2.0
   ```
3. **Build binaries** for all platforms
4. **Create GitHub release** with binaries attached
5. **Update documentation** site
6. **Announce** in discussions

## Getting Help

- **Questions**: Open a discussion on GitHub
- **Bugs**: Open an issue
- **Security**: Email maintainers directly (see SECURITY.md)
- **Chat**: (Coming soon)

## Recognition

Contributors will be:

- Listed in CONTRIBUTORS.md
- Mentioned in release notes
- Thanked publicly

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to vaultix! 🎉
