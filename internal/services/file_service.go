package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type FileService struct {
	BucketName  string
	MinioClient *minio.Client
}

func NewFileService(bucketName, minioUrl, minioUser, minioPasswd string) FileService {
	minioClient, err := minio.New(minioUrl, &minio.Options{
		Creds: credentials.NewStaticV4(minioUser, minioPasswd, ""),
	})
	if err != nil {
		fmt.Println(err)
		log.Fatalf("Failed to initialize minio client: %v", err)
	}
	return FileService{
		BucketName:  bucketName,
		MinioClient: minioClient,
	}
}

func (fs FileService) CreateBucketIfNotExists(ctx context.Context) bool {
	ok, err := fs.MinioClient.BucketExists(ctx, fs.BucketName)
	if err != nil {
		fmt.Println(err)
		log.Fatalf("Failed to acess minio: %v", err)
	}
	if !ok {
		if err := fs.MinioClient.MakeBucket(ctx, fs.BucketName, minio.MakeBucketOptions{}); err != nil {
			fmt.Println(err)
			log.Fatalf("Failed to create bucket: %v", err)
		}
		return true
	}
	return false
}

func (fs FileService) PutObject(ctx context.Context, fileName string, file io.Reader, size int64) (minio.UploadInfo, error) {
	info, err := fs.MinioClient.PutObject(ctx, fs.BucketName, fileName, file, size, minio.PutObjectOptions{})
	if err != nil {
		return minio.UploadInfo{}, err
	}
	return info, nil
}

func (fs FileService) GetSignedURL(ctx context.Context, key, bucket string) (*url.URL, error) {
	url, err := fs.MinioClient.PresignedGetObject(ctx, bucket, key, time.Minute*15, url.Values{})
	if err != nil {
		return nil, errors.New("failed to generate signed url")
	}

	return url, nil
}
