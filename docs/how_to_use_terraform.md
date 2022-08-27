#### Terraform code example

1. Add module execution to your TF code

```terraform
module "cognito_restore_lambda" {
  source = "git@github.com:kvendingoldo/aws-cognito-restore-lambda.git//files/terraform/module?ref=rc/0.1.0"

  blank_name = "${module.naming.common_name}-cognito-restore"
  tags       = local.tags

  image_uri = var.cognito_restore_lambda_image_uri
}
```

2. Specify variables

```terraform
variable "cognito_restore_lambda_image_uri" {
  default = "<YOUR_ACCOUNT_ID>.dkr.ecr.us-east-2.amazonaws.com/aws-cognito_restore_lambda:<VERSION>"
}
```