# vaultix

<p align="center">
  <img src="https://img.shields.io/github/v/release/Zayan-Mohamed/vaultix?style=for-the-badge" alt="Release"/>
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version"/>
  <img src="https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge" alt="License"/>
  <img src="https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-blue?style=for-the-badge" alt="Platform"/>
  <img src="https://img.shields.io/badge/Encryption-AES--256--GCM-red?style=for-the-badge&logo=gnuprivacyguard&logoColor=white" alt="Encryption"/>
  <img src="https://img.shields.io/github/actions/workflow/status/Zayan-Mohamed/vaultix/build.yml?style=for-the-badge&label=Build" alt="Build Status"/>
</p>

<p align="center">
  <strong>A cross-platform command-line tool for managing password-protected encrypted folders</strong>
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#installation">Installation</a> •
  <a href="#quick-start">Quick Start</a> •
  <a href="https://zayan-mohamed.github.io/vaultix">Documentation</a> •
  <a href="#security">Security</a> •
  <a href="#contributing">Contributing</a>
</p>

---

## 📖 Overview

vaultix is a secure, lightweight CLI tool that encrypts files in place using military-grade cryptography. No cloud, no services, no complexity—just strong encryption for your sensitive files.

### Key Highlights

- 🔒 **Strong Encryption**: AES-256-GCM with Argon2id key derivation
- 🚀 **Zero Dependencies**: Single static binary, no runtime requirements
- 💻 **Cross-Platform**: Linux, macOS, and Windows support
- 🎯 **Simple UX**: Intuitive commands with smart defaults
- 🔐 **No Password Storage**: Passwords exist only in memory
- 📦 **Portable**: Encrypted vaults work across all platforms

## ✨ Features

✅ **Automatic Encryption** - Initialize a vault and all files are encrypted instantly  
✅ **Master Key Encryption** - Random 256-bit master key protects all vault data  
✅ **Recovery Key Support** - Unlock vault if you forget your password  
✅ **Dual Unlock Methods** - Use password OR recovery key  
✅ **Fuzzy File Matching** - No need to type exact filenames  
✅ **Default to Current Directory** - Less typing, more doing  
✅ **Extract or Drop** - Extract files while keeping in vault, or drop them out  
✅ **Secure Deletion** - Original files are overwritten before deletion  
✅ **Hidden Metadata** - Even filenames are encrypted  
✅ **No Background Process** - Runs only when you invoke it

## Security Model

### Cryptography

vaultix uses a **master key encryption model** with industry-standard cryptographic primitives:

- **Master Key**: Random 256-bit key generated per vault (encrypted, never stored in plaintext)
- **Password Protection**: Master key encrypted with Argon2id-derived key (64MB memory, 1 iteration, 4 threads)
- **Recovery Key**: Random 256-bit key that can decrypt the master key (backup unlock method)
- **Data Encryption**: AES-256-GCM authenticated encryption for all vault data
- **Randomness**: Go's `crypto/rand` package for all cryptographic random generation

### Threat Model

vaultix protects against:

- Unauthorized access to files at rest (assuming strong password)
- Accidental exposure of file contents
- Casual inspection of encrypted data

vaultix does **not** protect against:

- Weak passwords (use a strong, unique password)
- Malware or keyloggers on your system
- Physical access to your computer while unlocked
- Attacks on the underlying operating system
- Side-channel attacks or memory analysis
- Coercion or legal compulsion

### Important Limitations

- **Dual authentication required**: Keep both password AND recovery key safe
- **No password reset**: If you lose BOTH password and recovery key, data is permanently lost
- **Recovery key is critical**: Store it safely (printed, secure password manager, etc.)
- **No automatic backups**: You are responsible for backing up your vaults
- **Single-user design**: No multi-user or sharing capabilities
- **Files only**: Cannot encrypt directories (add files individually)

## 📦 Installation

### Download Pre-built Binary (Recommended)

Download the latest release for your platform:

**Linux:**

```bash
curl -LO https://github.com/Zayan-Mohamed/vaultix/releases/latest/download/vaultix-linux-amd64
chmod +x vaultix-linux-amd64
sudo mv vaultix-linux-amd64 /usr/local/bin/vaultix
```

**macOS:**

```bash
curl -LO https://github.com/Zayan-Mohamed/vaultix/releases/latest/download/vaultix-darwin-amd64
chmod +x vaultix-darwin-amd64
sudo mv vaultix-darwin-amd64 /usr/local/bin/vaultix
```

**Windows (PowerShell as Admin):**

