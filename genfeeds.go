package eprinttools

import (
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

func generateGroupListJSON(cfg *Config, groupIDs []string, verbose bool) ([]map[string]interface{}, error) {
	groupList := []map[string]interface{}{}
	tot := len(groupIDs)
	t0 := time.Now()
	for i, groupID := range groupIDs {
		group, err := GetGroup(cfg, groupID)
		if err != nil {
			return nil, err
		}
		src, err := jsonEncode(group)
		if err != nil {
			return nil, err
		}
		m := map[string]interface{}{}
		err = jsonDecode(src, &m)
		if err != nil {
			return nil, err
		}
		for _, repo := range cfg.Repositories {
			repoName := repo.DefaultCollection
			// For each list type and combined retrieve
			// a sublist of eprint id from the appropriate
			// aggregated table
			aggregations, err := GetGroupAggregations(cfg, repoName, groupID)
			if err != nil {
				return nil, err
			}
			if aggregations != nil {
				m[repoName] = aggregations
			}
			if verbose {
				log.Printf("Processed %s", repoName)
			}
		}
		if len(m) > 0 {
			groupList = append(groupList, m)
			if verbose {
				log.Printf("added %s (%s)", groupID, progress(t0, i, tot))
			}
		}
	}
	if verbose {
		log.Printf("%d entries in group list %v", tot, time.Since(t0).Truncate(time.Second))
	}
	return groupList, nil
}

// GenerateGroupIDs returns a JSON document containing an array of
// group keys. Group keys are sorted alphabetically. Group keys are
// formed from the group field slugified.
func GenerateGroupIDs(cfg *Config, verbose bool) error {
	groupDir := path.Join(cfg.Htdocs, "groups")
	// NOTE: Is htdocs relative to project? If so handle that case
	if !(strings.HasPrefix(cfg.Htdocs, "/") || strings.HasPrefix(groupDir, cfg.ProjectDir)) {
		groupDir = path.Join(cfg.ProjectDir, groupDir)
	}
	if _, err := os.Stat(groupDir); os.IsNotExist(err) {
		if err := os.MkdirAll(groupDir, 0775); err != nil {
			return err
		}
	}
	fName := path.Join(groupDir, "index.json")

	// generate htdocs/groups/index.json
	groupIDs, err := GetGroupIDs(cfg)
	if err != nil {
		return err
	}
	if verbose {
		log.Printf("Writing %d group ids to %s", len(groupIDs), fName)
	}
	if err := jsonEncodeToFile(fName, groupIDs, 0664); err != nil {
		return err
	}

	// For each group in _groups, find the records that should be included
	groupList, err := generateGroupListJSON(cfg, groupIDs, verbose)
	if err != nil {
		return err
	}
	if verbose {
		log.Printf("Writing %d group list entries to %s", len(groupIDs), fName)
	}
	if err := jsonEncodeToFile(fName, groupList, 0664); err != nil {
		return err
	}

	// with sublists of docuument type, options (major/minor), degree
	// thesis/disertations.
	return nil
}

// GeneratePeopleIDs returns a JSON document contiainer an array
// of people ids. People keys are sorted alphabetically.
func GeneratePeopleIDs(cfg *Config, verbose bool) error {
	peopleDir := path.Join(cfg.Htdocs, "people")
	// NOTE: Is htdocs relative to project? If so handle that case
	if !(strings.HasPrefix(cfg.Htdocs, "/") || strings.HasPrefix(peopleDir, cfg.ProjectDir)) {
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
	fName := path.Join(peopleDir, "index.json")
	if verbose {
		log.Printf("Writing %d person ids to %s", len(personIDs), fName)
	}
	if err := jsonEncodeToFile(fName, personIDs, 0664); err != nil {
		return err
	}
	fName = path.Join(peopleDir, "people_list.json")
	peopleList := []*Person{}
	for _, personID := range personIDs {
		person, err := GetPerson(cfg, personID)
		if err != nil {
			return fmt.Errorf("failed to find %q in %q, %s", personID, cfg.JSONStore, err)
		} else if person != nil {
			peopleList = append(peopleList, person)
		}
	}
	if verbose {
		log.Printf("Writing %d people info %s", len(peopleList), fName)
	}
	if err := jsonEncodeToFile(path.Join(peopleDir, "people_list.json"), peopleList, 0664); err != nil {
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
	log.Printf("%s started %v", appName, t0.Format("2006-01-02 15:04:05"))
	if err := GenerateGroupIDs(cfg, verbose); err != nil {
		return err
	}
	if err := GeneratePeopleIDs(cfg, verbose); err != nil {
		return err
	}
	if verbose {
		log.Printf("%s run time %v", appName, time.Since(t0).Truncate(time.Second))
	}
	return nil
}
