# Biblio

WIP

## Benchmarks

*
Mackbook Pro
macOS Mojave version 10.14.6
MacBook Pro (Retina, 13-inch, Early 2015)
Processor 3.1 GHz Intel Core i7

$ go test -bench=Benchmark
*

| Operation | Input Size | biblio | [BobusSumisu](https://github.com/BobuSumisu/aho-corasick) | [anknown](https://github.com/anknown/ahocorasick) |
| - | - | - | - | - |
| `Compile` | 235886 patterns | **202864 µs** | 243112 µs | 1651203 µs |
| `Compile` | 23589 patterns  |  56626 µs |  **54502 µs** |  204531 µs |
| `Compile` | 2359 patterns   |   **9471 µs** |  11911 µs |   14307 µs |
| `Compile` | 236 patterns    |   1115 µs |   1714 µs |    **1069 µs** |
| `Compile` | 24 patterns     |    162 µs |    247 µs |     **123 µs** |
| `FindAll` | 3227439 bytes | **38018 µs** | 51468 µs | 209033 µs |
| `FindAll` | 318647 bytes  |  **3950 µs** |  4461 µs |  19219 µs |
| `FindAll` | 31626 bytes   |   **381 µs** |   391 µs |   1332 µs |
| `FindAll` | 3657 bytes    |    **34 µs** |    57 µs |    153 µs |

**NOTE**: Find Matches uses a state machine compiled from 2359 patterns.

**NOTE**: Two implementations were intentionally omitted, [cloudflare/ahocorasick](https://github.com/cloudflare/ahocorasick) and [iohub/ahocorasick](https://github.com/iohub/ahocorasick). The *cloudflare* implementation was omitted due to incorrectness. However, even still, it has a significantly slower (~10x) `Compile` time and a competitive `FindAll` time (partly due to it reporting far less results than actually exist). The *iohub* implementation was omitted since I was unable to get it working, it kept getting into some infinite loop. As well, it has an unnecessarily confusing API.
