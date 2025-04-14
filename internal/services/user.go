package services

import (
	"errors"

	"gorm.io/gorm"

	"demos/internal/models"
	"demos/pkg/utils"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// RegisterRequest 注册请求的数据结构
type RegisterRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// RegisterResponse 注册响应的数据结构
type RegisterResponse struct {
	ID       uint   `json:"id"`
	UserName string `json:"userName"`
	Email    string `json:"email"`
}

// Validate 验证注册请求数据
func (r *RegisterRequest) Validate() error {
	if r.UserName == "" || r.Password == "" || r.Email == "" {
		return errors.New("username, password and email cannot be empty")
	}
	return nil
}

// RegisterUser 注册用户
func (s *UserService) RegisterUser(req *RegisterRequest) (*RegisterResponse, error) {
	// 验证请求数据
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := s.db.Where("username = ?", req.UserName).First(&existingUser).Error; err == nil {
		return nil, errors.New("username already exists")
	}

	// 检查邮箱是否已存在
	if err := s.db.Where("email =?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)

	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &models.User{
		Username: req.UserName,
		Password: hashedPassword,
		Email:    req.Email,
		Role:     "1", // 默认角色为普通用户
		Sex:      "0", // 默认性别为未知
		Age:      0,   // 默认年龄为0
	}

	// 保存用户到数据库
	if err := s.db.Create(user).Error; err != nil {
		return nil, errors.New("failed to create user")
	}

	// 返回精简后的用户信息：删除密码字段
	return &RegisterResponse{
		ID:       user.ID,
		UserName: user.Username,
		Email:    user.Email,
	}, nil
}
