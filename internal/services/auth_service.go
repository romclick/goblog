package services

import (
	"fmt"
	"goblog/internal/models"
	"goblog/utils"

	"gorm.io/gorm"
)

type AuthService struct {
	db        *gorm.DB
	jwtSecret []byte
	jwtExpire int64
}

func NewAuthService(db *gorm.DB, jwtSecret []byte) *AuthService {
	return &AuthService{
		db:        db,
		jwtSecret: jwtSecret,
		jwtExpire: 3600,
	}
}

func (s *AuthService) Register(username, email, password string) (*models.User, error) {
	var count int64
	if err := s.db.Model(&models.User{}).
		Where("username = ?", username).
		Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, fmt.Errorf("user %s 已存在", username)
	}

	user := &models.User{
		Username: username,
		Email:    email,
		Password: password,
	}
	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Login(loginID, password string) (string, *models.User, error) {
	var user models.User
	if err := s.db.Where("username = ? or email = ?", loginID, loginID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil, fmt.Errorf("用户名/邮箱不存在")
		}
		return "", nil, err
	}
	if password != user.Password {
		return "", nil, fmt.Errorf("密码错误")
	}
	token, err := utils.GenerateToken(user.ID, 24*3600)
	if err != nil {
		return "", nil, fmt.Errorf("生成jwt令牌失败：%w", err)
	}

	return token, &user, nil
}
