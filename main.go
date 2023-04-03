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

var products = []Product{
	{ID: 1, Name: "Zapatillas Nike Air Max", Price: 14999},
	{ID: 2, Name: "Raqueta de tenis Wilson", Price: 9999},
	{ID: 3, Name: "Pelotas de tenis Penn", Price: 499},
}

func main() {
	router := gin.Default()

	router.GET("/products", func(c *gin.Context) {
		c.JSON(http.StatusOK, products)
	})

	router.POST("/orders", func(c *gin.Context) {
		var order Order
		c.BindJSON(&order)

		var total int
		var productsList []string
		for _, productID := range order.Products {
			for _, product := range products {
				if product.Name == productID {
					total += product.Price
					productsList = append(productsList, product.Name)
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "Order created successfully",
			"products": productsList,
			"total":    total,
		})
	})

	router.Run(":8080")
}
