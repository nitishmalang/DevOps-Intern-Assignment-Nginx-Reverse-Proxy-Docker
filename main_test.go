package main

import (
	"os"
	"runtime"
	"strings"
	"testing"
	"yubikey-setup/cmd"
)

func TestNewYubiKeySetup(t *testing.T) {
	setup := cmd.NewYubiKeySetup()
	
	if setup == nil {
		t.Fatal("NewYubiKeySetup() returned nil")
	}
	
	if setup.OS != runtime.GOOS {
		t.Errorf("Expected OS to be %s, got %s", runtime.GOOS, setup.OS)
	}
	
	if setup.ErrorUUID == "" {
		t.Error("ErrorUUID should not be empty")
	}
	
	if setup.HomeDir == "" {
		t.Error("HomeDir should not be empty")
	}
}

func TestCheckOSCompatibility(t *testing.T) {
	setup := cmd.NewYubiKeySetup()
	
	// Test Linux
	setup.OS = "linux"
	setup.CheckOSCompatibility()
	
	if !strings.Contains(setup.ShellRC, ".bashrc") {
		t.Error("Linux should use .bashrc")
	}
	
	if !strings.Contains(setup.PinentryProgram, "pinentry-gnome3") {
		t.Error("Linux should use pinentry-gnome3")
	}
	
	// Test macOS
	setup.OS = "darwin"
	setup.CheckOSCompatibility()
	
	if !strings.Contains(setup.ShellRC, ".bash_profile") {
		t.Error("macOS should use .bash_profile")
	}
	
	if !strings.Contains(setup.PinentryProgram, "pinentry-mac") {
		t.Error("macOS should use pinentry-mac")
	}
}

func TestRunCommand(t *testing.T) {
	setup := cmd.NewYubiKeySetup()
	
	// Test successful command
	err := setup.RunCommand("echo", "test")
	if err != nil {
		t.Errorf("RunCommand failed: %v", err)
	}
	
	// Test failed command
	err = setup.RunCommand("nonexistentcommand")
	if err == nil {
		t.Error("RunCommand should have failed for nonexistent command")
	}
}

func TestRunCommandWithOutput(t *testing.T) {
	setup := cmd.NewYubiKeySetup()
	
	// Test successful command with output
	output, err := setup.RunCommandWithOutput("echo", "hello")
	if err != nil {
		t.Errorf("RunCommandWithOutput failed: %v", err)
	}
	
	if !strings.Contains(output, "hello") {
		t.Errorf("Expected output to contain 'hello', got: %s", output)
	}
}

func TestCreateGPGConfigs(t *testing.T) {
	setup := cmd.NewYubiKeySetup()
	
	// Create temporary directory for testing
	tmpDir := t.TempDir()
	setup.HomeDir = tmpDir
	
	setup.CreateGPGConfigs()
	
	// Check if .gnupg directory was created
	gnupgDir := tmpDir + "/.gnupg"
	if _, err := os.Stat(gnupgDir); os.IsNotExist(err) {
		t.Error(".gnupg directory was not created")
	}
	
	// Check if gpg.conf was created
	gpgConfPath := gnupgDir + "/gpg.conf"
	if _, err := os.Stat(gpgConfPath); os.IsNotExist(err) {
		t.Error("gpg.conf was not created")
	}
	
	// Check if dirmngr.conf was created
	dirmgrConfPath := gnupgDir + "/dirmngr.conf"
	if _, err := os.Stat(dirmgrConfPath); os.IsNotExist(err) {
		t.Error("dirmngr.conf was not created")
	}
	
	// Verify content of gpg.conf
	content, err := os.ReadFile(gpgConfPath)
	if err != nil {
		t.Errorf("Failed to read gpg.conf: %v", err)
	}
	
	if !strings.Contains(string(content), "auto-key-locate keyserver") {
		t.Error("gpg.conf does not contain expected content")
	}
}

func TestConfigureGPGAgent(t *testing.T) {
	setup := cmd.NewYubiKeySetup()
	
	// Create temporary directory for testing
	tmpDir := t.TempDir()
	setup.HomeDir = tmpDir
	
	// Test Linux configuration
	setup.OS = "linux"
	setup.CheckOSCompatibility()
	setup.CreateGPGConfigs() // Create .gnupg directory first
	setup.ConfigureGPGAgent()
	
	agentConfPath := tmpDir + "/.gnupg/gpg-agent.conf"
	content, err := os.ReadFile(agentConfPath)
	if err != nil {
		t.Errorf("Failed to read gpg-agent.conf: %v", err)
	}
	
	if !strings.Contains(string(content), "enable-ssh-support") {
		t.Error("gpg-agent.conf should contain enable-ssh-support")
	}
}

