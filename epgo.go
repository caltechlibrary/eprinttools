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
	"time"

	// 3rd Party packages
	"github.com/boltdb/bolt"
)

const (
	// Version is the revision number for this implementation of epgo
	Version = "0.0.0"

	// Ascending sorts from lowest (oldest) to highest (newest)
	Ascending = iota
	// Descending sorts from highest (newest) to lowest (oldest)
	Descending = iota
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
)

func failCheck(err error, msg string) {
	if err != nil {
		log.Fatalf("%s\n", msg)
	}
}

// EPrintsAPI holds the basic connectin information to read the REST API for EPrints
type EPrintsAPI struct {
	URL       *url.URL `json:"base_url"`  // EPGO_BASE_URL
	DBName    string   `json:"dbname"`    // EPGO_DBNAME
	Htdocs    string   `json:"htdocs"`    // EPGO_HTDOCS
	Templates string   `json:"templates"` // EPGO_TEMPLATES
}

// Name returns the contents of eprint>creators>item>name as a struct
type Name struct {
	XMLName xml.Name `json:"-"`
	Given   string   `xml:"given"`
	Family  string   `xml:"family"`
}

// Record returns a structure that can be converted to JSON easily
type Record struct {
	XMLName            xml.Name `json:"-"`
	Title              string   `xml:"eprint>title" json:"title"`
	ID                 int      `xml:"eprint>eprintid" json:"id"`
	RevNumber          int      `xml:"eprint>rev_number" json:"rev_number"`
	UserID             int      `xml:"eprint>userid" json:"userid"`
	Dir                string   `xml:"eprint>dir" json:"eprint_dir"`
	Datestamp          string   `xml:"eprint>datestamp" json:"datestamp"`
	LastModified       string   `xml:"eprint>lastmod" json:"lastmod"`
	StatusChange       string   `xml:"eprint>status_changed" json:"status_changed"`
	Type               string   `xml:"eprint>type" json:"type"`
	MetadataVisibility string   `xml:"eprint>metadata_visibility" json:"metadata_visibility"`
	Creators           []*Name  `xml:"eprint>creators>item>name" json:"creators"`
	IsPublished        string   `xml:"eprint>ispublished" json:"ispublished"`
	Subjects           []string `xml:"eprint>subjects>item" json:"subjects"`
	FullTextStatus     string   `xml:"eprint>full_text_status" json:"full_text_status"`
	Date               string   `xml:"eprint>date" json:"data"`
	DateType           string   `xml:"eprint>date_type" json:"date_type"`
	Publication        string   `xml:"eprint>publication" json:"publication"`
	Volume             string   `xml:"eprint>volume" json:"volume"`
	Number             string   `xml:"eprint>number" json:"number"`
	PageRange          string   `xml:"eprint>pagerange" json:"pagerange"`
	IDNumber           string   `xml:"eprint>id_number" json:"id_number"`
	Referred           bool     `xml:"eprint>refereed" json:"refereed"`
	ISSN               string   `xml:"eprint>issn" json:"issn"`
	OfficialURL        string   `xml:"eprint>official_url" json:"official_url"`
}

type ePrintIDs struct {
	XMLName xml.Name `xml:"html" json:"-"`
	IDs     []string `xml:"body>ul>li>a" json:"ids"`
}

// String renders the Record as a pretty printed JSON object
func (r *Record) String() string {
	src, _ := json.MarshalIndent(r, "", "  ")
	return string(src)
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
	baseURL := os.Getenv("EPGO_BASE_URL")
	htdocs := os.Getenv("EPGO_HTDOCS")
	dbName := os.Getenv("EPGO_DBNAME")
	templates := os.Getenv("EPGO_TEMPLATES")

	if baseURL == "" {
		return nil, fmt.Errorf("Environment not configured, missing EPGO_BASE_URL")
	}
	api := new(EPrintsAPI)
	api.URL, err = url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("EPGO_BASE_URL malformed %s, %s", baseURL, err)
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

// ExportEPrints gets a list of all EPrints and then saves each record in a DB
func (api *EPrintsAPI) ExportEPrints() error {
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

	j := 0 // success count
	k := 0 // error count
	log.Printf("Exporting of %d uris", len(uris))
	for i, uri := range uris {
		rec, err := api.GetEPrint(uri)
		if err != nil {
			log.Printf("Failed to get %s, %s\n", uri, err)
			k++
		} else {
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
					}
					return err
				})
				if err != nil {
					log.Printf("Failed to save eprint %s, %s\n", uri, err)
					k++
				}
			}
		}
		if (i % 10) == 0 {
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
					results = append(results, rec)
					count--
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
					results = append(results, rec)
					count--
				}
				p++
			}
			return nil
		})
	}
	return results, err
}

// BuildSite generates a website based on the contents of the exported EPrints data.
// The site builder needs to know the name of the BoltDB, the root directory
// for the website and directory to find the templates
func (api *EPrintsAPI) BuildSite() error {
	// Collect the published records
	records, err := api.GetPublishedRecords(0, 100, Descending)
	if err != nil {
		return fmt.Errorf("Can't get published records, %s", err)
	}
	if len(records) == 0 {
		return fmt.Errorf("No published records found")
	}
	fmt.Printf("DEBUG records: %+v", records)
	// FIXME: Build feeds and HTML page for published EPrints
	// pageHTML, err := ioutil.ReadFile(path.Join(api.Templates, "page.html"))
	// if err != nil {
	// 	return fmt.Errorf("Can't open template %s/page.html, %s", api.Templates, err)
	// }
	// pageInclude, err := ioutil.ReadFile(path.Join(api.Templates, "page.include"))
	// if err != nil {
	// 	return fmt.Errorf("Can't open template %s/page.include, %s", api.Templates, err)
	// }
	// fmt.Printf("DEBUG pageHTML: %s\n", pageHTML)
	// fmt.Printf("DEBUG pageInclude: %s\n", pageInclude)
	//
	// pageTmpl, err := template.New("page.html").Parse(string(pageHTML))
	// if err != nil {
	// 	return fmt.Errorf("Can't parse %s/page.html, %s", api.Templates, err)
	// }
	// includeTmpl, err := template.Must(pageTmpl.Clone()).Parse(string(pageInclude))
	// if err != nil {
	// 	return fmt.Errorf("Can't parse %s/page.include, %s", api.Templates, err)
	// }
	// // Write out .html page
	// if err = pageTmpl.Execute(os.Stdout, records); err != nil {
	// 	return fmt.Errorf("Can't execute %s/page.html, %s", api.Templates, err)
	//
	// }
	// // Write out .include page
	// if err = includeTmpl.Execute(os.Stdout, records); err != nil {
	// 	return fmt.Errorf("Can't execute %s/page.include, %s", api.Templates, err)
	//
	// }
	return nil
}
