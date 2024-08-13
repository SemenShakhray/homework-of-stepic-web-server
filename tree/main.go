package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

func dirTree(out io.Writer, path string, printFiles bool) error {
	err := getDirTree(out, "", path, printFiles)
	if err != nil {
		return err
	}
	return nil
}

func getDirTree(out io.Writer, prefix, path string, printFiles bool) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("error read directory - %v", err)
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })
	var newFiles []fs.DirEntry
	for _, file := range files {
		if file.IsDir() {
			newFiles = append(newFiles, file)
		}
	}
	var graf, newPrefix string
	if printFiles {
		for i, file := range files {
			info, err := file.Info()
			if err != nil {
				return fmt.Errorf("error get info - %v", err)
			}
			if i == len(files)-1 {
				graf = "└───"
				newPrefix = prefix + "\t"
			} else {
				graf = "├───"
				newPrefix = prefix + "│\t"
			}
			if file.IsDir() {
				fmt.Fprintf(out, "%s%s%s\n", prefix, graf, file.Name())
				err = getDirTree(out, newPrefix, filepath.Join(path, file.Name()), printFiles)
				if err != nil {
					return fmt.Errorf("err - %v", err)
				}
			} else {
				if info.Size() != 0 {
					fmt.Fprintf(out, "%s%s%s (%db)\n", prefix, graf, file.Name(), info.Size())
				} else {
					fmt.Fprintf(out, "%s%s%s (empty)\n", prefix, graf, file.Name())
				}
			}

		}
	} else {
		for i, file := range newFiles {
			if i == len(newFiles)-1 {
				graf = "└───"
				newPrefix = prefix + "\t"
			} else {
				graf = "├───"
				newPrefix = prefix + "│\t"
			}
			if file.IsDir() {
				fmt.Fprintf(out, "%s%s%s\n", prefix, graf, file.Name())
				err = getDirTree(out, newPrefix, filepath.Join(path, file.Name()), printFiles)
				if err != nil {
					return fmt.Errorf("err - %v", err)
				}
			}
		}
	}
	return nil
}

func main() {
	// file:=os.Create("dire")
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
