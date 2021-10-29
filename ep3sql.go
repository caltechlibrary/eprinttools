package eprinttools

//
// ep3sql.go provides crosswalk methods to/from SQL
//
import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

/*
 * Column mapping for tables.
 */

// colExpr takes a column name, ifNull bool and default value.
// If the "ifNull" bool is true then the form expressed is
// `IFNULL(%s, %s) AS %s` otherwise just the column name
// is returned.
func colExpr(name string, ifNull bool, value string) string {
	if ifNull {
		return fmt.Sprintf(`IFNULL(%s, %s) AS %s`, name, value, name)
	}
	return name
}

// eprintToColumnsAndValues for a given EPrints struct generate a
// list of column names to query along with a recieving values array.
// Return a list of column names (with null handle and aliases) and values.
//
// The bool ifNull will control the type of expression of the column.
func eprintToColumnsAndValues(eprint *EPrint, columnsIn []string, ifNull bool) ([]string, []interface{}) {
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
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "eprint_status":
			values = append(values, &eprint.EPrintStatus)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "userid":
			values = append(values, &eprint.UserID)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "dir":
			values = append(values, &eprint.Dir)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "datestamp_year":
			values = append(values, &eprint.DatestampYear)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "datestamp_month":
			values = append(values, &eprint.DatestampMonth)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "datestamp_day":
			values = append(values, &eprint.DatestampDay)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "datestamp_hour":
			values = append(values, &eprint.DatestampHour)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "datestamp_minute":
			values = append(values, &eprint.DatestampMinute)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "datestamp_second":
			values = append(values, &eprint.DatestampSecond)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "lastmod_year":
			values = append(values, &eprint.LastModifiedYear)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "lastmod_month":
			values = append(values, &eprint.LastModifiedMonth)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "lastmod_day":
			values = append(values, &eprint.LastModifiedDay)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "lastmod_hour":
			values = append(values, &eprint.LastModifiedHour)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "lastmod_minute":
			values = append(values, &eprint.LastModifiedMinute)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "lastmod_second":
			values = append(values, &eprint.LastModifiedSecond)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "status_changed_year":
			values = append(values, &eprint.StatusChangedYear)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "status_changed_month":
			values = append(values, &eprint.StatusChangedMonth)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "status_changed_day":
			values = append(values, &eprint.StatusChangedDay)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "status_changed_hour":
			values = append(values, &eprint.StatusChangedHour)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "status_changed_minute":
			values = append(values, &eprint.StatusChangedMinute)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "status_changed_second":
			values = append(values, &eprint.StatusChangedSecond)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "type":
			values = append(values, &eprint.Type)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "metadata_visibility":
			values = append(values, &eprint.MetadataVisibility)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "title":
			values = append(values, &eprint.Title)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "ispublished":
			values = append(values, &eprint.IsPublished)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "full_text_status":
			values = append(values, &eprint.FullTextStatus)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "keywords":
			values = append(values, &eprint.Keywords)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "note":
			values = append(values, &eprint.Note)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "abstract":
			values = append(values, &eprint.Abstract)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "date_year":
			values = append(values, &eprint.DateYear)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "date_month":
			values = append(values, &eprint.DateMonth)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "date_day":
			values = append(values, &eprint.DateDay)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "date_type":
			values = append(values, &eprint.DateType)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "series":
			values = append(values, &eprint.Series)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "volume":
			values = append(values, &eprint.Volume)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "number":
			values = append(values, &eprint.Number)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "publication":
			values = append(values, &eprint.Publication)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "publisher":
			values = append(values, &eprint.Publisher)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "place_of_pub":
			values = append(values, &eprint.PlaceOfPub)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "edition":
			values = append(values, &eprint.Edition)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "pagerange":
			values = append(values, &eprint.PageRange)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "pages":
			values = append(values, &eprint.Pages)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "event_type":
			values = append(values, &eprint.EventType)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "event_title":
			values = append(values, &eprint.EventTitle)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "event_location":
			values = append(values, &eprint.EventLocation)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "event_dates":
			values = append(values, &eprint.EventDates)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "id_number":
			values = append(values, &eprint.IDNumber)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "refereed":
			values = append(values, &eprint.Refereed)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "isbn":
			values = append(values, &eprint.ISBN)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "issn":
			values = append(values, &eprint.ISSN)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "book_title":
			values = append(values, &eprint.BookTitle)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "official_url":
			values = append(values, &eprint.OfficialURL)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "alt_url":
			values = append(values, &eprint.AltURL)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "rights":
			values = append(values, &eprint.Rights)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "collection":
			values = append(values, &eprint.Collection)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "reviewer":
			values = append(values, &eprint.Reviewer)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "official_cit":
			values = append(values, &eprint.OfficialCitation)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "monograph_type":
			values = append(values, &eprint.MonographType)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "suggestions":
			values = append(values, &eprint.Suggestions)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "pres_type":
			values = append(values, &eprint.PresType)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "succeeds":
			values = append(values, &eprint.Succeeds)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "commentary":
			values = append(values, &eprint.Commentary)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "contact_email":
			values = append(values, &eprint.ContactEMail)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "fileinfo":
			values = append(values, &eprint.FileInfo)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "latitude":
			values = append(values, &eprint.Latitude)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0.0`))
		case "longitude":
			values = append(values, &eprint.Longitude)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0.0`))
		case "department":
			values = append(values, &eprint.Department)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "output_media":
			values = append(values, &eprint.OutputMedia)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "num_pieces":
			values = append(values, &eprint.NumPieces)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "composition_type":
			values = append(values, &eprint.CompositionType)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "data_type":
			values = append(values, &eprint.DataType)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "pedagogic_type":
			values = append(values, new(string))
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "learning_level":
			values = append(values, &eprint.LearningLevelText)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "completion_time":
			values = append(values, &eprint.CompletionTime)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "task_purpose":
			values = append(values, &eprint.TaskPurpose)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "doi":
			values = append(values, &eprint.DOI)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "pmc_id":
			values = append(values, &eprint.PMCID)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "pmid":
			values = append(values, &eprint.PMID)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "parent_url":
			values = append(values, &eprint.ParentURL)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "toc":
			values = append(values, &eprint.TOC)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "interviewer":
			values = append(values, &eprint.Interviewer)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "interviewdate":
			values = append(values, &eprint.InterviewDate)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "nonsubj_keywords":
			values = append(values, &eprint.NonSubjKeywords)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "season":
			values = append(values, &eprint.Season)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "classification_code":
			values = append(values, &eprint.ClassificationCode)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "sword_depositor":
			values = append(values, &eprint.SwordDepositor)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "sword_depository":
			values = append(values, &eprint.SwordDepository)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "sword_slug":
			values = append(values, &eprint.SwordSlug)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "importid":
			values = append(values, &eprint.ImportID)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "patent_applicant":
			values = append(values, &eprint.PatentApplicant)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "patent_number":
			values = append(values, &eprint.PatentNumber)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "institution":
			values = append(values, &eprint.Institution)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "thesis_type":
			values = append(values, &eprint.ThesisType)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "thesis_degree":
			values = append(values, &eprint.ThesisDegree)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "thesis_degree_grantor":
			values = append(values, &eprint.ThesisDegreeGrantor)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "thesis_degree_date_year":
			values = append(values, &eprint.ThesisDegreeDateYear)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_degree_date_month":
			values = append(values, &eprint.ThesisDegreeDateMonth)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_degree_date_day":
			values = append(values, &eprint.ThesisDegreeDateDay)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_submitted_date_year":
			values = append(values, &eprint.ThesisSubmittedDateYear)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_submitted_date_month":
			values = append(values, &eprint.ThesisSubmittedDateMonth)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_submitted_date_day":
			values = append(values, &eprint.ThesisSubmittedDateDay)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_defense_date":
			values = append(values, &eprint.ThesisDefenseDate)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "thesis_defense_date_year":
			values = append(values, &eprint.ThesisDefenseDateYear)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_defense_date_month":
			values = append(values, &eprint.ThesisDefenseDateMonth)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_defense_date_day":
			values = append(values, &eprint.ThesisDefenseDateDay)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_approved_date_year":
			values = append(values, &eprint.ThesisApprovedDateYear)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_approved_date_month":
			values = append(values, &eprint.ThesisApprovedDateMonth)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_approved_date_day":
			values = append(values, &eprint.ThesisApprovedDateDay)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_public_date_year":
			values = append(values, &eprint.ThesisPublicDateYear)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_public_date_month":
			values = append(values, &eprint.ThesisPublicDateMonth)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_public_date_day":
			values = append(values, &eprint.ThesisPublicDateDay)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_author_email":
			values = append(values, &eprint.ThesisAuthorEMail)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "hide_thesis_author_email":
			values = append(values, &eprint.HideThesisAuthorEMail)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "gradofc_approval_date":
			values = append(values, &eprint.GradOfficeApprovalDate)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "gradofc_approval_date_year":
			values = append(values, &eprint.GradOfficeApprovalDateYear)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "gradofc_approval_date_month":
			values = append(values, &eprint.GradOfficeApprovalDateMonth)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "gradofc_approval_date_day":
			values = append(values, &eprint.GradOfficeApprovalDateDay)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_awards":
			values = append(values, &eprint.ThesisAwards)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "review_status":
			values = append(values, &eprint.ReviewStatus)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "copyright_statement":
			values = append(values, &eprint.CopyrightStatement)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "source":
			values = append(values, &eprint.Source)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "replacedby":
			values = append(values, &eprint.ReplacedBy)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "item_issues_count":
			values = append(values, &eprint.ItemIssuesCount)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "errata":
			values = append(values, &eprint.ErrataText)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "coverage_dates":
			values = append(values, &eprint.CoverageDates)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "edit_lock_user":
			values = append(values, &eprint.EditLockUser)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "edit_lock_since":
			values = append(values, &eprint.EditLockSince)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "edit_lock_until":
			values = append(values, &eprint.EditLockUntil)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		// The follow values represent sub tables and processed separately.
		case "patent_classification":
			values = append(values, &eprint.PatentClassificationText)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		default:
			// Handle case where we have value that is unmapped or not available in EPrint struct
			log.Printf("could not map %q (%d) into EPrint struct", key, i)
		}

	}
	return columnsOut, values
}

