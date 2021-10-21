package eprinttools

//
// ep3sql.go provides crosswalk methods to/from SQL
//
import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

type EPrintTable struct {
}

// includesString takes an array of string and checks to see
// if the one provided is included.
func includesString(list []string, target string) bool {
	for _, val := range list {
		if val == target {
			return true
		}
	}
	return false
}

// eprintToColumnsAndValues for a given EPrints struct generate a
// list of column names to query along with a recieving values array.
// Return a list of column names (with null handle and aliases) and values.
func eprintToColumnsAndValues(eprint *EPrint, columnsIn []string) ([]string, []interface{}) {
	var (
		columnsOut []string
		values     []interface{}
	)
	for i, key := range columnsIn {
		switch key {
		case "eprintid":
			values = append(values, &eprint.EPrintID)
			columnsOut = append(columnsOut, key)
		case "rev_number":
			values = append(values, &eprint.RevNumber)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "eprint_status":
			values = append(values, &eprint.EPrintStatus)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "userid":
			values = append(values, &eprint.UserID)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "dir":
			values = append(values, &eprint.Dir)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "datestamp_year":
			values = append(values, &eprint.DatestampYear)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "datestamp_month":
			values = append(values, &eprint.DatestampMonth)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "datestamp_day":
			values = append(values, &eprint.DatestampDay)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "datestamp_hour":
			values = append(values, &eprint.DatestampHour)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "datestamp_minute":
			values = append(values, &eprint.DatestampMinute)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "datestamp_second":
			values = append(values, &eprint.DatestampSecond)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "lastmod_year":
			values = append(values, &eprint.LastModifiedYear)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "lastmod_month":
			values = append(values, &eprint.LastModifiedMonth)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "lastmod_day":
			values = append(values, &eprint.LastModifiedDay)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "lastmod_hour":
			values = append(values, &eprint.LastModifiedHour)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "lastmod_minute":
			values = append(values, &eprint.LastModifiedMinute)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "lastmod_second":
			values = append(values, &eprint.LastModifiedSecond)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "status_changed_year":
			values = append(values, &eprint.StatusChangedYear)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "status_changed_month":
			values = append(values, &eprint.StatusChangedMonth)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "status_changed_day":
			values = append(values, &eprint.StatusChangedDay)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "status_changed_hour":
			values = append(values, &eprint.StatusChangedHour)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "status_changed_minute":
			values = append(values, &eprint.StatusChangedMinute)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "status_changed_second":
			values = append(values, &eprint.StatusChangedSecond)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "type":
			values = append(values, &eprint.Type)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "metadata_visibility":
			values = append(values, &eprint.MetadataVisibility)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "title":
			values = append(values, &eprint.Title)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "ispublished":
			values = append(values, &eprint.IsPublished)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "full_text_status":
			values = append(values, &eprint.FullTextStatus)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "keywords":
			values = append(values, &eprint.Keywords)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "note":
			values = append(values, &eprint.Note)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "abstract":
			values = append(values, &eprint.Abstract)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "date_year":
			values = append(values, &eprint.DateYear)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "date_month":
			values = append(values, &eprint.DateMonth)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "date_day":
			values = append(values, &eprint.DateDay)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "date_type":
			values = append(values, &eprint.DateType)
			columnsOut = append(columnsOut, key)
		case "series":
			values = append(values, &eprint.Series)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "volume":
			values = append(values, &eprint.Volume)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "number":
			values = append(values, &eprint.Number)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "publication":
			values = append(values, &eprint.Publication)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "publisher":
			values = append(values, &eprint.Publisher)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "place_of_pub":
			values = append(values, &eprint.PlaceOfPub)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "edition":
			values = append(values, &eprint.Edition)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "pagerange":
			values = append(values, &eprint.PageRange)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "pages":
			values = append(values, &eprint.Pages)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "event_type":
			values = append(values, &eprint.EventType)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "event_title":
			values = append(values, &eprint.EventTitle)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "event_location":
			values = append(values, &eprint.EventLocation)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "event_dates":
			values = append(values, &eprint.EventDates)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "id_number":
			values = append(values, &eprint.IDNumber)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "refereed":
			values = append(values, &eprint.Refereed)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "isbn":
			values = append(values, &eprint.ISBN)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "issn":
			values = append(values, &eprint.ISSN)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "book_title":
			values = append(values, &eprint.BookTitle)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "official_url":
			values = append(values, &eprint.OfficialURL)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "alt_url":
			values = append(values, &eprint.AltURL)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "rights":
			values = append(values, &eprint.Rights)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "collection":
			values = append(values, &eprint.Collection)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "reviewer":
			values = append(values, &eprint.Reviewer)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "official_cit":
			values = append(values, &eprint.OfficialCitation)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "monograph_type":
			values = append(values, &eprint.MonographType)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "suggestions":
			values = append(values, &eprint.Suggestions)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "pres_type":
			values = append(values, &eprint.PresType)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "succeeds":
			values = append(values, &eprint.Succeeds)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "commentary":
			values = append(values, &eprint.Commentary)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "contact_email":
			values = append(values, &eprint.ContactEMail)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "fileinfo":
			values = append(values, &eprint.FileInfo)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "latitude":
			values = append(values, &eprint.Latitude)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "longitude":
			values = append(values, &eprint.Longitude)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "department":
			values = append(values, &eprint.Department)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "output_media":
			values = append(values, &eprint.OutputMedia)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "num_pieces":
			values = append(values, &eprint.NumPieces)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "composition_type":
			values = append(values, &eprint.CompositionType)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "data_type":
			values = append(values, &eprint.DataType)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "pedagogic_type":
			values = append(values, new(string))
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "learning_level":
			values = append(values, &eprint.LearningLevelText)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "completion_time":
			values = append(values, &eprint.CompletionTime)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "task_purpose":
			values = append(values, &eprint.TaskPurpose)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "doi":
			values = append(values, &eprint.DOI)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "pmc_id":
			values = append(values, &eprint.PMCID)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "pmid":
			values = append(values, &eprint.PMID)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "parent_url":
			values = append(values, &eprint.ParentURL)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "toc":
			values = append(values, &eprint.TOC)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "interviewer":
			values = append(values, &eprint.Interviewer)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "interviewdate":
			values = append(values, &eprint.InterviewDate)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "nonsubj_keywords":
			values = append(values, &eprint.NonSubjKeywords)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "season":
			values = append(values, &eprint.Season)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "classification_code":
			values = append(values, &eprint.ClassificationCode)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "sword_depositor":
			values = append(values, &eprint.SwordDepository)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s,'') AS %s`, key, key))
		case "sword_depository":
			values = append(values, &eprint.SwordDepository)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s,'') AS %s`, key, key))
		case "sword_slug":
			values = append(values, &eprint.SwordSlug)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s,'') AS %s`, key, key))
		case "importid":
			values = append(values, &eprint.ImportID)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s,0) AS %s`, key, key))
		case "patent_applicant":
			values = append(values, &eprint.PatentApplicant)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "patent_number":
			values = append(values, &eprint.PatentNumber)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "institution":
			values = append(values, &eprint.Institution)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "thesis_type":
			values = append(values, &eprint.ThesisType)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "thesis_degree":
			values = append(values, &eprint.ThesisDegree)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "thesis_degree_grantor":
			values = append(values, &eprint.ThesisDegreeGrantor)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "thesis_degree_date_year":
			values = append(values, &eprint.ThesisDegreeDateYear)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "thesis_degree_date_month":
			values = append(values, &eprint.ThesisDegreeDateMonth)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "thesis_degree_date_day":
			values = append(values, &eprint.ThesisDegreeDateDay)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "thesis_submitted_date_year":
			values = append(values, &eprint.ThesisSubmittedDateYear)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "thesis_submitted_date_month":
			values = append(values, &eprint.ThesisSubmittedDateMonth)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "thesis_submitted_date_day":
			values = append(values, &eprint.ThesisSubmittedDateDay)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "thesis_defense_date_year":
			values = append(values, &eprint.ThesisDefenseDateYear)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "thesis_defense_date_month":
			values = append(values, &eprint.ThesisDefenseDateMonth)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "thesis_defense_date_day":
			values = append(values, &eprint.ThesisDefenseDateDay)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "thesis_approved_date_year":
			values = append(values, &eprint.ThesisApprovedDateYear)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "thesis_approved_date_month":
			values = append(values, &eprint.ThesisApprovedDateMonth)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "thesis_approved_date_day":
			values = append(values, &eprint.ThesisApprovedDateDay)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "thesis_public_date_year":
			values = append(values, &eprint.ThesisPublicDateYear)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "thesis_public_date_month":
			values = append(values, &eprint.ThesisPublicDateMonth)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "thesis_public_date_day":
			values = append(values, &eprint.ThesisPublicDateDay)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, 0) AS %s`, key, key))
		case "thesis_author_email":
			values = append(values, &eprint.ThesisAuthorEMail)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "hide_thesis_author_email":
			values = append(values, &eprint.HideThesisAuthorEMail)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "gradofc_approval_date":
			values = append(values, &eprint.GradOfficeApprovalDate)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "thesis_awards":
			values = append(values, &eprint.ThesisAwards)
			columnsOut = append(columnsOut, key)
		case "review_status":
			values = append(values, &eprint.ReviewStatus)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "copyright_statement":
			values = append(values, &eprint.CopyrightStatement)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "source":
			values = append(values, &eprint.Source)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s, '') AS %s`, key, key))
		case "replacedby":
			values = append(values, &eprint.ReplacedBy)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s,0) AS %s`, key, key))
		case "item_issues_count":
			values = append(values, &eprint.ItemIssuesCount)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s,0) AS %s`, key, key))
		case "errata":
			values = append(values, &eprint.ErrataText)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s,'') AS %s`, key, key))
		case "edit_lock_user":
			values = append(values, &eprint.EditLockUser)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s,0) AS %s`, key, key))
		case "edit_lock_since":
			values = append(values, &eprint.EditLockSince)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s,0) AS %s`, key, key))
		case "edit_lock_until":
			values = append(values, &eprint.EditLockUntil)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s,0) AS %s`, key, key))
		// The follow values represent sub tables and processed separately.
		case "patent_classification":
			values = append(values, &eprint.PatentClassificationText)
			columnsOut = append(columnsOut, fmt.Sprintf(`IFNULL(%s,'') AS %s`, key, key))
		default:
			// Handle case where we have value that is unmapped or not available in EPrint struct
			log.Printf("could not map %q (%d) into EPrint struct", key, i)
		}

	}
	return columnsOut, values
}

