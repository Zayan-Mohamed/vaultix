package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Zayan-Mohamed/vaultix/internal/crypto"
	"github.com/Zayan-Mohamed/vaultix/internal/storage"
)

var (
	ErrFileAlreadyExists = errors.New("file already exists in vault")
	ErrFileNotFound      = errors.New("file not found in vault")
)

// Vault represents a secure vault instance
type Vault struct {
	rootPath string
}

// New creates a new vault instance at the given path
func New(rootPath string) *Vault {
	return &Vault{rootPath: rootPath}
}

// Initialize creates a new vault with the given password and encrypts all files in the directory
// Returns the recovery key that should be saved by the user
func (v *Vault) Initialize(password string) ([]byte, error) {
	// Get list of files to encrypt before creating vault structure
	filesToEncrypt, err := storage.ListDirectoryFiles(v.rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	// Initialize vault structure
	if err := storage.InitializeVault(v.rootPath); err != nil {
		return nil, err
	}

	// Generate master key (random 256-bit key)
	masterKey, err := crypto.GenerateMasterKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate master key: %w", err)
	}

	// Generate recovery key (random 256-bit key)
	recoveryKey, err := crypto.GenerateRecoveryKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate recovery key: %w", err)
	}

	// Generate salt for password-based key derivation
	salt, err := crypto.GenerateSalt()
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	if err := storage.WriteSalt(v.rootPath, salt); err != nil {
		return nil, err
	}

	// Encrypt master key with password-derived key
	encryptedMasterKey, err := crypto.EncryptMasterKey(masterKey, []byte(password), salt)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt master key with password: %w", err)
	}

	if err := storage.WriteMasterKey(v.rootPath, encryptedMasterKey); err != nil {
		return nil, err
	}

	// Encrypt master key with recovery key
	encryptedMasterKeyForRecovery, err := crypto.EncryptMasterKeyWithRecoveryKey(masterKey, recoveryKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt master key with recovery key: %w", err)
	}

	if err := storage.WriteRecoveryKey(v.rootPath, encryptedMasterKeyForRecovery); err != nil {
		return nil, err
	}

	// Encrypt the initial empty metadata with master key
	meta := &storage.VaultMetadata{
		Version: 1,
		Files:   []storage.FileMetadata{},
	}
	if err := v.writeMetadata(masterKey, meta); err != nil {
		return nil, fmt.Errorf("failed to write initial metadata: %w", err)
	}

	// Encrypt all files in the directory
	if len(filesToEncrypt) > 0 {
		for _, filePath := range filesToEncrypt {
			if err := v.addFileInternal(filePath, masterKey); err != nil {
				return nil, fmt.Errorf("failed to encrypt %s: %w", filePath, err)
			}
		}

		// Delete original files after successful encryption
		for _, filePath := range filesToEncrypt {
			if err := storage.SecureDelete(filePath); err != nil {
				// Log warning but don't fail - file is already encrypted
				fmt.Fprintf(os.Stderr, "warning: failed to delete %s: %v\n", filePath, err)
			}
		}
	}

	return recoveryKey, nil
}

// unlockWithPassword decrypts the master key using the password
func (v *Vault) unlockWithPassword(password string) ([]byte, error) {
	// Read salt
	salt, err := storage.ReadSalt(v.rootPath)
	if err != nil {
		return nil, err
	}

	// Read encrypted master key
	encryptedMasterKey, err := storage.ReadMasterKey(v.rootPath)
	if err != nil {
		return nil, err
	}

	// Decrypt master key
	masterKey, err := crypto.DecryptMasterKey(encryptedMasterKey, password, salt)
	if err != nil {
		return nil, err
	}

	return masterKey, nil
}

// UnlockWithRecoveryKey decrypts the master key using the recovery key
func (v *Vault) UnlockWithRecoveryKey(recoveryKey []byte) ([]byte, error) {
	// Read encrypted master key (for recovery)
	encryptedMasterKeyForRecovery, err := storage.ReadRecoveryKey(v.rootPath)
	if err != nil {
		return nil, err
	}

	// Decrypt master key
	masterKey, err := crypto.DecryptMasterKeyWithRecoveryKey(encryptedMasterKeyForRecovery, recoveryKey)
	if err != nil {
		return nil, err
	}

	return masterKey, nil
}

