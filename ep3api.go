/**
 * ep3api.go defines an extended EPrints web server.
 */
package eprinttools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

//
// Service configuration management
//

// Config holds a configuration file structure used by EPrints Extended API
// Configuration file is expected to be in JSON format.
type Config struct {
	// Hostname for running service
	Hostname string `json:"hostname"`

	// Repositories are defined by a REPO_ID (string)
	// that points at a MySQL Db connection string
	Repositories map[string]string `json:"repositories"`

	// Connections is a map to database connections
	Connections map[string]*sql.DB `json:"-"`

	// Routes holds the mapping of end points to repository id
	// instances.
	Routes map[string]map[string]func(http.ResponseWriter, *http.Request, string, []string) (int, error) `json:"-"`
}

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
			dt += "-01"
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

//
// End point documentation
//
func readmeDocument() string {
	return fmt.Sprintf(`
EPrints 3.3.x extended API, eprinttools version %s

EPrints extended API
====================

The EPrints software package from University of Southampton provides a rich internal Perl API along with a RESTful web API. But getting specific lists of EPrint IDs is challenging. This API addresses this. The API returns lists of EPrint IDs as a JSON array or plain text documentaiton (like this page).  The EPrint IDs lists are not sorted. An empty JSON array means no EPrints IDs are available for that type of request. 

List repositories available
---------------------------

The general structure of URLs in the extended API is in the form

    http://<HOSTNAME>:<PORT>/<REPO_ID>/<END_POINT>/<PARAMETERS>

- <HOSTNAME> is normally (recommended) "localhost"
- <PORT> defaults to 8484
- <REPO_ID> is the label used to reference the repository name
- <END_POINT> is the list of end points provided the service
- <PARAMATERS> are any needed values for the end point one per path part

To see a list of available repositories use the "/repositories" end point. E.g.

   curl http://localhost:8484/repositories


Unique IDs to EPrint IDs
------------------------

The following URL end points are intended to take one unique identifier and map that to one or more EPrint IDs. This can be done because each unique ID  targeted can be identified by querying a single table in EPrints.  In addition the scan can return the complete results since all EPrint IDs are integers and returning all EPrint IDs in any of our repositories is sufficiently small to be returned in a single HTTP request.

- "/<REPO_ID>/doi/<DOI>" with the adoption of EPrints "doi" field in the EPrint table it makes sense to have a quick translation of DOI to EPrint id for a given EPrints repository. 
- "/<REPO_ID>/creator-id/<CREATOR_ID>" scans the name creator id field associated with creators and returns a list of EPrint ID 
- "/<REPO_ID>/creator-orcid/<ORCID>" scans the "orcid" field associated with creators and returns a list of EPrint ID 
- "/<REPO_ID>/editor-id/<CREATOR_ID>" scans the name creator id field associated with editors and returns a list of EPrint ID 
- "/<REPO_ID>/contributor-id/<CONTRIBUTOR_ID>" scans the "id" field associated with a contributors and returns a list of EPrint ID 
- "/<REPO_ID>/advisor-id/<ADVISOR_ID>" scans the name advisor id field associated with advisors and returns a list of EPrint ID 
- "/<REPO_ID>/committee-id/<COMMITTEE_ID>" scans the committee id field associated with committee members and returns a list of EPrint ID
- "/<REPO_ID>/group-id/<GROUP_ID>" this scans group ID and returns a list of EPrint IDs associated with the group
- "/<REPO_ID>/grant-number/<GRANT_NUMBER>" returns a list of EPrint IDs associated with the grant number
- "/<REPO_ID>/creator-name/<FAMILY_NAME>/<GIVEN_NAME>" scans the name fields associated with creators and returns a list of EPrint ID 
- "/<REPO_ID>/editor-name/<FAMILY_NAME>/<GIVEN_NAME>" scans the family and given name field associated with a editors and returns a list of EPrint ID 
- "/<REPO_ID>/contributor-name/<FAMILY_NAME>/<GIVEN_NAME>" scans the family and given name field associated with a contributors and returns a list of EPrint ID 
- "/<REPO_ID>/advisor-name/<FAMILY_NAME>/<GIVEN_NAME>" scans the name fields associated with advisors returns a list of EPrint ID 
- "/<REPO_ID>/committee-name/<FAMILY_NAME>/<GIVEN_NAME>" scans the family and given name fields associated with committee members and returns a list of EPrint ID
- "/<REPO_ID>/pubmed/<PUBMED_ID>" returns a list of EPrint IDs associated with the PubMed ID
- "/<REPO_ID>/issn/<ISSN>" returns a list of EPrint IDs associated with the ISSN
- "/<REPO_ID>/isbn/<ISSN>" returns a list of EPrint IDs associated with the ISSN
- "/<REPO_ID>/patent-number/<PATENT_NUMBER>" returns a list of EPrint IDs associated with the patent number

Change Events
-------------

The follow API end points would facilitate faster updates to our feeds platform as well as allow us to create a separate public view of our EPrint repository content.

- "/<REPO_ID>/updated/<TIMESTAMP>/<TIMESTAMP>" returns a list of EPrint IDs updated starting at the first timestamp (timestamps should have a resolution to the minute, e.g. "YYYY-MM-DD HH:MM:SS") through inclusive of the second timestmap (if the second is omitted the timestamp is assumed to be "now")
- "/<REPO_ID>/deleted/<TIMESTAMP>/<TIMESTAMP>" through the returns a list of EPrint IDs deleted starting at first timestamp through inclusive of the second timestamp, if the second timestamp is omitted it is assumed to be "now"
- "/<REPO_ID>/pubdate/<APROX_DATESTAMP>/<APPOX_DATESTAMP>" this query scans the EPrint table for records with publication starts starting with the first approximate date through inclusive of the second approximate date. If the second date is omitted it is assumed to be "today". Approximate dates my be expressed just the year (starting with Jan 1, ending with Dec 31), just the year and month (starting with first day of month ending with the last day) or year, month and day. The end returns zero or more EPrint IDs.

`, Version)
}

