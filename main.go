package main

import (
	"debug/pe"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func readString(b []byte) (string, int) {
	sb := strings.Builder{}
	sb.Grow(len(b))

	offset := 0
	for {
		if b[offset] != 0 {
			sb.WriteByte(b[offset])
			offset += 2
		} else {
			break
		}
	}

	return sb.String(), offset + 2
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:", os.Args[0], "SETUP.EXE")
		return
	}

	exe, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer exe.Close()

	p, err := pe.NewFile(exe)
	if err != nil {
		panic(err)
	}

	endOfPE := int64(0)
	for _, section := range p.Sections {
		endOfSection := int64(section.Offset + section.Size)
		if endOfSection > endOfPE {
			endOfPE = endOfSection
		}
	}

	fmt.Printf("endOfPE = 0x%X\n", endOfPE)
	exe.Seek(endOfPE, io.SeekStart)

	packed, err := ioutil.ReadAll(exe)
	if err != nil {
		panic(err)
	}

	blockCount := int(binary.LittleEndian.Uint32(packed))
	fmt.Printf("blockCount = %d\n", blockCount)

	offset := 4
	for blockNum := 1; blockNum <= blockCount; blockNum++ {
		fmt.Printf("Block %d:\n", blockNum)

		shortName, strLen := readString(packed[offset:])
		fmt.Printf("  shortName = %q\n", shortName)
		offset += strLen

		fullName, strLen := readString(packed[offset:])
		fmt.Printf("  fullName = %q\n", fullName)
		offset += strLen

		versionStr, strLen := readString(packed[offset:])
		fmt.Printf("  versionStr = %q\n", versionStr)
		offset += strLen

		lengthStr, strLen := readString(packed[offset:])
		fmt.Printf("  lengthStr = %q\n", lengthStr)
		offset += strLen

		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(shortName, packed[offset:offset+length], 0644)
		if err != nil {
			panic(err)
		}

		offset += length
	}
}
