
USAGE
=====

    ep3apid [OPTIONS] [SETTINGS_FILENAME]

SYNOPSIS
--------

Run an extended EPrints 3.x web API

DETAIL
------

ep3apid can be run from the command line and the will create an
http web service on %!s(MISSING). The web service provides a limitted number of
end points providing eprint ids for content matched in EPrints's MySQL
databases.

The following URL end points are intended to take one unique identifier and map that to one or more EPrint IDs. This can be done because each unique ID  targeted can be identified by querying a single table in EPrints.  In addition the scan can return the complete results since all EPrint IDs are integers and returning all EPrint IDs in any of our repositories is sufficiently small to be returned in a single HTTP request.

### Unique ID to EPrint ID

- '/<REPO_ID>/doi/<DOI>' with the adoption of EPrints "doi" field in the EPrint table it makes sense to have a quick translation of DOI to EPrint id for a given EPrints repository. 
- '/<REPO_ID>/creator-id/<CREATOR_ID>' scans the name creator id field associated with creators and returns a list of EPrint ID 
- '/<REPO_ID>/creator-orcid/<ORCID>' scans the "orcid" field associated with creators and returns a list of EPrint ID 
- '/<REPO_ID>/editor-id/<CREATOR_ID>' scans the name creator id field associated with editors and returns a list of EPrint ID 
- '/<REPO_ID>/contributor-id/<CONTRIBUTOR_ID>' scans the "id" field associated with a contributors and returns a list of EPrint ID 
- '/<REPO_ID>/advisor-id/<ADVISOR_ID>' scans the name advisor id field associated with advisors and returns a list of EPrint ID 
- '/<REPO_ID>/committee-id/<COMMITTEE_ID>' scans the committee id field associated with committee members and returns a list of EPrint ID
- '/<REPO_ID>/group-id/<GROUP_ID>' this scans group ID and returns a list of EPrint IDs associated with the group
- '/<REPO_ID>/funder-id/<FUNDER_ID>' returns a list of EPrint IDs associated with the funder's ROR
- '/<REPO_ID>/grant-number/<GRANT_NUMBER>' returns a list of EPrint IDs associated with the grant number

### Change Events

The follow API end points would facilitate faster updates to our feeds platform as well as allow us to create a separate public view of our EPrint repository content.

- '/<REPO_ID>/updated/<TIMESTAMP>/<TIMESTAMP>' returns a list of EPrint IDs updated starting at the first timestamp (timestamps should have a resolution to the minute, e.g. "YYYY-MM-DD HH:MM:SS") through inclusive of the second timestmap (if the second is omitted the timestamp is assumed to be "now")
- '/<REPO_ID>/deleted/<TIMESTAMP>/<TIMESTAMP>' through the returns a list of EPrint IDs deleted starting at first timestamp through inclusive of the second timestamp, if the second timestamp is omitted it is assumed to be "now"
- '/<REPO_ID>/pubdate/<APROX_DATESTAMP>/<APPOX_DATESTMP>' this query scans the EPrint table for records with publication starts starting with the first approximate date through inclusive of the second approximate date. If the second date is omitted it is assumed to be "today". Approximate dates my be expressed just the year (starting with Jan 1, ending with Dec 31), just the year and month (starting with first day of month ending with the last day) or year, month and day. The end returns zero or more EPrint IDs.

```
  -h	Display this help message
  -help
    	Display this help message
  -license
    	Display software license
  -version
    	Display software version
```



Settings (configuration)
------------------------

To run the web service create a JSON file named settings.ini in the
current directory where you're invoking _ep3apid_ from. The web
service can be started with running

```
    ep3apid
```

or to load "settings.json" from the current work directory.

```
    ep3apid settings.json
```

