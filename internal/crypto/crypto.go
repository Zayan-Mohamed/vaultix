package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
)

const (
	// Argon2id parameters - secure defaults
	argonTime    = 1
	argonMemory  = 64 * 1024 // 64 MB
	argonThreads = 4
	keyLength    = 32 // AES-256 (256-bit)
	saltLength   = 32
	nonceLength  = 12 // GCM standard nonce size
)

var (
	ErrInvalidPassword   = errors.New("decryption failed: incorrect password")
	ErrCorruptedData     = errors.New("decryption failed: data is corrupted")
	ErrInvalidSaltLength = errors.New("invalid salt length")
)

// GenerateSalt creates a cryptographically secure random salt
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, saltLength)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	return salt, nil
}

// DeriveKey derives an encryption key from a password and salt using Argon2id
// This function is deterministic - same password + salt always produces same key
func DeriveKey(password string, salt []byte) ([]byte, error) {
	if len(salt) != saltLength {
		return nil, ErrInvalidSaltLength
	}

	// Argon2id is the recommended variant (hybrid of Argon2i and Argon2d)
	key := argon2.IDKey(
		[]byte(password),
		salt,
		argonTime,
		argonMemory,
		argonThreads,
		keyLength,
	)

	return key, nil
}

// Encrypt encrypts plaintext using AES-256-GCM with the provided key
// Returns: nonce + ciphertext + tag (all in one slice)
func Encrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Generate random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Encrypt and authenticate
	// GCM appends the authentication tag to the ciphertext
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return ciphertext, nil
}

// Decrypt decrypts ciphertext using AES-256-GCM with the provided key
// Expects: nonce + ciphertext + tag (as returned by Encrypt)
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
	if len(ciphertext) < nonceSize {
		return nil, ErrCorruptedData
	}

	// Extract nonce and actual ciphertext
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt and verify authentication tag
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		// GCM authentication failure could mean wrong password or corrupted data
		// We treat both as wrong password for security (don't leak info)
		return nil, ErrInvalidPassword
	}

	return plaintext, nil
}

// GenerateMasterKey generates a random 256-bit master key for the vault
// This key is used to encrypt all vault data
func GenerateMasterKey() ([]byte, error) {
	masterKey := make([]byte, keyLength) // 32 bytes = 256 bits
	if _, err := io.ReadFull(rand.Reader, masterKey); err != nil {
		return nil, fmt.Errorf("failed to generate master key: %w", err)
	}
	return masterKey, nil
}

// GenerateRecoveryKey generates a random 256-bit recovery key
// This key can be used to decrypt the master key as an alternative to password
func GenerateRecoveryKey() ([]byte, error) {
	recoveryKey := make([]byte, keyLength) // 32 bytes = 256 bits
	if _, err := io.ReadFull(rand.Reader, recoveryKey); err != nil {
		return nil, fmt.Errorf("failed to generate recovery key: %w", err)
	}
	return recoveryKey, nil
}

// EncryptMasterKey encrypts the master key using a password-derived key
// Returns the encrypted master key
func EncryptMasterKey(masterKey, password []byte, salt []byte) ([]byte, error) {
	// Derive key from password
	derivedKey, err := DeriveKey(string(password), salt)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key: %w", err)
	}

	// Encrypt master key
	encryptedMasterKey, err := Encrypt(masterKey, derivedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt master key: %w", err)
	}

	return encryptedMasterKey, nil
}

// DecryptMasterKey decrypts the master key using a password-derived key
// Returns the decrypted master key
func DecryptMasterKey(encryptedMasterKey []byte, password string, salt []byte) ([]byte, error) {
	// Derive key from password
	derivedKey, err := DeriveKey(password, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key: %w", err)
	}

	// Decrypt master key
	masterKey, err := Decrypt(encryptedMasterKey, derivedKey)
	if err != nil {
		return nil, err // Return original error (ErrInvalidPassword or ErrCorruptedData)
	}

	return masterKey, nil
}

// EncryptMasterKeyWithRecoveryKey encrypts the master key using the recovery key
// Returns the encrypted master key
func EncryptMasterKeyWithRecoveryKey(masterKey, recoveryKey []byte) ([]byte, error) {
	encryptedMasterKey, err := Encrypt(masterKey, recoveryKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt master key with recovery key: %w", err)
	}
	return encryptedMasterKey, nil
}

// DecryptMasterKeyWithRecoveryKey decrypts the master key using the recovery key
// Returns the decrypted master key
func DecryptMasterKeyWithRecoveryKey(encryptedMasterKey, recoveryKey []byte) ([]byte, error) {
	masterKey, err := Decrypt(encryptedMasterKey, recoveryKey)
	if err != nil {
		return nil, err // Return original error (ErrInvalidPassword or ErrCorruptedData)
	}
	return masterKey, nil
}

// EncodeRecoveryKeyHex encodes a recovery key as a hexadecimal string
// This format is easier to copy/paste and store
func EncodeRecoveryKeyHex(recoveryKey []byte) string {
	return hex.EncodeToString(recoveryKey)
}

// DecodeRecoveryKeyHex decodes a recovery key from a hexadecimal string
// Returns an error if the format is invalid
func DecodeRecoveryKeyHex(hexKey string) ([]byte, error) {
	recoveryKey, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, fmt.Errorf("invalid recovery key format: %w", err)
	}
	if len(recoveryKey) != keyLength {
		return nil, fmt.Errorf("invalid recovery key length: expected %d bytes, got %d", keyLength, len(recoveryKey))
	}
	return recoveryKey, nil
}

// FormatRecoveryKeyForDisplay formats a recovery key for human-readable display
// Splits the hex string into groups for easier reading
func FormatRecoveryKeyForDisplay(recoveryKey []byte) string {
	hexKey := EncodeRecoveryKeyHex(recoveryKey)
	// Split into groups of 8 characters for readability
	// Example: 12345678-90abcdef-12345678-90abcdef-...
	var formatted string
	for i := 0; i < len(hexKey); i += 8 {
		end := i + 8
		if end > len(hexKey) {
			end = len(hexKey)
		}
		if i > 0 {
			formatted += "-"
		}
		formatted += hexKey[i:end]
	}
	return formatted
}
