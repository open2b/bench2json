package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/cep21/benchparse"
)

func main() {
	var p, i string
	var h bool
	flag.StringVar(&p, "p", "", "filter the benchmarks by program name")
	flag.StringVar(&i, "i", "", "filter the benchmarks by interpreter, "+
		"accepted values are Scriggo,Yaegi, Tengo, GoLua, GopherLua.\n"+
		"\t\texample: -i=Scriggo,Yaegi")
	flag.BoolVar(&h, "h", false, "print this help")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "bench2json returns a json encoded representation of the benchmark data passed by command line\n"+
			"Benhmark names are expected to respect the following format:\n"+
			"Benchmark[INTERPRETER NAME]/[PROGRAM NAME].[FILE EXTENSION]-[PROC NUM]\n"+
			"the benchmarks are expected to be run with the -test.benchmem option\n")
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(os.Stderr, " -%s %s\n", f.Name, f.Usage)
		})
		fmt.Fprintf(os.Stderr, "\nUsage example: $ go test -bench=. -test.benchmem | bench2json\n")
	}

	flag.Parse()

	if h {
		flag.Usage()
		os.Exit(0)
	}
	var interpreters []string
	if i != "" {
		interpreters = strings.Split(i, ",")
		for i, interp := range interpreters {
			interpreters[i] = strings.ToLower(interp)
			switch interpreters[i] {
			case "golua":
				interpreters[i] = "go-lua"
			case "gopherlua":
				interpreters[i] = "gopher-lua"
			}
		}
	}

	res := encodeBenchmarkData(os.Stdin, options{
		Interpreters: interpreters,
		Program:      p,
	})
	fmt.Println(res)
}

type options struct {
	Interpreters []string
	Program      string
}

var benchmarkProgram = regexp.MustCompile(`Benchmark(.+)\/(.+)\..+-\d+`)

type Benchmark struct {
	Program string
	Results map[string]BenchmarkResult
}

type BenchmarkResult struct {
	Time   float64
	Allocs float64
}

const TimeColumn = 0
const AllocsColumn = 2

func encodeBenchmarkData(reader io.Reader, o options) string {
	d := benchparse.Decoder{}
	run, err := d.Decode(reader)
	if err != nil {
		panic(err)
	}

	benchmarkOf := map[string]Benchmark{}

	for _, result := range run.Results {
		res := benchmarkProgram.FindStringSubmatch(result.Name)
		interpreter := strings.ToLower(res[1])
		switch interpreter {
		case "golua":
			interpreter = "go-lua"
		case "gopherlua":
			interpreter = "gopher-lua"
		}
		program := res[2]

		if o.Program != "" && o.Program != program {
			continue
		}
		if len(o.Interpreters) > 0 {
			var inList bool
			for _, i := range o.Interpreters {
				if i == interpreter {
					inList = true
					break
				}
			}
			if !inList {
				continue
			}
		}
		if _, ok := benchmarkOf[program]; !ok {
			benchmarkOf[program] = Benchmark{
				Program: program,
				Results: map[string]BenchmarkResult{},
			}
		}
		benchmarkOf[program].Results[interpreter] = BenchmarkResult{
			Time:   result.Values[TimeColumn].Value,
			Allocs: result.Values[AllocsColumn].Value,
		}
	}

	benchmarks := make([]Benchmark, 0, len(benchmarkOf))
	for _, b := range benchmarkOf {
		benchmarks = append(benchmarks, b)
	}

	encoded, err := json.MarshalIndent(benchmarks, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(encoded)
}
