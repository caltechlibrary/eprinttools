//
// Package eprinttools is a collection of structures and functions for working with the E-Prints REST API
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
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
package eprinttools

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	//"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	// Caltech Library packages
	"github.com/caltechlibrary/rc"
)

const (
	// Version is the revision number for this implementation of epgo
	Version = `v0.0.17`

	// LicenseText holds the string for rendering License info on the command line
	LicenseText = `
%s %s

Copyright (c) 2017, Caltech
All rights not granted herein are expressly reserved by Caltech.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.`
)

// These are our main bucket and index buckets
var (
	// Primary collection
	ePrintBucket = []byte("eprints")
)

func failCheck(err error, msg string) {
	pid := os.Getpid()
	if err != nil {
		log.Fatalf("(pid: %d) %s\n", pid, msg)
	}
}

// EPrintsAPI holds the basic connectin information to read the REST API for EPrints
type EPrintsAPI struct {
	XMLName xml.Name `json:"-"`
	// EPRINT_URL
	URL *url.URL `xml:"epgo>eprint_url" json:"eprint_url"`
	// EPRINT_DATASET
	Dataset string `xml:"epgo>dataset" json:"dataset"`
	// EPRINT_AUTH_METHOD
	AuthType int
	// EPRINT_USERNAME
	Username string
	// EPRINT_PASSWORD
	Secret string
	// SuppressNote suppresses the Note field
	SuppressNote bool
}

// Person returns the contents of eprint>creators>item>name as a struct
type Person struct {
	XMLName xml.Name `json:"-"`
	Given   string   `xml:"name>given" json:"given"`
	Family  string   `xml:"name>family" json:"family"`
	ID      string   `xml:"id,omitempty" json:"id"`

	// Customizations for Caltech Library
	ORCID string `xml:"orcid,omitempty" json:"orcid,omitempty"`
	//EMail string `xml:"email,omitempty" json:"email,omitempty"`
	Role string `xml:"role,omitempty" json:"role,omitempty"`
}

// PersonList is an array of pointers to Person structs
type PersonList []*Person

// RelatedURL is a structure containing information about a relationship
type RelatedURL struct {
	XMLName     xml.Name `json:"-"`
	URL         string   `xml:"url" json:"url"`
	Type        string   `xml:"type" json:"type"`
	Description string   `xml:"description" json:"description"`
}

// NumberingSystem is a structure describing other numbering systems for record
type NumberingSystem struct {
	XMLName xml.Name `json:"-"`
	Name    string   `xml:"name" json:"name"`
	ID      string   `xml:"id" json:"id"`
}

// Funder is a structure describing a funding source for record
type Funder struct {
	XMLName     xml.Name `json:"-"`
	Agency      string   `xml:"agency" json:"agency"`
	GrantNumber string   `xml:"grant_number,omitempty" json:"grant_number"`
}

// FunderList is an array of pointers to Funder structs
type FunderList []*Funder

// File structures in Document
type File struct {
	XMLName   xml.Name `json:"-"`
	ID        string   `xml:"id,attr" json:"id"`
	FileID    int      `xml:"fileid" json:"fileid"`
	DatasetID string   `xml:"datasetid" json:"datasetid"`
	ObjectID  int      `xml:"objectid" json:"objectid"`
	Filename  string   `xml:"filename" json:"filename"`
	MimeType  string   `xml:"mime_type" json:"mime_type"`
	Hash      string   `xml:"hash" json:"hash"`
	HashType  string   `xml:"hash_type" json:"hash_type"`
	FileSize  int      `xml:"filesize" json:"filesize"`
	MTime     string   `xml:"mtime" json:"mtime"`
	URL       string   `xml:"url" json:"url"`
}

// Document structures inside a Record (i.e. <eprint>...<documents><document>...</document>...</documents>...</eprint>)
type Document struct {
	XMLName    xml.Name `json:"-"`
	XMLNS      string   `xml:"xmlns,attr,omitempty" json:"name_space,omitempty"`
	ID         string   `xml:"id,attr" json:"id"`
	DocID      int      `xml:"docid" json:"doc_id"`
	RevNumber  int      `xml:"rev_number" json:"rev_number,omitempty"`
	Files      []*File  `xml:"files>file" json:"files,omitempty"`
	EPrintID   int      `xml:"eprintid" json:"eprint_id"`
	Pos        int      `xml:"pos" json:"pos,omitempty"`
	Placement  int      `xml:"placement" json:"placement,omitempty"`
	MimeType   string   `xml:"mime_type" json:"mime_type"`
	Format     string   `xml:"format" json:"format"`
	FormatDesc string   `xml:"formatdesc,omitempty" json:"format_desc,omitempty"`
	Language   string   `xml:"language" json:"language"`
	Security   string   `xml:"security" json:"security"`
	License    string   `xml:"license" json:"license"`
	Main       string   `xml:"main" json:"main"`
	Content    string   `xml:"content" json:"content"`
	Relation   []*Item  `xml:"relation>item,omitempty" json:"relation,omitempty"`
}

