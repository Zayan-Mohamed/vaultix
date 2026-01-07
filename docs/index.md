# Welcome to Vaultix

<p align="center">
  <strong>A cross-platform command-line tool for managing password-protected encrypted folders</strong>
</p>

## What is Vaultix?

Vaultix is a secure, lightweight CLI tool that encrypts files in place using military-grade cryptography. No cloud, no services, no complexity—just strong encryption for your sensitive files.

## Key Features

🔒 **Strong Encryption**: AES-256-GCM with Argon2id key derivation  
🚀 **Zero Dependencies**: Single static binary, no runtime requirements  
💻 **Cross-Platform**: Linux, macOS, and Windows support  
🎯 **Simple UX**: Intuitive commands with smart defaults  
🔐 **No Password Storage**: Passwords exist only in memory  
📦 **Portable**: Encrypted vaults work across all platforms

## Quick Example

```bash
# Navigate to sensitive files
cd ~/my_secrets

# Initialize vault (auto-encrypts all files)
vaultix init
# Enter password: ****
# ✓ Vault initialized
# ✓ All files encrypted

# List files
vaultix list

# Extract a file
vaultix extract passwords

# Drop a file (extract and remove from vault)
vaultix drop api_keys
```

## Why Vaultix?

### Simple

One binary, no installation hassles, no background processes. Just run it when you need it.

### Secure

Uses proven cryptographic algorithms (Argon2id + AES-256-GCM) with secure defaults. No custom crypto, no shortcuts.

### Portable

Create a vault on Linux, access it on Windows, move it to macOS. It just works.

### Transparent

Open source, readable codebase, no hidden behavior. See exactly what it does.

## What Vaultix is NOT

Vaultix is **not**:

- ❌ A password manager
- ❌ A cloud sync tool
- ❌ A keychain replacement
- ❌ An OS-level security boundary
- ❌ A background service

It's a simple tool that does one thing well: encrypt your files.

## Getting Started

Ready to secure your files? Start with the [Installation Guide](installation.md) or jump straight to the [Quick Start](quickstart.md).

## Documentation Structure

- **Getting Started**: Installation, quick start, and basic usage
- **User Guide**: Detailed command reference and examples
- **Technical**: Architecture, security model, and cryptography details
- **Reference**: API documentation, FAQ, and troubleshooting
- **Contributing**: Development guide and contribution guidelines

## Community and Support

- 🐛 [Report Issues](https://github.com/Zayan-Mohamed/vaultix/issues)
- 💬 [Discussions](https://github.com/Zayan-Mohamed/vaultix/discussions)
- 📖 [Documentation](https://zayan-mohamed.github.io/vaultix)
- ⭐ [Star on GitHub](https://github.com/Zayan-Mohamed/vaultix)

## License

Vaultix is open source software licensed under the [MIT License](https://github.com/Zayan-Mohamed/vaultix/blob/main/LICENSE).
