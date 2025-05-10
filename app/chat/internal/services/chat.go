package services

import (
	"context"

	"connectrpc.com/connect"
	chatv1 "github.com/webbythien/monorepo/api/chat/v1"
	"github.com/webbythien/monorepo/api/chat/v1/chatv1connect"
)

var _ chatv1connect.ChatAPIHandler = &ChatAPI{}

type ChatAPI struct{}

func NewChatAPI() *ChatAPI {
	return &ChatAPI{}
}

func (c *ChatAPI) CreateRoom(ctx context.Context, req *connect.Request[chatv1.CreateRoomRequest]) (*connect.Response[chatv1.CreateRoomResponse], error) {
	return connect.NewResponse(&chatv1.CreateRoomResponse{
		Id: "abc",
	}), nil
}
