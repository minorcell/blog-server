package services

import (
	"demos/internal/models"
	"demos/pkg/utils"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

type RegisterRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type RegisterResponse struct {
	ID       uint   `json:"id"`
	UserName string `json:"userName"`
	Email    string `json:"email"`
}

func (r *RegisterRequest) Validate() error {
	if r.UserName == "" || r.Password == "" || r.Email == "" {
		return errors.New("username, password and email cannot be empty")
	}
	return nil
}

func (s *UserService) RegisterUser(req *RegisterRequest) (*RegisterResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	var existingUser models.User
	if err := s.db.Where("username = ?", req.UserName).First(&existingUser).Error; err == nil {
		return nil, errors.New("username already exists")
	}

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
		Role:     "1",
		Sex:      "0",
		Age:      0,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, errors.New("failed to create user")
	}

	return &RegisterResponse{
		ID:       user.ID,
		UserName: user.Username,
		Email:    user.Email,
	}, nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	UserName string `json:"userName"`
	Password string `json:"password"`
}

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

type LoginResponse struct {
	Token string             `json:"token"`
	User  *LoginUserResponse `json:"user"`
}

func (r *LoginRequest) Validate() error {
	if r.Email == "" && r.UserName == "" {
		return errors.New("email or username cannot be empty")
	}
	if r.Password == "" {
		return errors.New("password cannot be empty")
	}
	return nil
}

func (s *UserService) LoginUser(req *LoginRequest) (*LoginResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	var user models.User
	if err := s.db.Where("email = ?", req.Email).Or("username = ?", req.UserName).First(&user).Error; err != nil {
		return nil, errors.New("invalid email or username")
	}

	if !utils.ValidatePassword(req.Password, user.Password) {
		return nil, errors.New("invalid password")
	}

	token, err := utils.GenerateToken(int(user.ID), user.Role)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &LoginResponse{
		Token: token,
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

type GetUserInfoResponse struct {
	ID        uint      `json:"id"`
	UserName  string    `json:"userName"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Sex       string    `json:"sex"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type GetUserInfoRequest struct {
	ID uint `json:"id"`
}

func (s *UserService) GetUserInfo(req *gin.Context) (*GetUserInfoResponse, error) {
	quaryID := req.Query("id")
	if quaryID == "" {
		return nil, errors.New("id cannot be empty")
	}
	var user models.User
	if err := s.db.First(&user, quaryID).Error; err != nil {
		return nil, err
	}

	return &GetUserInfoResponse{
		ID:        user.ID,
		UserName:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Sex:       user.Sex,
		Age:       user.Age,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

type UpdateResponse struct {
	ID        uint      `json:"id"`
	UserName  string    `json:"userName"`
	Address   string    `json:"address"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Sex       string    `json:"sex"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateRequest struct {
	UserName string `json:"userName"`
	Email    string `json:"email"`
	Sex      string `json:"sex"`
	Age      int    `json:"age"`
	Address  string `json:"address"`
}

func (r *UpdateRequest) Validate() error {
	if r.UserName == "" && r.Email == "" && r.Sex == "" && r.Age == 0 && r.Address == "" {
		return errors.New("at least one field must be provided")
	}
	return nil
}

func (s *UserService) UpdateUser(req *gin.Context, user_id int) (*UpdateResponse, error) {
	var updateReq UpdateRequest
	if err := req.ShouldBindJSON(&updateReq); err!= nil {
		return nil, errors.New("invalid request body")
	}
	if err := updateReq.Validate(); err!= nil {
		return nil, err
	}

	var user models.User
	if err := s.db.First(&user, user_id).Error; err!= nil {
		return nil, errors.New("user not found")
	}

	if updateReq.UserName != "" && updateReq.UserName != user.Username {
		var existingUser models.User
		if err := s.db.Where("username = ? AND id != ?", updateReq.UserName, user_id).First(&existingUser).Error; err == nil {
			return nil, errors.New("username already exists")
		}
	}

	if updateReq.Email != "" && updateReq.Email != user.Email {
		var existingUser models.User
		if err := s.db.Where("email = ? AND id != ?", updateReq.Email, user_id).First(&existingUser).Error; err == nil {
			return nil, errors.New("email already exists")
		}
	}

	updates := make(map[string]interface{})

	if updateReq.UserName != "" {
		updates["username"] = updateReq.UserName
	}
	if updateReq.Email != "" {
		updates["email"] = updateReq.Email
	}
	if updateReq.Sex != "" {
		updates["sex"] = updateReq.Sex
	}
	if updateReq.Age != 0 {
		updates["age"] = updateReq.Age
	}
	if updateReq.Address != "" {
		updates["address"] = updateReq.Address
	}

	if err := s.db.Model(&user).Updates(updates).Error; err != nil {
		return nil, errors.New("failed to update user information")
	}
	return &UpdateResponse{
		ID:        user.ID,
		UserName:  user.Username,
		Email:     user.Email,
		Address:   user.Address,
		Role:      user.Role,
		Sex:       user.Sex,
		Age:       user.Age,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
