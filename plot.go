package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/minor-industries/bbqueue/database"
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
	probeNames, err := database.GetProbes(db)
	if err != nil {
		_ = c.Error(err)
		return
	}

	_, err = plotIt(db, probeNames)
	if err != nil {
		_ = c.Error(err)
		return
	}

	//svg, err := plotIt(data)
	//if err != nil {
	//	panic(err)
	//}
	//c.Data(200, "image/svg+xml", svg)
}

func plotIt(db *gorm.DB, probeNames []string) ([]byte, error) {
	p := plot.New()

	p.Title.Text = "Temperature"
	p.X.Label.Text = "t"
	p.Y.Label.Text = "Temp"

	xticks := plot.TimeTicks{Format: "2006-01-02\n15:04"}
	p.X.Tick.Marker = xticks

	var vs []interface{}

	for _, probeName := range probeNames {
		vs = append(vs, probeName)
		pts := plotter.XYs{}

		data, err := database.GetProbeData(db, probeName)
		if err != nil {
			return nil, errors.Wrap(err, "get probe data")
		}

		//t0 := data[0].Time
		for _, d := range data {
			pts = append(pts, plotter.XY{
				X: float64(d.Time.Unix()),
				Y: d.Temp,
			})
		}

		vs = append(vs, pts)
	}

	err := plotutil.AddLines(p, vs...)
	if err != nil {
		return nil, errors.Wrap(err, "add lines")
	}

	w, err := p.WriterTo(8*vg.Inch, 4*vg.Inch, "svg")
	buf := bytes.NewBuffer(nil)
	_, err = w.WriteTo(buf)
	if err != nil {
		return nil, errors.Wrap(err, "write to")
	}

	return buf.Bytes(), nil
}