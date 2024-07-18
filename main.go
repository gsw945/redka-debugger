package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/nalgeon/redka"
	_ "modernc.org/sqlite"
)

func printUsage() {
	path, err := os.Executable()
	if err != nil {
		panic(err)
	}
	Executable := filepath.Base(path)
	fmt.Println("Usage: " + Executable + " <Options>")
	fmt.Println("Options:")
	flag.PrintDefaults()
}

func loadDB(dbPath string) *redka.DB {
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

func redkaKeys(db *redka.DB, pattern string) {
	keys, err := db.Key().Keys(pattern)
	if err != nil {
		panic(err)
	}
	for _, key := range keys {
		fmt.Printf("key(ID=%d, Key=%s, Type=%v, Version=%d, ETime=%v, MTime=%v)\n", key.ID, key.Key, redkaTypeName(int(key.Type)), key.Version, key.ETime, key.MTime)
	}
}

func redkaTypeName(t int) string {
	switch t {
	case 0:
		return "Any(0)"
	case 1:
		return "String(1)"
	case 2:
		return "List(2)"
	case 3:
		return "Set(3)"
	case 4:
		return "Hash(4)"
	case 5:
		return "ZSet(5)"
	default:
		return "<Unknown>"
	}
}

func redkaGet(db *redka.DB, key string) {
	k, err := db.Key().Get(key)
	if err != nil {
		panic(err)
	}
	intType := int(k.Type)
	typeName := redkaTypeName(intType)
	fmt.Printf("key(ID=%d, Key=%s, Type=%v, Version=%d, ETime=%v, MTime=%v)\n", k.ID, k.Key, typeName, k.Version, k.ETime, k.MTime)
	switch intType {
	case 0:
		fmt.Println("not supported")
	case 1:
		valueString, err := db.Str().Get(k.Key)
		if err != nil {
			panic(err)
		}
		fmt.Printf("value: %s\n", valueString)
	case 2:
		valueList, err := db.List().Range(k.Key, 0, -1)
		if err != nil {
			panic(err)
		}
		for i, v := range valueList {
			fmt.Printf("value[%d]: %v\n", i, v)
		}
	case 3:
		valueSet, err := db.Set().Items(k.Key)
		if err != nil {
			panic(err)
		}
		for i, v := range valueSet {
			fmt.Printf("value[%d]: %v\n", i, v)
		}
	case 4:
		valueHash, err := db.Hash().Items(k.Key)
		if err != nil {
			panic(err)
		}
		for k, v := range valueHash {
			fmt.Printf("value[%s]: %v\n", k, v)
		}
	case 5:
		valueZSet, err := db.ZSet().Range(k.Key, 0, -1)
		if err != nil {
			panic(err)
		}
		for i, v := range valueZSet {
			fmt.Printf("value[%d]: %v\n", i, v)
		}
	default:
		fmt.Println("not supported")
	}
}

func main() {
	log.Println("reddka debugger")

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
	const valueUsage = "action: get value by key"
	flag.StringVar(&value, "value", defaultValue, valueUsage)
	flag.StringVar(&value, "v", defaultValue, valueUsage+" (shorthand)")

	flag.Parse()

	if database == "" {
		printUsage()
		os.Exit(1)
	}

	db := loadDB(database)
	defer db.Close()

	if Keys != "" {
		redkaKeys(db, Keys)
		return
	}

	if value != "" {
		redkaGet(db, value)
		return
	}

	fmt.Println("No actions specified")
	printUsage()
}
