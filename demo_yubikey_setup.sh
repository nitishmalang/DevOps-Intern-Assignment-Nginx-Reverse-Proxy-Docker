#!/bin/bash

# Demo Script for YubiKey Setup
# Shows user interaction flow without making actual changes

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}"
echo "╔══════════════════════════════════════════════════════════════════════════════╗"
echo "║                     YubiKey Setup Script - DEMO MODE                        ║"
echo "║                        EnableIT Employee Setup                              ║"
echo "╚══════════════════════════════════════════════════════════════════════════════╝"
echo -e "${NC}"

echo -e "${YELLOW}🚨 DEMO MODE: This script will show the user interaction flow without making changes${NC}"
echo

# Simulate OS detection
echo -e "${BLUE}[INFO] Checking operating system compatibility...${NC}"
sleep 1
echo -e "${GREEN}[SUCCESS] Linux detected - supported${NC}"
echo

# Simulate PIN check
echo "The script will ask for YubiKey PIN confirmation:"
echo -n "Do you have your YubiKey PIN? (y/n): "
read -n 1 -r pin_response
echo

if [[ ! $pin_response =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}You need your YubiKey PIN to proceed.${NC}"
    echo "Please contact Klavs or Ashish to get your PIN."
    echo "Reference: https://gitea.obmondo.com/EnableIT/pass"
    echo -e "${RED}[DEMO] Would exit here in real script${NC}"
else
    echo -e "${GREEN}[SUCCESS] PIN confirmed${NC}"
fi

echo

# Simulate GPG key path input
echo "The script will ask for GPG public key path:"
echo -n "Enter the path to your GPG public key (e.g., ~/abc.key): "
read gpg_path

if [[ -z "$gpg_path" ]]; then
    gpg_path="~/demo.key"
fi

echo -e "${BLUE}[INFO] Would validate GPG key at: $gpg_path${NC}"
echo -e "${GREEN}[SUCCESS] GPG public key found at: $gpg_path${NC}"
echo

# Show what the script would do
echo -e "${BLUE}[INFO] The script would now perform these steps:${NC}"
echo
echo "1. 📦 Install prerequisites (gnupg2, pinentry, etc.)"
echo "2. ⚙️  Configure GPG (~/.gnupg/gpg.conf, dirmngr.conf)"
echo "3. 🔧 Set up GPG agent (~/.gnupg/gpg-agent.conf)"
echo "4. 🐚 Configure shell environment (.bashrc/.bash_profile)"
echo "5. 🔐 Import GPG public key"
echo "6. 🔄 Restart GPG agent"
echo "7. 🎯 Check YubiKey detection"
echo "8. 🔒 Set key trust to ultimate (5)"
echo "9. 🧪 Test encryption/decryption"
echo "10. 📝 Configure Git signing"
echo "11. ✅ Provide final instructions"

echo
echo -e "${YELLOW}⚠️  During actual execution, you would need to:${NC}"
echo "• Insert your YubiKey"
echo "• Enter your PIN when prompted"
echo "• Touch your YubiKey when it blinks"
echo "• Provide sudo password for package installation"

echo
echo -e "${BLUE}[INFO] Final step instructions would include:${NC}"
echo
echo "1. Add SSH key to Gitea (from ssh-add -L output)"
echo "2. Add GPG key to Gitea (from gpg --export -a output)"
echo "3. Source shell configuration"
echo "4. Test Git signing with a commit"

echo
echo -e "${GREEN}╔══════════════════════════════════════════════════════════════════════════════╗"
echo "║                          Demo Complete!                                     ║"
echo "╚══════════════════════════════════════════════════════════════════════════════╝${NC}"
echo
echo "To run the actual setup, execute: ./yubikey_setup.sh"
echo "To run tests, execute: ./test_yubikey_setup.sh"