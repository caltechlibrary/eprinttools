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
package main

/**
 * ep3apid.go implements an HTTP/HTTPS extended web API as a standalone
 * service.
 */

import (
	"flag"
	"fmt"
	"os"
	"path"

	// Caltech Library Modules
	"github.com/caltechlibrary/eprinttools"
)

var (
	description = `
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
- '/{REPO_ID}/creator-name' returns a list of creator names (family, given) in repository
- '/{REPO_ID}/creator-name/{FAMILY}/{GIVEN}' returns a list of EPrint ID for the given creator using their family and given names
- '/{REPO_ID}/creator-orcid' return a list of "orcid" associated with creators in repository
- '/{REPO_ID}/creator-orcid/{ORCID}' scans the "orcid" field associated with creators and returns a list of EPrint ID
- '/{REPO_ID}/editor-id' returns a list of editor ids available in the EPrints repository
- '/{REPO_ID}/editor-id/{CREATOR_ID}' scans the name creator id field associated with editors and returns a list of EPrint ID
- '/{REPO_ID}/editor-name' returns a list of editor names (family, given) in repository
- '/{REPO_ID}/editor-name/{FAMILY}/{GIVEN}' returns a list of EPrint ID for the given editor using their family and given names
- '/{REPO_ID}/contributor-id' returns a list of contributor ids available in the eprints repository
- '/{REPO_ID}/contributor-id/{CONTRIBUTOR_ID}' scans the "id" field associated with a contributors and returns a list of EPrint ID
- '/{REPO_ID}/contributor-name' returns a list of contributor names (family, given) in repository
- '/{REPO_ID}/contributor-name/{FAMILY}/{GIVEN}' returns a list of EPrint ID for the given contributor using their family and given names
- '/{REPO_ID}/advisor-id' returns a list of advisor ids in the eprints repository
- '/{REPO_ID}/advisor-id/{ADVISOR_ID}' scans the name advisor id field associated with advisors and returns a list of EPrint ID
- '/{REPO_ID}/advisor-name' returns a list of advisor names (family, given) in repository
- '/{REPO_ID}/advisor-name/{FAMILY}/{GIVEN}' returns a list of EPrint ID for the given advisor using their family and given names
- '/{REPO_ID}/committee-id' returns a list of committee id in EPrints repository
- '/{REPO_ID}/committee-id/{COMMITTEE_ID}' scans the committee id field associated with committee members and returns a list of EPrint ID
- '/{REPO_ID}/committee-name' returns a list of committee members names (family, given) in repository
- '/{REPO_ID}/committee-name/{FAMILY}/{GIVEN}' returns a list of EPrint ID for the given committee member using their family and given names
- '/{REPO_ID}/corp-creator-id' returns a list of corp creator ids in the eprints repository
- '/{REPO_ID}/corp-creator-id/{CORP_CREATOR_ID}' returns the list of eprint id for the corporate creator id
- '/{REPO_ID}/corp-creator-uri' returns a list of corp creator uri in the eprints repository
- '/{REPO_ID}/corp-creator-uri/{CORP_CREATOR_URI}' returns the list of eprint id for the corporate creator's URI
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
- '/{REPO_ID}/year' - return a descending list of years containing record with a date type of "published".
- '/{REPO_ID}/year/{YEAR}' - return a list of eprintid for a given year contaning date type of "published".


Change Events
-------------

The follow API end points would facilitate faster updates to our feeds platform as well as allow us to create a separate public view of our EPrint repository content.

- '/{REPO_ID}/keys' returns complete list of EPrint ID in the repository
- '/{REPO_ID}/updated/{TIMESTAMP}/{TIMESTAMP}' returns a list of EPrint IDs updated starting at the first timestamp (timestamps should have a resolution to the minute, e.g. "YYYY-MM-DD HH:MM:SS") through inclusive of the second timestmap (if the second is omitted the timestamp is assumed to be "now")
- '/{REPO_ID}/deleted/{TIMESTAMP}/{TIMESTAMP}' through the returns a list of EPrint IDs deleted starting at first timestamp through inclusive of the second timestamp, if the second timestamp is omitted it is assumed to be "now"
- '/{REPO_ID}/pubdate/{APROX_DATESTAMP}/{APPOX_DATESTMP}' this query scans the EPrint table for records with publication starts starting with the first approximate date through inclusive of the second approximate date. If the second date is omitted it is assumed to be "today". Approximate dates my be expressed just the year (starting with Jan 1, ending with Dec 31), just the year and month (starting with first day of month ending with the last day) or year, month and day. The end returns zero or more EPrint IDs.

Read/Write API
--------------

As of __ep3apid__ version 1.0.3 a new set of end points exists for reading (retreiving EPrints XML) and writing (metadata import) of EPrints XML.  The extended API only supports working with EPrints metadata not directly with the documents or files associated with individual records.

The metadata import functionality is enabled per repository. It only supports importing records at this time.  Importing an EPrint XML document, which could containing multiple EPrint metadata records, is implemented purely using SQL statements and not the EPrints Perl API. This allows you (with the right MySQL configuration) to run the extended API on a different server without resorting to Perl.

- '/{REPO_ID}/eprint/{EPRINT_ID}' method GET with a content type of "application/json" (JSON of EPrint XML) or "application/xml" for EPrint XML
- '/{REPO_ID}/eprint-import' POST accepts EPrints XML with content type of "application/xml" or JSON of EPrints XML with content type "application/json". To enable this feature add the attribute '"write": true' to the repositories setting in settins.json.


settings.json (configuration)
-----------------------------

The JSON settings.json file should look something like "REPO_ID" would
be the name used in the __ep3apid__ to access a specific repsitory. The
"dsn" value should be replaced with an appropriate data source name to
access the MySQL database for the repository you're supporting. You can have many repositories configured in a single __ep3apid__ instance.

    {
        "repositories": {
            "REPO_ID": {
                "dsn": "DB_USER:SECRET@/DB_NAME",
                "base_url": "URL_TO_EPRINT_REPOSITORY",
                "write": false,
                "default_collection": "REPO_ID",
                "default_official_url": "PERMA_LINK_URL",
                "default_rights": "RIGHTS_STATEMENT_GOES_HERE",
				"default_refereed": "TRUE",
				"default_status": "inbox"
            },
            ... /* Additional repositories configured here */ ...
        }
    }

NOTE: The "default_collection", "default_official_url", "default_rights", "default_refereed", "default_status" are option configurations in the 'settings.json' file.


Options
-------

  -h	Display this help message
  -help
    	Display this help message
  -license
    	Display software license
  -version
    	Display software version

`

	examples = `

Example
-------

To run the web service create a JSON file named settings.ini in the
current directory where you're invoking _{app_name}_ from. The web
service can be started with running

    {app_name}

or to load "settings.json" from the current work directory.

    {app_name} settings.json

`

	license = `
{app_name} {version}

Copyright (c) 2021, Caltech
All rights not granted herein are expressly reserved by Caltech.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
`

	showHelp    bool
	showVersion bool
	showLicense bool

	debugLogs bool

	logFile string
)