func makePersonName(given string, family string, honourific string, lineage string) *Name {
	name := new(Name)
	isFlat := true
	if s := strings.TrimSpace(given); s != "" {
		name.Given = s
		isFlat = false
	}
	if s := strings.TrimSpace(family); s != "" {
		name.Family = s
		isFlat = false
	}
	if s := strings.TrimSpace(honourific); s != "" {
		name.Honourific = s
		isFlat = false
	}
	if s := strings.TrimSpace(lineage); s != "" {
		name.Lineage = s
		isFlat = false
	}
	if isFlat {
		return nil
	}
	return name
}

func makeTimestamp(year int, month int, day int, hour int, minute int, second int) string {
	return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC).Format("2006-01-02 15:04:05")
}

func makeDatestamp(year int, month int, day int) string {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
}

func makeApproxDate(year int, month int, day int) string {
	var layout string
	switch {
	case (year > 0) && (month > 0) && (day > 0):
		layout = "2006-01-02"
	case (year > 0) && (month > 0) && (day == 0):
		layout = "2006-01"
	case (year > 0) && (month == 0) && (day == 0):
		layout = "2006"
	}
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC).Format(layout)
}

func creatorsToEPrint(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *CreatorItemList {
	var (
		pos                                                       int
		creatorid, honourific, given, family, lineage, orcid, uri string
	)
	subTables := []string{}
	for table := range tables {
		if strings.HasPrefix(table, "eprint_creators") {
			subTables = append(subTables, table)
		}
	}
	if includesString(subTables, "eprint_creators_name") {
		// eprint_%_id is a known structure. eprintid, pos, creators_id
		creatorItemList := new(CreatorItemList)
		stmt := `SELECT pos, IFNULL(creators_name_honourific, '') AS creators_name_honourific, IFNULL(creators_name_given, '') AS creators_name_given, IFNULL(creators_name_family, '') AS creators_name_family, IFNULL(creators_name_lineage, '') AS creators_name_lineage FROM eprint_creators_name WHERE eprintid = ? ORDER BY eprintid, pos`
		rows, err := db.Query(stmt, eprintID)
		if err != nil {
			log.Printf("Could not scan eprint_creators_name for %d in %q, %s", eprintID, repoID, err)
		} else {
			i := 0
			for rows.Next() {
				if err := rows.Scan(&pos, &honourific, &given, &family, &lineage); err != nil {
					log.Printf("Could not scan eprint_creators_name for %d in %q, %s", eprintID, repoID, err)
				} else {
					item := new(Item)
					item.Pos = pos
					item.Name = makePersonName(given, family, honourific, lineage)
					creatorItemList.Items = append(creatorItemList.Items, item)
				}
				i++
			}
			rows.Close()
		}
		if includesString(subTables, "eprint_creators_id") {
			stmt := `SELECT pos, IFNULL(creators_id, "") AS creators_id FROM eprint_creators_id WHERE eprintid = ? ORDER BY eprintid, pos`
			rows, err = db.Query(stmt, eprintID)
			if err != nil {
				log.Printf("Could not scan eprint_creators_name for %d in %q, %s", eprintID, repoID, err)
			} else {
				i := 0
				for rows.Next() {
					if err := rows.Scan(&pos, &creatorid); err != nil {
						log.Printf("Could not scan eprint_creators_name for %d in %q, %s", eprintID, repoID, err)
					} else if i < len(creatorItemList.Items) {
						for _, item := range creatorItemList.Items {
							if item.Pos == pos {
								item.ID = creatorid
								break
							}
						}
					}
					i++
				}
				rows.Close()
			}
		}
		if includesString(subTables, "eprint_creators_orcid") {
			rows, err = db.Query(`SELECT pos, IFNULL(creators_orcid, '') AS orcid FROM eprint_creators_orcid WHERE eprintid = ? ORDER BY eprintid, pos`, eprintID)
			if err != nil {
				log.Printf("Could not scan eprint_creators_orcid for %d in %q, %s", eprintID, repoID, err)
			} else {
				i := 0
				for rows.Next() {
					if err := rows.Scan(&pos, &orcid); err != nil {
						log.Printf("Could not scan eprint_creators_orcid for %d in %q, %s", eprintID, repoID, err)
					} else if i < len(creatorItemList.Items) {
						if orcid != "" {
							for _, item := range creatorItemList.Items {
								if item.Pos == pos {
									item.ORCID = orcid
									break
								}
							}
						}
					}
					i++
				}
				rows.Close()
			}
		}
		if includesString(subTables, "eprint_creators_uri") {
			rows, err = db.Query(`SELECT pos, IFNULL(creators_uri, '') AS uri FROM eprint_creators_uri WHERE eprintid = ? ORDER BY eprintid, pos`, eprintID)
			if err != nil {
				log.Printf("Could not scan eprint_creators_uri for %d in %q, %s", eprintID, repoID, err)
			} else {
				i := 0
				for rows.Next() {
					if err := rows.Scan(&pos, &uri); err != nil {
						log.Printf("Could not scan eprint_creators_uri for %d in %q, %s", eprintID, repoID, err)
					} else if i < len(creatorItemList.Items) {
						for _, item := range creatorItemList.Items {
							if item.Pos == pos {
								item.URI = uri
								break
							}
						}
					}
					i++
				}
				rows.Close()
			}
		}
		return creatorItemList
	}
	return nil
}

func editorsToEPrint(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *EditorItemList {
	var (
		pos                                                      int
		editorid, honourific, given, family, lineage, orcid, uri string
	)
	subTables := []string{}
	for table := range tables {
		if strings.HasPrefix(table, "eprint_editors") {
			subTables = append(subTables, table)
		}
	}
	if includesString(subTables, "eprint_editors_name") {
		// eprint_%_id is a known structure. eprintid, pos, editors_id
		editorItemList := new(EditorItemList)
		stmt := `SELECT pos, IFNULL(editors_name_honourific, '') AS editors_name_honourific, IFNULL(editors_name_given, '') AS editors_name_given, IFNULL(editors_name_family, '') AS editors_name_family, IFNULL(editors_name_lineage, '') AS editors_name_lineage FROM eprint_editors_name WHERE eprintid = ? ORDER BY eprintid, pos`
		rows, err := db.Query(stmt, eprintID)
		if err != nil {
			log.Printf("Could not scan eprint_editors_id for %d in %q, %s", eprintID, repoID, err)
		} else {
			i := 0
			for rows.Next() {
				if err := rows.Scan(&pos, &honourific, &given, &family, &lineage); err != nil {
					log.Printf("Could not scan eprint_editors_id for %d in %q, %s", eprintID, repoID, err)
				} else {
					item := new(Item)
					item.Pos = pos
					item.Name = makePersonName(given, family, honourific, lineage)
					editorItemList.Items = append(editorItemList.Items, item)
				}
				i++
			}
			rows.Close()
		}
		if includesString(subTables, "eprint_editors_id") {
			stmt := `SELECT pos, IFNULL(editors_id, '') AS editors_id FROM eprint_editors_id WHERE eprintid = ? ORDER BY eprintid, pos`
			rows, err = db.Query(stmt, eprintID)
			if err != nil {
				log.Printf("Could not scan eprint_editors_name for %d in %q, %s", eprintID, repoID, err)
			} else {
				i := 0
				for rows.Next() {
					if err := rows.Scan(&pos, &editorid); err != nil {
						log.Printf("Could not scan eprint_editors_name for %d in %q, %s", eprintID, repoID, err)
					} else {
						if editorid != "" {
							for _, item := range editorItemList.Items {
								if item.Pos == pos {
									item.ID = editorid
									break
								}
							}
						}
					}
					i++
				}
				rows.Close()
			}
		}
		if includesString(subTables, "eprint_editors_orcid") {
			rows, err = db.Query(`SELECT pos, IFNULL(editors_orcid, '') AS orcid FROM eprint_editors_orcid WHERE eprintid = ? ORDER BY eprintid, pos`, eprintID)
			if err != nil {
				log.Printf("Could not scan eprint_editors_orcid for %d in %q, %s", eprintID, repoID, err)
			} else {
				i := 0
				for rows.Next() {
					if err := rows.Scan(&pos, &orcid); err != nil {
						log.Printf("Could not scan eprint_editors_orcid for %d in %q, %s", eprintID, repoID, err)
					} else {
						if orcid != "" {
							for _, item := range editorItemList.Items {
								if item.Pos == pos {
									item.ORCID = orcid
									break
								}
							}
						}
					}
					i++
				}
				rows.Close()
			}
		}
		if includesString(subTables, "eprint_editors_uri") {
			rows, err = db.Query(`SELECT pos, IFNULL(editors_uri, '') AS uri FROM eprint_editors_uri WHERE eprintid = ? ORDER BY eprintid, pos`, eprintID)
			if err != nil {
				log.Printf("Could not scan eprint_editors_uri for %d in %q, %s", eprintID, repoID, err)
			} else {
				i := 0
				for rows.Next() {
					if err := rows.Scan(&pos, &uri); err != nil {
						log.Printf("Could not scan eprint_editors_uri for %d in %q, %s", eprintID, repoID, err)
					} else {
						for _, item := range editorItemList.Items {
							if item.Pos == pos {
								item.URI = uri
								break
							}
						}
					}
					i++
				}
				rows.Close()
			}
		}
		if len(editorItemList.Items) > 0 {
			return editorItemList
		}
	}
	return nil
}

func contributorsToEPrint(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *ContributorItemList {
	var (
		pos                                                           int
		contributorid, honourific, given, family, lineage, orcid, uri string
	)
	subTables := []string{}
	for table := range tables {
		if strings.HasPrefix(table, "eprint_contributors") {
			subTables = append(subTables, table)
		}
	}
	if includesString(subTables, "eprint_contributors_name") {
		// eprint_%_id is a known structure. eprintid, pos, contributors_id
		contributorItemList := new(ContributorItemList)
		stmt := `SELECT pos, IFNULL(contributors_name_honourific, '') AS contributors_name_honourific, IFNULL(contributors_name_given, '') AS contributors_name_given, IFNULL(contributors_name_family, '') AS contributors_name_family, IFNULL(contributors_name_lineage, '') AS contributors_name_lineage FROM eprint_contributors_name WHERE eprintid = ? ORDER BY eprintid, pos`
		rows, err := db.Query(stmt, eprintID)
		if err != nil {
			log.Printf("Could not scan eprint_contributors_id for %d in %q, %s", eprintID, repoID, err)
		} else {
			i := 0
			for rows.Next() {
				if err := rows.Scan(&pos, &honourific, &given, &family, &lineage); err != nil {
					log.Printf("Could not scan eprint_contributors_id for %d in %q, %s", eprintID, repoID, err)
				} else {
					item := new(Item)
					item.Pos = pos
					item.Name = makePersonName(given, family, honourific, lineage)
					contributorItemList.Items = append(contributorItemList.Items, item)
				}
				i++
			}
			rows.Close()
		}
		if includesString(subTables, "eprint_contributors_id") {
			stmt := `SELECT pos, IFNULL(contributors_id, '') AS contributors_id FROM eprint_contributors_id WHERE eprintid = ? ORDER BY eprintid, pos`
			rows, err = db.Query(stmt, eprintID)
			if err != nil {
				log.Printf("Could not scan eprint_contributors_name for %d in %q, %s", eprintID, repoID, err)
			} else {
				i := 0
				for rows.Next() {
					if err := rows.Scan(&pos, &contributorid); err != nil {
						log.Printf("Could not scan eprint_contributors_name for %d in %q, %s", eprintID, repoID, err)
					} else {
						if contributorid != "" {
							for _, item := range contributorItemList.Items {
								if item.Pos == pos {
									item.ID = contributorid
									break
								}
							}
						}
					}
					i++
				}
				rows.Close()
			}
		}
		if includesString(subTables, "eprint_contributors_orcid") {
			rows, err = db.Query(`SELECT pos, IFNULL(contributors_orcid, '') AS orcid FROM eprint_contributors_orcid WHERE eprintid = ? ORDER BY eprintid, pos`, eprintID)
			if err != nil {
				log.Printf("Could not scan eprint_contributors_orcid for %d in %q, %s", eprintID, repoID, err)
			} else {
				i := 0
				for rows.Next() {
					if err := rows.Scan(&pos, &orcid); err != nil {
						log.Printf("Could not scan eprint_contributors_orcid for %d in %q, %s", eprintID, repoID, err)
					} else {
						if orcid != "" {
							for _, item := range contributorItemList.Items {
								if item.Pos == pos {
									item.ORCID = orcid
									break
								}
							}
						}
					}
					i++
				}
				rows.Close()
			}
		}
		if includesString(subTables, "eprint_contributors_uri") {
			rows, err = db.Query(`SELECT pos, IFNULL(contributors_uri, '') AS uri FROM eprint_contributors_uri WHERE eprintid = ? ORDER BY eprintid, pos`, eprintID)
			if err != nil {
				log.Printf("Could not scan eprint_contributors_uri for %d in %q, %s", eprintID, repoID, err)
			} else {
				i := 0
				for rows.Next() {
					if err := rows.Scan(&pos, &uri); err != nil {
						log.Printf("Could not scan eprint_contributors_uri for %d in %q, %s", eprintID, repoID, err)
					} else {
						if uri != "" {
							for _, item := range contributorItemList.Items {
								if item.Pos == pos {
									contributorItemList.Items[i].URI = uri
									break
								}
							}
						}
					}
					i++
				}
				rows.Close()
			}
		}
		if len(contributorItemList.Items) > 0 {
			return contributorItemList
		}
	}
	return nil
}

func fundersToEPrint(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *FunderItemList {
	var (
		pos                 int
		agency, grantNumber string
	)
	subTables := []string{}
	for table := range tables {
		if strings.HasPrefix(table, "eprint_funders") {
			subTables = append(subTables, table)
		}
	}
	if includesString(subTables, "eprint_funders_agency") {
		// eprint_%_id is a known structure. eprintid, pos, contributors_id
		funderItemList := new(FunderItemList)
		rows, err := db.Query(`SELECT pos, IFNULL(funders_agency, '') AS agency FROM eprint_funders_agency WHERE eprintid = ? ORDER BY eprintid, pos`, eprintID)
		if err != nil {
			log.Printf("Could not scan eprint_funders_agency for %d in %q, %s", eprintID, repoID, err)
		} else {
			i := 0
			for rows.Next() {
				if err := rows.Scan(&pos, &agency); err != nil {
					log.Printf("Could not scan eprint_funders_agency for %d in %q, %s", eprintID, repoID, err)
				} else {
					item := new(Item)
					item.Pos = pos
					if agency != "" {
						item.Agency = agency
					}
					funderItemList.Items = append(funderItemList.Items, item)
				}
				i++
			}
			rows.Close()
		}
		if includesString(subTables, "eprint_funders_grant_number") {
			rows, err = db.Query(`SELECT pos, IFNULL(funders_grant_number, '') AS grant_number FROM eprint_funders_grant_number WHERE eprintid = ? ORDER BY eprintid, pos`, eprintID)
			if err != nil {
				log.Printf("Could not scan eprint_funders_grant_number for %d in %q, %s", eprintID, repoID, err)
			} else {
				i := 0
				for rows.Next() {
					if err := rows.Scan(&pos, &grantNumber); err != nil {
						log.Printf("Could not scan eprint_contributors_name for %d in %q, %s", eprintID, repoID, err)
					} else {
						if grantNumber != "" {
							for _, item := range funderItemList.Items {
								if item.Pos == pos {
									item.GrantNumber = grantNumber
									break
								}
							}
						}
					}
					i++
				}
				rows.Close()
			}
		}
		if len(funderItemList.Items) > 0 {
			return funderItemList
		}
	}
	return nil
}

func localGroupToEPrint(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *LocalGroupItemList {
	var (
		localGroup string
	)
	_, ok := tables["eprint_local_group"]
	if ok {
		// eprint_%_id is a known structure. eprintid, pos, localGroups_id
		localGroupItemList := new(LocalGroupItemList)
		rows, err := db.Query(`SELECT local_group FROM eprint_local_group WHERE eprintid = ? ORDER BY eprintid, pos`, eprintID)
		if err != nil {
			log.Printf("Could not scan eprint_local_group for %d in %q, %s", eprintID, repoID, err)
		} else {
			i := 0
			for rows.Next() {
				if err := rows.Scan(&localGroup); err != nil {
					log.Printf("Could not scan eprint_local_group for %d in %q, %s", eprintID, repoID, err)
				} else {
					item := new(Item)
					if localGroup != "" {
						item.Value = localGroup
					} else {
						item.Value = ""
					}
					localGroupItemList.Items = append(localGroupItemList.Items, item)
				}
				i++
			}
			rows.Close()
		}
		if len(localGroupItemList.Items) > 0 {
			return localGroupItemList
		}
	}
	return nil
}

func documentToColumnsAndValues(document *Document, columns []string) ([]string, []interface{}) {
	columnsOut := []string{}
	values := []interface{}{}
	for i, key := range columns {
		switch key {
		case "docid":
			values = append(values, &document.DocID)
			columnsOut = append(columnsOut, key)
		case "eprintid":
			values = append(values, &document.EPrintID)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s, 0) AS %s", key, key))
		case "pos":
			values = append(values, &document.Pos)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s, 0) AS %s", key, key))
		case "rev_number":
			values = append(values, &document.RevNumber)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s, 0) AS %s", key, key))
		case "format":
			values = append(values, &document.Format)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s, '') AS %s", key, key))
		case "formatdesc":
			values = append(values, &document.FormatDesc)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s, '') AS %s", key, key))
		case "language":
			values = append(values, &document.Language)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s, '') AS %s", key, key))
		case "security":
			values = append(values, &document.Security)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s, '') AS %s", key, key))
		case "license":
			values = append(values, &document.License)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s, '') AS %s", key, key))
		case "main":
			values = append(values, &document.Main)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s, '') AS %s", key, key))
		case "date_embargo_year":
			values = append(values, &document.DateEmbargoYear)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s, 0) AS %s", key, key))
		case "date_embargo_month":
			values = append(values, &document.DateEmbargoMonth)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s, 0) AS %s", key, key))
		case "date_embargo_day":
			values = append(values, &document.DateEmbargoDay)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s, 0) AS %s", key, key))
		case "content":
			values = append(values, &document.Content)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s, '') AS %s", key, key))
		case "placement":
			values = append(values, &document.Placement)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s, 0) AS %s", key, key))
		case "mime_type":
			values = append(values, &document.MimeType)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s, '') AS %s", key, key))
		case "media_duration":
			values = append(values, &document.MediaDuration)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s, '') AS %s", key, key))
		case "media_audio_codec":
			values = append(values, &document.MediaAudioCodec)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s, '') AS %s", key, key))
		case "media_video_codec":
			values = append(values, &document.MediaVideoCodec)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s, '') AS %s", key, key))
		case "media_width":
			values = append(values, &document.MediaWidth)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s, 0) AS %s", key, key))
		case "media_height":
			values = append(values, &document.MediaHeight)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s, 0) AS %s", key, key))
		case "media_aspect_ratio":
			values = append(values, &document.MediaAspectRatio)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s, '') AS %s", key, key))
		case "media_sample_start":
			values = append(values, &document.MediaSampleStart)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s, '') AS %s", key, key))
		case "media_sample_stop":
			values = append(values, &document.MediaSampleStop)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s, '') AS %s", key, key))
		default:
			log.Printf("%q (%d) not found in document table", key, i)
			//return nil, nil, fmt.Errorf("%q (%d) not found in document", key, i)
		}
	}
	return columnsOut, values
}

