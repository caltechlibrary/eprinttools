package eprinttools

import (
	"fmt"
)

//
// End point documentation
//
func readmeDocument() string {
	return `
USAGE
=====

    ep3apid [OPTIONS] [SETTINGS_FILENAME]

SYNOPSIS
--------

Run an extended EPrints 3.x web API based on direct manipulation
of EPrint's MySQL database(s).


DETAIL
------

__ep3apid__ can be run from the command line and the will create an http web service. The web service provides a limitted number of end points providing eprint ids for content matched in EPrints's MySQL databases. You can configure it to provide read/write support to and from the MySQL databases used by EPrints.

The following URL end points are intended to take one unique identifier and map that to one or more EPrint IDs. This can be done because each unique ID  targeted can be identified by querying a single table in EPrints.  In addition the scan can return the complete results since all EPrint IDs are integers and returning all EPrint IDs in any of our repositories is sufficiently small to be returned in a single HTTP request.

Configuration information
-------------------------

There are two end points that give you information about what repositories are configured in for __ep3apid__ and what the database structure (tables and column names) for each configure repository.

- '/repositores' - returns a list of repositories configured for access by __ep3apid__
- '/repository/{REPO_ID}' returns the databases and columns of the repository indicated by "{REPO_ID}".


Unique ID to EPrint ID
----------------------

Unique ids maybe standards based (e.g. ORCID, DOI, ISSN, ISBN) or internal (e.g. group ids, funder ids)

- '/{REPO_ID}/doi/{DOI}' with the adoption of EPrints "doi" field in the EPrint table it makes sense to have a quick translation of DOI to EPrint id for a given EPrints repository. 
- '/{REPO_ID}/pmid/{PMID}' with the "pmid" field in the EPrint table, it refers to PubMed is an index of the biomedical literature.
- '/{REPO_ID}/pmcid/{PMCID}' with the "pmcid" field in the EPrint table, PMCID an Identifier to each full-text paper in PubMed Central Archive
- '/{REPO_ID}/creator-id' returns a list of creaator-id available in the eprints repository
- '/{REPO_ID}/creator-id/{CREATOR_ID}' scans the name creator id field associated with creators and returns a list of EPrint ID 
- '/{REPO_ID}/creator-orcid' return a list of "orcid" associated with creators in repository 
- '/{REPO_ID}/creator-orcid/{ORCID}' scans the "orcid" field associated with creators and returns a list of EPrint ID 
- '/{REPO_ID}/editor-id' returns a list of editor ids available in the EPrints repository
- '/{REPO_ID}/editor-id/{CREATOR_ID}' scans the name creator id field associated with editors and returns a list of EPrint ID 
- '/{REPO_ID}/contributor-id' returns a list of contributor ids available in the eprints repository
- '/{REPO_ID}/contributor-id/{CONTRIBUTOR_ID}' scans the "id" field associated with a contributors and returns a list of EPrint ID 
- '/{REPO_ID}/advisor-id' returns a list of advisor ids in the eprints repository
- '/{REPO_ID}/advisor-id/{ADVISOR_ID}' scans the name advisor id field associated with advisors and returns a list of EPrint ID 
- '/{REPO_ID}/committee-id' returns a list of committee id in EPrints repository
- '/{REPO_ID}/committee-id/{COMMITTEE_ID}' scans the committee id field associated with committee members and returns a list of EPrint ID
- '/{REPO_ID}/group-id' returns a list of group ids in EPrints repository
- '/{REPO_ID}/group-id/{GROUP_ID}' this scans group ID and returns a list of EPrint IDs associated with the group
- '/{REPO_ID}/funder-id' returns a list of funders in the EPrints repository
- '/{REPO_ID}/funder-id/{FUNDER_ID}' returns a list of EPrint IDs associated with the funder
- '/{REPO_ID}/grant-number' returns a list of grant numbers in EPrints repository
- '/{REPO_ID}/grant-number/{GRANT_NUMBER}' returns a list of EPrint IDs associated with the grant number
- '/{REPO_ID}/issn' - returns a list of ISSN in repository
- '/{REPO_ID}/issn/{ISSN}' - returns a list eprint id for ISSN in repository
- '/{REPO_ID}/isbn' - returns a list of ISBN in repository
- '/{REPO_ID}/isbn/{ISBN}' - returns a list eprint id for ISBN in repository
- '/{REPO_ID}/patent-number' - return a list of patent numbers in repository
- '/{REPO_ID}/patent-number/{PATENT_NUMBER}' - return a list eprint ids for patent number in repository
- '/{REPO_ID}/patent-applicant' - return a list of patent applicants in repository
- '/{REPO_ID}/patent-applicant/{PATENT_APPLICANT}' - return a list eprint ids for patent applicant in repository
- '/{REPO_ID}/patent-classification' - return a list of patent classificatins in repository
- '/{REPO_ID}/patent-classification/{PATENT_CLASSIFICATION}' - return a list eprint ids for patent classification in repository
- '/{REPO_ID}/patent-assignee' - return a list of patent assignee in repository
- '/{REPO_ID}/patent-assignee/{PATENT_ASSIGNEE}' - return a list eprint ids for patent assignee in repository


Change Events
-------------

The follow API end points would facilitate faster updates to our feeds platform as well as allow us to create a separate public view of our EPrint repository content.

- '/{REPO_ID}/keys' returns complete list of EPrint ID in the repository
- '/{REPO_ID}/updated/{TIMESTAMP}/{TIMESTAMP}' returns a list of EPrint IDs updated starting at the first timestamp (timestamps should have a resolution to the minute, e.g. "YYYY-MM-DD HH:MM:SS") through inclusive of the second timestmap (if the second is omitted the timestamp is assumed to be "now")
- '/{REPO_ID}/deleted/{TIMESTAMP}/{TIMESTAMP}' through the returns a list of EPrint IDs deleted starting at first timestamp through inclusive of the second timestamp, if the second timestamp is omitted it is assumed to be "now"
- '/{REPO_ID}/pubdate/{APROX_DATESTAMP}/{APPOX_DATESTMP}' this query scans the EPrint table for records with publication starts starting with the first approximate date through inclusive of the second approximate date. If the second date is omitted it is assumed to be "today". Approximate dates my be expressed just the year (starting with Jan 1, ending with Dec 31), just the year and month (starting with first day of month ending with the last day) or year, month and day. The end returns zero or more EPrint IDs.

Write API
---------

As of __ep3apid__ version 1.0.3 a new set of end points exists for reading and writing EPrints XML. This can be enabled per repository. It only supports interaction with one record at a time.  You can retrieve full EPrint XML using a GET request. This EPrint XML is generate dynamically based on the contents of the MySQL tables configured in the EPrints instance.  You can write EPrints records with a POST.  If the eprint ID furnished in the POST is 0 (zero) then a new record will be created. Otherwise the contents of the EPrint XML you post will replace the existing eprint record.  This transaction takes place only at the SQL level. None of EPrints's native Perl API is invoked. 

The end point is '/{REPO_ID}/eprint/{EPRINT_ID}' for EPrint XML.

To enable this feature add the attribute '"write": true' to the repositories setting in settins.json.

settings.json (configuration)
-----------------------------

To run the web service create a JSON file named settings.ini in the current directory where you're invoking __ep3apid__ from. The web service can be started with running

` + "```" + `
    ep3apid
` + "```" + `

or to load "settings.json" from the current work directory.

` + "```" + `
    ep3apid settings.json
` + "```" + `

The JSON settings.json file should look something like "REPO_ID" would
be the name used in the __ep3apid__ to access a specific repsitory. The
"dsn" value should be replaced with an appropriate data source name to
access the MySQL database for the repository you're supporting. You can have many repositories configured in a single __ep3apid__ instance.

` + "```" + `
    {
        "repositories": {
            "REPO_ID": {
                "dsn": "DB_USER:SECRET@/DB_NAME",
                "base_url": "URL_TO_EPRINT_REPOSITORY",
                "write": false
            },
            ... /* Additional repositories configured here */ ...
        }
    }
` + "```" + `

Options
-------

` + "```" + `
  -h	Display this help message
  -help
    	Display this help message
  -license
    	Display software license
  -version
    	Display software version
` + "```" + `


`
}

