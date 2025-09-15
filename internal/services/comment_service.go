package services

import (
	"fmt"
	"goblog/internal/models"

	"gorm.io/gorm"
)

type CommentService struct {
	db *gorm.DB
}

func NewCommentService(db *gorm.DB) *CommentService {
	return &CommentService{db: db}
}

func (c *CommentService) CreateComment(userID, postID uint, content string) (*models.Comment, error) {
	if content == "" {
		return nil, fmt.Errorf("评论不能为空")
	}

	comment := &models.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: content,
	}
	if err := c.db.Create(comment).Error; err != nil {
		return nil, fmt.Errorf("新建评论失败")
	}
	return comment, nil
}

func (c *CommentService) GetComments(postID uint) ([]models.Comment, error) {
	var comments []models.Comment
	if err := c.db.Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		return nil, fmt.Errorf("获取评论失败：%v", err)
	}
	return comments, nil
}

func (c *CommentService) DeleteComment(postID, userID, commentID uint) error {
	var comment models.Comment
	if err := c.db.Where("post_id = ? AND comment_id = ?", postID, commentID).Find(&comment).Error; err != nil {
		return fmt.Errorf("查询评论失败： %v", err)
	}

	if comment.PostID != postID && comment.UserID != userID {
		return fmt.Errorf("仅作者可删除评论")
	}

	if err := c.db.Delete(&comment).Error; err != nil {
		return fmt.Errorf("删除评论失败")
	}
	return nil
}