func fileToColumnsAndValues(file *File, columns []string) ([]string, []interface{}) {
	columnsOut := []string{}
	values := []interface{}{}
	for i, key := range columns {
		switch key {
		case "fileid":
			values = append(values, &file.FileID)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s,'') AS %s", key, key))
		case "datasetid":
			values = append(values, &file.DatasetID)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s,'') AS %s", key, key))
		case "objectid":
			values = append(values, &file.ObjectID)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s,'') AS %s", key, key))
		case "filename":
			values = append(values, &file.Filename)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s,'') AS %s", key, key))
		case "mime_type":
			values = append(values, &file.MimeType)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s,'') AS %s", key, key))
		case "hash":
			values = append(values, &file.Hash)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s,'') AS %s", key, key))
		case "hash_type":
			values = append(values, &file.HashType)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s,'') AS %s", key, key))
		case "filesize":
			values = append(values, &file.FileSize)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s,'') AS %s", key, key))
		case "mtime_year":
			values = append(values, &file.MTimeYear)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s,'') AS %s", key, key))
		case "mtime_month":
			values = append(values, &file.MTimeMonth)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s,'') AS %s", key, key))
		case "mtime_day":
			values = append(values, &file.MTimeDay)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s,'') AS %s", key, key))
		case "mtime_hour":
			values = append(values, &file.MTimeHour)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s,'') AS %s", key, key))
		case "mtime_minute":
			values = append(values, &file.MTimeMinute)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s,'') AS %s", key, key))
		case "mtime_second":
			values = append(values, &file.MTimeSecond)
			columnsOut = append(columnsOut, fmt.Sprintf("IFNULL(%s,'') AS %s", key, key))
		default:
			log.Printf("%q (%d) not found in file table", key, i)
		}
	}
	return columnsOut, values
}

