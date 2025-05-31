package models

import "time"

var _ = registerAutoMigrate(&Meeting{}, &Participant{})

type Meeting struct {
	ID              uint64        `gorm:"autoIncrement" db:"id" json:"id"`
	MeetingID       string        `gorm:"type:text;primaryKey;;not null" db:"meeting_id" json:"meeting_id"`
	HostUserID      string        `gorm:"type:text;not null" db:"host_user_id" json:"host_user_id"`
	Title           string        `gorm:"type:text" db:"title" json:"title"`
	ScheduledStart  time.Time     `gorm:"type:timestamptz" db:"scheduled_start" json:"scheduled_start"`
	ScheduledEnd    time.Time     `gorm:"type:timestamptz" db:"scheduled_end" json:"scheduled_end"`
	IsRecurring     bool          `gorm:"default:false" db:"is_recurring" json:"is_recurring"`
	CreatedAt       time.Time     `gorm:"autoCreateTime" db:"created_at" json:"created_at"`
	UpdatedAt       time.Time     `gorm:"autoUpdateTime" db:"updated_at" json:"updated_at"`
	DurationMinutes int64         `gorm:"type:BIGINT;not null" json:"duration_minutes,omitempty"`
	Status          MeetingStatus `gorm:"type:text;not null;default:in-progress" db:"status" json:"status"`

	// Participants []*Participant `gorm:"-" json:"participants,omitempty"`
}

func (m *Meeting) IsDone() bool {
	return m.Status == MeetingStatusClosed
}

func (Meeting) TableName() string {
	return "meetings"
}

type MeetingStatus string

const (
	MeetingStatusInProgress MeetingStatus = "in-progress"
	MeetingStatusClosed     MeetingStatus = "closed"
)

type Participant struct {
	ID        uint64          `gorm:"primaryKey;autoIncrement" db:"id" json:"id"`
	MeetingID string          `gorm:"type:varchar;not null;index:idx_meeting_user,unique" db:"meeting_id" json:"meeting_id"`
	UserID    string          `gorm:"type:varchar;not null;index:idx_meeting_user,unique" db:"user_id" json:"user_id"`
	Role      RoleParticipant `gorm:"type:text;default:participant" db:"role" json:"role"`
	JoinedAt  time.Time       `gorm:"type:timestamptz;autoCreateTime" db:"joined_at" json:"joined_at"`
	LeftAt    *time.Time      `gorm:"type:timestamptz" db:"left_at" json:"left_at,omitempty"`

	// Meeting *Meeting `gorm:"-" json:"meeting,omitempty"`
}

func (Participant) TableName() string {
	return "participants"
}

type RoleParticipant string

const (
	RoleParticipantParticipant RoleParticipant = "participant"
	RoleParticipantHost        RoleParticipant = "host"
)
