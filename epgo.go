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
	Version = "0.0.8"

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
	indexDelimiter   = "|"
	pubDatesBucket   = []byte("publicationDates")
	localGroupBucket = []byte("localGroup")
	orcidBucket      = []byte("orcid")

	//FIXME: Additional indexes might be useful.
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
	XMLName      xml.Name `json:"-"`
	URL          *url.URL `xml:"epgo>api_url" json:"api_url"`             // EPGO_API_URL
	DBName       string   `xml:"epgo>dbname" json:"dbname"`               // EPGO_DBNAME
	Htdocs       string   `xml:"epgo>htdocs" json:"htdocs"`               // EPGO_HTDOCS
	TemplatePath string   `xml:"epgo>template_path" json:"template_path"` // EPGO_TEMPLATES
	SiteURL      *url.URL `xml:"epgo>site_url" josn:"site_url"`           // EPGO_SITE_URL
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

// Config holds the common configuration elements needed to access and harvest data
type Config struct {
	XMLName xml.Name `json:"-"`
	// ApiURL is the root URL accessed to retrieve data from EPrints.
	ApiURL string `json:"api_url,required" xml:"api_url,required"`
	// Setup configuration for how the feed stores or reads harvested metadata. Becomes the primary bucket in BoltDB
	DBName string `json:"dbName,omitempty" xml:"dbName,omitempty"`
	// Site URL, URL of website hosting feeds and search service (not EPrints or another repository)
	SiteURL string `json:"site_url,omitempty" xml:"site_url,omitempty"`
	// HTDocs directory for the website hosting feeds and search service
	Htdocs string `json:"htdocs,omitempty" xml:"htdocs,omitempty"`
	// TemplatePath directory for generating website described by SiteURL and Htdocs
	TemplatePath string `json:"htdocs,omitempty" xml:"htdocs,omitempty"`
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

// MergeEnv merge environment variables into the configuration structure.
// options are
// + prefix - e.g. EPGO, name space before the first underscore in the envinronment
//      + prefix plus uppercase key forms the complete environment variable name
// + key - the field map (e.g ApiURL maps to API_URL in EPGO_API_URL for prefix EPGO)
// + proposedValue - the proposed value, usually the value from the flags passed in (an empty string means no value provided)
//
// returns an error if a problem is encountered.
func (cfg *Config) MergeEnv(prefix, key, proposedValue string) error {
	// Default is anything set in the environment
	val := os.Getenv(fmt.Sprintf("%s_%s", prefix, key))
	// Override with proposedValue provided (e.g. flag value from the command line)
	if proposedValue != "" {
		val = proposedValue
	}
	switch key {
	case "API_URL":
		cfg.ApiURL = val
	case "DBNAME":
		cfg.DBName = val
	case "SITE_URL":
		cfg.SiteURL = val
	case "HTDOCS":
		cfg.Htdocs = val
	case "TEMPLATE_PATH":
		cfg.TemplatePath = val
	default:
		return fmt.Errorf("%s isn't a known configuration option", key)
	}
	return nil
}

// New creates a new API instance
func New(cfg Config) (*EPrintsAPI, error) {
	var err error
	apiURL := cfg.ApiURL
	siteURL := cfg.SiteURL
	htdocs := cfg.Htdocs
	dbName := cfg.DBName
	templatePath := cfg.TemplatePath

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
	if dbName == "" {
		dbName = "eprints"
	}
	if templatePath == "" {
		templatePath = "templates"
	}
	api.Htdocs = htdocs
	api.DBName = dbName
	api.TemplatePath = templatePath
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

// initBuckets initializes expected buckets if necessary for boltdb
func initBuckets(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(ePrintBucket); err != nil {
			return fmt.Errorf("create bucket %s: %s", ePrintBucket, err)
		}
		if _, err := tx.CreateBucketIfNotExists(pubDatesBucket); err != nil {
			return fmt.Errorf("create bucket %s: %s", pubDatesBucket, err)
		}
		if _, err := tx.CreateBucketIfNotExists(localGroupBucket); err != nil {
			return fmt.Errorf("create bucket %s: %s", localGroupBucket, err)
		}
		if _, err := tx.CreateBucketIfNotExists(orcidBucket); err != nil {
			return fmt.Errorf("create bucket %s: %s", orcidBucket, err)
		}
		return nil
	})

}

