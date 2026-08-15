goamd64 := "v3"
goexperiment := "jsonv2,simd"

# List the available recipes.
default:
    @just --list

# Run the full five-sample benchmark suite and save its output.
bench output="bench.txt":
    cd bench && GOAMD64={{ goamd64 }} GOEXPERIMENT={{ goexperiment }} go test -test.fullpath=true -benchmem -run='^$' -count=5 -bench='^Benchmark' . | tee "../{{ output }}"

# Compile the three binary-size comparison commands in place.
build-cmd:
    cd bench && GOAMD64={{ goamd64 }} GOEXPERIMENT={{ goexperiment }} go build -trimpath -ldflags="-s -w" -o cmd/size-experiment/size-experiment ./cmd/size-experiment
    cd bench && GOAMD64={{ goamd64 }} GOEXPERIMENT={{ goexperiment }} go build -trimpath -ldflags="-s -w" -o cmd/size-sonic/size-sonic ./cmd/size-sonic
    cd bench && GOAMD64={{ goamd64 }} GOEXPERIMENT={{ goexperiment }} go build -trimpath -ldflags="-s -w" -o cmd/size-stdlib/size-stdlib ./cmd/size-stdlib

# Run the test suite with the benchmark build configuration.
test:
    GOAMD64={{ goamd64 }} GOEXPERIMENT={{ goexperiment }} go test ./...
