# API Reference

Internal API documentation for vaultix developers.

> **Note:** Vaultix is a CLI tool without a public API. This reference documents internal packages for contributors.

## Package: `vault`

Core vault operations and management.

### Type: `Vault`

```go
type Vault struct {
    Path string // Vault directory path
}
```

Represents a vaultix encrypted directory.

### Constructor Functions

#### `New(path string) *Vault`

Creates a new Vault instance.

**Parameters:**

- `path` - Directory path for the vault

**Returns:**

- `*Vault` - New vault instance

**Example:**

```go
v := vault.New("/home/user/my_vault")
```

---

### Methods

#### `Init(password string) error`

Initializes a new vault and encrypts all existing files.

**Parameters:**

- `password` - User's password for encryption

**Returns:**

- `error` - Error if initialization fails

**Behavior:**

1. Creates `.vaultix/` directory structure
2. Generates random salt
3. Scans directory for all files
4. Encrypts each file
5. Stores encrypted files in `.vaultix/objects/`
6. Securely deletes originals

**Errors:**

- `ErrVaultExists` - Vault already exists
- `ErrInvalidPath` - Invalid directory path
- Standard I/O errors

**Example:**

```go
v := vault.New("./my_vault")
err := v.Init("my_secure_password")
```

---

#### `AddFile(filename, password string) error`

Adds a file to an existing vault.

**Parameters:**

- `filename` - File to encrypt and add
- `password` - Vault password

**Returns:**

- `error` - Error if operation fails

**Behavior:**

1. Verifies password by decrypting metadata
2. Reads file contents
3. Encrypts file
4. Generates random file ID
5. Updates metadata
6. Securely deletes original

**Errors:**

- `ErrVaultNotFound` - Vault doesn't exist
- `ErrPasswordIncorrect` - Wrong password
- `ErrFileNotFound` - File doesn't exist
- Standard I/O errors

**Example:**

```go
v := vault.New("./my_vault")
err := v.AddFile("secret.txt", "my_secure_password")
```

---

#### `ListFiles(password string) ([]FileInfo, error)`

Lists all files in the vault.

**Parameters:**

- `password` - Vault password

**Returns:**

- `[]FileInfo` - List of file metadata
- `error` - Error if operation fails

**Behavior:**

1. Loads salt
2. Derives key from password
3. Decrypts metadata
4. Returns file list

**Errors:**

- `ErrVaultNotFound` - Vault doesn't exist
- `ErrPasswordIncorrect` - Wrong password

**Example:**

```go
v := vault.New("./my_vault")
files, err := v.ListFiles("my_secure_password")
for _, f := range files {
    fmt.Printf("%s (%d bytes)\n", f.Name, f.Size)
}
```

---

#### `ExtractFile(filename, password string) error`

Extracts (decrypts) a file from the vault.

**Parameters:**

- `filename` - File to extract (supports fuzzy matching)
- `password` - Vault password

**Returns:**

- `error` - Error if operation fails

**Behavior:**

1. Verifies password
2. Finds file by name (fuzzy match)
3. Reads encrypted file
4. Decrypts file
5. Writes to current directory

**Errors:**

- `ErrVaultNotFound` - Vault doesn't exist
- `ErrPasswordIncorrect` - Wrong password
- `ErrFileNotFound` - No matching file found
- `ErrMultipleMatches` - Ambiguous fuzzy match

**Example:**

```go
v := vault.New("./my_vault")
// Exact match
err := v.ExtractFile("document.pdf", "password")

// Fuzzy match
err = v.ExtractFile("doc", "password")  // Matches "document.pdf"
```

---

#### `ExtractAll(password string) error`

Extracts all files from the vault.

**Parameters:**

- `password` - Vault password

**Returns:**

- `error` - Error if operation fails

**Behavior:**

1. Lists all files
2. Extracts each file to current directory

**Errors:**

- Same as `ExtractFile`

**Example:**

```go
v := vault.New("./my_vault")
err := v.ExtractAll("my_secure_password")
```

---

#### `DropFile(filename, password string) error`

Extracts a file and removes it from the vault.

**Parameters:**

- `filename` - File to drop
- `password` - Vault password

**Returns:**

- `error` - Error if operation fails

**Behavior:**

1. Extracts file
2. Removes from vault
3. Updates metadata

**Errors:**

- Same as `ExtractFile` and `RemoveFile`

**Example:**

```go
v := vault.New("./my_vault")
err := v.DropFile("temp_file.txt", "password")
// temp_file.txt is now in current directory and removed from vault
```

---

#### `RemoveFile(filename, password string) error`

Removes a file from the vault without extracting.

**Parameters:**

- `filename` - File to remove
- `password` - Vault password

**Returns:**

- `error` - Error if operation fails

**Behavior:**

