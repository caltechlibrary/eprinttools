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
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	// Caltech Library Packages
	"github.com/caltechlibrary/eprinttools"
)

const (
	IsXML = iota
	IsJSON
)

var (
	description = `
USAGE
	{appName} 

{appName} is a command line program for pretty printing
EPrint XML. It can also convert EPrint XML to and from
JSON. By default it reads from standard input and writes to
standard out.

DESCRIPTION

{appName} EPrint XML (or JSON version) from
standard input and pretty prints the result to 
standard out. You can change output format XML 
and JSON by using either the '-xml' or '-json' 
option. The XML representation is based on EPrints 
3.x.  {appName} does NOT interact with the EPrints API 
only the the document presented via standard
input.
`

	examples = `
Pretty print EPrint XML as XML.

    {appName} < 123.xml

Pretty print from EPrint XML as JSON

    {appName} -json < 123.xml

Render EPrint JSON as EPrint XML.

    {appName} -xml < 123.json

{appName} will first parse the XML or JSON 
presented to it and pretty print the output 
in the desired format requested. If no 
format option chosen it will pretty print 
in the same format as input.
`

	license = `
{appName} {version}

Copyright (c) 2019, Caltech
All rights not granted herein are expressly reserved by Caltech.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
`

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
		src      []byte
		err      error
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

	// App Options
	flagSet.BoolVar(&asXML, "xml", false, "output EPrint XML")
	flagSet.BoolVar(&asJSON, "json", false, "output JSON version of EPrint XML")
	flagSet.BoolVar(&simplified, "s", false, "output simplified JSON version of EPrints XML")
	flagSet.BoolVar(&simplified, "simplified", false, "output simplified JSON version of EPrints XML")

	// We're ready to process args
	flagSet.Parse(os.Args[1:])
	args := flagSet.Args()

	if len(args) > 1 {
		inputFName = args[1]
	}
	if len(args) > 2 {
		outputFName = args[2]
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
			fmt.Fprintf(os.Stderr, "%s\n", err)
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

	// Read the file to []byte
	src, err = ioutil.ReadAll(in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	// Trim leading/trailing spaces
	src = bytes.TrimSpace(src)
	if len(src) == 0 {
		fmt.Fprintf(os.Stderr, "Nothing to parse\n")
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
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	for _, e := range obj.EPrint {
		e.SyntheticFields()
	}
	if asJSON == false && asXML == false {
		asXML = (inputFmt == IsXML)
	}
	if simplified {
		fmt.Printf("DEBUG generate simplified invenio style record\n")
		asXML = false
		asJSON = true
		if len(obj.EPrint) == 1 {
			sObject, err := eprinttools.CrosswalkEPrintToRecord(obj.EPrint[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
			src, err = json.MarshalIndent(sObject, "", "   ")
		} else {
			sObjects := []*eprinttools.Record{}
			for _, eprint := range obj.EPrint {
				obj, err := eprinttools.CrosswalkEPrintToRecord(eprint)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", err)
					os.Exit(1)
				}
				sObjects = append(sObjects, obj)
			}
			src, err = json.MarshalIndent(sObjects, "", "   ")
		}
	} else {
		// marshal pretty printed output based on options selected.
		if asXML {
			src, err = xml.MarshalIndent(obj, "", "  ")
		} else {
			src, err = json.MarshalIndent(obj, "", "   ")
		}
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	// Print the doc string if XML
	if asXML {
		fmt.Fprintf(out, "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n")
	}

	if newLine {
		fmt.Fprintf(out, "%s\n", src)
	} else {
		fmt.Fprintf(out, "%s", src)
	}
	os.Exit(0)
}