func repositoriesDocument() string {
	return `
Repositories (end point)
========================

This end point lists the repositories known to the __ep3apid__ service.

- '/repositories' returns a JSON array of repository names defined in settings.json
- '/repositories/help' returns this documentation.

Example
-------

In this example we assume the __ep3apid__ services is running on "localhost:8484" and is configured to support two repositories "lemurprints" and "test3". We are using curl to retrieve the data.

` + "```" + `shell
    curl -X GET http://localhost:8484/repositories
` + "```" + `

This should return a JSON array like

` + "```" + `json
    [
        "lemurprints",
        "test3"
    ]
` + "```" + `

`
}

func repositoryDocument() string {
	return `
Repository (end point)
----------------------

The end point executes a sequence of "SHOW" SQL statements to build a JSON object with table names as attributes pointing at an array of column names. This is suitable to determine on a per repository bases the related table and columnames representing an EPrint record.

- '/repository' return this documentation
- '/repository/{REPO_ID}' return the JSON representation
- '/repository/{REPO_ID/help' return this documentation

Example
-------


` + "```" + `shell
   curl http://localhost:8485/lemurprints/tables
` + "```" + `

Would return a JSON expression similar to the expression below.  The object has attributes that map to the EPrint talbles and for each table the attribute points at an array of column names.

The "eprint" table is the root of the object. Each other table is a sub object or array. Tables containing a "pos" field render as an array of objects (e.g. the "item" elements in the EPrint XML). If pos is missing then it is an object with attributes and values.

Each table is relatated by the "eprintid" column ("..." in the object below means the text was abbeviated). Sub tables are related by eprintid and pos columns. All metadata table names begin with "eprint%" or "document%".


` + "```" + `json
{
  "eprint": [ "abstract", "alt_url", "book_title", "collection", "commentary", "completion_time", "composition_type", "contact_email", "coverage_dates", "data_type", "date_day", "date_month", "date_type", "date_year", "datestamp_day", "datestamp_hour", "datestamp_minute", ... ],
  "eprint_accompaniment": [ "accompaniment", "eprintid", "pos" ],
  "eprint_alt_title": [ "alt_title", "eprintid", "pos" ],
  "eprint_conductors_id": [ "conductors_id", "eprintid", "pos" ],
  "eprint_conductors_name": [ "conductors_name_family", "conductors_name_given", "conductors_name_honourific", "conductors_name_lineage", "eprintid", "pos" ],
  "eprint_conductors_uri": [ "conductors_uri", "eprintid", "pos" ],
  "eprint_conf_creators_id": [ "conf_creators_id", "eprintid", "pos" ],
  "eprint_conf_creators_name": [ "conf_creators_name", "eprintid", "pos" ],
  "eprint_conf_creators_uri": [ "conf_creators_uri", "eprintid", "pos" ],
  "eprint_contributors_id": [ "contributors_id", "eprintid", "pos" ],
  "eprint_contributors_name": [ "contributors_name_family", "contributors_name_given", "contributors_name_honourific", "contributors_name_lineage", "eprintid", "pos" ],
  "eprint_contributors_type": [ "contributors_type", "eprintid", "pos" ], "eprint_contributors_uri": [ "contributors_uri", "eprintid", "pos" ],
  "eprint_copyright_holders": [ "copyright_holders", "eprintid", "pos" ],
  "eprint_corp_creators_id": [ "corp_creators_id", "eprintid", "pos" ],
  "eprint_corp_creators_name": [ "corp_creators_name", "eprintid", "pos" ],
  "eprint_corp_creators_uri": [ "corp_creators_uri", "eprintid", "pos" ],
  "eprint_creators_id": [ "creators_id", "eprintid", "pos" ],
  "eprint_creators_name": [ "creators_name_family", "creators_name_given", "creators_name_honourific", "creators_name_lineage", "eprintid", "pos" ],
  "eprint_creators_uri": [ "creators_uri", "eprintid", "pos" ],
  "eprint_divisions": [ "divisions", "eprintid", "pos" ],
  "eprint_editors_id": [ "editors_id", "eprintid", "pos" ],
  "eprint_editors_name": [ "editors_name_family", "editors_name_given", "editors_name_honourific", "editors_name_lineage", "eprintid", "pos" ],
  "eprint_editors_uri": [ "editors_uri", "eprintid", "pos" ],
  "eprint_exhibitors_id": [ "eprintid", "exhibitors_id", "pos" ],
  "eprint_exhibitors_name": [ "eprintid", "exhibitors_name_family", "exhibitors_name_given", "exhibitors_name_honourific", "exhibitors_name_lineage", "pos" ],
  "eprint_exhibitors_uri": [ "eprintid", "exhibitors_uri", "pos" ],
  "eprint_funders_agency": [ "eprintid", "funders_agency", "pos" ],
  "eprint_funders_grant_number": [ "eprintid", "funders_grant_number", "pos" ],
  ...
}
` + "```" + `

`
}

