//
// Package eprinttools is a collection of structures, functions and programs// for working with the EPrints XML and EPrints REST API
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2021, Caltech
// All rights not granted herein are expressly reserved by Caltech.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
package eprinttools

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"
)

//
// These default values are used when apply clsrules set
//
var (
	// DefaultCollection holds the default collection to use on deposit
	DefaultCollection string
	// DefaultRights sets the eprint.Rights to a default value on deposit
	DefaultRights string
	// DefaultOfficialURL holds a URL prefix for the persistent URL,
	// an ID Number will get added when generating per record
	// official_url values.
	DefaultOfficialURL string
	// DefaultRefereed
	DefaultRefereed string
	// DefaultStatus
	DefaultStatus string
)

//
// NOTE: This file contains the general structure in Caltech Libraries EPrints 3.x based repositories.
//

// EPrints is the high level XML you get from the REST API.
// E.g. curl -L -O https://eprints3.example.org/rest/eprint/1234.xml
// Then parse the 1234.xml document stucture.
type EPrints struct {
	XMLName xml.Name  `xml:"eprints" json:"-"`
	XMLNS   string    `xml:"xmlns,attr,omitempty" json:"xmlns,omitempty"`
	EPrint  []*EPrint `xml:"eprint" json:"eprint"`
}

