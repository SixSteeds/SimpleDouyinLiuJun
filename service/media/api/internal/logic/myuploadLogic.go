package logic

//import (
//	"context"
//	"doushen_by_liujun/internal/common"
//	"doushen_by_liujun/internal/gloabalType"
//	gloabalUtil "doushen_by_liujun/internal/util"
//	"doushen_by_liujun/service/media/api/internal/svc"
//	"doushen_by_liujun/service/media/api/internal/types"
//	"doushen_by_liujun/service/media/rpc/pb"
//	"doushen_by_liujun/service/media/rpc/util"
//	"encoding/json"
//	"fmt"
//	"github.com/google/uuid"
//	"github.com/zeromicro/go-zero/core/logx"
//	"github.com/zeromicro/go-zero/core/threading"
//	"strconv"
//	"time"
//)
//
//type UploadLogic struct {
//	logx.Logger
//	ctx    context.Context
//	svcCtx *svc.ServiceContext
//}
//
//func NewUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadLogic {
//	return &UploadLogic{
//		Logger: logx.WithContext(ctx),
//		ctx:    ctx,
//		svcCtx: svcCtx,
//	}
//}
//
//func (l *UploadLogic) Upload(req *types.UploadReq) (resp *types.UploadResp, err error) {
//
//	fmt.Println("进入上传api逻辑")
//	token, err := gloabalUtil.ParseToken(req.Token)
//	if err != nil {
//		return &types.UploadResp{
//			StatusMsg:  common.MapErrMsg(common.TOKEN_PARSE_ERROR),
//			StatusCode: common.TOKEN_PARSE_ERROR,
//		}, nil
//	}
//	//生成文件名
//	fileName := strconv.FormatInt(token.UserID, 10) + uuid.New().String()[:5]
//	data, err := gloabalUtil.NewSnowflake(common.MediaApiMachineId)
//	if err != nil {
//		l.Logger.Error("雪花算法出问题：", err)
//		return &types.UploadResp{
//			StatusMsg:  common.MapErrMsg(common.SERVER_COMMON_ERROR),
//			StatusCode: common.SERVER_COMMON_ERROR,
//		}, nil
//	}
//	snowId := data.Generate()
//
//	threading.GoSafe(func() {
//		err = util.Upload(l.ctx, req.Data, fileName)
//		fmt.Println("上传成功")
//		if err != nil {
//			//l.Logger.Error("上传出问题：", err)
//		}
//		// 抽取视频第 5 帧
//		//coverData, err := util2.GetFrame(fileName, 5)
//		err = util.GetFrameByDocker(fileName)
//		if err != nil {
//			//l.Logger.Error("抽帧封面出问题：", err)
//
//		}
//		// 上传封面
//		//util2.PutPicture(coverData, fileName)
//		fmt.Println("上传封面")
//		err = util.PutPictureByDocker(fileName)
//		if err != nil {
//			//l.Logger.Error("上传封面出问题：", err)
//
//		}
//		fmt.Println("上传封面")
//	})
//
//	//kafka 推送
//	ip := l.ctx.Value("ip")
//	ipString, ok := ip.(string)
//	message := gloabalType.UploadSuccessMessage{}
//	if ok {
//		fmt.Println("sdasdasdad")
//		message.IP = ipString
//		message.Uploadtime = time.Now()
//		message.UserId = token.UserID
//		message.PlayUrl = common.HTTP + common.MinIOEndPoint + "/" + common.MinIOVideoBucketName + "/" + fileName + ".mp4"
//		message.DataLen = int64(len(req.Data))
//		messageBytes, err := json.Marshal(message)
//		if err != nil {
//			l.Logger.Error("无法序列化 message 结构体为 JSON：", err)
//		}
//		if err := l.svcCtx.UploadPersistentKqPusherConf.Push(string(messageBytes)); err != nil {
//			l.Logger.Error("upload方法kafka日志处理失败")
//		}
//	} else {
//		l.Logger.Error("nginx出问题啦")
//	}
//
//	_, err = l.svcCtx.MediaRpcClient.SaveVideo(l.ctx, &pb.SaveVideoReq{
//		UserId:   token.UserID,
//		PlayUrl:  common.HTTP + common.MinIOEndPoint + "/" + common.MinIOVideoBucketName + "/" + fileName + ".mp4",
//		CoverUrl: common.HTTP + common.MinIOEndPoint + "/" + common.MinIOCoverBucketName + "/" + fileName + ".jpg",
//		Title:    req.Title,
//		Id:       snowId,
//	})
//	if err != nil {
//		return &types.UploadResp{
//			StatusMsg:  common.MapErrMsg(common.DB_ERROR),
//			StatusCode: common.DB_ERROR,
//		}, nil
//	}
//
//	_, err = l.svcCtx.RedisClient.Incr(common.CntCacheUserWorkPrefix + strconv.FormatInt(token.UserID, 10))
//	if err != nil {
//		return &types.UploadResp{
//			StatusMsg:  common.MapErrMsg(common.REDIS_ERROR),
//			StatusCode: common.REDIS_ERROR,
//		}, nil
//	}
//	err = l.svcCtx.RedisClient.SetCtx(l.ctx, common.VideoCache2User+strconv.FormatInt(token.UserID, 10), strconv.FormatInt(snowId, 10))
//	if err != nil {
//		return &types.UploadResp{
//			StatusMsg:  common.MapErrMsg(common.REDIS_ERROR),
//			StatusCode: common.REDIS_ERROR,
//		}, nil
//	}
//	return &types.UploadResp{
//		StatusMsg:  common.MapErrMsg(common.OK),
//		StatusCode: common.OK,
//	}, nil
//}
