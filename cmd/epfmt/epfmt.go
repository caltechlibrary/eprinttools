// Package eprinttools is a collection of structures, functions and programs// for working with the EPrints XML and EPrints REST API
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2021, Caltech
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
package main

//
// epfmt is a command line tool convert EPrints XML to/from JSON.
//

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	// Caltech Library Packages
	"github.com/caltechlibrary/eprinttools"
)

const (
	IsXML = iota
	IsJSON
)

var (
	helpText = `---
title: "{app_name} (1) user manual"
author: "R. S. Doiel"
pubDate: 2023-01-11
---

# NAME

{app_name}

# SYNOPSIS

{app_name} 

# DESCRIPTION

{app_name} is a command line program for pretty printing
EPrint XML. It can also convert EPrint XML to and from
JSON. By default it reads from standard input and writes to
standard out.

{app_name} EPrint XML (or JSON version) from
standard input and pretty prints the result to 
standard out. You can change output format XML 
and JSON by using either the '-xml' or '-json' 
option. The XML representation is based on EPrints 
3.x.  {app_name} does NOT interact with the EPrints API 
only the the document presented via standard
input.

# OPTIONS

-help
: display help

-license
: display license

-i, -input
: (string) input file name (read the URL connection string from the input file

-json
: output JSON version of EPrint XML

-nl, -newline
: if true add a trailing newline

-o, -output
: (string) output file name

-quiet
: suppress error messages

-s, -simplified
: output simplified JSON version of EPrints XML

-version
: display version

-xml
: output EPrint XML

# EXAMPLES

Pretty print EPrint XML as XML.

~~~
    {app_name} < 123.xml
~~~

Pretty print from EPrint XML as JSON

~~~
    {app_name} -json < 123.xml
~~~

Render EPrint JSON as EPrint XML.

~~~
    {app_name} -xml < 123.json
~~~

{app_name} will first parse the XML or JSON 
presented to it and pretty print the output 
in the desired format requested. If no 
format option chosen it will pretty print 
in the same format as input.

{app_name} {version}

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

	// App Options
	asJSON       bool
	asXML        bool
	asSimplified bool
)

func fmtTxt(src string, appName string, version string) string {
	return strings.ReplaceAll(strings.ReplaceAll(src, "{app_name}", appName), "{version}", version)
}

func main() {
	var (
		inputFmt int
		obj      *eprinttools.EPrints
		src      []byte
	)

	appName := path.Base(os.Args[0])

	// Standard Options
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.StringVar(&inputFName, "i", "", "input file name (read the URL connection string from the input file")
	flag.StringVar(&inputFName, "input", "", "input file name (read the URL connection string from the input file")
	flag.StringVar(&outputFName, "o", "", "output file name")
	flag.StringVar(&outputFName, "output", "", "output file name")
	flag.BoolVar(&quiet, "quiet", false, "suppress error messages")
	flag.BoolVar(&newLine, "nl", false, "if true add a trailing newline")
	flag.BoolVar(&newLine, "newline", false, "if true add a trailing newline")

	// App Options
	flag.BoolVar(&asXML, "xml", false, "output EPrint XML")
	flag.BoolVar(&asJSON, "json", false, "output JSON version of EPrint XML")
	flag.BoolVar(&asSimplified, "s", false, "output simplified JSON version of EPrints XML")
	flag.BoolVar(&asSimplified, "simplified", false, "output simplified JSON version of EPrints XML")

	// We're ready to process args
	flag.Parse()
	args := flag.Args()

	// Simplified output is aways JSON formatted.
	if asSimplified {
		asJSON = true
		asXML = false
	}

	// Setup IO
	var err error

	in := os.Stdin
	out := os.Stdout
	eout := os.Stderr


	if len(args) > 1 {
		inputFName = args[1]
	}
	if len(args) > 2 {
		outputFName = args[2]
	}

	if inputFName != "" && inputFName != "-" {
		in, err = os.Open(inputFName)
		if err != nil {
			fmt.Fprintln(eout, err)
			os.Exit(1)
		}
		defer in.Close()
	}

	if outputFName != "" && outputFName != "-" {
		out, err = os.Create(outputFName)
		if err != nil {
			fmt.Fprintln(eout, err)
			os.Exit(1)
		}
		defer out.Close()
	}

	// Handle options
	if showHelp {
		fmt.Fprintf(out, "%s\n", fmtTxt(helpText, appName, eprinttools.Version))
		os.Exit(0)
	}
	if showLicense {
		fmt.Fprintf(out, "%s\n", eprinttools.LicenseText)
		os.Exit(0)
	}
	if showVersion {
		fmt.Fprintf(out, "%s %s\n", appName, eprinttools.Version)
		os.Exit(0)
	}

	// Read the file to []byte
	src, err = ioutil.ReadAll(in)
	if err != nil {
		fmt.Fprintln(eout, err)
		os.Exit(1)
	}
	// Trim leading/trailing spaces
	src = bytes.TrimSpace(src)
	if len(src) == 0 {
		fmt.Fprintln(eout, "Nothing to parse")
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
		fmt.Fprintln(eout, err)
		os.Exit(1)
	}

	for _, e := range obj.EPrint {
		e.SyntheticFields()
	}
	if asJSON == false && asXML == false {
		asXML = (inputFmt == IsXML)
	}
	if asSimplified {
		if len(obj.EPrint) == 1 {
			sObject, err := eprinttools.CrosswalkEPrintToRecord(obj.EPrint[0])
			if err != nil {
				fmt.Fprintln(eout, err)
				os.Exit(1)
			}
			src, err = json.MarshalIndent(sObject, "", "   ")
		} else {
			sObjects := []*eprinttools.Record{}
			for _, eprint := range obj.EPrint {
				obj, err := eprinttools.CrosswalkEPrintToRecord(eprint)
				if err != nil {
					fmt.Fprintln(eout, err)
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
		fmt.Fprintln(eout, err)
		os.Exit(1)
	}
	// Print the doc string if XML
	if asXML {
		fmt.Fprintf(out, "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n")
	}

	eol := ""
	if newLine {
		eol = "\n"
	}
	fmt.Fprintf(out, "%s%s", src, eol)
	os.Exit(0)
}
