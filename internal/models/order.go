package models

import (
	"fmt"
	"time"
)

type Order struct {
	OrderUID          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string  `json:"transaction"`
	RequestID    string  `json:"request_id"`
	Currency     string  `json:"currency"`
	Provider     string  `json:"provider"`
	Amount       float64 `json:"amount"`
	PaymentDt    int     `json:"payment_dt"`
	Bank         string  `json:"bank"`
	DeliveryCost float64 `json:"delivery_cost"`
	GoodsTotal   float64 `json:"goods_total"`
	CustomFee    float64 `json:"custom_fee"`
}

type Item struct {
	ChrtID      int     `json:"chrt_id"`
	TrackNumber string  `json:"track_number"`
	Price       float64 `json:"price"`
	Rid         string  `json:"rid"`
	Name        string  `json:"name"`
	Sale        int     `json:"sale"`
	Size        string  `json:"size"`
	TotalPrice  float64 `json:"total_price"`
	NmID        int     `json:"nm_id"`
	Brand       string  `json:"brand"`
	Status      int     `json:"status"`
}

func ValidateOrder(order Order) error {
	if order.OrderUID == "" {
		return fmt.Errorf("order_uid is required")
	}
	if order.TrackNumber == "" {
		return fmt.Errorf("track_number is required")
	}
	if order.Entry == "" {
		return fmt.Errorf("entry is required")
	}
	if order.CustomerID == "" {
		return fmt.Errorf("customer_id is required")
	}
	if order.Delivery.Name == "" || order.Delivery.Phone == "" {
		return fmt.Errorf("delivery name and phone are required")
	}
	if order.Payment.Transaction == "" || order.Payment.Amount <= 0 {
		return fmt.Errorf("payment transaction and amount are required and must be > 0")
	}
	if len(order.Items) == 0 {
		return fmt.Errorf("at least one item is required")
	}
	for i, item := range order.Items {
		if item.Rid == "" || item.Name == "" {
			return fmt.Errorf("item[%d] missing required fields", i)
		}
	}
	return nil
}
