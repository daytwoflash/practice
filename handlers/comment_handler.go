package handlers

import (
	"net/http"
	"project/models"
	"project/services"
	"project/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentService *services.CommentService
	jwtSecret      []byte
}

func NewCommentHandler(postService *services.CommentService, jwtSecret []byte) *CommentHandler {
	return &CommentHandler{
		commentService: postService,
		jwtSecret:      jwtSecret,
	}
}

func (h *CommentHandler) GetCommentList(c *gin.Context) {
	// 博客 - 的 - 评论
	idStr := c.Param("id")
	postID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid Post ID")
		return
	}

	comments, appErr := h.commentService.GetCommentList(uint(postID))
	if appErr != nil {
		utils.HandleError(c, appErr)
		return
	}

	commentResps := make([]models.CommentResponse, len(comments))
	for i, comment := range comments {
		commentResps[i] = models.CommentResponse{
			UserID:  comment.UserID,
			PostID:  comment.PostID,
			Content: comment.Content,
		}
	}

	utils.Success(c, models.CommentListResponse{
		PostID:   uint(postID),
		Comments: commentResps,
	})

}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	// 解析 postID
	idStr := c.Param("id")
	postID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid Post ID")
		return
	}

	// 解析当前用户
	userID, ok := getAuthenticatedUserID(c)
	if !ok {
		return
	}

	// 解析请求
	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, parseValidationErrors(err))
		return
	}

	// 创建评论
	comment, appErr := h.commentService.CreateComment(uint(postID), userID, req)
	if appErr != nil {
		utils.HandleError(c, appErr)
		return
	}

	// 成功返回
	utils.Success(c, models.CommentResponse{
		ID:        comment.ID,
		UserID:    userID,
		PostID:    comment.PostID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
	})

}
