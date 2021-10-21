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
 * ep3api.go defines an extended EPrints web server.
 */

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const (
	// timestamp holds the Format of a MySQL time field
	timestamp = "2006-01-02 15:04:05"
	datestamp = "2006-01-02"
)

var (
	config *Config
)

//
// Date/Time expansion
//
func expandAproxDate(dt string, roundDown bool) string {
	switch len(dt) {
	case 4:
		if roundDown {
			dt += "-01-01"
		} else {
			dt += "-12-31"
		}
	case 7:
		if roundDown {
			dt += "-01 00:00:00"
		} else {
			//FIXME: need to handle 28/29/30/31 day mounths
			dt += "-31"
		}
	}
	return dt
}

//
// DB SQL functions.
//

// sqlQueryIDs takes a repostory ID, a SQL statement and applies
// the args returning a list of EPrint ID or error.
func sqlQueryIDs(repoID string, stmt string, args ...interface{}) ([]int, error) {
	if db, ok := config.Connections[repoID]; ok == true {
		rows, err := db.Query(stmt, args...)
		if err != nil {
			return nil, fmt.Errorf("ERROR: query error (%q), %s", repoID, err)
		}
		defer rows.Close()
		eprintid := 0
		eprintIDs := []int{}
		for rows.Next() {
			err := rows.Scan(&eprintid)
			if (err == nil) && (eprintid > 0) {
				eprintIDs = append(eprintIDs, eprintid)
			} else {
				return nil, fmt.Errorf("ERROR: scan error (%q), %s", repoID, err)
			}
		}
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("ERROR: rows error (%q), %s", repoID, err)
		}
		if err != nil {
			return nil, fmt.Errorf("ERROR: query error (%q), %s", repoID, err)
		}
		return eprintIDs, nil
	}
	return nil, fmt.Errorf("Bad Request")
}

//
// Package functions take the collected data and package into an HTTP
// response.
//

func packageIDs(w http.ResponseWriter, repoID string, eprintIDs []int, err error) (int, error) {
	if err != nil {
		log.Printf("ERROR: (%s) query error, %s", repoID, err)
		return 500, fmt.Errorf("Internal Server Error")
	}
	src, err := json.MarshalIndent(eprintIDs, "", "  ")
	if err != nil {
		log.Printf("ERROR: marshal error (%q), %s", repoID, err)
		return 500, fmt.Errorf("Internal Server Error")
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", src)
	return 200, nil
}

func packageDocument(w http.ResponseWriter, src string) (int, error) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, src)
	return 200, nil
}

func packageJSON(w http.ResponseWriter, repoID string, src []byte, err error) (int, error) {
	if err != nil {
		log.Printf("ERROR: (%s) package JSON error, %s", repoID, err)
		return 500, fmt.Errorf("Internal Server Error")
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", src)
	return 200, nil
}

func packageXML(w http.ResponseWriter, repoID string, src []byte, err error) (int, error) {
	if err != nil {
		log.Printf("ERROR: (%s) package JSON error, %s", repoID, err)
		return 500, fmt.Errorf("Internal Server Error")
	}
	w.Header().Set("Content-Type", "application/xml")
	fmt.Fprintln(w, `<?xml version='1.0' encoding='utf-8'?>`)
	fmt.Fprintf(w, "%s", src)
	return 200, nil
}

//
// EPrints Metadata Structure
//

// EPrintTablesAndColumns takes a DB connection and repoID then builds a map[string][]string{}
// structure representing the tables and their columns available in a EPrints Repository
func EPrintTablesAndColumns(db *sql.DB, repoID string) (map[string][]string, error) {
	data := map[string][]string{}
	stmt := `SHOW TABLES LIKE "eprint%"`
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("SQL(%q), %s", repoID, err)
	}
	tables := []string{}
	for rows.Next() {
		tableName := ""
		if err := rows.Scan(&tableName); err == nil {
			tables = append(tables, tableName)
		}
	}
	rows.Close()

	for _, tableName := range tables {
		data[tableName] = []string{}
		stmt := fmt.Sprintf(`SHOW COLUMNS IN %s`, tableName)
		cRows, err := db.Query(stmt)
		if err != nil {
			return nil, fmt.Errorf("SQL(%q), %s", repoID, err)
		}
		columns := []string{}
		var (
			colName, f1, f2, f3, f5 string
			f4                      interface{}
		)
		for cRows.Next() {
			//colName, f1, f2, f3, f4, f5 = &"", &"", &"", &"", nil, &""
			if err := cRows.Scan(&colName, &f1, &f2, &f3, &f4, &f5); err != nil {
				log.Printf("cRows.Scan() error: %s", err)
			} else {
				columns = append(columns, colName)
			}
		}
		data[tableName] = columns
		cRows.Close()
	}
	// We need to add the document set of tables too.
	stmt = `SHOW TABLES LIKE "document%"`
	rows, err = db.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("SQL(%q), %s", repoID, err)
	}
	tables = []string{}
	for rows.Next() {
		tableName := ""
		if err := rows.Scan(&tableName); err == nil {
			tables = append(tables, tableName)
		}
	}
	rows.Close()

	for _, tableName := range tables {
		data[tableName] = []string{}
		stmt := fmt.Sprintf(`SHOW COLUMNS IN %s`, tableName)
		cRows, err := db.Query(stmt)
		if err != nil {
			return nil, fmt.Errorf("SQL(%q), %s", repoID, err)
		}
		columns := []string{}
		var (
			colName, f1, f2, f3, f5 string
			f4                      interface{}
		)
		for cRows.Next() {
			//colName, f1, f2, f3, f4, f5 = &"", &"", &"", &"", nil, &""
			if err := cRows.Scan(&colName, &f1, &f2, &f3, &f4, &f5); err != nil {
				log.Printf("cRows.Scan() error: %s", err)
			} else {
				columns = append(columns, colName)
			}
		}
		data[tableName] = columns
		cRows.Close()
	}
	// We need to add the files set of tables too.
	stmt = `SHOW TABLES LIKE "file%"`
	rows, err = db.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("SQL(%q), %s", repoID, err)
	}
	tables = []string{}
	for rows.Next() {
		tableName := ""
		if err := rows.Scan(&tableName); err == nil {
			tables = append(tables, tableName)
		}
	}
	rows.Close()

	for _, tableName := range tables {
		data[tableName] = []string{}
		stmt := fmt.Sprintf(`SHOW COLUMNS IN %s`, tableName)
		cRows, err := db.Query(stmt)
		if err != nil {
			return nil, fmt.Errorf("SQL(%q), %s", repoID, err)
		}
		columns := []string{}
		var (
			colName, f1, f2, f3, f5 string
			f4                      interface{}
		)
		for cRows.Next() {
			//colName, f1, f2, f3, f4, f5 = &"", &"", &"", &"", nil, &""
			if err := cRows.Scan(&colName, &f1, &f2, &f3, &f4, &f5); err != nil {
				log.Printf("cRows.Scan() error: %s", err)
			} else {
				columns = append(columns, colName)
			}
		}
		data[tableName] = columns
		cRows.Close()
	}
	return data, nil
}

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

