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
 * simple.go implements an Invenio 3 like JSON representation
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
)

func CrosswalkEPrintToRecord(eprint *EPrint, rec *Record) error {
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
	/*
	if err := pidFromEPrint(eprint, rec); err != nil {
		return err
	}
	*/
	// Now finish simple record normalization ...
	if err := mapResourceType(rec); err != nil {
		return err
	}
	if err := simplifyCreators(rec); err != nil {
		return err
	}
	if err := simplifyContributors(rec); err != nil {
		return err
	}
	// FIXME: Map eprint record types to invenio RDM record types we've
	// decided on.
	// FIXME: Funders must have a title, could just copy in the funder
	// name for now.
	if err := simplifyFunding(rec); err != nil {
		return err
	}
	return nil
}

// simplifyCreators make sure the identifiers are mapped to Invenio-RDM
// identifiers.
func simplifyCreators(rec *Record) error {
	if rec.Metadata.Creators != nil && len(rec.Metadata.Creators) > 0 {
		creators := []*Creator{}
		for _, creator := range rec.Metadata.Creators {
			//FIXME: do I need to handle assigning a role here?
			appendCreator := false
			if creator.PersonOrOrg != nil && creator.PersonOrOrg.FamilyName != "" {
				if creator.PersonOrOrg.Identifiers != nil && len(creator.PersonOrOrg.Identifiers) > 0 {
					for _, identifier := range creator.PersonOrOrg.Identifiers {
						if identifier.Scheme == "creator_id" {
							identifier.Scheme = "clpid"
						}
					}
					appendCreator = true
				}
			}
			if appendCreator {
				creators = append(creators, creator)
			}
		}
		if len(creators) > 0 {
			rec.Metadata.Creators = creators
		}
	}
	return nil
}

// simplifyContributors make sure the identifiers are mapped to Invenio-RDM
// identifiers.
func simplifyContributors(rec *Record) error {
	if rec.Metadata.Contributors != nil && len(rec.Metadata.Contributors) > 0 {
		contributors := []*Creator{}
		for _, contributor := range rec.Metadata.Contributors {
			appendContributor := false
			if contributor.PersonOrOrg != nil && contributor.PersonOrOrg.FamilyName != "" {
				appendContributor = true
				if contributor.PersonOrOrg.Identifiers != nil && len(contributor.PersonOrOrg.Identifiers) > 0 {
					for _, identifier := range contributor.PersonOrOrg.Identifiers {
						switch identifier.Scheme {
							// conf_creators_id
							// corp_creators_id
							// creators_id 
							// exhibitors_id 
							// lyricists_id 
							// producers_id 
							case "thesis_advisor_id":
								identifier.Scheme = "clpid"
							case "thesis_committee_id":
								identifier.Scheme = "clpid"
							case "author_id": 
								identifier.Scheme = "clpid"
							case "editor_id": 
								identifier.Scheme = "clpid"
							case "contributor_id": 
								identifier.Scheme = "clpid"
						}
					}
				}
			}
			if appendContributor {
				if contributor.Role == nil {
					return fmt.Errorf("contributor is missing role -> %+v", contributor)
				}
				contributors = append(contributors, contributor)
			}
		}
		if len(contributors) > 0 {
			rec.Metadata.Contributors = contributors
		}
	}
	return nil
}


func simplifyFunding(rec *Record) error {
	if rec.Metadata.Funding != nil && len(rec.Metadata.Funding) > 0 {
		for _, funder := range rec.Metadata.Funding {
			if funder.Funder != nil {
				funder.Funder.Scheme = ""
			}
			if funder.Award != nil {
				if funder.Award.Number == "" {
					funder.Award = nil
				} else {
					//NOTE: funder.Award.Title is a struct in
					// Invenio-RDM like
					// ```
					//   title : { "lang": "en", "unavailable" }
					// ```
					// This needs to be normalized in the final
					// Python processing for importing into Invenio-RDM.
					funder.Award.Scheme = ""
				}
			}
		}
	}
	return nil
}

// mapResourceType maps the EPrints record types to a predetermined
// Invenio-RDM record type.
func mapResourceType(rec *Record) error {
	// FIXME: need to load this from a configuration file
	crosswalkResourceTypes := map[string]string{
		"article": "publication-article",
	}
	if rec.Metadata.ResourceType != nil {
		// Should always have an id for a reource_type
		if id, ok := rec.Metadata.ResourceType["id"]; ok {
			if val, hasID := crosswalkResourceTypes[id]; hasID {
				rec.Metadata.ResourceType["id"] = val
				//} else { // FIXME: I don't want to implement a full mapping yet.
				//	return fmt.Errorf("failed to find id %q in record type crosswalk", id)
			}
		}
	}
	return nil
}

