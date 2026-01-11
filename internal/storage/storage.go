package storage

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	vaultDirName        = ".vaultix"
	metaFileName        = "meta"
	saltFileName        = "salt"
	configFileName      = "config"
	objectsDirName      = "objects"
	masterKeyFileName   = "master.key"
	recoveryKeyFileName = "recovery.key"
)

var (
	ErrVaultExists    = errors.New("vault already exists at this location")
	ErrVaultNotFound  = errors.New("vault not found at this location")
	ErrFileNotInVault = errors.New("file not found in vault")
)

// VaultPaths holds all relevant paths for a vault
type VaultPaths struct {
	Root        string
	VaultDir    string
	Meta        string
	Salt        string
	Config      string
	Objects     string
	MasterKey   string
	RecoveryKey string
}

// FileMetadata stores information about an encrypted file
type FileMetadata struct {
	ID           string    `json:"id"`
	OriginalName string    `json:"original_name"`
	Size         int64     `json:"size"`
	ModTime      time.Time `json:"mod_time"`
	AddedAt      time.Time `json:"added_at"`
}

// VaultMetadata stores the list of all files in the vault
type VaultMetadata struct {
	Version int            `json:"version"`
	Files   []FileMetadata `json:"files"`
}

// GetVaultPaths returns the standard paths for a vault
func GetVaultPaths(rootPath string) VaultPaths {
	vaultDir := filepath.Join(rootPath, vaultDirName)
	return VaultPaths{
		Root:        rootPath,
		VaultDir:    vaultDir,
		Meta:        filepath.Join(vaultDir, metaFileName),
		Salt:        filepath.Join(vaultDir, saltFileName),
		Config:      filepath.Join(vaultDir, configFileName),
		Objects:     filepath.Join(vaultDir, objectsDirName),
		MasterKey:   filepath.Join(vaultDir, masterKeyFileName),
		RecoveryKey: filepath.Join(vaultDir, recoveryKeyFileName),
	}
}

// VaultExists checks if a vault exists at the given path
func VaultExists(rootPath string) bool {
	paths := GetVaultPaths(rootPath)
	info, err := os.Stat(paths.VaultDir)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// InitializeVault creates the vault directory structure
func InitializeVault(rootPath string) error {
	paths := GetVaultPaths(rootPath)

	if VaultExists(rootPath) {
		return ErrVaultExists
	}

	// Create vault directory
	if err := os.MkdirAll(paths.VaultDir, 0700); err != nil {
		return fmt.Errorf("failed to create vault directory: %w", err)
	}

	// Create objects directory
	if err := os.MkdirAll(paths.Objects, 0700); err != nil {
		return fmt.Errorf("failed to create objects directory: %w", err)
	}

	// Note: Metadata will be encrypted and written by the vault layer
	// We just create the directory structure here

	return nil
}

// ListDirectoryFiles returns a list of regular files in the directory (excluding hidden files and directories)
func ListDirectoryFiles(dirPath string) ([]string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		// Skip directories and hidden files (starting with .)
		if entry.IsDir() || entry.Name()[0] == '.' {
			continue
		}

		// Get full path
		fullPath := filepath.Join(dirPath, entry.Name())
		files = append(files, fullPath)
	}

	return files, nil
}

// WriteSalt stores the salt for the vault
func WriteSalt(rootPath string, salt []byte) error {
	paths := GetVaultPaths(rootPath)
	if err := os.WriteFile(paths.Salt, salt, 0600); err != nil {
		return fmt.Errorf("failed to write salt: %w", err)
	}
	return nil
}

// ReadSalt reads the salt from the vault
func ReadSalt(rootPath string) ([]byte, error) {
	paths := GetVaultPaths(rootPath)
	salt, err := os.ReadFile(paths.Salt)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrVaultNotFound
		}
		return nil, fmt.Errorf("failed to read salt: %w", err)
	}
	return salt, nil
}

// WriteMasterKey stores the encrypted master key
func WriteMasterKey(rootPath string, encryptedMasterKey []byte) error {
	paths := GetVaultPaths(rootPath)
	if err := os.WriteFile(paths.MasterKey, encryptedMasterKey, 0600); err != nil {
		return fmt.Errorf("failed to write master key: %w", err)
	}
	return nil
}

