package eprinttools

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
	"time"
)

var (
	err error
)

func objToString(obj interface{}) string {
	src, _ := json.MarshalIndent(obj, "", "    ")
	return string(src)
}

func assertOpenConnection(t *testing.T, config *Config, repoID string) {
	ds, ok := config.Repositories[repoID]
	if !ok {
		t.Skipf("can't find %q", repoID)
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

func clearTable(t *testing.T, config *Config, repoID string, tableName string) {
	db, ok := config.Connections[repoID]
	if !ok {
		t.Skipf("can't find connection for %q", repoID)
		t.SkipNow()
	}
	stmt := fmt.Sprintf(`DELETE FROM %s`, tableName)
	if _, err := db.Exec(stmt); err != nil {
		t.Errorf("Can't create table %q in %q, %s, ", tableName, repoID, err)
	}
}

func assertNameSame(t *testing.T, expected *Name, name *Name) {
	if expected == nil && name != nil {
		t.Errorf(`expected nil name, got %+v`, name)
		t.FailNow()
	}
	if expected != nil && name == nil {
		t.Errorf(`expected non-nil name for %+v`, expected)
		t.FailNow()

	}
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
	if strings.TrimSpace(expected.Value) != strings.TrimSpace(name.Value) {
		t.Errorf(`expected name value %q, got %q`, expected.Value, name.Value)
	}
}

func assertIntSame(t *testing.T, label string, expected int, got int) {
	if expected != got {
		t.Errorf(`expected %s %d, got %d`, label, expected, got)
	}
}

func assertFloat64Same(t *testing.T, label string, expected float64, got float64) {
	if expected != got {
		t.Errorf(`expected %s %f, got %f`, label, expected, got)
	}
}

func assertStringSame(t *testing.T, label string, expected string, got string) {
	if expected != got {
		t.Errorf(`expected %s %q, got %q`, label, expected, got)
	}
}

func assertItemSame(t *testing.T, label string, eprintid int, pos int, expected *Item, item *Item) {
	if expected.Name != nil {
		assertNameSame(t, expected.Name, item.Name)
	}
	if expected.Pos != item.Pos {
		src1 := objToString(expected)
		src2 := objToString(item)
		t.Logf(`
expected (pos: %d) -> %s
item (pos: %d) -> %s
`, expected.Pos, src1, item.Pos, src2)
	}
	assertIntSame(t, fmt.Sprintf(`(%s, %d, %d) item.Pos`, label, eprintid, pos), expected.Pos, item.Pos)
	assertStringSame(t, fmt.Sprintf(`(%s, %d, %d) item.ID`, label, eprintid, pos), expected.ID, item.ID)
	assertStringSame(t, fmt.Sprintf(`(%s, %d, %d) item.EMail`, label, eprintid, pos), expected.EMail, item.EMail)
	assertStringSame(t, fmt.Sprintf(`(%s, %d, %d) item.ShowEMail`, label, eprintid, pos), expected.ShowEMail, item.ShowEMail)
	assertStringSame(t, fmt.Sprintf(`(%s, %d, %d) item.Role`, label, eprintid, pos), expected.Role, item.Role)
	assertStringSame(t, fmt.Sprintf(`(%s, %d, %d) item.URL`, label, eprintid, pos), expected.URL, item.URL)
	assertStringSame(t, fmt.Sprintf(`(%s, %d, %d) item.Type`, label, eprintid, pos), expected.Type, item.Type)
	assertStringSame(t, fmt.Sprintf(`(%s, %d, %d) item.Description`, label, eprintid, pos), expected.Description, item.Description)
	assertStringSame(t, fmt.Sprintf(`(%s, %d, %d) item.Agency`, label, eprintid, pos), expected.Agency, item.Agency)
	assertStringSame(t, fmt.Sprintf(`(%s, %d, %d) item.GrantNumber`, label, eprintid, pos), expected.GrantNumber, item.GrantNumber)
	assertStringSame(t, fmt.Sprintf(`(%s, %d, %d) item.URI`, label, eprintid, pos), expected.URI, item.URI)
	assertStringSame(t, fmt.Sprintf(`(%s, %d, %d) item.ORCID`, label, eprintid, pos), expected.ORCID, item.ORCID)
	assertStringSame(t, fmt.Sprintf(`(%s, %d, %d) item.ROR`, label, eprintid, pos), expected.ROR, item.ROR)
	assertStringSame(t, fmt.Sprintf(`(%s, %d, %d) item.Timestamp`, label, eprintid, pos), expected.Timestamp, item.Timestamp)
	assertStringSame(t, fmt.Sprintf(`(%s, %d, %d) item.Status`, label, eprintid, pos), expected.Status, item.Status)
	assertStringSame(t, fmt.Sprintf(`(%s, %d, %d) item.ReportedBy`, label, eprintid, pos), expected.ReportedBy, item.ReportedBy)
	assertStringSame(t, fmt.Sprintf(`(%s, %d, %d) item.ResolvedBy`, label, eprintid, pos), expected.ResolvedBy, item.ResolvedBy)
	assertStringSame(t, fmt.Sprintf(`(%s, %d, %d) item.Comment`, label, eprintid, pos), expected.Comment, item.Comment)
	assertStringSame(t, fmt.Sprintf(`(%s, %d, %d) item.Value`, label, eprintid, pos), strings.TrimSpace(expected.Value), strings.TrimSpace(item.Value))
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
	assertIntSame(t, "RevNumber", expected.RevNumber, eprint.RevNumber)
	assertStringSame(t, "EPrintStatus", expected.EPrintStatus, eprint.EPrintStatus)
	assertIntSame(t, "UserID", expected.UserID, eprint.UserID)
	assertStringSame(t, "Datestamp", expected.Datestamp, eprint.Datestamp)
	assertIntSame(t, "DatestampYear", expected.DatestampYear, eprint.DatestampYear)
	assertIntSame(t, "DatestampMonth", expected.DatestampMonth, eprint.DatestampMonth)
	assertIntSame(t, "DatestampDay", expected.DatestampDay, eprint.DatestampDay)
	assertIntSame(t, "DatestampHour", expected.DatestampHour, eprint.DatestampHour)
	assertIntSame(t, "DatestampMinute", expected.DatestampMinute, eprint.DatestampMinute)
	assertIntSame(t, "DatestmapSecond", expected.DatestampSecond, eprint.DatestampSecond)
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

	assertStringSame(t, "MetadataVisilibity", expected.MetadataVisibility, eprint.MetadataVisibility)

	assertStringSame(t, "Title", expected.Title, eprint.Title)
	assertStringSame(t, "IsPublished", expected.IsPublished, eprint.IsPublished)
	assertStringSame(t, "FullTextStatus", expected.FullTextStatus, eprint.FullTextStatus)
	assertStringSame(t, "Keywords", expected.Keywords, eprint.Keywords)
	assertStringSame(t, "Note", expected.Note, eprint.Note)
	assertStringSame(t, "Abstract", expected.Abstract, eprint.Abstract)
	assertStringSame(t, "Date", expected.Date, eprint.Date)
	assertIntSame(t, "DateYear", expected.DateYear, eprint.DateYear)
	assertIntSame(t, "DateMonth", int(expected.DateMonth), int(eprint.DateMonth))
	assertIntSame(t, "DateDay", expected.DateDay, eprint.DateDay)
	assertStringSame(t, "DateType", expected.DateType, eprint.DateType)
	assertStringSame(t, "Series", expected.Series, eprint.Series)
	assertStringSame(t, "Publiction", expected.Publication, eprint.Publication)
	assertStringSame(t, "PlaceOfPub", expected.PlaceOfPub, eprint.PlaceOfPub)
	assertStringSame(t, "Edition", expected.Edition, eprint.Edition)
	assertStringSame(t, "PageRange", expected.PageRange, eprint.PageRange)
	assertIntSame(t, "Pages", expected.Pages, eprint.Pages)
	assertStringSame(t, "EventType", expected.EventType, eprint.EventType)
	assertStringSame(t, "EventTitle", expected.EventTitle, eprint.EventTitle)
	assertStringSame(t, "EventLocation", expected.EventLocation, eprint.EventLocation)
	assertStringSame(t, "EventDates", expected.EventDates, eprint.EventDates)
	assertStringSame(t, "IDNumber", expected.IDNumber, eprint.IDNumber)
	assertStringSame(t, "Refereed", expected.Refereed, eprint.Refereed)
	assertStringSame(t, "ISSN", expected.ISSN, eprint.ISSN)
	assertStringSame(t, "ISBN", expected.ISBN, eprint.ISBN)
	assertStringSame(t, "BookTitle", expected.BookTitle, eprint.BookTitle)
	assertStringSame(t, "OffialURL", expected.OfficialURL, eprint.OfficialURL)
	assertStringSame(t, "AltURL", expected.AltURL, eprint.AltURL)
	assertStringSame(t, "Rights", expected.Rights, eprint.Rights)
	assertStringSame(t, "Collection", expected.Collection, eprint.Collection)
	assertStringSame(t, "Reviwer", expected.Reviewer, eprint.Reviewer)
	assertStringSame(t, "OfficialCitation", expected.OfficialCitation, eprint.OfficialCitation)
	assertStringSame(t, "ErrataText", expected.ErrataText, eprint.ErrataText)
	assertStringSame(t, "MonographType", expected.MonographType, eprint.MonographType)
	assertStringSame(t, "Suggestions", expected.Suggestions, eprint.Suggestions)
	assertStringSame(t, "CoverageDates", expected.CoverageDates, eprint.CoverageDates)
	assertStringSame(t, "PresType", expected.PresType, eprint.PresType)
	assertIntSame(t, "Succeeds", expected.Succeeds, eprint.Succeeds)
	assertIntSame(t, "Commentary", expected.Commentary, eprint.Commentary)
	assertStringSame(t, "ContactEMail", expected.ContactEMail, eprint.ContactEMail)
	assertStringSame(t, "FileInfo", expected.FileInfo, eprint.FileInfo)
	assertFloat64Same(t, "Latitude", expected.Latitude, eprint.Latitude)
	assertFloat64Same(t, "Longitude", expected.Longitude, eprint.Longitude)
	assertIntSame(t, "ItemIssuesCount", expected.ItemIssuesCount, eprint.ItemIssuesCount)
	assertStringSame(t, "Department", expected.Department, eprint.Department)
	assertStringSame(t, "OutputMedia", expected.OutputMedia, eprint.OutputMedia)
	assertIntSame(t, "NumPieces", expected.NumPieces, eprint.NumPieces)
	assertStringSame(t, "CompositionType", expected.CompositionType, eprint.CompositionType)
	assertStringSame(t, "DataType", expected.DataType, eprint.DataType)
	assertStringSame(t, "PedagogicType", expected.PedagogicType, eprint.PedagogicType)
	assertStringSame(t, "CompletionTime", expected.CompletionTime, eprint.CompletionTime)
	assertStringSame(t, "TaskPurpose", expected.TaskPurpose, eprint.TaskPurpose)
	assertStringSame(t, "LearningLevelText", expected.LearningLevelText, eprint.LearningLevelText)
	assertStringSame(t, "DOI", expected.DOI, eprint.DOI)
	assertStringSame(t, "PMCID", expected.PMCID, eprint.PMCID)
	assertStringSame(t, "PMID", expected.PMID, eprint.PMID)
	assertStringSame(t, "ParentURL", expected.ParentURL, eprint.ParentURL)
	assertStringSame(t, "TOC", expected.TOC, eprint.TOC)
	assertStringSame(t, "Interviewer", expected.Interviewer, eprint.Interviewer)
	assertStringSame(t, "InterviewDate", expected.InterviewDate, eprint.InterviewDate)
	assertStringSame(t, "NonSubjKeywords", expected.NonSubjKeywords, eprint.NonSubjKeywords)
	assertStringSame(t, "Season", expected.Season, eprint.Season)
	assertStringSame(t, "ClassificationCode", expected.ClassificationCode, eprint.ClassificationCode)
	assertStringSame(t, "SwordDepository", expected.SwordDepository, eprint.SwordDepository)
	assertIntSame(t, "SwordDepositor", expected.SwordDepositor, eprint.SwordDepositor)
	assertStringSame(t, "SwordSlug", expected.SwordSlug, eprint.SwordSlug)
	assertIntSame(t, "ImportID", expected.ImportID, eprint.ImportID)
	assertStringSame(t, "PatentApplicant", expected.PatentApplicant, eprint.PatentApplicant)
	assertStringSame(t, "PatentNumber", expected.PatentNumber, eprint.PatentNumber)
	assertStringSame(t, "PatentClassificationText", expected.PatentClassificationText, eprint.PatentClassificationText)
	assertStringSame(t, "Institution", expected.Institution, eprint.Institution)
	assertStringSame(t, "ThesisType", expected.ThesisType, eprint.ThesisType)
	assertStringSame(t, "ThesisDegree", expected.ThesisDegree, eprint.ThesisDegree)
	assertStringSame(t, "ThesisDegreeGrantor", expected.ThesisDegreeGrantor, eprint.ThesisDegreeGrantor)
	assertStringSame(t, "ThesisDegreeDate", expected.ThesisDegreeDate, eprint.ThesisDegreeDate)
	assertIntSame(t, "ThesisDegreeDateYear", expected.ThesisDegreeDateYear, eprint.ThesisDegreeDateYear)
	assertIntSame(t, "ThesisDegreeDateMonth", int(expected.ThesisDegreeDateMonth), int(eprint.ThesisDegreeDateMonth))
	assertIntSame(t, "ThesisDegreeDateDay", expected.ThesisDegreeDateDay, eprint.ThesisDegreeDateDay)
	assertStringSame(t, "ThesisSubmittedDate", expected.ThesisSubmittedDate, eprint.ThesisSubmittedDate)
	assertIntSame(t, "ThesisSubmittedDateYear", expected.ThesisSubmittedDateYear, eprint.ThesisSubmittedDateYear)
	assertIntSame(t, "ThesisSubmittedDateMonth", int(expected.ThesisSubmittedDateMonth), int(eprint.ThesisSubmittedDateMonth))
	assertIntSame(t, "ThesisSubmittedDateDay", expected.ThesisSubmittedDateDay, eprint.ThesisSubmittedDateDay)
	assertStringSame(t, "ThesisDefenseDate", expected.ThesisDefenseDate, eprint.ThesisDefenseDate)
	assertIntSame(t, "ThesisDefenseDateYear", expected.ThesisDefenseDateYear, eprint.ThesisDefenseDateYear)
	assertIntSame(t, "ThesisDefenseDateMonth", int(expected.ThesisDefenseDateMonth), int(eprint.ThesisDefenseDateMonth))
	assertIntSame(t, "ThesisDefenseDateDay", expected.ThesisDefenseDateDay, eprint.ThesisDefenseDateDay)
	assertStringSame(t, "ThesisApprovedDate", expected.ThesisApprovedDate, eprint.ThesisApprovedDate)
	assertIntSame(t, "ThesisApprovedDateYear", expected.ThesisApprovedDateYear, eprint.ThesisApprovedDateYear)
	assertIntSame(t, "ThesisApprovedDateMonth", int(expected.ThesisApprovedDateMonth), int(eprint.ThesisApprovedDateMonth))
	assertIntSame(t, "ThesisApprovedDateDay", expected.ThesisApprovedDateDay, eprint.ThesisApprovedDateDay)
	assertStringSame(t, "ThesisPublicDate", expected.ThesisPublicDate, eprint.ThesisPublicDate)
	assertIntSame(t, "ThesisPublicDateYear", expected.ThesisPublicDateYear, eprint.ThesisPublicDateYear)
	assertIntSame(t, "ThesisPublicDateMonth", int(expected.ThesisPublicDateMonth), int(eprint.ThesisPublicDateMonth))
	assertIntSame(t, "ThesisPublicDateDay", expected.ThesisPublicDateDay, eprint.ThesisPublicDateDay)
	assertStringSame(t, "HideThesisAuthorEMail", expected.HideThesisAuthorEMail, eprint.HideThesisAuthorEMail)
	assertStringSame(t, "GradOfficeApprovalDate", expected.GradOfficeApprovalDate, eprint.GradOfficeApprovalDate)
	assertIntSame(t, "GradOfficeApprovalDateYear", expected.GradOfficeApprovalDateYear, eprint.GradOfficeApprovalDateYear)
	assertIntSame(t, "GradOfficeApprovalDateMonth", int(expected.GradOfficeApprovalDateMonth), int(eprint.GradOfficeApprovalDateMonth))
	assertIntSame(t, "GradOfficeApprovalDateDay", expected.GradOfficeApprovalDateDay, eprint.GradOfficeApprovalDateDay)
	assertStringSame(t, "ThesisAwards", expected.ThesisAwards, eprint.ThesisAwards)
	assertStringSame(t, "ReviewStatus", expected.ReviewStatus, eprint.ReviewStatus)
	assertStringSame(t, "CopyrightStatement", expected.CopyrightStatement, eprint.CopyrightStatement)
	assertStringSame(t, "Source", expected.Source, eprint.Source)
	assertIntSame(t, "ReplacedBy", expected.ReplacedBy, eprint.ReplacedBy)
	assertIntSame(t, "EditLockUser", expected.EditLockUser, eprint.EditLockUser)
	assertIntSame(t, "EditLockSince", expected.EditLockSince, eprint.EditLockSince)
	assertIntSame(t, "EditLockUntil", expected.EditLockUntil, eprint.EditLockUntil)

	if expected.Creators.Length() != eprint.Creators.Length() {
		t.Errorf(`expected eprint creators length %d, got %d`, expected.Creators.Length(), eprint.Creators.Length())
		src1 := objToString(expected.Creators)
		src2 := objToString(eprint.Creators)
		t.Logf(`
expected.Creators -> %s
eprint.Creators -> %s
`, src1, src2)
	} else {
		for i := 0; i < expected.Creators.Length(); i++ {
			expectedItem := expected.Creators.IndexOf(i)
			eprintItem := eprint.Creators.IndexOf(i)
			assertItemSame(t, "creators", expected.EPrintID, i, expectedItem, eprintItem)
		}
	}
	if expected.Editors.Length() != eprint.Editors.Length() {
		t.Errorf(`expected eprint editors length %d, got %d`, expected.Editors.Length(), eprint.Editors.Length())
		src1 := objToString(expected.Editors)
		src2 := objToString(eprint.Editors)
		t.Logf(`
expected.Editors -> %s
eprint.Editors -> %s
`, src1, src2)
	} else {
		for i := 0; i < expected.Editors.Length(); i++ {
			expectedItem := expected.Editors.IndexOf(i)
			eprintItem := eprint.Editors.IndexOf(i)
			assertItemSame(t, "editors", expected.EPrintID, i, expectedItem, eprintItem)
		}
	}
	if expected.Contributors.Length() != eprint.Contributors.Length() {
		t.Errorf(`expected eprint contributors length %d, got %d`, expected.Contributors.Length(), eprint.Contributors.Length())
		src1 := objToString(expected.Contributors)
		src2 := objToString(eprint.Contributors)
		t.Logf(`
expected.Contributors -> %s
eprint.Contributors -> %s
`, src1, src2)
	} else {
		for i := 0; i < expected.Contributors.Length(); i++ {
			expectedItem := expected.Contributors.IndexOf(i)
			eprintItem := eprint.Contributors.IndexOf(i)
			assertItemSame(t, "contributors", expected.EPrintID, i, expectedItem, eprintItem)
		}
	}
	if expected.CorpCreators.Length() != eprint.CorpCreators.Length() {
		t.Errorf(`expected eprint (eprintid %d) corp creators length %d, got %d`, expected.EPrintID, expected.CorpCreators.Length(), eprint.CorpCreators.Length())
		src1 := objToString(expected.CorpCreators)
		src2 := objToString(eprint.CorpCreators)
		t.Logf(`
expected.CorpCreators -> %s
eprint.CorpCreators -> %s
`, src1, src2)
	} else {
		for i := 0; i < expected.CorpCreators.Length(); i++ {
			expectedItem := expected.CorpCreators.IndexOf(i)
			eprintItem := eprint.CorpCreators.IndexOf(i)
			assertItemSame(t, "corp creators", expected.EPrintID, i, expectedItem, eprintItem)
		}
	}
	if expected.CorpContributors.Length() != eprint.CorpContributors.Length() {
		t.Errorf(`expected eprint corp contributors length %d, got %d`, expected.CorpContributors.Length(), eprint.CorpContributors.Length())
		src1 := objToString(expected.CorpContributors)
		src2 := objToString(eprint.CorpContributors)
		t.Logf(`
expected.CorpContributors -> %s
eprint.CorpContributors -> %s
`, src1, src2)
	} else {
		for i := 0; i < expected.CorpContributors.Length(); i++ {
			expectedItem := expected.CorpContributors.IndexOf(i)
			eprintItem := eprint.CorpContributors.IndexOf(i)
			assertItemSame(t, "corp contributors", expected.EPrintID, i, expectedItem, eprintItem)
		}
	}
	if expected.Lyricists.Length() != eprint.Lyricists.Length() {

		t.Errorf(`expected eprint (eprintid %d) lyricists length %d, got %d`, expected.EPrintID, expected.Lyricists.Length(), eprint.Lyricists.Length())
		src1 := objToString(expected.Lyricists)
		src2 := objToString(eprint.Lyricists)
		t.Logf(`
expected.Lyricists -> %s
eprint.Lyricists -> %s
`, src1, src2)
	} else {
		for i := 0; i < expected.Lyricists.Length(); i++ {
			expectedItem := expected.Lyricists.IndexOf(i)
			eprintItem := eprint.Lyricists.IndexOf(i)
			assertItemSame(t, "lyricist", expected.EPrintID, i, expectedItem, eprintItem)
		}
	}
	if expected.Producers.Length() != eprint.Producers.Length() {

		t.Errorf(`expected eprint (eprintid %d) producers length %d, got %d`, expected.EPrintID, expected.Producers.Length(), eprint.Producers.Length())
		src1 := objToString(expected.Producers)
		src2 := objToString(eprint.Producers)
		t.Logf(`
expected.Producers -> %s
eprint.Producers -> %s
`, src1, src2)
	} else {
		for i := 0; i < expected.Producers.Length(); i++ {
			expectedItem := expected.Producers.IndexOf(i)
			eprintItem := eprint.Producers.IndexOf(i)
			assertItemSame(t, "producers", expected.EPrintID, i, expectedItem, eprintItem)
		}
	}
	if expected.Conductors.Length() != eprint.Conductors.Length() {

		t.Errorf(`expected eprint (eprintid %d) conductors length %d, got %d`, expected.EPrintID, expected.Conductors.Length(), eprint.Conductors.Length())
		src1 := objToString(expected.Conductors)
		src2 := objToString(eprint.Conductors)
		t.Logf(`
expected.Conductors -> %s
eprint.Conductors -> %s
`, src1, src2)
	} else {
		for i := 0; i < expected.Conductors.Length(); i++ {
			expectedItem := expected.Conductors.IndexOf(i)
			eprintItem := eprint.Conductors.IndexOf(i)
			assertItemSame(t, "conductor", expected.EPrintID, i, expectedItem, eprintItem)
		}
	}
	if expected.Exhibitors.Length() != eprint.Exhibitors.Length() {

		t.Errorf(`expected eprint (eprintid %d) exhibitor length %d, got %d`, expected.EPrintID, expected.Exhibitors.Length(), eprint.Exhibitors.Length())
		src1 := objToString(expected.Exhibitors)
		src2 := objToString(eprint.Exhibitors)
		t.Logf(`
expected.Exhibitors -> %s
eprint.Exhitbitors -> %s
`, src1, src2)
	} else {
		for i := 0; i < expected.Exhibitors.Length(); i++ {
			expectedItem := expected.Exhibitors.IndexOf(i)
			eprintItem := eprint.Exhibitors.IndexOf(i)
			assertItemSame(t, "exhibitor", expected.EPrintID, i, expectedItem, eprintItem)
		}
	}
	if expected.ThesisAdvisor.Length() != eprint.ThesisAdvisor.Length() {

		t.Errorf(`expected eprint (eprintid %d) thesis advisor length %d, got %d`, expected.EPrintID, expected.ThesisAdvisor.Length(), eprint.ThesisAdvisor.Length())
		src1 := objToString(expected.ThesisAdvisor)
		src2 := objToString(eprint.ThesisAdvisor)
		t.Logf(`
expected.ThesisAdvisor -> %s
eprint.ThesisAdvisor -> %s
`, src1, src2)
	} else {
		for i := 0; i < expected.ThesisAdvisor.Length(); i++ {
			expectedItem := expected.ThesisAdvisor.IndexOf(i)
			eprintItem := eprint.ThesisAdvisor.IndexOf(i)
			assertItemSame(t, "thesis advisor", expected.EPrintID, i, expectedItem, eprintItem)
		}
	}
	if expected.ThesisCommittee.Length() != eprint.ThesisCommittee.Length() {

		t.Errorf(`expected eprint (eprintid %d) thesis committee length %d, got %d`, expected.EPrintID, expected.ThesisCommittee.Length(), eprint.ThesisCommittee.Length())
		src1 := objToString(expected.ThesisCommittee)
		src2 := objToString(eprint.ThesisCommittee)
		t.Logf(`
expected.ThesisCommittee -> %s
eprint.ThesisCommittee -> %s
`, src1, src2)
	} else {
		for i := 0; i < expected.ThesisCommittee.Length(); i++ {
			expectedItem := expected.ThesisCommittee.IndexOf(i)
			eprintItem := eprint.ThesisCommittee.IndexOf(i)
			assertItemSame(t, "thesis committee", expected.EPrintID, i, expectedItem, eprintItem)
		}
	}

	//   RelatedURL
	if expected.RelatedURL.Length() != eprint.RelatedURL.Length() {

		t.Errorf(`expected eprint (eprintid %d) RelatedURL length %d, got %d`, expected.EPrintID, expected.RelatedURL.Length(), eprint.RelatedURL.Length())
		src1 := objToString(expected.RelatedURL)
		src2 := objToString(eprint.RelatedURL)
		t.Logf(`
expected.RelatedURL -> %s
eprint.RelatedURL -> %s
`, src1, src2)
	} else {
		for i := 0; i < expected.RelatedURL.Length(); i++ {
			expectedItem := expected.RelatedURL.IndexOf(i)
			eprintItem := eprint.RelatedURL.IndexOf(i)
			assertItemSame(t, "RelatedURL", expected.EPrintID, i, expectedItem, eprintItem)
		}
	}

	//   ReferenceText
	if expected.ReferenceText.Length() != eprint.ReferenceText.Length() {

		t.Errorf(`expected eprint (eprintid %d) ReferenceText length %d, got %d`, expected.EPrintID, expected.ReferenceText.Length(), eprint.ReferenceText.Length())
		src1 := objToString(expected.ReferenceText)
		src2 := objToString(eprint.ReferenceText)
		t.Logf(`
expected.ReferenceText -> %s
eprint.ReferenceText -> %s
`, src1, src2)
	} else {
		for i := 0; i < expected.ReferenceText.Length(); i++ {
			expectedItem := expected.ReferenceText.IndexOf(i)
			eprintItem := eprint.ReferenceText.IndexOf(i)
			assertItemSame(t, "ReferenceText", expected.EPrintID, i, expectedItem, eprintItem)
		}
	}

	//   Projects
	if expected.Projects.Length() != eprint.Projects.Length() {

		t.Errorf(`expected eprint (eprintid %d) Projects length %d, got %d`, expected.EPrintID, expected.Projects.Length(), eprint.Projects.Length())
		src1 := objToString(expected.Projects)
		src2 := objToString(eprint.Projects)
		t.Logf(`
expected.Projects -> %s
eprint.Projects -> %s
`, src1, src2)
	} else {
		for i := 0; i < expected.Projects.Length(); i++ {
			expectedItem := expected.Projects.IndexOf(i)
			eprintItem := eprint.Projects.IndexOf(i)
			assertItemSame(t, "Projects", expected.EPrintID, i, expectedItem, eprintItem)
		}
	}

	//   Funders
	if expected.Funders.Length() != eprint.Funders.Length() {

		t.Errorf(`expected eprint (eprintid %d) Funders length %d, got %d`, expected.EPrintID, expected.Funders.Length(), eprint.Funders.Length())
		src1 := objToString(expected.Funders)
		src2 := objToString(eprint.Funders)
		t.Logf(`
expected.Funders -> %s
eprint.Funders -> %s
`, src1, src2)
	} else {
		for i := 0; i < expected.Funders.Length(); i++ {
			expectedItem := expected.Funders.IndexOf(i)
			eprintItem := eprint.Funders.IndexOf(i)
			assertItemSame(t, "Funders", expected.EPrintID, i, expectedItem, eprintItem)
		}
	}

	//   OtherNumberSystem
	if expected.OtherNumberingSystem.Length() != eprint.OtherNumberingSystem.Length() {

		t.Errorf(`expected eprint (eprintid %d) OtherNumberingSystem length %d, got %d`, expected.EPrintID, expected.OtherNumberingSystem.Length(), eprint.OtherNumberingSystem.Length())
		src1 := objToString(expected.OtherNumberingSystem)
		src2 := objToString(eprint.OtherNumberingSystem)
		t.Logf(`
expected.OtherNumberingSystem -> %s
eprint.OtherNumberingSystem -> %s
`, src1, src2)
	} else {
		for i := 0; i < expected.OtherNumberingSystem.Length(); i++ {
			expectedItem := expected.OtherNumberingSystem.IndexOf(i)
			eprintItem := eprint.OtherNumberingSystem.IndexOf(i)
			assertItemSame(t, "OtherNumberingSystem", expected.EPrintID, i, expectedItem, eprintItem)
		}
	}

	//   LocalGroup
	if expected.LocalGroup.Length() != eprint.LocalGroup.Length() {

		t.Errorf(`expected eprint (eprintid %d) LocalGroup length %d, got %d`, expected.EPrintID, expected.LocalGroup.Length(), eprint.LocalGroup.Length())
		src1 := objToString(expected.LocalGroup)
		src2 := objToString(eprint.LocalGroup)
		t.Logf(`
expected.LocalGroup -> %s
eprint.LocalGroup -> %s
`, src1, src2)
	} else {
		for i := 0; i < expected.LocalGroup.Length(); i++ {
			expectedItem := expected.LocalGroup.IndexOf(i)
			eprintItem := eprint.LocalGroup.IndexOf(i)
			assertItemSame(t, "LocalGroup", expected.EPrintID, i, expectedItem, eprintItem)
		}
	}

	//   Subjects
	if expected.Subjects.Length() != eprint.Subjects.Length() {

		t.Errorf(`expected eprint (eprintid %d) Subjects length %d, got %d`, expected.EPrintID, expected.Subjects.Length(), eprint.Subjects.Length())
		src1 := objToString(expected.Subjects)
		src2 := objToString(eprint.Subjects)
		t.Logf(`
expected.Subjects -> %s
eprint.Subjects -> %s
`, src1, src2)
	} else {
		for i := 0; i < expected.Subjects.Length(); i++ {
			expectedItem := expected.Subjects.IndexOf(i)
			eprintItem := eprint.Subjects.IndexOf(i)
			assertItemSame(t, "Subjects", expected.EPrintID, i, expectedItem, eprintItem)
		}
	}

	//   ItemIssues
	if expected.ItemIssues.Length() != eprint.ItemIssues.Length() {

		t.Errorf(`expected eprint (eprintid %d) ItemIssues length %d, got %d`, expected.EPrintID, expected.ItemIssues.Length(), eprint.ItemIssues.Length())
		src1 := objToString(expected.ItemIssues)
		src2 := objToString(eprint.ItemIssues)
		t.Logf(`
expected.ItemIssues -> %s
eprint.ItemIssues -> %s
`, src1, src2)
	} else {
		for i := 0; i < expected.ItemIssues.Length(); i++ {
			expectedItem := expected.ItemIssues.IndexOf(i)
			eprintItem := eprint.ItemIssues.IndexOf(i)
			assertItemSame(t, "ItemIssues", expected.EPrintID, i, expectedItem, eprintItem)
		}
	}

	//   Accompaniment
	if expected.Accompaniment.Length() != eprint.Accompaniment.Length() {

		t.Errorf(`expected eprint (eprintid %d) Accompaniment length %d, got %d`, expected.EPrintID, expected.Accompaniment.Length(), eprint.Accompaniment.Length())
		src1 := objToString(expected.Accompaniment)
		src2 := objToString(eprint.Accompaniment)
		t.Logf(`
expected.Accompaniment -> %s
eprint.Accompaniment -> %s
`, src1, src2)
	} else {
		for i := 0; i < expected.Accompaniment.Length(); i++ {
			expectedItem := expected.Accompaniment.IndexOf(i)
			eprintItem := eprint.Accompaniment.IndexOf(i)
			assertItemSame(t, "Accompaniment", expected.EPrintID, i, expectedItem, eprintItem)
		}
	}

	//   SkillArea
	if expected.SkillAreas.Length() != eprint.SkillAreas.Length() {

		t.Errorf(`expected eprint (eprintid %d) SkillAreas length %d, got %d`, expected.EPrintID, expected.SkillAreas.Length(), eprint.SkillAreas.Length())
		src1 := objToString(expected.SkillAreas)
		src2 := objToString(eprint.SkillAreas)
		t.Logf(`
expected.SkillAreas -> %s
eprint.SkillAreas -> %s
`, src1, src2)
	} else {
		for i := 0; i < expected.SkillAreas.Length(); i++ {
			expectedItem := expected.SkillAreas.IndexOf(i)
			eprintItem := eprint.SkillAreas.IndexOf(i)
			assertItemSame(t, "SkillAreas", expected.EPrintID, i, expectedItem, eprintItem)
		}
	}

	//   CopyrightHolders
	if expected.CopyrightHolders.Length() != eprint.CopyrightHolders.Length() {

		t.Errorf(`expected eprint (eprintid %d) CopyrightHolders length %d, got %d`, expected.EPrintID, expected.CopyrightHolders.Length(), eprint.CopyrightHolders.Length())
		src1 := objToString(expected.CopyrightHolders)
		src2 := objToString(eprint.CopyrightHolders)
		t.Logf(`
expected.CopyrightHolders -> %s
eprint.CopyrightHolders -> %s
`, src1, src2)
	} else {
		for i := 0; i < expected.CopyrightHolders.Length(); i++ {
			expectedItem := expected.CopyrightHolders.IndexOf(i)
			eprintItem := eprint.CopyrightHolders.IndexOf(i)
			assertItemSame(t, "CopyrightHolders", expected.EPrintID, i, expectedItem, eprintItem)
		}
	}

	//   Reference
	if expected.Reference.Length() != eprint.Reference.Length() {

		t.Errorf(`expected eprint (eprintid %d) Reference length %d, got %d`, expected.EPrintID, expected.Reference.Length(), eprint.Reference.Length())
		src1 := objToString(expected.Reference)
		src2 := objToString(eprint.Reference)
		t.Logf(`
expected.Reference -> %s
eprint.Reference -> %s
`, src1, src2)
	} else {
		for i := 0; i < expected.Reference.Length(); i++ {
			expectedItem := expected.Reference.IndexOf(i)
			eprintItem := eprint.Reference.IndexOf(i)
			assertItemSame(t, "Reference", expected.EPrintID, i, expectedItem, eprintItem)
		}
	}

	//   ConfCreators
	if expected.ConfCreators.Length() != eprint.ConfCreators.Length() {

		t.Errorf(`expected eprint (eprintid %d) ConfCreators length %d, got %d`, expected.EPrintID, expected.ConfCreators.Length(), eprint.ConfCreators.Length())
		src1 := objToString(expected.ConfCreators)
		src2 := objToString(eprint.ConfCreators)
		t.Logf(`
expected.ConfCreators -> %s
eprint.ConfCreators -> %s
`, src1, src2)
	} else {
		for i := 0; i < expected.ConfCreators.Length(); i++ {
			expectedItem := expected.ConfCreators.IndexOf(i)
			eprintItem := eprint.ConfCreators.IndexOf(i)
			assertItemSame(t, "ConfCreators", expected.EPrintID, i, expectedItem, eprintItem)
		}
	}

	//   AltTitle
	if expected.AltTitle.Length() != eprint.AltTitle.Length() {

		t.Errorf(`expected eprint (eprintid %d) AltTitle length %d, got %d`, expected.EPrintID, expected.AltTitle.Length(), eprint.AltTitle.Length())
		src1 := objToString(expected.AltTitle)
		src2 := objToString(eprint.AltTitle)
		t.Logf(`
expected.AltTitle -> %s
eprint.AltTitle -> %s
`, src1, src2)
	} else {
		for i := 0; i < expected.AltTitle.Length(); i++ {
			expectedItem := expected.AltTitle.IndexOf(i)
			eprintItem := eprint.AltTitle.IndexOf(i)
			assertItemSame(t, "AltTitle", expected.EPrintID, i, expectedItem, eprintItem)
		}
	}

	//   Shelves (NOTE: Not supporting Shelves)

	//   Relation
	if expected.Relation.Length() != eprint.Relation.Length() {

		t.Errorf(`expected eprint (eprintid %d) Relation length %d, got %d`, expected.EPrintID, expected.Relation.Length(), eprint.Relation.Length())
		src1 := objToString(expected.Relation)
		src2 := objToString(eprint.Relation)
		t.Logf(`
expected.Relation -> %s
eprint.Relation -> %s
`, src1, src2)
	} else {
		for i := 0; i < expected.Relation.Length(); i++ {
			expectedItem := expected.Relation.IndexOf(i)
			eprintItem := eprint.Relation.IndexOf(i)
			assertItemSame(t, "Relation", expected.EPrintID, i, expectedItem, eprintItem)
		}
	}

	//   PatentAssignee
	if expected.PatentAssignee.Length() != eprint.PatentAssignee.Length() {

		t.Errorf(`expected eprint (eprintid %d) PatentAssignee length %d, got %d`, expected.EPrintID, expected.PatentAssignee.Length(), eprint.PatentAssignee.Length())
		src1 := objToString(expected.PatentAssignee)
		src2 := objToString(eprint.PatentAssignee)
		t.Logf(`
expected.PatentAssignee -> %s
eprint.PatentAssignee -> %s
`, src1, src2)
	} else {
		for i := 0; i < expected.PatentAssignee.Length(); i++ {
			expectedItem := expected.PatentAssignee.IndexOf(i)
			eprintItem := eprint.PatentAssignee.IndexOf(i)
			assertItemSame(t, "PatentAssignee", expected.EPrintID, i, expectedItem, eprintItem)
		}
	}

	//   RelatedPatents
	if expected.RelatedPatents.Length() != eprint.RelatedPatents.Length() {

		t.Errorf(`expected eprint (eprintid %d) RelatedPatents length %d, got %d`, expected.EPrintID, expected.RelatedPatents.Length(), eprint.RelatedPatents.Length())
		src1 := objToString(expected.RelatedPatents)
		src2 := objToString(eprint.RelatedPatents)
		t.Logf(`
expected.RelatedPatents -> %s
eprint.RelatedPatents -> %s
`, src1, src2)
	} else {
		for i := 0; i < expected.RelatedPatents.Length(); i++ {
			expectedItem := expected.RelatedPatents.IndexOf(i)
			eprintItem := eprint.RelatedPatents.IndexOf(i)
			assertItemSame(t, "RelatedPatents", expected.EPrintID, i, expectedItem, eprintItem)
		}
	}

	//   Divisions
	if expected.Divisions.Length() != eprint.Divisions.Length() {

		t.Errorf(`expected eprint (eprintid %d) Divisions length %d, got %d`, expected.EPrintID, expected.Divisions.Length(), eprint.Divisions.Length())
		src1 := objToString(expected.Divisions)
		src2 := objToString(eprint.Divisions)
		t.Logf(`
expected.Divisions -> %s
eprint.Divisions -> %s
`, src1, src2)
	} else {
		for i := 0; i < expected.Divisions.Length(); i++ {
			expectedItem := expected.Divisions.IndexOf(i)
			eprintItem := eprint.Divisions.IndexOf(i)
			assertItemSame(t, "Divisions", expected.EPrintID, i, expectedItem, eprintItem)
		}
	}

	//   OptionMajor
	if expected.OptionMajor.Length() != eprint.OptionMajor.Length() {

		t.Errorf(`expected eprint (eprintid %d) OptionMajor length %d, got %d`, expected.EPrintID, expected.OptionMajor.Length(), eprint.OptionMajor.Length())
		src1 := objToString(expected.OptionMajor)
		src2 := objToString(eprint.OptionMajor)
		t.Logf(`
expected.OptionMajor -> %s
eprint.OptionMajor -> %s
`, src1, src2)
	} else {
		for i := 0; i < expected.OptionMajor.Length(); i++ {
			expectedItem := expected.OptionMajor.IndexOf(i)
			eprintItem := eprint.OptionMajor.IndexOf(i)
			assertItemSame(t, "OptionMajor", expected.EPrintID, i, expectedItem, eprintItem)
		}
	}

	//   OptionMinor
	if expected.OptionMinor.Length() != eprint.OptionMinor.Length() {

		t.Errorf(`expected eprint (eprintid %d) OptionMinor length %d, got %d`, expected.EPrintID, expected.OptionMinor.Length(), eprint.OptionMinor.Length())
		src1 := objToString(expected.OptionMinor)
		src2 := objToString(eprint.OptionMinor)
		t.Logf(`
expected.OptionMinor -> %s
eprint.OptionMinor -> %s
`, src1, src2)
	} else {
		for i := 0; i < expected.OptionMinor.Length(); i++ {
			expectedItem := expected.OptionMinor.IndexOf(i)
			eprintItem := eprint.OptionMinor.IndexOf(i)
			assertItemSame(t, "OptionMinor", expected.EPrintID, i, expectedItem, eprintItem)
		}
	}
}

//
// TestSQLCreateEPrint expects a writable
// "lemurprints" repository to be configured. If not
// the test will be skipped. The lemurprints repository
// database schema should match the schema of CaltechAUTHORS.
//
func TestSQLCreateEPrint(t *testing.T) {
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
	eprint.Title = `TestSQLCreateEPrint()`
	eprint.EPrintStatus = "archive"
	eprint.UserID = userID
	eprint.Datestamp = now.Format(timestamp)
	eprint.DatestampYear = year
	eprint.DatestampMonth = int(month)
	eprint.DatestampDay = day
	eprint.DatestampHour = hour
	eprint.DatestampMinute = minute
	eprint.DatestampSecond = second

	eprint.LastModified = now.Format(timestamp)
	eprint.LastModifiedYear = year
	eprint.LastModifiedMonth = int(month)
	eprint.LastModifiedDay = day
	eprint.LastModifiedHour = hour
	eprint.LastModifiedMinute = minute
	eprint.LastModifiedSecond = second

	eprint.Title = `This big SQL test`
	eprint.Abstract = `This is an example test recorded
generated in TestSQLCreateEPrint() in ep3sql_test.go.`

	eprint.Creators = new(CreatorItemList)
	item := new(Item)
	item.Name = new(Name)
	item.Name.Family = `Doe`
	item.Name.Given = `Jane`
	item.ID = `Doe-Jane`
	item.ORCID = `0000-0000-0000-0001`
	eprint.Creators.Append(item)

	eprint.Editors = new(EditorItemList)
	item = new(Item)
	item.Name = new(Name)
	item.Name.Family = `Doe`
	item.Name.Given = `Jill`
	item.ID = `Doe-Jill`
	item.ORCID = `0000-0000-0000-0002`
	eprint.Editors.Append(item)

	eprint.Contributors = new(ContributorItemList)
	item = new(Item)
	item.Name = new(Name)
	item.Name.Family = `Doe`
	item.Name.Given = `Jack`
	item.ID = `Doi-Jack`
	item.ORCID = `0000-0000-0000-0003`
	eprint.Contributors.Append(item)
	item = new(Item)
	item.Name = new(Name)
	item.Name.Family = `Doe`
	item.Name.Given = `Jaqualine`
	item.ID = `Doe-Jaqualine`
	item.ORCID = `0000-0000-0000-0004`
	eprint.Contributors.Append(item)

	eprint.CorpCreators = new(CorpCreatorItemList)
	item = new(Item)
	item.Name = new(Name)
	item.Name.Value = `Acme, Experimental Labratories`
	item.URI = `uri://example.library.edu/Acme-Experimental-Labrarories`
	eprint.CorpCreators.Append(item)

	eprint.CorpContributors = new(CorpContributorItemList)
	item = new(Item)
	item.Name = new(Name)
	item.Name.Value = `Lackey, Experimental Underwriters`
	item.URI = `uri://example.library.edu/Lackey-Experimental-Underwriters`
	eprint.CorpContributors.Append(item)

	eprint.Funders = new(FunderItemList)
	item = new(Item)
	item.Agency = `Digital Libraries Group`
	item.GrantNumber = `DLD-R-000000.007`
	item.ROR = `https://ror.org/XXXXXXX/XXXXXXXX.XXX`
	eprint.Funders.Append(item)

	eprint.Conductors = new(ConductorItemList)
	item = new(Item)
	item.Name = new(Name)
	item.Name.Given = `Quazi`
	item.Name.Family = `Moto`
	item.ID = `Moto-Quazi`
	item.ORCID = `9999-9999-9999-9990`
	eprint.Conductors.Append(item)

	eprint.Lyricists = new(LyricistItemList)
	item = new(Item)
	item.Name = new(Name)
	item.Name.Given = `Julia`
	item.Name.Family = `Childs`
	item.ID = `Childs-Julia`
	item.ORCID = `9999-1111-2222-3330`
	eprint.Lyricists.Append(item)

	eprint.Exhibitors = new(ExhibitorItemList)
	item = new(Item)
	item.Name = new(Name)
	item.Name.Family = `of the Rodent Lackey`
	item.Name.Given = `Museo`
	item.ID = `Museo-Rodent-Lackey`
	item.URI = `uri://Museo-Rodent-Lackey`
	eprint.Exhibitors.Append(item)

	eprint.Producers = new(ProducerItemList)
	item = new(Item)
	item.Name = new(Name)
	item.Name.Given = "Billiana"
	item.Name.Family = "Shakepole"
	item.ID = `Billiana-Shakepole`
	eprint.Producers.Append(item)

	eprint.Accompaniment = new(AccompanimentItemList)
	item = new(Item)
	item.Value = `The Julia Child Kitchen Quartz Trio`
	eprint.Accompaniment.Append(item)

	eprint.SkillAreas = new(SkillAreaItemList)
	item = new(Item)
	item.Value = `Visualization and Projection`
	eprint.SkillAreas.Append(item)

	eprint.CopyrightHolders = new(CopyrightHolderItemList)
	item = new(Item)
	item.Value = `James Dean and famous dead people, LLC`
	eprint.CopyrightHolders.Append(item)

	// NOTE: eprint.Relation is not used in our EPRint repositories

	eprint.RelatedURL = new(RelatedURLItemList)
	item = new(Item)
	item.URL = `http://doi.org/XXXX/XXXXXXX.02`
	item.Type = `doi`
	item.Description = `Figures 1 and 2`
	eprint.RelatedURL.Append(item)
	item = new(Item)
	item.URL = `http://doi.org/XXXX/XXXXXXX.03`
	item.Type = `doi`
	item.Description = `Map of region`
	eprint.RelatedURL.Append(item)

	eprint.ReferenceText = new(ReferenceTextItemList)
	item = new(Item)
	item.Value = `Evergreen, S. (2016). Effective data visualization: The right chart for the right data. SAGE.`
	eprint.ReferenceText.Append(item)
	item = new(Item)
	item.Value = `Nussbaumer, K. C. (2015). Storytelling with data: A data visualization guide for business professionals. Wiley & Sons.`
	eprint.ReferenceText.Append(item)
	item = new(Item)
	item.Value = `Few, S. (2012). Show me the numbers: Designing tables and graphs to enlighten. Analytics Press.`
	eprint.ReferenceText.Append(item)

	tablesAndColumns, err := GetTablesAndColumns(config, repoID)
	if err != nil {
		t.Errorf("Can't get tables and columns for %q, %s", repoID, err)
		t.FailNow()
	}
	for tableName := range tablesAndColumns {
		clearTable(t, config, repoID, tableName)
	}

	eprint.OtherNumberingSystem = new(OtherNumberingSystemItemList)
	item = new(Item)
	item.Name = new(Name)
	item.Name.Value = `Ballad of green fluids`
	item.ID = `1111-.II.II.VC`
	eprint.OtherNumberingSystem.Append(item)

	eprint.LocalGroup = new(LocalGroupItemList)
	item = new(Item)
	item.Value = `Hackers and Punters local 333`
	eprint.LocalGroup.Append(item)

	eprint.Subjects = new(SubjectItemList)
	item = new(Item)
	item.Value = `quandries`
	eprint.Subjects.Append(item)
	item = new(Item)
	item.Value = `mathmatics`
	eprint.Subjects.Append(item)

	eprint.PatentAssignee = new(PatentAssigneeItemList)
	item = new(Item)
	item.Value = `Philbert Stickey-Didgets`
	eprint.PatentAssignee.Append(item)

	eprint.RelatedPatents = new(RelatedPatentItemList)
	item = new(Item)
	item.Value = `WO-101010101`
	eprint.RelatedPatents.Append(item)

	eprint.Divisions = new(DivisionItemList)
	item = new(Item)
	item.Value = `Obfuscation`
	eprint.Divisions.Append(item)

	eprint.ThesisAdvisor = new(ThesisAdvisorItemList)
	item = new(Item)
	item.Name = new(Name)
	item.Name.Given = `Flake`
	item.Name.Family = `Agiven`
	item.ID = `Flakey-Chalk`
	eprint.ThesisAdvisor.Append(item)

	eprint.ThesisCommittee = new(ThesisCommitteeItemList)
	item = new(Item)
	item.Name = new(Name)
	item.Name.Given = `Pear`
	item.Name.Family = `Banana`
	item.Role = `chair`
	eprint.ThesisCommittee.Append(item)

	eprint.OptionMajor = new(OptionMajorItemList)
	item = new(Item)
	item.Value = `Dithering`
	eprint.OptionMajor.Append(item)

	eprint.OptionMinor = new(OptionMinorItemList)
	item = new(Item)
	item.Value = `Dilly-Dally`
	eprint.OptionMinor.Append(item)

	eprint.DateType = "publication"
	eprint.Date = now.Format(`2006-01-02`)
	eprint.DateYear = year
	eprint.DateMonth = int(month)
	eprint.DateDay = day

	eprint.EPrintStatus = `archive`
	eprint.StatusChanged = now.Format(timestamp)
	eprint.StatusChangedYear = year
	eprint.StatusChangedMonth = int(month)
	eprint.StatusChangedDay = day
	eprint.StatusChangedHour = hour
	eprint.StatusChangedMinute = minute
	eprint.StatusChangedSecond = second

	eprint.ThesisDegree = "BA"
	eprint.ThesisSubmittedDate = now.Format(datestamp)
	eprint.ThesisSubmittedDateYear = year
	eprint.ThesisSubmittedDateMonth = int(month)
	eprint.ThesisSubmittedDateDay = day
	eprint.ThesisDefenseDate = now.Format(datestamp)
	eprint.ThesisDefenseDateYear = year
	eprint.ThesisDefenseDateMonth = int(month)
	eprint.ThesisDefenseDateDay = day
	eprint.ThesisApprovedDate = now.Format(datestamp)
	eprint.ThesisApprovedDateYear = year
	eprint.ThesisApprovedDateMonth = int(month)
	eprint.ThesisApprovedDateDay = day
	eprint.ThesisPublicDate = now.Format(datestamp)
	eprint.ThesisPublicDateYear = year
	eprint.ThesisPublicDateMonth = int(month)
	eprint.ThesisPublicDateDay = day
	eprint.ThesisAwards = "The R. S. Doiel has a tremedious Ego award for 2021"
	eprint.ThesisDegreeGrantor = "Troope College if Rarified Design"
	eprint.ThesisAdvisor = new(ThesisAdvisorItemList)
	item = new(Item)
	item.Name = new(Name)
	item.Name.Family = "Doiel"
	item.Name.Given = "Mark"
	item.ORCID = "0000-0001-7321-1464"
	item.ID = "Doiel-Mark"
	eprint.ThesisAdvisor.Append(item)
	eprint.ThesisCommittee = new(ThesisCommitteeItemList)
	item = new(Item)
	item.Name = new(Name)
	item.Name.Family = "Doiel"
	item.Name.Given = "Mark"
	item.ORCID = "0000-0001-7321-1464"
	item.ID = "Doiel-Mark"
	eprint.ThesisCommittee.Append(item)

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

	id, err := SQLCreateEPrint(config, repoID, ds, eprint)
	if err != nil {
		t.Errorf("%s, %s", repoID, err)
		t.FailNow()
	}
	if id == 0 {
		t.Errorf("%s, failed to return non-zero eprint id", repoID)
		t.FailNow()
	}
	eprintCopy, err := SQLReadEPrint(config, repoID, baseURL, id)
	if err != nil {
		t.Errorf("%s, %s", repoID, err)
		t.FailNow()
	}
	assertEPrintSame(t, eprint, eprintCopy)
}

func TestImports(t *testing.T) {
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

	tablesAndColumns, err := GetTablesAndColumns(config, repoID)
	if err != nil {
		t.Errorf("Can't get tables and columns for %q, %s", repoID, err)
		t.FailNow()
	}
	for tableName := range tablesAndColumns {
		clearTable(t, config, repoID, tableName)
	}

	testEPrintXML := []string{
		"lemurprints-1.xml",
		"lemurprints-7.xml",
		"lemurprints-8.xml",
		"lemurprints-34.xml",
		"lemurprints-97.xml",
		"lemurprints-105.xml",
		"lemurprints-132.xml",
		"lemurprints-209.xml",
		"lemurprints-260.xml",
		"lemurprints-21235.xml",
		"lemurprints-8599.xml",
		"lemurprints-92759.xml",
	}
	reviewIds := []int{}
	reviewEPrints := new(EPrints)
	for _, fName := range testEPrintXML {
		fName := path.Join("srctest", fName)
		src, err := ioutil.ReadFile(fName)
		if err != nil {
			t.Errorf("Can't read %q, %s", fName, err)
			t.FailNow()
		}
		eprints := new(EPrints)
		if err := xml.Unmarshal(src, &eprints); err != nil {
			t.Errorf("Can't unmarshal XML for %q, %s", fName, err)
			t.FailNow()
		}
		// Modify the import ID before testing import
		for i := 0; i < eprints.Length(); i++ {
			eprints.EPrint[i].ImportID, eprints.EPrint[i].EPrintID = eprints.EPrint[i].EPrintID, 0
			reviewEPrints.Append(eprints.EPrint[i])
		}
		ids, err := ImportEPrints(config, repoID, ds, eprints)
		if err != nil {
			t.Errorf("Failed to create eprint for %q, %s", fName, err)
			t.FailNow()
		}
		if len(ids) == 0 {
			t.Errorf("Failed to generate new eprint id for %q", fName)
			t.FailNow()
		}
		reviewIds = append(reviewIds, ids...)
	}

	// Make sure we can read out the records that were imported.
	for i, id := range reviewIds {
		eprint := reviewEPrints.IndexOf(i)
		eprintCopy, err := SQLReadEPrint(config, repoID, baseURL, id)
		if err != nil {
			t.Errorf("%s, %s", repoID, err)
			t.FailNow()
		}
		assertEPrintSame(t, eprint, eprintCopy)
	}
}

//
// TestSQLReacEPrint tests read behaviors.
//
func TestSQLRead(t *testing.T) {
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

	// Read an impossible test record.
	eprintID := 123456790
	eprint, err := SQLReadEPrint(config, repoID, baseURL, eprintID)
	if err == nil {
		t.Errorf("Should have gotten an err for EPrint ID %d", eprintID)
	}
	if eprint != nil {
		t.Errorf("Expected a no eprint record for ID %d, got %+v", eprintID, eprint)
	}
}