// EPrint is the record contated in a EPrints XML document such as they used
// to store revisions.
type EPrint struct {
	XMLName             xml.Name         `json:"-"`
	ID                  string           `xml:"id,attr,omitempty" json:"id,omitempty"`
	EPrintID            int              `xml:"eprintid,omitempty" json:"eprint_id,omitempty"`
	RevNumber           int              `xml:"rev_number,omitempty" json:"rev_number,omitempty"`
	Documents           *DocumentList    `xml:"documents>document,omitempty" json:"documents,omitempty"`
	EPrintStatus        string           `xml:"eprint_status,omitempty" json:"eprint_status,omitempty"`
	UserID              int              `xml:"userid,omitempty" json:"userid,omitempty"`
	Dir                 string           `xml:"dir,omitempty" json:"dir,omitempty"`
	Datestamp           string           `xml:"datestamp,omitempty" json:"datestamp,omitempty"`
	DatestampYear       int              `xml:"-" json:"-"`
	DatestampMonth      int              `xml:"-" json:"-"`
	DatestampDay        int              `xml:"-" json:"-"`
	DatestampHour       int              `xml:"-" json:"-"`
	DatestampMinute     int              `xml:"-" json:"-"`
	DatestampSecond     int              `xml:"-" json:"-"`
	LastModified        string           `xml:"lastmod,omitempty" json:"lastmod,omitempty"`
	LastModifiedYear    int              `xml:"-" json:"-"`
	LastModifiedMonth   int              `xml:"-" json:"-"`
	LastModifiedDay     int              `xml:"-" json:"-"`
	LastModifiedHour    int              `xml:"-" json:"-"`
	LastModifiedMinute  int              `xml:"-" json:"-"`
	LastModifiedSecond  int              `xml:"-" json:"-"`
	StatusChanged       string           `xml:"status_changed,omitempty" json:"status_changed,omitempty"`
	StatusChangedYear   int              `xml:"-" json:"-"`
	StatusChangedMonth  int              `xml:"-" json:"-"`
	StatusChangedDay    int              `xml:"-" json:"-"`
	StatusChangedHour   int              `xml:"-" json:"-"`
	StatusChangedMinute int              `xml:"-" json:"-"`
	StatusChangedSecond int              `xml:"-" json:"-"`
	Type                string           `xml:"type,omitempty" json:"type,omitempty"`
	MetadataVisibility  string           `xml:"metadata_visibility,omitempty" json:"metadata_visibility,omitempty"`
	Creators            *CreatorItemList `xml:"creators,omitempty" json:"creators,omitempty"`
	Title               string           `xml:"title,omitempty" json:"title,omitempty"`
	IsPublished         string           `xml:"ispublished,omitempty" json:"ispublished,omitempty"`
	FullTextStatus      string           `xml:"full_text_status,omitempty" json:"full_text_status,omitempty"`
	Keywords            string           `xml:"keywords,omitempty" json:"keywords,omitempty"`
	//Keyword              *KeywordItemList              `xml:"-" json:""`
	Note                 string                        `xml:"note,omitempty" json:"note,omitempty"`
	Abstract             string                        `xml:"abstract,omitempty" json:"abstract,omitempty"`
	Date                 string                        `xml:"date,omitempty" json:"date,omitempty"`
	DateYear             int                           `xml:"-" json:"-"`
	DateMonth            int                           `xml:"-" json:"-"`
	DateDay              int                           `xml:"-" json:"-"`
	DateType             string                        `xml:"date_type,omitempty" json:"date_type,omitempty"`
	Series               string                        `xml:"series,omitempty" json:"series,omitempty"`
	Publication          string                        `xml:"publication,omitempty" json:"publication,omitempty"`
	Volume               string                        `xml:"volume,omitempty" json:"volume,omitempty"`
	Number               string                        `xml:"number,omitempty" json:"number,omitempty"`
	Publisher            string                        `xml:"publisher,omitempty" json:"publisher,omitempty"`
	PlaceOfPub           string                        `xml:"place_of_pub,omitempty" json:"place_of_pub,omitempty"`
	Edition              string                        `xml:"edition,omitempty" json:"edition,omitempty"`
	PageRange            string                        `xml:"pagerange,omitempty" json:"pagerange,omitempty"`
	Pages                int                           `xml:"pages,omitempty" json:"pages,omitempty"`
	EventType            string                        `xml:"event_type,omitempty" json:"event_type,omitempty"`
	EventTitle           string                        `xml:"event_title,omitempty" json:"event_title,omitempty"`
	EventLocation        string                        `xml:"event_location,omitempty" json:"event_location,omitempty"`
	EventDates           string                        `xml:"event_dates,omitempty" json:"event_dates,omitempty"`
	IDNumber             string                        `xml:"id_number,omitempty" json:"id_number,omitempty"`
	Refereed             string                        `xml:"refereed,omitempty" json:"refereed,omitempty"`
	ISBN                 string                        `xml:"isbn,omitempty" json:"isbn,omitempty"`
	ISSN                 string                        `xml:"issn,omitempty" json:"issn,omitempty"`
	BookTitle            string                        `xml:"book_title,omitempty" json:"book_title,omitempty"`
	Editors              *EditorItemList               `xml:"editors,omitempty" json:"editors,omitempty"`
	OfficialURL          string                        `xml:"official_url,omitempty" json:"official_url,omitempty"`
	AltURL               string                        `xml:"alt_url,omitempty" json:"alt_url,omitempty"`
	RelatedURL           *RelatedURLItemList           `xml:"related_url,omitempty" json:"related_url,omitempty"`
	ReferenceText        *ReferenceTextItemList        `xml:"referencetext,omitempty" json:"referencetext,omitempty"`
	Projects             *ProjectItemList              `xml:"projects,omitempty" json:"projects,omitempty"`
	Rights               string                        `xml:"rights,omitempty" json:"rights,omitempty"`
	Funders              *FunderItemList               `xml:"funders,omitempty" json:"funders,omitempty"`
	Collection           string                        `xml:"collection,omitempty" json:"collection,omitempty"`
	Reviewer             string                        `xml:"reviewer,omitempty" json:"reviewer,omitempty"`
	OfficialCitation     string                        `xml:"official_cit,omitempty" json:"official_cit,omitempty"`
	OtherNumberingSystem *OtherNumberingSystemItemList `xml:"other_numbering_system,omitempty" json:"other_numbering_system,omitempty"`
	LocalGroup           *LocalGroupItemList           `xml:"local_group,omitempty" json:"local_group,omitempty"`
	ErrataText           string                        `xml:"errata,omitempty" json:"errata,omitempty"`
	Contributors         *ContributorItemList          `xml:"contributors,omitempty" json:"contributors,omitempty"`
	MonographType        string                        `xml:"monograph_type,omitempty" json:"monograph_type,omitempty"`

	// Caltech Library uses suggestions as an internal note field (RSD, 2018-02-15)
	Suggestions string `xml:"suggestions,omitempty" json:"suggestions,omitempty"`

	// CaletchLN has a "coverage_dates" field in the eprint table.
	CoverageDates string `xml:"coverage_dates,omitempty" json:"coverage_dates,omitempty"`

	// NOTE: Misc fields discoverd exploring REST API records, not currently used at Caltech Library (RSD, 2018-01-02)
	Subjects     *SubjectItemList `xml:"subjects,omitempty" json:"subjects,omitempty"`
	PresType     string           `xml:"pres_type,omitempty" json:"presentation_type,omitempty"`
	Succeeds     int              `xml:"succeeds,omitempty" json:"succeeds,omitempty"`
	Commentary   int              `xml:"commentary,omitempty" json:"commentary,omitempty"`
	ContactEMail string           `xml:"contact_email,omitempty" json:"contect_email,omitempty"`
	// NOTE: EPrints XML doesn't include fileinfo
	FileInfo          string                   `xml:"-" json:"-"`
	Latitude          float64                  `xml:"latitude,omitempty" json:"latitude,omitempty"`
	Longitude         float64                  `xml:"longitude,omitempty" json:"longitude,omitempty"`
	ItemIssues        *ItemIssueItemList       `xml:"item_issues,omitempty" json:"item_issues,omitempty"`
	ItemIssuesCount   int                      `xml:"item_issues_count,omitempty" json:"item_issues_count,omitempty"`
	CorpCreators      *CorpCreatorItemList     `xml:"corp_creators,omitempty" json:"corp_creators,omitempty"`
	CorpContributors  *CorpContributorItemList `xml:"corp_contributors,omitempty" json:"corp_contributors,omitempty"`
	Department        string                   `xml:"department,omitempty" json:"department,omitempty"`
	OutputMedia       string                   `xml:"output_media,omitempty" json:"output_media,omitempty"`
	Exhibitors        *ExhibitorItemList       `xml:"exhibitors,omitempty" json:"exhibitors,omitempty"`
	NumPieces         int                      `xml:"num_pieces,omitempty" json:"num_pieces,omitempty"`
	CompositionType   string                   `xml:"composition_type,omitempty" json:"composition_type,omitempty"`
	Producers         *ProducerItemList        `xml:"producers,omitempty" json:"producers,omitempty"`
	Conductors        *ConductorItemList       `xml:"conductors,omitempty" json:"conductors,omitempty"`
	Lyricists         *LyricistItemList        `xml:"lyricists,omitempty" json:"lyricists,omitempty"`
	Accompaniment     *AccompanimentItemList   `xml:"accompaniment,omitempty" json:"accompaniment,omitempty"`
	DataType          string                   `xml:"data_type,omitempty" json:"data_type,omitempty"`
	PedagogicType     string                   `xml:"pedagogic_type,omitempty" json:"pedagogic_type,omitempty"`
	CompletionTime    string                   `xml:"completion_time,omitempty" json:"completion_time,omitempty"`
	TaskPurpose       string                   `xml:"task_purpose,omitempty" json:"task_purpose,omitempty"`
	SkillAreas        *SkillAreaItemList       `xml:"skill_areas,omitempty" json:"skill_areas,omitempty"`
	CopyrightHolders  *CopyrightHolderItemList `xml:"copyright_holders,omitempty" json:"copyright_holders,omitempty"`
	LearningLevelText string                   `xml:"learning_level,omitempty" json:"learning_level,omitempty"`
	//LearningLevel      *LearningLevelItemList   `xml:"-" json:"-"`
	DOI           string               `xml:"doi,omitempty" json:"doi,omitempty"`
	PMCID         string               `xml:"pmc_id,omitempty" json:"pmcid,omitempty"`
	PMID          string               `xml:"pmid,omitempty" json:"pmid,omitempty"`
	ParentURL     string               `xml:"parent_url,omitempty" json:"parent_url,omitempty"`
	Reference     *ReferenceItemList   `xml:"reference,omitempty" json:"reference,omitempty"`
	ConfCreators  *ConfCreatorItemList `xml:"conf_creators,omitempty" json:"conf_creators,omitempty"`
	AltTitle      *AltTitleItemList    `xml:"alt_title,omitempty" json:"alt_title,omitempty"`
	TOC           string               `xml:"toc,omitempty" json:"toc,omitempty"`
	Interviewer   string               `xml:"interviewer,omitempty" json:"interviewer,omitempty"`
	InterviewDate string               `xml:"interviewdate,omitempty" json:"interviewdate,omitempty"`
	//GScholar           *GScholarItemList    `xml:"gscholar,omitempty" json:"gscholar,omitempty"`
	NonSubjKeywords    string            `xml:"nonsubj_keywords,omitempty" json:"nonsubj_keywords,omitempty"`
	Season             string            `xml:"season,omitempty" json:"season,omitempty"`
	ClassificationCode string            `xml:"classification_code,omitempty" json:"classification_code,omitempty"`
	Shelves            *ShelfItemList    `xml:"shelves,omitempty" json:"shelves,omitempty"`
	Relation           *RelationItemList `xml:"relation,omitempty" json:"relation,omitempty"`

	// NOTE: Sword deposit fields
	SwordDepository string `xml:"sword_depository,omitempty" json:"sword_depository,omitempty"`
	SwordDepositor  int    `xml:"sword_depositor,omitempty" json:"sword_depositor,omitempty"`
	SwordSlug       string `xml:"sword_slug,omitempty" json:"sword_slug,omitempty"`
	ImportID        int    `xml:"importid,omitempty" json:"import_id,omitempty"`

	// Patent related fields
	PatentApplicant          string                  `xml:"patent_applicant,omitempty" json:"patent_applicant,omitempty"`
	PatentNumber             string                  `xml:"patent_number,omitempty" json:"patent_number,omitempty"`
	PatentAssignee           *PatentAssigneeItemList `xml:"patent_assignee,omitempty" json:"patent_assignee,omitempty"`
	PatentClassificationText string                  `xml:"-" json:"-"`
	//PatentClassification     *PatentClassificationItemList `xml:"patent_classification,omitempty" json:"patent_classification,omitempty"`
	RelatedPatents *RelatedPatentItemList `xml:"related_patents,omitempty" json:"related_patents,omitempty"`

	// Thesis oriented fields
	Divisions                *DivisionItemList        `xml:"divisions,omitemmpty" json:"divisions,omitempty"`
	Institution              string                   `xml:"institution,omitempty" json:"institution,omitempty"`
	ThesisType               string                   `xml:"thesis_type,omitempty" json:"thesis_type,omitempty"`
	ThesisAdvisor            *ThesisAdvisorItemList   `xml:"thesis_advisor,omitempty" json:"thesis_advisor,omitempty"`
	ThesisCommittee          *ThesisCommitteeItemList `xml:"thesis_committee,omitempty" json:"thesis_committee,omitempty"`
	ThesisDegree             string                   `xml:"thesis_degree,omitempty" json:"thesis_degree,omitempty"`
	ThesisDegreeGrantor      string                   `xml:"thesis_degree_grantor,omitempty" json:"thesis_degree_grantor,omitempty"`
	ThesisDegreeDate         string                   `xml:"thesis_degree_date,omitempty" json:"thesis_degree_date,omitempty"`
	ThesisDegreeDateYear     int                      `xml:"-" json:"-"`
	ThesisDegreeDateMonth    int                      `xml:"-" json:"-"`
	ThesisDegreeDateDay      int                      `xml:"-" json:"-"`
	ThesisSubmittedDate      string                   `xml:"thesis_submitted_date,omitempty" json:"thesis_submitted_date,omitempty"`
	ThesisSubmittedDateYear  int                      `xml:"-" json:"-"`
	ThesisSubmittedDateMonth int                      `xml:"-" json:"-"`
	ThesisSubmittedDateDay   int                      `xml:"-" json:"-"`
	ThesisDefenseDate        string                   `xml:"thesis_defense_date,omitempty" json:"thesis_defense_date,omitempty"`
	ThesisDefenseDateYear    int                      `xml:"-" json:"-"`
	ThesisDefenseDateMonth   int                      `xml:"-" json:"-"`
	ThesisDefenseDateDay     int                      `xml:"-" json:"-"`
	ThesisApprovedDate       string                   `xml:"thesis_approved_date,omitempty" json:"thesis_approved_date,omitempty"`
	ThesisApprovedDateYear   int                      `xml:"-" json:"-"`
	ThesisApprovedDateMonth  int                      `xml:"-" json:"-"`
	ThesisApprovedDateDay    int                      `xml:"-" json:"-"`
	ThesisPublicDate         string                   `xml:"thesis_public_date,omitempty" json:"thesis_public_date,omitempty"`
	ThesisPublicDateYear     int                      `xml:"-" json:"-"`
	ThesisPublicDateMonth    int                      `xml:"-" json:"-"`
	ThesisPublicDateDay      int                      `xml:"-" json:"-"`
	ThesisAuthorEMail        string                   `xml:"thesis_author_email,omitempty" json:"thesis_author_email,omitempty"`
	HideThesisAuthorEMail    string                   `xml:"hide_thesis_author_email,omitempty" json:"hide_thesis_author_email,omitempty"`
	// NOTE: GradOfficeApproval isn't output by CaltechTHESIS.
	GradOfficeApprovalDate      string               `xml:"gradofc_approval_date,omitempty" json:"gradofc_approval_date,omitempty"`
	GradOfficeApprovalDateYear  int                  `xml:"-" json:"-"`
	GradOfficeApprovalDateMonth int                  `xml:"-" json:"-"`
	GradOfficeApprovalDateDay   int                  `xml:"-" json:"-"`
	ThesisAwards                string               `xml:"thesis_awards,omitempty" json:"thesis_awards,omitempty"`
	ReviewStatus                string               `xml:"review_status,omitempty" json:"review_status,omitempty"`
	OptionMajor                 *OptionMajorItemList `xml:"option_major,omitempty" json:"option_major,omitempty"`
	OptionMinor                 *OptionMinorItemList `xml:"option_minor,omitempty" json:"option_major,omitempty"`
	CopyrightStatement          string               `xml:"copyright_statement,omitempty" json:"copyright_statement,omitempty"`

	// Custom fields from some EPrints repositories
	Source     string `xml:"source,omitempty" json:"source,omitempty"`
	ReplacedBy int    `xml:"replacedby,omitempty" json:"replacedby,omitempty"`

	// Edit Control Fields
	EditLockUser  int `xml:"-" json:"-"`
	EditLockSince int `xml:"-" json:"-"`
	EditLockUntil int `xml:"-" json:"-"`

	// Fields identified throw harvesting.
	ReferenceTextString string `xml:referencetext,omitempty" json:"referencetext,omitempty"`
	Language            string `xml:"language,omitempty" json:"language,omitempty"`

	// Synthetic fields are created to help in eventual migration of
	// EPrints field data to other JSON formats.
	PrimaryObject  map[string]interface{}   `xml:"-" json:"primary_object,omitempty"`
	RelatedObjects []map[string]interface{} `xml:"-" json:"related_objects,omitempty"`
}

