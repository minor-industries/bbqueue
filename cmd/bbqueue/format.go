package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"time"
)

var templateFuncs = map[string]any{
	"FmtTemp": func(temp float64) string {
		return fmt.Sprintf("%0.1f", temp*9/5+32)
	},
	"LastUpdated": func(t, now time.Time) string {
		return humanize.RelTime(t, now, "ago", "from now")
	},
}
