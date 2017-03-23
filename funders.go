package epgo

import (
	"fmt"
	"strings"

	// CaltechLibrary Packages
	"github.com/caltechlibrary/dataset"
)

// GetFunders returns a JSON list of unique Group names in index
func (api *EPrintsAPI) GetFunders() ([]string, error) {
	c, err := dataset.Open(api.Dataset)
	failCheck(err, fmt.Sprintf("GetFunders() %s, %s", api.Dataset, err))
	defer c.Close()

	sl, err := c.Select("funder")
	if err != nil {
		return nil, err
	}
	// Note: Aggregate the local group names
	funderNames := []string{}
	lastFunder := ""
	funderName := ""
	for _, id := range sl.List() {
		funderName = first(strings.Split(id, indexDelimiter))
		if funderName != lastFunder {
			funderNames = append(funderNames, funderName)
			lastFunder = funderName
		}
	}
	return funderNames, nil
}

// GetFunderPublications returns a list of EPrint records with funderName
func (api *EPrintsAPI) GetFunderPublications(funderName string, start, count int) ([]*Record, error) {
	c, err := dataset.Open(api.Dataset)
	failCheck(err, fmt.Sprintf("GetFunderPublications() %s, %s", api.Dataset, err))
	defer c.Close()

	sl, err := c.Select("funder")
	if err != nil {
		return nil, err
	}
	if count == -1 {
		count = len(sl.Keys) + 1
	}

	// Note: Filter for funderName, passing matching eprintIDs to getRecordList()
	ids := []string{}
	for _, id := range sl.List() {
		parts := strings.Split(id, indexDelimiter)
		grp := first(parts)
		if grp == funderName {
			eprintID := last(parts)
			ids = append(ids, eprintID)
		}
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("zero ids in select list")
	}
	return getRecordList(c, ids, start, count, func(rec *Record) bool {
		if rec.IsPublished == "pub" {
			return true
		}
		return false
	})
}

// GetFunderArticles returns a list of EPrint records with funderName
func (api *EPrintsAPI) GetFunderArticles(funderName string, start, count int) ([]*Record, error) {
	c, err := dataset.Open(api.Dataset)
	failCheck(err, fmt.Sprintf("GetFunderArticles() %s, %s", api.Dataset, err))
	defer c.Close()

	sl, err := c.Select("funder")
	if err != nil {
		return nil, err
	}
	if count == -1 {
		count = len(sl.Keys) + 1
	}

	// Note: Filter for funderName, passing matching eprintIDs to getRecordList()
	ids := []string{}
	for _, id := range sl.List() {
		parts := strings.Split(id, indexDelimiter)
		grp := first(parts)
		if grp == funderName {
			eprintID := last(parts)
			ids = append(ids, eprintID)
		}
	}
	return getRecordList(c, ids, start, count, func(rec *Record) bool {
		if rec.IsPublished == "pub" && rec.Type == "article" {
			return true
		}
		return false
	})
}
