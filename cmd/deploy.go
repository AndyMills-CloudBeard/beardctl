// cmd/deploy.go
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"beardctl/internal/ansible"
	"beardctl/internal/terraform"

	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy [path]",
	Short: "Deploy infrastructure with Terraform and configure EC2",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tfPath := args[0]
		fmt.Printf("🚀 Deploying with: %s\n", tfPath)

		if err := runTerraform(tfPath); err != nil {
			fmt.Printf("❌ Terraform failed: %v\n", err)
			os.Exit(1)
		}

		outputs, err := terraform.ParseOutputs(tfPath)
		if err != nil {
			fmt.Printf("❌ Failed to parse Terraform outputs: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\n✅ Terraform completed. Outputs:")
		for k, v := range outputs {
			fmt.Printf("  %s = %s\n", k, v)
		}

		// Extract first public IP address
		var ec2IP string
		if rawIPs, ok := outputs["ec2_public_ip"]; ok {
			clean := strings.Trim(rawIPs, "[]")
			ec2IP = strings.Split(clean, " ")[0]
			fmt.Printf("\n🌐 EC2 Public IP: %s\n", ec2IP)
		} else {
			fmt.Println("❌ Missing 'ec2_public_ip' output from Terraform")
			os.Exit(1)
		}

		// Extract ALB DNS, S3 Bucket, and Instance ID
		albDNS := outputs["alb_dns"]
		s3Bucket := outputs["s3_bucket"]
		instanceID := strings.Trim(outputs["instance_id"], "[]")

		// Wait for SSH to be ready
		if err := ansible.WaitForSSH(ec2IP, 2*time.Minute); err != nil {
			fmt.Printf("❌ SSH not ready: %v\n", err)
			os.Exit(1)
		}

		// Write inventory.ini and run playbook
		sshKeyPath := filepath.Join(os.Getenv("HOME"), ".ssh", "beardctl-key")
		if err := ansible.WriteInventory(ec2IP, sshKeyPath); err != nil {
			fmt.Printf("❌ Failed to write inventory: %v\n", err)
			os.Exit(1)
		}

		if err := ansible.RunPlaybook("playbook/nginx.yml"); err != nil {
			fmt.Printf("❌ Ansible playbook failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✅ EC2 nginx configuration complete with Ansible.")

		// Run monitor playbook with extra vars
		extraVars := fmt.Sprintf("alb_dns=%s instance_id=%s s3_bucket=%s", albDNS, instanceID, s3Bucket)
		if err := ansible.RunPlaybookWithVars("playbook/monitor.yml", extraVars); err != nil {
			fmt.Printf("❌ Ansible playbook failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✅ EC2 monitor app install with Ansible complete.")
		fmt.Printf("\n🌍 Deployment complete! Visit your site at: http://%s\n", albDNS)
	},
}

func runTerraform(path string) error {
	cmd := exec.Command("terraform", "init")
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("terraform", "apply", "-auto-approve")
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
