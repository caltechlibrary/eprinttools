#!/bin/bash

function usage() {
    APP_NAME=$(basename "$0")
    cat <<EOT
%${APP_NAME}(1) user manual
% R. S. Doiel
% 2022-10-20

# NAME

${APP_NAME}

# SYNOPSIS

${APP_NAME} OLD_ISSN NEW_ISSN NEW_JOURNAL_TITLE

# DESCRIPTION

${APP_NAME} will take the OLD_ISSN, NEW_ISSN and NEW_JOURNAL_TITLE
and generate a SQL statement suitable for run against an
Caltech Library EPrints repository database.

# EXAMPLES

From DR-456, reassign ISSN and Journal title for "Geochemistry,
Gephysics, Geosystems".

    ${APP_NAME} "1525-2027" "1525-2027" "Geochemistry, Geophysics, Geosystems"

EOT
}

#
# Main Processing
#
for ARG in "$@"; do
    case $ARG in
        -h|-help|--help)
        usage
        exit 0
        ;;
    esac
done

if [ "$#" != "3" ]; then
    echo "expected, '${APP_NAME} OLD_ISSN NEW_ISSN NEW_JOURNAL_TITLE', aborting"
    exit 1
fi

OLD_ISSN="$1"
NEW_ISSN="$2"
NEW_JOURNAL_TITLE="$3"

printf 'UPDATE eprint SET issn = "%s", publication = "%s" WHERE issn LIKE "%s";\n' "${NEW_ISSN}" "${NEW_JOURNAL_TITLE}" "${OLD_ISSN}"

