//
// Package epgo is a collection of structures and functions for working with the E-Prints REST API
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
package epgo

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"text/template"
	"time"

	// Caltech Library packages
	"github.com/caltechlibrary/bibtex"
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/dataset"
	"github.com/caltechlibrary/tmplfn"
)

// These are our main bucket and index buckets
var (
	// Primary collection
	ePrintBucket = []byte("eprints")

	// Select lists delimiter
	indexDelimiter = "|"
	// expected select lists used by epgo
	slNames = []string{
		"pubDate",
		"localGroup",
		"orcid",
		"isni",
		"funder",
		"grantNumber",
	}

	// TmplFuncs is the collected functions available in EPGO templates
	TmplFuncs = tmplfn.Join(tmplfn.TimeMap, tmplfn.PageMap)
)

func failCheck(err error, msg string) {
	if err != nil {
		log.Fatalf("%s\n", msg)
	}
}

// EPrintsAPI holds the basic connectin information to read the REST API for EPrints
type EPrintsAPI struct {
	XMLName        xml.Name `json:"-"`
	URL            *url.URL `xml:"epgo>api_url" json:"api_url"`                 // EPGO_API_URL
	Dataset        string   `xml:"epgo>dataset" json:"dataset"`                 // EPGO_DATASET
	BleveName      string   `xml:"epgo>bleve" json:"bleve"`                     // EPGO_BLEVE
	Htdocs         string   `xml:"epgo>htdocs" json:"htdocs"`                   // EPGO_HTDOCS
	TemplatePath   string   `xml:"epgo>template_path" json:"template_path"`     // EPGO_TEMPLATES
	SiteURL        *url.URL `xml:"epgo>site_url" json:"site_url"`               // EPGO_SITE_URL
	RepositoryPath string   `xml:"epgo>repository_path" json:"repository_path"` // EPGO_REPOSITORY_PATH
}

