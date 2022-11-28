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
CREATE DATABASE IF NOT EXISTS %s;
USE %s;

`
	table := `-- 
-- Table Schema generated for MySQL 8 by %s %s
-- for EPrint repository %s on %s
--
CREATE TABLE IF NOT EXISTS %s (
  id INTEGER NOT NULL PRIMARY KEY,
  src JSON,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  action VARCHAR(255) DEFAULT "",
  lastmod VARCHAR(255) DEFAULT "",
  status VARCHAR(255) DEFAULT ""
);
`
	cfg, err := LoadConfig(cfgName)
	if err != nil {
		return "", err
	}
	dbName, err := getDBName(cfg.JSONStore)
	if err != nil {
		return "", fmt.Errorf("Cannot parse jsonstore DSN in %s, %s", cfgName, err)
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

// RunHarvester will use the config file names by cfgName and
// the start and end time strings if set to retrieve all eprint
// records created or modified during that time sequence.
func RunHarvester(cfgName string, start string, end string) error {
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
	return harvest(cfg, start, end)
}

// harvest implements an harvest instance.
func harvest(cfg *Config, start string, end string) error {
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
		log.Printf("harvesting %s started", repoName)
		if err := harvestRepository(cfg, repoName, start, end); err != nil {
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
func harvestRepository(cfg *Config, repoName string, start string, end string) error {
	createdIDs, err := GetEPrintIDsInTimestampRange(cfg, repoName, "datestamp", start, end)
	if err != nil {
		return err
	}
	modifiedIDs, err := GetEPrintIDsInTimestampRange(cfg, repoName, "lastmod", start, end)
	if err != nil {
		return err
	}
	ids := append(createdIDs, modifiedIDs...)
	ids = getSortedUniqueIDs(ids)
	for i, id := range ids {
		// FIXME: do we want to show a progress bar or just errors?
		err := harvestEPrintRecord(cfg, repoName, id)
		if err != nil {
			log.Printf("Harvesting EPrint %d (%d/%d) failed, %s", id, i, len(ids), err)
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
	src, err := json.MarshalIndent(eprint, "", "    ")
	return SaveJSONDocument(cfg, repoName, eprintID, src, action, eprint.LastModified, eprint.EPrintStatus)
}
