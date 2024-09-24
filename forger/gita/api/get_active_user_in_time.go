package api

import (
	"fmt"
	"forger/db"
	"forger/gita/models"
	"forger/gita/utilis"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DateRange struct {
	StartDate string
	EndDate   string
}

func GetActiveUserInTime(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {

	istLocation, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return responseBuilder(0, nil, "Failed", err.Error())
	}

	svc := dynamodb.New(db.DB())

	// Get the current time and time 7 days ago in UTC
	now := time.Now().In(istLocation)
	// sevenDaysAgo := now.AddDate(0, 0, -7)
	const daysBack = 7
	var filterExpression string
	expressionAttributeValues := make(map[string]*dynamodb.AttributeValue)

	for i := 0; i < daysBack; i++ {

		startMinute := 0
		endMinute := 30
		if now.Minute() > 30 {
			startMinute = 31
			endMinute = 59
		}

		// Calculate the start and end dates for each range
		endDate := time.Date(now.Year(), now.Month(), now.Day()-i, now.Hour(), endMinute, 0, 0, istLocation)
		startDate := time.Date(now.Year(), now.Month(), now.Day()-i, now.Hour(), startMinute, 0, 0, istLocation)

		specificStartTime := startDate.In(istLocation).Format("2006-01-02T15:04:05Z07:00")
		specificEndTime := endDate.In(istLocation).Format("2006-01-02T15:04:05Z07:00")

		// Dynamically create placeholders like :start_date_1, :end_date_1
		startPlaceholder := fmt.Sprintf(":start_date_%d", i+1)
		endPlaceholder := fmt.Sprintf(":end_date_%d", i+1)

		// Append to filterExpression
		if i > 0 {
			filterExpression += " OR "
		}
		filterExpression += fmt.Sprintf("updated_at BETWEEN %s AND %s", startPlaceholder, endPlaceholder)

		// Add the values to ExpressionAttributeValues
		expressionAttributeValues[startPlaceholder] = &dynamodb.AttributeValue{
			S: aws.String(specificStartTime),
		}
		expressionAttributeValues[endPlaceholder] = &dynamodb.AttributeValue{
			S: aws.String(specificEndTime),
		}
	}

	// Define the scan input
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

	notificationTemplates := utilis.GetNotificationTemplates()
	log.Printf("Failed to send SNS push notification: %v", err)

	// return responseBuilder(1, expressionAttributeValues, "success", "")

	notificationSent := 0

	for _, user := range users {
		if user.ClientEndpoint == nil {
			continue
		}

		// Determine the number of days since the user's last update
		days, err := daysBetween(user.UpdatedAt)
		if err != nil {
			log.Printf("Failed to calculate days: %v", err)

		}

		// Get the appropriate notification template
		notification := notificationTemplates[days]

		var message string
		var data map[string]interface{}

		if user.LastRead != nil {
			lastRead := *user.LastRead
			trimmed := strings.TrimPrefix(lastRead, "BG")
			parts := strings.Split(trimmed, ".")

			if len(parts) != 2 {
				log.Printf("Invalid LastRead format: %s", lastRead)
				continue
			}

			chapterNo, err := strconv.Atoi(parts[0])
			if err != nil {
				log.Println("Error converting chapter number:", err)
				continue
			}

			verseNo, err := strconv.Atoi(parts[1])
			if err != nil {
				log.Println("Error converting verse number:", err)
				continue
			}

			data = map[string]interface{}{
				"screen": "/chaptersDetail",
				"arguments": map[string]interface{}{
					"chapter_no": chapterNo,
					"verse_no":   verseNo,
				},
			}
		}

		message, err = createMessage(notification.Title, notification.Body, data)
		if err != nil {
			log.Printf("Failed to create message: %v", err)
			continue
		}
		// Send the notification
		err = utilis.SendNotification(*user.ClientEndpoint, message)
		if err != nil {
			log.Printf("Failed to send SNS push notification: %v", err)
		} else {
			notificationSent++
		}

	}

	updateNotificationSent(svc, notificationSent)
	return responseBuilder(1, users, "success", "")

}

func daysBetween(updatedAt string) (int, error) {
	const layout = "2006-01-02T15:04:05Z07:00"
	// Parse `updated_at` string into time.Time
	updatedTime, err := time.Parse(layout, updatedAt)
	if err != nil {
		return 0, fmt.Errorf("failed to parse `updated_at`: %v", err)
	}

	// Get today's date
	today := time.Now()

	// Calculate the difference between today and the `updated_at` date
	duration := today.Sub(updatedTime)

	// Convert duration to days
	days := int(duration.Hours() / 24)

	return days, nil
}

func updateNotificationSent(svc *dynamodb.DynamoDB, notificationSent int) error {
	// Assuming the table name is 'analytics' and the primary key is 'date'
	tableName := "Analytics"

	// Calculate one day ago in the format YYYY-MM-DD
	date := time.Now().Format("2006-01-02")

	// Update expression to increment the notificationSent count
	updateExpression := "SET notification_sent = if_not_exists(notification_sent, :zero) + :inc"

	// Define the input for the UpdateItem API
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"date": {
				S: aws.String(date),
			},
		},
		UpdateExpression: aws.String(updateExpression),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":inc": {
				N: aws.String(fmt.Sprintf("%d", notificationSent)),
			},
			":zero": {
				N: aws.String("0"),
			},
		},
		ReturnValues: aws.String("UPDATED_NEW"),
	}

	// Call UpdateItem
	_, err := svc.UpdateItem(input)
	if err != nil {
		return fmt.Errorf("failed to update notificationSent count: %w", err)
	}

	return nil
}