func keysDocument(repoID string) string {
	return fmt.Sprintf(`'/%s/keys' returns a list of EPrint ID in the repository`, repoID)
}

func createdDocument(repoID string) string {
	return fmt.Sprintf(`'/%s/created/{TIMESTAMP}/{TIMESTAMP}' returns a list of EPrint IDs created starting at the first timestamp (timestamps should have a resolution to the minute, e.g. "YYYY-MM-DD HH:MM:SS") through inclusive of the second timestmap (if the second is omitted the timestamp is assumed to be "now")`, repoID)
}

func updatedDocument(repoID string) string {
	return fmt.Sprintf(`'/%s/updated/{TIMESTAMP}/{TIMESTAMP}' returns a list of EPrint IDs updated starting at the first timestamp (timestamps should have a resolution to the minute, e.g. "YYYY-MM-DD HH:MM:SS") through inclusive of the second timestmap (if the second is omitted the timestamp is assumed to be "now")`, repoID)
}

func deletedDocument(repoID string) string {
	return fmt.Sprintf(`'/%s/deleted/{TIMESTAMP}/{TIMESTAMP}' through the returns a list of EPrint IDs deleted starting at first timestamp through inclusive of the second timestamp, if the second timestamp is omitted it is assumed to be "now"`, repoID)
}

