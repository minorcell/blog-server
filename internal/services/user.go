package services

import (
	"demos/internal/models"
	"demos/pkg/utils"
	"errors"
	"time"

	"gorm.io/gorm"
)

// UserService 用户服务结构体，包含数据库连接
type UserService struct {
	db *gorm.DB // 使用GORM进行数据库操作
}

// NewUserService 创建UserService实例的构造函数
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// RegisterRequest 注册请求数据结构
type RegisterRequest struct {
	UserName string `json:"userName"` // 用户名
	Password string `json:"password"` // 密码(明文)
	Email    string `json:"email"`    // 邮箱
}

// RegisterResponse 注册响应数据结构
type RegisterResponse struct {
	ID       uint   `json:"id"`       // 用户ID
	UserName string `json:"userName"` // 用户名
	Email    string `json:"email"`    // 邮箱
}

// Validate 验证注册请求数据有效性
func (r *RegisterRequest) Validate() error {
	// 检查必填字段是否为空
	if r.UserName == "" || r.Password == "" || r.Email == "" {
		return errors.New("username, password and email cannot be empty")
	}
	return nil
}

// RegisterUser 用户注册业务逻辑
func (s *UserService) RegisterUser(req *RegisterRequest) (*RegisterResponse, error) {
	// 1. 验证请求数据
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// 2. 检查用户名是否已存在
	var existingUser models.User
	if err := s.db.Where("username = ?", req.UserName).First(&existingUser).Error; err == nil {
		return nil, errors.New("username already exists")
	}

	// 3. 检查邮箱是否已存在
	if err := s.db.Where("email =?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email already exists")
	}

	// 4. 密码加密处理
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// 5. 构建用户模型
	user := &models.User{
		Username: req.UserName,   // 用户名
		Password: hashedPassword, // 加密后的密码
		Email:    req.Email,      // 邮箱
		Role:     "1",            // 默认角色: 普通用户
		Sex:      "0",            // 默认性别: 未知
		Age:      0,              // 默认年龄: 0
	}

	// 6. 保存用户到数据库
	if err := s.db.Create(user).Error; err != nil {
		return nil, errors.New("failed to create user")
	}

	// 7. 返回响应数据(过滤敏感信息)
	return &RegisterResponse{
		ID:       user.ID,       // 用户ID
		UserName: user.Username, // 用户名
		Email:    user.Email,    // 邮箱
	}, nil
}

// LoginRequest 登录请求数据结构
type LoginRequest struct {
	Email    string `json:"email"`    // 邮箱(与用户名二选一)
	UserName string `json:"userName"` // 用户名(与邮箱二选一)
	Password string `json:"password"` // 密码(明文)
}

// LoginUserResponse 登录响应数据结构(过滤敏感信息)
type LoginUserResponse struct {
	ID        uint      `json:"id"`
	UserName  string    `json:"userName"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Sex       string    `json:"sex"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// LoginResponse 登录响应数据结构
type LoginResponse struct {
	Token string             `json:"token"` // JWT令牌
	User  *LoginUserResponse `json:"user"`  // 过滤后的用户信息
}

// Validate 验证登录请求数据有效性
func (r *LoginRequest) Validate() error {
	// 检查至少提供邮箱或用户名之一
	if r.Email == "" && r.UserName == "" {
		return errors.New("email or username cannot be empty")
	}
	// 检查密码是否为空
	if r.Password == "" {
		return errors.New("password cannot be empty")
	}
	return nil
}

// LoginUser 用户登录业务逻辑
func (s *UserService) LoginUser(req *LoginRequest) (*LoginResponse, error) {
	// 1. 验证请求数据
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// 2. 查询用户(通过邮箱或用户名)
	var user models.User
	if err := s.db.Where("email = ?", req.Email).Or("username = ?", req.UserName).First(&user).Error; err != nil {
		return nil, errors.New("invalid email or username")
	}

	// 3. 验证密码
	if !utils.ValidatePassword(req.Password, user.Password) {
		return nil, errors.New("invalid password")
	}

	// 4. 生成JWT令牌
	token, err := utils.GenerateToken(int(user.ID), user.Role)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// 5. 返回响应数据
	return &LoginResponse{
		Token: token, // JWT令牌
		User: &LoginUserResponse{
			ID:        user.ID,
			UserName:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			Sex:       user.Sex,
			Age:       user.Age,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}
