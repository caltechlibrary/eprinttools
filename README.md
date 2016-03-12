# epgo

An experimental Go wrapper around the E-Prints REST API. It provides for 
exporting eprint documents from the REST API into a BoltDB database.  In
the BoltDB database an index is created for recently published repository items.
A set of feed documents (json, rss and html include files) can then be render
based on the exported materials.


## Overview

 USAGE: epgo [OPTIONS] [EPRINT_URI]

 epgo wraps the REST API for E-Prints 3.3 or better. It can return a list of uri,
 a JSON view of the XML presentation as well as generates feeds and web pages.

 epgo can be configured with following environment variables

 + EPGO_BASE_URL (required) the URL to your E-Prints installation
 + EPGO_DBNAME   (required) the BoltDB name for exporting, site building, and content retrieval
 + EPGO_HTDOCS   (optional) the htdocs root for site building
 + EPGO_TEMPLATES (optional) the template directory to use for site building

 If EPRINT_URI is provided then an individual EPrint is return as
 a JSON structure (e.g. /rest/eprint/34.xml). Otherwise a list of EPrint paths are
 returned.

 OPTIONS

|         | Description |
+---------+-------------------------------------------------------------------+
| -api    | read the contents from the API without saving in the database     |
| -build  | build pages and feeds from database                               |
| -export | export EPrints to database                                        |
| -h      | display help info                                                 |
| -p      | pretty print JSON output                                          |
| -published-newest int | list the N newest published items                   |
| -published-oldest int | list the N oldest published items                   |
|    -v   | display version info                                              |