func pubdateDocument(repoID string) string {
	return fmt.Sprintf(`'/%s/pubdate/{APROX_DATESTAMP}/{APPOX_DATESTAMP}' this query scans the EPrint table for records with publication starts starting with the first approximate date through inclusive of the second approximate date. If the second date is omitted it is assumed to be "today". Approximate dates my be expressed just the year (starting with Jan 1, ending with Dec 31), just the year and month (starting with first day of month ending with the last day) or year, month and day. The end returns zero or more EPrint IDs.`, repoID)
}

func doiDocument(repoID string) string {
	return fmt.Sprintf(`'/%s/doi/{DOI}' with the adoption of EPrints "doi" field in the EPrint table it makes sense to have a quick translation of DOI to EPrint id for a given EPrints repository.`, repoID)
}

func creatorIDDocument(repoID string) string {
	return fmt.Sprintf(`
- '/%s/creator-id' returns a list of creaator-id available in the eprints repository
- '/%s/creator-id/{CREATOR_ID}' scans the name creator id field associated with creators and returns a list of EPrint ID 
`, repoID, repoID)
}

func creatorORCIDDocument(repoID string) string {
	return fmt.Sprintf(`
- '/%s/creator-orcid' return a list of "orcid" associated with creators in repository 
- '/%s/creator-orcid/{ORCID}' scans the "orcid" field associated with creators and returns a list of EPrint ID
`, repoID, repoID)
}

func editorIDDocument(repoID string) string {
	return fmt.Sprintf(`
- '/%s/editor-id' returns a list of editor ids available in the EPrints repository
- '/%s/editor-id/{CREATOR_ID}' scans the name creator id field associated with editors and returns a list of EPrint ID 
`, repoID, repoID)
}

func contributorIDDocument(repoID string) string {
	return fmt.Sprintf(`
- '/%s/contributor-id' returns a list of contributor ids available in the eprints repository
- '/%s/contributor-id/{CONTRIBUTOR_ID}' scans the "id" field associated with a contributors and returns a list of EPrint ID 
`, repoID, repoID)
}

func advisorIDDocument(repoID string) string {
	return fmt.Sprintf(`
- '/%s/advisor-id' returns a list of advisor ids in the eprints repository
- '/%s/advisor-id/{ADVISOR_ID}' scans the name advisor id field associated with advisors and returns a list of EPrint ID 
`, repoID, repoID)
}

func committeeIDDocument(repoID string) string {
	return fmt.Sprintf(`
- '/%s/committee-id' returns a list of committee id in EPrints repository
- '/%s/committee-id/{COMMITTEE_ID}' scans the committee id field associated with committee members and returns a list of EPrint ID
`, repoID, repoID)
}

func groupIDDocument(repoID string) string {
	return fmt.Sprintf(`
- '/%s/group-id' returns a list of group ids in EPrints repository
- '/%s/group-id/{GROUP_ID}' this scans group ID and returns a list of EPrint IDs associated with the group
`, repoID, repoID)
}

func funderIDDocument(repoID string) string {
	return fmt.Sprintf(`
- '/%s/funder-id' returns a list of funders in the EPrints repository
- '/%s/funder-id/{FUNDER_ID}' returns a list of EPrint IDs associated with the funder
`, repoID, repoID)
}

func grantNumberDocument(repoID string) string {
	return fmt.Sprintf(`
- '/%s/grant-number' returns a list of grant numbers in EPrints repository
- '/%s/grant-number/{GRANT_NUMBER}' returns a list of EPrint IDs associated with the grant number
`, repoID, repoID)
}

