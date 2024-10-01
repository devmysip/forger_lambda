package userwatch

import (
	"fmt"
	"forger/db"
	"forger/gita/constants"
	"forger/gita/models"
	"forger/gita/utilis"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func UpdateDailyAnalytics(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	svc := dynamodb.New(db.DB())

	istLocation, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Printf("Error loading location: %v", err)
		return utilis.ResponseBuilder(0, nil, "Internal Server Error", "Failed to parse chapter data: Unmarshal error")
	}

	// Get the current time in IST
	now := time.Now().In(istLocation)
	oneDayAgo := now.AddDate(0, 0, -1)

	updatedUser, err := getTodayActiveUser(svc, oneDayAgo, istLocation)
	if err != nil {
		log.Printf("Failed to unmarshal item: %v", err)
		return utilis.ResponseBuilder(0, nil, "Internal Server Error", "Failed to parse chapter data: Unmarshal error")
	}

	newUser, err := getTodayNewUser(svc, oneDayAgo, istLocation)
	if err != nil {
		log.Printf("Failed to unmarshal item: %v", err)
		return utilis.ResponseBuilder(0, nil, "Internal Server Error", "Failed to parse chapter data: Unmarshal error")
	}

	activity, err := getTodayUserActivity(svc, oneDayAgo, istLocation)
	if err != nil {
		log.Printf("Failed to unmarshal item: %v", err)
		return utilis.ResponseBuilder(0, nil, "Internal Server Error", err.Error())
	}

	// Prepare data for DynamoDB PutItem
	date := oneDayAgo.Format("2006-01-02")

	updateExpression := "SET active_users = :active_users, new_users = :new_users, user_activity = :user_activity"
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":active_users": {
			N: aws.String(fmt.Sprintf("%d", len(updatedUser))),
		},
		":new_users": {
			N: aws.String(fmt.Sprintf("%d", len(newUser))),
		},
		":user_activity": {
			N: aws.String(fmt.Sprintf("%d", len(activity))),
		},
	}

	updateItemInput := &dynamodb.UpdateItemInput{
		TableName: aws.String("Analytics"),
		Key: map[string]*dynamodb.AttributeValue{
			"date": {
				S: aws.String(date),
			},
		},
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: expressionAttributeValues,
	}

	// Execute the update operation
	_, err = svc.UpdateItem(updateItemInput)
	if err != nil {
		log.Printf("Failed to update item: %v", err)
		return utilis.ResponseBuilder(0, nil, "Internal Server Error", "Failed to update analytics data")
	}

	return utilis.ResponseBuilder(1, map[string]interface{}{
		"active_users":  len(updatedUser),
		"new_users":     len(newUser),
		"user_activity": len(activity),
		"date":          date,
	}, "success", "")
}

func getTodayActiveUser(svc *dynamodb.DynamoDB, now time.Time, istLocation *time.Location) ([]models.User, error) {

	// Start and end of the day in YYYY-MM-DD format
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, istLocation).Format("2006-01-02T15:04:05Z07:00")
	endDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, istLocation).Format("2006-01-02T15:04:05Z07:00")

	// Define the DynamoDB Scan input
	input := &dynamodb.ScanInput{
		TableName:        aws.String(constants.UserTable),
		FilterExpression: aws.String("updated_at BETWEEN :start_date AND :end_date"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":start_date": {
				S: aws.String(startDate),
			},
			":end_date": {
				S: aws.String(endDate),
			},
		},
		ProjectionExpression: aws.String("email, updated_at"),
	}

	// Scan the table
	result, err := svc.Scan(input)
	if err != nil {
		log.Printf("Error querying DynamoDB: %v", err)
		return nil, err
	}

	// Unmarshal the result into a slice of User models
	var users []models.User
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		log.Printf("Failed to unmarshal items: %v", err)
		return nil, err
	}

	// Return the list of active users
	return users, nil
}

func getTodayNewUser(svc *dynamodb.DynamoDB, now time.Time, istLocation *time.Location) ([]models.User, error) {

	// Start and end of the day in YYYY-MM-DD format
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, istLocation).Format("2006-01-02T15:04:05Z07:00")
	endDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, istLocation).Format("2006-01-02T15:04:05Z07:00")

	// Define the DynamoDB Scan input
	input := &dynamodb.ScanInput{
		TableName:        aws.String("User"),
		FilterExpression: aws.String("created_at BETWEEN :start_date AND :end_date"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":start_date": {
				S: aws.String(startDate),
			},
			":end_date": {
				S: aws.String(endDate),
			},
		},
		ProjectionExpression: aws.String("email, created_at"),
	}

	// Scan the table
	result, err := svc.Scan(input)
	if err != nil {
		log.Printf("Error querying DynamoDB: %v", err)
		return nil, err
	}

	// Unmarshal the result into a slice of User models
	var users []models.User
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		log.Printf("Failed to unmarshal items: %v", err)
		return nil, err
	}

	// Return the list of active users
	return users, nil
}

func getTodayUserActivity(svc *dynamodb.DynamoDB, now time.Time, istLocation *time.Location) ([]models.UserActivity, error) {

	// Get today's date in YYYY-MM-DD format
	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, istLocation).Format("2006-01-02")

	// Define the DynamoDB Scan input
	input := &dynamodb.ScanInput{
		TableName:        aws.String("UserActivity"),
		FilterExpression: aws.String("#d = :date"),
		ExpressionAttributeNames: map[string]*string{
			"#d": aws.String("date"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":date": {
				S: aws.String(date),
			},
		},
		ProjectionExpression: aws.String("email"), // Only return the email attribute
	}

	// Scan the table
	result, err := svc.Scan(input)
	if err != nil {
		log.Printf("Error querying DynamoDB: %v", err)
		return nil, err
	}

	// Unmarshal the result into a slice of User models
	var usersActivity []models.UserActivity
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &usersActivity)
	if err != nil {
		log.Printf("Failed to unmarshal items: %v", err)
		return nil, err
	}

	// Return the count of active users
	return usersActivity, nil
}
