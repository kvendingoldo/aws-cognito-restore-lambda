package config

import (
	"fmt"
	"github.com/guregu/null"
	"github.com/kvendingoldo/aws-cognito-restore-lambda/internal/types"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

type Config struct {
	AWSRegion string

	CognitoUserPoolID string
	CognitoRegion     string

	S3BucketName   string
	S3BucketRegion string

	BackupDirPath string

	RestoreUsers         null.Bool
	RestoreGroups        null.Bool
	CleanUpBeforeRestore null.Bool
}

//nolint:unparam
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

//nolint:gocyclo,maintidx
func New(eventRaw interface{}) (*Config, error) {
	var config = &Config{}

	var getFromEvent bool
	var event types.Event

	switch value := eventRaw.(type) {
	case types.Event:
		getFromEvent = true
		event = value
	default:
		getFromEvent = false
	}

	// Process AWSRegion
	if awsRegion := getEnv("AWS_REGION", ""); awsRegion != "" {
		config.AWSRegion = awsRegion
	} else {
		log.Warn("Environment variable AWS_REGION is empty")
	}
	if getFromEvent {
		if event.AWSRegion != "" {
			config.AWSRegion = event.AWSRegion
		} else {
			log.Warn("Event contains empty awsRegion variable")
		}
	}
	if config.AWSRegion == "" {
		return nil, fmt.Errorf("awsRegion is empty; Configure it via 'AWS_REGION' env variable OR pass in event body")
	}

	// Process CognitoRegion
	if cognitoRegion := getEnv("COGNITO_REGION", ""); cognitoRegion != "" {
		config.CognitoRegion = cognitoRegion
	} else {
		log.Warn("Environment variable 'COGNITO_REGION' is empty")
	}
	if getFromEvent {
		if event.CognitoRegion != "" {
			config.CognitoRegion = event.CognitoRegion
		} else {
			log.Warn("Event contains empty cognitoRegion variable")
		}
	}
	if config.CognitoRegion == "" {
		log.Warnf("cognitoRegion is empty; Default region %s will be used", config.AWSRegion)
		config.CognitoRegion = config.AWSRegion
	}

	// Process S3BucketRegion
	if bucketRegion := getEnv("S3_BUCKET_REGION", ""); bucketRegion != "" {
		config.S3BucketRegion = bucketRegion
	} else {
		log.Warn("Environment variable 'S3_BUCKET_REGION' is empty")
	}
	if getFromEvent {
		if event.S3BucketRegion != "" {
			config.S3BucketRegion = event.S3BucketRegion
		} else {
			log.Warn("Event contains empty s3BucketRegion variable")
		}
	}
	if config.S3BucketRegion == "" {
		log.Warnf("bucketRegion is empty; Default region %s will be used", config.AWSRegion)
		config.S3BucketRegion = config.AWSRegion
	}

	// Process CognitoUserPoolId
	if cognitoUserPoolID := getEnv("COGNITO_USER_POOL_ID", ""); cognitoUserPoolID != "" {
		config.CognitoUserPoolID = cognitoUserPoolID
	} else {
		log.Warn("Environment variable 'COGNITO_USER_POOL_ID' is empty")
	}
	if getFromEvent {
		if event.CognitoUserPoolID != "" {
			config.CognitoUserPoolID = event.CognitoUserPoolID
		} else {
			log.Warn("Event contains empty cognitoUserPoolId")
		}
	}
	if config.CognitoUserPoolID == "" {
		return nil, fmt.Errorf("cognitoUserPoolID is empty; Configure it via 'COGNITO_USER_POOL_ID' env variable OR pass in event body")
	}

	// Process S3BucketName
	if s3BucketName := getEnv("S3_BUCKET_NAME", ""); s3BucketName != "" {
		config.S3BucketName = s3BucketName
	} else {
		log.Warn("Environment variable 'S3_BUCKET_NAME' is empty")
	}
	if getFromEvent {
		if event.S3BucketName != "" {
			config.S3BucketName = event.S3BucketName
		} else {
			log.Warn("Event contains empty s3BucketName")
		}
	}
	if config.S3BucketName == "" {
		return nil, fmt.Errorf("BucketName is empty; Configure it via 'S3_BUCKET_NAME' env variable OR pass in event body")
	}

	// Process BackupDirPath
	if backupDirPath := getEnv("BACKUP_DIR_PATH", ""); backupDirPath != "" {
		config.BackupDirPath = backupDirPath
	} else {
		log.Warn("Environment variable 'BACKUP_DIR_PATH' is empty")
	}
	if getFromEvent {
		if event.BackupDirPath == "" {
			log.Warn("Event contains empty backupDirPath")
		} else {
			config.BackupDirPath = event.BackupDirPath
		}
	}
	if config.BackupDirPath == "" {
		return nil, fmt.Errorf("BackupDirPath is empty; Configure it via 'BACKUP_DIR_PATH' env variable OR pass in event body")
	}

	// Process RestoreUsers
	if restoreUsers := getEnv("RESTORE_USERS", ""); restoreUsers != "" {
		restoreUsersValue, err := strconv.ParseBool(restoreUsers)
		if err != nil {
			//nolint:stylecheck
			return nil, fmt.Errorf("Could not parse 'RESTORE_USERS' variable. Error: %w", err)
		}
		config.RestoreUsers = null.NewBool(restoreUsersValue, true)
	} else {
		log.Warn("Environment variable 'RESTORE_USERS' is empty")
	}
	if getFromEvent {
		if event.RestoreUsers.Valid {
			config.RestoreUsers = event.RestoreUsers
		}
	}
	if !config.RestoreUsers.Valid {
		log.Warn("restoreUsers is not specified; Default value 'false' will be used")
		config.RestoreUsers = null.NewBool(false, true)
	}

	// Process RestoreGroups
	if restoreGroups := getEnv("RESTORE_GROUPS", ""); restoreGroups != "" {
		restoreGroupsValue, err := strconv.ParseBool(restoreGroups)
		if err != nil {
			//nolint:stylecheck
			return nil, fmt.Errorf("Could not parse 'RESTORE_GROUPS' variable. Error: %w", err)
		}
		config.RestoreGroups = null.NewBool(restoreGroupsValue, true)
	} else {
		log.Warn("Environment variable 'RESTORE_GROUPS' is empty")
	}
	if getFromEvent {
		if event.RestoreGroups.Valid {
			config.RestoreGroups = event.RestoreGroups
		}
	}
	if !config.RestoreGroups.Valid {
		log.Warn("restoreGroups is not specified; Default value 'false' will be used")
		config.RestoreGroups = null.NewBool(false, true)
	}

	// Process cleanUpBeforeRestore
	if cleanUpBeforeRestore := getEnv("CLEANUP_BEFORE_RESTORE", ""); cleanUpBeforeRestore != "" {
		cleanUpBeforeRestoreValue, err := strconv.ParseBool(cleanUpBeforeRestore)
		if err != nil {
			//nolint:stylecheck
			return nil, fmt.Errorf("Could not parse 'CLEANUP_BEFORE_RESTORE' variable. Error: %w", err)
		}
		config.CleanUpBeforeRestore = null.NewBool(cleanUpBeforeRestoreValue, true)
	} else {
		log.Warn("Environment variable 'CLEANUP_BEFORE_RESTORE' is empty")
	}
	if getFromEvent {
		if event.CleanUpBeforeRestore.Valid {
			config.CleanUpBeforeRestore = event.CleanUpBeforeRestore
		}
	}
	if !config.CleanUpBeforeRestore.Valid {
		log.Warn("cleanUpBeforeRestore is not specified; Default value 'false' will be used")
		config.CleanUpBeforeRestore = null.NewBool(false, true)
	} else { //nolint:gocritic
		if config.CleanUpBeforeRestore.Bool {
			log.Warnf("Pay attention that CLEANUP_BEFORE_RESTORE is 'true'. It means that all data from %s userpool will be deleted before restore", config.CognitoUserPoolID)
		}
	}

	return config, nil
}
