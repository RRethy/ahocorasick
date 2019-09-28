# Biblio

## Usage

## Benchmarks

*macOS Mojave version 10.14.6*

*MacBook Pro (Retina, 13-inch, Early 2015)*

*Processor 3.1 GHz Intel Core i7*


```
$ go test -bench=Benchmark
```

The two basic operations are the compilation of the state machine from an array of patterns (`Compile`), and the usage of this state machine to find each pattern in text (`FindAll`). Other implementations call these operations under different names.

| Operation | Input Size | biblio | [BobuSumisu](https://github.com/BobuSumisu/aho-corasick) | [anknown](https://github.com/anknown/ahocorasick) |
| - | - | - | - | - |
| `Compile` | 235886 patterns | **202 ms** | 243 ms | 1651 ms |
| `Compile` | 23589 patterns  |  56 ms |  **54 ms** |  204 ms |
| `Compile` | 2359 patterns   |   **9471 µs** |  11911 µs |   14307 µs |
| `Compile` | 236 patterns    |   1115 µs |   1714 µs |    **1069 µs** |
| `Compile` | 24 patterns     |    162 µs |    247 µs |     **123 µs** |
| - | - | - | - | - |
| `FindAll` | 3227439 bytes | **38 ms** | 51 ms | 209 ms |
| `FindAll` | 318647 bytes  |  **3950 µs** |  4461 µs |  19219 µs |
| `FindAll` | 31626 bytes   |   **381 µs** |   391 µs |   1332 µs |
| `FindAll` | 3657 bytes    |    **34 µs** |    57 µs |    153 µs |

**NOTE**: `FindAll` uses a state machine compiled from 2359 patterns.

**NOTE**: `FindAll` time does **not** include the `Compile` time for the state machine.

**NOTE**: Two implementations were intentionally omitted, [cloudflare/ahocorasick](https://github.com/cloudflare/ahocorasick) and [iohub/ahocorasick](https://github.com/iohub/ahocorasick). The *cloudflare* implementation was omitted due to incorrectness. However, it has a significantly slower (~10x) `Compile` time and a competitive `FindAll` time (partly due to it reporting far less results than actually exist). The *iohub* implementation was omitted since I was unable to get it working, it kept timing out. As well, it has an unnecessarily confusing API.

## Implementation Details

**TODO**: Probably going to write a blog post. There are a lot of suddle things I did to improve the construction time of the double array trie.

## License

`MIT`
