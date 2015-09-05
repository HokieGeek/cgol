package life

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func readLife106(file *os.File) (func(Dimensions, Location) []Location, error) {
	/*
	   #Life 1.06
	   0 -1
	   1 0
	   -1 1
	   0 1
	   1 1
	*/
	return nil, errors.New("Still not processing Life 1.06 format")
}

func readLife105(file *os.File) (func(Dimensions, Location) []Location, error) {
	/*
	   #Life 1.05
	   #D This is a glider.  [comments]
	   #N  [optional. see below]
	   #P -1 -1  [can't quite get this...]
	   .*.
	   ..*
	   ***


	   #N == #R 23/3

	   #R 125/36 ->  the pattern should be run in a universe where 1, 2, or 5 neighbors are necessary for a cell's survival, and 3 or 6 neighbours allow a cell to come alive.
	*/
	return nil, errors.New("Still not processing Life 1.05 format")
}

func Read(path string) (func(Dimensions, Location) []Location, error) {
	fmt.Printf("Read(%s)\n", path)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(f)
	header, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	// fmt.Printf("Header: '%s'\n", header)

	switch strings.ToLower(strings.TrimSpace(header)) {
	case "#life 1.05":
		return readLife105(f)
	case "#life 1.06":
		return readLife106(f)
	default:
		return nil, errors.New("Did not recognize file format")
	}
}

func Write(func(Dimensions, Location) []Location, string) error {
	return nil
}
