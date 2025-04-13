package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 定义JWT相关错误
var (
	ErrInvalidToken = errors.New("invalid token")     // 无效的token错误
	ErrExpiredToken = errors.New("token has expired") // token已过期错误
)

// JwtClaims 定义JWT的claims结构
// 包含用户ID、角色和标准JWT声明
type JwtClaims struct {
	UserId               int `json:"user_id"` // 用户ID，用于标识用户
	Role                 int `json:"role"`    // 用户角色，用于权限控制
	jwt.RegisteredClaims     // 嵌入标准JWT声明(过期时间、签发时间等)
}

// GenerateToken 生成JWT token
// 参数:
//   - userId: 用户ID
//   - role: 用户角色
//
// 返回值:
//   - string: 生成的token字符串
//   - error: 生成过程中遇到的错误
func GenerateToken(userId int, role int) (string, error) {
	// 从环境变量中获取JWT密钥
	SecretKey := os.Getenv("JWT_SECRET")
	if len(SecretKey) == 0 {
		return "", errors.New("JWT_SECRET not set") // 密钥未设置错误
	}

	// 获取token过期时间配置，默认为24小时
	expireStrTime := os.Getenv("JWT_EXPIRES")
	if len(expireStrTime) == 0 {
		expireStrTime = "24h" // 默认过期时间
	}

	// 解析过期时间字符串为Duration类型
	expireTime, err := time.ParseDuration(expireStrTime)
	if err != nil {
		return "", err // 时间解析错误
	}

	// 创建JWT claims，包含用户信息和标准声明
	claims := &JwtClaims{
		UserId: userId,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireTime)), // 设置过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                 // 设置签发时间
		},
	}

	// 使用HS256算法创建token对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名并生成最终的token字符串
	return token.SignedString([]byte(SecretKey))
}

// ParseToken 解析JWT token
// 参数:
//   - tokenString: 需要解析的token字符串
//
// 返回值:
//   - *JwtClaims: 解析成功返回claims对象指针
//   - error: 解析过程中遇到的错误(密钥未设置/token过期/token无效等)
func ParseToken(tokenString string) (*JwtClaims, error) {
	// 从环境变量中获取JWT密钥
	SecretKey := os.Getenv("JWT_SECRET")
	if len(SecretKey) == 0 {
		return nil, errors.New("JWT_SECRET not set") // 密钥未设置错误
	}

	// 解析token字符串并验证签名
	// 使用JwtClaims结构体来解析claims
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法是否为HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(SecretKey), nil // 返回密钥用于签名验证
	})

	// 处理解析错误
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken // token过期错误
		}
		return nil, ErrInvalidToken // 其他无效token错误
	}

	// 类型断言验证claims类型和token有效性
	claims, ok := token.Claims.(*JwtClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken // claims类型不匹配或token无效
	}

	return claims, nil // 返回解析成功的claims
}
