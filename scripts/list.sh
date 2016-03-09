#!/bin/bash

if [ "$EPGO_URL" = "" ] || [ "$EPGO_USERNAME" = "" ]; then
    echo "Environment not configured."
    exit 1
fi

TARGET=""
if [ "$1" != "" ]; then
    TARGET=$1
fi

curl \
    -X GET \
    -u "$EPGO_USERNAME:$EPGO_PASSWORD" \
    -H "Accept: application/xml" \
    $EPGO_URL/rest/$TARGET
