package db

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_do(t *testing.T) {
	_, err := Get("sqlite.db")
	require.NoError(t, err)

}
