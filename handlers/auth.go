// 认证处理器 - 由 AaronChenH 维护
// 包含登录、注销、权限验证等功能
package handlers

import (
	"context"
	"game-admin/config"
	"game-admin/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest 登录请求结构
// 注意: 字段需要与前端表单保持一致
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

const (
	SecretKey = "your-secret-key" // 实际使用时应该放在配置文件中
)

// LoginHandler 处理用户登录
// @Summary 用户登录
// @Description 处理用户登录请求
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录凭证"
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Router /login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithError(err).Error("登录请求数据绑定失败")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	collection := config.GetCollection("users")
	var user models.User

	logrus.WithFields(logrus.Fields{
		"username": req.Username,
	}).Info("尝试登录")

	err := collection.FindOne(context.Background(), bson.M{"username": req.Username}).Decode(&user)
	if err != nil {
		logrus.WithError(err).WithField("username", req.Username).Error("查找用户失败")
		c.JSON(401, gin.H{"error": "用户名或密码错误"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logrus.WithError(err).WithField("username", req.Username).Error("密码验证失败")
		c.JSON(401, gin.H{"error": "用户名或密码错误"})
		return
	}

	if user.Status != 1 {
		c.JSON(403, gin.H{"error": "账号已被禁用"})
		return
	}

	// 生成JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		c.JSON(500, gin.H{"error": "生成token失败"})
		return
	}

	c.JSON(200, gin.H{
		"token": tokenString,
		"user":  user,
	})
}