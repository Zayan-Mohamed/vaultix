# Installation Guide

## Prerequisites

- **Go 1.19 or later** installed on your system
- **Git** to clone the repository
- **Internet connection** for downloading dependencies

## Platform-Specific Installation

### Linux

#### Option 1: Automated Installation (Recommended)

```bash
git clone https://github.com/Zayan-Mohamed/vaultix.git
cd vaultix
./install.sh
```

The script will:

- Build the binary
- Install to `/usr/local/bin/vaultix`
- Request sudo password if needed

#### Option 2: Manual Installation

```bash
git clone https://github.com/Zayan-Mohamed/vaultix.git
cd vaultix
go build -o vaultix
sudo mv vaultix /usr/local/bin/
sudo chmod +x /usr/local/bin/vaultix
```

#### Verify Installation

```bash
vaultix --help
```

---

### macOS

#### Option 1: Automated Installation (Recommended)

```bash
git clone https://github.com/Zayan-Mohamed/vaultix.git
cd vaultix
./install.sh
```

The script will:

- Build the binary
- Install to `/usr/local/bin/vaultix`
- Request sudo password if needed

#### Option 2: Manual Installation

```bash
git clone https://github.com/Zayan-Mohamed/vaultix.git
cd vaultix
go build -o vaultix
sudo mv vaultix /usr/local/bin/
sudo chmod +x /usr/local/bin/vaultix
```

#### Verify Installation

```bash
vaultix --help
```

---

### Windows

#### Option 1: Automated Installation (Recommended)

Open PowerShell and run:

```powershell
git clone https://github.com/Zayan-Mohamed/vaultix.git
cd vaultix
.\install.ps1
```

The script will:

- Build `vaultix.exe`
- Install to `%USERPROFILE%\bin\vaultix.exe`
- Add the directory to your PATH
- Prompt to restart terminal

**Important:** Restart your terminal after installation for PATH changes to take effect.

#### Option 2: Manual Installation

```powershell
git clone https://github.com/Zayan-Mohamed/vaultix.git
cd vaultix
go build -o vaultix.exe

# Create bin directory if it doesn't exist
New-Item -ItemType Directory -Path $env:USERPROFILE\bin -Force

# Copy binary
Copy-Item vaultix.exe -Destination $env:USERPROFILE\bin\vaultix.exe

# Add to PATH manually via System Settings
```

#### Verify Installation

Open a **new** terminal window:

```powershell
vaultix --help
```

---

## Installation Locations

| Platform | Default Location                | Requires Sudo/Admin |
| -------- | ------------------------------- | ------------------- |
| Linux    | `/usr/local/bin/vaultix`        | Yes                 |
| macOS    | `/usr/local/bin/vaultix`        | Yes                 |
| Windows  | `%USERPROFILE%\bin\vaultix.exe` | No                  |

---

## Troubleshooting

### "Go is not installed"

**Solution:** Install Go from https://golang.org/dl/

Verify installation:

```bash
go version
```

### "Command not found: vaultix" (after installation)

**Linux/macOS:**

- Verify `/usr/local/bin` is in your PATH:
  ```bash
  echo $PATH | grep /usr/local/bin
  ```
- If not, add it to your shell profile:
  ```bash
  echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.bashrc
  source ~/.bashrc
  ```

**Windows:**

- Restart your terminal/PowerShell
- Check if `%USERPROFILE%\bin` is in PATH:
  ```powershell
  $env:Path -split ';' | Select-String "bin"
  ```
- If not, add manually via System Settings → Environment Variables

### "Permission denied"

**Linux/macOS:**

- Make sure you have sudo privileges
- Or install to a local directory:
  ```bash
  mkdir -p ~/bin
  mv vaultix ~/bin/
  export PATH="$HOME/bin:$PATH"
  ```

### Build errors

**Solution:**

1. Update Go to the latest version
2. Ensure you're in the vaultix directory
3. Run `go mod tidy` and try again

---

## Alternative Installation Methods

### Using Make

If you have `make` installed:

```bash
git clone https://github.com/Zayan-Mohamed/vaultix.git
cd vaultix
make install
```

### From Release Binaries (when available)

Download pre-built binaries from GitHub Releases:

1. Go to https://github.com/Zayan-Mohamed/vaultix/releases
2. Download the binary for your platform
3. Extract and move to your PATH

---

## Uninstallation

### Linux/macOS

```bash
sudo rm /usr/local/bin/vaultix
```

### Windows

```powershell
Remove-Item $env:USERPROFILE\bin\vaultix.exe
```

Optionally remove from PATH via System Settings.

---

## Updating

To update to the latest version:

```bash
cd vaultix
git pull origin main
./install.sh  # or .\install.ps1 on Windows
```

This will rebuild and reinstall the latest version.

---

## Next Steps

After installation, verify everything works:

```bash
# Create a test vault
vaultix init ./test_vault

# Try adding a file
echo "test" > test.txt
vaultix add ./test_vault test.txt

# List files
vaultix list ./test_vault

# Clean up
rm -rf ./test_vault test.txt
```

See [README.md](README.md) for full documentation and usage examples.
