package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minor-industries/bbqueue/html"
	"github.com/minor-industries/bbqueue/radio"
	"github.com/minor-industries/rtgraph"
	"github.com/minor-industries/rtgraph/database/sqlite"
	"github.com/minor-industries/rtgraph/prom"
	"github.com/minor-industries/rtgraph/schema"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"strings"
	"time"
)

func Run() error {
	errCh := make(chan error)

	db, err := sqlite.Get(os.ExpandEnv("$HOME/bbqueue.db"))
	if err != nil {
		return errors.Wrap(err, "get database")
	}

	go db.RunWriter(errCh)

	graph, err := rtgraph.New(
		db,
		errCh,
		rtgraph.Opts{
			ExternalMetrics: prom.PublishPrometheusMetrics,
		},
	)
	if err != nil {
		return errors.Wrap(err, "new rtgraph")
	}

	router := gin.New()

	graph.SetupServer(router.Group("/rtgraph"))

	router.GET("/main.html", func(c *gin.Context) {
		c.FileFromFS("main.html", http.FS(html.FS))
	})

	go func() {
		gin.SetMode(gin.ReleaseMode)
		errCh <- router.Run("0.0.0.0:8076")
	}()

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
						return nil
					}
				}

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
