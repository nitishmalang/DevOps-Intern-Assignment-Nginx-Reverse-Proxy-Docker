#!/bin/bash

# YubiKey Employee Setup Script
# Based on: https://gitea.obmondo.com/EnableIT/wiki/src/branch/master/internal/yubikey-employee-setup.md

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Generate UUID for error tracking
ERROR_UUID=$(cat /proc/sys/kernel/random/uuid 2>/dev/null || python3 -c "import uuid; print(uuid.uuid4())" 2>/dev/null || echo "UNKNOWN-$(date +%s)")

# Logging function
log_error() {
    echo -e "${RED}[ERROR-$ERROR_UUID] $1${NC}" >&2
    exit 1
}

log_info() {
    echo -e "${BLUE}[INFO] $1${NC}"
}

log_success() {
    echo -e "${GREEN}[SUCCESS] $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}[WARNING] $1${NC}"
}

# Banner
echo -e "${BLUE}"
echo "╔══════════════════════════════════════════════════════════════════════════════╗"
echo "║                          YubiKey Setup Script                               ║"
echo "║                         EnableIT Employee Setup                             ║"
echo "╚══════════════════════════════════════════════════════════════════════════════╝"
echo -e "${NC}"

# Check OS compatibility
log_info "Checking operating system compatibility..."
OS=$(uname -s)
case "$OS" in
    "Linux")
        log_success "Linux detected - supported"
        SHELL_RC="$HOME/.bashrc"
        PINENTRY_PROGRAM="/usr/bin/pinentry-gnome3"
        SSH_AUTH_SOCK_PATH="/run/user/$(id -u)/gnupg/S.gpg-agent.ssh"
        ;;
    "Darwin")
        log_success "macOS detected - supported"
        SHELL_RC="$HOME/.bash_profile"
        PINENTRY_PROGRAM="/opt/homebrew/bin/pinentry-mac"
        SSH_AUTH_SOCK_PATH='$(gpgconf --list-dirs agent-ssh-socket)'
        ;;
    *)
        log_error "OS '$OS' is not supported. This script only works on Linux and macOS."
        ;;
esac

# Check for YubiKey PIN
echo
read -p "Do you have your YubiKey PIN? (y/n): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}You need your YubiKey PIN to proceed.${NC}"
    echo "Please contact Klavs or Ashish to get your PIN."
    echo "Reference: https://gitea.obmondo.com/EnableIT/pass"
    exit 1
fi

# Get GPG public key path
echo
read -p "Enter the path to your GPG public key (e.g., ~/abc.key): " GPG_KEY_PATH

# Expand tilde
GPG_KEY_PATH="${GPG_KEY_PATH/#\~/$HOME}"

# Validate GPG key file
if [[ ! -f "$GPG_KEY_PATH" ]]; then
    log_error "GPG public key file not found at '$GPG_KEY_PATH'. Please get it from https://gitea.obmondo.com/EnableIT/pass and try again."
fi

log_success "GPG public key found at: $GPG_KEY_PATH"

# Install prerequisites based on OS
log_info "Installing prerequisites..."
if [[ "$OS" == "Linux" ]]; then
    if ! command -v apt-get &> /dev/null; then
        log_error "apt-get not found. Please install prerequisites manually: gnupg2 gnupg-agent pinentry-curses scdaemon pcscd libusb-1.0-0-dev"
    fi
    
    if ! dpkg -l | grep -q gnupg2; then
        log_info "Installing Linux prerequisites..."
        sudo apt-get update || log_error "Failed to update package list"
        sudo apt-get install -y gnupg2 gnupg-agent pinentry-curses scdaemon pcscd libusb-1.0-0-dev pinentry-gnome3 || log_error "Failed to install prerequisites"
    fi
elif [[ "$OS" == "Darwin" ]]; then
    if ! command -v brew &> /dev/null; then
        log_error "Homebrew not found. Please install Homebrew first: https://brew.sh/"
    fi
    
    if ! brew list gnupg &> /dev/null; then
        log_info "Installing macOS prerequisites..."
        brew install gnupg yubikey-personalization pinentry-mac || log_error "Failed to install prerequisites"
    fi
fi

# Create .gnupg directory if it doesn't exist
mkdir -p ~/.gnupg
chmod 700 ~/.gnupg

