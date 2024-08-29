// internal/handlers/pair_handler.go

package handlers

import (
	"net/http"
	"pairs-trading-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PairHandler struct {
	DB *gorm.DB
}

func NewPairHandler(db *gorm.DB) *PairHandler {
	return &PairHandler{DB: db}
}

func (h *PairHandler) GetAllSuggestedPairs(c *gin.Context) {
	var pairs []models.SuggestedPair

	if err := h.DB.Find(&pairs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching suggested pairs"})
		return
	}

	c.JSON(http.StatusOK, pairs)
}

func (h *PairHandler) GetSuggestedPairByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	var pair models.SuggestedPair

	if err := h.DB.First(&pair, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Suggested pair not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching suggested pair"})
		}
		return
	}

	c.JSON(http.StatusOK, pair)
}
