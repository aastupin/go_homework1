package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

type FileDirStr struct {
	name  string
	isDir bool
	size  int64
	dirs  []FileDirStr
}

func (f FileDirStr) Sort() {
	sort.Slice(f.dirs, func(i, j int) bool {
		return f.dirs[i].name < f.dirs[j].name
	})
}

const newLevelSymbol = "├───"
const tabLevelSymbol = "│	"
const lastElementSymbol = "└───"

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(output io.Writer, path string, printFiles bool) error {
	abs, err := filepath.Abs(path)
	if err == nil {
		if path != abs+string(os.PathSeparator) {
			path = abs + string(os.PathSeparator)
		}
	}
	lastFolderName := filepath.Base(path)
	fileStr, err := getFileStruct(path, lastFolderName, printFiles)

	if err != nil {
		return err
	}

	PrintData(output, fileStr, 0, len(fileStr.dirs) == 0, printFiles, 0, "")
	return nil
}

func getFileStruct(path string, dirname string, printFiles bool) (FileDirStr, error) {
	resultFileDir := FileDirStr{name: dirname, isDir: true}

	files, err := os.ReadDir(path)
	if err != nil {
		return resultFileDir, err
	}

	for _, file := range files {
		subDirFilePath := path + string(os.PathSeparator) + file.Name()
		if file.IsDir() {
			subDir, subDirErr := getFileStruct(subDirFilePath, file.Name(), printFiles)
			if subDirErr != nil {
				return resultFileDir, subDirErr
			}

			resultFileDir.dirs = append(resultFileDir.dirs, subDir)
		} else {
			if printFiles {
				fileStat, err := os.Stat(subDirFilePath)
				if err != nil {
					return resultFileDir, err
				}
				resultFileDir.dirs = append(resultFileDir.dirs, FileDirStr{name: file.Name(), isDir: false, size: fileStat.Size()})
			}
		}
	}
	return resultFileDir, nil
}

func PrintData(output io.Writer, data FileDirStr, level int, lastDir bool, printFiles bool, lastDirCount int, prefixStr string) {
	if level > 0 {
		var symbolDraw string = newLevelSymbol
		if lastDir {
			symbolDraw = lastElementSymbol
		}
		var fileSizeStr string
		if !data.isDir {
			if data.size == 0 {
				fileSizeStr = " (empty)"
			} else {
				fileSizeStr = " (" + strconv.Itoa(int(data.size)) + "b)"
			}
		}
		fmt.Fprintln(output, prefixStr+symbolDraw+data.name+fileSizeStr)
	}
	if lastDir {
		lastDirCount = lastDirCount + 1
	}
	if len(data.dirs) > 0 {
		data.Sort()
		var newPrefix string = ""
		if level > 0 {
			if lastDir {
				newPrefix = "\t"
			} else {
				newPrefix = tabLevelSymbol
			}
		}
		for indx, dir := range data.dirs {
			if dir.isDir || printFiles {
				lastDirFlag := indx == len(data.dirs)-1
				PrintData(output, dir, level+1, lastDirFlag, printFiles, lastDirCount, prefixStr+newPrefix)
			}
		}
	}
}