// ListFilesWithMasterKey lists files using the master key directly (for recovery)
func (v *Vault) ListFilesWithMasterKey(masterKey []byte) ([]storage.FileMetadata, error) {
	meta, err := v.readMetadata(masterKey)
	if err != nil {
		return nil, err
	}
	return meta.Files, nil
}

// ExtractFileWithMasterKey extracts a file using the master key directly (for recovery)
func (v *Vault) ExtractFileWithMasterKey(masterKey []byte, fileName, destPath string) (string, error) {
	// Read and decrypt metadata
	meta, err := v.readMetadata(masterKey)
	if err != nil {
		return "", err
	}

	// Find the file with fuzzy matching
	fileMeta := findFileByName(meta.Files, fileName)
	if fileMeta == nil {
		return "", ErrFileNotFound
	}

	// Read encrypted object
	encryptedData, err := storage.ReadObject(v.rootPath, fileMeta.ID)
	if err != nil {
		return "", err
	}

	// Decrypt with master key
	plaintext, err := crypto.Decrypt(encryptedData, masterKey)
	if err != nil {
		return "", err
	}

	// Determine output path
	outputPath := destPath
	if destPath == "" || destPath == "." {
		outputPath = fileMeta.OriginalName
	}

	// Write decrypted file
	if err := storage.WritePlaintextFile(outputPath, plaintext, fileMeta.ModTime); err != nil {
		return "", err
	}

	return fileMeta.OriginalName, nil
}

// ExtractAllFilesWithMasterKey extracts all files using the master key directly (for recovery)
func (v *Vault) ExtractAllFilesWithMasterKey(masterKey []byte, destDir string) (int, error) {
	// Read and decrypt metadata
	meta, err := v.readMetadata(masterKey)
	if err != nil {
		return 0, err
	}

	if len(meta.Files) == 0 {
		return 0, nil
	}

	// Extract each file
	count := 0
	for _, fileMeta := range meta.Files {
		// Read encrypted object
		encryptedData, err := storage.ReadObject(v.rootPath, fileMeta.ID)
		if err != nil {
			return count, fmt.Errorf("failed to read %s: %w", fileMeta.OriginalName, err)
		}

		// Decrypt with master key
		plaintext, err := crypto.Decrypt(encryptedData, masterKey)
		if err != nil {
			return count, fmt.Errorf("failed to decrypt %s: %w", fileMeta.OriginalName, err)
		}

		// Determine output path
		outputPath := fileMeta.OriginalName
		if destDir != "" && destDir != "." {
			outputPath = filepath.Join(destDir, fileMeta.OriginalName)
		}

		// Write decrypted file
		if err := storage.WritePlaintextFile(outputPath, plaintext, fileMeta.ModTime); err != nil {
			return count, fmt.Errorf("failed to write %s: %w", fileMeta.OriginalName, err)
		}

		count++
	}

	return count, nil
}

// AddFile encrypts and adds a file to the vault
func (v *Vault) AddFile(password, filePath string) error {
	// Unlock vault with password
	masterKey, err := v.unlockWithPassword(password)
	if err != nil {
		return err
	}

	// Add the file using internal helper
	if err := v.addFileInternal(filePath, masterKey); err != nil {
		return err
	}

	// Securely delete the original file
	if err := storage.SecureDelete(filePath); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to securely delete original file: %v\n", err)
		// Continue anyway - the file was encrypted successfully
	}

	return nil
}

// addFileInternal is the internal implementation for adding files
func (v *Vault) addFileInternal(filePath string, masterKey []byte) error {
	// Read the file to be added
	data, info, err := storage.ReadPlaintextFile(filePath)
	if err != nil {
		return err
	}

	// Read and decrypt existing metadata
	meta, err := v.readMetadata(masterKey)
	if err != nil {
		return fmt.Errorf("failed to read metadata: %w", err)
	}

	// Check if file already exists
	fileName := filepath.Base(filePath)
	for _, f := range meta.Files {
		if f.OriginalName == fileName {
			return ErrFileAlreadyExists
		}
	}

	// Generate unique object ID
	objectID := storage.GenerateObjectID(fileName)

	// Encrypt file data with master key
	encryptedData, err := crypto.Encrypt(data, masterKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt file: %w", err)
	}

	// Write encrypted object
	if err := storage.WriteObject(v.rootPath, objectID, encryptedData); err != nil {
		return err
	}

	// Add to metadata
	fileMeta := storage.FileMetadata{
		ID:           objectID,
		OriginalName: fileName,
		Size:         info.Size(),
		ModTime:      info.ModTime(),
		AddedAt:      time.Now(),
	}
	meta.Files = append(meta.Files, fileMeta)

	// Encrypt and write updated metadata
	if err := v.writeMetadata(masterKey, meta); err != nil {
		// Try to clean up the object file
		storage.DeleteObject(v.rootPath, objectID)
		return fmt.Errorf("failed to update metadata: %w", err)
	}

	return nil
}

