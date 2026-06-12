package services

import (
	"errors"
	"project/models"
	"project/utils"

	"gorm.io/gorm"
)

type CommentService struct {
	db *gorm.DB
}

func NewCommentService(db *gorm.DB) *CommentService {
	return &CommentService{db: db}
}

func (s *CommentService) CreateComment(postID, userID uint, req models.CreateCommentRequest) (*models.Comment, error) {
	if err := s.db.First(&models.Post{}, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewAppError(404, "Post not found")
		}
		return nil, err
	}

	// 允许 用户 - 博客 重复评论
	comment := models.Comment{
		UserID:  userID,
		PostID:  postID,
		Content: req.Content,
	}

	if err := s.db.Create(&comment).Error; err != nil {
		return nil, err
	}

	return &comment, nil

}

func (s *CommentService) GetCommentList(postID uint) ([]models.Comment, error) {
	if err := s.db.First(&models.Post{}, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewAppError(404, "Post not found")
		}
		return nil, err
	}

	var comments []models.Comment
	if err := s.db.Find(&comments, "post_id = ?", postID).Error; err != nil {
		// 空时 - 正常
		return nil, err
	}

	return comments, nil
}
