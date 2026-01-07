# Examples

Real-world usage examples for common scenarios.

## Personal Documents

### Protecting Family Photos

```bash
# Create vault for photos
mkdir ~/photos_vault
cd ~/photos_vault

# Add your photos
cp ~/Pictures/family/*.jpg .

# Initialize vault
vaultix init
# Enter strong password

# Photos are now encrypted
ls -la
# drwx------  .vaultix/

# When you need them
vaultix extract vacation_2024.jpg
# Opens or extracts to working directory
```

### Tax Documents

```bash
# Annual tax folder
mkdir ~/taxes_2024
cd ~/taxes_2024

# Collect documents
mv ~/Downloads/w2.pdf .
mv ~/Downloads/1099.pdf .
mv ~/Documents/receipts.zip .

# Encrypt everything
vaultix init

# Later, during tax season
vaultix list
# w2.pdf
# 1099.pdf
# receipts.zip

vaultix extract w2
# Work with file

# Done? Remove extract
rm w2.pdf
```

## Professional Use Cases

### Source Code Repository

```bash
# Proprietary code
mkdir ~/work/private_client_project
cd ~/work/private_client_project

# Initialize git repo
git init

# Add code
echo "const API_KEY = 'secret123'" > config.js
echo "# Client Project" > README.md

# Encrypt repository
vaultix init

# Share encrypted version safely
tar czf project_encrypted.tar.gz .
# Can upload to untrusted storage

# Recipient extracts and decrypts
tar xzf project_encrypted.tar.gz
cd private_client_project
vaultix extract  # Gets all files
```

### API Keys and Credentials

```bash
# Credentials vault
mkdir ~/credentials
cd ~/credentials

# Create key files
cat > aws_keys.txt << EOF
AWS_ACCESS_KEY_ID=AKIA...
AWS_SECRET_ACCESS_KEY=secret...
EOF

cat > api_tokens.txt << EOF
GITHUB_TOKEN=ghp_...
STRIPE_KEY=sk_live_...
EOF

# Encrypt
vaultix init

# When deploying
vaultix extract aws_keys.txt
source aws_keys.txt
# Deploy application
rm aws_keys.txt  # Clean up
```

### Client Contracts

```bash
# Legal documents
mkdir ~/contracts/client_xyz
cd ~/contracts/client_xyz

# Add signed contracts
cp ~/Downloads/signed_nda.pdf .
cp ~/Downloads/signed_contract.pdf .
cp ~/Downloads/sow.pdf .

# Encrypt
vaultix init

# Reference later
vaultix list
vaultix extract signed_contract.pdf
```

## Development Workflows

### Environment Files

```bash
# Development project
cd ~/projects/webapp
cat > .env << EOF
DATABASE_URL=postgresql://...
SECRET_KEY=very_secret_key
API_TOKEN=prod_token
EOF

# Encrypt environment file
vaultix add .env

# In .gitignore
echo ".env" >> .gitignore

# Other developers
vaultix extract .env
# Enter password (shared securely)
npm run dev
```

### Configuration Management

```bash
# Production configs
mkdir ~/configs/production
cd ~/configs/production

# Create configs
cat > database.yml << EOF
host: prod-db.example.com
user: admin
password: super_secret
EOF

cat > redis.yml << EOF
host: redis.example.com
password: another_secret
EOF

# Encrypt all
vaultix init

# Deploy script
#!/bin/bash
cd ~/configs/production
vaultix extract database.yml
scp database.yml server:/etc/app/
rm database.yml
```

### Private Keys

```bash
# SSH keys vault
mkdir ~/.ssh/vaultix_keys
cd ~/.ssh/vaultix_keys

# Add private keys
cp ~/.ssh/client_deploy_key .
cp ~/.ssh/production_key .

# Encrypt
vaultix init

# When needed
vaultix extract client_deploy_key
chmod 600 client_deploy_key
ssh -i client_deploy_key user@server
rm client_deploy_key
```

## Research and Writing

### Research Notes

```bash
# Research project
mkdir ~/research/project_x
cd ~/research/project_x

# Create notes
vim notes_2024-01-15.md
vim draft_paper.docx
vim references.bib

# Encrypt
vaultix init

# Daily workflow
vaultix extract notes  # Fuzzy match
# Edit notes
vim notes_2024-01-15.md
# Update vault
vaultix add notes_2024-01-15.md
```

