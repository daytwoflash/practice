// todo 公共逻辑提取
package handlers

import (
	"net/http"
	"project/models"
	"project/services"
	"project/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService *services.PostService
	jwtSecret   []byte
}

func NewPostHandler(postService *services.PostService, jwtSecret []byte) *PostHandler {
	return &PostHandler{
		postService: postService,
		jwtSecret:   jwtSecret,
	}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	// === 用户 - 创建 - 博客 ===
	// 解析当前用户
	userID, ok := getAuthenticatedUserID(c)
	if !ok {
		return
	}

	// 解析请求
	var req models.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, parseValidationErrors(err))
		return
	}

	// 创建博客
	post, AppErr := h.postService.CreatePost(userID, req)
	if AppErr != nil {
		utils.HandleError(c, AppErr)
		return
	}

	// 成功返回
	utils.Success(c, models.PostResponse{
		PostID:    post.ID,
		UserID:    post.UserID,
		Title:     post.Title,
		Option:    "create",
		CreatedAt: post.CreatedAt,
	})
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	// === 用户 - 更新 - 博客（自己的）===
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
	var req models.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, parseValidationErrors(err))
		return
	}

	// 判断要更新的博客是否属于当前用户
	post, AppErr := h.postService.GetPostDetail(uint(postID))
	if AppErr != nil {
		utils.HandleError(c, AppErr)
		return
	}
	if post.UserID != userID {
		utils.Error(c, 403, "You can only update your own posts")
		return
	}

	// 更新博客
	post, AppErr = h.postService.UpdatePost(uint(postID), req)
	if AppErr != nil {
		utils.HandleError(c, AppErr)
		return
	}

	// 成功返回
	utils.Success(c, models.PostResponse{
		PostID:    post.ID,
		UserID:    post.UserID,
		Title:     post.Title,
		Option:    "update",
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	})

}

func (h *PostHandler) DeletePost(c *gin.Context) {
	// === 用户 - 删除 - 博客（自己的）===
	// 解析postID
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid Post ID")
		return
	}

	// 解析当前用户
	userID, ok := getAuthenticatedUserID(c)
	if !ok {
		return
	}

	// 判断要删除的博客是否属于当前用户
	post, AppErr := h.postService.GetPostDetail(uint(id))
	if AppErr != nil {
		utils.HandleError(c, AppErr)
		return
	}
	if post.UserID != userID {
		utils.Error(c, 403, "You can only delete your own posts")
		return
	}

	// 删除博客
	if appErr := h.postService.DeletePost(uint(id)); appErr != nil {
		utils.HandleError(c, appErr)
		return
	}

	// 成功返回
	utils.Success(c, models.PostResponse{
		PostID: post.ID,
		UserID: post.UserID,
		Title:  post.Title,
		Option: "delete",
	})

}

func (h *PostHandler) GetPostList(c *gin.Context) {
	// 读取博客列表
	posts, appErr := h.postService.GetPostList()
	if appErr != nil {
		utils.HandleError(c, appErr)
		return
	}

	postResps := make([]models.PostResponse, len(posts))
	for i, post := range posts {
		postResps[i] = models.PostResponse{
			PostID: post.ID,
			UserID: post.UserID,
			Title:  post.Title,
		}
	}

	utils.Success(c, postResps)
}

func (h *PostHandler) GetPostProfile(c *gin.Context) {
	// 解析博客ID
	idStr := c.Param("id")

	postID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid Post ID")
		return
	}

	// 读取指定博客
	post, appErr := h.postService.GetPostDetail(uint(postID))
	if appErr != nil {
		utils.HandleError(c, appErr)
		return
	}

	utils.Success(c, models.PostResponse{
		PostID:  post.ID,
		UserID:  post.UserID,
		Title:   post.Title,
		Content: post.Content,
	})
}
