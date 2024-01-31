package database

import (
	"fmt"
	"github.com/dustin/go-humanize"
	_ "github.com/dustin/go-humanize"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_do(t *testing.T) {
	db, err := Get("../sqlite3.db")
	require.NoError(t, err)

	var results []Measurement

	tx := db.Distinct("probe_id").Find(&results)
	require.NoError(t, tx.Error)

	for _, r := range results {
		fmt.Println(r.ProbeID)
	}
}

func Test_MostRecent(t *testing.T) {
	db, err := Get("../sqlite3.db")
	require.NoError(t, err)

	M, err := GetLatestTemps(db)
	require.NoError(t, err)

	now := time.Now()
	for _, m := range M {
		fmt.Println(m.Time, m.Temp, humanize.RelTime(m.Time, now, "ago", "from now"))
	}
}
