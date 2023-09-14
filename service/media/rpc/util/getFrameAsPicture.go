package util

import (
	"bytes"
	"context"
	"doushen_by_liujun/internal/common"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// ffmpeg URL-IO模式demo
var (
	client *minio.Client
)

func initMinio() {
	fmt.Println("创建MinIO连接")
	timeout := 60 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	//创建 MinIO 客户端对象
	minioClient, err := minio.New(common.MinIOEndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(common.MinIOAccessKey, common.MinIOSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln("创建 MinIO 客户端失败", err)
		return
	}
	client = minioClient
	log.Printf("创建 MinIO 客户端成功")

	exist, err := client.BucketExists(context.Background(), common.MinIOCoverBucketName)

	if err != nil {
		log.Fatalln("查找桶失败", err)
	}

	if !exist {
		err = client.MakeBucket(context.Background(), common.MinIOCoverBucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalln("创建桶失败", err)
		}
	}
	log.Printf("桶已创建")
	//设置成公开访问
	policy := `{"Version": "2012-10-17","Statement": [{"Action": ["s3:GetObject"],"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::` + common.MinIOCoverBucketName + `/*"],"Sid": ""}]}`
	err = client.SetBucketPolicy(ctx, common.MinIOCoverBucketName, policy)
	if err != nil {
		log.Fatalln("设置公开访问策略失败", err)
	}
}

// PutPicture 上传图片
//func PutPicture(buf *bytes.Buffer, inFileName string) {
//	//上传图片
//	initMinio()
//	saveName := inFileName + ".jpg"
//	_, err := client.PutObject(context.Background(),
//		common.MinIOCoverBucketName,
//		saveName,
//		buf,
//		int64(buf.Len()),
//		minio.PutObjectOptions{
//			ContentType: "image/jpeg",
//		})
//
//	if err != nil {
//		log.Fatalln("图片上传失败", err)
//		return
//	}
//	fmt.Println("图片上传成功")
//	//获取对象预签名url
//	URL, err := client.PresignedGetObject(context.Background(), common.MinIOCoverBucketName, saveName, time.Second*24*60*60, nil)
//	if err != nil {
//		log.Fatalln("获取url失败", err)
//		return
//	}
//	fmt.Println("URL获取成功,URL为：", URL)
//}

// GetFrame 抽帧
//func GetFrame(inFileName string, frameNum int) (*bytes.Buffer, error) {
//	log.Printf("进到抽帧函数")
//	inFileName = inFileName + ".mp4"
//	buf := bytes.NewBuffer(nil)
//	err := ffmpeg.Input(common.HTTP+common.MinIOEndPoint+"/"+common.MinIOVideoBucketName+"/"+inFileName).
//		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
//		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
//		WithOutput(buf, os.Stdout).
//		Run()
//	// 运行 ffmpeg 命令
//	if err != nil {
//		log.Fatalln("截取图片失败", err)
//		return nil, err
//	}
//	log.Printf("截取图片成功")
//	return buf, nil
//}

func GetFrameByDocker(inFileName string) error {
	dir, err := os.Getwd()
	if err != nil {
		// 处理错误
	}
	fmt.Println("当前工作目录：", dir)
	videoURL := common.HTTP + common.MinIOEndPoint + "/" + common.MinIOVideoBucketName + "/" + inFileName + ".mp4"
	//videoURL := "http://8.137.50.160:9000/dousheng-video/" + inFileName + ".mp4"
	fmt.Println("videoURL: ", videoURL)
	cmd := exec.Command(
		"docker",
		"run",
		"--rm",
		"-v",
		dir+":/data",
		"jrottenberg/ffmpeg",
		"-i",
		videoURL,
		"-vf",
		"select=eq(n\\,4)",
		"-vframes",
		"1",
		"/data/"+inFileName+".jpg",
	)
	// 执行命令，并等待命令完成
	err = cmd.Run()
	if err != nil {
		fmt.Println("docker抽帧命令执行命令出错：", err)
	} else {
		fmt.Println("docker抽帧命令执行成功！")
	}
	return err
}

func PutPictureByDocker(inFileName string) {
	initMinio()
	dir, err := os.Getwd()
	if err != nil {
		// 处理错误
		fmt.Println("获取当前工作目录失败")
	}
	fmt.Println("当前工作目录：", dir)

	pictureBytes, err := os.ReadFile(inFileName + ".jpg")
	if err != nil {
		fmt.Println("无法读取视频文件:", err)
		return
	}
	// 打印视频字节切片的长度
	fmt.Println("图片字节切片长度:", len(pictureBytes))
	saveName := inFileName + ".jpg"
	//上传图片
	_, err = client.PutObject(context.Background(),
		common.MinIOCoverBucketName,
		saveName,
		bytes.NewReader(pictureBytes),
		int64(len(pictureBytes)),
		minio.PutObjectOptions{
			ContentType: "picture/jpeg",
		})
	if err != nil {
		log.Fatalln("图片上传失败", err)
		return
	}
	fmt.Println("图片上传成功")
	//获取对象预签名url
	URL, err := client.PresignedGetObject(context.Background(), common.MinIOCoverBucketName, saveName, time.Second*24*60*60, nil)
	if err != nil {
		log.Fatalln("获取url失败", err)
		return
	}
	fmt.Println("URL获取成功,URL为：", URL)
	//删除本地图片
	err = os.Remove(inFileName + ".jpg")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("临时图片文件删除成功")
	}
}
