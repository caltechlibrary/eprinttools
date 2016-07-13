#!/bin/bash
#
# This script outputs contents from the EPrints 3.3 REST API
#
# Examples:
#     List the eprint URI
#
#           ./script/list.sh eprint/
#
#     List the eprint id number 3
#
#           ./script/list.sh eprint/3.xml
#
#     List subjects
#
#           ./script/list.sh subject/
#
#
if [ "$EPGO_API_URL" = "" ]; then
    echo "Environment not configured."
    exit 1
fi

TARGET=""
if [ "$1" != "" ]; then
    TARGET=$1
fi

if [ "$EPGO_USERNAME" != "" ] && [ "$EPGO_PASSWORD" != "" ]; then
    curl \
        -X GET \
        -u "$EPGO_USERNAME:$EPGO_PASSWORD" \
        $EPGO_API_URL/rest/$TARGET
else
    curl \
        -X GET \
        $EPGO_API_URL/rest/$TARGET
fi
