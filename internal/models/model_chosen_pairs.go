package models

import "github.com/google/uuid"

type ModelChosenPair struct {
	ID      uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Ticker1 string    `gorm:"column:ticker_1;type:varchar(10);not null" json:"ticker_1"`
	Ticker2 string    `gorm:"column:ticker_2;type:varchar(10);not null" json:"ticker_2"`
}

func (ModelChosenPair) TableName() string {
	return "gold.model_chosen_pairs"
}
