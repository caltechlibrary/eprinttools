/**
 * simplified presents an Invenio 3 like JSON representation of an EPrint
 * record. This is intended to make the development of V2 of feeds easier
 * for both our audience on internal programming needs.
 */
package eprinttools

import (
	"encoding/xml"
	"fmt"
)

// InvenioType is an Invenio 3 e.g. ResourceType, title type or language
type InvenioType struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Title string `json:"title,omitempty"`
}

type InvenioAltType struct {
	ID    string            `json:"id,omitempty"`
	Name  string            `json:"name,omitempty"`
	Title map[string]string `json:"title,omitempty"`
}

// Identifier holds an Identifier, e.g. ORCID, ROR, ISNI, GND
type Identifier struct {
	Scheme       string          `json:"scheme,omitempty"`
	Identifier   string          `json:"identifier,omitempty"`
	RelationType *InvenioAltType `json:"relation_type,omitempty"`
	ResourceType *InvenioAltType `json:"resource_type,omitempty"`
}

// PersonOrOrg holds either a person or corporate entity information.
type PersonOrOrg struct {
	// ID can be for a Person (e.g. Doiel-R-S) or Corporate entity
	ID string `json:"id,omitempty" xml:"id,omitempty"`

	// Type holds the type of Agent
	Type string `json:"type,omitempty"`

	// GivenName holds a peron's given name, e.g. Jane
	GivenName string `json:"given_name,omitempty" xml:"given_name,omitempty"`
	// FamilyName holds a person's family name, e.g. Doe
	FamilyName string `json:"family_name,omitempty" xml:"family_name,omitempty"`
	// Name holds a corporate name, e.g. The Unseen University
	Name string `json:"name,omitempty" xml:"name,omitempty"`

	// Identifiers holds a list of unique ID like ORCID, GND, ROR, ISNI
	Identifiers []*Identifier `json:"identifiers,omitempty"`
}

type Role struct {
}

type Affiliation struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Creator struct {
	PersonOrOrg *PersonOrOrg   `json:"person_or_org,omitempty"`
	Role        *Role          `json:"role,omitempty"`
	Affiliation []*Affiliation `json:"affiliation,omitempty"`
}

// Funder holds information for funding organizations
type Funder struct {
	XMLName     xml.Name `json:"-"`
	Name        string   `json:"name,omitempty" xml:"name,omitempty"`
	Description string   `json:"description,omitempty" xml:"description,omitempty"`
	GrantNumber string   `json:"grant_number,omitempty" xml:"grant_number,omitempty"`
	ROR         string   `json:"ror,omitempty" xml:"ror,omitempty"`
}

// ResourceURL
type ResourceURL struct {
	XMLName     xml.Name `json:"-"`
	Type        string   `json:"type,omitempty" xml:"type,omitempty"`
	Url         string   `json:"url" xml:"url"`
	Description string   `json:"description,omitempty" xml:"description,omitempty"`
}

// AltTitle is an Invenio 3 additional title type
type AltTitle struct {
	Title string       `json:"title,omitempty"`
	Type  *InvenioType `json:"type,omitempty"`
	Lang  *InvenioType `json:"lang,omitempty"`
}

// Subject isan Invenio 3 subject type
type Subject struct {
	Subject string `json:"subject,omitempty"`
	ID      string `json:"id,omitempty"`
}

// AltDescription  is an Invenio 3 alternate description
type AltDescription struct {
	Description string       `json:"description,omitempty"`
	Type        *InvenioType `json:"type,omitempty"`
	Lang        *InvenioType `json:"lang,omitempty"`
}

