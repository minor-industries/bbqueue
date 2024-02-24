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
	"io/fs"
	"os"
	"path/filepath"
	"strings"
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

	r.GET("/data.csv", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/plain")
		c.Status(200)
		content, err := getCSV(db, c)
		if err != nil {
			c.Error(err)
			return
		}
		_, _ = c.Writer.Write(content)
	})

	files(r,
		"dygraph.js", "application/javascript",
		"dygraph.css", "text/css",
		"chart.html", "text/html",
	)

	r.Run("0.0.0.0:8080")
}

func getCSV(db *gorm.DB, c *gin.Context) ([]byte, error) {
	lines := []string{
		"Date,Temp",
	}

	probes, err := database.GetProbes(db)
	if err != nil {
		return nil, errors.Wrap(err, "get probes")
	}

	fmt.Println(probes)

	after, err := time.Parse("2006-01-02", "1980-01-19")
	if err != nil {
		panic("bad time")
	}
	data, err := database.GetProbeData(db, probes[1], after, time.Now())
	for _, datum := range data {
		lines = append(lines, fmt.Sprintf("%s,%f",
			datum.Time.Format("2006/01/02 15:04:05"),
			datum.Temp*9/5+32,
		))
	}

	return []byte(strings.Join(lines, "\n")), nil
}

func files(r *gin.Engine, files ...string) {
	for i := 0; i < len(files); i += 2 {
		name := files[i]
		ct := files[i+1]
		r.GET("/"+name, func(c *gin.Context) {
			header := c.Writer.Header()
			header["Content-Type"] = []string{ct}
			content, err := fs.ReadFile(html.FS, name)
			if err != nil {
				c.Status(404)
				return
			}
			_, _ = c.Writer.Write(content)
		})
	}
}

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}
