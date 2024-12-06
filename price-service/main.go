package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CryptoPrice struct {
	Symbol    string    `json:"symbol" bson:"symbol"`
	Price     float64   `json:"price" bson:"price"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}

type PriceService struct {
	client     *resty.Client
	collection *mongo.Collection
}

func NewPriceService(collection *mongo.Collection) *PriceService {
	client := resty.New()
	client.SetHeader("X-CMC_PRO_API_KEY", os.Getenv("COINMARKETCAP_API_KEY"))
	return &PriceService{
		client:     client,
		collection: collection,
	}
}

func (s *PriceService) FetchPrices(symbols []string) ([]CryptoPrice, error) {
	// CoinMarketCap API endpoint
	url := "https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest"

	// Make API request
	resp, err := s.client.R().
		SetQueryParam("symbol", joinSymbols(symbols)).
		Get(url)

	if err != nil {
		return nil, err
	}

	// Parse response and store in MongoDB
	prices := make([]CryptoPrice, 0)
	// ... implement response parsing ...

	return prices, nil
}

func (s *PriceService) StorePrices(prices []CryptoPrice) error {
	documents := make([]interface{}, len(prices))
	for i, price := range prices {
		documents[i] = price
	}

	_, err := s.collection.InsertMany(context.Background(), documents)
	return err
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

	collection := client.Database("coinsight").Collection("prices")
	priceService := NewPriceService(collection)

	// Setup Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Get latest prices endpoint
	r.GET("/prices", func(c *gin.Context) {
		symbols := []string{"BTC", "ETH", "XRP", "DOGE"} // Default symbols
		prices, err := priceService.FetchPrices(symbols)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, prices)
	})

	// Get historical prices endpoint
	r.GET("/prices/historical", func(c *gin.Context) {
		symbol := c.Query("symbol")
		if symbol == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "symbol is required"})
			return
		}

		filter := bson.M{"symbol": symbol}
		cursor, err := collection.Find(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var prices []CryptoPrice
		if err := cursor.All(ctx, &prices); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, prices)
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}

func joinSymbols(symbols []string) string {
	result := ""
	for i, symbol := range symbols {
		if i > 0 {
			result += ","
		}
		result += symbol
	}
	return result
}
