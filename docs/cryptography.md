# Cryptography Details

Deep dive into the cryptographic primitives and implementation in vaultix.

## Overview

Vaultix uses a **master key encryption model** with **standard, well-audited cryptographic algorithms** from Go's crypto library. No custom cryptography is implemented.

### Master Key Architecture

Vaultix employs a layered encryption approach:

1. **Master Key (256-bit)**: Random key that encrypts all vault data
2. **Password-Encrypted Master Key**: Master key encrypted with Argon2id-derived key from user password
3. **Recovery-Key-Encrypted Master Key**: Master key encrypted with random recovery key (backup access)
4. **Data Encryption**: All files and metadata encrypted with the master key

This design provides:
- ✅ Dual unlock methods (password OR recovery key)
- ✅ No plaintext master key on disk
- ✅ Recovery option if password forgotten
- ✅ Fast re-keying possible (only re-encrypt master key, not all data)

## Algorithm Selection

### Encryption: AES-256-GCM

**Why AES-256?**

- ✅ Industry standard (NIST FIPS 197)
- ✅ Hardware acceleration on modern CPUs (AES-NI)
- ✅ Proven security (no practical attacks)
- ✅ 256-bit keys provide excellent security margin

**Why GCM mode?**

- ✅ Authenticated encryption (AEAD)
- ✅ Provides both confidentiality and integrity
- ✅ Detects tampering automatically
- ✅ Fast (hardware accelerated)
- ✅ Well-studied and standardized

**Alternatives considered and rejected:**

- ❌ AES-CBC: No authentication, vulnerable to padding oracle attacks
- ❌ AES-CTR: No authentication by itself
- ❌ ChaCha20-Poly1305: Good algorithm but less hardware support

### Key Derivation: Argon2id

**Why Argon2id?**

- ✅ Winner of Password Hashing Competition (2015)
- ✅ Resistant to GPU/ASIC attacks (memory-hard)
- ✅ Configurable time/memory cost
- ✅ Combined data-dependent and data-independent mixing
- ✅ Side-channel resistant

**Parameters used:**

```go
const (
    ArgonTime        = 3      // Number of iterations
    ArgonMemory      = 64*1024 // 64 MB memory
    ArgonParallelism = 4      // 4 threads
    ArgonKeyLength   = 32     // 32 bytes (256 bits)
)
```

**Why these parameters?**

- Time=3: ~500ms on modern CPU (usable but not instant)
- Memory=64MB: Expensive for attackers, reasonable for users
- Parallelism=4: Good for multi-core CPUs
- Key=32 bytes: AES-256 requirement

**Alternatives considered and rejected:**

- ❌ PBKDF2: Too fast, weak against GPU attacks
- ❌ bcrypt: Password length limit (72 bytes), not memory-hard
- ❌ scrypt: Good but Argon2 is newer and better

### Random Number Generation: crypto/rand

**Why crypto/rand?**

- ✅ Cryptographically secure (CSPRNG)
- ✅ Uses OS entropy source (/dev/urandom, CryptGenRandom, etc.)
- ✅ Non-deterministic
- ✅ Tested and audited

**Used for:**

- Salt generation
- Nonce generation
- File ID generation

**Never use:**

- ❌ math/rand: Deterministic, not cryptographically secure
- ❌ Time-based seeds: Predictable

## Detailed Cryptographic Operations

### Master Key Generation

```go
// Generate random 256-bit master key (performed once during vault init)
masterKey := make([]byte, 32) // 32 bytes = 256 bits
if _, err := rand.Read(masterKey); err != nil {
    return err
}
// Master key is NEVER stored in plaintext on disk
```

**Master Key Properties:**
- 256 bits (32 bytes) of cryptographically secure random data
- Generated once per vault during initialization
- Stored only in encrypted form (twice: password-encrypted and recovery-key-encrypted)
- All vault data (files + metadata) encrypted with this key

### Recovery Key Generation

```go
// Generate random 256-bit recovery key (performed once during vault init)
recoveryKey := make([]byte, 32) // 32 bytes = 256 bits
if _, err := rand.Read(recoveryKey); err != nil {
    return err
}
// Displayed to user ONCE as hex string with dashes
// Format: 5025f74e-c5d7a54a-7b99c87b-78cca1a0-...
```

**Recovery Key Properties:**
- 256 bits (32 bytes) of cryptographically secure random data
- Displayed once to user during initialization (save it!)
- Can decrypt the master key (alternative to password)
- Formatted as 8 groups of 8 hex characters for readability

### Key Derivation Process (Password to Key)

