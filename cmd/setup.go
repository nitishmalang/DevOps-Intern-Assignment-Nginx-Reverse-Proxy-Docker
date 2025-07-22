package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	gpgKeyPath string
	skipPin    bool
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Run the complete YubiKey setup process",
	Long: `Runs the complete YubiKey setup process including:

- OS compatibility check
- Prerequisites installation
- GPG configuration
- YubiKey detection and setup
- SSH authentication configuration
- Git signing setup

This is the main command that performs the full automated setup.`,
	Run: runSetup,
}

func runSetup(cmd *cobra.Command, args []string) {
	// Import the main setup logic
	setup := NewYubiKeySetup()
	
	if verbose {
		fmt.Println("Running in verbose mode...")
	}
	
	if dryRun {
		fmt.Println("🚨 DRY RUN MODE - No changes will be made")
		setup.DryRun = true
	}
	
	// Override GPG key path if provided
	if gpgKeyPath != "" {
		// Expand tilde
		if gpgKeyPath[0] == '~' {
			home, _ := os.UserHomeDir()
			gpgKeyPath = filepath.Join(home, gpgKeyPath[1:])
		}
		setup.GPGKeyPath = gpgKeyPath
	}
	
	// Skip PIN check if requested
	if skipPin {
		setup.SkipPinCheck = true
	}
	
	// Run the setup
	setup.Run()
}

func init() {
	rootCmd.AddCommand(setupCmd)
	
	// Setup-specific flags
	setupCmd.Flags().StringVarP(&gpgKeyPath, "gpg-key", "k", "", "path to GPG public key file")
	setupCmd.Flags().BoolVar(&skipPin, "skip-pin", false, "skip PIN verification (not recommended)")
}