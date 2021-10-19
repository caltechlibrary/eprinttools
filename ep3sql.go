package eprinttools

//
// ep3sql.go provides crosswalk methods to/from SQL
//
import (
	"fmt"
	"log"
	"strings"
)

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

// eprintToScanValues map the column names provided into the EPrint
// struct.
func eprintToScanValues(eprint *EPrint, columnNames []string) ([]interface{}, error) {
	values := make([]interface{}, len(columnNames))
	for i, columnName := range columnNames {
		switch columnName {
		case "eprintid":
			values[i] = &eprint.EPrintID
		case "rev_number":
			values[i] = &eprint.RevNumber
		case "eprint_status":
			values[i] = &eprint.EPrintStatus
		case "userid":
			values[i] = &eprint.UserID
		case "dir":
			values[i] = &eprint.Dir
		case "datestamp_year":
			values[i] = &eprint.DatestampYear
		case "datestamp_month":
			values[i] = &eprint.DatestampMonth
		case "datestamp_day":
			values[i] = &eprint.DatestampDay
		case "datestamp_hour":
			values[i] = &eprint.DatestampHour
		case "datestamp_minute":
			values[i] = &eprint.DatestampMinute
		case "datestamp_second":
			values[i] = &eprint.DatestampSecond
		case "lastmod_year":
			values[i] = &eprint.LastModifiedYear
		case "lastmod_month":
			values[i] = &eprint.LastModifiedMonth
		case "lastmod_day":
			values[i] = &eprint.LastModifiedDay
		case "lastmod_hour":
			values[i] = &eprint.LastModifiedHour
		case "lastmod_minute":
			values[i] = &eprint.LastModifiedMinute
		case "lastmod_second":
			values[i] = &eprint.LastModifiedSecond
		case "status_changed_year":
			values[i] = &eprint.StatusChangedYear
		case "status_changed_month":
			values[i] = &eprint.StatusChangedMonth
		case "status_changed_day":
			values[i] = &eprint.StatusChangedDay
		case "status_changed_hour":
			values[i] = &eprint.StatusChangedHour
		case "status_changed_minute":
			values[i] = &eprint.StatusChangedMinute
		case "status_changed_second":
			values[i] = &eprint.StatusChangedSecond
		case "type":
			values[i] = &eprint.Type
		case "metadata_visibility":
			values[i] = &eprint.MetadataVisibility
		case "title":
			values[i] = &eprint.Title
		case "ispublished":
			values[i] = &eprint.IsPublished
		case "full_text_status":
			values[i] = &eprint.FullTextStatus
		case "keywords":
			values[i] = &eprint.Keywords
		case "note":
			values[i] = &eprint.Note
		case "abstract":
			values[i] = &eprint.Abstract
		case "date_year":
			values[i] = &eprint.DateYear
		case "date_month":
			values[i] = &eprint.DateMonth
		case "date_day":
			values[i] = &eprint.DateDay
		case "date_type":
			values[i] = &eprint.DateType
		case "series":
			values[i] = &eprint.Series
		case "volume":
			values[i] = &eprint.Volume
		case "number":
			values[i] = &eprint.Number
		case "publication":
			values[i] = &eprint.Publication
		case "publisher":
			values[i] = &eprint.Publisher
		case "place_of_pub":
			values[i] = &eprint.PlaceOfPub
		case "edition":
			values[i] = &eprint.Edition
		case "pagerange":
			values[i] = &eprint.PageRange
		case "pages":
			values[i] = &eprint.Pages
		case "event_type":
			values[i] = &eprint.EventType
		case "event_title":
			values[i] = &eprint.EventTitle
		case "event_location":
			values[i] = &eprint.EventLocation
		case "event_dates":
			values[i] = &eprint.EventDates
		case "id_number":
			values[i] = &eprint.IDNumber
		case "refereed":
			values[i] = &eprint.Refereed
		case "isbn":
			values[i] = &eprint.ISBN
		case "issn":
			values[i] = &eprint.ISSN
		case "book_title":
			values[i] = &eprint.BookTitle
		case "official_url":
			values[i] = &eprint.OfficialURL
		case "alt_url":
			values[i] = &eprint.AltURL
		case "rights":
			values[i] = &eprint.Rights
		case "collection":
			values[i] = &eprint.Collection
		case "reviewer":
			values[i] = &eprint.Reviewer
		case "official_cit":
			values[i] = &eprint.OfficialCitation
		case "monograph_type":
			values[i] = &eprint.MonographType
		case "suggestions":
			values[i] = &eprint.Suggestions
		case "pres_type":
			values[i] = &eprint.PresType
		case "suceeds":
			values[i] = &eprint.Succeeds
		case "commentary":
			values[i] = &eprint.Commentary
		case "contact_email":
			values[i] = &eprint.ContactEMail
		case "fileinfo":
			values[i] = &eprint.FileInfo
		case "latitude":
			values[i] = &eprint.Latitude
		case "longitude":
			values[i] = &eprint.Longitude
		case "department":
			values[i] = &eprint.Department
		case "output_media":
			values[i] = &eprint.OutputMedia
		case "num_pieces":
			values[i] = &eprint.NumPieces
		case "composition_type":
			values[i] = &eprint.CompositionType
		case "data_type":
			values[i] = &eprint.DataType
		case "pedagogic_type":
			values[i] = &eprint.PedagogicType
		case "learning_level":
			values[i] = &eprint.LearningLevel
		case "completion_time":
			values[i] = &eprint.CompletionTime
		case "task_purpose":
			values[i] = &eprint.TaskPurpose
		case "doi":
			values[i] = &eprint.DOI
		case "pmc_id":
			values[i] = &eprint.PMCID
		case "pmid":
			values[i] = &eprint.PMID
		case "parent_url":
			values[i] = &eprint.ParentURL
		case "toc":
			values[i] = &eprint.TOC
		case "interviewer":
			values[i] = &eprint.Interviewer
		case "interviewdate":
			values[i] = &eprint.InterviewDate
		case "nonsubj_keywords":
			values[i] = &eprint.NonSubjKeywords
		case "season":
			values[i] = &eprint.Season
		case "classification_code":
			values[i] = &eprint.ClassificationCode
		case "sword_depositor":
			values[i] = &eprint.SwordDepositor
		case "sword_depository":
			values[i] = &eprint.SwordDepository
		case "sword_slug":
			values[i] = &eprint.SwordSlug
		case "importid":
			values[i] = &eprint.ImportID
		case "patent_applicant":
			values[i] = &eprint.PatentApplicant
		case "patent_number":
			values[i] = &eprint.PatentNumber
		case "institution":
			values[i] = &eprint.Institution
		case "thesis_type":
			values[i] = &eprint.ThesisType
		case "thesis_degree":
			values[i] = &eprint.ThesisDegree
		case "thesis_degree_grantor":
			values[i] = &eprint.ThesisDegreeGrantor
		case "thesis_degree_date_year":
			values[i] = &eprint.ThesisDegreeDateYear
		case "thesis_degree_date_month":
			values[i] = &eprint.ThesisDegreeDateMonth
		case "thesis_degree_date_day":
			values[i] = &eprint.ThesisDegreeDateDay
		case "thesis_submitted_date_year":
			values[i] = &eprint.ThesisSubmittedDateYear
		case "thesis_submitted_date_month":
			values[i] = &eprint.ThesisSubmittedDateMonth
		case "thesis_submitted_date_day":
			values[i] = &eprint.ThesisSubmittedDateDay
		case "thesis_defense_date_year":
			values[i] = &eprint.ThesisDefenseDateYear
		case "thesis_defense_date_month":
			values[i] = &eprint.ThesisDefenseDateMonth
		case "thesis_defense_date_day":
			values[i] = &eprint.ThesisDefenseDateDay
		case "thesis_approved_date_year":
			values[i] = &eprint.ThesisApprovedDateYear
		case "thesis_approved_date_month":
			values[i] = &eprint.ThesisApprovedDateMonth
		case "thesis_approved_date_day":
			values[i] = &eprint.ThesisApprovedDateDay
		case "thesis_public_date_year":
			values[i] = &eprint.ThesisPublicDateYear
		case "thesis_public_date_month":
			values[i] = &eprint.ThesisPublicDateMonth
		case "thesis_public_date_day":
			values[i] = &eprint.ThesisPublicDateDay
		case "thesis_author_email":
			values[i] = &eprint.ThesisAuthorEMail
		case "hide_thesis_author_email":
			values[i] = &eprint.HideThesisAuthorEMail
		case "gradofc_approval_date":
			values[i] = &eprint.GradOfficeApprovalDate
		case "thesis_awards":
			values[i] = &eprint.ThesisAwards
		case "review_status":
			values[i] = &eprint.ReviewStatus
		case "copyright_statement":
			values[i] = &eprint.CopyrightStatement
		case "source":
			values[i] = &eprint.Source
		case "succeeds":
			values[i] = &eprint.Succeeds
		case "replacedby":
			values[i] = &eprint.ReplacedBy
		case "item_issues_count":
			values[i] = &eprint.ItemIssuesCount
		case "errata":
			values[i] = &eprint.ErrataText
		case "edit_lock_user":
			values[i] = &eprint.EditLockUser
		case "edit_lock_since":
			values[i] = &eprint.EditLockSince
		case "edit_lock_until":
			values[i] = &eprint.EditLockUntil
		// The follow values represent sub tables and processed separately.
		case "patent_classification":
			values[i] = &eprint.PatentClassificationText
		default:
			// Handle case where we have value that is unmapped or not available in EPrint struct
			return nil, fmt.Errorf("could not map %q into EPrint struct", columnName)
		}
	}
	return values, nil
}

