# YubiKey Employee Setup Script

This script automates the YubiKey setup process for EnableIT employees on Linux and macOS systems, following the official documentation at [EnableIT YubiKey Setup Guide](https://gitea.obmondo.com/EnableIT/wiki/src/branch/master/internal/yubikey-employee-setup.md).

## 🎯 Purpose

Sets up YubiKey for:
- GPG signing and encryption
- SSH authentication  
- Git commit signing
- Integration with EnableIT infrastructure

## 📋 Prerequisites

Before running this script, ensure you have:

1. **YubiKey device** - Physical hardware token
2. **YubiKey PIN** - Contact Klavs or Ashish if you don't have it
3. **GPG public key** - Available at https://gitea.obmondo.com/EnableIT/pass
4. **Administrative access** - Script may need sudo for package installation

## 🚀 Usage

### Quick Start

```bash
# Make the script executable (if not already)
chmod +x yubikey_setup.sh

# Run the setup script
./yubikey_setup.sh
```

### Step-by-Step Process

1. **OS Detection**: Script automatically detects and validates your operating system
2. **PIN Verification**: Prompts to confirm you have your YubiKey PIN
3. **Key Path Input**: Enter path to your GPG public key file
4. **Automatic Setup**: Script handles all configuration automatically
5. **Testing**: Validates the setup with encryption/decryption tests
6. **Final Instructions**: Provides next steps for Gitea integration

## 🖥️ Supported Systems

- **Linux** (Ubuntu/Debian with apt-get)
- **macOS** (with Homebrew)

## 📁 What the Script Does

### Prerequisites Installation
- **Linux**: `gnupg2`, `gnupg-agent`, `pinentry-curses`, `scdaemon`, `pcscd`, `libusb-1.0-0-dev`, `pinentry-gnome3`
- **macOS**: `gnupg`, `yubikey-personalization`, `pinentry-mac`

### Configuration Files Created/Modified
- `~/.gnupg/gpg.conf` - GPG configuration
- `~/.gnupg/dirmngr.conf` - Key server configuration  
- `~/.gnupg/gpg-agent.conf` - GPG agent settings
- `~/.bashrc` (Linux) or `~/.bash_profile` (macOS) - Shell environment

### GPG Operations
- Import your public key
- Set key trust to ultimate (level 5)
- Configure GPG agent with SSH support
- Test encryption/decryption functionality

### Git Configuration
- Set signing key for commits
- Enable automatic commit signing

## 🔧 Error Handling

The script includes comprehensive error handling:

- **UUID tracking** - Each run generates a unique error ID for troubleshooting
- **Validation checks** - Verifies each step before proceeding
- **Graceful failures** - Stops on any critical error with clear messages
- **Common issue fixes** - Handles known problems automatically

### Common Issues and Solutions

| Issue | Solution |
|-------|----------|
| "Screen or window too small" | Script sets `GPG_TTY` and `pinentry-gnome3` |
| "No secret key" | Uses alternative pinentry configuration |
| "Key not changed" | Skips import if key already exists |
| YubiKey not detected | Prompts to replug device |

## 📊 Script Flow

```
1. OS Detection (Linux/macOS)
   ↓
2. PIN Verification  
   ↓
3. GPG Key Path Input
   ↓
4. Prerequisites Installation
   ↓
5. GPG Configuration
   ↓
6. Shell Environment Setup
   ↓
7. Key Import & Trust
   ↓
8. YubiKey Detection
   ↓
9. Encryption Testing
   ↓
10. Git Configuration
    ↓
11. Final Instructions
```

## 🔐 Security Features

- **PIN Protection**: Requires YubiKey PIN for operations
- **Key Validation**: Verifies GPG key integrity
- **Secure Defaults**: Uses recommended cipher and digest preferences
- **Agent Security**: Configures secure cache timeouts

## 📝 Post-Setup Steps

After successful script completion:

1. **Add SSH Key to Gitea**:
   - Copy output from `ssh-add -L`
   - Add to Gitea → Settings → SSH/GPG Keys

2. **Add GPG Key to Gitea**:
   - Copy output from `gpg --export -a <key-id>`
   - Add to Gitea → Settings → SSH/GPG Keys

3. **Source Shell Configuration**:
   ```bash
   source ~/.bashrc     # Linux
   source ~/.bash_profile  # macOS
   ```

4. **Test Git Signing**:
   ```bash
   git commit -S -m "test commit"
   ```

## 🐛 Troubleshooting

### If the script fails:

1. **Check the error UUID** - Note the tracking ID for support
2. **Verify prerequisites** - Ensure YubiKey and PIN are available
3. **Check permissions** - Ensure you can write to home directory
4. **Validate GPG key** - Confirm key file exists and is readable

### Getting Help

- **Internal Issues**: Contact Klavs or Ashish
- **Gitea Problems**: Contact your team
- **Script Bugs**: Reference the error UUID when reporting

## ⚠️ Important Notes

- **Backup Important**: The script modifies GPG and shell configurations
- **YubiKey Required**: Physical device must be present during setup
- **Network Access**: Script downloads packages and may access keyservers
- **Sudo Access**: May prompt for password during package installation

## 🔄 Re-running the Script

The script is designed to be idempotent - you can safely run it multiple times. It will:
- Skip already installed packages
- Detect existing configurations
- Only import keys if needed
- Preserve existing settings where appropriate

---

**Script Version**: 1.0  
**Based on**: [EnableIT YubiKey Employee Setup Documentation](https://gitea.obmondo.com/EnableIT/wiki/src/branch/master/internal/yubikey-employee-setup.md)  
**Support**: Contact EnableIT team for assistance