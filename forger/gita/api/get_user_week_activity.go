package api

import (
	"forger/db"
	"forger/gita/constants"
	"forger/gita/models"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func GetUserWeekActivity(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {

	svc := dynamodb.New(db.DB())

	email, err := headerHandler(request.Headers)
	if err != nil {
		log.Printf("Error extracting email: %s", err)
		return responseBuilder(0, nil, "Internal Server Error", "Failed to extract email from request")
	}

	currentTime := time.Now()
	weekday := int(currentTime.Weekday())

	// Calculate the start of the week (Monday)
	startOfWeek := currentTime.AddDate(0, 0, -weekday+1)
	// Calculate the end of the week (Sunday)
	endOfWeek := startOfWeek.AddDate(0, 0, 6)

	startOfWeekFormatted := startOfWeek.Format("2006-01-02")
	endOfWeekFormatted := endOfWeek.Format("2006-01-02")

	// Query DynamoDB for the user's activity within the date range
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(constants.UserActivityTable),
		KeyConditionExpression: aws.String("#email = :email AND #date BETWEEN :start_date AND :end_date"),
		ExpressionAttributeNames: map[string]*string{
			"#email": aws.String("email"),
			"#date":  aws.String("date"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":email":      {S: aws.String(email)},
			":start_date": {S: aws.String(startOfWeekFormatted)},
			":end_date":   {S: aws.String(endOfWeekFormatted)},
		},
	}

	result, err := svc.Query(queryInput)
	if err != nil {
		log.Printf("Error querying DynamoDB: %s", err)
		return responseBuilder(0, nil, "Internal Server Error", err.Error())
	}

	var userActivities []models.UserActivity
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &userActivities)
	if err != nil {
		log.Printf("Error unmarshalling result items: %s", err)
		return responseBuilder(0, nil, "Internal Server Error", "Failed to parse query results")
	}

	// Initialize a map to store activity by date
	activityMap := make(map[string][]struct {
		ChapterNo string `json:"chapter_no"`
		VerseNo   string `json:"verse_no"`
	})

	for _, activity := range userActivities {
		activityMap[activity.Date] = activity.Activity
	}

	weekActivity := make([]models.UserActivity, 7)
	today := time.Now().Format("2006-01-02")
	daysOfWeek := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	for i := 0; i < 7; i++ {
		date := startOfWeek.AddDate(0, 0, i).Format("2006-01-02")
		day := daysOfWeek[i]
		if date > today {
			weekActivity[i] = models.UserActivity{
				Date:     date,
				Day:      day,
				Activity: nil,
			}
		} else if activities, found := activityMap[date]; found {

			weekActivity[i] = models.UserActivity{
				Date:     date,
				Day:      day,
				Activity: activities,
			}
		} else {
			weekActivity[i] = models.UserActivity{
				Date: date,
				Day:  day,
				Activity: []struct {
					ChapterNo string `json:"chapter_no"`
					VerseNo   string `json:"verse_no"`
				}{},
			}
		}
	}
	return responseBuilder(1, weekActivity, "Success", "User activity retrieved successfully")
}
