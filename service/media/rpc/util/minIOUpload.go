package util

import (
	"bytes"
	"context"
	"doushen_by_liujun/internal/common"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

// Upload 视频上传
func Upload(ctx context.Context, video []byte, fileName string) error {
	useSSL := false
	//初始化客户端
	client, err := minio.New(common.UploadPath, &minio.Options{
		Creds:  credentials.NewStaticV4(common.MinIOAccessKey, common.MinIOSecretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return err
	}
	//检查桶是否存在，不存在就创建
	exists, err := client.BucketExists(ctx, common.MinIOBucketName)
	if err != nil {
		return err
	}
	if !exists {
		err = client.MakeBucket(ctx, common.MinIOBucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
		fmt.Printf("Bucket %s created successfully\n", common.MinIOBucketName)
	}
	//设置成公开访问
	policy := `{"Version": "2012-10-17","Statement": [{"Action": ["s3:GetObject"],"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::` + common.MinIOBucketName + `/*"],"Sid": ""}]}`
	err = client.SetBucketPolicy(ctx, common.MinIOBucketName, policy)
	if err != nil {
		return err
	}
	//字节流形式上传
	_, err = client.PutObject(ctx, common.MinIOBucketName, fileName, bytes.NewReader(video), int64(len(video)), minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	//文件形式上传
	//_, err = client.FPutObject(ctx, bucketName, fileName, filePath, minio.PutObjectOptions{})

	log.Println("文件访问地址为")
	objectURL := "http://" + common.UploadPath + "/" + common.MinIOBucketName + "/" + fileName
	log.Println(objectURL)
	return nil
}
