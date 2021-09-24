//
// eprintxml2json.go - converts EPrints XML to JSON
//
package main

import (
"path"
"flag"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	// Caltech Library packages
	"github.com/caltechlibrary/eprinttools"
)

var (
	description = `
USAGE
	{appName} [OPTIONS] EPRINT_XML_FILENAME

SYNOPSIS

{appName} converts EPrintXML documents to JSON

DETAIL

{appName} converts EPrintXML like that retrieved from the
EPrint 3.x REST API to JSON. If no filename is provided on
the command line then standard input is used to read the EPrint
XML. If the EPrint XML isn't understood then an error message
will be written and an exit code of 1 used to close the process
otherwise the process will render JSON to standard out.
`

	examples = `Converting a document, eprints-dump.xml, to JSON.

    {appName} eprints-dump.xml

Or

    cat eprints-dump.xml | {appName}

`
	license = `
{appName} {version}

Copyright (c) 2018, Caltech
All rights not granted herein are expressly reserved by Caltech.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
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
	inputFName       string
	outputFName      string
	prettyPrint      bool
)

func main() {
	var (
		err error
	)
	appName := path.Base(os.Args[0])
	flagSet := flag.NewFlagSet(appName, flag.ContinueOnError)

	// Standard Options
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.BoolVar(&showLicense, "license", false, "display license")
	flagSet.BoolVar(&showVersion, "version", false, "display version")

	flagSet.StringVar(&inputFName, "i", "", "input file name (read the URL connection string from the input file")
	flagSet.StringVar(&inputFName, "input", "", "input file name (read the URL connection string from the input file")
	flagSet.StringVar(&outputFName, "o", "", "output file name")
	flagSet.StringVar(&outputFName, "output", "", "output file name")

	flagSet.BoolVar(&quiet, "quiet", false, "suppress error messages")
	flagSet.BoolVar(&newLine, "nl", false, "if true add a trailing newline")
	flagSet.BoolVar(&newLine, "newline", false, "if true add a trailing newline")
	flagSet.BoolVar(&prettyPrint, "p", true, "pretty print output")
	flagSet.BoolVar(&prettyPrint, "pretty", true, "pretty print output")

	// We're ready to process args
	flagSet.Parse(os.Args)
	args := flagSet.Args()

	if len(args) > 0 {
		inputFName = args[0]
	}
	// Setup IO
	in := os.Stdin
	out := os.Stdout

	if inputFName != "" {
		if in, err = os.Open(inputFName); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		defer in.Close()
	}
	if outputFName != "" {
		if out, err = os.Create(outputFName); err != nil {
			fmt.Fprintf(os.Stdout, "%s\n", err)
			os.Exit(1)
		}
		defer out.Close()
	}

	// Handle options
	if showHelp {
		eprinttools.DisplayUsage(out, appName, flagSet, description, examples, license)
		os.Exit(0)
	}
	if showLicense {
		eprinttools.DisplayLicense(out, appName, license)
		os.Exit(0)
	}
	if showVersion {
		eprinttools.DisplayVersion(out, appName)
		os.Exit(0)
	}

	src, err := ioutil.ReadAll(in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	data := new(eprinttools.EPrints)
	err = xml.Unmarshal(src, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	//NOTE: populate the synthetic fields
	for _, e := range data.EPrint {
		e.SyntheticFields()
	}
	if prettyPrint {
		src, err = json.MarshalIndent(data, "", "    ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	} else {
		src, err = json.Marshal(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	}
	if newLine == true {
		fmt.Fprintf(out, "%s\n", src)
	} else {
		fmt.Fprintf(out, "%s", src)
	}
}
