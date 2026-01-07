# Frequently Asked Questions

Common questions about vaultix and encrypted folders.

## General

### What is vaultix?

Vaultix is a command-line tool that encrypts files in a directory using password-based encryption. It's designed for protecting sensitive files on your local computer.

### Is vaultix free?

Yes, vaultix is open-source software released under the MIT License. It's completely free to use, modify, and distribute.

### What platforms does vaultix support?

- ✅ Linux (all distributions)
- ✅ macOS (Intel and Apple Silicon)
- ✅ Windows (Windows 10+)

### Is vaultix a password manager?

No. Vaultix encrypts files and folders. For managing passwords, use dedicated password managers like:

- Bitwarden
- 1Password
- KeePassXC

### Can I use vaultix on multiple computers?

Yes! Your encrypted vault is just a normal folder. Copy it to any computer and use vaultix with the same password.

## Security

### How secure is vaultix?

Vaultix uses industry-standard cryptography:

- **Encryption**: AES-256-GCM
- **Key derivation**: Argon2id
- **Randomness**: crypto/rand (cryptographically secure)

These are the same algorithms used by banks and military applications.

### Can anyone break the encryption?

With a strong password, no. AES-256 would take billions of years to brute-force with current technology.

**However**, weak passwords can be cracked quickly. Use strong passwords!

### What if I forget my password?

**Your files are permanently unrecoverable.** There is no password reset, no backdoor, no recovery mechanism. This is by design.

**Prevention**:

- Use a password manager
- Write password down and store securely
- Test password immediately after creating vault

### Can someone access my files if they steal my laptop?

If your laptop is stolen:

- ✅ Vaultix-encrypted files are safe (with strong password)
- ❌ Extracted (decrypted) files are readable
- ❌ Files currently being edited are readable

**Additional protection**:

- Use full-disk encryption (LUKS, FileVault, BitLocker)
- Lock your computer when away
- Don't leave decrypted files lying around

### Is the `.vaultix` folder hidden for security?

No. Hiding `.vaultix` is purely for cleanliness, not security. The **encryption** protects your files, not hiding them.

Anyone can see the `.vaultix` folder exists, but they can't decrypt its contents without your password.

### Can I safely upload encrypted vaults to cloud storage?

Yes! Encrypted vaults are safe to store in:

- Dropbox
- Google Drive
- OneDrive
- Any cloud service

**But**:

- ❌ Never upload decrypted files to untrusted storage
- ✅ Use unique, strong passwords for sensitive vaults
- ✅ Consider using dedicated encrypted cloud storage (Proton Drive, etc.)

### What if someone installs a keylogger?

Vaultix can't protect against keyloggers. If your system is compromised:

- Attacker can capture your password as you type
- Attacker can read files while you work on them

**Protection**:

- Keep your OS and antivirus updated
- Don't run untrusted software
- Use hardware security keys for important accounts
- Consider using an air-gapped computer for extremely sensitive data

## Usage

### Do I need to type my password every time?

Yes. Vaultix has no "unlock" state. You enter your password each time you run a command.

This is intentional - it ensures:

- Files are only accessible when you explicitly decrypt them
- No background daemon that could be attacked
- No risk of leaving vault "unlocked"

### Can I change my password?

Not currently. To change password:

```bash
# 1. Extract all files
vaultix extract

# 2. Delete vault
rm -rf .vaultix

# 3. Re-initialize with new password
vaultix init
```

### Can multiple people share a vault?

Technically yes, if they share the password. But this is not recommended:

**Problems**:

- Everyone knows the password
- Can't revoke access individually
- No audit trail of who accessed what

**Better alternatives**:

- Give each person their own vault copy
- Use proper access control systems for teams
- Consider enterprise encryption solutions for organizations

### How do I organize multiple vaults?

```bash
~/vaults/
├── personal/       # Personal documents
├── work/           # Work files
├── financial/      # Tax documents
└── projects/
    ├── project_a/
    └── project_b/
```

Each vault can have a different password.

### Why can't I extract files to a different directory?

By design, `extract` outputs to the current directory. This prevents accidentally extracting sensitive files to insecure locations.

**Workaround**:

```bash
# Extract in vault directory
cd ~/vault
vaultix extract file.txt

# Move to target location
mv file.txt ~/destination/
```

### Can I encrypt a file without moving it to a vault directory?

No. Vaultix operates on directories, not individual files.

**Alternative workflow**:

```bash
# Create vault
mkdir ~/secure
cd ~/secure

# Copy file
cp ~/Documents/sensitive.pdf .

# Encrypt
vaultix init
```

## Technical

### What happens to my original files?

When you run `vaultix init` or `add`:

1. Original file is read
2. Encrypted copy is stored in `.vaultix/objects/`
3. **Original file is securely deleted**

The original is gone. Only the encrypted version remains.

### Can I recover deleted files?

No. When vaultix deletes originals or when you use `clear`, files are securely overwritten.

**But**: If you have backups or disk snapshots, deleted files might be recoverable from there.

### Why are encrypted filenames random?

```bash
.vaultix/objects/
├── 3f9a2c1d.enc
└── 91bd77aa.enc
```

This prevents information leakage. If filenames were preserved, an attacker could learn:

- How many files you have
- File sizes
- File names (potentially sensitive)

Original filenames are encrypted in the metadata.

### How much space does encryption add?

Minimal overhead:

- ~32 bytes per file (nonce + tag)
- Small metadata file

A 1 MB file becomes approximately 1.000032 MB encrypted.

### Can I encrypt very large files?

Currently, vaultix loads entire files into memory. This means:

- 100 MB file: Fine
- 1 GB file: Probably okay
- 10 GB file: May cause issues

