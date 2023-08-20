package logic

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/service/user/rpc/internal/svc"
	"doushen_by_liujun/service/user/rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"strconv"
)

type GetUserinfoByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserinfoByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserinfoByIdLogic {
	return &GetUserinfoByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserinfoByIdLogic) GetUserinfoById(in *pb.GetUserinfoByIdReq) (*pb.GetUserinfoByIdResp, error) {
	// todo: add your logic here and delete this line
	//userId查id这个人

	info, err := l.svcCtx.UserinfoModel.FindOne(l.ctx, in.Id, in.UserID)
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.KqPusherClient.Push("user_rpc_getUserinfoByIdLogic_GetUserinfoById_FindOne_false"); err != nil {
		log.Fatal(err)
	}
	redisClient := l.svcCtx.RedisClient
	followKey := common.FollowNum + strconv.Itoa(int(in.Id))
	followerKey := common.FollowerNum + strconv.Itoa(int(in.Id))
	followRecord, _ := redisClient.GetCtx(l.ctx, followKey)
	followNum := 0
	followerNum := 0
	expiration := 3600 //秒
	workCount := 0
	favoriteCount := 0
	totalFavorited := 0
	workCountRecord, _ := redisClient.GetCtx(l.ctx, common.CntCacheUserWorkPrefix+strconv.Itoa(int(in.Id)))
	if len(workCountRecord) != 0 { //等于0 代表没有记录，直接赋值0
		//有记录
		workCount, _ = strconv.Atoi(workCountRecord)
	}
	favoriteCountRecord, _ := redisClient.GetCtx(l.ctx, common.CntCacheUserLikePrefix+strconv.Itoa(int(in.Id)))
	if len(favoriteCountRecord) != 0 { //等于0 代表没有记录，直接赋值0
		//有记录
		favoriteCount, _ = strconv.Atoi(favoriteCountRecord)
	}
	totalFavoritedRecord, _ := redisClient.GetCtx(l.ctx, common.CntCacheUserLikedPrefix+strconv.Itoa(int(in.Id)))
	if len(totalFavoritedRecord) != 0 { //等于0 代表没有记录，直接赋值0
		//有记录
		totalFavorited, _ = strconv.Atoi(totalFavoritedRecord)
	}
	if len(followRecord) == 0 {
		//没有记录，去查表
		num, err := l.svcCtx.FollowsModel.FindFollowsCount(l.ctx, in.Id)
		num = num - 1 //剪掉自己
		if err != nil {
			return nil, err
		}
		_ = redisClient.SetexCtx(l.ctx, followKey, strconv.Itoa(num), expiration)
		followNum = num
	} else {
		//有记录
		followNum, _ = strconv.Atoi(followRecord)
	}
	followerRecord, _ := redisClient.GetCtx(l.ctx, followerKey)
	if len(followerRecord) == 0 {
		//没有记录，去查表
		num, err := l.svcCtx.FollowsModel.FindFollowersCount(l.ctx, in.Id)
		num = num - 1 //剪掉自己
		if err != nil {
			return nil, err
		}
		_ = redisClient.SetexCtx(l.ctx, followerKey, strconv.Itoa(num), expiration)
		followerNum = num
	} else {
		//有记录
		followerNum, _ = strconv.Atoi(followerRecord)
	}
	userInfo := pb.Userinfo{
		Id:              info.Id,
		FollowCount:     int64(followNum),
		FollowerCount:   int64(followerNum),
		IsFollow:        info.IsFollow,
		Username:        info.Username.String,
		Avatar:          info.Avatar.String,
		BackgroundImage: info.BackgroundImage.String,
		Signature:       info.Signature.String,
		WorkCount:       int64(workCount),
		FavoriteCount:   int64(favoriteCount),
		TotalFavorited:  int64(totalFavorited),
	}
	if err := l.svcCtx.KqPusherClient.Push("user_rpc_getUserinfoByIdLogic_GetUserinfoById_success"); err != nil {
		log.Fatal(err)
	}
	return &pb.GetUserinfoByIdResp{
		Userinfo: &userInfo,
	}, nil
}
