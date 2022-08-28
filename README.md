# aws-cognito-restore-lambda

## Overview
Users of AWS Cognito periodically face issues with backups. The Cognito user pool backup mechanism is not embedded into AWS.
We may use this Lambda to restore users from backups created by the [aws-cognito-backup-lambda](https://github.com/kvendingoldo/aws-cognito-backup-lambda) project.

## Documentation
You can review the following documents on the Lambda to learn more:
* [How to use the Lambda inside of AWS](docs/how_to_use_aws.md)
* [How to use the Lambda locally](docs/how_to_use_locally.md)
* [How to use Terraform automation](docs/how_to_use_terraform.md)
* [Labmda's environment variables](docs/environment_variables.md)