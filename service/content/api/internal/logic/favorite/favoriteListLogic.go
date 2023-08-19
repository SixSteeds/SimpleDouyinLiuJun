package favorite

import (
	"context"
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/service/content/rpc/pb"
	"fmt"
	"log"

	"doushen_by_liujun/service/content/api/internal/svc"
	"doushen_by_liujun/service/content/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FavoriteListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFavoriteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteListLogic {
	return &FavoriteListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FavoriteListLogic) FavoriteList(req *types.FavoriteListReq) (resp *types.FavoriteListResp, err error) {

	//1.根据 user_id 查询 favorite 表，返回点赞的所有 video_id
	favoriteListResp, err := l.svcCtx.ContentRpcClient.SearchFavorite(l.ctx, &pb.SearchFavoriteReq{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.KqPusherClient.Push("content_api_favorite_favoriteListLogic_FavoriteList_SearchFavorite_false"); err != nil {
		log.Fatal(err)
	}

	// TODO 2.依次根据 video_id 查询 video 表，返回点赞的所有视频（需要调别人的接口）
	var favoriteList = favoriteListResp.GetFavorite()
	var videoList []*pb.Video
	for _, item := range favoriteList {
		l.svcCtx.ContentRpcClient.GetVideoById(l.ctx, &pb.GetVideoByIdReq{
			Id: item.VideoId,
		})
		videoList = append(videoList, &pb.Video{})
	}

	logx.Error("api-查询用户点赞列表成功")
	fmt.Print(favoriteList)
	if err := l.svcCtx.KqPusherClient.Push("content_api_favorite_favoriteListLogic_FavoriteList_success"); err != nil {
		log.Fatal(err)
	}
	return &types.FavoriteListResp{
		StatusCode: common.OK,
		StatusMsg:  common.MapErrMsg(common.OK),
		// TODO vedio_list
	}, nil
}
