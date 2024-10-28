package utils

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func UploadToS3(fileHeader *multipart.FileHeader, showId string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("REGION")))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Initialize the S3 client
	s3Client := s3.NewFromConfig(cfg)
	bucketName := os.Getenv("BUCKET_NAME")
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	key := fmt.Sprintf("uploads/%s/%s", showId, fileHeader.Filename)

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	})
	if err != nil {
		return "", err
	}
	fileUrl := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, os.Getenv("REGION"), key)
	return fileUrl, nil
}
