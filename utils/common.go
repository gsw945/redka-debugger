package utils

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func PrintUsage() {
	path, err := os.Executable()
	if err != nil {
		panic(err)
	}
	Executable := filepath.Base(path)
	fmt.Println("Usage: " + Executable + " <Options>")
	fmt.Println("Options:")
	flag.PrintDefaults()
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
