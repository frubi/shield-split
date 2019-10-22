package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:", os.Args[0], "SETUP.EXE")
		return
	}

	exe, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	magic := []byte{'I', 'S', 'c', '(', '`', 0x09}
	tbl := make([]int, 0)

	for offset := 0; offset < len(exe)-len(magic); offset++ {
		if bytes.HasPrefix(exe[offset:], magic) {
			fmt.Printf("found file at %d (0x%X)\n", offset, offset)
			tbl = append(tbl, offset)
		}
	}

	if len(tbl) != 3 {
		fmt.Println("Invalid number of files:", len(tbl))
		return
	}

	files := []struct {
		name  string
		start int
		end   int
	}{
		{"data1.hdr", tbl[2], len(exe)},
		{"data1.cab", tbl[0], tbl[1]},
		{"data2.cab", tbl[1], tbl[2]},
	}

	for _, file := range files {
		fmt.Println("Writing", file.name)
		err = ioutil.WriteFile(file.name, exe[file.start:file.end], 0644)
		if err != nil {
			panic(err)
		}
	}
}
