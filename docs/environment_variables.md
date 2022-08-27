## Environment variables

* FORMATTER_TYPE
    * Description: formatter type for Lambda's logs
    * Possible values: JSON | TEXT
    * Required: no

* MODE
    * Description: application mode
    * Possible values: cloud | local
    * Required: yes

* LOG_LEVEL
    * Description: Lambda's log level
    * Possible values: panic|fatal|error|warn|info|debug|trace
    * Required: no

* AWS_REGION
    * Description: Default AWS Region. Inside of Lambda it's setting automatically by AWS
    * Possible values: <any valid AWS region>
    * Required: yes

* COGNITO_USER_POOL_ID
    * Description: Cognito user pool ID
    * Possible values: <any valid ID>
    * Required: yes

* COGNITO_REGION
    * Description: AWS Region for Cognito client; If this variable is omitted value from AWS REGION will be used
    * Possible values: <any valid AWS region>
    * Required: no

* S3_BUCKET_NAME
    * Description: S3 bucket name
    * Possible values: <any valid bucket name>
    * Required: no

* S3_BUCKET_REGION
    * Description: AWS Region for S3 client; If this variable is omitted value from AWS REGION will be used
    * Possible values: <any valid AWS region>
    * Required: no

* BACKUP_DIR_PATH
    * Description: Backup path directory inside S3 bucket
    * Possible values: <any string>
    * Required: yes

* RESTORE_USERS
    * Description: Should users be restored from backup?
    * Possible values: true|false
    * Required: no

* RESTORE_GROUPS
    * Description: Should groups be restored from backup?
    * Possible values: true|false
    * Required: no

* CLEANUP_BEFORE_RESTORE
    * Description: Should Cognito user pool be cleaned up before restore?
    * Possible values: true|false
    * Required: no