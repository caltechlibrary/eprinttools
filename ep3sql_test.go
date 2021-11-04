package eprinttools

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"
)

var (
	err error
)

func assertOpenConnection(t *testing.T, config *Config, repoID string) {
	ds, ok := config.Repositories[repoID]
	if !ok {
		t.Skipf("can't fund %q", repoID)
		t.SkipNow()
	}
	if config.Connections == nil {
		if err := OpenConnections(config); err != nil {
			t.Skipf("config has uninitialized connections")
			t.SkipNow()
		}
	} else {
		_, ok = config.Connections[repoID]
		if ok {
			db := config.Connections[repoID]
			if err := db.Ping(); err != nil {
				db, err = sql.Open("mysql", ds.DSN)
				if err != nil {
					t.Skipf(`could not re-open %q for %s, %s`, ds.DSN, repoID, err)
					t.SkipNow()
				}
				config.Connections[repoID] = db
			}
		} else {
			db, err := sql.Open("mysql", ds.DSN)
			if err != nil {
				t.Skipf(`could not open %q for %s, %s`, ds.DSN, repoID, err)
				t.SkipNow()
			}
			config.Connections[repoID] = db
		}
	}
}

func assertCloseConnection(t *testing.T, config *Config, repoID string) {
	if err := CloseConnections(config); err != nil {
		t.Errorf("Close error, %s", err)
		t.FailNow()
	}
}

func assertNameSame(t *testing.T, expected *Name, name *Name) {
	if expected.Family != name.Family {
		t.Errorf(`expected name family %s, got %s`, expected.Family, name.Family)
	}
	if expected.Given != name.Given {
		t.Errorf(`expected name given %s, got %s`, expected.Family, name.Family)
	}
	if expected.ID != name.ID {
		t.Errorf(`expected name id %s, got %s`, expected.ID, name.ID)
	}
	if expected.ORCID != name.ORCID {
		t.Errorf(`expected name orcid %s, got %s`, expected.ORCID, name.ORCID)
	}
	if expected.Honourific != name.Honourific {
		t.Errorf(`expected name honourific %s, got %s`, expected.Honourific, name.Honourific)
	}
	if expected.Lineage != name.Lineage {
		t.Errorf(`expected name lineage %s, got %s`, expected.Lineage, name.Lineage)
	}
	if expected.Value != name.Value {
		t.Errorf(`expected name value %s, got %s`, expected.Value, name.Value)
	}
}

func assertIntSame(t *testing.T, label string, expected int, got int) {
	if expected != got {
		t.Errorf(`expected %s %d, got %d`, label, expected, got)
	}
}

func assertStringSame(t *testing.T, label string, expected string, got string) {
	if expected != got {
		t.Errorf(`expected %s %s, got %s`, label, expected, got)
	}
}

func assertItemSame(t *testing.T, label string, expected *Item, item *Item) {
	assertNameSame(t, expected.Name, item.Name)
	assertIntSame(t, "item.Pos", expected.Pos, item.Pos)
	assertStringSame(t, "item.ID", expected.ID, item.ID)
	assertStringSame(t, "item.EMail", expected.EMail, item.EMail)
	assertStringSame(t, "item.ShowEMail", expected.ShowEMail, item.ShowEMail)
	assertStringSame(t, "item.Role", expected.Role, item.Role)
	assertStringSame(t, "item.URL", expected.URL, item.URL)
	assertStringSame(t, "item.Type", expected.Type, item.Type)
	assertStringSame(t, "item.Description", expected.Description, item.Description)
	assertStringSame(t, "item.Agency", expected.Agency, item.Agency)
	assertStringSame(t, "item.GrantNumber", expected.GrantNumber, item.GrantNumber)
	assertStringSame(t, "item.URI", expected.URI, item.URI)
	assertStringSame(t, "item.ORCID", expected.ORCID, item.ORCID)
	assertStringSame(t, "item.ROR", expected.ROR, item.ROR)
	assertStringSame(t, "item.Timestamp", expected.Timestamp, item.Timestamp)
	assertStringSame(t, "item.Status", expected.Status, item.Status)
	assertStringSame(t, "item.ReportedBy", expected.ReportedBy, item.ReportedBy)
	assertStringSame(t, "item.ResolvedBy", expected.ResolvedBy, item.ResolvedBy)
	assertStringSame(t, "item.Comment", expected.Comment, item.Comment)
	assertStringSame(t, "item.Value", expected.Value, item.Value)
}

