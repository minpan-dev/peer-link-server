package model

import "time"

type Room struct {
	ID              uint      `json:"id"            gorm:"primaryKey;autoIncrement"`
	Name            string    `json:"name"          gorm:"type:varchar(100);uniqueIndex;not null"`
	DisplayName     string    `json:"display_name"  gorm:"type:varchar(255)"`
	MaxParticipants int       `json:"max_participants" gorm:"default:20"`
	CreatedByID     uint      `json:"created_by_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (Room) TableName() string { return "rooms" }
