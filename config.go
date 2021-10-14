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

	// Repositories are defined by a REPO_ID (string)
	// that points at a MySQL Db connection string
	Repositories map[string]*DataSource `json:"repositories"`

	// Connections is a map to database connections
	Connections map[string]*sql.DB `json:"-"`

	// Routes holds the mapping of end points to repository id
	// instances.
	Routes map[string]map[string]func(http.ResponseWriter, *http.Request, string, []string) (int, error) `json:"-"`
}

// DataSource can contain one or more types of datasources. E.g.
// E.g. dsn for MySQL connections and also data for REST API access.
type DataSource struct {
	// DSN is used to connect to a MySQL style DB.
	DSN string `json:"dsn,omitempty"`
	// Rest is used to connect to EPrints REST API
	// NOTE: assumes Basic Auth for authentication
	RestAPI string `json:"rest,omitempty"`
	// ReadWrite enables the API read/write API for creating
	// or replacing EPrint records via SQL database calls.
	// The default value is false.
	ReadWrite bool `json:"create" default:"false"`
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
		if err = json.Unmarshal(src, &config); err != nil {
			return nil, fmt.Errorf("Unmarshaling %q failed, %s", fname, err)
		}
		if config.Hostname == "" {
			config.Hostname = "localhost:8484"
		}
	}
	return config, nil
}
