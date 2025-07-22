package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check system prerequisites and YubiKey status",
	Long: `Performs system checks to validate that all prerequisites are met:

- Operating system compatibility
- Required packages installation
- GPG configuration status
- YubiKey detection
- SSH agent status

This command is useful for troubleshooting and verifying the setup.`,
	Run: runCheck,
}

func runCheck(cmd *cobra.Command, args []string) {
	fmt.Println(color.BlueString("🔍 System Prerequisites Check"))
	fmt.Println(color.BlueString("=============================="))
	fmt.Println()
	
	checkOS()
	checkPrerequisites()
	checkGPGSetup()
	checkYubiKey()
	checkSSHAgent()
}

func checkOS() {
	fmt.Print("Operating System: ")
	switch runtime.GOOS {
	case "linux":
		fmt.Println(color.GreenString("✅ Linux (supported)"))
	case "darwin":
		fmt.Println(color.GreenString("✅ macOS (supported)"))
	default:
		fmt.Println(color.RedString("❌ %s (not supported)", runtime.GOOS))
	}
}

func checkPrerequisites() {
	fmt.Println("\nPrerequisites:")
	
	packages := map[string][]string{
		"linux":  {"gpg", "gpg2", "gpg-agent", "pinentry-gnome3"},
		"darwin": {"gpg", "gpg2", "gpg-agent", "pinentry-mac"},
	}
	
	osPackages := packages[runtime.GOOS]
	if osPackages == nil {
		fmt.Println(color.YellowString("⚠️  Cannot check prerequisites for this OS"))
		return
	}
	
	for _, pkg := range osPackages {
		if checkCommand(pkg) {
			fmt.Printf("  %s %s\n", color.GreenString("✅"), pkg)
		} else {
			fmt.Printf("  %s %s\n", color.RedString("❌"), pkg)
		}
	}
}

func checkGPGSetup() {
	fmt.Println("\nGPG Configuration:")
	
	// Check .gnupg directory
	home := getUserHome()
	gnupgDir := home + "/.gnupg"
	
	if fileExists(gnupgDir) {
		fmt.Printf("  %s .gnupg directory\n", color.GreenString("✅"))
	} else {
		fmt.Printf("  %s .gnupg directory\n", color.RedString("❌"))
	}
	
	// Check config files
	configs := []string{"gpg.conf", "gpg-agent.conf", "dirmngr.conf"}
	for _, config := range configs {
		path := gnupgDir + "/" + config
		if fileExists(path) {
			fmt.Printf("  %s %s\n", color.GreenString("✅"), config)
		} else {
			fmt.Printf("  %s %s\n", color.YellowString("⚠️"), config)
		}
	}
}

func checkYubiKey() {
	fmt.Println("\nYubiKey Status:")
	
	// Check card status
	if err := exec.Command("gpg2", "--card-status").Run(); err == nil {
		fmt.Printf("  %s YubiKey detected\n", color.GreenString("✅"))
	} else {
		fmt.Printf("  %s YubiKey not detected\n", color.RedString("❌"))
		return
	}
	
	// Check if keys are available
	output, err := exec.Command("gpg", "--list-secret-keys").Output()
	if err == nil && len(output) > 0 {
		fmt.Printf("  %s GPG secret keys found\n", color.GreenString("✅"))
	} else {
		fmt.Printf("  %s No GPG secret keys\n", color.YellowString("⚠️"))
	}
}

func checkSSHAgent() {
	fmt.Println("\nSSH Agent:")
	
	// Check ssh-add
	output, err := exec.Command("ssh-add", "-L").Output()
	if err == nil && len(output) > 0 {
		fmt.Printf("  %s SSH agent running with keys\n", color.GreenString("✅"))
		
		// Check for YubiKey card number
		if contains(string(output), "cardno:") {
			fmt.Printf("  %s YubiKey SSH key detected\n", color.GreenString("✅"))
		} else {
			fmt.Printf("  %s No YubiKey SSH key\n", color.YellowString("⚠️"))
		}
	} else {
		fmt.Printf("  %s SSH agent not running or no keys\n", color.RedString("❌"))
	}
}

// Helper functions
func checkCommand(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func fileExists(path string) bool {
	_, err := exec.Command("test", "-f", path).Output()
	return err == nil
}

func getUserHome() string {
	output, err := exec.Command("sh", "-c", "echo $HOME").Output()
	if err != nil {
		return ""
	}
	return string(output)[:len(output)-1] // Remove newline
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr || 
		   len(s) > len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func init() {
	rootCmd.AddCommand(checkCmd)
}