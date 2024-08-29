package handlers

import (
	"net/http"
	"pairs-trading-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TickerHandler struct {
	DB *gorm.DB
}

func NewTickerHandler(db *gorm.DB) *TickerHandler {
	return &TickerHandler{DB: db}
}

func (h *TickerHandler) GetTickerDetails(c *gin.Context) {
	ticker := c.Param("ticker")

	var tickerDetails models.Ticker
	result := h.DB.Where("ticker = ?", ticker).First(&tickerDetails)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ticker not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve ticker details"})
		return
	}

	c.JSON(http.StatusOK, tickerDetails)
}

func (h *TickerHandler) GetETFDailyOHLC(c *gin.Context) {
	ticker := c.Param("ticker")

	var etfData []models.ETFDailyOHLC
	result := h.DB.Where("ticker = ?", ticker).Find(&etfData)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve ETF OHLC data"})
		return
	}

	if len(etfData) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No ETF OHLC data found for the given ticker"})
		return
	}

	c.JSON(http.StatusOK, etfData)
}
