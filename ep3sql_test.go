package eprinttools

import (
	"database/sql"
	"testing"
	"time"
)

func assertOpenDBConnection(t *testing.T, repoID string, ds *DataSource) {
	_, ok := config.Connections[repoID]
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

func assertCloseDBConnection(t *testing.T, repoID string) {
	db, ok := config.Connections[repoID]
	if !ok {
		t.Errorf("Could not find connection for %s", repoID)
	}
	if err := db.Close(); err != nil {
		t.Errorf("Failed closing db conncetion for %s, %s", repoID, err)
	}
	t.Logf("DB connection for %s closed", repoID)
}

//
// TestCrosswalkEPrintToSQLCreate expects a writable
// "lemurprints" repository to be configured. If not
// the test will be skipped. The lemurprints repository
// database schema should match the schema of CaltechAUTHORS.
//
func TestCrosswalkEPrintToSQLCreate(t *testing.T) {
	repoID := `lemurprints`
	ds, ok := config.Repositories[repoID]
	if !ok || ds.Write == false {
		t.Skipf(`%s not available for testing`, repoID)
		t.SkipNow()
	}
	assertOpenDBConnection(t, repoID, ds)
	defer assertCloseDBConnection(t, repoID)

	now := time.Now()
	year, month, day := now.Date()

	eprint := new(EPrint)
	eprint.EPrintID = 0
	eprint.Title = `TestCrosswalkEPrintToSQLCreate()`
	eprint.EPrintStatus = "archive"
	eprint.UserID = 1
	eprint.Datestamp = now.Format(`2006-01-02`)
	eprint.DatestampYear = year
	eprint.DatestampMonth = int(month)
	eprint.DatestampDay = day
	eprint.DatestampHour = 0
	eprint.DatestampMinute = 0
	eprint.DatestampSecond = 0
	eprint.Abstract = `This is an example test recorded
generated in TestCrosswalkEPrintToSQLCreate() in ep3sql_test.go.
`

	eprint.Creators = new(CreatorItemList)
	item := new(Item)
	item.Name = new(Name)
	item.Name.Family = `Doe`
	item.Name.Given = `Jane`
	item.Name.ID = `Doe-Jane`
	eprint.Creators.Append(item)
	item = new(Item)
	item.Name = new(Name)
	item.Name.Family = `Doe`
	item.Name.Given = `Jill`
	item.Name.ID = `Doe-Jill`
	eprint.Editors = new(EditorItemList)
	eprint.Editors.Append(item)
	item = new(Item)
	item.Name = new(Name)
	item.Name.Family = `Doe`
	item.Name.Given = `Jack`
	item.Name.ID = `Doi-Jack`
	eprint.Contributors = new(ContributorItemList)
	item = new(Item)
	item.Name = new(Name)
	item.Name.Value = `Acme, Experimental Labratories`
	eprint.Contributors.Append(item)
	eprint.CorpCreators = new(CorpCreatorItemList)
	eprint.DateType = "publication"
	eprint.Date = now.Format(`2006-01-02`)
	eprint.DateYear = year
	eprint.DateMonth = int(month)
	eprint.DateDay = day
	eprint.LastModified = now.Format(`2006-01-02`)
	eprint.LastModifiedYear = year
	eprint.LastModifiedMonth = int(month)
	eprint.LastModifiedDay = day

	id, err := CrosswalkEPrintToSQLCreate(repoID, eprint)
	if err != nil {
		t.Errorf("%s, %s", repoID, err)
		t.FailNow()
	}
	if id == 0 {
		t.Errorf("%s, failed to return non-zero eprint id", repoID)
		t.FailNow()
	}
}
