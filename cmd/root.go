package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version of the application
	Version = "1.0.0"
	
	// Global flags
	verbose bool
	dryRun  bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "yubikey-setup",
	Short: "YubiKey Employee Setup Tool for EnableIT",
	Long: `A comprehensive YubiKey setup tool for EnableIT employees that automates:

- OS compatibility checking (Linux/macOS)
- GPG configuration and key management
- YubiKey detection and validation
- SSH authentication setup
- Git signing configuration
- Shell environment configuration

Based on: https://gitea.obmondo.com/EnableIT/wiki/src/branch/master/internal/yubikey-employee-setup.md`,
	Version: Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "show what would be done without making changes")
	
	// Set version template
	rootCmd.SetVersionTemplate(`{{printf "%s version %s\n" .Name .Version}}`)
}