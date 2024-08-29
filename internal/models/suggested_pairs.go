package models

import (
	"time"

	"github.com/google/uuid"
)

type SuggestedPair struct {
	ID                  uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Ticker1             string    `gorm:"column:ticker_1;type:varchar(10);not null" json:"ticker_1"`
	Ticker2             string    `gorm:"column:ticker_2;type:varchar(10);not null" json:"ticker_2"`
	SpreadIntercept     float64   `json:"spread_intercept"`
	SpreadSlope         float64   `json:"spread_slope"`
	CointegrationPValue float64   `json:"cointegration_p_value"`
	HalfLife            float64   `json:"half_life"`
	MeanCrossings       float64   `json:"mean_crossings"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (SuggestedPair) TableName() string {
	return "gold.suggested_pairs"
}
