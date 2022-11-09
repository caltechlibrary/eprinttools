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
package eprinttools

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	// Golang optional libraries
	"golang.org/x/crypto/ssh/terminal"
)

func DisplayLicense(out io.Writer, appName string) {
	fmt.Fprintf(out, "%s\n", LicenseText)
	fmt.Fprintf(out, "%s %s\n", appName, Version)
}

func DisplayVersion(out io.Writer, appName string) {
	fmt.Fprintf(out, "%s %s\n", appName, Version)
}

func DisplayUsage(out io.Writer, appName string, flagSet *flag.FlagSet, description string, examples string) {
	// Convert {app_name} and {version} in description
	if description != "" {
		fmt.Fprintf(out, strings.ReplaceAll(description, "{app_name}", appName))
	}
	flagSet.SetOutput(out)
	flagSet.PrintDefaults()

	if examples != "" {
		fmt.Fprintf(out, strings.ReplaceAll(examples, "{app_name}", appName))
	}
	DisplayLicense(out, appName)
}

func RunEPrintsRESTClient(out io.Writer, getURL string, auth string, username string, secret string, options map[string]bool) int {
	var (
		password                                                        string
		src                                                             []byte
		raw, passwordPrompt, getDocument, asSimplified, asJSON, newLine bool
	)
	for k, v := range options {
		switch k {
		case "raw":
			raw = v
		case "passwordPrompt":
			passwordPrompt = v
		case "getDocument":
			getDocument = v
		case "asSimplified":
			asSimplified = v
		case "asJSON":
			asJSON = v
		case "newLine":
			newLine = v
		}
	}
	u, err := url.Parse(getURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return 1
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
	// Finally check to see if we can read it from the environment.
	if username == "" {
		username = os.Getenv("EPRINT_USER")
	}
	if password == "" {
		password = os.Getenv("EPRINT_PASSWORD")
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
	req.Header.Set("User-Agent", fmt.Sprintf("eprinttools %s", Version))
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return 1
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		src, err = ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return 1
		}
	} else {
		fmt.Fprintf(os.Stderr, "%s for %s", res.Status, getURL)
		return 1
	}
	if len(bytes.TrimSpace(src)) == 0 {
		return 0
	}
	if raw {
		if newLine {
			fmt.Fprintf(out, "%s\n", src)
		} else {
			fmt.Fprintf(out, "%s", src)
		}
		return 0
	}

	switch {
	case getDocument:
		docName := path.Base(u.Path)
		err = ioutil.WriteFile(docName, src, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return 1
		}
		fmt.Fprintf(out, "retrieved %s\n", docName)
		return 0
	case u.Path == "/rest/eprint/":
		data := EPrintsDataSet{}
		err = xml.Unmarshal(src, &data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return 1
		}
		if asJSON || asSimplified {
			src, err = json.MarshalIndent(data, "", "   ")
		} else {
			fmt.Fprintf(out, "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n")
			src, err = xml.MarshalIndent(data, "", "  ")
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return 1
		}
	default:
		data := EPrints{}
		err = xml.Unmarshal(src, &data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return 1
		}
		for _, e := range data.EPrint {
			e.SyntheticFields()
		}
		if asSimplified {
			if sObj, err := CrosswalkEPrintToRecord(data.EPrint[0]); err != nil {
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
			return 1
		}
	}

	if newLine {
		fmt.Fprintf(out, "%s\n", src)
	} else {
		fmt.Fprintf(out, "%s", src)
	}
	return 0
}

func RunExtendedAPIClient(out io.Writer, args []string, asJSON bool, verbose bool) int {
	var (
		parts  []string
		src    []byte
		getURL string
	)
	parts = append(parts, `http://localhost:8484`)
	for _, arg := range args {
		if strings.Contains(arg, `?`) {
			parts = append(parts, arg)
		} else {
			parts = append(parts, url.PathEscape(arg))
		}
	}
	getURL = strings.Join(parts, "/")
	if verbose {
		fmt.Fprintf(os.Stderr, "Extended API URL: %s\n", getURL)
	}
	req, err := http.NewRequest(`GET`, getURL, nil)
	if asJSON {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("User-Agent", fmt.Sprintf("eprinttools %s", Version))
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return 1
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		src, err = ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return 1
		}
	} else {
		fmt.Fprintf(os.Stderr, "%s for %s\n", res.Status, getURL)
		return 1
	}
	if len(bytes.TrimSpace(src)) == 0 {
		return 0
	}
	fmt.Fprintf(out, "%s\n", src)
	return 0
}
