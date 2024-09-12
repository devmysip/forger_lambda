package utilis

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go/aws"
)

func SendNotification(clientEndpoint, message string) error {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		fmt.Println(err)
		return err

	}

	snsClient := sns.NewFromConfig(sdkConfig)

	input := &sns.PublishInput{
		Message:          aws.String(message),
		TargetArn:        aws.String(clientEndpoint),
		MessageStructure: aws.String("json"),
	}

	// Send the push notification
	_, err = snsClient.Publish(context.TODO(), input)
	if err != nil {
		log.Printf("Failed to send SNS push notification: %v", err)
		return err
	}

	return nil
}
