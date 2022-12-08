package eprinttools

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

//
//
// The "genfeeds" creates the directory structures, JSON and
// non-templated Markdown documents in the htdocs directory.
//

// GenerateGroupIDs returns a JSON document containing an array of
// group keys. Group keys are sorted alphabetically. Group keys are
// formed from the group field slugified.
func GenerateGroupIDs(cfg *Config, repoName string, verbose bool) error {
	if cfg.Jdb == nil {
		if err := OpenJSONStore(cfg); err != nil {
			return err
		}
	}
	groupDir := path.Join(cfg.Htdocs, "groups")
	// NOTE: Is htdocs relative to project? If so handle that case
	if ! (strings.HasPrefix(cfg.Htdocs, "/") || strings.HasPrefix(groupDir, cfg.ProjectDir)){
		groupDir = path.Join(cfg.ProjectDir, groupDir)
	}
	if _, err := os.Stat(groupDir); os.IsNotExist(err) {
		if err := os.MkdirAll(groupDir, 0775); err != nil {
			return err
		}
	}

	// generate htdocs/groups/index.json
	groupIDs, err := GetGroupIDs(cfg)
	if err != nil {
		return err
	}
	src, err := json.MarshalIndent(groupIDs, "", "    ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(groupDir, "index.json"), src, 0664); err != nil {
		return err
	}

	// For each group in _groups, find the records that should be included
	// with sublists of docuument type, options (major/minor), degree
	// thesis/disertations.
	return nil
}

// GeneratePeopleIDs returns a JSON document contiainer an array
// of people ids. People keys are sorted alphabetically.
func GeneratePeopleIDs(cfg *Config, repoName string, verbose bool) error {
	if cfg.Jdb == nil {
		if err := OpenJSONStore(cfg); err != nil {
			return err
		}
	}
	peopleDir := path.Join(cfg.Htdocs, "people")
	// NOTE: Is htdocs relative to project? If so handle that case
	if ! (strings.HasPrefix(cfg.Htdocs, "/") || strings.HasPrefix(peopleDir, cfg.ProjectDir)){
		peopleDir = path.Join(cfg.ProjectDir, peopleDir)
	}
	if _, err := os.Stat(peopleDir); os.IsNotExist(err) {
		if err := os.MkdirAll(peopleDir, 0775); err != nil {
			return err
		}
	}

	// generate htdocs/people/index.json
	personIDs, err := GetPersonIDs(cfg)
	if err != nil {
		return err
	}
	src, err := json.MarshalIndent(personIDs, "", "    ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(peopleDir, "index.json"), src, 0664); err != nil {
		return err
	}

	// For each person in _people, find the records that should be included
	// e.g. creator, editor, contributor, advisor, committee member.
	return nil
}

// RunGenfeeds will use the config file names by cfgName and
// render all the directorys, JSON documents and non-templated
// markdown content needed for a feeds v1.1 website in the htdocs
// directory indicated in the configuration file.
func RunGenfeeds(cfgName string, verbose bool) error {
	t0 := time.Now()
	appName := path.Base(os.Args[0])
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
	log.Printf("%s started %v", appName, t0)
	for repoName := range cfg.Repositories {
		if err := GenerateGroupIDs(cfg, repoName, verbose); err != nil {
			return err
		}
		if err := GeneratePeopleIDs(cfg, repoName, verbose); err != nil {
			return err
		}
	}
	if verbose {
		log.Printf("Genfeeds run time %v", time.Since(t0).Truncate(time.Second))
	}
	return nil
}
