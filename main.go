package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var cbPath string = "~/.cx/clipboard"

func cut(file string, cbFile string) {
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

	all := flag.Bool("a", false, "paste all clipboard files into current dir")
	flag.Parse()

	var files []string
	if *all {
		entries, err := os.ReadDir(cbPath)
		if err != nil {
			panic(err)
		}
		for _, e := range entries {
			files = append(files, e.Name())
		}
	} else {
		arg := flag.Args()
		files = arg[0:]
	}

	for _, file := range files {
		cbFile := cbFile(cbPath, file)

		if fileExists(file) {
			// cut
			cut(file, cbFile)
			continue
		}

		if fileExists(cbFile) {
			// paste
			paste(cbFile, file)
			continue
		}

		fmt.Printf("Error: %s not found in current dir or clipboard\n", file)
	}
}
