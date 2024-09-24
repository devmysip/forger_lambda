package api

import (
	"fmt"
	"forger/db"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func UpdateNotificationReadCounter(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	email, err := headerHandler(request.Headers)
	if err != nil {
		log.Printf("Error extracting email: %s", err)
		return responseBuilder(0, nil, "Internal Server Error", "Failed to extract email from request")
	}

	attributeName := "notification_clicked"
	svc := dynamodb.New(db.DB())

	// Update expression to increment the notification_clicked attribute
	updateExpression := fmt.Sprintf("SET %s = if_not_exists(%s, :start) + :incr", attributeName, attributeName)

	input := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":incr": {
				N: aws.String("1"),
			},
			":start": {
				N: aws.String("0"),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		TableName:        aws.String("User"),
		UpdateExpression: aws.String(updateExpression),
	}

	_, err = svc.UpdateItem(input)
	if err != nil {
		log.Printf("Error updating notification counter: %s", err)
		return responseBuilder(0, nil, "Internal Server Error", "Failed to update notification counter")
	}

	date := time.Now().Format("2006-01-02")

	// Define the input for the UpdateItem API
	inputAnalytics := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"date": {
				S: aws.String(date),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":incr": {
				N: aws.String("1"),
			},
			":start": {
				N: aws.String("0"),
			},
		},
		ReturnValues: aws.String("UPDATED_NEW"),
		TableName: aws.String("Analytics"),
		UpdateExpression: aws.String(updateExpression),

	}

	_, err = svc.UpdateItem(inputAnalytics)
	if err != nil {
		log.Printf("Error updating notification counter: %s", err)
		return responseBuilder(0, nil, "Internal Server Error", "Failed to update notification counter")
	}

	return responseBuilder(1, nil, "Success", "Notification counter updated successfully")
}
