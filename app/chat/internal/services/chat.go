package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"connectrpc.com/connect"
	chatv1 "github.com/monorepo/api/chat/v1"
	"github.com/monorepo/api/chat/v1/chatv1connect"
	"github.com/monorepo/app/chat/internal/id"
	"github.com/monorepo/app/chat/internal/models"
	"github.com/monorepo/app/chat/internal/repositories"
	eventChatV1 "github.com/monorepo/event/chat/v1"
	"github.com/monorepo/pkg/emitter"
	"gorm.io/gorm"
)

var _ chatv1connect.ChatAPIHandler = &ChatAPI{} // Implement all methods from chatv1connect.ChatAPIHandler

type ChatAPI struct {
	meetingIDGenerator *id.MeetingIDGenerator
	meetingStore       repositories.Meeting
	emitter            *emitter.Emitter
}

func NewChatAPI(meetingStore repositories.Meeting, emitter *emitter.Emitter) *ChatAPI {
	return &ChatAPI{
		meetingIDGenerator: id.NewMeetingIDGenerator(),
		meetingStore:       meetingStore,
		emitter:            emitter,
	}
}

func (c *ChatAPI) UserCreateMeeting(ctx context.Context, request *connect.Request[chatv1.UserCreateMeetingRequest]) (*connect.Response[chatv1.UserCreateMeetingResponse], error) {
	var (
		req   = request.Msg
		title = req.GetTitle()
	)
	meetingID := c.meetingIDGenerator.GenerateID()

	durationMinutes := int64(60) // FIXME: get from tier user
	scheduledStart := time.Now()
	scheduledEnd := scheduledStart.Add(time.Duration(durationMinutes) * time.Minute)
	meeting := &models.Meeting{
		Title:           title,
		MeetingID:       meetingID,
		DurationMinutes: durationMinutes, // FIXME: get from tier user
		Status:          models.MeetingStatusInProgress,
		HostUserID:      "test-user-id-create", // FIXME: get from context
		ScheduledStart:  scheduledStart,
		ScheduledEnd:    scheduledEnd,
	}

	if err := c.meetingStore.CreateMeeting(ctx, meeting); err != nil {
		return nil, err
	}

	return connect.NewResponse(&chatv1.UserCreateMeetingResponse{
		MeetingId: meetingID,
	}), nil
}

func (c *ChatAPI) UserJoinMeeting(ctx context.Context, request *connect.Request[chatv1.UserJoinMeetingRequest]) (*connect.Response[chatv1.UserJoinMeetingResponse], error) {
	var (
		req       = request.Msg
		meetingID = req.GetMeetingId()
		sdpOffer  = req.GetSdpOffer()
	)

	fmt.Println("sdpOffer", sdpOffer)
	meeting, err := c.meetingStore.GetMeetingByID(ctx, meetingID)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, connect.NewError(connect.CodeNotFound, errors.New("meeting not found")) // TODO: return biz error
	case err != nil:
		return nil, err
	default:
		break
	}

	if meeting.IsDone() {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("meeting is ended")) // TODO: return biz error
	}

	participant := &models.Participant{
		MeetingID: meetingID,
		UserID:    "test-user-id-join", // FIXME: get from context
		Role:      models.RoleParticipantParticipant,
		JoinedAt:  time.Now(),
	}

	if err := c.meetingStore.CreateOrUpdateParticipant(ctx, participant); err != nil {
		return nil, err
	}

	c.emitter.In(meetingID).Emit(eventChatV1.Event_EVENT_SOCKET_JOIN_MEETING.String(), "Go Emitter: Hello, World!")
	fmt.Println("eventChatV1.Event_EVENT_SOCKET_JOIN_MEETING: ", eventChatV1.Event_EVENT_SOCKET_JOIN_MEETING)
	return connect.NewResponse(&chatv1.UserJoinMeetingResponse{
		Message: "ok",
	}), nil
}
