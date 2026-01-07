# Troubleshooting

Solutions to common problems with vaultix.

## Installation Issues

### "Command not found: vaultix"

**Symptoms:**

```bash
$ vaultix init
bash: vaultix: command not found
```

**Cause:** Binary not in PATH or not installed.

**Solutions:**

**1. Check if binary exists:**

```bash
# Find vaultix binary
which vaultix
ls -l ~/bin/vaultix
```

**2. Add to PATH:**

```bash
# Temporary (current session)
export PATH="$PATH:$HOME/bin"

# Permanent (add to ~/.bashrc or ~/.zshrc)
echo 'export PATH="$PATH:$HOME/bin"' >> ~/.bashrc
source ~/.bashrc
```

**3. Reinstall:**

```bash
# Download and install
curl -LO https://github.com/Zayan-Mohamed/vaultix/releases/latest/download/vaultix-linux-amd64
chmod +x vaultix-linux-amd64
mv vaultix-linux-amd64 ~/bin/vaultix
```

---

### "Permission denied" when running vaultix

**Symptoms:**

```bash
$ ./vaultix
bash: ./vaultix: Permission denied
```

**Cause:** Binary not executable.

**Solution:**

```bash
chmod +x vaultix
./vaultix --version
```

---

### Build fails with "Go version too old"

**Symptoms:**

```
go: go.mod requires go >= 1.21
```

**Cause:** Go version too old.

**Solution:**

```bash
# Check Go version
go version

# Update Go
# Ubuntu/Debian
sudo snap install go --classic

# macOS
brew install go

# Windows
# Download from https://golang.org/dl/
```

---

## Vault Operation Issues

### "Vault not found"

**Symptoms:**

```bash
$ vaultix list
Error: vault not found
```

**Cause:** No `.vaultix` directory in current location.

**Solutions:**

**1. Check if vault exists:**

```bash
ls -la .vaultix
```

**2. Initialize vault:**

```bash
vaultix init
```

**3. Navigate to correct directory:**

```bash
cd /path/to/vault
vaultix list
```

---

### "Password incorrect" (but you're sure it's right)

**Symptoms:**

```bash
$ vaultix list
Enter password: ****
Error: password incorrect
```

**Possible causes:**

**1. Typing error:**

- Extra spaces
- Caps Lock enabled
- Wrong keyboard layout
- Auto-correct changed password

**Solution:**

```bash
# Type password in text editor to verify
echo "password123" > /tmp/pass.txt
cat /tmp/pass.txt | vaultix list
rm /tmp/pass.txt  # Clean up
```

**2. Different password used during init:**

**Solution:**

- Try to remember which password you used
- Restore from backup if available
- **No recovery possible if password forgotten**

**3. Vault corrupted:**

**Solution:**

- Restore from backup
- See "Vault Corrupted" section below

---

### "Vault already exists"

**Symptoms:**

```bash
$ vaultix init
Error: vault already exists
```

**Cause:** `.vaultix` directory already present.

**Solutions:**

**1. List existing vault:**

```bash
vaultix list
```

**2. Remove vault (DESTRUCTIVE):**

```bash
# WARNING: This deletes all encrypted files!
rm -rf .vaultix
vaultix init
```

**3. Use different directory:**

```bash
mkdir new_vault
cd new_vault
vaultix init
```

---

### "File not found in vault"

**Symptoms:**

```bash
$ vaultix extract document.pdf
Error: file not found in vault
```

**Causes & Solutions:**

**1. Typo in filename:**

```bash
# List all files
vaultix list

# Use fuzzy match
vaultix extract doc  # Will match "document.pdf"
```

**2. File never added:**

```bash
# Add file first
vaultix add document.pdf
vaultix extract document.pdf
```

**3. File was removed:**

```bash
# Check if file exists
vaultix list | grep document
```

---

### "Multiple files match"

**Symptoms:**

```bash
$ vaultix extract doc
Error: multiple files match 'doc': document.pdf, docs.txt
```

**Cause:** Fuzzy match is ambiguous.

**Solution:** Be more specific:

```bash
# Use more characters
vaultix extract docum  # Matches only "document.pdf"

# Or use exact name
vaultix extract document.pdf
```

---

## Performance Issues

### vaultix is very slow

**Symptoms:**

- Commands take many seconds
- High CPU usage
- System feels sluggish

**Causes & Solutions:**

**1. Large vault (many files):**

```bash
# Check number of files
vaultix list | wc -l

# Solution: Split into multiple vaults
mkdir vault1 vault2
# Move some files to vault2
```

**2. Slow storage:**

