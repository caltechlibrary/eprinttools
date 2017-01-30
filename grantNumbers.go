package epgo

import (
	"fmt"
	"strings"

	// CaltechLibrary Packages
	"github.com/caltechlibrary/dataset"
)

// GetGrantNumbersByFunder returns a JSON list of unique Group names in index
func (api *EPrintsAPI) GetGrantNumbersByFunder(direction int) ([]string, error) {
	c, err := dataset.Open(api.Dataset)
	failCheck(err, fmt.Sprintf("GetGrantNumbers() %s, %s", api.Dataset, err))
	defer c.Close()

	sl, err := c.Select("funder")
	if err != nil {
		return nil, err
	}
	sl.Sort(direction)

	// Note: Aggregate the local group names
	grantNumbersByFunder := []string{}
	lastKey := ""
	curKey := ""
	for _, id := range sl.List() {
		kys := strings.Split(id, indexDelimiter)
		funderName := first(kys)
		grantNo := second(kys)
		//NOTE: Since some ageny funded pubs lack Grant No, skip those (They will be listed under funder)
		if len(grantNo) > 0 && len(funderName) > 0 {
			curKey = fmt.Sprintf("%s%s%s", funderName, indexDelimiter, grantNo)

			if curKey != lastKey {
				lastKey = curKey
				grantNumbersByFunder = append(grantNumbersByFunder, curKey)
			}
		}
	}
	return grantNumbersByFunder, nil
}

// GetGrantNumberPublications returns a list of EPrint records with funderName
func (api *EPrintsAPI) GetGrantNumberPublications(funderName string, grantNumber string, start, count, direction int) ([]*Record, error) {
	c, err := dataset.Open(api.Dataset)
	failCheck(err, fmt.Sprintf("GetGrantNumberPublications() %s, %s", api.Dataset, err))
	defer c.Close()

	// Note: Filter for funderName/Grant Number, passing matching eprintIDs to getRecordList()
	ids, err := api.GetIDsBySelectList("funder", direction, func(s string) bool {
		parts := strings.Split(s, indexDelimiter)
		if funderName == first(parts) && grantNumber == second(parts) {
			return true
		}
		return false
	})
	if err != nil {
		return nil, err
	}
	return getRecordList(c, ids, start, count, func(rec *Record) bool {
		if rec.IsPublished == "pub" {
			return true
		}
		return false
	})
}

// GetGrantNumberArticles returns a list of EPrint records with funderName
func (api *EPrintsAPI) GetGrantNumberArticles(funderName string, grantNumber string, start, count, direction int) ([]*Record, error) {
	c, err := dataset.Open(api.Dataset)
	failCheck(err, fmt.Sprintf("GetGrantNumberArticles() %s, %s", api.Dataset, err))
	defer c.Close()

	// Note: Filter for funderName/GrantNumber, passing matching eprintIDs to getRecordList()
	ids, err := api.GetIDsBySelectList("funder", direction, func(s string) bool {
		parts := strings.Split(s, indexDelimiter)
		if funderName == first(parts) && grantNumber == second(parts) {
			return true
		}
		return false
	})
	if err != nil {
		return nil, err
	}
	return getRecordList(c, ids, start, count, func(rec *Record) bool {
		if rec.IsPublished == "pub" && rec.Type == "article" {
			return true
		}
		return false
	})
}
