// Code generated by goctl. DO NOT EDIT.
// Source: user.proto

package user

import (
	"context"

	"doushen_by_liujun/service/user/rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	AddFollowsReq       = pb.AddFollowsReq
	AddFollowsResp      = pb.AddFollowsResp
	AddUserinfoReq      = pb.AddUserinfoReq
	AddUserinfoResp     = pb.AddUserinfoResp
	DelFollowsReq       = pb.DelFollowsReq
	DelFollowsResp      = pb.DelFollowsResp
	DelUserinfoReq      = pb.DelUserinfoReq
	DelUserinfoResp     = pb.DelUserinfoResp
	Follows             = pb.Follows
	GetFollowsByIdReq   = pb.GetFollowsByIdReq
	GetFollowsByIdResp  = pb.GetFollowsByIdResp
	GetUserinfoByIdReq  = pb.GetUserinfoByIdReq
	GetUserinfoByIdResp = pb.GetUserinfoByIdResp
	SearchFollowsReq    = pb.SearchFollowsReq
	SearchFollowsResp   = pb.SearchFollowsResp
	SearchUserinfoReq   = pb.SearchUserinfoReq
	SearchUserinfoResp  = pb.SearchUserinfoResp
	UpdateFollowsReq    = pb.UpdateFollowsReq
	UpdateFollowsResp   = pb.UpdateFollowsResp
	UpdateUserinfoReq   = pb.UpdateUserinfoReq
	UpdateUserinfoResp  = pb.UpdateUserinfoResp
	Userinfo            = pb.Userinfo

	User interface {
		// -----------------------鐢ㄦ埛鍩烘湰淇℃伅-----------------------
		AddFollows(ctx context.Context, in *AddFollowsReq, opts ...grpc.CallOption) (*AddFollowsResp, error)
		UpdateFollows(ctx context.Context, in *UpdateFollowsReq, opts ...grpc.CallOption) (*UpdateFollowsResp, error)
		DelFollows(ctx context.Context, in *DelFollowsReq, opts ...grpc.CallOption) (*DelFollowsResp, error)
		GetFollowsById(ctx context.Context, in *GetFollowsByIdReq, opts ...grpc.CallOption) (*GetFollowsByIdResp, error)
		SearchFollows(ctx context.Context, in *SearchFollowsReq, opts ...grpc.CallOption) (*SearchFollowsResp, error)
		// -----------------------鐢ㄦ埛鍩烘湰淇℃伅-----------------------
		AddUserinfo(ctx context.Context, in *AddUserinfoReq, opts ...grpc.CallOption) (*AddUserinfoResp, error)
		UpdateUserinfo(ctx context.Context, in *UpdateUserinfoReq, opts ...grpc.CallOption) (*UpdateUserinfoResp, error)
		DelUserinfo(ctx context.Context, in *DelUserinfoReq, opts ...grpc.CallOption) (*DelUserinfoResp, error)
		GetUserinfoById(ctx context.Context, in *GetUserinfoByIdReq, opts ...grpc.CallOption) (*GetUserinfoByIdResp, error)
		SearchUserinfo(ctx context.Context, in *SearchUserinfoReq, opts ...grpc.CallOption) (*SearchUserinfoResp, error)
	}

	defaultUser struct {
		cli zrpc.Client
	}
)

func NewUser(cli zrpc.Client) User {
	return &defaultUser{
		cli: cli,
	}
}

// -----------------------鐢ㄦ埛鍩烘湰淇℃伅-----------------------
func (m *defaultUser) AddFollows(ctx context.Context, in *AddFollowsReq, opts ...grpc.CallOption) (*AddFollowsResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.AddFollows(ctx, in, opts...)
}

func (m *defaultUser) UpdateFollows(ctx context.Context, in *UpdateFollowsReq, opts ...grpc.CallOption) (*UpdateFollowsResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.UpdateFollows(ctx, in, opts...)
}

func (m *defaultUser) DelFollows(ctx context.Context, in *DelFollowsReq, opts ...grpc.CallOption) (*DelFollowsResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.DelFollows(ctx, in, opts...)
}

func (m *defaultUser) GetFollowsById(ctx context.Context, in *GetFollowsByIdReq, opts ...grpc.CallOption) (*GetFollowsByIdResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.GetFollowsById(ctx, in, opts...)
}

func (m *defaultUser) SearchFollows(ctx context.Context, in *SearchFollowsReq, opts ...grpc.CallOption) (*SearchFollowsResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.SearchFollows(ctx, in, opts...)
}

// -----------------------鐢ㄦ埛鍩烘湰淇℃伅-----------------------
func (m *defaultUser) AddUserinfo(ctx context.Context, in *AddUserinfoReq, opts ...grpc.CallOption) (*AddUserinfoResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.AddUserinfo(ctx, in, opts...)
}

func (m *defaultUser) UpdateUserinfo(ctx context.Context, in *UpdateUserinfoReq, opts ...grpc.CallOption) (*UpdateUserinfoResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.UpdateUserinfo(ctx, in, opts...)
}

func (m *defaultUser) DelUserinfo(ctx context.Context, in *DelUserinfoReq, opts ...grpc.CallOption) (*DelUserinfoResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.DelUserinfo(ctx, in, opts...)
}

func (m *defaultUser) GetUserinfoById(ctx context.Context, in *GetUserinfoByIdReq, opts ...grpc.CallOption) (*GetUserinfoByIdResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.GetUserinfoById(ctx, in, opts...)
}

func (m *defaultUser) SearchUserinfo(ctx context.Context, in *SearchUserinfoReq, opts ...grpc.CallOption) (*SearchUserinfoResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.SearchUserinfo(ctx, in, opts...)
}
