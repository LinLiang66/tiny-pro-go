package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"tiny-admin-api-serve/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

// AuthMiddleware 鉴权中间件结构体
type AuthMiddleware struct {
	secretKey []byte
	issuer    string
}

var Auth AuthMiddleware

// UserClaims 用户自定义声明结构体
type UserClaims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// 初始化全局AuthMiddleware
func init() {
	Auth = AuthMiddleware{
		secretKey: []byte(viper.GetString("jwt.secret")),
		issuer:    viper.GetString("jwt.app_name"),
	}
}

// AuthRequired 鉴权拦截器
func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token format required"})
			c.Abort()
			return
		}

		claims, err := m.parseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			c.Abort()
			return
		}

		// 检查token是否在黑名单中(已注销)
		ctx := context.Background()
		blacklisted := utils.Redis.Exists(ctx, fmt.Sprintf("blacklist:%s", tokenString))
		if blacklisted {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("claims", claims)

		c.Next()
	}
}

// parseToken 解析并验证JWT token
func (m *AuthMiddleware) parseToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	// 验证签发者
	if claims.Issuer != m.issuer {
		return nil, errors.New("invalid issuer")
	}

	return claims, nil
}

// GenerateToken 生成JWT token
func (m *AuthMiddleware) GenerateToken(userID int64, email, role string) (string, error) {
	claims := &UserClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24小时过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        generateUniqueID(), // 可以使用UUID库生成唯一ID
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secretKey)
}

// Logout 用户登出，将token加入黑名单
func (m *AuthMiddleware) Logout(tokenString string) error {
	ctx := context.Background()

	// 解析token获取过期时间
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return m.secretKey, nil
	})

	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("missing expiration time")
	}

	// 计算剩余过期时间
	ttl := time.Until(time.Unix(int64(exp), 0))

	// 将token加入Redis黑名单，设置与token相同的过期时间
	return utils.Redis.SetExpire(ctx, fmt.Sprintf("blacklist:%s", tokenString), ttl)
}

// RefreshToken 刷新token
func (m *AuthMiddleware) RefreshToken(oldTokenString string) (string, error) {
	claims, err := m.parseToken(oldTokenString)
	if err != nil {
		return "", err
	}

	// 生成新token
	newToken, err := m.GenerateToken(claims.UserID, claims.Email, claims.Role)
	if err != nil {
		return "", err
	}

	// 将旧token加入黑名单
	err = m.Logout(oldTokenString)
	if err != nil {
		return "", err
	}
	return newToken, nil
}

// generateUniqueID 生成唯一ID（简化版，实际可使用UUID）
func generateUniqueID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
