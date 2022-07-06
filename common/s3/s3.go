package s3

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Uploader struct {
	*s3manager.Uploader
}

func (s *S3Uploader) Upload(ctx context.Context, fileName, bucket string, content []byte) error {
	contentType := http.DetectContentType(content)
	// TODO: check if need to set expiry time and ACL
	// My guess definitely yes
	uploadInput := &s3manager.UploadInput{
		Key:         aws.String(fileName),
		Body:        bytes.NewReader(content),
		Bucket:      aws.String(bucket),
		ContentType: &contentType,
	}

	loc, err := s.UploadWithContext(ctx, uploadInput)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, loc.Location)
	return nil
}

func NewS3Uploader() (*S3Uploader, error) {
	accessKeyID := os.Getenv("ACCESSKEYID")
	secretAccessKey := os.Getenv("SECRETACCESSKEY")
	token := ""
	region := os.Getenv("REGION")
	creds := credentials.NewStaticCredentials(accessKeyID, secretAccessKey, token)
	_, err := creds.Get()
	if err != nil {
		return nil, err
	}
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: creds,
	})
	if err != nil {
		return nil, err
	}

	return &S3Uploader{
		s3manager.NewUploader(sess),
	}, nil
}
