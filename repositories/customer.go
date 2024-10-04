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
	GetCustomer(mobile string) (*models.Customer, error)
	GetAllCustomers() ([]models.Customer, error)
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

func (r *DynamoDBCustomerRepository) GetCustomer(mobile string) (*models.Customer, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Mobile": {
				S: aws.String(mobile),
			},
		},
	}

	result, err := db.DynamoDB.GetItem(input)
	if err != nil {
		logger.Error("Failed to get item from DynamoDB", zap.Error(err))
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	customer := &models.Customer{}
	err = dynamodbattribute.UnmarshalMap(result.Item, customer)
	if err != nil {
		logger.Error("Failed to unmarshal DynamoDB item", zap.Error(err))
		return nil, err
	}

	return customer, nil
}

func (r *DynamoDBCustomerRepository) GetAllCustomers() ([]models.Customer, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	}

	result, err := db.DynamoDB.Scan(input)
	if err != nil {
		logger.Error("Failed to scan items from DynamoDB", zap.Error(err))
		return nil, err
	}

	var customers []models.Customer
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &customers)
	if err != nil {
		logger.Error("Failed to unmarshal DynamoDB items", zap.Error(err))
		return nil, err
	}

	return customers, nil
}
