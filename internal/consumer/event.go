package consumer

import "time"

type OrderEvent struct {
	OrderID    int       `json:"order_id"`
	UserID     int       `json:"user_id"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}
