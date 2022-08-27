package cloud

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	CognitoClient *cognitoidentityprovider.Client
	S3Client      *s3.Client
}

func New(ctx context.Context, cognitoRegion, s3BucketRegion string) (*Client, error) {
	cognitoCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(cognitoRegion))
	if err != nil {
		return nil, err
	}

	s3Cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(s3BucketRegion))
	if err != nil {
		return nil, err
	}

	return &Client{
		CognitoClient: cognitoidentityprovider.NewFromConfig(cognitoCfg),
		S3Client:      s3.NewFromConfig(s3Cfg),
	}, nil
}
