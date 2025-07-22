# YubiKey Setup Tool (Go Version)

A comprehensive YubiKey setup automation tool for EnableIT employees, written in Go. This tool automates the complete YubiKey configuration process for both Linux and macOS systems.

## 🎯 Features

- **OS Compatibility**: Supports Linux and macOS with appropriate configurations
- **PIN Verification**: Ensures user has YubiKey PIN before proceeding
- **GPG Configuration**: Automated GPG setup with security best practices
- **YubiKey Integration**: Complete YubiKey detection and configuration
- **SSH Authentication**: Sets up SSH authentication via YubiKey
- **Git Signing**: Configures automatic Git commit signing
- **Error Handling**: Comprehensive error tracking with UUIDs
- **CLI Interface**: Modern command-line interface with multiple commands
- **Dry Run Mode**: Test what changes would be made without executing them

## 📋 Prerequisites

- **YubiKey device** - Physical hardware token
- **YubiKey PIN** - Contact Klavs or Ashish if you don't have it
- **GPG public key** - Available at https://gitea.obmondo.com/EnableIT/pass
- **Go 1.21+** - For building from source
- **Administrative access** - May need sudo for package installation

## 🚀 Installation

### Option 1: Download Binary (Recommended)
```bash
# Download latest release for your platform
wget https://github.com/enableit/yubikey-setup/releases/latest/download/yubikey-setup-linux-amd64
chmod +x yubikey-setup-linux-amd64
sudo mv yubikey-setup-linux-amd64 /usr/local/bin/yubikey-setup
```

### Option 2: Build from Source
```bash
# Clone the repository
git clone https://github.com/enableit/yubikey-setup.git
cd yubikey-setup

# Install dependencies
go mod download

# Build the binary
go build -o yubikey-setup .

# Install globally (optional)
sudo cp yubikey-setup /usr/local/bin/
```

## 📖 Usage

### Quick Start
```bash
# Run the complete setup process
yubikey-setup setup

# Run with specific GPG key path
yubikey-setup setup --gpg-key ~/my-key.asc

# Dry run to see what would be done
yubikey-setup setup --dry-run

# Check system prerequisites
yubikey-setup check

# Get help
yubikey-setup --help
```

### Available Commands

#### `setup` - Complete YubiKey Setup
Runs the full automated setup process:
```bash
# Basic usage
yubikey-setup setup

# With options
yubikey-setup setup --gpg-key ~/abc.key --skip-pin --dry-run -v
```

**Flags:**
- `--gpg-key, -k string` - Path to GPG public key file
- `--skip-pin` - Skip PIN verification (not recommended)
- `--dry-run` - Show what would be done without making changes
- `--verbose, -v` - Verbose output

#### `check` - System Prerequisites Check
Validates system requirements and current setup:
```bash
yubikey-setup check
```

Checks:
- Operating system compatibility
- Required packages installation
- GPG configuration status
- YubiKey detection
- SSH agent status

#### Global Flags
- `--verbose, -v` - Enable verbose output
- `--dry-run` - Dry run mode for all commands
- `--help, -h` - Show help
- `--version` - Show version information

## 🔧 What the Tool Does

### 1. System Validation
- Detects OS (Linux/macOS) and validates compatibility
- Checks for required packages and dependencies
- Validates YubiKey PIN availability

### 2. Prerequisites Installation
**Linux (Ubuntu/Debian):**
- `gnupg2`, `gnupg-agent`, `pinentry-curses`
- `scdaemon`, `pcscd`, `libusb-1.0-0-dev`
- `pinentry-gnome3`

**macOS (Homebrew):**
- `gnupg`, `yubikey-personalization`
- `pinentry-mac`

### 3. GPG Configuration
Creates and configures:
- `~/.gnupg/gpg.conf` - GPG settings with security preferences
- `~/.gnupg/dirmngr.conf` - Keyserver configuration
- `~/.gnupg/gpg-agent.conf` - Agent settings with SSH support

### 4. Shell Environment
- Configures `SSH_AUTH_SOCK` in `.bashrc` (Linux) or `.bash_profile` (macOS)
- Unsets conflicting environment variables

### 5. YubiKey Setup
- Imports GPG public key
- Restarts GPG agent
- Validates YubiKey detection
- Sets key trust to ultimate (level 5)
- Tests encryption/decryption functionality

### 6. Git Integration
- Configures global Git signing key
- Enables automatic commit signing

## 🛡️ Error Handling

The tool includes comprehensive error handling:

