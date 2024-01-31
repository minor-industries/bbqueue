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
	"os"
	"path/filepath"
	"time"
)

func setup() (string, error) {
	datadir := os.ExpandEnv("$HOME/.bbqueue")

	err := os.MkdirAll(datadir, 0o700)
	if err != nil {
		return "", errors.Wrap(err, "make data dir")
	}

	return datadir, nil
}

func run() error {
	datadir, err := setup()
	if err != nil {
		return errors.Wrap(err, "setup")
	}

	db, err := database.Get(filepath.Join(datadir, "bbqueue.db"))
	if err != nil {
		return errors.Wrap(err, "get db")
	}

	go server(db)

	for {
		err := radio.Poll(func(probeName string, temp float64) error {
			//fmt.Println("callback", probeName, temp)

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
	gin.SetMode(gin.ReleaseMode)
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
