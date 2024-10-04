package repositories

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/spf13/viper"
	"github.com/timam/uttarawave-finance-backend/models"
	"github.com/timam/uttarawave-finance-backend/pkg/db"
	"github.com/timam/uttarawave-finance-backend/pkg/logger"
	"go.uber.org/zap"
	"os"
)

type CustomerRepository interface {
	CreateCustomer(customer *models.Customer) error
}

type DynamoDBCustomerRepository struct {
	tableName string
}

func NewDynamoDBCustomerRepository() *DynamoDBCustomerRepository {
	tableName := viper.GetString("database.dynamodb.tables.customer")
	env := os.Getenv("ENV")
	if env == "dev" {
		tableName = "dev-" + tableName
	} else if env == "prod" {
		tableName = "prod-" + tableName
	}
	return &DynamoDBCustomerRepository{tableName: tableName}
}

func (r *DynamoDBCustomerRepository) CreateCustomer(customer *models.Customer) error {
	av, err := dynamodbattribute.MarshalMap(customer)
	if err != nil {
		logger.Error("Failed to marshal customer", zap.Error(err))
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(r.tableName),
	}

	_, err = db.DynamoDB.PutItem(input)
	if err != nil {
		logger.Error("Failed to put item in DynamoDB", zap.Error(err))
		return err
	}

	logger.Info("Customer created")

	return nil
}
