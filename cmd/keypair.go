// cmd/keypair.go
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var keyName string

var keypairCmd = &cobra.Command{
	Use:   "setup-keypair",
	Short: "Generate SSH keypair and upload public key to AWS",
	Run: func(cmd *cobra.Command, args []string) {
		keyPath := filepath.Join(os.Getenv("HOME"), ".ssh", keyName)
		pubKeyPath := keyPath + ".pub"

		if _, err := os.Stat(keyPath); err == nil {
			fmt.Printf("🔐 SSH key '%s' already exists, skipping generation.\n", keyPath)
		} else {
			fmt.Println("🔧 Generating SSH keypair...")
			gen := exec.Command("ssh-keygen", "-t", "rsa", "-b", "4096", "-f", keyPath, "-C", keyName, "-N", "")
			gen.Stdout = os.Stdout
			gen.Stderr = os.Stderr
			if err := gen.Run(); err != nil {
				fmt.Printf("❌ Failed to generate SSH key: %v\n", err)
				os.Exit(1)
			}
		}

		fmt.Println("☁️ Uploading public key to AWS EC2 key pairs...")
		upload := exec.Command("aws", "ec2", "import-key-pair",
			"--key-name", keyName,
			"--public-key-material", fmt.Sprintf("fileb://%s", pubKeyPath))
		upload.Stdout = os.Stdout
		upload.Stderr = os.Stderr
		if err := upload.Run(); err != nil {
			fmt.Printf("❌ Failed to upload SSH public key to AWS: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✅ Keypair setup complete.")
	},
}

func init() {
	keypairCmd.Flags().StringVar(&keyName, "name", "beardctl-key", "Name of the SSH keypair")
	envCmd.AddCommand(keypairCmd)
}