// DocumentList is an array of pointers to Document structs
type DocumentList []*Document

type ePrintIDs struct {
	XMLName xml.Name `xml:"html" json:"-"`
	IDs     []string `xml:"body>ul>li>a" json:"ids"`
}

func normalizeDate(in string) string {
	var (
		x   int
		err error
	)
	parts := strings.Split(in, "-")
	if len(parts) == 1 {
		parts = append(parts, "01")
	}
	if len(parts) == 2 {
		parts = append(parts, "01")
	}
	for i := 0; i < len(parts); i++ {
		x, err = strconv.Atoi(parts[i])
		if err != nil {
			x = 1
		}
		if i == 0 {
			parts[i] = fmt.Sprintf("%0.4d", x)
		} else {
			parts[i] = fmt.Sprintf("%0.2d", x)
		}
	}
	return strings.Join(parts, "-")
}

// Pick the first element in an array of strings
func first(s []string) string {
	if len(s) > 0 {
		return s[0]
	}
	return ""
}

// Pick the second element in an array of strings
func second(s []string) string {
	if len(s) > 1 {
		return s[1]
	}
	return ""
}

// Pick the list element in an array of strings
func last(s []string) string {
	l := len(s) - 1
	if l >= 0 {
		return s[l]
	}
	return ""
}

// New creates a new API instance
func New(eprintURL, datasetName string, suppressNote bool, authMethod, userName, userSecret string) (*EPrintsAPI, error) {
	var err error

	// Setup required
	api := new(EPrintsAPI)
	api.SuppressNote = suppressNote
	if eprintURL == "" {
		eprintURL = "http://localhost:8080"
	}
	api.URL, err = url.Parse(eprintURL)
	if err != nil {
		return nil, fmt.Errorf("eprint url is malformed %s, %s", eprintURL, err)
	}
	if datasetName == "" {
		return nil, fmt.Errorf("Must have a non-empty dataset collection name")
	}
	api.Dataset = datasetName

	// Handle Optional authentication settings
	switch authMethod {
	case "basic":
		api.AuthType = rc.BasicAuth
	case "oauth":
		api.AuthType = rc.OAuth
	case "shib":
		api.AuthType = rc.Shibboleth
	default:
		api.AuthType = rc.AuthNone
	}
	api.Username = userName
	api.Secret = userSecret
	return api, nil
}

// ListEPrintsURI returns a list of eprint record ids from the EPrints REST API
func (api *EPrintsAPI) ListEPrintsURI() ([]string, error) {
	var (
		results []string
	)

	workingURL, err := url.Parse(api.URL.String())
	if err != nil {
		return nil, err
	}
	if workingURL.Path == "" {
		workingURL.Path = path.Join("rest", "eprint") + "/"
	} else {
		p := api.URL.Path
		workingURL.Path = path.Join(p, "rest", "eprint") + "/"
	}
	// Switch to use Rest Client Wrapper
	rest, err := rc.New(workingURL.String(), api.AuthType, api.Username, api.Secret)
	if err != nil {
		return nil, fmt.Errorf("requesting %s, %s", workingURL.String(), err)
	}
	content, err := rest.Request("GET", workingURL.Path, map[string]string{})
	if err != nil {
		return nil, fmt.Errorf("requested %s, %s", workingURL.String(), err)
	}
	eIDs := new(ePrintIDs)
	err = xml.Unmarshal(content, &eIDs)
	if err != nil {
		return nil, err
	}
	// Build a list of Unique IDs in a map, then convert unique querys to results array
	m := make(map[string]bool)
	for _, val := range eIDs.IDs {
		if strings.HasSuffix(val, ".xml") == true {
			uri := "/" + path.Join("rest", "eprint", val)
			if _, hasID := m[uri]; hasID == false {
				// Save the new ID found
				m[uri] = true
				// Only store Unique IDs in result
				results = append(results, uri)
			}
		}
	}
	return results, nil
}

