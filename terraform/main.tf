# EC2 instance
module "ec2_beard" {
  #   source = "github.com/Coalfire-CF/terraform-aws-ec2?ref=v2.0.12"
  source = "github.com/AndyMills-CloudBeard/terraform-aws-ec2"

  name = var.instance_name

  ami               = data.aws_ami.ami.id
  ec2_instance_type = var.ec2_instance_type
  instance_count    = var.instance_count

  iam_policies = [
    aws_iam_policy.s3_put_object.arn,
    aws_iam_policy.ec2_self_reboot.arn
  ]

  associate_public_ip = true

  vpc_id     = aws_vpc.main.id
  subnet_ids = [aws_subnet.public[0].id]

  ec2_key_pair    = var.key_name
  ebs_kms_key_arn = local.kms_key_arn

  # Storage
  root_volume_size = var.instance_volume_size

  # Security Group Rules
  ingress_rules = {
    "alb_http" = {
      ip_protocol                  = "tcp"
      from_port                    = "80"
      to_port                      = "80"
      referenced_security_group_id = aws_security_group.alb_sg.id
      description                  = "Allow HTTP from ALB"
    }

    "remote_ssh" = {
      ip_protocol = "tcp"
      from_port   = "22"
      to_port     = "22"
      cidr_ipv4   = "0.0.0.0/0"
      description = "Allow HTTP from ALB"
    }
  }

  egress_rules = {
    "allow_all_egress" = {
      ip_protocol = "-1"
      from_port   = "0"
      to_port     = "0"
      cidr_ipv4   = "0.0.0.0/0"
      description = "Allow all egress"
    }
  }

  # Tagging
  global_tags = {}
}

resource "random_id" "bucket_suffix" {
  byte_length = 4
}

# Create error log bucket
module "s3_bucket_errorlog" {
  source = "github.com/Coalfire-CF/terraform-aws-s3?ref=v1.0.4"

  name                                 = "beardctl-s3-errorlog-bucket-${random_id.bucket_suffix.hex}"
  enable_lifecycle_configuration_rules = false
  enable_kms                           = true
  enable_server_side_encryption        = true
  kms_master_key_id                    = local.kms_key_arn
}
