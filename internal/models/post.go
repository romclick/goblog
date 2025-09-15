package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title          string `gorm:"type:varchar(255)"`
	Content        string `gorm:"type:text"`
	Comment_status bool   `gorm:"not null;default:false"`
	UserID         uint   `gorm:"not null"`
	User           User   `gorm:"foreignkey:UserID"`
}

func (post *Post) BeforeCreate(tx *gorm.DB) error {
	if post.Title == "" {
		return fmt.Errorf("文章标题不能为空")
	}
	if len(post.Title) > 300 {
		return fmt.Errorf("文章标题不能超过100字")
	}
	if post.Content == "" {
		return fmt.Errorf("文章内容不能为空")
	}
	return nil
}

func (p *Post) AfterCreate(tx *gorm.DB) error {
	result := tx.Model(&User{}).
		Where("userid = ?", p.UserID).
		Update("post_count", gorm.Expr("post_count + ?", 1))
	if result.Error != nil {
		return fmt.Errorf("更新文章数量失败", result.Error)
	}
	return nil
}

func (p *Post) AfterDelete(tx *gorm.DB) error {
	if err := tx.Where("post_id = ?", p.ID).Delete(&Post{}).Error; err != nil {
		return fmt.Errorf("删除文章评论失败", err)
	}

	return tx.Model(&User{}).
		Where("userid = ?", p.UserID).
		Update("post_count", gorm.Expr("post_count - ?", 1)).Error
}
