package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var InputFile string = `-`
var InputFormat Format = FormatUnknown
var OutputFile string = `-`
var SearchPath string = `.`
var OutputTemplate string = `{{ . | yaml }}`
var OutputTemplateFile string = ``

func getOpts() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `A tool for querying and reformatting structured files.

Usage %s [options] filename [outputfile]

OPTIONS:
  --in -i
    The structured file to read, or - for STDIN (the default).
    This can also be specified with the first positional argument.

  --out -o
    The file to write out, or - for STDOUT (the default).
    This can also be specified with the second positional argument.
  
  --search -s 
    The search path. This is a little like JQuery. It consists mainly of
    dot-delimited path elements.
      Example: locale.en.errors.0.text

    Brackets can be used if the element might have conflicting syntax.
      Example: keys["key with spaces and dot."].value

    A few functions have been defined. They are "length()", "yaml()" and
    "json()".
      Example: local.en.errors.yaml()

    The special value "[*]" will resolve to all the values of an array
    or dictionary.
      Example: contacts[*].name

    You can also do simple tests that consist of a path, a comparison,
    and a value.
      Example: contacts[zip_code == "90210"].name
      Example: contacts[phones.length() > 2].name

  --template -t
    A Go Text Template to render the results. The template will be
    repeated for each result, so if you used [*] anywhere in your search
    pattern, the entire template will be repeated in the output. You may
    use Masterminds Sprig functions as well as the "yaml" and "json"
    functions.
      Example: The secret is {{ .client_secret | squote }}

    Currently, the only way to enter a carriage return is with {{ '\x0A' }}.
    So if you're expecting multiple documents, you may want to use
    the --template-file option.
      Example: ./stool -s '[*][client_title == "ekg"].client_id' \
                   -t "ID: {{ . }}{{ "'"\x0A"'" }}" \
                   rockauth_clients.yml

  --template-file -T
    For longer templates, you may specify a file to read instead of
    putting the template on the command line.

  --format -f
    The input file format. If the program cannot guess the file format,
    you may specify it as either "json" or "yaml".
`, os.Args[0])
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
