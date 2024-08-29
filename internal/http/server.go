package http

import (
	"os"
	"pairs-trading-backend/internal/auth"
	"pairs-trading-backend/internal/config"
	"pairs-trading-backend/internal/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewServer(cfg *config.Config, db *gorm.DB) *gin.Engine {
	router := gin.Default()

	authHandler := handlers.NewAuthHandler(db, cfg)
	tickerHandler := handlers.NewTickerHandler(db)
	pairHandler := handlers.NewPairHandler(db)
	modelHandler := handlers.NewModelHandler(db, os.Getenv("QUANT_SERVICE_HOST"))
	modelChosenPairsHandler := handlers.NewModelChosenPairsHandler(db)

	router.GET("/auth/google/login", authHandler.GoogleLogin)
	router.GET("/auth/google/callback", authHandler.GoogleCallback)

	protected := router.Group("/api")
	protected.Use(auth.AuthMiddleware(cfg.JWTSecret))
	{
		protected.GET("/ticker/:ticker/details", tickerHandler.GetTickerDetails)
		protected.GET("/ticker/:ticker/daily-ohlc", tickerHandler.GetETFDailyOHLC)
		protected.GET("/ticker/:ticker/news-mentions", tickerHandler.GetTickerNews)

		protected.GET("/pairs", pairHandler.GetAllSuggestedPairs)
		protected.GET("/pairs/:id", pairHandler.GetSuggestedPairByID)

		protected.POST("/ml/fit-rlrt", modelHandler.ComputeRLRT)

		protected.GET("/model-chosen-pairs", modelChosenPairsHandler.GetAllModelChosenPairs)
		protected.GET("/model-chosen-pairs/trades", modelChosenPairsHandler.GetModelChosenPairTradesByPair)
		protected.GET("/model-chosen-pairs/metrics", modelChosenPairsHandler.GetModelChosenPairMetricsByPair)
	}

	return router
}
