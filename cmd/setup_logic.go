package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/google/uuid"
)

// Colors for output
var (
	red    = color.New(color.FgRed).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	blue   = color.New(color.FgBlue).SprintFunc()
	cyan   = color.New(color.FgCyan).SprintFunc()
)

// YubiKeySetup represents the main setup configuration
type YubiKeySetup struct {
	OS              string
	ErrorUUID       string
	ShellRC         string
	PinentryProgram string
	SSHAuthSockPath string
	GPGKeyPath      string
	PGPKeyID        string
	HomeDir         string
	DryRun          bool
	SkipPinCheck    bool
}

// NewYubiKeySetup creates a new setup instance
func NewYubiKeySetup() *YubiKeySetup {
	currentUser, _ := user.Current()
	return &YubiKeySetup{
		OS:        runtime.GOOS,
		ErrorUUID: uuid.New().String(),
		HomeDir:   currentUser.HomeDir,
	}
}

// Banner displays the application banner
func (y *YubiKeySetup) Banner() {
	fmt.Println(blue("╔══════════════════════════════════════════════════════════════════════════════╗"))
	fmt.Println(blue("║                          YubiKey Setup Tool                                 ║"))
	fmt.Println(blue("║                         EnableIT Employee Setup                             ║"))
	fmt.Println(blue("╚══════════════════════════════════════════════════════════════════════════════╝"))
	fmt.Println()
}

// LogError logs an error and exits
func (y *YubiKeySetup) LogError(msg string) {
	fmt.Printf("%s [ERROR-%s] %s\n", red("❌"), y.ErrorUUID, msg)
	if !y.DryRun {
		os.Exit(1)
	}
}

// LogInfo logs an info message
func (y *YubiKeySetup) LogInfo(msg string) {
	fmt.Printf("%s [INFO] %s\n", blue("ℹ️"), msg)
}

// LogSuccess logs a success message
func (y *YubiKeySetup) LogSuccess(msg string) {
	fmt.Printf("%s [SUCCESS] %s\n", green("✅"), msg)
}

// LogWarning logs a warning message
func (y *YubiKeySetup) LogWarning(msg string) {
	fmt.Printf("%s [WARNING] %s\n", yellow("⚠️"), msg)
}

// CheckOSCompatibility checks if the OS is supported
func (y *YubiKeySetup) CheckOSCompatibility() {
	y.LogInfo("Checking operating system compatibility...")
	
	switch y.OS {
	case "linux":
		y.LogSuccess("Linux detected - supported")
		y.ShellRC = filepath.Join(y.HomeDir, ".bashrc")
		y.PinentryProgram = "/usr/bin/pinentry-gnome3"
		currentUser, _ := user.Current()
		y.SSHAuthSockPath = fmt.Sprintf("/run/user/%s/gnupg/S.gpg-agent.ssh", currentUser.Uid)
	case "darwin":
		y.LogSuccess("macOS detected - supported")
		y.ShellRC = filepath.Join(y.HomeDir, ".bash_profile")
		y.PinentryProgram = "/opt/homebrew/bin/pinentry-mac"
		y.SSHAuthSockPath = "$(gpgconf --list-dirs agent-ssh-socket)"
	default:
		y.LogError(fmt.Sprintf("OS '%s' is not supported. This tool only works on Linux and macOS.", y.OS))
	}
}

// CheckYubiKeyPIN asks user if they have their YubiKey PIN
func (y *YubiKeySetup) CheckYubiKeyPIN() {
	if y.SkipPinCheck {
		y.LogWarning("Skipping PIN check (not recommended)")
		return
	}
	
	fmt.Print("Do you have your YubiKey PIN? (y/n): ")
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	
	if response != "y" && response != "yes" {
		fmt.Println(yellow("You need your YubiKey PIN to proceed."))
		fmt.Println("Please contact Klavs or Ashish to get your PIN.")
		fmt.Println("Reference: https://gitea.obmondo.com/EnableIT/pass")
		if !y.DryRun {
			os.Exit(1)
		}
	}
}

