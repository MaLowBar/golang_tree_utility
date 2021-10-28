package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)


func printDir(out io.Writer, path, prefix string, needFiles, isPrevLast bool, level *int, seps map[int]bool) error {
	rawFiles, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	*level += 1
	var files []os.DirEntry

	if !needFiles {
		for _, file := range rawFiles {
			if file.IsDir() {
				files = append(files, file)
			}
		}
	} else {
		files = rawFiles
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
	if !isPrevLast && *level != 1 {
		seps[*level-1] = true
	}
	prefix = strings.Replace(prefix, "│", " ", -1)
	temp := []rune(prefix)
	for key := range seps {
		temp[2*(key-1)] = '│'
	}
	prefix = string(temp)

	for idx, file := range files {
		var isLast = false
		if idx == len(files)-1 {
			prefix = strings.Replace(prefix, `├`, `└`, 1)
			isLast = true
			delete(seps, *level)
		} else {
			prefix = strings.Replace(prefix, `└`, `├`, 1)
		}

		if file.IsDir() {
			strOut := fmt.Sprintf("%v%v\n", prefix, file.Name())
			strOut = strings.Replace(strOut, " ", "", -1)
			out.Write([]byte(strOut))
			printDir(out, path+"/"+file.Name(), " \t"+prefix, needFiles, isLast, level, seps)
		} else {
			inf, err := file.Info()
			if err != nil {
				log.Fatal(err)
			}
			fileSize := `empty`
			if inf.Size() > 0 {
				fileSize = strconv.Itoa(int(inf.Size())) + `b`
			}
			strOut := fmt.Sprintf("%v%v (%v)\n", prefix, inf.Name(), fileSize)
			spacesCount := strings.Count(strOut, " ") - 1
			strOut = strings.Replace(strOut, " ", "", spacesCount)
			out.Write([]byte(strOut))
		}
	}
	*level -= 1
	return nil
}

func dirTree(out io.Writer, path string, printFilse bool) error {
	lvl := new(int)
	separs := make(map[int]bool)
	return printDir(out, path, `├───`, printFilse, false, lvl, separs)
}

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
