# Usage Guide

Comprehensive guide to using Vaultix effectively.

## Basic Concepts

### What is a Vault?

A vault is a directory containing:

- A hidden `.vaultix/` folder with encrypted data
- Your encrypted files stored securely

### How Files Are Stored

```
my_vault/
└── .vaultix/
    ├── salt          # Random data for key derivation
    ├── meta          # Encrypted filenames and metadata
    └── objects/
        ├── a1b2c3d4.enc  # Your encrypted files
        └── e5f6g7h8.enc  # (with randomized names)
```

## Command Workflows

### Creating a New Vault

**Scenario**: You have a directory with sensitive files you want to encrypt.

```bash
cd ~/Documents/taxes_2024
ls
# tax_return.pdf  receipts.xlsx  statements.pdf

vaultix init
# Enter password: ****
# Confirm password: ****
# ✓ Vault initialized
# ✓ All files encrypted
# ✓ Original files deleted

ls
# .vaultix/  (only this remains)
```

**What happened:**

1. Vaultix scanned for all regular files
2. Created `.vaultix/` structure
3. Encrypted each file with your password
4. Securely deleted the originals

### Working with Files

#### List Files

```bash
vaultix list

# Output:
# Files in vault (3):
#   tax_return.pdf (245 KB, modified: 2024-03-15)
#   receipts.xlsx (128 KB, modified: 2024-03-10)
#   statements.pdf (512 KB, modified: 2024-03-01)
```

#### Extract a File

```bash
# Extract specific file
vaultix extract tax_return.pdf

# Now you can view it
cat tax_return.pdf
```

**Important**: The file is still in the vault! `extract` is non-destructive.

#### Add New Files

```bash
# Create a new sensitive file
echo "Bank account: 123456" > account.txt

# Add to vault
vaultix add account.txt

# Original is automatically deleted after encryption
ls account.txt
# ls: cannot access 'account.txt': No such file or directory
```

#### Remove Files from Vault

Two options:

**Option 1: Drop (extract + remove)**

```bash
vaultix drop old_tax_return.pdf
# ✓ Dropped: old_tax_return.pdf
# File is now decrypted AND removed from vault
```

**Option 2: Remove (delete without extracting)**

```bash
vaultix remove junk_file.pdf
# File removed from vault (no extraction)
```

## Advanced Usage

### Fuzzy File Matching

You don't need to type exact filenames:

```bash
# Vault contains: "2024_tax_return_final_v2.pdf"

vaultix extract 2024        # ✓ Matches
vaultix extract tax         # ✓ Matches
vaultix extract TAX         # ✓ Case-insensitive
vaultix extract final       # ✓ Matches
vaultix extract return.pdf  # ✓ Matches
```

**How it works:**

1. Exact match tried first
2. Case-insensitive exact match
3. Partial match (contains)

### Batch Operations

#### Extract Everything

```bash
# Extract all files to current directory
vaultix extract

# Extract to specific directory
vaultix extract . /tmp/decrypted/
```

#### Drop Everything

```bash
# Extract all files and clear vault
vaultix drop

# ⚠️  This empties the vault after extracting all files
```

#### Clear Vault (No Extraction)

```bash
vaultix clear
# ⚠️  This will DELETE all files from the vault WITHOUT extracting them.
# Continue? (yes/no): yes
# ✓ Vault cleared
```

### Working with Multiple Vaults

You can manage multiple vaults:

```bash
# Vault 1: Personal documents
cd ~/personal
vaultix list

# Vault 2: Work files
cd ~/work
vaultix list

# Vault 3: Old archives
vaultix list ~/archives
```

Each vault has its own:

- Encryption password
- Salt
- Encrypted files

### Directory Structure

#### Default Behavior (Current Directory)

```bash
cd my_vault
vaultix init      # Creates .vaultix in current dir
vaultix list      # Lists current vault
vaultix add file  # Adds to current vault
```