- '/{REPO_ID}/doi/{DOI}' with the adoption of EPrints "doi" field in the EPrint table it makes sense to have a quick translation of DOI to EPrint id for a given EPrints repository. 
- '/{REPO_ID}/creator-id/{CREATOR_ID}' scans the name creator id field associated with creators and returns a list of EPrint ID 
- '/{REPO_ID}/creator-orcid/{ORCID}' scans the "orcid" field associated with creators and returns a list of EPrint ID 
- '/{REPO_ID}/editor-id/{CREATOR_ID}' scans the name creator id field associated with editors and returns a list of EPrint ID 
- '/{REPO_ID}/contributor-id/{CONTRIBUTOR_ID}' scans the "id" field associated with a contributors and returns a list of EPrint ID 
- '/{REPO_ID}/advisor-id/{ADVISOR_ID}' scans the name advisor id field associated with advisors and returns a list of EPrint ID 
- '/{REPO_ID}/committee-id/{COMMITTEE_ID}' scans the committee id field associated with committee members and returns a list of EPrint ID
- '/{REPO_ID}/group-id/{GROUP_ID}' this scans group ID and returns a list of EPrint IDs associated with the group
- '/{REPO_ID}/funder-id/{FUNDER_ID}' returns a list of EPrint IDs associated with the funder's ROR
- '/{REPO_ID}/grant-number/{GRANT_NUMBER}' returns a list of EPrint IDs associated with the grant number

Change Events
-------------

The follow API end points would facilitate faster updates to our feeds platform as well as allow us to create a separate public view of our EPrint repository content.

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

To run the web service create a JSON file named settings.ini in the current directory where you're invoking _ep3apid_ from. The web service can be started with running

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
	return fmt.Sprintf(`'/%s/creator-id/{CREATOR_ID}' scans the name creator id field associated with creators and returns a list of EPrint ID`, repoID)
}

func creatorORCIDDocument(repoID string) string {
	return fmt.Sprintf(`'/%s/creator-orcid/{ORCID}' scans the "orcid" field associated with creators and returns a list of EPrint ID`, repoID)
}

func editorIDDocument(repoID string) string {
	return fmt.Sprintf(`'/%s/editor-id/{CREATOR_ID}' scans the name creator id field associated with editors and returns a list of EPrint ID`, repoID)
}

func contributorIDDocument(repoID string) string {
	return fmt.Sprintf(`'/%s/contributor-id/{CONTRIBUTOR_ID}' scans the "id" field associated with a contributors and returns a list of EPrint ID`, repoID)
}