// ListModifiedEPrintURI return a list of modifed EPrint URI (eprint_ids) in start and end times
func (api *EPrintsAPI) ListModifiedEPrintURI(start, end time.Time, verbose bool) ([]string, error) {
	var (
		results []string
	)

	pid := os.Getpid()
	now := time.Now()
	t0 := now
	t1 := now
	if verbose == true {
		log.Printf("(pid: %d) Getting EPrints Ids", pid)
	}
	uris, err := api.ListEPrintsURI()
	if err != nil {
		return nil, err
	}
	if verbose == true {
		now = time.Now()
		log.Printf("(pid: %d) Retrieved %d ids, %s", pid, len(uris), now.Sub(t0).Round(time.Second))
	}

	workingURL, err := url.Parse(api.URL.String())
	if err != nil {
		return nil, err
	}
	if workingURL.Path == "" {
		workingURL.Path = path.Join("rest", "eprint") + "/"
	} else {
		p := workingURL.Path
		workingURL.Path = path.Join(p, "rest", "eprint") + "/"
	}

	if verbose == true {
		log.Printf("(pid: %d) Filtering EPrints ids by modification dates, %s to %s", pid, start.Format("2006-01-02"), end.Format("2006-01-02"))
	}
	total := len(uris)
	lastI := total - 1
	u := workingURL
	client := &http.Client{}
	for i, uri := range uris {
		u.Path = strings.TrimSuffix(uri, ".xml") + "/lastmod.txt"
		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", fmt.Sprintf("eprinttools %s", Version))
		if res, err := client.Do(req); err == nil {
			if buf, err := ioutil.ReadAll(res.Body); err == nil {
				res.Body.Close()
				datestring := fmt.Sprintf("%s", buf)
				if len(datestring) > 9 {
					datestring = datestring[0:10]
				}
				if dt, err := time.Parse("2006-01-02", datestring); err == nil && dt.Unix() >= start.Unix() && dt.Unix() <= end.Unix() {
					results = append(results, uri)
				}
			}
		}
		if verbose == true {
			now = time.Now()
			if i == lastI {
				log.Printf("(pid: %d) %d/%d ids checked, batch time %s, running time %s", pid, total, total, now.Sub(t1).Round(time.Second), now.Sub(t0).Round(time.Second))
				t1 = now
			} else if (i % 1000) == 0 {
				log.Printf("(pid: %d) %d/%d ids checked, batch time %s, running time %s", pid, i, total, now.Sub(t1).Round(time.Second), now.Sub(t0).Round(time.Second))
				t1 = now
			}
		}
	}
	if verbose == true {
		now = time.Now()
		log.Printf("(pid: %d) %d records in modified range, running time %s", pid, len(results), now.Sub(t0).Round(time.Second))
	}
	return results, nil
}

// GetEPrint retrieves an EPrint record via REST API
// Returns a EPrint structure, the raw XML and an error value.
func (api *EPrintsAPI) GetEPrint(uri string) (*EPrint, []byte, error) {
	workingURL, err := url.Parse(api.URL.String())
	if err != nil {
		return nil, nil, err
	}
	if workingURL.Path == "" {
		workingURL.Path = uri
	} else {
		p := api.URL.Path
		workingURL.Path = path.Join(p, uri)
	}

	// Switch to use Rest Client Wrapper
	rest, err := rc.New(workingURL.String(), api.AuthType, api.Username, api.Secret)
	if err != nil {
		return nil, nil, fmt.Errorf("requesting %s, %s", workingURL.String(), err)
	}
	content, err := rest.Request("GET", workingURL.Path, map[string]string{})
	if err != nil {
		return nil, nil, err
	}

	eprints := new(EPrints)
	err = xml.Unmarshal(content, &eprints)
	if err != nil {
		return nil, content, err
	}
	if len(eprints.EPrint) == 1 {
		if api.SuppressNote {
			eprints.EPrint[0].Note = ""
		}
		return eprints.EPrint[0], content, nil

	}
	if len(eprints.EPrint) > 1 {
		return nil, content, fmt.Errorf("Expected only one eprint for %s", uri)
	}
	return nil, content, fmt.Errorf("Expected an eprint for %s", uri)
}

// ToNames takes an array of pointers to Person and returns a list of names (family, given)
func (persons PersonList) ToNames() []string {
	var result []string

	for _, person := range persons {
		result = append(result, fmt.Sprintf("%s, %s", person.Family, person.Given))
	}
	return result
}

// ToORCIDs takes an an array of pointers to Person and returns a list of ORCID ids
func (persons PersonList) ToORCIDs() []string {
	var result []string

	for _, person := range persons {
		result = append(result, person.ORCID)
	}

	return result
}

// ToAgencies takes an array of pointers to Funders and returns a list of Agency names
func (funders FunderList) ToAgencies() []string {
	var result []string

	for _, funder := range funders {
		result = append(result, funder.Agency)
	}

	return result
}

// ToGrantNumbers takes an array of pointers to Funders and returns a list of Agency names
func (funders FunderList) ToGrantNumbers() []string {
	var result []string

	for _, funder := range funders {
		result = append(result, funder.GrantNumber)
	}

	return result
}

func (record *EPrint) PubDate() string {
	if record.DateType == "published" {
		return record.Date
	}
	return ""
}

func (person *Person) String() string {
	src, err := json.Marshal(person)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s", src)
}
