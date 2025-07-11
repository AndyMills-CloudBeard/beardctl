// cmd/destroy.go
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var destroyCmd = &cobra.Command{
	Use:   "destroy [path]",
	Short: "Tear down Terraform-managed infrastructure",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tfPath := args[0]
		fmt.Printf("🧨 Destroying infrastructure in: %s\n", tfPath)

		if err := runTerraformDestroy(tfPath); err != nil {
			fmt.Printf("❌ Terraform destroy failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✅ Infrastructure destroyed successfully.")
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
}

func runTerraformDestroy(path string) error {
	cmd := exec.Command("terraform", "init")
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("terraform", "destroy", "-auto-approve")
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
