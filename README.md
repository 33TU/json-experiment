# json-experiment

An experimental, performance-focused JSON marshaler for Go.

## Benchmarks

Benchmark 7 compares matching output-ownership contracts:

- Owned output: `Marshal`, Go's `encoding/json`, and
  `sonic.ConfigFastest.Marshal`.
- Reusable output: `MarshalAppend`, Go's experimental `encoding/json/v2` using
  `MarshalWrite`, and Sonic's `EncodeInto` API.

The append path can avoid output allocation when the destination has sufficient
capacity. The charts show owned and reusable APIs side by side, but scale each
workload independently so both small values and large composite values remain
readable. Exact median latency and allocation counts are printed on every bar.

### Benchmark 7

This is a five-run median on Go 1.26 using `GOAMD64=v3` with the JSON v2 and
SIMD experiments enabled, measured on an AMD Ryzen 9 9950X3D. The suite
includes maps, primitive slices, numbers, struct variants, UTF-8 validation,
and standard marshaling interfaces.

#### Values and structs

![Benchmark 7 values and structs](assets/benchmarks/benchmark7-values.svg)

#### Maps and field options

![Benchmark 7 maps and field options](assets/benchmarks/benchmark7-maps-options.svg)

#### UTF-8 and marshaling interfaces

![Benchmark 7 UTF-8 and marshaling interfaces](assets/benchmarks/benchmark7-utf8-interfaces.svg)

```sh
just bench
```

Five-run median latency for owned-output APIs (lower is better):

| Workload                        |  Marshal | encoding/json | Sonic Marshal |
| ------------------------------- | -------: | ------------: | ------------: |
| `map[string]int`                | 149.3 ns |      589.3 ns |      213.2 ns |
| `map[string][]int`              | 167.2 ns |      707.0 ns |      225.1 ns |
| `map[string]any`                | 219.1 ns |       1352 ns |      323.5 ns |
| `[]int`                         | 68.23 ns |      165.4 ns |      84.74 ns |
| `float32`                       | 41.16 ns |      80.34 ns |      54.73 ns |
| `float64`                       | 57.02 ns |      104.8 ns |      68.41 ns |
| mixed struct                    | 211.6 ns |      726.3 ns |      243.0 ns |
| struct slice with metadata maps | 542.1 ns |       2072 ns |      664.3 ns |
| quoted struct fields            | 139.9 ns |      660.1 ns |      133.6 ns |
| `omitempty` / `omitzero`        | 90.15 ns |      291.4 ns |      105.4 ns |
| `omitzero` fields               | 80.93 ns |      260.4 ns |      178.3 ns |
| UTF-8 validation: ASCII         | 188.1 ns |      538.5 ns |      255.2 ns |
| UTF-8 validation: Unicode       | 144.2 ns |      427.8 ns |      189.6 ns |
| UTF-8 validation: invalid byte  | 268.8 ns |      667.7 ns |      882.8 ns |
| `encoding.TextMarshaler`        | 35.47 ns |      91.40 ns |      54.68 ns |
| `json.Marshaler`                | 30.07 ns |      70.72 ns |      50.26 ns |
| integer `json.Marshaler`        | 29.90 ns |      77.54 ns |      50.35 ns |

Five-run median latency for reusable-output APIs (lower is better):

| Workload                        | MarshalAppend | json/v2 Write | Sonic EncodeInto |
| ------------------------------- | ------------: | ------------: | ---------------: |
| `map[string]int`                |      111.8 ns |      403.9 ns |         168.4 ns |
| `map[string][]int`              |      127.5 ns |      488.1 ns |         179.6 ns |
| `map[string]any`                |      167.2 ns |       1066 ns |         272.2 ns |
| `[]int`                         |      41.26 ns |      148.0 ns |         57.53 ns |
| `float32`                       |      26.66 ns |      77.85 ns |         40.04 ns |
| `float64`                       |      37.44 ns |      88.63 ns |         44.22 ns |
| mixed struct                    |      152.4 ns |      582.6 ns |         206.2 ns |
| struct slice with metadata maps |      441.2 ns |       1603 ns |         553.2 ns |
| quoted struct fields            |      97.38 ns |      332.9 ns |         94.09 ns |
| `omitempty` / `omitzero`        |      58.69 ns |      304.8 ns |         73.67 ns |
| `omitzero` fields               |      50.32 ns |      243.5 ns |         143.8 ns |
| UTF-8 validation: ASCII         |      43.56 ns |      406.7 ns |         64.87 ns |
| UTF-8 validation: Unicode       |      55.12 ns |      333.4 ns |         71.00 ns |
| UTF-8 validation: invalid byte  |      74.31 ns |      497.3 ns |         879.6 ns |
| `encoding.TextMarshaler`        |      18.35 ns |      74.97 ns |         41.94 ns |
| `json.Marshaler`                |      13.18 ns |      68.56 ns |         52.08 ns |
| integer `json.Marshaler`        |      13.05 ns |      70.84 ns |         52.66 ns |

`MarshalAppend` records zero allocations for ordinary values and one allocation
when invoking the marshaling interfaces. Allocation counts for every
implementation are included directly in the charts and raw output.

