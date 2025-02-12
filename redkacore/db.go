package redkacore

import (
	"path/filepath"

	"github.com/nalgeon/redka"
	_ "modernc.org/sqlite"
)

func LoadDB(dbPath string) *redka.DB {
	absPath, err := filepath.Abs(dbPath)
	if err != nil {
		panic(err)
	}
	opts := redka.Options{
		DriverName: "sqlite",
	}
	db, err := redka.Open(absPath, &opts)
	if err != nil {
		panic(err)
	}
	return db
}
