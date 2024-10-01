package userengagement

import (
	"fmt"
	"forger/db"
	"forger/gita/constants"
	"forger/gita/models"
	s3services "forger/gita/s3_services"
	"forger/gita/utilis"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func GetActiveUserInDays(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {

	type RequestBody struct {
		Days           int   `json:"days"`
		UploadFileToS3 *bool `json:"upload_file_to_s3"`
	}

	_, err := utilis.HeaderHandler(request.Headers)
	if err != nil {
		return utilis.ResponseBuilder(0, nil, "Internal Server Error", "Failed to extract email from request")
	}

	body, err := utilis.DecodeAndUnmarshal[RequestBody](request)
	if err != nil {
		return utilis.ResponseBuilder(0, nil, "Bad Request", "Failed to parse request body")
	}

	svc := dynamodb.New(db.DB())

	now := utilis.GetLocalTime()
	endDate := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location()).Format("2006-01-02")
	startDate := time.Date(now.Year(), now.Month(), now.Day()-body.Days, now.Hour(), 0, 0, 0, now.Location()).Format("2006-01-02")

	scanInput := &dynamodb.ScanInput{
		TableName:        aws.String("User"),
		FilterExpression: aws.String("updated_at BETWEEN :start_date AND :end_date"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":start_date": {
				S: aws.String(startDate),
			},
			":end_date": {
				S: aws.String(endDate),
			},
		},
		ProjectionExpression: aws.String("email, client_endpoint, updated_at, display_name, last_read"),
	}

	result, err := svc.Scan(scanInput)
	if err != nil {
		fmt.Println("Error scanning table:", err)
		return utilis.ResponseBuilder(0, nil, "Failed", err.Error())
	}

	var users []models.User
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return utilis.ResponseBuilder(0, users, "Failed", "")
	}

	if body.UploadFileToS3 != nil && *body.UploadFileToS3 {
		tempFile, err := utilis.CreateTempJSONFile(users)
		if err != nil {
			return utilis.ResponseBuilder(0, nil, "Failed to create temp file", err.Error())
		}

		defer tempFile.Close()

		objectKey := fmt.Sprintf("%s/%s.json", constants.DailyNotifictionUserBucketDirectory, now.Format("2006-01-02"))
		if err := s3services.UploadFileToS3(tempFile.Name(), constants.GitaSarathiBucket, objectKey); err != nil {
			return utilis.ResponseBuilder(0, nil, "Failed to upload file to S3", err.Error())
		}

		return utilis.ResponseBuilder(1, nil, "success", "")
	}

	return utilis.ResponseBuilder(1, users, "success", "")

}