func documentToColumnsAndValues(document *Document, columns []string, ifNull bool) ([]string, []interface{}) {
	columnsOut := []string{}
	values := []interface{}{}
	for i, key := range columns {
		switch key {
		case "docid":
			values = append(values, &document.DocID)
			columnsOut = append(columnsOut, key)
		case "eprintid":
			values = append(values, &document.EPrintID)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "pos":
			values = append(values, &document.Pos)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "rev_number":
			values = append(values, &document.RevNumber)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "format":
			values = append(values, &document.Format)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "formatdesc":
			values = append(values, &document.FormatDesc)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "language":
			values = append(values, &document.Language)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "security":
			values = append(values, &document.Security)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "license":
			values = append(values, &document.License)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "main":
			values = append(values, &document.Main)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "date_embargo_year":
			values = append(values, &document.DateEmbargoYear)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "date_embargo_month":
			values = append(values, &document.DateEmbargoMonth)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "date_embargo_day":
			values = append(values, &document.DateEmbargoDay)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "content":
			values = append(values, &document.Content)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "placement":
			values = append(values, &document.Placement)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "mime_type":
			values = append(values, &document.MimeType)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "media_duration":
			values = append(values, &document.MediaDuration)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "media_audio_codec":
			values = append(values, &document.MediaAudioCodec)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "media_video_codec":
			values = append(values, &document.MediaVideoCodec)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "media_width":
			values = append(values, &document.MediaWidth)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "media_height":
			values = append(values, &document.MediaHeight)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "media_aspect_ratio":
			values = append(values, &document.MediaAspectRatio)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "media_sample_start":
			values = append(values, &document.MediaSampleStart)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "media_sample_stop":
			values = append(values, &document.MediaSampleStop)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		default:
			log.Printf("%q (%d) not found in document table", key, i)
		}
	}
	return columnsOut, values
}

func fileToColumnsAndValues(file *File, columns []string, ifNull bool) ([]string, []interface{}) {
	columnsOut := []string{}
	values := []interface{}{}
	for i, key := range columns {
		switch key {
		case "fileid":
			values = append(values, &file.FileID)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "datasetid":
			values = append(values, &file.DatasetID)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "objectid":
			values = append(values, &file.ObjectID)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "filename":
			values = append(values, &file.Filename)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "mime_type":
			values = append(values, &file.MimeType)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "hash":
			values = append(values, &file.Hash)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "hash_type":
			values = append(values, &file.HashType)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "filesize":
			values = append(values, &file.FileSize)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "mtime_year":
			values = append(values, &file.MTimeYear)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "mtime_month":
			values = append(values, &file.MTimeMonth)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "mtime_day":
			values = append(values, &file.MTimeDay)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "mtime_hour":
			values = append(values, &file.MTimeHour)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "mtime_minute":
			values = append(values, &file.MTimeMinute)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "mtime_second":
			values = append(values, &file.MTimeSecond)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		default:
			log.Printf("%q (%d) not found in file table", key, i)
		}
	}
	return columnsOut, values
}

