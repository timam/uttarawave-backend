package db

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/spf13/viper"
	"github.com/timam/uttarawave-backend/pkg/logger"
	"go.uber.org/zap"
)

var DynamoDB *dynamodb.DynamoDB

func InitializeDynamoDB() error {
	region := viper.GetString("database.dynamodb.region")
	isDev := os.Getenv("IS_DEV") == "true"

	var creds *credentials.Credentials
	if isDev {
		creds = credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			os.Getenv("AWS_SESSION_TOKEN"),
		)
	} else {
		creds = credentials.NewEnvCredentials()
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: creds,
	})

	if err != nil {
		logger.Error("Failed to create session", zap.Error(err))
		return fmt.Errorf("failed to create session: %v", err)
	}

	DynamoDB = dynamodb.New(sess)

	return nil
}
