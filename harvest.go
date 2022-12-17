package eprinttools

//
// This package implements an EPrints 3.x harvester storing
// Each harvested EPrints repository's JSON metadata in a
// MySQL 8 table using JSON column type.
//

import (
	"encoding/csv"
	"fmt"
	"io"
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

// HarvesterDBSchema returns SQL statements for creating
// the tables and database for the harvester based on the
// initialization file provides.
func HarvesterDBSchema(cfgName string) (string, error) {
	now := time.Now()
	appName := os.Args[0]
	database := `--
--
-- Database generated for MySQL 8 by %s %s
-- using %s on %s
--
DROP DATABASE IF EXISTS %s;
CREATE DATABASE IF NOT EXISTS %s;
USE %s;

-- Table Schema generated for MySQL 8
-- for all EPrint repositories as _aggregate_creator
CREATE TABLE IF NOT EXISTS _aggregate_creator (
	aggregate_id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	repository VARCHAR(256),
	collection VARCHAR(256),
	eprintid INT,
	person_id VARCHAR(256) DEFAULT "",
  	pubdate VARCHAR(256) DEFAULT "",
  	is_public BOOLEAN DEFAULT FALSE,
  	record_type VARCHAR(256) DEFAULT "",
    thesis_type VARCHAR(256) DEFAULT ""
);
CREATE INDEX _aggregrate_creator_person_id_i ON _aggregate_creator (person_id ASC);
CREATE INDEX _aggregrate_creator_pubdate_i ON _aggregate_creator (pubdate DESC);

-- Table Schema generated for MySQL 8
-- for all EPrint repositories as _aggregate_editor
CREATE TABLE IF NOT EXISTS _aggregate_editor (
	aggregate_id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	repository VARCHAR(256),
	collection VARCHAR(256),
	eprintid INT,
	person_id VARCHAR(256) DEFAULT "",
  	pubdate VARCHAR(256) DEFAULT "",
  	is_public BOOLEAN DEFAULT FALSE,
  	record_type VARCHAR(256) DEFAULT "",
    thesis_type VARCHAR(256) DEFAULT ""
);
CREATE INDEX _aggregrate_editor_person_id_i ON _aggregate_editor (person_id ASC);
CREATE INDEX _aggregrate_editor_pubdate_i ON _aggregate_editor (pubdate DESC);

-- Table Schema generated for MySQL 8
-- for all EPrint repositories as _aggregate_contributor
CREATE TABLE IF NOT EXISTS _aggregate_contributor (
	aggregate_id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	repository VARCHAR(256),
	collection VARCHAR(256),
	eprintid INTEGER,
	person_id VARCHAR(256) DEFAULT "",
  	pubdate VARCHAR(256) DEFAULT "",
  	is_public BOOLEAN DEFAULT FALSE,
  	record_type VARCHAR(256) DEFAULT "",
      thesis_type VARCHAR(256) DEFAULT ""
);
CREATE INDEX _aggregrate_contributor_person_id_i ON _aggregate_contributor (person_id ASC);
CREATE INDEX _aggregrate_contributor_pubdate_i ON _aggregate_contributor (pubdate DESC);

-- Table Schema generated for MySQL 8
-- for all EPrint repositories as _aggregate_advisor
CREATE TABLE IF NOT EXISTS _aggregate_advisor (
	aggregate_id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	repository VARCHAR(256),
	collection VARCHAR(256),
	eprintid INTEGER,
	person_id VARCHAR(256) DEFAULT "",
  	pubdate VARCHAR(256) DEFAULT "",
  	is_public BOOLEAN DEFAULT FALSE,
  	record_type VARCHAR(256) DEFAULT "",
    thesis_type VARCHAR(256) DEFAULT ""
);
CREATE INDEX _aggregrate_advisor_person_id_i ON _aggregate_advisor (person_id ASC);
CREATE INDEX _aggregrate_advisor_pubdate_i ON _aggregate_advisor (pubdate DESC);

-- Table Schema generated for MySQL 8
-- for all EPrint repositories as _aggregate_committee
CREATE TABLE IF NOT EXISTS _aggregate_committee (
	aggregate_id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	repository VARCHAR(256),
	collection VARCHAR(256),
	eprintid INTEGER,
	person_id VARCHAR(256) DEFAULT "",
  	pubdate VARCHAR(256) DEFAULT "",
  	is_public BOOLEAN DEFAULT FALSE,
  	record_type VARCHAR(256) DEFAULT "",
    thesis_type VARCHAR(256) DEFAULT ""
);
CREATE INDEX _aggregate_committee_person_id_i ON _aggregate_committee (person_id ASC);
CREATE INDEX _aggregate_committee_pubdate_i ON _aggregate_committee (pubdate DESC);

-- Table Schema generated for MySQL 8
-- for all EPrint repositories as _aggregate_option_major
CREATE TABLE IF NOT EXISTS _aggregate_option_major (
	aggregate_id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	repository VARCHAR(256),
	collection VARCHAR(256),
	eprintid INTEGER,
  	pubdate VARCHAR(256) DEFAULT "",
  	is_public BOOLEAN DEFAULT FALSE,
  	record_type VARCHAR(256) DEFAULT "",
    thesis_type VARCHAR(256) DEFAULT "",
    local_option VARCHAR(256) DEFAULT ""
);
CREATE INDEX _aggregate_option_major_pubdate_i ON _aggregate_option_major (pubdate DESC);

-- Table Schema generated for MySQL 8
-- for all EPrint repositories as _aggregate_option_minor
CREATE TABLE IF NOT EXISTS _aggregate_option_minor (
	aggregate_id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	repository VARCHAR(256),
	collection VARCHAR(256),
	eprintid INTEGER,
  	pubdate VARCHAR(256) DEFAULT "",
  	is_public BOOLEAN DEFAULT FALSE,
  	record_type VARCHAR(256) DEFAULT "",
    thesis_type VARCHAR(256) DEFAULT "",
    local_option VARCHAR(256) DEFAULT ""
);
CREATE INDEX _aggregate_option_minor_pubdate_i ON _aggregate_option_minor (pubdate DESC);

-- Table Schema generated for MySQL 8
-- for all EPrint repositories as _aggregate_groups
CREATE TABLE IF NOT EXISTS _aggregate_groups (
	aggregate_id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	repository VARCHAR(256),
	collection VARCHAR(256),
	eprintid INTEGER,
  	pubdate VARCHAR(256) DEFAULT "",
  	is_public BOOLEAN DEFAULT FALSE,
  	record_type VARCHAR(256) DEFAULT "",
    thesis_type VARCHAR(256) DEFAULT "",
	name VARCHAR(1024) DEFAULT "",
    alternative VARCHAR(1024) DEFAULT "",
    group_id VARCHAR(256) DEFAULT ""
);
CREATE INDEX _aggregate_groups_pubdate_i ON _aggregate_groups (pubdate DESC);

-- Table Schema generate for MySQL 8
-- for external People list (e.g. from people.csv)
CREATE TABLE IF NOT EXISTS _people (
    person_id VARCHAR(256) NOT NULL PRIMARY KEY,
    cl_people_id VARCHAR(256) DEFAULT "",
    family_name VARCHAR(256) DEFAULT "",
    given_name VARCHAR(256) DEFAULT "",
    sort_name VARCHAR(256) DEFAULT "",
    thesis_id VARCHAR(256) DEFAULT "",
    advisor_id VARCHAR(256) DEFAULT "",
    authors_id VARCHAR(256) DEFAULT "",
    editor_id VARCHAR(256) DEFAULT "",
    contributor_id VARCHAR(256) DEFAULT "",
    archivesspace_id VARCHAR(256) DEFAULT "",
    directory_id VARCHAR(256) DEFAULT "",
    viaf_id VARCHAR(256) DEFAULT "",
    lcnaf VARCHAR(256) DEFAULT "",
    isni VARCHAR(256) DEFAULT "",
    wikidata VARCHAR(256) DEFAULT "",
    snac VARCHAR(256) DEFAULT "",
    orcid VARCHAR(256) DEFAULT "",
    image VARCHAR(1024) DEFAULT "",
    educated_at TEXT,
    caltech BOOLEAN DEFAULT FALSE,
    jpl BOOLEAN DEFAULT FALSE,
    faculty BOOLEAN DEFAULT FALSE,
    alumn BOOLEAN DEFAULT FALSE,
    status VARCHAR(256) DEFAULT "",
    directory_person_type VARCHAR(1024) DEFAULT "",
    title VARCHAR(1024) DEFAULT "",
    bio TEXT,
    division VARCHAR(256) DEFAULT "",
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
CREATE INDEX _people_sort_name_i ON _people (sort_name ASC);

CREATE TABLE IF NOT EXISTS _groups (
    group_id VARCHAR(256) NOT NULL PRIMARY KEY,
    name VARCHAR(256) DEFAULT "",
    alternative  VARCHAR(256) DEFAULT "",
    email VARCHAR(256) DEFAULT "",
    date VARCHAR(256) DEFAULT "",
    description TEXT,
    start VARCHAR(256) DEFAULT "",
    approx_start VARCHAR(256) DEFAULT "",
    activity VARCHAR(256) DEFAULT "",
    end VARCHAR(256) DEFAULT "",
    approx_end VARCHAR(256) DEFAULT "",
    website VARCHAR(256) DEFAULT "",
    pi VARCHAR(256) DEFAULT "",
    parent VARCHAR(256) DEFAULT "",
    prefix VARCHAR(256) DEFAULT "",
    grid VARCHAR(256) DEFAULT "",
    isni VARCHAR(256) DEFAULT "",
    ringold VARCHAR(256) DEFAULT "",
    viaf VARCHAR(256) DEFAULT "",
    ror VARCHAR(256) DEFAULT "",
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
CREATE INDEX _groups_name_i ON _groups (name ASC);

`
	table := `--
-- Table Schema generated for MySQL 8 by %s %s
-- for EPrint repository %s on %s
--
CREATE TABLE IF NOT EXISTS %s (
  id INTEGER NOT NULL PRIMARY KEY,
  src JSON,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  action VARCHAR(256) DEFAULT "",
  created VARCHAR(256) DEFAULT "",
  lastmod VARCHAR(256) DEFAULT "",
  pubdate VARCHAR(256) DEFAULT "",
  is_public BOOLEAN DEFAULT FALSE,
  record_type VARCHAR(256) DEFAULT "",
  thesis_type VARCHAR(256) DEFAULT "",
  status VARCHAR(256) DEFAULT ""
);
CREATE INDEX %s_pubdate_i ON %s (pubdate DESC);
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
		src = append(src, fmt.Sprintf(table, appName, Version, repoName, now.Format(mysqlTimeFmt), repoName, repoName, repoName))
	}
	return strings.Join(src, "\n"), nil
}

// RunHarvester will use the config file names by cfgName and
// the start and end time strings if set to retrieve all eprint
// records created or modified during that time sequence.
func RunHarvester(cfgName string, start string, end string, verbose bool) error {
	if cfgName == "" {
		return fmt.Errorf("configuration filename missing")
	}
	now := time.Now()
	// Read in the configuration for this harvester instance.
	cfg, err := LoadConfig(cfgName)
	if err != nil {
		return err
	}
	if cfg == nil {
		return fmt.Errorf("could not create a configuration object")
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

// RunHarvestRepoID will use the config file names by cfgName and
// the repository id, the start and end time strings if set to retrieve all eprint
// records created or modified during that time sequence for that repository.
func RunHarvestRepoID(cfgName string, repoName, start string, end string, verbose bool) error {
	if cfgName == "" {
		return fmt.Errorf("configuration filename missing")
	}
	now := time.Now()
	// Read in the configuration for this harvester instance.
	cfg, err := LoadConfig(cfgName)
	if err != nil {
		return err
	}
	if cfg == nil {
		return fmt.Errorf("could not create a configuration object")
	}
	if start == "" {
		// Pick a start date/time before EPrints existed.
		start = "2000-01-01 00:00:00"
	}
	if end == "" {
		// Pick the current date time when harvest is starting
		end = now.Format(mysqlTimeFmt)
	}
	if err := OpenConnections(cfg); err != nil {
		return err
	}
	defer CloseConnections(cfg)
	if err := OpenJSONStore(cfg); err != nil {
		return err
	}
	defer cfg.Jdb.Close()
	log.Printf("harvesting %s started", repoName)
	if err := harvestRepository(cfg, repoName, start, end, verbose); err != nil {
		log.Printf("harvesting %s aborted", repoName)
		return err
	}
	log.Printf("harvesting %s completed", repoName)
	return nil
}

// harvest implements an harvest instance.
func harvest(cfg *Config, start string, end string, verbose bool) error {
	// Open the repository collections, JSON store with aggregation tables
	if err := OpenConnections(cfg); err != nil {
		return err
	}
	defer CloseConnections(cfg)
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
	tot := len(ids)
	modValue := calcModValue(tot)
	t0 := time.Now()
	if verbose {
		log.Printf("Processing %d unique keys", tot)
	}

	// FIXME: This code is called to deal with the person id name collisions
	// in our EPrints repositories due to lack of common person objects.
	// This calls code in person_id_remapping.go.
	/*
		if (repoName == "caltechauthors") {
			if err := loadPersonIDMapping(cfg, "caltechauthors"); err != nil {
				return err
			}
		}
	*/
	if repoName == "caltechthesis" {
		if err := loadPersonIDMapping(cfg, "caltechthesis"); err != nil {
			return err
		}
	}

	for i, id := range ids {
		err := harvestEPrintRecord(cfg, repoName, id)
		if err != nil {
			log.Printf("Harvesting EPrint %d (%s) failed, %s", id, progress(t0, i, tot), err)
		}
		if verbose && ((i % modValue) == 0) {
			log.Printf("Harvested EPrint %d (%s)", id, progress(t0, i, tot))
		}
	}
	log.Printf("Harvested %q in %v", repoName, time.Since(t0).Truncate(time.Second))
	return nil
}

// harvestEPrintRecord takes a configuration with open database connections
// a repository name and EPrint record ID and populates a JSON datastore
// of harvested EPrints.
func harvestEPrintRecord(cfg *Config, repoName string, eprintID int) error {
	ds, ok := cfg.Repositories[repoName]
	if !ok {
		return fmt.Errorf("data source not found for %q looking up eprint %d", repoName, eprintID)
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
	src, _ := jsonEncode(eprint)
	err = SaveJSONDocument(cfg, repoName, eprintID, src, action, eprint.Datestamp, eprint.LastModified, eprint.PubDate(), eprint.EPrintStatus, eprint.IsPublic(), eprint.Type, eprint.ThesisType)
	if err != nil {
		return err
	}
	// Since we can save the JSON recordd, need to aggregate the contents of it.
	aggregateEPrintRecord(cfg, repoName, eprintID, eprint)
	return err
}

// aggregatePersons aggregates by the person related roles, e.g. creator, editor, contributor, advisor, committee memember
func aggregatePersons(cfg *Config, repoName string, collection string, role string, eprintID int, recordType string, thesisType string, isPublic bool, pubDate string, personIDs []string) {
	// Normalize the person id from the crosswalks using in repositories.
	deleteStmt := fmt.Sprintf(`DELETE FROM _aggregate_%s WHERE repository = ? AND collection = ? AND eprintid = ?`, role)
	cfg.Jdb.Exec(deleteStmt)
	if len(personIDs) > 0 {
		insertStmt := fmt.Sprintf(`INSERT INTO _aggregate_%s (repository, collection, eprintid, person_id, record_type, thesis_type, is_public, pubdate) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, role)
		for _, personID := range personIDs {
			if _, err := cfg.Jdb.Exec(insertStmt, repoName, collection, eprintID, personID, recordType, thesisType, isPublic, pubDate); err != nil {
				log.Printf("WARNING: failed aggregreatePersons(cfg, %q, %q, %q, %d, %q, %q, %t, %q, %q): %s", repoName, collection, role, eprintID, recordType, thesisType, isPublic, pubDate, role, err)
			}
		}
	}
}