func advisorIDDocument(repoID string) string {
	return fmt.Sprintf(`'/%s/advisor-id/{ADVISOR_ID}' scans the name advisor id field associated with advisors and returns a list of EPrint ID`, repoID)
}

func committeeIDDocument(repoID string) string {
	return fmt.Sprintf(`'/%s/committee-id/{COMMITTEE_ID}' scans the committee id field associated with committee members and returns a list of EPrint ID`, repoID)
}

func groupIDDocument(repoID string) string {
	return fmt.Sprintf(`'/%s/group-id/{GROUP_ID}' this scans group ID and returns a list of EPrint IDs associated with the group`, repoID)
}

func grantNumberDocument(repoID string) string {
	return fmt.Sprintf(`'/%s/grant-number/{GRANT_NUMBER}' returns a list of EPrint IDs associated with the grant number`, repoID)
}

func creatorNameDocument(repoID string) string {
	return fmt.Sprintf(`'/%s/creator-name/{FAMILY_NAME}/{GIVEN_NAME}' scans the name fields associated with creators and returns a list of EPrint ID `, repoID)
}
func editorNameDocument(repoID string) string {
	return fmt.Sprintf(`'/%s/editor-name/{FAMILY_NAME}/{GIVEN_NAME}' scans the family and given name field associated with a editors and returns a list of EPrint ID`, repoID)
}

func contributorNameDocument(repoID string) string {
	return fmt.Sprintf(`'/%s/contributor-name/{FAMILY_NAME}/{GIVEN_NAME}' scans the family and given name field associated with a contributors and returns a list of EPrint ID`, repoID)
}

func advisorNameDocument(repoID string) string {
	return fmt.Sprintf(`'/%s/advisor-name/{FAMILY_NAME}/{GIVEN_NAME}' scans the name fields associated with advisors returns a list of EPrint ID`, repoID)
}
func committeeNameDocument(repoID string) string {
	return fmt.Sprintf(`'/%s/committee-name/{FAMILY_NAME}/{GIVEN_NAME}' scans the family and given name fields associated with committee members and returns a list of EPrint ID`, repoID)
}

func pubmedDocument(repoID string) string {
	return fmt.Sprintf(`'/%s/pubmed/{PUBMED_ID}' returns a list of EPrint IDs associated with the PubMed ID`, repoID)
}

func issnDocument(repoID string) string {
	return fmt.Sprintf(`'/%s/issn/{ISSN}' returns a list of EPrint IDs associated with the ISSN`, repoID)
}

func isbnDocument(repoID string) string {
	return fmt.Sprintf(`'/%s/isbn/{ISSN}' returns a list of EPrint IDs associated with the ISSN`, repoID)
}

func patentNumberDocument(repoID string) string {
	return fmt.Sprintf(`'/%s/patent-number/{PATENT_NUMBER}' returns a list of EPrint IDs associated with the patent number`, repoID)
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

//
// Expose the repositories available
//

func repositoriesEndPoint(w http.ResponseWriter, r *http.Request) (int, error) {
	if strings.HasSuffix(r.URL.Path, "/help") {
		return packageDocument(w, repositoryDocument())
	}
	repositories := []string{}
	for repository, _ := range config.Repositories {
		repositories = append(repositories, repository)
	}
	src, err := json.MarshalIndent(repositories, "", "  ")
	if err != nil {
		return 500, fmt.Errorf("Internal Server Error, %s", err)
	}
	return packageDocument(w, fmt.Sprintf("%s", src))
}

func repositoryEndPoint(w http.ResponseWriter, r *http.Request, repoID string) (int, error) {
	if repoID == "" || strings.HasSuffix(r.URL.Path, "/help") {
		return packageDocument(w, repositoryDocument())
	}
	if db, ok := config.Connections[repoID]; ok == true {
		data, err := EPrintTablesAndColumns(db, repoID)
		src, err := json.MarshalIndent(data, "", "    ")
		return packageJSON(w, repoID, src, err)
	}
	return 404, fmt.Errorf("not found")
}

//
// End Point handles (route as defined `/{REPO_ID}/{END-POINT}/{ARGS}`)
//

func createdEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 {
		return packageDocument(w, updatedDocument(repoID))
	}
	if (len(args) < 1) || (len(args) > 2) {
		return 400, fmt.Errorf("Bad Request")
	}
	var err error
	end := time.Now()
	start := time.Now()
	if (len(args) > 0) && (args[0] != "now") {
		start, err = time.Parse(timestamp, args[0])
		if err != nil {
			return 400, fmt.Errorf("Bad Request, (start) %s", err)
		}
	}
	if (len(args) > 1) && (args[1] != "now") {
		end, err = time.Parse(timestamp, args[1])
		if err != nil {
			return 400, fmt.Errorf("Bad Request, (end) %s", err)
		}
	}
	eprintIDs, err := sqlQueryIDs(repoID,
		`SELECT eprintid FROM eprint WHERE 
(CONCAT(datestamp_year, "-", 
LPAD(IFNULL(datestamp_month, 1), 2, "0"), "-", 
LPAD(IFNULL(datestamp_day, 1), 2, "0"), " ", 
LPAD(IFNULL(datestamp_hour, 0), 2, "0"), ":", 
LPAD(IFNULL(datestamp_minute, 0), 2, "0"), ":", 
LPAD(IFNULL(datestamp_second, 0), 2, "0")) >= ?) AND 
(CONCAT(datestamp_year, "-", 
LPAD(IFNULL(datestamp_month, 12), 2, "0"), "-", 
LPAD(IFNULL(datestamp_day, 28), 2, "0"), " ", 
LPAD(IFNULL(datestamp_hour, 23), 2, "0"), ":", 
LPAD(IFNULL(datestamp_minute, 59), 2, "0"), ":", 
LPAD(IFNULL(datestamp_second, 59), 2, "0")) <= ?)`,
		start.Format(timestamp), end.Format(timestamp))
	return packageIDs(w, repoID, eprintIDs, err)
}

func updatedEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 {
		return packageDocument(w, updatedDocument(repoID))
	}
	if (len(args) < 1) || (len(args) > 2) {
		return 400, fmt.Errorf("Bad Request")
	}
	var err error
	end := time.Now()
	start := time.Now()
	if (len(args) > 0) && (args[0] != "now") {
		start, err = time.Parse(timestamp, args[0])
		if err != nil {
			return 400, fmt.Errorf("Bad Request, (start) %s", err)
		}
	}
	if (len(args) > 1) && (args[1] != "now") {
		end, err = time.Parse(timestamp, args[1])
		if err != nil {
			return 400, fmt.Errorf("Bad Request, (end) %s", err)
		}
	}
	eprintIDs, err := sqlQueryIDs(repoID,
		`SELECT eprintid FROM eprint WHERE 
(CONCAT(lastmod_year, "-", 
LPAD(IFNULL(lastmod_month, 1), 2, "0"), "-", 
LPAD(IFNULL(lastmod_day, 1), 2, "0"), " ", 
LPAD(IFNULL(lastmod_hour, 0), 2, "0"), ":", 
LPAD(IFNULL(lastmod_minute, 0), 2, "0"), ":", 
LPAD(IFNULL(lastmod_second, 0), 2, "0")) >= ?) AND 
(CONCAT(lastmod_year, "-", 
LPAD(IFNULL(lastmod_month, 12), 2, "0"), "-", 
LPAD(IFNULL(lastmod_day, 28), 2, "0"), " ", 
LPAD(IFNULL(lastmod_hour, 23), 2, "0"), ":", 
LPAD(IFNULL(lastmod_minute, 59), 2, "0"), ":", 
LPAD(IFNULL(lastmod_second, 59), 2, "0")) <= ?)`,
		start.Format(timestamp), end.Format(timestamp))
	return packageIDs(w, repoID, eprintIDs, err)
}

func deletedEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 {
		return packageDocument(w, deletedDocument(repoID))
	}
	if (len(args) < 1) || (len(args) > 2) {
		return 400, fmt.Errorf("Bad Request")
	}
	var err error
	end := time.Now()
	start := time.Now()
	if (len(args) > 0) && (args[0] != "now") {
		start, err = time.Parse(timestamp, args[0])
		if err != nil {
			return 400, fmt.Errorf("Bad Request, (start) %s", err)
		}
	}
	if (len(args) > 1) && (args[1] != "now") {
		end, err = time.Parse(timestamp, args[1])
		if err != nil {
			return 400, fmt.Errorf("Bad Request, (end) %s", err)
		}
	}
	eprintIDs, err := sqlQueryIDs(repoID,
		`SELECT eprintid FROM eprint WHERE (eprint_status = "deletion") AND 
(CONCAT(lastmod_year, "-", 
LPAD(IFNULL(lastmod_month, 1), 2, "0"), "-", 
LPAD(IFNULL(lastmod_day, 1), 2, "0"), " ", 
LPAD(IFNULL(lastmod_hour, 0), 2, "0"), ":", 
LPAD(IFNULL(lastmod_minute, 0), 2, "0"), ":", 
LPAD(IFNULL(lastmod_second, 0), 2, "0")) >= ?) AND 
(CONCAT(lastmod_year, "-", 
LPAD(IFNULL(lastmod_month, 12), 2, "0"), "-", 
LPAD(IFNULL(lastmod_day, 28), 2, "0"), " ", 
LPAD(IFNULL(lastmod_hour, 23), 2, "0"), ":", 
LPAD(IFNULL(lastmod_minute, 59), 2, "0"), ":",
LPAD(IFNULL(lastmod_second, 59), 2, "0")) <= ?)`,
		start.Format(timestamp), end.Format(timestamp))
	return packageIDs(w, repoID, eprintIDs, err)
}

func pubdateEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 {
		return packageDocument(w, pubdateDocument(repoID))
	}
	if (len(args) < 1) || (len(args) > 2) {
		return 400, fmt.Errorf("Bad Request")
	}
	var (
		err error
		dt  string
	)
	end := time.Now()
	start := time.Now()
	if (len(args) > 0) && (args[0] != "now") {
		dt = expandAproxDate(args[0], true)
		start, err = time.Parse(datestamp, dt)
		if err != nil {
			return 400, fmt.Errorf("Bad Request, (start date) %s", err)
		}
	}
	if (len(args) > 1) && (args[1] != "now") {
		dt = expandAproxDate(args[1], false)
		end, err = time.Parse(datestamp, dt)
		if err != nil {
			return 400, fmt.Errorf("Bad Request, (end date) %s", err)
		}
	}
	eprintIDs, err := sqlQueryIDs(repoID,
		`SELECT eprintid FROM eprint
WHERE ((date_type) = "published") AND 
(CONCAT(date_year, "-", 
LPAD(IFNULL(date_month, 1), 2, "0"), "-", 
LPAD(IFNULL(date_day, 1), 2, "0")) >= ?) AND 
(CONCAT(date_year, "-", 
LPAD(IFNULL(date_month, 12), 2, "0"), "-", 
LPAD(IFNULL(date_day, 28), 2, "0")) <= ?)`,
		start.Format(datestamp), end.Format(datestamp))
	return packageIDs(w, repoID, eprintIDs, err)
}

func doiEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 {
		return packageDocument(w, doiDocument(repoID))
	}
	if len(args) < 1 {
		return 400, fmt.Errorf("Bad Request")
	}
	var err error
	doi := strings.Join(args, "/")
	//FIXME: Should validate DOI format ...
	eprintIDs, err := sqlQueryIDs(repoID, `SELECT eprintid FROM eprint WHERE doi = ?`, doi)
	return packageIDs(w, repoID, eprintIDs, err)
}

func creatorIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 {
		return packageDocument(w, creatorIDDocument(repoID))
	}
	if len(args) != 1 {
		return 400, fmt.Errorf("Bad Request")
	}
	var err error
	eprintIDs, err := sqlQueryIDs(repoID, `SELECT eprintid FROM eprint_creators_id WHERE creators_id = ?`, args[0])
	return packageIDs(w, repoID, eprintIDs, err)
}

func creatorORCIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 {
		return packageDocument(w, creatorORCIDDocument(repoID))
	}
	if len(args) != 1 {
		return 400, fmt.Errorf("Bad Request")
	}
	var err error
	//FIXME: Should validate ORCID format ...
	eprintIDs, err := sqlQueryIDs(repoID, `SELECT eprintid FROM eprint_creators_orcid WHERE creators_orcid = ?`, args[0])
	return packageIDs(w, repoID, eprintIDs, err)
}

func editorIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 {
		return packageDocument(w, editorIDDocument(repoID))
	} else if len(args) != 1 {
		return 400, fmt.Errorf("Bad Request")
	}
	var err error
	eprintIDs, err := sqlQueryIDs(repoID, `SELECT eprintid FROM eprint_editors_id WHERE editors_id = ?`, args[0])
	return packageIDs(w, repoID, eprintIDs, err)
}

func contributorIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 {
		return packageDocument(w, contributorIDDocument(repoID))
	}
	if len(args) != 1 {
		return 400, fmt.Errorf("Bad Request")
	}
	var err error
	eprintIDs, err := sqlQueryIDs(repoID, `SELECT eprintid FROM eprint_contibutors_id WHERE contibutors_id = ?`, args[0])
	return packageIDs(w, repoID, eprintIDs, err)
}

func advisorIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 {
		return packageDocument(w, advisorIDDocument(repoID))
	}
	if len(args) != 1 {
		return 400, fmt.Errorf("Bad Request")
	}
	var err error
	eprintIDs, err := sqlQueryIDs(repoID, `SELECT eprintid FROM eprint_thesis_advisor_id WHERE thesis_advisor_id = ?`, args[0])
	return packageIDs(w, repoID, eprintIDs, err)
}

func committeeIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 {
		return packageDocument(w, committeeIDDocument(repoID))
	}
	if len(args) != 1 {
		return 400, fmt.Errorf("Bad Request")
	}
	var err error
	eprintIDs, err := sqlQueryIDs(repoID, `SELECT eprintid FROM eprint_thesis_committee_id WHERE thesis_committee_id = ?`, args[0])
	return packageIDs(w, repoID, eprintIDs, err)
}

func groupIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 {
		return packageDocument(w, groupIDDocument(repoID))
	}
	if len(args) != 1 {
		return 400, fmt.Errorf("Bad Request")
	}
	var err error
	eprintIDs, err := sqlQueryIDs(repoID, `SELECT eprintid FROM eprint_local_group WHERE local_group = ?`, args[0])
	return packageIDs(w, repoID, eprintIDs, err)
}

func grantNumberEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 {
		return packageDocument(w, grantNumberDocument(repoID))
	}
	if len(args) != 1 {
		return 400, fmt.Errorf("Bad Request")
	}
	var err error
	eprintIDs, err := sqlQueryIDs(repoID, `SELECT eprintid FROM eprint_funders_grant_number WHERE funders_grant_number = ?`, args[0])
	return packageIDs(w, repoID, eprintIDs, err)
}

func creatorNameEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 {
		return packageDocument(w, creatorNameDocument(repoID))
	}
	if len(args) != 2 {
		return 400, fmt.Errorf("Bad Request")
	}
	var err error
	eprintIDs, err := sqlQueryIDs(repoID, `SELECT eprintid FROM eprint_creator_name WHERE family_name = ? AND given_name = ?`, args[0], args[1])
	return packageIDs(w, repoID, eprintIDs, err)
}

func editorNameEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	//- `/<REPO_ID>/editor-name/<FAMILY_NAME>/<GIVEN_NAME>` scans the family and given name field associated with a editors and returns a list of EPrint ID
	if len(args) == 0 {
		return packageDocument(w, editorNameDocument(repoID))
	}
	if len(args) != 2 {
		return 400, fmt.Errorf("Bad Request")
	}
	var err error
	eprintIDs, err := sqlQueryIDs(repoID, `SELECT eprintid FROM eprint_editor_name WHERE family_name = ? AND given_name = ?`, args[0], args[1])
	return packageIDs(w, repoID, eprintIDs, err)
}

func contributorNameEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 {
		return packageDocument(w, contributorNameDocument(repoID))
	}
	if len(args) != 2 {
		return 400, fmt.Errorf("Bad Request")
	}
	var err error
	eprintIDs, err := sqlQueryIDs(repoID, `SELECT eprintid FROM eprint_contributors_name WHERE family_name = ? AND given_name = ?`, args[0], args[1])
	return packageIDs(w, repoID, eprintIDs, err)
}

func advisorNameEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 {
		return packageDocument(w, advisorNameDocument(repoID))
	}
	if len(args) != 2 {
		return 400, fmt.Errorf("Bad Request")
	}
	var err error
	eprintIDs, err := sqlQueryIDs(repoID, `SELECT eprintid FROM eprint_thesis_advisor_name WHERE family_name = ? AND given_name = ?`, args[0], args[1])
	return packageIDs(w, repoID, eprintIDs, err)
}

func committeeNameEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 {
		return packageDocument(w, committeeNameDocument(repoID))
	}
	if len(args) != 2 {
		return 400, fmt.Errorf("Bad Request")
	}
	var err error
	eprintIDs, err := sqlQueryIDs(repoID, `SELECT eprintid FROM eprint_thesis_committee_name WHERE family_name = ? AND given_name = ?`, args[0], args[1])
	return packageIDs(w, repoID, eprintIDs, err)
}

func pubmedEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 {
		return packageDocument(w, pubmedDocument(repoID))
	}
	if len(args) != 1 {
		return 400, fmt.Errorf("Bad Request")
	}
	var err error
	eprintIDs, err := sqlQueryIDs(repoID, `SELECT eprintid FROM eprint WHERE pmc_id = ?`, args[0])
	return packageIDs(w, repoID, eprintIDs, err)
}

func issnEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 {
		return packageDocument(w, issnDocument(repoID))
	}
	if len(args) != 1 {
		return 400, fmt.Errorf("Bad Request")
	}
	var err error
	eprintIDs, err := sqlQueryIDs(repoID, `SELECT eprintid FROM eprint WHERE issn = ?`, args[0])
	return packageIDs(w, repoID, eprintIDs, err)
}

func isbnEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 {
		return packageDocument(w, isbnDocument(repoID))
	}
	if len(args) != 1 {
		return 400, fmt.Errorf("Bad Request")
	}
	var err error
	eprintIDs, err := sqlQueryIDs(repoID, `SELECT eprintid FROM eprint WHERE isbn = ?`, args[0])
	return packageIDs(w, repoID, eprintIDs, err)
}

func patentNumberEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 {
		return packageDocument(w, patentNumberDocument(repoID))
	}
	if len(args) != 1 {
		return 400, fmt.Errorf("Bad Request")
	}
	var err error
	eprintIDs, err := sqlQueryIDs(repoID, `SELECT eprintid FROM eprint WHERE patent_number = ?`, args[0])
	return packageIDs(w, repoID, eprintIDs, err)
}

