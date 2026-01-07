package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/Zayan-Mohamed/vaultix/internal/storage"
	"github.com/Zayan-Mohamed/vaultix/internal/vault"
	"golang.org/x/term"
)

// readPassword reads a password from stdin without echoing
func readPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	password, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // Print newline after password input
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}
	return string(password), nil
}

// Init initializes a new vault at the specified path
func Init(args []string) error {
	vaultPath := "."
	if len(args) >= 1 {
		vaultPath = args[0]
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(vaultPath)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	// Check if vault already exists
	if storage.VaultExists(absPath) {
		return fmt.Errorf("vault already exists at: %s", absPath)
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(absPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Read password
	password, err := readPassword("Enter password: ")
	if err != nil {
		return err
	}

	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	// Confirm password
	confirmPassword, err := readPassword("Confirm password: ")
	if err != nil {
		return err
	}

	if password != confirmPassword {
		return fmt.Errorf("passwords do not match")
	}

	// Initialize vault
	v := vault.New(absPath)
	fmt.Println("Initializing vault and encrypting existing files...")
	if err := v.Initialize(password); err != nil {
		return fmt.Errorf("failed to initialize vault: %w", err)
	}

	fmt.Printf("✓ Vault initialized at: %s\n", absPath)
	fmt.Println("✓ All files have been encrypted")
	fmt.Println("✓ Original plaintext files have been securely deleted")
	return nil
}

// Add encrypts and adds a file to the vault
func Add(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: vaultix add <file> [vault-path]")
	}

	filePath := args[0]
	vaultPath := "."
	if len(args) >= 2 {
		vaultPath = args[1]
	}

	// Convert to absolute paths
	absVaultPath, err := filepath.Abs(vaultPath)
	if err != nil {
		return fmt.Errorf("invalid vault path: %w", err)
	}

	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}

	// Check if vault exists
	if !storage.VaultExists(absVaultPath) {
		return fmt.Errorf("vault not found at: %s", absVaultPath)
	}

	// Check if file exists
	if _, err := os.Stat(absFilePath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", absFilePath)
	}

	// Read password
	password, err := readPassword("Enter vault password: ")
	if err != nil {
		return err
	}

	// Add file
	v := vault.New(absVaultPath)
	if err := v.AddFile(password, absFilePath); err != nil {
		return fmt.Errorf("failed to add file: %w", err)
	}

	fileName := filepath.Base(absFilePath)
	fmt.Printf("File added: %s\n", fileName)
	return nil
}

// List displays all files in the vault
func List(args []string) error {
	vaultPath := "."
	if len(args) >= 1 {
		vaultPath = args[0]
	}

	// Convert to absolute path
	absVaultPath, err := filepath.Abs(vaultPath)
	if err != nil {
		return fmt.Errorf("invalid vault path: %w", err)
	}

	// Check if vault exists
	if !storage.VaultExists(absVaultPath) {
		return fmt.Errorf("vault not found at: %s", absVaultPath)
	}

	// Read password
	password, err := readPassword("Enter vault password: ")
	if err != nil {
		return err
	}

	// List files
	v := vault.New(absVaultPath)
	files, err := v.ListFiles(password)
	if err != nil {
		return fmt.Errorf("failed to list files: %w", err)
	}

	if len(files) == 0 {
		fmt.Println("Vault is empty")
		return nil
	}

	fmt.Printf("Files in vault (%d):\n", len(files))
	for _, f := range files {
		fmt.Printf("  %s (%d bytes, modified: %s)\n",
			f.OriginalName,
			f.Size,
			f.ModTime.Format("2006-01-02 15:04:05"))
	}

	return nil
}

// Extract decrypts and extracts a file from the vault
func Extract(args []string) error {
	vaultPath := "."
	fileName := ""
	outputPath := ""

	// Parse arguments flexibly
	if len(args) >= 1 {
		// If first arg is a vault path (contains / or is .), use it
		if args[0] == "." || args[0] == ".." || filepath.IsAbs(args[0]) || filepath.Dir(args[0]) != "." {
			vaultPath = args[0]
			if len(args) >= 2 {
				fileName = args[1]
			}
			if len(args) >= 3 {
				outputPath = args[2]
			}
		} else {
			// First arg is the filename
			fileName = args[0]
			if len(args) >= 2 {
				outputPath = args[1]
			}
		}
	}

	// Convert to absolute path
	absVaultPath, err := filepath.Abs(vaultPath)
	if err != nil {
		return fmt.Errorf("invalid vault path: %w", err)
	}

	// Check if vault exists
	if !storage.VaultExists(absVaultPath) {
		return fmt.Errorf("vault not found at: %s", absVaultPath)
	}

	// Read password
	password, err := readPassword("Enter vault password: ")
	if err != nil {
		return err
	}

	// Extract file(s)
	v := vault.New(absVaultPath)

	// If no filename specified, extract all files
	if fileName == "" {
		count, err := v.ExtractAllFiles(password, outputPath)
		if err != nil {
			return fmt.Errorf("failed to extract files: %w", err)
		}
		fmt.Printf("✓ Extracted %d file(s)\n", count)
		return nil
	}

	// Extract single file with fuzzy matching
	actualFileName, err := v.ExtractFile(password, fileName, outputPath)
	if err != nil {
		return fmt.Errorf("failed to extract file: %w", err)
	}

	fmt.Printf("✓ File extracted: %s\n", actualFileName)
	return nil
}

// Drop extracts and removes file(s) from the vault (destructive operation)
func Drop(args []string) error {
	vaultPath := "."
	fileName := ""
	outputPath := ""

	// Parse arguments flexibly (same as Extract)
	if len(args) >= 1 {
		if args[0] == "." || args[0] == ".." || filepath.IsAbs(args[0]) || filepath.Dir(args[0]) != "." {
			vaultPath = args[0]
			if len(args) >= 2 {
				fileName = args[1]
			}
			if len(args) >= 3 {
				outputPath = args[2]
			}
		} else {
			fileName = args[0]
			if len(args) >= 2 {
				outputPath = args[1]
			}
		}
	}

	// Convert to absolute path
	absVaultPath, err := filepath.Abs(vaultPath)
	if err != nil {
		return fmt.Errorf("invalid vault path: %w", err)
	}

	// Check if vault exists
	if !storage.VaultExists(absVaultPath) {
		return fmt.Errorf("vault not found at: %s", absVaultPath)
	}

	// Read password
	password, err := readPassword("Enter vault password: ")
	if err != nil {
		return err
	}

	// Drop file(s)
	v := vault.New(absVaultPath)

	// If no filename specified, drop all files
	if fileName == "" {
		count, err := v.DropAllFiles(password, outputPath)
		if err != nil {
			return fmt.Errorf("failed to drop files: %w", err)
		}
		fmt.Printf("✓ Dropped %d file(s) from vault\n", count)
		return nil
	}

	// Drop single file
	actualFileName, err := v.DropFile(password, fileName, outputPath)
	if err != nil {
		return fmt.Errorf("failed to drop file: %w", err)
	}

	fmt.Printf("✓ Dropped: %s (extracted and removed from vault)\n", actualFileName)
	return nil
}

// Clear removes all files from the vault without extracting them
func Clear(args []string) error {
	vaultPath := "."
	if len(args) >= 1 {
		vaultPath = args[0]
	}

	// Convert to absolute path
	absVaultPath, err := filepath.Abs(vaultPath)
	if err != nil {
		return fmt.Errorf("invalid vault path: %w", err)
	}

	// Check if vault exists
	if !storage.VaultExists(absVaultPath) {
		return fmt.Errorf("vault not found at: %s", absVaultPath)
	}

	// Read password
	password, err := readPassword("Enter vault password: ")
	if err != nil {
		return err
	}

	// Confirm dangerous operation
	fmt.Print("⚠️  This will DELETE all files from the vault WITHOUT extracting them. Continue? (yes/no): ")
	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "yes" {
		return fmt.Errorf("operation cancelled")
	}

	// Clear vault
	v := vault.New(absVaultPath)
	if err := v.ClearVault(password); err != nil {
		return fmt.Errorf("failed to clear vault: %w", err)
	}

	fmt.Println("✓ Vault cleared (all files removed)")
	return nil
}

// Remove removes a file from the vault
func Remove(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: vaultix remove <file> [vault-path]")
	}

	fileName := args[0]
	vaultPath := "."
	if len(args) >= 2 {
		vaultPath = args[1]
	}

	// Convert to absolute path
	absVaultPath, err := filepath.Abs(vaultPath)
	if err != nil {
		return fmt.Errorf("invalid vault path: %w", err)
	}

	// Check if vault exists
	if !storage.VaultExists(absVaultPath) {
		return fmt.Errorf("vault not found at: %s", absVaultPath)
	}

	// Read password
	password, err := readPassword("Enter vault password: ")
	if err != nil {
		return err
	}

	// Remove file
	v := vault.New(absVaultPath)
	if err := v.RemoveFile(password, fileName); err != nil {
		return fmt.Errorf("failed to remove file: %w", err)
	}

	fmt.Printf("File removed: %s\n", fileName)
	return nil
}

