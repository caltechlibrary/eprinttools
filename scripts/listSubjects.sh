#!/bin/bash

if [ "$EPGO_URL" = "" ] || [ "$EPGO_USERNAME" = "" ]; then
    echo "Environment not configured."
    exit 1
fi

TARGET=""
if [ "$1" != "" ]; then
    TARGET=$1
fi

curl -u $EPGO_USERNAME:$EPGO_PASSWORD -X GET $EPGO_URL/rest/subject/$TARGET