//
// Record End Point is experimental and may not make it to the
// release version of eprinttools. It accepts a EPrint ID and returns
// a simplfied JSON object.
//
func recordEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 || strings.HasSuffix(r.URL.Path, "/help") {
		return packageDocument(w, recordDocument(repoID))
	}
	if len(args) != 1 {
		return 400, fmt.Errorf("Bad Request")
	}
	baseURL := ""
	dataSource, ok := config.Repositories[repoID]
	if ok {
		baseURL = dataSource.BaseURL
	}
	eprintID, err := strconv.Atoi(args[0])
	if err != nil {
		return 400, fmt.Errorf("Bad Request, eprint id invalid, %s", err)
	}
	eprint, err := CrosswalkSQLToEPrint(repoID, baseURL, eprintID)
	if err != nil {
		log.Printf("CrosswalkSQLToEPrint Error: %s\n", err)
		return 500, fmt.Errorf("Internal Server Error")
	}
	//FIXME: we should just be a simple JSON from SQL ...
	simple, err := CrosswalkEPrintToRecord(eprint)
	if err != nil {
		return 500, fmt.Errorf("Internal Server Error")
	}
	src, err := json.MarshalIndent(simple, "", "    ")
	return packageJSON(w, repoID, src, err)
}

//
// EPrint End Point is an experimental read/write end point provided
// in the extended EPrint API.  It reads/writes EPrint data structures
// based on SQL calls to the MySQL database for a given EPrints repository.
// Note this end point has to be enabled for the repository in the
// configuraiton.
//
func eprintEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 || repoID == "" || strings.HasSuffix(r.URL.Path, "/help") {
		return packageDocument(w, eprintDocument(repoID))
	}
	contentType := r.Header.Get("Content-Type")
	writeAccess := false
	dataSource, ok := config.Repositories[repoID]
	if ok == true {
		writeAccess = dataSource.Write
	} else {
		log.Printf("Data Source not found for %q", repoID)
		return 404, fmt.Errorf("not found")
	}
	if (r.Method != "GET" && r.Method != "POST") ||
		(r.Method == "POST" && writeAccess == false) {
		return 405, fmt.Errorf("method not allowed")
	}
	if len(args) != 1 {
		return 400, fmt.Errorf("bad request, missing eprint id in path")
	}
	eprintID, err := strconv.Atoi(args[0])
	if err != nil {
		return 400, fmt.Errorf("bad request, eprint id %q not valid", eprintID)
	}
	if r.Method == "POST" {
		return 500, fmt.Errorf("not implemented, POST")
	}
	eprint, err := CrosswalkSQLToEPrint(repoID, dataSource.BaseURL, eprintID)
	if err != nil {
		return 404, fmt.Errorf("not found, %s", err)

	}
	switch contentType {
	case "application/json":
		src, err := json.MarshalIndent(eprint, "", "    ")
		return packageJSON(w, repoID, src, err)
	case "application/xml":
		eprints := new(EPrints)
		eprints.AddEPrint(eprint)
		src, err := xml.MarshalIndent(eprints, "", "    ")
		return packageXML(w, repoID, src, err)
	case "":
		eprints := new(EPrints)
		eprints.AddEPrint(eprint)
		src, err := xml.MarshalIndent(eprints, "", "    ")
		return packageXML(w, repoID, src, err)
	default:
		return 415, fmt.Errorf("unsupported media type, %q", contentType)
	}
}

// The following define the API as a service handling errors,
// routes and logging.
//
func logRequest(r *http.Request, status int, err error) {
	q := r.URL.Query()
	errStr := "OK"
	if err != nil {
		errStr = fmt.Sprintf("%s", err)
	}
	if len(q) > 0 {
		log.Printf("%s %s RemoteAddr: %s UserAgent: %s Query: %+v Response: %d %s", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), q, status, errStr)
	} else {
		log.Printf("%s %s RemoteAddr: %s UserAgent: %s Response: %d %s", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), status, errStr)
	}
}

func handleError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, `ERROR: %d %s`, statusCode, err)
}

func routeEndPoints(w http.ResponseWriter, r *http.Request) (int, error) {
	parts := strings.Split(r.URL.Path, "/")
	// Create args by removing empty strings from path parts
	repoID, endPoint, args := "", "", []string{}
	for _, arg := range parts {
		if arg != "" {
			// FIXME: URL decode that path part
			args = append(args, arg)
		}
	}
	if len(args) == 0 {
		return packageDocument(w, readmeDocument())
	}
	if len(args) == 1 {
		return packageDocument(w, strings.ReplaceAll(readmeDocument(), "<REPO_ID>", args[0]))
	}
	// Expected URL structure of `/<REPO_ID>/<END_POINT>/<ARGS>`
	if len(args) == 2 {
		repoID, endPoint, args = args[0], args[1], []string{}
	} else {
		repoID, endPoint, args = args[0], args[1], args[2:]
	}
	if routes, hasRepo := config.Routes[repoID]; hasRepo == true {
		// Confirm we have a route
		if fn, hasRoute := routes[endPoint]; hasRoute == true {
			// Confirm we have a DB connection
			if _, hasConnection := config.Connections[repoID]; hasConnection == true {
				return fn(w, r, repoID, args)
			}
		}
	}
	return 404, fmt.Errorf("Not Found")
}

