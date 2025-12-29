package utils

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	client *minio.Client
	bucket string
}

func NewMinioClient(endpoint, accessKey, secretKey, bucket string) (*MinioClient, error) {
	var lastErr error

	for i := 0; i < 10; i++ {
		client, err := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
			Secure: false,
		})
		if err == nil {
			ctx := context.Background()

			if _, err = client.ListBuckets(ctx); err == nil {
				exists, err := client.BucketExists(ctx, bucket)
				if err == nil {
					if !exists {
						err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
					}
					if err == nil {
						log.Println("MinIO ready")
						return &MinioClient{
							client: client,
							bucket: bucket,
						}, nil
					}
				}
			}
		}

		lastErr = err
		log.Printf("Waiting for MinIO (%d/10): %v", i+1, lastErr)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("minio not available after retries: %w", lastErr)
}

func (m *MinioClient) PutObject(ctx context.Context, name string, data []byte) error {
	_, err := m.client.PutObject(
		ctx,
		m.bucket,
		name,
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{
			ContentType: "application/json",
		},
	)
	return err
}
