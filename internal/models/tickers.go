package models

import "time"

type Ticker struct {
	Ticker    string    `gorm:"type:varchar(10);primary_key" json:"ticker"`
	Name      string    `json:"name"`
	Active    bool      `json:"active"`
	Type      string    `gorm:"index" json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Ticker) TableName() string {
	return "gold.tickers"
}
