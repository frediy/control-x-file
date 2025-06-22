package main

import (
	"fmt"
	"os"
	"strings"
)

var cbPath string = "~/.cx/clipboard"

func cut(file string, cbPath string) {
	cbFile := cbFile(cbPath, file)

	fmt.Printf("cut %s %s\n", file, abbreviateHomeDir(cbFile))

	err := os.Link(file, cbFile)
	if err != nil {
		panic(err)
	}

	err = os.Remove(file)
	if err != nil {
		panic(err)
	}
}

func paste(cbFile string, file string) {
	// get relative path
	fmt.Printf("paste %s %s\n", abbreviateHomeDir(cbFile), file)

	err := os.Link(cbFile, file)
	if err != nil {
		panic(err)
	}

	err = os.Remove(cbFile)
	if err != nil {
		panic(err)
	}
}

func cbFile(cbPath string, file string) string {
	return fmt.Sprintf("%s/%s", cbPath, file)
}

func expandHomeDir(path string) string {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		path = strings.Replace(path, "~", "", 1)
		return fmt.Sprintf("%s%s", home, path)
	}
	return path
}

func abbreviateHomeDir(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	if strings.HasPrefix(path, home) {
		path = strings.Replace(path, home, "~", 1)
		return path
	}
	return path
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func main() {
	cbPath = expandHomeDir(cbPath)

	arg := os.Args
	file := arg[1]

	if fileExists(file) {
		// cut
		cut(file, cbPath)
		return
	}

	cbFile := cbFile(cbPath, file)

	if fileExists(cbFile) {
		// paste
		paste(cbFile, file)
		return
	}

	fmt.Printf("Error: %s not found in current dir or clipboard", file)
}
