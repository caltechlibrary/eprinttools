#!/bin/bash

#
# Requirements MySQL Shell and Stephen Dolan's jq for formatting json output.
#

function showEPrintColumns() {
    DB_NAME="${1}"
    TABLE_NAME="${2}"
    mysql --raw --batch --execute \
      "SHOW COLUMNS IN ${TABLE_NAME} IN ${DB_NAME}" |\
         cut  -f 1 | tail -n +2 | sort
}

function showEPrintTables() {
    DB__NAME="${1}"
    mysql --raw --batch --execute \
      "SHOW TABLES IN ${DB_NAME} LIKE 'eprint%'" \
          | cut -f 1 | tail -n +2 | sort
}

#
# Main processing
#
if [ "$1" = "" ]; then
    APP_NAME=$(basename $0)
    cat <<EOF
USAGE 

    ${APP_NAME} EPRINT_DB_NAME

This script will generate a JSON map to each EPrint
table and table's columns. This can be used to determine
what to crosswalk for any specific EPrint repository.

NOTE: It only maps the tables related to the EPrint XML
record.

Requires the MySQL Shell and Stephen Dolan's jq to work.

The MySQL Shell needs to be configure to access the
MySQL database 

EOF
    exit 1
fi

for DB_NAME in $1; do
    echo -n '{'
    TABLES=$(showEPrintTables "${DB_NAME}")
    DELIM=""
    for TABLE_NAME in ${TABLES}; do
        if [ "$DELIM" != "" ]; then
            echo "${DELIM}"
        else
            DELIM=","
            echo ''
        fi
        echo -n "\"${TABLE_NAME}\":"
        echo -n '['
        for COLUMN in $(showEPrintColumns ${DB_NAME} ${TABLE_NAME}); do
            echo -n "\"${COLUMN}\" "
        done | sed -E 's/" "/", "/g'
        echo -n ']'
    done
echo -n '}'
done | jq .
