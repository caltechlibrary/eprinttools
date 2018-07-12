#!/bin/bash

if [ "$#" != "2" ]; then
	cat <<EOT

  USAGE: $(basename $0) REPO_URL DATASET

  If the EPrints repository is installed at http://authors.example.org then

      $(basename $0) http://authors.example.org author_id_to_orcid

EOT
	exit 1
fi

T=$(which eputil)
if [ "$T" = "" ]; then
	echo "Can't find eputil, check your path or install eprinttools from https://github.com/caltechlibrary/eprinttools/releases/latest"
	exit 1
fi

REPO="$1"
DATASET="$2"
for EPRINT_ID in $(eputil -url "${REPO}/rest/eprint/" -ids | jsonrange -values); do
	echo "Fetching eprint $EPRINT_ID"
	NO_CREATORS=$(eputil -quiet -url "${REPO}/rest/eprint/${EPRINT_ID}/creators/size.txt")
	if [ "$?" != "" ] && [ "$NO_CREATORS" != "0" ] && [ "$NO_CREATORS" != "" ]; then
		for I in $(range 1 "$NO_CREATORS"); do
			# Step 1, get the creator id
			CREATOR_ID=$(eputil -url "${REPO}/rest/eprint/${EPRINT_ID}/creators/${I}/id.txt")
			if [ "$?" = "0" ] && [ "${CREATOR_ID}" != "" ]; then
				# Step 2, have we harvested this creator before?
				REC=$(dataset -quiet -c "${DATASET}" read "${CREATOR_ID}")
				if [ "$?" != "0" ] && [ "$REC" = "" ]; then
					# Step 3, see if we have an orcid
					ORCID=$(eputil -url "${REPO}/rest/eprint/${EPRINT_ID}/creators/${I}/orcid.txt")
					if [ "$?" = "0" ] && [ "${ORCID}" != "" ]; then
                        # Step 4, save the creator id and ORCID
						echo "   Saving ${I} of ${NO_CREATORS}. CREATOR: ${CREATOR_ID} --> ${ORCID}"
						cat <<EOT | dataset -c "${DATASET}" -i - create "${CREATOR_ID}"
{"creator_id":"${CREATOR_ID}","orcid":"${ORCID}"}
EOT
					fi
				fi
			fi
		done
	fi
done
