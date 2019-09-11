package main

import (
	"fmt"
	"github.com/rrethy/biblio"
)

func main() {
	// g 103
	// o 111
	// d 100
	// e 101
	// -102 -109 - - - - - -108
	// 0    1    2 10 - - - - - 1
	patterns := [][]byte{[]byte("go"), []byte("goo"), []byte("good"), []byte("oode")}
	// output := []biblio.Match{{[]byte("the"), 2}, {[]byte("they"), 3}, {[]byte("theyre"), 5}, {[]byte("go"), 14}, {[]byte("goo"), 15}, {[]byte("good"), 16}, {[]byte("oode"), 17}, {[]byte("te"), 20}, {[]byte("tea"), 21}, {[]byte("team"), 22}}
	text := []byte("theyre not a goode team")

	m := biblio.Compile(patterns)
	matches := m.FindAll(text)

	for _, match := range matches {
		fmt.Printf("%d - %s\n", match.Index, string(match.Word))
	}
	fmt.Printf("%v\n", m.Base)
	fmt.Printf("%v\n", m.Check)
	fmt.Printf("%v\n", m.Fail)
	for state, words := range m.Output {
		fmt.Println(state)
		for _, word := range words {
			fmt.Printf("%s\n", word)
		}
	}
}
