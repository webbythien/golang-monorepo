package repositories

import (
	"context"

	"github.com/monorepo/app/chat/internal/models"
	"gorm.io/gorm"
)

type Meeting interface {
	CreateMeeting(ctx context.Context, meeting *models.Meeting) error
}

type MeetingStore struct {
	db *gorm.DB
}

var _ Meeting = &MeetingStore{}

func NewMeetingStore(db *gorm.DB) *MeetingStore {
	return &MeetingStore{db: db}
}

func (s *MeetingStore) CreateMeeting(ctx context.Context, meeting *models.Meeting) error {
	return s.db.Create(meeting).Error
}
