#!/bin/bash

echo "🧹 Cleaning up beardctl deployment..."

# Step 1: Find and empty the S3 bucket
bucket=$(aws s3api list-buckets --query "Buckets[].Name" --output text | tr '\t' '\n' | grep '^beardctl-s3-errorlog-bucket-' | head -n1)

if [ -z "$bucket" ]; then
  echo "❌ No beardctl error log bucket found."
else
  echo "🪣 Target bucket: $bucket"
  aws s3api list-object-versions --bucket "$bucket" --output json | \
    jq -r '.Versions[]?, .DeleteMarkers[]? | [.Key, .VersionId] | @tsv' | \
    while IFS=$'\t' read -r key version; do
      aws s3api delete-object --bucket "$bucket" --key "$key" --version-id "$version"
    done
fi

# Step 2: Destroy Terraform infrastructure
echo "☠️  Destroying infrastructure..."
beardctl destroy terraform

# Step 3: Remove beardctl binary
echo "🗑️  Removing beardctl binary..."
sudo rm -f /usr/local/bin/beardctl

# Step 4: Remove SSH key files
echo "🔐 Removing local SSH keypair..."
sudo rm -f /Users/$(whoami)/.ssh/beardctl-key.pub
sudo rm -f /Users/$(whoami)/.ssh/beardctl-key

# Step 5: Delete AWS key pair
echo "🔑 Deleting AWS key-pair..."
aws ec2 delete-key-pair --key-name beardctl-key

echo "Removing GitHub Repo..."
cd .. && rm -rf beardctl

echo "✅ Cleanup complete!"
