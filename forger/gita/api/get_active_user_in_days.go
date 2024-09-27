package api

import (
	"encoding/json"
	"fmt"
	"forger/db"
	"forger/gita/models"
	s3services "forger/gita/s3_services"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func GetActiveUserInDays(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {

	type RequestBody struct {
		Days int `json:"days"`
	}

	email, err := headerHandler(request.Headers)
	if err != nil {
		log.Printf("Error extracting email: %s", email)

		return responseBuilder(0, nil, "Internal Server Error", "Failed to extract email from request")
	}

	body, err := decodeAndUnmarshal[RequestBody](request)
	if err != nil {
		log.Printf("Error decoding and unmarshalling request body: %s", err)
		return responseBuilder(0, nil, "Bad Request", "Failed to parse request body")
	}

	// Load IST timezone
	istLocation, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return responseBuilder(0, nil, "Failed", err.Error())
	}

	svc := dynamodb.New(db.DB())
	now := time.Now().In(istLocation)

	endDate := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, istLocation).Format("2006-01-02")
	startDate := time.Date(now.Year(), now.Month(), now.Day()-body.Days, now.Hour(), 0, 0, 0, istLocation).Format("2006-01-02")

	filterExpression := "updated_at BETWEEN :start_date AND :end_date"
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":start_date": {
			S: aws.String(startDate),
		},
		":end_date": {
			S: aws.String(endDate),
		},
	}

	scanInput := &dynamodb.ScanInput{
		TableName:                 aws.String("User"),
		FilterExpression:          aws.String(filterExpression),
		ExpressionAttributeValues: expressionAttributeValues,
		ProjectionExpression:      aws.String("email, client_endpoint, updated_at, display_name, last_read"),
	}

	// Perform the scan operation
	result, err := svc.Scan(scanInput)
	if err != nil {
		fmt.Println("Error scanning table:", err)
		return responseBuilder(0, result, "Failed", "")
	}

	// Parse the result into a slice of users
	var users []models.User
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		fmt.Println("Error unmarshalling result:", err)
		return responseBuilder(0, users, "Failed", "")
	}
	// Create a temporary JSON file
	tempFile, err := os.CreateTemp("", "active_users_*.json")
	if err != nil {
		log.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	if err := json.NewEncoder(tempFile).Encode(users); err != nil {
		log.Fatalf("Failed to write users to temporary file: %v", err)
	}

	bucket := "com.gitasarathi"
	objectKey := fmt.Sprintf("%s.json", now.Format("2006-01-02"))

	if err := s3services.UploadFileToS3(tempFile.Name(), bucket, objectKey); err != nil {
		log.Fatalf("Failed to upload file: %v", err)
	}

	return responseBuilder(1, map[string]interface{}{
	
		"user_count": len(users)}, "success", "")

}
