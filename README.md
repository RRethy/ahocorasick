# Biblio

:zap: The fastest Golang implementation of the Aho-Corasick string matching algorithm, bar none.

## Usage

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
$ go test -bench=Benchmark
```

The two basic operations are the compilation of the state machine from an array of patterns (`Compile`), and the usage of this state machine to find each pattern in text (`FindAll`). Other implementations call these operations under different names.

| Operation | Input Size | rrethy/biblio | [BobuSumisu/aho-corasick](https://github.com/BobuSumisu/aho-corasick) | [anknown/ahocorasick](https://github.com/anknown/ahocorasick) |
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

### Other Implementations

Two implementations were intentionally omitted, [cloudflare/ahocorasick](https://github.com/cloudflare/ahocorasick) and [iohub/ahocorasick](https://github.com/iohub/ahocorasick). There are existing benchmarks [here](https://github.com/BobuSumisu/aho-corasick) and [here](https://github.com/anknown/ahocorasick) which have these implementations in their benchmark comparisons.

The *cloudflare* implementation was omitted due to incorrectness, it does not find all instances of the compiled patterns in the text. Even still, it has a significantly slower (~50x) `Compile` time and a slightly slower `FindAll` time (partly due to it reporting far less results than actually exist).

The *iohub* implementation was omitted since I was unable to get it working as it seemed to time out during compilation. As well, it has an unnecessarily confusing API.

## Implementation Details

**TODO**

## License

`MIT`
