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

func generateGroupDir(cfg *Config, groupID string, group *Group, m map[string]interface{}) error {
	groupDir := path.Join(cfg.Htdocs, "groups", groupID)
	// NOTE: Is htdocs relative to project? If so handle that case
	if !(strings.HasPrefix(cfg.Htdocs, "/") || strings.HasPrefix(groupDir, cfg.ProjectDir)) {
		groupDir = path.Join(cfg.ProjectDir, groupDir)
	}
	if _, err := os.Stat(groupDir); os.IsNotExist(err) {
		if err := os.MkdirAll(groupDir, 0775); err != nil {
			return err
		}
	}
	// FIXME: group.json should be renamed index.json, this file contains
	// all the metadata to write index.md, index.html and icnlude.include
	fName := path.Join(groupDir, "group.json")
	if err := jsonEncodeToFile(fName, m, 0664); err != nil {
		return err
	}
	for repoName, v := range m {
		// FIXME: We are using a switch to check the type because
		// the v1.0 feeds doesn't have a defined "aggregation" attribute.
		switch v.(type) {
		case map[string][]int:
			aMap := v.(map[string][]int)
			for recordType, ids := range aMap {
				docs := []*EPrint{}
				for _, id := range ids {
					eprint := new(EPrint)
					if err := GetDocumentAsEPrint(cfg, repoName, id, eprint); err != nil {
						return err
					}
					docs = append(docs, eprint)
				}
				fName = path.Join(groupDir, fmt.Sprintf("%s-%s.json", repoName, recordType))
				if err := jsonEncodeToFile(fName, docs, 0664); err != nil {
					return err
				}

			}
		}
	}
	return nil
}

func generateGroupListAndDir(cfg *Config, groupIDs []string, verbose bool) ([]map[string]interface{}, error) {
	groupList := []map[string]interface{}{}
	tot := len(groupIDs)
	t0 := time.Now()
	for i, groupID := range groupIDs {
		// FIXME: In feeds v1.0 aggregrated repositories info is mixed
		// into the attributes of the group's metadata, this prevents
		// us from using a struct representing Group or People
		// aggregations. In v1.1 or v1.2 this should change. Aggregations
		// should fall under an aggregated attribute which is a
		// map[string][]int type (or map[string][]string type when we
		// move to Invenio-RDM.
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
		// For each list type and combined retrieve
		// a sublist of eprint id from the appropriate
		// aggregated table
		aggregations, err := GetGroupAggregations(cfg, groupID)
		if err != nil {
			return nil, err
		}
		// Copy the aggregations into the map we'll use to represent a list of group info
		// for each repository in the system.
		hasAggregation := false
		for k, v := range aggregations {
			if combined, ok := v["combined"]; ok && len(combined) > 0 {
				m[k] = v
				hasAggregation = true
			}
		}
		if hasAggregation {
			if err := generateGroupDir(cfg, groupID, group, m); err != nil {
				log.Printf("failed to generate directory for %s, %s", groupID, err)
			}
			groupList = append(groupList, m)
			if verbose {
				log.Printf("processed %s (%s)", groupID, progress(t0, i, tot))
			}
		} else {
			log.Printf("%s did not have any aggregations", groupID)
		}
	}
	if verbose {
		log.Printf("%d entries in group list %v", tot, time.Since(t0).Truncate(time.Second))
	}
	return groupList, nil
}

