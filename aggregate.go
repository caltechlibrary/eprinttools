package eprinttools

import (
	"database/sql"
	"sort"
	"log"
	"time"
	"os"
	"strings"
	"fmt"
)

/*
This file should define and orchestrate aggregating EPrint content
into "view" (e.g. by record type, by pub date, by author, by editor, etc).
It needs to work across all our repositories and be general enough for
an easy fit when we've migrated to a new repository system.
*/

// Aggregation holds the connection to the SQL/JSON store where aggregations are
// maintained.
type Aggregation struct {
	DSN string  `json:"dsn,omitempty"`
	db  *sql.DB `json:"-"`
}

type AggregateItem struct {
	Repository  string `json:"repository"`
	EPrintID    int    `json:"eprintid"`
	ViewName        string `json:"view_name"`
	Created     string `json:"created,omitempty"`
	PubDate     string `json:"pubdate,omitempty"`
	IsPublic    string `json:"is_public,omitempty"`
	Citation map[string]interface{} `json:"citation,omitempty"`
	FilterObject string `json:"filter_object,omitempty"`
	SortObject string `json:"sort_object,omitempty"`
}

// AggregationInitDB returns SQL statements for creating
// the tables and database for the harvester based on the
// initialization file provides.
func AggregationInitDB(cfgName string) (string, error) {
	now := time.Now()
	appName := os.Args[0]
	database := `--
--
-- Database generated for MySQL 8 by %s %s
-- using %s on %s
--
CREATE DATABASE IF NOT EXISTS %s;
USE %s;

`

	table := `--
-- Table Schema generated for MySQL 8 by %s %s
-- for EPrint aggregeation %s on %s
--
CREATE TABLE IF NOT EXISTS %s (
	aggregate_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	repository VARCHAR(255),
	eprintid INT,
    view_name VARCHAR(255) DEFAULT "",
  	created VARCHAR(255) DEFAULT "",
  	pubdate VARCHAR(255) DEFAULT "",
    is_public BOOLEAN DEFAULT FALSE,
	citation_object JSON DEFAULT "",
	filter_object JSON DEFAULT "",
	sort_object JSON DEFAULT ""
);
`
	cfg, err := LoadConfig(cfgName)
	if err != nil {
		return "", err
	}
	dbName, err := getDBName(cfg.AggregationStore)
	if err != nil {
		return "", fmt.Errorf("Cannot parse aggregation DSN in %s, %s", cfgName, err)
	}
	src := []string{}
	if dbName != "" {
		src = append(src, fmt.Sprintf(database, appName, Version, cfgName, now.Format(mysqlTimeFmt), dbName, dbName))
	}
	for repoName := range cfg.Repositories {
		src = append(src, fmt.Sprintf(table, appName, Version, repoName, now.Format(mysqlTimeFmt), repoName))
	}
	return strings.Join(src, "\n"), nil
}

// RunAggregator will use the config file names by cfgName and
// the start and end time strings if set to retrieve all eprint
// records created or modified during that time sequence.
func RunAggregator(cfgName string, start string, end string, verbose bool) error {
        now := time.Now()
        // Read in the configuration for this harvester instance.
        cfg, err := LoadConfig(cfgName)
        if err != nil {
                return err
        }
        if cfg == nil {
                return fmt.Errorf("Could not create a configuration object")
        }
        if start == "" {
                // Pick a start date/time before EPrints existed.
                start = "2000-01-01 00:00:00"
        }
        if end == "" {
                // Pick the current date time when harvest is starting
                end = now.Format(mysqlTimeFmt)
        }
        return aggregateRepositories(cfg, start, end, verbose)
}

