# bench2json

bench2json is a tool to transform the output of go benchmarks to a json ready to be embedded on the scriggo site.

Benchmark names are expected to respect the following format:

Benchmark[INTERPRETER NAME]/[PROGRAM NAME].[FILE EXTENSION]-[PROC NUM]

the benchmarks are expected to be run with the -test.benchmem option.

### Installation

```shell
go install github.com/open2b/bench2json
```

### Usage
* **-h** prints the help
* **-i** filter the benchmarks by interpreter, accepted values are Scriggo, Yaegi, Tengo, GoLua, GopherLua and Goja.

     Eg: -i=Scriggo,Yaegi
* **-p** filter the benchmarks by program name

**Usage example**: 

```
go test -bench=. -test.benchmem | bench2json
```