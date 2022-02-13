package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/whiskey-back/internal/config"
	"github.com/whiskey-back/internal/types"
	"mime/multipart"
)

type FileStorage struct {
	client    *minio.Client
	Bucket    string
	BadBucket string
}

func NewMinio(cfg *config.Storage) (*FileStorage, error) {
	client, err := minio.New(
		cfg.Host,
		&minio.Options{Creds: credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""), Secure: cfg.SSL})

	if err != nil {
		return nil, err
	}

	return &FileStorage{
		client:    client,
		Bucket:    cfg.Bucket,
		BadBucket: cfg.BucketBad,
	}, nil
}

func contentTypeGet(file *multipart.FileHeader) (string, error) {
	mime := file.Header.Get("Content-Type")

	return mime, nil
}

func (fs *FileStorage) Add(file *multipart.FileHeader, bucket, userID string) (*types.File, error) {
	mime, err := contentTypeGet(file)
	if err != nil {
		return nil, err
	}

	f, err := file.Open()
	if err != nil {
		return nil, err
	}

	uploadInfo, err := fs.client.PutObject(context.Background(), bucket, fmt.Sprintf("%s/%s", userID, file.Filename), f, file.Size, minio.PutObjectOptions{
		ContentType: mime,
	})
	if err != nil {
		return nil, err
	}

	fileObject := &types.File{
		Length:   uploadInfo.Size,
		Name:     file.Filename,
		MimeType: mime,
		Bucket:   uploadInfo.Bucket,
		Object:   uploadInfo.Key,
		UserId:   userID,
	}

	return fileObject, nil
}

func (fs *FileStorage) Delete(file *types.File, bucket string) error {

	err := fs.client.RemoveObject(context.Background(), bucket, file.Object, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}
