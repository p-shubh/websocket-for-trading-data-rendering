package webSocketOperations

import (
	"fmt"
	"time"
	"trade/pkg/db"
)

func GetWebSocketDbOperation() RateStats {
	endTime := time.Now()
	startTime := endTime.Add(-10 * time.Minute)

	// var data []models.Payload
	var data RateStats

	result := db.DB.Raw("SELECT MAX(rate) AS max_rate, MIN(rate) AS min_rate, SUM(volume) AS sum_volume FROM payloads where timestamp >= '" + startTime.Format("2006-01-02 15:04:05") + "' AND timestamp <= '" + endTime.Format("2006-01-02 15:04:05") + "'").Debug().Scan(&data)
	if result.Error != nil {
		// handle error print the error
		fmt.Println("GetWebSocketDbOperation : ", result.Error)
		return RateStats{}
	}

	return data

}
