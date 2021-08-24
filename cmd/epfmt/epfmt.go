//
// epfmt is a command line tool convert EPrints XML to/from JSON.
//
// @author R. S. Doiel, <rsdoiel@library.caltech.edu>
//
// Copyright (c) 2019, Caltech
// All rights not granted herein are expressly reserved by Caltech.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	// Caltech Library Packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/eprinttools"
)

const (
	IsXML = iota
	IsJSON
)

var (
	synopsis = []byte(`
_epfmt_ is a command line program for 
pretty printing EPrint XML. It can also convert
EPrint XML to and from JSON.
`)
	description = []byte(`
_epfmt_ parses EPrint XML (or JSON version) from
standard input and pretty prints the result to 
standard out. You can change output format XML 
and JSON by using either the '-xml' or '-json' 
option. The XML representation is based on EPrints 
3.x.  _epfmt_ does NOT interact with the EPrints API 
only the the document presented via standard
input.
`)

	examples = []byte(`
Pretty print EPrint XML as XML.

` + "```" + `
    epfmt < 123.xml
` + "```" + `

Pretty print from EPrint XML as JSON

` + "```" + `
    epfmt -json < 123.xml
` + "```" + `

Render EPrint JSON as EPrint XML.

` + "```" + `
    epfmt -xml < 123.json
` + "```" + `

_epfmt_ will first parse the XML or JSON 
presented to it and pretty print the output 
in the desired format requested. If no 
format option chosen it will pretty print 
in the same format as input.
`)

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
	simplified       bool

	// App Options
	asJSON bool
	asXML  bool
)

func main() {
	var (
		inputFmt int
		obj      *eprinttools.EPrints
		sObj     []*eprinttools.SimplePrint
		src      []byte
		err      error
	)

	app := cli.NewCli(eprinttools.Version)

	// Add Help Docs
	app.AddHelp("synopsis", synopsis)
	app.AddHelp("description", description)
	app.AddHelp("examples", examples)

	// Standard Options
	app.BoolVar(&showHelp, "h,help", false, "display help")
	app.BoolVar(&showLicense, "l,license", false, "display license")
	app.BoolVar(&showVersion, "v,version", false, "display version")
	app.BoolVar(&showExamples, "e,examples", false, "display examples")
	app.StringVar(&inputFName, "i,input", "", "input file name (read the URL connection string from the input file")
	app.StringVar(&outputFName, "o,output", "", "output file name")
	app.BoolVar(&quiet, "quiet", false, "suppress error messages")
	app.BoolVar(&newLine, "nl,newline", false, "if true add a trailing newline")
	app.BoolVar(&generateMarkdown, "generate-markdown", false, "generate Markdown documentation")
	app.BoolVar(&generateManPage, "generate-manpage", false, "generate man page")

	// App Options
	app.BoolVar(&asXML, "xml", false, "output EPrint XML")
	app.BoolVar(&asJSON, "json", false, "output JSON version of EPrint XML")
	app.BoolVar(&simplified, "s,simplified", false, "output simplified JSON version of EPrints XML")

	// We're ready to process args
	app.Parse()
	args := app.Args()

	if len(args) > 1 {
		inputFName = args[1]
	}
	if len(args) > 2 {
		outputFName = args[2]
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

	// Read the file to []byte
	src, err = ioutil.ReadAll(app.In)
	if err != nil {
		fmt.Fprintf(app.Eout, "%s\n", err)
		os.Exit(1)
	}
	// Trim leading/trailing spaces
	src = bytes.TrimSpace(src)
	if len(src) == 0 {
		fmt.Fprintf(app.Eout, "Nothing to parse\n")
		os.Exit(1)
	}

	// Check if JSON or XML
	if bytes.HasPrefix(src, []byte("{")) {
		// Unmarshal as JSON
		inputFmt = IsJSON
		err = json.Unmarshal(src, &obj)
	} else {
		// Unmarshal as EPrintXML
		inputFmt = IsXML
		err = xml.Unmarshal(src, &obj)
	}
	if err != nil {
		fmt.Fprintf(app.Eout, "%s\n", err)
		os.Exit(1)
	}

	for _, e := range obj.EPrint {
		e.SyntheticFields()
	}
	if asJSON == false && asXML == false {
		asXML = (inputFmt == IsXML)
	}
	if simplified {
		asXML = false
		asJSON = true
		sObj, err = eprinttools.SimplifyEPrints(obj)
		if err != nil {
			fmt.Fprintf(app.Eout, "%s\n", err)
			os.Exit(1)
		}
		src, err = json.MarshalIndent(sObj, "", "   ")
	} else {
		// marshal pretty printed output based on options selected.
		if asXML {
			src, err = xml.MarshalIndent(obj, "", "  ")
		} else {
			src, err = json.MarshalIndent(obj, "", "   ")
		}
	}
	if err != nil {
		fmt.Fprintf(app.Eout, "%s\n", err)
		os.Exit(1)
	}
	// Print the doc string if XML
	if asXML {
		fmt.Fprintf(app.Out, "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n")
	}

	if newLine {
		fmt.Fprintf(app.Out, "%s\n", src)
	} else {
		fmt.Fprintf(app.Out, "%s", src)
	}
	os.Exit(0)
}
