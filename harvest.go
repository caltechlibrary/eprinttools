package eprinttools

//
// This package implements an EPrints 3.x harvester storing
// Each harvested EPrints repository's JSON metadata in a
// MySQL 8 table using JSON column type.
//

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	// Aliasing mysql driver to get to ParseDSN
	mysqlDriver "github.com/go-sql-driver/mysql"
)

const (
	mysqlTimeFmt = "2006-01-02 15:04:05"
)

// getDBName uses the ParseDSN function from the MySQL driver to
// return the DB name.
func getDBName(dsn string) (string, error) {
	cfg, err := mysqlDriver.ParseDSN(dsn)
	if err != nil {
		return "", err
	}
	return cfg.DBName, nil
}

// HarvesterInitDB returns SQL statements for creating
// the tables and database for the harvester based on the
// initialization file provides.
func HarvesterInitDB(cfgName string) (string, error) {
	now := time.Now()
	appName := os.Args[0]
	database := `--
--
-- Database generated for MySQL 8 by %s %s
-- using %s on %s
--
-- DROP DATABASE IF EXISTS %s;
CREATE DATABASE IF NOT EXISTS %s;
USE %s;

-- Table Schema generated for MySQL 8
-- for all EPrint repositories as _aggregate_creator
CREATE TABLE IF NOT EXISTS _aggregate_creator (
	aggregate_id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	repository VARCHAR(255),
	collection VARCHAR(255),
	eprintid INT,
	person_id VARCHAR(255) DEFAULT ""
);

-- Table Schema generated for MySQL 8
-- for all EPrint repositories as _aggregate_editor
CREATE TABLE IF NOT EXISTS _aggregate_editor (
	aggregate_id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	repository VARCHAR(255),
	collection VARCHAR(255),
	eprintid INT,
	person_id VARCHAR(255) DEFAULT ""
);

-- Table Schema generated for MySQL 8
-- for all EPrint repositories as _aggregate_contributor
CREATE TABLE IF NOT EXISTS _aggregate_contributor (
	aggregate_id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	repository VARCHAR(255),
	collection VARCHAR(255),
	eprintid INTEGER,
	person_id VARCHAR(255) DEFAULT ""
);

-- Table Schema generated for MySQL 8
-- for all EPrint repositories as _aggregate_advisor
CREATE TABLE IF NOT EXISTS _aggregate_advisor (
	aggregate_id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	repository VARCHAR(255),
	collection VARCHAR(255),
	eprintid INTEGER,
	person_id VARCHAR(255) DEFAULT ""
);

-- Table Schema generated for MySQL 8
-- for all EPrint repositories as _aggregate_committee
CREATE TABLE IF NOT EXISTS _aggregate_committee (
	aggregate_id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	repository VARCHAR(255),
	collection VARCHAR(255),
	eprintid INTEGER,
	person_id VARCHAR(255) DEFAULT ""
);

-- Table Schema generated for MySQL 8
-- for all EPrint repositories as _aggregate_option_major
CREATE TABLE IF NOT EXISTS _aggregate_option_major (
	aggregate_id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	repository VARCHAR(255),
	collection VARCHAR(255),
	eprintid INTEGER,
    local_option VARCHAR(255) DEFAULT ""
);

-- Table Schema generated for MySQL 8
-- for all EPrint repositories as _aggregate_option_minor
CREATE TABLE IF NOT EXISTS _aggregate_option_minor (
	aggregate_id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	repository VARCHAR(255),
	collection VARCHAR(255),
	eprintid INTEGER,
    local_option VARCHAR(255) DEFAULT ""
);

-- Table Schema generated for MySQL 8
-- for all EPrint repositories as _aggregate_group
CREATE TABLE IF NOT EXISTS _aggregate_group (
	aggregate_id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	repository VARCHAR(255),
	collection VARCHAR(255),
	eprintid INTEGER,
	local_group VARCHAR(255) DEFAULT ""
);
`
	table := `--
-- Table Schema generated for MySQL 8 by %s %s
-- for EPrint repository %s on %s
--
CREATE TABLE IF NOT EXISTS %s (
  id INTEGER NOT NULL PRIMARY KEY,
  src JSON,
  updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  action VARCHAR(255) DEFAULT "",
  created VARCHAR(255) DEFAULT "",
  lastmod VARCHAR(255) DEFAULT "",
  pubdate VARCHAR(255) DEFAULT "",
  is_public BOOLEAN DEFAULT FALSE,
  record_type VARCHAR(255) DEFAULT "",
  status VARCHAR(255) DEFAULT ""
);
`
	cfg, err := LoadConfig(cfgName)
	if err != nil {
		return "", err
	}
	dbName, err := getDBName(cfg.JSONStore)
	if err != nil {
		return "", fmt.Errorf("cannot parse jsonstore DSN in %s, %s", cfgName, err)
	}
	src := []string{}
	if dbName != "" {
		src = append(src, fmt.Sprintf(database, appName, Version, cfgName, now.Format(mysqlTimeFmt), dbName, dbName, dbName))
	}
	for repoName := range cfg.Repositories {
		src = append(src, fmt.Sprintf(table, appName, Version, repoName, now.Format(mysqlTimeFmt), repoName))
	}
	return strings.Join(src, "\n"), nil
}