```go
// 1. Generate random 32-byte salt (only once during init)
salt := make([]byte, 32)
if _, err := rand.Read(salt); err != nil {
    return err
}

// 2. Derive key from password + salt using Argon2id
derivedKey := argon2.IDKey(
    []byte(password),  // User's password
    salt,              // Random salt
    1,                 // Time cost (iterations)
    64*1024,           // Memory cost (64 MB)
    4,                 // Parallelism (threads)
    32,                // Key length (256 bits)
)

// 3. Use derived key to encrypt/decrypt the master key
```

**Why salt?**

- Prevents rainbow table attacks
- Makes each vault's derived key unique
- Stored unencrypted (not secret data)

**Salt properties:**

- 32 bytes (256 bits)
- Randomly generated
- Unique per vault
- Stored in `.vaultix/salt`

### Master Key Encryption (with Password)

```go
// 1. Derive key from user's password
derivedKey := argon2.IDKey(password, salt, 1, 64*1024, 4, 32)

// 2. Encrypt master key with derived key
encryptedMasterKey := AES256GCM_Encrypt(masterKey, derivedKey)

// 3. Store encrypted master key
// File: .vaultix/master.key
// Contains: [nonce || encrypted_master_key || auth_tag]
```

### Master Key Encryption (with Recovery Key)

```go
// 1. Use recovery key directly (no derivation needed)
// 2. Encrypt master key with recovery key
encryptedMasterKeyForRecovery := AES256GCM_Encrypt(masterKey, recoveryKey)

// 3. Store encrypted master key
// File: .vaultix/recovery.key
// Contains: [nonce || encrypted_master_key || auth_tag]
```

### Unlocking the Vault (with Password)

```go
// 1. Read salt and encrypted master key
salt := ReadFile(".vaultix/salt")
encryptedMasterKey := ReadFile(".vaultix/master.key")

// 2. Derive key from password
derivedKey := argon2.IDKey(password, salt, 1, 64*1024, 4, 32)

// 3. Decrypt master key
masterKey := AES256GCM_Decrypt(encryptedMasterKey, derivedKey)
// If wrong password: decryption fails with authentication error

// 4. Use master key to decrypt vault data
```

### Unlocking the Vault (with Recovery Key)

```go
// 1. Read encrypted master key (recovery version)
encryptedMasterKeyForRecovery := ReadFile(".vaultix/recovery.key")

// 2. Decrypt master key using recovery key directly
masterKey := AES256GCM_Decrypt(encryptedMasterKeyForRecovery, recoveryKey)
// If wrong recovery key: decryption fails with authentication error

// 3. Use master key to decrypt vault data
```

### Encryption Process (File Data)

```go
// 1. Create AES cipher with MASTER KEY (not password!)
block, err := aes.NewCipher(masterKey) // masterKey is 32 bytes

// 2. Create GCM mode wrapper
gcm, err := cipher.NewGCM(block)

// 3. Generate random nonce (12 bytes for GCM)
nonce := make([]byte, gcm.NonceSize()) // 12 bytes
rand.Read(nonce)

// 4. Encrypt plaintext
ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
// Result: [nonce || encrypted_data || auth_tag]

// 5. Write ciphertext to disk
```

**Output format:**

```
[12-byte nonce][encrypted data][16-byte auth tag]
```

**Why include nonce in output?**

- Nonce must be unique for each encryption
- Needed for decryption
- Not secret (but must not be reused!)

### Decryption Process

```go
// 1. Read ciphertext from disk
ciphertext, err := os.ReadFile(path)

// 2. Create AES cipher with MASTER KEY
block, err := aes.NewCipher(masterKey)

// 3. Create GCM mode wrapper
gcm, err := cipher.NewGCM(block)

// 4. Split nonce from ciphertext
nonceSize := gcm.NonceSize() // 12
nonce := ciphertext[:nonceSize]
encrypted := ciphertext[nonceSize:]

// 5. Decrypt and verify
plaintext, err := gcm.Open(nil, nonce, encrypted, nil)
if err != nil {
    // Wrong password OR corrupted data OR tampered data
    return ErrDecryptionFailed
}

// 6. Return plaintext
```

**GCM authentication:**

- Verifies data hasn't been modified
- Detects bitflips, truncation, etc.
- Fails decryption if auth tag doesn't match

### Metadata Encryption

```go
// 1. Marshal metadata to JSON
jsonData, err := json.Marshal(metadata)

// 2. Encrypt JSON using MASTER KEY
encryptedMeta := gcm.Seal(nonce, nonce, jsonData, masterKey)

// 3. Write to .vaultix/meta
```

**Why encrypt metadata?**

