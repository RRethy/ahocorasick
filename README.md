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

## TODO

* add benchmarks comparing biblio against other implementations such as `anknown/ahocorasick` and `cloudflare/ahocorasick`. The anknown implementation is faster at creating the state machine but much slower when parsing. The cloudflare implementation is faster at parsing, and can be either slower or faster at creating the state machine depending on input. Overall, the anknown implementation is fastest if it is simply a compile+parse which is shown in the benchmarks on his README, but it is slower the more parses occur.
* try using byte data type instead of rune to improve efficiency.
* reduce the amount of copying done when compiling the state machine. It has a noticeable penalty on massive input dictionary with many unique characters.
