package eprinttools

/**
 * This file implements an intermediate metadata representation
 * suitable for moving data between EPrints 3.3 and Invenio-RDM 11.
 *
 * See documentation and example on Invenio's structured data:
 *
 * - https://inveniordm.docs.cern.ch/reference/metadata/
 * - https://github.com/caltechlibrary/caltechdata_api/blob/ce16c6856eb7f6424db65c1b06de741bbcaee2c8/tests/conftest.py#L147
 *
 */

import (
	"encoding/json"
	"time"
)

//
// Top Level Elements
//

// Record implements the top level Invenio 3 record structure
type Record struct {
	Schema       string                           `json:"$schema,omitempty"`
	ID           string                           `json:"id"`                  // Interneral persistent identifier for a specific version.
	PID          map[string]interface{}           `json:"pid,omitempty"`       // Interneral persistent identifier for a specific version.
	Parent       *RecordIdentifier                `json:"parent"`              // The internal persistent identifier for ALL versions.
	ExternalPIDs map[string]*PersistentIdentifier `json:"pids,omitempty"`      // System-managed external persistent identifiers (DOI, Handles, OAI-PMH identifiers)
	RecordAccess *RecordAccess                    `json:"access,omitempty"`    // Access control for record
	Metadata     *Metadata                        `json:"metadata"`            // Descriptive metadata for the resource
	Files        *Files                           `json:"files,omitempty"`     // Associated files information.
	Tombstone    *Tombstone                       `json:"tombstone,omitempty"` // Tombstone (deasscession) information.
	Created      time.Time                        `json:"created"`             // create time for record
	Updated      time.Time                        `json:"updated"`             // modified time for record

	// Annotation this is where I'm going to map custom fields.
	Annoations map[string]interface{} `json:"annotations,omitempty"`
}

//
// Second level Elements
//

// RecordIdentifier implements the scheme of "parent", a persistant
// identifier to the record.
type RecordIdentifier struct {
	ID     string  `json:"id"`               // The identifier of the parent record
	Access *Access `json:"access,omitempty"` // Access details for the record as a whole
}

// PersistentIdentifier holds an Identifier, e.g. ORCID, ROR, ISNI, GND
type PersistentIdentifier struct {
	Identifier string `json:"identifier,omitempty"` // The identifier value
	Provider   string `json:"provider,omitempty"`   // The provider idenitifier used internally by the system
	Client     string `json:"client,omitempty"`     // The client identifier used for connecting with an external registration service.
}

// RecordAccess implements a datastructure used by Invenio 3 to
// control record level accesss, e.g. in the REST API.
type RecordAccess struct {
	Record  string   `json:"record,omitempty"`  // "public" or "restricted. Read access to the record.
	Files   string   `json:"files,omitempty"`   // "public" or "restricted". Read access to the record's files.
	Embargo *Embargo `json:"embargo,omitempty"` // Embargo options for the record.
}

// Metadata holds the primary metadata about the record. This
// is where most of the EPrints 3.3.x data is mapped into.
type Metadata struct {
	ResourceType           map[string]string    `json:"resource_type,omitempty"` // Resource type id from the controlled vocabulary.
	Creators               []*Creator           `json:"creators,omitempty"`      //list of creator information (person or organization)
	Title                  string               `json:"title"`
	PublicationDate        string               `json:"publication_date,omitempty"`
	AdditionalTitles       []*TitleDetail       `json:"additional_titles,omitempty"`
	Description            string               `json:"description,omitempty"`
	AdditionalDescriptions []*Description       `json:"additional_descriptions,omitempty"`
	Rights                 []*Right             `json:"rights,omitempty"`
	Contributors           []*Creator           `json:"contributors,omitempty"`
	Subjects               []*Subject           `json:"subjects,omitempty"`
	Languages              []*map[string]string `json:"languages,omitempty"`
	Dates                  []*DateType          `json:"dates,omitempty"`
	Version                string               `json:"version,omitempty"`
	Publisher              string               `json:"publisher,omitempty"`
	Identifiers            []*Identifier        `json:"identifier,omitempty"`
	Funding                []*Funder            `json:"funding,omitempty"`

	/*
		// Extended  is where I am putting important
		// EPrint XML fields that don't clearly map.
		Extended map[string]*interface{} `json:"extended,omitempty"`
	*/
}

// Files
type Files struct {
	Enabled        bool                    `json:"enabled,omitempty"`
	Entries        map[string]*Entry       `json:"entries,omitempty"`
	DefaultPreview string                  `json:"default_preview,omitempty"`
	Sizes          []string                `json:"sizes,omitempty"`
	Formats        []string                `json:"formats,omitempty"`
	Locations      map[string]*interface{} `json:"locations,omitempty"`
}

type Entry struct {
	BucketID     string `json:"bucket_id,omitempty"`
	VersionID    string `json:"version_id,omitempty"`
	FileID       string `json:"file_id,omitempty"`
	Backend      string `json:"backend,omitempty"`
	StorageClass string `json:"storage_class,omitempty"`
	Key          string `json:"key,omitempty"`
	MimeType     string `json:"mimetype,omitempty"`
	Size         int    `json:"size,omitempty"`
	CheckSum     string `json:"checksum,omitempty"`
}

