// cmd/env.go
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var macInstall bool

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Environment setup and diagnostics",
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for required tools and AWS credentials",
	Run: func(cmd *cobra.Command, args []string) {
		check("terraform")
		check("ansible-playbook")
		check("aws")

		// Check AWS creds
		if os.Getenv("AWS_ACCESS_KEY_ID") != "" || os.Getenv("AWS_PROFILE") != "" {
			fmt.Println("✅ AWS credentials found.")
		} else {
			fmt.Println("❌ AWS credentials or profile not found. Set AWS_ACCESS_KEY_ID or AWS_PROFILE.")
		}
	},
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Guide user to install missing tools and configure AWS",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("⚙️ Setup guidance:")
		fmt.Println("- Install Terraform: https://developer.hashicorp.com/terraform/install")
		fmt.Println("- Install Ansible: https://docs.ansible.com/ansible/latest/installation_guide/")
		fmt.Println("- Install AWS CLI: https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html")
		fmt.Println("- Configure AWS credentials: aws configure OR export AWS_PROFILE=your-profile")

		if macInstall {
			fmt.Println("\n🍎 Installing missing tools using Homebrew...")
			installIfMissing("terraform")
			installIfMissing("ansible")
			installIfMissing("awscli")
		}
	},
}

func init() {
	setupCmd.Flags().BoolVar(&macInstall, "mac", false, "Install tools using Homebrew on macOS")
	envCmd.AddCommand(checkCmd)
	envCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(envCmd)
}

func check(tool string) {
	_, err := exec.LookPath(tool)
	if err != nil {
		fmt.Printf("❌ %s not found in PATH\n", tool)
	} else {
		fmt.Printf("✅ %s is installed\n", tool)
	}
}

func installIfMissing(tool string) {
	_, err := exec.LookPath(tool)
	if err != nil {
		fmt.Printf("📦 Installing %s via Homebrew...\n", tool)
		installBrewPackage(tool)
	} else {
		fmt.Printf("✅ %s already installed. Skipping.\n", tool)
	}
}

func installBrewPackage(pkg string) {
	cmd := exec.Command("brew", "install", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Failed to install %s via Homebrew: %v\n", pkg, err)
	} else {
		fmt.Printf("✅ Installed %s\n", pkg)
	}
}