#### Explicit Paths

```bash
# Work from anywhere
vaultix list ~/my_vault
vaultix extract secret ~/my_vault
vaultix add newfile.txt ~/my_vault
```

## Common Patterns

### Daily Workflow

```bash
# Morning: Extract files you need
cd ~/work_vault
vaultix extract project_plan
vaultix extract api_keys

# ... do your work ...

# Evening: Add new/modified files
vaultix add updated_plan.pdf
vaultix add new_credentials.txt

# Clean up extracted files
rm project_plan.pdf api_keys.json
```

### Secure File Transfer

```bash
# Sender
cd ~/sensitive_docs
vaultix init
# Creates encrypted vault

# Transfer the entire directory
scp -r ~/sensitive_docs user@server:/tmp/

# Receiver
cd /tmp/sensitive_docs
vaultix extract
# Enter password (communicated securely!)
```

### Archiving Old Files

```bash
# Create archive vault
mkdir ~/archive_2023
cd ~/archive_2023

# Copy old files
cp ~/Documents/old_project/* .

# Encrypt everything
vaultix init

# Move to backup location
mv ~/archive_2023 /mnt/backup/
```

### Rotating Passwords

Currently, Vaultix doesn't support password changes. To rotate:

```bash
# Extract everything
vaultix extract

# Delete vault
rm -rf .vaultix

# Re-initialize with new password
vaultix init
# (Creates new vault with new password)
```

## Tips and Tricks

### Quick Reference Card

Save this to `~/.vaultix_commands`:

```
# Vaultix Quick Reference
init     - Encrypt directory
list     - Show files
extract  - Get file (keeps in vault)
drop     - Get file (removes from vault)
add      - Encrypt new file
remove   - Delete from vault
clear    - Delete all from vault
```

### Bash Aliases

Add to `~/.bashrc`:

```bash
alias vi='vaultix init'
alias vl='vaultix list'
alias ve='vaultix extract'
alias vd='vaultix drop'
alias va='vaultix add'
```

### Password Manager Integration

Use `pass` (Unix password manager):

```bash
# Store vault password
pass insert vaultix/my_vault

# Use it
vaultix list
# Password: $(pass vaultix/my_vault)
```

### Backup Script

```bash
#!/bin/bash
# backup-vault.sh

VAULT_DIR="$HOME/my_vault"
BACKUP_DIR="/mnt/backup/vaults"
DATE=$(date +%Y%m%d)

# Copy entire vault (encrypted)
cp -r "$VAULT_DIR" "$BACKUP_DIR/my_vault_$DATE"

echo "✓ Vault backed up to $BACKUP_DIR/my_vault_$DATE"
```

### Verification Script

```bash
#!/bin/bash
# verify-vault.sh

cd ~/my_vault

echo "Files in vault:"
vaultix list

echo ""
echo "Vault structure:"
ls -lR .vaultix/
```

## Troubleshooting

### "Vault not found"

```bash
# Check if .vaultix exists
ls -la .vaultix

# If not, you're in wrong directory
pwd
```

### "Failed to decrypt"

- Wrong password
- Corrupted vault
- Vault from different vaultix version

### "File not found"

- Check spelling: `vaultix list`
- Use fuzzy matching: `vaultix extract part_of_name`

### "Permission denied"

```bash
# Fix .vaultix permissions
chmod 700 .vaultix
chmod 600 .vaultix/*
```

## Best Practices

1. **Always backup vaults** before experimenting
2. **Test your password** immediately after creating a vault
3. **Use strong passwords** (16+ characters)
4. **Clean up extracted files** when done
5. **Keep vaults on encrypted drives** when possible
6. **Don't mix vault and work directories** (extract to separate folder)

## Next Steps

- Read [Commands Reference](commands.md) for complete command documentation
- Check [Security Model](security.md) to understand protection limits
- See [Examples](examples.md) for real-world usage scenarios