func documentIDToFiles(repoID string, baseURL string, eprintID int, documentID int, db *sql.DB, tables map[string][]string) []*File {
	//FIXME: Need to figure out if I need to pay attention to
	// file_copies_plugin and file_copies_sourceid tables. This appear
	// to be related to the storage manager.  They don't appear to
	// be reflected in the EPrint XML. files_copies_plugin has a single
	// unique column called copies_pluginid which has a single value
	// of "Storage::Local" for all rows (about the same number of rows
	// as in the file table.  In the table file_copies_sourceid the column
	// copies_sourceid appears to be base filenames like in the
	// file.filename where file.datasetid = "document". I'm ignoring in
	// this processing for file table for now.
	tableName := "file"
	columns, ok := tables[tableName]
	if ok {
		files := []*File{}
		file := new(File)
		columnSQL, values := fileToColumnsAndValues(file, columns)
		stmt := fmt.Sprintf(`SELECT %s FROM %s WHERE datasetid = 'document' AND objectid = ?`, strings.Join(columnSQL, ", "), tableName)
		rows, err := db.Query(stmt, documentID)
		if err != nil {
			log.Printf("Query failed %q for %d in %q, doc ID %d, , %q,  %s", tableName, eprintID, repoID, documentID, stmt, err)
		} else {
			i := 0
			for rows.Next() {
				file = new(File)
				_, values = fileToColumnsAndValues(file, columns)
				if err := rows.Scan(values...); err != nil {
					log.Printf("Could not scan %q for %d in %q, doc ID %d, %s", tableName, eprintID, repoID, documentID, err)
				} else {
					file.ID = fmt.Sprintf("%s/id/file/%d", baseURL, file.FileID)
					file.MTime = makeTimestamp(file.MTimeYear, file.MTimeMonth, file.MTimeDay, file.MTimeHour, file.MTimeMinute, file.MTimeSecond)
					file.URL = fmt.Sprintf("%s/%d/%d/%s", baseURL, eprintID, i+1, file.Filename)
					files = append(files, file)
				}
				i++
			}
		}
		if len(files) > 0 {
			return files
		}
	}
	return nil
}