// PubDate returns the publication date or empty string
func (eprint *EPrint) PubDate() string {
	if eprint.DateType == "published" {
		return eprint.Date
	}
	return ""
}

// Item is a generic type used by various fields (e.g. Creator, Division, OptionMajor)
type Item struct {
	XMLName     xml.Name `xml:"item" json:"-"`
	Name        *Name    `xml:"name,omitempty" json:"name,omitempty"`
	Pos         int      `xml:"-" json:"-"`
	ID          string   `xml:"id,omitempty" json:"id,omitempty"`
	EMail       string   `xml:"email,omitempty" json:"email,omitempty"`
	ShowEMail   string   `xml:"show_email,omitempty" json:"show_email,omitempty"`
	Role        string   `xml:"role,omitempty" json:"role,omitempty"`
	URL         string   `xml:"url,omitempty" json:"url,omitempty"`
	Type        string   `xml:"type,omitempty" json:"type,omitempty"`
	Description string   `xml:"description,omitempty" json:"description,omitempty"`
	Agency      string   `xml:"agency,omitempty" json:"agency,omitempty"`
	GrantNumber string   `xml:"grant_number,omitempty" json:"grant_number,omitempty"`
	URI         string   `xml:"uri,omitempty" json:"uri,omitempty"`
	ORCID       string   `xml:"orcid,omitempty" json:"orcid,omitempty"`
	ROR         string   `xml:"ror,omitempty" json:"ror,omitempty"`
	Timestamp   string   `xml:"timestamp,omitempty" json:"timestamp,omitempty"`
	Status      string   `xml:"status,omitempty" json:"status,omitempty"`
	ReportedBy  string   `xml:"reported_by,omitempty" json:"reported_by,omitempty"`
	ResolvedBy  string   `xml:"resolved_by,omitempty" json:"resolved_by,omitempty"`
	Comment     string   `xml:"comment,omitempty" json:"comment,omitempty"`
	Value       string   `xml:",chardata" json:"value,omitempty"`
}

