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
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	// 3rd Party packages
	"github.com/boltdb/bolt"
)

const (
	// Version is the revision number for this implementation of epgo
	Version = "0.0.5"

	// Ascending sorts from lowest (oldest) to highest (newest)
	Ascending = iota
	// Descending sorts from highest (newest) to lowest (oldest)
	Descending = iota

	// EPrintsExportBatchSize sets the summary output frequency when exporting content from E-Prints
	EPrintsExportBatchSize = 1000

	// DefaultFeedSize sets the default size of rss, JSON, HTML include and index lists
	DefaultFeedSize = 25
)

// These are our main bucket and index buckets
var (
	// Primary collection
	ePrintBucket = []byte("eprints")

	// Indexes available
	indexDelimiter = "|"
	pubDatesBucket = []byte("publicationDates")
	// publicationsBucket  = []byte("publications")
	// titlesBucket        = []byte("titles")
	// subjectsBucket      = []byte("subjects")
	// authors             = []byte("authors")
	// additionDatesBucket = []byte("additionsDates")

	// EPTmplFuncs provides a common set of functions available to templates
	EPTmplFuncs = template.FuncMap{
		"year": func(s string) string {
			var (
				dt  time.Time
				err error
			)
			if s == "now" {
				dt = time.Now()
			} else {
				dt, err = time.Parse("2006-01-02", normalizeDate(s))
				if err != nil {
					return ""
				}
			}
			return dt.Format("2006")
		},
		"rfc3339": func(s string) string {
			var (
				dt  time.Time
				err error
			)
			if s == "now" {
				dt = time.Now()
			} else {
				dt, err = time.Parse("2006-01-02", normalizeDate(s))
				if err != nil {
					return ""
				}
			}
			return dt.Format(time.RFC3339)
		},
		"rfc1123": func(s string) string {
			var (
				dt  time.Time
				err error
			)
			if s == "now" {
				dt = time.Now()
			} else {
				dt, err = time.Parse("2006-01-02", normalizeDate(s))
				if err != nil {
					return ""
				}
			}
			return dt.Format(time.RFC1123)
		},
		"rfc1123z": func(s string) string {
			var (
				dt  time.Time
				err error
			)
			if s == "now" {
				dt = time.Now()
			} else {
				dt, err = time.Parse("2006-01-02", normalizeDate(s))
				if err != nil {
					return ""
				}
			}
			return dt.Format(time.RFC1123Z)
		},
		"rfc822z": func(s string) string {
			var (
				dt  time.Time
				err error
			)
			if s == "now" {
				dt = time.Now()
			} else {
				dt, err = time.Parse("2006-01-02", normalizeDate(s))
				if err != nil {
					return ""
				}
			}
			return dt.Format(time.RFC822Z)
		},
		"rfc822": func(s string) string {
			var (
				dt  time.Time
				err error
			)
			if s == "now" {
				dt = time.Now()
			} else {
				dt, err = time.Parse("2006-01-02", normalizeDate(s))
				if err != nil {
					return ""
				}
			}
			return dt.Format(time.RFC822)
		},
	}
)

func failCheck(err error, msg string) {
	if err != nil {
		log.Fatalf("%s\n", msg)
	}
}

// EPrintsAPI holds the basic connectin information to read the REST API for EPrints
type EPrintsAPI struct {
	URL       *url.URL `xml:"epgo>api_url" json:"api_url"`     // EPGO_API_URL
	DBName    string   `xml:"epgo>dbname" json:"dbname"`       // EPGO_DBNAME
	Htdocs    string   `xml:"epgo>htdocs" json:"htdocs"`       // EPGO_HTDOCS
	Templates string   `xml:"epgo>templates" json:"templates"` // EPGO_TEMPLATES
	SiteURL   *url.URL `xml:"epgo>site_url" josn:"site_url"`   // EPGO_SITE_URL
}

// Person returns the contents of eprint>creators>item>name as a struct
type Person struct {
	XMLName xml.Name `json:"-"`
	Given   string   `xml:"name>given" json:"given"`
	Family  string   `xml:"name>family" json:"family"`
	ID      string   `xml:"id,omitempty" json:"id,omitempty"`
	ORCID   string   `xml:"orcid,omitempty" json:"orcid,omitempty"`
}

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
	GrantNumber string   `xml:"grant_number,omitempty" json:"grant_number,omitempty"`
}

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