func aggregateOptions(cfg *Config, repoName string, collection string, tableName string, eprintID int, recordType string, thesisType string, isPublic bool, pubDate string, options []string) {
	deleteStmt := fmt.Sprintf(`DELETE FROM %s WHERE repository = ? AND collection = ? AND eprintid = ?`, tableName)
	cfg.Jdb.Exec(deleteStmt)
	if len(options) > 0 {
		insertStmt := fmt.Sprintf(`INSERT INTO %s (repository, collection, eprintid, record_type, thesis_type, is_public, pubdate, local_option) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, tableName)
		for _, option := range options {
			if _, err := cfg.Jdb.Exec(insertStmt, repoName, collection, eprintID, recordType, thesisType, isPublic, pubDate, option); err != nil {
				log.Printf("WARNING: failed aggregateOptions(cfg, %q, %q, %d, %q, %q, %t, %q, %q): %s", repoName, collection, eprintID, recordType, thesisType, isPublic, pubDate, option, err)
			}
		}
	}
}

// aggregateGroup aggregates a single group by group_id
func aggregateGroup(cfg *Config, repoName string, collection string, eprintID int, groupID string, recordType string, thesisType string, isPublic bool, pubDate string) error {
	insertStmt := `INSERT INTO _aggregate_groups (repository, collection, eprintid, record_type, thesis_type, is_public, pubdate, group_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	if _, err := cfg.Jdb.Exec(insertStmt, repoName, collection, eprintID, recordType, thesisType, isPublic, pubDate, groupID); err != nil {
		return err
	}
	return nil
}

// aggregateGroups aggregates a list of groups by group name
func aggregateGroups(cfg *Config, repoName string, collection string, eprintID int, recordType string, thesisType string, isPublic bool, pubDate string, groups []string) {
	deleteStmt := `DELETE FROM _aggregate_groups WHERE repository = ? AND collection = ? AND eprintid = ?`
	cfg.Jdb.Exec(deleteStmt)
	if len(groups) > 0 {
		for _, groupName := range groups {
			// We need to check to see if group name is defined in groups.csv (i.e. _groups)
			groupID, err := GetGroupIDByName(cfg, groupName)
			if err != nil {
				log.Printf("Skipping, %q (%s eprintid %d) not found in _groups table, query error %s", groupName, repoName, eprintID, err)
				continue
			}
			if groupID == "" {
				log.Printf("Skipping, %q (%s eprintid %d) not found in _groups table, returned empty groupID", groupName, repoName, eprintID)
				continue
			}
			if err := aggregateGroup(cfg, repoName, collection, eprintID, groupID, recordType, thesisType, isPublic, pubDate); err != nil {
				log.Printf("WARNING: failed aggregateGroup(cfg, %q, %q, %d, %q, %q, %q,%t, %q): %s", repoName, collection, eprintID, groupID, recordType, thesisType, isPublic, pubDate, err)
			}
		}
	}
}

// aggregateEPrintRecord takes a configuration, repository name, eprintID and struct
// then performances an analysis of the record aggregrating it's component parts.
func aggregateEPrintRecord(cfg *Config, repoName string, eprintID int, eprint *EPrint) {
	collection := eprint.Collection
	recordType := eprint.Type
	thesisType := eprint.ThesisType
	isPublic := eprint.IsPublic()
	pubDate := eprint.PubDate()
	// Clear the creators aggregation for this eprint record

	if personIDs := eprint.Creators.GetIDs(); len(personIDs) > 0 {
		// FIXME: This code is called to deal with the person id name collisions
		// in our EPrints repositories due to lack of common person objects.
		// This calls code in person_id_remapping.go.
		if repoName == "caltechthesis" {
			personIDs = normalizePersonIDs(repoName, personIDs, "thesis_id")
		}
		if repoName == "caltechauthors" {
			personIDs = normalizePersonIDs(repoName, personIDs, "authors_id")
		}
		aggregatePersons(cfg, repoName, collection, "creator", eprintID, recordType, thesisType, isPublic, pubDate, personIDs)
	}
	if personIDs := eprint.Editors.GetIDs(); len(personIDs) > 0 {
		aggregatePersons(cfg, repoName, collection, "editor", eprintID, recordType, thesisType, isPublic, pubDate, personIDs)
	}
	if personIDs := eprint.Contributors.GetIDs(); len(personIDs) > 0 {
		aggregatePersons(cfg, repoName, collection, "contributor", eprintID, recordType, thesisType, isPublic, pubDate, personIDs)
	}
	if personIDs := eprint.ThesisAdvisor.GetIDs(); len(personIDs) > 0 {
		// FIXME: This code is called to deal with the person id name collisions
		// in our EPrints repositories due to lack of common person objects.
		// This calls code in person_id_remapping.go.
		if repoName == "caltechthesis" {
			personIDs = normalizePersonIDs(repoName, personIDs, "advisor_id")
		}
		aggregatePersons(cfg, repoName, collection, "advisor", eprintID, recordType, thesisType, isPublic, pubDate, personIDs)
	}
	if personIDs := eprint.ThesisCommittee.GetIDs(); len(personIDs) > 0 {
		aggregatePersons(cfg, repoName, collection, "committee", eprintID, recordType, thesisType, isPublic, pubDate, personIDs)
	}
	if options := eprint.OptionMajor.GetOptions(); len(options) > 0 {
		aggregateOptions(cfg, repoName, collection, "_aggregate_option_major", eprintID, recordType, thesisType, isPublic, pubDate, options)
	}
	if options := eprint.OptionMinor.GetOptions(); len(options) > 0 {
		aggregateOptions(cfg, repoName, collection, "_aggregate_option_minor", eprintID, recordType, thesisType, isPublic, pubDate, options)
	}
	if groups := eprint.LocalGroup.GetGroups(); len(groups) > 0 {
		aggregateGroups(cfg, repoName, collection, eprintID, recordType, thesisType, isPublic, pubDate, groups)
	}
}

//
// EPrints does not maintain a "person" or "group" object.
// This is maintained externally an a CSV file for each.
// This file can match against various ids fields in EPrints
// record to create a curated person/group view.
//

// Person holds the data structure representing the general person
// information and the crosswalk IDs maintained in the people.csv file.
type Person struct {
	PersonID            string    `json:"id"`
	CLPeopleID          string    `json:"cl_people_id,omitempty"`
	FamilyName          string    `json:"family_name,omitempty"`
	GivenName           string    `json:"given_name,omitempty"`
	SortName            string    `json:"sort_name,omitempty"`
	ThesisID            string    `json:"thesis_id,omitempty"`
	AdvisorID           string    `json:"advisor_id,omitempty"`
	AuthorsID           string    `json:"authors_id,omitempty"`
	ArchivesSpaceID     string    `json:"archivesspace_id,omitempty"`
	DirectoryID         string    `json:"directory_id,omitempty"`
	VIAF                string    `json:"viaf_id,omitempty"`
	LCNAF               string    `json:"lcnaf,omitempty"`
	ISNI                string    `json:"isni,omitempty"`
	Wikidata            string    `json:"wikidata,omitempty"`
	SNAC                string    `json:"snac,omitempty"`
	ORCID               string    `json:"orcid,omitempty"`
	Image               string    `json:"image,omitempty"`
	EducatedAt          string    `json:"educated_at,omitempty"`
	Caltech             bool      `json:"caltech,omitempty"`
	JPL                 bool      `json:"jpl,omitempty"`
	Faculty             bool      `json:"faculty,omitempty"`
	Alumn               bool      `json:"alumn,omitempty"`
	Status              string    `json:"status,omitempty"`
	DirectoryPersonType string    `json:"directory_person_type,omitempty"`
	Title               string    `json:"title,omitempty"`
	Bio                 string    `json:"bio,omitempty"`
	Division            string    `json:"division,omitempty"`
	Updated             time.Time `json:"updated"`
}

// Group holds the data structure presenting the group information
// and the crossswalk IDs maintained in groups.csv
type Group struct {
	GroupID     string    `json:"key"`
	Name        string    `json:"name,omitempty"`
	Alternative string    `json:"alternative"`
	EMail       string    `json:"email,omitempty"`
	Date        string    `json:"date,omitempty"`
	Description string    `json:"description,omitempty"`
	Start       string    `json:"start,omitempty"`
	ApproxStart string    `json:"approx_start,omitempty"`
	Activity    string    `json:"activity,omitempty"`
	End         string    `json:"end,omitempty"`
	ApproxEnd   string    `json:"approx_end,omitempty"`
	Website     string    `json:"website,omitempty"`
	PI          string    `json:"pi,omitempty"`
	Parent      string    `json:"parent,omitempty"`
	Prefix      string    `json:"prefix,omitempty"`
	GRID        string    `json:"grid,omitempty"`
	ISNI        string    `json:"isni,omitempty"`
	RinGold     string    `json:"ringold,omitempty"`
	VIAF        string    `json:"viaf,omitempty"`
	ROR         string    `json:"ror,omitempty"`
	Updated     time.Time `json:"updated,omitempty"`
}

func aToBool(s string) bool {
	s = strings.ToLower(s)
	if strings.HasPrefix(s, "t") || strings.HasPrefix(s, "1") {
		return true
	}
	return false
}

func columnsToPerson(columns []string, row []string) (*Person, error) {
	if len(columns) != len(row) {
		return nil, fmt.Errorf("mismatched columns in row")
	}
	person := new(Person)
	for i, field := range columns {
		switch field {
		case "cl_people_id":
			person.CLPeopleID = strings.TrimSpace(row[i])
		case "family_name":
			person.FamilyName = strings.TrimSpace(row[i])
		case "given_name":
			person.GivenName = strings.TrimSpace(row[i])
		case "thesis_id":
			person.ThesisID = strings.TrimSpace(row[i])
		case "advisor_id":
			person.AdvisorID = strings.TrimSpace(row[i])
		case "authors_id":
			person.AuthorsID = strings.TrimSpace(row[i])
		case "archivesspace_id":
			person.ArchivesSpaceID = strings.TrimSpace(row[i])
		case "directory_id":
			person.DirectoryID = strings.TrimSpace(row[i])
		case "directory_person_type":
			person.DirectoryPersonType = strings.TrimSpace(row[i])
		case "viaf_id":
			person.VIAF = strings.TrimSpace(row[i])
		case "lcnaf":
			person.LCNAF = strings.TrimSpace(row[i])
		case "isni":
			person.ISNI = strings.TrimSpace(row[i])
		case "wikidata":
			person.Wikidata = strings.TrimSpace(row[i])
		case "snac":
			person.SNAC = strings.TrimSpace(row[i])
		case "orcid":
			person.ORCID = strings.TrimSpace(row[i])
		case "image":
			person.Image = strings.TrimSpace(row[i])
		case "educated_at":
			person.EducatedAt = strings.TrimSpace(row[i])
		case "caltech":
			person.Caltech = aToBool(row[i])
		case "jpl":
			person.JPL = aToBool(row[i])
		case "faculty":
			person.Faculty = aToBool(row[i])
		case "alumn":
			person.Alumn = aToBool(row[i])
		case "status":
			person.Status = strings.TrimSpace(row[i])
		case "directory_perosn_type":
			person.DirectoryPersonType = strings.TrimSpace(row[i])
		case "title":
			person.Title = strings.TrimSpace(row[i])
		case "bio":
			person.Bio = strings.TrimSpace(row[i])
		case "division":
			person.Division = strings.TrimSpace(row[i])
		case "updated":
			person.Updated = time.Now()
		case "authors_count":
			// Skip, we ignore this column
		case "editor_count":
			// Skip, we ignore this column
		case "thesis_count":
			// Skip, we ignore this column
		case "data_count":
			// Skip, we ignore this column
		case "advisor_count":
			// Skip, we ignore this column
		case "committee_count":
			// Skip, we ignore this columns
		default:
			return nil, fmt.Errorf("failed to map column (%d: %q)", i, field)
		}
	}
	// Handle personID, in version >= 2 PersonID should be lowercase
	person.PersonID = person.CLPeopleID
	// Handle creating SortName
	person.SortName = fmt.Sprintf("%s, %s", person.FamilyName, person.GivenName)
	return person, nil
}

func ReadPersonCSV(fName string, verbose bool) ([]*Person, error) {
	fp, err := os.Open(fName)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	r := csv.NewReader(fp)
	columns := []string{}
	people := []*Person{}
	i := 0
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if i == 0 {
			columns = append(columns, row...)
		} else {
			person, err := columnsToPerson(columns, row)
			if err != nil {
				log.Printf("could not convert row (%d) %+v, %s", i+1, row, err)
			} else if people != nil {
				people = append(people, person)
			}
		}
		i++
	}
	if verbose {
		log.Printf("%d people added from %s", i, fName)
	}
	return people, nil
}

func columnsToGroup(columns []string, row []string) (*Group, error) {
	if len(columns) != len(row) {
		return nil, fmt.Errorf("mismatched columns in row")
	}
	group := new(Group)
	for i, field := range columns {
		switch field {
		case "key":
			group.GroupID = strings.TrimSpace(row[i])
		case "name":
			group.Name = strings.TrimSpace(row[i])
		case "alternative":
			group.Alternative = strings.TrimSpace(row[i])
		case "email":
			group.EMail = strings.TrimSpace(row[i])
		case "date":
			group.Date = strings.TrimSpace(row[i])
		case "description":
			group.Description = strings.TrimSpace(row[i])
		case "start":
			group.Start = strings.TrimSpace(row[i])
		case "approx_start":
			group.ApproxStart = strings.TrimSpace(row[i])
		case "activity":
			group.Activity = strings.TrimSpace(row[i])
		case "end":
			group.End = strings.TrimSpace(row[i])
		case "approx_end":
			group.ApproxEnd = strings.TrimSpace(row[i])
		case "website":
			group.Website = strings.TrimSpace(row[i])
		case "pi":
			group.PI = strings.TrimSpace(row[i])
		case "parent":
			group.Parent = strings.TrimSpace(row[i])
		case "prefix":
			group.Prefix = strings.TrimSpace(row[i])
		case "grid":
			group.GRID = strings.TrimSpace(row[i])
		case "isni":
			group.ISNI = strings.TrimSpace(row[i])
		case "ringold":
			group.RinGold = strings.TrimSpace(row[i])
		case "viaf":
			group.VIAF = strings.TrimSpace(row[i])
		case "ror":
			group.ROR = strings.TrimSpace(row[i])
		case "updated":
			group.Updated = time.Now()
		default:
			return nil, fmt.Errorf("failed to map column (%d: %q)", i, field)
		}
	}
	return group, nil
}

func ReadGroupsCSV(fName string, verbose bool) ([]*Group, error) {
	fp, err := os.Open(fName)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	r := csv.NewReader(fp)
	columns := []string{}
	groups := []*Group{}
	i := 0
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if i == 0 {
			columns = append(columns, row...)
		} else {
			group, err := columnsToGroup(columns, row)
			if err != nil {
				log.Printf("could not convert row (%d) %+v", i+1, row)
			} else if group != nil {
				groups = append(groups, group)
			}
		}
		i++
	}
	if verbose {
		log.Printf("%d groups added from %s", i, fName)
	}
	return groups, nil
}

// RunHarvestGroups loads CSV files containing people and group
// crosswalk tables. These synthesize Person and Group objects not present
// in EPrints.
func RunHarvestGroups(cfgName string, verbose bool) error {
	if cfgName == "" {
		return fmt.Errorf("configuration filename missing")
	}
	// Read in the configuration for this harvester instance.
	cfg, err := LoadConfig(cfgName)
	if err != nil {
		return err
	}
	if cfg == nil {
		return fmt.Errorf("could not create a configuration object")
	}
	if err := OpenJSONStore(cfg); err != nil {
		return err
	}
	defer cfg.Jdb.Close()
	if cfg.GroupsCSV != "" {
		groups, err := ReadGroupsCSV(cfg.GroupsCSV, verbose)
		if err != nil {
			return err
		}
		tot := len(groups)
		modValue := calcModValue(tot)
		t0 := time.Now()
		for i, group := range groups {
			if err := SaveGroupJSON(cfg, group); err != nil {
				log.Printf("failed to save group record, %s", err)
			}
			if verbose && ((i % modValue) == 0) {
				log.Printf("loaded %d groups (%s)", i, progress(t0, i, tot))
			}
		}
		if verbose {
			log.Printf("loaded %d groups in %v", tot, time.Since(t0).Truncate(time.Second))
		}
	}
	return nil
}

// RunHarvestPeople loads CSV files containing people and group
// crosswalk tables. These synthesize Person and Group objects not present
// in EPrints.
func RunHarvestPeople(cfgName string, verbose bool) error {
	if cfgName == "" {
		return fmt.Errorf("configuration filename missing")
	}
	// Read in the configuration for this harvester instance.
	cfg, err := LoadConfig(cfgName)
	if err != nil {
		return err
	}
	if cfg == nil {
		return fmt.Errorf("could not create a configuration object")
	}
	if err := OpenJSONStore(cfg); err != nil {
		return err
	}
	defer cfg.Jdb.Close()
	if cfg.PeopleCSV != "" {
		people, err := ReadPersonCSV(cfg.PeopleCSV, verbose)
		if err != nil {
			return err
		}
		tot := len(people)
		modValue := calcModValue(tot)
		t0 := time.Now()
		for i, person := range people {
			if err := SavePersonJSON(cfg, person); err != nil {
				log.Printf("failed to save person record, %s", err)
			}
			if verbose && ((i % modValue) == 0) {
				log.Printf("loaded %d people (%s)", i, progress(t0, i, tot))
			}
		}
		if verbose {
			log.Printf("loaded %d people in %v", tot, time.Since(t0).Truncate(time.Second))
		}
	}
	return nil
}
