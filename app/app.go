package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minor-industries/bbqueue/html"
	"github.com/minor-industries/bbqueue/radio"
	"github.com/minor-industries/rtgraph"
	"github.com/minor-industries/rtgraph/database"
	"github.com/pkg/errors"
	"os"
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

func Run() error {
	errCh := make(chan error)

	db, err := database.Get(os.ExpandEnv("$HOME/bbqueue.db"))
	if err != nil {
		return errors.Wrap(err, "get database")
	}

	graph, err := rtgraph.New(
		&database.Backend{DB: db},
		errCh,
		[]string{
			"bbqueue_bbq01_bbq",
			"bbqueue_bbq01_meat",
			"bbqueue_bbq01_voltage",
		},
	)
	if err != nil {
		return errors.Wrap(err, "new rtgraph")
	}

	graph.StaticFiles(html.FS,
		"index.html", "text/html",
	)

	go func() {
		gin.SetMode(gin.ReleaseMode)
		errCh <- graph.RunServer("0.0.0.0:8076")
	}()

	go func() {
		for {
			err := radio.Poll(func(probeName string, temp float64) error {
				//fmt.Println("callback", probeName, temp)
				now := time.Now()
				name := strings.ReplaceAll(probeName, "-", "_")
				err := graph.CreateValue("bbqueue_"+name, now, temp)
				return errors.Wrap(err, "create")
			})

			if err != nil {
				fmt.Println(errors.Wrap(err, "poll"))
			}
		}
	}()

	return <-errCh
}
