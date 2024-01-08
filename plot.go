package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gorm.io/gorm"
)

func plotHandler(db *gorm.DB, c *gin.Context) {
	//months := getMonthsParam(c, 3)
	//
	//data, err := getData(dbmap, months)
	//if err != nil {
	//	_ = c.Error(err)
	//	return
	//}

	svg, err := plotIt(data)
	if err != nil {
		panic(err)
	}
	c.Data(200, "image/svg+xml", svg)
}

func plotIt() {
	p := plot.New()

	p.Title.Text = "Temperature"
	p.X.Label.Text = "t"
	p.Y.Label.Text = "Temp"

	var pts plotter.XYs

	xticks := plot.TimeTicks{Format: "2006-01-02\n15:04"}

	for _, w := range data {
		pts = append(pts, plotter.XY{
			X: float64(w.T.Unix()),
			Y: w.Weight,
		})
	}

	p.X.Tick.Marker = xticks

	err := plotutil.AddLinePoints(p,
		"weight (kg)", pts,
	)
	if err != nil {
		return nil, errors.Wrap(err, "add line points")
	}

	w, err := p.WriterTo(8*vg.Inch, 4*vg.Inch, "svg")
	buf := bytes.NewBuffer(nil)
	_, err = w.WriteTo(buf)
	if err != nil {
		return nil, errors.Wrap(err, "write to")
	}

	return buf.Bytes(), nil
}
