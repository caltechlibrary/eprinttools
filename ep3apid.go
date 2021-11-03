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
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
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
	lg     *log.Logger
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
// Package functions take the collected data and package into an HTTP
// response.
//

func packageIntIDs(w http.ResponseWriter, repoID string, values []int, err error) (int, error) {
	if err != nil {
		lg.Printf("ERROR: (%s) query error, %s", repoID, err)
		return 500, fmt.Errorf("Internal Server Error")
	}
	src, err := json.MarshalIndent(values, "", "  ")
	if err != nil {
		lg.Printf("ERROR: marshal error (%q), %s", repoID, err)
		return 500, fmt.Errorf("Internal Server Error")
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", src)
	return 200, nil
}

func packageStringIDs(w http.ResponseWriter, repoID string, values []string, err error) (int, error) {
	if err != nil {
		lg.Printf("ERROR: (%s) query error, %s", repoID, err)
		return 500, fmt.Errorf("Internal Server Error")
	}
	src, err := json.MarshalIndent(values, "", "  ")
	if err != nil {
		lg.Printf("ERROR: marshal error (%q), %s", repoID, err)
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
		lg.Printf("ERROR: (%s) package JSON error, %s", repoID, err)
		return 500, fmt.Errorf("Internal Server Error")
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", src)
	return 200, nil
}

func packageXML(w http.ResponseWriter, repoID string, src []byte, err error) (int, error) {
	if err != nil {
		lg.Printf("ERROR: (%s) package JSON error, %s", repoID, err)
		return 500, fmt.Errorf("Internal Server Error")
	}
	w.Header().Set("Content-Type", "application/xml")
	fmt.Fprintln(w, `<?xml version="1.0" encoding="utf-8"?>`)
	fmt.Fprintf(w, "%s", src)
	return 200, nil
}

// Given a POST Request, read the Request body and
// populate an EPrints structure or return an error.
func unpackageEPrintsPOST(r *http.Request) (*EPrints, error) {
	var (
		eprints *EPrints
		err     error
		src     []byte
	)
	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "application/json":
		if src, err = ioutil.ReadAll(r.Body); err != nil {
			return nil, fmt.Errorf("Failed to read, %s", err)
		}
		err = json.Unmarshal(src, &eprints)
	case "application/xml":
		if src, err = ioutil.ReadAll(r.Body); err != nil {
			return nil, fmt.Errorf("Failed to read, %s", err)
		}
		err = xml.Unmarshal(src, &eprints)
	default:
		return nil, fmt.Errorf("%s not supported", contentType)
	}
	return eprints, err
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

func keysEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) != 0 || strings.HasSuffix(r.URL.Path, "/help") {
		return packageDocument(w, keysDocument(repoID))
	}
	eprintIDs, err := sqlQueryIntIDs(repoID, `SELECT eprintid FROM eprint`)
	return packageIntIDs(w, repoID, eprintIDs, err)
}

func createdEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 || strings.HasSuffix(r.URL.Path, `/help`) {
		return packageDocument(w, createdDocument(repoID))
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
	eprintIDs, err := sqlQueryIntIDs(repoID,
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
	return packageIntIDs(w, repoID, eprintIDs, err)
}

func updatedEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 || strings.HasSuffix(r.URL.Path, `/help`) {
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
	eprintIDs, err := sqlQueryIntIDs(repoID,
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
	return packageIntIDs(w, repoID, eprintIDs, err)
}

func deletedEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 || strings.HasSuffix(r.URL.Path, `/help`) {
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
	eprintIDs, err := sqlQueryIntIDs(repoID,
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
	return packageIntIDs(w, repoID, eprintIDs, err)
}

func pubdateEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 || strings.HasSuffix(r.URL.Path, `/help`) {
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
	eprintIDs, err := sqlQueryIntIDs(repoID,
		`SELECT eprintid FROM eprint
WHERE ((date_type) = "published") AND 
(CONCAT(date_year, "-", 
LPAD(IFNULL(date_month, 1), 2, "0"), "-", 
LPAD(IFNULL(date_day, 1), 2, "0")) >= ?) AND 
(CONCAT(date_year, "-", 
LPAD(IFNULL(date_month, 12), 2, "0"), "-", 
LPAD(IFNULL(date_day, 28), 2, "0")) <= ?)`,
		start.Format(datestamp), end.Format(datestamp))
	return packageIntIDs(w, repoID, eprintIDs, err)
}

func doiEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 || strings.HasSuffix(r.URL.Path, `/help`) {
		return packageDocument(w, creatorIDDocument(repoID))
	}
	if len(args) < 1 {
		return 400, fmt.Errorf("Bad Request")
	}
	var err error
	doi := strings.Join(args, "/")
	//FIXME: Should validate DOI format ...
	eprintIDs, err := sqlQueryIntIDs(repoID, `SELECT eprintid FROM eprint WHERE doi = ?`, doi)
	return packageIntIDs(w, repoID, eprintIDs, err)
}

func creatorIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return packageDocument(w, creatorIDDocument(repoID))
	}
	if len(args) == 0 {
		values, err := sqlQueryStringIDs(repoID, `SELECT creators_id
FROM eprint_creators_id
WHERE creators_id IS NOT NULL
GROUP BY creators_id ORDER BY creators_id`)
		return packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := sqlQueryIntIDs(repoID, `SELECT eprintid
FROM eprint_creators_id WHERE creators_id = ?`, args[0])
	return packageIntIDs(w, repoID, eprintIDs, err)
}

func creatorORCIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return packageDocument(w, creatorORCIDDocument(repoID))
	}
	if len(args) == 0 {
		values, err := sqlQueryStringIDs(repoID, `SELECT creators_orcid
FROM eprint_creators_orcid
WHERE creators_orcid IS NOT NULL
GROUP BY creators_orcid ORDER BY creators_orcid`)
		return packageStringIDs(w, repoID, values, err)
	}
	//FIXME: Should validate ORCID format ...
	eprintIDs, err := sqlQueryIntIDs(repoID, `SELECT eprintid
FROM eprint_creators_orcid WHERE creators_orcid = ?`, args[0])
	return packageIntIDs(w, repoID, eprintIDs, err)
}

func editorIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return packageDocument(w, editorIDDocument(repoID))
	}
	if len(args) == 0 {
		values, err := sqlQueryStringIDs(repoID, `SELECT editors_id
FROM eprint_editors_id
WHERE editors_id IS NOT NULL
GROUP BY editors_id ORDER BY editors_id`)
		return packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := sqlQueryIntIDs(repoID, `SELECT eprintid
FROM eprint_editors_id WHERE editors_id = ?`, args[0])
	return packageIntIDs(w, repoID, eprintIDs, err)
}

func contributorIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return packageDocument(w, contributorIDDocument(repoID))
	}
	if len(args) == 0 {
		values, err := sqlQueryStringIDs(repoID, `SELECT contributors_id
FROM eprint_contributors_id
WHERE contributors_id IS NOT NULL
GROUP BY contributors_id ORDER BY contributors_id`)
		return packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := sqlQueryIntIDs(repoID, `SELECT eprintid FROM eprint_contibutors_id WHERE contibutors_id = ?`, args[0])
	return packageIntIDs(w, repoID, eprintIDs, err)
}

func advisorIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return packageDocument(w, advisorIDDocument(repoID))
	}
	if len(args) == 0 {
		values, err := sqlQueryStringIDs(repoID, `SELECT thesis_advisor_id
FROM eprint_thesis_advisor_id
WHERE thesis_advisor_id IS NOT NULL
GROUP BY thesis_advisor_id ORDER BY thesis_advisor_id`)
		return packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := sqlQueryIntIDs(repoID, `SELECT eprintid FROM eprint_thesis_advisor_id WHERE thesis_advisor_id = ?`, args[0])
	return packageIntIDs(w, repoID, eprintIDs, err)
}

func committeeIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return packageDocument(w, committeeIDDocument(repoID))
	}
	if len(args) == 0 {
		values, err := sqlQueryStringIDs(repoID, `SELECT thesis_committee_id
FROM eprint_thesis_committee_id
WHERE thesis_committee_id IS NOT NULL
GROUP BY thesis_committee_id ORDER BY thesis_committee_id`)
		return packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := sqlQueryIntIDs(repoID, `SELECT eprintid FROM eprint_thesis_committee_id WHERE thesis_committee_id = ?`, args[0])
	return packageIntIDs(w, repoID, eprintIDs, err)
}

func groupIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return packageDocument(w, groupIDDocument(repoID))
	}
	if len(args) == 0 {
		values, err := sqlQueryStringIDs(repoID, `SELECT local_group
FROM eprint_local_group WHERE local_group IS NOT NULL GROUP BY local_group ORDER BY local_group`)
		return packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := sqlQueryIntIDs(repoID, `SELECT eprintid FROM eprint_local_group WHERE local_group = ?`, args[0])
	return packageIntIDs(w, repoID, eprintIDs, err)
}

func funderIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return packageDocument(w, funderIDDocument(repoID))
	}
	if len(args) == 0 {
		values, err := sqlQueryStringIDs(repoID, `SELECT funders_agency FROM eprint_funders_agency WHERE funders_agency IS NOT NULL GROUP BY funders_agency ORDER BY funders_agency`)
		return packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := sqlQueryIntIDs(repoID, `SELECT eprintid FROM eprint_funders_agency WHERE funders_agency = ?`, args[0])
	return packageIntIDs(w, repoID, eprintIDs, err)
}

func grantNumberEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return packageDocument(w, grantNumberDocument(repoID))
	}
	if len(args) == 0 {
		values, err := sqlQueryStringIDs(repoID, `SELECT funders_grant_number FROM eprint_funders_grant_number WHERE funders_grant_number IS NOT NULL GROUP BY funders_grant_number ORDER BY funders_grant_number`)
		return packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := sqlQueryIntIDs(repoID, `SELECT eprintid FROM eprint_funders_grant_number WHERE funders_grant_number = ?`, args[0])
	return packageIntIDs(w, repoID, eprintIDs, err)
}

func creatorNameEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 || strings.HasSuffix(r.URL.Path, `/help`) {
		return packageDocument(w, creatorNameDocument(repoID))
	}
	if len(args) != 2 {
		return 400, fmt.Errorf("Bad Request")
	}
	var err error
	eprintIDs, err := sqlQueryIntIDs(repoID, `SELECT eprintid FROM eprint_creator_name WHERE family_name = ? AND given_name = ?`, args[0], args[1])
	return packageIntIDs(w, repoID, eprintIDs, err)
}

func editorNameEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	//- `/<REPO_ID>/editor-name/<FAMILY_NAME>/<GIVEN_NAME>` scans the family and given name field associated with a editors and returns a list of EPrint ID
	if len(args) == 0 || strings.HasSuffix(r.URL.Path, `/help`) {
		return packageDocument(w, editorNameDocument(repoID))
	}
	if len(args) != 2 {
		return 400, fmt.Errorf("Bad Request")
	}
	var err error
	eprintIDs, err := sqlQueryIntIDs(repoID, `SELECT eprintid FROM eprint_editor_name WHERE family_name = ? AND given_name = ?`, args[0], args[1])
	return packageIntIDs(w, repoID, eprintIDs, err)
}

func contributorNameEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 || strings.HasSuffix(r.URL.Path, `/help`) {
		return packageDocument(w, contributorNameDocument(repoID))
	}
	if len(args) != 2 {
		return 400, fmt.Errorf("Bad Request")
	}
	var err error
	eprintIDs, err := sqlQueryIntIDs(repoID, `SELECT eprintid FROM eprint_contributors_name WHERE family_name = ? AND given_name = ?`, args[0], args[1])
	return packageIntIDs(w, repoID, eprintIDs, err)
}

func advisorNameEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 || strings.HasSuffix(r.URL.Path, `/help`) {
		return packageDocument(w, advisorNameDocument(repoID))
	}
	if len(args) != 2 {
		return 400, fmt.Errorf("Bad Request")
	}
	var err error
	eprintIDs, err := sqlQueryIntIDs(repoID, `SELECT eprintid FROM eprint_thesis_advisor_name WHERE family_name = ? AND given_name = ?`, args[0], args[1])
	return packageIntIDs(w, repoID, eprintIDs, err)
}

func committeeNameEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 || strings.HasSuffix(r.URL.Path, `/help`) {
		return packageDocument(w, committeeNameDocument(repoID))
	}
	if len(args) != 2 {
		return 400, fmt.Errorf("Bad Request")
	}
	var err error
	eprintIDs, err := sqlQueryIntIDs(repoID, `SELECT eprintid FROM eprint_thesis_committee_name WHERE family_name = ? AND given_name = ?`, args[0], args[1])
	return packageIntIDs(w, repoID, eprintIDs, err)
}

func pubmedIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 || strings.HasSuffix(r.URL.Path, `/help`) {
		return packageDocument(w, pubmedIDDocument(repoID))
	}
	if len(args) != 1 {
		return 400, fmt.Errorf("Bad Request")
	}
	var err error
	eprintIDs, err := sqlQueryIntIDs(repoID, `SELECT eprintid FROM eprint WHERE pmid = ?`, args[0])
	return packageIntIDs(w, repoID, eprintIDs, err)
}

func pubmedCentralIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 || strings.HasSuffix(r.URL.Path, `/help`) {
		return packageDocument(w, pubmedCentralIDDocument(repoID))
	}
	if len(args) != 1 {
		return 400, fmt.Errorf("Bad Request")
	}
	var err error
	eprintIDs, err := sqlQueryIntIDs(repoID, `SELECT eprintid FROM eprint WHERE pmc_id = ?`, args[0])
	return packageIntIDs(w, repoID, eprintIDs, err)
}

func issnEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return packageDocument(w, issnDocument(repoID))
	}
	if len(args) == 0 {
		values, err := sqlQueryStringIDs(repoID, `SELECT issn FROM eprint WHERE issn IS NOT NULL GROUP BY issn ORDER BY issn`)
		return packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := sqlQueryIntIDs(repoID, `SELECT eprintid FROM eprint WHERE issn = ?`, args[0])
	return packageIntIDs(w, repoID, eprintIDs, err)
}

func isbnEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return packageDocument(w, isbnDocument(repoID))
	}
	if len(args) == 0 {
		values, err := sqlQueryStringIDs(repoID, `SELECT isbn FROM eprint WHERE isbn IS NOT NULL GROUP BY isbn ORDER BY isbn`)
		return packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := sqlQueryIntIDs(repoID, `SELECT eprintid FROM eprint WHERE isbn = ?`, args[0])
	return packageIntIDs(w, repoID, eprintIDs, err)
}

func patentApplicantEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return packageDocument(w, patentApplicantDocument(repoID))
	}
	if len(args) == 0 {
		values, err := sqlQueryStringIDs(repoID, `SELECT patent_applicant FROM eprint WHERE patent_applicant IS NOT NULL GROUP BY patent_applicant ORDER BY patent_applicant`)
		return packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := sqlQueryIntIDs(repoID, `SELECT eprintid FROM eprint WHERE patent_applicant = ?`, args[0])
	return packageIntIDs(w, repoID, eprintIDs, err)
}

func patentNumberEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return packageDocument(w, patentNumberDocument(repoID))
	}
	if len(args) == 0 {
		values, err := sqlQueryStringIDs(repoID, `SELECT patent_number FROM eprint WHERE patent_number IS NOT NULL GROUP BY patent_number ORDER BY patent_number`)
		return packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := sqlQueryIntIDs(repoID, `SELECT eprintid FROM eprint WHERE patent_number = ?`, args[0])
	return packageIntIDs(w, repoID, eprintIDs, err)
}

func patentClassificationEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return packageDocument(w, patentClassificationDocument(repoID))
	}
	if len(args) == 0 {
		values, err := sqlQueryStringIDs(repoID, `SELECT patent_classification FROM eprint WHERE patent_classification IS NOT NULL GROUP BY patent_classification ORDER BY patent_classification`)
		return packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := sqlQueryIntIDs(repoID, `SELECT eprintid FROM eprint WHERE patent_classification = ?`, args[0])
	return packageIntIDs(w, repoID, eprintIDs, err)
}

func patentAssigneeEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return packageDocument(w, patentAssigneeDocument(repoID))
	}
	if len(args) != 1 {
		values, err := sqlQueryStringIDs(repoID, `SELECT patent_assignee FROM eprint_patent_assignee WHERE patent_assignee IS NOT NULL GROUP BY patent_assignee ORDER BY patent_assignee`)
		return packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := sqlQueryIntIDs(repoID, `SELECT eprintid FROM eprint_patent_assignee WHERE patent_assignee = ?`, args[0])
	return packageIntIDs(w, repoID, eprintIDs, err)
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
	db, _ := config.Connections[repoID]
	eprint, err := CrosswalkSQLToEPrint(db, repoID, baseURL, eprintID)
	if err != nil {
		lg.Printf("CrosswalkSQLToEPrint Error: %s\n", err)
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
// EPrint End Point is an experimental read end point provided
// in the extended EPrint API.  It reads EPrint data structures
// based on SQL calls to the MySQL database for a given EPrints
// repository. It accepts two content types - "application/xml"
// which returns EPrints XML or "application/json" which returns
// a JSON version of the EPrints XML.
//
func eprintEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 || repoID == "" || strings.HasSuffix(r.URL.Path, "/help") {
		return packageDocument(w, eprintReadWriteDocument(repoID))
	}
	contentType := r.Header.Get("Content-Type")
	dataSource, ok := config.Repositories[repoID]
	if ok == false {
		lg.Printf("Data Source not found for %q", repoID)
		return 404, fmt.Errorf("not found")
	}
	if r.Method != "GET" {
		return 405, fmt.Errorf("method not allowed")
	}
	if len(args) != 1 {
		return 400, fmt.Errorf("bad request, missing eprint id in path")
	}
	eprintID, err := strconv.Atoi(args[0])
	if err != nil {
		return 400, fmt.Errorf("bad request, eprint id (%s) %q not valid", repoID, eprintID)
	}
	db, _ := config.Connections[repoID]
	eprint, err := CrosswalkSQLToEPrint(db, repoID, dataSource.BaseURL, eprintID)
	if err != nil {
		return 404, fmt.Errorf("not found, %s", err)

	}
	switch contentType {
	case "application/json":
		src, err := json.MarshalIndent(eprint, "", "    ")
		return packageJSON(w, repoID, src, err)
	case "application/xml":
		eprints := NewEPrints()
		eprints.XMLNS = `http://eprints.org/ep2/data/2.0`
		eprints.Append(eprint)
		src, err := xml.MarshalIndent(eprints, "", "  ")
		return packageXML(w, repoID, src, err)
	case "":
		eprints := NewEPrints()
		eprints.XMLNS = `http://eprints.org/ep2/data/2.0`
		eprints.Append(eprint)
		src, err := xml.MarshalIndent(eprints, "", "  ")
		return packageXML(w, repoID, src, err)
	default:
		return 415, fmt.Errorf("unsupported media type, %q", contentType)
	}
}

//
// EPrint Import End Point is an experimental write end point provided
// in the extended EPrint API.  It Accepts EPrints XML or the JSON
// expression of the EPrint XML and creates new EPrints medata records
//
// This end requires a POST method.  It accepts POST content encoded
// as either "application/json" or "application/xml".
//
// NOTE: this end point has to be enabled in the settings.json file
// defining the individual repository support. "write" needs to be
// set to true.
func eprintImportEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if repoID == "" {
		repoID = `{REPO_ID}`
		return packageDocument(w, eprintReadWriteDocument(repoID))
	}
	if r.Method == "GET" || strings.HasSuffix(r.URL.Path, "/help") {
		return packageDocument(w, eprintReadWriteDocument(repoID))
	}
	writeAccess := false
	dataSource, ok := config.Repositories[repoID]
	if ok == true {
		writeAccess = dataSource.Write
	} else {
		lg.Printf("Data Source not found for %q", repoID)
		return 404, fmt.Errorf("not found")
	}
	if r.Method != "POST" || writeAccess == false {
		return 405, fmt.Errorf("method not allowed")
	}
	// Check to see if we have application/xml or application/json
	// Get data from post
	eprints, err := unpackageEPrintsPOST(r)
	if err != nil {
		return 400, fmt.Errorf("bad request, POST failed (%s), %s", repoID, err)
	}
	ids := []int{}
	db, _ := config.Connections[repoID]
	ids, err = ImportEPrints(db, repoID, dataSource, eprints)
	if err != nil {
		return 400, fmt.Errorf("bad request, create EPrint failed, %s", err)
	}
	return packageIntIDs(w, repoID, ids, err)
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
		lg.Printf("%s %s RemoteAddr: %s UserAgent: %s Query: %+v Response: %d %s", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), q, status, errStr)
	} else {
		lg.Printf("%s %s RemoteAddr: %s UserAgent: %s Response: %d %s", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), status, errStr)
	}
}

