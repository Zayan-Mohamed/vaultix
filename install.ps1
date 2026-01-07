# vaultix installation script for Windows
# Run this script in PowerShell with: .\install.ps1

$ErrorActionPreference = "Stop"

Write-Host "======================================" -ForegroundColor Cyan
Write-Host "vaultix Installation Script (Windows)" -ForegroundColor Cyan
Write-Host "======================================" -ForegroundColor Cyan
Write-Host ""

# Check if Go is installed
try {
    $goVersion = go version
    Write-Host "Go version: $goVersion" -ForegroundColor Green
    Write-Host ""
} catch {
    Write-Host "Error: Go is not installed" -ForegroundColor Red
    Write-Host "Please install Go from https://golang.org/dl/" -ForegroundColor Yellow
    exit 1
}

# Build the binary
Write-Host "Building vaultix.exe..." -ForegroundColor Cyan
try {
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    go build -o vaultix.exe -ldflags="-s -w" .
    Write-Host "✓ Build successful" -ForegroundColor Green
} catch {
    Write-Host "✗ Build failed" -ForegroundColor Red
    exit 1
}
Write-Host ""

# Get the binary size
$binarySize = (Get-Item vaultix.exe).Length / 1MB
Write-Host "Binary size: $([math]::Round($binarySize, 2)) MB" -ForegroundColor Cyan
Write-Host ""

# Determine installation directory
# Option 1: Install to user's local bin (recommended, no admin required)
$userBin = "$env:USERPROFILE\bin"
$installDir = $userBin

# Create user bin directory if it doesn't exist
if (-not (Test-Path $installDir)) {
    Write-Host "Creating directory: $installDir" -ForegroundColor Cyan
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null
}

# Install the binary
Write-Host "Installing vaultix.exe to $installDir..." -ForegroundColor Cyan
try {
    Copy-Item vaultix.exe -Destination "$installDir\vaultix.exe" -Force
    Write-Host "✓ Installation successful" -ForegroundColor Green
} catch {
    Write-Host "✗ Installation failed: $_" -ForegroundColor Red
    exit 1
}

# Check if the directory is in PATH
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notlike "*$installDir*") {
    Write-Host ""
    Write-Host "Adding $installDir to PATH..." -ForegroundColor Cyan
    try {
        $newPath = "$userPath;$installDir"
        [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
        Write-Host "✓ PATH updated" -ForegroundColor Green
        Write-Host ""
        Write-Host "IMPORTANT: Please restart your terminal for PATH changes to take effect" -ForegroundColor Yellow
    } catch {
        Write-Host "⚠ Could not update PATH automatically" -ForegroundColor Yellow
        Write-Host "Please add $installDir to your PATH manually" -ForegroundColor Yellow
    }
} else {
    Write-Host "✓ $installDir is already in PATH" -ForegroundColor Green
}

# Clean up build artifact in current directory
if (Test-Path "vaultix.exe" -and (Resolve-Path "vaultix.exe").Path -ne (Resolve-Path "$installDir\vaultix.exe").Path) {
    Remove-Item vaultix.exe -Force
}

Write-Host ""
Write-Host "======================================" -ForegroundColor Cyan
Write-Host "Installation Complete!" -ForegroundColor Green
Write-Host "======================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "vaultix is now installed at: $installDir\vaultix.exe" -ForegroundColor Cyan
Write-Host ""

if ($userPath -notlike "*$installDir*") {
    Write-Host "⚠ Please restart your terminal, then verify installation:" -ForegroundColor Yellow
} else {
    Write-Host "Verify installation (in a new terminal):" -ForegroundColor Cyan
}

Write-Host "  vaultix --help" -ForegroundColor White
Write-Host ""
Write-Host "Get started:" -ForegroundColor Cyan
Write-Host "  vaultix init .\my_vault" -ForegroundColor White
Write-Host ""
Write-Host "For more information, see README.md" -ForegroundColor Cyan
Write-Host ""

# Alternative installation note
Write-Host "Alternative: Install to Program Files (requires admin)" -ForegroundColor DarkGray
Write-Host "  Run PowerShell as Administrator and change `$installDir to:" -ForegroundColor DarkGray
Write-Host "  `$env:ProgramFiles\vaultix" -ForegroundColor DarkGray
