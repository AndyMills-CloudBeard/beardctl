// internal/terraform/output.go
package terraform

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// OutputResult represents the structure of each Terraform output value
type OutputResult struct {
	Value interface{} `json:"value"`
}

// ParseOutputs runs `terraform output -json` in the given path and returns a flat map
func ParseOutputs(tfDir string) (map[string]string, error) {
	cmd := exec.Command("terraform", "output", "-json")
	cmd.Dir = tfDir
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run terraform output: %w", err)
	}

	// Parse JSON output
	parsed := map[string]OutputResult{}
	if err := json.Unmarshal(output, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse terraform output: %w", err)
	}

	// Flatten the map
	results := make(map[string]string)
	for k, v := range parsed {
		results[k] = fmt.Sprintf("%v", v.Value)
	}

	return results, nil
}
