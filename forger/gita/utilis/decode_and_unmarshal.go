package utilis

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

func DecodeAndUnmarshal[T any](request events.APIGatewayProxyRequest) (T, error) {
	var targetStruct T

	// Decode the base64-encoded body
	decodedBody, err := base64.StdEncoding.DecodeString(request.Body)
	if err != nil {
		return targetStruct, fmt.Errorf("error decoding base64 body: %w", err)
	}

	// Unmarshal the decoded body into the target struct
	err = json.Unmarshal(decodedBody, &targetStruct)
	if err != nil {
		return targetStruct, fmt.Errorf("error unmarshalling JSON body to struct: %w", err)
	}

	return targetStruct, nil
}
