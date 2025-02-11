package main

import (
	"fmt"
	"game-admin/config"
	"game-admin/handlers"
	"game-admin/middleware"
	"io"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 初始化日志配置
func initLogger() {
	// 创建日志目录
	if err := os.MkdirAll("logs", 0755); err != nil {
		panic(fmt.Sprintf("创建日志目录失败: %v", err))
	}

	// 设置日志格式
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	// 创建不同级别的日志文件
	infoWriter := &lumberjack.Logger{
		Filename:   path.Join("logs", "info.log"),
		MaxSize:    100,  // 单个文件最大尺寸，单位 MB
		MaxBackups: 30,   // 最多保留 30 个备份
		MaxAge:     7,    // 最多保留 7 天
		Compress:   true, // 是否压缩
	}

	errorWriter := &lumberjack.Logger{
		Filename:   path.Join("logs", "error.log"),
		MaxSize:    100,
		MaxBackups: 30,
		MaxAge:     7,
		Compress:   true,
	}

	// 创建 Hook
	logrus.AddHook(&LogHook{
		writers: map[logrus.Level]io.Writer{
			logrus.InfoLevel:  infoWriter,
			logrus.WarnLevel:  errorWriter,
			logrus.ErrorLevel: errorWriter,
			logrus.FatalLevel: errorWriter,
			logrus.PanicLevel: errorWriter,
		},
	})

	// 设置日志级别
	logrus.SetLevel(logrus.InfoLevel)
}

// LogHook 自定义 Hook
type LogHook struct {
	writers map[logrus.Level]io.Writer
}

func (hook *LogHook) Fire(entry *logrus.Entry) error {
	if writer, ok := hook.writers[entry.Level]; ok {
		line, err := entry.String()
		if err != nil {
			return err
		}
		_, err = writer.Write([]byte(line))
		return err
	}
	return nil
}

func (hook *LogHook) Levels() []logrus.Level {
	var levels []logrus.Level
	for level := range hook.writers {
		levels = append(levels, level)
	}
	return levels
}

// 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		logrus.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"method":       reqMethod,
			"uri":          reqUri,
		}).Info("HTTP请求")
	}
}

func main() {
	// 初始化日志
	initLogger()
	logrus.Info("日志系统初始化完成")

	// 初始化MongoDB连接
	config.InitMongoDB()

	r := gin.Default()

	// 使用日志中间件
	r.Use(LoggerMiddleware())

	// 加载HTML模板
	r.LoadHTMLGlob("templates/*")

	// 根路由 - 提供登录页面
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	// 登录相关路由
	r.POST("/login", handlers.Login)

	// 需要认证的路由
	api := r.Group("")
	api.Use(middleware.AuthMiddleware())
	{
		// 用户管理
		api.POST("/users", middleware.CheckPermission("manage_users"), handlers.CreateUser)
		api.GET("/users", middleware.CheckPermission("manage_users"), handlers.GetUsers)
		api.PUT("/users/:id/status", middleware.CheckPermission("manage_users"), handlers.UpdateUserStatus)
		api.PUT("/users/:id/password", middleware.CheckPermission("manage_users"), handlers.ResetPassword)

		// 玩家相关
		api.GET("/player/info/:id", middleware.CheckPermission("player_info"), handlers.GetPlayerInfo)
		api.POST("/player/item", middleware.CheckPermission("give_item"), handlers.GivePlayerItem)

		// 系统相关
		api.POST("/system/execute-script", middleware.CheckPermission("execute_script"), handlers.ExecuteScript)
	}

	// 添加用户管理页面路由
	r.GET("/user-management", func(c *gin.Context) {
		c.HTML(200, "user_management.html", nil)
	})

	r.Run(":8080")

	logrus.Info("服务器启动成功，端口：8080")
}
