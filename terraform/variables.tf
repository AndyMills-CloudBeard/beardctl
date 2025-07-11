variable "instance_name" {
  description = "The name of the ec2 instance"
  type        = string
}

variable "ec2_instance_type" {
  description = "The type of the ec2 instance"
  type        = string
}

variable "instance_count" {
  description = "The type of the ec2 instance"
  type        = number
}

variable "key_name" {
  description = "The ssh key name for the ec2 instance"
  type        = string
}

variable "instance_volume_size" {
  description = "The size of the ec2 instance ebs volume"
  type        = number
}

variable "vpc_name" {
  description = "The name of the VPC"
  type        = string
}

variable "cidr" {
  description = "The CIDR range of the VPC"
  type        = string
}
