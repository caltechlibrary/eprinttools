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
	XMLName              xml.Name                      `json:"-"`
	ID                   string                        `xml:"id,attr,omitempty" json:"id,omitempty"`
	EPrintID             int                           `xml:"eprintid,omitempty" json:"eprint_id,omitempty"`
	RevNumber            int                           `xml:"rev_number,omitempty" json:"rev_number,omitempty"`
	Documents            *DocumentList                 `xml:"documents>document,omitempty" json:"documents,omitempty"`
	EPrintStatus         string                        `xml:"eprint_status,omitempty" json:"eprint_status,omitempty"`
	UserID               int                           `xml:"userid,omitempty" json:"userid,omitempty"`
	Dir                  string                        `xml:"dir,omitempty" json:"dir,omitempty"`
	Datestamp            string                        `xml:"datestamp,omitempty" json:"datestamp,omitempty"`
	DatestampYear        int                           `xml:"-" json:"-"`
	DatestampMonth       int                           `xml:"-" json:"-"`
	DatestampDay         int                           `xml:"-" json:"-"`
	DatestampHour        int                           `xml:"-" json:"-"`
	DatestampMinute      int                           `xml:"-" json:"-"`
	DatestampSecond      int                           `xml:"-" json:"-"`
	LastModified         string                        `xml:"lastmod,omitempty" json:"lastmod,omitempty"`
	LastModifiedYear     int                           `xml:"-" json:"-"`
	LastModifiedMonth    int                           `xml:"-" json:"-"`
	LastModifiedDay      int                           `xml:"-" json:"-"`
	LastModifiedHour     int                           `xml:"-" json:"-"`
	LastModifiedMinute   int                           `xml:"-" json:"-"`
	LastModifiedSecond   int                           `xml:"-" json:"-"`
	StatusChanged        string                        `xml:"status_changed,omitempty" json:"status_changed,omitempty"`
	StatusChangedYear    int                           `xml:"-" json:"-"`
	StatusChangedMonth   int                           `xml:"-" json:"-"`
	StatusChangedDay     int                           `xml:"-" json:"-"`
	StatusChangedHour    int                           `xml:"-" json:"-"`
	StatusChangedMinute  int                           `xml:"-" json:"-"`
	StatusChangedSecond  int                           `xml:"-" json:"-"`
	Type                 string                        `xml:"type,omitempty" json:"type,omitempty"`
	MetadataVisibility   string                        `xml:"metadata_visibility,omitempty" json:"metadata_visibility,omitempty"`
	Creators             *CreatorItemList              `xml:"creators,omitempty" json:"creators,omitempty"`
	Title                string                        `xml:"title,omitempty" json:"title,omitempty"`
	IsPublished          string                        `xml:"ispublished,omitempty" json:"ispublished,omitempty"`
	FullTextStatus       string                        `xml:"full_text_status,omitempty" json:"full_text_status,omitempty"`
	Keywords             string                        `xml:"keywords,omitempty" json:"keywords,omitempty"`
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
	Pages                string                        `xml:"pages,omitempty" json:"pages,omitempty"`
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

	// NOTE: Misc fields discoverd exploring REST API records, not currently used at Caltech Library (RSD, 2018-01-02)
	Subjects           *SubjectItemList         `xml:"subjects,omitempty" json:"subjects,omitempty"`
	PresType           string                   `xml:"pres_type,omitempty" json:"presentation_type,omitempty"`
	Succeeds           string                   `xml:"succeeds,omitempty" json:"succeeds,omitempty"`
	Commentary         string                   `xml:"commentary,omitempty" json:"commentary,omitempty"`
	ContactEMail       string                   `xml:"contact_email,omitempty" json:"contect_email,omitempty"`
	FileInfo           string                   `xml:"fileinfo,omitempty" json:"file_info,omitempty"`
	Latitude           string                   `xml:"latitude,omitempty" json:"latitude,omitempty"`
	Longitude          string                   `xml:"longitude,omitempty" json:"longitude,omitempty"`
	ItemIssues         *ItemIssueItemList       `xml:"item_issues,omitempty" json:"item_issues,omitempty"`
	ItemIssuesCount    int                      `xml:"item_issues_count,omitempty" json:"item_issues_count,omitempty"`
	CorpCreators       *CorpCreatorItemList     `xml:"corp_creators,omitempty" json:"corp_creators,omitempty"`
	CorpContributors   *CorpContributorItemList `xml:"corp_contributors,omitempty" json:"corp_contributors,omitempty"`
	Department         string                   `xml:"department,omitempty" json:"department,omitempty"`
	OutputMedia        string                   `xml:"output_media,omitempty" json:"output_media,omitempty"`
	Exhibitors         *ExhibitorItemList       `xml:"exhibitors,omitempty" json:"exhibitors,omitempty"`
	NumPieces          string                   `xml:"num_pieces,omitempty" json:"num_pieces,omitempty"`
	CompositionType    string                   `xml:"composition_type,omitempty" json:"composition_type,omitempty"`
	Producers          *ProducerItemList        `xml:"producers,omitempty" json:"producers,omitempty"`
	Conductors         *ConductorItemList       `xml:"conductors,omitempty" json:"conductors,omitempty"`
	Lyricists          *LyricistItemList        `xml:"lyricists,omitempty" json:"lyricists,omitempty"`
	Accompaniment      *AccompanimentItemList   `xml:"accompaniment,omitempty" json:"accompaniment,omitempty"`
	DataType           string                   `xml:"data_type,omitempty" json:"data_type,omitempty"`
	PedagogicType      string                   `xml:"pedagogic_type,omitempty" json:"pedagogic_type,omitempty"`
	CompletionTime     string                   `xml:"completion_time,omitempty" json:"completion_time,omitempty"`
	TaskPurpose        string                   `xml:"task_purpose,omitempty" json:"task_purpose,omitempty"`
	SkillAreas         *SkillAreaItemList       `xml:"skill_areas,omitempty" json:"skill_areas,omitempty"`
	CopyrightHolders   *CopyrightHolderItemList `xml:"copyright_holders,omitempty" json:"copyright_holders,omitempty"`
	LearningLevelText  string                   `xml:"learning_level,omitempty" json:"learning_level,omitempty"`
	LearningLevel      *LearningLevelItemList   `xml:"-" json:"-"`
	DOI                string                   `xml:"doi,omitempty" json:"doi,omitempty"`
	PMCID              string                   `xml:"pmc_id,omitempty" json:"pmcid,omitempty"`
	PMID               string                   `xml:"pmid,omitempty" json:"pmid,omitempty"`
	ParentURL          string                   `xml:"parent_url,omitempty" json:"parent_url,omitempty"`
	Reference          *ReferenceItemList       `xml:"reference,omitempty" json:"reference,omitempty"`
	ConfCreators       *ConfCreatorItemList     `xml:"conf_creators,omitempty" json:"conf_creators,omitempty"`
	AltTitle           *AltTitleItemList        `xml:"alt_title,omitempty" json:"alt_title,omitempty"`
	TOC                string                   `xml:"toc,omitempty" json:"toc,omitempty"`
	Interviewer        string                   `xml:"interviewer,omitempty" json:"interviewer,omitempty"`
	InterviewDate      string                   `xml:"interviewdate,omitempty" json:"interviewdate,omitempty"`
	GScholar           *GScholarItemList        `xml:"gscholar,omitempty" json:"gscholar,omitempty"`
	NonSubjKeywords    string                   `xml:"nonsubj_keywords,omitempty" json:"nonsubj_keywords,omitempty"`
	Season             string                   `xml:"season,omitempty" json:"season,omitempty"`
	ClassificationCode string                   `xml:"classification_code,omitempty" json:"classification_code,omitempty"`
	Shelves            *ShelfItemList           `xml:"shelves,omitempty" json:"shelves,omitempty"`

	// NOTE: Sword deposit fields
	SwordDepository string `xml:"sword_depository,omitempty" json:"sword_depository,omitempty"`
	SwordDepositor  string `xml:"sword_depositor,omitempty" json:"sword_depositor,omitempty"`
	SwordSlug       string `xml:"sword_slug,omitempty" json:"sword_slug,omitempty"`
	ImportID        int    `xml:"importid,omitempty" json:"import_id,omitempty"`

	// Patent related fields
	PatentApplicant          string                        `xml:"patent_applicant,omitempty" json:"patent_applicant,omitempty"`
	PatentNumber             string                        `xml:"patent_number,omitempty" json:"patent_number,omitempty"`
	PatentAssignee           *PatentAssigneeItemList       `xml:"patent_assignee,omitempty" json:"patent_assignee,omitempty"`
	PatentClassificationText string                        `xml:"-" json:"-"`
	PatentClassification     *PatentClassificationItemList `xml:"patent_classification,omitempty" json:"patent_classification,omitempty"`
	RelatedPatents           *RelatedPatentItemList        `xml:"related_patents,omitempty" json:"related_patents,omitempty"`

	// Thesis oriented fields
	Divisions                   *DivisionItemList        `xml:"divisions,omitemmpty" json:"divisions,omitempty"`
	Institution                 string                   `xml:"institution,omitempty" json:"institution,omitempty"`
	ThesisType                  string                   `xml:"thesis_type,omitempty" json:"thesis_type,omitempty"`
	ThesisAdvisor               *ThesisAdvisorItemList   `xml:"thesis_advisor,omitempty" json:"thesis_advisor,omitempty"`
	ThesisCommittee             *ThesisCommitteeItemList `xml:"thesis_committee,omitempty" json:"thesis_committee,omitempty"`
	ThesisDegree                string                   `xml:"thesis_degree,omitempty" json:"thesis_degree,omitempty"`
	ThesisDegreeGrantor         string                   `xml:"thesis_degree_grantor,omitempty" json:"thesis_degree_grantor,omitempty"`
	ThesisDegreeDate            string                   `xml:"thesis_degree_date,omitempty" json:"thesis_degree_date,omitempty"`
	ThesisDegreeDateYear        int                      `xml:"-" json:"-"`
	ThesisDegreeDateMonth       int                      `xml:"-" json:"-"`
	ThesisDegreeDateDay         int                      `xml:"-" json:"-"`
	ThesisSubmittedDate         string                   `xml:"thesis_submit_date,omitempty" json:"thesis_submit_date,omitempty"`
	ThesisSubmittedDateYear     int                      `xml:"-" json:"-"`
	ThesisSubmittedDateMonth    int                      `xml:"-" json:"-"`
	ThesisSubmittedDateDay      int                      `xml:"-" json:"-"`
	ThesisDefenseDate           string                   `xml:"thesis_defense_date,omitempty" json:"thesis_defense_date,omitempty"`
	ThesisDefenseDateYear       int                      `xml:"-" json:"-"`
	ThesisDefenseDateMonth      int                      `xml:"-" json:"-"`
	ThesisDefenseDateDay        int                      `xml:"-" json:"-"`
	ThesisApprovedDate          string                   `xml:"thesis_approved_date,omitempty" json:"thesis_approved_date,omitempty"`
	ThesisApprovedDateYear      int                      `xml:"-" json:"-"`
	ThesisApprovedDateMonth     int                      `xml:"-" json:"-"`
	ThesisApprovedDateDay       int                      `xml:"-" json:"-"`
	ThesisPublicDate            string                   `xml:"thesis_public_date,omitempty" json:"thesis_public_date,omitempty"`
	ThesisPublicDateYear        int                      `xml:"-" json:"-"`
	ThesisPublicDateMonth       int                      `xml:"-" json:"-"`
	ThesisPublicDateDay         int                      `xml:"-" json:"-"`
	ThesisAuthorEMail           string                   `xml:"thesis_author_email,omitempty" json:"thesis_author_email,omitempty"`
	HideThesisAuthorEMail       string                   `xml:"hide_thesis_author_email,omitempty" json:"hide_thesis_author_email,omitempty"`
	GradOfficeApprovalDate      string                   `xml:"gradofc_approval_date,omitempty" json:"gradofc_approval_date,omitempty"`
	GradOfficeApprovalDateYear  int                      `xml:"-" json:"-"`
	GradOfficeApprovalDateMonth int                      `xml:"-" json:"-"`
	GradOfficeApprovalDateDay   int                      `xml:"-" json:"-"`
	ThesisAwards                string                   `xml:"thesis_awards,omitempty" json:"thesis_awards,omitempty"`
	ReviewStatus                string                   `xml:"review_status,omitempty" json:"review_status,omitempty"`
	OptionMajor                 *OptionMajorItemList     `xml:"option_major,omitempty" json:"option_major,omitempty"`
	OptionMinor                 *OptionMinorItemList     `xml:"option_minor,omitempty" json:"option_major,omitempty"`
	CopyrightStatement          string                   `xml:"copyright_statement,omitempty" json:"copyright_statement,omitempty"`

	// Custom fields from some EPrints repositories
	Source     string `xml:"source,omitempty" json:"source,omitempty"`
	ReplacedBy int    `xml:"replacedby,omitempty" json:"replacedby,omitempty"`

	// Edit Control Fields
	EditLockUser  int `xml:"-" json:"-"`
	EditLockSince int `xml:"-" json:"-"`
	EditLockUntil int `xml:"-" json:"-"`

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
	err := json.Unmarshal(src, &m)
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
	//SetAttributeOf(int, string, string) bool
}

