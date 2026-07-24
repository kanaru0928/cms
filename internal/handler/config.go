package handler

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port int
	JWKSURL string
	S3BucketName string
	S3Prefix string
	DynamoDBTableName string
}

func LoadConfig() (*Config, error) {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return nil, fmt.Errorf("invalid PORT environment variable: %v", err)
	}
	jwksURL := os.Getenv("JWKS_URL")
	s3BucketName := os.Getenv("S3_BUCKET_NAME")
	s3Prefix := os.Getenv("S3_PREFIX")
	dynamoDBTableName := os.Getenv("DYNAMODB_TABLE_NAME")

	if port == 0 {
		port = 8080
	}
	if jwksURL == "" {
		return nil, fmt.Errorf("JWKS_URL environment variable is not set")
	}
	if s3BucketName == "" {
		return nil, fmt.Errorf("S3_BUCKET_NAME environment variable is not set")
	}
	if s3Prefix == "" {
		return nil, fmt.Errorf("S3_PREFIX environment variable is not set")
	}
	if dynamoDBTableName == "" {
		return nil, fmt.Errorf("DYNAMODB_TABLE_NAME environment variable is not set")
	}

	return &Config{
		Port: port,
		JWKSURL: jwksURL,
		S3BucketName: s3BucketName,
		DynamoDBTableName: dynamoDBTableName,
	}, nil
}
