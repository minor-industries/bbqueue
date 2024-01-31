package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minor-industries/bbqueue/database"
	"github.com/minor-industries/bbqueue/html"
	"github.com/minor-industries/bbqueue/radio"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"html/template"
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

	templ := template.Must(template.New("").Funcs(templateFuncs).ParseFS(html.FS, "*.html"))
	r.SetHTMLTemplate(templ)

	r.GET("/", func(c *gin.Context) {
		data, err := database.GetLatestTemps(db)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.HTML(200, "index.html", map[string]any{
			"data": data,
			"now":  time.Now(),
		})
	})

	r.GET("/plot.svg", func(c *gin.Context) {
		plotHandler(db, c)
	})

	r.Run("0.0.0.0:8080")
}

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}
