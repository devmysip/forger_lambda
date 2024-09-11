package api

import (
	"fmt"
	"forger/db"
	"forger/gita/models"
	"forger/gita/utilis"
	"log"
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
	for _, user := range users {
		if user.ClientEndpoint == nil {
			continue
		}

		message := ""
		if user.LastRead != nil {
			lastRead := *user.LastRead
			trimmed := strings.TrimPrefix(lastRead, "BG")
			parts := strings.Split(trimmed, ".")

			if len(parts) != 2 {
				continue
			}

			chapterNo := parts[0]
			verseNo := parts[1]

			data := map[string]interface{}{
				"screen": "/chaptersDetail",
				"arguments": map[string]interface{}{
					"chapter_no": chapterNo,
					"verse_no":   verseNo,
				},
			}

			message, err = createMessage("Hello", "World", data)
			if err != nil {
				log.Printf("Failed to create message: %v", err)
				continue
			}
		}

		err := utilis.SendNotification(*user.ClientEndpoint, message)
		if err != nil {
			log.Printf("Failed to send SNS push notification: %v", err)
		}
	}

	return responseBuilder(1, users, "success", "")

}