func checkConfig(cfg *eprinttools.Config) {
	ok := true
	if cfg.Hostname == "" {
		fmt.Fprintf(os.Stderr, `"hostname" not set`)
		ok = false
	}
	if len(cfg.Repositories) == 0 {
		fmt.Fprintf(os.Stderr, `"repositories" not set`)
	}
	if !ok {
		os.Exit(1)
	}
}

func main() {
	appName := path.Base(os.Args[0])
	/* Process command line options */
	flagSet := flag.NewFlagSet(appName, flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "help", false, "Display this help message")
	flagSet.BoolVar(&showVersion, "version", false, "Display software version")
	flagSet.BoolVar(&showLicense, "license", false, "Display software license")

	flagSet.Parse(os.Args[1:])
	args := flagSet.Args()

	if showHelp {
		eprinttools.DisplayUsage(os.Stdout, appName, flagSet, description, examples)
		os.Exit(0)
	}
	if showVersion {
		eprinttools.DisplayVersion(os.Stdout, appName)
		os.Exit(0)
	}
	if showLicense {
		eprinttools.DisplayLicense(os.Stdout, appName)
		os.Exit(0)
	}

	/* Looking settings.json */
	settings := "settings.json"
	if len(args) > 0 {
		settings = args[0]
	}

	/* Initialize Extended API web service */
	api := new(eprinttools.EP3API)
	if err := api.InitExtendedAPI(settings); err != nil {
		fmt.Fprintf(os.Stderr, "InitExtendedAPI(%q) failed, %s\n", settings, err)
		os.Exit(1)
	}
	/* Run Extended API web service */
	if err := api.RunExtendedAPI(appName, settings); err != nil {
		fmt.Fprintf(os.Stderr, "RunExtendedAPI(%q) failed, %s\n", appName, err)
		os.Exit(1)
	}
}
