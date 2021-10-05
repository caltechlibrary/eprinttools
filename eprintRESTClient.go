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
package eprinttools

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	AuthNone = iota
	BasicAuth

	// Basic configurations that don't need to change much.
	timeoutSeconds               = 10
	maxConsecutiveFailedRequests = 10
)

func restClient(urlEndPoint string) ([]byte, error) {
	var (
		username string
		password string
		auth     string
		src      []byte
	)
	u, err := url.Parse(urlEndPoint)
	if err != nil {
		return nil, fmt.Errorf("%q, %s,", urlEndPoint, err)
	}
	username, password, auth = "", "", "basic"
	if userinfo := u.User; userinfo != nil {
		username = userinfo.Username()
		if secret, isSet := userinfo.Password(); isSet {
			password = secret
		}
	}

	// NOTE: We build our client request object so we can
	// set authentication if necessary.
	req, err := http.NewRequest("GET", urlEndPoint, nil)
	switch strings.ToLower(auth) {
	case "basic":
		req.SetBasicAuth(username, password)
	}
	req.Header.Set("User-Agent", Version)
	client := &http.Client{
		Timeout: timeoutSeconds * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		src, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("%s for %s", res.Status, urlEndPoint)
	}
	if len(bytes.TrimSpace(src)) == 0 {
		return nil, fmt.Errorf("No data")
	}
	return src, nil
}

// GetEPrint fetches a single EPrint record via the EPrint REST API.
func GetEPrint(baseURL string, eprintID int) (*EPrints, error) {
	endPoint := fmt.Sprintf("%s/rest/eprint/%d.xml", baseURL, eprintID)
	src, err := restClient(endPoint)
	if err != nil {
		return nil, err
	}
	data := new(EPrints)
	err = xml.Unmarshal(src, &data)
	if err != nil {
		return nil, err
	}
	for _, e := range data.EPrint {
		e.SyntheticFields()
	}
	return data, nil
}

// GetKeys returns a list of eprint record ids from the EPrints REST API
func GetKeys(baseURL string) ([]int, error) {
	type eprintKeyPage struct {
		XMLName xml.Name `xml:"html"`
		Anchors []string `xml:"body>ul>li>a"`
	}

	endPoint := fmt.Sprintf("%s/rest/eprint/", baseURL)
	src, err := restClient(endPoint)
	if err != nil {
		return nil, err
	}
	keysPage := new(eprintKeyPage)
	err = xml.Unmarshal(src, &keysPage)
	if err != nil {
		return nil, err
	}
	// Build a list of Unique IDs in a map, then convert unique querys to results array
	results := []int{}
	for _, val := range keysPage.Anchors {
		if strings.HasSuffix(val, ".xml") == true {
			eprintID, err := strconv.Atoi(strings.TrimSuffix(val, ".xml"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not extract eprint ID from %s\n", val)
				continue
			}
			results = append(results, eprintID)
		}
	}
	return results, nil
}
