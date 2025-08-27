package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var cbPath string = "~/.cx/clipboard"

func cut(fromFile string, toFile string) {
	err := os.Link(fromFile, toFile)
	if err != nil {
		panic(err)
	}

	err = os.Remove(fromFile)
	if err != nil {
		panic(err)
	}
}

func clipboardFile(cbPath string, file string) string {
	return fmt.Sprintf("%s/%s", cbPath, file)
}

func relpathFromClipboardFile(cbPath, file string) string {
	return strings.Replace(file, cbPath+"/", "", 1)
}

func workdirFile(wdPath string, file string) string {
	return fmt.Sprintf("%s/%s", wdPath, file)
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

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func pathIsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	return fileInfo.IsDir()
}

func relativeFile(absFile string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	relFile, isRelative := strings.CutPrefix(absFile, wd+"/")
	if !isRelative {
		return "", fmt.Errorf("%s is not a relative file in current working dir %s", absFile, wd)
	}

	return relFile, nil
}

func getWdPath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return wd
}

func ensureDestinationPath(basePath string, file string) {
	// full path path without file
	pathComponents := strings.Split(file, "/")
	pathComponents = pathComponents[0 : len(pathComponents)-1]
	path := strings.Join(pathComponents, "/")
	fullPath := basePath + "/" + path

	// create if not exists
	_, err := os.Stat(fullPath)
	if err != nil {

		os.MkdirAll(fullPath, 0777)
	}
}

func main() {
	cbPath = expandHomeDir(cbPath)
	wdPath := getWdPath()

	all := flag.Bool("a", false, "paste all clipboard files into current dir")
	flag.Parse()

	var files []string
	if *all {
		// TODO: support pasting -a when subdirs exist on clipboard
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
		if strings.HasPrefix(file, "/") {
			relFile, err := relativeFile(file)
			file = relFile
			if err != nil {
				panic(err)
			}
		}

		cbFile := clipboardFile(cbPath, file)

		if pathExists(file) {
			// fmt.Println(pathExists(file))
			if pathIsDir(file) {
				// fmt.Println(pathIsDir(file))
				err := filepath.WalkDir(file, func(ipath string, d fs.DirEntry, err error) error {
					// fmt.Println(ipath, d.IsDir())
					if !d.IsDir() {
						cipath := clipboardFile(cbPath, ipath)
						// fmt.Printf("cut %s %s\n", ipath, abbreviateHomeDir(cipath))
						ensureDestinationPath(cbPath, ipath)
						cut(ipath, cipath)
					}
					return nil
				})
				if err != nil {
					panic(err)
				}
				// fmt.Println("os.Remove " + file)
				os.RemoveAll(file)
			} else {
				// cut
				// fmt.Printf("cut %s %s\n", file, abbreviateHomeDir(cbFile))
				ensureDestinationPath(cbPath, file)
				cut(file, cbFile)
			}
			continue
		}

		if pathExists(cbFile) {
			// fmt.Println(pathExists(cbFile))
			if pathIsDir(cbFile) {
				// fmt.Println(pathIsDir(cbFile))
				err := filepath.WalkDir(cbFile, func(cipath string, d fs.DirEntry, err error) error { // should iterate over cipath, not ipath
					// TODO: write method to convert clipboardFile to relativePath and call on cipath to get ipath
					// fmt.Println(cipath, d.IsDir())
					if !d.IsDir() {
						// fmt.Println(cipath)
						// cipath := clipboardFile(cbPath, ipath)
						rpath := relpathFromClipboardFile(cbPath, cipath)
						// fmt.Println("-")
						// fmt.Println(cipath)
						// fmt.Println(abbreviateHomeDir(cipath))
						// fmt.Println(rpath)
						// fmt.Printf("paste %s %s\n", abbreviateHomeDir(cipath), rpath)
						// fmt.Println("-")
						ensureDestinationPath(wdPath, rpath)
						cut(cipath, rpath)
					}
					return nil
				})
				if err != nil {
					panic(err)
				}
				// fmt.Println("os.Remove " + cbFile)
				os.RemoveAll(cbFile)
			} else {
				// cut
				// fmt.Printf("cut %s %s\n", file, abbreviateHomeDir(cbFile))
				ensureDestinationPath(wdPath, file)
				cut(cbFile, file)
			}
			continue
		}

		fmt.Printf("Error: %s not found in current dir or clipboard\n", file)
	}
}
