package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	// Caltech Library Packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/eprinttools"
)

var (
	description = []byte(`
	epparsexml parses XML content retrieved from disc or the EPrints API. It will 
	render JSON if the XML is valid otherwise return errors.
`)

	examples = []byte(`
Parse an EPrint reversion XML file.

	epparsexml -revision revision/2.xml
`)

	// Standard Options
	showHelp             bool
	showLicense          bool
	showVersion          bool
	showExamples         bool
	newLine              bool
	quiet                bool
	verbose              bool
	generateMarkdownDocs bool
	inputFName           string
	outputFName          string

	// App Options
	eprints bool
	eprint  bool
)

func main() {
	app := cli.NewCli(eprinttools.Version)

	// Add Help Docs
	app.AddHelp("description", description)
	app.AddHelp("examples", examples)

	// Standard Options
	app.BoolVar(&showHelp, "h,help", false, "display help")
	app.BoolVar(&showLicense, "l,license", false, "display license")
	app.BoolVar(&showVersion, "v,version", false, "display version")
	app.BoolVar(&showExamples, "e,examples", false, "display examples")
	app.StringVar(&inputFName, "i,input", "", "input file name")
	app.StringVar(&outputFName, "o,output", "", "output file name")
	app.BoolVar(&quiet, "quiet", false, "suppress error messages")
	app.BoolVar(&newLine, "nl,newline", false, "if true add a trailing newline")
	app.BoolVar(&generateMarkdownDocs, "generate-markdown-docs", false, "output documentation in Markdown")

	// App Options
	app.BoolVar(&eprints, "document,eprints", false, "parse an eprints (e.g. rest response) document")
	app.BoolVar(&eprint, "revision,eprint", false, "parse a eprint (revision) document")

	// We're ready to process args
	app.Parse()
	args := app.Args()

	// Setup IO
	var err error

	app.Eout = os.Stderr
	app.In, err = cli.Open(inputFName, os.Stdin)
	cli.ExitOnError(app.Eout, err, quiet)
	defer cli.CloseFile(inputFName, app.In)

	app.Out, err = cli.Create(outputFName, os.Stdout)
	cli.ExitOnError(app.Eout, err, quiet)
	defer cli.CloseFile(outputFName, app.Out)

	// Handle options
	if generateMarkdownDocs {
		app.GenerateMarkdownDocs(app.Out)
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
	cli.ExitOnError(app.Eout, err, quiet)

	switch {
	case eprints:
		data := eprinttools.EPrints{}
		err = xml.Unmarshal(src, &data)
		cli.ExitOnError(app.Eout, err, quiet)

		src, err = json.MarshalIndent(data, "", " ")
		cli.ExitOnError(app.Eout, err, quiet)
	case eprint:
		data := eprinttools.EPrint{}
		err = xml.Unmarshal(src, &data)
		cli.ExitOnError(app.Eout, err, quiet)

		src, err = json.MarshalIndent(data, "", " ")
		cli.ExitOnError(app.Eout, err, quiet)
	}

	fmt.Fprintf(os.Stdout, "%s\n", src)
	os.Exit(0)
}
