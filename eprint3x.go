// eprinttools.go is a package for working with EPrints 3.x REST API as well as XML artifacts on disc.
//
// @author R. S. Doiel, <rsdoiel@library.caltech.edu>
//
// Copyright (c) 2018, Caltech
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
	"log"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	// Caltech Library packages
	"github.com/caltechlibrary/rc"
)

const (
	maxConsecutiveFailedRequests = 10
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
	DateStamp            string                        `xml:"datestamp,omitempty" json:"datestamp,omitempty"`
	LastModified         string                        `xml:"lastmod,omitempty" json:"lastmod,omitempty"`
	StatusChanged        string                        `xml:"status_changed,omitempty" json:"status_changed,omitempty"`
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
	RelatedURL           *RelatedURLItemList           `xml:"related_url,omitempty" json:"related_url,omitempty"`
	ReferenceText        *ReferenceTextItemList        `xml:"referencetext,omitempty" json:"referencetext,omitempty"`
	Projects             *ProjectItemList              `xml:"projects,omitempty" json:"projects,omitempty"`
	Rights               string                        `xml:"rights,omitempty" json:"rights,omitempty"`
	Funders              *FunderItemList               `xml:"funders,omitempty" json:"funders,omitempty"`
	Collection           string                        `xml:"collection,omitempty" json:"collection,omitempty"`
	Reviewer             string                        `xml:"reviewer,omitempty" json:"reviewer,omitempty"`
	OfficeCitation       string                        `xml:"official_cit,omitempty" json:"official_cit,omitempty"`
	OtherNumberingSystem *OtherNumberingSystemItemList `xml:"other_numbering_system,omitempty" json:"other_numbering_system,omitempty"`
	LocalGroup           *LocalGroupItemList           `xml:"local_group,omitempty" json:"local_group,omitempty"`
	Errata               *ErrataItemList               `xml:"errata,omitempty" json:"errata,omitempty"`
	Contributors         *ContributorItemList          `xml:"contributors,omitempty" json:"contributors,omitempty"`
	MonographType        string                        `xml:"monograph_type,omitempty" json:"monograph_type,omitempty"`

	// Caltech Library uses suggestions as an internal note field (RSD, 2018-02-15)
	Suggestions string            `xml:"suggestions,omitempty" json:"suggestions,omitempty"`
	OtherURL    *OtherURLItemList `xml:"other_url,omitempty" json:"other_url,omitempty"`

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
	LearningLevel      *LearningLevelItemList   `xml:"learning_level,omitempty" json:"learning_level,omitempty"`
	DOI                string                   `xml:"doi,omitempty" json:"doi,omitempty"`
	PMCID              string                   `xml:"pmc_id,omitempty" json:"pmc_id,omitempty"`
	PMID               string                   `xml:"pmid,omitempty" json:"pmid,omitempty"`
	ParentURL          string                   `xml:"parent_url,omitempty" json:"parent_url,omitempty"`
	AltURL             string                   `xml:"alt_url,omitempty" json:"alt_url,omitempty"`
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
	SwordSlug       string `xml:"sword_slug,omitempty" json:"sword_slug,omitempty"`
	ImportID        string `xml:"importid,omitempty" json:"import_id,omitempty"`

	// Patent related fields
	PatentApplicant      string                        `xml:"patent_applicant,omitempty" json:"patent_applicant,omitempty"`
	PatentNumber         string                        `xml:"patent_number,omitempty" json:"patent_number,omitempty"`
	PatentAssignee       *PatentAssigneeItemList       `xml:"patent_assignee,omitempty" json:"patent_assignee,omitempty"`
	PatentClassification *PatentClassificationItemList `xml:"patent_classification,omitempty" json:"patent_classification,omitempty"`
	RelatedPatents       *RelatedPatentItemList        `xml:"related_patents,omitempty" json:"related_patents,omitempty"`

	// Thesis oriented fields
	Divisions              *DivisionItemList        `xml:"divisions,omitemmpty" json:"divisions,omitempty"`
	Institution            string                   `xml:"institution,omitempty" json:"institution,omitempty"`
	ThesisType             string                   `xml:"thesis_type,omitempty" json:"thesis_type,omitempty"`
	ThesisAdvisor          *ThesisAdvisorItemList   `xml:"thesis_advisor,omitempty" json:"thesis_advisor,omitempty"`
	ThesisCommittee        *ThesisCommitteeItemList `xml:"thesis_committee,omitempty" json:"thesis_committee,omitempty"`
	ThesisDegree           string                   `xml:"thesis_degree,omitempty" json:"thesis_degree,omitempty"`
	ThesisDegreeGrantor    string                   `xml:"thesis_degree_grantor,omitempty" json:"thesis_degree_grantor,omitempty"`
	ThesisDegreeDate       string                   `xml:"thesis_degree_date,omitempty" json:"thesis_degree_date,omitempty"`
	ThesisSubmittedDate    string                   `xml:"thesis_submit_date,omitempty" json:"thesis_submit_date,omitempty"`
	ThesisDefenseDate      string                   `xml:"thesis_defense_date,omitempty" json:"thesis_defense_date,omitempty"`
	ThesisApprovedDate     string                   `xml:"thesis_approved_date,omitempty" json:"thesis_approved_date,omitempty"`
	ThesisPublicDate       string                   `xml:"thesis_public_date,omitempty" json:"thesis_public_date,omitempty"`
	ThesisAuthorEMail      string                   `xml:"thesis_author_email,omitempty" json:"thesis_author_email,omitempty"`
	HideThesisAuthorEMail  string                   `xml:"hide_thesis_author_email,omitempty" json:"hide_thesis_author_email,omitempty"`
	GradOfficeApprovalDate string                   `xml:"gradofc_approval_date,omitempty" json:"gradofc_approval_date,omitempty"`
	ThesisAwards           string                   `xml:"thesis_awards,omitempty" json:"thesis_awards,omitempty"`
	ReviewStatus           string                   `xml:"review_status,omitempty" json:"review_status,omitempty"`
	OptionMajor            *OptionMajorItemList     `xml:"option_major,omitempty" json:"option_major,omitempty"`
	CopyrightStatement     string                   `xml:"copyright_statement,omitempty" json:"copyright_statement,omitempty"`
}