// RunHarvester will use the config file names by cfgName and
// the start and end time strings if set to retrieve all eprint
// records created or modified during that time sequence.
func RunHarvester(cfgName string, start string, end string, verbose bool) error {
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
	return harvest(cfg, start, end, verbose)
}

// harvest implements an harvest instance.
func harvest(cfg *Config, start string, end string, verbose bool) error {
	// Open the repository collections, JSON store with aggregation tables
	if err := OpenConnections(cfg); err != nil {
		return err
	}
	if err := OpenJSONStore(cfg); err != nil {
		return err
	}
	defer cfg.Jdb.Close()
	repoNames := []string{}
	for repoName := range cfg.Connections {
		repoNames = append(repoNames, repoName)
	}
	sort.Strings(repoNames)
	for _, repoName := range repoNames {
		//FIXME: we could use a go routine to support concurrent harvests.
		log.Printf("harvesting %s started", repoName)
		if err := harvestRepository(cfg, repoName, start, end, verbose); err != nil {
			log.Printf("harvesting %s aborted", repoName)
			return err
		}
		log.Printf("harvesting %s completed", repoName)
	}
	return CloseConnections(cfg)
}

// getSortedUniqueIDs takes an array of ints, makes a copy and
// returns a new sorted unique list of ints
func getSortedUniqueIDs(ids []int) []int {
	copiedIDs := ids[:]
	sort.Ints(copiedIDs)
	uniqueIDs := []int{}
	last := 0
	for _, val := range ids {
		if (val > 0) && ((last == 0) || (val > last)) {
			uniqueIDs = append(uniqueIDs, val)
			last = val
		}
	}
	return uniqueIDs
}

// harvestRepository takes a configuration with open database connections.
// a repository name (i.e. eprint database name) along with a start and
// end timestamp. It harvests records created/modified from the repository
// in time range.
func harvestRepository(cfg *Config, repoName string, start string, end string, verbose bool) error {
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
		err := harvestEPrintRecord(cfg, repoName, id)
		if err != nil {
			log.Printf("Harvesting EPrint %d (%d/%d) failed, %s", id, i, len(ids), err)
		}
		if verbose && ((i % 1000) == 0) {
			log.Printf("Harvested EPrint %d (%d/%d)", id, i, len(ids))
		}
	}
	return nil
}

// harvestEPrintRecord takes a configuration with open database connections
// a repository name and EPrint record ID and populates a JSON datastore
// of harvested EPrints.
func harvestEPrintRecord(cfg *Config, repoName string, eprintID int) error {
	ds, ok := cfg.Repositories[repoName]
	if !ok {
		return fmt.Errorf("Data Source not found for %q looking up eprint %d", repoName, eprintID)
	}
	eprint, err := SQLReadEPrint(cfg, repoName, ds.BaseURL, eprintID)
	if err != nil {
		return err
	}
	action := "created"
	if eprint.Datestamp != eprint.LastModified {
		action = "updated"
	}
	if eprint.EPrintStatus == "deletion" {
		action = "deleted"
	}
	src, _ := json.MarshalIndent(eprint, "", "    ")
	err = SaveJSONDocument(cfg, repoName, eprintID, src, action, eprint.Datestamp, eprint.LastModified, eprint.PubDate(), eprint.EPrintStatus, eprint.IsPublic(), eprint.Type)
	if err != nil {
		return err
	}
	// Since we can save the JSON recordd, need to aggregate the contents of it.
	aggregateEPrintRecord(cfg, repoName, eprintID, eprint)
	return err
}