// Tombstone
type Tombstone struct {
	Reason    string    `json:"reason,omitempty"`
	Category  string    `json:"category,omitempty"`
	RemovedBy *User     `json:"removed_by,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

//
// Third/Fourth Level Elements
//

// Access is a third level element used by PersistentIdenitifier to
// describe access ownership of the record.
type Access struct {
	OwnedBy []*User `json:"owned_by,omitempty"`
}

// User is a data structured used in Access to describe record
// ownership or user actions.
type User struct {
	User        int    `json:"user,omitempty"`         // User (integer) identifier
	DisplayName string `json:"display_name,omitempty"` // This is my field to quickly associate the internal integer user id with a name for reporting and display.
	Email       string `json:"email,omitempty"`        // This is my field to quickly display a concact email associated with the integer user id.
}

// Embargo is a third level element used by RecordAccess to describe
// the embargo status of a record.
type Embargo struct {
	Active bool   `json:"active,omitempty"` // boolean, is the record under an embargo or not.
	Until  string `json:"until,omitempty"`  // Required if active true. ISO date string. When to lift the embargo. e.g. "2100-10-01"
	Reason string `json:"reason,omitempty"` // Explanation for the embargo
}

//
// Third level elements used in Metadata data structures
//

// Creator of a record's object
type Creator struct {
	PersonOrOrg  *PersonOrOrg   `json:"person_or_org,omitempty"` // The person or organization.
	Role         *Role          `json:"role,omitempty"`          // The role of the person or organization selected from a customizable controlled vocabularly.
	Affiliations []*Affiliation `json:"affiliations,omitempty"`  // Affiliations if `PersonOrOrg.Type` is personal.
}

// Role is an object describing a relationship to authorship
type Role struct {
	ID    string            `json:"id,omitempty"`
	Title string            `json:"title,omitempty"`
	Props map[string]string `json:"props,omitempty"`
}

// PersonOrOrg holds either a person or corporate entity information
// for the creators associated with the record.
type PersonOrOrg struct {
	ID   string `json:"cl_identifier,omitempty"` // The Caltech Library internal person or organizational identifier used to cross walk data across library systems. (this is not part of Invenion 3)
	Type string `json:"type,omitempty"`          // The type of name. Either "personal" or "organizational".

	GivenName  string `json:"given_name,omitempty" xml:"given_name,omitempty"`   // GivenName holds a peron's given name, e.g. Jane
	FamilyName string `json:"family_name,omitempty" xml:"family_name,omitempty"` // FamilyName holds a person's family name, e.g. Doe
	Name       string `json:"name,omitempty" xml:"name,omitempty"`               // Name holds a corporate name, e.g. The Unseen University

	// Identifiers holds a list of unique ID like ORCID, GND, ROR, ISNI
	Identifiers []*Identifier `json:"identifiers,omitempty"`
}

// Affiliation describes how a person or organization is affialated
// for the purpose of the record.
type Affiliation struct {
	ID   string `json:"id,omitempty"`   // The organizational or institutional id from the controlled vocabularly
	Name string `json:"name,omitempty"` // The name of the organization or institution
}

// Identifier holds an Identifier, e.g. ORCID, ROR, ISNI, GND
// for a person for organization it holds GRID, ROR. etc.
type Identifier struct {
	Scheme       string      `json:"scheme,omitempty"`
	Name         string      `json:"name,omitempty"`
	Title        string      `json:"title,omitempty"`
	Number       string      `json:"number,omitempty"`
	Identifier   string      `json:"identifier,omitempty"`
	RelationType *TypeDetail `json:"relation_type,omitempty"`
	ResourceType *TypeDetail `json:"resource_type,omitempty"`
}

// Type is an Invenio 3 e.g. ResourceType, title type or language
type Type struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Title string `json:"title,omitempty"`
}

// TitleDetail is used by AdditionalTitles in Metadata.
type TitleDetail struct {
	Title string `json:"title,omitempty"`
	Type  *Type  `json:"type,omitempty"`
	Lang  *Type  `json:"lang,omitempty"`
}

// Description holds additional descriptions in Metadata
// element. e.g. language versions of Abstract, etc.
type Description struct {
	Description string `json:"description,omitempty"`
	Type        *Type  `json:"type,omitempty"`
	Lang        *Type  `json:"lang,omitempty"`
}

// Right holds a specific Rights element for the Metadata's
// list of Rights.
//
// NOTE: for REST API lookup by ID or Title (but not both) should
// be supported at the same end point. I.e. they both must be unique
// with in their set of field values.
type Right struct {
	ID          string       `json:"id,omitempty"`          // Identifier value
	Title       *TitleDetail `json:"title,omitempty"`       // Localized human readable title e.g., `{"en": "The ACME Corporation License."}`.
	Description *Description `json:"description,omitempty"` // Localized license description text e.g., `{"en":"This license ..."}`.
	Link        string       `json:"link,omitempty"`        // Link to full license.
}

// Subject element holds one of a list of subjects
// in the Metadata element.
type Subject struct {
	Subject string `json:"subject,omitempty"`
	ID      string `json:"id,omitempty"`
}

// DateType holds Invenio dates used in Metadata element.
type DateType struct {
	Date        string `json:"date,omitempty"`
	Type        *Type  `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
}

// Funder holds funding information for funding organizations in Metadata
type Funder struct {
	Funder    *Identifier   `json:"funder,omitempty"`
	Award     *Identifier   `json:"award,omitempty"`
	Reference []*Identifier `json:"references,omitempty"`
}

//
// Additional fourth level elements
//

// Type is an alternate expression of a type where title is map
// with additional info like language. It is used to describe relationships
// and resources in Identifiers. It is a variation of Type.
type TypeDetail struct {
	ID    string            `json:"id,omitempty"`
	Name  string            `json:"name,omitempty"`
	Title map[string]string `json:"title,omitempty"`
}

func (rec *Record) ToString() []byte {
	src, _ := json.MarshalIndent(rec, "", "    ")
	return src
}
