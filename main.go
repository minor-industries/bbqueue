package main

import (
	"fmt"
	"github.com/minor-industries/bbqueue/database"
	"github.com/minor-industries/bbqueue/radio"
	"github.com/pkg/errors"
	"time"
)

func run() error {
	db, err := database.Get("sqlite3.db")
	if err != nil {
		return errors.Wrap(err, "get db")
	}

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

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}
