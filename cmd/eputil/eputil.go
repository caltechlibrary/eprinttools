//
// eputil is a command line tool for interacting with the EPrints REST API. Currently it supports harvesting REST API content as JSON or XML.
//
// @author R. S. Doiel, <rsdoiel@library.caltech.edu>
//
// Copyright (c) 2017, Caltech
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
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	// Caltech Library Packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/eprinttools"
)

var (
	description = []byte(`
	eputil parses XML content retrieved from disc or the EPrints API. It will 
	render JSON if the XML is valid otherwise return errors.
`)

	examples = []byte(`
Fetch an EPrints document as JSON from a URL for an EPrint with an id of 123

    eputil -get https://eprints.example.org/rest/eprint/123.xml -json

Fetch an EPrints document as XML from a URL for an EPrint with an id of 123

    eputil -get https://eprints.example.org/rest/eprint/123.xml

Fetch the creators.xml as JSON for an EPrint with the id of 123.

    eputil -get https://eprints.example.org/rest/eprint/123/creators.xml -json

Parse an EPrint reversion XML document

    eputil -i revision/2.xml -eprint

Get a JSON array of eprint ids from the REST API

    eputil -url https://eprints.example.org/rest/eprint/ -ids
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
	user    string
	getURL  string
	putURL  string
	postURL string
	eprints bool
	eprint  bool
	getIDs  bool
	asJSON  bool
)

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	req.Header.Add("Authorization", "Basic "+basicAuth("username1", "password123"))
	return nil
}

func main() {
	var (
		src      []byte
		err      error
		username string
		password string
	)

	app := cli.NewCli(eprinttools.Version)

	// Add Help Docs
	app.AddHelp("description", description)
	app.AddHelp("examples", examples)

	// Add Environment variables
	app.EnvStringVar(&user, "EP_USER", "", "Sets the eprint USERNAME:PASSWORD to use with URL basic auth")

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
	app.StringVar(&getURL, "get,url", "", "do an HTTP GET to fetch the XML from the URL then parse")
	app.StringVar(&postURL, "post", "", "do an HTTP POST to add/update a record (e.g. update a URL endpoint")
	app.StringVar(&putURL, "put", "", "do an HTTP POST to add/update a record (e.g. update a URL endpoint")
	app.BoolVar(&eprints, "document,eprints", false, "parse an eprints (e.g. rest response) document")
	app.BoolVar(&eprint, "revision,eprint", false, "parse a eprint (revision) document")
	app.BoolVar(&asJSON, "json", false, "attempt to parse XML into generaic JSON structure")
	app.BoolVar(&getIDs, "ids", false, "get a list of doc paths (e.g. ids or sub-fields depending on the URL provided")
	app.StringVar(&user, "u,user", "", "set the username:password for authentication authenticated access")

	// We're ready to process args
	app.Parse()
	args := app.Args()

	// Setup IO
	app.Eout = os.Stderr
	if getURL == "" {
		app.In, err = cli.Open(inputFName, os.Stdin)
		cli.ExitOnError(app.Eout, err, quiet)
		defer cli.CloseFile(inputFName, app.In)
	}

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

	if user != "" {
		if strings.Contains(user, ":") {
			p := strings.SplitN(user, ":", 2)
			username, password = p[0], p[1]
		} else {
			username = user
		}
	}

	if getURL == "" {
		src, err = ioutil.ReadAll(app.In)
		cli.ExitOnError(app.Eout, err, quiet)
	} else {
		// NOTE: We build our client request object so we can
		// set authentication if necessary.
		req, err := http.NewRequest("GET", getURL, nil)
		if user != "" {
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
	}
	if len(bytes.TrimSpace(src)) == 0 {
		os.Exit(0)
	}

	if postURL == "" && putURL == "" {
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
		case asJSON:
			data := eprinttools.Generic{}
			err = xml.Unmarshal(src, &data)
			cli.ExitOnError(app.Eout, err, quiet)

			src, err = json.MarshalIndent(data, "", " ")
			cli.ExitOnError(app.Eout, err, quiet)
		case getIDs:
			data := eprinttools.EPrintsDataSet{}
			err = xml.Unmarshal(src, &data)
			cli.ExitOnError(app.Eout, err, quiet)
			src, err = json.MarshalIndent(data, "", " ")
		default:
			// Don't do anything, just return the raw XML
		}

		if newLine {
			fmt.Fprintf(os.Stdout, "%s\n", src)
		} else {
			fmt.Fprintf(os.Stdout, "%s", src)
		}
		os.Exit(0)
	} else {
		switch {
		case putURL != "":
			fmt.Fprintf(app.Out, "PUT: %s\nDATA: %q\n", putURL, src)
			// NOTE: We build our client request object so we can
			// set authentication if necessary.
			req, err := http.NewRequest("PUT", putURL, strings.NewReader(fmt.Sprintf("%s", src)))
			if username != "" {
				fmt.Printf("DEBUG username %q, password %q\n", username, password)
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
				if newLine {
					fmt.Fprintf(app.Out, "%s\n", src)
				} else {
					fmt.Fprintf(app.Out, "%s", src)
				}
			} else {
				cli.ExitOnError(app.Eout, fmt.Errorf("%s for %s", res.Status, putURL), quiet)
			}
		case postURL != "":
			fmt.Fprintf(app.Out, "POST: %s\nDATA: %q\n", postURL, src)
			// NOTE: We build our client request object so we can
			// set authentication if necessary.
			req, err := http.NewRequest("POST", postURL, strings.NewReader(fmt.Sprintf("%s", src)))
			if user != "" {
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
				if newLine {
					fmt.Fprintf(app.Out, "%s\n", src)
				} else {
					fmt.Fprintf(app.Out, "%s", src)
				}
			} else {
				cli.ExitOnError(app.Eout, fmt.Errorf("%s for %s", res.Status, postURL), quiet)
			}
		default:
			fmt.Fprintf(app.Eout, "HTTP Method not implemented.")
		}
	}
}
