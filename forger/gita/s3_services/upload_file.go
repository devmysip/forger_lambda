package s3services

import (
	"fmt"
	"forger/db"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func UploadFileToS3(filePath string, bucket string, objectKey string) error {

	s3Client := s3.New(db.DB())

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	input := &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(objectKey),
		Body:          file,
		ContentType:   aws.String("application/octet-stream"),
		ContentLength: aws.Int64(fileInfo.Size()),
	}

	_, err = s3Client.PutObject(input)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully uploaded %s to %s/%s\n", filePath, bucket, objectKey)
	return nil
}