```powershell
Invoke-WebRequest -Uri "https://github.com/Zayan-Mohamed/vaultix/releases/latest/download/vaultix-windows-amd64.exe" -OutFile "vaultix.exe"
Move-Item vaultix.exe C:\Windows\System32\
```

### Build from Source

Requires Go 1.21 or higher:

```bash
git clone https://github.com/Zayan-Mohamed/vaultix.git
cd vaultix
go build -o vaultix
```

Then move the binary to your PATH:

**Linux/macOS:**

```bash
sudo mv vaultix /usr/local/bin/
```

**Windows:**

```powershell
Move-Item vaultix.exe C:\Windows\System32\
```

---

## 🚀 Quick Start

```bash
cd ~/my_secrets

# Initialize vault (encrypts all files automatically)
vaultix init
# Enter password: ****
# Confirm password: ****
# ✓ Vault initialized
# ✓ All files encrypted
# ✓ Original files securely deleted
#
# ⚠️  IMPORTANT: RECOVERY KEY
# Your recovery key: 5025f74e-c5d7a54a-7b99c87b-78cca1a0-...
# Save this recovery key in a secure location!
# It can unlock your vault if you forget your password.

# List encrypted files
vaultix list
# Files in vault (3):
#   passwords.txt
#   api_keys.json
#   private_key.pem

# Extract a file (keeps in vault)
vaultix extract passwords
# ✓ File extracted: passwords.txt

# Drop a file (extracts and removes from vault)
vaultix drop api_keys
# ✓ Dropped: api_keys.json (extracted and removed from vault)

# Extract all files
vaultix extract
# ✓ Extracted 3 file(s)
```

## 📚 Usage

### Commands

| Command          | Description                            | Example                  |
| ---------------- | -------------------------------------- | ------------------------ |
| `init [path]`    | Initialize vault and encrypt all files | `vaultix init`           |
| `add <file>`     | Add file to vault                      | `vaultix add secret.txt` |
| `list [path]`    | List encrypted files                   | `vaultix list`           |
| `extract [file]` | Extract file(s), keeps in vault        | `vaultix extract`        |
| `drop <file>`    | Extract and remove from vault          | `vaultix drop secret`    |
| `remove <file>`  | Remove file from vault (no extract)    | `vaultix remove old.txt` |
| `clear [path]`   | Remove all files from vault            | `vaultix clear`          |
| `recover [file]` | Unlock vault using recovery key        | `vaultix recover`        |

> 💡 **Pro Tip**: Most commands default to current directory, so you rarely need to specify paths!

### Advanced Usage

```bash
# Fuzzy file matching (case-insensitive)
vaultix extract SECRET    # Matches "secret_document.pdf"
vaultix extract api       # Matches "api_keys.json"

# Extract all to specific directory
vaultix extract . /tmp/decrypted/

# Work with specific vault path
vaultix list ~/other_vault
vaultix extract document ~/other_vault
```

---

## 🏗️ How It Works

When you initialize a vault at a path (e.g., `./my_secrets`), vaultix creates a hidden `.vaultix/` directory inside:

```
my_secrets/
└── .vaultix/
    ├── meta          # Encrypted metadata (filenames, sizes, timestamps)
    ├── salt          # Random salt for password-based key derivation
    ├── master.key    # Master key encrypted with password-derived key
    ├── recovery.key  # Master key encrypted with recovery key
    └── objects/
        ├── 3f9a2c1d.enc  # Encrypted file data
        └── 91bd77aa.enc  # Encrypted file data
```

### Security Details

1. **Master key encryption**: A random 256-bit master key encrypts all vault data
2. **Dual unlock methods**: Master key can be decrypted with password OR recovery key
3. **No plaintext keys**: Master key never stored in plaintext on disk
4. **No passwords stored**: Your password exists only in memory during operations
5. **Encrypted metadata**: Even filenames are encrypted with the master key
6. **Obfuscated object names**: Encrypted files have random IDs
7. **Salt per vault**: Each vault has a unique random salt
8. **Authentication**: AES-GCM provides both encryption and integrity verification

### Authentication and Unlock

Password/recovery key correctness is verified by successful decryption of the master key. There are no stored password hashes. This means:

- Incorrect password/recovery key = decryption failure
- No way to test credentials without attempting decryption
- Recovery key provides backup access if password is forgotten
- If you lose BOTH password AND recovery key, data is permanently lost

---

## 💡 Best Practices

### Password and Recovery Key Management

**Password Selection** - Use a strong, unique password:

- ✅ At least 16 characters
- ✅ Mix of letters, numbers, and symbols
- ✅ Not used anywhere else
- ✅ Not easily guessable