func handleError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "text/plain")
	http.Error(w, fmt.Sprintf(`%s`, err), statusCode)
	//fmt.Fprintf(w, `ERROR: %d %s`, statusCode, err)
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

// hasColumn takes a table map (key = table name, value is
// array of column names) and return true if table and column found.
func hasColumn(tableMap map[string][]string, tableName string, columnName string) bool {
	if columns, ok := tableMap[tableName]; ok {
		for _, value := range columns {
			if value == columnName {
				return true
			}
		}
	}
	return false
}

// Shutdown shutdowns the EPrints extended API web service started with
// RunExtendedAPI.
func Shutdown(appName string, sigName string) int {
	exitCode := 0
	pid := os.Getpid()
	lg.Printf(`Received signal %s`, sigName)
	lg.Printf(`Closing database connections %s pid: %d`, appName, pid)
	if err := CloseConnections(config); err != nil {
		exitCode = 1
	}
	lg.Printf(`Shutdown completed %s pid: %d exit code: %d `, appName, pid, exitCode)
	return exitCode
}

// Reload performs a Shutdown and an init after re-reading
// in the settings.json file.
func Reload(appName string, sigName string, settings string) error {
	exitCode := Shutdown(appName, sigName)
	if exitCode != 0 {
		return fmt.Errorf("Reload failed, could not shutdown the current processes")
	}
	lg.Printf("Restarting %s using %s", appName, settings)
	return InitExtendedAPI(settings)
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
	/* Setup logging */
	if config.Logfile == `` {
		lg = log.Default()
	} else {
		// Append or create a new log file
		lp, err := os.OpenFile(config.Logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return err
		}
		lg = log.New(lp, ``, log.LstdFlags)
		defer func() {
			lp.Close()
			lg = log.Default()
		}()
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
	// This is a map standard endpoints and end point handlers.
	// NOTE: Eventually this should evolving into a registration
	// style design pattern based on unique fields in repositories.
	routes := map[string]func(http.ResponseWriter, *http.Request, string, []string) (int, error){
		// Standard fields
		"keys":             keysEndPoint,
		"created":          createdEndPoint,
		"updated":          updatedEndPoint,
		"deleted":          deletedEndPoint,
		"pubdate":          pubdateEndPoint,
		"doi":              doiEndPoint,
		"record":           recordEndPoint,
		"eprint":           eprintEndPoint,
		"eprint-import":    eprintImportEndPoint,
		"creator-id":       creatorIDEndPoint,
		"creator-orcid":    creatorORCIDEndPoint,
		"editor-id":        editorIDEndPoint,
		"contributor-id":   contributorIDEndPoint,
		"advisor-id":       advisorIDEndPoint,
		"committee-id":     committeeIDEndPoint,
		"group-id":         groupIDEndPoint,
		"funder-id":        funderIDEndPoint,
		"grant-number":     grantNumberEndPoint,
		"creator-name":     creatorNameEndPoint,
		"editor-name":      editorNameEndPoint,
		"contributor-name": contributorNameEndPoint,
		"advisor-name":     advisorNameEndPoint,
		"commitee-name":    committeeNameEndPoint,
		"issn":             issnEndPoint,
		"isbn":             isbnEndPoint,
	}

	/* NOTE: We need a DB connection to MySQL for each
	   EPrints repository supported by the API
	   for access to MySQL */
	if err := OpenConnections(config); err != nil {
		return err
	}

	for repoID, dataSource := range config.Repositories {
		// Add routes (end points) for the target repository
		for route, fn := range routes {
			if config.Routes[repoID] == nil {
				config.Routes[repoID] = map[string]func(http.ResponseWriter, *http.Request, string, []string) (int, error){}
			}
			config.Routes[repoID][route] = fn
		}
		// NOTE: make sure each end point is supported by repository
		// e.g. CaltechTHESIS doens't have "patent_number",
		// "patent_classification", "parent_assignee", "pmc_id",
		// or "pmid".
		if hasColumn(dataSource.TableMap, "eprint", "pmc_id") {
			config.Routes[repoID]["pmcid"] = pubmedCentralIDEndPoint
		}
		if hasColumn(dataSource.TableMap, "eprint", "pmid") {
			config.Routes[repoID]["pmid"] = pubmedIDEndPoint
		}
		if hasColumn(dataSource.TableMap, "eprint", "patent_applicant") {
			config.Routes[repoID]["patent-applicant"] = patentApplicantEndPoint
		}
		if hasColumn(dataSource.TableMap, "eprint", "patent_number") {
			config.Routes[repoID]["patent-number"] = patentNumberEndPoint
		}
		if hasColumn(dataSource.TableMap, "eprint", "patent_classification") {
			config.Routes[repoID]["patent-classification"] = patentClassificationEndPoint
		}
		if hasColumn(dataSource.TableMap, "eprint_patent_assignee", "patent_assignee") {
			config.Routes[repoID]["patent-assignee"] = patentAssigneeEndPoint
		}
	}
	return nil
}

func RunExtendedAPI(appName string, settings string) error {
	/* Setup web server */
	lg.Printf(`
%s %s

Using configuration %s

Process id: %d

EPrints 3.3.x Extended API

Listening on http://%s

Press ctl-c to terminate.
`, appName, Version, settings, os.Getpid(), config.Hostname)

	/* Listen for Ctr-C */
	processControl := make(chan os.Signal, 1)
	signal.Notify(processControl, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
	go func() {
		sig := <-processControl
		switch sig {
		case syscall.SIGINT:
			os.Exit(Shutdown(appName, sig.String()))
		case syscall.SIGTERM:
			os.Exit(Shutdown(appName, sig.String()))
		case syscall.SIGHUP:
			if err := Reload(appName, sig.String(), settings); err != nil {
				lg.Println(err)
				os.Exit(1)
			}
		default:
			os.Exit(Shutdown(appName, sig.String()))
		}
	}()

	http.HandleFunc("/", api)
	return http.ListenAndServe(config.Hostname, nil)
}