func updatedDocument(repoID string) string {
	return fmt.Sprintf(`"/%s/updated/<TIMESTAMP>/<TIMESTAMP>" returns a list of EPrint IDs updated starting at the first timestamp (timestamps should have a resolution to the minute, e.g. "YYYY-MM-DD HH:MM:SS") through inclusive of the second timestmap (if the second is omitted the timestamp is assumed to be "now")`, repoID)
}

func deletedDocument(repoID string) string {
	return fmt.Sprintf(`"/%s/deleted/<TIMESTAMP>/<TIMESTAMP>" through the returns a list of EPrint IDs deleted starting at first timestamp through inclusive of the second timestamp, if the second timestamp is omitted it is assumed to be "now"`, repoID)
}

func pubdateDocument(repoID string) string {
	return fmt.Sprintf(`"/%s/pubdate/<APROX_DATESTAMP>/<APPOX_DATESTAMP>" this query scans the EPrint table for records with publication starts starting with the first approximate date through inclusive of the second approximate date. If the second date is omitted it is assumed to be "today". Approximate dates my be expressed just the year (starting with Jan 1, ending with Dec 31), just the year and month (starting with first day of month ending with the last day) or year, month and day. The end returns zero or more EPrint IDs.`, repoID)
}

func doiDocument(repoID string) string {
	return fmt.Sprintf(`"/%s/doi/<DOI>" with the adoption of EPrints "doi" field in the EPrint table it makes sense to have a quick translation of DOI to EPrint id for a given EPrints repository.`, repoID)
}

func creatorIDDocument(repoID string) string {
	return fmt.Sprintf(`"/%s/creator-id/<CREATOR_ID>" scans the name creator id field associated with creators and returns a list of EPrint ID`, repoID)
}

func creatorORCIDDocument(repoID string) string {
	return fmt.Sprintf(`"/%s/creator-orcid/<ORCID>" scans the "orcid" field associated with creators and returns a list of EPrint ID`, repoID)
}

func editorIDDocument(repoID string) string {
	return fmt.Sprintf(`"/%s/editor-id/<CREATOR_ID>" scans the name creator id field associated with editors and returns a list of EPrint ID`, repoID)
}

