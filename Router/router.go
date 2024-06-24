package router

import (
	"log"
	webSocketOperations "trade/WebSocketOperations"

	"github.com/gin-gonic/gin"
)

func Router() {
	r := gin.Default()
	gin.SetMode(gin.DebugMode)
	// set gin method
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	WebSocketRouting(r)
	if err := r.Run(":9090"); err != nil {
		log.Println()
		return
	}
}

func WebSocketRouting(c *gin.Engine) {
	// ws://localhost:8080/ws
	c.GET("/CoinMarketHistory", webSocketOperations.CoinMarketHistory)
	c.GET("/WebSocketForCalGetCoinMarketHistory", webSocketOperations.WebSocketForGetCoinMarketHistory)
	// c.GET("/GetCoinMarketHistory", webSocketOperations.GetCoinMarketHistory)
	//REST API
	c.GET("/ApiGetCoinMarketHistory", webSocketOperations.ApiForGetCoinMarketHistory)

}
