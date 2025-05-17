package id

import (
	"math/rand"
	"strings"
	"time"
)

type MeetingIDGenerator struct{}

func NewMeetingIDGenerator() *MeetingIDGenerator {
	return &MeetingIDGenerator{}
}

func (gen *MeetingIDGenerator) GenerateID() string {
	const charset = "abcdefghijklmnopqrstuvwxyz"

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	generateSegment := func(length int) string {
		b := make([]byte, length)
		for i := range b {
			b[i] = charset[r.Intn(len(charset))]
		}
		return string(b)
	}

	return strings.Join([]string{generateSegment(3), generateSegment(4), generateSegment(3)}, "-")
}
