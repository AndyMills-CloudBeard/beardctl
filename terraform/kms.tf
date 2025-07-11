resource "aws_kms_key" "main" {
  description         = "Key for EBS and VPC flow logs"
  enable_key_rotation = true

  tags = {
    Name        = "beardctl-kms-key"
    Environment = "mgmt"
  }
}
