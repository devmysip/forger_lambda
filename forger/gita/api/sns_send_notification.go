package api

import (
	"forger/gita/utilis"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

func SNSSendNotification(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {

	type RequestBody struct {
		ClientEndpoint string `json:"client_endpoint"`
	}

	email, err := headerHandler(request.Headers)
	if err != nil {
		log.Printf("Error extracting email: %s", email)

		return responseBuilder(0, nil, "Unauthorised User", err.Error())
	}

	body, err := decodeAndUnmarshal[RequestBody](request)
	if err != nil {
		log.Printf("Error decoding and unmarshalling request body: %s", err)
		return responseBuilder(0, nil, "Failed to parse request body", err.Error())
	}

	data := map[string]interface{}{
		"screen": "/chaptersDetail",
		"arguments": map[string]interface{}{
			"chapter_no": 1,
			"verse_no":   1,
		},
	}

	message, err := createMessage("Hello", "World", data)
	if err != nil {
		log.Printf("Failed to send SNS push notification: %v", err)
		return responseBuilder(0, nil, "Failed to send SNS push notification", err.Error())
	}

	err = utilis.SendNotification(body.ClientEndpoint, message)
	if err != nil {
		log.Printf("Failed to send SNS push notification: %v", err)
		return responseBuilder(0, nil, "Failed to send SNS push notification", err.Error())
	}

	return responseBuilder(1, message, "success", "")
}