/*
// PIDFromEPrint crosswalks the PID from an EPrint record.
func pidFromEPrint(eprint *EPrint, rec *Record) error {
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
*/

// parentFromEPrint crosswalks the Perent unique ID from EPrint record.
func parentFromEPrint(eprint *EPrint, rec *Record) error {
	if eprint.Reviewer != "" {
		parent := new(RecordIdentifier)
		parent.ID = fmt.Sprintf("%s:%d", eprint.Collection, eprint.EPrintID)
		parent.Access = new(Access)
		ownedBy := new(User)
		ownedBy.User = eprint.UserID
		ownedBy.DisplayName = eprint.Reviewer
		parent.Access.OwnedBy = append(parent.Access.OwnedBy, ownedBy)
		rec.Parent = parent
	} else {
		rec.Parent = nil
	}
	return nil
}

// externalPIDFromEPrint aggregates all the external identifiers
// from the EPrint record into Record
func externalPIDFromEPrint(eprint *EPrint, rec *Record) error {
	rec.ExternalPIDs = map[string]*PersistentIdentifier{}
	// Pickup DOI
	if eprint.DOI != "" {
		pid := new(PersistentIdentifier)
		pid.Identifier = eprint.DOI
		pid.Provider = "datacite"
		pid.Client = ""
		rec.ExternalPIDs["doi"] = pid
	}
	// Pickup ISSN
	if eprint.ISBN != "" {
		pid := new(PersistentIdentifier)
		pid.Identifier = eprint.ISSN
		pid.Provider = ""
		pid.Client = ""
		rec.ExternalPIDs["issn"] = pid
	}
	// Pickup ISBN
	if eprint.ISBN != "" {
		pid := new(PersistentIdentifier)
		pid.Identifier = eprint.ISBN
		pid.Provider = ""
		pid.Client = ""
		rec.ExternalPIDs["isbn"] = pid
	}
	//FIXME: figure out if we have other persistent identifiers
	//scattered in the EPrints XML and map them.
	return nil
}

