#
# Lambda
#
resource "aws_lambda_function" "main" {
  function_name = var.blank_name
  description   = var.description
  tags          = var.tags

  role         = var.create_iam_role ? aws_iam_role.main[0].arn : var.iam_role_arn
  image_uri    = var.image_uri
  package_type = "Image"
  timeout      = var.timeout
  memory_size  = var.memory_size

  vpc_config {
    subnet_ids         = var.subnet_ids
    security_group_ids = var.security_group_ids
  }

  environment {
    variables = merge(var.environ, { "MODE" : "cloud", "FORMATTER_TYPE" : "JSON" })
  }
}

#
# Cloudwatch group
#
resource "aws_cloudwatch_log_group" "main" {
  name = "/aws/lambda/${var.blank_name}"
  tags = var.tags

  retention_in_days = var.log_retention_in_days
  depends_on        = [aws_lambda_function.main]
}