/*
 * Document and files models
 */
func documentIDToFiles(repoID string, baseURL string, eprintID int, documentID int, pos int, db *sql.DB, tables map[string][]string) []*File {
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
		columnSQL, values := fileToColumnsAndValues(file, columns, true)
		//FIXME: This needs to be an ordered list.
		stmt := fmt.Sprintf(`SELECT %s FROM %s WHERE datasetid = 'document' AND objectid = ?`, strings.Join(columnSQL, ", "), tableName)
		rows, err := db.Query(stmt, documentID)
		if err != nil {
			log.Printf("Query failed %q for %d in %q, doc ID %d, , %q,  %s", tableName, eprintID, repoID, documentID, stmt, err)
		} else {
			i := 0
			for rows.Next() {
				file = new(File)
				_, values = fileToColumnsAndValues(file, columns, true)
				if err := rows.Scan(values...); err != nil {
					log.Printf("Could not scan %q for %d in %q, doc ID %d, %s", tableName, eprintID, repoID, documentID, err)
				} else {
					file.ID = fmt.Sprintf("%s/id/file/%d", baseURL, file.FileID)
					file.MTime = makeTimestamp(file.MTimeYear, file.MTimeMonth, file.MTimeDay, file.MTimeHour, file.MTimeMinute, file.MTimeSecond)
					file.URL = fmt.Sprintf("%s/%d/%d/%s", baseURL, eprintID, pos, file.Filename)
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

func documentIDToRelation(repoID string, baseURL string, documentID int, pos int, db *sql.DB, tables map[string][]string) *ItemList {
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
					item.URI = fmt.Sprintf(`%s%s`, baseURL, relationURI)
					itemList.Append(item)
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
		columnSQL, values := documentToColumnsAndValues(document, columns, true)
		stmt := fmt.Sprintf(`SELECT %s FROM %s WHERE eprintid = ? ORDER BY eprintid ASC, pos ASC, rev_number DESC`, strings.Join(columnSQL, ", "), tableName)
		rows, err := db.Query(stmt, eprintID)
		if err != nil {
			log.Printf("Query failed %q for %d in %q, %q,  %s", tableName, eprintID, repoID, stmt, err)
		} else {
			i := 0
			for rows.Next() {
				document = new(Document)
				_, values = documentToColumnsAndValues(document, columns, true)
				if err := rows.Scan(values...); err != nil {
					log.Printf("Could not scan %q for %d in %q, %s", tableName, eprintID, repoID, err)
				} else {
					document.ID = fmt.Sprintf("%s/id/document/%d", baseURL, document.DocID)
					document.Files = documentIDToFiles(repoID, baseURL, eprintID, document.DocID, document.Pos, db, tables)
					document.Relation = documentIDToRelation(repoID, baseURL, document.DocID, document.Pos, db, tables)
					documentList.Append(document)
				}
				i++
			}
			rows.Close()
		}

		// NOTE: The document_permission_group the table is empty in our repositories
		if documentList.Length() > 0 {
			// Attach files to documents
			for i := 0; i < documentList.Length(); i++ {
				document := documentList.IndexOf(i)
				files := documentIDToFiles(repoID, baseURL, eprintID, document.DocID, document.Pos, db, tables)
				if (files != nil) && (len(files) > 0) {
					document.Files = files
				}
			}
			return documentList
		}
	}
	return nil
}

/*
 * Common models and help functions
 */

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

// NOTE: Do to the funky way MySQL and EPrints works with UTF
// I can't use the time package to build my formatted date string.
// An odd edge case can make the year, month or day off by one.

func makeTimestamp(year int, month int, day int, hour int, minute int, second int) string {
	if (year > 0) && (month > 0) && (day > 0) {
		return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, minute, second)
	}
	return ""
}

func makeDatestamp(year int, month int, day int) string {
	if (year > 0) && (month > 0) && (day > 0) {
		return fmt.Sprintf("%d-%02d-%02d", year, month, day)
	}
	return ""
}

func makeApproxDate(year int, month int, day int) string {
	switch {
	case (year > 0) && (month > 0) && (day > 0):
		return fmt.Sprintf("%d-%02d-%02d", year, month, day)
	case (year > 0) && (month > 0) && (day == 0):
		return fmt.Sprintf("%d-%02d", year, month)
	case (year > 0) && (month == 0) && (day == 0):
		return fmt.Sprintf("%d", year)
	default:
		return ""
	}
}

/*
 * PersonItemList model
 */
func eprintIDToPersonItemList(db *sql.DB, tables map[string][]string, repoID string, eprintID int, tablePrefix string, itemList ItemsInterface) int {
	var (
		pos                                       int
		value, honourific, given, family, lineage string
	)
	tableName := tablePrefix + "_name"
	_, ok := tables[tableName]
	if ok {
		columnPrefix := strings.TrimPrefix(tableName, `eprint_`)
		stmt := fmt.Sprintf(`SELECT pos, IFNULL(%s_honourific, '') AS honourific, IFNULL(%s_given, '') AS given, IFNULL(%s_family, '') AS family, IFNULL(%s_lineage, '') AS lineage FROM %s WHERE eprintid = ? ORDER BY eprintid, pos`, columnPrefix, columnPrefix, columnPrefix, columnPrefix, tableName)
		rows, err := db.Query(stmt, eprintID)
		if err != nil {
			log.Printf("Could not query %s for %d in %q, %s", tableName, eprintID, repoID, err)
		} else {
			i := 0
			for rows.Next() {
				if err := rows.Scan(&pos, &honourific, &given, &family, &lineage); err != nil {
					log.Printf("Could not scan %s for %d in %q, %s", tableName, eprintID, repoID, err)
				} else {
					item := new(Item)
					item.Pos = pos
					item.Name = makePersonName(given, family, honourific, lineage)
					itemList.Append(item)
				}
				i++
			}
			rows.Close()

			if itemList.Length() > 0 {
				tablesAndColumn := map[string][2]string{}
				columnPrefix := strings.TrimPrefix(tablePrefix, `eprint_`)
				for _, suffix := range []string{"id", "orcid", "uri", "role", "email", "show_email"} {
					key := fmt.Sprintf("%s_%s", tablePrefix, suffix)
					tablesAndColumn[key] = [2]string{
						// Column Name
						fmt.Sprintf("%s_%s", columnPrefix, suffix),
						// Column Alias
						suffix,
					}
				}
				for tableName, columnNames := range tablesAndColumn {
					columnName, columnAlias := columnNames[0], columnNames[1]
					_, ok := tables[tableName]
					if ok {
						stmt := fmt.Sprintf(`SELECT pos, IFNULL(%s, "") AS %s FROM %s WHERE eprintid = ? ORDER BY eprintid, pos`, columnName, columnAlias, tableName)
						rows, err = db.Query(stmt, eprintID)
						if err != nil {
							log.Printf("Could not query (%d) %s for %d in %q, %s", i, tableName, eprintID, repoID, err)
						} else {
							i := 0
							for rows.Next() {
								if err := rows.Scan(&pos, &value); err != nil {
									log.Printf("Could not scan (%d) %s for %d in %q, %s", i, tableName, eprintID, repoID, err)
								} else {
									itemList.SetAttributeOf(pos, columnAlias, value)
								}
								i++
							}
							rows.Close()
						}
					}
				}
			}
		}
	}
	return itemList.Length()
}

func eprintIDToCreators(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *CreatorItemList {
	tablePrefix := `eprint_creators`
	itemList := new(CreatorItemList)
	if count := eprintIDToPersonItemList(db, tables, repoID, eprintID, tablePrefix, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToEditors(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *EditorItemList {
	tablePrefix := `eprint_editors`
	itemList := new(EditorItemList)
	if count := eprintIDToPersonItemList(db, tables, repoID, eprintID, tablePrefix, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToContributors(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *ContributorItemList {
	tablePrefix := `eprint_constibutors`
	itemList := new(ContributorItemList)
	if count := eprintIDToPersonItemList(db, tables, repoID, eprintID, tablePrefix, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToExhibitors(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *ExhibitorItemList {
	tablePrefix := `eprint_exhibitors`
	itemList := new(ExhibitorItemList)
	if count := eprintIDToPersonItemList(db, tables, repoID, eprintID, tablePrefix, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToProducers(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *ProducerItemList {
	tablePrefix := `eprint_producers`
	itemList := new(ProducerItemList)
	if count := eprintIDToPersonItemList(db, tables, repoID, eprintID, tablePrefix, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToConductors(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *ConductorItemList {
	tablePrefix := `eprint_conductors`
	itemList := new(ConductorItemList)
	if count := eprintIDToPersonItemList(db, tables, repoID, eprintID, tablePrefix, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToLyricists(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *LyricistItemList {
	tablePrefix := `eprint_lyricists_name`
	itemList := new(LyricistItemList)
	if count := eprintIDToPersonItemList(db, tables, repoID, eprintID, tablePrefix, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToThesisAdvisors(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *ThesisAdvisorItemList {
	tablePrefix := `eprint_thesis_advisor`
	itemList := new(ThesisAdvisorItemList)
	if count := eprintIDToPersonItemList(db, tables, repoID, eprintID, tablePrefix, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToThesisCommittee(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *ThesisCommitteeItemList {
	tablePrefix := `eprint_thesis_committee`
	itemList := new(ThesisCommitteeItemList)
	if count := eprintIDToPersonItemList(db, tables, repoID, eprintID, tablePrefix, itemList); count > 0 {
		return itemList
	}
	return nil
}

/*
 * SimpleItemList model
 */

func eprintIDToSimpleItemList(db *sql.DB, tables map[string][]string, repoID string, eprintID int, tableName string, itemList ItemsInterface) int {
	columnName := strings.TrimPrefix(tableName, `eprint_`)
	var (
		pos   int
		value string
	)
	_, ok := tables[tableName]
	if ok {
		stmt := fmt.Sprintf(`SELECT pos, IFNULL(%s, '') AS %s FROM %s WHERE eprintid = ? ORDER BY eprintid, pos`, columnName, columnName, tableName)
		rows, err := db.Query(stmt, eprintID)
		if err != nil {
			log.Printf("Could not query %s for %d in %q, %s", tableName, eprintID, repoID, err)
		} else {
			i := 0
			for rows.Next() {
				if err := rows.Scan(&pos, &value); err != nil {
					log.Printf("Could not scan %s for %d in %q, %s", tableName, eprintID, repoID, err)
				} else {
					item := new(Item)
					item.Pos = pos
					if value != "" {
						item.Value = value
					}
					itemList.Append(item)
				}
				i++
			}
			rows.Close()
		}
	}
	return itemList.Length()
}

func eprintIDToLocalGroup(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *LocalGroupItemList {
	tableName := `eprint_local_group`
	itemList := new(LocalGroupItemList)
	if count := eprintIDToSimpleItemList(db, tables, repoID, eprintID, tableName, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToReferenceText(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *ReferenceTextItemList {
	tableName := `eprint_referencetext`
	itemList := new(ReferenceTextItemList)
	if count := eprintIDToSimpleItemList(db, tables, repoID, eprintID, tableName, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToProjects(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *ProjectItemList {
	tableName := `eprint_projects`
	itemList := new(ProjectItemList)
	if count := eprintIDToSimpleItemList(db, tables, repoID, eprintID, tableName, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToSubjects(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *SubjectItemList {
	tableName := `eprint_subjects`
	itemList := new(SubjectItemList)
	if count := eprintIDToSimpleItemList(db, tables, repoID, eprintID, tableName, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToAccompaniment(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *AccompanimentItemList {
	tableName := `eprint_accompaniment`
	itemList := new(AccompanimentItemList)
	if count := eprintIDToSimpleItemList(db, tables, repoID, eprintID, tableName, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToSkillAreas(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *SkillAreaItemList {
	tableName := `eprint_skill_area`
	itemList := new(SkillAreaItemList)
	if count := eprintIDToSimpleItemList(db, tables, repoID, eprintID, tableName, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToCopyrightHolders(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *CopyrightHolderItemList {
	tableName := `eprint_copyright_holders`
	itemList := new(CopyrightHolderItemList)
	if count := eprintIDToSimpleItemList(db, tables, repoID, eprintID, tableName, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToReference(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *ReferenceItemList {
	tableName := `eprint_reference`
	itemList := new(ReferenceItemList)
	if count := eprintIDToSimpleItemList(db, tables, repoID, eprintID, tableName, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToAltTitle(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *AltTitleItemList {
	tableName := `eprint_alt_title`
	itemList := new(AltTitleItemList)
	if count := eprintIDToSimpleItemList(db, tables, repoID, eprintID, tableName, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToPatentAssignee(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *PatentAssigneeItemList {
	tableName := `eprint_patent_assignee`
	itemList := new(PatentAssigneeItemList)
	if count := eprintIDToSimpleItemList(db, tables, repoID, eprintID, tableName, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToRelatedPatents(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *RelatedPatentItemList {
	tableName := `eprint_related_patents`
	itemList := new(RelatedPatentItemList)
	if count := eprintIDToSimpleItemList(db, tables, repoID, eprintID, tableName, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToDivisions(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *DivisionItemList {
	tableName := `eprint_divisions`
	itemList := new(DivisionItemList)
	if count := eprintIDToSimpleItemList(db, tables, repoID, eprintID, tableName, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToOptionMajor(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *OptionMajorItemList {
	tableName := `eprint_option_major`
	itemList := new(OptionMajorItemList)
	if count := eprintIDToSimpleItemList(db, tables, repoID, eprintID, tableName, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToOptionMinor(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *OptionMinorItemList {
	tableName := `eprint_option_minor`
	itemList := new(OptionMinorItemList)
	if count := eprintIDToSimpleItemList(db, tables, repoID, eprintID, tableName, itemList); count > 0 {
		return itemList
	}
	return nil
}

/*
 * Hetrogenous models
 */

func eprintIDToConfCreators(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *ConfCreatorItemList {
	var (
		pos   int
		value string
	)
	tableName := `eprint_conf_creators_name`
	columnName := `conf_creators_name`
	_, ok := tables[tableName]
	if ok {
		itemList := new(ConfCreatorItemList)
		stmt := fmt.Sprintf(`SELECT pos, IFNULL(%s, '') AS %s
FROM %s WHERE eprintid = ? ORDER BY eprintid, pos`, columnName, columnName, tableName)
		rows, err := db.Query(stmt, eprintID)
		if err != nil {
			log.Printf("Could not query %s for %d in %q, %s", tableName, eprintID, repoID, err)
		} else {
			i := 0
			for rows.Next() {
				if err := rows.Scan(&pos, &value); err != nil {
					log.Printf("Could not scan %s for %d in %q, %s", tableName, eprintID, repoID, err)
				} else {
					item := new(Item)
					item.Pos = pos
					item.Name = new(Name)
					item.Name.Value = value
					itemList.Append(item)
				}
				i++
			}
			rows.Close()

			if itemList.Length() > 0 {
				tablesAndColumn := map[string]string{
					"eprint_conf_creators_id":  "conf_creators_id",
					"eprint_conf_creators_ror": "conf_creators_ror",
					"eprint_conf_creators_uri": "conf_creators_uri",
					"eprint_conf_creators":     "conf_creators",
				}
				for tableName, columnName := range tablesAndColumn {
					_, ok := tables[tableName]
					if ok {
						stmt := fmt.Sprintf(`SELECT pos, IFNULL(%s, "") AS %s FROM %s WHERE eprintid = ? ORDER BY eprintid, pos`, columnName, columnName, tableName)
						rows, err = db.Query(stmt, eprintID)
						if err != nil {
							log.Printf("Could not query (%d) %s for %d in %q, %s", i, tableName, eprintID, repoID, err)
						} else {
							i := 0
							for rows.Next() {
								if err := rows.Scan(&pos, &value); err != nil {
									log.Printf("Could not scan (%d) %s for %d in %q, %s", i, tableName, eprintID, repoID, err)
								} else {
									for _, item := range itemList.Items {
										if item.Pos == pos && value != "" {
											switch columnName {
											case "conf_creators_id":
												item.ID = value
											case "conf_creators_ror":
												item.ROR = value
											case "conf_creators_uri":
												item.URI = value
											case "conf_creators":
												item.Name = new(Name)
												item.Name.Value = value
											}
											break
										}
									}
								}
								i++
							}
							rows.Close()
						}
					}
					return itemList
				}
			}
		}
	}
	return nil
}

func eprintIDToCorpCreators(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *CorpCreatorItemList {
	var (
		pos   int
		value string
	)
	tableName := `eprint_corp_creators_name`
	columnName := `corp_creators_name`
	_, ok := tables[tableName]
	if ok {
		itemList := new(CorpCreatorItemList)
		stmt := fmt.Sprintf(`SELECT pos, IFNULL(%s, '') AS %s
FROM %s WHERE eprintid = ? ORDER BY eprintid, pos`, columnName, columnName, tableName)
		rows, err := db.Query(stmt, eprintID)
		if err != nil {
			log.Printf("Could not query %s for %d in %q, %s", tableName, eprintID, repoID, err)
		} else {
			i := 0
			for rows.Next() {
				if err := rows.Scan(&pos, &value); err != nil {
					log.Printf("Could not scan %s for %d in %q, %s", tableName, eprintID, repoID, err)
				} else {
					item := new(Item)
					item.Pos = pos
					item.Name = new(Name)
					item.Name.Value = value
					itemList.Append(item)
				}
				i++
			}
			rows.Close()

			if itemList.Length() > 0 {
				tablesAndColumn := map[string]string{
					"eprint_corp_creators_id":  "corp_creators_id",
					"eprint_corp_creators_ror": "corp_creators_ror",
					"eprint_corp_creators_uri": "corp_creators_uri",
					"eprint_corp_creators":     "corp_creators",
				}
				for tableName, columnName := range tablesAndColumn {
					_, ok := tables[tableName]
					if ok {
						stmt := fmt.Sprintf(`SELECT pos, IFNULL(%s, "") AS %s FROM %s WHERE eprintid = ? ORDER BY eprintid, pos`, columnName, columnName, tableName)
						rows, err = db.Query(stmt, eprintID)
						if err != nil {
							log.Printf("Could not query (%d) %s for %d in %q, %s", i, tableName, eprintID, repoID, err)
						} else {
							i := 0
							for rows.Next() {
								if err := rows.Scan(&pos, &value); err != nil {
									log.Printf("Could not scan (%d) %s for %d in %q, %s", i, tableName, eprintID, repoID, err)
								} else {
									for _, item := range itemList.Items {
										if item.Pos == pos && value != "" {
											switch columnName {
											case "corp_creators_id":
												item.ID = value
											case "corp_creators_ror":
												item.ROR = value
											case "corp_creators_uri":
												item.URI = value
											case "corp_creators":
												item.Name = new(Name)
												item.Name.Value = value
											}
											break
										}
									}
								}
								i++
							}
							rows.Close()
						}
					}
					return itemList
				}
			}
		}
	}
	return nil
}

func eprintIDToFunders(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *FunderItemList {
	var (
		pos   int
		value string
	)
	tableName := `eprint_funders_agency`
	columnName := `funders_agency`
	_, ok := tables[tableName]
	if ok {
		// eprint_%_id is a known structure. eprintid, pos, contributors_id
		itemList := new(FunderItemList)
		stmt := fmt.Sprintf(`SELECT pos, IFNULL(%s, '') AS %s FROM %s WHERE eprintid = ? ORDER BY eprintid, pos`, columnName, columnName, tableName)
		rows, err := db.Query(stmt, eprintID)
		if err != nil {
			log.Printf("Could not query %s for %d in %q, %s", tableName, eprintID, repoID, err)
		} else {
			i := 0
			for rows.Next() {
				if err := rows.Scan(&pos, &value); err != nil {
					log.Printf("Could not scan %s for %d in %q, %s", tableName, eprintID, repoID, err)
				} else {
					item := new(Item)
					item.Pos = pos
					if value != "" {
						item.Agency = value
					}
					itemList.Append(item)
				}
				i++
			}
			rows.Close()
		}
		tablesAndColumns := map[string]string{
			"eprint_funders_grant_number": "funders_grant_number",
			"eprint_funders_ror":          "funders_ror",
		}
		if itemList.Length() > 0 {
			for tableName, columnName := range tablesAndColumns {
				if _, ok := tables[tableName]; ok {
					stmt := fmt.Sprintf(`SELECT pos, IFNULL(%s, '') AS %s FROM %s WHERE eprintid = ? ORDER BY eprintid, pos`, columnName, columnName, tableName)
					rows, err = db.Query(stmt, eprintID)
					if err != nil {
						log.Printf("Could not query %s for %d in %q, %s", tableName, eprintID, repoID, err)
					} else {
						i := 0
						for rows.Next() {
							if err := rows.Scan(&pos, &value); err != nil {
								log.Printf("Could not scan (%d) %s for %d in %q, %s", i, tableName, eprintID, repoID, err)
							} else {
								if value != "" {
									for _, item := range itemList.Items {
										if item.Pos == pos {
											switch columnName {
											case "funders_grant_number":
												item.GrantNumber = value
											case "funders_ror":
												item.ROR = value
											}
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
			}
			return itemList
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
					relatedURLs.Append(item)
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

func eprintIDToOtherNumberingSystem(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *OtherNumberingSystemItemList {
	tableNames := []string{`eprint_other_numbering_system_name`, `eprint_other_numbering_system_id`}
	ok := true
	for _, tableName := range tableNames {
		if _, hasTable := tables[tableName]; hasTable == false {
			ok = false
			break
		}
	}
	if ok {
		itemList := new(OtherNumberingSystemItemList)
		stmt := fmt.Sprintf(`
SELECT %s.pos AS pos, IFNULL(other_numbering_system_name, '') AS name, IFNULL(other_numbering_system_id, '') AS systemid FROM %s JOIN %s ON (%s.eprintid = %s.eprintid AND %s.pos = %s.pos) WHERE %s.eprintid = ?`, tableNames[0], tableNames[0], tableNames[1], tableNames[0], tableNames[1], tableNames[0], tableNames[1], tableNames[0])
		rows, err := db.Query(stmt, eprintID)
		if err != nil {
			log.Printf("Could not query %d in %q, %s", eprintID, repoID, err)
		} else {
			var (
				pos      int
				name, id string
			)
			i := 0
			for rows.Next() {
				if err := rows.Scan(&pos, &name, &id); err != nil {
					log.Printf("Could not scan (%d) %d for %s, %s", i, eprintID, repoID, err)
				} else {
					item := new(Item)
					item.Pos = pos
					item.ID = id
					if name != "" {
						item.Name = new(Name)
						item.Name.Value = name
					}
					itemList.Append(item)
				}
				i++
			}
			rows.Close()
			if itemList.Length() > 0 {
				return itemList
			}
		}
	}
	return nil
}

func eprintIDToItemIssues(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *ItemIssueItemList {
	tableName := `eprint_item_issues_timestamp`
	_, ok := tables[tableName]
	if ok {
		var (
			year, month, day, hour, minute, second, pos int
			value                                       string
		)
		itemList := new(ItemIssueItemList)
		stmt := fmt.Sprintf(`SELECT pos, 
IFNULL(item_issues_timestamp_year, 0) AS year, 
IFNULL(item_issues_timestamp_month, 0) AS month,
IFNULL(item_issues_timestamp_day, 0) AS day,
IFNULL(item_issues_timestamp_hour, 0) AS hour,
IFNULL(item_issues_timestamp_minute, 0) AS minute,
IFNULL(item_issues_timestamp_second, 0) AS second
FROM %s WHERE eprintid = ? ORDER BY eprintid, pos`, tableName)
		rows, err := db.Query(stmt, eprintID)
		if err != nil {
			log.Printf("Could not query %d in %q, %s", eprintID, repoID, err)
		} else {
			i := 0
			for rows.Next() {
				if err := rows.Scan(&pos, &year, &month, &day, &hour, &minute, &second); err != nil {
					log.Printf("Could not scan (%d) %d for %s, %s", i, eprintID, repoID, err)
				} else {
					item := new(Item)
					item.Pos = pos
					item.Timestamp = makeTimestamp(year, month, day, hour, minute, second)
					itemList.Append(item)
				}
				i++
			}
			rows.Close()
		}
		if itemList.Length() > 0 {
			tablesAndColumn := map[string]string{
				"eprint_item_issues_type":        "item_issues_type",
				"eprint_item_issues_status":      "item_issues_status",
				"eprint_item_issues_description": "item_issues_description",
				"eprint_item_issues_id":          "item_issues_id",
				"eprint_item_issues_resolved_by": "item_issues_resolved_by",
				"eprint_item_issues_reported_by": "item_issues_reported_by",

				"eprint_item_issues_comment": "item_issues_comment",
			}
			for tableName, columnName := range tablesAndColumn {
				if _, ok := tables[tableName]; ok {
					stmt = fmt.Sprintf(`SELECT pos, IFNULL(%s, '') AS %s FROM %s WHERE eprintid = ? ORDER BY eprintid, pos`, columnName, columnName, tableName)
					rows, err := db.Query(stmt, eprintID)
					if err != nil {
						log.Printf("Could not query %d in %q, %s", eprintID, repoID, err)
					} else {
						i := 0
						for rows.Next() {
							if err := rows.Scan(&pos, &value); err != nil {
								log.Printf("Could not scan (%d) %d for %s, %s", i, eprintID, repoID, err)
							} else {
								for _, item := range itemList.Items {
									if item.Pos == pos && strings.TrimSpace(value) != "" {
										switch strings.TrimPrefix(columnName, "item_issues_") {
										case "type":
											item.Type = value
										case "status":
											item.Status = value
										case "description":
											item.Description = value
										case "id":
											item.ID = value
										case "resolved_by":
											item.ResolvedBy = value
										case "reported_by":
											item.ReportedBy = value
										case "comment":
											item.Comment = value
										}
									}
								}
							}
							i++
						}
						rows.Close()
					}
				}
			}
			return itemList
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
	columnSQL, values := eprintToColumnsAndValues(eprint, columns, true)

	// NOTE: With the "values" pointer array setup the query can be built
	// and executed in the usually SQL fashion.
	stmt := fmt.Sprintf(`SELECT %s FROM eprint WHERE eprintid = ? LIMIT 1`, strings.Join(columnSQL, `, `))
	rows, err := db.Query(stmt, eprintID)
	if err != nil {
		return nil, fmt.Errorf(`ERROR: query error (%q, %q), %s`, repoID, stmt, err)
	}
	for rows.Next() {
		// NOTE: Because values array holds the addresses into our
		// EPrint struct the "Scan" does the actual mapping.
		// This makes it sorta "auto-magical"
		if err := rows.Scan(values...); err != nil {
			log.Printf(`Could not read eprint table, %s`, err)
		}
	}
	rows.Close()

	// Normalize fields inferred from MySQL database tables.
	eprint.ID = fmt.Sprintf(`%s/id/eprint/%d`, baseURL, eprint.EPrintID)
	eprint.LastModified = makeTimestamp(eprint.LastModifiedYear, eprint.LastModifiedMonth, eprint.LastModifiedDay, eprint.LastModifiedHour, eprint.LastModifiedMinute, eprint.LastModifiedSecond)
	// NOTE: EPrint XML uses a datestamp for output but tracks a timestamp.
	//eprint.Datestamp = makeTimestamp(eprint.DatestampYear, eprint.DatestampMonth, eprint.DatestampDay, eprint.DatestampHour, eprint.DatestampMinute, eprint.DatestampSecond)
	eprint.Datestamp = makeDatestamp(eprint.DatestampYear, eprint.DatestampMonth, eprint.DatestampDay)
	eprint.StatusChanged = makeTimestamp(eprint.StatusChangedYear, eprint.StatusChangedMonth, eprint.StatusChangedDay, eprint.StatusChangedHour, eprint.StatusChangedMinute, eprint.StatusChangedSecond)
	eprint.Date = makeApproxDate(eprint.DateYear, eprint.DateMonth, eprint.DateDay)

	// Used in CaltechTHESIS
	eprint.ThesisSubmittedDate = makeDatestamp(eprint.ThesisSubmittedDateYear, eprint.ThesisSubmittedDateMonth, eprint.ThesisSubmittedDateDay)
	eprint.ThesisDefenseDate = makeDatestamp(eprint.ThesisDefenseDateYear, eprint.ThesisDefenseDateMonth, eprint.ThesisDefenseDateDay)
	eprint.ThesisApprovedDate = makeDatestamp(eprint.ThesisApprovedDateYear, eprint.ThesisApprovedDateMonth, eprint.ThesisApprovedDateDay)
	eprint.ThesisPublicDate = makeDatestamp(eprint.ThesisPublicDateYear, eprint.ThesisPublicDateMonth, eprint.ThesisPublicDateDay)
	eprint.ThesisDegreeDate = makeDatestamp(eprint.ThesisDegreeDateYear, eprint.ThesisDegreeDateMonth, eprint.ThesisDegreeDateDay)
	eprint.GradOfficeApprovalDate = makeDatestamp(eprint.GradOfficeApprovalDateYear, eprint.GradOfficeApprovalDateMonth, eprint.GradOfficeApprovalDateDay)

	// CreatorsItemList
	eprint.Creators = eprintIDToCreators(repoID, eprintID, db, tables)
	// EditorsItemList
	eprint.Editors = eprintIDToEditors(repoID, eprintID, db, tables)
	// ContributorsItemList
	eprint.Contributors = eprintIDToContributors(repoID, eprintID, db, tables)

	// CorpCreators
	eprint.CorpCreators = eprintIDToCorpCreators(repoID, eprintID, db, tables)

	// LocalGroupItemList (SimpleItemList)
	eprint.LocalGroup = eprintIDToLocalGroup(repoID, eprintID, db, tables)
	// FundersItemList (custom)
	eprint.Funders = eprintIDToFunders(repoID, eprintID, db, tables)
	// Documents (*DocumentList)
	eprint.Documents = eprintIDToDocumentList(repoID, baseURL, eprintID, db, tables)
	// RelatedURLs List
	eprint.RelatedURL = eprintIDToRelatedURL(repoID, baseURL, eprintID, db, tables)
	// ReferenceText (item list)
	eprint.ReferenceText = eprintIDToReferenceText(repoID, eprintID, db, tables)
	// Projects
	eprint.Projects = eprintIDToProjects(repoID, eprintID, db, tables)
	// OtherNumberingSystem (item list)
	eprint.OtherNumberingSystem = eprintIDToOtherNumberingSystem(repoID, eprintID, db, tables)
	// Subjects List
	eprint.Subjects = eprintIDToSubjects(repoID, eprintID, db, tables)
	// ItemIssues
	eprint.ItemIssues = eprintIDToItemIssues(repoID, eprintID, db, tables)

	// Exhibitors
	eprint.Exhibitors = eprintIDToExhibitors(repoID, eprintID, db, tables)
	// Producers
	eprint.Producers = eprintIDToProducers(repoID, eprintID, db, tables)
	// Conductors
	eprint.Conductors = eprintIDToConductors(repoID, eprintID, db, tables)

	// Lyricists
	eprint.Lyricists = eprintIDToLyricists(repoID, eprintID, db, tables)

	// Accompaniment
	eprint.Accompaniment = eprintIDToAccompaniment(repoID, eprintID, db, tables)
	// SkillAreas
	eprint.SkillAreas = eprintIDToSkillAreas(repoID, eprintID, db, tables)
	// CopyrightHolders
	eprint.CopyrightHolders = eprintIDToCopyrightHolders(repoID, eprintID, db, tables)
	// Reference
	eprint.Reference = eprintIDToReference(repoID, eprintID, db, tables)

	// ConfCreators
	eprint.ConfCreators = eprintIDToConfCreators(repoID, eprintID, db, tables)
	// AltTitle
	eprint.AltTitle = eprintIDToAltTitle(repoID, eprintID, db, tables)
	// PatentAssignee
	eprint.PatentAssignee = eprintIDToPatentAssignee(repoID, eprintID, db, tables)
	// RelatedPatents
	eprint.RelatedPatents = eprintIDToRelatedPatents(repoID, eprintID, db, tables)
	// Divisions
	eprint.Divisions = eprintIDToDivisions(repoID, eprintID, db, tables)
	// ThesisAdvisor
	eprint.ThesisAdvisor = eprintIDToThesisAdvisors(repoID, eprintID, db, tables)
	// ThesisCommittee
	eprint.ThesisCommittee = eprintIDToThesisCommittee(repoID, eprintID, db, tables)

	// OptionMajor
	eprint.OptionMajor = eprintIDToOptionMajor(repoID, eprintID, db, tables)
	// OptionMinor
	eprint.OptionMinor = eprintIDToOptionMinor(repoID, eprintID, db, tables)

	/*************************************************************
		NOTE: These are notes about possible original implementation
		errors or elements that did not surprive the upgrade to
		EPrints 3.3.16

		eprint.LearningLevels (not an item list in EPrints) using LearningLevelText
		GScholar, skipping not an item list, a 2010 plugin for EPRints 3.2.
		eprint.GScholar = eprintIDToGScholar(repoID, eprintID, db, tables)
		Shelves, a plugin, not replicating, not an item list
		eprint.Shelves = eprintIDToSchelves(repoID, eprintID, db, tables)
		eprint.PatentClassification is not not an item list, using eprint.PatentClassificationText
		eprint.OtherURL appears to be an extraneous
		eprint.CorpContributors apears to be an extraneous
	*************************************************************/

	return eprint, nil
}

// qmList generates an array of string where each element holds "?".
func qmList(length int) []string {
	list := []string{}
	for i := 0; i < length; i++ {
		list = append(list, `?`)
	}
	return list
}

// insertItemList takes an repoID, table name, list of columns and
// an EPrint datastructure then generates and executes a series of
// INSERT statement to create an Item List for the given table.
func insertItemList(repoID string, tableName string, columns []string, eprint *EPrint) error {
	return fmt.Errorf(`insertItemList() not implemented`)
}

// CrosswalkEPrintToSQLCreate will read a EPrint structure and
// generate SQL INSERT, REPLACE and DELETE statements
// suitable for creating a new EPrint record in the repository.
func CrosswalkEPrintToSQLCreate(repoID string, eprint *EPrint) (int, error) {
	var (
		err      error
		eprintID int
	)
	db, ok := config.Connections[repoID]
	if ok == false {
		return 0, fmt.Errorf(`%s is not configured`, repoID)
	}
	dsn, ok := config.Repositories[repoID]
	if ok == false || dsn.Write == false {
		return 0, fmt.Errorf(`%s is not configured for creating eprint records`, repoID)
	}
	// If eprint id is zero generate a sequence of INSERT statements
	// for the record. Others use generate the appropriate History
	// records and then delete insert the new record.
	tableName := `eprint`

	if columns, ok := dsn.TableMap[tableName]; ok {
		// Step one, as a transaction get the highest eprintid, add one and
		// generate an empty eprint row.
		stmt := fmt.Sprintf(`INSERT INTO eprint (eprintid) SELECT IFNULL((SELECT (eprintid + 1) FROM eprint ORDER BY eprintid DESC LIMIT 1), 1)`)
		_, err := db.Exec(stmt)
		if err != nil {
			return 0, fmt.Errorf(`SQL error, %q, %s`, stmt, err)
		}
		stmt = fmt.Sprintf(`SELECT eprintid FROM eprint ORDER BY eprintid DESC LIMIT 1`)
		rows, err := db.Query(stmt)
		if err != nil {
			return 0, fmt.Errorf(`SQL error, %q, %s`, stmt, err)
		}
		id := 0
		for rows.Next() {
			if err := rows.Scan(&id); err != nil {
				return 0, fmt.Errorf(`Could not calculate the new eprintid value, %s`, err)
			}
		}
		rows.Close()
		if err != nil {
			return 0, fmt.Errorf(`SQL failed to get insert id, %s`, err)
		}
		log.Printf("DEBUG retrieved the last insert ID of %d", id)
		eprint.EPrintID = id
		// Step two, write the rest of the date into the main table.
		columnsSQL, values := eprintToColumnsAndValues(eprint, columns, false)
		stmt = fmt.Sprintf(`REPLACE INTO %s (%s) VALUES (%s)`,
			tableName,
			strings.Join(columnsSQL, `, `),
			strings.Join(qmList(len(columnsSQL)), `, `))
		_, err = db.Exec(stmt, values...)
		if err != nil {
			return 0, fmt.Errorf(`SQL error, %q, %s`, stmt, err)
		}
	}
	log.Printf(`DEBUG we should now have eprintid %d %s table populated for %s`, eprint.EPrintID, tableName, repoID)
	if eprint.EPrintID != 0 {
		for tableName, columns := range dsn.TableMap {
			// Handle the remaining tables, i.e. skip eprint table.
			switch {
			case tableName == `eprint`:
			// Skip eprint table, we've already processed it
			case strings.HasPrefix(tableName, `document`):
				log.Printf(`FIXME %s columns: %s`, tableName, strings.Join(columns, `, `))
			case strings.HasPrefix(tableName, `file`):
				log.Printf(`FIXME %s columns: %s`, tableName, strings.Join(columns, `, `))
			default:
				log.Printf(`DEBUG attempted to insert an ItemList %s, %s`, tableName, strings.Join(columns, `, `))
				// Insert new rows in associated table
				if err := insertItemList(repoID, tableName, columns, eprint); err != nil {
					return eprint.EPrintID, fmt.Errorf(`Failed to insert eprintid %d in %s for %s, %s`, eprint.EPrintID, tableName, repoID, err)
				}
			}
		}
		return eprintID, nil
	}
	return 0, err
}

// ImportEPrints take an repository id, eprints structure.
// It is a re-implementation of EPrints Perl Import EPrint XML.
//
// ImportEPrints will create EPrint records using
// CrosswalkEPrintToSQLCreate function.
//
// ImportEPrints returns a list of EPrint IDs created
// if successful and an error if something goes wrong.
func ImportEPrints(repoID string, eprints *EPrints) ([]int, error) {
	var importErrors error
	ids := []int{}

	// Check to make sure updates are allowed if non-Zero
	// eprint ids present.
	for _, eprint := range eprints.EPrint {
		if eprint.EPrintID != 0 {
			return nil, fmt.Errorf("Create failed eprint id %d in %s", eprint.EPrintID, repoID)
		}
	}
	for _, eprint := range eprints.EPrint {
		id, err := CrosswalkEPrintToSQLCreate(repoID, eprint)
		if err != nil {
			if importErrors == nil {
				importErrors = err
			} else {
				importErrors = fmt.Errorf("%s; %s", importErrors, err)
			}
		}
		ids = append(ids, id)
	}
	return ids, importErrors
}
