#!/bin/bash

#
# Harvest the ORCID profiles via the ORCID API
#
function HarvestProfiles() {
	template="$1"
	profilePath="$2"
	for ITEM in $(find "$profilePath" -depth 1 -type d | sed -E 's/\// /g' | cut -d\  -f 3); do
		echo "Building orcid profile for $ITEM"
		orcid -profile "$ITEM" >"$profilePath/$ITEM/orcid-profile.json"
		# FIXME: I need to pickup the template set from the environment
		mkpage ".=$profilePath/$ITEM/orcid-profile.json" "$template" >"$profilePath/$ITEM/index.html"
	done
}

#
# Check to make sure all the environment vars needed are defined before running harvest
#
function CheckEnv() {
	VarName="$1"
	if [ "$2" = "" ]; then
		echo "Missing $VarName"
		exit 1
	fi
}

#
# Main processing
#
CheckEnv "ORCID_API_URL" "$ORCID_API_URL"
CheckEnv "ORCID_CLIENT_ID" "$ORCID_CLIENT_ID"
CheckEnv "ORCID_CLIENT_SECRET" "$ORCID_CLIENT_SECRET"
if [ ! -d htdocs/person ]; then
	echo "Can't find htdocs/person"
	exit 1
fi
# FIXME: Need to check to make sure template is available too and set from the environment 
if [ ! -f templates/default/orcid-profile.tmpl ]; then
	echo "Can't find templates/default/orcid-profile.tmpl"
	exit 1
fi
HarvestProfiles templates/default/orcid-profile.tmpl htdocs/person
