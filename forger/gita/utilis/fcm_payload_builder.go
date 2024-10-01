package utilis

import (
	"encoding/json"
	"fmt"
)

func FCMPayloadBuilder(title, body string, data map[string]interface{}) (string, error) {

	// message := `{
	//   "GCM": "{\"notification\":{\"title\":\"Hello\",\"body\":\"World\"},\"data\":{\"screen\":\"/chaptersDetail\",\"arguments\":{\"chapter_no\":1,\"verse_no\":1}}}",
	// "APNS": "{\"aps\": {\"alert\": \"%s\", \"sound\": \"default\"}}"
	// 	}`

	// Convert the data map to a JSON string
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "nil", err
	}

	// Create the GCM and APNS strings
	gcm := fmt.Sprintf(`{"notification":{"title":"%s","body":"%s"},"data":%s}`, title, body, string(dataJson))
	apns := fmt.Sprintf(`{"aps": {"alert": "%s", "sound": "default"}}`, body)

	// Format the full message
	message := fmt.Sprintf(`{
	  "GCM": %q,
	  "APNS": %q
	}`, gcm, apns)

	return message, nil
}