// PrintUsage prints the usage information
func PrintUsage() {
	fmt.Println("vaultix - Secure encrypted folder management")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  vaultix init [path]              Initialize vault (defaults to current directory)")
	fmt.Println("  vaultix add <file> [vault]       Add a file to the vault (defaults to current)")
	fmt.Println("  vaultix list [vault]             List files in the vault (defaults to current)")
	fmt.Println("  vaultix extract [file] [vault]   Extract file(s) - keeps in vault")
	fmt.Println("  vaultix drop [file] [vault]      Extract file(s) and remove from vault")
	fmt.Println("  vaultix remove <file> [vault]    Remove a file from vault (no extraction)")
	fmt.Println("  vaultix clear [vault]            Remove ALL files from vault (no extraction)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  cd my_secrets && vaultix init    # Encrypt all files in current directory")
	fmt.Println("  vaultix add newfile.txt          # Add file to current vault")
	fmt.Println("  vaultix list                     # List files in current vault")
	fmt.Println("  vaultix extract                  # Extract ALL files (keeps in vault)")
	fmt.Println("  vaultix extract secret           # Extract one file (keeps in vault)")
	fmt.Println("  vaultix drop secret              # Extract and remove from vault")
	fmt.Println("  vaultix drop                     # Extract all and clear vault")
	fmt.Println("  vaultix remove old.txt           # Remove file (no extraction)")
	fmt.Println("  vaultix clear                    # Remove all (no extraction, asks confirm)")
}
