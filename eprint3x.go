// eprinttools.go is a package for working with EPrints 3.x REST API as well as XML artifacts on disc.
//
// @author R. S. Doiel, <rsdoiel@library.caltech.edu>
//
// Copyright (c) 2017, Caltech
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
	"encoding/json"
	"encoding/xml"
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
	XMLName xml.Name `xml:"eprints" json:"-"`
	XMLNS   string   `xml:"xmlns,attr,omitempty" json:"xmlns,omitempty"`
	EPrint  *EPrint  `xml:"eprint" json:"eprint"`
}

// EPrint is the record contated in a EPrints XML document such as they used
// to store revisions.
type EPrint struct {
	XMLName              xml.Name    `json:"-"`
	XMLNS                string      `xml:"xmlns,attr,omitempty" json:"xmlns,omitempty"`
	ID                   string      `xml:"id,attr,omitempty" json:"id,omitempty"`
	EPrintID             int         `xml:"eprintid,omitempty" json:"eprint_id,omitempty"`
	RevNumber            int         `xml:"rev_number,omitempty" json:"rev_number,omitempty"`
	Documents            []*Document `xml:"documents>document,omitempty" json:"documents,omitempty"`
	EPrintStatus         string      `xml:"eprint_status,omitempty" json:"eprint_status,omitempty"`
	UserID               int         `xml:"userid,omitempty" json:"userid,omitempty"`
	Dir                  string      `xml:"dir,omitempty" json:"dir,omitempty"`
	DateStamp            string      `xml:"datestamp,omitempty" json:"datestamp,omitempty"`
	LastModified         string      `xml:"lastmod,omitempty" json:"lastmod,omitempty"`
	StatusChanged        string      `xml:"status_changed,omitempty" json:"status_changed,omitempty"`
	Type                 string      `xml:"type,omitempty" json:"type,omitempty"`
	MetadataVisibility   string      `xml:"metadata_visibility,omitempty" json:"metadata_visibility,omitempty"`
	Creators             []*Item     `xml:"creators>item,omitempty" json:"creators,omitempty"`
	Title                string      `xml:"title,omitempty" json:"title,omitempty"`
	IsPublished          string      `xml:"ispublished,omitempty" json:"ispublished,omitempty"`
	FullTextStatus       string      `xml:"full_text_status,omitempty" json:"full_text_status,omitempty"`
	Keywords             string      `xml:"keywords,omitempty" json:"keywords,omitempty"`
	Note                 string      `xml:"note,omitempty" json:"note,omitempty"`
	Abstract             string      `xml:"abstract,omitempty" json:"abstract,omitempty"`
	Date                 string      `xml:"date,omitempty" json:"date,omitempty"`
	DateType             string      `xml:"date_type,omitempty" json:"date_type,omitempty"`
	Series               string      `xml:"series,omitempty" json:"series,omitempty"`
	Publication          string      `xml:"publication,omitempty" json:"publication,omitempty"`
	Volumne              string      `xml:"volumne,omitempty" json:"volumne,omitempty"`
	Publisher            string      `xml:"publisher,omitempty" json:"publisher,omitempty"`
	PlaceOfPub           string      `xml:"place_of_pub,omitempty" json:"place_of_pub,omitempty"`
	Edition              string      `xml:"edition,omitempty" json:"edition,omitempty"`
	PageRange            string      `xml:"pagerange,omitempty" json:"pagerange,omitempty"`
	Pages                string      `xml:"pages,omitempty" json:"pages,omitempty"`
	EventTitle           string      `xml:"event_title,omitempty" json:"event_title,omitempty"`
	EventLocation        string      `xml:"event_location,omitempty" json:"event_location,omitempty"`
	EventDates           string      `xml:"event_dates,omitempty" json:"event_dates,omitempty"`
	IDNumber             string      `xml:"id_number,omitempty" json:"id_number,omitempty"`
	Refereed             string      `xml:"refereed,omitempty" json:"refereed,omitempty"`
	ISBN                 string      `xml:"isbn,omitempty" json:"isbn,omitempty"`
	ISSN                 string      `xml:"issn,omitempty" json:"issn,omitempty"`
	BookTitle            string      `xml:"book_title,omitempty" json:"book_title,omitempty"`
	Editors              []*Item     `xml:"editors>item,omitempty" json:"editors,omitempty"`
	OfficialURL          string      `xml:"official_url,omitempty" json:"official_url,omitempty"`
	RelatedURL           []*Item     `xml:"related_url>item,omitempty" json:"related_url,omitempty"`
	ReferenceText        []*Item     `xml:"referencetext>item,omitempty" json:"referencetext,omitempty"`
	Projects             []*Item     `xml:"projects>item,omitempty" json:"projects,omitempty"`
	Rights               string      `xml:"rights,omitempty" json:"rights,omitempty"`
	Funders              []*Item     `xml:"funders>item,omitempty" json:"funders,omitempty"`
	Collection           string      `xml:"collection,omitempty" json:"collection,omitempty"`
	Reviewer             string      `xml:"reviewer,omitempty" json:"reviewer,omitempty"`
	OfficeCitation       string      `xml:"official_cit,omitempty" json:"official_cit,omitempty"`
	OtherNumberingSystem []*Item     `xml:"other_numbering_system>item,omitempty" json:"other_numbering_system,omitempty"`
	LocalGroup           []*Item     `xml:"local_group>item,omitempty" json:"local_group,omitempty"`
	Errata               []*Item     `xml:"errata>item,omitempty" json:"errata,omitempty"`
	Contributors         []*Item     `xml:"contributors>item,omitempty" json:"contributors,omitempty"`
	MonographType        string      `xml:"monograph_type,omitempty" json:"monograph_type,omitempty"`

	// NOTE: Misc fields discoverd exploring REST API records, not currently used at Caltech Library (RSD, 2018-01-02)
	Subjects           []*Item `xml:"subjects>item,omitempty" json:"subjects,omitempty"`
	PresType           string  `xml:"pres_type,omitempty" json:"presentation_type,omitempty"`
	Suggestions        string  `xml:"suggestions,omitempty" json:"suggestions,omitempty"`
	Succeeds           string  `xml:"succeeds,omitempty" json:"succeeds,omitempty"`
	Commentary         string  `xml:"commentary,omitempty" json:"commentary,omitempty"`
	ContactEMail       string  `xml:"contact_email,omitempty" json:"contect_email,omitempty"`
	FileInfo           string  `xml:"fileinfo,omitempty" json:"file_info,omitempty"`
	Latitude           string  `xml:"latitude,omitempty" json:"latitude,omitempty"`
	Longitude          string  `xml:"longitude,omitempty" json:"longitude,omitempty"`
	ItemIssues         []*Item `xml:"item_issues>item,omitempty" json:"item_issues,omitempty"`
	ItemIssuesCount    int     `xml:"item_issues_count,omitempty" json:"item_issues_count,omitempty"`
	CorpCreators       []*Item `xml:"corp_creators>item,omitempty" json:"corp_creators,omitempty"`
	Department         string  `xml:"department,omitempty" json:"department,omitempty"`
	OutputMedia        string  `xml:"output_media,omitempty" json:"output_media,omitempty"`
	Exhibitors         []*Item `xml:"exhibitors,omitempty" json:"exhibitors,omitempty"`
	NumPieces          string  `xml:"num_pieces,omitempty" json:"num_pieces,omitempty"`
	CompositionType    string  `xml:"composition_type,omitempty" json:"composition_type,omitempty"`
	Producers          []*Item `xml:"producers>item,omitempty" json:"producers,omitempty"`
	Conductors         []*Item `xml:"conductors>item,omitempty" json:"conductors,omitempty"`
	Lyricists          []*Item `xml:"lyricists>item,omitempty" json:"lyricists,omitempty"`
	Accompaniment      []*Item `xml:"accompaniment>item,omitempty" json:"accompaniment,omitempty"`
	DataType           string  `xml:"data_type,omitempty" json:"data_type,omitempty"`
	PedagogicType      string  `xml:"pedagogic_type,omitempty" json:"pedagogic_type,omitempty"`
	CompletionTime     string  `xml:"completion_time,omitempty" json:"completion_time,omitempty"`
	TaskPurpose        string  `xml:"task_purpose,omitempty" json:"task_purpose,omitempty"`
	SkillAreas         []*Item `xml:"skill_areas>item,omitempty" json:"skill_areas,omitempty"`
	CopyrightHolders   []*Item `xml:"copyright_holders>item,omitempty" json:"copyright_holders,omitempty"`
	LearningLevel      []*Item `xml:"learning_level>item,omitempty" json:"learning_level,omitempty"`
	DOI                string  `xml:"doi,omitempty" json:"doi,omitempty"`
	PMCID              string  `xml:"pmc_id,omitempty" json:"pmc_id,omitempty"`
	PMID               string  `xml:"pmid,omitempty" json:"pmid,omitempty"`
	ParentURL          string  `xml:"parent_url,omitempty" json:"parent_url,omitempty"`
	AltURL             string  `xml:"alt_url,omitempty" json:"alt_url,omitempty"`
	Reference          []*Item `xml:"reference>item,omitempty" json:"reference,omitempty"`
	ConfCreators       []*Item `xml:"conf_creators>item,omitempty" json:"conf_creators,omitempty"`
	AltTitle           []*Item `xml:"alt_title>item,omitempty" json:"alt_title,omitempty"`
	TOC                string  `xml:"toc,omitempty" json:"toc,omitempty"`
	Interviewer        string  `xml:"interviewer,omitempty" json:"interviewer,omitempty"`
	InterviewDate      string  `xml:"interviewdate,omitempty" json:"interviewdate,omitempty"`
	GScholar           []*Item `xml:"gscholar,omitempty" json:"gscholar,omitempty"`
	NonSubjKeywords    string  `xml:"nonsubj_keywords,omitempty" json:"nonsubj_keywords,omitempty"`
	Season             string  `xml:"season,omitempty" json:"season,omitempty"`
	ClassificationCode string  `xml:"classification_code,omitempty" json:"classification_code,omitempty"`
	Shelves            []*Item `xml:"shelves>item,omitempty" json:"shelves,omitempty"`

	// NOTE: Sword deposit fields
	SwordDepository string `xml:"sword_depository,omitempty" json:"sword_depository,omitempty"`
	SwordSlug       string `xml:"sword_slug,omitempty" json:"sword_slug,omitempty"`
	ImportID        string `xml:"importid,omitempty" json:"import_id,omitempty"`

	// Patent related fields
	PatentApplicant      string  `xml:"patent_applicant,omitempty" json:"patent_applicant,omitempty"`
	PatentNumber         string  `xml:"patent_number,omitempty" json:"patent_number,omitempty"`
	PatentAssignee       []*Item `xml:"patent_assignee>item,omitempty" json:"patent_assignee,omitempty"`
	PatentClassification string  `xml:"patent_classification>item,omitempty" json:"patent_classification,omitempty"`
	RelatedPatents       []*Item `xml:"related_patents>item,omitempty" json:"related_patents,omitempty"`

	// Thesis oriented fields
	Divisions              []*Item `xml:"divisions>item,omitemmpty" json:"divisions,omitempty"`
	Institution            string  `xml:"institution,omitempty" json:"institution,omitempty"`
	ThesisType             string  `xml:"thesis_type,omitempty" json:"thesis_type,omitempty"`
	ThesisAdvisor          []*Item `xml:"thesis_advisor>item,omitempty" json:"thesis_advisor,omitempty"`
	ThesisCommittee        []*Item `xml:"thesis_committee>item,omitempty" json:"thesis_committee,omitempty"`
	ThesisDegree           string  `xml:"thesis_degree,omitempty" json:"thesis_degree,omitempty"`
	ThesisDegreeGrantor    string  `xml:"thesis_degree_grantor,omitempty" json:"thesis_degree_grantor,omitempty"`
	ThesisDegreeDate       string  `xml:"thesis_degree_date,omitempty" json:"thesis_degree_date,omitempty"`
	ThesisSubmittedDate    string  `xml:"thesis_submit_date,omitempty" json:"thesis_submit_date,omitempty"`
	ThesisDefenseDate      string  `xml:"thesis_defense_date,omitempty" json:"thesis_defense_date,omitempty"`
	ThesisApprovedDate     string  `xml:"thesis_approved_date,omitempty" json:"thesis_approved_date,omitempty"`
	ThesisPublicDate       string  `xml:"thesis_public_date,omitempty" json:"thesis_public_date,omitempty"`
	ThesisAuthorEMail      string  `xml:"thesis_author_email,omitempty" json:"thesis_author_email,omitempty"`
	HideThesisAuthorEMail  string  `xml:"hide_thesis_author_email,omitempty" json:"hide_thesis_author_email,omitempty"`
	GradOfficeApprovalDate string  `xml:"gradofc_approval_date,omitempty" json:"gradofc_approval_date,omitempty"`
	ThesisAwards           string  `xml:"thesis_awards,omitempty" json:"thesis_awards,omitempty"`
	ReviewStatus           string  `xml:"review_status,omitempty" json:"review_status,omitempty"`
	OptionMajor            []*Item `xml:"option_major>item,omitempty" json:"option_major,omitempty"`
	CopyrightStatement     string  `xml:"copyright_statement,omitempty" json:"copyright_statement,omitempty"`
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

// Name handles the "name" types found in Items.
type Name struct {
	XMLName xml.Name `json:"-"`
	Family  string   `xml:"family,omitempty" json:"family,omitempty"`
	Given   string   `xml:"given,omitempty" json:"given,omitempty"`
	ID      string   `xml:"id,omitempty" json:"id,omitempty"`
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

// Generic attempts to parse a generic XML fragment into an array of JSON items
type Generic struct {
	XMLName   xml.Name
	EPrints   *EPrints    `xml:"eprints,omitempty" json:"eprints,omitempty"`
	EPrint    *EPrint     `xml:"eprint,omitempty" json:"eprint,omitempty"`
	Documents []*Document `xml:"document,omitempty" json:"documents,omitempty"`
	Files     []*File     `xml:"file,omitempty" json:"files,omitempty"`
	Items     []*Item     `xml:"item" json:"items,omitempty"`
	Value     string      `xml:",chardata" json:"value,omitempty"`
}

// MarshalJSON() normalizes Generic into a sensible JSON structure
func (g *Generic) MarshalJSON() ([]byte, error) {
	var (
		k string
	)
	m := map[string]interface{}{}

	if g.XMLName.Local != "" {
		m["xml_name"] = g.XMLName
		k = g.XMLName.Local
	} else {
		k = "items"
	}
	flatten := true
	switch {
	case g.EPrints != nil:
		m[k] = g.EPrints
		flatten = false
	case g.EPrint != nil:
		m[k] = g.EPrint
		flatten = false
	case len(g.Documents) > 0:
		m[k] = g.Documents
		flatten = false
	case len(g.Files) > 0:
		m[k] = g.Files
		flatten = false
	case len(g.Items) > 0:
		m[k] = g.Items
		flatten = false
	}
	if flatten {
		s := strings.TrimSpace(g.Value)
		return json.Marshal(s)
	}
	return json.Marshal(m)
}

func (name *Name) String() string {
	src, err := json.Marshal(name)
	if err != nil {
		return ""
	}
	return string(src)
}