// ListFiles returns the list of files in the vault
func (v *Vault) ListFiles(password string) ([]storage.FileMetadata, error) {
	// Unlock vault with password
	masterKey, err := v.unlockWithPassword(password)
	if err != nil {
		return nil, err
	}

	// Read and decrypt metadata
	meta, err := v.readMetadata(masterKey)
	if err != nil {
		return nil, err
	}

	return meta.Files, nil
}

// ExtractFile decrypts and extracts a file from the vault
// Returns the actual filename that was matched (for fuzzy matching)
func (v *Vault) ExtractFile(password, fileName, destPath string) (string, error) {
	// Unlock vault with password
	masterKey, err := v.unlockWithPassword(password)
	if err != nil {
		return "", err
	}

	// Read and decrypt metadata
	meta, err := v.readMetadata(masterKey)
	if err != nil {
		return "", err
	}

	// Find the file with fuzzy matching
	fileMeta := findFileByName(meta.Files, fileName)

	if fileMeta == nil {
		return "", ErrFileNotFound
	}

	// Read encrypted object
	encryptedData, err := storage.ReadObject(v.rootPath, fileMeta.ID)
	if err != nil {
		return "", err
	}

	// Decrypt with master key
	plaintext, err := crypto.Decrypt(encryptedData, masterKey)
	if err != nil {
		return "", err
	}

	// Determine output path - use original filename, not the query
	outputPath := destPath
	if destPath == "" || destPath == "." {
		outputPath = fileMeta.OriginalName
	}

	// Write decrypted file
	if err := storage.WritePlaintextFile(outputPath, plaintext, fileMeta.ModTime); err != nil {
		return "", err
	}

	return fileMeta.OriginalName, nil
}

// ExtractAllFiles decrypts and extracts all files from the vault
func (v *Vault) ExtractAllFiles(password, destDir string) (int, error) {
	// Unlock vault with password
	masterKey, err := v.unlockWithPassword(password)
	if err != nil {
		return 0, err
	}

	// Read and decrypt metadata
	meta, err := v.readMetadata(masterKey)
	if err != nil {
		return 0, err
	}

	if len(meta.Files) == 0 {
		return 0, nil
	}

	// Extract each file
	count := 0
	for _, fileMeta := range meta.Files {
		// Read encrypted object
		encryptedData, err := storage.ReadObject(v.rootPath, fileMeta.ID)
		if err != nil {
			return count, fmt.Errorf("failed to read %s: %w", fileMeta.OriginalName, err)
		}

		// Decrypt with master key
		plaintext, err := crypto.Decrypt(encryptedData, masterKey)
		if err != nil {
			return count, fmt.Errorf("failed to decrypt %s: %w", fileMeta.OriginalName, err)
		}

		// Determine output path
		outputPath := fileMeta.OriginalName
		if destDir != "" && destDir != "." {
			outputPath = filepath.Join(destDir, fileMeta.OriginalName)
		}

		// Write decrypted file
		if err := storage.WritePlaintextFile(outputPath, plaintext, fileMeta.ModTime); err != nil {
			return count, fmt.Errorf("failed to write %s: %w", fileMeta.OriginalName, err)
		}

		count++
	}

	return count, nil
}

// DropFile extracts a file and then removes it from the vault
func (v *Vault) DropFile(password, fileName, destPath string) (string, error) {
	// First extract the file
	actualFileName, err := v.ExtractFile(password, fileName, destPath)
	if err != nil {
		return "", err
	}

	// Then remove it from the vault
	if err := v.RemoveFile(password, actualFileName); err != nil {
		return "", fmt.Errorf("extracted but failed to remove from vault: %w", err)
	}

	return actualFileName, nil
}