// Person returns the contents of eprint>creators>item>name as a struct
type Person struct {
	XMLName xml.Name `json:"-"`
	Given   string   `xml:"name>given" json:"given"`
	Family  string   `xml:"name>family" json:"family"`
	ID      string   `xml:"id,omitempty" json:"id"`
	ORCID   string   `xml:"orcid,omitempty" json:"orcid"`
	ISNI    string   `xml:"isni,omitempty" json:"isni"`
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

// Document structures in Record
type Document struct {
	XMLName   xml.Name `json:"-"`
	ID        string   `xml:"id,attr" json:"id"`
	DocID     int      `xml:"docid" json:"docid"`
	RevNumber int      `xml:"rev_number" json:"rev_number"`
	Files     []*File  `xml:"files>file" json:"files"`
	EPrintID  int      `xml:"eprintid" json:"eprintid"`
	Pos       int      `xml:"pos" json:"pos"`
	Placement int      `xml:"placement" json:"placement"`
	MimeType  string   `xml:"mime_type" json:"mime_type"`
	Format    string   `xml:"format" json:"format"`
	Language  string   `xml:"language" json:"language"`
	Security  string   `xml:"security" json:"security"`
	License   string   `xml:"license" json:"license"`
	Main      string   `xml:"main" json:"main"`
	Content   string   `xml:"content" json:"content"`
}

// DocumentList is an array of pointers to Document structs
type DocumentList []*Document

// Record returns a structure that can be converted to JSON easily
type Record struct {
	XMLName              xml.Name           `json:"-"`
	Title                string             `xml:"eprint>title" json:"title"`
	URI                  string             `json:"uri"`
	Abstract             string             `xml:"eprint>abstract" json:"abstract"`
	Documents            DocumentList       `xml:"eprint>documents>document" json:"documents"`
	Note                 string             `xml:"eprint>note" json:"note"`
	ID                   int                `xml:"eprint>eprintid" json:"id"`
	RevNumber            int                `xml:"eprint>rev_number" json:"rev_number"`
	UserID               int                `xml:"eprint>userid" json:"userid"`
	Dir                  string             `xml:"eprint>dir" json:"eprint_dir"`
	Datestamp            string             `xml:"eprint>datestamp" json:"datestamp"`
	LastModified         string             `xml:"eprint>lastmod" json:"lastmod"`
	StatusChange         string             `xml:"eprint>status_changed" json:"status_changed"`
	Type                 string             `xml:"eprint>type" json:"type"`
	MetadataVisibility   string             `xml:"eprint>metadata_visibility" json:"metadata_visibility"`
	Creators             PersonList         `xml:"eprint>creators>item" json:"creators"`
	IsPublished          string             `xml:"eprint>ispublished" json:"ispublished"`
	Subjects             []string           `xml:"eprint>subjects>item" json:"subjects"`
	FullTextStatus       string             `xml:"eprint>full_text_status" json:"full_text_status"`
	Keywords             string             `xml:"eprint>keywords" json:"keywords"`
	Date                 string             `xml:"eprint>date" json:"date"`
	DateType             string             `xml:"eprint>date_type" json:"date_type"`
	Publication          string             `xml:"eprint>publication" json:"publication"`
	Volume               string             `xml:"eprint>volume" json:"volume"`
	Number               string             `xml:"eprint>number" json:"number"`
	PageRange            string             `xml:"eprint>pagerange" json:"pagerange"`
	IDNumber             string             `xml:"eprint>id_number" json:"id_number"`
	Referred             bool               `xml:"eprint>refereed" json:"refereed"`
	ISSN                 string             `xml:"eprint>issn" json:"issn"`
	OfficialURL          string             `xml:"eprint>official_url" json:"official_url"`
	RelatedURL           []*RelatedURL      `xml:"eprint>related_url>item" json:"related_url"`
	ReferenceText        []string           `xml:"eprint>referencetext>item" json:"referencetext"`
	Rights               string             `xml:"eprint>rights" json:"rights"`
	OfficialCitation     string             `xml:"eprint>official_cit" json:"official_citation"`
	OtherNumberingSystem []*NumberingSystem `xml:"eprint>other_numbering_system>item,omitempty" json:"other_numbering_system"`
	Funders              FunderList         `xml:"eprint>funders>item" json:"funders"`
	Collection           string             `xml:"eprint>collection" json:"collection"`
	Reviewer             string             `xml:"eprint>reviewer" json:"reviewer"`
	LocalGroup           []string           `xml:"eprint>local_group>item" json:"local_group"`
}

type ePrintIDs struct {
	XMLName xml.Name `xml:"html" json:"-"`
	IDs     []string `xml:"body>ul>li>a" json:"ids"`
}

func normalizeDate(in string) string {
	parts := strings.Split(in, "-")
	if len(parts) == 1 {
		parts = append(parts, "01")
		parts = append(parts, "01")
	}
	if len(parts) == 2 {
		parts = append(parts, "01")
	}
	for i := 0; i < len(parts); i++ {
		x, err := strconv.Atoi(parts[i])
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

// Turn a string into a URL friendly path part
func slugify(s string) (string, error) {
	if len(s) > 250 {
		return "", fmt.Errorf("string to long (%d), %q", len(s), s)
	}
	return url.QueryEscape(s), nil
}

// ToBibTeXElement takes an epgo.Record and turns it into a bibtex.Element record
func (rec *Record) ToBibTeXElement() *bibtex.Element {
	bib := &bibtex.Element{}
	bib.Set("type", rec.Type)
	bib.Set("id", fmt.Sprintf("eprint-%d", rec.ID))
	bib.Set("title", rec.Title)
	if len(rec.Abstract) > 0 {
		bib.Set("abstract", rec.Abstract)
	}
	if rec.DateType == "pub" {
		dt, err := time.Parse("2006-01-02", rec.Date)
		if err != nil {
			bib.Set("year", dt.Format("2006"))
			bib.Set("month", dt.Format("January"))
		}
	}
	if len(rec.PageRange) > 0 {
		bib.Set("pages", rec.PageRange)
	}
	/*
		if len(rec.Note) > 0 {
			bib.Set("note", rec.Note)
		}
	*/
	if len(rec.Creators) > 0 {
		people := []string{}
		for _, person := range rec.Creators {
			people = append(people, fmt.Sprintf("%s, %s", person.Family, person.Given))
		}
		bib.Set("author", strings.Join(people, " and "))
	}
	switch rec.Type {
	case "article":
		bib.Set("journal", rec.Publication)
	case "book":
		bib.Set("publisher", rec.Publication)
	}
	if len(rec.Volume) > 0 {
		bib.Set("volume", rec.Volume)
	}
	if len(rec.Number) > 0 {
		bib.Set("number", rec.Number)
	}
	return bib
}

// New creates a new API instance
func New(cfg *cli.Config) (*EPrintsAPI, error) {
	var err error
	apiURL := cfg.Get("api_url")
	siteURL := cfg.Get("site_url")
	htdocs := cfg.Get("htdocs")
	datasetName := cfg.Get("dataset")
	bleveName := cfg.Get("bleve")
	templatePath := cfg.Get("template_path")
	repositoryPath := cfg.Get("repository_path")

	if apiURL == "" {
		return nil, fmt.Errorf("Environment not configured, missing eprint api url")
	}
	api := new(EPrintsAPI)
	api.URL, err = url.Parse(apiURL)
	if err != nil {
		return nil, fmt.Errorf("api url is malformed %s, %s", apiURL, err)
	}
	api.SiteURL, err = url.Parse(siteURL)
	if err != nil {
		return nil, fmt.Errorf("site url malformed %s, %s", siteURL, err)
	}
	if htdocs == "" {
		htdocs = "htdocs"
	}
	if datasetName == "" {
		datasetName = "eprints"
	}
	if bleveName == "" {
		bleveName = "eprints.bleve"
	}
	if templatePath == "" {
		templatePath = "templates"
	}
	if repositoryPath == "" {
		repositoryPath = "repository"
	}
	api.Htdocs = htdocs
	api.Dataset = datasetName
	api.TemplatePath = templatePath
	api.BleveName = bleveName
	api.RepositoryPath = repositoryPath
	return api, nil
}

type byURI []string

func (s byURI) Len() int {
	return len(s)
}

func (s byURI) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byURI) Less(i, j int) bool {
	s1 := strings.TrimSuffix(path.Base(s[i]), path.Ext(s[i]))
	s2 := strings.TrimSuffix(path.Base(s[j]), path.Ext(s[j]))
	a1, err := strconv.Atoi(s1)
	if err != nil {
		return false
	}
	a2, err := strconv.Atoi(s2)
	if err != nil {
		return false
	}
	return a1 > a2
}

// ListEPrintsURI returns a list of eprint record ids from the EPrints REST API
func (api *EPrintsAPI) ListEPrintsURI() ([]string, error) {
	var (
		results []string
	)

	api.URL.Path = path.Join("rest", "eprint") + "/"
	resp, err := http.Get(api.URL.String())
	if err != nil {
		return nil, fmt.Errorf("requested %s, %s", api.URL.String(), err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http error %s, %s", api.URL.String(), resp.Status)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("content can't be read %s, %s", api.URL.String(), err)
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

// GetEPrint retrieves an EPrint record via REST API and returns a Record structure and error
func (api *EPrintsAPI) GetEPrint(uri string) (*Record, error) {
	api.URL.Path = uri
	resp, err := http.Get(api.URL.String())
	if err != nil {
		return nil, fmt.Errorf("requested %s, %s", api.URL.String(), err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http error %s, %s", api.URL.String(), resp.Status)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("content can't be read %s, %s", api.URL.String(), err)
	}
	rec := new(Record)
	err = xml.Unmarshal(content, &rec)
	if err != nil {
		return nil, err
	}
	return rec, nil
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

// ToISNIs takes an array of pointers to Person and returns a list of ISNI ids
func (persons PersonList) ToISNIs() []string {
	var result []string

	for _, person := range persons {
		result = append(result, person.ISNI)
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

func (record *Record) PubDate() string {
	if record.DateType == "published" {
		return record.Date
	}
	return ""
}

// ListURI returns a list of eprint record ids from the dataset
func (api *EPrintsAPI) ListURI(start, count int) ([]string, error) {
	c, err := dataset.Open(api.Dataset)
	failCheck(err, fmt.Sprintf("ListURI() %s, %s", api.Dataset, err))
	defer c.Close()

	ids := c.Keys()
	results := []string{}
	for i := start; count > 0; count-- {
		results = append(results, ids[i])
	}
	return results, nil
}

// Get retrieves an EPrint record from the dataset
func (api *EPrintsAPI) Get(uri string) (*Record, error) {
	c, err := dataset.Open(api.Dataset)
	failCheck(err, fmt.Sprintf("Get() %s, %s", api.Dataset, err))
	defer c.Close()

	record := new(Record)
	if err := c.Read(uri, record); err != nil {
		return nil, err
	}
	return record, nil
}

// GetIDsBySelectList returns a list of ePrint IDs from a select list filterd by filterFn
func (api *EPrintsAPI) GetIDsBySelectList(slName string, direction int, filterFn func(s string) bool) ([]string, error) {
	c, err := dataset.Open(api.Dataset)
	failCheck(err, fmt.Sprintf("GetIDs() %s, %s", api.Dataset, err))
	defer c.Close()

	sl, err := c.Select(slName)
	if err != nil {
		return nil, err
	}
	sl.Sort(direction)
	ids := []string{}
	for _, id := range sl.List() {
		if filterFn(id) == true {
			ids = append(ids, last(strings.Split(id, indexDelimiter)))
		}
	}
	return ids, err
}

// getRecordList takes a list of ePrint IDs and filters for start and end count return an array of records
func getRecordList(c *dataset.Collection, ePrintIDs []string, start int, count int, filterFn func(*Record) bool) ([]*Record, error) {
	results := []*Record{}
	i := 0
	for _, id := range ePrintIDs {
		rec := new(Record)
		if err := c.Read(id, &rec); err != nil {
			return results, err
		}
		if filterFn(rec) == true {
			if i >= start {
				results = append(results, rec)
			}
			i++
			count--
			if count <= 0 {
				return results, nil
			}
		}
	}
	return results, nil
}

// GetAllRecords reads and returns all records sorted by Publication Date
// returning an array of keys in  ascending or decending order
func (api *EPrintsAPI) GetAllRecords(direction int) ([]*Record, error) {
	c, err := dataset.Open(api.Dataset)
	failCheck(err, fmt.Sprintf("GetAllRecords() %s, %s", api.Dataset, err))
	defer c.Close()

	ids, err := api.GetIDsBySelectList("pubDate", direction, func(s string) bool {
		return true
	})
	if err != nil {
		return nil, err
	}
	// Build a select list in descending publication order
	return getRecordList(c, ids, 0, -1, func(rec *Record) bool {
		return true
	})
}

// GetPublications reads the index for published content and returns a populated
// array of records found in index in ascending or decending order
func (api *EPrintsAPI) GetPublications(start, count, direction int) ([]*Record, error) {
	c, err := dataset.Open(api.Dataset)
	failCheck(err, fmt.Sprintf("GetPublications() %s, %s", api.Dataset, err))
	defer c.Close()

	ids, err := api.GetIDsBySelectList("pubDate", direction, func(s string) bool {
		return true
	})
	if err != nil {
		return nil, err
	}
	return getRecordList(c, ids, start, count, func(rec *Record) bool {
		if rec.IsPublished == "pub" {
			return true
		}
		return false
	})
}

// GetArticles reads the index for published content and returns a populated
// array of records found in index in decending order
func (api *EPrintsAPI) GetArticles(start, count, direction int) ([]*Record, error) {
	c, err := dataset.Open(api.Dataset)
	failCheck(err, fmt.Sprintf("GetArticles() %s, %s", api.Dataset, err))
	defer c.Close()

	ids, err := api.GetIDsBySelectList("pubDate", direction, func(s string) bool {
		return true
	})
	if err != nil {
		return nil, err
	}
	return getRecordList(c, ids, start, count, func(rec *Record) bool {
		if rec.IsPublished == "pub" && rec.Type == "article" {
			return true
		}
		return false
	})
}

// GetLocalGroups returns a JSON list of unique Group names in index
func (api *EPrintsAPI) GetLocalGroups(direction int) ([]string, error) {
	c, err := dataset.Open(api.Dataset)
	failCheck(err, fmt.Sprintf("GetLocalGroups() %s, %s", api.Dataset, err))
	defer c.Close()

	sl, err := c.Select("localGroup")
	if err != nil {
		return nil, err
	}
	sl.Sort(direction)

	// Note: Aggregate the local group names
	groupNames := []string{}
	lastGroup := ""
	groupName := ""
	for _, id := range sl.List() {
		groupName = first(strings.Split(id, indexDelimiter))
		if groupName != lastGroup {
			groupNames = append(groupNames, groupName)
			lastGroup = groupName
		}
	}
	return groupNames, nil
}

// GetLocalGroupPublications returns a list of EPrint records with groupName
func (api *EPrintsAPI) GetLocalGroupPublications(groupName string, start, count, direction int) ([]*Record, error) {
	c, err := dataset.Open(api.Dataset)
	failCheck(err, fmt.Sprintf("GetLocalGroupPublications() %s, %s", api.Dataset, err))
	defer c.Close()

	// Note: Filter for groupName, passing matching eprintIDs to getRecordList()
	ids, err := api.GetIDsBySelectList("localGroup", direction, func(s string) bool {
		parts := strings.Split(s, indexDelimiter)
		grp := first(parts)
		if groupName == grp {
			return true
		}
		return false
	})
	if err != nil {
		return nil, err
	}
	return getRecordList(c, ids, start, count, func(rec *Record) bool {
		if rec.IsPublished == "pub" {
			return true
		}
		return false
	})
}

// GetLocalGroupArticles returns a list of EPrint records with groupName
func (api *EPrintsAPI) GetLocalGroupArticles(groupName string, start, count, direction int) ([]*Record, error) {
	c, err := dataset.Open(api.Dataset)
	failCheck(err, fmt.Sprintf("GetLocalGroupArticles() %s, %s", api.Dataset, err))
	defer c.Close()

	// Note: Filter for groupName, passing matching eprintIDs to getRecordList()
	ids, err := api.GetIDsBySelectList("localGroup", direction, func(s string) bool {
		parts := strings.Split(s, indexDelimiter)
		grp := first(parts)
		if groupName == grp {
			return true
		}
		return false
	})
	if err != nil {
		return nil, err
	}
	return getRecordList(c, ids, start, count, func(rec *Record) bool {
		if rec.IsPublished == "pub" && rec.Type == "article" {
			return true
		}
		return false
	})
}

// GetORCIDs (or ISNI) returns a list unique of ORCID/ISNI IDs in index
func (api *EPrintsAPI) GetORCIDs(direction int) ([]string, error) {
	c, err := dataset.Open(api.Dataset)
	failCheck(err, fmt.Sprintf("GetORCIDs() %s, %s", api.Dataset, err))
	defer c.Close()

	sl, err := c.Select("orcid")
	if err != nil {
		return nil, err
	}
	sl.Sort(direction)

	// Note: Filter for orcid, passing matching eprintIDs to getRecordList()
	orcids := []string{}
	lastORCID := ""
	for _, id := range sl.List() {
		orcid := first(strings.Split(id, indexDelimiter))
		if orcid != lastORCID {
			lastORCID = orcid
			orcids = append(orcids, orcid)
		}
	}
	return orcids, nil
}

// GetORCIDPublications returns a list of EPrint records with a given ORCID
func (api *EPrintsAPI) GetORCIDPublications(orcid string, start, count, direction int) ([]*Record, error) {
	c, err := dataset.Open(api.Dataset)
	failCheck(err, fmt.Sprintf("GetORCIDPublications() %s, %s", api.Dataset, err))
	defer c.Close()

	// Note: Filter for orcid, passing matching eprintIDs to getRecordList()
	ids, err := api.GetIDsBySelectList("orcid", direction, func(s string) bool {
		key := first(strings.Split(s, indexDelimiter))
		if orcid == key {
			return true
		}
		return false
	})
	if err != nil {
		return nil, err
	}
	return getRecordList(c, ids, start, count, func(rec *Record) bool {
		if rec.IsPublished == "pub" {
			return true
		}
		return false
	})
}

// GetORCIDArticles returns a list of EPrint records with a given ORCID
func (api *EPrintsAPI) GetORCIDArticles(orcid string, start, count, direction int) ([]*Record, error) {
	c, err := dataset.Open(api.Dataset)
	failCheck(err, fmt.Sprintf("GetORCIDArticles() %s, %s", api.Dataset, err))
	defer c.Close()

	// Note: Filter for orcid, passing matching eprintIDs to getRecordList()
	ids, err := api.GetIDsBySelectList("orcid", direction, func(s string) bool {
		key := first(strings.Split(s, indexDelimiter))
		if orcid == key {
			return true
		}
		return false
	})
	if err != nil {
		return nil, err
	}
	return getRecordList(c, ids, start, count, func(rec *Record) bool {
		if rec.IsPublished == "pub" && rec.Type == "article" {
			return true
		}
		return false
	})
}

// RenderEPrint writes a single EPrint record to disc.
func (api *EPrintsAPI) RenderEPrint(basepath string, record *Record) error {
	// Convert record to JSON
	src, err := json.Marshal(record)
	if err != nil {
		return err
	}
	fname := path.Join(basepath, fmt.Sprintf("%d.json", record.ID))
	return ioutil.WriteFile(fname, src, 0664)
}

// RenderDocuments writes JSON, HTML, include and rss to the directory indicated by docpath
func (api *EPrintsAPI) RenderDocuments(docTitle, docDescription, docpath string, records []*Record) error {
	// Create the the directory part of docpath if neccessary
	if _, err := os.Open(path.Join(api.Htdocs, docpath)); err != nil && os.IsNotExist(err) == true {
		os.MkdirAll(path.Join(api.Htdocs, path.Dir(docpath)), 0775)
	}

	//NOTE: create a data wrapper for HTML page creation
	pageData := &struct {
		Version        string
		Basepath       string
		ApiURL         string
		SiteURL        string
		DocTitle       string
		DocDescription string
		Records        []*Record
	}{
		Version:        Version,
		Basepath:       docpath,
		ApiURL:         api.URL.String(),
		SiteURL:        api.SiteURL.String(),
		DocTitle:       docTitle,
		DocDescription: docDescription,
		Records:        records,
	}

	// Writing JSON file
	fname := path.Join(api.Htdocs, docpath+".json")
	src, err := json.Marshal(records)
	if err != nil {
		return fmt.Errorf("Can't convert records to JSON %s, %s", fname, err)
	}
	err = ioutil.WriteFile(fname, src, 0664)
	if err != nil {
		return fmt.Errorf("Can't write %s, %s", fname, err)
	}

	// Write out RSS 2.0 file
	fname = path.Join(api.TemplatePath, "rss.xml")
	rss20, err := ioutil.ReadFile(fname)
	if err != nil {
		return fmt.Errorf("Can't open template %s, %s", fname, err)
	}
	rssTmpl, err := template.New("rss").Funcs(TmplFuncs).Parse(string(rss20))
	if err != nil {
		return fmt.Errorf("Can't convert records to RSS %s, %s", fname, err)
	}
	fname = path.Join(api.Htdocs, docpath) + ".rss"
	out, err := os.Create(fname)
	if err != nil {
		return fmt.Errorf("Can't write %s, %s", fname, err)
	}
	if err := rssTmpl.Execute(out, pageData); err != nil {
		return fmt.Errorf("Can't render %s, %s", fname, err)
	}
	out.Close()

	// Write out BibTeX file.
	bibDoc := []string{}
	for _, rec := range records {
		bibDoc = append(bibDoc, rec.ToBibTeXElement().String())
	}
	fname = path.Join(api.Htdocs, docpath+".bib")
	err = ioutil.WriteFile(fname, []byte(strings.Join(bibDoc, "\n\n")), 0664)
	if err != nil {
		return fmt.Errorf("Can't write %s, %s", fname, err)
	}

	// Write out include file
	fname = path.Join(api.TemplatePath, "page.include")
	pageInclude, err := ioutil.ReadFile(fname)
	if err != nil {
		return fmt.Errorf("Can't open template %s, %s", fname, err)
	}
	pageIncludeTmpl, err := template.New("page.include").Funcs(TmplFuncs).Parse(string(pageInclude))
	if err != nil {
		return fmt.Errorf("Can't parse %s, %s", fname, err)
	}
	fname = path.Join(api.Htdocs, docpath+".include")
	out, err = os.Create(fname)
	if err != nil {
		return fmt.Errorf("Can't write %s, %s", fname, err)
	}
	//log.Printf("Writing %s", fname)
	if err := pageIncludeTmpl.Execute(out, pageData); err != nil {
		return fmt.Errorf("Can't render %s, %s", fname, err)
	}
	out.Close()

	pageHTMLTmpl, err := template.New("page.html").Funcs(TmplFuncs).ParseFiles(
		path.Join(api.TemplatePath, "page.include"),
		path.Join(api.TemplatePath, "page.html"),
	)
	if err != nil {
		return fmt.Errorf("Can't parse %s, %s", fname, err)
	}
	fname = path.Join(api.Htdocs, docpath+".html")
	out, err = os.Create(fname)
	if err != nil {
		return fmt.Errorf("Can't write %s, %s", fname, err)
	}

	//log.Printf("Writing %s", fname)
	if err := pageHTMLTmpl.Execute(out, pageData); err != nil {
		return fmt.Errorf("Can't render %s, %s", fname, err)
	}
	out.Close()

	return nil
}

// BuildPages generates webpages based on the contents of the exported EPrints data.
// The site builder needs to know the name of the BoltDB, the root directory
// for the website and directory to find the templates
func (api *EPrintsAPI) BuildPages(feedSize int, title, target string, filter func(*EPrintsAPI, int, int, int) ([]*Record, error)) error {
	if feedSize < 1 {
		feedSize = DefaultFeedSize
	}
	// Collect the published records
	docPath := path.Join(api.Htdocs, target)
	log.Printf("Building %s", docPath)
	records, err := filter(api, 0, feedSize, Descending)
	if err != nil {
		return fmt.Errorf("Can't get records for %q %s, %s", title, docPath, err)
	}
	if len(records) == 0 {
		log.Printf("No records found for %q %s", title, docPath)
	} else {
		log.Printf("%d records found.", len(records))
	}
	if err := api.RenderDocuments(title, fmt.Sprintf("Building pages 0 to %d descending", feedSize), target, records); err != nil {
		return fmt.Errorf("%q %s error, %s", title, docPath, err)
	}
	return nil
}

func (api *EPrintsAPI) BuildEPrintMirror() error {
	// checkPath checks  and creates a path if needed
	checkPath := func(p string) error {
		_, err := os.Stat(p)
		if os.IsExist(err) == true {
			return nil
		}
		return os.MkdirAll(p, 0775)
	}

	ids, err := api.GetIDsBySelectList("keys", Descending, func(s string) bool {
		return true
	})
	if err != nil {
		return err
	}

	//FIXME: this should really copy the data/collection with a filter for "published" items

	// Setup subdirs to hold all the individual eprint records.
	keys := []string{}
	subdir := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	q := len(subdir)
	// Make subdirs as needed
	for _, p := range subdir {
		checkPath(path.Join(api.Htdocs, api.RepositoryPath, p))
	}
	total := len(ids)
	i := 0
	for _, uri := range ids {
		record, err := api.Get(uri)
		if err != nil {
			return err
		}
		basepath := path.Join(api.Htdocs, api.RepositoryPath, subdir[i%q])
		err = api.RenderEPrint(basepath, record)
		if err != nil {
			return err
		}
		//NOTE: We only save the path relative to the web docroot.
		keys = append(keys, path.Join(api.RepositoryPath, subdir[i%q], fmt.Sprintf("%d.json", record.ID)))
		if (i % 1000) == 0 {
			log.Printf("%d of %d records written", i, total)
		}
		i++
	}
	log.Printf("%d of %d records written", i, total)
	src, err := json.Marshal(keys)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(api.Htdocs, api.RepositoryPath, "eprints.json"), src, 0664)
}

// BuildSelectLists iterates over the exported data and creates fresh selectLists
func (api *EPrintsAPI) BuildSelectLists() error {
	c, err := dataset.Create(api.Dataset, dataset.GenerateBucketNames(dataset.DefaultAlphabet, 2))
	failCheck(err, fmt.Sprintf("BuildSelectLists() %s, %s", api.Dataset, err))
	defer c.Close()

	sLists := map[string]*dataset.SelectList{}
	// Clear the select lists
	log.Println("Clearing select lists")
	for _, name := range slNames {
		c.Clear(name)
		sLists[name], err = c.Select(name)
		if err != nil {
			return err
		}
	}

	// Now iterate over the records and populate fresh select lists
	for i, ky := range c.Keys() {
		rec := new(Record)
		err := c.Read(ky, &rec)
		if err != nil {
			return err
		}
		// Update pubDate select list
		dt := normalizeDate(rec.Date)
		if rec.DateType == "published" && rec.Date != "" {
			sLists["pubDate"].Push(fmt.Sprintf("%s%s%d", dt, indexDelimiter, rec.ID))
		}

		// Update localGroup select list
		if len(rec.LocalGroup) > 0 {
			for _, grp := range rec.LocalGroup {
				grp = strings.TrimSpace(grp)
				if len(grp) > 0 {
					sLists["localGroup"].Push(fmt.Sprintf("%s%s%s%s%d", grp, indexDelimiter, dt, indexDelimiter, rec.ID))
				}
			}
		}
		// Update orcids, isnis and authors select list
		if len(rec.Creators) > 0 {
			for _, person := range rec.Creators {
				orcid := strings.TrimSpace(person.ORCID)
				isni := strings.TrimSpace(person.ISNI)

				// Update orcid select list
				if len(orcid) > 0 {
					sLists["orcid"].Push(fmt.Sprintf("%s%s%s%s%d", orcid, indexDelimiter, dt, indexDelimiter, rec.ID))
				}

				// Update isni select list
				if len(isni) > 0 {
					sLists["isni"].Push(fmt.Sprintf("%s%s%s%s%d", isni, indexDelimiter, dt, indexDelimiter, rec.ID))
				}
			}
		}

		// Add funders and grantNumbers to select lists
		if len(rec.Funders) > 0 {
			for _, funder := range rec.Funders {
				funderName := strings.TrimSpace(funder.Agency)
				grantNumber := strings.TrimSpace(funder.GrantNumber)

				// Update funder select list
				if len(funderName) > 0 {
					sLists["funder"].Push(fmt.Sprintf("%s%s%s%s%d", funderName, indexDelimiter, dt, indexDelimiter, rec.ID))
				}
				if len(funderName) > 0 && len(grantNumber) > 0 {
					sLists["grantNumber"].Push(fmt.Sprintf("%s%s%s%s%s%s%d", funderName, indexDelimiter, grantNumber, indexDelimiter, dt, indexDelimiter, rec.ID))
				}
			}
		}
		if (i % 1000) == 0 {
			log.Printf("%d recs processed", i)
		}
	}
	return nil
}

// BuildSite generates a website based on the contents of the exported EPrints data.
// The site builder needs to know the name of the BoltDB, the root directory
// for the website and directory to find the templates
func (api *EPrintsAPI) BuildSite(feedSize int, buildEPrintMirror bool) error {
	var err error

	if feedSize < 1 {
		feedSize = DefaultFeedSize
	}

	// FIXME: This could be replaced by copying all the records in dataset/COLLECTION
	// that are public and published.
	if buildEPrintMirror == true {
		// Build mirror of repository content.
		log.Printf("Mirroring eprint records")
		err = api.BuildEPrintMirror()
		if err != nil {
			return nil
		}
	}

	// Collect the recent publications (all types)
	log.Printf("Building Recently Published")
	err = api.BuildPages(feedSize, "Recently Published", path.Join("recent", "publications"), func(api *EPrintsAPI, start, count, direction int) ([]*Record, error) {
		return api.GetPublications(0, feedSize, Descending)
	})
	if err != nil {
		return err
	}

	// Collect the rencently published  articles
	log.Printf("Building Recent Articles")
	err = api.BuildPages(feedSize, "Recent Articles", path.Join("recent", "articles"), func(api *EPrintsAPI, start, count, direction int) ([]*Record, error) {
		return api.GetArticles(start, count, Descending)
	})
	if err != nil {
		return err
	}

	// Collect EPrints by orcid ID and publish
	log.Printf("Building Person (orcid) works")
	orcids, err := api.GetORCIDs(Ascending)
	if err != nil {
		return err
	}
	log.Printf("Found %d orcids", len(orcids))
	for _, orcid := range orcids {
		// Build a list of recent ORCID Publications
		err = api.BuildPages(-1, fmt.Sprintf("ORCID: %s", orcid), path.Join("person", fmt.Sprintf("%s", orcid), "recent", "publications"), func(api *EPrintsAPI, start, count, direction int) ([]*Record, error) {
			return api.GetORCIDPublications(orcid, start, count, Descending)
		})
		if err != nil {
			return err
		}
		// Build complete list for each orcid
		err = api.BuildPages(-1, fmt.Sprintf("ORCID: %s", orcid), path.Join("person", fmt.Sprintf("%s", orcid), "publications"), func(api *EPrintsAPI, start, count, direction int) ([]*Record, error) {
			return api.GetORCIDPublications(orcid, 0, -1, Descending)
		})
		if err != nil {
			return err
		}
		// Build a list of recent ORCID Articles
		err = api.BuildPages(-1, fmt.Sprintf("ORCID: %s", orcid), path.Join("person", fmt.Sprintf("%s", orcid), "recent", "articles"), func(api *EPrintsAPI, start, count, direction int) ([]*Record, error) {
			return api.GetORCIDArticles(orcid, start, count, Descending)
		})
		if err != nil {
			return err
		}
		// Build complete list of articels for each ORCID
		err = api.BuildPages(-1, fmt.Sprintf("ORCID: %s", orcid), path.Join("person", fmt.Sprintf("%s", orcid), "articles"), func(api *EPrintsAPI, start, count, direction int) ([]*Record, error) {
			return api.GetORCIDArticles(orcid, 0, -1, Descending)
		})
		if err != nil {
			return err
		}
	}

	// Collect EPrints by Group/Affiliation
	log.Printf("Building Local Groups")
	groupNames, err := api.GetLocalGroups(Ascending)
	if err != nil {
		return err
	}
	log.Printf("Found %d groups", len(groupNames))
	for _, groupName := range groupNames {
		// Build recently for each affiliation
		slug, err := slugify(groupName)
		if err != nil {
			log.Printf("Skipping %q, %s\n", groupName, err)
		} else {
			err = api.BuildPages(-1, fmt.Sprintf("%s", groupName), path.Join("affiliation", slug, "recent", "publications"), func(api *EPrintsAPI, start, count, direction int) ([]*Record, error) {
				return api.GetLocalGroupPublications(groupName, start, count, Descending)
			})
			if err != nil {
				return err
			}
			// Build complete list for each affiliation
			err = api.BuildPages(-1, fmt.Sprintf("%s", groupName), path.Join("affiliation", slug, "publications"), func(api *EPrintsAPI, start, count, direction int) ([]*Record, error) {
				return api.GetLocalGroupPublications(groupName, 0, -1, Descending)
			})
			if err != nil {
				return err
			}
			// Build recent articles for each affiliation
			err = api.BuildPages(-1, fmt.Sprintf("%s", groupName), path.Join("affiliation", slug, "recent", "articles"), func(api *EPrintsAPI, start, count, direction int) ([]*Record, error) {
				return api.GetLocalGroupArticles(groupName, start, count, Descending)
			})
			if err != nil {
				return err
			}
			// Build complete list of articles for each affiliation
			err = api.BuildPages(-1, fmt.Sprintf("%s", groupName), path.Join("affiliation", slug, "articles"), func(api *EPrintsAPI, start, count, direction int) ([]*Record, error) {
				return api.GetLocalGroupArticles(groupName, 0, -1, Descending)
			})
			if err != nil {
				return err
			}
		}
	}

	// Collect EPrints by Funders
	log.Printf("Building Funders")
	funderNames, err := api.GetFunders(Ascending)
	if err != nil {
		return err
	}
	log.Printf("Found %s funders", len(funderNames))
	for _, funderName := range funderNames {
		slug, err := slugify(funderName)
		if err != nil {
			log.Printf("Skipping %q, %s\n", funderName, err)
		} else {
			// Build recently for each funder
			err = api.BuildPages(-1, fmt.Sprintf("%s", funderName), path.Join("funder", slug, "recent", "publications"), func(api *EPrintsAPI, start, count, direction int) ([]*Record, error) {
				return api.GetFunderPublications(funderName, start, count, Descending)
			})
			if err != nil {
				return err
			}
			// Build complete list for each funder
			err = api.BuildPages(-1, fmt.Sprintf("%s", funderName), path.Join("funder", slug, "publications"), func(api *EPrintsAPI, start, count, direction int) ([]*Record, error) {
				return api.GetFunderPublications(funderName, 0, -1, Descending)
			})
			if err != nil {
				return err
			}
			// Build recent articles for each funder
			err = api.BuildPages(-1, fmt.Sprintf("%s", funderName), path.Join("funder", slug, "recent", "articles"), func(api *EPrintsAPI, start, count, direction int) ([]*Record, error) {
				return api.GetFunderArticles(funderName, start, count, Descending)
			})
			if err != nil {
				return err
			}
			// Build complete list of articles for each funder
			err = api.BuildPages(-1, fmt.Sprintf("%s", funderName), path.Join("funder", slug, "articles"), func(api *EPrintsAPI, start, count, direction int) ([]*Record, error) {
				return api.GetFunderArticles(funderName, 0, -1, Descending)
			})
			if err != nil {
				return err
			}
		}
	}

	// Collect EPrints by Funders/Grant Number
	funderGrantNames, err := api.GetGrantNumbersByFunder(Ascending)
	if err != nil {
		return err
	}
	log.Printf("Found %s funders/grant numbers", len(funderGrantNames))
	for _, funderGrantName := range funderGrantNames {
		parts := strings.Split(funderGrantName, indexDelimiter)
		funderName := first(parts)
		grantNumber := second(parts)
		slug, err := slugify(funderName)
		if err != nil {
			log.Printf("Skipping %q, %s\n", funderName, err)
		} else {
			// Build recently for each funder
			err = api.BuildPages(-1, funderName, path.Join("funder", slug, "grant", grantNumber, "recent", "publications"), func(api *EPrintsAPI, start, count, direction int) ([]*Record, error) {
				return api.GetGrantNumberPublications(funderName, grantNumber, start, count, Descending)
			})
			if err != nil {
				return err
			}
			// Build complete list for each funder
			err = api.BuildPages(-1, funderName, path.Join("funder", slug, "grant", grantNumber, "publications"), func(api *EPrintsAPI, start, count, direction int) ([]*Record, error) {
				return api.GetGrantNumberPublications(funderName, grantNumber, 0, -1, Descending)
			})
			if err != nil {
				return err
			}
			// Build recent articles for each funder
			err = api.BuildPages(-1, funderName, path.Join("funder", slug, "grant", grantNumber, "recent", "articles"), func(api *EPrintsAPI, start, count, direction int) ([]*Record, error) {
				return api.GetGrantNumberArticles(funderName, grantNumber, start, count, Descending)
			})
			if err != nil {
				return err
			}
			// Build complete list of articles for each funder
			err = api.BuildPages(-1, funderName, path.Join("funder", slug, "grant", grantNumber, "articles"), func(api *EPrintsAPI, start, count, direction int) ([]*Record, error) {
				return api.GetGrantNumberArticles(funderName, grantNumber, 0, -1, Descending)
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
