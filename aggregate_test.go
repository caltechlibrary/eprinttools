package eprinttools

/*
This file should test our definitions and orchestratation for
aggregating EPrint content into "views" (e.g. by record type,
by pub date, by author, by editor, etc).
It needs to work across all our repositories and be general enough for
an easy fit when we've migrated to a new repository system.
*/

import (
	"testing"
)

func TestTestAggregate(t *testing.T) {
	t.Errorf("TestAggregate() not implemented")
}

func TestAggregationByType(t *testing.T) {
	t.Errorf("TestAggregationByPubDate not implemented")
}

func TestAggregationByPubDate(t *testing.T) {
	t.Errorf("TestAggregationByPubDate not implemented")
}

func TestAggregatorByDateAdded(t *testing.T) {
	t.Errorf("TestAggregateByDateAdded not implemeneted")
}

func TestAggregationByPersons(t *testing.T) {
	/* this should return list of cannonical person ids, this list can then be
	   used to create individual lists by roles */
}

func TestAggregateByPersonAndRole(t *testing.T) {
	/* This aggregation should be broken down by person and role, it should include the following
	   sub lists
	   - author
	   - editor
	   - contributor
	   - advisor
	   - commitee memeber
	   - * would represent all roles (i.e. a combined list)
	*/
	t.Errorf("TestAggregatedByPerson not implemented")
}

func TestAggregateByPublication(t *testing.T) {
	t.Errorf("TestAggregateByPublication not implemented")
}

func TestAggregateByGroup(t *testing.T) {
	t.Errorf("TestAggregatedByGroup not implemented")
}

func TestAggregateBySubject(t *testing.T) {
	t.Errorf("TestAggregateBySubject not implemented")
}

func TestAggregateByCategory(t *testing.T) {
	t.Errorf("TestAggregateByCategory not implemented")
}

func TestAggregateByConference(t *testing.T) {
	t.Errorf("TestAggregatedByConference not implemented")
}

func TestAggregateByCollection(t *testing.T) {
	t.Errorf("TestAggregateByCollection")
}

func TestAggregateByAuthor(t *testing.T) {
	t.Errorf("TestAggregatedByAuthor not implemented")
}

func TestAggregateByEditor(t *testing.T) {
	t.Errorf("TestAggregatedByEditor not implemented")
}

func TestAggregateByContributor(t *testing.T) {
	t.Errorf("TestAggregatedByContributor not implemented")
}

func TestAggregateByAdvisor(t *testing.T) {
	t.Errorf("TestAggregatedByAdvisor not implemented")
}

func TestAggregateByCommitteeMember(t *testing.T) {
	t.Errorf("TestAggregatedByCommitteeMember not implemented")
}

func TestAggregateByDegreeThesisType(t *testing.T) {
	t.Errorf("TestAggregatedByThesisType not implemented")
}

func TestAggregateByOption(t *testing.T) {
	t.Errorf("TestAggregatedByOption not implemented")
}