// Record returns a structure that can be converted to JSON easily
type Record struct {
	XMLName              xml.Name           `json:"-"`
	Title                string             `xml:"eprint>title" json:"title"`
	URI                  string             `json:"uri"`
	Abstract             string             `xml:"eprint>abstract" json:"abstract"`
	Documents            []*Document        `xml:"eprint>documents>document" json:"documents"`
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
	Creators             []*Person          `xml:"eprint>creators>item" json:"creators"`
	IsPublished          string             `xml:"eprint>ispublished" json:"ispublished"`
	Subjects             []string           `xml:"eprint>subjects>item" json:"subjects"`
	FullTextStatus       string             `xml:"eprint>full_text_status" json:"full_text_status"`
	Keywords             string             `xml:"eprint>keywords" json:"keywords"`
	Date                 string             `xml:"eprint>date" json:"data"`
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
	OtherNumberingSystem []*NumberingSystem `xml:"eprint>other_numbering_system>item,omitempty" json:"other_numbering_system,omitempty"`
	Funders              []*Funder          `xml:"eprint>funders>item" json:"funders"`
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

// New creates a new API instance
func New() (*EPrintsAPI, error) {
	var err error
	apiURL := os.Getenv("EPGO_API_URL")
	siteURL := os.Getenv("EPGO_SITE_URL")
	htdocs := os.Getenv("EPGO_HTDOCS")
	dbName := os.Getenv("EPGO_DBNAME")
	templates := os.Getenv("EPGO_TEMPLATES")

	if apiURL == "" {
		return nil, fmt.Errorf("Environment not configured, missing EPGO_API_URL")
	}
	api := new(EPrintsAPI)
	api.URL, err = url.Parse(apiURL)
	if err != nil {
		return nil, fmt.Errorf("EPGO_API_URL malformed %s, %s", apiURL, err)
	}
	api.SiteURL, err = url.Parse(siteURL)
	if err != nil {
		return nil, fmt.Errorf("EPGO_SITE_URL malformed %s, %s", siteURL, err)
	}
	if htdocs == "" {
		htdocs = "htdocs"
	}
	if dbName == "" {
		dbName = "eprints"
	}
	if templates == "" {
		templates = "templates"
	}
	api.Htdocs = htdocs
	api.DBName = dbName
	api.Templates = templates
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

// ListEPrintsURI returns a list of eprint record ids
func (api *EPrintsAPI) ListEPrintsURI() ([]string, error) {
	var results []string

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
	for _, val := range eIDs.IDs {
		if strings.HasSuffix(val, ".xml") == true {
			results = append(results, "/"+path.Join("rest", "eprint", val))
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

// ExportEPrints from highest ID to lowest for cnt. Saves each record in a DB and indexes published ones
func (api *EPrintsAPI) ExportEPrints(cnt int) error {
	db, err := bolt.Open(api.DBName, 0660, &bolt.Options{Timeout: 1 * time.Second, ReadOnly: false})
	failCheck(err, fmt.Sprintf("Export %s failed to open db, %s", api.DBName, err))
	defer db.Close()
	// Make sure we have a buckets to store things in
	db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(ePrintBucket); err != nil {
			return fmt.Errorf("create bucket %s: %s", ePrintBucket, err)
		}
		if _, err := tx.CreateBucketIfNotExists(pubDatesBucket); err != nil {
			return fmt.Errorf("create bucket %s: %s", pubDatesBucket, err)
		}
		return nil
	})

	uris, err := api.ListEPrintsURI()
	failCheck(err, fmt.Sprintf("Export %s failed, %s", api.URL.String(), err))
	//NOTE: I am sorting the URI by decscending ID number so that the newest articles
	// are exported first
	sort.Sort(byURI(uris))

	uriCnt := len(uris)
	if cnt < 0 {
		cnt = uriCnt
	}
	j := 0 // success count
	k := 0 // error count
	log.Printf("Exporting %d of %d uris", cnt, uriCnt)
	for i := 0; i < uriCnt && cnt > 0; i++ {
		uri := uris[i]
		rec, err := api.GetEPrint(uri)
		if err != nil {
			log.Printf("Failed to get %s, %s\n", uri, err)
			k++
		} else {
			rec.URI = strings.TrimPrefix(strings.TrimSuffix(uri, path.Ext(".xml")), "/rest")
			src, err := json.Marshal(rec)
			if err != nil {
				log.Printf("json.Marshal() failed on %s, %s", uri, err)
				k++
			} else {
				err := db.Update(func(tx *bolt.Tx) error {
					b := tx.Bucket(ePrintBucket)
					err := b.Put([]byte(uri), src)
					if err == nil {
						// See if we need to add this to the publicationDates index
						if rec.DateType == "published" && rec.Date != "" {
							idx := tx.Bucket(pubDatesBucket)
							dt := normalizeDate(rec.Date)
							err = idx.Put([]byte(fmt.Sprintf("%s%s%s", dt, indexDelimiter, uri)), []byte(uri))
						}
						j++
						cnt--
					}
					return err
				})
				if err != nil {
					log.Printf("Failed to save eprint %s, %s\n", uri, err)
					k++
				}
			}
		}
		if (i % EPrintsExportBatchSize) == 0 {
			log.Printf("%d uri processed, %d exported, %d unexported", i+1, j, k)
		}
	}
	log.Printf("%d uri processed, %d exported, %d unexported", len(uris), j, k)
	return nil
}

// ListURI returns a list of eprint record ids from the database
func (api *EPrintsAPI) ListURI(start, count int) ([]string, error) {
	var results []string
	db, err := bolt.Open(api.DBName, 0660, &bolt.Options{Timeout: 1 * time.Second, ReadOnly: true})
	failCheck(err, fmt.Sprintf("ListURI %s failed to open db, %s", api.DBName, err))
	defer db.Close()
	// Make sure we have a buckets to store things in
	db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(ePrintBucket); err != nil {
			return fmt.Errorf("create bucket %s: %s", ePrintBucket, err)
		}
		if _, err := tx.CreateBucketIfNotExists(pubDatesBucket); err != nil {
			return fmt.Errorf("create bucket %s: %s", pubDatesBucket, err)
		}
		return nil
	})
	err = db.View(func(tx *bolt.Tx) error {
		recs := tx.Bucket(ePrintBucket)
		c := recs.Cursor()
		p := 0
		for uri, _ := c.First(); uri != nil && count > 0; uri, _ = c.Next() {
			if p >= start {
				results = append(results, string(uri))
				count--
			}
			p++
		}
		return nil
	})
	return results, err
}

// Get retrieves an EPrint record from the database
func (api *EPrintsAPI) Get(uri string) (*Record, error) {
	db, err := bolt.Open(api.DBName, 0660, &bolt.Options{Timeout: 1 * time.Second, ReadOnly: true})
	failCheck(err, fmt.Sprintf("Get(%q) %s failed to open db, %s", uri, api.DBName, err))
	defer db.Close()
	// Make sure we have a buckets to store things in
	db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(ePrintBucket); err != nil {
			return fmt.Errorf("create bucket %s: %s", ePrintBucket, err)
		}
		if _, err := tx.CreateBucketIfNotExists(pubDatesBucket); err != nil {
			return fmt.Errorf("create bucket %s: %s", pubDatesBucket, err)
		}
		return nil
	})
	record := new(Record)
	err = db.View(func(tx *bolt.Tx) error {
		recs := tx.Bucket(ePrintBucket)
		src := recs.Get([]byte(uri))
		err := json.Unmarshal(src, &record)
		if err != nil {
			return err
		}
		return nil
	})
	return record, err
}