// Item is a generic type used by various fields (e.g. Creator, Division, OptionMajor)
type Item struct {
	XMLName     xml.Name `xml:"item" json:"-"`
	Name        *Name    `xml:"name,omitempty" json:"name,omitempty"`
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
	fmt.Printf("DEBUG *Item src %s\n", src)
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

// ItemList holds an array of items (e.g. creators, related urls, etc)
//type ItemList []*Item

// CreatorItemList holds a list of authors
type CreatorItemList struct {
	XMLName xml.Name `xml:"creators" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the Creator list and returns the new count of items
func (creatorItemList *CreatorItemList) AddItem(item *Item) int {
	creatorItemList.Items = append(creatorItemList.Items, item)
	return len(creatorItemList.Items)
}

// EditorItemList holds a list of editors
type EditorItemList struct {
	XMLName xml.Name `xml:"editors" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the Editor item list and returns the new count of items
func (editorItemList *EditorItemList) AddItem(item *Item) int {
	editorItemList.Items = append(editorItemList.Items, item)
	return len(editorItemList.Items)
}

// RelatedURLItemList holds the related URLs (e.g. doi, aux material doi)
type RelatedURLItemList struct {
	XMLName xml.Name `xml:"related_url" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the related url item list and returns the new count of items
func (relatedURLItemList *RelatedURLItemList) AddItem(item *Item) int {
	relatedURLItemList.Items = append(relatedURLItemList.Items, item)
	return len(relatedURLItemList.Items)
}

// OtherURLItemList is a legacy Caltech Library field, old records have
// it new records use RelatedURLItemList
// RelatedURLItemList holds the related URLs (e.g. doi, aux material doi)
type OtherURLItemList struct {
	XMLName xml.Name `xml:"other_url" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the "other" url item list and returns the new count of items, this is a legacy Caltech Library-ism in EPrints
func (otherURLItemList *OtherURLItemList) AddItem(item *Item) int {
	otherURLItemList.Items = append(otherURLItemList.Items, item)
	return len(otherURLItemList.Items)
}

// ReferenceTextItemList
type ReferenceTextItemList struct {
	XMLName xml.Name `xml:"referencetext" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the reference text url item list and returns the new count of items
func (referenceTextItemList *ReferenceTextItemList) AddItem(item *Item) int {
	referenceTextItemList.Items = append(referenceTextItemList.Items, item)
	return len(referenceTextItemList.Items)
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
		referenceTextItemList.AddItem(item)
	}
	return err
}

// ProjectItemList
type ProjectItemList struct {
	XMLName xml.Name `xml:"projects" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the project item list and returns the new count of items
func (projectItemList *ProjectItemList) AddItem(item *Item) int {
	projectItemList.Items = append(projectItemList.Items, item)
	return len(projectItemList.Items)
}

// FunderItemList
type FunderItemList struct {
	XMLName xml.Name `xml:"funders" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the funder item list and returns the new count of items
func (funderItemList *FunderItemList) AddItem(item *Item) int {
	funderItemList.Items = append(funderItemList.Items, item)
	return len(funderItemList.Items)
}

// LocalGroupItemList holds the related URLs (e.g. doi, aux material doi)
type LocalGroupItemList struct {
	XMLName xml.Name `xml:"local_group" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the local group item list and returns the new count of items
func (localGroupItemList *LocalGroupItemList) AddItem(item *Item) int {
	localGroupItemList.Items = append(localGroupItemList.Items, item)
	return len(localGroupItemList.Items)
}

// OtherNumberingSystemItemList
type OtherNumberingSystemItemList struct {
	XMLName xml.Name `xml:"other_numbering_system" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the other numbering system item list and returns the new count of items
func (otherNumberingSystemItemList *OtherNumberingSystemItemList) AddItem(item *Item) int {
	otherNumberingSystemItemList.Items = append(otherNumberingSystemItemList.Items, item)
	return len(otherNumberingSystemItemList.Items)
}

// ErrataItemList
type ErrataItemList struct {
	XMLName xml.Name `xml:"errata" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the errata item list and returns the new count of items
func (errataItemList *ErrataItemList) AddItem(item *Item) int {
	errataItemList.Items = append(errataItemList.Items, item)
	return len(errataItemList.Items)
}

// ContributorItemList
type ContributorItemList struct {
	XMLName xml.Name `xml:"contributors" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the contributor item list and returns the new count of items
func (contributorItemList *ContributorItemList) AddItem(item *Item) int {
	contributorItemList.Items = append(contributorItemList.Items, item)
	return len(contributorItemList.Items)
}

// SubjectItemList
type SubjectItemList struct {
	XMLName xml.Name `xml:"subjects" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the subject item list and returns the new count of items
func (subjectItemList *SubjectItemList) AddItem(item *Item) int {
	subjectItemList.Items = append(subjectItemList.Items, item)
	return len(subjectItemList.Items)
}

// ItemIssueItemList
type ItemIssueItemList struct {
	XMLName xml.Name `xml:"item_issues" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the issue item list and returns the new count of items
func (issueItemList *ItemIssueItemList) AddItem(item *Item) int {
	issueItemList.Items = append(issueItemList.Items, item)
	return len(issueItemList.Items)
}

// CorpCreatorItemList
type CorpCreatorItemList struct {
	XMLName xml.Name `json:"-"` //`xml:"corp_creators" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the corp creator item list and returns the new count of items
func (corpCreatorItemList *CorpCreatorItemList) AddItem(item *Item) int {
	corpCreatorItemList.Items = append(corpCreatorItemList.Items, item)
	return len(corpCreatorItemList.Items)
}

// ExhibitorItemList
type ExhibitorItemList struct {
	XMLName xml.Name `xml:"exhibitors" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the exhibitor item list and returns the new count of items
func (exhibitorItemList *ExhibitorItemList) AddItem(item *Item) int {
	exhibitorItemList.Items = append(exhibitorItemList.Items, item)
	return len(exhibitorItemList.Items)
}

// ProducerItemList
type ProducerItemList struct {
	XMLName xml.Name `xml:"producers" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the producer item list and returns the new count of items
func (producerItemList *ProducerItemList) AddItem(item *Item) int {
	producerItemList.Items = append(producerItemList.Items, item)
	return len(producerItemList.Items)
}

// ConductorItemList
type ConductorItemList struct {
	XMLName xml.Name `xml:"conductors" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the conductor item list and returns the new count of items
func (conductorItemList *ConductorItemList) AddItem(item *Item) int {
	conductorItemList.Items = append(conductorItemList.Items, item)
	return len(conductorItemList.Items)
}

// LyricistItemList
type LyricistItemList struct {
	XMLName xml.Name `xml:"lyricists" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the lyricist item list and returns the new count of items
func (lyricistItemList *LyricistItemList) AddItem(item *Item) int {
	lyricistItemList.Items = append(lyricistItemList.Items, item)
	return len(lyricistItemList.Items)
}

// OptionMajorItemList
type OptionMajorItemList struct {
	XMLName xml.Name `xml:"option_major" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the option major item list and returns the new count of items
func (optionMajorItemList *OptionMajorItemList) AddItem(item *Item) int {
	optionMajorItemList.Items = append(optionMajorItemList.Items, item)
	return len(optionMajorItemList.Items)
}

// ThesisCommitteeItemList
type ThesisCommitteeItemList struct {
	XMLName xml.Name `xml:"thesis_committee" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the thesis committee item list and returns the new count of items
func (thesisCommitteeItemList *ThesisCommitteeItemList) AddItem(item *Item) int {
	thesisCommitteeItemList.Items = append(thesisCommitteeItemList.Items, item)
	return len(thesisCommitteeItemList.Items)
}

// ThesisAdvisorItemList
type ThesisAdvisorItemList struct {
	XMLName xml.Name `xml:"thesis_advisor" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the thesis advisor item list and returns the new count of items
func (thesisAdvisorItemList *ThesisAdvisorItemList) AddItem(item *Item) int {
	thesisAdvisorItemList.Items = append(thesisAdvisorItemList.Items, item)
	return len(thesisAdvisorItemList.Items)
}

// DivisionItemList
type DivisionItemList struct {
	XMLName xml.Name `xml:"divisions" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the division item list and returns the new count of items
func (divisionItemList *DivisionItemList) AddItem(item *Item) int {
	divisionItemList.Items = append(divisionItemList.Items, item)
	return len(divisionItemList.Items)
}

// RelatedPatentItemList
type RelatedPatentItemList struct {
	XMLName xml.Name `xml:"related_patents" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the related patent item list and returns the new count of items
func (relatedPatentItemList *RelatedPatentItemList) AddItem(item *Item) int {
	relatedPatentItemList.Items = append(relatedPatentItemList.Items, item)
	return len(relatedPatentItemList.Items)
}

// PatentClassificationItemList
type PatentClassificationItemList struct {
	XMLName xml.Name `xml:"patent_classification" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the patent classification item list and returns the new count of items
func (patentClassificationItemList *PatentClassificationItemList) AddItem(item *Item) int {
	patentClassificationItemList.Items = append(patentClassificationItemList.Items, item)
	return len(patentClassificationItemList.Items)
}

// PatentAssigneeItemList
type PatentAssigneeItemList struct {
	XMLName xml.Name `xml:"patent_assignee" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the patent assignee item list and returns the new count of items
func (patentAssigneeItemList *PatentAssigneeItemList) AddItem(item *Item) int {
	patentAssigneeItemList.Items = append(patentAssigneeItemList.Items, item)
	return len(patentAssigneeItemList.Items)
}

// ShelfItemList
type ShelfItemList struct {
	XMLName xml.Name `xml:"shelves" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the shelf item list and returns the new count of items
func (shelfItemList *ShelfItemList) AddItem(item *Item) int {
	shelfItemList.Items = append(shelfItemList.Items, item)
	return len(shelfItemList.Items)
}

// GScholarItemList
type GScholarItemList struct {
	XMLName xml.Name `xml:"gscholar" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the gScholar item list and returns the new count of items
func (gScholarItemList *GScholarItemList) AddItem(item *Item) int {
	gScholarItemList.Items = append(gScholarItemList.Items, item)
	return len(gScholarItemList.Items)
}

// AltTitleItemList
type AltTitleItemList struct {
	XMLName xml.Name `xml:"alt_title" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the altTitle item list and returns the new count of items
func (altTitleItemList *AltTitleItemList) AddItem(item *Item) int {
	altTitleItemList.Items = append(altTitleItemList.Items, item)
	return len(altTitleItemList.Items)
}

// ConfCreatorItemList
type ConfCreatorItemList struct {
	XMLName xml.Name `xml:"conf_creators" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the confCreator item list and returns the new count of items
func (confCreatorItemList *ConfCreatorItemList) AddItem(item *Item) int {
	confCreatorItemList.Items = append(confCreatorItemList.Items, item)
	return len(confCreatorItemList.Items)
}

// ReferenceItemList
type ReferenceItemList struct {
	XMLName xml.Name `xml:"reference" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the reference item list and returns the new count of items
func (referenceItemList *ReferenceItemList) AddItem(item *Item) int {
	referenceItemList.Items = append(referenceItemList.Items, item)
	return len(referenceItemList.Items)
}

// LearningLevelItemList
type LearningLevelItemList struct {
	XMLName xml.Name `xml:"learning_level" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the learningLevel item list and returns the new count of items
func (learningLevelItemList *LearningLevelItemList) AddItem(item *Item) int {
	learningLevelItemList.Items = append(learningLevelItemList.Items, item)
	return len(learningLevelItemList.Items)
}

// CopyrightHolderItemList
type CopyrightHolderItemList struct {
	XMLName xml.Name `xml:"copyright_holders" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the copyrightHolder item list and returns the new count of items
func (copyrightHolderItemList *CopyrightHolderItemList) AddItem(item *Item) int {
	copyrightHolderItemList.Items = append(copyrightHolderItemList.Items, item)
	return len(copyrightHolderItemList.Items)
}

// SkillAreaItemList
type SkillAreaItemList struct {
	XMLName xml.Name `xml:"skill_areas" json:"-"`
	Items   []*Item  `xml:"item,omitempty" jsons:"item,omitempty"`
}

// AddItem adds an item to the skillArea item list and returns the new count of items
func (skillAreaItemList *SkillAreaItemList) AddItem(item *Item) int {
	skillAreaItemList.Items = append(skillAreaItemList.Items, item)
	return len(skillAreaItemList.Items)
}

// AccompanimentItemList
type AccompanimentItemList struct {
	XMLName xml.Name `xml:"accompaniment" json:"-"`
	Items   []*Item  `xml:"item,omitempty" json:"items,omitempty"`
}

// AddItem adds an item to the accompaniment item list and returns the new count of items
func (accompanimentItemList *AccompanimentItemList) AddItem(item *Item) int {
	accompanimentItemList.Items = append(accompanimentItemList.Items, item)
	return len(accompanimentItemList.Items)
}

// Name handles the "name" types found in Items.
type Name struct {
	XMLName xml.Name `json:"-"`
	Family  string   `xml:"family,omitempty" json:"family,omitempty"`
	Given   string   `xml:"given,omitempty" json:"given,omitempty"`
	ID      string   `xml:"id,omitempty" json:"id,omitempty"`
	ORCID   string   `xml:"orcid,omitempty" json:"orcid,omitempty"`
	Value   string   `xml:",chardata" json:"value,omitempty"`
}

// MarshalJSON() is a custom JSON marshaler for Name
func (name *Name) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{}
	flatten := true
	if s := strings.TrimSpace(name.Family); s != "" {
		m["family"] = s
		flatten = false
	}

	if s := strings.TrimSpace(name.Given); s != "" {
		m["given"] = s
		flatten = false
	}
	if s := strings.TrimSpace(name.Value); s != "" {
		if flatten == true {
			return json.Marshal(s)
		}
		m["value"] = s
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

// AddEPrint appends an EPrint struct to an EPrints struct returning the count of attached eprints
func (eprints *EPrints) AddEPrint(eprint *EPrint) int {
	eprints.EPrint = append(eprints.EPrint, eprint)
	return len(eprints.EPrint)
}

// GetEPrints retrieves an EPrint record (e.g. via REST API)
// A populated EPrints structure, the raw XML and an error.
func GetEPrints(baseURL string, authType int, username string, secret string, key string) (*EPrints, []byte, error) {
	workURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, nil, err
	}
	if workURL.Path == "" {
		workURL.Path = path.Join("rest", "eprint") + "/" + key + ".xml"
	} else {
		p := workURL.Path
		workURL.Path = path.Join(p, "rest", "eprint") + "/" + key + ".xml"
	}

	// Switch to use Rest Client Wrapper
	rest, err := rc.New(workURL.String(), authType, username, secret)
	if err != nil {
		return nil, nil, err
	}
	content, err := rest.Request("GET", workURL.Path, map[string]string{})
	if err != nil {
		return nil, nil, err
	}

	rec := new(EPrints)
	err = xml.Unmarshal(content, &rec)
	if err != nil {
		return nil, content, err
	}
	if len(rec.EPrint) > 0 && (rec.EPrint[0].EPrintStatus == "deletion" || rec.EPrint[0].EPrintStatus == "inbox" || rec.EPrint[0].EPrintStatus == "buffer") {
		return rec, content, fmt.Errorf("WARNING status %s %s", rec.EPrint[0].ID, rec.EPrint[0].EPrintStatus)
	}
	return rec, content, nil
}

// GetKeys returns a list of eprint record ids from the EPrints REST API
func GetKeys(baseURL string, authType int, username string, secret string) ([]string, error) {
	var (
		results []string
	)

	workURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	if workURL.Path == "" {
		workURL.Path = path.Join("rest", "eprint") + "/"
	} else {
		p := workURL.Path
		workURL.Path = path.Join(p, "rest", "eprint") + "/"
	}
	// Switch to use Rest Client Wrapper
	rest, err := rc.New(workURL.String(), authType, username, secret)
	if err != nil {
		return nil, fmt.Errorf("requesting %s, %s", workURL.String(), err)
	}
	content, err := rest.Request("GET", workURL.Path, map[string]string{})
	if err != nil {
		return nil, fmt.Errorf("requested %s, %s", workURL.String(), err)
	}
	eIDs := new(ePrintIDs)
	err = xml.Unmarshal(content, &eIDs)
	if err != nil {
		return nil, err
	}
	// Build a list of Unique IDs in a map, then convert unique querys to results array
	m := make(map[string]bool)
	for _, val := range eIDs.IDs {
		if strings.HasSuffix(val, ".xml") == true {
			eprintID := strings.TrimSuffix(val, ".xml")
			if _, hasID := m[eprintID]; hasID == false {
				// Save the new ID found
				m[eprintID] = true
				// Only store Unique IDs in result
				results = append(results, eprintID)
			}
		}
	}
	return results, nil
}

// GetModifiedKeys returns a list of eprint record ids from the EPrints REST API that match the modification date range
func GetModifiedKeys(baseURL string, authType int, username string, secret string, start time.Time, end time.Time, verbose bool) ([]string, error) {
	var (
		results []string
	)
	// need to calculate the base restDocPath
	workURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	if workURL.Path == "" {
		workURL.Path = path.Join("rest", "eprint") + "/"
	} else {
		p := workURL.Path
		workURL.Path = path.Join(p, "rest", "eprint") + "/"
	}
	restDocPath := workURL.Path
	// Switch to use Rest Client Wrapper
	rest, err := rc.New(baseURL, authType, username, secret)
	if err != nil {
		return nil, fmt.Errorf("requesting %s, %s", baseURL, err)
	}

	// Pass baseURL to GetKeys(), get key list then filter for modified times.
	pid := os.Getpid()
	keys, err := GetKeys(baseURL, authType, username, secret)
	// NOTE: consecutiveFailedCount tracks repeated failures
	// e.g. You need to authenticate with the server to get
	// modified information.
	consecutiveFailedCount := 0
	for _, key := range keys {
		// form a request to the REST API for just the modified date
		docPath := path.Join(restDocPath, key, "lastmod.txt")
		lastModified, err := rest.Request("GET", docPath, map[string]string{})
		if err != nil {
			if verbose == true {
				log.Printf("(pid: %d) request failed, %s", pid, err)
			}
			consecutiveFailedCount++
			if consecutiveFailedCount >= maxConsecutiveFailedRequests {
				return results, err
			}
		} else {
			consecutiveFailedCount = 0
			datestring := fmt.Sprintf("%s", lastModified)
			if len(datestring) > 9 {
				datestring = datestring[0:10]
			}
			// Parse the modified date and compare to our range
			if dt, err := time.Parse("2006-01-02", datestring); err == nil && dt.Unix() >= start.Unix() && dt.Unix() <= end.Unix() {
				// If range is OK then add the key to results
				results = append(results, key)
			}
		}
	}
	return results, nil
}

// File structures in Document
type File struct {
	XMLName   xml.Name `json:"-"`
	ID        string   `xml:"id,attr" json:"id"`
	FileID    int      `xml:"fileid" json:"fileid"`
	DatasetID string   `xml:"datasetid" json:"datasetid"`
	ObjectID  int      `xml:"objectid" json:"objectid"`
	Filename  string   `xml:"filename" json:"filename"`
	MimeType  string   `xml:"mime_type" json:"mime_type"`
	Hash      string   `xml:"hash,omitempty" json:"hash,omitempty"`
	HashType  string   `xml:"hash_type,omitempty" json:"hash_type,omitempty"`
	FileSize  int      `xml:"filesize" json:"filesize"`
	MTime     string   `xml:"mtime" json:"mtime"`
	URL       string   `xml:"url" json:"url"`
}

// Document structures inside a Record (i.e. <eprint>...<documents><document>...</document>...</documents>...</eprint>)
type Document struct {
	XMLName    xml.Name `json:"-"`
	ID         string   `xml:"id,attr" json:"id"`
	DocID      int      `xml:"docid" json:"doc_id"`
	RevNumber  int      `xml:"rev_number" json:"rev_number,omitempty"`
	Files      []*File  `xml:"files>file" json:"files,omitempty"`
	EPrintID   int      `xml:"eprintid" json:"eprint_id"`
	Pos        int      `xml:"pos" json:"pos,omitempty"`
	Placement  int      `xml:"placement,omitempty" json:"placement,omitempty"`
	MimeType   string   `xml:"mime_type" json:"mime_type"`
	Format     string   `xml:"format" json:"format"`
	FormatDesc string   `xml:"formatdesc,omitempty" json:"format_desc,omitempty"`
	Language   string   `xml:"language,omitempty" json:"language,omitempty"`
	Security   string   `xml:"security" json:"security"`
	License    string   `xml:"license" json:"license"`
	Main       string   `xml:"main" json:"main"`
	Content    string   `xml:"content,omitempty" json:"content,omitempty"`
	Relation   []*Item  `xml:"relation>item,omitempty" json:"relation,omitempty"`
}

// DocumentList is an array of pointers to Document structs
type DocumentList []*Document

// AddDocument adds a document to the documents list and returns the new count of items
func (documentList DocumentList) AddDocument(document *Document) int {
	documentList = append(documentList, document)
	return len(documentList)
}

// Length returns the length of DocumentList
func (documentList DocumentList) Length() int {
	return len(documentList)
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