// SetAttribute takes a lower case string and value and sets
// the attribute of the related item.
func (item *Item) SetAttribute(key string, value interface{}) bool {
	switch value.(type) {
	case string:
		value = strings.TrimSpace(value.(string))
	}
	switch strings.ToLower(key) {
	case "name":
		item.Name = value.(*Name)
		return true
	case "pos":
		item.Pos = value.(int)
		return true
	case "id":
		item.ID = value.(string)
		return true
	case "email":
		item.EMail = value.(string)
		return true
	case "showemail":
		item.ShowEMail = value.(string)
		return true
	case "show_email":
		item.ShowEMail = value.(string)
		return true
	case "role":
		item.Role = value.(string)
		return true
	case "url":
		item.URL = value.(string)
		return true
	case "type":
		item.Type = value.(string)
		return true
	case "description":
		item.Description = value.(string)
		return true
	case "agency":
		item.Agency = value.(string)
		return true
	case "grantnumber":
		item.GrantNumber = value.(string)
		return true
	case "grant_number":
		item.GrantNumber = value.(string)
		return true
	case "uri":
		item.URI = value.(string)
		return true
	case "orcid":
		item.ORCID = value.(string)
		return true
	case "ror":
		item.ROR = value.(string)
		return true
	case "timestamp":
		item.Timestamp = value.(string)
		return true
	case "status":
		item.Status = value.(string)
		return true
	case "reportedby":
		item.ReportedBy = value.(string)
		return true
	case "reported_by":
		item.ReportedBy = value.(string)
		return true
	case "resolvedby":
		item.ResolvedBy = value.(string)
		return true
	case "resolved_by":
		item.ResolvedBy = value.(string)
		return true
	case "comment":
		item.Comment = value.(string)
		return true
	case "value":
		item.Value = value.(string)
		return true
	case "":
		item.Value = value.(string)
		return true
	}
	return false
}

// MarshalJSON() is a custom JSON marshaler for Item
func (item *Item) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{}
	flatten := true
	if item.Name != nil {
		m["name"] = item.Name
		flatten = false
	}
	if strings.TrimSpace(item.ID) != "" {
		m["id"] = item.ID
		flatten = false
	}
	if strings.TrimSpace(item.EMail) != "" {
		m["email"] = item.EMail
		flatten = false
	}
	if strings.TrimSpace(item.ShowEMail) != "" {
		m["show_email"] = item.ShowEMail
		flatten = false
	}
	if strings.TrimSpace(item.Role) != "" {
		m["role"] = item.Role
		flatten = false
	}
	if strings.TrimSpace(item.URL) != "" {
		m["url"] = item.URL
		flatten = false
	}
	if strings.TrimSpace(item.Type) != "" {
		m["type"] = item.Type
		flatten = false
	}
	if strings.TrimSpace(item.Description) != "" {
		m["description"] = item.Description
		flatten = false
	}
	if strings.TrimSpace(item.Agency) != "" {
		m["agency"] = item.Agency
		flatten = false
	}
	if strings.TrimSpace(item.GrantNumber) != "" {
		m["grant_number"] = item.GrantNumber
		flatten = false
	}
	if strings.TrimSpace(item.URI) != "" {
		m["uri"] = item.URI
		flatten = false
	}
	if s := strings.TrimSpace(item.ORCID); s != "" {
		//FIXME: should validate the orcid to avoid legacy data issues
		m["orcid"] = s
		flatten = false
	}
	if s := strings.TrimSpace(item.Value); s != "" {
		if flatten == true {
			return json.Marshal(s)
		}
		m["value"] = s
	}
	return json.Marshal(m)
}

func (item *Item) UnmarshalJSON(src []byte) error {
	if bytes.HasPrefix(src, []byte(`"`)) && bytes.HasSuffix(src, []byte(`"`)) {
		item.Value = fmt.Sprintf("%s", bytes.Trim(src, `"`))
		return nil
	}
	m := make(map[string]interface{})
	err := jsonDecode(src, &m)
	if err != nil {
		return err
	}
	for key, value := range m {
		switch key {
		case "name":
			name := new(Name)
			switch value.(type) {
			case string:
				name.Value = value.(string)
			case map[string]interface{}:
				m := value.(map[string]interface{})
				if family, ok := m["family"]; ok == true {
					name.Family = family.(string)
				}
				if given, ok := m["given"]; ok == true {
					name.Given = given.(string)
				}
				if id, ok := m["id"]; ok == true {
					name.ID = id.(string)
				}
				if orcid, ok := m["orcid"]; ok == true {
					name.ID = orcid.(string)
				}
			}
			item.Name = name
		case "id":
			item.ID = value.(string)
		case "email":
			item.EMail = value.(string)
		case "show_email":
			item.ShowEMail = value.(string)
		case "role":
			item.Role = value.(string)
		case "url":
			item.URL = value.(string)
		case "type":
			item.Type = value.(string)
		case "description":
			item.Description = value.(string)
		case "agency":
			item.Agency = value.(string)
		case "grant_number":
			item.GrantNumber = value.(string)
		case "uri":
			item.URI = value.(string)
		case "orcid":
			item.ORCID = value.(string)
		case "value":
			item.Value = value.(string)
		}
	}
	return err
}

// ItemsInterface describes a common set of operations on an item list.
type ItemsInterface interface {
	// Append an item to an ItemList
	Append(*Item) int
	// Length returns the item count
	Length() int
	// IndexOf returns Item or nil
	IndexOf(int) *Item
	// SetAttributOf at index position sets an item's attribute
	SetAttributeOf(int, string, interface{}) bool
}

