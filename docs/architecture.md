# Architecture

Understanding vaultix's internal structure and design.

## High-Level Overview

```
┌─────────────────────────────────────────┐
│           User (CLI)                     │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│         Command Layer (cmd/)             │
│  init, add, list, extract, drop, etc.   │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│        Vault Layer (vault/)              │
│  High-level vault operations             │
└──────────┬─────────┬────────────────────┘
           │         │
           ▼         ▼
    ┌──────────┐  ┌─────────────┐
    │ Crypto   │  │  Storage    │
    │ (crypto/)│  │ (storage/)  │
    └──────────┘  └─────────────┘
```

## Directory Structure

```
vaultix/
├── main.go              # Entry point
├── cmd/                 # CLI commands
│   ├── init.go
│   ├── add.go
│   ├── list.go
│   ├── extract.go
│   ├── drop.go
│   ├── remove.go
│   └── clear.go
├── vault/               # Vault management
│   ├── vault.go         # Core vault operations
│   ├── metadata.go      # Metadata handling
│   └── file.go          # File operations
├── crypto/              # Cryptographic operations
│   ├── encrypt.go
│   ├── decrypt.go
│   └── key.go
├── storage/             # File system operations
│   ├── storage.go
│   └── paths.go
└── docs/                # Documentation
```

## Component Responsibilities

### Main Entry Point

**File**: `main.go`

**Responsibilities**:

- Parse command-line arguments
- Route to appropriate command handler
- Handle global flags (--version, --help)
- Set up error handling

**Example**:

```go
func main() {
    if len(os.Args) < 2 {
        showUsage()
        os.Exit(1)
    }

    command := os.Args[1]

    switch command {
    case "init":
        cmd.Init(os.Args[2:])
    case "add":
        cmd.Add(os.Args[2:])
    // ...
    }
}
```

### Command Layer (`cmd/`)

**Responsibilities**:

- Parse command-specific arguments
- Validate user input
- Prompt for password
- Call vault layer operations
- Display results to user
- Handle errors and provide user-friendly messages

**Example** (`cmd/init.go`):

```go
func Init(args []string) error {
    // Parse arguments
    path := parsePathOrDefault(args)

    // Validate directory
    if err := validateDirectory(path); err != nil {
        return err
    }

    // Get password from user
    password, err := promptPassword()
    if err != nil {
        return err
    }

    // Initialize vault
    v := vault.New(path)
    return v.Init(password)
}
```

**Commands**:

- `init.go` - Initialize new vault
- `add.go` - Add files to vault
- `list.go` - List encrypted files
- `extract.go` - Decrypt and extract files
- `drop.go` - Extract and remove files
- `remove.go` - Remove files from vault
- `clear.go` - Clear entire vault

### Vault Layer (`vault/`)

**Responsibilities**:

- High-level vault operations
- Manage vault state
- Coordinate between crypto and storage
- Handle metadata
- Enforce vault invariants

**Core Types**:

```go
// Vault represents an encrypted directory
type Vault struct {
    Path     string
    Metadata *Metadata
}

// Metadata stores encrypted file information
type Metadata struct {
    Files []FileMetadata
}

// FileMetadata describes an encrypted file
type FileMetadata struct {
    ID       string    // Random identifier
    Name     string    // Original filename
    Size     int64     // Original size
    Modified time.Time // Last modified
}
```

**Key Functions**:

```go
// Initialize new vault
func (v *Vault) Init(password string) error

// Add file to vault
func (v *Vault) AddFile(path string, password string) error

// List all files
func (v *Vault) ListFiles(password string) ([]FileMetadata, error)

// Extract file
func (v *Vault) ExtractFile(id string, password string) error

// Remove file
func (v *Vault) RemoveFile(id string, password string) error
```

### Crypto Layer (`crypto/`)

**Responsibilities**:

- All cryptographic operations
- Key derivation
- Encryption/decryption
- No business logic - pure crypto functions

**Key Functions**:

```go
// Derive encryption key from password
func DeriveKey(password string, salt []byte) []byte {
    return argon2.IDKey(
        []byte(password),
        salt,
        3,      // time parameter
        64*1024, // memory parameter (64 MB)
        4,      // parallelism
        32,     // key length (AES-256)
    )
}

// Encrypt data using AES-256-GCM
func Encrypt(plaintext, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := rand.Read(nonce); err != nil {
        return nil, err
    }

    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
    return ciphertext, nil
}

// Decrypt data using AES-256-GCM
func Decrypt(ciphertext, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    nonceSize := gcm.NonceSize()
    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, ErrDecryptionFailed
    }

    return plaintext, nil
}
```

