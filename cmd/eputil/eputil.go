//
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
//
package main

//
// eputil is a command line tool for interacting with the EPrints REST API. Currently it supports harvesting REST API content as JSON or XML.
//

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	// Golang optional libraries
	"golang.org/x/crypto/ssh/terminal"

	// Caltech Library Packages
	"github.com/caltechlibrary/eprinttools"
)

var (
	description = `
USAGE
	{app_name} [OPTIONS] EPRINT_REST_URL
	
SYNOPSIS

{app_name} is a command line program for exploring 
EPrint REST API and EPrint XML document structure
in XML as well as JSON.

DETAIL

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
		src []byte
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

	// App Options
	flagSet.BoolVar(&raw, "raw", false, "get the raw EPrint REST API response")
	flagSet.BoolVar(&asJSON, "json", false, "attempt to parse XML into generaic JSON structure")
	flagSet.StringVar(&username, "u", "", "set the username for authenticated access")
	flagSet.StringVar(&username, "un", "", "set the username for authenticated access")
	flagSet.StringVar(&username, "user", "", "set the username for authenticated access")
	flagSet.StringVar(&username, "username", "", "set the username for authenticated access")
	flagSet.BoolVar(&passwordPrompt, "password", false, "Prompt for the password for authenticated access")
	flagSet.StringVar(&auth, "auth", "basic", "set the authentication type for access, default is basic")
	flagSet.BoolVar(&getDocument, "document", false, "Retrieve a document from the provided url")
	flagSet.BoolVar(&asSimplified, "s", false, "Return the object in a simplified JSON data structure.")
	flagSet.BoolVar(&asSimplified, "simple", false, "Return the object in a simplified JSON data structure.")

	// We're ready to process args
	flagSet.Parse(os.Args[1:])
	args := flagSet.Args()

	if len(args) > 0 {
		getURL = args[0]
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

	/*
		if getURL == "" {
			in, err = cli.Open(inputFName, os.Stdin)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
			defer in.Close()
		}
	*/

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

	if getURL == "" {
		eprinttools.DisplayUsage(os.Stderr, appName, flagSet, description, examples, license)
		os.Exit(1)
	}

	u, err := url.Parse(getURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	if passwordPrompt {
		fmt.Fprintf(out, "Please type the password for accessing\n%s\n", getURL)
		if src, err := terminal.ReadPassword(0); err == nil {
			password = fmt.Sprintf("%s", src)
		}
	}
	if userinfo := u.User; userinfo != nil {
		username = userinfo.Username()
		if secret, isSet := userinfo.Password(); isSet {
			password = secret
		}
		if auth == "" {
			auth = "basic"
		}
	}

	// NOTE: We build our client request object so we can
	// set authentication if necessary.
	req, err := http.NewRequest("GET", getURL, nil)
	switch strings.ToLower(auth) {
	case "basic":
		req.SetBasicAuth(username, password)
	case "basic_auth":
		req.SetBasicAuth(username, password)
	}
	req.Header.Set("User-Agent", eprinttools.Version)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		src, err = ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Fprintf(os.Stderr, "%s for %s", res.Status, getURL)
		os.Exit(1)
	}
	if len(bytes.TrimSpace(src)) == 0 {
		os.Exit(0)
	}
	if raw {
		if newLine {
			fmt.Fprintf(out, "%s\n", src)
		} else {
			fmt.Fprintf(out, "%s", src)
		}
		os.Exit(0)
	}

	switch {
	case getDocument:
		docName := path.Base(u.Path)
		err = ioutil.WriteFile(docName, src, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(out, "retrieved %s\n", docName)
		os.Exit(0)
	case u.Path == "/rest/eprint/":
		data := eprinttools.EPrintsDataSet{}
		err = xml.Unmarshal(src, &data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		if asJSON || asSimplified {
			src, err = json.MarshalIndent(data, "", "   ")
		} else {
			fmt.Fprintf(out, "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n")
			src, err = xml.MarshalIndent(data, "", "  ")
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	default:
		data := eprinttools.EPrints{}
		err = xml.Unmarshal(src, &data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		for _, e := range data.EPrint {
			e.SyntheticFields()
		}
		if asSimplified {
			if sObj, err := eprinttools.CrosswalkEPrintToRecord(data.EPrint[0]); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
			} else {
				src, err = json.MarshalIndent(sObj, "", "   ")
			}
		} else if asJSON {
			src, err = json.MarshalIndent(data, "", "   ")
		} else {
			fmt.Fprintf(out, "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n")
			src, err = xml.MarshalIndent(data, "", "  ")
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	}

	if newLine {
		fmt.Fprintf(out, "%s\n", src)
	} else {
		fmt.Fprintf(out, "%s", src)
	}
	os.Exit(0)
}
