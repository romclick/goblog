package services

import (
	"errors"
	"fmt"
	"goblog/internal/models"

	"gorm.io/gorm"
)

type PostService struct {
	db *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{db: db}
}

func (s *PostService) CreatePost(userID uint, title, content string) (*models.Post, error) {
	if title == "" || content == "" {
		return nil, fmt.Errorf("文章标题/内容不能为空")
	}

	post := &models.Post{
		UserID:  userID,
		Title:   title,
		Content: content,
	}
	if err := s.db.Create(post).Error; err != nil {
		return nil, fmt.Errorf("创建文章失败： %w", err)
	}
	return post, nil
}

func (s *PostService) UpdatePost(userID uint, postID uint, title, content string) (*models.Post, error) {
	var post models.Post
	if err := s.db.Where("id = ?", postID).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("文章不存在")
		}
		return nil, fmt.Errorf("查询文章失败： %w", err)
	}
	if post.UserID != userID {
		return nil, fmt.Errorf("无权更新该文章")
	}

	if err := s.db.Model(&post).Updates(map[string]interface{}{
		"title":   title,
		"content": content,
	}).Error; err != nil {
		return nil, fmt.Errorf("更新文章失败 %w", err)
	}
	return &post, nil
}

func (s *PostService) DeletePost(userID uint, postID uint) error {
	var post models.Post
	if err := s.db.Where("id = ?", postID).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("文章不存在")
		}
		return fmt.Errorf("查询文章失败 %w", err)
	}

	if post.UserID != userID {
		return fmt.Errorf("仅作者可删除")
	}

	if err := s.db.Delete(&post).Error; err != nil {
		return fmt.Errorf("删除文章失败 %w", err)
	}
	return nil
}

func (s *PostService) GetPost(postID uint) (*models.Post, error) {
	var post models.Post
	if err := s.db.Where("id = ?", postID).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("文章不存在")
		}
		return nil, fmt.Errorf("查询失败 %w", err)
	}
	return &post, nil
}
