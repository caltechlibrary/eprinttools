package eprinttools

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// DigitalObject describes the digital material associated with
// the EPrint record.
type DigitalObject struct {
	XMLName  xml.Name `json:"-"`
	EPrintID string   `json:"eprint_id,omitempty" xml:"eprint_id,omitempty"`
	Content  string   `json:"content,omitempty" xml:"content,omitempty"`
	Basename string   `json:"basename,omitempty" xml:"basename,omitempty"`
	FileSize int      `json:"file_size,omitempty" xml:"file_size,omitempty"`
	License  string   `json:"license,omitempty" xml:"license,omitempty"`
	MimeType string   `json:"mime_type,omitempty" xml:"mime_type,omitempty"`
	Url      string   `json:"url,omitempty" xml:"url,omitempty"`
	Version  string   `json:"version,omitempty" xml:"version,omitempty"`
}

// Agent holds either a person or corporate entity information.
type Agent struct {
	XMLName xml.Name `json:"-"`
	// Type can be "person" or "corporate"
	Type string `json:"type,omitempty" xml:"type,omitempty"`
	// ID can be for a Person (e.g. Doiel-R-S) or Corporate entity
	ID string `json:"id,omitempty" xml:"id,omitempty"`

	// Name holds a corporate name, e.g. The Unseen University
	Name string `json:"name,omitempty" xml:"name,omitempty"`
	// GivenName holds a peron's given name, e.g. Jane
	GivenName string `json:"given_name,omitempty" xml:"given_name,omitempty"`
	// FamilyName holds a person's family name, e.g. Doe
	FamilyName  string `json:"family_name,omitempty" xml:"family_name,omitempty"`
	ORCID       string `json:"orcid,omitempty" xml:"orcid,omitempty"`
	ROR         string `json:"ror,omitempty" xml:"ror,omitempty"`
	Affiliation string `json:"affiliation,omitempty" xml:"affiliation,omitempty"`
	Role        string `json:"role,omitempty" xml:"role,omitempty"`
	EMail       string `json:"email,omitempty" xml:"email,omitempty"`
	Website     string `json:"website,omitempty" xml:"website,omitempty"`
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

// SimplePrint is an indivudal eprint record optimize for ingest by
// using in search, e.g. Solr 8.9.0.
type SimplePrint struct {
	// General fields
	XMLName      xml.Name `json:"-"`
	EPrintID     string   `json:"eprint_id" xml:"eprint_id"`
	Collection   string   `json:"collection,omitempty" xml:"collection,omitempty"`
	Type         string   `json:"type,omitempty" xml:"type,omitempty"`
	Title        string   `json:"title" xml:"title"`
	Creators     []*Agent `json:"creators,omitempty" xml:"creators>creator,omitempty"`
	Contributors []*Agent `json:"contributors,omitempty" xml:"contributors>cotributor,omitempty"`
	Editors      []*Agent `json:"editors,omitempty" xml:"editors>editor,omitempty"`
	Committee    []*Agent `json:"committee,omitempty" xml:"committee>member,omitempty"`
	Advisors     []*Agent `json:"advisors,omitempty" xml:"advisors>advisor,omitempty"`

	Funders                  []*Funder        `json:"funders,omitempty" xml:"funders>funder,omitempty"`
	Abstract                 string           `json:"abstract,omitempty" xml:"abstract,omitempty"`
	PrimaryObject            *DigitalObject   `json:"primary_object,omitempty" xml:"primary_object,omitempty"`
	RelatedObjects           []*DigitalObject `json:"related_objects,omitempty" xml:"related_objects>related_object,omitempty"`
	RelatedURLs              []*ResourceURL   `json:"related_urls,omitempty" xml:"related_urls>related_url,omitempty"`
	Publisher                string           `json:"publisher,omitempty" xml:"publisher,omitempty"`
	Publication              string           `json:"publication,omitempty" xml:"journal_or_pub_title,omitempty"`
	PlaceOfPublication       string           `json:"place_of_publication,omitempty" xml:"place_of_publication,omitempty"`
	Edition                  string           `json:"edition,omitempty" xml:"edition,omitempty"`
	BookTitle                string           `json:"book_title,omitempty" xml:"book_title,omitempty"`
	Series                   string           `json:"series,omitempty" xml:"series,omitempty"`
	Volume                   string           `json:"volume,omitempty" xml:"volume,omitempty"`
	Number                   string           `json:"number,omitempty" xml:"number,omitempty"`
	Refereed                 bool             `json:"refereed,omitempty" xml:"refereed,omitmpety"`
	Department               string           `json:"department,omitempty" xml:"department,omitempty"`
	Group                    string           `json:"group,omitempty" xml:"group,omitempty"`
	OtherNumberingSystemName string           `json:"other_numbering_system_name,omitempty" xml:"other_numbering_system_name,omitempty"`
	OtherNumberingSystemID   string           `json:"other_numbering_system_id,omitempty"`
	DOI                      string           `json:"doi,omitempty" xml:"doi,omitempty"`
	ISSN                     string           `json:"issn,omitempty" xml:"issn,omitempty"`
	ISBN                     string           `json:"isbn,omitempty" xml:"isbn,omitempty"`
	PubCentralID             string           `json:"pub_central_id,omitempty" xml:"pub_central_id,omitempty"`
	Created                  string           `json:"created,omitempty" xml:"created,omitempty"`
	Updated                  string           `json:"updated,omitempty" xml:"updated,omitempty"`
	PubDate                  string           `json:"pub_date,omitempty" xml:"pub_date,omitempty"`
	Status                   string           `json:"status" xml:"status"`
	//FIXME: Eprints stores the numeric id, we need a name or username to populate Username
	Username       string   `json:"username,omitempty" xml:"username,omitempty"`
	FullTextStatus string   `json:"full_text_status,omitempty" xml:"full_text_status,omitempty"`
	Notes          string   `json:"note,omitempty" xml:"note,omitempty"`
	Keywords       []string `json:"keywords,omitempty" xml:"keywords>keyword,omitempty"`
	Subject        []string `json:"subject,omitempty" xml:"subjects>subject,omitempty"`

	// Patent oriented fields
	PatentApplicant      string                    `json:"patent_applicant,omitempty" xml:"patent_applicant,omitempty"`
	PatentNumber         string                    `json:"patent_number,omitempty" xml:"patent_number,omitempty"`
	PatentAssignee       []*Agent                  `json:"patent_assignee,omitempty" xml:"patent_assignees>patent_assignee,omitempty"`
	PatentClassification []*map[string]interface{} `json:"patent_classification,omitempty" xml:"patent_classifications>patent_classification,omitempty"`
	RelatedPatents       []*map[string]interface{} `json:"related_patents,omitempty" xml:"related_patents>related_patent,omitempty"`

	// Thesis oriented fields
	Divisions              []string `json:"divisions,omitempty" xml:"divisions,omitemmpty"`
	Institution            string   `json:"institution,omitempty" xml:"institution,omitempty"`
	ThesisType             string   `json:"thesis_type,omitempty" xml:"thesis_type,omitempty"`
	ThesisDegree           string   `json:"thesis_degree,omitempty" xml:"thesis_degree,omitempty"`
	ThesisDegreeGrantor    string   `json:"thesis_degree_grantor,omitempty" xml:"thesis_degree_grantor,omitempty"`
	ThesisDegreeDate       string   `json:"thesis_degree_date,omitempty" xml:"thesis_degree_date,omitempty"`
	ThesisSubmittedDate    string   `json:"thesis_submit_date,omitempty" xml:"thesis_submit_date,omitempty"`
	ThesisDefenseDate      string   `json:"thesis_defense_date,omitempty" xml:"thesis_defense_date,omitempty"`
	ThesisApprovedDate     string   `json:"thesis_approved_date,omitempty" xml:"thesis_approved_date,omitempty"`
	ThesisPublicDate       string   `json:"thesis_public_date,omitempty" xml:"thesis_public_date,omitempty"`
	ThesisAuthorEMail      string   `json:"thesis_author_email,omitempty" xml:"thesis_author_email,omitempty"`
	HideThesisAuthorEMail  string   `json:"hide_thesis_author_email,omitempty" xml:"hide_thesis_author_email,omitempty"`
	GradOfficeApprovalDate string   `json:"gradofc_approval_date,omitempty" xml:"gradofc_approval_date,omitempty"`
	ThesisAwards           string   `json:"thesis_awards,omitempty" xml:"thesis_awards,omitempty"`
	ReviewStatus           string   `json:"review_status,omitempty" xml:"review_status,omitempty"`
	OptionMajor            []string `json:"option_major,omitempty" xml:"option_major,omitempty"`
	OptionMinor            []string `json:"option_minor,omitempty" xml:"option_minor,omitempty"`
	CopyrightStatement     string   `json:"copyright_statement,omitempty" xml:"copyright_statement,omitempty"`
}

func MapObjectToDigitalObject(mapObject map[string]interface{}) (*DigitalObject, error) {
	dObject := new(DigitalObject)
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
func Simplify(eprint *EPrint) (*SimplePrint, error) {
	simple := new(SimplePrint)
	simple.EPrintID = fmt.Sprintf("%d", eprint.EPrintID)
	simple.Status = eprint.EPrintStatus
	simple.Collection = eprint.Collection
	simple.Type = eprint.Type
	simple.Title = eprint.Title
	simple.Abstract = eprint.Abstract
	simple.Creators = []*Agent{}
	simple.Contributors = []*Agent{}
	simple.Editors = []*Agent{}
	simple.Committee = []*Agent{}
	simple.Advisors = []*Agent{}
	if eprint.Creators != nil {
		for _, creator := range eprint.Creators.Items {
			agent := new(Agent)
			agent.Type = "person"
			agent.ID = creator.ID
			agent.Role = "creator"
			agent.FamilyName = creator.Name.Family
			agent.GivenName = creator.Name.Given
			agent.ORCID = creator.ORCID
			simple.Creators = append(simple.Creators, agent)
		}
	}
	if eprint.CorpCreators != nil {
		for _, creator := range eprint.CorpCreators.Items {
			agent := new(Agent)
			agent.Type = "corporate"
			agent.ID = creator.ID
			agent.Role = "creator"
			//FIXME: Need to figure out where I record actual CorpCreator name
			agent.Name = creator.Value
			simple.Creators = append(simple.Creators, agent)
		}
	}
	if eprint.ConfCreators != nil {
		for _, creator := range eprint.ConfCreators.Items {
			agent := new(Agent)
			agent.Type = "conference"
			agent.Role = "creator"
			agent.ID = creator.ID
			//FIXME: Need to figure out where I record actual CorpCreator name
			agent.Name = creator.Value
			simple.Creators = append(simple.Creators, agent)
		}
	}
	if eprint.Contributors != nil {
		for _, contributor := range eprint.Contributors.Items {
			agent := new(Agent)
			agent.Type = "person"
			agent.Role = "contributor"
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
			agent.Type = "person"
			agent.Role = "editor"
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
func SimplifyEPrints(eprints *EPrints) ([]*SimplePrint, error) {
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