// GetPublishedRecords reads the index for published content and returns a populated
// array of records found in index in decending order
func (api *EPrintsAPI) GetPublishedRecords(start, count, direction int) ([]*Record, error) {
	db, err := bolt.Open(api.DBName, 0660, &bolt.Options{Timeout: 1 * time.Second, ReadOnly: true})
	failCheck(err, fmt.Sprintf("GetPulishedRecords() %s failed to open db, %s", api.DBName, err))
	defer db.Close()
	// Make sure we have a buckets to store things in
	db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(ePrintBucket); err != nil {
			return fmt.Errorf("create bucket %s: %s", ePrintBucket, err)
		}
		if _, err := tx.CreateBucketIfNotExists(pubDatesBucket); err != nil {
			return fmt.Errorf("create bucket %s: %s", pubDatesBucket, err)
		}
		return nil
	})

	//	var records []Record
	var (
		results []*Record
	)
	switch direction {
	case Ascending:
		err = db.View(func(tx *bolt.Tx) error {
			recs := tx.Bucket(ePrintBucket)
			idx := tx.Bucket(pubDatesBucket)
			c := idx.Cursor()
			p := 0
			for k, uri := c.First(); k != nil && count > 0; k, uri = c.Next() {
				if p >= start {
					rec := new(Record)
					src := recs.Get([]byte(uri))
					err := json.Unmarshal(src, rec)
					if err != nil {
						return fmt.Errorf("Can't unmarshal %s, %s", uri, err)
					}
					if rec.IsPublished == "pub" {
						results = append(results, rec)
						count--
					}
				}
				p++
			}
			return nil
		})
	case Descending:
		err = db.View(func(tx *bolt.Tx) error {
			recs := tx.Bucket(ePrintBucket)
			idx := tx.Bucket(pubDatesBucket)
			c := idx.Cursor()
			p := 0
			for k, uri := c.Last(); k != nil && count > 0; k, uri = c.Prev() {
				if p >= start {
					rec := new(Record)
					src := recs.Get([]byte(uri))
					err := json.Unmarshal(src, rec)
					if err != nil {
						return fmt.Errorf("Can't unmarshal %s, %s", uri, err)
					}
					if rec.IsPublished == "pub" {
						results = append(results, rec)
						count--
					}
				}
				p++
			}
			return nil
		})
	}
	return results, err
}

