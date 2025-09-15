package handlers

import (
	"goblog/internal/services"
	"goblog/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService *services.PostService
}

// 创建文章接口处理器
func NewPostHandler(postService *services.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

// 创建文章
func (h *PostHandler) CreatePost(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		response.Error(c, http.StatusUnauthorized, 401, "请先登录")
		return
	}
	uid, _ := userID.(uint)

	type CreatePostRequest struct {
		Title   string `json:"title" binding:"required,max=255"`
		Content string `json:"content" binding:"required"`
	}
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, 422, err.Error())
		return
	}

	post, err := h.postService.CreatePost(uid, req.Title, req.Content)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.Success(c, post)
}

// 获取文章详情接口
func (h *PostHandler) GetPostByID(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusUnprocessableEntity, 422, err.Error())
	}

	post, err := h.postService.GetPost(uint(postID))
	if err != nil {
		if err.Error() == "sql: no rows in result set 文章不存在" {
			response.Error(c, http.StatusUnprocessableEntity, 422, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.Success(c, post)
}
