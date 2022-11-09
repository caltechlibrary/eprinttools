// Package eprinttools is a collection of structures and functions for working with the EPrints XML and EPrints REST API
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
// eputil is a command line tool for interacting with the EPrints REST API. Currently it supports harvesting REST API content as JSON or XML.
//

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	// Caltech Library Packages
	"github.com/caltechlibrary/eprinttools"
)

var (
	description = `
USAGE
	{app_name} [OPTIONS] EPRINT_REST_URL

SYNOPSIS

{app_name} is a command line program for exploring
the EPrints Extended API (provided by ep3apid) or EPrint's
own REST API.  Records are returned in either JSON or EPrints XML.
Lists of eprint ids are returned in JSON.

DETAIL FOR EPrints Extended API

The extended API is expected to be present on the local machine
at http://localhost:8484.  {app_name} will convert the command line
parameters into the appropriate URL encoding the command line as
necessary and return the values from the Extended API end points.

The format of the command working with the EPrints extended API
is ` + "`" + `{app_name} REPO_ID END_POINT_NAME [PARAMETER ...]` + "`" + `
You must specify the repository id in the command. E.g.

    {app_name} caltechauthors keys
	{app_name} caltechauthors doi
	{app_name} caltechauthors doi "10.5062/F4NP22DV"
	{app_name} caltechauthors creator-name "Morrell" "Thomas"
	{app_name} caltechauthors grant-number
	{app_name} caltechauthors grant-number "kzcx3-sw-147"
	{app_name} caltechauthors eprint 18339
	{app_name} -json caltechauthors eprint 18339

See website for a full list of available end points.

    https://caltechlibrary.github.io/eprinttools/docs/ep3apid.html

DETAIL FOR EPrints REST API

{app_name} parses XML content retrieved from
EPrints 3.x. REST API. It will render
results in JSON or XML.  With the ` + "`" + `-raw` + "`" + `
option you can get the unmodified EPrintXML from the
REST API otherwise the XML is parsed before final
rendering as JSON or XML. It requires a basic knowledge
of the layout of EPrint 3.x's REST API. It supports
both unauthenticated and basic authentication access
to the API. The REST API authentication mechanism
appears indepent of the primary website authentication
setup of the installed EPrints (at least at Caltech
Library). See the examples to start exploring the API.
`

	examples = `
Fetch the raw unmarshaled EPrint XML via the
EPrint REST API for id 123.

    {app_name} -raw https://example.org/rest/eprint/123.xml

Fetch the EPrint XML marshaled as XML using the
EPrints REST API for id 123.

    {app_name} https://example.org/rest/eprint/123.xml

Fetch the EPrint XML marshaled as JSON using the
EPrints REST API for id 123.

    {app_name} -json https://example.org/rest/eprint/123.xml

Get a JSON array of eprint ids from the REST API

    {app_name} -json https://example.org/rest/eprint/

Get the last modified date for id 123 from REST API

    {app_name} -raw https://example.org/rest/eprint/123/lastmod.txt

If the EPrint REST API is protected by basic authentication
you can pass the username and password via command line
options. You will be prompted for the password value.
or via the URL.  In this example the username is
"user" and password is "secret". In this example you will
be prompted to enter your secret.

    {app_name} -username=user -password \
      https://example.org/rest/eprint/123.xml

You can also pass the username and secret via the URL
but this leaves you vunerable to the password being recorded
in your command history or if another person has access to
the process table. You SHOULD NOT use this approach on a
shared machine!

    {app_name} https://user:secret@example.org/rest/eprint/123.xml

Getting IDs doesn't typically require authentication but seeing
specific records may depending on the roles and security
setup implemented in the EPrint instance.

Supported Environment Variables

    EPRINT_USER     sets the default username used by eputil
	EPRINT_PASSWORD sets the default password used by eputil
	EPRINT_BASE_URL sets the default base URL to access the
	                EPrints REST API

`

	// Standard Options
	showHelp    bool
	showLicense bool
	showVersion bool
	newLine     bool
	quiet       bool
	verbose     bool
	inputFName  string
	outputFName string

	// App Options
	username       string
	passwordPrompt bool
	password       string
	auth           string
	asJSON         bool
	raw            bool
	getURL         string
	getDocument    bool
	asSimplified   bool
)

