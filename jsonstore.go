//
// jsonstore.go holds the operations for openning and close a JSON
// Document Store currently implemented in MySQL 8 using JSON columns.
//
package eprinttools

import (
	"database/sql"
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

// SaveJSONDocument takes a configuration, repoName, eprint id as integer and
// JSON source saving it to the appropriate JSON table.
func SaveJSONDocument(cfg *Config, repoName string, id int, src []byte, action string, lastmod string, status string) error {
	type Doc struct {
		ID           int    `json:"id"`
		Src          []byte `json:"src"`
		Action       string `json:"action,omitempty"`
		LastModified string `json:"lastmod,omitempty"`
		Status       string `json:"status,omitempty"`
	}

	stmt := fmt.Sprintf(`REPLACE INTO %s (id, src, action, lastmod, status) VALUES (?, ?, ?, ?, ?)`, repoName)
	doc := new(Doc)
	doc.ID = id
	doc.Src = src
	doc.Action = action
	doc.LastModified = lastmod
	doc.Status = status
	if cfg.Jdb == nil {
		OpenJSONStore(cfg)
	}
	_, err := cfg.Jdb.Exec(stmt, doc.ID, doc.Src, doc.Action, doc.LastModified, doc.Status)
	if err != nil {
		return fmt.Errorf("sql failed for %d in %s, %s", id, repoName, err)
	}
	return nil
}

// ReadJSONDocument takes a configuration, repoName, eprint id and returns the JSON
// source document.
func GetJSONDocument(cfg *Config, repoName string, id int) ([]byte, error) {
	if cfg.Jdb == nil {
		OpenJSONStore(cfg)
	}
	stmt := fmt.Sprintf("SELECT id, src, action, created, updated FROM %s WHERE id = ?", repoName)
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
