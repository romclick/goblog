package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Content string `gorm:"not null"`
	UserID  uint   `gorm:"not null"`
	User    User   `gorm:"foreignkey:UserID"`
	PostID  uint   `gorm:"not null"`
	Post    Post   `gorm:"foreignkey:PostID"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	if c.Content != "" {
		return fmt.Errorf("评论不能为空")
	}

	var postCount int64
	if err := tx.Model(&Post{}).Where("post_id = ?", c.PostID).Count(&postCount).Error; err != nil {
		return fmt.Errorf("校验文章存在失败", err)
	}
	if postCount == 0 {
		return fmt.Errorf("文章不存在")
	}

	var userCount int64
	if err := tx.Model(&User{}).Where("user_id = ?", c.UserID).Count(&userCount).Error; err != nil {
		return fmt.Errorf("校验用户失败", err)
	}
	if userCount == 0 {
		return fmt.Errorf("用户不存在", err)
	}
	return nil
}

func (c *Comment) AfterCreate(tx *gorm.DB) (err error) {
	return tx.Model(&Post{}).
		Where("id = ?", c.PostID).
		Update("comment_status", "有评论").Error
}

func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	var commentCount int64
	if err := tx.Model(&Comment{}).
		Where("post_id = ? and delete_at is null", c.PostID).
		Count(&commentCount).Error; err != nil {
		return fmt.Errorf("统计剩余评论数失败", err)
	}
	status := "无评论"
	if commentCount > 0 {
		status = "有评论"
	}
	return tx.Model(&Post{}).
		Where("id = ?", c.PostID).
		Update("comment_status", status).Error
}