**Security Principles**:

- Never log sensitive data
- Clear sensitive data from memory when done
- Use crypto/rand for all random data
- No custom cryptography
- Fail secure (return error, don't continue)

### Storage Layer (`storage/`)

**Responsibilities**:

- File system operations
- Path handling
- Secure file deletion
- Cross-platform compatibility

**Key Functions**:

```go
// Read file from disk
func ReadFile(path string) ([]byte, error)

// Write file to disk
func WriteFile(path string, data []byte) error

// Securely delete file (overwrite then remove)
func SecureDelete(path string) error {
    // Open file
    f, err := os.OpenFile(path, os.O_WRONLY, 0)
    if err != nil {
        return err
    }
    defer f.Close()

    // Get file size
    info, err := f.Stat()
    if err != nil {
        return err
    }

    // Overwrite with random data
    random := make([]byte, info.Size())
    rand.Read(random)
    f.Write(random)

    // Close and remove
    f.Close()
    return os.Remove(path)
}

// Get vault directory paths
func GetVaultPaths(basePath string) *Paths {
    return &Paths{
        Base:    basePath,
        Vaultix: filepath.Join(basePath, ".vaultix"),
        Objects: filepath.Join(basePath, ".vaultix", "objects"),
        Salt:    filepath.Join(basePath, ".vaultix", "salt"),
        Meta:    filepath.Join(basePath, ".vaultix", "meta"),
    }
}
```

## Data Flow

### Initialization Flow

```
User runs: vaultix init

1. [cmd/init.go]
   ├─ Parse arguments (path)
   ├─ Prompt for password
   └─ Call vault.Init()

2. [vault/vault.go]
   ├─ Create .vaultix/ directory structure
   ├─ Generate random salt
   ├─ Derive key from password + salt
   ├─ Scan directory for files
   └─ For each file:
       ├─ Encrypt file
       ├─ Add to metadata
       └─ Secure delete original

3. [crypto/encrypt.go]
   ├─ Read plaintext
   ├─ Encrypt with AES-256-GCM
   └─ Return ciphertext

4. [storage/storage.go]
   ├─ Write encrypted data to objects/
   ├─ Write salt to .vaultix/salt
   ├─ Encrypt and write metadata
   └─ Secure delete originals
```

### Add File Flow

```
User runs: vaultix add file.txt

1. [cmd/add.go]
   ├─ Parse filename
   ├─ Verify file exists
   ├─ Prompt for password
   └─ Call vault.AddFile()

2. [vault/vault.go]
   ├─ Load salt
   ├─ Derive key
   ├─ Decrypt metadata (verify password)
   ├─ Generate random file ID
   ├─ Encrypt file
   ├─ Update metadata
   ├─ Save metadata
   └─ Secure delete original

3. [crypto/encrypt.go]
   ├─ Read file contents
   ├─ Encrypt with key
   └─ Return encrypted data

4. [storage/storage.go]
   ├─ Write to objects/<id>.enc
   ├─ Encrypt metadata
   ├─ Write updated metadata
   └─ Secure delete original file
```

### Extract File Flow

```
User runs: vaultix extract file.txt

1. [cmd/extract.go]
   ├─ Parse filename (fuzzy match)
   ├─ Prompt for password
   └─ Call vault.ExtractFile()

2. [vault/vault.go]
   ├─ Load salt
   ├─ Derive key
   ├─ Decrypt metadata (verify password)
   ├─ Find file by name (fuzzy match)
   ├─ Read encrypted file from objects/
   ├─ Decrypt file
   └─ Write to current directory

3. [crypto/decrypt.go]
   ├─ Verify key length
   ├─ Decrypt with AES-256-GCM
   └─ Return plaintext

4. [storage/storage.go]
   ├─ Read from objects/<id>.enc
   ├─ Decrypt
   └─ Write to working directory
```

## On-Disk Format

### Vault Structure

```
my_vault/
└── .vaultix/
    ├── salt          # 32 bytes random salt
    ├── meta          # Encrypted metadata JSON
    ├── config        # Vault configuration (future use)
    └── objects/
        ├── 3f9a2c1d.enc
        ├── 91bd77aa.enc
        └── ...
```

### Salt File Format

```
Fixed 32 bytes of random data
Used for Argon2id key derivation
Not encrypted (no sensitive data)
```

### Metadata File Format

```
Encrypted JSON containing:
{
  "version": 1,
  "files": [
    {
      "id": "3f9a2c1d",
      "name": "document.pdf",
      "size": 1048576,
      "modified": "2024-01-15T10:30:00Z"
    }
  ]
}

Encrypted with derived key (AES-256-GCM)
```

### Encrypted File Format

```
[nonce (12 bytes)][encrypted data][auth tag (16 bytes)]

- Nonce: Random 12 bytes (GCM standard)
- Encrypted data: AES-256-GCM output
- Auth tag: GCM authentication tag
```

## Error Handling

### Error Types

```go
// ErrPasswordIncorrect indicates wrong password
var ErrPasswordIncorrect = errors.New("password incorrect")

// ErrVaultNotFound indicates no vault exists
var ErrVaultNotFound = errors.New("vault not found")

// ErrFileNotFound indicates file doesn't exist in vault
var ErrFileNotFound = errors.New("file not found")

// ErrVaultCorrupted indicates vault data is corrupted
var ErrVaultCorrupted = errors.New("vault corrupted")
```

### Error Flow

```
[Low Level]
crypto.Decrypt() → vault.LoadMetadata() → cmd.List()
                                           └─> Display to user

Errors are wrapped with context:
fmt.Errorf("failed to decrypt metadata: %w", err)
```

## Security Considerations

### Memory Management

```go
// Clear sensitive data from memory
defer func() {
    for i := range password {
        password[i] = 0
    }
    for i := range key {
        key[i] = 0
    }
}()
```

### File Operations

- Original files are securely deleted (overwritten before removal)
- Temporary files are cleaned up even on error
- Encrypted data is written atomically (write to temp, then rename)

### Error Messages

```go
// Good: No information leakage
return errors.New("decryption failed")

// Bad: Leaks file information
return fmt.Errorf("failed to decrypt %s: wrong password", filename)
```

## Testing Strategy

### Unit Tests

Each layer has isolated tests:

```go
// crypto package tests
func TestEncryptDecrypt(t *testing.T) {
    key := make([]byte, 32)
    rand.Read(key)

    plaintext := []byte("secret data")

    ciphertext, err := Encrypt(plaintext, key)
    require.NoError(t, err)

    decrypted, err := Decrypt(ciphertext, key)
    require.NoError(t, err)

    assert.Equal(t, plaintext, decrypted)
}
```

### Integration Tests

Test full command flows:

```bash
#!/bin/bash
# test_workflow.sh

# Initialize vault
echo "password" | vaultix init
test $? -eq 0 || exit 1

# Add file
echo "secret" > test.txt
echo "password" | vaultix add test.txt
test $? -eq 0 || exit 1

# Verify encrypted
test ! -f test.txt || exit 1
test -d .vaultix || exit 1

# Extract
echo "password" | vaultix extract test.txt
test $? -eq 0 || exit 1
test -f test.txt || exit 1
```

## Performance Considerations

### Memory Usage

- Files are loaded entirely into memory
- Metadata is loaded entirely into memory
- For large files (>1GB), consider splitting

### CPU Usage

- Argon2id is intentionally CPU/memory intensive
- Adjust parameters in `crypto/key.go` if needed (carefully!)

### Disk I/O

- Sequential reads/writes preferred
- Atomic writes using temp files + rename

## Future Extensions

### Pluggable Storage

```go
type Storage interface {
    Read(id string) ([]byte, error)
    Write(id string, data []byte) error
    Delete(id string) error
    List() ([]string, error)
}

// Could support:
// - Local filesystem (current)
// - Cloud storage (S3, etc.)
// - Database backends
```

### Streaming Encryption

```go
// For large files
func EncryptStream(r io.Reader, w io.Writer, key []byte) error
```

### Multiple Vaults

```go
// Vault manager
type VaultManager struct {
    vaults map[string]*Vault
}

func (vm *VaultManager) Open(path string) (*Vault, error)
func (vm *VaultManager) List() []string
```

---

This architecture balances simplicity, security, and maintainability while leaving room for future enhancements.
