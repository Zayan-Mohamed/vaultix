package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"

	"golang.org/x/crypto/argon2"
)

const (
	// Argon2id parameters - secure defaults
	argonTime    = 1
	argonMemory  = 64 * 1024 // 64 MB
	argonThreads = 4
	keyLength    = 32 // AES-256
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
