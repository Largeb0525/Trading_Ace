package server

import (
	"log"
	"net/http"

	"github.com/Largeb0525/Trading_Ace/eth"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func StartServer() {
	go eth.ListenToContractEvents()
	go ProcessSharePoolTicker()
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	r.POST("/Campaign", CreateCampaignHandler)
	r.GET("/user/task/status", GetUserTaskStatusHandler)
	r.GET("/user/points", GetUserPointsHistoryHandler)

	port := viper.GetString("server.port")
	err := r.Run(":" + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
