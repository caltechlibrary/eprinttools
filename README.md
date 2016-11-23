[![Go Report Card](http://goreportcard.com/badge/caltechlibrary/epgo)](http://goreportcard.com/report/caltechlibrary/epgo)

# epgo

_epgo_ is a command line utility utilizing EPrints' REST API to produce alternative
feeds and formats. Currently it supports generating a feed of repository items based
on publication dates.

## Overview

USAGE: epgo [OPTIONS] [EPRINT_URI|JAVASCRIPT_FILES]

_epgo_ wraps the REST API for E-Prints 3.3 or better. It can return a list of uri,
a JSON view of the XML presentation as well as generates feeds and web pages.

_epgo_ can be configured with following environment variables

+ EPGO_API_URL (required) the URL to your E-Prints installation
+ EPGO_DBNAME   (required) the BoltDB name for exporting, site building, and content retrieval
+ EPGO_SITE_URL (optional) the website URL (might be the same as E-Prints)
+ EPGO_HTDOCS   (optional) the htdocs root for site building
+ EPGO_TEMPLATE_PATH (optional) the template directory to use for site building

If EPRINT_URI is provided then an individual EPrint is return as a JSON structure
(e.g. /rest/eprint/34.xml). Otherwise a list of EPrint paths are returned.


| Options               | Description                                         |
|-----------------------|-----------------------------------------------------|
| -api	                | display EPrint REST API response                    |
| -dbname               | BoltDB name                                         |
| -export int           | export N EPrints to local database, if N is negative export all EPrints |
| -read-api             | read the contents from the API without saving in the database |
| -feed-size int        | the number of items included in generated feeds     |
| -published-newest int | list the N newest published records                 |
| -published-oldest int | list the N oldest published records                 |
| -articles-newest int  | list the N newest articles                          |
| -articles-oldest int  | list the N oldest articles                          |
|                       |                                                     |
| -p                    | pretty print JSON output                            |
|                       |                                                     |
| -h                    |  display help info                                  |
| -l                    |  show license information                           |
| -v                    |  display version info                               |




