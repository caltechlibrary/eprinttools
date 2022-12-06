package eprinttools

import (
	"fmt"
	"log"
	"os"
	"path"
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
	return fmt.Errorf("GenerateGroupIDs() not implemented")
}

// GeneratePeopleIDs returns a JSON document contiainer an array
// of people ids. People keys are sorted alphabetically.
func GeneratePeopleIDs(cfg *Config, repoName string, verbose bool) error {
	return fmt.Errorf("GeneratePeopleIDs() not implemented")
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
