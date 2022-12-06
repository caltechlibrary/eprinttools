package eprinttools

//
// Service configuration management
//

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Config holds a configuration file structure used by EPrints Extended API
// Configuration file is expected to be in JSON format.
type Config struct {
	// Hostname for running service
	Hostname string `json:"hostname"`

	// BaseURL is the base URL passed to any client configuration
	BaseURL string `json:"base_url"`

	// Logfile
	Logfile string `json:"logfile,omitempty"`

	// Repositories are defined by a REPO_ID (string)
	// that points at a MySQL Db connection string
	Repositories map[string]*DataSource `json:"eprint_repositories"`

	// Connections is a map to database connections
	Connections map[string]*sql.DB `json:"-"`

	// JSONStore is the name of a MySQL 8 database in DNS format
	// that holds for each repository. The tables have two columns
	// eprint id (INTEGER) and document (JSON COLUMNS).
	// The JSONStore is where data is harvested into and where it is
	// staged for writing out to a published Object store like S3.
	JSONStore string `json:"jsonstore"`

	// Jdb holds the MySQL connector to the jsonstore
	Jdb *sql.DB `json:"-"`

	// Routes holds the mapping of end points to repository id
	// instances.
	Routes map[string]map[string]func(http.ResponseWriter, *http.Request, string, []string) (int, error) `json:"-"`

	// ProjectDir is the directory where you stage harvested content
	ProjectDir string `json:"project_dir, omitempty"`

	// Htdocs is the directory where aggregated information and
	// website content is generated to after running the harvester.
	Htdocs string `json:"htdocs,omitempty"`

	// PandocServer is the URL to the Pandoc server
	// E.g. localhost:8080
	PandocServer string `json:"pandoc_server,omitempty"`
}

// DataSource can contain one or more types of datasources. E.g.
// E.g. dsn for MySQL connections and also data for REST API access.
type DataSource struct {
	// DSN is used to connect to a MySQL style DB.
	DSN string `json:"dsn,omitempty"`

	// BaseURL is the URL to use in constructing eprint id, document id
	// and file id attribute strings.
	BaseURL string `json:"base_url,omitempty"`

	// Rest is used to connect to EPrints REST API
	// NOTE: assumes Basic Auth for authentication
	RestAPI string `json:"rest,omitempty"`

	// Write enables the write API for creating
	// or replacing EPrint records via SQL database calls.
	// The default value is false.
	Write bool `json:"write" default:"false"`

	// DefaultCollection
	DefaultCollection string `json:"default_collection,omitempty"`

	// DefaultOfficialURL
	DefaultOfficialURL string `json:"default_official_url,omitempty"`

	// DefaultRights (i.e. usage statement)
	DefaultRights string `json:"default_rights,omitempty"`

	// DefaultIsRefereed (i.e. refereed field applied for "article" types.
	// Caltech Library defaults this to "TRUE" for type "article".
	DefaultRefereed string `json:"default_refereed,omitempty"`

	// DefaultStatus is the eprint.eprint_status value to set by default
	// on creating new eprint records. Normally this is "inbox" or "buffer"
	DefaultStatus string `json:"default_status,omitempty"`

	// StripTags bool is true then an EPrint Abstract will have XML/HTML tags
	// stripped on import.
	StripTags bool `json:"strip_tags,omitempty"`

	// TableMap holds the mapping of tables and columns for
	// the repository presented.
	TableMap map[string][]string `json:"tables,omitempty"`

	// PublicOnly is a boolean indicating if the "harvested" content
	// should be restricted to public records.
	PublicOnly bool `json:"is_public,omitempty"`
}

func DefaultConfig() []byte {
	config := new(Config)
	config.Hostname = "localhost:8484"
	config.BaseURL = "http//localhost:8484"
	config.JSONStore = "$DB_USER:$DB_PASSWORD@/collections"
	config.Htdocs = "htdocs"
	config.ProjectDir = "."
	config.PandocServer = "localhost:3030"
	repo := new(DataSource)
	repo.DSN = `$DB_USER:$DB_PASSWORD@/authors`
	repo.BaseURL = `http://authors.example.edu`
	repo.Write = false
	repo.DefaultCollection = `authors`
	repo.DefaultRights = "No commercial reproduction, distribution, display or performance rights in this work are provided."
	repo.DefaultRefereed = "TRUE"
	repo.DefaultStatus = "inbox"
	repo.StripTags = true
	repo.PublicOnly = true
	config.Repositories = map[string]*DataSource{
		"authors": repo,
	}
	src, _ := json.MarshalIndent(config, "", "     ")
	return src
}

// LoadConfig reads a JSON file and returns a Config structure
// or error.
func LoadConfig(fname string) (*Config, error) {
	config := new(Config)
	config.Repositories = map[string]*DataSource{}
	if src, err := ioutil.ReadFile(fname); err != nil {
		return nil, err
	} else {
		// Since we should be OK, unmarshal in into active config
		if err = jsonDecode(src, &config); err != nil {
			return nil, fmt.Errorf("Unmarshaling %q failed, %s", fname, err)
		}
		if config.Hostname == "" {
			config.Hostname = "localhost:8484"
		}
		if config.BaseURL == "" {
			config.BaseURL = fmt.Sprintf("http://%s", config.Hostname)
		}
		if config.Htdocs == "" {
			config.Htdocs = "htdocs"
		}
	}
	return config, nil
}
