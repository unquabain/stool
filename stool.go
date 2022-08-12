package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

//go:embed usage.txt
var Usage string

var InputFile string = `-`
var InputFormat Format = FormatUnknown
var OutputFile string = `-`
var SearchPath string = `.`
var OutputTemplate string = `{{ . | yaml }}`
var OutputTemplateFile string = ``

func getOpts() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, Usage, os.Args[0])
	}
	flag.StringVar(&InputFile, `in`, InputFile, `the file to read or - for STDIN`)
	flag.StringVar(&InputFile, `i`, InputFile, `the file to read or - for STDIN`)
	format := flag.String(`format`, ``, `the format of the input file; yaml|json anything else will try to auto-detect`)
	flag.StringVar(format, `f`, ``, `the format of the input file; yaml|json anything else will try to auto-detect`)
	flag.StringVar(&OutputFile, `out`, OutputFile, `the file to write to or - for STDOUT`)
	flag.StringVar(&OutputFile, `o`, OutputFile, `the file to write to or - for STDOUT`)
	flag.StringVar(&SearchPath, `search`, SearchPath, `a path to search the input data before rendering`)
	flag.StringVar(&SearchPath, `s`, SearchPath, `a path to search the input data before rendering`)
	flag.StringVar(&OutputTemplate, `template`, OutputTemplate, `a go template to use to render the output`)
	flag.StringVar(&OutputTemplate, `t`, OutputTemplate, `a go template to use to render the output`)
	flag.StringVar(&OutputTemplateFile, `template-file`, OutputTemplateFile, `read the template from this file instead of the command line`)
	flag.StringVar(&OutputTemplateFile, `T`, OutputTemplateFile, `read the template from this file instead of the command line`)
	flag.Parse()

	switch strings.ToLower(*format) {
	case `json`, `j`, `js`:
		InputFormat = FormatJSON
	case `yaml`, `yml`, `y`:
		InputFormat = FormatYAML
	}
	if infile := flag.Arg(0); InputFile == `-` && infile != `` {
		InputFile = infile
	}
	if outfile := flag.Arg(1); OutputFile == `-` && outfile != `` {
		OutputFile = outfile
	}
}

func main() {
	getOpts()
	data, err := Parse(InputFile, InputFormat)
	if err != nil {
		log.Print(err)
		os.Exit(-1)
	}
	filtered, err := Evaluate(data, SearchPath)
	if err != nil {
		log.Print(err)
		os.Exit(-1)
	}
	tmplt, err := GetTemplate(OutputTemplate, OutputTemplateFile)
	if err != nil {
		log.Print(err)
		os.Exit(-1)
	}
	buff := new(bytes.Buffer)
	if err := tmplt.Execute(buff, filtered); err != nil {
		log.Print(err)
		os.Exit(-1)
	}
	if OutputFile == `-` {
		fmt.Printf(`%s`, string(buff.Bytes()))
		return
	}
	if err := ioutil.WriteFile(OutputFile, buff.Bytes(), 0644); err != nil {
		log.Print(err)
		os.Exit(-1)
	}
}