// GetGPGKeyPath gets and validates the GPG public key path
func (y *YubiKeySetup) GetGPGKeyPath() {
	if y.GPGKeyPath != "" {
		y.LogInfo(fmt.Sprintf("Using provided GPG key path: %s", y.GPGKeyPath))
	} else {
		fmt.Print("Enter the path to your GPG public key (e.g., ~/abc.key): ")
		reader := bufio.NewReader(os.Stdin)
		path, _ := reader.ReadString('\n')
		path = strings.TrimSpace(path)
		
		// Expand tilde
		if strings.HasPrefix(path, "~/") {
			path = filepath.Join(y.HomeDir, path[2:])
		}
		
		y.GPGKeyPath = path
	}
	
	// Check if file exists
	if _, err := os.Stat(y.GPGKeyPath); os.IsNotExist(err) {
		y.LogError(fmt.Sprintf("GPG public key file not found at '%s'. Please get it from https://gitea.obmondo.com/EnableIT/pass and try again.", y.GPGKeyPath))
	}
	
	y.LogSuccess(fmt.Sprintf("GPG public key found at: %s", y.GPGKeyPath))
}

// RunCommand executes a command and returns error if it fails
func (y *YubiKeySetup) RunCommand(name string, args ...string) error {
	if y.DryRun {
		y.LogInfo(fmt.Sprintf("DRY RUN: Would execute: %s %s", name, strings.Join(args, " ")))
		return nil
	}
	cmd := exec.Command(name, args...)
	return cmd.Run()
}

// RunCommandWithOutput executes a command and returns output
func (y *YubiKeySetup) RunCommandWithOutput(name string, args ...string) (string, error) {
	if y.DryRun {
		y.LogInfo(fmt.Sprintf("DRY RUN: Would execute: %s %s", name, strings.Join(args, " ")))
		return "dry-run-output", nil
	}
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	return string(output), err
}

// InstallPrerequisites installs required packages
func (y *YubiKeySetup) InstallPrerequisites() {
	y.LogInfo("Installing prerequisites...")
	
	switch y.OS {
	case "linux":
		// Check if apt-get exists
		if err := y.RunCommand("which", "apt-get"); err != nil {
			y.LogError("apt-get not found. Please install prerequisites manually: gnupg2 gnupg-agent pinentry-curses scdaemon pcscd libusb-1.0-0-dev")
		}
		
		// Check if gnupg2 is installed
		if err := y.RunCommand("dpkg", "-l", "gnupg2"); err != nil {
			y.LogInfo("Installing Linux prerequisites...")
			if err := y.RunCommand("sudo", "apt-get", "update"); err != nil {
				y.LogError("Failed to update package list")
			}
			packages := []string{"gnupg2", "gnupg-agent", "pinentry-curses", "scdaemon", "pcscd", "libusb-1.0-0-dev", "pinentry-gnome3"}
			args := append([]string{"apt-get", "install", "-y"}, packages...)
			if err := y.RunCommand("sudo", args...); err != nil {
				y.LogError("Failed to install prerequisites")
			}
		}
		
	case "darwin":
		// Check if brew exists
		if err := y.RunCommand("which", "brew"); err != nil {
			y.LogError("Homebrew not found. Please install Homebrew first: https://brew.sh/")
		}
		
		// Check if gnupg is installed
		if err := y.RunCommand("brew", "list", "gnupg"); err != nil {
			y.LogInfo("Installing macOS prerequisites...")
			packages := []string{"gnupg", "yubikey-personalization", "pinentry-mac"}
			for _, pkg := range packages {
				if err := y.RunCommand("brew", "install", pkg); err != nil {
					y.LogError("Failed to install prerequisites")
				}
			}
		}
	}
}