// aggregatePersons aggregates by the person related roles, e.g. creator, editor, contributor, advisor, committee memember
func aggregatePersons(cfg *Config, repoName string, collection string, tableName string, eprintID int, personIDs []string) {
	deleteStmt := fmt.Sprintf(`DELETE FROM %s WHERE repository = ? AND collection = ? AND eprintid = ?`, tableName)
	cfg.Jdb.Exec(deleteStmt)
	if len(personIDs) > 0 {
		insertStmt := fmt.Sprintf(`INSERT INTO %s (repository, collection, eprintid, person_id) VALUES (?, ?, ?, ?)`, tableName)
		for _, personID := range personIDs {
			if _, err := cfg.Jdb.Exec(insertStmt, repoName, collection, eprintID, personID); err != nil {
				log.Printf("WARNING: failed aggregreatePersons(cfg, %q, %q, %d, %q): %s", repoName, collection, eprintID, personID, err)
			}
		}
	}
}

func aggregateOptions(cfg *Config, repoName string, collection string, tableName string, eprintID int, options []string) {
	deleteStmt := fmt.Sprintf(`DELETE FROM %s WHERE repository = ? AND collection = ? AND eprintid = ?`, tableName)
	cfg.Jdb.Exec(deleteStmt)
	if len(options) > 0 {
		insertStmt := fmt.Sprintf(`INSERT INTO %s (repository, collection, eprintid, local_option) VALUES (?, ?, ?, ?)`, tableName)
		for _, option := range options {
			if _, err := cfg.Jdb.Exec(insertStmt, repoName, collection, eprintID, option); err != nil {
				log.Printf("WARNING: failed aggregateOptions(cfg, %q, %q, %d, %q): %s", repoName, collection, eprintID, option, err)
			}
		}
	}
}

// aggregateGroup aggregates by the by group
func aggregateGroup(cfg *Config, repoName string, collection string, tableName string, eprintID int, groups []string) {
	deleteStmt := fmt.Sprintf(`DELETE FROM %s WHERE repository = ? AND collection = ? AND eprintid = ?`, tableName)
	cfg.Jdb.Exec(deleteStmt)
	if len(groups) > 0 {
		insertStmt := fmt.Sprintf(`INSERT INTO %s (repository, collection, eprintid, local_group) VALUES (?, ?, ?, ?)`, tableName)
		for _, groupName := range groups {
			if _, err := cfg.Jdb.Exec(insertStmt, repoName, collection, eprintID, groupName); err != nil {
				log.Printf("WARNING: failed aggregateGroup(cfg, %q, %q, %d, %q): %s", repoName, collection, eprintID, groupName, err)
			}
		}
	}
}

// aggregateEPrintRecord takes a configuration, repository name, eprintID and struct
// then performances an analysis of the record aggregrating it's component parts.
func aggregateEPrintRecord(cfg *Config, repoName string, eprintID int, eprint *EPrint) {
	collection := eprint.Collection
	// Clear the creators aggregation for this eprint record
	if personIDs := eprint.Creators.GetIDs(); len(personIDs) > 0 {
		aggregatePersons(cfg, repoName, collection, "_aggregate_creator", eprintID, personIDs)
	}
	if personIDs := eprint.Editors.GetIDs(); len(personIDs) > 0 {
		aggregatePersons(cfg, repoName, collection, "_aggregate_editor", eprintID, personIDs)
	}
	if personIDs := eprint.Contributors.GetIDs(); len(personIDs) > 0 {
		aggregatePersons(cfg, repoName, collection, "_aggregate_contributor", eprintID, personIDs)
	}
	if personIDs := eprint.ThesisAdvisor.GetIDs(); len(personIDs) > 0 {
		aggregatePersons(cfg, repoName, collection, "_aggregate_advisor", eprintID, personIDs)
	}
	if personIDs := eprint.ThesisCommittee.GetIDs(); len(personIDs) > 0 {
		aggregatePersons(cfg, repoName, collection, "_aggregate_committee", eprintID, personIDs)
	}
	if options := eprint.OptionMajor.GetOptions(); len(options) > 0 {
		aggregateOptions(cfg, repoName, collection, "_aggregate_option_major", eprintID, options)
	}
	if options := eprint.OptionMinor.GetOptions(); len(options) > 0 {
		aggregateOptions(cfg, repoName, collection, "_aggregate_option_minor", eprintID, options)
	}
	if groups := eprint.LocalGroup.GetGroups(); len(groups) > 0 {
		aggregateGroup(cfg, repoName, collection, "_aggregate_group", eprintID, groups)
	}
}
