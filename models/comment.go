package models

import "time"

type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index; not null"`
	PostID    uint      ` json:"post_id" gorm:"index; not null"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// todo:去除
type CommentResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	PostID    uint      `json:"post_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentListResponse struct {
	PostID   uint              `json:"post_id"`
	Comments []CommentResponse `json:"comments"`
}
