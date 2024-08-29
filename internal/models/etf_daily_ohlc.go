package models

import (
	"time"
)

type ETFDailyOHLC struct {
	Ticker              string     `gorm:"primaryKey;type:varchar(10)" json:"ticker"`
	Name                string     `json:"name"`
	Open                float64    `json:"open"`
	High                float64    `json:"high"`
	Low                 float64    `json:"low"`
	Close               float64    `json:"close"`
	Date                time.Time  `gorm:"primaryKey;type:date" json:"date"`
	Volume              float64    `json:"volume"`
	VolumeWeightedPrice float64    `json:"volume_weighted_price"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	DeletedAt           *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (ETFDailyOHLC) TableName() string {
	return "gold.etf_daily_ohlc"
}
