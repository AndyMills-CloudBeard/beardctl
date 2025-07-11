data "aws_iam_policy_document" "s3_put_object" {
  statement {
    effect = "Allow"
    actions = [
      "s3:PutObject"
    ]
    resources = [
      "${module.s3_bucket_errorlog.arn}/*"
    ]
  }

  statement {
    effect = "Allow"
    actions = [
      "kms:GenerateDataKey",
      "kms:Encrypt",
      "kms:Decrypt",
      "kms:DescribeKey"
    ]
    resources = [
      local.kms_key_arn
    ]
  }
}

resource "aws_iam_policy" "s3_put_object" {
  name   = "beardctl-s3-put-object"
  policy = data.aws_iam_policy_document.s3_put_object.json
}

data "aws_iam_policy_document" "ec2_self_reboot" {
  statement {
    effect = "Allow"
    actions = [
      "ec2:RebootInstances"
    ]
    resources = [
      "arn:aws:ec2:us-east-1:${data.aws_caller_identity.current.account_id}:instance/*"
    ]
    condition {
      test     = "StringEquals"
      variable = "ec2:ResourceTag/Name"
      values   = [var.instance_name]
    }
  }
}

resource "aws_iam_policy" "ec2_self_reboot" {
  name   = "beardctl-ec2-self-reboot"
  policy = data.aws_iam_policy_document.ec2_self_reboot.json
}
