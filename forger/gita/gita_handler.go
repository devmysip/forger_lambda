package gita

import (
	"forger/db"
	"forger/gita/api"
	snsservice "forger/gita/sns_service"
	userwatch "forger/gita/user_watch"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func GitaHandler(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {

	svc := dynamodb.New(db.DB())

	if strings.Contains(request.Path, "/gita/createUser") {
		return api.CreateUser(request, svc)
	}

	if strings.Contains(request.Path, "/gita/user") {
		return api.GetUser(request, svc)
	}

	if strings.Contains(request.Path, "/gita/updateRead") {
		return userwatch.UpdateUserRead(request, svc)
	}

	if strings.Contains(request.Path, "/gita/updateFCM") {
		return api.UpdateFCMToken(request, svc)
	}

	if strings.Contains(request.Path, "/gita/updateUserActivity") {
		return userwatch.UpdateUserActivity(request)
	}

	if strings.Contains(request.Path, "/gita/getUserWeekActivity") {
		return api.GetUserWeekActivity(request)
	}

	if strings.Contains(request.Path, "/gita/chapter") {
		return api.GetChapter(request)
	}

	if strings.Contains(request.Path, "/gita/verse") {
		return api.GetVerse(request)
	}

	if strings.Contains(request.Path, "/gita/snsCreate") {
		return snsservice.SNSCreateClientEndpoint(request)
	}

	if strings.Contains(request.Path, "/gita/snsUpdate") {
		return snsservice.SNSUpdateClientEndpoint(request)
	}

	if strings.Contains(request.Path, "/gita/snsSendNotification") {
		return snsservice.SNSSendNotification(request)
	}

	if strings.Contains(request.Path, "/gita/getActiveUserInTime") {
		return api.GetActiveUserInTime(request)
	}

	if strings.Contains(request.Path, "/gita/updateNotificationReadCounter") {
		return userwatch.UpdateNotificationReadCounter(request)
	}

	if strings.Contains(request.Path, "/gita/updateDailyAnalytics") {
		return userwatch.UpdateDailyAnalytics(request)
	}

	if strings.Contains(request.Path, "/gita/getActiveUserInDays") {
		return api.GetActiveUserInDays(request)
	}

	if strings.Contains(request.Path, "/gita/sendDailyNotification") {
		return api.SendDailyNotification(request)
	}

	return events.APIGatewayProxyResponse{
		Body:       "No Gita Path Found",
		StatusCode: http.StatusInternalServerError,
	}
}
