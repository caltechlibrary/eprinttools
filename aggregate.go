package eprinttools

import (
	"database/sql"
	"fmt"
)

/*
The "aggregator" operates by generating JSON lists and documents based on the
content previously harvested. During the harvesting process each eprint record
is analyized and content saved to a table that represents the supported aggregations.

Note all our EPrints repositories need to have been harvested as well as CaltechPEOPLE
before the "Aggretator" lists are retrieved and the JSON documents are generated.
*/

// Aggregation holds the connection to the SQL/JSON store where aggregations are
// maintained.
type Aggregation struct {
	DSN string  `json:"dsn,omitempty"`
	db  *sql.DB `json:"-"`
}

// Aggregate takes an eprint record and generates the related aggregations
// that the eprint record could fall into.
//
// ```
//
//	aggregation, err := eprinttools.OpenAggregation(cfg)
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
	// First prune any references to the EPrint record
	// SQL: DELETE FROM %q WHERE repository = ? AND eprintid = ?;

	// Analyze EPrint record generating the aggregation rows
	// for each of our List types.

	// 1. Generate rows per eprint id for each *_id of creators, editors, contributors, advisor, committee member
	// 2. Generate row per eprint id for each "local group"
	// 3. Generate row per eprint id for each "funder"

	// Each row needs to contain document type, creation date (dated added),  pub date, collection

	return fmt.Errorf("aggregation.Aggregate() not implemented")
}

func (aggretation *Aggregation) ListByType() ([]byte, error) {
	return nil, fmt.Errorf("aggregation.ByType() not implemented")
}

func (aggregation *Aggregation) ListByPubDate() ([]byte, error) {
	return nil, fmt.Errorf("aggregation.ByPubDate() not implemented")
}

func (aggregation *Aggregation) ListByDateAdded() ([]byte, error) {
	return nil, fmt.Errorf("aggregation.ByDateAdded() not implemeneted")
}

func (aggregation *Aggregation) ListByPerson() ([]byte, error) {
	/* this should return list of cannonical person ids, this list can then be used to create individual lists by roles */
	return nil, fmt.Errorf("aggregation.ByPerson() not implemented")
}

func (aggregation *Aggregation) ListrByPersonAndRole() ([]byte, error) {
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

func (aggregation *Aggregation) ListByPublication() ([]byte, error) {
	return nil, fmt.Errorf("aggregation.ByPublication() not implemented")
}

func (aggregation *Aggregation) ListByGroup() ([]byte, error) {
	return nil, fmt.Errorf("aggregation.ByGroup() not implemented")
}

func (aggregation *Aggregation) ListBySubject() ([]byte, error) {
	return nil, fmt.Errorf("aggregeation.BySubject() not implemented")
}

func (aggregation *Aggregation) ListByCategory() ([]byte, error) {
	return nil, fmt.Errorf("aggregation.ByCategory() not implemented")
}

func (aggregation *Aggregation) ListByConference() ([]byte, error) {
	return nil, fmt.Errorf("aggregation.ByConference not implemented")
}

func (aggregation *Aggregation) ListByCollection() ([]byte, error) {
	return nil, fmt.Errorf("aggregation.ByCollection() not implemented")
}

func (aggregation *Aggregation) ListByAuthor() ([]byte, error) {
	return nil, fmt.Errorf("aggregation.ByAuthor() not implemented")
}

func (aggregation *Aggregation) ListByEditor() ([]byte, error) {
	return nil, fmt.Errorf("aggregation.ByEditor() not implemented")
}

func (aggregation *Aggregation) ListByContributor() ([]byte, error) {
	return nil, fmt.Errorf("aggregation.ByContributor() not implemented")
}

func (aggregation *Aggregation) ListByAdvisor() ([]byte, error) {
	return nil, fmt.Errorf("aggregation.ByAdvisor() not implemented")
}

func (aggregation *Aggregation) ListByCommitteeMember() ([]byte, error) {
	return nil, fmt.Errorf("aggregation.ByCommitteeMember() not implemented")
}

func (aggregation *Aggregation) ListByDegreeThesisType() ([]byte, error) {
	return nil, fmt.Errorf("aggregation.ByThesisType() not implemented")
}

func (aggregation *Aggregation) ListByOption() ([]byte, error) {
	return nil, fmt.Errorf("aggregation.ByOption() not implemented")
}