```bash
# Check disk speed
dd if=/dev/zero of=test.tmp bs=1M count=100
rm test.tmp

# Solution: Move vault to SSD if possible
```

**3. Limited RAM:**

```bash
# Check available memory
free -h

# Solution: Close other applications
```

---

### "Out of memory" error

**Symptoms:**

```bash
$ vaultix add large_file.bin
fatal error: runtime: out of memory
```

**Cause:** File too large for available RAM.

**Solutions:**

**1. Close other applications:**

```bash
# Free up memory
killall chrome
```

**2. Split large file:**

```bash
# Split into chunks
split -b 100M large_file.bin chunk_

# Add chunks
vaultix add chunk_*

# Later: extract and combine
vaultix extract
cat chunk_* > large_file.bin
```

**3. Increase system RAM:**

- Add more physical RAM
- Increase swap space (Linux)

---

## Data Corruption Issues

### "Vault corrupted"

**Symptoms:**

```bash
$ vaultix list
Error: vault corrupted: unexpected end of file
```

**Possible causes:**

- Disk errors
- Incomplete write (power loss, crash)
- Manual editing of `.vaultix` files

**Solutions:**

**1. Check disk health:**

```bash
# Linux
sudo smartctl -a /dev/sda

# macOS
diskutil verifyDisk disk0

# Windows
chkdsk C: /F
```

**2. Restore from backup:**

```bash
# Copy backup
cp -r /backup/my_vault ~/my_vault_restored
cd ~/my_vault_restored

# Verify
vaultix list
```

**3. Attempt recovery:**

```bash
# Check vault structure
ls -la .vaultix/
# Should have: salt, meta, objects/

# If salt is missing, vault is unrecoverable
# If meta is corrupted, try older backup
```

**4. Extract what you can:**

```bash
# Try to list files
vaultix list 2>&1 | tee errors.log

# Try to extract individually
for file in $(vaultix list 2>/dev/null); do
    vaultix extract "$file" || echo "Failed: $file"
done
```

---

### Encrypted files missing

**Symptoms:**

```bash
$ vaultix list
Error: encrypted file 3f9a2c1d.enc not found
```

**Cause:** `.vaultix/objects/` files deleted.

**Solutions:**

**1. Restore from backup:**

```bash
cp -r /backup/vault/.vaultix .
```

**2. If no backup, data is lost:**

- Encrypted files can't be recovered without the `.enc` files
- This is why backups are critical!

---

## Platform-Specific Issues

### Linux

#### "Operation not permitted" on encrypted filesystems

**Symptoms:**

```bash
$ vaultix init
Error: operation not permitted
```

**Cause:** Writing to encrypted filesystem (eCryptfs, EncFS).

**Solution:**

```bash
# Move vault to non-encrypted location
mkdir ~/vault
cd ~/vault
vaultix init

# Or adjust filesystem permissions
```

#### Secure delete doesn't work on ext4

**Cause:** ext4 file system doesn't guarantee overwrite.

**Solution:**

- Use encrypted disk (LUKS)
- Or use SSD TRIM feature
- Or use dedicated secure delete tools

---

### macOS

#### "Disk full" with plenty of space available

**Symptoms:**

```bash
$ df -h
/dev/disk1  500GB  250GB  250GB  50%  /Users

$ vaultix add large_file
Error: no space left on device
```

**Cause:** Spotlight indexing or Time Machine snapshots.

**Solution:**

```bash
# Disable Spotlight indexing for vault
mdutil -i off /path/to/vault

# Or exclude from Time Machine
tmutil addexclusion /path/to/vault
```

#### Keychain interference

**Symptoms:**
Password prompt appears multiple times.

**Solution:**

- vaultix doesn't use Keychain
- Check for other password manager interference

---

### Windows

#### "Access denied" in Program Files

**Symptoms:**

```
Error: access denied
```

**Cause:** Program Files requires admin rights.

**Solution:**

```powershell
# Use user directory
cd %USERPROFILE%\Documents
mkdir vault
cd vault
vaultix init
```

#### Antivirus quarantines vaultix

**Symptoms:**
Vaultix binary deleted or blocked.

**Solution:**

- Add vaultix to antivirus exclusions
- Download from official GitHub releases only

#### Hidden files not visible

**Symptoms:**
Can't see `.vaultix` directory in Explorer.

**Solution:**

```powershell
# Show hidden files
attrib -h .vaultix

# Or enable "Show hidden files" in Explorer
```

---

## Security Concerns

### "Someone accessed my vault!"

**Steps to take:**

**1. Change password immediately:**

