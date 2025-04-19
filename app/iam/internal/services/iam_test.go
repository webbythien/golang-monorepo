package services

import (
	"context"
	iamv1 "monorepo/api/iam/v1"
	"monorepo/api/iam/v1/iamv1connect"

	"connectrpc.com/connect"
)

var _ iamv1connect.SecurityTokenAPIHandler = &IamTest{}

type IamTest struct {
}

func NewIamTest() *IamTest {
	return &IamTest{}
}

func (i *IamTest) TestApiGenProto(ctx context.Context, c *connect.Request[iamv1.TestApiGenProtoRequest]) (*connect.Response[iamv1.TestApiGenProtoResponse], error) {

	return nil, nil
}
