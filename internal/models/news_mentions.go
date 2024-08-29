package models

import (
	"time"

	"github.com/google/uuid"
)

type NewsMention struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	Ticker       string     `gorm:"type:varchar(10);not null;index" json:"ticker"`
	Publisher    string     `json:"publisher"`
	Author       string     `json:"author"`
	SourceURL    string     `json:"source_url"`
	ShortSummary string     `json:"short_summary"`
	PublishTime  *time.Time `json:"publish_time,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (NewsMention) TableName() string {
	return "gold.news_mentions"
}
