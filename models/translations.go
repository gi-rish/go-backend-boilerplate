package models

import (
	"time"

	"github.com/google/uuid"
)

type Translation struct {
	Id           uuid.UUID  `gorm:"primaryKey;column:id;not null;default:uuid_generate_v4()"`
	Code         string     `gorm:"code;not null;"`
	Key          string     `gorm:"key;not null;"`
	Translations string     `gorm:"translations;not null;"`
	DateCreated  *time.Time `gorm:"date_created;default:current_timestamp"`
	DateUpdated  *time.Time `gorm:"date_updated;default:current_timestamp"`
}
