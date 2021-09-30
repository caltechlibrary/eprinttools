/**
 * ep3apid.go implements an HTTP/HTTPS extended web API as a standalone
 * service.
 */
package main

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

    {app_name} [OPTIONS] [SETTINGS_FILENAME]

SYNOPSIS

Run an extended EPrints 3.x web API

DETAIL

{app_name} can be run from the command line and the will create an
http web service on %s. The web service provides a limitted number of
end points providing eprint ids for content matched in EPrints's MySQL
databases.

The following URL end points are intended to take one unique identifier and map that to one or more EPrint IDs. This can be done because each unique ID  targeted can be identified by querying a single table in EPrints.  In addition the scan can return the complete results since all EPrint IDs are integers and returning all EPrint IDs in any of our repositories is sufficiently small to be returned in a single HTTP request.

Unique ID to EPrint ID
----------------------

- "/<REPO_ID>/doi/<DOI>" with the adoption of EPrints "doi" field in the EPrint table it makes sense to have a quick translation of DOI to EPrint id for a given EPrints repository. 
- "/<REPO_ID>/creator-id/<CREATOR_ID>" scans the name creator id field associated with creators and returns a list of EPrint ID 
- "/<REPO_ID>/creator-orcid/<ORCID>" scans the "orcid" field associated with creators and returns a list of EPrint ID 
- "/<REPO_ID>/editor-id/<CREATOR_ID>" scans the name creator id field associated with editors and returns a list of EPrint ID 
- "/<REPO_ID>/contributor-id/<CONTRIBUTOR_ID>" scans the "id" field associated with a contributors and returns a list of EPrint ID 
- "/<REPO_ID>/advisor-id/<ADVISOR_ID>" scans the name advisor id field associated with advisors and returns a list of EPrint ID 
- "/<REPO_ID>/committee-id/<COMMITTEE_ID>" scans the committee id field associated with committee members and returns a list of EPrint ID
- "/<REPO_ID>/group-id/<GROUP_ID>" this scans group ID and returns a list of EPrint IDs associated with the group
- "/<REPO_ID>/funder-id/<FUNDER_ID>" returns a list of EPrint IDs associated with the funder's ROR
- "/<REPO_ID>/grant-number/<GRANT_NUMBER>" returns a list of EPrint IDs associated with the grant number

Change Events
-------------

The follow API end points would facilitate faster updates to our feeds platform as well as allow us to create a separate public view of our EPrint repository content.

- "/<REPO_ID>/updated/<TIMESTAMP>/<TIMESTAMP>" returns a list of EPrint IDs updated starting at the first timestamp (timestamps should have a resolution to the minute, e.g. "YYYY-MM-DD HH:MM:SS") through inclusive of the second timestmap (if the second is omitted the timestamp is assumed to be "now")
- "/<REPO_ID>/deleted/<TIMESTAMP>/<TIMESTAMP>" through the returns a list of EPrint IDs deleted starting at first timestamp through inclusive of the second timestamp, if the second timestamp is omitted it is assumed to be "now"
- "/<REPO_ID>/pubdate/<APROX_DATESTAMP>/<APPOX_DATESTMP>" this query scans the EPrint table for records with publication starts starting with the first approximate date through inclusive of the second approximate date. If the second date is omitted it is assumed to be "today". Approximate dates my be expressed just the year (starting with Jan 1, ending with Dec 31), just the year and month (starting with first day of month ending with the last day) or year, month and day. The end returns zero or more EPrint IDs.

`

	examples = `
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
	flag.BoolVar(&showHelp, "h", false, "Display this help message")
	flag.BoolVar(&showHelp, "help", false, "Display this help message")
	flag.BoolVar(&showVersion, "version", false, "Display software version")
	flag.BoolVar(&showLicense, "license", false, "Display software license")

	flagSet.Parse(os.Args[1:])
	args := flagSet.Args()

	if showHelp {
		eprinttools.DisplayUsage(os.Stdout, appName, flagSet, description, examples, license)
		os.Exit(0)
	}
	if showVersion {
		eprinttools.DisplayVersion(os.Stdout, appName)
		os.Exit(0)
	}
	if showLicense {
		eprinttools.DisplayLicense(os.Stdout, appName, license)
		os.Exit(0)
	}

	/* Looking settings.json */
	settings := "settings.json"
	if len(args) > 0 {
		settings = args[0]
	}

	/* Initialize Extended API web service */
	if err := eprinttools.InitExtendedAPI(settings); err != nil {
		fmt.Fprintf(os.Stderr, "InitExtendedAPI(%q) Failed, %s", settings, err)
		os.Exit(1)
	}
	/* Run Extended API web service */
	if err := eprinttools.RunExtendedAPI(appName); err != nil {
		fmt.Fprintf(os.Stderr, "RunExtendedAPI(%q) failed, %s", appName, err)
		os.Exit(1)
	}
}
