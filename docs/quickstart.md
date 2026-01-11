# Quick Start

This guide will help you get started with Vaultix in under 5 minutes.

## Your First Vault

### Step 1: Create a Directory

```bash
mkdir ~/my_secrets
cd ~/my_secrets
```

### Step 2: Add Some Files

```bash
echo "My secret password" > passwords.txt
echo "API_KEY=abc123" > api_keys.env
echo "Private notes" > notes.md
```

### Step 3: Initialize the Vault

```bash
vaultix init
```

You'll be prompted for a password:

```
Enter password: ****
Confirm password: ****
Initializing vault and encrypting existing files...
✓ Vault initialized at: /home/user/my_secrets
✓ All files have been encrypted
✓ Original plaintext files have been securely deleted

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  IMPORTANT: RECOVERY KEY
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Your recovery key (save this in a secure location):

  5025f74e-c5d7a54a-7b99c87b-78cca1a0-61854d30-fb0d2783-a9df7067-b67ad345

This recovery key can unlock your vault if you forget your password.
Store it safely - if you lose both your password AND recovery key,
your vault will be permanently unrecoverable.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

!!! warning "Critical - Save Your Recovery Key!"
    **Copy and save this recovery key immediately!** It's displayed only once during initialization.
    - Print it and store in a safe
    - Save to a password manager
    - Store in a different location from your vault
    - **Never** lose both password AND recovery key

!!! info "Dual Unlock Methods"
    You can unlock your vault with either:
    - Your password (what you just set)
    - Your recovery key (the hex string above)

### Step 4: Verify Files are Encrypted

```bash
ls -la
```

You should only see:

```
drwx------ 3 user user 4096 Jan  7 17:00 .vaultix
```

Your original files are gone—they're now safely encrypted inside `.vaultix/`!

### Step 5: List Encrypted Files

```bash
vaultix list
```

Enter your password, and you'll see:

```
Files in vault (3):
  passwords.txt (20 bytes, modified: 2026-01-07 17:00:00)
  api_keys.env (15 bytes, modified: 2026-01-07 17:00:00)
  notes.md (13 bytes, modified: 2026-01-07 17:00:00)
```

### Step 6: Extract a File

```bash
vaultix extract passwords
```

The file is extracted with its original name:

```bash
cat passwords.txt
# My secret password
```

!!! note "Non-Destructive"
`extract` keeps the file in the vault. Use `drop` to remove it after extracting.

## Common Workflows

### Secure Existing Directory

```bash
cd ~/Documents/sensitive
vaultix init
# Done! All files are now encrypted
```

### Add New Files to Vault

```bash
# Create a new file
echo "New secret" > newsecret.txt

# Add it to the vault
vaultix add newsecret.txt

# File is automatically encrypted and deleted
```

### Extract All Files

```bash
vaultix extract
# Extracts everything to current directory
```

### Drop a File (Extract and Remove)

```bash
vaultix drop old_file.txt
# Extracts the file AND removes it from vault
```

### Fuzzy File Matching

You don't need to type exact filenames:

```bash
# Vault contains "secret_passwords_2024.txt"
vaultix extract secret      # ✓ Matches
vaultix extract PASSWORD    # ✓ Case-insensitive
vaultix extract 2024        # ✓ Partial match
```

## Understanding the Vault Structure

After initialization, your directory looks like:

```
my_secrets/
└── .vaultix/
    ├── salt          # Random salt for key derivation
    ├── meta          # Encrypted metadata (filenames, sizes)
    └── objects/
        ├── a1b2c3d4.enc
        └── e5f6g7h8.enc
```

- **salt**: Random data used for secure key derivation
- **meta**: Encrypted list of files and their metadata
- **objects/**: Encrypted file contents with random names

!!! danger "Backup Warning"
If you lose the `.vaultix/` directory, you lose your data permanently!

## Next Steps

Now that you understand the basics:

- Check out [Commands](commands.md) for all available commands
- Read [Best Practices](best-practices.md) for security tips
- See [Examples](examples.md) for advanced usage scenarios

## Quick Reference

| Task         | Command                  |
| ------------ | ------------------------ |
| Create vault | `vaultix init`           |
| Add file     | `vaultix add <file>`     |
| List files   | `vaultix list`           |
| Extract file | `vaultix extract <file>` |
| Extract all  | `vaultix extract`        |
| Drop file    | `vaultix drop <file>`    |
| Remove file  | `vaultix remove <file>`  |
| Clear vault  | `vaultix clear`          |
