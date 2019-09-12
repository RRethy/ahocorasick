package main

import (
	"fmt"
	"github.com/rrethy/biblio"
)

func main() {
	patterns := [][]byte{}
	text := []byte("potholderz by MF DOOM hot shit aw shit 锅 持有人")

	m := biblio.Compile(patterns)
	m.FindAll(text)

	// for _, match := range matches {
	// 	fmt.Printf("%d - %s\n", match.Index, string(match.Word))
	// }
	fmt.Printf("%v\n", m.Base)
	fmt.Printf("%v\n", m.Check)
	fmt.Printf("%v\n", m.Fail)
}
