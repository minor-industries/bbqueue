package database

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
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
