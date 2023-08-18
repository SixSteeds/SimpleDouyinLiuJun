package logic

import (
	"context"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/media/api/internal/svc"
	"doushen_by_liujun/service/media/api/internal/types"
	"doushen_by_liujun/service/media/rpc/pb"
	"errors"
	"fmt"
	"github.com/google/uuid"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadLogic {
	return &UploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadLogic) Upload(req *types.UploadReq) (resp *types.UploadResp, err error) {
	// todo: add your logic here and delete this line
	fmt.Println("uploadHandler err1")
	fmt.Println(req.Token)
	token, err := util.ParseToken(req.Token)
	if err != nil {
		return &types.UploadResp{
			StatusMsg:  "token解析失败",
			StatusCode: 500,
		}, nil
	}
	fmt.Println("token解析成功", token)
	//生成文件名
	fileName := req.Title + uuid.New().String() + ".mp4"
	_, err = util.Upload(l.ctx, req.Data, fileName)
	fmt.Println("上传成功")

	if err != nil {
		return &types.UploadResp{
			StatusMsg:  "上传失败",
			StatusCode: 500,
		}, nil
	}
	data, err := util.NewSnowflake(3)
	if err != nil {
		l.Logger.Info("雪花算法报错", err)
		return nil, errors.New("雪花算法报错")
	}
	_, err = l.svcCtx.MediaRpcClient.SaveVideo(l.ctx, &pb.SaveVideoReq{
		UserId:   token.UserID,
		PlayUrl:  "xxx",
		CoverUrl: "xxx",
		Title:    req.Title,
		Id:       data.Generate(),
	})
	if err != nil {
		return nil, err
	}
	return &types.UploadResp{
		StatusMsg:  "上传成功",
		StatusCode: 0,
	}, nil
}