// CreatorItemList holds a list of authors
type CreatorItemList struct {
	XMLName xml.Name `xml:"creators" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the Creator list and returns the new count of items
func (creatorItemList *CreatorItemList) Append(item *Item) int {
	creatorItemList.Items = append(creatorItemList.Items, item)
	return len(creatorItemList.Items)
}

// Length() returns the numnber of items in the list
func (creatorItemList *CreatorItemList) Length() int {
	return len(creatorItemList.Items)
}

// IndexOf() returns an item in the list or nil
func (creatorItemList *CreatorItemList) IndexOf(i int) *Item {
	if i >= 0 && i < creatorItemList.Length() {
		return creatorItemList.Items[i]
	}
	return nil
}

// EditorItemList holds a list of editors
type EditorItemList struct {
	XMLName xml.Name `xml:"editors" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the Editor item list and returns the new count of items
func (editorItemList *EditorItemList) Append(item *Item) int {
	editorItemList.Items = append(editorItemList.Items, item)
	return len(editorItemList.Items)
}

// Length() returns the number of items in the list
func (editorItemList *EditorItemList) Length() int {
	return len(editorItemList.Items)
}

// IndexOf() returns an item in the list or nil
func (editorItemList *EditorItemList) IndexOf(i int) *Item {
	if i >= 0 && i < editorItemList.Length() {
		return editorItemList.Items[i]
	}
	return nil
}

// RelatedURLItemList holds the related URLs (e.g. doi, aux material doi)
type RelatedURLItemList struct {
	XMLName xml.Name `xml:"related_url" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the related url item list and returns the new count of items
func (relatedURLItemList *RelatedURLItemList) Append(item *Item) int {
	relatedURLItemList.Items = append(relatedURLItemList.Items, item)
	return len(relatedURLItemList.Items)
}

// Length() returns the number of items in the list
func (relatedURLItemList *RelatedURLItemList) Length() int {
	return len(relatedURLItemList.Items)
}

// IndexOf() returns an item in the list or nil
func (relatedURLItemList *RelatedURLItemList) IndexOf(i int) *Item {
	if i >= 0 && i < relatedURLItemList.Length() {
		return relatedURLItemList.Items[i]
	}
	return nil
}

// ReferenceTextItemList
type ReferenceTextItemList struct {
	XMLName xml.Name `xml:"referencetext" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the reference text url item list and returns the new count of items
func (referenceTextItemList *ReferenceTextItemList) Append(item *Item) int {
	referenceTextItemList.Items = append(referenceTextItemList.Items, item)
	return len(referenceTextItemList.Items)
}

// Length returns the length of an ReferenceTextItemList
func (referenceTextItemList *ReferenceTextItemList) Length() int {
	return len(referenceTextItemList.Items)
}

// IndexOf returns an Item or nil
func (referenceTextItemList *ReferenceTextItemList) IndexOf(i int) *Item {
	if i >= 0 && i < referenceTextItemList.Length() {
		return referenceTextItemList.Items[i]
	}
	return nil
}

// UnmarshJSON takes a reference text list of item and returns
// an appropriately values to assigned struct.
func (referenceTextItemList *ReferenceTextItemList) UnmarshalJSON(src []byte) error {
	var values []string

	m := make(map[string][]interface{})
	err := json.Unmarshal(src, &m)
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
		item.Value = value
		referenceTextItemList.Append(item)
	}
	return err
}

// ProjectItemList
type ProjectItemList struct {
	XMLName xml.Name `xml:"projects" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the project item list and returns the new count of items
func (projectItemList *ProjectItemList) Append(item *Item) int {
	projectItemList.Items = append(projectItemList.Items, item)
	return len(projectItemList.Items)
}

// Length() returns the length of the item list
func (projectItemList *ProjectItemList) Length() int {
	return len(projectItemList.Items)
}

// IndexOf returns an item or nil
func (projectItemList *ProjectItemList) IndexOf(i int) *Item {
	if i >= 0 && i < projectItemList.Length() {
		return projectItemList.Items[i]
	}
	return nil
}

// FunderItemList
type FunderItemList struct {
	XMLName xml.Name `xml:"funders" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the funder item list and returns the new count of items
func (funderItemList *FunderItemList) Append(item *Item) int {
	funderItemList.Items = append(funderItemList.Items, item)
	return len(funderItemList.Items)
}

// Length of item list
func (funderItemList *FunderItemList) Length() int {
	return len(funderItemList.Items)
}

// IndexOf returns an item or nil
func (funderItemList *FunderItemList) IndexOf(i int) *Item {
	if i >= 0 && i < funderItemList.Length() {
		return funderItemList.Items[i]
	}
	return nil
}

// LocalGroupItemList holds the related URLs (e.g. doi, aux material doi)
type LocalGroupItemList struct {
	XMLName xml.Name `xml:"local_group" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the local group item list and returns the new count of items
func (localGroupItemList *LocalGroupItemList) Append(item *Item) int {
	localGroupItemList.Items = append(localGroupItemList.Items, item)
	return len(localGroupItemList.Items)
}

// Length returns length of item list
func (localGroupItemList *LocalGroupItemList) Length() int {
	return len(localGroupItemList.Items)
}

// IndexOf returns an item or nil
func (localGroupItemList *LocalGroupItemList) IndexOf(i int) *Item {
	if i >= 0 && i < localGroupItemList.Length() {
		return localGroupItemList.Items[i]
	}
	return nil
}

// OtherNumberingSystemItemList
type OtherNumberingSystemItemList struct {
	XMLName xml.Name `xml:"other_numbering_system" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the other numbering system item list and returns the new count of items
func (otherNumberingSystemItemList *OtherNumberingSystemItemList) Append(item *Item) int {
	otherNumberingSystemItemList.Items = append(otherNumberingSystemItemList.Items, item)
	return len(otherNumberingSystemItemList.Items)
}

// Length returns the length of the item
func (otherNumberingSystemItemList *OtherNumberingSystemItemList) Length() int {
	return len(otherNumberingSystemItemList.Items)
}

// IndexOf return an item or nil
func (otherNumberingSystemItemList *OtherNumberingSystemItemList) IndexOf(i int) *Item {
	if i >= 0 && i < otherNumberingSystemItemList.Length() {
		return otherNumberingSystemItemList.Items[i]
	}
	return nil
}

// ContributorItemList
type ContributorItemList struct {
	XMLName xml.Name `xml:"contributors" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the contributor item list and returns the new count of items
func (contributorItemList *ContributorItemList) Append(item *Item) int {
	contributorItemList.Items = append(contributorItemList.Items, item)
	return len(contributorItemList.Items)
}

// Length returns the number of items in list
func (contributorItemList *ContributorItemList) Length() int {
	return len(contributorItemList.Items)
}

// IndexOf returns an item or nil
func (contributorItemList *ContributorItemList) IndexOf(i int) *Item {
	if i >= 0 && i < contributorItemList.Length() {
		return contributorItemList.Items[i]
	}
	return nil
}

// SubjectItemList
type SubjectItemList struct {
	XMLName xml.Name `xml:"subjects" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the subject item list and returns the new count of items
func (subjectItemList *SubjectItemList) Append(item *Item) int {
	subjectItemList.Items = append(subjectItemList.Items, item)
	return len(subjectItemList.Items)
}

// Length returns number of items in list
func (subjectItemList *SubjectItemList) Length() int {
	return len(subjectItemList.Items)
}

// IndexOf returns an item or nil
func (subjectItemList *SubjectItemList) IndexOf(i int) *Item {
	if i >= 0 && i < subjectItemList.Length() {
		return subjectItemList.Items[i]
	}
	return nil
}

// ItemIssueItemList
type ItemIssueItemList struct {
	XMLName xml.Name `xml:"item_issues" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the issue item list and returns the new count of items
func (issueItemList *ItemIssueItemList) Append(item *Item) int {
	issueItemList.Items = append(issueItemList.Items, item)
	return len(issueItemList.Items)
}

// Lengths returns the number of items in the list
func (issueItemList *ItemIssueItemList) Length() int {
	return len(issueItemList.Items)
}

// IndexOf returns an item or nil
func (issueItemList *ItemIssueItemList) IndexOf(i int) *Item {
	if i >= 0 && i < issueItemList.Length() {
		return issueItemList.Items[i]
	}
	return nil
}

// CorpCreatorItemList
type CorpCreatorItemList struct {
	XMLName xml.Name `json:"-"` //`xml:"corp_creators" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the corp creator item list and returns the new count of items
func (corpCreatorItemList *CorpCreatorItemList) Append(item *Item) int {
	corpCreatorItemList.Items = append(corpCreatorItemList.Items, item)
	return len(corpCreatorItemList.Items)
}

// Length returns count of items in list
func (corpCreatorItemList *CorpCreatorItemList) Length() int {
	return len(corpCreatorItemList.Items)
}

// IndexOf returns an item or nil
func (corpCreatorItemList *CorpCreatorItemList) IndexOf(i int) *Item {
	if i >= 0 && i < corpCreatorItemList.Length() {
		return corpCreatorItemList.Items[i]
	}
	return nil
}

// CorpContributorItemList (not used in EPrints, but used in Invenio)
type CorpContributorItemList struct {
	XMLName xml.Name `json:"-"` //`xml:"corp_contributors" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the corp creator item list and returns the new count of items
func (corpContributorItemList *CorpContributorItemList) Append(item *Item) int {
	corpContributorItemList.Items = append(corpContributorItemList.Items, item)
	return len(corpContributorItemList.Items)
}

// Length returns count of items in list
func (corpContributorItemList *CorpContributorItemList) Length() int {
	return len(corpContributorItemList.Items)
}

// IndexOf returns an items or nil
func (corpContributorItemList *CorpContributorItemList) IndexOf(i int) *Item {
	if i >= 0 && i < corpContributorItemList.Length() {
		return corpContributorItemList.Items[i]
	}
	return nil
}

// ExhibitorItemList
type ExhibitorItemList struct {
	XMLName xml.Name `xml:"exhibitors" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the exhibitor item list and returns the new count of items
func (exhibitorItemList *ExhibitorItemList) Append(item *Item) int {
	exhibitorItemList.Items = append(exhibitorItemList.Items, item)
	return len(exhibitorItemList.Items)
}

// Length returns count of items
func (exhibitorItemList *ExhibitorItemList) Length() int {
	return len(exhibitorItemList.Items)
}

// IndexOf returns an item or nil
func (exhibitorItemList *ExhibitorItemList) IndexOf(i int) *Item {
	if i >= 0 && i < exhibitorItemList.Length() {
		return exhibitorItemList.Items[i]
	}
	return nil
}

// ProducerItemList
type ProducerItemList struct {
	XMLName xml.Name `xml:"producers" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the producer item list and returns the new count of items
func (producerItemList *ProducerItemList) Append(item *Item) int {
	producerItemList.Items = append(producerItemList.Items, item)
	return len(producerItemList.Items)
}

// Length return count of items
func (producerItemList *ProducerItemList) Length() int {
	return len(producerItemList.Items)
}

// IndexOf returns an item or nil
func (producerItemList *ProducerItemList) IndexOf(i int) *Item {
	if i >= 0 && i < producerItemList.Length() {
		return producerItemList.Items[i]
	}
	return nil
}

// ConductorItemList
type ConductorItemList struct {
	XMLName xml.Name `xml:"conductors" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the conductor item list and returns the new count of items
func (conductorItemList *ConductorItemList) Append(item *Item) int {
	conductorItemList.Items = append(conductorItemList.Items, item)
	return len(conductorItemList.Items)
}

// Length returns count of items
func (conductorItemList *ConductorItemList) Length() int {
	return len(conductorItemList.Items)
}

// IndexOf return an item or nil
func (conductorItemList *ConductorItemList) IndexOf(i int) *Item {
	if i >= 0 && i < conductorItemList.Length() {
		return conductorItemList.Items[i]
	}
	return nil
}

// LyricistItemList
type LyricistItemList struct {
	XMLName xml.Name `xml:"lyricists" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the lyricist item list and returns the new count of items
func (lyricistItemList *LyricistItemList) Append(item *Item) int {
	lyricistItemList.Items = append(lyricistItemList.Items, item)
	return len(lyricistItemList.Items)
}

// Length return count of items
func (lyricistItemList *LyricistItemList) Length() int {
	return len(lyricistItemList.Items)
}

// IndexOf return item or nil
func (lyricistItemList *LyricistItemList) IndexOf(i int) *Item {
	if i >= 0 && i < lyricistItemList.Length() {
		return lyricistItemList.Items[i]
	}
	return nil
}

// OptionMajorItemList
type OptionMajorItemList struct {
	XMLName xml.Name `xml:"option_major" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the option major item list and returns the new count of items
func (optionMajorItemList *OptionMajorItemList) Append(item *Item) int {
	optionMajorItemList.Items = append(optionMajorItemList.Items, item)
	return len(optionMajorItemList.Items)
}

// Length return count of items
func (optionMajorItemList *OptionMajorItemList) Length() int {
	return len(optionMajorItemList.Items)
}

// IndexOf return an item or nil
func (optionMajorItemList *OptionMajorItemList) IndexOf(i int) *Item {
	if 0 >= i && i < optionMajorItemList.Length() {
		return optionMajorItemList.Items[i]
	}
	return nil
}

// OptionMinorItemList
type OptionMinorItemList struct {
	XMLName xml.Name `xml:"option_minor" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the option minor item list and returns the new count of items
func (optionMinorItemList *OptionMinorItemList) Append(item *Item) int {
	optionMinorItemList.Items = append(optionMinorItemList.Items, item)
	return len(optionMinorItemList.Items)
}

// Length return count of items
func (optionMinorItemList *OptionMinorItemList) Length() int {
	return len(optionMinorItemList.Items)
}

// IndexOf return an item or nil
func (optionMinorItemList *OptionMinorItemList) IndexOf(i int) *Item {
	if i >= 0 && i < optionMinorItemList.Length() {
		return optionMinorItemList.Items[i]
	}

	return nil
}

// ThesisCommitteeItemList
type ThesisCommitteeItemList struct {
	XMLName xml.Name `xml:"thesis_committee" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the thesis committee item list and returns the new count of items
func (thesisCommitteeItemList *ThesisCommitteeItemList) Append(item *Item) int {
	thesisCommitteeItemList.Items = append(thesisCommitteeItemList.Items, item)
	return len(thesisCommitteeItemList.Items)
}

// Length return count of items in list
func (thesisCommitteeItemList *ThesisCommitteeItemList) Length() int {
	return len(thesisCommitteeItemList.Items)
}

// IndexOf return an item or nil
func (thesisCommitteeItemList *ThesisCommitteeItemList) IndexOf(i int) *Item {
	if i >= 0 && i < thesisCommitteeItemList.Length() {
		return thesisCommitteeItemList.Items[i]
	}
	return nil
}

// ThesisAdvisorItemList
type ThesisAdvisorItemList struct {
	XMLName xml.Name `xml:"thesis_advisor" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the thesis advisor item list and returns the new count of items
func (thesisAdvisorItemList *ThesisAdvisorItemList) Append(item *Item) int {
	thesisAdvisorItemList.Items = append(thesisAdvisorItemList.Items, item)
	return len(thesisAdvisorItemList.Items)
}

// Length return count of items
func (thesisAdvisorItemList *ThesisAdvisorItemList) Length() int {
	return len(thesisAdvisorItemList.Items)
}

// IndexOf return an item or nil
func (thesisAdvisorItemList *ThesisAdvisorItemList) IndexOf(i int) *Item {
	if i >= 0 && i < thesisAdvisorItemList.Length() {
		return thesisAdvisorItemList.Items[i]
	}
	return nil
}

// DivisionItemList
type DivisionItemList struct {
	XMLName xml.Name `xml:"divisions" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the division item list and returns the new count of items
func (divisionItemList *DivisionItemList) Append(item *Item) int {
	divisionItemList.Items = append(divisionItemList.Items, item)
	return len(divisionItemList.Items)
}

// Length return a count of items
func (divisionItemList *DivisionItemList) Length() int {
	return len(divisionItemList.Items)
}

// IndexOf returns an item or nil
func (divisionItemList *DivisionItemList) IndexOf(i int) *Item {
	if i >= 0 && i < divisionItemList.Length() {
		return divisionItemList.Items[i]
	}
	return nil
}

// RelatedPatentItemList
type RelatedPatentItemList struct {
	XMLName xml.Name `xml:"related_patents" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the related patent item list and returns the new count of items
func (relatedPatentItemList *RelatedPatentItemList) Append(item *Item) int {
	relatedPatentItemList.Items = append(relatedPatentItemList.Items, item)
	return len(relatedPatentItemList.Items)
}

// Length return count of items
func (relatedPatentItemList *RelatedPatentItemList) Length() int {
	return len(relatedPatentItemList.Items)
}

// IndexOf return an item or nil
func (relatedPatentItemList *RelatedPatentItemList) IndexOf(i int) *Item {
	if i >= 0 && i < relatedPatentItemList.Length() {
		return relatedPatentItemList.Items[i]
	}
	return nil
}

// PatentClassificationItemList
type PatentClassificationItemList struct {
	XMLName xml.Name `xml:"patent_classification" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the patent classification item list and returns the new count of items
func (patentClassificationItemList *PatentClassificationItemList) Append(item *Item) int {
	patentClassificationItemList.Items = append(patentClassificationItemList.Items, item)
	return len(patentClassificationItemList.Items)
}

// Length return an item count
func (patentClassificationItemList *PatentClassificationItemList) Length() int {
	return len(patentClassificationItemList.Items)
}

// IndexOf return an item or nil
func (patentClassificationItemList *PatentClassificationItemList) IndexOf(i int) *Item {
	if i >= 0 && i < patentClassificationItemList.Length() {
		return patentClassificationItemList.Items[i]
	}
	return nil
}

// PatentAssigneeItemList
type PatentAssigneeItemList struct {
	XMLName xml.Name `xml:"patent_assignee" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the patent assignee item list and returns the new count of items
func (patentAssigneeItemList *PatentAssigneeItemList) Append(item *Item) int {
	patentAssigneeItemList.Items = append(patentAssigneeItemList.Items, item)
	return len(patentAssigneeItemList.Items)
}

// Length return item count
func (patentAssigneeItemList *PatentAssigneeItemList) Length() int {
	return len(patentAssigneeItemList.Items)
}

// IndexOf return item or nil
func (patentAssigneeItemList *PatentAssigneeItemList) IndexOf(i int) *Item {
	if i >= 0 && i < patentAssigneeItemList.Length() {
		return patentAssigneeItemList.Items[i]
	}
	return nil
}

// ShelfItemList
type ShelfItemList struct {
	XMLName xml.Name `xml:"shelves" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the shelf item list and returns the new count of items
func (shelfItemList *ShelfItemList) Append(item *Item) int {
	shelfItemList.Items = append(shelfItemList.Items, item)
	return len(shelfItemList.Items)
}

// Length return item count
func (shelfItemList *ShelfItemList) Length() int {
	return len(shelfItemList.Items)
}

// IndexOf return item or nil
func (shelfItemList *ShelfItemList) IndexOf(i int) *Item {
	if i >= 0 && i <= shelfItemList.Length() {
		return shelfItemList.Items[i]
	}
	return nil
}

// GScholarItemList
type GScholarItemList struct {
	XMLName xml.Name `xml:"gscholar" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the gScholar item list and returns the new count of items
func (gScholarItemList *GScholarItemList) Append(item *Item) int {
	gScholarItemList.Items = append(gScholarItemList.Items, item)
	return len(gScholarItemList.Items)
}

// Length return item count
func (gScholarItemList *GScholarItemList) Length() int {
	return len(gScholarItemList.Items)
}

// IndexOf return item or nil
func (gScholarItemList *GScholarItemList) IndexOf(i int) *Item {
	if i >= 0 && i < gScholarItemList.Length() {
		return gScholarItemList.Items[i]
	}
	return nil
}

// AltTitleItemList
type AltTitleItemList struct {
	XMLName xml.Name `xml:"alt_title" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the altTitle item list and returns the new count of items
func (altTitleItemList *AltTitleItemList) Append(item *Item) int {
	altTitleItemList.Items = append(altTitleItemList.Items, item)
	return len(altTitleItemList.Items)
}

// Length return item count
func (altTitleItemList *AltTitleItemList) Length() int {
	return len(altTitleItemList.Items)
}

// IndexOf return item or nil
func (altTitleItemList *AltTitleItemList) IndexOf(i int) *Item {
	if i >= 0 && i < altTitleItemList.Length() {
		return altTitleItemList.Items[i]
	}
	return nil
}

// ConfCreatorItemList
type ConfCreatorItemList struct {
	XMLName xml.Name `xml:"conf_creators" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the confCreator item list and returns the new count of items
func (confCreatorItemList *ConfCreatorItemList) Append(item *Item) int {
	confCreatorItemList.Items = append(confCreatorItemList.Items, item)
	return len(confCreatorItemList.Items)
}

// Length return item count
func (confCreatorItemList *ConfCreatorItemList) Length() int {
	return len(confCreatorItemList.Items)
}

// IndexOf return item or nil
func (confCreatorItemList *ConfCreatorItemList) IndexOf(i int) *Item {
	if i >= 0 && i < confCreatorItemList.Length() {
		return confCreatorItemList.Items[i]
	}
	return nil
}

// ReferenceItemList
type ReferenceItemList struct {
	XMLName xml.Name `xml:"reference" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the reference item list and returns the new count of items
func (referenceItemList *ReferenceItemList) Append(item *Item) int {
	referenceItemList.Items = append(referenceItemList.Items, item)
	return len(referenceItemList.Items)
}

// Length return item count
func (referenceItemList *ReferenceItemList) Length() int {
	return len(referenceItemList.Items)
}

// IndexOf return item or nil
func (referenceItemList *ReferenceItemList) IndexOf(i int) *Item {
	if i >= 0 && i < referenceItemList.Length() {
		return referenceItemList.Items[i]
	}
	return nil
}

// LearningLevelItemList
type LearningLevelItemList struct {
	XMLName xml.Name `xml:"learning_level" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the learningLevel item list and returns the new count of items
func (learningLevelItemList *LearningLevelItemList) Append(item *Item) int {
	learningLevelItemList.Items = append(learningLevelItemList.Items, item)
	return len(learningLevelItemList.Items)
}

// Length return item count
func (learningLevelItemList *LearningLevelItemList) Length() int {
	return len(learningLevelItemList.Items)
}

// IndexOf return item or nil
func (learningLevelItemList *LearningLevelItemList) IndexOf(i int) *Item {
	if i >= 0 && i < learningLevelItemList.Length() {
		return learningLevelItemList.Items[i]
	}
	return nil
}

// CopyrightHolderItemList
type CopyrightHolderItemList struct {
	XMLName xml.Name `xml:"copyright_holders" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the copyrightHolder item list and returns the new count of items
func (copyrightHolderItemList *CopyrightHolderItemList) Append(item *Item) int {
	copyrightHolderItemList.Items = append(copyrightHolderItemList.Items, item)
	return len(copyrightHolderItemList.Items)
}

// Length return item count
func (copyrightHolderItemList *CopyrightHolderItemList) Length() int {
	return len(copyrightHolderItemList.Items)
}

// IndexOf return item or nil
func (copyrightHolderItemList *CopyrightHolderItemList) IndexOf(i int) *Item {
	if i >= 0 && i < copyrightHolderItemList.Length() {
		return copyrightHolderItemList.Items[i]
	}
	return nil
}

// SkillAreaItemList
type SkillAreaItemList struct {
	XMLName xml.Name `xml:"skill_areas" json:"-"`
	Items   []*Item  `xml:"item,omitempty" jsons:"item,omitempty"`
}

// Append adds an item to the skillArea item list and returns the new count of items
func (skillAreaItemList *SkillAreaItemList) Append(item *Item) int {
	skillAreaItemList.Items = append(skillAreaItemList.Items, item)
	return len(skillAreaItemList.Items)
}

// Length return item count
func (skillAreaItemList *SkillAreaItemList) Length() int {
	return len(skillAreaItemList.Items)
}

// IndexOf return item or nil
func (skillAreaItemList *SkillAreaItemList) IndexOf(i int) *Item {
	if i >= 0 && i < skillAreaItemList.Length() {
		return skillAreaItemList.Items[i]
	}
	return nil
}

// AccompanimentItemList
type AccompanimentItemList struct {
	XMLName xml.Name `xml:"accompaniment" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// Append adds an item to the accompaniment item list and returns the new count of items
func (accompanimentItemList *AccompanimentItemList) Append(item *Item) int {
	accompanimentItemList.Items = append(accompanimentItemList.Items, item)
	return len(accompanimentItemList.Items)
}

// Length return item count
func (accompanimentItemList *AccompanimentItemList) Length() int {
	return len(accompanimentItemList.Items)
}

// IndexOf return item or nil
func (accompanimentItemList *AccompanimentItemList) IndexOf(i int) *Item {
	if i >= 0 && i < accompanimentItemList.Length() {
		return accompanimentItemList.Items[i]
	}
	return nil
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

// Append an EPrint struct to an EPrints struct returning the count of attached eprints
func (eprints *EPrints) Append(eprint *EPrint) int {
	eprints.EPrint = append(eprints.EPrint, eprint)
	return len(eprints.EPrint)
}

// Length returns EPrint count
func (eprints *EPrints) Length() int {
	return len(eprints.EPrint)
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
	MimeType         string   `xml:"mime_type" json:"mime_type"`
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

	Content  string    `xml:"content,omitempty" json:"content,omitempty"`
	Relation *ItemList `xml:"relation>item,omitempty" json:"relation,omitempty"`
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
	return len(*documentList)
}

// GetDocument takes a position (zero based) and returns the *Document
// in the DocumentList.
func (documentList DocumentList) IndexOf(i int) *Document {
	if i < 0 || i >= len(documentList) {
		return nil
	}
	return documentList[i]
}

// ItemList is an array of pointers to Item structs
type ItemList []*Item

// Append adds a Item to the relation list and returns the new count of items
func (itemList *ItemList) Append(item *Item) int {
	*itemList = append(*itemList, item)
	return len(*itemList)
}

// Length returns the length of ItemList
func (itemList *ItemList) Length() int {
	return len(*itemList)
}

// IndexOf return item or nil
func (itemList *ItemList) IndexOf(i int) *Item {
	if i >= 0 && i < itemList.Length() {
		for j, item := range *itemList {
			if j == i {
				return item
			}
		}
	}
	return nil
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
