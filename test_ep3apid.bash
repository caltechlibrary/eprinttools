#!/bin/bash

#
# Test ep3apid support
#

function post_file () {
    FNAME="${1}"
    echo "POSTing ${FNAME}"
    curl -X POST -H 'Content-Type: application/xml' \
         http://localhost:8484/lemurprints/eprint-import \
         --data-binary "@${FNAME}" >> test-responses.txt
}

echo 'Testing ep3apid'
if [[ -f "test-settings.json" && -f "bin/ep3apid" ]]; then
   echo 'Resetting lemurprints database'
   mysql lemurprints < srctest/lemurprints-setup-schema.sql
   echo 'Starting epi3apid and waiting 30 seconds'
   ./bin/ep3apid test-settings.json &
   PID=$!
   # Add content to empty repository database
   sleep 30
   echo 'Adding content'
   let err_count=0
   for FNAME in $(ls -1 srctest/lemurprints-import-api-*.xml); do
       if ! post_file $FNAME; then
          echo "Failed to POST ${FNAME}"
          let err_count++
       fi
   done
   echo "Kill process $PID"
   kill $PID
fi
if [[ "$err_count" != "0" ]]; then
    echo "Error count: $err_count"
else
    echo 'OK, passed cli tests'
fi