// DropAllFiles extracts all files and then removes them from the vault
func (v *Vault) DropAllFiles(password, destDir string) (int, error) {
	// First extract all files
	count, err := v.ExtractAllFiles(password, destDir)
	if err != nil {
		return 0, err
	}

	// Then clear the vault
	if err := v.ClearVault(password); err != nil {
		return count, fmt.Errorf("extracted %d files but failed to clear vault: %w", count, err)
	}

	return count, nil
}

// ClearVault removes all files from the vault without extracting them
func (v *Vault) ClearVault(password string) error {
	// Unlock vault with password
	masterKey, err := v.unlockWithPassword(password)
	if err != nil {
		return err
	}

	// Read and decrypt metadata
	meta, err := v.readMetadata(masterKey)
	if err != nil {
		return err
	}

	// Delete all objects
	for _, f := range meta.Files {
		if err := storage.DeleteObject(v.rootPath, f.ID); err != nil {
			return fmt.Errorf("failed to delete %s: %w", f.OriginalName, err)
		}
	}

	// Clear metadata
	meta.Files = []storage.FileMetadata{}
	if err := v.writeMetadata(masterKey, meta); err != nil {
		return err
	}

	return nil
}

// RemoveFile removes a file from the vault
func (v *Vault) RemoveFile(password, fileName string) error {
	// Unlock vault with password
	masterKey, err := v.unlockWithPassword(password)
	if err != nil {
		return err
	}

	// Read and decrypt metadata
	meta, err := v.readMetadata(masterKey)
	if err != nil {
		return err
	}

	// Find and remove the file with fuzzy matching
	var objectID string
	found := false
	newFiles := make([]storage.FileMetadata, 0, len(meta.Files))
	for _, f := range meta.Files {
		if matchesFileName(f.OriginalName, fileName) {
			objectID = f.ID
			found = true
		} else {
			newFiles = append(newFiles, f)
		}
	}

	if !found {
		return ErrFileNotFound
	}

	// Delete the object file
	if err := storage.DeleteObject(v.rootPath, objectID); err != nil {
		return err
	}

	// Update metadata
	meta.Files = newFiles
	if err := v.writeMetadata(masterKey, meta); err != nil {
		return fmt.Errorf("failed to update metadata: %w", err)
	}

	return nil
}

// readMetadata reads and decrypts the vault metadata
func (v *Vault) readMetadata(key []byte) (*storage.VaultMetadata, error) {
	encryptedMeta, err := storage.ReadMetadata(v.rootPath)
	if err != nil {
		return nil, err
	}

	// Decrypt metadata
	plainMeta, err := crypto.Decrypt(encryptedMeta, key)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON
	var meta storage.VaultMetadata
	if err := json.Unmarshal(plainMeta, &meta); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	return &meta, nil
}

// findFileByName performs fuzzy matching to find a file
// Tries: exact match, case-insensitive, contains, case-insensitive contains
func findFileByName(files []storage.FileMetadata, query string) *storage.FileMetadata {
	lowerQuery := strings.ToLower(query)

	// First pass: exact match
	for i := range files {
		if files[i].OriginalName == query {
			return &files[i]
		}
	}

	// Second pass: case-insensitive exact match
	for i := range files {
		if strings.ToLower(files[i].OriginalName) == lowerQuery {
			return &files[i]
		}
	}

	// Third pass: contains (case-insensitive)
	for i := range files {
		if strings.Contains(strings.ToLower(files[i].OriginalName), lowerQuery) {
			return &files[i]
		}
	}

	return nil
}

// matchesFileName checks if a filename matches the query (fuzzy matching)
func matchesFileName(filename, query string) bool {
	return strings.EqualFold(filename, query) ||
		strings.Contains(strings.ToLower(filename), strings.ToLower(query))
}

// writeMetadata encrypts and writes the vault metadata
func (v *Vault) writeMetadata(key []byte, meta *storage.VaultMetadata) error {
	// Marshal to JSON
	plainMeta, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("failed to serialize metadata: %w", err)
	}

	// Encrypt metadata
	encryptedMeta, err := crypto.Encrypt(plainMeta, key)
	if err != nil {
		return fmt.Errorf("failed to encrypt metadata: %w", err)
	}

	// Write to disk
	return storage.WriteMetadata(v.rootPath, encryptedMeta)
}
