# bench2json

bench2json outputs a json encoded representation of the benchmark data passed by command line.

Benhmark names are expected to respect the following format:

Benchmark[INTERPRETER NAME]/[PROGRAM NAME].[FILE EXTENSION]-[PROC NUM]

the benchmarks are expected to be run with the -test.benchmem option.

### Usage
* **-h** prints the help
* **-i** filter the benchmarks by interpreter, accepted values are Scriggo, Yaegi, Tengo, GoLua, GopherLua.

     Eg: -i=Scriggo,Yaegi
* **-p** filter the benchmarks by program name

**Usage example**: $ go test -bench=. -test.benchmem | bench2json
