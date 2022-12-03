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

// GenerateRepositorySite renders an EPrint repository static site's
// directory structure, JSON documents and non-templated Markdown documents.
func GenerateRepositorySite(cfg *Config, repoName string, verbose bool) error {
	return fmt.Errorf("GenRepositorySite(cfg, %q, %T) not implemented", repoName, verbose)
}

// GenrateDataset renders a dataset collection holding the metadata
// from a previusouly harvested repository.
func GenerateDataset(cfg *Config, repoName string, verbose bool) error {
	return fmt.Errorf("GenerateDatasets(cfg, %q, %T) not implemented", repoName, verbose)
}

//
// Generate Aggregated views
//

// GeneratePeople generates the directory and JSON documents for
// aggregation of people related views across all our repositories.
func GeneratePeople(cfg *Config, verbose bool) error {
	return fmt.Errorf("GeneratePeople(cfg, %T) not implemented", verbose)
}

// GeneratePeople generates the directory and JSON documents for
// aggregation of group related views across all our EPrint repositories.
func GenerateGroups(cfg *Config, verbose bool) error {
	return fmt.Errorf("GenerateGroups(cfg, %T) not implemented", verbose)
}

// GenerateRecent generates the directory and JSON documents for all
// document types aggregated across all our EPrint repositories
func GenerateRecent(cfg *Config, verbose bool) error {
	return fmt.Errorf("GenerateGroups(cfg, %T) not implemented", verbose)
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
		return fmt.Errorf("Could not create a configuration object")
	}
	log.Printf("%s started %v", appName, t0)
	for repoName := range cfg.Repositories {
		if err := GenerateDataset(cfg, repoName, verbose); err != nil {
			return err
		}
		if err := GenerateRepositorySite(cfg, repoName, verbose); err != nil {
			return err
		}
	}
	if err := GenerateRecent(cfg, verbose); err != nil {
		return err
	}
	if err := GenerateGroups(cfg, verbose); err != nil {
		return err
	}
	if err := GeneratePeople(cfg, verbose); err != nil {
		return err
	}
	t1 := time.Now()
	log.Printf("%s finished %v", appName, t1)
	log.Printf("%s total run time %v", t1.Sub(t0))
	return nil
}
