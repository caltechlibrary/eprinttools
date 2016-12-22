#!/bin/bash

function findEPrintsByPersonID() {
    person=$1
    findfile -s json htdocs/CaltechAUTHORS | while read FNAME; do 
      grep -l "$person" htdocs/CaltechAUTHORS/$FNAME; 
    done > $person.dat
}

function isMissingORCID() {
    person=$1
    cat $person.dat | while read P; do 
        EID=$(basename -s ".json" $P)
        jq -c '.creators[] | {"eprint": "'$EID'", "id": .id, "orcid": .orcid}' $P | grep $person 
    done | while read DATA; do
        jq -c '[.eprint, .id, .orcid] | join(", ")' | sed -E 's/"//g'
    done
}

# Make sure we have an id to work with
if [ "$1" = "" ]; then
    echo "USAGE $0 PERSON_ID_STRING"
    exit 1
fi

# Build the data file if needed of records to scan
if [ ! -f "$1.dat" ]; then
    findEPrintsByPersonID $1
fi
# Scan for the records and their orcid values
if [ "$2" != "" ]; then
    isMissingORCID $1 > $2
else
    isMissingORCID $1
fi

