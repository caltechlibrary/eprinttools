//
// jsonstore.go holds the operations for openning and close a JSON
// Document Store currently implemented in MySQL 8 using JSON columns.
//
package eprinttools

import (
	"database/sql"
	"encoding/json"
	"fmt"

	// MySQL database support
	_ "github.com/go-sql-driver/mysql"
)

// OpenJSONStore
func OpenJSONStore(config *Config) error {
	if config.JSONStore == "" {
		return fmt.Errorf("JSONStore is not set")
	} else {
		// Setup DB connection for target repository
		db, err := sql.Open("mysql", config.JSONStore)
		if err != nil {
			return fmt.Errorf("Could not open MySQL connection for %s, %s", config.JSONStore, err)
		}
		config.Jdb = db
	}
	return nil
}

// CloseJSONStore
func CloseJSONStore(config *Config) error {
	if config.Jdb != nil {
		if err := config.Jdb.Close(); err != nil {
			return fmt.Errorf("Failed to close %s, %s", config.JSONStore, err)
		}
	}
	return nil
}

// Doc holds the data structure of our jsonstore row.
type Doc struct {
		ID           int    `json:"id"`
		Src          []byte `json:"src"`
		Action       string `json:"action,omitempty"`
		Created      string `json:"created,omitempty"`
		LastModified string `json:"lastmod,omitempty"`
		PubDate      string `json:"pubDate,omitempty"`
		Status       string `json:"status,omitempty"`
		IsPublic     bool `json:"is_public,omitempty"`
		RecordType string `json:"record_type,omitempty"`
}

// SaveJSONDocument takes a configuration, repoName, eprint id as integer and
// JSON source saving it to the appropriate JSON table.
func SaveJSONDocument(cfg *Config, repoName string, id int, src []byte, action string, created string, lastmod string, pubdate string, status string, isPublic bool, recordType string) error {

	stmt := fmt.Sprintf(`REPLACE INTO %s (id, src, action, created, lastmod, pubdate, status, is_public, record_type) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, repoName)
	doc := new(Doc)
	doc.ID = id
	doc.Src = src
	doc.Action = action
	doc.Created = created
	doc.LastModified = lastmod
	doc.PubDate = pubdate
	doc.Status = status
	doc.IsPublic = isPublic
	doc.RecordType = recordType
	if cfg.Jdb == nil {
		OpenJSONStore(cfg)
	}
	_, err := cfg.Jdb.Exec(stmt, doc.ID, doc.Src, doc.Action, doc.Created, doc.LastModified, doc.PubDate, doc.Status, doc.IsPublic, doc.RecordType)
	if err != nil {
		return fmt.Errorf("sql failed for %d in %s, %s", id, repoName, err)
	}
	return nil
}

// GetJSONDocument takes a configuration, repoName, eprint id and returns
// the JSON source document.
func GetJSONDocument(cfg *Config, repoName string, id int) ([]byte, error) {
	if cfg.Jdb == nil {
		OpenJSONStore(cfg)
	}
	stmt := fmt.Sprintf("SELECT src FROM %s WHERE id = ? LIMIT 1", repoName)
	rows, err := cfg.Jdb.Query(stmt, id)
	if err != nil {
		return nil, fmt.Errorf("Failed to get %s id %d, %s", repoName, id, err)
	}
	var src []byte
	for rows.Next() {
		if err := rows.Scan(src); err != nil {
			return nil, fmt.Errorf("Failed to get row in %s for id %d, %s", repoName, id, err)
		}
	}
	rows.Close()
	return src, nil
}


// GetJSONRow takes a configuration, repoName, eprint id and returns
// the table row as JSON source.
func GetJSONRow(cfg *Config, repoName string, id int) ([]byte, error) {
	if cfg.Jdb == nil {
		OpenJSONStore(cfg)
	}
	stmt := fmt.Sprintf("SELECT id, src, action, created, lastmod, pubdate, status, is_public, record_type FROM %s WHERE id = ?", repoName)
	rows, err := cfg.Jdb.Query(stmt, id)
	if err != nil {
		return nil, fmt.Errorf("Failed to get %s id %d, %s", repoName, id, err)
	}
	doc := new(Doc)
	for rows.Next() {
		if err := rows.Scan(doc.ID, doc.Src, doc.Action, doc.Created, doc.LastModified, doc.PubDate, doc.Status, doc.IsPublic, doc.RecordType); err != nil {
			return nil, fmt.Errorf("Failed to get row in %s for id %d, %s", repoName, id, err)
		}
	}
	rows.Close()
	return json.MarshalIndent(doc, "", "    ")
}

