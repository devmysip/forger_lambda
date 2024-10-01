package utilis

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func ResponseBuilder(status int, result interface{}, message, errorMessage string) events.APIGatewayProxyResponse {
	response := map[string]interface{}{
		"status":  status,
		"result":  result,
		"message": message,
	}

	if status == 0 {
		response["error_message"] = errorMessage
	}

	responseBody, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: string(responseBody),
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(responseBody),
	}
}
