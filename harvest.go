package eprinttools

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"
)

const (
	mysqlTimeFmt = "2006-01-02 15:04:05"
)

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
	for repoName := range cfg.Connections {
		//FIXME: we could use a go routine to support concurrent harvests.
		log.Printf("harvesting %s started", repoName)
		if err := harvestRepository(cfg, repoName, start, end); err != nil {
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
	src, err := json.MarshalIndent(eprint, "", "    ")
	return SaveJSONDocument(cfg, repoName, eprintID, src)
}
