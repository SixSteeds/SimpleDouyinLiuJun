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
func Upload(noUse context.Context, video []byte, fileName string) error {
	timeout := 60 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	fmt.Println("进来MINIO啦！！！！！！！")
	useSSL := false
	//初始化客户端
	client, err := minio.NewCore(common.UploadPath, &minio.Options{
		Creds:  credentials.NewStaticV4(common.MinIOAccessKey, common.MinIOSecretKey, ""),
		Secure: useSSL,
	})
	fmt.Println("正确1")
	if err != nil {
		fmt.Println("错误1")
		return err
	}
	//检查桶是否存在，不存在就创建
	exists, err := client.BucketExists(ctx, common.MinIOBucketName)
	fmt.Println("正确2")
	if err != nil {
		fmt.Println("错误2")
		fmt.Println(err)
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
	fmt.Println("正确3")
	if err != nil {
		fmt.Println("错误3")
		return err
	}
	fmt.Println("视频大小", len(video), 5*1024*1024)
	if len(video) <= 5*1024*1024 {
		fmt.Println("普通上传啦！！！！！！！")
		// 文件大小小于等于 5MB，直接上传整个视频
		_, err = client.PutObject(ctx, common.MinIOBucketName, fileName, bytes.NewReader(video), int64(len(video)), "", "", minio.PutObjectOptions{ContentType: "video/mp4"})
		if err != nil {
			fmt.Println("普通上传出错", err)
			return err
		}
	} else {
		// 文件大小大于 5MB，使用分片上传
		fmt.Println("分片上传啦！！！！！！！")
		err = uploadInParts(ctx, client, video, common.MinIOBucketName, fileName)
		if err != nil {
			fmt.Println("分片上传出错", err)
			return err
		}
	}
	//文件形式上传
	//_, err = client.FPutObject(ctx, bucketName, fileName, filePath, minio.PutObjectOptions{})
	log.Println("文件访问地址为")
	objectURL := "http://" + common.UploadPath + "/" + common.MinIOBucketName + "/" + fileName
	log.Println(objectURL)
	fmt.Println("MinIO上传成功啦")
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
		part, err := client.PutObjectPart(ctx, common.MinIOBucketName, fileName, uploadID, partNumber, bytes.NewReader(partData), int64(len(partData)), minio.PutObjectPartOptions{})
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

	fmt.Printf("Uploaded %s using multipart upload\n", fileName)
	return nil
}
