package logic

import (
	"context"
	"doushen_by_liujun/service/content/rpc/internal/svc"
	"doushen_by_liujun/service/content/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFeedListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFeedListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFeedListLogic {
	return &GetFeedListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFeedListLogic) GetFeedList(in *pb.FeedListReq) (*pb.FeedListResp, error) {
	// todo: add your logic here and delete this line
	feedList, err := l.svcCtx.VideoModel.GetFeedList(l.ctx, in.UserId, in.LatestTime, in.Size)
	// 将feedlist中的userId全部拿出来转换为一个数组
	var userIds []int64
	for _, feed := range feedList {
		userIds = append(userIds, feed.UserId)
	}
	// 通过userIds获取到所有的user信息

	if err != nil {
		return nil, err
	}

	return &pb.FeedListResp{}, nil
}
