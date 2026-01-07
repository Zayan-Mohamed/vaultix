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
✅ **Fuzzy File Matching** - No need to type exact filenames  
✅ **Default to Current Directory** - Less typing, more doing  
✅ **Extract or Drop** - Extract files while keeping in vault, or drop them out  
✅ **Secure Deletion** - Original files are overwritten before deletion  
✅ **Hidden Metadata** - Even filenames are encrypted  
✅ **No Background Process** - Runs only when you invoke it

## Security Model

### Cryptography

vaultix uses industry-standard cryptographic primitives:

- **Key Derivation**: Argon2id with 64MB memory, 1 iteration, 4 threads
- **Encryption**: AES-256-GCM (authenticated encryption)
- **Randomness**: Go's `crypto/rand` package

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

- **Password-only security**: Your vault is only as secure as your password
- **No password recovery**: Forget your password = lose your data permanently
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

sudo mv vaultix /usr/local/bin/  # Linux/macOS
```

---

## cd ~/my_secrets

# Initialize vault (encrypts all files automatically)
vaultix init
# Enter password: ****
# Confirm password: ****
# ✓ Vault initialized
# ✓ All files encrypted
# ✓ Original files securely deleted

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
| `list`           | List encrypted files                   | `vaultix list`           |
| `extract [file]` | Extract file(s), keeps in vault        | `vaultix extract`        |
| `drop [file]`    | Extract and remove from vault          | `vaultix drop secret`    |

| `🏗️ Architecture

### Pro [path]`    | List encrypted files                   | `vaultix list`           |
| `extract [file]` | Extract file(s), keeps in vault        | `vaultix extract`        |
| `drop <file>`    | Extract and remove from vault          | `vaultix drop secret`    |
| `remove <file>`  | Remove file from vault (no extract)    | `vaultix remove old.txt` |
| `clear [path]`   | Remove all files from vault            | `vaultix clear`          |

> 💡 **Pro Tip**: Most commands default to current directory, so you rarely need to specify paths!

---

## 🔒 Security
│   ├── crypto/       # Cryptographic operations (Argon2id, AES-GCM)
│   ├── storage/      # File system operations
│   ├── vault/        # Business logic layer
│   └── cli/          # Command-line interface
├── .github/          # GitHub configuration
├── docs/             # MkDocs documentation
├── main.go           # Entry point
└── install.sh        # Installation script
```

### How It Wfile>`| Remove file without extracting |`vaultix remove old.txt` |

| `clear` | Remove all files (with confirmation) | `vaultix clear` |

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

For complete usage examples, see [EXAMPLES.md](EXAMPLES.md)

```bash
vaultix remove ./my_secrets document.pdf
```

This permanently deletes the encrypted file from the vault.

## How it works

When you initialize a vault at a path (e.g., `./my_secrets`), vaultix creates a hidden `.vaultix/` directory inside:

```
my_secrets/
└── .vaultix/
    ├── meta          # Encrypted metadata (filenames, sizes, timestamps)
    ├── salt          # Random salt for key derivation
    └── objects/
        ├── 3f9a2c1d.enc
        └── 91bd77aa.enc
```

### Security Details

1. **No passwords stored**: Your password exists only in memory during operations
2. **Encrypted metadata**: Even filenames are encrypted
3. **Obfuscated object names**: Encrypted files have random IDs
4. **Salt per vault**: Each vault has a unique random salt
5. **Authentication**: AES-GCM provides both encryption and integrity verification

### Password Verification

Password correctness is verified by successful decryption of the metadata. There are no stored password hashes. This means:

- Incorrect password = decryption failure
- No way to test passwords without attempting decryption
- No way to recover from a forgotten password

## Best Practices

### Password Selection

Use a strong, unique password:

- At least 16 characters
- Mix of letters, numbers, and symbols
- Not used anywhere else
- Not easily guessable

Consider using a password manager to generate and store your vault password.

### Backup Strategy

- The entire vault directory (including `.vaultix/`) must be backed up
- Test your backups by extracting files from backup copies
- Encrypted vaults are safe to backup to cloud storage
- Losing `.vaultix/` = permanent data loss

### File Management

- Original files are **not** automatically deleted when added to vault
- Delete originals manually if needed (consider secure deletion tools)
- Keep temporary extractions out of the vault directory
- Don't extract sensitive files to public/shared directories

### Operational Security

- Don't enter passwords where they might be logged
- Don't use vaultix over remote connections without encryption
- Close your terminal after vault operations
- Consider using full-disk encryption alongside vaultix
- Be aware of swap files and hibernation dumps

## Platform Support

vaultix works identically on:
📖 Documentation

Full documentation is available at: [https://zayan-mohamed.github.io/vaultix](https://zayan-mohamed.github.io/vaultix)

## 🔧 Development

### Building from Source

⚠️ Disclaimer

This software is provided as-is, without any warranties. While vaultix uses well-established cryptographic libraries and follows security best practices, it has not undergone formal security auditing. Use at your own risk.

**Remember**: The security of your vault depends entirely on your password strength and operational security practices.

## 🙏 Acknowledgments

- Built with [Go](https://golang.org/)
- Uses [Argon2](https://github.com/P-H-C/phc-winner-argon2) for key derivation
- Inspired by the need for simple, secure file encryption

## 📊 Project Stats

![GitHub code size](https://img.shields.io/github/languages/code-size/Zayan-Mohamed/vaultix?style=flat-square)
![GitHub go.mod version](https://img.shields.io/github/go-mod/go-version/Zayan-Mohamed/vaultix?style=flat-square)
![Lines of code](https://img.shields.io/tokei/lines/github/Zayan-Mohamed/vaultix?style=flat-square)

---

<p align="center">Made with ❤️ for security-conscious developers</p>
```

---

## 🔒 Security

### Cryptography

- **Key Derivation**: Argon2id with 64MB memory, 3 iterations, 4 threads
- **Encryption**: AES-256-GCM (authenticated encryption)
- **Randomness**: Go's `crypto/rand` package

### Threat Model

**Protects Against:**
- ✅ Unauthorized access to files at rest
- ✅ Accidental exposure of file contents
- ✅ Cloud storage providers reading your data

**Does NOT Protect Against:**
- ❌ Weak passwords (use 16+ characters!)
- ❌ Malware or keyloggers on your system
- ❌ Physical access while computer is unlocked
- ❌ Coercion or legal compulsion

### Important Limitations

- **Password-only security**: Your vault is only as secure as your password
- **No password recovery**: Forget your password = lose your data permanently
- **No automatic backups**: You are responsible for backing up your vaults

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

# Run linters
go vet ./...
staticcheck ./...
```

---

## 📄 License

MIT License - See [LICENSE](LICENSE) file for details.

---

## ⚠️ Disclaimer

This software is provided as-is, without any warranties. While vaultix uses well-established cryptographic libraries, it has not undergone formal security auditing. Use at your own risk.

**Remember**: Your vault's security depends entirely on your password strength and operational security practices.

---

<p align="center">Made with ❤️ for security-conscious developers</p>
<p align="center">⭐ Star this repo if you find it useful!</p>
