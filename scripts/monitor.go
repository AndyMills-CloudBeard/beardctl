// File: monitor.go
package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	checkInterval   = 5 * time.Second
	maxFailures     = 3
	expectedContent = "Welcome to the Beard"
	logFilePath     = "/var/log/beardctl-monitor.log"
)

var (
	instanceID string
	bucketName string
	albDNS     string
)

func logToFile(msg string) {
	f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening log file: %v", err)
		return
	}
	defer f.Close()
	logger := log.New(f, "", log.LstdFlags)
	logger.Println(msg)
}

func rebootInstance(cfg aws.Config) {
	ec2Client := ec2.NewFromConfig(cfg)
	_, err := ec2Client.RebootInstances(context.TODO(), &ec2.RebootInstancesInput{
		InstanceIds: []string{instanceID},
	})
	if err != nil {
		logToFile(fmt.Sprintf("❌ Failed to reboot instance %s: %v", instanceID, err))
	} else {
		logToFile(fmt.Sprintf("🔁 Rebooted instance %s after repeated failures.", instanceID))
	}
}

func uploadLogToS3(cfg aws.Config) {
	file, err := os.Open(logFilePath)
	if err != nil {
		log.Printf("Error opening log file for upload: %v", err)
		return
	}
	defer file.Close()

	s3Client := s3.NewFromConfig(cfg)
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fmt.Sprintf("logs/beardctl-monitor-%d.log", time.Now().Unix())),
		Body:   file,
	})
	if err != nil {
		logToFile(fmt.Sprintf("❌ Failed to upload log to S3: %v", err))
	} else {
		logToFile("📤 Log successfully uploaded to S3.")
	}
}

func checkHealth() bool {
	resp, err := http.Get("http://" + albDNS)
	if err != nil {
		logToFile(fmt.Sprintf("❌ HTTP request failed: %v", err))
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logToFile(fmt.Sprintf("❌ Unexpected status code: %d", resp.StatusCode))
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logToFile(fmt.Sprintf("❌ Failed to read response body: %v", err))
		return false
	}

	if !bytes.Contains(body, []byte(expectedContent)) {
		logToFile("❌ Expected content not found in response.")
		return false
	}

	return true
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: ./monitor <alb_dns> <instance_id> <s3_bucket>")
		os.Exit(1)
	}

	albDNS = os.Args[1]
	instanceID = os.Args[2]
	bucketName = os.Args[3]

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("Unable to load AWS config: %v", err)
	}

	failureCount := 0

	for {
		healthy := checkHealth()
		if !healthy {
			failureCount++
			if failureCount >= maxFailures {
				logToFile("⚠️ Detected 3 consecutive failures. Attempting reboot...")
				rebootInstance(cfg)
				uploadLogToS3(cfg)
				failureCount = 0
			}
		} else {
			failureCount = 0
		}

		time.Sleep(checkInterval)
	}
}
