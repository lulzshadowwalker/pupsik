package types

import (
	"time"

	"github.com/google/uuid"
)

type ImageStatus string

const (
	ImageStatusPending  ImageStatus = "PENDING"
	ImageStatusError    ImageStatus = "ERROR"
	ImageStatusFinished ImageStatus = "FINISHED"
)

type Image struct {
	ID        int
	UserID    uuid.UUID
	BatchID   uuid.UUID
	Prompt    string
	Status    ImageStatus
	URL       string
	CreatedAt time.Time
	DeletedAt time.Time
}
