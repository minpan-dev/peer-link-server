package model

import "time"

type User struct {
	ID        uint      `json:"id"         gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name"       gorm:"type:varchar(100);not null"`
	Email     string    `json:"email"      gorm:"type:varchar(255);uniqueIndex;not null"`
	Password  string    `json:"-"          gorm:"type:varchar(255);not null"` // 永不序列化
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string { return "users" }
