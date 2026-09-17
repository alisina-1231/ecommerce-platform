package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func NewDynamoDBClient(ctx context.Context) (*dynamodb.Client, error) {
	region := os.Getenv("AWS_REGION")

	if region == "" {
		region = "us-east-1"
	}

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(region),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	log.Printf("DynamoDB region: %s", region)

	return dynamodb.NewFromConfig(cfg), nil
}

func CheckConnection(
	ctx context.Context,
	client *dynamodb.Client,
	tableName string,
) error {
	_, err := client.DescribeTable(
		ctx,
		&dynamodb.DescribeTableInput{
			TableName: aws.String(tableName),
		},
	)

	if err != nil {
		return fmt.Errorf(
			"failed to connect to DynamoDB table %s: %w",
			tableName,
			err,
		)
	}

	return nil
}