// CreateGPGConfigs creates GPG configuration files
func (y *YubiKeySetup) CreateGPGConfigs() {
	y.LogInfo("Configuring GPG...")
	
	gnupgDir := filepath.Join(y.HomeDir, ".gnupg")
	if y.DryRun {
		y.LogInfo(fmt.Sprintf("DRY RUN: Would create directory: %s", gnupgDir))
	} else {
		if err := os.MkdirAll(gnupgDir, 0700); err != nil {
			y.LogError("Failed to create .gnupg directory")
		}
	}
	
	// Create gpg.conf
	gpgConf := `auto-key-locate keyserver
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
require-cross-certification`
	
	gpgConfPath := filepath.Join(gnupgDir, "gpg.conf")
	if y.DryRun {
		y.LogInfo(fmt.Sprintf("DRY RUN: Would write GPG config to: %s", gpgConfPath))
	} else {
		if err := os.WriteFile(gpgConfPath, []byte(gpgConf), 0644); err != nil {
			y.LogError("Failed to write GPG configuration")
		}
	}
	
	// Create dirmngr.conf
	dirmgrConf := `keyserver hkp://jirk5u4osbsr34t5.onion
keyserver hkp://keys.gnupg.net
honor-http-proxy
hkp-cacert /etc/sks-keyservers.netCA.pem`
	
	dirmgrConfPath := filepath.Join(gnupgDir, "dirmngr.conf")
	if y.DryRun {
		y.LogInfo(fmt.Sprintf("DRY RUN: Would write dirmngr config to: %s", dirmgrConfPath))
	} else {
		if err := os.WriteFile(dirmgrConfPath, []byte(dirmgrConf), 0644); err != nil {
			y.LogError("Failed to write dirmngr configuration")
		}
	}
}

// ConfigureGPGAgent configures the GPG agent
func (y *YubiKeySetup) ConfigureGPGAgent() {
	y.LogInfo("Configuring GPG agent...")
	gnupgDir := filepath.Join(y.HomeDir, ".gnupg")
	
	var agentConf string
	switch y.OS {
	case "linux":
		currentUser, _ := user.Current()
		agentConf = fmt.Sprintf(`# enables SSH support (ssh-agent)
enable-ssh-support
#remote
extra-socket /run/user/%s/gnupg/S.gpg-agent-extra
# default cache timeout of 600 seconds
default-cache-ttl 600
max-cache-ttl 7200`, currentUser.Uid)
		
	case "darwin":
		// Check if pinentry-mac exists
		if !y.DryRun {
			if _, err := os.Stat(y.PinentryProgram); os.IsNotExist(err) {
				y.LogError(fmt.Sprintf("pinentry-mac not found at %s. Please ensure it's installed correctly.", y.PinentryProgram))
			}
		}
		
		agentConf = fmt.Sprintf(`pinentry-program %s
enable-ssh-support
# default cache timeout of 600 seconds
default-cache-ttl 600
max-cache-ttl 7200`, y.PinentryProgram)
	}
	
	agentConfPath := filepath.Join(gnupgDir, "gpg-agent.conf")
	if y.DryRun {
		y.LogInfo(fmt.Sprintf("DRY RUN: Would write GPG agent config to: %s", agentConfPath))
	} else {
		if err := os.WriteFile(agentConfPath, []byte(agentConf), 0644); err != nil {
			y.LogError("Failed to write GPG agent configuration")
		}
	}
}

