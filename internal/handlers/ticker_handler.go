package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"pairs-trading-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TickerDetailsResponse struct {
	models.Ticker
	SMA  float64 `json:"sma"`
	EMA  float64 `json:"ema"`
	RSI  float64 `json:"rsi"`
	MACD float64 `json:"macd"`
}

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

	extendedDetails := TickerDetailsResponse{
		Ticker: tickerDetails,
	}

	polygonAPIKey := os.Getenv("POLYGON_API_KEY")
	if polygonAPIKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Polygon API key not set"})
		return
	}

	baseURL := "https://api.polygon.io/v1/indicators"

	indicators := map[string]*float64{
		"sma":  &extendedDetails.SMA,
		"ema":  &extendedDetails.EMA,
		"rsi":  &extendedDetails.RSI,
		"macd": &extendedDetails.MACD,
	}

	for indicator, value := range indicators {
		url := fmt.Sprintf("%s/%s/%s?timespan=day&adjusted=true&window=30&series_type=close&order=desc&limit=1&apiKey=%s", baseURL, indicator, ticker, polygonAPIKey)
		if indicator == "macd" {
			url = fmt.Sprintf("%s/macd/%s?timespan=day&adjusted=true&short_window=12&long_window=26&signal_window=9&series_type=close&order=desc&limit=1&apiKey=%s", baseURL, ticker, polygonAPIKey)
		}

		indicatorValue, err := fetchIndicator(url)
		if err != nil {
			log.Printf("Error fetching %s: %v", indicator, err)
			continue
		}
		*value = indicatorValue
	}

	c.JSON(http.StatusOK, extendedDetails)
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

func (h *TickerHandler) GetTickerNews(c *gin.Context) {
	ticker := c.Param("ticker")

	var newsMentions []models.NewsMention
	result := h.DB.Where("ticker = ?", ticker).Find(&newsMentions)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve ticker news mentions"})
		return
	}

	if len(newsMentions) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No news mentions found for the given ticker"})
		return
	}

	c.JSON(http.StatusOK, newsMentions)
}

func fetchIndicator(url string) (float64, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Results struct {
			Values []struct {
				Value float64 `json:"value"`
			} `json:"values"`
		} `json:"results"`
		Status string `json:"status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	if result.Status != "OK" || len(result.Results.Values) == 0 {
		return 0, fmt.Errorf("invalid response or no data available")
	}

	return result.Results.Values[0].Value, nil
}