- Protects filenames (could be sensitive)
- Protects file sizes
- Protects modification times
- Verifies password correctness

## Security Properties

### Confidentiality

**What's protected:**

- ✅ File contents
- ✅ Filenames
- ✅ File sizes
- ✅ Modification times
- ✅ Number of files (hidden in metadata)

**What's NOT protected:**

- ❌ Vault existence (`.vaultix/` is visible)
- ❌ Approximate vault size
- ❌ When vault was last modified
- ❌ Number of encrypted file objects

### Integrity

**GCM provides:**

- ✅ Authentication of ciphertext
- ✅ Detection of any modification
- ✅ Detection of truncation
- ✅ Prevention of ciphertext manipulation

**Example:**

```go
// If attacker modifies even one bit:
ciphertext[100] ^= 0x01

// Decryption will fail:
plaintext, err := gcm.Open(...)
// err != nil (authentication failed)
```

### Forward Secrecy

❌ **Not provided.**

If password is compromised:

- All past encrypted data can be decrypted
- All future encrypted data can be decrypted

**Mitigation:**

- Use strong passwords
- Change password if compromise suspected
- Consider separate vaults for time-sensitive data

## Attack Resistance

### Brute Force Attacks

**Scenario:** Attacker has vault and tries passwords.

**Protection:** Argon2id makes each attempt expensive

- Time: ~500ms per attempt
- Memory: 64 MB per attempt
- Parallelization limited by memory bandwidth

**Comparison:**

```
PBKDF2 (old standard):
- 1 billion attempts/second on GPU
- Weak 8-char password cracked: instantly

Argon2id (vaultix):
- ~2000 attempts/second on high-end GPU
- Weak 8-char password cracked: hours to days
- Strong 16-char password: infeasible
```

**Recommendation:** Use 16+ character passwords

### Dictionary Attacks

**Scenario:** Attacker tries common passwords.

**Protection:** Argon2id slows down attempts

**Mitigation:**

- Don't use dictionary words
- Use password manager generated passwords
- Use passphrases: "correct horse battery staple"

### Rainbow Table Attacks

**Scenario:** Attacker pre-computes password hashes.

**Protection:** Random salt makes pre-computation useless

- Each vault has unique salt
- Must compute from scratch for each vault

### Side-Channel Attacks

**Timing attacks:**

- ❌ Partially vulnerable: String comparison for filenames
- ✅ GCM is constant-time authenticated decryption
- ✅ Argon2id is side-channel resistant

**Power analysis:**

- ❓ Not applicable (requires physical access to hardware)

**Cache-timing:**

- ✅ AES-NI resistant
- ✅ Argon2id resistant (memory-hard)

### Known-Plaintext Attacks

**Scenario:** Attacker knows some plaintext/ciphertext pairs.

**Protection:** AES-256-GCM is resistant

- Knowing plaintext doesn't help recover key
- Each file has unique nonce (no pattern reuse)

### Chosen-Ciphertext Attacks

**Scenario:** Attacker can get files decrypted.

**Protection:** GCM authentication prevents manipulation

- Can't create valid ciphertexts without key
- Tampering detected before plaintext revealed

### Quantum Computing

**Current status:**

