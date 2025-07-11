locals {
  public_subnets = [
    "10.0.1.0/24",
    "10.0.3.0/24"
  ]

  private_subnets = [
    "10.0.2.0/24",
    "10.0.4.0/24"
  ]

  public_subnet_tags = {
    0 = "public"
    1 = "public"
  }

  private_subnet_tags = {
    0 = "private"
    1 = "private"
  }
}

locals {
  kms_key_arn = aws_kms_key.main.arn
}