# Write GPG configuration
log_info "Configuring GPG..."
cat > ~/.gnupg/gpg.conf << 'EOF' || log_error "Failed to write GPG configuration"
auto-key-locate keyserver
keyserver hkps://hkps.pool.sks-keyservers.net
keyserver-options no-honor-keyserver-url
personal-cipher-preferences AES256 AES192 AES CAST5
personal-digest-preferences SHA512 SHA384 SHA256 SHA224
default-preference-list SHA512 SHA384 SHA256 SHA224 AES256 AES192 AES CAST5 ZLIB BZIP2 ZIP Uncompressed
cert-digest-algo SHA512
s2k-cipher-algo AES256
s2k-digest-algo SHA512
charset utf-8
fixed-list-mode
no-comments
no-emit-version
keyid-format 0xlong
list-options show-uid-validity
verify-options show-uid-validity
with-fingerprint
use-agent
require-cross-certification
EOF

cat > ~/.gnupg/dirmngr.conf << 'EOF' || log_error "Failed to write dirmngr configuration"
keyserver hkp://jirk5u4osbsr34t5.onion
keyserver hkp://keys.gnupg.net
honor-http-proxy
hkp-cacert /etc/sks-keyservers.netCA.pem
EOF

# Configure GPG agent
log_info "Configuring GPG agent..."
if [[ "$OS" == "Linux" ]]; then
    cat > ~/.gnupg/gpg-agent.conf << EOF || log_error "Failed to write GPG agent configuration"
# enables SSH support (ssh-agent)
enable-ssh-support
#remote
extra-socket /run/user/\$(id -u)/gnupg/S.gpg-agent-extra
# default cache timeout of 600 seconds
default-cache-ttl 600
max-cache-ttl 7200
EOF
elif [[ "$OS" == "Darwin" ]]; then
    # Check if pinentry-mac exists
    if [[ ! -f "$PINENTRY_PROGRAM" ]]; then
        log_error "pinentry-mac not found at $PINENTRY_PROGRAM. Please ensure it's installed correctly."
    fi
    
    cat > ~/.gnupg/gpg-agent.conf << EOF || log_error "Failed to write GPG agent configuration"
pinentry-program $PINENTRY_PROGRAM
enable-ssh-support
# default cache timeout of 600 seconds
default-cache-ttl 600
max-cache-ttl 7200
EOF
fi

# Configure shell environment
log_info "Configuring shell environment..."
if [[ "$OS" == "Linux" ]]; then
    if ! grep -q "SSH_AUTH_SOCK=/run/user" "$SHELL_RC" 2>/dev/null; then
        echo "SSH_AUTH_SOCK=/run/user/\$(id -u)/gnupg/S.gpg-agent.ssh" >> "$SHELL_RC" || log_error "Failed to update $SHELL_RC"
    fi
elif [[ "$OS" == "Darwin" ]]; then
    if ! grep -q "SSH_AUTH_SOCK.*gpgconf" "$SHELL_RC" 2>/dev/null; then
        echo "export SSH_AUTH_SOCK=\$(gpgconf --list-dirs agent-ssh-socket)" >> "$SHELL_RC" || log_error "Failed to update $SHELL_RC"
    fi
fi

# Check and unset existing SSH_AUTH_SOCK
log_info "Checking for existing SSH_AUTH_SOCK..."
if env | grep -q SSH_AUTH_SOCK; then
    log_warning "Found existing SSH_AUTH_SOCK, unsetting it..."
    unset SSH_AUTH_SOCK || log_error "Failed to unset SSH_AUTH_SOCK"
fi

# Import GPG key
log_info "Importing GPG public key..."
# Check if key already exists
if gpg --list-keys | grep -q "0x3996B9E90711DD51"; then
    log_warning "Key 0x3996B9E90711DD51 already exists, skipping import"
else
    if ! gpg --import "$GPG_KEY_PATH"; then
        if echo "$?" | grep -q "key not changed"; then
            log_warning "Key already imported (key not changed)"
        else
            log_error "Failed to import GPG key from $GPG_KEY_PATH"
        fi
    fi
fi

# Restart GPG agent
log_info "Restarting GPG agent..."
if pgrep gpg-agent > /dev/null; then
    pkill gpg-agent || log_error "Failed to kill GPG agent"
    sleep 2
fi

# Verify GPG agent is not running
if pgrep gpg-agent > /dev/null; then
    log_error "GPG agent still running after kill attempt"
fi

# Start GPG agent fresh
gpg-agent --daemon > /dev/null 2>&1

# Check YubiKey detection
log_info "Checking YubiKey detection..."
echo "Please insert your YubiKey if not already inserted and press Enter to continue..."
read -r