// CrosswalkSQLToEPrint expects a repository map and EPrint ID
// and will generate a series of SELECT statements populating
// a new EPrint struct or return an error
func CrosswalkSQLToEPrint(repoID string, eprintID int) (*EPrint, error) {
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

	//FIXME: we need to figure out what to do where one repository has a
	// column name matching a Item List.

	// Map the list of addresses of the EPrint structure to
	// the values array of pointers.
	// Build a placehold for results, a list of pointers
	values, err := eprintToScanValues(eprint, columns)
	if err != nil {
		return nil, fmt.Errorf("internal server error, eprint to scan values, %s", err)
	}
	if len(values) != len(columns) {
		return nil, fmt.Errorf("internal server error, could not be placeholder for results")
	}

	// NOTE: With the "values" pointer array setup the query can be built
	// and executed in the usually SQL fashion.
	stmt := fmt.Sprintf("SELECT %s FROM eprint WHERE eprintid = ? LIMIT 1", strings.Join(columns, ", "))
	rows, err := db.Query(stmt, eprintID)
	if err != nil {
		return nil, fmt.Errorf("ERROR: query error (%q, %q), %s", repoID, stmt, err)
	}
	for rows.Next() {
		// NOTE: Because values array holds the addresses into our
		// EPrint struct the "Scan" does the actual mapping.
		if err := rows.Scan(values...); err != nil {
			log.Printf("Could not read columns, %s", err)
		}
	}
	rows.Close()

	// CreatorsItemList
	subTables := []string{}
	for table := range tables {
		if strings.HasPrefix(table, "eprint_creators") {
			subTables = append(subTables, table)
		}
	}
	if includesString(subTables, "eprint_creators_id") {
		var (
			creatorid, honourific, given, family, lineage, orcid, uri string
		)
		// eprint_%_id is a known structure. eprintid, pos, creators_id
		creatorItemList := new(CreatorItemList)
		rows, err = db.Query(`SELECT creators_id FROM eprint_creators_id WHERE eprintid = ? ORDER BY eprintid, pos`, eprintID)
		if err != nil {
			log.Printf("Could not scan eprint_creators_id for %d in %q, %s", eprintID, repoID, err)
		} else {
			i := 0
			for rows.Next() {
				if err := rows.Scan(&creatorid); err != nil {
					log.Printf("Could not scan eprint_creators_id for %d in %q, %s", eprintID, repoID, err)
				} else {
					item := new(Item)
					if creatorid != "" {
						item.ID = creatorid
					} else {
						item.ID = ""
					}
					creatorItemList.Items = append(creatorItemList.Items, item)
				}
				i++
			}
			rows.Close()
			log.Printf("DEBUG creatorsItems found %d", len(creatorItemList.Items))
		}
		if includesString(subTables, "eprint_creators_name") {
			rows, err = db.Query(`SELECT creators_name_honourific, creators_name_given, creators_name_family, creators_name_lineage FROM eprint_creators_name WHERE eprintid = ? ORDER BY eprintid, pos`, eprintID)
			if err != nil {
				log.Printf("Could not scan eprint_creators_name for %d in %q, %s", eprintID, repoID, err)
			} else {
				i := 0
				for rows.Next() {
					if err := rows.Scan(&honourific, &given, &family, &lineage); err != nil {
						log.Printf("Could not scan eprint_creators_name for %d in %q, %s", eprintID, repoID, err)
					} else if i < len(creatorItemList.Items) {
						creatorItemList.Items[i].Name = new(Name)
						creatorItemList.Items[i].Name.Given = given
						creatorItemList.Items[i].Name.Family = family
					}
					i++
				}
				rows.Close()
			}
		}
		if includesString(subTables, "eprint_creators_orcids") {
			rows, err = db.Query(`SELECT creators_orcid WHERE eprintid = ? ORDER BY eprintid, pos`, eprintID)
			if err != nil {
				log.Printf("Could not scan eprint_creators_orcids for %d in %q, %s", eprintID, repoID, err)
			} else {
				i := 0
				for rows.Next() {
					if err := rows.Scan(&orcid); err != nil {
						log.Printf("Could not scan eprint_creators_orcid for %d in %q, %s", eprintID, repoID, err)
					} else if i < len(creatorItemList.Items) {
						if orcid != "" {
							creatorItemList.Items[i].ORCID = orcid
						}
					}
					i++
				}
				rows.Close()
			}
		}
		if includesString(subTables, "eprint_creators_uri") {
			rows, err = db.Query(`SELECT creators_uri FROM eprint_creators_uri WHERE eprintid = ? ORDER BY eprintid, pos`, eprintID)
			if err != nil {
				log.Printf("Could not scan eprint_creators_uri for %d in %q, %s", eprintID, repoID, err)
			} else {
				i := 0
				for rows.Next() {
					if err := rows.Scan(&uri); err != nil {
						log.Printf("Could not scan eprint_creators_orcid for %d in %q, %s", eprintID, repoID, err)
					} else if i < len(creatorItemList.Items) {
						creatorItemList.Items[i].URI = uri
					}
					i++
				}
				rows.Close()
			}
		}
		eprint.Creators = creatorItemList
	}

	// EditorsItemList
	// ContributorsItemList
	// LocalGroupItemList
	// ErrataItemList (or ErrataText)

	//FIXME: Need to map the values from the various date/time stamps to
	//there datetime/date fields string.

	return eprint, nil
}

// CrosswalkEPrintToSQL will read an EPrint structure and
// generate SQL INSERT (eprint.ID == 0) or REPLACE statements
// suitable for updating an EPrint record in the repository.
func CrosswalkEPrintToSQL(eprint *EPrint) ([]byte, error) {
	return nil, fmt.Errorf("Not implemented")
}
