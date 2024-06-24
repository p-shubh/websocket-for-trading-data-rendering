package webSocketOperations

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
	"trade/pkg/db"
	"trade/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// var upgrader1 = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

func CoinMarketHistory(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to set websocket upgrade:", err)
		return
	}
	defer conn.Close()

	internalConn, _, err := websocket.DefaultDialer.Dial("wss://stream.bit24hr.in/coin_market_history/", nil)
	if err != nil {
		log.Println("Failed to connect to internal websocket:", err)
		return
	}
	defer internalConn.Close()

	payload := map[string]string{"coin": "BTC"}
	if err := internalConn.WriteJSON(payload); err != nil {
		log.Println("Failed to send payload to internal websocket:", err)
		return
	}
	var currentTime time.Time

	messageChan := make(chan []byte)

	go func() {
		defer close(messageChan)
		for {
			_, message, err := internalConn.ReadMessage()
			if err != nil {
				log.Println("Error reading message from internal websocket:", err)
				return
			}
			messageChan <- message
		}
	}()

	for message := range messageChan {

		var payloads []models.Payload
		if err := json.Unmarshal(message, &payloads); err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		// for _, payload := range payloads {

		// 	var currentTime string

		// 	if len(currentTime) == 0 {
		// 		currentTime = "2006-01-02 15:04:05"
		// 	}

		// 	timestamp, err := time.Parse("2006-01-02 15:04:05", payload.Timestamp)
		// 	if err != nil {
		// 		log.Println("Error parsing timestamp:", err)
		// 		continue
		// 	}
		// 	fmt.Println("Time stamp : ", timestamp)
		// 	payload.Timestamp = timestamp.Format("2006-01-02 15:04:05")
		// 	log.Printf("%+v\n", payload)
		// 	time1, _ := time.Parse("2006-01-02 15:04:05", payload.Timestamp)
		// 	time2, _ := time.Parse("2006-01-02 15:04:05", currentTime)
		// 	if time1 > time2 {
		// 		if err := db.DB.Create(&payload).Error; err != nil {
		// 			log.Println("Error saving to database:", err)
		// 		}
		// 		currentTime = payload.Timestamp
		// 	}
		// }
		for _, payload := range payloads {
			// Parse the timestamp
			timestamp, err := time.Parse("2006-01-02 15:04:05", payload.Timestamp)
			if err != nil {
				log.Println("Error parsing timestamp:", err)
				continue
			}

			// Check if the timestamp is greater than the current timestamp
			if timestamp.After(currentTime) {
				currentTime = timestamp

				// Insert the payload record in the database
				if err := db.DB.Create(&payload).Error; err != nil {
					log.Println("Error saving to database:", err)
				}
			}

			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("Error writing message to external websocket:", err)
				return
			}
		}
	}
}

func GetCoinMarketHistory(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to set websocket upgrade:", err)
		return
	}
	defer conn.Close()

	// Start a new goroutine to handle the websocket logic
	go func() {
		for {
			endTime := time.Now()
			startTime := endTime.Add(-time.Minute)

			// Query to find highest and lowest rates within the time range
			var maxRate, minRate float64
			if err := db.DB.Model(&models.Payload{}).
				Select("MAX(rate) as max_rate, MIN(rate) as min_rate").
				Where("timestamp >= ? AND timestamp <= ?", startTime, endTime).
				Group("timestamp").
				Group("rate").
				Order("timestamp").
				Limit(1).
				First(&struct {
					MaxRate float64
					MinRate float64
				}{}).Error; err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					log.Fatalln("Error querying database:", err)
				}
				continue
			}

			// Prepare the response
			response := struct {
				MaxRate float64 `json:"max_rate"`
				MinRate float64 `json:"min_rate"`
			}{
				MaxRate: maxRate,
				MinRate: minRate,
			}

			// Convert struct to bytes
			jsonData, err := json.Marshal(response)
			if err != nil {
				fmt.Println("Error marshalling JSON:", err)
				return
			}

			// Send response to client
			if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
				log.Println("Error sending response to client:", err)
				continue
			}

			// Sleep for 1 minute before allowing the next request
			time.Sleep(time.Minute)
			// }
		}
	}()
}

// RateStats to hold the result of the query
type RateStats struct {
	MaxRate   float64 `json:"max_Rate"`
	MinRate   float64 `json:"min_rate"`
	SumVolume float64 `json:"sum_volume"`
}

func WebSocketForGetCoinMarketHistory(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to set websocket upgrade:", err)
		return
	}
	defer conn.Close()

	internalConn, _, err := websocket.DefaultDialer.Dial("wss://stream.bit24hr.in/coin_market_history/", nil)
	if err != nil {
		log.Println("Failed to connect to internal websocket:", err)
		return
	}
	defer internalConn.Close()

	payload := map[string]string{"coin": "BTC"}
	if err := internalConn.WriteJSON(payload); err != nil {
		log.Println("Failed to send payload to internal websocket:", err)
		return
	}
	var currentTime time.Time

	messageChan := make(chan []byte)

	go func() {
		defer close(messageChan)
		for {
			_, message, err := internalConn.ReadMessage()
			if err != nil {
				log.Println("Error reading message from internal websocket:", err)
				return
			}
			messageChan <- message
		}
	}()

	for message := range messageChan {

		var payloads []models.Payload
		if err := json.Unmarshal(message, &payloads); err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}
		for _, payload := range payloads {
			// Parse the timestamp
			timestamp, err := time.Parse("2006-01-02 15:04:05", payload.Timestamp)
			if err != nil {
				log.Println("Error parsing timestamp:", err)
				continue
			}

			// Check if the timestamp is greater than the current timestamp
			if timestamp.After(currentTime) {
				currentTime = timestamp

				// Insert the payload record in the database
				if err := db.DB.Create(&payload).Error; err != nil {
					log.Println("Error saving to database:", err)
				}
			}

			// Example: Send data every 5 seconds
			ticker := time.NewTicker(5 * time.Second)
			// defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					// Get data (e.g., from Redis)
					data := GetWebSocketDbOperation()

					reqBodyBytes := new(bytes.Buffer)
					json.NewEncoder(reqBodyBytes).Encode(data)

					conn.SetWriteDeadline(time.Now().Add(5 * time.Second))

					if err := conn.WriteMessage(websocket.TextMessage, reqBodyBytes.Bytes()); err != nil {
						fmt.Println("Failed to write the response messages : ", err.Error())
					}
				}
			}

			// if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			// 	log.Println("Error writing message to external websocket:", err)
			// 	return
			// }
		}
	}
}

func ApiForGetCoinMarketHistory(c *gin.Context) {

	// Prepare the response
	data := GetWebSocketDbOperation()

	c.JSON(200, gin.H{
		"data": data,
	})
}
