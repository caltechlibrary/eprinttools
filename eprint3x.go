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
	"encoding/xml"
)

//
// NOTE: This file contains the general structure in Caltech Libraries EPrints 3.x based repositories.
//

// EPrints is the high level XML you get from the REST API.
// E.g. curl -L -O https://eprints3.example.org/rest/eprint/1234.xml
// Then parse the 1234.xml document stucture.
type EPrints struct {
	XMLName xml.Name `json:"-"`
	XMLNS   string   `xml:"xmlns,attr,omitempty" json:"xmlns,omitempty"`
	EPrint  *EPrint  `xml:"eprint" json:"eprint"`
}

// EPrint is the record contated in a EPrints XML document such as they used
// to store revisions.
type EPrint struct {
	XMLName            xml.Name    `json:"-"`
	XMLNS              string      `xml:"xmlns,attr,omitempty" json:"xmlns,omitempty"`
	ID                 string      `xml:"id,attr,omitempty" json:"id,omitempty"`
	EPrintID           string      `xml:"eprintid,omitempty" json:"eprint_id,omitempty"`
	RevNumber          int         `xml:"rev_number,omitempty" json:"rev_number,omitempty"`
	Documents          []*Document `xml:"documents>document,omitempty" json:"documents,omitempty"`
	EPrintStatus       string      `xml:"eprint_status,omitempty" json:"eprint_status,omitempty"`
	UserID             int         `xml:"userid,omitempty" json:"userid,omitempty"`
	Dir                string      `xml:"dir,omitempty" json:"dir,omitempty"`
	DateStamp          string      `xml:"datestamp,omitempty" json:"datestamp,omitempty"`
	LastModified       string      `xml:"lastmod,omitempty" json:"lastmod,omitempty"`
	StatusChanged      string      `xml:"status_changed,omitempty" json:"status_changed,omitempty"`
	Type               string      `xml:"type,omitempty" json:"type,omitempty"`
	MetadataVisibility string      `xml:"metadata_visibility,omitempty" json:"metadata_visibility,omitempty"`
	Creators           []*Item     `xml:"creators>item,omitempty" json:"creators,omitempty"`
	Title              string      `xml:"title,omitempty" json:"title,omitempty"`
	IsPublished        string      `xml:"ispublished,omitempty" json:"ispublished,omitempty"`
	FullTextStatus     string      `xml:"full_text_status,omitempty" json:"full_text_status,omitempty"`
	Keywords           string      `xml:"keywords,omitempty" json:"keywords,omitempty"`
	Note               string      `xml:"note,omitempty" json:"note,omitempty"`
	Abstract           string      `xml:"abstract,omitempty" json:"abstract,omitempty"`
	Date               string      `xml:"date,omitempty" json:"date,omitempty"`
	DateType           string      `xml:"date_type,omitempty" json:"date_type,omitempty"`
	Publication        string      `xml:"publication,omitempty" json:"publication,omitempty"`
	Volumne            string      `xml:"volumne,omitempty" json:"volumne,omitempty"`
	Publisher          string      `xml:"publisher,omitempty" json:"publisher,omitempty"`
	PageRange          string      `xml:"pagerange,omitempty" json:"pagerange,omitempty"`
	IDNumber           string      `xml:"id_number,omitempty" json:"id_number,omitempty"`
	Refereed           bool        `xml:"refereed,omitempty" json:"refereed,omitempty"`
	OfficialURL        string      `xml:"official_url,omitempty" json:"official_url,omitempty"`
	RelatedURL         []*Item     `xml:"related_url>item,omitempty" json:"related_url,omitempty"`
	ReferenceText      []*Item     `xml:"referencetext>item,omitempty" json:"referencetext,omitempty"`
	Rights             string      `xml:"rights,omitempty" json:"rights,omitempty"`
	Funders            []*Item     `xml:"funders>item,omitempty" json:"funders,omitempty"`
	Collection         string      `xml:"collection,omitempty" json:"collection,omitempty"`
	Reviewer           string      `xml:"reviewer,omitempty" json:"reviewer,omitempty"`

	// Thesis oriented fields
	ItemIssuesCount     int     `xml:"item_issues_count,omitempty" json:"item_issues_count,omitempty"`
	Divisions           []*Item `xml:"divisions>item,omitemmpty" json:"divisions,omitempty"`
	Institution         string  `xml:"institution,omitempty" json:"institution,omitempty"`
	ThesisType          string  `xml:"thesis_type,omitempty" json:"thesis_type,omitempty"`
	ThesisAdvisor       []*Item `xml:"thesis_advisor>item,omitempty" json:"thesis_advisor,omitempty"`
	ThesisCommittee     []*Item `xml:"thesis_committee>item,omitempty" json:"thesis_committee,omitempty"`
	ThesisDegree        string  `xml:"thesis_degree,omitempty" json:"thesis_degree,omitempty"`
	ThesisDegreeGrantor string  `xml:"thesis_degree_grantor,omitempty" json:"thesis_degree_grantor,omitempty"`
	ThesisSubmittedDate string  `xml:"thesis_submit_date,omitempty" json:"thesis_submit_date,omitempty"`
	ThesisDefenseDate   string  `xml:"thesis_defense_date,omitempty" json:"thesis_defense_date,omitempty"`
	ThesisApprovedDate  string  `xml:"thesis_approved_date,omitempty" json:"thesis_approved_date,omitempty"`
	ReviewStatus        string  `xml:"review_status,omitempty" json:"review_status,omitempty"`
	OptionMajor         []*Item `xml:"option_major>item,omitempty" json:"option_major,omitempty"`
	CopyrightStatement  string  `xml:"copyright_statement,omitempty" json:"copyright_statement,omitempty"`
}

// Item is a generic type used by various fields (e.g. Creator, Division, OptionMajor)
type Item struct {
	XMLName     xml.Name `xml:"item" json:"-"`
	Name        *Name    `xml:"name,omitempty" json:"name,omitempty"`
	ID          string   `xml:"id,omitempty" json:"id,omitempty"`
	EMail       string   `xml:"email,omitempty" json:"email,omitempty"`
	ShowEMail   bool     `xml:"show_email,omitempty" json:"show_email,omitempty"`
	Role        string   `xml:"role,omitempty" json:"role,omitempty"`
	URL         string   `xml:"url,omitempty" json:"url,omitempty"`
	Type        string   `xml:"type,omitempty" json:"type,omitempty"`
	Description string   `xml:"description,omitempty" json:"description,omitempty"`
	Agency      string   `xml:"agency,omitempty" json:"agency,omitempty"`
	GrantNumber string   `xml:"grant_number,omitempty" json:"grant_number,omitempty"`
	URI         string   `xml:"uri,omitempty" json:"uri,omitempty"`
	ORCID       string   `xml:"orcid,omitempty" json:"orcid,omitempty"`
	InnerText   string   `xml:",chardata" json:"inner_text,omitempty"`
}

type Name struct {
	XMLName   xml.Name `json:"-"`
	Family    string   `xml:"family,omitempty" json:"family,omitempty"`
	Given     string   `xml:"given,omitempty" json:"given,omitempty"`
	ID        string   `xml:"id,omitempty" json:"id,omitempty"`
	InnerText string   `xml:",chardata" json:"inner_text,omitempty"`
}
