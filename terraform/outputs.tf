output "al2023_ami_id" {
  value = data.aws_ami.ami.id
}

output "kms_key_arn" {
  value = aws_kms_key.main.arn
}

output "ec2_public_ip" {
  value = module.ec2_beard.public_ips
}

output "alb_dns" {
  description = "The DNS name of the ALB"
  value       = aws_lb.http.dns_name
}

output "s3_bucket" {
  value       = module.s3_bucket_errorlog.id
  description = "The name of the error log S3 bucket"
}

output "instance_id" {
  description = "The AWS Instance id created"
  value       = module.ec2_beard.instance_id
}
