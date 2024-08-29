package models

type ModelChosenPairMetrics struct {
	Ticker1          string  `gorm:"column:ticker_1;type:varchar(10);not null" json:"ticker_1"`
	Ticker2          string  `gorm:"column:ticker_2;type:varchar(10);not null" json:"ticker_2"`
	TotalReturn      float64 `json:"total_return"`
	AnnualizedReturn float64 `json:"annualized_return"`
	MaxDrawdown      float64 `json:"max_drawdown"`
}

func (ModelChosenPairMetrics) TableName() string {
	return "gold.model_success_metrics"
}
