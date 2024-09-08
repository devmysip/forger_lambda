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

const platformApplicationArn = "arn:aws:sns:ap-south-1:533267104785:app/GCM/com.gitasarathi"

func SNSCreateClientEndpoint(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {

	type RequestBody struct {
		FCMToken string `json:"fcm_token"`
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

	input := &sns.CreatePlatformEndpointInput{
		Token:                  aws.String(body.FCMToken),
		PlatformApplicationArn: aws.String(platformApplicationArn),
	}

	resp, err := snsClient.CreatePlatformEndpoint(context.TODO(), input)
	if err != nil {
		log.Printf("Error: %v", err)
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
			":client_endpoint": {
				S: aws.String(*resp.EndpointArn),
			},
			":updated_at": {
				S: aws.String(time.Now().Format(time.RFC3339)),
			},
		},
		TableName:        aws.String("User"),
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set fcm_token = :fcm_token, client_endpoint =:client_endpoint, updated_at = :updated_at"),
	}

	_, err = svc.UpdateItem(query)
	if err != nil {
		log.Printf("Error updating FCM token: %s", err)
		return responseBuilder(0, nil, "Internal Server Error", err.Error())
	}

	fmt.Println("The ARN of the endpoint is", *resp.EndpointArn)
	return responseBuilder(1, *resp.EndpointArn, "Success", "")

}
