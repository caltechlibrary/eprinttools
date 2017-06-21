#!/bin/bash

# Setup the necessary environment variables for using ep
# EP_API_URL - base URL for your eprints repository (e.g. http://lemurprints.local)
export EP_API_URL=http://lemurprints.local
# EP_SITE_URL - the website url (might be the same as your eprints repository)
export EP_SITE_URL=http://lemurprints.local
# EP_DBNAME - the database name used in exporting or building
export EP_DBNAME=$HOME/eprint-data.db
# EP_HTDOCS - the htdocs directory where files are written in building
export EP_HTDOCS=/var/local/www/htdocs
# EP_TEMPLATES - the template directory holding the templates for building
export EP_TEMPLATES=$HOME/src/github.com/caltechlibrary/ep/templates/default

#
# The first time you setup feed you want to export all EPrints:
#$HOME/bin/ep -export -1
#

# Export the N highest published IDs
$HOME/bin/ep -export 1000

# Build the feed from export
$HOME/bin/ep -build

