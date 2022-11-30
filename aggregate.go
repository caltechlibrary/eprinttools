package eprinttools

import (
	"database/sql"
	"fmt"
)

/*
This file should define and orchestrate aggregating EPrint content
into "view" (e.g. by record type, by pub date, by author, by editor, etc).
It needs to work across all our repositories and be general enough for
an easy fit when we've migrated to a new repository system.
*/

// Aggregation holds the connection to the SQL/JSON store where aggregations are
// maintained.
type Aggregation struct {
	DSN string  `json:"dsn,omitempty"`
	db  *sql.DB `json:"-"`
}

type AggregateItem struct {
	Name        string `json:"name"`
	Link        string `json:"link,omitempty"`
	Description string `json:"description"`
	Repository  string `json:"repository,omitempty"`
	EPrintID    int    `json:"eprintid"`
	Created     string `json:"created,omitempty"`
	PubDate     string `json:"pubdate,omitempty"`
	IsPublic    string `json:"is_public,omitempty"`
}

// OpenAggregation takes a filename for JSON settings file and an attribute name holding
// the DSN for maintaing aggregations in a SQL/JSON store.
func OpenAggregation(cfgName string) (*Aggregation, error) {
    cfg, err := LoadConfig(cfgName)
    if err != nil {
        return nil, err
    }
	aggregation := new(Aggregation)
	aggregation.DSN = cfg.AggregationStore
	if cfg.Adb == nil {
    	dbName, err := getDBName(cfg.AggregationStore)
		if err != nil {
        	return nil, fmt.Errorf("Cannot parse aggregation DSN in %s, %s", cfgName, err)
		}
		db, err := sql.Open("mysql", dbName)
    	if err != nil {
        	return nil, fmt.Errorf("Cannot open %s, %s", dbName, err)
    	}
		cfg.Adb = db
	}
	aggregation.db = cfg.Adb
	return aggregation, nil
}

// Close closes an aggregation.
func (aggregation *Aggregation) Close() error {
	aggregation.DSN = ""
	if aggregation.db != nil {
		return aggregation.db.Close()
	}
	return nil
}

// Aggregate takes an eprint record and generates the related aggregations
// that the eprint record could fall into.
//
// ```
//
//	aggregation, err := eprinttools.OpenAggregation("settings.json")
//	if err err != nil {
//	      // ... handle error ...
//	}
//	defer aggregation.Close()
//	eprint := eprinttools.GetEPrint("http://eprint.example.edu", 12122)
//	if err = aggregation.Aggregate(eprint); err != nil {
//	      // ... handle error ...
//	}
//
// ```
func (aggregation *Aggregation) Aggregate(eprint *EPrint) error {
	// Analyze EPrint record generating the aggregation rows
	// appropirately.

	return fmt.Errorf("aggregation.Aggregate() not implemented")
}

func (aggretation *Aggregation) ByType() ([]*AggregateItem, error) {
	return nil, fmt.Errorf("aggregation.ByType() not implemented")
}

func (aggregation *Aggregation) ByPubDate() ([]*AggregateItem, error) {
	return nil, fmt.Errorf("aggregation.ByPubDate() not implemented")
}

func (aggregation *Aggregation) ByDateAdded() ([]*AggregateItem, error)  {
	return nil, fmt.Errorf("aggregation.ByDateAdded() not implemeneted")
}

func (aggregation *Aggregation) ByPerson() ([]*AggregateItem, error) {
	/* this should return list of cannonical person ids, this list can then be
	   used to create individual lists by roles */
	   return nil, fmt.Errorf("aggregation.ByPerson() not implemented")
}

func (aggregation *Aggregation) PersonAndRole() ([]*AggregateItem, error) {
	/* This aggregation should be broken down by person and role, it should include the following
	   sub lists
	   - author
	   - editor
	   - contributor
	   - advisor
	   - commitee memeber
	   - * would represent all roles (i.e. a combined list)
	*/
	return nil, fmt.Errorf("aggregation.ByPersonAndRole() not implemented")
}

func (aggregation *Aggregation) ByPublication() ([]*AggregateItem, error) {
	return nil, fmt.Errorf("aggregation.ByPublication() not implemented")
}

func (aggregation *Aggregation) ByGroup() ([]*AggregateItem, error) {
	return nil, fmt.Errorf("aggregation.ByGroup() not implemented")
}

func (aggregation *Aggregation) BySubject() ([]*AggregateItem, error) {
	return nil, fmt.Errorf("aggregeation.BySubject() not implemented")
}

func (aggregation *Aggregation) ByCategory() ([]*AggregateItem, error) {
	return nil, fmt.Errorf("aggregation.ByCategory() not implemented")
}

func (aggregation *Aggregation)ByConference() ([]*AggregateItem, error) {
	return nil, fmt.Errorf("aggregation.ByConference not implemented")
}

func (aggregation *Aggregation)ByCollection() ([]*AggregateItem, error) {
	return nil, fmt.Errorf("aggregation.ByCollection() not implemented")
}

func (aggregation *Aggregation)ByAuthor() ([]*AggregateItem, error) {
	return nil, fmt.Errorf("aggregation.ByAuthor() not implemented")
}

func (aggregation *Aggregation)ByEditor() ([]*AggregateItem, error) {
	return nil, fmt.Errorf("aggregation.ByEditor() not implemented")
}

func (aggregation *Aggregation)ByContributor() ([]*AggregateItem, error) {
	return nil, fmt.Errorf("aggregation.ByContributor() not implemented")
}

func (aggregation *Aggregation)ByAdvisor() ([]*AggregateItem, error) {
	return nil, fmt.Errorf("aggregation.ByAdvisor() not implemented")
}

func (aggregation *Aggregation)ByCommitteeMember() ([]*AggregateItem, error) {
	return nil, fmt.Errorf("aggregation.ByCommitteeMember() not implemented")
}

func (aggregation *Aggregation)ByDegreeThesisType() ([]*AggregateItem, error) {
	return nil, fmt.Errorf("aggregation.ByThesisType() not implemented")
}

func (aggregation *Aggregation)ByOption() ([]*AggregateItem, error) {
	return nil, fmt.Errorf("aggregation.ByOption() not implemented")
}
