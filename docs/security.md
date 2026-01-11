# Security Model

Understanding Vaultix's security guarantees and limitations.

## Architecture Overview

Vaultix uses a **master key encryption model**:

1. **Master Key (256-bit random)**: Encrypts all vault data
2. **Password Protection**: Master key encrypted with Argon2id-derived key
3. **Recovery Key (256-bit random)**: Alternative way to decrypt master key
4. **No Plaintext Storage**: Master key never stored unencrypted on disk

This architecture provides:
- Dual unlock methods (password OR recovery key)
- Fast password changes (only re-encrypt master key, not all data)
- Recovery option if password forgotten
- Defense in depth (multiple encryption layers)

## Threat Model

### What Vaultix Protects Against ✓

- **Unauthorized file access**: Files are encrypted at rest with AES-256-GCM
- **Casual snooping**: Encrypted data is unreadable without password or recovery key
- **Filename leakage**: Original filenames are encrypted in metadata
- **Data tampering**: GCM provides authentication, detecting modifications
- **Password loss**: Recovery key provides backup access method

### What Vaultix Does NOT Protect Against ✗

- **Weak passwords**: A guessable password defeats all encryption
- **Lost recovery key**: If you lose BOTH password AND recovery key, data is permanently lost
- **Recovery key exposure**: Anyone with recovery key can unlock vault
- **Keyloggers/malware**: If your system is compromised, passwords can be captured
- **Memory attacks**: Decrypted data exists in memory during operations
- **Physical access**: Someone with your password or recovery key and physical access can decrypt
- **Legal compulsion**: Courts can order you to provide passwords/recovery keys
- **Side-channel attacks**: Advanced attacks on the cryptographic implementation

## Cryptographic Primitives

### Master Key

**Algorithm**: Cryptographically Secure Random Number Generator (CSPRNG)

**Size**: 256 bits (32 bytes)

**Purpose**: 
- Encrypts all vault data (files + metadata)
- Never stored in plaintext
- Encrypted twice: once with password-derived key, once with recovery key

**Generation**:
```go
masterKey := make([]byte, 32)
crypto/rand.Read(masterKey)
```

### Recovery Key

**Algorithm**: Cryptographically Secure Random Number Generator (CSPRNG)

**Size**: 256 bits (32 bytes)

**Purpose**:
- Alternative method to unlock vault
- Can decrypt the master key
- Displayed once during initialization

**Format**: Hexadecimal string with dashes for readability
```
5025f74e-c5d7a54a-7b99c87b-78cca1a0-61854d30-fb0d2783-a9df7067-b67ad345
```

### Key Derivation: Argon2id

**Algorithm**: Argon2id (winner of Password Hashing Competition)

**Parameters**:

- Memory: 64 MB
- Iterations: 1
- Parallelism: 4 threads
- Output: 32 bytes (256 bits)

**Why Argon2id?**

- Resistant to GPU/ASIC attacks
- Protects against side-channel attacks
- Recommended by OWASP
- Balanced between Argon2i (side-channel resistant) and Argon2d (GPU resistant)

### Encryption: AES-256-GCM

**Algorithm**: Advanced Encryption Standard with Galois/Counter Mode

**Key size**: 256 bits  
**Nonce size**: 96 bits (12 bytes)  
**Authentication tag**: 128 bits (16 bytes)

**Why AES-256-GCM?**

- Industry standard, extensively analyzed
- Provides both confidentiality and authenticity
- Authenticated encryption prevents tampering
- Hardware acceleration available on modern CPUs
- NIST approved

### Random Number Generation

**Source**: Go's `crypto/rand` package

Uses OS-provided cryptographically secure random number generators:

- Linux: `/dev/urandom`
- macOS: `SecRandomCopyBytes`
- Windows: `CryptGenRandom`

## Password Handling

### Password Flow

1. **Input**: Password entered via terminal (no echo)
2. **Derivation**: Argon2id generates 256-bit key
3. **Usage**: Key encrypts/decrypts data
4. **Cleanup**: Key zeroed in memory after use

### No Password Storage

Vaultix **never** stores:

- ✗ Your password
- ✗ Password hashes
- ✗ Password hints
- ✗ Recovery keys

Password correctness is verified by attempting to decrypt the metadata. Incorrect password = decryption failure.

### Password Requirements

Vaultix enforces:

- Minimum length: 1 character (but please use more!)
- Maximum length: No limit

**Recommended**:

- At least 16 characters
- Mix of uppercase, lowercase, numbers, symbols
- Use a password manager
- Don't reuse passwords
- Consider a passphrase (e.g., "correct horse battery staple")

## Data Flow

### Encryption Process

```
Plaintext File
    ↓
Read into memory
    ↓
Generate random nonce
    ↓
AES-256-GCM encryption with derived key
    ↓
Encrypted data + authentication tag
    ↓
Write to .vaultix/objects/
    ↓
Secure delete original file
```

### Decryption Process

```
Encrypted object
    ↓
Read from .vaultix/objects/
    ↓
Extract nonce from ciphertext
    ↓
AES-256-GCM decryption with derived key
    ↓
Verify authentication tag
    ↓
Plaintext data
    ↓
Write to output file
```

## Metadata Security

### What's in Metadata

- Original filenames
- File sizes
- Modification timestamps
- Object IDs (encrypted file references)

### Metadata Protection

- **Encrypted**: Metadata is encrypted with AES-256-GCM
- **Authenticated**: Tampering is detected
- **Single file**: All metadata in one encrypted blob

**Why encrypt metadata?**

Filenames can reveal sensitive information:

- `tax_return_2024.pdf` → Financial data
- `medical_records.txt` → Health information
- `job_applications.docx` → Employment status

## Secure Deletion

When files are deleted (after encryption or with `drop`/`remove`):

1. **Overwrite**: File contents overwritten with random data
2. **Delete**: File unlinked from filesystem

**Limitations**:

- SSDs may not physically overwrite due to wear leveling
- Copy-on-write filesystems (Btrfs, ZFS) may keep copies
- Filesystem journaling may preserve data
- Swap/hibernation files may contain plaintext

**Recommendation**: Use full-disk encryption (LUKS, FileVault, BitLocker) alongside Vaultix.

## Attack Scenarios

### Scenario 1: Stolen Laptop

**Attack**: Thief gets your laptop with encrypted vault

**Protection**:

- ✓ Files are encrypted with AES-256
- ✓ Decryption requires password
- ✓ Brute-force is slow (Argon2id)

**Outcome**: Data is safe if password is strong

---

### Scenario 2: Cloud Backup

**Attack**: Cloud provider compromised, vault backup leaked

**Protection**:

- ✓ Vault is fully encrypted
- ✓ Metadata is encrypted
- ✓ Object names don't reveal content

**Outcome**: Data is safe (same as stolen laptop)

---

### Scenario 3: Malware on System

**Attack**: Keylogger captures password while you use vault

**Protection**:

- ✗ Password is captured as you type
- ✗ Decrypted files can be read from memory
- ✗ Extracted files can be stolen

**Outcome**: Vaultix cannot protect against compromised systems

**Mitigation**: Keep your system clean, use antivirus, practice safe computing

---

### Scenario 4: Physical Access While Unlocked

**Attack**: Someone accesses your computer while you're away

**Protection**:

- ✗ Files may be extracted
- ✗ Password may be in command history
- ✗ Decrypted files may be on disk

**Outcome**: Lock your computer when away

**Mitigation**: Use screen lock, log out, close terminal after vault operations

---

### Scenario 5: Weak Password

**Attack**: Attacker brute-forces your password

**Protection**:

- ⚠️ Argon2id slows down attacks
- ✗ Weak passwords are still crackable

**Outcome**: Security depends on password strength

**Mitigation**: Use strong, unique passwords (16+ characters)

## Security Best Practices

### DO ✓

- Use strong, unique passwords
- Store vaults on encrypted drives
- Lock your computer when away
- Keep your OS and software updated
- Use a password manager
- Make encrypted backups
- Test your backups regularly

### DON'T ✗

- Reuse passwords
- Store passwords in plaintext
- Leave decrypted files lying around
- Use vaultix over unencrypted connections
- Trust public computers
- Forget to lock your screen

## Auditing

Vaultix has **not** undergone formal security auditing. The code is open source for community review, but no independent security firm has assessed it.

**Use at your own risk**.

## Reporting Security Issues

If you discover a security vulnerability:

1. **DO NOT** open a public issue
2. Email security concerns privately
3. Include details and reproduction steps
4. Allow time for fix before disclosure

## Conclusion

Vaultix provides strong cryptographic protection for files at rest. However, it's not a silver bullet:

- 🔐 Strong encryption protects data from unauthorized access
- 🔑 Security depends on password strength
- 💻 Cannot protect against compromised systems
- 🎯 Best used alongside other security measures

Think of Vaultix as **one layer** in your security strategy, not the only layer.