// GetPublishedArticles reads the index for published content and returns a populated
// array of records found in index in decending order
func (api *EPrintsAPI) GetPublishedArticles(start, count, direction int) ([]*Record, error) {
	db, err := bolt.Open(api.DBName, 0660, &bolt.Options{Timeout: 1 * time.Second, ReadOnly: true})
	failCheck(err, fmt.Sprintf("GetPublishedArticles() %s failed to open db, %s", api.DBName, err))
	defer db.Close()
	// Make sure we have a buckets to store things in
	db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(ePrintBucket); err != nil {
			return fmt.Errorf("create bucket %s: %s", ePrintBucket, err)
		}
		if _, err := tx.CreateBucketIfNotExists(pubDatesBucket); err != nil {
			return fmt.Errorf("create bucket %s: %s", pubDatesBucket, err)
		}
		return nil
	})

	//	var records []Record
	var (
		results []*Record
	)
	switch direction {
	case Ascending:
		err = db.View(func(tx *bolt.Tx) error {
			recs := tx.Bucket(ePrintBucket)
			idx := tx.Bucket(pubDatesBucket)
			c := idx.Cursor()
			p := 0
			for k, uri := c.First(); k != nil && count > 0; k, uri = c.Next() {
				if p >= start {
					rec := new(Record)
					src := recs.Get([]byte(uri))
					err := json.Unmarshal(src, rec)
					if err != nil {
						return fmt.Errorf("Can't unmarshal %s, %s", uri, err)
					}
					if rec.Type == "article" && rec.IsPublished == "pub" {
						results = append(results, rec)
						count--
					}
				}
				p++
			}
			return nil
		})
	case Descending:
		err = db.View(func(tx *bolt.Tx) error {
			recs := tx.Bucket(ePrintBucket)
			idx := tx.Bucket(pubDatesBucket)
			c := idx.Cursor()
			p := 0
			for k, uri := c.Last(); k != nil && count > 0; k, uri = c.Prev() {
				if p >= start {
					rec := new(Record)
					src := recs.Get([]byte(uri))
					err := json.Unmarshal(src, rec)
					if err != nil {
						return fmt.Errorf("Can't unmarshal %s, %s", uri, err)
					}
					if rec.Type == "article" && rec.IsPublished == "pub" {
						results = append(results, rec)
						count--
					}
				}
				p++
			}
			return nil
		})
	}
	return results, err
}