func documentIDToRelation(repoID string, documentID int, pos int, db *sql.DB, tables map[string][]string) *ItemList {
	typeTable := "document_relation_type"
	_, okTypeTable := tables[typeTable]
	uriTable := "document_relation_uri"
	_, okUriTable := tables[uriTable]

	if okTypeTable && okUriTable {
		itemList := new(ItemList)
		stmt := fmt.Sprintf(`SELECT document_relation_type.relation_type, document_relation_uri.relation_uri FROM %s JOIN %s ON ((%s.docid = %s.docid) AND (%s.pos = %s.pos)) WHERE (%s.docid = ?)`, typeTable, uriTable, typeTable, uriTable, typeTable, uriTable, typeTable)
		rows, err := db.Query(stmt, documentID)
		if err != nil {
			log.Printf("Query failed %q, doc id %d, pos %d, %s", stmt, documentID, pos, err)
		} else {
			i := 0
			for rows.Next() {
				var (
					relationType, relationURI string
				)
				if err := rows.Scan(&relationType, &relationURI); err != nil {
					log.Printf("Could not scan relation type and relation uri (%d), %q join %q, doc id %d and pos %d, %s", i, typeTable, uriTable, documentID, pos, err)
				} else {
					item := new(Item)
					item.Type = relationType
					item.URI = relationURI
					itemList.AddItem(item)
				}
				i++
			}
			if itemList.Length() > 0 {
				return itemList
			}
		}
	}
	return nil
}

