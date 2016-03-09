package epgo

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

const (
	// Version is the revision number for this implementation of epgo
	Version = "0.0.0"
)

func failCheck(err error, msg string) {
	if err != nil {
		log.Fatalf("%s\n", msg)
	}
}

// EPrintsAPI holds the basic connectin information to read the REST API for EPrints
type EPrintsAPI struct {
	URL      *url.URL `json:"base_url"`
	Username string   `json:"username"`
	Password string   `json:"password"`
}

// Name returns the contents of eprint>creators>item>name as a struct
type Name struct {
	XMLName xml.Name
	Given   string `xml:"given"`
	Family  string `xml:"family"`
}

// Record returns a structure that can be converted to JSON easily
type Record struct {
	XMLName            xml.Name
	Title              string   `xml:"eprint>title"`
	ID                 int      `xml:"eprint>eprintid"`
	RevNumber          int      `xml:"eprint>rev_number"`
	UserID             int      `xml:"eprint>userid"`
	Dir                string   `xml:"eprint>dir"`
	Datestamp          string   `xml:"eprint>datestamp"`
	LastModified       string   `xml:"eprint>lastmod"`
	StatusChange       string   `xml:"eprint>status_changed"`
	Type               string   `xml:"eprint>type"`
	MetadataVisibility string   `xml:"eprint>metadata_visibility"`
	Creators           []*Name  `xml:"eprint>creators>item>name"`
	IsPublished        string   `xml:"eprint>ispublished"`
	Subjects           []string `xml:"eprint>subjects>item"`
	FullTextStatus     string   `xml:"eprint>full_text_status"`
	Date               string   `xml:"eprint>date"`
	DateType           string   `xml:"eprint>date_type"`
	Publication        string   `xml:"eprint>publication"`
	Volume             string   `xml:"eprint>volume"`
	Number             string   `xml:"eprint>number"`
	PageRange          string   `xml:"eprint>pagerange"`
	IDNumber           string   `xml:"eprint>id_number"`
	Referred           bool     `xml:"eprint>refereed"`
	ISSN               string   `xml:"eprint>issn"`
	OfficialURL        string   `xml:"eprint>official_url"`
}

type ePrintIDs struct {
	XMLName xml.Name `xml:"html"`
	IDs     []string `xml:"body>ul>li>a"`
}

// New creates a new API instance
func New() (*EPrintsAPI, error) {
	var err error
	baseURL := os.Getenv("EPGO_BASE_URL")
	username := os.Getenv("EPGO_USERNAME")
	password := os.Getenv("EPGO_PASSWORD")
	if baseURL == "" {
		return nil, fmt.Errorf("Environment not configured, missing EPGO_BASE_URL")
	}
	api := new(EPrintsAPI)
	api.URL, err = url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("EPGO_BASE_URL malformed %s, %s", baseURL, err)
	}
	api.Username = username
	api.Password = password
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