// RenderDocuments writes JSON, HTML, include and rss to the directory indicated by basepath
func (api *EPrintsAPI) RenderDocuments(docTitle, docDescription, basepath string, records []*Record) error {
	// Create the basepath if neccessary
	if _, err := os.Open(path.Join(api.Htdocs, basepath)); err != nil && os.IsNotExist(err) == true {
		os.MkdirAll(path.Join(api.Htdocs, basepath), 0775)
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
		Basepath:       basepath,
		ApiURL:         api.URL.String(),
		SiteURL:        api.SiteURL.String(),
		DocTitle:       docTitle,
		DocDescription: docDescription,
		Records:        records,
	}

	// Writing JSON file
	fname := path.Join(api.Htdocs, basepath, "index.json")
	src, err := json.Marshal(records)
	if err != nil {
		return fmt.Errorf("Can't convert records to JSON %s, %s", fname, err)
	}
	err = ioutil.WriteFile(fname, src, 0664)
	if err != nil {
		return fmt.Errorf("Can't write %s, %s", fname, err)
	}
	// Write out RSS 2.0 file
	fname = path.Join(api.Templates, "rss.xml")
	rss20, err := ioutil.ReadFile(fname)
	if err != nil {
		return fmt.Errorf("Can't open template %s, %s", fname, err)
	}
	rssTmpl, err := template.New("rss").Funcs(EPTmplFuncs).Parse(string(rss20))
	if err != nil {
		return fmt.Errorf("Can't convert records to RSS %s, %s", fname, err)
	}
	fname = path.Join(api.Htdocs, basepath, "rss.xml")
	out, err := os.Create(fname)
	if err != nil {
		return fmt.Errorf("Can't write %s, %s", fname, err)
	}
	if err := rssTmpl.Execute(out, pageData); err != nil {
		return fmt.Errorf("Can't render %s, %s", fname, err)
	}
	out.Close()

	// Write out include file
	fname = path.Join(api.Templates, "page.include")
	pageInclude, err := ioutil.ReadFile(fname)
	if err != nil {
		return fmt.Errorf("Can't open template %s, %s", fname, err)
	}
	pageIncludeTmpl, err := template.New("page.include").Funcs(EPTmplFuncs).Parse(string(pageInclude))
	if err != nil {
		return fmt.Errorf("Can't parse %s, %s", fname, err)
	}
	fname = path.Join(api.Htdocs, basepath, "index.include")
	out, err = os.Create(fname)
	if err != nil {
		return fmt.Errorf("Can't write %s, %s", fname, err)
	}
	log.Printf("Writing %s", fname)
	if err := pageIncludeTmpl.Execute(out, pageData); err != nil {
		return fmt.Errorf("Can't render %s, %s", fname, err)
	}
	out.Close()

	pageHTMLTmpl, err := template.New("page.html").Funcs(EPTmplFuncs).ParseFiles(
		path.Join(api.Templates, "page.include"),
		path.Join(api.Templates, "page.html"),
	)
	if err != nil {
		return fmt.Errorf("Can't parse %s, %s", fname, err)
	}
	fname = path.Join(api.Htdocs, basepath, "index.html")
	out, err = os.Create(fname)
	if err != nil {
		return fmt.Errorf("Can't write %s, %s", fname, err)
	}

	log.Printf("Writing %s", fname)
	if err := pageHTMLTmpl.Execute(out, pageData); err != nil {
		return fmt.Errorf("Can't render %s, %s", fname, err)
	}
	out.Close()
	return nil
}

// BuildSite generates a website based on the contents of the exported EPrints data.
// The site builder needs to know the name of the BoltDB, the root directory
// for the website and directory to find the templates
func (api *EPrintsAPI) BuildSite(feedSize int) error {
	if feedSize < 1 {
		feedSize = DefaultFeedSize
	}
	// Collect the published records
	log.Printf("Building recently-published")
	records, err := api.GetPublishedRecords(0, feedSize, Descending)
	if err != nil {
		return fmt.Errorf("Can't get published records, %s", err)
	}
	if len(records) == 0 {
		return fmt.Errorf("No published records found")
	}
	log.Printf("%d records found.", len(records))
	if err := api.RenderDocuments("Recently Published", "Recently published items including chapters, collections and articles.", "recently-published", records); err != nil {
		return fmt.Errorf("recently published error, %s", err)
	}
	log.Printf("Building recent-articles")
	// Collect the published articles
	records, err = api.GetPublishedArticles(0, feedSize, Descending)
	log.Printf("%d records found.", len(records))
	if err != nil {
		return fmt.Errorf("Can't get published articles, %s", err)
	}
	if len(records) == 0 {
		return fmt.Errorf("No published articles found")
	}
	if err := api.RenderDocuments("Recent Articles", "Recently published articles.", "recent-articles", records); err != nil {
		return fmt.Errorf("recent articles error, %s", err)
	}
	// FIXME: Should build entire site with searchable content
	return nil
}
