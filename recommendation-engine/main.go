package main

import (
	"context"
	"log"
	"math"
	"net/http"
	"os"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gonum.org/v1/gonum/stat"
)

type Portfolio struct {
	UserID    string                  `json:"user_id" bson:"user_id"`
	Holdings  map[string]float64      `json:"holdings" bson:"holdings"`
	Risk      float64                 `json:"risk" bson:"risk"`
	Returns   []float64               `json:"returns" bson:"returns"`
	Metadata  map[string]interface{}  `json:"metadata" bson:"metadata"`
}

type Recommendation struct {
	Action     string  `json:"action"`
	Symbol     string  `json:"symbol"`
	Percentage float64 `json:"percentage"`
	Reason     string  `json:"reason"`
}

type RecommendationEngine struct {
	collection *mongo.Collection
}

func NewRecommendationEngine(collection *mongo.Collection) *RecommendationEngine {
	return &RecommendationEngine{
		collection: collection,
	}
}

func (e *RecommendationEngine) calculatePortfolioMetrics(portfolio Portfolio) (float64, float64) {
	if len(portfolio.Returns) < 2 {
		return 0, 0
	}

	// Calculate mean return
	mean := stat.Mean(portfolio.Returns, nil)

	// Calculate volatility (standard deviation)
	variance := stat.Variance(portfolio.Returns, nil)
	volatility := math.Sqrt(variance)

	return mean, volatility
}

func (e *RecommendationEngine) generateRecommendations(portfolio Portfolio) []Recommendation {
	mean, volatility := e.calculatePortfolioMetrics(portfolio)
	recommendations := make([]Recommendation, 0)

	// Simple recommendation logic based on risk-return profile
	if volatility > 0.2 { // High volatility
		recommendations = append(recommendations, Recommendation{
			Action:     "REDUCE",
			Symbol:     findMostVolatileHolding(portfolio),
			Percentage: 10,
			Reason:     "Portfolio volatility is high. Consider reducing exposure to volatile assets.",
		})
	}

	if mean < 0.05 { // Low returns
		recommendations = append(recommendations, Recommendation{
			Action:     "ADD",
			Symbol:     "BTC", // Default to Bitcoin as a suggestion
			Percentage: 5,
			Reason:     "Portfolio returns are below target. Consider adding established cryptocurrencies.",
		})
	}

	// Diversification check
	if len(portfolio.Holdings) < 3 {
		recommendations = append(recommendations, Recommendation{
			Action:     "ADD",
			Symbol:     "ETH", // Ethereum as a diversification suggestion
			Percentage: 10,
			Reason:     "Portfolio could benefit from increased diversification.",
		})
	}

	return recommendations
}

func findMostVolatileHolding(portfolio Portfolio) string {
	// In a real implementation, this would analyze historical price data
	// For now, return the largest holding
	var maxHolding string
	var maxAmount float64

	for symbol, amount := range portfolio.Holdings {
		if amount > maxAmount {
			maxAmount = amount
			maxHolding = symbol
		}
	}

	return maxHolding
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Connect to MongoDB
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database("coinsight").Collection("portfolios")
	engine := NewRecommendationEngine(collection)

	// Setup Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Get recommendations endpoint
	r.GET("/recommendations/:user_id", func(c *gin.Context) {
		userID := c.Param("user_id")

		var portfolio Portfolio
		err := collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&portfolio)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Portfolio not found"})
			return
		}

		recommendations := engine.generateRecommendations(portfolio)
		c.JSON(http.StatusOK, recommendations)
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	r.Run(":" + port)
}
