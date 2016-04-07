#!/bin/bash

# Setup the necessary environment variables for using epgo
# EPGO_API_URL - base URL for your eprints repository (e.g. http://lemurprints.local)
export EPGO_API_URL=http://lemurprints.local
# EPGO_SITE_URL - the website url (might be the same as your eprints repository)
export EPGO_SITE_URL=http://lemurprints.local
# EPGO_DBNAME - the database name used in exporting or building
export EPGO_DBNAME=$HOME/eprint-data.db
# EPGO_HTDOCS - the htdocs directory where files are written in building
export EPGO_HTDOCS=/var/local/www/htdocs
# EPGO_TEMPLATES - the template directory holding the templates for building
export EPGO_TEMPLATES=$HOME/src/github.com/caltechlibrary/epgo/templates/default

#
# The first time you setup feed you want to export all EPrints:
#$HOME/bin/epgo -export -1
#

# Export the N highest published IDs
$HOME/bin/epgo -export 1000

# Build the feed from export
$HOME/bin/epgo -build

