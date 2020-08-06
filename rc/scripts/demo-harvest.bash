#!/bin/bash

function harvest_sample() {
	COLLECTION="$1"
	SAMPLE_SIZE="$2"
	START=$(pwd)
	# Source our configuration
	if [ -f "etc/${COLLECTION}.bash" ]; then
		. "etc/${COLLECTION}.bash"
	else
		echo "Can't find etc/${COLLECTION}.bash holding configuration"
		exit 1
	fi
	mkdir -p "${COLLECTION}-sample"
	cd "${COLLECTION}-sample"
	echo "Getting Keys via the list of anchors in the HTML Response"
	apiexplorer "$RC_API_URL/rest/eprint/" | grep '\.xml' | cut -d\> -f 3 | sed -E 's/\.xml<\/a//g' | sort -r > eprint-keys.txt
	echo "Getting the first 10 records in the keys collection"
	head -n "${SAMPLE_SIZE}" eprint-keys.txt | while read ID; do
		echo "Getting the XML version of the ${ID}"
		apiexplorer "${RC_API_URL}/rest/eprint/${ID}.xml" >"${ID}.xml"
		echo "Getting the JSON version of the ${ID}"
		apiexplorer -as-json "${RC_API_URL}/rest/eprint/${ID}.xml" >"${ID}.json"
	done

	cd "$START"
}

function ingest_sample() {
	COLLECTION="$1"
	dataset init "$COLLECTION"
	for FNAME in $(findfile -s .json "${COLLECTION}-sample"); do
		dataset -c "${COLLECTION}" -i "${COLLECTION}-sample/$FNAME" create "$(basename "$FNAME" ".json")"
	done
}

if [ "$1" = "" ] || [ "$2" = "" ]; then
	echo "USAGE: $(basename ${0}) COLLECTION_NAME SAMPLE_SIZE"
	exit 1
fi

harvest_sample "$1" "$2"
ingest_sample "$1"