// ConfigureShellEnvironment configures shell environment
func (y *YubiKeySetup) ConfigureShellEnvironment() {
	y.LogInfo("Configuring shell environment...")
	
	var envLine string
	switch y.OS {
	case "linux":
		currentUser, _ := user.Current()
		envLine = fmt.Sprintf("SSH_AUTH_SOCK=/run/user/%s/gnupg/S.gpg-agent.ssh", currentUser.Uid)
	case "darwin":
		envLine = "export SSH_AUTH_SOCK=$(gpgconf --list-dirs agent-ssh-socket)"
	}
	
	if y.DryRun {
		y.LogInfo(fmt.Sprintf("DRY RUN: Would add to %s: %s", y.ShellRC, envLine))
		return
	}
	
	// Check if line already exists
	if content, err := os.ReadFile(y.ShellRC); err == nil {
		if strings.Contains(string(content), "SSH_AUTH_SOCK") && strings.Contains(string(content), "gnupg") {
			return // Already configured
		}
	}
	
	// Append to shell RC file
	file, err := os.OpenFile(y.ShellRC, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		y.LogError(fmt.Sprintf("Failed to update %s", y.ShellRC))
	}
	defer file.Close()
	
	if _, err := file.WriteString(fmt.Sprintf("\n%s\n", envLine)); err != nil {
		y.LogError(fmt.Sprintf("Failed to update %s", y.ShellRC))
	}
}

// UnsetSSHAuthSock unsets existing SSH_AUTH_SOCK
func (y *YubiKeySetup) UnsetSSHAuthSock() {
	y.LogInfo("Checking for existing SSH_AUTH_SOCK...")
	if _, exists := os.LookupEnv("SSH_AUTH_SOCK"); exists {
		y.LogWarning("Found existing SSH_AUTH_SOCK, unsetting it...")
		if !y.DryRun {
			os.Unsetenv("SSH_AUTH_SOCK")
		}
	}
}

// ImportGPGKey imports the GPG public key
func (y *YubiKeySetup) ImportGPGKey() {
	y.LogInfo("Importing GPG public key...")
	
	// Check if key already exists
	if output, err := y.RunCommandWithOutput("gpg", "--list-keys"); err == nil {
		if strings.Contains(output, "0x3996B9E90711DD51") {
			y.LogWarning("Key 0x3996B9E90711DD51 already exists, skipping import")
			return
		}
	}
	
	if err := y.RunCommand("gpg", "--import", y.GPGKeyPath); err != nil {
		y.LogError(fmt.Sprintf("Failed to import GPG key from %s", y.GPGKeyPath))
	}
}

// RestartGPGAgent restarts the GPG agent
func (y *YubiKeySetup) RestartGPGAgent() {
	y.LogInfo("Restarting GPG agent...")
	
	// Kill existing agent
	y.RunCommand("pkill", "gpg-agent")
	if !y.DryRun {
		time.Sleep(2 * time.Second)
	}
	
	// Verify it's not running
	if err := y.RunCommand("pgrep", "gpg-agent"); err == nil && !y.DryRun {
		y.LogError("GPG agent still running after kill attempt")
	}
	
	// Start fresh agent
	y.RunCommand("gpg-agent", "--daemon")
}

// CheckYubiKeyDetection checks if YubiKey is detected
func (y *YubiKeySetup) CheckYubiKeyDetection() {
	y.LogInfo("Checking YubiKey detection...")
	if !y.DryRun {
		fmt.Println("Please insert your YubiKey if not already inserted and press Enter to continue...")
		bufio.NewReader(os.Stdin).ReadString('\n')
	}
	
	if err := y.RunCommand("gpg2", "--card-status"); err != nil && !y.DryRun {
		y.LogError("YubiKey not detected. Please replug your YubiKey and try again.")
	}
	
	// For Linux, set pinentry
	if y.OS == "linux" {
		if err := y.RunCommand("which", "update-alternatives"); err == nil {
			y.LogInfo("Setting pinentry to pinentry-gnome3...")
			y.RunCommand("sudo", "update-alternatives", "--install", "/usr/bin/pinentry", "pinentry", "/usr/bin/pinentry-gnome3", "1")
			y.RunCommand("sudo", "update-alternatives", "--set", "pinentry", "/usr/bin/pinentry-gnome3")
		}
	}
}

// CheckSSHSupport checks SSH support
func (y *YubiKeySetup) CheckSSHSupport() {
	y.LogInfo("Checking SSH support...")
	output, err := y.RunCommandWithOutput("ssh-add", "-L")
	if err != nil || (!strings.Contains(output, "cardno:") && !y.DryRun) {
		y.LogError("YubiKey card number not found in ssh-add -L output. Please replug YubiKey and restart the script.")
	}
}

