package utilis

import (
	"fmt"
	"forger/db"
	"forger/gita/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func UpdateTime(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {

	svc := dynamodb.New(db.DB())

	scanInput := &dynamodb.ScanInput{
		TableName:            aws.String("User"),
		ProjectionExpression: aws.String("email, client_endpoint, updated_at, display_name, last_read"),
	}

	result, err := svc.Scan(scanInput)
	if err != nil {
		fmt.Println("Error scanning table:", err)
		return ResponseBuilder(0, nil, "Failed", err.Error())
	}

	var users []models.User
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return ResponseBuilder(0, users, "Failed", "")
	}

	// for _, user := range users {

	// }

	return ResponseBuilder(1, users, "Success", "")
}