func main() {
	var (
		err error
	)
	appName := path.Base(os.Args[0])
	username = os.Getenv("EPRINT_USER")
	password = os.Getenv("EPRINT_PASSWORD")
	getURL = os.Getenv("EPRINT_BASE_URL")
	// Make sure our EPRINT_BASE_URL includes the rest path if set.
	if getURL != "" && strings.HasSuffix(getURL, "/rest/eprint") == false {
		getURL = fmt.Sprintf("%s/rest/eprint", getURL)
	}

	flagSet := flag.NewFlagSet(appName, flag.ContinueOnError)

	// Standard Options
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.BoolVar(&showLicense, "license", false, "display license")
	flagSet.BoolVar(&showVersion, "version", false, "display version")
	flagSet.BoolVar(&verbose, "verbose", false, "verbose output")
	flagSet.StringVar(&inputFName, "i", "", "input file name (read the URL connection string from the input file")
	flagSet.StringVar(&inputFName, "input", "", "input file name (read the URL connection string from the input file")
	flagSet.StringVar(&outputFName, "o", "", "output file name")
	flagSet.StringVar(&outputFName, "output", "", "output file name")
	flagSet.BoolVar(&quiet, "quiet", false, "suppress error messages")
	flagSet.BoolVar(&newLine, "nl", false, "if true add a trailing newline")
	flagSet.BoolVar(&newLine, "newline", false, "if true add a trailing newline")

	// App Options
	flagSet.BoolVar(&raw, "raw", false, "get the raw EPrint REST API response")
	flagSet.BoolVar(&asJSON, "json", false, "attempt to parse XML into generaic JSON structure")
	flagSet.StringVar(&username, "u", username, "set the username for authenticated access")
	flagSet.StringVar(&username, "un", username, "set the username for authenticated access")
	flagSet.StringVar(&username, "user", username, "set the username for authenticated access")
	flagSet.StringVar(&username, "username", username, "set the username for authenticated access")
	flagSet.BoolVar(&passwordPrompt, "password", false, "Prompt for the password for authenticated access")
	flagSet.StringVar(&auth, "auth", "basic", "set the authentication type for access, default is basic")
	flagSet.BoolVar(&getDocument, "document", false, "Retrieve a document from the provided url")
	flagSet.BoolVar(&asSimplified, "s", false, "Return the object in a simplified JSON data structure.")
	flagSet.BoolVar(&asSimplified, "simple", false, "Return the object in a simplified JSON data structure.")

	// We're ready to process args
	flagSet.Parse(os.Args[1:])
	args := flagSet.Args()

	if len(args) > 0 && strings.Contains(args[0], "://") {
		getURL = args[0]
	} else {
		getURL = ""
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
		eprinttools.DisplayUsage(out, appName, flagSet, description, examples)
		os.Exit(0)
	}
	if showLicense {
		eprinttools.DisplayLicense(out, appName)
		os.Exit(0)
	}
	if showVersion {
		eprinttools.DisplayVersion(out, appName)
		os.Exit(0)
	}

	options := map[string]bool{
		"newLine":        newLine,
		"passwordPrompt": passwordPrompt,
		"getDocument":    getDocument,
		"asJSON":         asJSON,
		"asSimplified":   asSimplified,
		"verbose":        verbose,
	}
	if getURL == "" {
		os.Exit(eprinttools.RunExtendedAPIClient(out, args, asJSON, verbose))
	} else {
		os.Exit(eprinttools.RunEPrintsRESTClient(out, getURL, auth, username, password, options))
	}
}
