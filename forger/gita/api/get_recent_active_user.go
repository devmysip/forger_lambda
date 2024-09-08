package api

import (
	"fmt"
	"forger/db"
	"forger/gita/models"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func GetRecentActiveUser(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	svc := dynamodb.New(db.DB())

	// Calculate the date range for the past 7 days
	now := time.Now()
	sevenDaysAgo := now.AddDate(0, 0, -7)

	// Format dates to match your DynamoDB attribute format
	nowFormatted := now.Format("2006-01-02T15:04:05Z")
	sevenDaysAgoFormatted := sevenDaysAgo.Format("2006-01-02T15:04:05Z")

	// Define the scan input
	scanInput := &dynamodb.ScanInput{
		TableName:        aws.String("User"),
		FilterExpression: aws.String("updated_at BETWEEN :start_date AND :end_date"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":start_date": {
				S: aws.String(sevenDaysAgoFormatted),
			},
			":end_date": {
				S: aws.String(nowFormatted),
			},
		},
		ProjectionExpression: aws.String("email, fcm_token, updated_at, display_name"),
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

	// Return the response
	return responseBuilder(1, users, fmt.Sprintf("%s - %s", nowFormatted, sevenDaysAgoFormatted), "")

}
