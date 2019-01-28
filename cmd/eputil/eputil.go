//
// eputil is a command line tool for interacting with the EPrints REST API. Currently it supports harvesting REST API content as JSON or XML.
//
// @author R. S. Doiel, <rsdoiel@library.caltech.edu>
//
// Copyright (c) 2018, Caltech
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
	"net/http"
	"net/url"
	"os"
	"strings"

	// Caltech Library Packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/eprinttools"
)

var (
	synopsis = []byte(`
_eputil_ is a command line program for exploring 
EPrint REST API and EPrint XML document structure
in XML as well as JSON.
`)
	description = []byte(`
_eputil_ parses XML content retrieved from 
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
`)

	examples = []byte(`
Fetch the raw unmarshaled EPrint XML via the 
EPrint REST API for id 123.

` + "```" + `
    eputil -raw https://example.org/rest/eprint/123.xml
` + "```" + `

Fetch the EPrint XML marshaled as XML using the 
EPrints REST API for id 123.

` + "```" + `
    eputil https://example.org/rest/eprint/123.xml 
` + "```" + `

Fetch the EPrint XML marshaled as JSON using the
EPrints REST API for id 123.

` + "```" + `
    eputil -json https://example.org/rest/eprint/123.xml
` + "```" + `

Get a JSON array of eprint ids from the REST API

` + "```" + `
    eputil -json https://example.org/rest/eprint/ 
` + "```" + `

Get the last modified date for id 123 from REST API

` + "```" + `
    eputil -raw https://example.org/rest/eprint/123/lastmod.txt 
` + "```" + `

If the EPrint REST API is protected by basic auth then
you can pass the username and password via the URL.
In this example the username is "user" and password is
"secret".

` + "```" + `
    eputil https://user:secret@example.org/rest/eprint/123.xml
` + "```" + `

Getting IDs doesn't typically require authentication but seeing
specific records may depending on the roles and security
setup implemented in the EPrint instance.

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
	//inputFName       string
	outputFName string

	// App Options
	username string
	password string
	auth     string
	asJSON   bool
	raw      bool
	getURL   string
)

func main() {
	var (
		src []byte
		err error
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
	//app.StringVar(&inputFName, "i,input", "", "input file name (read the URL connection string from the input file")
	app.StringVar(&outputFName, "o,output", "", "output file name")
	app.BoolVar(&quiet, "quiet", false, "suppress error messages")
	app.BoolVar(&newLine, "nl,newline", false, "if true add a trailing newline")
	app.BoolVar(&generateMarkdown, "generate-markdown", false, "generate Markdown documentation")
	app.BoolVar(&generateManPage, "generate-manpage", false, "generate man page")

	// App Options
	app.BoolVar(&raw, "raw", false, "get the raw EPrint REST API response")
	app.BoolVar(&asJSON, "json", false, "attempt to parse XML into generaic JSON structure")
	app.StringVar(&username, "u,un,user,username", "", "set the username for authenticated access")
	app.StringVar(&password, "pw,password", "", "set the password for authenticated access")
	app.StringVar(&auth, "auth", "", "set the authentication type for access")

	// We're ready to process args
	app.Parse()
	args := app.Args()

	if len(args) > 0 {
		getURL = args[0]
	}

	// Setup IO
	app.Eout = os.Stderr
	/*
		if getURL == "" {
			app.In, err = cli.Open(inputFName, os.Stdin)
			cli.ExitOnError(app.Eout, err, quiet)
			defer cli.CloseFile(inputFName, app.In)
		}
	*/

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

	if getURL == "" {
		app.Usage(app.Eout)
		os.Exit(1)
	}

	u, err := url.Parse(getURL)
	if err != nil {
		fmt.Fprintf(app.Eout, "%s\n", err)
		os.Exit(1)
	}
	if userinfo := u.User; userinfo != nil {
		username = userinfo.Username()
		if secret, isSet := userinfo.Password(); isSet {
			fmt.Printf("DEBUG is secret URL encoded? %q\n", secret)
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
	req.Header.Set("User-Agent", app.Version())
	client := &http.Client{}
	res, err := client.Do(req)
	cli.ExitOnError(app.Eout, err, quiet)
	defer res.Body.Close()
	if res.StatusCode == 200 {
		src, err = ioutil.ReadAll(res.Body)
		cli.ExitOnError(app.Eout, err, quiet)
	} else {
		cli.ExitOnError(app.Eout, fmt.Errorf("%s for %s", res.Status, getURL), quiet)
	}
	if len(bytes.TrimSpace(src)) == 0 {
		os.Exit(0)
	}
	if raw {
		if newLine {
			fmt.Fprintf(app.Out, "%s\n", src)
		} else {
			fmt.Fprintf(app.Out, "%s", src)
		}
		os.Exit(0)
	}

	switch {
	case u.Path == "/rest/eprint/":
		data := eprinttools.EPrintsDataSet{}
		err = xml.Unmarshal(src, &data)
		if asJSON {
			src, err = json.MarshalIndent(data, "", "   ")
		} else {
			fmt.Fprintf(app.Out, "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n")
			src, err = xml.MarshalIndent(data, "", "  ")
		}
		cli.ExitOnError(app.Eout, err, quiet)
	default:
		data := eprinttools.EPrints{}
		err = xml.Unmarshal(src, &data)
		cli.ExitOnError(app.Eout, err, quiet)
		if asJSON {
			src, err = json.MarshalIndent(data, "", "   ")
		} else {
			fmt.Fprintf(app.Out, "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n")
			src, err = xml.MarshalIndent(data, "", "  ")
		}
		cli.ExitOnError(app.Eout, err, quiet)
	}

	if newLine {
		fmt.Fprintf(app.Out, "%s\n", src)
	} else {
		fmt.Fprintf(app.Out, "%s", src)
	}
	os.Exit(0)
}