**Workaround for large files**:

```bash
# Split large file
split -b 100M huge_file.bin chunk_

# Encrypt chunks
vaultix add chunk_*

# Later: extract and combine
vaultix extract
cat chunk_* > huge_file.bin
```

### Does vaultix compress files?

No. Compression is not built-in.

**To compress before encrypting**:

```bash
tar czf archive.tar.gz my_files/
vaultix add archive.tar.gz
```

### What's inside the `.vaultix` folder?

```bash
.vaultix/
├── salt           # Random salt for key derivation
├── meta           # Encrypted metadata (filenames, sizes)
├── config         # Vault configuration
└── objects/
    └── *.enc      # Encrypted file contents
```

### Can I edit the `.vaultix` files manually?

**Don't**. Editing `.vaultix` contents will corrupt your vault.

If you need to modify files:

1. Extract file
2. Edit extracted copy
3. Add modified file back

## Troubleshooting

### "Password incorrect" but I know it's right

Possible causes:

- Extra spaces in password
- Caps Lock is on
- Different keyboard layout
- Typed too quickly (missed characters)

**Solutions**:

- Type password in text editor first, then copy/paste
- Try typing slowly and carefully
- Check for Unicode characters (if you copy/pasted)

### "Vault corrupted" error

Possible causes:

- `.vaultix` folder was modified
- Disk errors
- Interrupted encryption process

**Solutions**:

- Restore from backup
- Run filesystem check (`fsck`, `chkdsk`)
- If you have backups, restore the entire vault directory

### "Out of memory" when adding large files

Your file is too large for available RAM.

**Solutions**:

- Close other applications
- Add more RAM
- Split file into smaller chunks
- Use a computer with more memory

### Vaultix is slow on my vault

Large vaults (1000+ files) may be slow because:

- Metadata must be decrypted for every operation
- All filenames are searched for fuzzy matching

**Solutions**:

- Split into multiple smaller vaults
- Use exact filenames (faster than fuzzy)
- Upgrade to faster storage (SSD)

### Can't find vaultix command after installation

Check your `$PATH`:

```bash
# Linux/macOS
echo $PATH
export PATH="$PATH:$HOME/bin"

# Windows
echo %PATH%
set PATH=%PATH%;C:\tools
```

## Comparison

### Vaultix vs. VeraCrypt?

**VeraCrypt**:

- Creates encrypted container files
- Mounts as virtual drive
- OS-level encryption

**Vaultix**:

- Encrypts individual files
- CLI tool, no mounting
- File-level encryption

Use **VeraCrypt** if you want drive-level encryption.
Use **vaultix** if you want simple file encryption.

### Vaultix vs. 7-Zip with password?

**7-Zip**:

- Creates compressed archives
- Encryption is weaker (AES-256 but with PBKDF2)
- Requires extracting entire archive

**Vaultix**:

- Individual file access
- Stronger key derivation (Argon2id)
- Extract only what you need

Use **7-Zip** if you need compression + encryption for archival.
Use **vaultix** if you need a working encrypted directory.

### Vaultix vs. GPG?

**GPG**:

- Industry standard
- Public key cryptography
- Complex to use

**Vaultix**:

- Simple password-based
- Easy to use
- No key management

Use **GPG** if you need public key encryption or digital signatures.
Use **vaultix** if you want simple password-based encryption.

### Vaultix vs. BitLocker/FileVault?

**BitLocker/FileVault**:

- Full-disk encryption
- OS-integrated
- Transparent to applications

**Vaultix**:

- Directory-level
- Cross-platform
- Explicit encryption/decryption

Use **BitLocker/FileVault** as your base layer.
Use **vaultix** for additional protection of specific directories.

## Getting Help

### Where can I report bugs?

GitHub Issues: https://github.com/zayan-mohamed/vaultix/issues

Please include:

- Vaultix version (`vaultix --version`)
- Operating system
- Steps to reproduce
- Error messages

### Where can I request features?

GitHub Discussions: https://github.com/zayan-mohamed/vaultix/discussions

### Is there a mailing list or chat?

Not yet. For now, use GitHub Discussions.

### How can I contribute?

See [contributing.md](contributing.md) for details.

### I found a security vulnerability!

**Don't open a public issue.**

Email maintainers directly (see SECURITY.md).

## Other Questions

### Why is it called "vaultix"?

- **Vault**: Secure storage
- **ix**: Unix/Linux tradition (suffix for tools)
- Short, memorable, CLI-friendly

### Who maintains vaultix?

Vaultix is maintained by [@zayan-mohamed](https://github.com/zayan-mohamed) and contributors.

### What's the project roadmap?

See [GitHub Issues](https://github.com/zayan-mohamed/vaultix/issues) and [Discussions](https://github.com/zayan-mohamed/vaultix/discussions) for planned features.

### Can I use vaultix in commercial projects?

Yes! Vaultix is MIT licensed. You can use it commercially, modify it, and distribute it. Attribution is appreciated but not required.

### Does vaultix have telemetry or phone home?

**No.** Vaultix:

- Doesn't collect any data
- Doesn't send anything over the network
- Doesn't check for updates automatically
- Doesn't require an account or registration

### Is vaultix audited?

Not yet. Vaultix uses audited standard library cryptography (Go's `crypto` packages), but the vaultix code itself has not undergone a professional security audit.

If you'd like to sponsor an audit, please reach out!

---

**Didn't find your answer?**

- Check the [documentation](index.md)
- Search [GitHub Issues](https://github.com/zayan-mohamed/vaultix/issues)
- Ask in [GitHub Discussions](https://github.com/zayan-mohamed/vaultix/discussions)
