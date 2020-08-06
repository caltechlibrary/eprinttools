#!/bin/bash

START=$(pwd)
# Source our configuration
if [ -f etc/eprint-api.bash ]; then
    . etc/eprint-api.bash
else
        echo "Can't find etc/eprint-api.bash holding configuration"
        exit 1
fi
mkdir -p demotest
cd demotest
echo "Getting Keys via the list of anchors in the HTML Response"
apiexplorer "$RC_API_URL/rest/eprint/" | grep '\.xml' | cut -d\> -f 3 | sed -E 's/\.xml<\/a//g' > eprint-keys.txt
echo "Getting the first 10 records in the keys collection"
head -n 10 eprint-keys.txt | while read ID; do
    echo "Getting the XML version of the ${ID}"
    apiexplorer "${RC_API_URL}/rest/eprint/${ID}.xml" > "${ID}.xml"
    echo "Getting the JSON version of the ${ID}"
    apiexplorer -as-json "${RC_API_URL}/rest/eprint/${ID}.xml" > "${ID}.json"
done
cd "$START"
