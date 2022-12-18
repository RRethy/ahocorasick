# Biblio

The fastest Golang implementation of the Aho-Corasick algorithm for string-searching.

## Usage

```bash
go get github.com/rrethy/biblio@v1.0.0
```

[Documentation](https://godoc.org/github.com/RRethy/biblio)

```go
m := CompileByteSlices([][]byte{
  []byte("he"),
  []byte("she"),
  []byte("his"),
  []byte("hers"),
  []byte("she"),
})
m.FindAllByteSlice([]byte("ushers")) // => { "she" 1 }, { "he" 2 }, { "hers" 2 }

m := CompileStrings([]string{
  "he",
  "she",
  "his",
  "hers",
  "she",
)
m.FindAllString("ushers") // => { "she" 1 }, { "he" 2 }, { "hers" 2 }
```

## Benchmarks

*macOS Mojave version 10.14.6*

*MacBook Pro (Retina, 13-inch, Early 2015)*

*Processor 3.1 GHz Intel Core i7*


```
$ git co d7354e5e7912add9c2c602aae74c508bca3b2f4d; go test -bench=Benchmark
```

The two basic operations are the compilation of the state machine from an array of patterns (`Compile`), and the usage of this state machine to find each pattern in text (`FindAll`). Other implementations call these operations under different names.

| Operation | Input Size | rrethy/biblio | [BobuSumisu/aho-corasick](https://github.com/BobuSumisu/aho-corasick) | [anknown/ahocorasick](https://github.com/anknown/ahocorasick) |
| - | - | - | - | - |
| - | - | Double-Array Trie | LinkedList Trie | Double-Array Trie |
| - | - | - | - | - |
| `Compile` | 235886 patterns | **133 ms** | 214 ms | 1408 ms |
| `Compile` | 23589 patterns  | **20 ms** | 50 ms  | 137 ms |
| `Compile` | 2359 patterns   | **3320 µs** | 11026 µs | 10506 µs |
| `Compile` | 236 patterns    | **229 µs**| 1377 µs| 867s µs |
| `Compile` | 24 patterns     | **43 µs**| 144 µs| 82s µs |
| - | - | - | - | - |
| `FindAll` | 3227439 bytes | **36 ms** | 38 ms | 116 ms |
| `FindAll` | 318647 bytes  | **3641 µs** | 3764 µs | 11335 µs |
| `FindAll` | 31626 bytes   | **359 µs** | 370 µs | 1103 µs |
| `FindAll` | 3657 bytes    | **31 µs** | 40 µs | 131 µs |

**NOTE**: `FindAll` uses a state machine compiled from 2359 patterns.

**NOTE**: `FindAll` time does **not** include the `Compile` time for the state machine.

### Reference Papers

[1] A. V. Aho, M. J. Corasick, "Efficient String Matching: An Aid to Bibliographic Search," Communications of the ACM, vol. 18, no. 6, pp. 333-340, June 1975.

[2] J.I. Aoe, "An Efficient Digital Search Algorithm by Using a Doble-Array Structure," IEEE Transactions on Software Engineering, vol. 15, no. 9, pp. 1066-1077, September 1989.

[3] J.I. Aoe, K. Morimoto, T. Sato, "An Efficient Implementation of Trie Stuctures," Software - Practice and Experience, vol. 22, no.9, pp. 695-721, September 1992.

## License

`MIT`