```bash
# Extract all files
vaultix extract

# Remove old vault
rm -rf .vaultix

# Reinitialize with new password
vaultix init
```

**2. Check for malware:**

```bash
# Run antivirus scan
# Check for keyloggers
# Review running processes
```

**3. Review file access logs:**

```bash
# Linux
sudo journalctl -xe | grep vaultix

# Check file access times
stat .vaultix/meta
```

**4. Move vault to secure location:**

- Encrypted disk
- Offline storage
- Air-gapped computer (if very sensitive)

---

### Password compromised

**What to do:**

**1. Assume all encrypted data is compromised**

**2. Change password:**

```bash
vaultix extract        # Get all files
rm -rf .vaultix        # Delete vault
vaultix init           # New vault with new password
```

**3. Rotate sensitive data:**

- Change API keys
- Change other passwords stored in vault
- Update security questions

---

## Advanced Troubleshooting

### Enable debug output

```bash
# Set debug environment variable (if implemented)
export VAULTIX_DEBUG=1
vaultix list

# Or use Go's debugging
go run -race main.go list
```

### Check vault structure manually

```bash
# Verify .vaultix structure
tree .vaultix
# Should see:
# .vaultix/
# ├── salt
# ├── meta
# ├── config
# └── objects/
#     ├── *.enc files

# Check file sizes
du -sh .vaultix/objects/*

# Verify salt exists and is 32 bytes
ls -l .vaultix/salt
# Should be exactly 32 bytes
```

### Test with minimal vault

```bash
# Create test vault with one small file
mkdir /tmp/test_vault
cd /tmp/test_vault
echo "test" > test.txt
vaultix init

# If this works, original vault may be corrupted
```

### Recover from git

```bash
# If vault is in git repo
git log -- .vaultix/
git checkout HEAD~5 -- .vaultix/  # Go back 5 commits
vaultix list  # Try password
```

---

## Getting Help

### Before asking for help, collect:

1. **Vaultix version:**

   ```bash
   vaultix --version
   ```

2. **Operating system:**

   ```bash
   uname -a  # Linux/macOS
   systeminfo | findstr /B /C:"OS"  # Windows
   ```

3. **Error message:**

   ```bash
   vaultix list 2>&1 | tee error.log
   ```

4. **Vault structure (NO passwords!):**

   ```bash
   tree .vaultix
   ls -lah .vaultix/
   ```

5. **Steps to reproduce:**
   - What command did you run?
   - What did you expect?
   - What actually happened?

### Where to ask:

- **GitHub Issues:** https://github.com/Zayan-Mohamed/vaultix/issues
- **GitHub Discussions:** https://github.com/Zayan-Mohamed/vaultix/discussions

### What NOT to share:

- ❌ Passwords
- ❌ Encrypted file contents
- ❌ Personal file names
- ❌ Salt or encryption keys

---

## Emergency Recovery

### Lost password

**Bad news:** **No recovery possible.**

Vaultix uses strong encryption with no backdoor or password reset.

**Prevention:**

- Use password manager
- Write down password and store securely
- Test password immediately after creating vault

---

### Corrupted vault, no backup

**Options:**

**1. Try partial recovery:**

```bash
# List what can be recovered
vaultix list 2>&1 | grep -v Error

# Extract files one by one
for file in file1.txt file2.pdf; do
    vaultix extract "$file" || continue
done
```

**2. Forensic recovery:**

- Use `testdisk` or `photorec` to recover deleted `.enc` files
- Professional data recovery service (expensive)

**3. Accept data loss:**

- If vault is unrecoverable, data is lost
- Emphasizes importance of backups

---

### Accidentally cleared vault

**If you ran `vaultix clear`:**

**Immediate actions:**

1. **Don't write to disk** (increases recovery chance)
2. Use file recovery tools:

   ```bash
   # Linux
   sudo extundelete /dev/sda1 --restore-directory .vaultix/objects

   # macOS
   diskutil list
   # Use Time Machine or third-party recovery

   # Windows
   # Use Recuva or similar tools
   ```

**Future prevention:**

- Always backup before destructive operations
- Use confirmation prompts
- Test on copies, not originals

---

## Still Stuck?

If nothing here helps:

1. **Search existing issues:** https://github.com/Zayan-Mohamed/vaultix/issues
2. **Ask in discussions:** https://github.com/Zayan-Mohamed/vaultix/discussions
3. **Open new issue:** Provide all information listed in "Getting Help"

**For security vulnerabilities:** Email maintainers directly (see SECURITY.md), don't open public issue.

---

Remember: Backups are your best troubleshooting tool! 💾
