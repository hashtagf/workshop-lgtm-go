package product

import (
	"fmt"
	"net/http"
	"shared"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type ProductHandler struct {
	Client *mongo.Client
}

func (h *ProductHandler) GetAllProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a new span for this request
		span1Ctx, span := shared.StartNewSpan(c.Request.Context(), "service-c", "GetProducts")
		defer span.End()

		// Simulate some work
		time.Sleep(100 * time.Millisecond)

		// Create a span for the business logic
		span2Ctx, businessLogicSpan := shared.StartNewSpan(span1Ctx, "service-c", "BusinessLogic")
		defer businessLogicSpan.End()

		// Get products from mongoDB
		var products []Product
		cursor, err := h.Client.Database("test").Collection("products").Find(span2Ctx, bson.M{})
		if err != nil {
			fmt.Print("Failed to fetch products: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
			return
		}
		defer cursor.Close(span2Ctx)

		if err := cursor.All(span2Ctx, &products); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
			return
		}

		c.JSON(http.StatusOK, products)
	}
}
