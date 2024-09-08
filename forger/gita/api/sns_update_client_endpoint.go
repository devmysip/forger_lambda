package api

import (
	"context"
	"fmt"
	"forger/db"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func SNSUpdateClientEndpoint(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {

	type RequestBody struct {
		FCMToken   string `json:"fcm_token"`
		EnpointARN string `json:"client_endpoint"`
	}

	email, err := headerHandler(request.Headers)
	if err != nil {
		log.Printf("Error extracting email: %s", email)

		return responseBuilder(0, nil, "Internal Server Error", "Failed to extract email from request")
	}

	body, err := decodeAndUnmarshal[RequestBody](request)
	if err != nil {
		log.Printf("Error decoding and unmarshalling request body: %s", err)
		return responseBuilder(0, nil, "Bad Request", "Failed to parse request body")
	}

	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		fmt.Println(err)
		return responseBuilder(0, nil, "Bad Request", "Failed to parse request body")

	}
	snsClient := sns.NewFromConfig(sdkConfig)

	// Assuming you get the endpoint ARN and new token from the request body or other source
	newToken := body.FCMToken

	input := &sns.SetEndpointAttributesInput{
		EndpointArn: aws.String(body.EnpointARN),
		Attributes: map[string]string{
			"Token": newToken,
		},
	}

	_, err = snsClient.SetEndpointAttributes(context.TODO(), input)
	if err != nil {
		log.Printf("Error updating endpoint: %v", err)
		return responseBuilder(0, nil, "Bad Request", err.Error())

	}
	svc := dynamodb.New(db.DB())

	query := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":fcm_token": {
				S: aws.String(body.FCMToken),
			},
			":updated_at": {
				S: aws.String(time.Now().Format(time.RFC3339)),
			},
		},
		TableName:        aws.String("User"),
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set fcm_token = :fcm_token, updated_at = :updated_at"),
	}
	_, err = svc.UpdateItem(query)
	if err != nil {
		log.Printf("Error updating FCM token: %s", err)
		return responseBuilder(0, nil, "Internal Server Error", "Failed to update FCM token")
	}

	fmt.Println("Successfully updated the endpoint")
	return responseBuilder(1, nil, "Endpoint updated successfully", "")

}