// GenerateGroupFeed returns a JSON document containing an array of
// group keys. Group keys are sorted alphabetically. Group keys are
// formed from the group field slugified.
func GenerateGroupFeed(cfg *Config, verbose bool) error {
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
	fName = path.Join(groupDir, "group_list.json")
	groupList, err := generateGroupListAndDir(cfg, groupIDs, verbose)
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

// GeneratePeopleFeed returns a JSON document contiainer an array
// of people ids. People keys are sorted alphabetically.
func GeneratePeopleFeed(cfg *Config, verbose bool) error {
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
	peopleList := []*Person{}
	tot := len(personIDs)
	modValue := calcModValue(tot)
	t0 := time.Now()
	for i, personID := range personIDs {
		person, err := GetPerson(cfg, personID)
		if err != nil {
			return fmt.Errorf("failed to find %q in %q, %s", personID, cfg.JSONStore, err)
		}
		if person != nil {
			// For each person in _people, find the records that should be included
			// e.g. creator, editor, contributor, advisor, committee member.
			includePerson := false
			dName := path.Join(peopleDir, personID)
			if _, err := os.Stat(dName); os.IsNotExist(err) {
				if err := os.MkdirAll(dName, 0775); err != nil {
					return err
				}
			}
			// NOTE: each role is stored in a separate table for performance reasons. The table name is `_aggregate_<ROLE>`
			for _, role := range []string{ "creator", "contributor", "editor", "advisor", "committee" } {
				aMap, err := GetPersonByRoleAggregations(cfg, person, role)
				if err != nil {
					log.Printf("skipping %q for %q, %s", personID, role, err)
					continue
				} 
				if len(aMap) > 0 {
					includePerson = true
					fName = path.Join(dName, fmt.Sprintf("%s.json", role))
					if err := jsonEncodeToFile(fName, aMap, 0664); err != nil {
						return err
					}
					for repoName, rMap := range aMap {
						for recordType, ids := range rMap {
							fName = path.Join(dName, fmt.Sprintf("%s-%s-%s", repoName, recordType, role))
							records := []*EPrint{}
							for _, eprintid := range ids {
								eprint := new(EPrint)
								if err := GetDocumentAsEPrint(cfg, repoName, eprintid, eprint); err != nil {
									return err
								}
								records = append(records, eprint)

							}
							if err := jsonEncodeToFile(fName, records, 0664); err != nil {
								return err
							}
						}
					}
				}
			}
			if includePerson {
				peopleList = append(peopleList, person)
			} else {
				log.Printf("skipped %q, no aggregations found for roles, probably and id mismatch", personID)
			}
		}
		if verbose && ((i % modValue) == 0) {
			log.Printf("processed %s in _people, (%s)", personID, progress(t0, i, tot))
		}
	}
	if verbose {
		log.Printf("%d people updated in _people (%s)", tot, time.Since(t0).Truncate(time.Second))
		log.Printf("Writing %d people info %s", len(peopleList), fName)
	}
	fName = path.Join(peopleDir, "people_list.json")
	if err := jsonEncodeToFile(path.Join(peopleDir, "people_list.json"), peopleList, 0664); err != nil {
		return err
	}
	return nil
}

// RunGenPeople will use the config file names by cfgName and
// render all the directorys, JSON documents and non-templated
// markdown content needed for a feeds v1.1 website in the htdocs
// directory indicated in the configuration file.
func RunGenPeople(cfgName string, verbose bool) error {
	if cfgName == "" {
		return fmt.Errorf("Configuration filename missing")
	}
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
	if err := GeneratePeopleFeed(cfg, verbose); err != nil {
		return err
	}
	if verbose {
		log.Printf("%s run time %v", appName, time.Since(t0).Truncate(time.Second))
	}
	return nil
}

// RunGenGroups will use the config file names by cfgName and
// render all the directorys, JSON documents and non-templated
// markdown content needed for a feeds v1.1 website in the htdocs
// directory indicated in the configuration file.
func RunGenGroups(cfgName string, verbose bool) error {
	if cfgName == "" {
		return fmt.Errorf("Configuration filename missing")
	}
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
	if err := GenerateGroupFeed(cfg, verbose); err != nil {
		return err
	}
	if verbose {
		log.Printf("%s run time %v", appName, time.Since(t0).Truncate(time.Second))
	}
	return nil
}
// RunGenfeeds will use the config file names by cfgName and
// render all the directorys, JSON documents and non-templated
// markdown content needed for a feeds v1.1 website in the htdocs
// directory indicated in the configuration file.
func RunGenfeeds(cfgName string, verbose bool) error {
	if cfgName == "" {
		return fmt.Errorf("Configuration filename missing")
	}
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
	if err := GeneratePeopleFeed(cfg, verbose); err != nil {
		return err
	}
	if err := GenerateGroupFeed(cfg, verbose); err != nil {
		return err
	}
	if verbose {
		log.Printf("%s run time %v", appName, time.Since(t0).Truncate(time.Second))
	}
	return nil
}
