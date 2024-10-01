package userengagement

import (
	"fmt"
	"forger/gita/constants"
	"forger/gita/models"
	"forger/gita/utilis"
	"sort"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func UpdateUserRead(request events.APIGatewayProxyRequest, svc *dynamodb.DynamoDB) events.APIGatewayProxyResponse {

	email, err := utilis.HeaderHandler(request.Headers)
	if err != nil {
		return utilis.ResponseBuilder(0, nil, "Internal Server Error", "Failed to extract email from request")
	}

	updateRead, err := utilis.DecodeAndUnmarshal[models.UpdateRead](request)
	if err != nil {
		return utilis.ResponseBuilder(0, nil, "Bad Request", "Failed to parse request body")
	}

	// Fetch current user data from DynamoDB
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(constants.UserTable),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
	})
	if err != nil {
		return utilis.ResponseBuilder(0, nil, "Internal Server Error", "Failed to fetch user data")
	}

	if result.Item == nil {
		return utilis.ResponseBuilder(0, nil, "Not Found", "User not found")
	}

	// Unmarshal the user's data into a User struct
	var user models.User
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return utilis.ResponseBuilder(0, nil, "Internal Server Error", "Failed to unmarshal user data")
	}

	user.UpdatedAt = utilis.GetCurrentTime()

	for i := range user.Reads {
		if user.Reads[i].Chapter == updateRead.ChapterNo {
			// Check if the verse number is already present
			verseExists := false
			for _, verse := range user.Reads[i].Verses {
				if verse == updateRead.VerseNo {
					verseExists = true
					break
				}
			}

			// Append the verse number if it does not exist

			lastRead := fmt.Sprintf("BG%d.%d", updateRead.ChapterNo, updateRead.VerseNo)
			if !verseExists {
				user.Reads[i].Verses = append(user.Reads[i].Verses, updateRead.VerseNo)
				sort.Ints(user.Reads[i].Verses)
				user.LastRead = &lastRead
				user.Reads[i].Progress = len(user.Reads[i].Verses) * 100 / models.GitaChapters[i]
			}
			break
		}
	}

	

	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return utilis.ResponseBuilder(0, nil, "Internal Server Error", "Failed to marshal updated user data")
	}

	av["email"] = &dynamodb.AttributeValue{
		S: aws.String(user.Email),
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(constants.UserTable),
		Item:      av,
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return utilis.ResponseBuilder(0, nil, "Internal Server Error", "Failed to update user data")
	}

	UpdateUserActivity(request)



	return utilis.ResponseBuilder(1, user, "Success", "User reads updated successfully")
}
