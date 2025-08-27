package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"
)

var cbPath string = "~/.cx/clipboard"

func cut(fromFile string, toFile string, keepFromFile bool) {
	err := os.Link(fromFile, toFile)
	if err != nil {
		panic(err)
	}

	if !keepFromFile {
		err = os.Remove(fromFile)
		if err != nil {
			panic(err)
		}
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
func cutToClipboard(cbFile string, file string, keepSource bool) {
	if pathIsDir(file) {
		err := filepath.WalkDir(file, func(ipath string, d fs.DirEntry, err error) error {
			if !d.IsDir() {
				cipath := clipboardFile(cbPath, ipath)
				ensureDestinationPath(cbPath, ipath)
				cut(ipath, cipath, keepSource)
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
		if !keepSource {
			os.RemoveAll(file)
		}
	} else {
		ensureDestinationPath(cbPath, file)
		cut(file, cbFile, keepSource)
	}
}

func pasteFromClipboard(cbFile string, wdPath string, file string, keepSource bool) {
	if pathIsDir(cbFile) {
		err := filepath.WalkDir(cbFile, func(cipath string, d fs.DirEntry, err error) error { // should iterate over cipath, not ipath
			if !d.IsDir() {
				rpath := relpathFromClipboardFile(cbPath, cipath)
				ensureDestinationPath(wdPath, rpath)
				cut(cipath, rpath, keepSource)
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
		if !keepSource {
			os.RemoveAll(cbFile)
		}

	} else {
		ensureDestinationPath(wdPath, file)
		cut(cbFile, file, keepSource)
	}
}

func main() {
	cbPath = expandHomeDir(cbPath)
	wdPath := getWdPath()

	all := flag.BoolP("all", "a", false, "paste all clipboard paths into current dir")
	keep := flag.BoolP("keep", "k", false, "keep paths in workdir after cut and in clipboard after paste")
	help := flag.BoolP("help", "h", false, "show help")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [<options>] [<path>...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n  Detectes whether [<path>...] is in current dir or clipboard.")
		fmt.Fprintf(os.Stderr, "\n  Cuts [<path>...] in current dir recursively to clipboard.")
		fmt.Fprintf(os.Stderr, "\n  Pastes [<path>...] recursively from clipboard to current dir.\n")
		fmt.Fprintln(os.Stderr, "\noptions:")
		flag.PrintDefaults()
	}

	flag.Parse()

	// TODO: ensure clipboard path exists
	// TODO: support overwrite clashes between clipboard and wd link fails with panic
	// TODO: add clipoard -d to delete all clipboard contents
	// TODO: add clipoard -l to list files

	// TODO: add clipoard -p for enforce paste
	// TODO: add clipoard -c for enforce cut

	var files []string

	arg := flag.Args()
	files = arg[0:]

	if (len(files) == 0 && !*all) || *help {
		flag.Usage()
		os.Exit(0)
	}

	keepSource := *keep

	if *all {
		entries, err := os.ReadDir(cbPath)

		if err != nil {
			panic(err)
		}
		for _, e := range entries {
			file := e.Name()
			cbFile := clipboardFile(cbPath, file)
			pasteFromClipboard(cbFile, wdPath, file, keepSource)
		}
		os.Exit(0)
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
			cutToClipboard(cbFile, file, keepSource)
			continue
		}

		if pathExists(cbFile) {
			pasteFromClipboard(cbFile, wdPath, file, keepSource)
			continue
		}

		fmt.Printf("Error: %s not found in current dir or clipboard\n", file)
	}
}
