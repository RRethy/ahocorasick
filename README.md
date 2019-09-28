# Biblio

:zap: Lightning fast implementation of the Aho-Corasick string matching algorithm using a Double Array Trie. All characters from all languages are supported since matching is done on a byte level.

## Usage

**TODO**: Link to GoDoc

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
$ go test -bench=Benchmark
```

The two basic operations are the compilation of the state machine from an array of patterns (`Compile`), and the usage of this state machine to find each pattern in text (`FindAll`). Other implementations call these operations under different names.

| Operation | Input Size | rrethy/biblio | [BobuSumisu/aho-corasick](https://github.com/BobuSumisu/aho-corasick) | [anknown/ahocorasick](https://github.com/anknown/ahocorasick) |
| - | - | - | - | - |
| `Compile` | 235886 patterns | **204 ms**    | 213 ms    | 1391 ms |
| `Compile` | 23589 patterns  |  **44 ms**    |  50 ms    |  136 ms |
| `Compile` | 2359 patterns   |   **8500 µs** |  11441 µs |   10577 µs |
| `Compile` | 236 patterns    |   1050 µs |   1514 µs |     **860 µs** |
| `Compile` | 24 patterns     |    148 µs |    139 µs |      **82 µs** |
| - | - | - | - | - |
| `FindAll` | 3227439 bytes | **35 ms**   | 37 ms   | 117 ms |
| `FindAll` | 318647 bytes  |  **3543 µs** |  3767 µs |  11367 µs |
| `FindAll` | 31626 bytes   |   **354 µs** |   373 µs |   1109 µs |
| `FindAll` | 3657 bytes    |    **32 µs** |    40 µs |    128 µs |

**NOTE**: `FindAll` uses a state machine compiled from 2359 patterns.

**NOTE**: `FindAll` time does **not** include the `Compile` time for the state machine.

### Other Implementations

Two implementations were intentionally omitted, [cloudflare/ahocorasick](https://github.com/cloudflare/ahocorasick) and [iohub/ahocorasick](https://github.com/iohub/ahocorasick). There are existing benchmarks [here](https://github.com/BobuSumisu/aho-corasick) and [here](https://github.com/anknown/ahocorasick) which have these implementations in their benchmark comparisons.

The *cloudflare* implementation was omitted due to incorrectness, it does not find all instances of the compiled patterns in the text. Even still, it has a significantly slower (~10x) `Compile` time and a competitive `FindAll` time (partly due to it reporting far less results than actually exist).

The *iohub* implementation was omitted since I was unable to get it working as it seemed to time out during compilation. As well, it has an unnecessarily confusing API.

## Implementation Details

**TODO**: Probably going to write a blog post. There are a lot of suddle things I did to improve the construction time of the double array trie.

## License

`MIT`
