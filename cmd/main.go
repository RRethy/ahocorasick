package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/rrethy/biblio"
	"io/ioutil"
	"os"
)

var (
	patsFname = flag.String("patternsFile", "", "file to read patterns from")
	fname     = flag.String("file", "", "file to parse")
)

func main() {
	flag.Parse()

	if len(*patsFname) == 0 || len(*fname) == 0 {
		return
	}

	patsFile, err := os.Open(*patsFname)
	defer patsFile.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	var patterns [][]byte
	scanner := bufio.NewScanner(patsFile)
	for scanner.Scan() {
		patterns = append(patterns, scanner.Bytes())
	}

	fileContents, err := ioutil.ReadFile(*fname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	m := biblio.CompileByteSlices(patterns)
	matches := m.FindAllByteSlice(fileContents)
	fmt.Println(len(matches))
}
