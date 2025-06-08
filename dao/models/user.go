package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name         string `gorm:"uniqueIndex"`
	PasswordHash string
}

type Session struct {
	ID     uint64 `gorm:"primaryKey;autoIncrement"`
	UserID uint64 `gorm:"index"`
	UUID   string `gorm:"type:char(36);uniqueIndex"`
	Hint   string
}

type Message struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	SessionID uint64 `gorm:"index"` // 可选索引，便于查询

	Role    string
	Content string
}
