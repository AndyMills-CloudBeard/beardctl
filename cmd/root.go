// cmd/root.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "beardctl",
	Short: "A CLI tool to deploy Terraform infrastructure and configure EC2",
	Long:  `beardctl wraps Terraform and Ansible to deploy and manage secure infrastructure in AWS.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
