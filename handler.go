package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type Order struct {
	ID       int      `json:"id"`
	Products []string `json:"products"`
	Total    int      `json:"total"`
}

type OrderResponse struct {
	Message  string   `json:"message"`
	Products []string `json:"products"`
	Total    int      `json:"total"`
}

// Catalog is the in-memory product store. Kept as a package-level var to
// preserve the original service behaviour; tests swap it via resetStore.
var products = []Product{
	{ID: 1, Name: "Zapatillas Nike Air Max", Price: 14999},
	{ID: 2, Name: "Raqueta de tenis Wilson", Price: 9999},
	{ID: 3, Name: "Pelotas de tenis Penn", Price: 499},
}

// calculateOrder resolves the product names in the request to a total and
// returns the matched products. Unknown products are ignored; duplicates
// are counted once each occurrence (original behaviour preserved).
func calculateOrder(requested []string, catalog []Product) (names []string, total int) {
	for _, wanted := range requested {
		for _, p := range catalog {
			if p.Name == wanted {
				total += p.Price
				names = append(names, p.Name)
			}
		}
	}
	return names, total
}

func listProducts(c *gin.Context) {
	c.JSON(http.StatusOK, products)
}

func createOrder(c *gin.Context) {
	var order Order
	// The original used c.BindJSON and ignored the error, replying 200 with
	// an empty order for malformed payloads. Bind explicitly and fail 400.
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body: " + err.Error()})
		return
	}
	names, total := calculateOrder(order.Products, products)
	c.JSON(http.StatusOK, OrderResponse{
		Message:  "Order created successfully",
		Products: names,
		Total:    total,
	})
}

// newRouter wires the HTTP handlers (separated from main so tests can
// exercise it in-process).
func newRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.GET("/products", listProducts)
	router.POST("/orders", createOrder)

	return router
}
