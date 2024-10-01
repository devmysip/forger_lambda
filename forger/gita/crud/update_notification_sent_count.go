package crud

import (
	"fmt"
	"forger/gita/constants"
	"forger/gita/utilis"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func UpdateNotificationSentCount(svc *dynamodb.DynamoDB, notificationSent int) error {

	date := utilis.GetLocalTime().Format("2006-01-02")

	tableName := constants.AnalyticsTable

	updateExpression := "SET notification_sent = if_not_exists(notification_sent, :zero) + :inc"

	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":inc": {
			N: aws.String(fmt.Sprintf("%d", notificationSent)),
		},
		":zero": {
			N: aws.String("0"),
		},
	}

	key := map[string]*dynamodb.AttributeValue{
		"date": {
			S: aws.String(date),
		},
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       key,
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: expressionAttributeValues,
		ReturnValues:              aws.String("UPDATED_NEW"),
	}

	_, err := svc.UpdateItem(input)
	if err != nil {
		return fmt.Errorf("failed to update notificationSent count: %w", err)
	}

	return nil
}