// aggregateRepositories implements an aggreation instance.
func aggregateRepositories(cfg *Config, start string, end string, verbose bool) error {
	aggregation, err := OpenAggregation(cfg)
	if err != nil {
		return err
	}
	defer aggregation.Close()
        if err := OpenConnections(cfg); err != nil {
                return err
        }
        repoNames := []string{}
        for repoName := range cfg.Connections {
                repoNames = append(repoNames, repoName)
        }
        sort.Strings(repoNames)
        for _, repoName := range repoNames {
                //FIXME: we could use a go routine to support concurrent harvests.
                log.Printf("aggretating %s started", repoName)
                if err := aggregateRepository(aggregation, cfg, repoName, start, end, verbose); err != nil {
                        log.Printf("aggregating %s aborted", repoName)
                        return err
                }
                log.Printf("aggregating %s completed", repoName)
        }
        return CloseConnections(cfg)
}

// aggregateRepository takes a configuration with open database connections.
// a repository name (i.e. eprint database name) along with a start and
// end timestamp. It harvests records created/modified from the repository
// in time range.
func aggregateRepository(aggregation *Aggregation, cfg *Config, repoName string, start string, end string, verbose bool) error {
        createdIDs, err := GetEPrintIDsInTimestampRange(cfg, repoName, "datestamp", start, end)
        if err != nil {
                return err
        }
        if verbose {
                log.Printf("Retrieved %d keys based on creation date", len(createdIDs))
        }
        modifiedIDs, err := GetEPrintIDsInTimestampRange(cfg, repoName, "lastmod", start, end)
        if err != nil {
                return err
        }
        if verbose {
                log.Printf("Retrieved %d keys based on modified date", len(modifiedIDs))
        }
        ids := append(createdIDs, modifiedIDs...)
        ids = getSortedUniqueIDs(ids)
        if verbose {
                log.Printf("Processing %d unique keys", len(ids))
        }
        for i, id := range ids {
                // FIXME: do we want to show a progress bar or just errors?
                err := aggregateEPrintRecord(aggregation, cfg, repoName, id)
                if err != nil {
                        log.Printf("Harvesting EPrint %d (%d/%d) failed, %s", id, i, len(ids), err)
                }
                if verbose && ((i % 1000) == 0) {
                        log.Printf("Harvested EPrint %d (%d/%d)", id, i, len(ids))
                }
        }
        return nil
}

// aggregateEPrintRecord takes a configuration with open database connections
// a repository name and EPrint record ID and populates a aggretation datastore.
func aggregateEPrintRecord(aggregation *Aggregation, cfg *Config, repoName string, eprintID int) error {
        ds, ok := cfg.Repositories[repoName]
        if !ok {
                return fmt.Errorf("Data Source not found for %q looking up eprint %d", repoName, eprintID)
        }
        eprint, err := SQLReadEPrint(cfg, repoName, ds.BaseURL, eprintID)
        if err != nil {
                return err
        }
		return aggregation.Aggregate(eprint)
}



// OpenAggregation takes a filename for JSON settings file and an attribute name holding
// the DSN for maintaing aggregations in a SQL/JSON store.
func OpenAggregation(cfg *Config) (*Aggregation, error) {
	aggregation := new(Aggregation)
	aggregation.DSN = cfg.AggregationStore
	if cfg.Adb == nil {
		dbName, err := getDBName(cfg.AggregationStore)
		if err != nil {
			return nil, fmt.Errorf("Cannot parse aggregation DSN in %s, %s", aggregation.DSN, err)
		}
		db, err := sql.Open("mysql", dbName)
		if err != nil {
			return nil, fmt.Errorf("Cannot open %s, %s", dbName, err)
		}
		cfg.Adb = db
	}
	aggregation.db = cfg.Adb
	return aggregation, nil
}

// Close closes an aggregation.
func (aggregation *Aggregation) Close() error {
	aggregation.DSN = ""
	if aggregation.db != nil {
		return aggregation.db.Close()
	}
	return nil
}

