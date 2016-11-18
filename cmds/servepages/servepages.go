//
// Package epgo is a collection of structures and functions for working with the E-Prints REST API
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2016, Caltech
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
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"text/template"

	// 3rd Party packages
	"github.com/blevesearch/bleve"

	// Caltech Library packages
	"github.com/caltechlibrary/epgo"
)

var (
	description = `
 USAGE: %s [OPTIONS]

 OVERVIEW

	%s a webserver for explosing EPrints as HTML pages,  HTML .include pages, JSON and BibTeX formats.

 OPTIONS
`
	configuration = `

 CONFIGURATION

    %s can be configured through setting the following environment
	variables-

    EPGO_BLEVE    this is the Bleve index that drives the search feature.

    EPGO_HTDOCS    this is the directory where the HTML files are written.

	EPGO_SITE_URL  this is the website URL that the public will use

    EPGO_TEMPLATE_PATH  this is the directory that contains the templates
                   used used to generate the static content of the website.
`

	license = `
%s %s

Copyright (c) 2016, Caltech
All rights not granted herein are expressly reserved by Caltech.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

* Neither the name of epgo nor the names of its
  contributors may be used to endorse or promote products derived from
  this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
`
	showHelp    bool
	showVersion bool
	showLicense bool

	htdocs       string
	dbName       string
	bleveName    string
	templatePath string
	apiURL       string
	siteURL      string

	index bleve.Index
)

// QueryOptions holds the support query terms expected in either a GET or POST
// to the webserver
type QueryOptions struct {
	// Bleve specific properties
	Explain    bool              `json:"explain"`
	FilterTerm map[string]string `json:"filter_term,omitempty"`
	Type       string            `json:"type,omitempty"`

	// Unified search form properties, works for both Basic and Advanced search
	Method string `json:"method"`
	Action string `json:"action"`

	// This holds the query fields submitted
	Q         string `json:"q"`
	QExact    string `json:"q_exact"`
	QExcluded string `json:"q_excluded"`
	QRequired string `json:"q_required"`
	Size      int    `json:"size"`
	From      int    `json:"from"`
	AllIDs    bool   `json:"all_ids"`

	// Results olds the submitted query results
	Total           int    `json:"total"`
	DetailsBaseURI  string `json:"details_base_uri"`
	QueryURLEncoded string
	Request         *bleve.SearchRequest
	Results         *bleve.SearchResult
}

// Parses the submitted map for query options setting the internals of the QueryOptions structure
func (q *QueryOptions) Parse(m map[string]interface{}) error {
	var (
		err error
	)
	// raw is a tempory data structure to sanitize the
	// form request submitted via the query.
	raw := new(QueryOptions)
	isQuery := false

	src, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("Can't marshal %+v, %s", m, err)
	}
	err = json.Unmarshal(src, &raw)
	if err != nil {
		return fmt.Errorf("Can't unmarshal %s, %s", src, err)
	}
	if len(raw.Q) > 0 {
		q.Q = raw.Q
		isQuery = true
	}
	if len(raw.QExact) > 0 {
		q.QExact = raw.QExact
		isQuery = true
	}
	if len(raw.QExcluded) > 0 {
		q.QExcluded = q.QExact
	}
	if len(raw.QRequired) > 0 {
		q.QRequired = raw.QRequired
		isQuery = true
	}

	if isQuery == false {
		return fmt.Errorf("Missing query value fields")
	}

	if raw.AllIDs == true {
		q.AllIDs = true
	}

	//Note: if q.Size is not set by the query request pick a nice default value
	if raw.Size <= 1 {
		q.Size = 10
	} else {
		q.Size = raw.Size
	}
	if raw.From < 0 {
		q.From = 0
	} else {
		q.From = raw.From
	}
	return nil
}

func uInt64ToInt(u uint64) (int, error) {
	return strconv.Atoi(fmt.Sprintf("%d", u))
}

