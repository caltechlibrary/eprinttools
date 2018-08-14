//
// eprintxml2json.go - converts EPrints XML to JSON
//
package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	// Caltech Library packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/eprinttools"
)

var (
	synopsis = `_eprintxml2json_ convers EPrintXML documents to JSON`

	description = `_eprintxml2json_ converts EPrintXML like that
retrieved from the EPrint 3.x REST API to JSON. If no filename
is provided on the command line then standard input is used
to read the EPrint XML. If the EPrint XML isn't understood
then an error message will be written and an exit code of 1
used to close the process otherwise the process will render
JSON to standard out.
`

	examples = `Converting a document, eprints-dump.xml, to JSON.

` + "```" + `
    eprintxml2json eprints-dump.xml
` + "```" + `

Or

` + "```" + `
    cat eprints-dump.xml | eprintxml2json 
` + "```" + `

`

	// Standard Options
	// Standard Options
	showHelp         bool
	showLicense      bool
	showVersion      bool
	showExamples     bool
	newLine          bool
	quiet            bool
	verbose          bool
	generateMarkdown bool
	generateManPage  bool
	inputFName       string
	outputFName      string
	prettyPrint      bool
)

func main() {
	var (
		err error
	)
	app := cli.NewCli(eprinttools.Version)

	app.AddParams("[input filename]")

	// Add Help
	app.AddHelp("synopsis", []byte(synopsis))
	app.AddHelp("description", []byte(description))
	app.AddHelp("examples", []byte(examples))

	// Standard Options
	app.BoolVar(&showHelp, "h,help", false, "display help")
	app.BoolVar(&showLicense, "l,license", false, "display license")
	app.BoolVar(&showVersion, "v,version", false, "display version")
	app.BoolVar(&showExamples, "e,examples", false, "display examples")
	//app.StringVar(&inputFName, "i,input", "", "input file name (read the URL connection string from the input file")
	app.StringVar(&outputFName, "o,output", "", "output file name")
	app.BoolVar(&quiet, "quiet", false, "suppress error messages")
	app.BoolVar(&newLine, "nl,newline", false, "if true add a trailing newline")
	app.BoolVar(&generateMarkdown, "generate-markdown", false, "generate Markdown documentation")
	app.BoolVar(&generateManPage, "generate-manpage", false, "generate man page")
	app.BoolVar(&prettyPrint, "p,pretty", true, "pretty print output")

	// We're ready to process args
	app.Parse()
	args := app.Args()

	if len(args) > 0 {
		inputFName = args[0]
	}
	// Setup IO
	app.Eout = os.Stderr

	app.In, err = cli.Open(inputFName, os.Stdin)
	cli.ExitOnError(app.Eout, err, quiet)
	defer cli.CloseFile(inputFName, app.In)

	app.Out, err = cli.Create(outputFName, os.Stdout)
	cli.ExitOnError(app.Eout, err, quiet)
	defer cli.CloseFile(outputFName, app.Out)

	// Handle options
	if generateMarkdown {
		app.GenerateMarkdown(app.Out)
		os.Exit(0)
	}
	if generateManPage {
		app.GenerateManPage(app.Out)
		os.Exit(0)
	}
	if showHelp || showExamples {
		if len(args) > 0 {
			fmt.Fprintf(app.Out, app.Help(args...))
		} else {
			app.Usage(app.Out)
		}
		os.Exit(0)
	}
	if showLicense {
		fmt.Fprintln(app.Out, app.License())
		os.Exit(0)
	}
	if showVersion {
		fmt.Fprintln(app.Out, app.Version())
		os.Exit(0)
	}

	src, err := ioutil.ReadAll(app.In)
	if err != nil {
		fmt.Fprintf(app.Eout, "%s\n", err)
		os.Exit(1)
	}
	data := new(eprinttools.EPrints)
	err = xml.Unmarshal(src, &data)
	if err != nil {
		fmt.Fprintf(app.Eout, "%s\n", err)
		os.Exit(1)
	}
	if prettyPrint {
		src, err = json.MarshalIndent(data, "", "    ")
		if err != nil {
			fmt.Fprintf(app.Eout, "%s\n", err)
			os.Exit(1)
		}
	} else {
		src, err = json.Marshal(data)
		if err != nil {
			fmt.Fprintf(app.Eout, "%s\n", err)
			os.Exit(1)
		}
	}
	if newLine == true {
		fmt.Fprintf(app.Out, "%s\n", src)
	} else {
		fmt.Fprintf(app.Out, "%s", src)
	}
}
