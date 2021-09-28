/**
 * simplified presents an Invenio 3 like JSON representation of an EPrint
 * record.
 */
package eprinttools

import (
	"encoding/xml"
	"fmt"
	"strings"
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

// Metadata is an indivudal eprint record optimize for ingest by
// using in Invenio or Solr 8.9.0.
type Matadata struct {
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
	PatentAssignee       []*Agent                  `json:"patent_assignee,omitempty"`
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

func MapObjectToMetadata(mapObject map[string]interface{}) (*Metadata, error) {
	dObject := new(Metadata)
	foundContent := false
	for k, v := range mapObject {
		switch k {
		case "basename":
			dObject.Basename = v.(string)
			foundContent = true
		case "content":
			dObject.Content = v.(string)
			foundContent = true
		case "filesize":
			dObject.FileSize = v.(int)
			foundContent = true
		case "license":
			dObject.License = v.(string)
		case "mime_type":
			dObject.MimeType = v.(string)
			foundContent = true
		case "url":
			dObject.Url = v.(string)
			foundContent = true
		case "version":
			dObject.Version = v.(string)
			foundContent = true
		}
	}
	if !foundContent {
		return nil, fmt.Errorf("No digital object attributes found")
	}
	return dObject, nil
}

// Simplify take a single EPrint struct and converts it to
// an SimplePrint structure.
func Simplify(eprint *EPrint) (*Metadata, error) {
	simple := new(SimplePrint)
	simple.EPrintID = fmt.Sprintf("%d", eprint.EPrintID)
	simple.Status = eprint.EPrintStatus
	simple.Collection = eprint.Collection
	simple.Type = eprint.Type
	simple.Title = eprint.Title
	simple.Abstract = eprint.Abstract
	simple.Creators = []*Agent{}
	simple.CorpCreators = []*Agent{}
	simple.ConfCreators = []*Agent{}
	simple.Contributors = []*Agent{}
	simple.Editors = []*Agent{}
	simple.Committee = []*Agent{}
	simple.Advisors = []*Agent{}
	if eprint.Creators != nil {
		for _, creator := range eprint.Creators.Items {
			agent := new(Agent)
			agent.ID = creator.ID
			agent.FamilyName = creator.Name.Family
			agent.GivenName = creator.Name.Given
			agent.ORCID = creator.ORCID
			simple.Creators = append(simple.Creators, agent)
		}
	}
	if eprint.CorpCreators != nil {
		for _, creator := range eprint.CorpCreators.Items {
			agent := new(Agent)
			agent.ID = creator.ID
			//FIXME: Need to figure out where I record actual CorpCreator name
			agent.Name = creator.Value
			simple.CorpCreators = append(simple.CorpCreators, agent)
		}
	}
	if eprint.ConfCreators != nil {
		for _, creator := range eprint.ConfCreators.Items {
			agent := new(Agent)
			agent.ID = creator.ID
			//FIXME: Need to figure out where I record actual CorpCreator name
			agent.Name = creator.Value
			simple.ConfCreators = append(simple.ConfCreators, agent)
		}
	}
	if eprint.Contributors != nil {
		for _, contributor := range eprint.Contributors.Items {
			agent := new(Agent)
			agent.ID = contributor.ID
			agent.FamilyName = contributor.Name.Family
			agent.GivenName = contributor.Name.Given
			agent.ORCID = contributor.ORCID
			simple.Contributors = append(simple.Contributors, agent)
		}
	}
	if eprint.Editors != nil {
		for _, editor := range eprint.Editors.Items {
			agent := new(Agent)
			agent.ID = editor.ID
			agent.FamilyName = editor.Name.Family
			agent.GivenName = editor.Name.Given
			agent.ORCID = editor.ORCID
			simple.Editors = append(simple.Editors, agent)
		}
	}
	if eprint.PrimaryObject != nil {
		if dObj, err := MapObjectToDigitalObject(eprint.PrimaryObject); err == nil {
			simple.PrimaryObject = dObj
		}
	}
	if eprint.RelatedObjects != nil {
		simple.RelatedObjects = []*DigitalObject{}
		for _, obj := range eprint.RelatedObjects {
			if sObj, err := MapObjectToDigitalObject(obj); err == nil {
				simple.RelatedObjects = append(simple.RelatedObjects, sObj)
			}
		}
	}
	if eprint.Funders != nil {
		for _, item := range eprint.Funders.Items {
			funder := new(Funder)
			funder.Name = item.Agency
			funder.Description = item.Description
			funder.GrantNumber = item.GrantNumber
			//FIXME: need to determine FinderID and ROR if possible.
			simple.Funders = append(simple.Funders, funder)
		}
	}
	if eprint.DOI != "" {
		simple.DOI = eprint.DOI
	}
	if eprint.ISBN != "" {
		simple.ISBN = eprint.ISBN
	}
	if eprint.ISSN != "" {
		simple.ISSN = eprint.ISSN
	}
	if eprint.RelatedURL != nil {
		for _, item := range eprint.RelatedURL.Items {
			resource := new(ResourceURL)
			resource.Type = item.Type
			resource.Description = item.Description
			resource.Url = item.URL
			simple.RelatedURLs = append(simple.RelatedURLs, resource)
		}
	}
	if eprint.FullTextStatus != "" {
		simple.FullTextStatus = eprint.FullTextStatus
	}
	if eprint.Publisher != "" {
		simple.Publisher = eprint.Publisher
	}
	if eprint.Publication != "" {
		simple.Publication = eprint.Publication
	}
	if eprint.PlaceOfPub != "" {
		simple.PlaceOfPublication = eprint.PlaceOfPub
	}
	if eprint.BookTitle != "" {
		simple.BookTitle = eprint.BookTitle
	}
	if eprint.Edition != "" {
		simple.Edition = eprint.Edition
	}
	if eprint.Series != "" {
		simple.Series = eprint.Series
	}
	if eprint.Volume != "" {
		simple.Volume = eprint.Volume
	}
	if eprint.Number != "" {
		simple.Number = eprint.Number
	}
	if eprint.Refereed != "" {
		if strings.ToLower(eprint.Refereed) == "true" {
			simple.Refereed = true
		} else {
			simple.Refereed = false
		}
	}
	if eprint.Department != "" {
		simple.Department = eprint.Department
	}
	if eprint.Divisions != nil {
		for _, division := range eprint.Divisions.Items {
			simple.Divisions = append(simple.Divisions, division.Value)
		}
	}
	if eprint.OptionMajor != nil {
		for _, option := range eprint.OptionMajor.Items {
			simple.OptionMajor = append(simple.OptionMajor, option.Value)
		}
	}
	if eprint.OptionMinor != nil {
		for _, option := range eprint.OptionMinor.Items {
			simple.OptionMinor = append(simple.OptionMinor, option.Value)
		}
	}
	if eprint.CopyrightStatement != "" {
		simple.CopyrightStatement = eprint.CopyrightStatement
	}
	return simple, nil
}

// SimplifyEPrints takes an EPrints struct and converts it to an array SimplePrint
// structure.
func SimplifyEPrints(eprints *EPrints) ([]*Metadata, error) {
	var simpleList []*SimplePrint
	for i, eprint := range eprints.EPrint {
		if simple, err := Simplify(eprint); err != nil {
			return nil, fmt.Errorf("EPrint simplification (%d) error, %s", i, err)
		} else {
			simpleList = append(simpleList, simple)
		}
	}
	return simpleList, nil
}
