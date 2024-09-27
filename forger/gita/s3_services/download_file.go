package s3services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"forger/db"
	"io"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func decodeAndUnmarshal[T any](request events.APIGatewayProxyRequest) (T, error) {
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
func DownloadFileFromS3[T any](bucket, objKey string) (T, error) {

	s3Client := s3.New(db.DB())
	var model T
	// Define GetObjectInput
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objKey),
	}

	// Fetch the object from S3
	result, err := s3Client.GetObject(input)
	if err != nil {
		return model, err
	}
	defer result.Body.Close()

	// Create a temp file to write the content
	tmpFile, err := os.CreateTemp("", "s3file-*")
	if err != nil {
		return model, err
	}

	// Copy the S3 object to the temp file
	_, err = io.Copy(tmpFile, result.Body)
	if err != nil {
		return model, err
	}

	// Seek to the beginning of the file before returning
	_, err = tmpFile.Seek(0, 0)
	if err != nil {
		return model, err
	}

	// Read the content of the file
	var buf bytes.Buffer
	_, err = io.Copy(&buf, tmpFile)
	if err != nil {
		return model, err
	}

	err = json.Unmarshal(buf.Bytes(), &model)
	if err != nil {
		return model, err
	}

	return model, nil
}
