package services

import (
	"errors"
	"project/models"
	"project/utils"

	"gorm.io/gorm"
)

type PostService struct {
	db *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{db: db}
}

func (s *PostService) CreatePost(userID uint, req models.CreatePostRequest) (*models.Post, error) {
	// 默认可创建相同标题和内容的博客
	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}

	if err := s.db.Create(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

func (s *PostService) UpdatePost(postID uint, req models.UpdatePostRequest) (*models.Post, error) {
	var post models.Post

	result := s.db.Model(&post).
		Where("id = ?", postID).
		Updates(models.Post{Title: req.Title, Content: req.Content})

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, utils.NewAppError(404, "Post not found")
	}

	if err := s.db.First(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewAppError(404, "Post not found")
		}
		return nil, err
	}

	return &post, nil

}

func (s *PostService) DeletePost(postID uint) error {
	result := s.db.Delete(&models.Post{}, postID)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return utils.NewAppError(404, "Post not found")
	}

	return nil
}

func (s *PostService) GetPostDetail(postID uint) (*models.Post, error) {
	var post models.Post
	if err := s.db.First(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewAppError(404, "Post not found")
		}
		return nil, err
	}

	return &post, nil

}

func (s *PostService) GetPostList() ([]models.Post, error) {
	var posts []models.Post
	if err := s.db.Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil

}