// GetPGPKeyID gets the PGP key ID
func (y *YubiKeySetup) GetPGPKeyID() {
	y.LogInfo("Getting PGP Key ID...")
	output, err := y.RunCommandWithOutput("gpg2", "--list-keys", "--keyid-format", "0xlong")
	if err != nil && !y.DryRun {
		y.LogError("Could not list GPG keys")
	}
	
	if y.DryRun {
		y.PGPKeyID = "0x1234567890ABCDEF"
	} else {
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "pub") {
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					keyParts := strings.Split(parts[1], "/")
					if len(keyParts) >= 2 {
						y.PGPKeyID = strings.Fields(keyParts[1])[0]
						break
					}
				}
			}
		}
	}
	
	if y.PGPKeyID == "" && !y.DryRun {
		y.LogError("Could not determine PGP Key ID")
	}
	
	y.LogSuccess(fmt.Sprintf("PGP Key ID: %s", y.PGPKeyID))
}

// SetKeyTrust sets the key trust to ultimate
func (y *YubiKeySetup) SetKeyTrust() {
	y.LogInfo("Setting key trust to ultimate...")
	fmt.Printf("Setting trust level to 5 (ultimate) for key %s\n", y.PGPKeyID)
	
	if y.DryRun {
		y.LogInfo("DRY RUN: Would set key trust to ultimate")
		return
	}
	
	cmd := exec.Command("gpg", "--edit-key", y.PGPKeyID, "--command-fd", "0")
	cmd.Stdin = strings.NewReader("trust\n5\ny\nsave\n")
	if err := cmd.Run(); err != nil {
		// Try with email
		currentUser, _ := user.Current()
		userEmail := fmt.Sprintf("%s@obmondo.com", currentUser.Username)
		y.LogWarning(fmt.Sprintf("Trying with email: %s", userEmail))
		
		cmd = exec.Command("gpg", "--edit-key", userEmail, "--command-fd", "0")
		cmd.Stdin = strings.NewReader("trust\n5\ny\nsave\n")
		if err := cmd.Run(); err != nil {
			y.LogError("Failed to set key trust")
		}
	}
}

// TestEncryption tests encryption and decryption
func (y *YubiKeySetup) TestEncryption() {
	y.LogInfo("Testing encryption and decryption...")
	if !y.DryRun {
		fmt.Println("Please enter your PIN when prompted and touch your YubiKey when it blinks...")
	}
	
	if y.DryRun {
		y.LogInfo("DRY RUN: Would test encryption/decryption")
		y.LogSuccess("Encryption/decryption test passed!")
		return
	}
	
	// Set GPG_TTY
	os.Setenv("GPG_TTY", "/dev/tty")
	
	// Test with key ID
	cmd1 := exec.Command("uname", "-a")
	cmd2 := exec.Command("gpg2", "--encrypt", "--armor", "--recipient", y.PGPKeyID)
	cmd3 := exec.Command("gpg2", "--decrypt")
	
	cmd2.Stdin, _ = cmd1.StdoutPipe()
	cmd3.Stdin, _ = cmd2.StdoutPipe()
	
	cmd1.Start()
	cmd2.Start()
	cmd3.Start()
	
	if err := cmd3.Wait(); err != nil {
		// Try with email
		currentUser, _ := user.Current()
		userEmail := fmt.Sprintf("%s@obmondo.com", currentUser.Username)
		y.LogWarning(fmt.Sprintf("Trying encryption test with email: %s", userEmail))
		
		cmd1 = exec.Command("uname", "-a")
		cmd2 = exec.Command("gpg2", "--encrypt", "--armor", "--recipient", userEmail)
		cmd3 = exec.Command("gpg2", "--decrypt")
		
		cmd2.Stdin, _ = cmd1.StdoutPipe()
		cmd3.Stdin, _ = cmd2.StdoutPipe()
		
		cmd1.Start()
		cmd2.Start()
		cmd3.Start()
		
		if err := cmd3.Wait(); err != nil {
			y.LogError("Encryption/decryption test failed. Make sure pinentry-gnome3 is set and GPG_TTY is exported.")
		}
	}
	
	y.LogSuccess("Encryption/decryption test passed!")
}

