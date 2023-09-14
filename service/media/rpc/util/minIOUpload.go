package util

import (
	"bytes"
	"context"
	"doushen_by_liujun/internal/common"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"math"
	"time"
)

// Upload 视频上传
func Upload(_ context.Context, video []byte, fileName string) error {
	fileName = fileName + ".mp4"
	timeout := 60 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	//初始化客户端
	client, err := minio.NewCore(common.MinIOEndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(common.MinIOAccessKey, common.MinIOSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return err
	}
	//检查桶是否存在，不存在就创建
	exists, err := client.BucketExists(ctx, common.MinIOVideoBucketName)
	if err != nil {
		return err
	}
	if !exists {
		err = client.MakeBucket(ctx, common.MinIOVideoBucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
		fmt.Printf("Bucket %s created successfully\n", common.MinIOVideoBucketName)
	}
	//设置成公开访问
	policy := `{"Version": "2012-10-17","Statement": [{"Action": ["s3:GetObject"],"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::` + common.MinIOVideoBucketName + `/*"],"Sid": ""}]}`
	err = client.SetBucketPolicy(ctx, common.MinIOVideoBucketName, policy)
	if err != nil {
		return err
	}
	if len(video) <= 5*1024*1024 {
		// 文件大小小于等于 5MB，直接上传整个视频
		_, err = client.PutObject(ctx, common.MinIOVideoBucketName, fileName, bytes.NewReader(video), int64(len(video)), "", "", minio.PutObjectOptions{ContentType: "video/mp4"})
		if err != nil {
			return err
		}
	} else {
		// 文件大小大于 5MB，使用分片上传
		err = uploadInParts(ctx, client, video, common.MinIOVideoBucketName, fileName)
		if err != nil {
			return err
		}
	}
	//文件形式上传
	//_, err = client.FPutObject(ctx, bucketName, fileName, filePath, minio.PutObjectOptions{})
	objectURL := common.HTTP + common.MinIOEndPoint + "/" + common.MinIOVideoBucketName + "/" + fileName
	log.Println(objectURL)
	return nil
}
func uploadInParts(ctx context.Context, client *minio.Core, videoData []byte, bucketName, fileName string) error {

	partSize := int64(5 * 1024 * 1024) // 5MB
	totalParts := int(math.Ceil(float64(len(videoData)) / float64(partSize)))

	// 开始分片上传
	uploadID, err := client.NewMultipartUpload(ctx, bucketName, fileName, minio.PutObjectOptions{ContentType: "video/mp4"})
	if err != nil {
		return err
	}

	var uploadedParts []minio.CompletePart
	for partNumber := 1; partNumber <= totalParts; partNumber++ {
		start := (partNumber - 1) * int(partSize)
		end := int64(start) + partSize
		if end > int64(len(videoData)) {
			end = int64(len(videoData))
		}

		partData := videoData[start:end]
		// 上传分片
		part, err := client.PutObjectPart(ctx, common.MinIOVideoBucketName, fileName, uploadID, partNumber, bytes.NewReader(partData), int64(len(partData)), minio.PutObjectPartOptions{})
		if err != nil {
			return err
		}
		uploadedParts = append(uploadedParts, minio.CompletePart{PartNumber: partNumber, ETag: part.ETag})
	}

	// 完成分片上传
	_, err = client.CompleteMultipartUpload(ctx, bucketName, fileName, uploadID, uploadedParts, minio.PutObjectOptions{ContentType: "video/mp4"})
	if err != nil {
		return err
	}

	return nil
}