if ! gpg2 --card-status &> /dev/null; then
    log_error "YubiKey not detected. Please replug your YubiKey and try again."
fi

# For Linux, ensure pinentry-gnome3 is set
if [[ "$OS" == "Linux" ]]; then
    if command -v update-alternatives &> /dev/null; then
        log_info "Setting pinentry to pinentry-gnome3..."
        sudo update-alternatives --install /usr/bin/pinentry pinentry /usr/bin/pinentry-gnome3 1
        sudo update-alternatives --set pinentry /usr/bin/pinentry-gnome3
    fi
fi

# Check SSH support
log_info "Checking SSH support..."
if ! ssh-add -L | grep -q "cardno:"; then
    log_error "YubiKey card number not found in ssh-add -L output. Please replug YubiKey and restart the script."
fi

# Get PGP Key ID
log_info "Getting PGP Key ID..."
export PGP_KEY_ID=$(gpg2 --list-keys --keyid-format 0xlong | grep -E "^pub" | cut -d "/" -f 2 | cut -d " " -f 1 | head -1)
if [[ -z "$PGP_KEY_ID" ]]; then
    log_error "Could not determine PGP Key ID"
fi
log_success "PGP Key ID: $PGP_KEY_ID"

# Set key trust
log_info "Setting key trust to ultimate..."
echo "Setting trust level to 5 (ultimate) for key $PGP_KEY_ID"
echo -e "trust\n5\ny\nsave\n" | gpg --edit-key "$PGP_KEY_ID" --command-fd 0 || {
    # Try with email if key ID doesn't work
    USER_EMAIL=$(whoami)@obmondo.com
    log_warning "Trying with email: $USER_EMAIL"
    echo -e "trust\n5\ny\nsave\n" | gpg --edit-key "$USER_EMAIL" --command-fd 0 || log_error "Failed to set key trust"
}

# Set GPG_TTY for testing
export GPG_TTY=$(tty)

# Test encryption/decryption
log_info "Testing encryption and decryption..."
echo "Please enter your PIN when prompted and touch your YubiKey when it blinks..."

TEST_RESULT=$(uname -a | gpg2 --encrypt --armor --recipient "$PGP_KEY_ID" 2>/dev/null | gpg2 --decrypt 2>/dev/null)
if [[ $? -ne 0 ]] || [[ -z "$TEST_RESULT" ]]; then
    # Try with email instead
    USER_EMAIL=$(whoami)@obmondo.com
    log_warning "Trying encryption test with email: $USER_EMAIL"
    TEST_RESULT=$(uname -a | gpg2 --encrypt --armor --recipient "$USER_EMAIL" 2>/dev/null | gpg2 --decrypt 2>/dev/null)
    if [[ $? -ne 0 ]] || [[ -z "$TEST_RESULT" ]]; then
        log_error "Encryption/decryption test failed. Make sure pinentry-gnome3 is set and GPG_TTY is exported."
    fi
fi

log_success "Encryption/decryption test passed!"

# Configure Git signing
log_info "Configuring Git signing..."
if ! command -v git &> /dev/null; then
    log_error "Git is not installed. Please install Git first."
fi

git config --global user.signingkey "$PGP_KEY_ID" || log_error "Failed to set Git signing key"
git config --global commit.gpgsign true || log_error "Failed to enable Git commit signing"

log_success "Git signing configured successfully!"

# Final instructions
echo
echo -e "${GREEN}╔══════════════════════════════════════════════════════════════════════════════╗"
echo "║                          Setup Complete!                                    ║"
echo "╚══════════════════════════════════════════════════════════════════════════════╝${NC}"
echo
log_info "Next steps:"
echo "1. Add your SSH public key to Gitea:"
echo "   - Go to Gitea -> Settings -> SSH/GPG Keys -> Manage SSH Keys"
echo "   - Use this key:"
echo
ssh-add -L | head -1
echo
echo "2. Add your GPG public key to Gitea:"
echo "   - Go to Gitea -> Settings -> SSH/GPG Keys -> Manage GPG Keys"
echo "   - Use this key:"
echo
gpg --export -a "$PGP_KEY_ID" | head -10
echo "   [... truncated for display ...]"
echo
echo "3. Source your shell configuration:"
echo "   source $SHELL_RC"
echo
echo "4. Test Git signing in a repository:"
echo "   git commit -S -m 'test commit'"
echo
echo "If you encounter any issues with Gitea, please contact your team."
echo
log_success "Setup completed successfully! Error tracking ID: $ERROR_UUID"