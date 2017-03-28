#!/bin/bash

function findEPrintsByPersonID() {
	person="$1"
	findfile -s json htdocs/CaltechAUTHORS | while read FNAME; do
		grep -l "$person" "htdocs/CaltechAUTHORS/$FNAME"
	done >"${person}.dat"
}

function isMissingORCID() {
	person="$1"
	while read P; do
		EID=$(basename -s ".json" "$P")
		jq -c '.creators[] | {"url":"'"$EPGO_API_URL"'/'"$EID"'/","eprint_record": "'"$EID"'", "id": .id, "orcid": .orcid} | join(", ")' "$P" \
			| sed -E 's/"//g' \
			| grep "$person"
	done <"${person}.dat"
}

#
# Main processing
#

# If $EPGO_API_URL is missing stop
if [ "$EPGO_API_URL" = "" ]; then
	echo "Missing EPGO_API_URL environment setting"
	exit 1
fi
# Make sure we have an id to work with
if [ "$1" = "" ]; then
	echo "USAGE $0 PERSON_ID_STRING"
	exit 1
fi

# Build the data file if needed of records to scan
if [ ! -f "$1.keys" ]; then
	findEPrintsByPersonID "$1"
fi
# Scan for the records and their orcid values
if [ "$2" != "" ]; then
	isMissingORCID "$1" >"$2"
else
	isMissingORCID "$1"
fi
