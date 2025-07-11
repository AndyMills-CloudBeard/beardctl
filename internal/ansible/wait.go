// internal/ansible/wait.go
package ansible

import (
	"fmt"
	"net"
	"time"
)

// WaitForSSH tries to connect to port 22 on the EC2 instance until it's ready
func WaitForSSH(ip string, timeout time.Duration) error {
	fmt.Printf("⏳ Waiting for SSH to become available on %s:22...\n", ip)
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:22", ip), 3*time.Second)
		if err == nil {
			_ = conn.Close()
			fmt.Println("✅ SSH is now available.")
			return nil
		}
		fmt.Print(".")
		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("timeout waiting for SSH on %s", ip)
}
