package userengagement

import (
	"fmt"
	"forger/db"
	"forger/gita/constants"
	"forger/gita/crud"
	"forger/gita/models"
	s3services "forger/gita/s3_services"
	"forger/gita/utilis"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func SendDailyNotification(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {

	now := utilis.GetLocalTime()

	bucket := constants.GitaSarathiBucket
	objectKey := fmt.Sprintf("%s/%s.json", constants.DailyNotifictionUserBucketDirectory, now.Format("2006-01-02"))

	users, err := s3services.DownloadFileFromS3[[]models.User](bucket, objectKey)
	if err != nil {
		return utilis.ResponseBuilder(0, users, "Failed to write users to temporary file", "")
	}

	filterUser, err := _filterUser(users)
	if err != nil {
		return utilis.ResponseBuilder(0, users, "Failed to write users to temporary file", "")
	}

	_sendNotificationToClients(filterUser)

	response := map[string]interface{}{
		"users":     users,
		"filter":    filterUser,
		"objectKey": objectKey,
	}

	return utilis.ResponseBuilder(1, response, "success", "")

}

func _sendNotificationToClients(users []models.User) {
	svc := dynamodb.New(db.DB())
	notificationTemplates := utilis.GetNotificationTemplates()
	notificationSent := 0

	for _, user := range users {
		if user.ClientEndpoint == nil {
			continue
		}

		days, err := utilis.DaysSinceDate(user.UpdatedAt)
		if err != nil {
			log.Printf("Failed to calculate days: %v", err)

		}

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

		message, err = utilis.FCMPayloadBuilder(notification.Title, notification.Body, data)
		if err != nil {
			log.Printf("Failed to create message: %v", err)
			continue
		}

		err = utilis.SendNotification(*user.ClientEndpoint, message)
		if err != nil {
			log.Printf("Failed to send SNS push notification: %v", err)
			continue
		}

		notificationSent++

	}

	crud.UpdateNotificationSentCount(svc, notificationSent)
}

func _filterUser(users []models.User) ([]models.User, error) {

	var filterdUser []models.User
	now := utilis.GetLocalTime()

	for _, user := range users {

		startMinute := 0
		endMinute := 30
		if now.Minute() > 30 {
			startMinute = 31
			endMinute = 59
		}

		endDate := time.Date(now.Year(), now.Month(), now.Day()-1, now.Hour(), endMinute, 0, 0, now.Location())
		startDate := time.Date(now.Year(), now.Month(), now.Day()-1, now.Hour(), startMinute, 0, 0, now.Location())

		specificStartTime := startDate.In(now.Location()).Format("2006-01-02T15:04:05Z07:00")
		specificEndTime := endDate.In(now.Location()).Format("2006-01-02T15:04:05Z07:00")

		if user.UpdatedAt >= specificStartTime && user.UpdatedAt <= specificEndTime {
			filterdUser = append(filterdUser, user)
		}

	}

	return filterdUser, nil
}