func contributorIDDocument(repoID string) string {
	return fmt.Sprintf(`"/%s/contributor-id/<CONTRIBUTOR_ID>" scans the "id" field associated with a contributors and returns a list of EPrint ID`, repoID)
}

func advisorIDDocument(repoID string) string {
	return fmt.Sprintf(`"/%s/advisor-id/<ADVISOR_ID>" scans the name advisor id field associated with advisors and returns a list of EPrint ID`, repoID)
}

func committeeIDDocument(repoID string) string {
	return fmt.Sprintf(`"/%s/committee-id/<COMMITTEE_ID>" scans the committee id field associated with committee members and returns a list of EPrint ID`, repoID)
}

func groupIDDocument(repoID string) string {
	return fmt.Sprintf(`"/%s/group-id/<GROUP_ID>" this scans group ID and returns a list of EPrint IDs associated with the group`, repoID)
}

func grantNumberDocument(repoID string) string {
	return fmt.Sprintf(`"/%s/grant-number/<GRANT_NUMBER>" returns a list of EPrint IDs associated with the grant number`, repoID)
}

func creatorNameDocument(repoID string) string {
	return fmt.Sprintf(`"/%s/creator-name/<FAMILY_NAME>/<GIVEN_NAME>" scans the name fields associated with creators and returns a list of EPrint ID `, repoID)
}
func editorNameDocument(repoID string) string {
	return fmt.Sprintf(`"/<REPO_ID>/editor-name/<FAMILY_NAME>/<GIVEN_NAME>" scans the family and given name field associated with a editors and returns a list of EPrint ID`, repoID)
}

func contributorNameDocument(repoID string) string {
	return fmt.Sprintf(`"/<REPO_ID>/contributor-name/<FAMILY_NAME>/<GIVEN_NAME>" scans the family and given name field associated with a contributors and returns a list of EPrint ID`, repoID)
}

func advisorNameDocument(repoID string) string {
	return fmt.Sprintf(`"/%s/advisor-name/<FAMILY_NAME>/<GIVEN_NAME>" scans the name fields associated with advisors returns a list of EPrint ID`, repoID)
}
func committeeNameDocument(repoID string) string {
	return fmt.Sprintf(`"/%s/committee-name/<FAMILY_NAME>/<GIVEN_NAME>" scans the family and given name fields associated with committee members and returns a list of EPrint ID`, repoID)
}

func pubmedDocument(repoID string) string {
	return fmt.Sprintf(`"/%s/pubmed/<PUBMED_ID>" returns a list of EPrint IDs associated with the PubMed ID`, repoID)
}

func issnDocument(repoID string) string {
	return fmt.Sprintf(`"/%s/issn/<ISSN>" returns a list of EPrint IDs associated with the ISSN`, repoID)
}

func isbnDocument(repoID string) string {
	return fmt.Sprintf(`"/%s/isbn/<ISSN>" returns a list of EPrint IDs associated with the ISSN`, repoID)
}

func patentNumberDocument(repoID string) string {
	return fmt.Sprintf(`"/%s/patent-number/<PATENT_NUMBER>" returns a list of EPrint IDs associated with the patent number`, repoID)
}

//
// End Point handles (route as defined `/<REPO_ID>/<END-POINT>/<ARGS>`)
//

func repositoriesEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	repositories := []string{}
	for repository, _ := range config.Repositories {
		repositories = append(repositories, repository)
	}
	return packageDocument(w, fmt.Sprintf(`
Available repositories:

- %s
`, strings.Join(repositories, "\n- ")))
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
			return 400, fmt.Errorf("Bad Request, %s", err)
		}
	}
	if (len(args) > 1) && (args[1] != "now") {
		end, err = time.Parse(timestamp, args[1])
		if err != nil {
			return 400, fmt.Errorf("Bad Request, %s", err)
		}
	}
	eprintIDs, err := sqlQueryIDs(repoID, `SELECT eprintid FROM eprint WHERE (lastmod_year >= ? AND lastmod_year <= ?) AND (lastmod_month >= ? AND lastmod_month <= ?) AND (lastmod_day >= ? AND lastmod_day <= ?) AND (lastmod_hour >= ? AND lastmod_hour <= ?) AND (lastmod_minute >= ? AND lastmod_minute <=?)`, start.Year(), end.Year(), start.Month(), end.Month(), start.Day(), end.Day(), start.Hour(), end.Hour(), start.Minute(), end.Minute())
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
			return 400, fmt.Errorf("Bad Request, %s", err)
		}
	}
	if (len(args) > 1) && (args[1] != "now") {
		end, err = time.Parse(timestamp, args[1])
		if err != nil {
			return 400, fmt.Errorf("Bad Request, %s", err)
		}
	}
	eprintIDs, err := sqlQueryIDs(repoID, `SELECT eprintid FROM eprint WHERE (lastmod_year >= ? AND lastmod_year <= ?) AND (lastmod_month >= ? AND lastmod_month <= ?) AND (lastmod_day >= ? AND lastmod_day <= ?) AND (lastmod_hour >= ? AND lastmod_hour <= ?) AND (lastmod_minute >= ? AND lastmod_minute <=?) AND (eprint_status = 'deletion')`, start.Year(), end.Year(), start.Month(), end.Month(), start.Day(), end.Day(), start.Hour(), end.Hour(), start.Minute(), end.Minute())
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
			return 400, fmt.Errorf("Bad Request, %s", err)
		}
	}
	if (len(args) > 1) && (args[1] != "now") {
		dt = expandAproxDate(args[1], false)
		end, err = time.Parse(datestamp, dt)
		if err != nil {
			return 400, fmt.Errorf("Bad Request, %s", err)
		}
	}
	eprintIDs, err := sqlQueryIDs(repoID, `SELECT eprintid FROM eprint WHERE (date_type = "published") and (date_year >= ? AND date_year <= ?) AND (date_month >= ? AND date_month <= ?) AND (date_day >= ? AND date_day <= ?)`, start.Year(), end.Year(), start.Month(), end.Month(), start.Day(), end.Day())
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
	eprintIDs, err := sqlQueryIDs(repoID, `SELECT eprintid FROM eprint_funder_agency_grant_number WHERE funders_grant_number = ?`, args[0])
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
	//FIXME: the API should reject requests that are not application/json or text/plain since that is all we provide.
	if r.Method != "GET" {
		statusCode, err = 405, fmt.Errorf("Method Not Allowed")
		handleError(w, statusCode, err)
	} else {
		switch r.URL.Path {
		case "/favicon.ico":
			statusCode, err = 200, nil
			fmt.Fprintf(w, "")
			//statusCode, err = 404, fmt.Errorf("Not Found")
			//handleError(w, statusCode, err)
		default:
			statusCode, err = routeEndPoints(w, r)
			if err != nil {
				handleError(w, statusCode, err)
			}
		}
	}
	logRequest(r, statusCode, err)
}

func loadConfig(fname string) error {
	config = new(Config)
	config.Repositories = map[string]string{}
	if src, err := ioutil.ReadFile(fname); err != nil {
		return err
	} else {
		// Since we should be OK, unmarshal in into active config
		if err = json.Unmarshal(src, &config); err != nil {
			return fmt.Errorf("Unmarshaling %q failed, %s", fname, err)
		}
		if config.Hostname == "" {
			config.Hostname = "localhost:8484"
		}
	}
	return nil
}

func InitExtendedAPI(settings string) error {
	var err error
	if err = loadConfig(settings); err != nil {
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
	}

	/* NOTE: We need a DB connection to MySQL for each
	   EPrints repository supported by the API
	   for access to MySQL */
	for repoID, dataSourceName := range config.Repositories {
		// Setup DB connection for target repository
		if db, err := sql.Open("mysql", dataSourceName); err != nil {
			return fmt.Errorf("Could not open MySQL conncetion for %s, %s", repoID, err)
		} else {
			//log.Printf("Setting  DB connection to %q", repoID)
			//db.Ping()
			config.Connections[repoID] = db
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
