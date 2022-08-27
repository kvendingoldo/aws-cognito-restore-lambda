resource "aws_iam_role" "main" {
  count = var.create_iam_role ? 1 : 0

  name               = var.blank_name
  description        = "IAM role for for Lambda ${var.blank_name}"
  tags               = var.tags
  assume_role_policy = <<-POLICY
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Principal": {
          "Service": "lambda.amazonaws.com"
        },
        "Effect": "Allow",
        "Sid": ""
      }
    ]
  }
  POLICY
}

#
# VPC permissions
#
resource "aws_iam_role_policy_attachment" "vpc_permissions" {
  count      = (length(var.subnet_ids) != 0 && var.create_iam_role) ? 1 : 0

  role       = aws_iam_role.main[0].name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

#
# Logging policy
#
resource "aws_iam_policy" "logging" {
  count = var.create_iam_role ? 1 : 0

  name        = format("%s-%s", var.blank_name, "logging")
  path        = "/"
  description = "IAM policy for logging from a lambda"

  policy = <<-POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*",
      "Effect": "Allow"
    }
  ]
}
POLICY
}
resource "aws_iam_role_policy_attachment" "logging" {
  count = var.create_iam_role ? 1 : 0

  role       = aws_iam_role.main[0].name
  policy_arn = aws_iam_policy.logging[0].arn
}

#
# Cognito policy
#
resource "aws_iam_policy" "cognito" {
  count = var.create_iam_role ? 1 : 0

  name        = format("%s-%s", var.blank_name, "cognito")
  path        = "/"
  description = "IAM policy for working with Cognito from the Lambda"

  policy = <<-POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "cognito-identity:*",
        "cognito-idp:*",
        "cognito-sync:*"
      ],
      "Resource": "*",
      "Effect": "Allow"
    }
  ]
}
POLICY
}
resource "aws_iam_role_policy_attachment" "cognito" {
  count = var.create_iam_role ? 1 : 0

  role       = aws_iam_role.main[0].name
  policy_arn = aws_iam_policy.cognito[0].arn
}

#
# S3 policy
#
resource "aws_iam_policy" "s3" {
  count = var.create_iam_role ? 1 : 0

  name        = format("%s-%s", var.blank_name, "s3")
  path        = "/"
  description = "IAM policy for working with S3 from the Lambda"

  policy = <<-POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "s3:*"
      ],
      "Resource": "*",
      "Effect": "Allow"
    }
  ]
}
POLICY
}
resource "aws_iam_role_policy_attachment" "route53" {
  count = var.create_iam_role ? 1 : 0

  role       = aws_iam_role.main[0].name
  policy_arn = aws_iam_policy.s3[0].arn
}


