#!/bin/bash

#
# Run some command line tests to confirm working cli
#

if [[ -d "testout" ]]; then
    rm -fR testout
fi
mkdir testout

#
# Jira Issue DR-40
#
echo 'Testing DR-40'
./bin/doi2eprintxml -i testdata/DR-40-test.txt > testout/dr40-eprint.xml
T=$(grep '<date>1942-07</date>' testout/dr40-eprint.xml)
if [[ "${T}" != "" ]]; then
    echo "expected '<data>1942-07</date>', got '${T}'"
    exit 1
fi

#
# Jira Iusse DR-43
#
echo 'Testing DR-43'
./bin/doi2eprintxml -i testdata/DR-43-test.txt > testout/dr43-eprint.xml
T=$(grep '<date>2001-01</date>' testout/dr43-eprint.xml)
if [[ "${T}" != "" ]]; then
    echo "expected '<data>2001-01</date>', got '${T}'"
    exit 1
fi




#
# Jira Issue DR-45
#
echo 'Testing DR-45 (this takes a while)'
./bin/doi2eprintxml -i testdata/DR-45-test.txt > testout/dr45-eprint.xml
if [[ "$?" != "0" ]]; then
    echo ''
    echo 'Testing doi2eprintxml DR-45 issue failed.'
    exit 1
fi

#
# Jira Issue DR-59
#
echo 'Testing DR-59'
./bin/doi2eprintxml -i testdata/DR-59-test.txt > testout/dr-59-eprint.xml
if [[ "$?" != "0" ]]; then
    echo ''
    echo 'Testing doi2eprintxml DR-45 issue failed.'
    exit 1
fi

#
# Jira Issue DR-141
#
echo 'Testing DR-141'
./bin/doi2eprintxml -i testdata/DR-141-test.txt > testout/dr-141-eprint.xml
if [[ "$?" != "0" ]]; then
    echo ''
    echo 'Testing doi2eprintxml DR-141 issue failed.'
    exit 1
fi

#
# Test ep3apid support
#

function post_file () {
    USERID="1" # Using Admin userid 1 for tests.
    FNAME="${1}"
    echo "POSTing ${FNAME}"
    curl -X POST -H 'Content-Type: application/xml' \
         "http://localhost:8484/lemurprints/eprint-import/$USERID" \
         --data-binary "@${FNAME}" >/dev/null
}

echo 'Testing ep3apid'
if [[ -f "test-settings.json" && -f "bin/ep3apid" ]]; then
   echo 'Resetting lemurprints database'
   mysql lemurprints < srctest/lemurprints-setup-schema.sql
   echo 'Starting epi3apid and waiting 20 seconds'
   ./bin/ep3apid test-settings.json &
   PID=$!
   # Add content to empty repository database
   sleep 20
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
    exit 1
fi

echo 'OK, passed cli tests'
