package ansible

import (
	"fmt"
	"os"
	"os/exec"
)

// RunPlaybookWithVars runs an Ansible playbook with extra vars
func RunPlaybookWithVars(playbook string, extraVars string) error {
	cmd := exec.Command("ansible-playbook", "-i", "inventory.ini", playbook, "--extra-vars", extraVars)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Printf("📦 Running Ansible playbook: %s with vars: %s\n", playbook, extraVars)
	return cmd.Run()
}