// ReadMasterKey reads the encrypted master key
func ReadMasterKey(rootPath string) ([]byte, error) {
	paths := GetVaultPaths(rootPath)
	encryptedMasterKey, err := os.ReadFile(paths.MasterKey)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrVaultNotFound
		}
		return nil, fmt.Errorf("failed to read master key: %w", err)
	}
	return encryptedMasterKey, nil
}

// WriteRecoveryKey stores the encrypted master key (encrypted with recovery key)
func WriteRecoveryKey(rootPath string, encryptedMasterKeyForRecovery []byte) error {
	paths := GetVaultPaths(rootPath)
	if err := os.WriteFile(paths.RecoveryKey, encryptedMasterKeyForRecovery, 0600); err != nil {
		return fmt.Errorf("failed to write recovery key file: %w", err)
	}
	return nil
}

// ReadRecoveryKey reads the encrypted master key (for recovery key unlock)
func ReadRecoveryKey(rootPath string) ([]byte, error) {
	paths := GetVaultPaths(rootPath)
	encryptedMasterKeyForRecovery, err := os.ReadFile(paths.RecoveryKey)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrVaultNotFound
		}
		return nil, fmt.Errorf("failed to read recovery key file: %w", err)
	}
	return encryptedMasterKeyForRecovery, nil
}

// ReadMetadata reads and returns the encrypted metadata
func ReadMetadata(rootPath string) ([]byte, error) {
	paths := GetVaultPaths(rootPath)
	data, err := os.ReadFile(paths.Meta)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrVaultNotFound
		}
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}
	return data, nil
}

// WriteMetadata writes encrypted metadata to disk
func WriteMetadata(rootPath string, data []byte) error {
	paths := GetVaultPaths(rootPath)
	if err := os.WriteFile(paths.Meta, data, 0600); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}
	return nil
}

// GenerateObjectID creates a unique ID for an object based on name and timestamp
func GenerateObjectID(originalName string) string {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s-%d", originalName, time.Now().UnixNano())))
	return hex.EncodeToString(hash[:8])
}

// GetObjectPath returns the path for an encrypted object
func GetObjectPath(rootPath, objectID string) string {
	paths := GetVaultPaths(rootPath)
	return filepath.Join(paths.Objects, objectID+".enc")
}

// WriteObject writes encrypted data to an object file
func WriteObject(rootPath, objectID string, data []byte) error {
	objectPath := GetObjectPath(rootPath, objectID)
	if err := os.WriteFile(objectPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write object: %w", err)
	}
	return nil
}

// ReadObject reads encrypted data from an object file
func ReadObject(rootPath, objectID string) ([]byte, error) {
	objectPath := GetObjectPath(rootPath, objectID)
	data, err := os.ReadFile(objectPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotInVault
		}
		return nil, fmt.Errorf("failed to read object: %w", err)
	}
	return data, nil
}

// DeleteObject removes an encrypted object file
func DeleteObject(rootPath, objectID string) error {
	objectPath := GetObjectPath(rootPath, objectID)
	if err := os.Remove(objectPath); err != nil {
		if os.IsNotExist(err) {
			return ErrFileNotInVault
		}
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}

// ReadPlaintextFile reads a file from disk (for adding to vault)
func ReadPlaintextFile(filePath string) ([]byte, os.FileInfo, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to stat file: %w", err)
	}

	if info.IsDir() {
		return nil, nil, errors.New("cannot add directory, only files are supported")
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, info, nil
}

// WritePlaintextFile writes decrypted data to disk
func WritePlaintextFile(filePath string, data []byte, modTime time.Time) error {
	// Ensure parent directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	// Restore modification time
	if err := os.Chtimes(filePath, modTime, modTime); err != nil {
		// Non-fatal - just log and continue
		fmt.Fprintf(os.Stderr, "warning: failed to restore modification time: %v\n", err)
	}

	return nil
}

// SecureDelete overwrites a file before deletion (best effort)
func SecureDelete(filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	// Open file for writing
	file, err := os.OpenFile(filePath, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()

	// Overwrite with random data
	size := info.Size()
	buf := make([]byte, 4096)
	for size > 0 {
		n := int64(len(buf))
		if n > size {
			n = size
		}
		// Read random data from crypto/rand
		if _, err := io.ReadFull(rand.Reader, buf[:n]); err != nil {
			// If random read fails, just use zeros
			for i := range buf[:n] {
				buf[i] = 0
			}
		}
		if _, err := file.Write(buf[:n]); err != nil {
			return err
		}
		size -= n
	}

	file.Close()

	// Now delete the file
	return os.Remove(filePath)
}
