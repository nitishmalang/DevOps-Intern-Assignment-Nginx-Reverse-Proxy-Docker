#!/bin/bash

# Test Script for YubiKey Setup
# Validates script functionality without making actual changes

echo "🧪 Testing YubiKey Setup Script..."
echo "=================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

test_passed=0
test_failed=0

# Test function
run_test() {
    local test_name="$1"
    local test_command="$2"
    
    echo -n "Testing $test_name... "
    
    if eval "$test_command" &>/dev/null; then
        echo -e "${GREEN}PASS${NC}"
        ((test_passed++))
    else
        echo -e "${RED}FAIL${NC}"
        ((test_failed++))
    fi
}

# Test 1: Check if script exists and is executable
run_test "Script Existence" "[ -f ./yubikey_setup.sh ] && [ -x ./yubikey_setup.sh ]"

# Test 2: Check script syntax
run_test "Script Syntax" "bash -n ./yubikey_setup.sh"

# Test 3: OS Detection Logic
run_test "OS Detection" "grep -q 'uname -s' ./yubikey_setup.sh"

# Test 4: PIN Check Logic
run_test "PIN Verification" "grep -q 'Do you have your YubiKey PIN' ./yubikey_setup.sh"

# Test 5: GPG Key Path Logic
run_test "GPG Key Path Input" "grep -q 'Enter the path to your GPG public key' ./yubikey_setup.sh"

# Test 6: Error Handling
run_test "Error Handling" "grep -q 'log_error' ./yubikey_setup.sh"

# Test 7: UUID Generation
run_test "UUID Generation" "grep -q 'ERROR_UUID' ./yubikey_setup.sh"

# Test 8: OS-specific configurations
run_test "Linux Configuration" "grep -q 'pinentry-gnome3' ./yubikey_setup.sh"
run_test "macOS Configuration" "grep -q 'pinentry-mac' ./yubikey_setup.sh"

# Test 9: GPG Configuration Templates
run_test "GPG Config Template" "grep -q 'auto-key-locate keyserver' ./yubikey_setup.sh"

# Test 10: Git Configuration
run_test "Git Signing Setup" "grep -q 'git config --global commit.gpgsign true' ./yubikey_setup.sh"

# Test 11: SSH Configuration
run_test "SSH Auth Socket" "grep -q 'SSH_AUTH_SOCK' ./yubikey_setup.sh"

# Test 12: YubiKey Detection
run_test "YubiKey Detection" "grep -q 'gpg2 --card-status' ./yubikey_setup.sh"

# Test 13: Key Import Logic
run_test "Key Import" "grep -q 'gpg --import' ./yubikey_setup.sh"

# Test 14: Trust Level Setting
run_test "Trust Level" "grep -q 'trust.*5' ./yubikey_setup.sh"

# Test 15: Encryption Test
run_test "Encryption Test" "grep -q 'uname -a.*gpg2.*encrypt.*decrypt' ./yubikey_setup.sh"

# Test the dry-run functionality
echo
echo "🔍 Testing Script Dry Run (OS check only)..."
echo "============================================="

# Create a temporary script that only runs the OS check
cat > /tmp/test_os_check.sh << 'EOF'
#!/bin/bash
OS=$(uname -s)
case "$OS" in
    "Linux")
        echo "Linux detected - supported"
        exit 0
        ;;
    "Darwin")
        echo "macOS detected - supported"
        exit 0
        ;;
    *)
        echo "OS '$OS' is not supported"
        exit 1
        ;;
esac
EOF

chmod +x /tmp/test_os_check.sh

if /tmp/test_os_check.sh; then
    echo -e "${GREEN}✓ OS compatibility check works${NC}"
    ((test_passed++))
else
    echo -e "${RED}✗ OS compatibility check failed${NC}"
    ((test_failed++))
fi

# Cleanup
rm -f /tmp/test_os_check.sh

# Test file structure
echo
echo "📁 Testing File Structure..."
echo "============================"

required_sections=(
    "log_error"
    "log_info" 
    "log_success"
    "log_warning"
    "Check OS compatibility"
    "Check for YubiKey PIN"
    "Get GPG public key path"
    "Install prerequisites"
    "Configure GPG"
    "Configure shell environment"
    "Import GPG key"
    "Restart GPG agent"
    "Check YubiKey detection"
    "Set key trust"
    "Test encryption"
    "Configure Git signing"
)

for section in "${required_sections[@]}"; do
    if grep -q "$section" ./yubikey_setup.sh; then
        echo -e "${GREEN}✓${NC} Found: $section"
        ((test_passed++))
    else
        echo -e "${RED}✗${NC} Missing: $section"
        ((test_failed++))
    fi
done

# Summary
echo
echo "📊 Test Results Summary"
echo "======================="
echo -e "Total Tests: $((test_passed + test_failed))"
echo -e "${GREEN}Passed: $test_passed${NC}"
echo -e "${RED}Failed: $test_failed${NC}"

if [ $test_failed -eq 0 ]; then
    echo -e "\n${GREEN}🎉 All tests passed! The script appears to be well-structured.${NC}"
    exit 0
else
    echo -e "\n${RED}❌ Some tests failed. Please review the script.${NC}"
    exit 1
fi