The custom `MarshalerAppend` precedence benchmark only applies to this package,
so it is omitted from the comparison charts. Its median results were 5.893 ns
with zero allocations for `MarshalAppend`, and 22.75 ns with one allocation for
`Marshal`.

The complete Benchmark 7 output is available in
[`bench7.txt`](assets/benchmarks/raw/bench7.txt). Earlier benchmark SVGs and raw
outputs remain under [`assets/benchmarks`](assets/benchmarks), with older
commentary preserved in [`README-old.md`](README-old.md).

### Complete Benchmark 7 overview

![Complete Benchmark 7 comparison](assets/benchmarks/benchmark7.svg)

---

Results vary by hardware, Go version, and workload. Run the benchmarks on the
target system before drawing conclusions for a particular application.

## Large-struct benchmark

This focused benchmark encodes a 4,740-byte payload containing nested structs,
slices, maps, field options, ASCII, and Unicode. Results are five-run medians on
an AMD Ryzen 9 9950X3D using Go 1.26, `GOAMD64=v3`, and the JSON v2 and SIMD
experiments.

![Large nested struct benchmark](assets/benchmarks/large-struct.svg)

| Output contract | Encoder          |      Latency |    Throughput | Allocated bytes | Allocations |
| --------------- | ---------------- | -----------: | ------------: | --------------: | ----------: |
| Reusable        | MarshalAppend    | **2.473 µs** | **1916 MB/s** |         **0 B** |       **0** |
| Owned           | Marshal          | **3.124 µs** | **1517 MB/s** |      **4874 B** |       **1** |
| Reusable        | Sonic EncodeInto |     3.640 µs |     1302 MB/s |          2162 B |          18 |
| Owned           | Sonic Marshal    |     4.230 µs |     1120 MB/s |          7260 B |          19 |
| Reusable        | json/v2 Write    |     9.034 µs |      525 MB/s |          1651 B |          34 |
| Owned           | encoding/json    |    11.331 µs |      418 MB/s |          7203 B |          81 |

For matching output contracts, `Marshal` is 1.35× faster than Sonic Marshal,
while `MarshalAppend` is 1.47× faster than Sonic EncodeInto and remains
allocation-free once the destination buffer has sufficient reusable capacity.
The complete five-run output is included in
[`bench7.txt`](assets/benchmarks/raw/bench7.txt).

### Parallel scaling

The same payload was also measured with `b.RunParallel`, giving each worker its
own reusable destination buffer. The table reports aggregate throughput from
three-run medians at each `GOMAXPROCS` setting; 16 corresponds to the physical
core count of the 9950X3D, while 32 includes SMT threads.

![Large nested struct parallel scaling](assets/benchmarks/large-struct-parallel.svg)

| CPUs |  MarshalAppend | Sonic EncodeInto | Reusable advantage |       Marshal | Sonic Marshal | Owned advantage |
| ---: | -------------: | ---------------: | -----------------: | ------------: | ------------: | --------------: |
|    1 |  **1.84 GB/s** |        1.12 GB/s |          **1.64×** | **1.39 GB/s** |     0.99 GB/s |       **1.40×** |
|    2 |  **3.57 GB/s** |        2.09 GB/s |          **1.71×** | **2.46 GB/s** |     1.69 GB/s |       **1.46×** |
|    4 |  **7.06 GB/s** |        4.57 GB/s |          **1.54×** | **4.44 GB/s** |     2.97 GB/s |       **1.49×** |
|    8 | **13.75 GB/s** |        7.07 GB/s |          **1.94×** | **6.11 GB/s** |     4.48 GB/s |       **1.36×** |
|   16 | **25.12 GB/s** |        8.81 GB/s |          **2.85×** | **7.83 GB/s** |     4.92 GB/s |       **1.59×** |
|   32 | **29.06 GB/s** |       10.22 GB/s |          **2.84×** | **9.20 GB/s** |     5.59 GB/s |       **1.65×** |

`MarshalAppend` scales by 13.65× from one to 16 CPUs and remains at zero bytes
and zero allocations per operation throughout. Sonic EncodeInto scales by 7.87×
over the same range and records 18 allocations per operation. These are
aggregate throughput measurements rather than per-request tail-latency results.

At 16 CPUs, json/v2 Write reaches 4.52 GB/s and `encoding/json` reaches 2.42
GB/s, making `MarshalAppend` 5.6× and 10.4× faster, respectively. The
standard-library series are omitted from the chart to keep the closer Sonic
comparison readable, but their complete measurements remain in the raw output.
The complete three-run scaling output is available in
[`bench-large-struct-parallel.txt`](assets/benchmarks/raw/bench-large-struct-parallel.txt).

```sh
cd bench
GOAMD64=v3 GOEXPERIMENT=jsonv2,simd go test -run='^$' -benchmem \
  -count=3 -benchtime=300ms -cpu=1,2,4,8,16,32 \
  -bench='^BenchmarkMarshalLargeStruct/parallel' .
```

## License

json-experiment is available under the [MIT License](LICENSE).