func eprintIDToDocumentList(repoID string, baseURL string, eprintID int, db *sql.DB, tables map[string][]string) *DocumentList {
	tableName := "document"
	columns, ok := tables[tableName]
	if ok {
		documentList := new(DocumentList)
		document := new(Document)
		// NOTE: Bind the values in document to the values array used by
		// rows.Scan().
		columnSQL, values := documentToColumnsAndValues(document, columns)
		stmt := fmt.Sprintf(`SELECT %s FROM %s WHERE eprintid = ? ORDER BY eprintid ASC, pos ASC, rev_number DESC`, strings.Join(columnSQL, ", "), tableName)
		rows, err := db.Query(stmt, eprintID)
		if err != nil {
			log.Printf("Query failed %q for %d in %q, %q,  %s", tableName, eprintID, repoID, stmt, err)
		} else {
			i := 0
			for rows.Next() {
				document = new(Document)
				_, values = documentToColumnsAndValues(document, columns)
				if err := rows.Scan(values...); err != nil {
					log.Printf("Could not scan %q for %d in %q, %s", tableName, eprintID, repoID, err)
				} else {
					document.ID = fmt.Sprintf("%s/id/document/%d", baseURL, document.DocID)
					document.Files = documentIDToFiles(repoID, baseURL, eprintID, document.DocID, db, tables)
					document.Relation = documentIDToRelation(repoID, document.DocID, document.Pos, db, tables)
					documentList.AddDocument(document)
				}
				i++
			}
			rows.Close()
		}

		// FIXME How to Handle document_permission_group, the table appears to be empty?

		if documentList.Length() > 0 {
			// Attach files to documents
			for i := 0; i < documentList.Length(); i++ {
				document := documentList.IndexOf(i)
				files := documentIDToFiles(repoID, baseURL, eprintID, document.DocID, db, tables)
				if (files != nil) && (len(files) > 0) {
					document.Files = files
				}
			}
			return documentList
		}
	}
	return nil
}

