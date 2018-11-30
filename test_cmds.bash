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
echo 'Testing DR-40 (this takes a file)'
./bin/doi2eprintxml -i testdata/DR-40-test.txt > testout/dr40-eprint.xml
T=$(grep '<date>2005-03-19</date>' testout/dr40-eprint.xml)
if [[ "${T}" != "" ]]; then
    echo "expected '', got '${T}'"
    exit 1
fi
exit 0 # DEBUG


#
# Jira Issue DR-45
#
echo 'Testing DR-45 (this takes a file)'
./bin/doi2eprintxml -i testdata/DR-45-test.txt > testout/dr45-eprint.xml
if [[ "$?" != "0" ]]; then
    echo ''
    echo 'Testing doi2eprintxml DR-45 issue failed.'
    exit 1
fi


echo 'OK, passed cli tests'
