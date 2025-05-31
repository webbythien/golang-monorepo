package repositories

import (
	"context"

	"github.com/monorepo/app/chat/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Meeting interface {
	CreateMeeting(ctx context.Context, meeting *models.Meeting) error
	GetMeetingByID(ctx context.Context, meetingID string) (*models.Meeting, error)
	CreateOrUpdateParticipant(ctx context.Context, participant *models.Participant) error
}

type MeetingStore struct {
	db *gorm.DB
}

var _ Meeting = &MeetingStore{}

func NewMeetingStore(db *gorm.DB) *MeetingStore {
	return &MeetingStore{db: db}
}

func (s *MeetingStore) CreateMeeting(ctx context.Context, meeting *models.Meeting) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(meeting).Error; err != nil {
			return err
		}

		participant := &models.Participant{
			MeetingID: meeting.MeetingID,
			UserID:    meeting.HostUserID,
			Role:      models.RoleParticipantHost,
		}

		if err := tx.Create(participant).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *MeetingStore) GetMeetingByID(ctx context.Context, meetingID string) (*models.Meeting, error) {
	var meeting models.Meeting
	if err := s.db.WithContext(ctx).Where("meeting_id = ?", meetingID).First(&meeting).Error; err != nil {
		return nil, err
	}
	return &meeting, nil
}

func (s *MeetingStore) CreateOrUpdateParticipant(ctx context.Context, participant *models.Participant) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "meeting_id"}, {Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"joined_at", "left_at"}),
		}).Create(participant).Error
	})
}