// AttachSearchResults sets the value of the SearchResults field in the QueryOptions struct.
func (q *QueryOptions) AttachSearchResults(sr *bleve.SearchResult) {
	if sr != nil {
		q.Results = sr
		q.Total, _ = uInt64ToInt(sr.Total)
		q.Request = sr.Request
	} else {
		q.Total = 0
	}

	v := url.Values{}
	if q.AllIDs == true {
		v.Add("all_ids", "true")
	}
	v.Add("size", fmt.Sprintf("%d", q.Size))
	v.Add("from", fmt.Sprintf("%d", q.From))
	v.Add("total", fmt.Sprintf("%d", q.Total))
	v.Add("q", q.Q)
	v.Add("q_required", q.QRequired)
	v.Add("q_exact", q.QExact)
	v.Add("q_excluded", q.QExcluded)
	q.QueryURLEncoded = v.Encode()
}

func resultsHandler(w http.ResponseWriter, r *http.Request) {
	urlQuery := r.URL.Query()
	err := r.ParseForm()
	if err != nil {
		responseLogger(r, http.StatusBadRequest, err)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error in POST: %s", err)))
		return
	}

	// Collect the submissions fields.
	submission := make(map[string]interface{})
	// Basic Search results
	if r.Method == "GET" {
		for k, v := range urlQuery {
			if k == "all_ids" {
				if b, err := strconv.ParseBool(strings.Join(v, "")); err == nil {
					submission[k] = b
				}
			} else if k == "from" || k == "size" || k == "total" {
				if i, err := strconv.Atoi(strings.Join(v, "")); err == nil {
					submission[k] = i
				}
			} else if k == "q" || k == "q_exact" || k == "q_excluded" || k == "q_required" {
				submission[k] = strings.Join(v, "")
			}
		}
	}

	// Advanced Search results
	if r.Method == "POST" {
		for k, v := range r.Form {
			if k == "all_ids" {
				if b, err := strconv.ParseBool(strings.Join(v, "")); err == nil {
					submission[k] = b
				}
			} else if k == "from" || k == "size" || k == "total" {
				if i, err := strconv.Atoi(strings.Join(v, "")); err == nil {
					submission[k] = i
				}
			} else if k == "q" || k == "q_exact" || k == "q_excluded" || k == "q_required" {
				submission[k] = strings.Join(v, "")
			}
		}
	}

	q := new(QueryOptions)
	err = q.Parse(submission)
	if err != nil {
		responseLogger(r, http.StatusBadRequest, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}

	//
	// Note: Adding logic to handle basic and advanced search...
	//
	// q           NewQueryStringQuery
	// q_required  NewQueryStringQuery with a + prefix for each strings.Fields(q_required) value
	// q_exact     NewMatchPhraseQuery
	// q_excluded  NewQueryStringQuery with a - prefix for each strings.Feilds(q_excluded) value
	//
	var conQry []bleve.Query

	if q.Q != "" {
		conQry = append(conQry, bleve.NewQueryStringQuery(q.Q))
	}
	if q.QExact != "" {
		conQry = append(conQry, bleve.NewMatchPhraseQuery(q.QExact))
	}
	var terms []string
	for _, s := range strings.Fields(q.QRequired) {
		terms = append(terms, fmt.Sprintf("+%s", strings.TrimSpace(s)))
	}
	for _, s := range strings.Fields(q.QExcluded) {
		terms = append(terms, fmt.Sprintf("-%s", strings.TrimSpace(s)))
	}
	if len(terms) > 0 {
		qString := strings.Join(terms, " ")
		conQry = append(conQry, bleve.NewQueryStringQuery(qString))
	}

	qry := bleve.NewConjunctionQuery(conQry)
	if q.Size == 0 {
		q.Size = 10
	}
	searchRequest := bleve.NewSearchRequestOptions(qry, q.Size, q.From, q.Explain)
	if searchRequest == nil {
		responseLogger(r, http.StatusBadRequest, fmt.Errorf("Can't build new search request options %+v, %s", qry, err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}

	searchRequest.Highlight = bleve.NewHighlight()
	searchRequest.Highlight.AddField("title")
	searchRequest.Highlight.AddField("content_description")
	searchRequest.Highlight.AddField("subjects")
	searchRequest.Highlight.AddField("subjects_function")
	searchRequest.Highlight.AddField("subjects_topical")
	searchRequest.Highlight.AddField("extents")

	subjectFacet := bleve.NewFacetRequest("subjects", 3)
	searchRequest.AddFacet("subjects", subjectFacet)

	subjectTopicalFacet := bleve.NewFacetRequest("subjects_topical", 3)
	searchRequest.AddFacet("subjects_topical", subjectTopicalFacet)

	subjectFunctionFacet := bleve.NewFacetRequest("subjects_function", 3)
	searchRequest.AddFacet("subjects_function", subjectFunctionFacet)

	// Return all fields
	searchRequest.Fields = []string{}

	searchResults, err := index.Search(searchRequest)
	if err != nil {
		responseLogger(r, http.StatusInternalServerError, fmt.Errorf("Bleve results error %v, %s", qry, err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}

	// q (QueryOptions) performs double duty as both the structure for query submission as well
	// as carring the results to support paging and other types of navigation through
	// the query set. Results are a query with the bleve.SearchReults merged
	q.AttachSearchResults(searchResults)
	pageHTML := path.Join(templatePath, "results.html")
	pageInclude := path.Join(templatePath, "results.include")

	// Load my templates and setup to execute them
	tmpl, err := epgo.AssembleTemplate(pageHTML, pageInclude)
	if err != nil {
		responseLogger(r, http.StatusInternalServerError, fmt.Errorf("Template Errors: %s, %s, %s\n", pageHTML, pageInclude, err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Template errors: %s", err)))
		return
	}

	// Render the page
	w.Header().Set("Content-Type", "text/html")
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, q)
	if err != nil {
		responseLogger(r, http.StatusInternalServerError, fmt.Errorf("Can't render %s, %s, %s", pageHTML, pageInclude, err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Template error"))
		return
	}
	//NOTE: This bit of ugliness is here because I need to allow <mark> elements and ellipis in the results fragments
	w.Write(bytes.Replace(bytes.Replace(bytes.Replace(buf.Bytes(), []byte("&lt;mark&gt;"), []byte("<mark>"), -1), []byte("&lt;/mark&gt;"), []byte("</mark>"), -1), []byte(`…`), []byte(`&hellip;`), -1))
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	//logRequest(r)
	// If GET with Query String or POST pass to results handler
	// else display Basic Search Form
	query := r.URL.Query()
	if r.Method == "POST" || len(query) > 0 {
		resultsHandler(w, r)
		return
	}

	// Shared form data fields for a New Search.
	formData := struct {
		URI string
	}{
		URI: "/",
	}

	// Handle the basic or advanced search form requests.
	var (
		tmpl *template.Template
		err  error
	)
	pageHTML := path.Join(templatePath, "search.html")
	pageInclude := path.Join(templatePath, "search.include")
	w.Header().Set("Content-Type", "text/html")
	if strings.HasPrefix(r.URL.Path, "/search/") == true {
		formData.URI = "/search/"
		tmpl, err = epgo.AssembleTemplate(pageHTML, pageInclude)
		if err != nil {
			fmt.Printf("Can't read search templates %s, %s, %s", pageHTML, pageInclude, err)
			return
		}
	}

	err = tmpl.Execute(w, formData)
	if err != nil {
		responseLogger(r, http.StatusInternalServerError, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//FIXME: add the response status returned.
		q := r.URL.Query()
		if len(q) > 0 {
			log.Printf("Request: %s Path: %s RemoteAddr: %s UserAgent: %s Query: %+v\n", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), q)
		} else {
			log.Printf("Request: %s Path: %s RemoteAddr: %s UserAgent: %s\n", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
		}
		next.ServeHTTP(w, r)
	})
}

func responseLogger(r *http.Request, status int, err error) {
	q := r.URL.Query()
	if len(q) > 0 {
		log.Printf("Response: %s Path: %s RemoteAddr: %s UserAgent: %s Query: %+v Status: %d, %s %q\n", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), q, status, http.StatusText(status), err)
	} else {
		log.Printf("Response: %s Path: %s RemoteAddr: %s UserAgent: %s Status: %d, %s %q\n", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), status, http.StatusText(status), err)
	}
}

// isMultiViewPath checks to see if the path requested behaves like an Apache MultiView request
func isMultiViewPath(p string) bool {
	// check to see if p plus .html extension exists
	fname := fmt.Sprintf("%s.html", p)
	if _, err := os.Stat(path.Join(htdocs, fname)); err == nil {
		return true
	}
	return false
}

func multiViewPath(p string) string {
	return fmt.Sprintf("%s.html", p)
}

func customRoutes(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/search/") == true {
			searchHandler(w, r)
			return
		}
		// NOTE: The default static file server doesn't seem send the correct mimetype for RSS and JSON responses.

		// If this is a MultiViews style request (i.e. missing .html) then update r.URL.Path
		if isMultiViewPath(r.URL.Path) == true {
			p := multiViewPath(r.URL.Path)
			r.URL.Path = p
		}
		// If we make it this far, fall back to the default handler
		next.ServeHTTP(w, r)
	})
}

func usage(appName, version string) {
	fmt.Printf(description, appName, appName)
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("\t-%s\t%s\n", f.Name, f.Usage)
	})
	fmt.Printf(configuration, appName)
	fmt.Printf("%s %s\n", appName, version)
	os.Exit(0)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	flag.BoolVar(&showHelp, "h", false, "display help")
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showVersion, "v", false, "display version")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.BoolVar(&showLicense, "l", false, "display license")
	flag.BoolVar(&showLicense, "license", false, "display license")

	flag.StringVar(&htdocs, "htdocs", "", "specify where to write the HTML files to")
	flag.StringVar(&bleveName, "bleve", "", "the Bleve index/db name")
	flag.StringVar(&siteURL, "site-url", "", "the website url")
	flag.StringVar(&templatePath, "template-path", "", "specify where to read the templates from")
}

