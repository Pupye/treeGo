package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

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

func dirTree(out io.Writer, path string, printFiles bool) (err error) {
	if path == "." {
		path = ""
	}
	var workingPath string
	var tree string

	if workingPath, err = os.Getwd(); err != nil {
		return err
	}
	path = workingPath + "/" + path
	if !printFiles {
		if tree, err = constructTree(path, ""); err != nil {
			return err
		}
	} else {
		if tree, err = constructTreeFullTree(path, ""); err != nil {
			return err
		}
	}
	fmt.Fprint(out, tree)
	return
}

func constructTree(path string, tab string) (result string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	names, err := f.Readdir(-1)
	if err != nil {
		return
	}
	dirNames := filter(names)
	sort.Strings(dirNames)
	for i, name := range dirNames {
		var subtree string
		if i == len(dirNames)-1 {
			if subtree, err = constructTree(path+"/"+name, tab+"\t"); err != nil {
				return
			}
			result += tab + "└───" + name + "\n" + subtree
		} else {
			if subtree, err = constructTree(path+"/"+name, tab+"│\t"); err != nil {
				return
			}
			result += tab + "├───" + name + "\n" + subtree
		}
	}
	return
}

func constructTreeFullTree(path string, tab string) (result string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	names, err := f.Readdir(-1)
	if err != nil {
		return
	}
	sort.Slice(names, func(i, j int) bool {
		return names[i].Name() < names[j].Name()
	})
	for i, name := range names {
		var subtree string
		if name.IsDir() {
			if i == len(names)-1 {
				if subtree, err = constructTreeFullTree(path+"/"+name.Name(), tab+"\t"); err != nil {
					return
				}
				result += tab + "└───" + name.Name() + "\n" + subtree
			} else {
				if subtree, err = constructTreeFullTree(path+"/"+name.Name(), tab+"│\t"); err != nil {
					return
				}
				result += tab + "├───" + name.Name() + "\n" + subtree
			}
		} else {
			size := getFileSize(name)
			if i == len(names)-1 {
				result += tab + "└───" + name.Name() + " " + size + "\n"
			} else {
				result += tab + "├───" + name.Name() + " " + size + "\n"
			}
		}
	}
	return
}

func getFileSize(info os.FileInfo) (result string) {
	if info.Size() == 0 {
		result = "(empty)"
	} else {
		result = fmt.Sprintf("(%db)", info.Size())
	}
	return
}

func filter(files []os.FileInfo) (result []string) {
	for _, file := range files {
		if file.IsDir() {
			result = append(result, file.Name())
		}
	}
	return result
}
