package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minor-industries/bbqueue/database"
	"github.com/minor-industries/bbqueue/radio"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

func run() error {
	db, err := database.Get("sqlite3.db")
	if err != nil {
		return errors.Wrap(err, "get db")
	}

	go server(db)

	for {
		err := radio.Poll(func(probeName string, temp float64) error {
			fmt.Println("callback", probeName, temp)

			result := db.Create(&database.Measurement{
				ProbeID: probeName,
				Time:    time.Now(),
				Temp:    temp,
			})

			return errors.Wrap(result.Error, "create")
		})

		if err != nil {
			fmt.Println("poll error:", err)
		}
	}
}

func server(db *gorm.DB) {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}
