package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"redka-debugger/gui"
	"redka-debugger/redkacore"
	"redka-debugger/utils"
	"strings"
)

func OldMain() {
	var database string
	var defaultDatabase = ""
	const databaseUsage = "database file path"
	flag.StringVar(&database, "database", defaultDatabase, databaseUsage)
	flag.StringVar(&database, "d", defaultDatabase, databaseUsage+" (shorthand)")

	var Keys string
	var defaultKeys = ""
	const keysUsage = "action: list keys (support wildcard by '*')"
	flag.StringVar(&Keys, "keys", defaultKeys, keysUsage)
	flag.StringVar(&Keys, "k", defaultKeys, keysUsage+" (shorthand)")

	var value string
	var defaultValue = ""
	const valueUsage = "action: 1. the key to get value when not write action, 2. the key (value is hash, key format is 'key#field') to write when write action"
	flag.StringVar(&value, "value", defaultValue, valueUsage)
	flag.StringVar(&value, "v", defaultValue, valueUsage+" (shorthand)")

	var write string
	var defaultWrite = ""
	const writeUsage = "action: the value to write"
	flag.StringVar(&write, "write", defaultWrite, writeUsage)
	flag.StringVar(&write, "w", defaultWrite, writeUsage+" (shorthand)")

	flag.Parse()

	if database == "" {
		utils.PrintUsage()
		os.Exit(1)
	}

	db := redkacore.LoadDB(database)
	defer db.Close()

	if Keys != "" {
		redkacore.RedkaKeys(db, Keys)
		return
	}

	if value != "" {
		if write != "" {
			keyField := strings.Split(value, "#")
			if len(keyField) == 2 {
				redkacore.RedkaSet(db, keyField[0], write, keyField[1])
			} else {
				redkacore.RedkaSet(db, keyField[0], write, "")
			}
			return
		}
		redkacore.RedkaGet(db, value)
		return
	}

	fmt.Println("No actions specified")
	utils.PrintUsage()
}

func main() {
	log.Println("reddka debugger")

	gui.TuiDemo2()
}
