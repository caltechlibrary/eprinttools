/**
 * simplified presents an Invenio 3 like JSON representation of an EPrint
 * record. This is intended to make the development of V2 of feeds easier
 * for both our audience on internal programming needs.
 *
 * See documentation and example on Invenio's structured data:
 *
 * - https://inveniordm.docs.cern.ch/reference/metadata/
 * - https://github.com/caltechlibrary/caltechdata_api/blob/ce16c6856eb7f6424db65c1b06de741bbcaee2c8/tests/conftest.py#L147
 *
 */
package eprinttools

import (
	"encoding/xml"
	"time"
)

// Invenio Record structure
type InvenioRecord struct {
	ID           string                           `json:"id"`
	PID          map[string]interface{}           `json:"pid,omitempty"`
	Parent       *RecordIdentifier                `json:"parent"`
	ExternalPIDs map[string]*PersistentIdentifier `json:"pids,omitempty"`
	Access       *RecordAccess                    `json:"access,omitempty"`
	Metadata     *Metadata                        `json:"metadata"`
	Files        *Files                           `json:"files"`
	Tombstone    *Tombstone                       `json:"tombstone,omitempty"`
	Created      time.Time                        `json:"created"`
	Updated      time.Time                        `json:"updated"`
}

// RecordIdentifier
type RecordIdentifier struct {
	ID     string  `json:"id"`
	Access *Access `json:"access,omitempty"`
}

// Access
type Access struct {
	OwnedBy []*Owner `json:"owned_by,omitempty"`
}

// Owner
type Owner struct {
	User int `json:"user,omitempty"`
}

// PersistentIdentifier holds an Identifier, e.g. ORCID, ROR, ISNI, GND
type PersistentIdentifier struct {
	Identifier string `json:"identifier,omitempty"`
	Provider   string `json:"provider,omitempty"`
	Client     string `json:"client,omitempty"`
}

// RecordAccess
type RecordAccess struct {
	Record  string   `json:"record,omitempty"`
	Files   string   `json:"files,omitempty"`
	Embargo *Embargo `json:"embargo,omitempty"`
}

// Embargo
type Embargo struct {
	Active bool   `json:"active,omitempty"`
	Until  string `json:"until,omitempty"`
	Reason string `json:"reason,omitempty"`
}

type Metadata struct {
	// ID is based on 10. Resource Type from DataCite
	ResourceType     map[string]string    `json:"resource_type,omitempty"`
	Creators         []*Creator           `jons:"creators,omitempty"`
	Title            string               `json:"title"`
	PublicationDate  string               `json:"publication_date,omitempty"`
	AdditionalTitles []*AltTitle          `json:"additional_titles,omitempty"`
	Description      string               `json:"description,omitempty"`
	AdditoinalDesc   []*AltDescription    `json:"additional_description,omitempty"`
	Rights           []*Right             `json:"rights,omitempty"`
	Contributors     []*Creator           `json:"contributors,omitempty"`
	Subjects         []*Subject           `json:"subjects,omitempty"`
	Languages        []*map[string]string `json:"languages,omitempty"`
	Dates            []*InvenioDate       `json:"dates,omitempty"`
	Version          string               `json:"version,omitempty"`
	Publisher        string               `json:"publisher,omitempty"`
	AlterIdentifiers []*Identifier        `json:"identifier,omitempty"`
	Funding          []*Funder            `json:"funding,omitempty"`
}

// Creator of a record's object
type Creator struct {
	PersonOrOrg  *PersonOrOrg   `json:"person_or_org,omitempty"`
	Role         string         `json:"role,omitempty"`
	Affiliations []*Affiliation `json:"affiliations,omitempty"`
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

type Affiliation struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Files
type Files struct {
	Enabled        bool                    `json:"enabled,omitempty"`
	Entries        map[string]*interface{} `json:"entries,omitempty"`
	DefaultPreview string                  `json:"default_preview,omitempty"`
}

// Tombstone
type Tombstone struct {
	Reason    string    `json:"reason,omitempty"`
	Category  string    `json:"category,omitempty"`
	RemovedBy *Owner    `json:"removed_by,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

// InvenioDate
type InvenioDate struct {
	Date        string          `json:"date,omitempty"`
	Type        *InvenioAltType `json:"type,omitempty"`
	Description string          `json:"description,omitempty"`
}

// InvenioType is an Invenio 3 e.g. ResourceType, title type or language
type InvenioType struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Title string `json:"title,omitempty"`
}

// InvenioAltType is an alternate expression of a type where title is map
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

// Funder holds information for funding organizations
type Funder struct {
	Funder map[string]string `json:"funder,omitempty"`
	Award  map[string]string `json:"award,omitempty"`
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
