package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Expense struct {
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
}

var db *sql.DB

func main() {

	connStr := "user=postgres password=hiddenfornow dbname=finance-tracker sslmode=disable"

	var err error

	db, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to PostgreSQL successfully!")

	r := gin.Default()

	r.POST("/add-expense", func(c *gin.Context) {

		var expense Expense

		if err := c.ShouldBindJSON(&expense); err != nil {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid JSON",
			})

			return
		}

		_, err := db.Exec(
			"INSERT INTO expenses (category, amount) VALUES ($1, $2)",
			expense.Category,
			expense.Amount,
		)

		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Expense saved to database!",
		})
	})

	r.Run(":8080")
}