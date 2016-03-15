# epgo

epgo is a command line utility utilizing EPrints' REST API to produce alternative
feeds and formats. Currently it supports generating a feed of repository items based
on publication dates.

## Overview

USAGE: epgo [OPTIONS] [EPRINT_URI]

epgo wraps the REST API for E-Prints 3.3 or better. It can return a list of uri,
a JSON view of the XML presentation as well as generates feeds and web pages.

epgo can be configured with following environment variables

+ EPGO_API_URL (required) the URL to your E-Prints installation
+ EPGO_DBNAME   (required) the BoltDB name for exporting, site building, and content retrieval
+ EPGO_SITE_URL (optional) the website URL (might be the same as E-Prints)
+ EPGO_HTDOCS   (optional) the htdocs root for site building
+ EPGO_TEMPLATES (optional) the template directory to use for site building

If EPRINT_URI is provided then an individual EPrint is return as a JSON structure
(e.g. /rest/eprint/34.xml). Otherwise a list of EPrint paths are returned.


| Options | Description |
|---------|---------------------------------------------------------------------|
| -api    | read the contents from the API without saving in the database       |
| -build  | build pages and feeds from database                                 |
| -feed-size int | the number of items contained in a feed like recent articles |
| -export | export EPrints to database                                          |
| -h      | display help info                                                   |
| -p      | pretty print JSON output                                            |
| -published-newest int | list the N newest published items                     |
| -published-oldest int | list the N oldest published items                     |
| -articles-newest int  | list the N newest articles                            |
| -articles-oldest int  | list the N oldest articles                            |
|    -v   | display version info                                                |