func assertEPrintSame(t *testing.T, expected *EPrint, eprint *EPrint) {
	if (expected == nil) && (eprint != nil) {
		t.Errorf("Expected nil eprint")
		t.FailNow()
	}
	if (expected != nil) && (eprint == nil) {
		t.Errorf(`Expected not nil eprint`)
		t.FailNow()
	}
	assertIntSame(t, "EPrintID", expected.EPrintID, eprint.EPrintID)
	assertStringSame(t, "Title", expected.Title, eprint.Title)
	assertStringSame(t, "EPrintStatus", expected.EPrintStatus, eprint.EPrintStatus)
	assertIntSame(t, "UserID", expected.UserID, eprint.UserID)
	assertStringSame(t, "Datestamp", expected.Datestamp, eprint.Datestamp)
	assertIntSame(t, "DatestampYear", expected.DatestampYear, eprint.DatestampYear)
	assertIntSame(t, "DatestampMonth", expected.DatestampMonth, eprint.DatestampMonth)
	assertIntSame(t, "DatestampDay", expected.DatestampDay, eprint.DatestampDay)
	assertIntSame(t, "DatestampHour", expected.DatestampHour, eprint.DatestampHour)
	assertIntSame(t, "DatestampMinute", expected.DatestampMinute, eprint.DatestampMinute)
	assertIntSame(t, "DatestmapSecond", expected.DatestampSecond, eprint.DatestampSecond)
	assertStringSame(t, "Abstract", expected.Abstract, eprint.Abstract)
	if expected.Creators.Length() != eprint.Creators.Length() {
		t.Errorf(`expected eprint creators length %d, got %d`, expected.Creators.Length(), eprint.Creators.Length())
	} else {
		for i := 0; i < expected.Creators.Length(); i++ {
			t.Logf("checking creator %d", i)
			expectedItem := expected.Creators.IndexOf(i)
			eprintItem := eprint.Creators.IndexOf(i)
			assertItemSame(t, fmt.Sprintf("creators[%d]", i), expectedItem, eprintItem)
		}
	}
	if expected.Editors.Length() != eprint.Editors.Length() {
		t.Errorf(`expected eprint editors length %d, got %d`, expected.Editors.Length(), eprint.Editors.Length())
	} else {
		for i := 0; i < expected.Editors.Length(); i++ {
			t.Logf("checking editor %d", i)
			expectedItem := expected.Editors.IndexOf(i)
			eprintItem := eprint.Editors.IndexOf(i)
			assertItemSame(t, fmt.Sprintf("editors[%d]", i), expectedItem, eprintItem)
		}
	}
	if expected.Contributors.Length() != eprint.Contributors.Length() {
		t.Errorf(`expected eprint contributors length %d, got %d`, expected.Contributors.Length(), eprint.Contributors.Length())
	} else {
		for i := 0; i < expected.Contributors.Length(); i++ {
			t.Logf("checking contributor %d", i)
			expectedItem := expected.Contributors.IndexOf(i)
			eprintItem := eprint.Contributors.IndexOf(i)
			assertItemSame(t, fmt.Sprintf("contributors[%d]", i), expectedItem, eprintItem)
		}
	}
	if expected.CorpCreators.Length() != eprint.CorpCreators.Length() {
		t.Errorf(`expected eprint corp creators length %d, got %d`, expected.CorpCreators.Length(), eprint.CorpCreators.Length())
	} else {
		for i := 0; i < expected.CorpCreators.Length(); i++ {
			t.Logf("checking corp creator %d", i)
			expectedItem := expected.CorpCreators.IndexOf(i)
			eprintItem := eprint.CorpCreators.IndexOf(i)
			assertItemSame(t, fmt.Sprintf("corp creators[%d]", i), expectedItem, eprintItem)
		}
	}
	assertStringSame(t, "DateType", expected.DateType, eprint.DateType)
	assertStringSame(t, "Date", expected.Date, eprint.Date)
	assertIntSame(t, "DateYear", expected.DateYear, eprint.DateYear)
	assertIntSame(t, "DateMonth", expected.DateMonth, eprint.DateMonth)
	assertIntSame(t, "DateDay", expected.DateDay, eprint.DateDay)
	assertStringSame(t, "LastModified", expected.LastModified, eprint.LastModified)
	assertIntSame(t, "LastModifiedYear", expected.LastModifiedYear, eprint.LastModifiedYear)
	assertIntSame(t, "LastModifiedMonth", expected.LastModifiedMonth, eprint.LastModifiedMonth)
	assertIntSame(t, "LastModifiedDay", expected.LastModifiedDay, eprint.LastModifiedDay)
	assertIntSame(t, "LastModifiedHour", expected.LastModifiedHour, eprint.LastModifiedHour)
	assertIntSame(t, "LastModifiedMinute", expected.LastModifiedMinute, eprint.LastModifiedMinute)
	assertIntSame(t, "LastModifiedSecond", expected.LastModifiedSecond, eprint.LastModifiedSecond)

	assertStringSame(t, "StatusChanged", expected.StatusChanged, eprint.StatusChanged)
	assertIntSame(t, "StatusChangedYear", expected.StatusChangedYear, eprint.StatusChangedYear)
	assertIntSame(t, "StatusChangedMonth", expected.StatusChangedMonth, eprint.StatusChangedMonth)
	assertIntSame(t, "StatusChangedDay", expected.StatusChangedDay, eprint.StatusChangedDay)
	assertIntSame(t, "StatusChangedHour", expected.StatusChangedHour, eprint.StatusChangedHour)
	assertIntSame(t, "StatusChangedMinute", expected.StatusChangedMinute, eprint.StatusChangedMinute)
	assertIntSame(t, "StatusChangedSecond", expected.StatusChangedSecond, eprint.StatusChangedSecond)
	//FIXME: check the rest of the fields.
	t.Errorf("Additional field tests need to be implemented")
}

