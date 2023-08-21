package logic

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/util"
	"doushen_by_liujun/service/media/api/internal/svc"
	"doushen_by_liujun/service/media/api/internal/types"
	"doushen_by_liujun/service/media/rpc/pb"
	util2 "doushen_by_liujun/service/media/rpc/util"
	"fmt"
	"github.com/google/uuid"
	"strconv"

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

	fmt.Println("进入上传api逻辑")
	token, err := util.ParseToken(req.Token)
	if err != nil {
		return &types.UploadResp{
			StatusMsg:  "token解析失败",
			StatusCode: 500,
		}, nil
	}
	//生成文件名
	fileName := "doushen/" + strconv.FormatInt(token.UserID, 10) + uuid.New().String()[:5] + ".mp4"
	err = util2.Upload(l.ctx, req.Data, fileName)
	fmt.Println("上传成功")
	if err != nil {
		return &types.UploadResp{
			StatusMsg:  "上传失败",
			StatusCode: 500,
		}, nil
	}
	data, err := util.NewSnowflake(common.MediaApiMachineId)
	if err != nil {
		return &types.UploadResp{
			StatusMsg:  "雪花算法报错",
			StatusCode: 500,
		}, nil
	}
	snowId := data.Generate()
	_, err = l.svcCtx.MediaRpcClient.SaveVideo(l.ctx, &pb.SaveVideoReq{
		UserId:   token.UserID,
		PlayUrl:  common.QiliuyunDomain + fileName,
		CoverUrl: common.DefaultBackImage,
		Title:    req.Title,
		Id:       snowId,
	})
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.RedisClient.Incr(common.CntCacheUserWorkPrefix + strconv.FormatInt(token.UserID, 10))
	if err != nil {
		return &types.UploadResp{
			StatusMsg:  common.MapErrMsg(common.REDIS_ERROR),
			StatusCode: common.REDIS_ERROR,
		}, nil
	}
	err = l.svcCtx.RedisClient.SetCtx(l.ctx, common.VideoCache2User+strconv.FormatInt(token.UserID, 10), strconv.FormatInt(snowId, 10))
	if err != nil {
		return &types.UploadResp{
			StatusMsg:  common.MapErrMsg(common.REDIS_ERROR),
			StatusCode: common.REDIS_ERROR,
		}, nil
	}
	return &types.UploadResp{
		StatusMsg:  "上传成功",
		StatusCode: common.OK,
	}, nil
}
