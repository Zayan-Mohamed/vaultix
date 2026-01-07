# Best Practices

Guidelines for using Vaultix securely and effectively.

## Password Security

### Choosing a Strong Password

✓ **DO:**

- Use at least 16 characters
- Mix uppercase, lowercase, numbers, symbols
- Use a password manager to generate passwords
- Consider passphrases: "correct horse battery staple"
- Make it unique (don't reuse)

✗ **DON'T:**

- Use dictionary words
- Use personal information (birthdays, names)
- Use common passwords ("password123")
- Reuse passwords from other services
- Use passwords shorter than 12 characters

### Password Storage

✓ **DO:**

- Store vault passwords in a password manager
- Use hardware keys for password manager (YubiKey)
- Write down password and store in physical safe
- Use encrypted password databases

✗ **DON'T:**

- Store passwords in plaintext files
- Email passwords to yourself
- Share passwords over unsecured channels
- Store passwords in browser autofill

## File Management

### Before Encryption

✓ **DO:**

- Make backups before first encryption
- Verify files are complete and not corrupted
- Test password immediately after init
- Document what's in the vault

✗ **DON'T:**

- Encrypt your only copy
- Forget what password you used
- Encrypt system files
- Encrypt files you can't afford to lose

### During Use

✓ **DO:**

- Extract files to private directories
- Delete extracted files when done
- Use secure deletion tools for sensitive extracts
- Keep vault on encrypted filesystem

✗ **DON'T:**

- Extract to public/shared folders
- Leave decrypted files indefinitely
- Extract to cloud-synced directories
- Work directly in vault directory

### File Organization

```
Good structure:
~/vaults/
  ├── personal/
  │   └── .vaultix/
  ├── work/
  │   └── .vaultix/
  └── archive/
      └── .vaultix/

Bad structure:
~/Documents/
  ├── file1.txt
  ├── file2.pdf
  ├── .vaultix/        # Don't mix vault and regular files
  └── normal_doc.docx
```

## Backup Strategy

### What to Backup

✓ **Backup:** The entire vault directory (including `.vaultix/`)
✓ **Backup:** To multiple locations
✓ **Backup:** Encrypted vaults (safe for cloud storage)
✓ **Backup:** Regularly (automated schedule)

✗ **Don't backup:** Just the `.vaultix/` folder (need original directory too)
✗ **Don't backup:** Decrypted files to untrusted storage

### Backup Methods

**Local Backup:**

```bash
# Copy vault to external drive
cp -r ~/my_vault /mnt/backup/my_vault_$(date +%Y%m%d)

# Or use rsync
rsync -av ~/my_vault /mnt/backup/
```

**Cloud Backup:**

```bash
# Encrypted vaults are safe for cloud
rclone sync ~/my_vault remote:backups/my_vault

# Or tar + upload
tar czf my_vault.tar.gz ~/my_vault
aws s3 cp my_vault.tar.gz s3://my-bucket/backups/
```

**Verify Backups:**

```bash
# Test that backup is extractable
cd /tmp/test_restore
cp -r /mnt/backup/my_vault .
cd my_vault
vaultix list  # Enter password
```

## Operational Security

### System Security

✓ **DO:**

- Keep your OS updated
- Use antivirus/antimalware
- Enable firewall
- Use full-disk encryption
- Lock screen when away

✗ **DON'T:**

- Run untrusted software
- Disable security features
- Use admin/root unnecessarily
- Leave computer unlocked

### Network Security

✓ **DO:**

- Use VPN on public WiFi
- Use encrypted connections (HTTPS, SSH)
- Verify file integrity after transfer

✗ **DON'T:**

- Enter passwords on public WiFi
- Transfer vaults over unencrypted connections
- Use vaultix on shared/public computers

### Physical Security

✓ **DO:**

- Lock your computer when away
- Store backups in secure locations
- Encrypt backup drives
- Shred paper copies of passwords

✗ **DON'T:**

- Leave laptop unattended
- Store backups in obvious places
- Write passwords on sticky notes
- Leave vault passwords visible

## Workflow Best Practices

### Daily Workflow

**Morning:**

```bash
cd ~/work_vault
vaultix extract project_files
# Work on extracted files
```

**Evening:**

```bash
# Add updated files
vaultix add updated_file.pdf

# Clean up extracts
rm -f project_files/
```

### Project Workflow

**Starting Project:**

```bash
mkdir ~/projects/secret_project
cd ~/projects/secret_project
# Add initial files
vaultix init
```

**During Project:**

```bash
# Extract what you need
vaultix extract spec.pdf

# Modify
vim spec.pdf

# Re-add
vaultix add spec.pdf
```

**Ending Project:**

```bash
# Extract everything
vaultix extract

# Move out of vault
mv ~/projects/secret_project ~/archive/

# Clear vault
cd ~/projects/secret_project
vaultix clear
```

## Multi-Vault Management

### Organizing Vaults

```bash
~/vaults/
├── personal/       # Personal documents
├── work/           # Work files
├── financial/      # Tax, banking
├── projects/
│   ├── project_a/
│   └── project_b/
└── archive/
    ├── 2023/
    └── 2024/
```

### Password Strategy

**Option 1: One Master Password**

- Use same strong password for all vaults
- Easier to remember
- Higher risk if compromised

**Option 2: Different Passwords**

- Unique password per vault
- Better security
- Use password manager to track

**Option 3: Hierarchical**

- Weak password for low-security vaults
- Strong password for sensitive vaults
- Balance security and convenience

## Common Mistakes

### Mistake 1: Weak Passwords

❌ **Bad:**

```
Password: password123
```

✅ **Good:**

```
Password: Tr0ub4dor&3-correct-horse-battery
```

### Mistake 2: No Backups

❌ **Bad:**

```bash
# Only copy in vault
vaultix init
# Oops, hard drive died!
```

✅ **Good:**

```bash
vaultix init
# Backup vault
cp -r ~/vault /mnt/backup/
```

### Mistake 3: Extracting to Public Folders

❌ **Bad:**

```bash
cd ~/vault
vaultix extract passwords.txt
# Extracted to ~/vault/passwords.txt (visible!)
```

✅ **Good:**

```bash
cd ~/vault
vaultix extract passwords.txt
mv passwords.txt ~/private/temp/
# Work in private directory
```

### Mistake 4: Forgetting Password

❌ **Bad:**

```bash
vaultix init
# Enter password: ****
# (forget password)
# Files lost forever!
```

✅ **Good:**

```bash
vaultix init
# Enter password: <from password manager>
# Confirm: <paste from password manager>
vaultix list  # Test immediately
```

### Mistake 5: Mixing Vault and Work Directory

❌ **Bad:**

```bash
~/Documents/
├── .vaultix/
├── decrypted_file.txt  # Plaintext!
└── work_in_progress.pdf
```

✅ **Good:**

```bash
~/vault/
└── .vaultix/

~/work/
├── decrypted_file.txt  # Extracted here
└── work_in_progress.pdf
```

## Performance Tips

### Large Files

- Vaultix loads entire files into memory
- Splitting large files can improve performance
- Consider compressing before encryption

### Many Files

- Group related files in subdirectories
- Zip directories before adding to vault
- Use separate vaults for different projects

### SSD Optimization

```bash
# SSDs may not securely delete
# Use full-disk encryption + vaultix
# Or use secure delete tools:
shred -vfz -n 10 sensitive_file.txt
```

## Emergency Procedures

### Forgotten Password

**No recovery possible.** Prevention:

1. Use password manager
2. Write down and store securely
3. Test password immediately after creating vault

### Corrupted Vault

```bash
# Check vault structure
ls -la .vaultix/
# Should have: salt, meta, objects/

# Try listing files
vaultix list
# If it works, extract everything immediately

# If corrupted, restore from backup
cp -r /mnt/backup/my_vault ~/my_vault_restored
```

### Compromised Password

```bash
# Extract all files immediately
vaultix extract

# Create new vault with new password
rm -rf .vaultix
vaultix init
# Use NEW password

# Re-add files
vaultix add *
```

### Lost Backup

**Prevention is key:**

- Multiple backup locations
- Test backups regularly
- Automated backup schedule
- Off-site backups

## Compliance and Legal

### Data Retention

- Know your data retention requirements
- Don't over-retain sensitive data
- Use `vaultix clear` for permanent deletion
- Consider regulatory requirements (GDPR, HIPAA)

### Audit Trail

Vaultix doesn't log operations. If you need audit trails:

```bash
# Wrap commands with logging
echo "$(date): vaultix list" >> ~/.vaultix_audit.log
vaultix list
```

### Legal Considerations

- Encryption may be regulated in some jurisdictions
- You may be compelled to provide passwords
- Export controls may apply
- Consult legal counsel for compliance

## Conclusion

Security is a process, not a product. Vaultix is one tool in your security toolkit:

- ✓ Use strong passwords
- ✓ Make backups
- ✓ Follow operational security practices
- ✓ Keep systems updated
- ✓ Think before you act

Stay safe! 🔒
