## How to use it locally

1. Look at [environments variables](environment_variables.md) and set at least required variables

### example:
```shell
export AWS_REGION="us-east-2"
export MODE=local
export COGNITO_USER_POOL_ID="ap-southeast-2_EPyUfpQq7"
export S3_BUCKET_NAME="mybuckettest"
export BACKUP_DIR_PATH="2022-08-26T20:16:37Z"
export RESTORE_USERS=true 
export CLEANUP_BEFORE_RESTORE=true 
```

2. Run lambda locally

```sh
go run main.go
```

