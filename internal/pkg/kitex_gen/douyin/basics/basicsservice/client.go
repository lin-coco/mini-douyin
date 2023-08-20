// Code generated by Kitex v0.6.2. DO NOT EDIT.

package basicsservice

import (
	"context"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
	"mini-douyin/internal/pkg/kitex_gen/douyin/basics"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	GetUserInfoById(ctx context.Context, Req *basics.GetUserRequest, callOptions ...callopt.Option) (r *basics.GetUserResponse, err error)
	CreateUser(ctx context.Context, Req *basics.CreateUserRequest, callOptions ...callopt.Option) (r *basics.CreateUserResponse, err error)
	CheckUser(ctx context.Context, Req *basics.CheckUserRequest, callOptions ...callopt.Option) (r *basics.CheckUserResponse, err error)
	GetVideoInfoById(ctx context.Context, Req *basics.GetVideoByIdRequest, callOptions ...callopt.Option) (r *basics.GetVideoByIdResponse, err error)
	GetVideo(ctx context.Context, Req *basics.GetVideoRequest, callOptions ...callopt.Option) (r *basics.GetVideoResponse, err error)
	UploadVideo(ctx context.Context, Req *basics.UploadVideoRequest, callOptions ...callopt.Option) (r *basics.UploadVideoResponse, err error)
	GetVideosByUserId(ctx context.Context, Req *basics.GetVideosByUserIdRequest, callOptions ...callopt.Option) (r *basics.GetVideosByUserIdResponse, err error)
	GetVideoListByIds(ctx context.Context, Req *basics.GetVideoListByIdsRequest, callOptions ...callopt.Option) (r *basics.GetVideoListByIdsResponse, err error)
	GetUserListByIds(ctx context.Context, Req *basics.GetUserListByIdsRequest, callOptions ...callopt.Option) (r *basics.GetUserListByIdsResponse, err error)
	GetVideoCount(ctx context.Context, Req *basics.VideoCountRequest, callOptions ...callopt.Option) (r *basics.VideoCountResponse, err error)
	GetUserVideoIds(ctx context.Context, Req *basics.UserVideoIdsRequest, callOptions ...callopt.Option) (r *basics.UserVideoIdsResponse, err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfo(), options...)
	if err != nil {
		return nil, err
	}
	return &kBasicsServiceClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kBasicsServiceClient struct {
	*kClient
}

func (p *kBasicsServiceClient) GetUserInfoById(ctx context.Context, Req *basics.GetUserRequest, callOptions ...callopt.Option) (r *basics.GetUserResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetUserInfoById(ctx, Req)
}

func (p *kBasicsServiceClient) CreateUser(ctx context.Context, Req *basics.CreateUserRequest, callOptions ...callopt.Option) (r *basics.CreateUserResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.CreateUser(ctx, Req)
}

func (p *kBasicsServiceClient) CheckUser(ctx context.Context, Req *basics.CheckUserRequest, callOptions ...callopt.Option) (r *basics.CheckUserResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.CheckUser(ctx, Req)
}

func (p *kBasicsServiceClient) GetVideoInfoById(ctx context.Context, Req *basics.GetVideoByIdRequest, callOptions ...callopt.Option) (r *basics.GetVideoByIdResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetVideoInfoById(ctx, Req)
}

func (p *kBasicsServiceClient) GetVideo(ctx context.Context, Req *basics.GetVideoRequest, callOptions ...callopt.Option) (r *basics.GetVideoResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetVideo(ctx, Req)
}

func (p *kBasicsServiceClient) UploadVideo(ctx context.Context, Req *basics.UploadVideoRequest, callOptions ...callopt.Option) (r *basics.UploadVideoResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.UploadVideo(ctx, Req)
}

func (p *kBasicsServiceClient) GetVideosByUserId(ctx context.Context, Req *basics.GetVideosByUserIdRequest, callOptions ...callopt.Option) (r *basics.GetVideosByUserIdResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetVideosByUserId(ctx, Req)
}

func (p *kBasicsServiceClient) GetVideoListByIds(ctx context.Context, Req *basics.GetVideoListByIdsRequest, callOptions ...callopt.Option) (r *basics.GetVideoListByIdsResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetVideoListByIds(ctx, Req)
}

func (p *kBasicsServiceClient) GetUserListByIds(ctx context.Context, Req *basics.GetUserListByIdsRequest, callOptions ...callopt.Option) (r *basics.GetUserListByIdsResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetUserListByIds(ctx, Req)
}

func (p *kBasicsServiceClient) GetVideoCount(ctx context.Context, Req *basics.VideoCountRequest, callOptions ...callopt.Option) (r *basics.VideoCountResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetVideoCount(ctx, Req)
}

func (p *kBasicsServiceClient) GetUserVideoIds(ctx context.Context, Req *basics.UserVideoIdsRequest, callOptions ...callopt.Option) (r *basics.UserVideoIdsResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetUserVideoIds(ctx, Req)
}