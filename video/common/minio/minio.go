package minio

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"strconv"
	"videoService/config"
)

var MinioCli = NewMinioClientNoErr(config.C.MinIoConfig.Endpoint, config.C.MinIoConfig.AccessKey, config.C.MinIoConfig.SecretKey, config.C.MinIoConfig.UseSSL)

func NewMinioClient(endpoint, accessKey, secretKey string, useSSL bool) (*MinioClient, error) {

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	return &MinioClient{minioClient}, err
}
func NewMinioClientNoErr(endpoint, accessKey, secretKey string, useSSL bool) *MinioClient {

	// Initialize minio client object.
	minioClient, _ := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	return &MinioClient{minioClient}
}

type MinioClient struct {
	c *minio.Client
}

func (c *MinioClient) Get(
	ctx context.Context,
	bucket string,
	filename string) bool {
	object, err := c.c.GetObject(ctx, bucket, filename, minio.GetObjectOptions{})
	if err != nil {
		log.Println(err)
		return false
	}
	stat, err := object.Stat()
	if err != nil {
		log.Println(err)
		return false
	}
	return stat.Key != ""
}
func (mc *MinioClient) Put(
	ctx context.Context,
	bucketName string,
	fileName string,
	data []byte,
	size int64,
	contentType string,
) (minio.UploadInfo, error) {
	object, err := mc.c.PutObject(ctx, bucketName, fileName, bytes.NewBuffer(data), size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return object, err
}

func (mc *MinioClient) Compose(
	ctx context.Context,
	bucketName string,
	fileName string,
	chunkNumber int,
) (minio.UploadInfo, error) {
	dst := minio.CopyDestOptions{
		Bucket: bucketName,
		Object: fileName,
	}
	var srcs []minio.CopySrcOptions
	for i := 1; i <= chunkNumber; i++ {
		formatInd := strconv.FormatInt(int64(i), 10)
		src := minio.CopySrcOptions{
			Bucket: bucketName,
			Object: fileName + "_" + formatInd,
		}
		srcs = append(srcs, src)
	}
	object, err := mc.c.ComposeObject(
		ctx,
		dst,
		srcs...)
	return object, err
}