func main() {
	var (
		cfg epgo.Config
		err error
	)

	appName := path.Base(os.Args[0])
	flag.Parse()
	if showHelp == true {
		usage(appName, epgo.Version)
	}
	if showVersion == true {
		fmt.Printf("%s %s\n", appName, epgo.Version)
		os.Exit(0)
	}
	if showLicense == true {
		fmt.Printf(license, appName, epgo.Version)
		os.Exit(0)
	}

	// Check to see we can merge the required fields are merged.
	check(cfg.MergeEnv("EPGO", "HTDOCS", htdocs))
	check(cfg.MergeEnv("EPGO", "TEMPLATE_PATH", templatePath))
	check(cfg.MergeEnv("EPGO", "SITE_URL", siteURL))
	check(cfg.MergeEnv("EPGO", "BLEVE", bleveName))

	if cfg.Htdocs != "" {
		if _, err := os.Stat(htdocs); os.IsNotExist(err) {
			os.MkdirAll(htdocs, 0775)
		}
	}
	htdocs = cfg.Get("htdocs")
	templatePath = cfg.Get("template_path")
	log.Printf("DEBUG templatePath: %q", templatePath)
	bleveName = cfg.Get("bleve")
	siteURL = cfg.Get("site_url")

	//NOTE: Need to get hostname and port from siteURL
	u, err := url.Parse(siteURL)
	check(err)

	//
	// Run the webserver and search service
	//
	log.Printf("%s %s\n", appName, epgo.Version)

	// Wake up our search engine
	log.Println("Initializing search services")
	index, err = bleve.Open(bleveName)
	if err != nil {
		log.Fatalf("Can't open Bleve index %q, %s", bleveName, err)
	}
	defer index.Close()

	// Send static file request to the default handler,
	// search routes are handled by middleware customRoutes()
	http.Handle("/", http.FileServer(http.Dir(htdocs)))

	log.Printf("Listening on %s\n", u.String())
	err = http.ListenAndServe(u.Host, requestLogger(customRoutes(http.DefaultServeMux)))
	if err != nil {
		log.Fatal(err)
	}
}
