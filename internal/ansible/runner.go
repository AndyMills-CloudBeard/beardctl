// internal/ansible/runner.go
package ansible

import (
	"fmt"
	"os"
	"os/exec"
)

// WriteInventory creates an inventory.ini file from the EC2 public/private IP
func WriteInventory(ip string, privateKeyPath string) error {
	content := fmt.Sprintf(`[web]
%s ansible_user=ec2-user ansible_ssh_private_key_file=%s ansible_ssh_extra_args='-o StrictHostKeyChecking=no'
`, ip, privateKeyPath)
	return os.WriteFile("inventory.ini", []byte(content), 0644)
}

// RunPlaybook executes the Ansible playbook using the generated inventory
func RunPlaybook(playbookPath string) error {
	cmd := exec.Command("ansible-playbook", "-i", "inventory.ini", playbookPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("📦 Running Ansible playbook...")
	return cmd.Run()
}
