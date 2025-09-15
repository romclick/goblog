package handlers

import (
	"goblog/internal/services"
	"goblog/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentService *services.CommentService
}

func NewCommentHandler(commentService *services.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

// 创建评论接口
func (h *CommentHandler) CreateComment(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		response.Error(c, http.StatusUnauthorized, 401, "请先登录")
	}
	uid, _ := userID.(uint)

	if !ok {
		response.Error(c, http.StatusUnauthorized, 400, "用户ID格式错误")
		return
	}

	postIDStr := c.PostForm("post_id")
	postID, _ := strconv.Atoi(postIDStr)

	pid := uint(postID)

	type CreateCommentRequest struct {
		Content string `json:"content"`
	}

	var req CreateCommentRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, 422, err.Error())
		return
	}

	comment, err := h.commentService.CreateComment(uid, pid, req.Content)
	if err != nil {
		response.Error(c, http.StatusUnprocessableEntity, 500, "创建评论失败")
		return
	}
	response.Success(c, comment)
}

func (h *CommentHandler) GetComment(c *gin.Context) {
	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusUnprocessableEntity, 422, err.Error())
	}

	comment, err := h.commentService.GetComments(uint(commentID))
	if err != nil {
		if err.Error() == "sql: no rows in result set 暂无评论" {
			response.Error(c, http.StatusUnprocessableEntity, 422, err.Error())
		}
		response.Error(c, http.StatusUnprocessableEntity, 422, err.Error())
		return
	}
	response.Success(c, comment)
}
