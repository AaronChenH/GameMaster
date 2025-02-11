package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ScriptRequest struct {
	ScriptName string                 `json:"script_name" binding:"required"`
	Params     map[string]interface{} `json:"params"`
}

func ExecuteScript(c *gin.Context) {
	var req ScriptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// 这里添加执行脚本的逻辑
	
	c.JSON(http.StatusOK, gin.H{
		"message":     "脚本执行成功",
		"script_name": req.ScriptName,
		"params":      req.Params,
	})
} 