// CreatorItemList holds a list of authors
type CreatorItemList struct {
	XMLName xml.Name `xml:"creators" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the Creator list and returns the new count of items
func (itemList *CreatorItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length() returns item count
func (itemList *CreatorItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf() returns item or nil
func (itemList *CreatorItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *CreatorItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 {
		if i >= itemList.Length() {
			for j := itemList.Length(); j <= i; j++ {
				item := new(Item)
				itemList.Append(item)
			}
		}
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// EditorItemList holds a list of editors
type EditorItemList struct {
	XMLName xml.Name `xml:"editors" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the Editor item list and returns the new count of items
func (itemList *EditorItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length() returns the number of items in the list
func (itemList *EditorItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf() returns an item in the list or nil
func (itemList *EditorItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *EditorItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// RelatedURLItemList holds the related URLs (e.g. doi, aux material doi)
type RelatedURLItemList struct {
	XMLName xml.Name `xml:"related_url" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append an item to the related url item list
func (itemList *RelatedURLItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length() returns item count
func (itemList *RelatedURLItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf() returns item or nil
func (itemList *RelatedURLItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *RelatedURLItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// ReferenceTextItemList
type ReferenceTextItemList struct {
	XMLName xml.Name `xml:"referencetext" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the reference text url item list and returns the new count of items
func (itemList *ReferenceTextItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length returns the length of an ReferenceTextItemList
func (itemList *ReferenceTextItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf returns an Item or nil
func (itemList *ReferenceTextItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *ReferenceTextItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// UnmarshJSON takes a reference text list of item and returns
// an appropriately values to assigned struct.
func (itemList *ReferenceTextItemList) UnmarshalJSON(src []byte) error {
	var values []string

	m := make(map[string][]interface{})
	err := jsonDecode(src, &m)
	if err != nil {
		return err
	}
	if itemList, ok := m["items"]; ok == true {
		for _, item := range itemList {
			switch item.(type) {
			case string:
				values = append(values, item.(string))
			}
		}
	}
	for _, value := range values {
		item := new(Item)
		item.SetAttribute("value", value)
		itemList.Append(item)
	}
	return err
}

// ProjectItemList
type ProjectItemList struct {
	XMLName xml.Name `xml:"projects" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the project item list and returns the new count of items
func (itemList *ProjectItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length() returns the length of the item list
func (itemList *ProjectItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf returns an item or nil
func (itemList *ProjectItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *ProjectItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// FunderItemList
type FunderItemList struct {
	XMLName xml.Name `xml:"funders" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the funder item list and returns the new count of items
func (itemList *FunderItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length of item list
func (itemList *FunderItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf returns an item or nil
func (itemList *FunderItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *FunderItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// LocalGroupItemList holds the related URLs (e.g. doi, aux material doi)
type LocalGroupItemList struct {
	XMLName xml.Name `xml:"local_group" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the local group item list and returns the new count of items
func (itemList *LocalGroupItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length returns length of item list
func (itemList *LocalGroupItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf returns an item or nil
func (itemList *LocalGroupItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *LocalGroupItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// OtherNumberingSystemItemList
type OtherNumberingSystemItemList struct {
	XMLName xml.Name `xml:"other_numbering_system" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the other numbering system item list and returns the new count of items
func (itemList *OtherNumberingSystemItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length returns the length of the item
func (itemList *OtherNumberingSystemItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf return an item or nil
func (itemList *OtherNumberingSystemItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *OtherNumberingSystemItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// ContributorItemList
type ContributorItemList struct {
	XMLName xml.Name `xml:"contributors" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the contributor item list and returns the new count of items
func (itemList *ContributorItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length returns the number of items in list
func (itemList *ContributorItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf returns an item or nil
func (itemList *ContributorItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *ContributorItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// SubjectItemList
type SubjectItemList struct {
	XMLName xml.Name `xml:"subjects" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the subject item list and returns the new count of items
func (itemList *SubjectItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length returns number of items in list
func (itemList *SubjectItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf returns an item or nil
func (itemList *SubjectItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *SubjectItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// KeywordItemList
type KeywordItemList struct {
	XMLName xml.Name `xml:"keywords" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the subject item list and returns the new count of items
func (itemList *KeywordItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length returns number of items in list
func (itemList *KeywordItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf returns an item or nil
func (itemList *KeywordItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *KeywordItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// RelationItemList is an array of pointers to Item structs
type RelationItemList struct {
	XMLName xml.Name `xml:"relation" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the subject item list and returns the new count of items
func (itemList *RelationItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length returns number of items in list
func (itemList *RelationItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf returns an item or nil
func (itemList *RelationItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *RelationItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// ItemIssueItemList
type ItemIssueItemList struct {
	XMLName xml.Name `xml:"item_issues" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the issue item list and returns the new count of items
func (itemList *ItemIssueItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Lengths returns the number of items in the list
func (itemList *ItemIssueItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf returns an item or nil
func (itemList *ItemIssueItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *ItemIssueItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// CorpCreatorItemList
type CorpCreatorItemList struct {
	XMLName xml.Name `xml:"corp_creators" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the corp creator item list and returns the new count of items
func (itemList *CorpCreatorItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length returns count of items in list
func (itemList *CorpCreatorItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf returns an item or nil
func (itemList *CorpCreatorItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *CorpCreatorItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// CorpContributorItemList (not used in EPrints, but used in Invenio)
type CorpContributorItemList struct {
	XMLName xml.Name `xml:"corp_contributors" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the corp creator item list and returns the new count of items
func (itemList *CorpContributorItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length returns count of items in list
func (itemList *CorpContributorItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf returns an items or nil
func (itemList *CorpContributorItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *CorpContributorItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// ExhibitorItemList
type ExhibitorItemList struct {
	XMLName xml.Name `xml:"exhibitors" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the exhibitor item list and returns the new count of items
func (itemList *ExhibitorItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length returns count of items
func (itemList *ExhibitorItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf returns an item or nil
func (itemList *ExhibitorItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *ExhibitorItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// ProducerItemList
type ProducerItemList struct {
	XMLName xml.Name `xml:"producers" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the producer item list and returns the new count of items
func (itemList *ProducerItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length return count of items
func (itemList *ProducerItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf returns an item or nil
func (itemList *ProducerItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *ProducerItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// ConductorItemList
type ConductorItemList struct {
	XMLName xml.Name `xml:"conductors" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the conductor item list and returns the new count of items
func (itemList *ConductorItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length returns count of items
func (itemList *ConductorItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf return an item or nil
func (itemList *ConductorItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *ConductorItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// LyricistItemList
type LyricistItemList struct {
	XMLName xml.Name `xml:"lyricists" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the lyricist item list and returns the new count of items
func (itemList *LyricistItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length return count of items
func (itemList *LyricistItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf return item or nil
func (itemList *LyricistItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *LyricistItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// OptionMajorItemList
type OptionMajorItemList struct {
	XMLName xml.Name `xml:"option_major" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the option major item list and returns the new count of items
func (itemList *OptionMajorItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length return count of items
func (itemList *OptionMajorItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf return an item or nil
func (itemList *OptionMajorItemList) IndexOf(i int) *Item {
	if 0 >= i && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *OptionMajorItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// OptionMinorItemList
type OptionMinorItemList struct {
	XMLName xml.Name `xml:"option_minor" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the option minor item list and returns the new count of items
func (itemList *OptionMinorItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length return count of items
func (itemList *OptionMinorItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf return an item or nil
func (itemList *OptionMinorItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *OptionMinorItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// ThesisCommitteeItemList
type ThesisCommitteeItemList struct {
	XMLName xml.Name `xml:"thesis_committee" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the thesis committee item list and returns the new count of items
func (itemList *ThesisCommitteeItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length return count of items in list
func (itemList *ThesisCommitteeItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf return an item or nil
func (itemList *ThesisCommitteeItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *ThesisCommitteeItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// ThesisAdvisorItemList
type ThesisAdvisorItemList struct {
	XMLName xml.Name `xml:"thesis_advisor" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the thesis advisor item list and returns the new count of items
func (itemList *ThesisAdvisorItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length return count of items
func (itemList *ThesisAdvisorItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf return an item or nil
func (itemList *ThesisAdvisorItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *ThesisAdvisorItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// DivisionItemList
type DivisionItemList struct {
	XMLName xml.Name `xml:"divisions" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the division item list and returns the new count of items
func (itemList *DivisionItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length return a count of items
func (itemList *DivisionItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf returns an item or nil
func (itemList *DivisionItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *DivisionItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// RelatedPatentItemList
type RelatedPatentItemList struct {
	XMLName xml.Name `xml:"related_patents" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the related patent item list and returns the new count of items
func (itemList *RelatedPatentItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length return count of items
func (itemList *RelatedPatentItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf return an item or nil
func (itemList *RelatedPatentItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *RelatedPatentItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// PatentClassificationItemList
type PatentClassificationItemList struct {
	XMLName xml.Name `xml:"patent_classification" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the patent classification item list and returns the new count of items
func (itemList *PatentClassificationItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length return an item count
func (itemList *PatentClassificationItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf return an item or nil
func (itemList *PatentClassificationItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *PatentClassificationItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// PatentAssigneeItemList
type PatentAssigneeItemList struct {
	XMLName xml.Name `xml:"patent_assignee" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the patent assignee item list and returns the new count of items
func (itemList *PatentAssigneeItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length return item count
func (itemList *PatentAssigneeItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf return item or nil
func (itemList *PatentAssigneeItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *PatentAssigneeItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// ShelfItemList
type ShelfItemList struct {
	XMLName xml.Name `xml:"shelves" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the shelf item list and returns the new count of items
func (itemList *ShelfItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length return item count
func (itemList *ShelfItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf return item or nil
func (itemList *ShelfItemList) IndexOf(i int) *Item {
	if i >= 0 && i <= itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *ShelfItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// GScholarItemList
type GScholarItemList struct {
	XMLName xml.Name `xml:"gscholar" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the gScholar item list and returns the new count of items
func (itemList *GScholarItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length return item count
func (itemList *GScholarItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf return item or nil
func (itemList *GScholarItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *GScholarItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// AltTitleItemList
type AltTitleItemList struct {
	XMLName xml.Name `xml:"alt_title" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the altTitle item list and returns the new count of items
func (itemList *AltTitleItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length return item count
func (itemList *AltTitleItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf return item or nil
func (itemList *AltTitleItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *AltTitleItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// ConfCreatorItemList
type ConfCreatorItemList struct {
	XMLName xml.Name `xml:"conf_creators" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the confCreator item list and returns the new count of items
func (itemList *ConfCreatorItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length return item count
func (itemList *ConfCreatorItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf return item or nil
func (itemList *ConfCreatorItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *ConfCreatorItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// ReferenceItemList
type ReferenceItemList struct {
	XMLName xml.Name `xml:"reference" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the reference item list and returns the new count of items
func (itemList *ReferenceItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length return item count
func (itemList *ReferenceItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf return item or nil
func (itemList *ReferenceItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *ReferenceItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// LearningLevelItemList
type LearningLevelItemList struct {
	XMLName xml.Name `xml:"learning_level" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the learningLevel item list and returns the new count of items
func (itemList *LearningLevelItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length return item count
func (itemList *LearningLevelItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf return item or nil
func (itemList *LearningLevelItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *LearningLevelItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// CopyrightHolderItemList
type CopyrightHolderItemList struct {
	XMLName xml.Name `xml:"copyright_holders" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the copyrightHolder item list and returns the new count of items
func (itemList *CopyrightHolderItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length return item count
func (itemList *CopyrightHolderItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf return item or nil
func (itemList *CopyrightHolderItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *CopyrightHolderItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// SkillAreaItemList
type SkillAreaItemList struct {
	XMLName xml.Name `xml:"skill_areas" json:"-"`
	Items   []*Item  `xml:"item,omitempty" jsons:"item,omitempty"`
}

// Append adds an item to the skillArea item list and returns the new count of items
func (itemList *SkillAreaItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length return item count
func (itemList *SkillAreaItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf return item or nil
func (itemList *SkillAreaItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *SkillAreaItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// AccompanimentItemList
type AccompanimentItemList struct {
	XMLName xml.Name `xml:"accompaniment" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the accompaniment item list and returns the new count of items
func (itemList *AccompanimentItemList) Append(item *Item) int {
	itemList.Items = append(itemList.Items, item)
	return len(itemList.Items)
}

// Length return item count
func (itemList *AccompanimentItemList) Length() int {
	if itemList != nil {
		return len(itemList.Items)
	}
	return 0
}

// IndexOf return item or nil
func (itemList *AccompanimentItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i]
	}
	return nil
}

// SetAttributeOf at pos set item attribute return success
func (itemList *AccompanimentItemList) SetAttributeOf(i int, key string, value interface{}) bool {
	if i >= 0 && i < itemList.Length() {
		return itemList.Items[i].SetAttribute(key, value)
	}
	return false
}

// Name handles the "name" types found in Items.
type Name struct {
	XMLName    xml.Name `json:"-"`
	Family     string   `xml:"family,omitempty" json:"family,omitempty"`
	Given      string   `xml:"given,omitempty" json:"given,omitempty"`
	ID         string   `xml:"id,omitempty" json:"id,omitempty"`
	ORCID      string   `xml:"orcid,omitempty" json:"orcid,omitempty"`
	Honourific string   `xml:"honourific,omitempty" json:"honourific,omitempty"`
	Lineage    string   `xml:"lineage,omitempty" json:"lineage,omitempty"`
	Value      string   `xml:",chardata" json:"value,omitempty"`
}

func (name *Name) SetAttribute(key string, value string) bool {
	switch strings.ToLower(key) {
	case "Family":
		name.Family = value
		return true
	case "Given":
		name.Given = value
		return true
	case "id":
		name.ID = value
		return true
	case "orcid":
		name.ORCID = value
		return true
	case "honourific":
		name.Honourific = value
		return true
	case "lineage ":
		name.Lineage = value
		return true
	case "value":
		name.Value = value
		return true
	case "":
		name.Value = value
		return true
	}
	return false
}

func nameToMap(name *Name) map[string]string {
	m := map[string]string{}
	if s := strings.TrimSpace(name.Family); s != "" {
		m["family"] = s
	}
	if s := strings.TrimSpace(name.Given); s != "" {
		m["given"] = s
	}
	if s := strings.TrimSpace(name.Honourific); s != "" {
		m["honourific"] = s
	}
	if s := strings.TrimSpace(name.Lineage); s != "" {
		m["lineage"] = s
	}
	if s := strings.TrimSpace(name.Value); s != "" {
		m["value"] = s
	}
	return m
}

// MarshalJSON() is a custom JSON marshaler for Name
func (name *Name) MarshalJSON() ([]byte, error) {
	m := nameToMap(name)
	if value, flatten := m["value"]; flatten {
		return json.Marshal(value)
	}
	return json.Marshal(m)
}

// EPrintsDataSet is a struct for parsing the HTML page that returns a list of available EPrint IDs with links.
type EPrintsDataSet struct {
	XMLName xml.Name `xml:"html" json:"-"`
	Paths   []string `xml:"body>ul>li>a,omitempty" json:"paths"`
}

// MarshalJSON() renders the EPrintsDataSet HTML/XML as a list of ids
func (epds EPrintsDataSet) MarshalJSON() ([]byte, error) {
	ids := []int{}
	for _, p := range epds.Paths {
		if strings.HasSuffix(p, "/") {
			s := strings.TrimSuffix(p, "/")
			if i, err := strconv.Atoi(s); err == nil {
				ids = append(ids, i)
			}
		}
	}
	return json.Marshal(ids)
}

// String() returns a json marshaled *Name as a string
func (name *Name) String() string {
	src, err := json.Marshal(name)
	if err != nil {
		return ""
	}
	return string(src)
}

// NewEPrints returns a *EPrint with the name space set.
func NewEPrints() *EPrints {
	eprints := new(EPrints)
	eprints.XMLNS = `http://eprints.org/ep2/data/2.0`
	return eprints
}

// Append an EPrint struct to an EPrints struct returning the count of attached eprints
func (eprints *EPrints) Append(eprint *EPrint) int {
	if eprints.XMLNS == "" {
		eprints.XMLNS = `http://eprints.org/ep2/data/2.0`
	}
	eprints.EPrint = append(eprints.EPrint, eprint)
	return len(eprints.EPrint)
}

// Length returns EPrint count
func (eprints *EPrints) Length() int {
	if eprints != nil {
		return len(eprints.EPrint)
	}
	return 0
}

// IndexOf returns an EPrint or nil
func (eprints *EPrints) IndexOf(i int) *EPrint {
	if i >= 0 && i < eprints.Length() {
		return eprints.EPrint[i]
	}
	return nil
}

// File structures in Document
type File struct {
	XMLName     xml.Name `json:"-"`
	ID          string   `xml:"id,attr" json:"id"`
	FileID      int      `xml:"fileid" json:"fileid"`
	DatasetID   string   `xml:"datasetid" json:"datasetid"`
	ObjectID    int      `xml:"objectid" json:"objectid"`
	Filename    string   `xml:"filename" json:"filename"`
	MimeType    string   `xml:"mime_type" json:"mime_type"`
	Hash        string   `xml:"hash,omitempty" json:"hash,omitempty"`
	HashType    string   `xml:"hash_type,omitempty" json:"hash_type,omitempty"`
	FileSize    int      `xml:"filesize" json:"filesize"`
	MTime       string   `xml:"mtime" json:"mtime"`
	MTimeYear   int      `xml:"-" json:"-"`
	MTimeMonth  int      `xml:"-" json:"-"`
	MTimeDay    int      `xml:"-" json:"-"`
	MTimeHour   int      `xml:"-" json:"-"`
	MTimeMinute int      `xml:"-" json:"-"`
	MTimeSecond int      `xml:"-" json:"-"`
	URL         string   `xml:"url" json:"url"`

	// Additional fields found with working with our smaller repositories
	PronomID                string `xml:"pronomid,omitempty" json:"pronomID,omitempty"`
	ClassificationDateYear  int    `xml:"classification_date_year,omitempty" json:"classification_date_year,omitempty"`
	ClassificationDateMonth int    `xml:"classification_date_month,omitempty" json:"classification_date_month,omitempty"`
	ClassificationDateDay   int    `xml:"classification_date_day,omitempty" json:"classification_date_day,omitempty"`

	ClassificationDateHour int `xml:"classification_date_hour,omitempty" json:"
classification_date_hour,omitempty"`

	ClassificationDateMinute int `xml:"classification_date_minute,omitempty" json:"classification_date_minute,omitempty"`

	ClassificationDateSecond int `xml:"classification_date_second,omitempty" json:"classification_date_second,omitempty"`

	ClassificationQuality string `xml:"classification_quality,omitempty" json:"classification_quality,omitempty"`
}

// Document structures inside a Record (i.e. <eprint>...<documents><document>...</document>...</documents>...</eprint>)
type Document struct {
	XMLName          xml.Name `json:"-"`
	ID               string   `xml:"id,attr" json:"id"`
	DocID            int      `xml:"docid" json:"doc_id"`
	RevNumber        int      `xml:"rev_number" json:"rev_number,omitempty"`
	Files            []*File  `xml:"files>file" json:"files,omitempty"`
	EPrintID         int      `xml:"eprintid" json:"eprint_id"`
	Pos              int      `xml:"pos" json:"pos,omitempty"`
	Placement        int      `xml:"placement,omitempty" json:"placement,omitempty"`
	MimeType         string   `xml:"mime_type,omitempty" json:"mime_type,omitempty"`
	Format           string   `xml:"format" json:"format"`
	FormatDesc       string   `xml:"formatdesc,omitempty" json:"format_desc,omitempty"`
	Language         string   `xml:"language,omitempty" json:"language,omitempty"`
	Security         string   `xml:"security" json:"security"`
	License          string   `xml:"license" json:"license"`
	Main             string   `xml:"main" json:"main"`
	DateEmbargo      string   `xml:"date_embargo,omitempty" json:"date_embargo,omitempty"`
	DateEmbargoYear  int      `xml:"-" json:"-"`
	DateEmbargoMonth int      `xml:"-" json:"-"`
	DateEmbargoDay   int      `xml:"-" json:"-"`

	MediaDuration    string `xml:"media_duration,omitempty" json:"media_duration,omitempty"`
	MediaAudioCodec  string `xml:"media_audio_codec,omitempty" json:"media_audio_code,omitempty"`
	MediaVideoCodec  string `xml:"media_video_codec,omitempty" json:"media_video_code,omitempty"`
	MediaWidth       int    `xml:"media_width,omitempty" json:"media_width,omitempty"`
	MediaHeight      int    `xml:"media_height,omitempty" json:"media_height,omitempty"`
	MediaAspectRatio string `xml:"media_aspect_ratio,omitempty" json:"media_aspect_ratio,omitempty"`
	MediaSampleStart string `xml:"media_sample_start,omitempty" json:"media_sample_start,omitempty"`
	MediaSampleStop  string `xml:"media_sample_stop,omitempty" json:"media_sample_stop,omitempty"`

	Content  string            `xml:"content,omitempty" json:"content,omitempty"`
	Relation *RelationItemList `xml:"relation,omitempty" json:"relation,omitempty"`
}

// DocumentList is an array of pointers to Document structs
type DocumentList []*Document

// Append adds a document to the documents list and returns the new count of items
func (documentList *DocumentList) Append(document *Document) int {
	*documentList = append(*documentList, document)
	return len(*documentList)
}

// Length returns the length of DocumentList
func (documentList *DocumentList) Length() int {
	if documentList != nil {
		return len(*documentList)
	}
	return 0
}

// GetDocument takes a position (zero based) and returns the *Document
// in the DocumentList.
func (documentList DocumentList) IndexOf(i int) *Document {
	if i < 0 || i >= len(documentList) {
		return nil
	}
	return documentList[i]
}

// ePrintIDs is a struct for parsing the ids page of EPrints REST API
type ePrintIDs struct {
	XMLName xml.Name `xml:"html" json:"-"`
	IDs     []string `xml:"body>ul>li>a" json:"ids"`
}

// SyntheticFields renders analyzes an EPrint object
// and populates or updates any synthetic fields like
// primary_object and related_object.
func (e *EPrint) SyntheticFields() {
	// Render PrimaryObject and RelatedObjects fields
	e.PrimaryObject = make(map[string]interface{})
	e.RelatedObjects = []map[string]interface{}{}
	if e.Documents != nil {
		docCnt := e.Documents.Length()
		for i := 0; i < docCnt; i++ {
			doc := e.Documents.IndexOf(i)
			if doc.Security == "public" && doc.Main != "indexcodes.txt" {
				obj := make(map[string]interface{})
				obj["basename"] = doc.Main
				obj["url"] = fmt.Sprintf("%s/%d/%s", strings.Replace(e.ID, "/id/eprint", "", 1), doc.Pos, doc.Main)
				obj["mime_type"] = doc.MimeType
				obj["content"] = doc.Content
				obj["license"] = doc.License
				if doc.Files != nil {
					for _, fObj := range doc.Files {
						if fObj.Filename == doc.Main {
							obj["filesize"] = fObj.FileSize
							break
						}
					}
				}
				obj["version"] = fmt.Sprintf("v%d.0.0", doc.RevNumber)
				if (doc.Placement == 1 || doc.Pos == 1) && doc.Content != "supplemental" {
					e.PrimaryObject = obj
				} else {
					e.RelatedObjects = append(e.RelatedObjects, obj)
				}
			}
		}
	}
}

// EPrintUser is a struct for representing a user in a EPrint repository.
// NOTE: it does not represent all user fields and attributes.
type EPrintUser struct {
	XMLName   xml.Name `xml:"user" json:"-"`
	UserID    int      `xml:"userid" json:"userid"`
	Username  string   `xml:"username" json:"username"`
	Type      string   `xml:"type" json:"type"`
	Name      *Name    `xml:"name,omitempty" json:"name,omitempty"`
	EMail     string   `xml:"email,omitempty" json:"email,omitempty"`
	HideEMail bool     `xml:"hideemail,omitempty" json:"hideemail,omitempty"`
	Joined    string   `xml:"joined,omitempty" json:"joined,omitempty"`
	Dept      string   `xml:"dept,omitempty" json:"dept,omitempty"`
	Org       string   `xml:"org,omitempty" json:"org,omitempty"`
	Address   string   `xml:"address,omitempty" json:"address,omitempty"`
	Country   string   `xml:"country,omitempty" json:"country,omitempty"`
}

// SetDefaults sets the default values for DefaultCollection, DefaultRights, Default
func SetDefaults(collection string, rights string, officialURL string, refereed string, status string) {
	DefaultCollection = collection
	DefaultRights = rights
	DefaultOfficialURL = officialURL
	DefaultRefereed = refereed
	DefaultStatus = status
}

// GenerateIDNumber generates a unique ID number based on the
// instance of generation and the default collection name.
// I.e. COLLECTION_NAME:DATESTAMP-NANOSECOND
//
// NOTE: You need to have set the values for DefaultCollection and
// DefaultOfficialURL before calling this function.
func GenerateIDNumber(eprint *EPrint) string {
	collection := DefaultCollection
	if eprint.Collection != "" {
		collection = eprint.Collection
	}
	now := time.Now()
	return fmt.Sprintf(`%s/%s:%s-%d`, DefaultOfficialURL, collection, now.Format("20060102"), now.Nanosecond())
}

// GenerateImportID generates a unique ID number based on the
// instance of generation and the default collection name.
// I.e. PREFIX:DATESTAMP-NANOSECOND.USER_ID
//
// NOTE: You need to have set the values for DefaultCollection and
// DefaultOfficialURL before calling this function.
func GenerateImportID(prefix string, eprint *EPrint) string {
	now := time.Now()
	return fmt.Sprintf(`%s:%s-%d.%d`, prefix, now.Format("20060102"), now.Nanosecond(), eprint.UserID)
}

// GenerateOfficialURL generates an OfficalURL (i.e.
//   idNumber string appended to OfficialURLPrefix)
func GenerateOfficialURL(eprint *EPrint) string {
	/* IDNumber and OfficialURL (resolver URL) are the same value. */
	idNumber := eprint.IDNumber
	if idNumber == "" {
		idNumber = GenerateIDNumber(eprint)
	}
	return fmt.Sprintf(`%s`, idNumber)
}
