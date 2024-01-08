package main

import (
	"github.com/minor-industries/bbqueue/database"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_plotIt(t *testing.T) {
	db, err := database.Get("sqlite3.db")
	require.NoError(t, err)

	svg, err := plotIt(db, []string{"bbq01-bbq", "bbq01-meat"})
	require.NoError(t, err)

	err = os.WriteFile("plot.svg", svg, 0o640)
	require.NoError(t, err)
}
