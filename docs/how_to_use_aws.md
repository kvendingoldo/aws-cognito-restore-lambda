## How to use it within AWS

1. Lambda image should be pulled from docker hub and pushed into your personal ECR repository; AWS Lambda is not able to
   work with any other docker registry except ECR.
2. Apply TF module into your infrastructure

### How to trigger lambda manually via UI

1. Go to Lambda function that has been created via Terraform -> Tests
2. Fill "Test Event" and click "Test"

```
{
  "awsRegion": "<AWS_REGION | optional>",
  
  "cognitoUserPoolId": "<COGNITO_USER_POOL_ID>",
  "cognitoRegion": "<COGNITO_POOL_AWS_REGION | AWS_REGION will be used if this var is omitted>",
  
  "s3BucketName": "<S3_BUCKET_NAME>",
  "s3BucketRegion": "<S3_AWS_REGION | AWS_REGION will be used if this var is omitted>",
  
  "backupDirPath": "<backup path inside of S3 bucket>",
  
  "restoreUsers": <true | false>,
  "restoreGroups": <true | false>,
  "cleanUpBeforeRestore": <true | false>
}
```

#### Example #1:

```json
{
  "awsRegion": "us-west-2",
  "cognitoUserPoolId": "ap-southeast-2_EPyUfpQq7",
  "s3BucketName": "mybuckettest",
  "backupDirPath": "2022-08-26T20:16:37Z/",
  "restoreUsers": true,
  "restoreGroups": true,
  "cleanUpBeforeRestore": true
}
```