1. Verifies password
2. Finds file by name
3. Securely deletes encrypted file
4. Updates metadata

**Errors:**

- `ErrVaultNotFound` - Vault doesn't exist
- `ErrPasswordIncorrect` - Wrong password
- `ErrFileNotFound` - No matching file found

**Example:**

```go
v := vault.New("./my_vault")
err := v.RemoveFile("old_file.txt", "password")
```

---

#### `Clear(password string) error`

Removes all files from the vault.

**Parameters:**

- `password` - Vault password

**Returns:**

- `error` - Error if operation fails

**Behavior:**

1. Verifies password
2. Removes all files from vault
3. Clears metadata

**Warning:** This is destructive and irreversible!

**Example:**

```go
v := vault.New("./my_vault")
err := v.Clear("my_secure_password")
```

---

## Package: `crypto`

Cryptographic operations.

### Functions

#### `DeriveKey(password string, salt []byte) []byte`

Derives an encryption key from a password using Argon2id.

**Parameters:**

- `password` - User's password
- `salt` - 32-byte random salt

**Returns:**

- `[]byte` - 32-byte AES-256 key

**Algorithm:**

```go
argon2.IDKey(
    []byte(password),
    salt,
    3,       // Time cost
    64*1024, // Memory cost (64 MB)
    4,       // Parallelism
    32,      // Key length
)
```

**Example:**

```go
salt := make([]byte, 32)
rand.Read(salt)
key := crypto.DeriveKey("password", salt)
```

---

#### `Encrypt(plaintext, key []byte) ([]byte, error)`

Encrypts data using AES-256-GCM.

**Parameters:**

- `plaintext` - Data to encrypt
- `key` - 32-byte encryption key

**Returns:**

- `[]byte` - Encrypted data (nonce + ciphertext + auth tag)
- `error` - Error if encryption fails

**Format:**

```
[12-byte nonce][encrypted data][16-byte auth tag]
```

**Example:**

```go
plaintext := []byte("secret data")
key := make([]byte, 32)
rand.Read(key)

ciphertext, err := crypto.Encrypt(plaintext, key)
```

---

#### `Decrypt(ciphertext, key []byte) ([]byte, error)`

Decrypts data using AES-256-GCM.

**Parameters:**

- `ciphertext` - Encrypted data (from `Encrypt`)
- `key` - 32-byte encryption key

**Returns:**

- `[]byte` - Decrypted plaintext
- `error` - Error if decryption fails or authentication fails

**Errors:**

- `ErrDecryptionFailed` - Wrong key OR corrupted/tampered data

**Example:**

```go
plaintext, err := crypto.Decrypt(ciphertext, key)
if err != nil {
    // Wrong password or corrupted data
}
```

---

#### `GenerateSalt() []byte`

Generates a random 32-byte salt.

**Returns:**

- `[]byte` - 32-byte random salt

**Example:**

```go
salt := crypto.GenerateSalt()
```

---

## Package: `storage`

File system operations.

### Functions

#### `ReadFile(path string) ([]byte, error)`

Reads a file from disk.

**Parameters:**

- `path` - File path

**Returns:**

- `[]byte` - File contents
- `error` - Error if read fails

---

#### `WriteFile(path string, data []byte) error`

Writes data to a file atomically.

**Parameters:**

- `path` - File path
- `data` - Data to write

**Returns:**

- `error` - Error if write fails

**Behavior:**

- Writes to temporary file
- Renames to target path (atomic)

---

#### `SecureDelete(path string) error`

Securely deletes a file (overwrite then remove).

**Parameters:**

- `path` - File to delete

**Returns:**

- `error` - Error if deletion fails

**Behavior:**

1. Opens file
2. Overwrites with random data
3. Syncs to disk
4. Removes file

---

#### `GetVaultPaths(base string) *Paths`

Returns vault directory paths.

**Parameters:**

- `base` - Vault base directory

**Returns:**

- `*Paths` - Struct with all vault paths

**Example:**

```go
paths := storage.GetVaultPaths("/home/user/vault")
// paths.Vaultix = "/home/user/vault/.vaultix"
// paths.Objects = "/home/user/vault/.vaultix/objects"
// paths.Salt    = "/home/user/vault/.vaultix/salt"
// paths.Meta    = "/home/user/vault/.vaultix/meta"
```

---

## Types

### `FileInfo`

```go
type FileInfo struct {
    ID       string    // Random identifier
    Name     string    // Original filename
    Size     int64     // File size in bytes
    Modified time.Time // Last modified time
}
```

Represents metadata for an encrypted file.

---

### `Metadata`

```go
type Metadata struct {
    Version int        // Metadata format version
    Files   []FileInfo // List of encrypted files
}
```

Vault metadata structure (stored encrypted).

---

### `Paths`

