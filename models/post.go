package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"size:30;not null"`
	Content   string
	UserID    uint      `gorm:"index"`
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,max=30"`
	Content string `json:"content" binding:"required"`
}

type UpdatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type PostResponse struct {
	PostID uint   `json:"post_id"`
	Title  string `json:"title"`
	// todo:返回文章摘要
	UserID    uint      `json:"user_id"` // todo:理应返回用户可读信息
	Content   string    `json:"content,omitempty"`
	Option    string    `json:"option,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
