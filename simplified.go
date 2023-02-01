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
	"fmt"
	"strings"
	"time"

	// Caltech Library packages
	"github.com/caltechlibrary/simplified"
)

func CrosswalkEPrintToRecord(eprint *EPrint, rec *simplified.Record) error {
	rec.Schema = `local://records/record-v2.0.0.json`
	rec.ID = fmt.Sprintf("%s:%d", eprint.Collection, eprint.EPrintID)

	if err := parentFromEPrint(eprint, rec); err != nil {
		return err
	}
	if err := externalPIDFromEPrint(eprint, rec); err != nil {
		return err
	}
	if err := recordAccessFromEPrint(eprint, rec); err != nil {
		return err
	}

	if err := metadataFromEPrint(eprint, rec); err != nil {
		return err
	}
	if err := filesFromEPrint(eprint, rec); err != nil {
		return err
	}

	if eprint.EPrintStatus == "deletion" {
		if err := tombstoneFromEPrint(eprint, rec); err != nil {
			return err
		}
	}

	if err := createdUpdatedFromEPrint(eprint, rec); err != nil {
		return err
	}
	if err := pidFromEPrint(eprint, rec); err != nil {
		return err
	}
	return nil
}

// PIDFromEPrint crosswalks the PID from an EPrint record.
func pidFromEPrint(eprint *EPrint, rec *simplified.Record) error {
	data := map[string]interface{}{}
	src := fmt.Sprintf(`{
"id": %d,
"pid": { "eprint": "eprintid"}
}`, eprint.EPrintID)
	err := jsonDecode([]byte(src), &data)
	if err != nil {
		return fmt.Errorf("Cannot generate PID from EPrint %d", eprint.EPrintID)
	}
	rec.PID = data
	return nil
}

// parentFromEPrint crosswalks the Perent unique ID from EPrint record.
func parentFromEPrint(eprint *EPrint, rec *simplified.Record) error {
	parent := new(simplified.RecordIdentifier)
	parent.ID = fmt.Sprintf("%s:%d", eprint.Collection, eprint.EPrintID)
	parent.Access = new(simplified.Access)
	ownedBy := new(simplified.User)
	ownedBy.User = eprint.UserID
	ownedBy.DisplayName = eprint.Reviewer
	parent.Access.OwnedBy = append(parent.Access.OwnedBy, ownedBy)
	rec.Parent = parent
	return nil
}

