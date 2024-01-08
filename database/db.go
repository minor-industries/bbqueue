package database

import (
	"github.com/glebarez/sqlite"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type Measurement struct {
	gorm.Model
	ProbeID string
	Time    time.Time
	Temp    float64
}

func Get(filename string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "open")
	}

	err = db.AutoMigrate(&Measurement{})
	if err != nil {
		return nil, errors.Wrap(err, "migrate measurement")
	}

	return db, nil
}

func GetProbes(db *gorm.DB) ([]string, error) {
	var m []Measurement

	tx := db.Distinct("probe_id").Order("probe_id").Find(&m)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "find")
	}

	var result []string
	for _, r := range m {
		result = append(result, r.ProbeID)
	}

	return result, nil
}

func GetProbeData(
	db *gorm.DB,
	probeName string,
	after time.Time,
	before time.Time,
) ([]Measurement, error) {
	var m []Measurement

	tx := db.Where(
		"probe_id = ? and time >= ? and time <= ?",
		probeName,
		after,
		before,
	).Order("time asc").Find(&m)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "find")
	}

	return m, nil
}
