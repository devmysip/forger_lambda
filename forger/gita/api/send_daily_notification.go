package api

import (
	"fmt"
	"forger/db"
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

	istLocation, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return responseBuilder(0, nil, "Failed", err.Error())
	}

	now := time.Now().In(istLocation)

	bucket := "com.gitasarathi"
	objectKey := fmt.Sprintf("%s.json", now.Format("2006-01-02"))

	users, err := s3services.DownloadFileFromS3[[]models.User](bucket, objectKey)
	if err != nil {
		return responseBuilder(0, users, "Failed to write users to temporary file", "")
	}

	filterUser, err := filterUser(users)
	if err != nil {
		return responseBuilder(0, users, "Failed to write users to temporary file", "")
	}

	sendDailyNotification(filterUser)

	return responseBuilder(1, map[string]interface{}{
		"users":     users,
		"filter":    filterUser,
		"objectKey": objectKey,
	}, "success", "")

}

func sendDailyNotification(users []models.User) {
	svc := dynamodb.New(db.DB())
	notificationTemplates := utilis.GetNotificationTemplates()
	notificationSent := 0

	for _, user := range users {
		if user.ClientEndpoint == nil {
			continue
		}

		days, err := daysBetween(user.UpdatedAt)
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

		message, err = createMessage(notification.Title, notification.Body, data)
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

	updateNotificationSent(svc, notificationSent)

}

func filterUser(users []models.User) ([]models.User, error) {

	var filterdUser []models.User
	istLocation, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return filterdUser, err
	}

	now := time.Now().In(istLocation)

	for _, user := range users {

		startMinute := 0
		endMinute := 30
		if now.Minute() > 30 {
			startMinute = 31
			endMinute = 59
		}

		endDate := time.Date(now.Year(), now.Month(), now.Day()-1, now.Hour(), endMinute, 0, 0, istLocation)
		startDate := time.Date(now.Year(), now.Month(), now.Day()-1, now.Hour(), startMinute, 0, 0, istLocation)

		specificStartTime := startDate.In(istLocation).Format("2006-01-02T15:04:05Z07:00")
		specificEndTime := endDate.In(istLocation).Format("2006-01-02T15:04:05Z07:00")

		if user.UpdatedAt >= specificStartTime && user.UpdatedAt <= specificEndTime {
			filterdUser = append(filterdUser, user)
		}

	}

	return filterdUser, nil
}
