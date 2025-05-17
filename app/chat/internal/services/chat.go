package services

import (
	"context"

	"connectrpc.com/connect"
	chatv1 "github.com/monorepo/api/chat/v1"
	"github.com/monorepo/api/chat/v1/chatv1connect"
	"github.com/monorepo/app/chat/internal/id"
	"github.com/monorepo/app/chat/internal/models"
	"github.com/monorepo/app/chat/internal/repositories"
)

var _ chatv1connect.ChatAPIHandler = &ChatAPI{}

type ChatAPI struct {
	meetingIDGenerator *id.MeetingIDGenerator
	meetingStore       repositories.Meeting
}

func NewChatAPI(meetingStore repositories.Meeting) *ChatAPI {
	return &ChatAPI{
		meetingIDGenerator: id.NewMeetingIDGenerator(),
		meetingStore:       meetingStore,
	}
}

func (c *ChatAPI) CreateRoom(ctx context.Context, request *connect.Request[chatv1.CreateRoomRequest]) (*connect.Response[chatv1.CreateRoomResponse], error) {
	var (
		req   = request.Msg
		title = req.GetTitle()
	)
	meetingID := c.meetingIDGenerator.GenerateID()

	meeting := &models.Meeting{
		Title:           title,
		MeetingCode:     meetingID,
		DurationMinutes: 60, // FIXME: get from tier user
		Status:          models.MeetingStatusInProgress,
		HostUserID:      "test-user-id", // FIXME: get from context
	}

	if err := c.meetingStore.CreateMeeting(ctx, meeting); err != nil {
		return nil, err
	}

	return connect.NewResponse(&chatv1.CreateRoomResponse{
		MeetingCode: meetingID,
	}), nil
}
