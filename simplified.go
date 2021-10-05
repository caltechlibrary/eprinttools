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

/**
 * simplified.go implements an Invenio 3 like JSON representation
 * of an EPrint record. This is intended to make the development of
 * V2 of feeds easier for both our audience and in for our internal
 * programming needs.
 *
 * See documentation and example on Invenio's structured data:
 *
 * - https://inveniordm.docs.cern.ch/reference/metadata/
 * - https://github.com/caltechlibrary/caltechdata_api/blob/ce16c6856eb7f6424db65c1b06de741bbcaee2c8/tests/conftest.py#L147
 *
 */

import (
	"encoding/json"
	"fmt"
	"strings"
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
	Files        *Files                           `json:"files"`               // Associated files information.
	Tombstone    *Tombstone                       `json:"tombstone,omitempty"` // Tombstone (deasscession) information.
	Created      time.Time                        `json:"created"`             // create time for record
	Updated      time.Time                        `json:"updated"`             // modified time for record
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
	Creators               []*Creator           `jons:"creators,omitempty"`      //list of creator information (person or organization)
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
	Role         string         `json:"role,omitempty"`          // The role of the person or organization selected from a customizable controlled vocabularly.
	Affiliations []*Affiliation `json:"affiliations,omitempty"`  // Affiliations if `PersonOrOrg.Type` is personal.
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
	ID          string `json:"id,omitempty"`          // Identifier value
	Title       string `json:"title,omitempty"`       // Localized human readable title e.g., `{"en": "The ACME Corporation License."}`.
	Description string `json:"description,omitempty"` // Localized license description text e.g., `{"en":"This license ..."}`.
	Link        string `json:"link,omitempty"`        // Link to full license.
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
	Funder    []*Identifier `json:"funder,omitempty"`
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

//
// The following have not been implemented in current release (2021-10-01)
// of Invenio 3.
//

// Files.Sizes (e.g. bytes, pages, inches, etc.) or duration (extent), e.g.
// hours, minutes, days, etc., of a resource.
// This structure is compatible with 13. Size in DataCite
//
// e.g. `{ "sizes": [ "11 pages" ], }`
//

// Files.Format is technical format of the resource. This structure is
// compatible with 14. Format in DataCite.
//
// e.g. `{ "formats": [ "application/pdf" ], }`
//

// Location is compatible with 18. GeoLocation in DataCite Metadata Schema.

//
// EPrints/Simplified model handling
//

// RecordFromEPrint takes an EPrint structure and crosswalks it into
// the Record structure.
func CrosswalkEPrintToRecord(eprint *EPrint) (*Record, error) {
	rec := new(Record)
	rec.Schema = `local://records/record-v2.0.0.json`
	rec.ID = fmt.Sprintf("%s:%d", eprint.Collection, eprint.EPrintID)

	if err := rec.parentFromEPrint(eprint); err != nil {
		return rec, err
	}
	if err := rec.externalPIDFromEPrint(eprint); err != nil {
		return rec, err
	}
	if err := rec.recordAccessFromEPrint(eprint); err != nil {
		return rec, err
	}

	if err := rec.metadataFromEPrint(eprint); err != nil {
		return rec, err
	}
	if err := rec.filesFromEPrint(eprint); err != nil {
		return rec, err
	}

	if eprint.EPrintStatus == "deletion" {
		if err := rec.tombstoneFromEPrint(eprint); err != nil {
			return rec, err
		}
	}

	if err := rec.createdUpdatedFromEPrint(eprint); err != nil {
		return rec, err
	}
	if err := rec.pidFromEPrint(eprint); err != nil {
		return rec, err
	}
	return rec, nil
}

// PIDFromEPrint crosswalks the PID from an EPrint record.
func (rec *Record) pidFromEPrint(eprint *EPrint) error {
	data := map[string]interface{}{}
	src := fmt.Sprintf(`{
"id": %d,
"pid": { "eprint": "eprintid"}
}`, eprint.EPrintID)
	err := json.Unmarshal([]byte(src), &data)
	if err != nil {
		return fmt.Errorf("Cannot generate PID from EPrint %d", eprint.EPrintID)
	}
	rec.PID = data
	return nil
}

// parentFromEPrint crosswalks the Perent unique ID from EPrint record.
func (rec *Record) parentFromEPrint(eprint *EPrint) error {
	parent := new(RecordIdentifier)
	parent.ID = fmt.Sprintf("%s:%d", eprint.Collection, eprint.EPrintID)
	parent.Access = new(Access)
	ownedBy := new(User)
	ownedBy.User = eprint.UserID
	ownedBy.DisplayName = eprint.Reviewer
	parent.Access.OwnedBy = append(parent.Access.OwnedBy, ownedBy)
	rec.Parent = parent
	return nil
}

// externalPIDFromEPrint aggregates all the external identifiers
// from the EPrint record into Record
func (rec *Record) externalPIDFromEPrint(eprint *EPrint) error {
	rec.ExternalPIDs = map[string]*PersistentIdentifier{}
	// Pickup DOI
	if eprint.DOI != "" {
		pid := new(PersistentIdentifier)
		pid.Identifier = eprint.DOI
		pid.Provider = "datacite" // FIXME: should be DataCite or CrossRef
		pid.Client = ""           // FIXME: need to find out client string
		rec.ExternalPIDs["doi"] = pid
	}
	// Pickup ISSN
	if eprint.ISBN != "" {
		pid := new(PersistentIdentifier)
		pid.Identifier = eprint.ISSN
		pid.Provider = "" // FIXME: Need to find out identifier string
		pid.Client = ""   // FIXME: need to find out client string
		rec.ExternalPIDs["ISSN"] = pid
	}
	// Pickup ISBN
	if eprint.ISBN != "" {
		pid := new(PersistentIdentifier)
		pid.Identifier = eprint.ISBN
		pid.Provider = "" // FIXME: Need to find out identifier string
		pid.Client = ""   // FIXME: need to find out client string
		rec.ExternalPIDs["ISBN"] = pid
	}
	//FIXME: figure out if we have other persistent identifiers
	//scattered in the EPrints XML and map them.
	return nil
}

// recordAccessFromEPrint extracts access permissions from the EPrint
func (rec *Record) recordAccessFromEPrint(eprint *EPrint) error {
	isPublic := true
	if (eprint.ReviewStatus == "review") ||
		(eprint.ReviewStatus == "withheld") ||
		(eprint.ReviewStatus == "gradoffice") ||
		(eprint.ReviewStatus == "notapproved") {
		isPublic = false
	}
	if eprint.EPrintStatus != "archive" || eprint.MetadataVisibility != "show" {
		isPublic = false
	}
	rec.RecordAccess = new(RecordAccess)
	if isPublic {
		rec.RecordAccess.Record = "public"
	} else {
		rec.RecordAccess.Record = "restricted"
	}
	// Need to make sure record is not embargoed
	for _, doc := range *eprint.Documents {
		if doc.DateEmbargo != "" {
			embargo := new(Embargo)
			embargo.Until = doc.DateEmbargo
			embargo.Reason = eprint.Suggestions
			if doc.Security == "internal" {
				embargo.Active = true
			} else {
				embargo.Active = false
			}
			rec.RecordAccess.Embargo = embargo
			break
		}
	}
	return nil
}

func creatorFromItem(item *Item, objType string, objRole string, objIdType string) *Creator {
	person := new(PersonOrOrg)
	person.Type = objType
	if item.Name != nil {
		person.FamilyName = item.Name.Family
		person.GivenName = item.Name.Given
	}
	if item.ORCID != "" {
		identifier := new(Identifier)
		identifier.Scheme = "ORCID"
		identifier.Identifier = item.ORCID
		person.Identifiers = append(person.Identifiers, identifier)
	}
	if item.ID != "" {
		identifier := new(Identifier)
		identifier.Scheme = objIdType
		identifier.Identifier = item.ID
		person.Identifiers = append(person.Identifiers, identifier)
	}
	creator := new(Creator)
	creator.PersonOrOrg = person
	creator.Role = objRole

	return creator
}

func dateTypeFromTimestamp(dtType string, timestamp string, description string) *DateType {
	dt := new(DateType)
	dt.Type = new(Type)
	dt.Type.ID = dtType
	dt.Type.Title = dtType
	dt.Description = description
	if len(timestamp) > 9 {
		dt.Date = timestamp[0:10]
	} else {
		dt.Date = timestamp
	}
	return dt
}

func mkSimpleIdentifier(scheme, value string) *Identifier {
	identifier := new(Identifier)
	identifier.Scheme = scheme
	identifier.Identifier = value
	return identifier
}

func funderFromItem(item *Item) *Funder {
	funder := new(Funder)
	if item.GrantNumber != "" {
		funder.Award = new(Identifier)
		funder.Award.Number = item.GrantNumber
		funder.Award.Scheme = "eprints_grant_number"
	}
	if item.Agency != "" {
		org := new(Identifier)
		org.Name = item.Agency
		org.Scheme = "eprints_agency"
		funder.Funder = append(funder.Funder, org)
	}
	return funder
}

// metadataFromEPrint extracts metadata from the EPrint record
func (rec *Record) metadataFromEPrint(eprint *EPrint) error {
	metadata := new(Metadata)
	metadata.ResourceType = map[string]string{}
	metadata.ResourceType["id"] = eprint.Type
	// NOTE: Creators get listed in the citation, Contributors do not.
	if (eprint.Creators != nil) && (eprint.Creators.Items != nil) {
		for _, item := range eprint.Creators.Items {
			creator := creatorFromItem(item, "person", "creator", "creator_id")
			metadata.Creators = append(metadata.Creators, creator)
		}
	}
	if (eprint.CorpCreators != nil) && (eprint.CorpCreators.Items != nil) {
		for _, item := range eprint.CorpCreators.Items {
			creator := creatorFromItem(item, "organization", "corporate_creator", "organization_id")
			metadata.Creators = append(metadata.Creators, creator)
		}
	}
	if (eprint.Contributors != nil) && (eprint.Contributors.Items != nil) {
		for _, item := range eprint.Contributors.Items {
			contributor := creatorFromItem(item, "person", "contributor", "contributor_id")
			metadata.Contributors = append(metadata.Contributors, contributor)
		}
	}
	if (eprint.CorpContributors != nil) && (eprint.CorpContributors.Items != nil) {
		for _, item := range eprint.CorpContributors.Items {
			contributor := creatorFromItem(item, "organization", "corporate_contributor", "organization_id")
			metadata.Contributors = append(metadata.Contributors, contributor)
		}
	}
	if (eprint.Editors != nil) && (eprint.Editors.Items != nil) {
		for _, item := range eprint.Editors.Items {
			editor := creatorFromItem(item, "person", "editor", "editor_id")
			metadata.Contributors = append(metadata.Contributors, editor)
		}
	}
	if (eprint.ThesisAdvisor != nil) && (eprint.ThesisAdvisor.Items != nil) {
		for _, item := range eprint.ThesisAdvisor.Items {
			advisor := creatorFromItem(item, "person", "thesis_advisor", "thesis_advisor_id")
			metadata.Contributors = append(metadata.Contributors, advisor)
		}
	}
	if (eprint.ThesisCommittee != nil) && (eprint.ThesisCommittee.Items != nil) {
		for _, item := range eprint.ThesisCommittee.Items {
			committee := creatorFromItem(item, "person", "thesis_committee", "thesis_committee_id")
			metadata.Contributors = append(metadata.Contributors, committee)
		}
	}
	metadata.Title = eprint.Title
	if (eprint.AltTitle != nil) && (eprint.AltTitle.Items != nil) {
		for _, item := range eprint.AltTitle.Items {
			title := new(TitleDetail)
			title.Title = item.Value
			metadata.AdditionalTitles = append(metadata.AdditionalTitles, title)
		}
	}
	metadata.Description = eprint.Abstract
	metadata.PublicationDate = eprint.PubDate()

	// Rights are scattered in several EPrints fields, they need to
	// be evaluated to create a "Rights" object used in DataCite/Invenio
	addRights := false
	rights := new(Right)
	if eprint.Rights != "" {
		addRights = true
		rights.Description = eprint.Rights
	}
	// Figure out if our copyright information is in the Note field.
	if (eprint.Note != "") && (strings.Contains(eprint.Note, "©") || strings.Contains(eprint.Note, "copyright") || strings.Contains(eprint.Note, "(c)")) {
		addRights = true
		rights.Description = eprint.Note
	}
	if addRights {
		metadata.Rights = append(metadata.Rights, rights)
	}
	if eprint.CopyrightStatement != "" {
		rights := new(Right)
		rights.Description = eprint.CopyrightStatement
		metadata.Rights = append(metadata.Rights, rights)
	}
	// FIXME: work with Tom to sort out how "Rights" and document level
	// copyright info should work.

	if (eprint.Subjects != nil) && (eprint.Subjects.Items != nil) {
		for _, item := range eprint.Subjects.Items {
			subject := new(Subject)
			subject.Subject = item.Value
			metadata.Subjects = append(metadata.Subjects, subject)
		}
	}

	// FIXME: Work with Tom to figure out correct mapping of rights from EPrints XML
	// FIXME: Language appears to be at the "document" level, not record level

	// Dates are scattered through the primary eprint table.
	if (eprint.DateType != "published") && (eprint.Date != "") {
		metadata.Dates = append(metadata.Dates, dateTypeFromTimestamp("pub_date", eprint.Date, "Publication Date"))
	}
	if eprint.DateStamp != "" {
		metadata.Dates = append(metadata.Dates, dateTypeFromTimestamp("created", eprint.DateStamp, "Created from EPrint's datestamp field"))
	}
	if eprint.LastModified != "" {
		metadata.Dates = append(metadata.Dates, dateTypeFromTimestamp("updated", eprint.LastModified, "Created from EPrint's last_modified field"))
	}
	// FIXME: is this date reflect when it changes status or when it was made available?
	if eprint.StatusChanged != "" {
		metadata.Dates = append(metadata.Dates, dateTypeFromTimestamp("status_changed", eprint.StatusChanged, "Created from EPrint's status_changed field"))
	}
	if eprint.RevNumber != 0 {
		metadata.Version = fmt.Sprintf("v0.0.%d", eprint.RevNumber)
	}
	if eprint.Publisher != "" {
		metadata.Publisher = eprint.Publisher
	}
	if eprint.DOI != "" {
		metadata.Identifiers = append(metadata.Identifiers, mkSimpleIdentifier("DOI", eprint.DOI))
	}
	if eprint.ISBN != "" {
		metadata.Identifiers = append(metadata.Identifiers, mkSimpleIdentifier("ISBN", eprint.ISBN))
	}
	if eprint.ISSN != "" {
		metadata.Identifiers = append(metadata.Identifiers, mkSimpleIdentifier("ISSN", eprint.ISSN))
	}
	if eprint.PMCID != "" {
		metadata.Identifiers = append(metadata.Identifiers, mkSimpleIdentifier("PMCID", eprint.PMCID))
	}
	if (eprint.Funders != nil) && (eprint.Funders.Items != nil) {
		for _, item := range eprint.Funders.Items {
			metadata.Funding = append(metadata.Funding, funderFromItem(item))
		}
	}
	rec.Metadata = metadata
	return nil
}

// filesFromEPrint extracts all the file specific metadata from the
// EPrint record
func (rec *Record) filesFromEPrint(eprint *EPrint) error {
	// crosswalk Files from EPrints DocumentList
	if (eprint.Documents != nil) && (eprint.Documents.Length() > 0) {
		rec.Files = new(Files)
		rec.Files.Enabled = true
		rec.Files.Entries = map[string]*Entry{}
		for i := 0; i < eprint.Documents.Length(); i++ {
			doc := eprint.Documents.IndexOf(i)
			if len(doc.Files) > 0 {
				for _, docFile := range doc.Files {
					entry := new(Entry)
					entry.FileID = docFile.URL
					entry.Size = docFile.FileSize
					entry.MimeType = docFile.MimeType
					if docFile.Hash != "" {
						entry.CheckSum = fmt.Sprintf("%s:%s", strings.ToLower(docFile.HashType), docFile.Hash)
					}
					rec.Files.Entries[docFile.Filename] = entry
					if strings.HasPrefix(docFile.Filename, "preview") {
						rec.Files.DefaultPreview = docFile.Filename
					}
				}
			}
		}
	}
	return nil
}

// tombstoneFromEPrint builds a tombstone is the EPrint record
// eprint_status is deletion.
func (rec *Record) tombstoneFromEPrint(eprint *EPrint) error {
	// FIXME: crosswalk Tombstone
	if eprint.EPrintStatus == "deletion" {
		tombstone := new(Tombstone)
		tombstone.RemovedBy = new(User)
		tombstone.RemovedBy.DisplayName = eprint.Reviewer
		tombstone.RemovedBy.User = eprint.UserID
		if eprint.Suggestions != "" {
			tombstone.Reason = eprint.Suggestions
		}
		rec.Tombstone = tombstone
	}
	return nil
}

// createdUpdatedFromEPrint extracts
func (rec *Record) createdUpdatedFromEPrint(eprint *EPrint) error {
	var (
		created, updated time.Time
		err              error
		tmFmt            string
	)
	// crosswalk Created date
	if len(eprint.DateStamp) > 0 {
		tmFmt = timestamp
		if len(eprint.DateStamp) < 11 {
			tmFmt = datestamp
		}
		created, err = time.Parse(tmFmt, eprint.DateStamp)
		if err != nil {
			return fmt.Errorf("Error parsing datestamp, %s", err)
		}
		rec.Created = created
	}
	if len(eprint.LastModified) > 0 {
		tmFmt = timestamp
		if len(eprint.LastModified) == 10 {
			tmFmt = datestamp
		}
		updated, err = time.Parse(tmFmt, eprint.LastModified)
		if err != nil {
			return fmt.Errorf("Error parsing last modified date, %s", err)
		}
		rec.Updated = updated
	}
	return nil
}

func (rec *Record) ToString() []byte {
	src, _ := json.MarshalIndent(rec, "", "    ")
	return src
}