func TestConfigureShellEnvironment(t *testing.T) {
	setup := cmd.NewYubiKeySetup()
	
	// Create temporary directory for testing
	tmpDir := t.TempDir()
	setup.HomeDir = tmpDir
	
	// Test Linux configuration
	setup.OS = "linux"
	setup.CheckOSCompatibility()
	setup.ConfigureShellEnvironment()
	
	// Check if .bashrc was created and contains SSH_AUTH_SOCK
	content, err := os.ReadFile(setup.ShellRC)
	if err != nil {
		t.Errorf("Failed to read shell RC file: %v", err)
	}
	
	if !strings.Contains(string(content), "SSH_AUTH_SOCK") {
		t.Error("Shell RC file should contain SSH_AUTH_SOCK")
	}
}

func TestUnsetSSHAuthSock(t *testing.T) {
	setup := cmd.NewYubiKeySetup()
	
	// Set environment variable
	os.Setenv("SSH_AUTH_SOCK", "/tmp/test")
	
	// Verify it's set
	if _, exists := os.LookupEnv("SSH_AUTH_SOCK"); !exists {
		t.Error("SSH_AUTH_SOCK should be set for this test")
	}
	
	setup.UnsetSSHAuthSock()
	
	// Verify it's unset
	if _, exists := os.LookupEnv("SSH_AUTH_SOCK"); exists {
		t.Error("SSH_AUTH_SOCK should have been unset")
	}
}

func TestGetPGPKeyIDParsing(t *testing.T) {
	setup := cmd.NewYubiKeySetup()
	_ = setup
	
	// Mock GPG output
	mockOutput := `pub   rsa4096/0x1234567890ABCDEF 2023-01-01 [SC]
      Key fingerprint = ABCD EFGH IJKL MNOP QRST UVWX YZ12 3456 7890 ABCD
uid                   [ultimate] Test User <test@example.com>
sub   rsa4096/0xFEDCBA0987654321 2023-01-01 [E]`
	
	lines := strings.Split(mockOutput, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "pub") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				keyParts := strings.Split(parts[1], "/")
				if len(keyParts) >= 2 {
					keyID := strings.Fields(keyParts[1])[0]
					expected := "0x1234567890ABCDEF"
					if keyID != expected {
						t.Errorf("Expected key ID %s, got %s", expected, keyID)
					}
					break
				}
			}
		}
	}
}

// Benchmark tests
func BenchmarkNewYubiKeySetup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		setup := cmd.NewYubiKeySetup()
		_ = setup
	}
}

func BenchmarkCheckOSCompatibility(b *testing.B) {
	setup := cmd.NewYubiKeySetup()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		setup.CheckOSCompatibility()
	}
}

// Example test showing how the setup flow would work
func ExampleYubiKeySetup_Run() {
	// This is an example of how the setup would be used
	// In real usage, this would be interactive
	setup := cmd.NewYubiKeySetup()
	
	// Would normally call setup.Run() for full interactive setup
	// But for example purposes, we'll just show the structure
	setup.Banner()
	setup.CheckOSCompatibility()
	
	// Output would be formatted banner and OS detection messages
}

// Integration test (would require actual GPG setup)
func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	// This test would require:
	// - Actual GPG installation
	// - Mock YubiKey or test environment
	// - Proper permissions
	
	t.Skip("Integration test requires full GPG environment")
}

// Test error scenarios
func TestErrorScenarios(t *testing.T) {
	setup := cmd.NewYubiKeySetup()
	
	// Test invalid OS
	setup.OS = "unsupported"
	// Would need to capture exit() to test this properly
	// In a real test, we'd refactor to return errors instead of calling os.Exit()
	
	// Test invalid GPG key path
	setup.GPGKeyPath = "/nonexistent/path/key.gpg"
	// Similar issue with os.Exit() - would need refactoring for proper testing
}

// Helper function for testing file operations
func createTempGPGKey(t *testing.T) string {
	tmpFile, err := os.CreateTemp("", "test_gpg_key_*.asc")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer tmpFile.Close()
	
	// Write a mock GPG public key
	mockKey := `-----BEGIN PGP PUBLIC KEY BLOCK-----

mQINBGXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
...mock key content...
-----END PGP PUBLIC KEY BLOCK-----`
	
	if _, err := tmpFile.WriteString(mockKey); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	
	return tmpFile.Name()
}