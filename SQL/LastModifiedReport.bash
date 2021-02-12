
TODAY=$(date +%Y-%m-%d)

# Describe script usage
function usage() {
    prog=$(basename $0)
    cat<<EOF

USAGE: ${prog} EPRINT_DATABASE_NAME REL_DAYS

This script generates a CSV file of keys and modified datestamps
from an EPrints database.

EXAMPLE:

    ${prog} caltechauthors -7

This would generate a file 

    caltechauthors-lastmod-${TODAY}.csv

containing keys for records modified or created
in the last seven days.
 
EOF
}

# Get the Last Modified report records for the last week.
function run_last_modified_report() {
    DB_NAME="${1}"
    REL_DAYS="${2}"
    START_DAY=$(date -d "${REL_DAYS}days")
    START=$(date -d "${REL_DAYS}days" +%Y,%m,%d)
    echo "Report period: ${START_DAY} to ${TODAY}"
    mysql $1 --batch --execute "CALL Last_Modified_Report(${START})" |\
          tr '\t' ',' > "${DB_NAME}-lastmod-${TODAY}.csv"
    echo "Completed. $(date)"
}


if [ ! "$1" ] || [ ! "$2" ]; then
    usage
    exit 1
fi
if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    usage
    exit 0
fi
run_last_modified_report "$1"

