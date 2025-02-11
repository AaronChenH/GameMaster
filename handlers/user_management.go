package handlers

import (
	"context"
	"game-admin/config"
	"game-admin/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// GetUsers 获取所有用户列表
func GetUsers(c *gin.Context) {
	collection := config.GetCollection("users")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表失败"})
		return
	}
	defer cursor.Close(context.Background())

	var users []models.User
	if err = cursor.All(context.Background(), &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解析用户数据失败"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// UpdateUserStatus 更新用户状态
func UpdateUserStatus(c *gin.Context) {
	userId := c.Param("id")
	var req struct {
		Status int `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	collection := config.GetCollection("users")
	update := bson.M{
		"$set": bson.M{
			"status":     req.Status,
			"updated_at": time.Now(),
		},
	}

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户状态失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// ResetPassword 重置用户密码
func ResetPassword(c *gin.Context) {
	userId := c.Param("id")
	var req struct {
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	collection := config.GetCollection("users")
	update := bson.M{
		"$set": bson.M{
			"password":   string(hashedPassword),
			"updated_at": time.Now(),
		},
	}

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "重置密码失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码重置成功"})
}

// CreateUser 创建新用户
func CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithError(err).Error("创建用户请求数据绑定失败")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logrus.WithFields(logrus.Fields{
		"username": req.Username,
		"role":     req.Role,
	}).Info("收到创建用户请求")

	// 检查用户名是否已存在
	collection := config.GetCollection("users")
	count, _ := collection.CountDocuments(context.Background(), bson.M{"username": req.Username})
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 创建用户
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "创建成功"})
}

// CreateNewUser 创建新用户
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

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 删除指定用户，管理员账户不能删除
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Success 200 {object} Response
// @Failure 400,403,404 {object} ErrorResponse
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	// 获取目标用户
	collection := config.GetCollection("users")
	var targetUser models.User
	objID, _ := primitive.ObjectIDFromHex(userID)
	err := collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&targetUser)
	if err != nil {
		c.JSON(404, gin.H{"error": "用户不存在"})
		return
	}

	// 检查是否是管理员账号
	if targetUser.Role == "admin" {
		c.JSON(403, gin.H{"error": "不能删除管理员账号"})
		return
	}

	// 删除用户
	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(500, gin.H{"error": "删除用户失败"})
		return
	}

	logrus.WithFields(logrus.Fields{
		"operator":    c.GetString("username"),
		"target_user": userID,
	}).Info("删除用户操作")

	c.JSON(200, gin.H{"message": "用户删除成功"})
}