type Right struct {
	ID          string `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Link        string `json:"link,omitempty"`
}

type DateType struct {
	Date        string       `json:"date,omitempty"`
	Type        *InvenioType `json:"type,omitempty"`
	Description string       `json:"description,omitempty"`
}

// Size
type Size struct {
}

// Format
type Format struct {
}

// Location
type Location struct {
}

// FileType
type FileType struct {
}

// Tumbstone
type Tumbstone struct {
}

// Metadata is an indivudal eprint record optimize for ingest by
// using in Invenio 3
type Metadata struct {
	// General fields
	XMLName      xml.Name     `json:"-"`
	EPrintID     string       `json:"eprint_id,omitempty"`
	Collection   string       `json:"collection,omitempty"`
	EPrintType   string       `json:"eprint_type,omitempty" xml:"type,omitempty"`
	ResourceType *InvenioType `json:"resource_type,omitempty"`

	Title          string            `json:"title" xml:"title"`
	AltTitles      []*AltTitle       `json:"additional_titles"`
	Creators       []*Creator        `json:"creators,omitempty"`
	Contributors   []*Creator        `json:"contributors,omitempty"`
	Subjects       []*Subject        `json:"subjects,omitempty"`
	Languages      []*InvenioType    `json:"languages,omitempty"`
	Dates          []*DateType       `json:"dates,omitempty"`
	Version        []string          `json:"version,omitempty"`
	Publisher      string            `json:"publisher,omitempty" xml:"publisher,omitempty"`
	AltIdentifiers []*Identifier     `json:"additional_identifiers,omitempty"`
	Publication    string            `json:"publication,omitempty"`
	PubDate        string            `json:"publication_date,omitempty"`
	Description    string            `json:"description,omitempty" xml:"abstract,omitempty"`
	AltDescription []*AltDescription `json:"additional_descriptions,omitempty"`
	Rights         []*Right          `json:"rights,omitempty"`
	Sizes          []*Size           `json:"sizes,omitempty"`
	Formats        []*Format         `json:"formats,omitempty"`
	Locations      []*Location       `json:"locations,omitempty"`

	Funders []*Funder `json:"funder,omitempty"`

	Files     []*FileType `json:"files,omitempty"`
	Tumbstone *Tumbstone  `json:"tumbstone,omitempty"`

	PlaceOfPublication       string `json:"place_of_publication,omitempty"`
	Edition                  string `json:"edition,omitempty"`
	BookTitle                string `json:"book_title,omitempty"`
	Series                   string `json:"series,omitempty"`
	Volume                   string `json:"volume,omitempty"`
	Number                   string `json:"number,omitempty"`
	Refereed                 bool   `json:"refereed,omitempty"`
	Department               string `json:"department,omitempty"`
	Group                    string `json:"group,omitempty"`
	OtherNumberingSystemName string `json:"other_numbering_system_name,omitempty"`
	OtherNumberingSystemID   string `json:"other_numbering_system_id,omitempty"`
	Created                  string `json:"created,omitempty"`
	Updated                  string `json:"updated,omitempty"`
	Status                   string `json:"status"`
	//FIXME: Eprints stores the numeric id, we need a name or username to populate Username
	Username       string   `json:"username,omitempty"`
	FullTextStatus string   `json:"full_text_status,omitempty"`
	Notes          string   `json:"note,omitempty"`
	Keywords       []string `json:"keywords,omitempty"`

	// Patent oriented fields
	PatentApplicant      string                    `json:"patent_applicant,omitempty"`
	PatentNumber         string                    `json:"patent_number,omitempty"`
	PatentAssignee       []*PersonOrOrg            `json:"patent_assignee,omitempty"`
	PatentClassification []*map[string]interface{} `json:"patent_classification,omitempty"`
	RelatedPatents       []*map[string]interface{} `json:"related_patents,omitempty"`

	// Thesis oriented fields
	Divisions              []string `json:"divisions,omitempty"`
	Institution            string   `json:"institution,omitempty"`
	ThesisType             string   `json:"thesis_type,omitempty"`
	ThesisDegree           string   `json:"thesis_degree,omitempty"`
	ThesisDegreeGrantor    string   `json:"thesis_degree_grantor,omitempty"`
	ThesisDegreeDate       string   `json:"thesis_degree_date,omitempty"`
	ThesisSubmittedDate    string   `json:"thesis_submit_date,omitempty"`
	ThesisDefenseDate      string   `json:"thesis_defense_date,omitempty"`
	ThesisApprovedDate     string   `json:"thesis_approved_date,omitempty"`
	ThesisPublicDate       string   `json:"thesis_public_date,omitempty"`
	ThesisAuthorEMail      string   `json:"thesis_author_email,omitempty"`
	HideThesisAuthorEMail  string   `json:"hide_thesis_author_email,omitempty"`
	GradOfficeApprovalDate string   `json:"gradofc_approval_date,omitempty"`
	ThesisAwards           string   `json:"thesis_awards,omitempty"`
	ReviewStatus           string   `json:"review_status,omitempty"`
	OptionMajor            []string `json:"option_major,omitempty"`
	OptionMinor            []string `json:"option_minor,omitempty"`
	CopyrightStatement     string   `json:"copyright_statement,omitempty"`
}

func MapEPrintToMetadata(mapObject map[string]interface{}) (*Metadata, error) {
	return nil, fmt.Errorf("MapEPrintToMetadata() not implemented")
}

// Simplify take a single EPrint struct and converts it to
// an SimplePrint structure.
func Simplify(eprint *EPrint) (*Metadata, error) {
	simple := new(Metadata)
	simple.EPrintID = fmt.Sprintf("%d", eprint.EPrintID)
	simple.Status = eprint.EPrintStatus
	simple.Collection = eprint.Collection
	simple.Title = eprint.Title
	simple.Description = eprint.Abstract
	simple.Creators = []*Creator{}
	if eprint.Creators != nil {
		//FIXME: map Creates to simple.Creators
	}
	if eprint.CorpCreators != nil {
		//FIXME: map Creates to simple.Creators
	}
	if eprint.ConfCreators != nil {
		//FIXME: map Creates to simple.Creators
	}
	if eprint.Contributors != nil {
		//FIXME: map Creates to simple.Creators
	}
	if eprint.Editors != nil {
		//FIXME: map Creates to simple.Creators
	}
	if eprint.PrimaryObject != nil {
		// FIXME: map into simple's Invenio Metadata model
	}
	if eprint.RelatedObjects != nil {
		// FIXME: map into simple's Invenio Metadata model
	}
	if eprint.Funders != nil {
		// FIXME: map to simple's funder model
	}
	if eprint.DOI != "" {
		// FIXME: map to simple's identifiers list
	}
	if eprint.ISBN != "" {
		// FIXME: map to simple's identifiers list
	}
	if eprint.ISSN != "" {
		// FIXME: map to simple's identifiers list
	}
	if eprint.RelatedURL != nil {
		// FIXME: map into simple's Invenio Metadata model
	}
	if eprint.FullTextStatus != "" {
		// FIXME: map into simple's Invenio Metadata model
	}
	if eprint.Publisher != "" {
		// FIXME: map into simple's Invenio Metadata model
	}
	if eprint.Publication != "" {
		// FIXME: map into simple's Invenio Metadata model
	}
	if eprint.PlaceOfPub != "" {
		// FIXME: map into simple's Invenio Metadata model
	}
	if eprint.BookTitle != "" {
		// FIXME: map into simple's Invenio Metadata model
	}
	if eprint.Edition != "" {
		// FIXME: map into simple's Invenio Metadata model
	}
	if eprint.Series != "" {
		// FIXME: map into simple's Invenio Metadata model
	}
	if eprint.Volume != "" {
		// FIXME: map into simple's Invenio Metadata model
	}
	if eprint.Number != "" {
		// FIXME: map into simple's Invenio Metadata model
	}
	if eprint.Refereed != "" {
		// FIXME: map into simple's Invenio Metadata model
	}
	if eprint.Department != "" {
		// FIXME: map into simple's Invenio Metadata model
	}
	if eprint.Divisions != nil {
		// FIXME: map into simple's Invenio Metadata model
	}
	if eprint.OptionMajor != nil {
		// FIXME: map into simple's Invenio Metadata model
	}
	if eprint.OptionMinor != nil {
		// FIXME: map into simple's Invenio Metadata model
	}
	if eprint.CopyrightStatement != "" {
		// FIXME: map into simple's Invenio Metadata model
	}
	return simple, nil
}

// SimplifyEPrints takes an EPrints struct and converts it to an array SimplePrint
// structure.
func SimplifyEPrints(eprints *EPrints) ([]*Metadata, error) {
	var simpleList []*Metadata
	for i, eprint := range eprints.EPrint {
		if simple, err := Simplify(eprint); err != nil {
			return nil, fmt.Errorf("EPrint simplification (%d) error, %s", i, err)
		} else {
			simpleList = append(simpleList, simple)
		}
	}
	return simpleList, nil
}
