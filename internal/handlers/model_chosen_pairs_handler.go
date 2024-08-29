package handlers

import (
	"math"
	"net/http"
	"pairs-trading-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ModelChosenPairsHandler struct {
	DB *gorm.DB
}

func NewModelChosenPairsHandler(db *gorm.DB) *ModelChosenPairsHandler {
	return &ModelChosenPairsHandler{DB: db}
}

func (h *ModelChosenPairsHandler) GetModelChosenPairTradesByPair(c *gin.Context) {
	var trades []models.ModelChosenPairsTrade

	ticker1 := c.Query("ticker_1")
	ticker2 := c.Query("ticker_2")

	if err := h.DB.Where("ticker_1 = ? AND ticker_2 = ?", ticker1, ticker2).Find(&trades).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching trades"})
		return
	}

	for i := range trades {
		trades[i].Spread = handleNaN(trades[i].Spread)
		trades[i].PredictedTrend = handleNaN(trades[i].PredictedTrend)
		trades[i].Position = handleNaN(trades[i].Position)
		trades[i].Budget = handleNaN(trades[i].Budget)
	}

	c.JSON(http.StatusOK, trades)
}

func (h *ModelChosenPairsHandler) GetModelChosenPairMetricsByPair(c *gin.Context) {
	var metrics models.ModelChosenPairMetrics

	ticker1 := c.Query("ticker_1")
	ticker2 := c.Query("ticker_2")

	if err := h.DB.Where("ticker_1 = ? AND ticker_2 = ?", ticker1, ticker2).First(&metrics).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Metrics not found for the given pair"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching metrics"})
		}
		return
	}

	metrics.TotalReturn = handleNaN(metrics.TotalReturn)
	metrics.AnnualizedReturn = handleNaN(metrics.AnnualizedReturn)
	metrics.MaxDrawdown = handleNaN(metrics.MaxDrawdown)

	c.JSON(http.StatusOK, metrics)
}

func (h *ModelChosenPairsHandler) GetAllModelChosenPairs(c *gin.Context) {
	var pairs []models.ModelChosenPair

	if err := h.DB.Find(&pairs).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "No pairs found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching model chosen pairs"})
		}
		return
	}

	c.JSON(http.StatusOK, pairs)
}

func handleNaN(f float64) float64 {
	if math.IsNaN(f) {
		return 0
	}
	return f
}
