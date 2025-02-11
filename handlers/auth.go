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

func CreateNewUser(c *gin.Context) {
	// 检查当前用户是否是管理员
	if c.GetString("role") != "admin" {
		c.JSON(403, gin.H{"error": "无权限"})
		return
	}

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名是否已存在
	collection := config.GetCollection("users")
	count, _ := collection.CountDocuments(context.Background(), bson.M{"username": req.Username})
	if count > 0 {
		c.JSON(400, gin.H{"error": "用户名已存在"})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "密码加密失败"})
		return
	}

	user := models.User{
		Username:  req.Username,
		Password:  string(hashedPassword),
		Role:      req.Role,
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		c.JSON(500, gin.H{"error": "创建用户失败"})
		return
	}

	c.JSON(200, gin.H{"message": "创建成功"})
}
