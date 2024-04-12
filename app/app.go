package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minor-industries/bbqueue/html"
	"github.com/minor-industries/bbqueue/radio"
	"github.com/minor-industries/rtgraph"
	"github.com/minor-industries/rtgraph/database"
	"github.com/minor-industries/rtgraph/schema"
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
			"bbqueue_bbq01_rssi",
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

	t0 := time.Now()
	go func() {
		lastSeen := map[string]schema.Value{}
		for {
			err := radio.Poll(func(seriesName string, value float64) error {
				now := time.Now()
				fullName := "bbqueue_" + strings.ReplaceAll(seriesName, "-", "_")

				last, ok := lastSeen[fullName]
				lastSeen[fullName] = schema.Value{
					Timestamp: now,
					Value:     value,
				}
				if ok {
					dt := now.Sub(last.Timestamp)
					if dt < time.Second {
						fmt.Printf("    %s %0.02f %f\n", fullName, now.Sub(t0).Seconds(), value)
						return nil
					}
				}

				fmt.Printf("add %s %0.02f %f\n", fullName, now.Sub(t0).Seconds(), value)
				err := graph.CreateValue(fullName, now, value)
				return errors.Wrap(err, "create")
			})

			if err != nil {
				fmt.Println(errors.Wrap(err, "poll"))
			}
		}
	}()

	return <-errCh
}
