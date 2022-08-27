package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/kvendingoldo/aws-cognito-restore-lambda/internal/cloud"
	"github.com/kvendingoldo/aws-cognito-restore-lambda/internal/config"
	log "github.com/sirupsen/logrus"
	"io"
	"time"
)

func getDataFromS3(ctx context.Context, client *cloud.Client, bucketName, keyName string) ([]byte, error) {
	obj, err := client.S3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName),
	})

	if err != nil {
		log.Errorf("Failed to get %s object data from %s bucket", keyName, bucketName)
		return nil, err
	}

	data, err := io.ReadAll(obj.Body)
	if err != nil {
		log.Errorf("Failed to convert %s object data to bytes", keyName)
		return nil, err
	}

	return data, err
}

func Execute(ctx context.Context, config config.Config) error {
	client, err := cloud.New(context.TODO(), config.CognitoRegion, config.S3BucketRegion)
	if err != nil {
		return fmt.Errorf("Could not create AWS client. Error: %s", err)
	}

	if config.CleanUpBeforeRestore.Bool {
		users, err := client.CognitoClient.ListUsers(ctx, &cognitoidentityprovider.ListUsersInput{
			UserPoolId: aws.String(config.CognitoUserPoolId),
		})
		if err != nil {
			return fmt.Errorf("[cleanup] Failed to get list of cognito users. Error: %s", err)
		}

		for _, user := range users.Users {
			_, err := client.CognitoClient.AdminDeleteUser(
				ctx, &cognitoidentityprovider.AdminDeleteUserInput{
					UserPoolId: aws.String(config.CognitoUserPoolId),
					Username:   user.Username,
				},
			)
			if err != nil {
				log.Errorf("[cleanup] Failed to user %s. Error: %s", *user.Username, err)
			} else {
				log.Debugf("User %s has been successefully deleted from %s userpool", *user.Username, config.CognitoUserPoolId)
			}
		}

		groups, err := client.CognitoClient.ListGroups(ctx, &cognitoidentityprovider.ListGroupsInput{
			UserPoolId: aws.String(config.CognitoUserPoolId),
		})
		if err != nil {
			return fmt.Errorf("[cleanup] Failed to get list of cognito groups. Error: %s", err)
		}

		for _, group := range groups.Groups {
			_, err := client.CognitoClient.DeleteGroup(
				ctx, &cognitoidentityprovider.DeleteGroupInput{
					UserPoolId: aws.String(config.CognitoUserPoolId),
					GroupName:  group.GroupName,
				},
			)

			if err != nil {
				log.Errorf("[cleanup] Failed to delete group %s. Error: %s", *group.GroupName, err)
			} else {
				log.Debugf("Group %s has been successefully deleted from %s userpool", *group.GroupName, config.CognitoUserPoolId)
			}
		}

		time.Sleep(3 * time.Second)
		log.Infof("User pool %s has been successfully cleaned up", config.CognitoUserPoolId)
	}

	if config.RestoreUsers.Bool {
		data, err := getDataFromS3(ctx, client, config.S3BucketName, fmt.Sprintf("%s/users.json", config.BackupDirPath))
		if err != nil {
			return fmt.Errorf("Failed to get users backup data from S3. Error: %s", err)
		} else {
			log.Debugf("%s/users.json data has been received successfully from S3", config.BackupDirPath)
		}

		var users cognitoidentityprovider.ListUsersOutput
		err = json.Unmarshal(data, &users)
		if err != nil {
			return fmt.Errorf("Failed to unmarshal users backup data. Error: %s", err)
		} else {
			log.Debug("users data has been unmarshalled successfully")
		}

		for _, user := range users.Users {
			var userAttributes []types.AttributeType
			var userName *string

			for _, attribute := range user.Attributes {
				if *attribute.Name == "email" {
					userName = attribute.Value
				}

				// NOTE: it's unique value for each Cognito UserPool and for that reason we do not need to
				//       include it to restore attributes
				if *attribute.Name != "sub" {
					userAttributes = append(userAttributes, attribute)
				}
			}

			_, err := client.CognitoClient.AdminCreateUser(
				ctx, &cognitoidentityprovider.AdminCreateUserInput{
					UserPoolId:     aws.String(config.CognitoUserPoolId),
					Username:       userName,
					UserAttributes: userAttributes,
				},
			)
			if err != nil {
				return fmt.Errorf("Failed to restore user %s. Error: %s", *user.Username, err)
			}
		}
	}

	if config.RestoreGroups.Bool {
		data, err := getDataFromS3(ctx, client, config.S3BucketName, fmt.Sprintf("%s/groups.json", config.BackupDirPath))
		if err != nil {
			return fmt.Errorf("Failed to get groups backup data from S3. Error: %s", err)
		} else {
			log.Debugf("%s/groups.json data has been received successfully from S3", config.BackupDirPath)
		}

		var groups cognitoidentityprovider.ListGroupsOutput
		err = json.Unmarshal(data, &groups)
		if err != nil {
			return fmt.Errorf("Failed to unmarshal groups backup data. Error: %s", err)

		} else {
			log.Debug("groups data has been unmarshalled successfully")
		}

		for _, group := range groups.Groups {
			_, err := client.CognitoClient.CreateGroup(
				ctx, &cognitoidentityprovider.CreateGroupInput{
					UserPoolId: aws.String(config.CognitoUserPoolId),
					GroupName:  group.GroupName,
				},
			)
			if err != nil {
				return fmt.Errorf("Failed to restore group %s. Error: %s", *group.GroupName, err)
			}
		}
	}

	return nil
}