Consider using a password manager to generate and store your vault password.

**Recovery Key Storage** - Your recovery key is displayed ONCE during vault initialization:

- ✅ Print it and store in a safe location
- ✅ Save to a password manager as a secure note
- ✅ Store in a separate secure location from your vault
- ⚠️ Never store recovery key inside the vault itself
- ⚠️ If you lose both password AND recovery key, data is permanently lost

### Backup Strategy

- The entire vault directory (including `.vaultix/`) must be backed up
- Test your backups by extracting files from backup copies
- Encrypted vaults are safe to backup to cloud storage
- ⚠️ Losing `.vaultix/` = permanent data loss

### File Management

- When adding files with `add`, original files are NOT automatically deleted
- Use secure deletion tools if you need to remove originals
- Keep temporary extractions out of the vault directory
- Don't extract sensitive files to public/shared directories

### Operational Security

- Don't enter passwords where they might be logged
- Don't use vaultix over remote connections without encryption
- Close your terminal after vault operations
- Consider using full-disk encryption alongside vaultix
- Be aware of swap files and hibernation dumps

---

## 🌐 Platform Support

vaultix works identically on **Linux**, **macOS**, and **Windows**.

---

## 📖 Documentation

Full documentation is available at: **[https://zayan-mohamed.github.io/vaultix](https://zayan-mohamed.github.io/vaultix)**

- 📘 [Installation Guide](https://zayan-mohamed.github.io/vaultix/installation/)
- 🚀 [Quick Start](https://zayan-mohamed.github.io/vaultix/quickstart/)
- 📚 [Command Reference](https://zayan-mohamed.github.io/vaultix/commands/)
- 💡 [Examples](https://zayan-mohamed.github.io/vaultix/examples/)
- 🔐 [Security Model](https://zayan-mohamed.github.io/vaultix/security/)
- 🏗️ [Architecture](https://zayan-mohamed.github.io/vaultix/architecture/)

---

## 🔧 Development

### Building from Source

```bash
# Clone the repository
git clone https://github.com/Zayan-Mohamed/vaultix.git
cd vaultix

# Build
go build -o vaultix

# Run tests
go test ./...

# Run linters
go vet ./...
```

### Project Architecture

```
vaultix/
├── internal/
│   ├── crypto/    # Cryptographic operations (Argon2id, AES-GCM)
│   ├── storage/   # File system operations
│   ├── vault/     # Business logic layer
│   └── cli/       # Command-line interface
├── docs/          # MkDocs documentation
└── main.go        # Entry point
```

---

## 🤝 Contributing

Contributions are welcome! Please read our [Contributing Guide](https://zayan-mohamed.github.io/vaultix/contributing/) first.

### Quick Start for Contributors

```bash
# Fork and clone
git clone https://github.com/YOUR_USERNAME/vaultix.git
cd vaultix

# Build
go build -o vaultix

# Run tests
go test ./...

# Lint code
go vet ./...
```

---

## 📄 License

MIT License - See [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- Built with [Go](https://golang.org/)
- Uses [Argon2](https://github.com/P-H-C/phc-winner-argon2) for key derivation
- Inspired by the need for simple, secure file encryption

---

## ⚠️ Disclaimer

This software is provided as-is, without any warranties. While vaultix uses well-established cryptographic libraries and follows security best practices, it has not undergone formal security auditing. Use at your own risk.

**Remember**: Your vault's security depends entirely on your password strength and operational security practices.

---

<div align="center">

![GitHub code size](https://img.shields.io/github/languages/code-size/Zayan-Mohamed/vaultix?style=flat-square&label=Code%20Size)
![GitHub go.mod version](https://img.shields.io/github/go-mod/go-version/Zayan-Mohamed/vaultix?style=flat-square&label=Go%20Version)
![Lines of code](https://img.shields.io/endpoint?url=https://ghloc.vercel.app/api/Zayan-Mohamed/vaultix/badge?filter=.go$&style=flat-square&label=Lines%20of%20Code)
![Total Files](https://img.shields.io/github/directory-file-count/Zayan-Mohamed/vaultix?style=flat-square&label=Total%20Files)
![GitHub repo size](https://img.shields.io/github/repo-size/Zayan-Mohamed/vaultix?style=flat-square&label=Repo%20Size)
![Last Commit](https://img.shields.io/github/last-commit/Zayan-Mohamed/vaultix?style=flat-square&label=Last%20Commit)

</div>



<p align="center">Made with ❤️ for security-conscious developers</p>
<p align="center">⭐ Star this repo if you find it useful!</p>