// Aggregate takes an eprint record and generates the related aggregations
// that the eprint record could fall into.
//
// ```
//
//	aggregation, err := eprinttools.OpenAggregation(cfg)
//	if err err != nil {
//	      // ... handle error ...
//	}
//	defer aggregation.Close()
//	eprint := eprinttools.GetEPrint("http://eprint.example.edu", 12122)
//	if err = aggregation.Aggregate(eprint); err != nil {
//	      // ... handle error ...
//	}
//
// ```
func (aggregation *Aggregation) Aggregate(eprint *EPrint) error {
	// First prune any references to the EPrint record 
	// SQL: DELETE FROM %q WHERE repository = ? AND eprintid = ?;

	// Analyze EPrint record generating the aggregation rows
	// for each of our List types.

	return fmt.Errorf("aggregation.Aggregate() not implemented")
}

func (aggretation *Aggregation) ListByType() ([]*AggregateItem, error) {
	return nil, fmt.Errorf("aggregation.ByType() not implemented")
}

func (aggregation *Aggregation) ListByPubDate() ([]*AggregateItem, error) {
	return nil, fmt.Errorf("aggregation.ByPubDate() not implemented")
}

func (aggregation *Aggregation) ListByDateAdded() ([]*AggregateItem, error) {
	return nil, fmt.Errorf("aggregation.ByDateAdded() not implemeneted")
}

func (aggregation *Aggregation) ListByPerson() ([]*AggregateItem, error) {
	/* this should return list of cannonical person ids, this list can then be used to create individual lists by roles */
	return nil, fmt.Errorf("aggregation.ByPerson() not implemented")
}

func (aggregation *Aggregation) ListrByPersonAndRole() ([]*AggregateItem, error) {
	/* This aggregation should be broken down by person and role, it should include the following
	   sub lists
	   - author
	   - editor
	   - contributor
	   - advisor
	   - commitee memeber
	   - * would represent all roles (i.e. a combined list)
	*/
	return nil, fmt.Errorf("aggregation.ByPersonAndRole() not implemented")
}

func (aggregation *Aggregation) ListByPublication() ([]*AggregateItem, error) {
	return nil, fmt.Errorf("aggregation.ByPublication() not implemented")
}

func (aggregation *Aggregation) ListByGroup() ([]*AggregateItem, error) {
	return nil, fmt.Errorf("aggregation.ByGroup() not implemented")
}

func (aggregation *Aggregation) ListBySubject() ([]*AggregateItem, error) {
	return nil, fmt.Errorf("aggregeation.BySubject() not implemented")
}

func (aggregation *Aggregation) ListByCategory() ([]*AggregateItem, error) {
	return nil, fmt.Errorf("aggregation.ByCategory() not implemented")
}

func (aggregation *Aggregation) ListByConference() ([]*AggregateItem, error) {
	return nil, fmt.Errorf("aggregation.ByConference not implemented")
}

func (aggregation *Aggregation) ListByCollection() ([]*AggregateItem, error) {
	return nil, fmt.Errorf("aggregation.ByCollection() not implemented")
}

func (aggregation *Aggregation) ListByAuthor() ([]*AggregateItem, error) {
	return nil, fmt.Errorf("aggregation.ByAuthor() not implemented")
}

func (aggregation *Aggregation) ListByEditor() ([]*AggregateItem, error) {
	return nil, fmt.Errorf("aggregation.ByEditor() not implemented")
}

func (aggregation *Aggregation) ListByContributor() ([]*AggregateItem, error) {
	return nil, fmt.Errorf("aggregation.ByContributor() not implemented")
}

func (aggregation *Aggregation) ListByAdvisor() ([]*AggregateItem, error) {
	return nil, fmt.Errorf("aggregation.ByAdvisor() not implemented")
}

func (aggregation *Aggregation) ListByCommitteeMember() ([]*AggregateItem, error) {
	return nil, fmt.Errorf("aggregation.ByCommitteeMember() not implemented")
}

func (aggregation *Aggregation) ListByDegreeThesisType() ([]*AggregateItem, error) {
	return nil, fmt.Errorf("aggregation.ByThesisType() not implemented")
}

func (aggregation *Aggregation) ListByOption() ([]*AggregateItem, error) {
	return nil, fmt.Errorf("aggregation.ByOption() not implemented")
}


