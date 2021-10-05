//
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
//
package main

//
// doi2eprintsxml.go is a Caltech Library centric command line utility
// to query CrossRef API and DataCite API for metadata and
// return the results as an EPrints XML file suitable for importing
// into EPrints.
//

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	// Caltech Library packages
	"github.com/caltechlibrary/crossrefapi"
	"github.com/caltechlibrary/dataciteapi"
	"github.com/caltechlibrary/eprinttools"
	"github.com/caltechlibrary/eprinttools/clsrules"
)

var (
	description = `
USAGE
	{app_name} [OPTIONS] DOI

SYNOPSIS

{app_name} is a Caltech Library centric application that
takes one or more DOI, queries the CrossRef API
and if that fails the DataCite API and returns an
EPrints XML document suitable for import into
EPrints. The DOI can be in either their canonical
form or URL form (e.g. "10.1021/acsami.7b15651" or
"https://doi.org/10.1021/acsami.7b15651").

`

	examples = `
Example generating an EPrintsXML for one DOI

	{app_name} "10.1021/acsami.7b15651" > article.xml

Example generating an EPrintsXML for two DOI

	{app_name} "10.1021/acsami.7b15651" "10.1093/mnras/stu2495" > articles.xml

Example processing a list of DOIs in a text file into
an XML document called "import-articles.xml".

	{app_name} -i doi-list.txt -o import-articles.xml
`

	license = `
{app_name} {version}

Copyright (c) 2021, Caltech
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
	generateMarkdown bool
	generateManPage  bool
	inputFName       string
	outputFName      string
	quiet            bool

	// App specific options
	apiEPrintsURL                         string
	mailto                                string
	crossrefOnly                          bool
	dataciteOnly                          bool
	useCaltechLibrarySpecificRules        bool
	use_1_0_0_CaltechLibrarySpecificRules bool
	asJSON                                bool
)

func main() {
	appName := path.Base(os.Args[0])

	// Standard Options
	flagSet := flag.NewFlagSet(appName, flag.ContinueOnError)

	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.BoolVar(&showLicense, "license", false, "display license")
	flagSet.BoolVar(&showVersion, "version", false, "display app version")
	flagSet.StringVar(&inputFName, "i", "", "set input filename")
	flagSet.StringVar(&inputFName, "input", "", "set input filename")
	flagSet.BoolVar(&quiet, "quiet", false, "set quiet output")

	// Application Options
	flagSet.StringVar(&apiEPrintsURL, "eprints-url", "", "Sets the EPRints API URL")
	flagSet.BoolVar(&crossrefOnly, "c", false, "only search CrossRef API for DOI records")
	flagSet.BoolVar(&crossrefOnly, "crossref", false, "only search CrossRef API for DOI records")
	flagSet.BoolVar(&dataciteOnly, "d", false, "only search DataCite API for DOI records")
	flagSet.BoolVar(&dataciteOnly, "datacite", false, "only search DataCite API for DOI records")
	flagSet.BoolVar(&useCaltechLibrarySpecificRules, "clsrules", false, "Apply current Caltech Library Specific Rules to EPrintXML output")
	flagSet.BoolVar(&use_1_0_0_CaltechLibrarySpecificRules, "v1.0.0-clsrules", false, "Apply v1.0.0 Caltech Library Specific Rules to EPrintXML output")
	flagSet.BoolVar(&asJSON, "json", false, "output EPrint structure as JSON")

	//FIXME: Need to come up with a better way of setting this,
	// perhaps a config mode and save the setting in
	// $HOME/etc/${AppName}.json
	flagSet.StringVar(&mailto, "m,mailto", "helpdesk@library.caltech.edu", "set the mailto value for CrossRef API access")

	flagSet.Parse(os.Args[1:])
	args := flagSet.Args()

	if showHelp {
		eprinttools.DisplayUsage(os.Stdout, appName, flagSet, description, examples, license)
		os.Exit(1)
	}

	if showLicense {
		eprinttools.DisplayLicense(os.Stdout, appName, license)
		os.Exit(0)
	}

	if showVersion {
		eprinttools.DisplayVersion(os.Stdout, appName)
		os.Exit(0)
	}

	if len(args) < 1 && inputFName == "" {
		eprinttools.DisplayUsage(os.Stderr, appName, flagSet, description, examples, license)
		os.Exit(1)
	}

	// Setup I/O
	var (
		err error
	)
	out := os.Stdout
	in := os.Stdin
	if outputFName != "" {
		if out, err = os.Create(outputFName); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		defer out.Close()
	}

	if inputFName != "" {
		if in, err = os.Open(inputFName); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		defer in.Close()
	}

	if inputFName != "" {
		src, err := ioutil.ReadAll(in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		for _, line := range strings.Split(fmt.Sprintf("%s", src), "\n") {
			arg := strings.TrimSpace(line)
			if len(arg) > 0 {
				args = append(args, arg)
			}
		}
	}

	// NOTE: OK we're ready to run our conversions
	eprintsList := new(eprinttools.EPrints)
	//FIXME: If crossrefapi returns 404 then we need to
	// query the dataciteapi
	apiCrossRef, err := crossrefapi.NewCrossRefClient(appName, mailto)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	apiDataCite, err := dataciteapi.NewDataCiteClient(appName, mailto)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	//NOTE: need to support processing one or more DOI
	for _, doi := range args {
		switch {
		case crossrefOnly:
			obj, err := apiCrossRef.Works(doi)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR (CrossRef API) %q, %s\n", doi, err)
				os.Exit(1)
			}
			if apiCrossRef.StatusCode == 200 {
				// NOTE: First we see if we can get a CrossRef record
				eprint, err := eprinttools.CrossRefWorksToEPrint(obj)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR (CrossRef to EPrintXML): skipping %q, %s\n", doi, err)
				} else {
					eprintsList.AddEPrint(eprint)
				}
			} else {
				fmt.Fprintf(os.Stderr, "WARNING (CrossRef API): %q, %s\n", doi, apiCrossRef.Status)
			}
		case dataciteOnly:
			obj, err := apiDataCite.Works(doi)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR (DataCite API): %q, %s\n", doi, err)
				os.Exit(1)
			}
			if apiDataCite.StatusCode == 200 {
				eprint, err := eprinttools.DataCiteWorksToEPrint(obj)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR (DataCite to EPrintXML): skipping %q, %s\n", doi, err)
				} else {
					eprintsList.AddEPrint(eprint)
				}
			} else {
				fmt.Fprintf(os.Stderr, "WARNING (DataCite API): %q, %s\n", doi, apiDataCite.Status)
			}
		default:
			// NOTE: just done for readability for flagging failed lookups
			isCrossRefDOI := false
			isDataCiteDOI := false

			obj, err := apiCrossRef.Works(doi)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR (CrossRef API): %q, %s\n", doi, err)
				os.Exit(1)
			}
			if apiCrossRef.StatusCode == 200 {
				isCrossRefDOI = true
				// NOTE: First we see if we can get a CrossRef record
				eprint, err := eprinttools.CrossRefWorksToEPrint(obj)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR (CrossRef to EPrintXML): skipping %q, %s\n", doi, err)
				} else {
					eprintsList.AddEPrint(eprint)
				}
			}

			// NOTE: We try DataCite's API as a fallback when CrossRef fials...
			if isCrossRefDOI == false {
				obj, err := apiDataCite.Works(doi)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR (DataCite API): %q, %s\n", doi, err)
					os.Exit(1)
				}
				if apiDataCite.StatusCode == 200 {
					isDataCiteDOI = true
					eprint, err := eprinttools.DataCiteWorksToEPrint(obj)
					if err != nil {
						fmt.Fprintf(os.Stderr, "ERROR (DataCite to EPrintXML): skipping %q, %s\n", doi, err)
					} else {
						eprintsList.AddEPrint(eprint)
					}
				}
			}
			if isCrossRefDOI == false && isDataCiteDOI == false {
				fmt.Fprintf(os.Stderr, "WARNING: %s not found in CrossRef or DataCite API lookup", doi)
			}
		}
	}
	// NOTE: We have an option to apply Caltech Library Special Rules
	// before marshaling our results...
	if useCaltechLibrarySpecificRules {
		eprintsList, err = clsrules.Apply(eprintsList)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	}
	if use_1_0_0_CaltechLibrarySpecificRules {
		eprintsList, err = clsrules.Apply1_0_0(eprintsList)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	}
	if asJSON {
		src, err := json.MarshalIndent(eprintsList, "", "   ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "%s\n", src)
		os.Exit(0)
	}
	src, err := xml.MarshalIndent(eprintsList, "", "   ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "%s\n", src)
}