### Journal

```bash
# Private journal
mkdir ~/journal
cd ~/journal

# Daily entries
date=$(date +%Y-%m-%d)
vim "entry_$date.md"

# Encrypt each entry
vaultix add "entry_$date.md"

# Read old entries
vaultix list | grep 2024-01
vaultix extract entry_2024-01-15.md
```

## Media Management

### Video Projects

```bash
# Video editing project
mkdir ~/videos/client_commercial
cd ~/videos/client_commercial

# Raw footage (large files)
# Tip: Compress first
tar czf raw_footage.tar.gz footage/
vaultix add raw_footage.tar.gz
rm -rf footage/

# When editing
vaultix extract raw_footage
tar xzf raw_footage.tar.gz
# Edit in your NLE
```

### Music Production

```bash
# Unreleased tracks
mkdir ~/music/unreleased
cd ~/music/unreleased

# Add stems and projects
cp -r ~/Production/new_song_project .

# Encrypt
vaultix init

# Collaboration
vaultix extract new_song_project
# Work in DAW
vaultix add new_song_project  # Update
```

## System Administration

### Server Backups

```bash
# Backup sensitive configs
mkdir ~/backups/server_configs_2024-01-15
cd ~/backups/server_configs_2024-01-15

# Pull configs
scp server:/etc/nginx/ssl/* .
scp server:/etc/app/secrets.yml .

# Encrypt
vaultix init

# Upload to cloud storage (safe when encrypted)
rclone sync . remote:encrypted_backups/2024-01-15/
```

### Password Database

```bash
# Offline password backup
mkdir ~/password_vault
cd ~/password_vault

# Export from password manager
# KeePass: File → Export
# 1Password: Export all items
mv ~/Downloads/passwords_export.csv .

# Encrypt immediately
vaultix init

# Delete export
rm ~/Downloads/passwords_export.csv

# Store backup
cp -r ~/password_vault /mnt/backup/
```

## Financial Records

### Cryptocurrency Wallets

```bash
# Wallet backups
mkdir ~/crypto/wallet_backups
cd ~/crypto/wallet_backups

# Export wallet seeds/keys
cat > btc_wallet.txt << EOF
Seed: word1 word2 ... word12
Private Key: 5K...
EOF

# Encrypt immediately
vaultix init

# Multiple secure locations
cp -r . /mnt/backup1/crypto/
cp -r . /mnt/backup2/crypto/
```

### Investment Records

```bash
# Brokerage statements
mkdir ~/finance/investments_2024
cd ~/finance/investments_2024

# Add statements
mv ~/Downloads/*_statement.pdf .

# Encrypt
vaultix init

# Tax time
vaultix extract
# Work with accountant
vaultix clear  # Clean up after tax season
```

## Academic Work

### Thesis/Dissertation

```bash
# PhD thesis
mkdir ~/phd/thesis
cd ~/phd/thesis

# Working drafts
vim chapter1.tex
vim chapter2.tex
vim references.bib

# Daily backup routine
vaultix init
vaultix add *.tex *.bib

# On different computer
vaultix list
vaultix extract chapter1
```

### Exam Materials

```bash
# Exam questions (instructors)
mkdir ~/teaching/cs101/exams
cd ~/teaching/cs101/exams

# Create exam
vim midterm_2024.pdf

# Encrypt until exam day
vaultix init

# Exam day
vaultix extract midterm
# Print and distribute
```

## Healthcare

### Medical Records

```bash
# Personal health
mkdir ~/health/medical_records
cd ~/health/medical_records

# Scan documents
mv ~/Downloads/lab_results_2024.pdf .
mv ~/Downloads/prescription.pdf .

# Encrypt
vaultix init

# Doctor visit
vaultix extract lab_results
# Show to doctor
rm lab_results_2024.pdf
```

## Business Operations

### Payroll

```bash
# Employee payroll
mkdir ~/business/payroll/2024-01
cd ~/business/payroll/2024-01

# Generate payroll
vim salaries.csv
vim tax_withholding.xlsx

# Encrypt immediately
vaultix init

# Processing day
vaultix extract
# Process payments
vaultix clear  # Remove after processing
```