//
// TestCrosswalkEPrintToSQLCreate expects a writable
// "lemurprints" repository to be configured. If not
// the test will be skipped. The lemurprints repository
// database schema should match the schema of CaltechAUTHORS.
//
func TestCrosswalkEPrintToSQLCreate(t *testing.T) {
	var err error
	fName := `test-settings.json`
	repoID := `lemurprints`
	config, err := LoadConfig(fName)
	if err != nil {
		t.Errorf("Cailed to reload %q, %s", fName, err)
	}
	ds, ok := config.Repositories[repoID]
	if ds == nil || ok == false || ds.Write == false {
		t.Skipf(`%s not available for testing`, repoID)
		t.SkipNow()
	}
	baseURL := ds.BaseURL
	assertOpenConnection(t, config, repoID)
	defer assertCloseConnection(t, config, repoID)

	// FIXME: Cleanup any data associated with the test repository

	userID := os.Getuid()
	now := time.Now()
	year, month, day := now.Date()
	hour, minute, second := now.Hour(), now.Minute(), now.Second()
	idNumber := fmt.Sprintf(`DLD-TEST:%d%d%d%d%d%d.%d`, year, month, day, hour, minute, second, userID)

	eprint := new(EPrint)
	eprint.EPrintID = 0
	eprint.Title = `TestCrosswalkEPrintToSQLCreate()`
	eprint.EPrintStatus = "archive"
	eprint.UserID = userID
	eprint.Datestamp = now.Format(`2006-01-02 15:04:05`)
	eprint.DatestampYear = year
	eprint.DatestampMonth = int(month)
	eprint.DatestampDay = day
	eprint.DatestampHour = hour
	eprint.DatestampMinute = minute
	eprint.DatestampSecond = second
	eprint.Abstract = `This is an example test recorded
generated in TestCrosswalkEPrintToSQLCreate() in ep3sql_test.go.`

	eprint.Creators = new(CreatorItemList)
	item := new(Item)
	item.Name = new(Name)
	item.Name.Family = `Doe`
	item.Name.Given = `Jane`
	item.Name.ID = `Doe-Jane`
	item.ORCID = `0000-0000-0000-0000`
	eprint.Creators.Append(item)

	eprint.Editors = new(EditorItemList)
	item = new(Item)
	item.Name = new(Name)
	item.Name.Family = `Doe`
	item.Name.Given = `Jill`
	item.Name.ID = `Doe-Jill`
	item.ORCID = `0000-0000-0000-0001`
	eprint.Editors.Append(item)

	eprint.Contributors = new(ContributorItemList)
	item = new(Item)
	item.Name = new(Name)
	item.Name.Family = `Doe`
	item.Name.Given = `Jack`
	item.Name.ID = `Doi-Jack`
	item.ORCID = `0000-0000-0000-0002`
	eprint.Contributors.Append(item)
	item = new(Item)
	item.Name = new(Name)
	item.Name.Family = `Doe`
	item.Name.Given = `Jaqualine`
	item.Name.ID = `Doe-Jaqualine`
	item.ORCID = `0000-0000-0000-0003`
	eprint.Contributors.Append(item)

	eprint.CorpCreators = new(CorpCreatorItemList)
	item = new(Item)
	item.Name = new(Name)
	item.Name.Value = `Acme, Experimental Labratories`
	item.URI = `uri://example.library.edu/Acme-Experimental-Labrarories`
	eprint.CorpCreators.Append(item)

	eprint.Funders = new(FunderItemList)
	item = new(Item)
	item.Agency = `Digital Libraries Group`
	item.GrantNumber = `DLD-R-000000.007`
	eprint.Funders.Append(item)

	eprint.DateType = "publication"
	eprint.Date = now.Format(`2006-01-02`)
	eprint.DateYear = year
	eprint.DateMonth = int(month)
	eprint.DateDay = day

	eprint.LastModified = now.Format(`2006-01-02 15:04:05`)
	eprint.LastModifiedYear = year
	eprint.LastModifiedMonth = int(month)
	eprint.LastModifiedDay = day
	eprint.LastModifiedHour = hour
	eprint.LastModifiedMinute = minute
	eprint.LastModifiedSecond = second

	eprint.EPrintStatus = `archive`
	eprint.StatusChanged = now.Format(`2006-01-02 15:04:05`)
	eprint.StatusChangedYear = year
	eprint.StatusChangedMonth = int(month)
	eprint.StatusChangedDay = day
	eprint.StatusChangedHour = hour
	eprint.StatusChangedMinute = minute
	eprint.StatusChangedSecond = second

	eprint.DOI = fmt.Sprintf(`0000.00/%s`, idNumber)
	eprint.RevNumber = 1
	eprint.MetadataVisibility = `show`
	eprint.FullTextStatus = `public`
	eprint.Type = `article`
	eprint.IsPublished = `pub`
	eprint.Keywords = `EPrints, Golang, API, Testing`
	eprint.Note = `This is a test record, simulating an article`
	eprint.Suggestions = `This is where suggestions go`
	eprint.Publication = `DLD Software Testing and Development`
	eprint.Volume = `1`
	eprint.Number = `2`
	eprint.Pages = 3
	eprint.PageRange = `15 - 18`
	eprint.PlaceOfPub = `Los Angeles, California, USA`
	eprint.Edition = `1st`
	eprint.Refereed = `TRUE`
	eprint.Series = `Software Testing Practice`
	eprint.IDNumber = idNumber
	eprint.OfficialURL = fmt.Sprintf(`https://resolver.example.edu/%s`, idNumber)
	eprint.Publisher = `The Library`
	eprint.ISSN = `0000-0000`
	eprint.Rights = `No commercial reproduction, distribution, display or performance rights in this work are provided.`
	eprint.OfficialCitation = fmt.Sprintf(`Doe, Jaqualine;&quot;An Alternate API for EPrints&quot;in DLD Software Testing and Development; vol. 1, pp. 15-18, doi %s`, eprint.DOI)
	eprint.Collection = `Test Records`
	eprint.Reviewer = `George Harrison`

	id, err := CrosswalkEPrintToSQLCreate(config, repoID, ds, eprint)
	if err != nil {
		t.Errorf("%s, %s", repoID, err)
		t.FailNow()
	}
	if id == 0 {
		t.Errorf("%s, failed to return non-zero eprint id", repoID)
		t.FailNow()
	}
	eprintCopy, err := CrosswalkSQLToEPrint(config, repoID, baseURL, id)
	if err != nil {
		t.Errorf("%s, %s", repoID, err)
		t.FailNow()
	}
	assertEPrintSame(t, eprint, eprintCopy)
}
