# biblio

Golang implementation of the aho-corasick string matching algorithm with full UTF-8 support.

## Usage

```go
import "github.com/rrethy/biblio"

// compile the dictionary of words
bib := biblio.Compile("mm", "food", "doom")

// find all matches in the string
matches := bib.FindAll("MM..Food MF DOOM - mm..food mf doom")
// => { "mm" 19 }, { "food" 23 }, { "doom" 31 }
```

## Install

```
go get -u github.com/rrethy/biblio
```