func eprintIDToRelatedURL(repoID string, baseURL string, eprintID int, db *sql.DB, tables map[string][]string) *RelatedURLItemList {
	tableName := `eprint_related_url_url`
	_, ok := tables[tableName]
	if ok {
		relatedURLs := new(RelatedURLItemList)
		stmt := `SELECT eprint_related_url_url.pos, IFNULL(related_url_url, '') AS url,
IFNULL(related_url_type, '') AS url_type,
IFNULL(related_url_description, '') AS url_description
FROM eprint_related_url_url LEFT JOIN (eprint_related_url_type, eprint_related_url_description) 
ON ((eprint_related_url_url.eprintid = eprint_related_url_type.eprintid) AND
(eprint_related_url_url.eprintid = eprint_related_url_description.eprintid) AND
(eprint_related_url_url.pos = eprint_related_url_type.pos) AND
(eprint_related_url_url.pos = eprint_related_url_description.pos)
) 
WHERE eprint_related_url_url.eprintid = ?
ORDER BY eprint_related_url_url.eprintid, eprint_related_url_url.pos`
		rows, err := db.Query(stmt, eprintID)
		if err != nil {
			log.Printf("Could not query %d in %q, %s", eprintID, repoID, err)
		} else {
			i := 0
			for rows.Next() {
				var (
					pos                            int
					url, url_type, url_description string
				)
				if err := rows.Scan(&pos, &url, &url_type, &url_description); err != nil {
					log.Printf("Could not scan (%d) %d for %s, %s", i, eprintID, repoID, err)
				} else {
					item := new(Item)
					item.Pos = pos
					item.URL = strings.TrimSpace(url)
					item.Type = strings.TrimSpace(url_type)
					item.Description = strings.TrimSpace(url_description)
					relatedURLs.AddItem(item)
				}
				i++
			}
			rows.Close()
			if len(relatedURLs.Items) > 0 {
				return relatedURLs
			}
		}
	}

	return nil
}