// ConfigureGitSigning configures Git signing
func (y *YubiKeySetup) ConfigureGitSigning() {
	y.LogInfo("Configuring Git signing...")
	
	if err := y.RunCommand("which", "git"); err != nil {
		y.LogError("Git is not installed. Please install Git first.")
	}
	
	if err := y.RunCommand("git", "config", "--global", "user.signingkey", y.PGPKeyID); err != nil {
		y.LogError("Failed to set Git signing key")
	}
	
	if err := y.RunCommand("git", "config", "--global", "commit.gpgsign", "true"); err != nil {
		y.LogError("Failed to enable Git commit signing")
	}
	
	y.LogSuccess("Git signing configured successfully!")
}

// ShowFinalInstructions shows final setup instructions
func (y *YubiKeySetup) ShowFinalInstructions() {
	fmt.Println()
	fmt.Println(green("╔══════════════════════════════════════════════════════════════════════════════╗"))
	fmt.Println(green("║                          Setup Complete!                                    ║"))
	fmt.Println(green("╚══════════════════════════════════════════════════════════════════════════════╝"))
	fmt.Println()
	
	y.LogInfo("Next steps:")
	fmt.Println("1. Add your SSH public key to Gitea:")
	fmt.Println("   - Go to Gitea -> Settings -> SSH/GPG Keys -> Manage SSH Keys")
	fmt.Println("   - Use this key:")
	fmt.Println()
	
	if output, err := y.RunCommandWithOutput("ssh-add", "-L"); err == nil {
		lines := strings.Split(output, "\n")
		if len(lines) > 0 {
			fmt.Println(lines[0])
		}
	}
	
	fmt.Println()
	fmt.Println("2. Add your GPG public key to Gitea:")
	fmt.Println("   - Go to Gitea -> Settings -> SSH/GPG Keys -> Manage GPG Keys")
	fmt.Println("   - Use this key:")
	fmt.Println()
	
	if output, err := y.RunCommandWithOutput("gpg", "--export", "-a", y.PGPKeyID); err == nil {
		lines := strings.Split(output, "\n")
		for i, line := range lines {
			if i < 10 {
				fmt.Println(line)
			}
		}
		fmt.Println("   [... truncated for display ...]")
	}
	
	fmt.Println()
	fmt.Printf("3. Source your shell configuration:\n")
	fmt.Printf("   source %s\n", y.ShellRC)
	fmt.Println()
	fmt.Println("4. Test Git signing in a repository:")
	fmt.Println("   git commit -S -m 'test commit'")
	fmt.Println()
	fmt.Println("If you encounter any issues with Gitea, please contact your team.")
	fmt.Println()
	y.LogSuccess(fmt.Sprintf("Setup completed successfully! Error tracking ID: %s", y.ErrorUUID))
}

// Run executes the complete setup process
func (y *YubiKeySetup) Run() {
	y.Banner()
	y.CheckOSCompatibility()
	y.CheckYubiKeyPIN()
	y.GetGPGKeyPath()
	y.InstallPrerequisites()
	y.CreateGPGConfigs()
	y.ConfigureGPGAgent()
	y.ConfigureShellEnvironment()
	y.UnsetSSHAuthSock()
	y.ImportGPGKey()
	y.RestartGPGAgent()
	y.CheckYubiKeyDetection()
	y.CheckSSHSupport()
	y.GetPGPKeyID()
	y.SetKeyTrust()
	y.TestEncryption()
	y.ConfigureGitSigning()
	y.ShowFinalInstructions()
}