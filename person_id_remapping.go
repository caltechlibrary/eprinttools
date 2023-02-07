package eprinttools

import (
	"strings"
	"log"
)

//
// This file is present to deal with the cl_people_id mapping collision
// in CaltechTHESIS and CaltechAUTHORS. After they are normalized to a
// common ID this code can be removed here and in harvest.go
//
var (
	thesisCreatorMap map[string]string
	thesisAdvisorMap map[string]string
	authorsCreatorMap map[string]string
)

// loadPersonIDMapping reads the _person table and creates maps for the
// id values based on the repository name (i.e. caltechauthors, caltechthesis)
func loadPersonIDMapping(cfg *Config, repoName string) error {
	var (
		personID string
		thesisID string
		advisorID string
		authorsID string
	)
	if repoName == "caltechthesis" {
		// Reset the maps and populate them
		thesisCreatorMap = map[string]string{}
		thesisAdvisorMap = map[string]string{}
		stmt := `SELECT person_id, thesis_id, advisor_id FROM _people ORDER BY person_id`
		rows, err := cfg.Jdb.Query(stmt)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&personID, &thesisID, &advisorID); err != nil {
				return err
			}
			if (personID != "") && (thesisID != "") && (strings.Compare(personID, thesisID) != 0) {
				thesisCreatorMap[thesisID] = personID
			}
			if (personID != "") && (advisorID != "") && (strings.Compare(personID, advisorID) != 0) {
				thesisAdvisorMap[advisorID] = personID
			}
		}
		log.Printf("%d ids remaped for %q thesis_id", len(thesisCreatorMap), repoName)
		log.Printf("%d ids remaped for %q advisor_id", len(thesisAdvisorMap), repoName)
		return nil
	}
	if repoName == "caltechauthors" {
		authorsCreatorMap = map[string]string{}
		stmt := `SELECT person_id, authors_id FROM _people ORDER BY person_id`
		rows, err := cfg.Jdb.Query(stmt)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&personID, &thesisID, &advisorID); err != nil {
				return err
			}
			if (personID != "") && (authorsID != "") && (strings.Compare(personID, authorsID) != 0){
				authorsCreatorMap[authorsID] = personID
			}
		}
		log.Printf("%d ids remaped for %q authors_id", len(authorsCreatorMap), repoName)
		return nil
	}
	return nil
}

// remapIDs in personIDs returning noramlize values in []string.
func remapIDs(personIDs []string, idMap map[string]string, label string) []string {
	var (
		newID string
		ok    bool
		cnt int
	)
	newPersonIDs := []string{}
	for _, personID := range personIDs {
		if newID, ok = idMap[personID]; ok {
			/*
			if label == "thesis_id" { // DEBUG
				log.Printf("DEBUG mapping (%q) %q -> %q", label, personID, newID)
			} // DEBUG
			*/
			newPersonIDs = append(newPersonIDs, newID)
			cnt++
		} else {
			newPersonIDs = append(newPersonIDs, personID)
		}
	}
	return newPersonIDs
}

// noramlizePersonIDs is a temporary fix until we can normalize the data in 
// CaltechAUTHORS/CaltechTHESIS around cl_people_id value.
func normalizePersonIDs(repoName string, personIDs []string, role string) []string {
	if strings.Compare(repoName, "caltechthesis") == 0 {
		switch role {
		case "thesis_id":
			return remapIDs(personIDs, thesisCreatorMap, role)
		case "advisor_id":
			return remapIDs(personIDs, thesisAdvisorMap, role)
		}
	}
	if strings.Compare(repoName, "caltechauthors") == 0 {
		switch role {
		case "authors_id":
			return remapIDs(personIDs, authorsCreatorMap, role)
		}
	}
	return personIDs
}