// recordAccessFromEPrint extracts access permissions from the EPrint
func recordAccessFromEPrint(eprint *EPrint, rec *Record) error {
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
	// By default lets assume the files are restricted.
	rec.RecordAccess.Files = "resticted"
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
				embargo := new(Embargo)
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

func uriToContributorType(role_uri string) string {
    roles := map[string]string{
        // Article Author
        "http://coda.library.caltech.edu/ARA": "author_section",
        // Astronaut
        "http://coda.library.caltech.edu/AST": "astronaut",
        // Author of afterword, colophon, etc.
        "http://www.loc.gov/loc.terms/relators/AFT": "aft",
        // Bibliographic antecedent
        "http://www.loc.gov/loc.terms/relators/ANT": "ant",
        // Author in quotations or text abstracts
        "http://www.loc.gov/loc.terms/relators/AQT": "aqt",
        // Screenwriter
        "http://www.loc.gov/loc.terms/relators/AUS": "screenwriter",
        // Author, joint author
        "http://www.loc.gov/loc.terms/relators/AUT": "author",
        // Collaborator
        "http://www.loc.gov/loc.terms/relators/CLB": "collaborator",
        // Compiler
        "http://www.loc.gov/loc.terms/relators/COM": "compiler",
        // Contributor
        "http://www.loc.gov/loc.terms/relators/CTB": "contributor",
        // Directory
        "http://www.loc.gov/loc.terms/relators/DRT": "director",
        // Editor
        "http://www.loc.gov/loc.terms/relators/EDT": "editor",
        // Narrator
        "http://www.loc.gov/loc.terms/relators/NRT": "narrator",
        // Other
        "http://www.loc.gov/loc.terms/relators/OTH": "other",
        // Publishing director
        "http://www.loc.gov/loc.terms/relators/PBD": "publishing_director",
        // Programmer
        "http://www.loc.gov/loc.terms/relators/PRG": "programmer",
        // Reviewer
        "http://www.loc.gov/loc.terms/relators/REV": "reviewer",
        // Research team member
        "http://www.loc.gov/loc.terms/relators/RTM": "research_team",
        // Speaker
        "http://www.loc.gov/loc.terms/relators/SPK": "speaker",
        // Teacher
        "http://www.loc.gov/loc.terms/relators/TCH": "teacher",
        // Translator
        "http://www.loc.gov/loc.terms/relators/TRL": "translator",
    }
	if val, ok := roles[role_uri]; ok {
		return val
	}
	return "contributor"
}

func creatorFromItem(item *Item, objType string, objRoleSrc string, objIdType string) *Creator {
	person := new(PersonOrOrg)
	person.Type = objType
	if item.Name != nil {
		person.FamilyName = item.Name.Family
		person.GivenName = item.Name.Given
	}
	if item.ORCID != "" {
		identifier := new(Identifier)
		identifier.Scheme = "orcid"
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
	switch objRoleSrc{
		case "contributor":
			//NOTE: for contributors we need to map the type as LOC URI
			// to a person's role.
			contributorType := "contributor"
			if item != nil {
				contributorType = uriToContributorType(item.Type)
			}
			creator.Role = &Role{ ID: contributorType }
		case "editor":
			creator.Role = &Role{ ID: "editor" }
	}
	// FIXME: For Creators we skip adding the role and affiliation for now,
	// it break RDM.
	//creator.PersonOrOrg.Role = &Role{ ID: objRoleSrc }
	//FIXME: Need to map affiliation here when we're ready.
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
	identifier.Scheme = strings.ToLower(scheme)
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
		funder.Funder = org
	}
	return funder
}

// metadataFromEPrint extracts metadata from the EPrint record
func metadataFromEPrint(eprint *EPrint, rec *Record) error {
	metadata := new(Metadata)
	metadata.ResourceType = map[string]string{}
	metadata.ResourceType["id"] = eprint.Type
	// NOTE: Creators get listed in the citation, Contributors do not.
	if (eprint.Creators != nil) && (eprint.Creators.Items != nil) {
		for _, item := range eprint.Creators.Items {
			creator := creatorFromItem(item, "personal", "creator", "creator_id")
			metadata.Creators = append(metadata.Creators, creator)
		}
	}
	if (eprint.CorpCreators != nil) && (eprint.CorpCreators.Items != nil) {
		for _, item := range eprint.CorpCreators.Items {
			creator := creatorFromItem(item, "organizational", "corporate_creator", "organization_id")
			metadata.Creators = append(metadata.Creators, creator)
		}
	}
	if (eprint.Contributors != nil) && (eprint.Contributors.Items != nil) {
		for _, item := range eprint.Contributors.Items {
			contributor := creatorFromItem(item, "personal", "contributor", "contributor_id")

			metadata.Contributors = append(metadata.Contributors, contributor)
		}
	}
	if (eprint.CorpContributors != nil) && (eprint.CorpContributors.Items != nil) {
		for _, item := range eprint.CorpContributors.Items {
			contributor := creatorFromItem(item, "organizational", "corporate_contributor", "organization_id")
			metadata.Contributors = append(metadata.Contributors, contributor)
		}
	}
	if (eprint.Editors != nil) && (eprint.Editors.Items != nil) {
		for _, item := range eprint.Editors.Items {
			editor := creatorFromItem(item, "personal", "editor", "editor_id")
			metadata.Contributors = append(metadata.Contributors, editor)
		}
	}
	if (eprint.ThesisAdvisor != nil) && (eprint.ThesisAdvisor.Items != nil) {
		for _, item := range eprint.ThesisAdvisor.Items {
			advisor := creatorFromItem(item, "personal", "thesis_advisor", "thesis_advisor_id")
			metadata.Contributors = append(metadata.Contributors, advisor)
		}
	}
	if (eprint.ThesisCommittee != nil) && (eprint.ThesisCommittee.Items != nil) {
		for _, item := range eprint.ThesisCommittee.Items {
			committee := creatorFromItem(item, "personal", "thesis_committee", "thesis_committee_id")
			metadata.Contributors = append(metadata.Contributors, committee)
		}
	}
	metadata.Title = eprint.Title
	if (eprint.AltTitle != nil) && (eprint.AltTitle.Items != nil) {
		for _, item := range eprint.AltTitle.Items {
			if strings.TrimSpace(item.Value) != "" {
				title := new(TitleDetail)
				title.Title = item.Value
				metadata.AdditionalTitles = append(metadata.AdditionalTitles, title)
			}
		}
	}
	if eprint.Abstract != "" {
		metadata.Description = eprint.Abstract
	}
	metadata.PublicationDate = eprint.PubDate()

	// Rights are scattered in several EPrints fields, they need to
	// be evaluated to create a "Rights" object used in DataCite/Invenio
	addRights := false
	rights := new(Right)
	if eprint.Rights != "" {
		addRights = true
		rights.Description = &Description{
			Description: eprint.Rights,
		}
	}
	// Figure out if our copyright information is in the Note field.
	if (eprint.Note != "") && (strings.Contains(eprint.Note, "©") || strings.Contains(eprint.Note, "copyright") || strings.Contains(eprint.Note, "(c)")) {
		addRights = true
		rights.Description = &Description{
			Description: fmt.Sprintf("%s", eprint.Note),
		}
	}
	if addRights {
		metadata.Rights = append(metadata.Rights, rights)
	}
	if eprint.CopyrightStatement != "" {
		rights := new(Right)
		rights.Description = &Description{
			Description: eprint.CopyrightStatement,
		}
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
	if eprint.Datestamp != "" {
		metadata.Dates = append(metadata.Dates, dateTypeFromTimestamp("created", eprint.Datestamp, "Created from EPrint's datestamp field"))
	}
	if eprint.LastModified != "" {
		metadata.Dates = append(metadata.Dates, dateTypeFromTimestamp("updated", eprint.LastModified, "Created from EPrint's last_modified field"))
	}
	/*
		// status_changed is not a date type in Invenio-RDM, might be mapped
		// into available object.
		// FIXME: is this date reflect when it changes status or when it was made available?
		if eprint.StatusChanged != "" {
			metadata.Dates = append(metadata.Dates, dateTypeFromTimestamp("status_changed", eprint.StatusChanged, "Created from EPrint's status_changed field"))
		}
	*/
	if eprint.RevNumber != 0 {
		metadata.Version = fmt.Sprintf("v%d", eprint.RevNumber)
	}
	if eprint.Publisher != "" {
		metadata.Publisher = eprint.Publisher
	} else if eprint.Publication != "" {
		metadata.Publisher = eprint.Publication
	} else if eprint.DOI == "" {
		metadata.Publisher = "Caltech Library"
	}
	
	// Pickup EPrint ID as "external identifier" in .metadata.identifier
	if eprint.EPrintID > 0 {
		metadata.Identifiers = append(metadata.Identifiers, mkSimpleIdentifier("eprintid", fmt.Sprintf("%d", eprint.EPrintID)))
	}

	if eprint.DOI != "" {
		metadata.Identifiers = append(metadata.Identifiers, mkSimpleIdentifier("doi", eprint.DOI))
	}
	if eprint.ISBN != "" {
		metadata.Identifiers = append(metadata.Identifiers, mkSimpleIdentifier("isbn", eprint.ISBN))
	}
	if eprint.ISSN != "" {
		metadata.Identifiers = append(metadata.Identifiers, mkSimpleIdentifier("issn", eprint.ISSN))
	}
	if eprint.PMCID != "" {
		metadata.Identifiers = append(metadata.Identifiers, mkSimpleIdentifier("pmcid", eprint.PMCID))
	}
	if (eprint.Funders != nil) && (eprint.Funders.Items != nil) {
		for _, item := range eprint.Funders.Items {
			if item.Agency != "" {
				metadata.Funding = append(metadata.Funding, funderFromItem(item))
			}
		}
	}
	rec.Metadata = metadata
	return nil
}

// filesFromEPrint extracts all the file specific metadata from the
// EPrint record
func filesFromEPrint(eprint *EPrint, rec *Record) error {
	// crosswalk Files from EPrints DocumentList
	if (eprint != nil) && (eprint.Documents != nil) && (eprint.Documents.Length() > 0) {
		addFiles := false
		files := new(Files)
		files.Order = []string{}
		files.Enabled = true
		files.Entries = map[string]*Entry{}
		for i := 0; i < eprint.Documents.Length(); i++ {
			doc := eprint.Documents.IndexOf(i)
			if len(doc.Files) > 0 {
				for _, docFile := range doc.Files {
					addFiles = true
					entry := new(Entry)
					entry.FileID = docFile.URL
					entry.Size = docFile.FileSize
					entry.MimeType = docFile.MimeType
					if docFile.Hash != "" {
						entry.CheckSum = fmt.Sprintf("%s:%s", strings.ToLower(docFile.HashType), docFile.Hash)
					}
					files.Entries[docFile.Filename] = entry
					if strings.HasPrefix(docFile.Filename, "preview") {
						files.DefaultPreview = docFile.Filename
					}
				}
			}
		}
		if addFiles {
			rec.Files = files
		} else {
			rec.Files = nil
		}
	}
	return nil
}

// tombstoneFromEPrint builds a tombstone is the EPrint record
// eprint_status is deletion.
func tombstoneFromEPrint(eprint *EPrint, rec *Record) error {
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
func createdUpdatedFromEPrint(eprint *EPrint, rec *Record) error {
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