func creatorNameDocument(repoID string) string {
	return fmt.Sprintf(`
- '/%s/creator-name/{FAMILY_NAME}/{GIVEN_NAME}' scans the name fields associated with creators and returns a list of EPrint ID
`, repoID)
}
func editorNameDocument(repoID string) string {
	return fmt.Sprintf(`
- '/%s/editor-name/{FAMILY_NAME}/{GIVEN_NAME}' scans the family and given name field associated with a editors and returns a list of EPrint ID
`, repoID)
}

func contributorNameDocument(repoID string) string {
	return fmt.Sprintf(`
- '/%s/contributor-name/{FAMILY_NAME}/{GIVEN_NAME}' scans the family and given name field associated with a contributors and returns a list of EPrint ID
`, repoID)
}

func advisorNameDocument(repoID string) string {
	return fmt.Sprintf(`
- '/%s/advisor-name/{FAMILY_NAME}/{GIVEN_NAME}' scans the name fields associated with advisors returns a list of EPrint ID
`, repoID)
}
func committeeNameDocument(repoID string) string {
	return fmt.Sprintf(`
- '/%s/committee-name/{FAMILY_NAME}/{GIVEN_NAME}' scans the family and given name fields associated with committee members and returns a list of EPrint ID
`, repoID)
}

func pubmedIDDocument(repoID string) string {
	return fmt.Sprintf(`
- '/%s/pmid/{PMID}' with the "pmid" field in the EPrint table, it refers to PubMed is an index of the biomedical literature.
`, repoID)
}

func pubmedCentralIDDocument(repoID string) string {
	return fmt.Sprintf(`
- '/%s/pmcid/{PMCID}' with the "pmcid" field in the EPrint table, PMCID an Identifier to each full-text paper in PubMed Central Archive
`, repoID)
}

func issnDocument(repoID string) string {
	return fmt.Sprintf(`
- '/%s/issn' - returns a list of ISSN in repository
- '/%s/issn/{ISSN}' - returns a list eprint id for ISSN in repository
`, repoID, repoID)
}

func isbnDocument(repoID string) string {
	return fmt.Sprintf(`
- '/%s/isbn' - returns a list of ISBN in repository
- '/%s/isbn/{ISBN}' - returns a list eprint id for ISBN in repository
`, repoID, repoID)
}

func patentNumberDocument(repoID string) string {
	return fmt.Sprintf(`
- '/%s/patent-number' - return a list of patent numbers in repository
- '/%s/patent-number/{PATENT_NUMBER}' - return a list eprint ids for patent number in repository
`, repoID, repoID)
}

func patentApplicantDocument(repoID string) string {
	return fmt.Sprintf(`
- '/%s/patent-applicant' - return a list of patent applicants in repository
- '/%s/patent-applicant/{PATENT_APPLICANT}' - return a list eprint ids for patent applicant in repository
`, repoID, repoID)
}

func patentClassificationDocument(repoID string) string {
	return fmt.Sprintf(`
- '/%s/patent-classification' - return a list of patent classificatins in repository
- '/%s/patent-classification/{PATENT_CLASSIFICATION}' - return a list eprint ids for patent classification in repository
`, repoID, repoID)
}

func patentAssigneeDocument(repoID string) string {
	return fmt.Sprintf(`
- '/{REPO_ID}/patent-assignee' - return a list of patent assignee in repository
- '/{REPO_ID}/patent-assignee/{PATENT_ASSIGNEE}' - return a list eprint ids for patent assignee in repository
`, repoID, repoID)
}

func recordDocument(repoID string) string {
	return fmt.Sprintf(`Simplified Record
-----------------

This version of the API includes a simplified JSON record view. The
JSON represents the JSON model used in DataCite and InvenioRDMs.

- '/%s/record/{EPRINT_ID}' returns a complex JSON object representing the EPrint record identified by {EPRINT_ID}.
`, repoID)
}

func eprintDocument(repoID string) string {
	return fmt.Sprintf(`
EPrint Record
-------------

If the read/write is enabled for %q then this end point will
access a retrieve an EPrint record as EPrint XML via http GET.
You can create or replace an EPrint record using a POST containing
valid EPrint XML. Not the transaction with the repository is only
done via SQL. There is no facility to upload PDFs or other digital
files.  It does NOT use the Perl API.

GET:

- '/%s/eprint/{EPRINT_ID}' will retrieve an existing EPrint record as EPrint XML by building up an eprint record via SQL queries.

POST:

- '/%s/eprint/{EPRINT_ID}' will replace an existing EPrint, POST must be valid EPrint XML with a content type of "application/xml".
- '/%s/eprint/0' will create a new EPrint record. The POST must be valid EPrint XML with a content type of "application/xml".

`, repoID, repoID, repoID, repoID)
}
