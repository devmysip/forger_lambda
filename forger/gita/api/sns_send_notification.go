package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go/aws"
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

	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		fmt.Println(err)
		return responseBuilder(0, nil, "Bad Request", err.Error())

	}

	snsClient := sns.NewFromConfig(sdkConfig)

	message := `{

  "GCM": "{\"notification\":{\"title\":\"Hello\",\"body\":\"World\"},\"data\":{\"screen\":\"/chaptersDetail\",\"arguments\":{\"chapter_no\":1,\"verse_no\":1}}}",
		"APNS": "{\"aps\": {\"alert\": \"%s\", \"sound\": \"default\"}}"
	}`
	// message, err := CreateDynamicMessage("Notificatuion", "Hii", "Read", "{\"notification\": {\"title\": \"Notification Title\", \"body\": \"kjwhfkjqwgjk\"}}")
	// if err != nil {
	// 	fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
	// 	fmt.Println(err)
	// 	return responseBuilder(0, nil, "Bad Request", err.Error())

	// }

	// Build the publish input
	input := &sns.PublishInput{
		Message:          aws.String(message),
		TargetArn:        aws.String(body.ClientEndpoint),
		MessageStructure: aws.String("json"),
	}

	// Send the push notification
	_, err = snsClient.Publish(context.TODO(), input)
	if err != nil {
		log.Printf("Failed to send SNS push notification: %v", err)
		return responseBuilder(0, nil, "Failed to send SNS push notification", err.Error())
	}

	return responseBuilder(1, message, "success", "")
}

type Notification struct {
	Default string            `json:"default"`
	GCM     map[string]string `json:"GCM"`
	APNS    map[string]string `json:"APNS"`
}

// Function to create dynamic notification message
func CreateDynamicMessage(defaultMessage, gcmTitle, gcmBody, gcmData string) (string, error) {
	// GCM payload in JSON format, including both "notification" and "data"
	gcmPayload := map[string]string{
		"notification": fmt.Sprintf("{\"title\": \"%s\", \"body\": \"%s\"}", gcmTitle, gcmBody),
		"data":         gcmData, // This will send the additional data payload
	}

	// APNS payload in JSON format
	apnsPayload := map[string]string{
		"APNS": "{\"aps\": {\"alert\": \"skjfskjgfd\", \"sound\": \"default\"}}",
	}

	// Create the notification object
	notification := Notification{
		Default: defaultMessage,
		GCM:     gcmPayload,
		APNS:    apnsPayload,
	}

	// Convert the struct to a JSON string
	message, err := json.Marshal(notification)
	if err != nil {
		return "", err
	}

	return string(message), nil
}
