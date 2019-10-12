package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/rrethy/biblio"
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
		patterns = append(patterns, []byte(scanner.Text()))
	}
	m := biblio.CompileByteSlices(patterns)

	linesFile, err := os.Open(*fname)
	defer linesFile.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	scanner = bufio.NewScanner(linesFile)
	matchedLines := 0
	for scanner.Scan() {
		matches := m.FindAllByteSlice(scanner.Bytes())
		if len(matches) > 0 {
			// for _, match := range matches {
			// fmt.Println(match)
			// }
			matchedLines++
		}
	}
	fmt.Println(matchedLines)
}