- AES-256: ~128-bit post-quantum security (Grover's algorithm)
- Still very strong (2^128 operations infeasible)

**Future-proofing:**

- May need post-quantum key derivation eventually
- AES-256 should remain secure for decades

## Implementation Details

### Memory Handling

```go
// Clear sensitive data after use
defer func() {
    // Overwrite password
    for i := range password {
        password[i] = 0
    }

    // Overwrite key
    for i := range key {
        key[i] = 0
    }
}()
```

**Why?**

- Prevents recovery from memory dumps
- Reduces window for memory-reading malware
- Defense in depth

**Limitations:**

- Go garbage collector may have copies
- Swapped memory may persist
- Not perfect but better than nothing

### Nonce Management

```go
// Generate random nonce for each encryption
nonce := make([]byte, 12)
rand.Read(nonce)
```

**Critical: Never reuse nonces!**

With GCM, reusing nonce with same key:

- ❌ Breaks confidentiality
- ❌ Breaks authentication
- ❌ Allows key recovery

**How we avoid this:**

- Random nonce for each file
- 12 bytes = 96 bits
- Collision probability: negligible

### Error Handling

```go
// Don't leak information in errors
if err := gcm.Open(...); err != nil {
    // Good: Generic error
    return errors.New("decryption failed")

    // Bad: Leaks information
    // return fmt.Errorf("wrong password for %s", filename)
}
```

**Why?**

- Prevents information leakage
- Doesn't reveal vault contents
- Constant-time-ish failure

## Comparison with Alternatives

### vs. GPG

| Feature            | Vaultix     | GPG             |
| ------------------ | ----------- | --------------- |
| **Encryption**     | AES-256-GCM | AES, 3DES, etc. |
| **Key derivation** | Argon2id    | S2K (weaker)    |
| **Public key**     | No          | Yes             |
| **Ease of use**    | High        | Low             |
| **File-level**     | Yes         | Yes             |

**Use vaultix:** Simple password-based encryption
**Use GPG:** Public key crypto, digital signatures

### vs. VeraCrypt

| Feature                   | Vaultix     | VeraCrypt             |
| ------------------------- | ----------- | --------------------- |
| **Encryption**            | AES-256-GCM | AES, Serpent, Twofish |
| **Key derivation**        | Argon2id    | PBKDF2/SHA-512        |
| **Container**             | Directory   | File container        |
| **Mounting**              | No          | Yes                   |
| **Plausible deniability** | No          | Yes                   |

**Use vaultix:** File-level encryption, CLI workflow
**Use VeraCrypt:** Container-based, GUI, hidden volumes

### vs. age

| Feature            | Vaultix     | age               |
| ------------------ | ----------- | ----------------- |
| **Encryption**     | AES-256-GCM | ChaCha20-Poly1305 |
| **Key derivation** | Argon2id    | scrypt            |
| **Public key**     | No          | Yes               |
| **Directory**      | Yes         | No (file-by-file) |

**Use vaultix:** Directory-based workflow
**Use age:** Single-file encryption, public key crypto

## Best Practices

### Choosing Passwords

**Good password characteristics:**

- ✅ Length: 16+ characters
- ✅ Entropy: Random characters or diceware
- ✅ Uniqueness: Don't reuse
- ✅ Storage: Password manager

**Examples:**

```
Weak (DON'T USE):
- password123
- qwerty
- letmein
- Password1!

Strong (GOOD):
- Generated: 7wT$9#xK2mP@5nL&8qR
- Passphrase: correct-horse-battery-staple-purple-elephant
- Diceware: ware-klaxon-pride-ache-agony-brunt
```

### Key Rotation

**When to rotate:**

- Password potentially compromised
- Regular schedule (e.g., yearly)
- After employee leaves (for work vaults)

**How to rotate:**

```bash
# 1. Extract all files
vaultix extract

# 2. Remove vault
rm -rf .vaultix

# 3. Re-initialize with new password
vaultix init
# (Enter new password)
```

### Multiple Vaults

**Consider separate vaults for:**

- Different sensitivity levels
- Different time periods
- Different teams/projects
- Different backup schedules

**Example structure:**

```
~/vaults/
├── personal/      # Password A
├── work/          # Password B
├── financial/     # Password C (strongest)
└── archive/       # Password A (same as personal)
```

## Threat Model

### What Vaultix Protects Against

✅ **Protects against:**

- Stolen laptop (with strong password)
- Unauthorized physical access
- Cloud storage providers reading data
- Network sniffing (of vault files)
- Malware reading encrypted files on disk

❌ **Does NOT protect against:**

- Keyloggers (captures password)
- Screen recorders (sees decrypted files)
- Memory dumps while files decrypted
- Malware with root/admin access
- $5 wrench attack (coercion)
- Quantum computers (in far future)

### Trust Assumptions

**You must trust:**

- Your OS and kernel
- Go's crypto library
- The computer's CPU (no hardware backdoors)
- Your password manager (if used)
- Yourself (not to forget password)

**You don't need to trust:**

- Cloud storage providers
- Network infrastructure
- Backup services
- Other users on same system

## Cryptographic Verification

### How to Verify Implementation

```bash
# 1. Check Go crypto usage
grep -r "crypto/" *.go

# Should see:
# - crypto/aes
# - crypto/cipher
# - crypto/rand
# - golang.org/x/crypto/argon2

# 2. Check for weak crypto
grep -r "math/rand" *.go  # Should return nothing

# 3. Run tests
go test -v ./crypto
```

### Audit Checklist

- [ ] AES-256 used (not AES-128)
- [ ] GCM mode used (authenticated encryption)
- [ ] Argon2id for key derivation (not PBKDF2/bcrypt)
- [ ] crypto/rand for all randomness
- [ ] Unique nonce per encryption
- [ ] Salt stored and used correctly
- [ ] No custom crypto implementations
- [ ] Sensitive data cleared from memory
- [ ] No password logging

---

**Bottom line:** Vaultix uses battle-tested, industry-standard cryptography. No novel algorithms, no custom implementations, just proven security.
