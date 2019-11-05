package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
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

var bStart, vert, hor, bEnd = "├", "│", "───", "└"

func dirTree(out io.Writer, path string, printFiles bool) error {
	return dirTreeFormat(out, path, printFiles, "")
}

func dirTreeFormat(out io.Writer, path string, printFiles bool, prepend string) error {
	infos, e := ioutil.ReadDir(path)
	if e != nil {
		return (e)
	}

	if !printFiles {
		infos = filter(infos, func(f os.FileInfo) bool { return f.IsDir() })
	}

	for idx, info := range infos {
		isLast := idx == len(infos)-1
		if info.IsDir() || printFiles {
			fmt.Fprintln(out, formatName(info, prepend, isLast))
		}

		if info.IsDir() {
			var p string
			if isLast {
				p = ""
			} else {
				p = vert
			}

			dirTreeFormat(out, fmt.Sprintf("%s%c%s", path, os.PathSeparator, info.Name()), printFiles,
				fmt.Sprintf("%s%s\t", prepend, p))
		}
	}
	return nil
}

func filter(vs []os.FileInfo, f func(os.FileInfo) bool) []os.FileInfo {
	vsf := make([]os.FileInfo, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func formatName(fileInfo os.FileInfo, prepend string, isLast bool) string {
	var startChar string
	if isLast {
		startChar = bEnd
	} else {
		startChar = bStart
	}
	res := prepend + startChar + hor + fileInfo.Name()
	if !fileInfo.IsDir() {
		size := fileInfo.Size()
		if size == 0 {
			res += " (empty)"
		} else {
			res += fmt.Sprintf(" (%vb)", size)
		}
	}
	return res
}