// CrosswalkSQLToEPrint expects a repository map and EPrint ID
// and will generate a series of SELECT statements populating
// a new EPrint struct or return an error
func CrosswalkSQLToEPrint(repoID string, baseURL string, eprintID int) (*EPrint, error) {
	var (
		tables  map[string][]string
		columns []string
	)
	db, ok := config.Connections[repoID]
	if ok == false {
		return nil, fmt.Errorf("not found, %q not known", repoID)
	}
	if eprintID == 0 {
		return nil, fmt.Errorf("not found, %d not in %q", eprintID, repoID)
	}
	_, ok = config.Repositories[repoID]
	if ok == false {
		return nil, fmt.Errorf("not found, %q not defined", repoID)
	}
	tables = config.Repositories[repoID].TableMap
	columns, ok = tables["eprint"]
	if ok == false {
		return nil, fmt.Errorf("not found, %q eprint table not defined", repoID)
	}

	// NOTE: since the specific subset of columns in a repository
	// are known only at run time we need to setup a generic pointer
	// array for the scan results based on our newly allocated
	// EPrint struct.

	eprint := new(EPrint) // Generate an empty EPrint struct

	//NOTE: The data is littered with NULLs in EPrints. We need to
	// generate both a map of values into the EPrint stucture and
	// aggregated the SQL Column definitions to deal with the NULL
	// values.
	columnSQL, values := eprintToColumnsAndValues(eprint, columns)

	// NOTE: With the "values" pointer array setup the query can be built
	// and executed in the usually SQL fashion.
	stmt := fmt.Sprintf("SELECT %s FROM eprint WHERE eprintid = ? LIMIT 1", strings.Join(columnSQL, ", "))
	rows, err := db.Query(stmt, eprintID)
	if err != nil {
		return nil, fmt.Errorf("ERROR: query error (%q, %q), %s", repoID, stmt, err)
	}
	for rows.Next() {
		// NOTE: Because values array holds the addresses into our
		// EPrint struct the "Scan" does the actual mapping.
		// This makes it sorta "auto-magical"
		if err := rows.Scan(values...); err != nil {
			log.Printf("Could not read eprint table, %s", err)
		}
	}
	rows.Close()

	// Normalize fields inferred from MySQL database tables.
	eprint.ID = fmt.Sprintf("%s/id/eprint/%d", baseURL, eprint.EPrintID)
	eprint.LastModified = makeTimestamp(eprint.LastModifiedYear, eprint.LastModifiedMonth, eprint.LastModifiedDay, eprint.LastModifiedHour, eprint.LastModifiedMinute, eprint.LastModifiedSecond)
	eprint.Datestamp = makeTimestamp(eprint.DatestampYear, eprint.DatestampMonth, eprint.DatestampDay, eprint.DatestampHour, eprint.DatestampMinute, eprint.DatestampSecond)
	eprint.StatusChanged = makeTimestamp(eprint.StatusChangedYear, eprint.StatusChangedMonth, eprint.StatusChangedDay, eprint.StatusChangedHour, eprint.StatusChangedMinute, eprint.StatusChangedSecond)
	eprint.Date = makeApproxDate(eprint.DateYear, eprint.DateMonth, eprint.DateDay)

	// CreatorsItemList
	eprint.Creators = creatorsToEPrint(repoID, eprintID, db, tables)
	// EditorsItemList
	eprint.Editors = editorsToEPrint(repoID, eprintID, db, tables)
	// ContributorsItemList
	eprint.Contributors = contributorsToEPrint(repoID, eprintID, db, tables)
	// LocalGroupItemList
	eprint.LocalGroup = localGroupToEPrint(repoID, eprintID, db, tables)
	// FundersItemList
	eprint.Funders = fundersToEPrint(repoID, eprintID, db, tables)
	// Documents (*DocumentList)
	eprint.Documents = eprintIDToDocumentList(repoID, baseURL, eprintID, db, tables)
	// RelatedURLs List
	eprint.RelatedURL = eprintIDToRelatedURL(repoID, baseURL, eprintID, db, tables)

	//FIXME: add various related table data.
	// ReferenceText (item list)
	// Projects
	// OtherNumberingSystem (item list)
	// ErrataItemList (or ErrataText)
	// OtherURL
	// Subjects List
	// ItemIssues
	// CorpCreators
	// CorpContributors
	// Exhibitors
	// Producers
	// Conductors
	// Lyricists
	// Accompaniment
	// SkillAreas
	// CopyrightHolders
	// LearningLevels
	// Reference
	// ConfCreators
	// AltTitle
	// GScholar
	// Shelves
	// PatentAssignee
	// PatentClassification
	// RelatedPatents
	// Divisions
	// ThesisAdvisor
	// ThesisCommittee
	// OptionMajor
	// OptionMinor

	// FIXME: Need to attach id attribute (string of eprintid)

	// FIXME: Need to map some of the variant *Text fields
	// into Item lists to make things consistant. (these are in
	// caltechln and caltechoh.
	//
	// Convert the following
	// - ErrataText to ErrataItemList
	//    - caltechln, caltechoh, caltechconf, caltechcampuspubs, caltechauthors, caltechthesis. Missing in calteches.
	// - LearningLevelText to LearningLevelItemList
	// - PatentClassificationText to PatentClassificationItemList

	return eprint, nil
}

// CrosswalkEPrintToSQL will read an EPrint structure and
// generate SQL INSERT (eprint.ID == 0) or REPLACE statements
// suitable for updating an EPrint record in the repository.
func CrosswalkEPrintToSQL(eprint *EPrint) ([]byte, error) {
	return nil, fmt.Errorf("Not implemented")
}
