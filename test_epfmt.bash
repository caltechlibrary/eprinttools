#!/bin/bash

if [[ "${1}" == "" ]]; then
	echo "Missing URL for test repository"
	echo "e.g. https://username:secret@eprint.example.edu"
	exit 1
fi

HAS_JSONRANGE=$(which jsonrange)
if [[ "${HAS_JSONRANGE}" == "" ]]; then
	echo "jsonrange cli needed to run these tests"
	exit 1
fi

if [[ ! -d "testout" ]]; then
	mkdir testout
fi
EP_API="${1}"
if [[ ! -f "testout/sample.keys" ]]; then
    echo "Generating getting keys for sample"
    bin/eputil -json "${EP_API}/rest/eprint/" | jsonrange -values >testout/t.keys

    echo "Generating 5% sample"
    awk 'BEGIN {srand()} !/^$/ { if (rand() <= .05) print $0}' testout/t.keys >testout/sample.keys
else
    echo "Using existing sample.keys"
fi

if [[ ! -s "testout/sample.keys" ]]; then
	echo "Failed to generate a sample of keys from testout/t.keys"
	exit 1
fi
echo -n "Sample size "
wc -l testout/sample.keys

echo "Harvesting records"
while read -r KEY; do
	if [[ "${KEY}" != "" ]]; then
		if [[ ! -f "testout/${KEY}.xml" ]]; then
			if bin/eputil "${EP_API}/rest/eprint/${KEY}.xml" >"testout/${KEY}.xml"; then
				if [[ -s "testout/${KEY}.xml" ]]; then
					echo -n "."
					bin/eputil -json "${EP_API}/rest/eprint/${KEY}.xml" >"testout/${KEY}.json"
				else
					echo "Skipping ${KEY}, empty record"
					rm "testout/${KEY}.xml"
				fi
			else
				echo " Skipping ${KEY}, error code"
				rm "testout/${KEY}.xml"
			fi
		fi
	fi
done <testout/sample.keys
echo ""
echo "Harvest completed."
echo ""

echo "Running epfmt tests on XML sources"
findfile -s .xml testout | grep -E '^[0-9]+\.xml$' | while read -r FNAME; do
	KEY=$(basename "${FNAME}" ".xml")
	if bin/epfmt <"testout/${KEY}.xml" >"testout/${KEY}_t1.xml"; then
		echo -n "."
	else
		echo ""
		echo " Failed on testout/${KEY}.xml to generate testout/${KEY}_t1.xml"
		exit 1
	fi
	if bin/epfmt -json <"testout/${KEY}.xml" >"testout/${KEY}_t2.json"; then
		echo -n "."
	else
		echo ""
		echo " Failed on testout/${KEY}.xml to generate testout/${KEY}_t2.json"
		exit 1
	fi
done

echo ""
echo "Running epfmt tests on JSON sources"
findfile -s .json testout | grep -E '^[0-9]+\.json$' | while read -r FNAME; do
	KEY=$(basename "${FNAME}" ".json")
	if bin/epfmt <"testout/${KEY}.json" >"testout/${KEY}_t3.json"; then
		echo -n "."
	else
		echo ""
		echo " Failed on testout/${KEY}.json to generate testout/${KEY}_t3.json"
		exit 1
	fi
	if bin/epfmt -xml <"testout/${KEY}.json" >"testout/${KEY}_t4.json"; then
		echo -n "."
	else
		echo ""
		echo " Failed on testout/${KEY}.json to generate testout/${KEY}_t4.json"
		exit 1
	fi
done

echo ""
echo "All Done!"