func api(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		statusCode int
	)
	if r.Method != "GET" && r.Method != "POST" {
		statusCode, err = 405, fmt.Errorf("Method Not Allowed")
		handleError(w, statusCode, err)
	} else {
		switch {
		case r.URL.Path == "/favicon.ico":
			statusCode, err = 200, nil
			fmt.Fprintf(w, "")
			//statusCode, err = 404, fmt.Errorf("Not Found")
			//handleError(w, statusCode, err)
		case strings.HasPrefix(r.URL.Path, "/repositories"):
			statusCode, err = repositoriesEndPoint(w, r)
			if err != nil {
				handleError(w, statusCode, err)
			}
		case strings.HasPrefix(r.URL.Path, "/repository"):
			repoID := strings.TrimPrefix(r.URL.Path, "/repository/")
			statusCode, err = repositoryEndPoint(w, r, repoID)
			if err != nil {
				handleError(w, statusCode, err)
			}
		default:
			statusCode, err = routeEndPoints(w, r)
			if err != nil {
				handleError(w, statusCode, err)
			}
		}
	}
	logRequest(r, statusCode, err)
}

func InitExtendedAPI(settings string) error {
	var err error
	// NOTE: This reads the settings file and creates a global
	// config object.
	config, err = LoadConfig(settings)
	if err != nil {
		return fmt.Errorf("Failed to load %q, %s", settings, err)
	}
	if config == nil {
		return fmt.Errorf("Missing configuration")
	}
	if config.Hostname == "" {
		return fmt.Errorf("Hostings hostname for service")
	}
	if config.Repositories == nil || len(config.Repositories) < 1 {
		return fmt.Errorf(`Missing "repositories" configuration`)
	}
	if config.Connections == nil {
		config.Connections = map[string]*sql.DB{}
	}
	if config.Routes == nil {
		config.Routes = map[string]map[string]func(http.ResponseWriter, *http.Request, string, []string) (int, error){}
	}
	// This is a map endpoints and point handlers.
	// This implements the registration design pattern
	routes := map[string]func(http.ResponseWriter, *http.Request, string, []string) (int, error){
		"created":          createdEndPoint,
		"updated":          updatedEndPoint,
		"deleted":          deletedEndPoint,
		"pubdate":          pubdateEndPoint,
		"doi":              doiEndPoint,
		"creator-id":       creatorIDEndPoint,
		"creator-orcid":    creatorORCIDEndPoint,
		"editor-id":        editorIDEndPoint,
		"contributor-id":   contributorIDEndPoint,
		"advisor-id":       advisorIDEndPoint,
		"committee-id":     committeeIDEndPoint,
		"group-id":         groupIDEndPoint,
		"grant-number":     grantNumberEndPoint,
		"creator-name":     creatorNameEndPoint,
		"editor-name":      editorNameEndPoint,
		"contributor-name": contributorNameEndPoint,
		"advisor-name":     advisorNameEndPoint,
		"commitee-name":    committeeNameEndPoint,
		"pubmed":           pubmedEndPoint,
		"issn":             issnEndPoint,
		"isbn":             isbnEndPoint,
		"patent-number":    patentNumberEndPoint,
		"record":           recordEndPoint,
		"eprint":           eprintEndPoint,
	}

	/* NOTE: We need a DB connection to MySQL for each
	   EPrints repository supported by the API
	   for access to MySQL */
	for repoID, dataSource := range config.Repositories {
		dataSourceName := dataSource.DSN
		// Setup DB connection for target repository
		db, err := sql.Open("mysql", dataSourceName)
		if err != nil {
			return fmt.Errorf("Could not open MySQL connection for %s, %s", repoID, err)
		}
		//log.Printf("Setting  DB connection to %q", repoID)
		//db.Ping()
		config.Connections[repoID] = db
		dataSource.TableMap, err = EPrintTablesAndColumns(db, repoID)
		if err != nil {
			return fmt.Errorf("failed to map table and columns for %q, %s", repoID, err)
		}
		// Add routes (end points) for the target repository
		for route, fn := range routes {
			if config.Routes[repoID] == nil {
				config.Routes[repoID] = map[string]func(http.ResponseWriter, *http.Request, string, []string) (int, error){}
			}
			config.Routes[repoID][route] = fn
		}
	}
	return nil
}

func RunExtendedAPI(appName string) error {
	/* Setup web server */
	log.Printf(`
%s %s

EPrints 3.3.x Extended API

Listening on http://%s

Press ctl-c to terminate.
`, appName, Version, config.Hostname)
	http.HandleFunc("/", api)
	return http.ListenAndServe(config.Hostname, nil)
}
