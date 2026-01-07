# Commands Reference

Complete reference for all Vaultix commands.

## Command Overview

| Command   | Description                         | Destructive |
| --------- | ----------------------------------- | ----------- |
| `init`    | Initialize vault and encrypt files  | ✓           |
| `add`     | Add file to vault                   | ✓           |
| `list`    | List encrypted files                | ✗           |
| `extract` | Extract file(s), keeps in vault     | ✗           |
| `drop`    | Extract and remove from vault       | ✓           |
| `remove`  | Remove file without extracting      | ✓           |
| `clear`   | Remove all files without extracting | ✓           |

## init

Initialize a new vault and automatically encrypt all files in the directory.

### Syntax

```bash
vaultix init [path]
```

### Parameters

- `path` (optional): Directory path. Defaults to current directory (`.`)

### Behavior

1. Scans directory for regular files (excludes hidden files and directories)
2. Creates `.vaultix/` structure
3. Encrypts all discovered files
4. Securely deletes original plaintext files
5. Reports success

### Examples

```bash
# Initialize in current directory
vaultix init

# Initialize specific path
vaultix init ~/secrets

# Initialize and encrypt all files
cd ~/Documents/sensitive
vaultix init
```

!!! warning "Destructive Operation"
Original files are permanently deleted after encryption. Make sure you have backups!

---

## add

Add a file to an existing vault.

### Syntax

```bash
vaultix add <file> [vault-path]
```

### Parameters

- `file` (required): Path to file to add
- `vault-path` (optional): Vault directory. Defaults to current directory (`.`)

### Behavior

1. Reads the specified file
2. Encrypts the file
3. Adds to vault metadata
4. Securely deletes the original file

### Examples

```bash
# Add file to current vault
vaultix add secret.txt

# Add file to specific vault
vaultix add document.pdf ~/my_vault

# Add file from another directory
vaultix add /tmp/keys.pem
```

---

## list

List all encrypted files in the vault.

### Syntax

```bash
vaultix list [vault-path]
```

### Parameters

- `vault-path` (optional): Vault directory. Defaults to current directory (`.`)

### Output

Shows:

- Filename
- Size in bytes
- Modification timestamp

### Examples

```bash
# List files in current vault
vaultix list

# List files in specific vault
vaultix list ~/other_vault
```

```
Files in vault (3):
  passwords.txt (1024 bytes, modified: 2026-01-07 15:30:00)
  api_keys.json (512 bytes, modified: 2026-01-07 14:20:00)
  notes.md (2048 bytes, modified: 2026-01-07 16:00:00)
```

---

## extract

Decrypt and extract file(s) from the vault. Files remain in the vault after extraction.

### Syntax

```bash
vaultix extract [file] [vault-path]
```

### Parameters

- `file` (optional): Filename or partial match. Omit to extract all files
- `vault-path` (optional): Vault directory. Defaults to current directory (`.`)

### Fuzzy Matching

Supports intelligent file matching:

1. **Exact match**: `passwords.txt`
2. **Case-insensitive**: `PASSWORDS.txt`
3. **Partial match**: `pass` matches `passwords.txt`

### Examples

```bash
# Extract all files
vaultix extract

# Extract specific file (exact match)
vaultix extract passwords.txt

# Extract with fuzzy matching
vaultix extract pass            # Matches "passwords.txt"
vaultix extract API             # Matches "api_keys.json"

# Extract to specific directory
vaultix extract . /tmp/output

# Extract from specific vault
vaultix extract secret ~/vault
```

---

## drop

Extract file(s) and remove them from the vault. This is a destructive operation.

### Syntax

```bash
vaultix drop [file] [vault-path]
```

### Parameters

- `file` (optional): Filename or partial match. Omit to drop all files
- `vault-path` (optional): Vault directory. Defaults to current directory (`.`)

### Behavior

1. Extracts the file(s)
2. Removes from vault
3. Deletes encrypted objects

### Examples

```bash
# Drop one file (extract and remove)
vaultix drop old_password.txt

# Drop with fuzzy matching
vaultix drop secret

# Drop all files (extracts everything and clears vault)
vaultix drop
```

!!! danger "Destructive"
Files are removed from the vault after extraction. Make sure you have the decrypted files before closing your terminal!

---

## remove

Remove a file from the vault **without** extracting it. Permanently deletes the encrypted data.

### Syntax

```bash
vaultix remove <file> [vault-path]
```

### Parameters

- `file` (required): Filename or partial match
- `vault-path` (optional): Vault directory. Defaults to current directory (`.`)

### Examples

```bash
# Remove file from current vault
vaultix remove old_file.txt

# Remove with fuzzy matching
vaultix remove OLD

# Remove from specific vault
vaultix remove file.txt ~/vault
```

!!! warning "No Extraction"
File is deleted without extracting. Use `drop` if you want to extract first.

---

## clear

Remove **all** files from the vault without extracting them. Requires confirmation.

### Syntax

```bash
vaultix clear [vault-path]
```

### Parameters

- `vault-path` (optional): Vault directory. Defaults to current directory (`.`)

### Behavior

1. Prompts for password
2. Asks for confirmation (`yes` to proceed)
3. Deletes all encrypted objects
4. Clears metadata

### Examples

```bash
# Clear current vault
vaultix clear
# ⚠️  This will DELETE all files from the vault WITHOUT extracting them. Continue? (yes/no): yes
# ✓ Vault cleared (all files removed)

# Clear specific vault
vaultix clear ~/vault
```

!!! danger "Extremely Destructive"
All encrypted data is permanently deleted. Cannot be undone!

---

## Common Patterns

### Secure a Directory

```bash
cd ~/sensitive
vaultix init
```

### Daily Workflow

```bash
# See what's in the vault
vaultix list

# Extract what you need
vaultix extract passwords

# Add new files
vaultix add newsecret.txt

# Clean up old files
vaultix remove old.txt
```

### Move Files Out

```bash
# Extract all files
vaultix extract

# Then clear the vault
vaultix clear
```

Or use drop:

```bash
# Extract all and clear in one command
vaultix drop
```

### Fuzzy Extraction

```bash
# Instead of typing "super_secret_passwords_2024.txt"
vaultix extract 2024
# or
vaultix extract PASSWORD
# or
vaultix extract super
```

## Password Security

All commands that access the vault require password entry. Passwords are:

- ✓ Never stored on disk
- ✓ Never logged
- ✓ Only exist in memory during operation
- ✓ Hidden during input (not echoed to terminal)

!!! tip "Password Management"
Consider using a password manager to generate and store your vault password.