// ExportEPrints from highest ID to lowest for cnt. Saves each record in a DB and indexes published ones
func (api *EPrintsAPI) ExportEPrints(count int) error {
	db, err := bolt.Open(api.DBName, 0660, &bolt.Options{Timeout: 1 * time.Second, ReadOnly: false})
	failCheck(err, fmt.Sprintf("Export %s failed to open db, %s", api.DBName, err))
	defer db.Close()

	// Make sure we have a buckets to store things in
	err = initBuckets(db)
	failCheck(err, fmt.Sprintf("Export %s failed to initialize buckets, %s", api.DBName, err))

	uris, err := api.ListEPrintsURI()
	failCheck(err, fmt.Sprintf("Export %s failed, %s", api.URL.String(), err))

	//NOTE: I am sorting the URI by decscending ID number so that the newest articles
	// are exported first
	sort.Sort(byURI(uris))

	uriCount := len(uris)
	if count < 0 {
		count = uriCount
	}
	j := 0 // success count
	k := 0 // error count
	log.Printf("Exporting %d of %d uris", count, uriCount)
	for i := 0; i < uriCount && i < count; i++ {
		uri := uris[i]
		rec, err := api.GetEPrint(uri)
		if err != nil {
			log.Printf("Failed to get %s, %s\n", uri, err)
			k++
		} else {
			rec.URI = strings.TrimPrefix(strings.TrimSuffix(uri, ".xml"), "/rest")
			src, err := json.Marshal(rec)
			if err != nil {
				log.Printf("json.Marshal() failed on %s, %s", uri, err)
				k++
			} else {
				err := db.Update(func(tx *bolt.Tx) error {
					var errs []string
					b := tx.Bucket(ePrintBucket)
					err := b.Put([]byte(rec.URI), src)
					if err == nil {
						// See if we need to add this to the publicationDates index
						if rec.DateType == "published" && rec.Date != "" {
							idx := tx.Bucket(pubDatesBucket)
							dt := normalizeDate(rec.Date)
							err = idx.Put([]byte(fmt.Sprintf("%s%s%s", dt, indexDelimiter, rec.URI)), []byte(rec.URI))
							if err != nil {
								errs = append(errs, fmt.Sprintf("%s", err))
							}
						}
						if len(rec.LocalGroup) > 0 {
							for _, grp := range rec.LocalGroup {
								grp = strings.TrimSpace(grp)
								if len(grp) > 0 {
									idx := tx.Bucket(localGroupBucket)
									err = idx.Put([]byte(fmt.Sprintf("%s%s%s", grp, indexDelimiter, rec.URI)), []byte(rec.URI))
									if err != nil {
										errs = append(errs, fmt.Sprintf("%s", err))
									}
								}
							}
						}
						if len(rec.Creators) > 0 {
							for _, person := range rec.Creators {
								orcid := strings.TrimSpace(person.ORCID)
								if len(orcid) > 0 {
									idx := tx.Bucket(orcidBucket)
									err := idx.Put([]byte(fmt.Sprintf("%s%s%s", orcid, indexDelimiter, rec.URI)), []byte(rec.URI))
									if err != nil {
										errs = append(errs, fmt.Sprintf("%s", err))
									}
								}
							}
						}
						j++
					}
					if len(errs) > 0 {
						return fmt.Errorf("%s", strings.Join(errs, "; "))
					}
					return nil
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

	err = db.View(func(tx *bolt.Tx) error {
		recs := tx.Bucket(ePrintBucket)
		c := recs.Cursor()
		p := 0
		if count < 0 {
			bStats := recs.Stats()
			count = bStats.KeyN
		}
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
			if count < 0 {
				bStats := idx.Stats()
				count = bStats.KeyN
			}
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
			if count < 0 {
				bStats := idx.Stats()
				count = bStats.KeyN
			}
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
			if count < 0 {
				bStats := idx.Stats()
				count = bStats.KeyN
			}
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
			if count < 0 {
				bStats := idx.Stats()
				count = bStats.KeyN
			}
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

// Utility methods used by the LocalGroup and ORCID index related functions
func appendToList(list []string, term string) []string {
	for _, element := range list {
		if strings.Compare(element, term) == 0 {
			return list
		}
	}
	return append(list, term)
}

func firstTerm(s, delimiter string) string {
	r := strings.SplitN(s, delimiter, 2)
	if len(r) > 0 {
		return strings.TrimSpace(r[0])
	}
	return ""
}

// GetLocalGroups returns a JSON list of unique Group names in index
func (api *EPrintsAPI) GetLocalGroups(start, count, direction int) ([]string, error) {
	groupNames := []string{}
	db, err := bolt.Open(api.DBName, 0660, &bolt.Options{Timeout: 1 * time.Second, ReadOnly: true})
	failCheck(err, fmt.Sprintf("GetLocalGroups() %s failed to open db, %s", api.DBName, err))
	defer db.Close()

	switch direction {
	case Ascending:
		err = db.View(func(tx *bolt.Tx) error {
			idx := tx.Bucket(localGroupBucket)
			c := idx.Cursor()
			p := 0
			if count < 0 {
				bStats := idx.Stats()
				count = bStats.KeyN
			}
			for k, _ := c.First(); k != nil && count > 0; k, _ = c.Next() {
				if p >= start {
					grp := firstTerm(fmt.Sprintf("%s", k), indexDelimiter)
					groupNames = appendToList(groupNames, grp)
					count--
				}
				p++
			}
			return nil
		})
	case Descending:
		err = db.View(func(tx *bolt.Tx) error {
			idx := tx.Bucket(localGroupBucket)
			c := idx.Cursor()
			p := 0
			if count < 0 {
				bStats := idx.Stats()
				count = bStats.KeyN
			}
			for k, _ := c.Last(); k != nil && count > 0; k, _ = c.Prev() {
				if p >= start {
					grp := firstTerm(fmt.Sprintf("%s", k), indexDelimiter)
					groupNames = appendToList(groupNames, grp)
					count--
				}
				p++
			}
			return nil
		})
	}
	if err != nil {
		return groupNames, err
	}
	return groupNames, nil
}

// GetLocalGroupRecords returns a list of EPrint records with groupName
func (api *EPrintsAPI) GetLocalGroupRecords(groupName string, start, count, direction int) ([]*Record, error) {
	results := []*Record{}

	db, err := bolt.Open(api.DBName, 0660, &bolt.Options{Timeout: 1 * time.Second, ReadOnly: true})
	failCheck(err, fmt.Sprintf("GetLocalGroupRecords() %s failed to open db, %s", api.DBName, err))
	defer db.Close()

	switch direction {
	case Ascending:
		err = db.View(func(tx *bolt.Tx) error {
			recs := tx.Bucket(ePrintBucket)
			idx := tx.Bucket(localGroupBucket)
			c := idx.Cursor()
			p := 0
			if count < 0 {
				bStats := idx.Stats()
				count = bStats.KeyN
			}
			for k, uri := c.First(); k != nil && count > 0; k, uri = c.Next() {
				if p >= start {
					grp := firstTerm(fmt.Sprintf("%s", k), indexDelimiter)
					if strings.Compare(grp, groupName) == 0 {
						rec := new(Record)
						src := recs.Get([]byte(uri))
						err := json.Unmarshal(src, rec)
						if err != nil {
							return fmt.Errorf("Can't unmarshal %s, %s", uri, err)
						}
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
			idx := tx.Bucket(localGroupBucket)
			c := idx.Cursor()
			p := 0
			if count < 0 {
				bStats := idx.Stats()
				count = bStats.KeyN
			}
			for k, uri := c.Last(); k != nil && count > 0; k, uri = c.Prev() {
				if p >= start {
					grp := firstTerm(fmt.Sprintf("%s", k), indexDelimiter)
					if strings.Compare(grp, groupName) == 0 {
						rec := new(Record)
						src := recs.Get([]byte(uri))
						err := json.Unmarshal(src, rec)
						if err != nil {
							return fmt.Errorf("Can't unmarshal %s, %s", uri, err)
						}
						results = append(results, rec)
						count--
					}
				}
				p++
			}
			return nil
		})
	}
	if err != nil {
		return results, err
	}
	return results, nil
}

// GetORCIDs returns a list unique of ORCID IDs in index
func (api *EPrintsAPI) GetORCIDs(start, count, direction int) ([]string, error) {
	orcids := []string{}
	db, err := bolt.Open(api.DBName, 0660, &bolt.Options{Timeout: 1 * time.Second, ReadOnly: true})
	failCheck(err, fmt.Sprintf("GetORCIDs() %s failed to open db, %s", api.DBName, err))
	defer db.Close()

	switch direction {
	case Ascending:
		err = db.View(func(tx *bolt.Tx) error {
			idx := tx.Bucket(orcidBucket)
			c := idx.Cursor()
			p := 0
			if count < 0 {
				bStats := idx.Stats()
				count = bStats.KeyN
			}
			for k, _ := c.First(); k != nil && count > 0; k, _ = c.Next() {
				if p >= start {
					orcid := firstTerm(fmt.Sprintf("%s", k), indexDelimiter)
					orcids = appendToList(orcids, orcid)
					count--
				}
				p++
			}
			return nil
		})
	case Descending:
		err = db.View(func(tx *bolt.Tx) error {
			idx := tx.Bucket(orcidBucket)
			c := idx.Cursor()
			p := 0
			if count < 0 {
				bStats := idx.Stats()
				count = bStats.KeyN
			}
			for k, _ := c.Last(); k != nil && count > 0; k, _ = c.Prev() {
				if p >= start {
					orcid := firstTerm(fmt.Sprintf("%s", k), indexDelimiter)
					orcids = appendToList(orcids, orcid)
					count--
				}
				p++
			}
			return nil
		})
	}
	if err != nil {
		return orcids, err
	}
	return orcids, nil
}

// GetORCIDRecords returns a list of EPrint records with a given ORCID
func (api *EPrintsAPI) GetORCIDRecords(orcid string, start, count, direction int) ([]*Record, error) {
	results := []*Record{}

	db, err := bolt.Open(api.DBName, 0660, &bolt.Options{Timeout: 1 * time.Second, ReadOnly: true})
	failCheck(err, fmt.Sprintf("GetORCIDRecords() %s failed to open db, %s", api.DBName, err))
	defer db.Close()

	switch direction {
	case Ascending:
		err = db.View(func(tx *bolt.Tx) error {
			recs := tx.Bucket(ePrintBucket)
			idx := tx.Bucket(orcidBucket)
			c := idx.Cursor()
			p := 0
			if count < 0 {
				bStats := idx.Stats()
				count = bStats.KeyN
			}
			for k, uri := c.First(); k != nil && count > 0; k, uri = c.Next() {
				if p >= start {
					term := firstTerm(fmt.Sprintf("%s", k), indexDelimiter)
					if strings.Compare(term, orcid) == 0 {
						rec := new(Record)
						src := recs.Get([]byte(uri))
						err := json.Unmarshal(src, rec)
						if err != nil {
							return fmt.Errorf("Can't unmarshal %s, %s", uri, err)
						}
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
			idx := tx.Bucket(orcidBucket)
			c := idx.Cursor()
			p := 0
			if count < 0 {
				bStats := idx.Stats()
				count = bStats.KeyN
			}
			for k, uri := c.Last(); k != nil && count > 0; k, uri = c.Prev() {
				if p >= start {
					term := firstTerm(fmt.Sprintf("%s", k), indexDelimiter)
					if strings.Compare(term, orcid) == 0 {
						rec := new(Record)
						src := recs.Get([]byte(uri))
						err := json.Unmarshal(src, rec)
						if err != nil {
							return fmt.Errorf("Can't unmarshal %s, %s", uri, err)
						}
						results = append(results, rec)
						count--
					}
				}
				p++
			}
			return nil
		})
	}
	if err != nil {
		return results, err
	}
	return results, nil
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
	fname = path.Join(api.TemplatePath, "rss.xml")
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
	fname = path.Join(api.TemplatePath, "page.include")
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
		path.Join(api.TemplatePath, "page.include"),
		path.Join(api.TemplatePath, "page.html"),
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

// BuildPages generates a webpages based on the contents of the exported EPrints data.
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
		return fmt.Errorf("No records found for %q %s", title, docPath)
	}
	log.Printf("%d records found.", len(records))
	if err := api.RenderDocuments(title, fmt.Sprintf("Building pages 0 to %d descending", feedSize), target, records); err != nil {
		return fmt.Errorf("%q %s error, %s", title, docPath, err)
	}
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
	log.Printf("Building Recently Published")
	err := api.BuildPages(feedSize, "Recently Published", "recent/published", func(api *EPrintsAPI, start, count, direction int) ([]*Record, error) {
		return api.GetPublishedRecords(0, feedSize, Descending)
	})
	if err != nil {
		return err
	}

	// Collect the published articles
	log.Printf("Building Recent Articles")
	err = api.BuildPages(feedSize, "Recent Articles", "recent/articles", func(api *EPrintsAPI, start, count, direction int) ([]*Record, error) {
		return api.GetPublishedArticles(start, count, direction)
	})
	if err != nil {
		return err
	}

	// Collect EPrints by orcid ID and publish
	log.Printf("Building ORCID works")
	orcids, err := api.GetORCIDs(0, -1, Ascending)
	if err != nil {
		return err
	}
	log.Printf("Found %d orcids", len(orcids))
	for _, orcid := range orcids {
		err = api.BuildPages(feedSize, fmt.Sprintf("ORCID: %s", orcid), fmt.Sprintf("orcid/%s", orcid), func(api *EPrintsAPI, start, count, direction int) ([]*Record, error) {
			return api.GetORCIDRecords(orcid, start, count, direction)
		})
		if err != nil {
			return err
		}
	}

	// Collect EPrints by Group/Affiliation
	log.Printf("Building Local Groups")
	groupNames, err := api.GetLocalGroups(0, -1, Ascending)
	if err != nil {
		return err
	}
	log.Printf("Found %d groups", len(groupNames))
	for _, groupName := range groupNames {
		err = api.BuildPages(feedSize, fmt.Sprintf("%s", groupName), fmt.Sprintf("affiliation/%s", groupName), func(api *EPrintsAPI, start, count, direction int) ([]*Record, error) {
			return api.GetLocalGroupRecords(groupName, start, count, direction)
		})
		if err != nil {
			return err
		}
	}
	return nil
}