```go
type Paths struct {
    Base    string // Vault base directory
    Vaultix string // .vaultix directory
    Objects string // objects directory
    Salt    string // salt file
    Meta    string // metadata file
    Config  string // config file
}
```

Vault directory structure paths.

---

## Error Types

### Standard Errors

```go
var (
    // ErrVaultNotFound indicates no vault exists at path
    ErrVaultNotFound = errors.New("vault not found")

    // ErrVaultExists indicates vault already initialized
    ErrVaultExists = errors.New("vault already exists")

    // ErrPasswordIncorrect indicates wrong password provided
    ErrPasswordIncorrect = errors.New("password incorrect")

    // ErrFileNotFound indicates file doesn't exist in vault
    ErrFileNotFound = errors.New("file not found in vault")

    // ErrMultipleMatches indicates ambiguous fuzzy match
    ErrMultipleMatches = errors.New("multiple files match")

    // ErrVaultCorrupted indicates corrupted vault data
    ErrVaultCorrupted = errors.New("vault corrupted")

    // ErrDecryptionFailed indicates decryption/authentication failed
    ErrDecryptionFailed = errors.New("decryption failed")
)
```

---

## Constants

### Cryptographic Parameters

```go
const (
    // Argon2 key derivation parameters
    ArgonTime        = 3      // Time cost
    ArgonMemory      = 64*1024 // Memory cost (64 MB)
    ArgonParallelism = 4      // Parallelism
    ArgonKeyLength   = 32     // Key length (bytes)

    // Salt size
    SaltSize = 32 // bytes

    // AES key size
    AESKeySize = 32 // bytes (AES-256)
)
```

---

## Usage Examples

### Complete Workflow

```go
package main

import (
    "fmt"
    "vaultix/vault"
)

func main() {
    // Create vault instance
    v := vault.New("./my_vault")

    // Initialize vault
    if err := v.Init("my_password"); err != nil {
        panic(err)
    }
    fmt.Println("Vault initialized")

    // Add file
    if err := v.AddFile("secret.txt", "my_password"); err != nil {
        panic(err)
    }
    fmt.Println("File added")

    // List files
    files, err := v.ListFiles("my_password")
    if err != nil {
        panic(err)
    }

    for _, f := range files {
        fmt.Printf("- %s (%d bytes)\n", f.Name, f.Size)
    }

    // Extract file
    if err := v.ExtractFile("secret.txt", "my_password"); err != nil {
        panic(err)
    }
    fmt.Println("File extracted")
}
```

### Error Handling

```go
files, err := v.ListFiles(password)
if err != nil {
    switch err {
    case vault.ErrVaultNotFound:
        fmt.Println("No vault found. Run 'vaultix init' first.")
    case vault.ErrPasswordIncorrect:
        fmt.Println("Incorrect password.")
    default:
        fmt.Printf("Error: %v\n", err)
    }
    return
}
```

### Fuzzy Matching

```go
// Exact match
err := v.ExtractFile("document.pdf", password)

// Case-insensitive match
err = v.ExtractFile("DOCUMENT.pdf", password) // Matches "document.pdf"

// Substring match
err = v.ExtractFile("doc", password) // Matches "document.pdf"

// Multiple matches (error)
err = v.ExtractFile("doc", password)
if err == vault.ErrMultipleMatches {
    fmt.Println("Multiple files match 'doc'. Be more specific.")
}
```

---

## Testing Utilities

### Test Helpers

```go
// Create temporary test vault
func createTestVault(t *testing.T) (*vault.Vault, string) {
    tmpDir := t.TempDir()
    v := vault.New(tmpDir)

    // Create test file
    testFile := filepath.Join(tmpDir, "test.txt")
    os.WriteFile(testFile, []byte("test data"), 0644)

    // Initialize vault
    err := v.Init("test_password")
    require.NoError(t, err)

    return v, tmpDir
}

// Test vault operations
func TestVaultOperations(t *testing.T) {
    v, dir := createTestVault(t)
    defer os.RemoveAll(dir)

    // Test list
    files, err := v.ListFiles("test_password")
    require.NoError(t, err)
    assert.Len(t, files, 1)

    // Test extract
    err = v.ExtractFile("test.txt", "test_password")
    require.NoError(t, err)
}
```

---

## Performance Considerations

### Benchmarks

```go
func BenchmarkEncrypt(b *testing.B) {
    key := make([]byte, 32)
    rand.Read(key)
    data := make([]byte, 1024*1024) // 1 MB

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        crypto.Encrypt(data, key)
    }
}
```

### Typical Performance

- **Key derivation (Argon2id):** ~500ms
- **Encryption (AES-256-GCM):** ~200 MB/s
- **Decryption (AES-256-GCM):** ~200 MB/s

---

This API is subject to change. Check the source code for the most up-to-date information.
