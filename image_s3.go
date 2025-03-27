package epubtomd

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gabriel-vasile/mimetype"
)

type S3ImageHandler struct {
	BucketName      string `json:"bucket_name"`
	AccountID       string `json:"account_id"`
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	Endpoint        string `json:"end_point"`
	Domain          string `json:"domain"`
	client          *s3.Client
	f               fs.FS
}

func NewS3ImageHandler(f fs.FS, bucketName, accountId, accessKeyId, accessKeySecret, endpoint, domain string) (ImageHandler, error) {
	customTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}
	customHTTPClient := &http.Client{
		Transport: customTransport,
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion("auto"),
		config.WithHTTPClient(customHTTPClient),
		config.WithRetryMaxAttempts(2),
	)
	if err != nil {
		log.Println(bucketName, accountId, accessKeyId, accessKeySecret, err)
		return nil, fmt.Errorf("loading AWS S3 bucket err: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = &endpoint
	})
	return &S3ImageHandler{
		BucketName:      bucketName,
		AccountID:       accountId,
		AccessKeyID:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		Domain:          domain,
		client:          client,
		f:               f,
	}, nil
}

func (c *S3ImageHandler) CopyImage(srcImagePath string, remoteFilename string) (string, error) {
	return c.CopyWithRename(srcImagePath, func(b []byte) string {
		return remoteFilename
	})
}

func (c *S3ImageHandler) CopyWithRename(srcImagePath string, namePathFn func(b []byte) string) (string, error) {
	file, err := c.f.Open(srcImagePath)
	if err != nil {
		return "", fmt.Errorf("error opening file %s: %w", srcImagePath, err)
	}

	mimeType, err := mimetype.DetectReader(file)

	if err != nil {
		return "", fmt.Errorf("detect file %s error: %w", srcImagePath, err)
	}

	body, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("error reading file %s: %w", srcImagePath, err)
	}
	remoteFilename := namePathFn(body)

	putObjectInput := &s3.PutObjectInput{
		Bucket:      &c.BucketName,
		Key:         aws.String(remoteFilename),
		Body:        file,
		ContentType: aws.String(mimeType.String()),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	_, pErr := c.client.PutObject(ctx, putObjectInput)
	if pErr != nil {
		return "", fmt.Errorf("upload file %s error: %w", remoteFilename, pErr)
	}
	return c.Domain + strings.ReplaceAll("/"+remoteFilename, "//", "/"), nil
}