### Customer Database

```bash
# Client information
mkdir ~/business/clients
cd ~/business/clients

# Export from CRM
mv ~/Downloads/clients_export.csv .

# Encrypt
vaultix init

# Marketing campaign
vaultix extract clients_export
# Import to email tool
rm clients_export.csv
```

## Automation Examples

### Batch Processing

```bash
#!/bin/bash
# encrypt_all_projects.sh

for project in ~/projects/*/; do
    cd "$project"
    if [ ! -d ".vaultix" ]; then
        echo "Encrypting $project"
        vaultix init < password.txt
    fi
done
```

### Backup Script

```bash
#!/bin/bash
# daily_vault_backup.sh

VAULT_DIR=~/important_vault
BACKUP_DIR=/mnt/backup/vaults
DATE=$(date +%Y%m%d)

# Copy vault
cp -r "$VAULT_DIR" "$BACKUP_DIR/vault_$DATE"

# Keep only last 7 days
find "$BACKUP_DIR" -type d -name "vault_*" -mtime +7 -exec rm -rf {} \;

echo "Backup completed: vault_$DATE"
```

### Extract and Process

```bash
#!/bin/bash
# process_vault_files.sh

cd ~/data_vault

# Extract all data files
vaultix extract

# Process
python3 analyze_data.py *.csv > report.txt

# Clean up
rm *.csv

# Store report
vaultix add report.txt
```

## Integration Examples

### With Git

```bash
# Track encrypted vault in git
cd ~/project
vaultix init

git add .vaultix/
git commit -m "Add encrypted data"
git push

# Teammate clones
git clone repo
cd project
vaultix list  # Enter shared password
```

### With Cloud Storage

```bash
# Sync encrypted vault
cd ~/vault
vaultix init

# Setup rclone
rclone config

# Sync
rclone sync . remote:vault/

# On other computer
rclone sync remote:vault/ ~/vault
vaultix extract
```

### With Docker

```dockerfile
# Dockerfile with vaultix
FROM golang:1.21
COPY vaultix /usr/local/bin/
COPY vault/ /vault/
WORKDIR /vault
CMD ["vaultix", "list"]
```

## Tips and Tricks

### Quick Access Alias

```bash
# Add to ~/.bashrc or ~/.zshrc
alias vault='cd ~/vault && vaultix'
alias vlist='cd ~/vault && vaultix list'
alias vadd='cd ~/vault && vaultix add'
alias vget='cd ~/vault && vaultix extract'
```

### Password from Environment

```bash
# DON'T DO THIS IN PRODUCTION
export VAULT_PASSWORD="my_password"
echo "$VAULT_PASSWORD" | vaultix list
```

### Preview Without Extracting

```bash
# For text files
vaultix extract readme.txt
cat readme.txt
rm readme.txt

# Or
vaultix extract readme.txt | less
```

### Compress Before Encrypting

```bash
# Large directories
tar czf data.tar.gz data/
vaultix add data.tar.gz
rm data.tar.gz

# Extract later
vaultix extract data.tar.gz
tar xzf data.tar.gz
```

## Common Patterns

### Temporary Extraction

```bash
# Extract to temp, use, delete
cd ~/vault
TEMP=$(mktemp -d)
vaultix extract sensitive.pdf
mv sensitive.pdf "$TEMP/"
# Work with file
cd "$TEMP"
# Done
rm -rf "$TEMP"
```

### Verify Before Clearing

```bash
# Before permanent deletion
vaultix list > files_list.txt
vaultix extract  # Get all files
# Verify extracted files
ls -l
# If satisfied
vaultix clear
```

### Version Control for Vaults

```bash
# Keep vault versions
mkdir ~/vault_versions
cd ~/my_vault

# Before major changes
tar czf ~/vault_versions/vault_$(date +%Y%m%d).tar.gz .

# Make changes
vaultix add new_files.zip

# If something goes wrong
tar xzf ~/vault_versions/vault_20240115.tar.gz
```

---

These examples demonstrate vaultix in various real-world scenarios. Adapt them to your specific needs while following security best practices.