// externalPIDFromEPrint aggregates all the external identifiers
// from the EPrint record into Record
func externalPIDFromEPrint(eprint *EPrint, rec *simplified.Record) error {
	rec.ExternalPIDs = map[string]*simplified.PersistentIdentifier{}
	// Pickup DOI
	if eprint.DOI != "" {
		pid := new(simplified.PersistentIdentifier)
		pid.Identifier = eprint.DOI
		pid.Provider = "datacite" // FIXME: should be DataCite or CrossRef
		pid.Client = ""           // FIXME: need to find out client string
		rec.ExternalPIDs["doi"] = pid
	}
	// Pickup ISSN
	if eprint.ISBN != "" {
		pid := new(simplified.PersistentIdentifier)
		pid.Identifier = eprint.ISSN
		pid.Provider = "" // FIXME: Need to find out identifier string
		pid.Client = ""   // FIXME: need to find out client string
		rec.ExternalPIDs["ISSN"] = pid
	}
	// Pickup ISBN
	if eprint.ISBN != "" {
		pid := new(simplified.PersistentIdentifier)
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
func recordAccessFromEPrint(eprint *EPrint, rec *simplified.Record) error {
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
	rec.RecordAccess = new(simplified.RecordAccess)
	if isPublic {
		rec.RecordAccess.Record = "public"
	} else {
		rec.RecordAccess.Record = "restricted"
	}
	// Need to make sure record is not embargoed
	if eprint.Documents != nil {
		for i := 0; i < eprint.Documents.Length(); i++ {
			doc := eprint.Documents.IndexOf(i)
			if doc.DateEmbargo != "" {
				embargo := new(simplified.Embargo)
				embargo.Until = doc.DateEmbargo
				if eprint.Suggestions != "" {
					embargo.Reason = eprint.Suggestions
				}
				if doc.Security == "internal" {
					embargo.Active = true
				} else {
					embargo.Active = false
				}
				rec.RecordAccess.Embargo = embargo
				break
			}
		}
	}
	return nil
}

func creatorFromItem(item *Item, objType string, objRole string, objIdType string) *simplified.Creator {
	person := new(simplified.PersonOrOrg)
	person.Type = objType
	if item.Name != nil {
		person.FamilyName = item.Name.Family
		person.GivenName = item.Name.Given
	}
	if item.ORCID != "" {
		identifier := new(simplified.Identifier)
		identifier.Scheme = "ORCID"
		identifier.Identifier = item.ORCID
		person.Identifiers = append(person.Identifiers, identifier)
	}
	if item.ID != "" {
		identifier := new(simplified.Identifier)
		identifier.Scheme = objIdType
		identifier.Identifier = item.ID
		person.Identifiers = append(person.Identifiers, identifier)
	}
	creator := new(simplified.Creator)
	creator.PersonOrOrg = person
	creator.Role = objRole

	return creator
}

func dateTypeFromTimestamp(dtType string, timestamp string, description string) *simplified.DateType {
	dt := new(simplified.DateType)
	dt.Type = new(simplified.Type)
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

func mkSimpleIdentifier(scheme, value string) *simplified.Identifier {
	identifier := new(simplified.Identifier)
	identifier.Scheme = scheme
	identifier.Identifier = value
	return identifier
}

func funderFromItem(item *Item) *simplified.Funder {
	funder := new(simplified.Funder)
	if item.GrantNumber != "" {
		funder.Award = new(simplified.Identifier)
		funder.Award.Number = item.GrantNumber
		funder.Award.Scheme = "eprints_grant_number"
	}
	if item.Agency != "" {
		org := new(simplified.Identifier)
		org.Name = item.Agency
		org.Scheme = "eprints_agency"
		funder.Funder = append(funder.Funder, org)
	}
	return funder
}

// metadataFromEPrint extracts metadata from the EPrint record
func metadataFromEPrint(eprint *EPrint, rec *simplified.Record) error {
	metadata := new(simplified.Metadata)
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
			title := new(simplified.TitleDetail)
			title.Title = item.Value
			metadata.AdditionalTitles = append(metadata.AdditionalTitles, title)
		}
	}
	if eprint.Abstract != "" {
		metadata.Description = eprint.Abstract
	}
	metadata.PublicationDate = eprint.PubDate()

	// Rights are scattered in several EPrints fields, they need to
	// be evaluated to create a "Rights" object used in DataCite/Invenio
	addRights := false
	rights := new(simplified.Right)
	if eprint.Rights != "" {
		addRights = true
		rights.Description = eprint.Rights
	}
	// Figure out if our copyright information is in the Note field.
	if (eprint.Note != "") && (strings.Contains(eprint.Note, "©") || strings.Contains(eprint.Note, "copyright") || strings.Contains(eprint.Note, "(c)")) {
		addRights = true
		rights.Description = fmt.Sprintf("%s", eprint.Note)
	}
	if addRights {
		metadata.Rights = append(metadata.Rights, rights)
	}
	if eprint.CopyrightStatement != "" {
		rights := new(simplified.Right)
		rights.Description = eprint.CopyrightStatement
		metadata.Rights = append(metadata.Rights, rights)
	}
	// FIXME: work with Tom to sort out how "Rights" and document level
	// copyright info should work.

	if (eprint.Subjects != nil) && (eprint.Subjects.Items != nil) {
		for _, item := range eprint.Subjects.Items {
			subject := new(simplified.Subject)
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
	if eprint.Datestamp != "" {
		metadata.Dates = append(metadata.Dates, dateTypeFromTimestamp("created", eprint.Datestamp, "Created from EPrint's datestamp field"))
	}
	if eprint.LastModified != "" {
		metadata.Dates = append(metadata.Dates, dateTypeFromTimestamp("updated", eprint.LastModified, "Created from EPrint's last_modified field"))
	}
	// FIXME: is this date reflect when it changes status or when it was made available?
	if eprint.StatusChanged != "" {
		metadata.Dates = append(metadata.Dates, dateTypeFromTimestamp("status_changed", eprint.StatusChanged, "Created from EPrint's status_changed field"))
	}
	if eprint.RevNumber != 0 {
		metadata.Version = fmt.Sprintf("v%d", eprint.RevNumber)
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
func filesFromEPrint(eprint *EPrint, rec *simplified.Record) error {
	// crosswalk Files from EPrints DocumentList
	if (eprint != nil) && (eprint.Documents != nil) && (eprint.Documents.Length() > 0) {
		rec.Files = new(simplified.Files)
		rec.Files.Enabled = true
		rec.Files.Entries = map[string]*simplified.Entry{}
		for i := 0; i < eprint.Documents.Length(); i++ {
			doc := eprint.Documents.IndexOf(i)
			if len(doc.Files) > 0 {
				for _, docFile := range doc.Files {
					entry := new(simplified.Entry)
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
func tombstoneFromEPrint(eprint *EPrint, rec *simplified.Record) error {
	// FIXME: crosswalk Tombstone
	if eprint.EPrintStatus == "deletion" {
		tombstone := new(simplified.Tombstone)
		tombstone.RemovedBy = new(simplified.User)
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
func createdUpdatedFromEPrint(eprint *EPrint, rec *simplified.Record) error {
	var (
		created, updated time.Time
		err              error
		tmFmt            string
	)
	// crosswalk Created date
	if len(eprint.Datestamp) > 0 {
		tmFmt = timestamp
		if len(eprint.Datestamp) < 11 {
			tmFmt = datestamp
		}
		created, err = time.Parse(tmFmt, eprint.Datestamp)
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
