#!/bin/bash
#
# This script outputs contents from the EPrints 3.3 REST API
#
# Examples:
#     List the eprint URI
#
#           ./script/list.bash eprint/
#
#     List the eprint id number 3
#
#           ./script/list.bash eprint/3.xml
#
#     List subjects
#
#           ./script/list.bash subject/
#
#
if [ "$EPRINT_URL" = "" ]; then
    echo "Environment not configured."
    exit 1
fi

TARGET=""
if [ "$1" != "" ]; then
    TARGET=$1
fi

if [ "$EP_USERNAME" != "" ] && [ "$EP_PASSWORD" != "" ]; then
    curl \
        -X GET \
        -u "$EP_USERNAME:$EP_PASSWORD" \
        $EP_EPRINT_URL/rest/$TARGET
else
    curl \
        -X GET \
        $EP_EPRINT_URL/rest/$TARGET
fi
