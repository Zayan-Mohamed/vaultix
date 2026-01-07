# Installation

## Prerequisites

Before installing Vaultix, ensure you have:

- **Go 1.21 or higher** (only for building from source)
- **Git** (for cloning the repository)

## Quick Install (Recommended)

=== "Linux & macOS"

    ```bash
    git clone https://github.com/Zayan-Mohamed/vaultix.git
    cd vaultix
    ./install.sh
    ```

    The script will:

    - Build the binary
    - Install to `/usr/local/bin/vaultix`
    - Make it executable
    - Verify the installation

=== "Windows"

    ```powershell
    git clone https://github.com/Zayan-Mohamed/vaultix.git
    cd vaultix
    .\install.ps1
    ```

    The script will:

    - Build the binary
    - Install to `%USERPROFILE%\bin\vaultix.exe`
    - Add to PATH
    - Verify the installation

## Manual Installation

### Build from Source

```bash
git clone https://github.com/Zayan-Mohamed/vaultix.git
cd vaultix
go build -o vaultix
```

### Install Binary

=== "Linux & macOS"

    ```bash
    # Copy to system directory
    sudo cp vaultix /usr/local/bin/

    # Make executable
    sudo chmod +x /usr/local/bin/vaultix

    # Verify
    vaultix help
    ```

=== "Windows"

    ```powershell
    # Create bin directory if it doesn't exist
    New-Item -ItemType Directory -Force -Path $env:USERPROFILE\bin

    # Copy binary
    Copy-Item vaultix.exe $env:USERPROFILE\bin\

    # Add to PATH (if not already)
    $oldPath = [Environment]::GetEnvironmentVariable('Path', 'User')
    $newPath = "$oldPath;$env:USERPROFILE\bin"
    [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')

    # Verify (restart terminal first)
    vaultix help
    ```

## Verify Installation

After installation, verify that Vaultix is working:

```bash
vaultix help
```

You should see the help output with all available commands.

## Update Vaultix

To update to the latest version:

```bash
cd vaultix
git pull origin main
go build -o vaultix

# Linux/macOS
sudo cp vaultix /usr/local/bin/

# Windows
Copy-Item vaultix.exe $env:USERPROFILE\bin\
```

## Uninstall

=== "Linux & macOS"

    ```bash
    sudo rm /usr/local/bin/vaultix
    ```

=== "Windows"

    ```powershell
    Remove-Item $env:USERPROFILE\bin\vaultix.exe
    ```

## Troubleshooting

### Command Not Found

If you get "command not found" after installation:

**Linux/macOS**: Make sure `/usr/local/bin` is in your PATH:

```bash
echo $PATH | grep -q "/usr/local/bin" && echo "✓ In PATH" || echo "✗ Not in PATH"
```

**Windows**: Restart your terminal/PowerShell after running the install script.

### Permission Denied

**Linux/macOS**: The binary needs execute permissions:

```bash
sudo chmod +x /usr/local/bin/vaultix
```

**Windows**: Run PowerShell as Administrator when installing.

### Go Not Found

If you don't have Go installed:

- **Linux**: `sudo apt install golang-go` or `sudo yum install golang`
- **macOS**: `brew install go`
- **Windows**: Download from [golang.org](https://golang.org/dl/)

## Next Steps

Now that Vaultix is installed, check out the [Quick Start Guide](quickstart.md) to create your first vault!
