package database

import (
	"fmt"
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

	fmt.Print(db)

	return db, nil
}
