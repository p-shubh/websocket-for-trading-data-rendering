package models

type Payload struct {
	Rate      float64 `json:"rate" gorm:"column:rate"`
	Volume    float64 `json:"volume" gorm:"column:volume"`
	OrderType string  `json:"order_type" gorm:"column:order_type"`
	Timestamp string  `json:"timestamp" gorm:"column:timestamp;type:timestamp"`
}