- **UUID Tracking**: Each run generates a unique error ID for troubleshooting
- **Validation Checks**: Verifies each step before proceeding
- **Graceful Failures**: Stops on critical errors with clear messages
- **Common Issue Fixes**: Handles known problems automatically

### Common Issues Handled
| Issue | Solution |
|-------|----------|
| "Screen or window too small" | Sets `GPG_TTY` and `pinentry-gnome3` |
| "No secret key" | Uses alternative pinentry configuration |
| "Key not changed" | Skips import if key already exists |
| YubiKey not detected | Prompts to replug device |

## 📊 Example Output

```bash
$ yubikey-setup setup
╔══════════════════════════════════════════════════════════════════════════════╗
║                          YubiKey Setup Tool                                 ║
║                         EnableIT Employee Setup                             ║
╚══════════════════════════════════════════════════════════════════════════════╝

ℹ️ [INFO] Checking operating system compatibility...
✅ [SUCCESS] Linux detected - supported
Do you have your YubiKey PIN? (y/n): y
Enter the path to your GPG public key (e.g., ~/abc.key): ~/my-key.asc
✅ [SUCCESS] GPG public key found at: /home/user/my-key.asc
ℹ️ [INFO] Installing prerequisites...
ℹ️ [INFO] Configuring GPG...
ℹ️ [INFO] Configuring GPG agent...
...
✅ [SUCCESS] Setup completed successfully! Error tracking ID: 12345678-1234-5678-9abc-123456789012
```

## 🧪 Testing

```bash
# Run unit tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestNewYubiKeySetup

# Run benchmarks
go test -bench=.
```

## 🔄 Development

### Project Structure
```
yubikey-setup/
├── main.go              # Application entry point
├── go.mod               # Go module definition
├── go.sum               # Dependency checksums
├── cmd/                 # CLI commands
│   ├── root.go          # Root command definition
│   ├── setup.go         # Setup command
│   ├── check.go         # Check command
│   └── setup_logic.go   # Core setup functionality
├── main_test.go         # Test files
└── README.md            # This file
```

### Adding New Commands
1. Create a new file in `cmd/` directory
2. Define the command using Cobra
3. Add the command to root in `init()` function
4. Implement the command logic

### Building for Different Platforms
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o yubikey-setup-linux-amd64

# macOS
GOOS=darwin GOARCH=amd64 go build -o yubikey-setup-darwin-amd64

# Windows (if needed)
GOOS=windows GOARCH=amd64 go build -o yubikey-setup-windows-amd64.exe
```

## 📝 Post-Setup Steps

After successful completion:

1. **Add SSH Key to Gitea**:
   ```bash
   ssh-add -L
   ```
   Copy the output and add to Gitea → Settings → SSH/GPG Keys

2. **Add GPG Key to Gitea**:
   ```bash
   gpg --export -a <key-id>
   ```
   Copy the output and add to Gitea → Settings → SSH/GPG Keys

3. **Source Shell Configuration**:
   ```bash
   source ~/.bashrc      # Linux
   source ~/.bash_profile # macOS
   ```

4. **Test Git Signing**:
   ```bash
   git commit -S -m "test commit"
   ```

## ❗ Troubleshooting

### Common Issues

1. **GPG Agent Issues**:
   ```bash
   pkill gpg-agent
   gpg-agent --daemon
   ```

2. **YubiKey Not Detected**:
   - Replug the YubiKey
   - Check `gpg2 --card-status`
   - Ensure proper drivers are installed

3. **Permission Issues**:
   - Check file permissions on `.gnupg` directory
   - Ensure `~/.gnupg` is mode 700

4. **Pinentry Issues**:
   ```bash
   export GPG_TTY=$(tty)
   sudo update-alternatives --config pinentry
   ```

### Getting Help

- **Internal Issues**: Contact Klavs or Ashish
- **Gitea Problems**: Contact your team
- **Tool Bugs**: Reference the error UUID when reporting

## 📚 References

- [EnableIT YubiKey Employee Setup Guide](https://gitea.obmondo.com/EnableIT/wiki/src/branch/master/internal/yubikey-employee-setup.md)
- [EnableIT Pass Repository](https://gitea.obmondo.com/EnableIT/pass)
- [Cobra CLI Library](https://github.com/spf13/cobra)

## 📄 License

This tool is for internal EnableIT use only.

---

**Version**: 1.0.0  
**Language**: Go 1.21+  
**Platforms**: Linux, macOS  
**Status**: Production Ready