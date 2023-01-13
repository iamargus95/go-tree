package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

type Options struct {
	Indent    string
	OutputBuf *bytes.Buffer
	ShowFiles bool
}

func (o *Options) tree(path string, indent string, isLast bool) {
	file, err := os.Open(path)
	if err != nil {
		o.OutputBuf.WriteString(err.Error() + "\n")
		return
	}
	defer file.Close()

	files, err := file.Readdir(-1)
	if err != nil {
		o.OutputBuf.WriteString(err.Error() + "\n")
		return
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for i, file := range files {
		if file.IsDir() {
			if i == len(files)-1 {
				o.OutputBuf.WriteString(indent + "└── " + file.Name() + "\n")
				o.tree(filepath.Join(path, file.Name()), indent+o.Indent+"   ", true)
			} else {
				o.OutputBuf.WriteString(indent + "├── " + file.Name() + "\n")
				o.tree(filepath.Join(path, file.Name()), indent+"│   "+o.Indent, false)
			}
		} else {
			if i == len(files)-1 {
				o.OutputBuf.WriteString(indent + "└── " + file.Name() + "\n")
			} else {
				o.OutputBuf.WriteString(indent + "├── " + file.Name() + "\n")
			}
		}
	}
}

func (o *Options) parseByFlags(path, indent string, file fs.FileInfo) {
	if o.ShowFiles {
		o.tree(filepath.Join(path, file.Name()), indent+o.Indent+"   ", true)
	}
}

func (o *Options) printTree(path string) {
	o.OutputBuf.WriteString(path + "\n")
	o.tree(path, o.Indent, false)
	o.OutputBuf.WriteString(
		strconv.Itoa(
			o.countDirs(path)-1,
		) + " directories , " + strconv.Itoa(
			o.countFiles(path),
		) + " files",
	)
}

var showFiles bool

func init() {
	flag.BoolVar(&showFiles, "f", false, "show full path")
	flag.Parse()
}

func main() {
	buf := &bytes.Buffer{}
	options := &Options{Indent: "    ", OutputBuf: buf, ShowFiles: showFiles}
	options.printTree(flag.Arg(0))
	fmt.Println(buf.String())
}

func (o *Options) countDirs(path string) int {
	count := 0
	filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if info.IsDir() {
			count++
		}
		return nil
	})
	return count
}

func (o *Options) countFiles(path string) int {
	count := 0
	filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			count++
		}
		return nil
	})
	return count
}
