package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PlayerInfo struct {
	PlayerID  string `json:"player_id"`
	Nickname  string `json:"nickname"`
	Level     int    `json:"level"`
	VipLevel  int    `json:"vip_level"`
	Items     []Item `json:"items"`
	CreatedAt string `json:"created_at"`
}

type Item struct {
	ItemID string `json:"item_id"`
	Amount int    `json:"amount"`
}

type GiveItemRequest struct {
	PlayerID string `json:"player_id" binding:"required"`
	ItemID   string `json:"item_id" binding:"required"`
	Amount   int    `json:"amount" binding:"required"`
}

const (
	GameServerURL = "http://game-server:8080" // 游戏服务器地址
)

func GetPlayerInfo(c *gin.Context) {
	playerID := c.Param("id")

	// 从游戏服务器获取玩家数据
	resp, err := http.Get(fmt.Sprintf("%s/api/player/%s", GameServerURL, playerID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "游戏服务器连接失败"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "玩家不存在"})
		return
	}

	var player PlayerInfo
	if err := json.NewDecoder(resp.Body).Decode(&player); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据解析失败"})
		return
	}

	c.JSON(http.StatusOK, player)
}

func GivePlayerItem(c *gin.Context) {
	var req GiveItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用游戏服务器接口发放道具
	itemReq, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "请求数据错误"})
		return
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/api/player/give-item", GameServerURL),
		"application/json",
		bytes.NewBuffer(itemReq),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "游戏服务器连接失败"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": "道具发放失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "道具发放成功",
		"data":    req,
	})
}
