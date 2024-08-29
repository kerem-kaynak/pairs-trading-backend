package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"gonum.org/v1/gonum/stat"
	"gorm.io/gorm"
)

type ModelHandler struct {
	DB              *gorm.DB
	QuantServiceURL string
}

type PairRequest struct {
	Ticker1 string `json:"ticker_1" binding:"required"`
	Ticker2 string `json:"ticker_2" binding:"required"`
}

type OHLCData struct {
	Date  time.Time
	Close float64
}

type SpreadData struct {
	Date   string  `json:"date"`
	Spread float64 `json:"spread"`
}

type TrendResponse struct {
	Date       string  `json:"date"`
	Trend      string  `json:"trend"`
	Confidence float64 `json:"confidence"`
}

func NewModelHandler(db *gorm.DB, quantServiceURL string) *ModelHandler {
	return &ModelHandler{
		DB:              db,
		QuantServiceURL: quantServiceURL,
	}
}

func (h *ModelHandler) ComputeRLRT(c *gin.Context) {
	var req PairRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ohlcData1, err := h.fetchOHLCData(req.Ticker1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error fetching data for %s: %v", req.Ticker1, err)})
		return
	}

	ohlcData2, err := h.fetchOHLCData(req.Ticker2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error fetching data for %s: %v", req.Ticker2, err)})
		return
	}

	alignedData := alignData(ohlcData1, ohlcData2)

	spread := calculateSpread(alignedData)

	rlrtResponse, err := h.sendToQuantService(spread)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error from trend service: %v", err)})
		return
	}

	c.JSON(http.StatusOK, rlrtResponse)
}

func (h *ModelHandler) fetchOHLCData(ticker string) ([]OHLCData, error) {
	var data []OHLCData
	result := h.DB.Table("gold.etf_daily_ohlc").
		Select("date, close").
		Where("ticker = ?", ticker).
		Order("date DESC").
		Limit(30).
		Find(&data)

	if result.Error != nil {
		return nil, result.Error
	}

	return data, nil
}

func alignData(data1, data2 []OHLCData) map[time.Time][2]float64 {
	aligned := make(map[time.Time][2]float64)
	for _, d1 := range data1 {
		for _, d2 := range data2 {
			if d1.Date == d2.Date {
				aligned[d1.Date] = [2]float64{d1.Close, d2.Close}
				break
			}
		}
	}
	return aligned
}

func calculateSpread(data map[time.Time][2]float64) []SpreadData {
	var x, y []float64
	var dates []time.Time

	for date, prices := range data {
		x = append(x, prices[0])
		y = append(y, prices[1])
		dates = append(dates, date)
	}

	alpha, beta := stat.LinearRegression(x, y, nil, false)

	var spread []SpreadData
	for i, date := range dates {
		spreadValue := y[i] - (alpha + beta*x[i])
		spread = append(spread, SpreadData{
			Date:   date.Format("2006-01-02"),
			Spread: spreadValue,
		})
	}

	sort.Slice(spread, func(i, j int) bool {
		return spread[i].Date < spread[j].Date
	})

	return spread
}

func (h *ModelHandler) sendToQuantService(spread []SpreadData) ([]TrendResponse, error) {
	if len(spread) < 10 {
		return nil, fmt.Errorf("insufficient data points for RLRT, minimum required: 10, got: %d", len(spread))
	}

	payload := map[string]interface{}{
		"data": spread,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload: %v", err)
	}

	req, err := http.NewRequest("POST", h.QuantServiceURL+"/ml/rlrt", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	apiKey := os.Getenv("QUANT_SERVICE_API_KEY")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request to trend service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("trend service returned non-OK status: %d", resp.StatusCode)
	}

	var trendResponse []TrendResponse
	if err := json.NewDecoder(resp.Body).Decode(&trendResponse); err != nil {
		return nil, fmt.Errorf("error decoding trend service response: %v", err)
	}

	return trendResponse, nil
}
