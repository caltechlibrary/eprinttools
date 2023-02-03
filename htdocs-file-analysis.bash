#!/bin/bash

APP_NAME=$(basename "$0")

VERSION="0.0.1"

function usage() {
	cat <<EOT
---
title: "${APP_NAME} (1) user manual"
author: "R. S. Doiel"
pubDate: 2023-02-03
---

# NAME

${APP_NAME}

# SYNOPSIS

${APP_NAME} DIRNAME

# DESCRIPTION

${APP_NAME} generates data files for types of file in a
directory structure such as the htdocs directory generated
for feeds.library.caltech.edu.

# OPTIONS

-help
: display help

-license
: display license

-version
: display version

# EXAMPLES

~~~
cd /Sites/feeds.library.caltech.edu
${APP_NAME} htdocs
~~~

EOT
}

function license() {
	cat <<EOT

Copyright (c) 2023, Caltech
All rights not granted herein are expressly reserved by Caltech.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice,
this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
this list of conditions and the following disclaimer in the documentation
and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors
may be used to endorse or promote products derived from this software without
specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.

EOT
	
}


function analyse_files() {
	HTDOCS="$1"
	DNAME=$(realpath "$HTDOCS")
	DNAME=$(basename "${DNAME}")
	echo "Run analysis on ${HTDOCS} folder name: ${DNAME}"
	find "$HTDOCS" -type f >"${DNAME}-files.txt"
	for EXT in json keys md bib rss html include; do
		RPT_NAME_1="${DNAME}-${EXT}.txt"
		RPT_NAME_2="${DNAME}-${EXT}-types.txt"
		echo "Generating ${RPT_NAME_1}"
	    grep -E "\.${EXT}\$" "${DNAME}-files.txt" >"${RPT_NAME_1}"
		echo "Generating ${RPT_NAME_2} (can take a while)"
	    while read -r FNAME; do
			basename "${FNAME}"
		done <"${RPT_NAME_1}" | sort -u >"${RPT_NAME_2}"
    done	
}

#
# Main processing
#

# Process options and run analysis
for ARG in "$@"; do
   case "${ARG}" in
      -h|-help|--help)
      usage
      exit 0
      ;;
      -version)
      echo "${APP_NAME} ${VERSION}"
      exit 0
      ;;
      -license)
      license
      exit 0
      ;;
      *)
	  if [ -d "$ARG" ]; then
		  analyse_files "${ARG}"
		  exit 0
      fi
   esac
done

echo "Expected a file path, see ${APP_NAME} -help"
exit 1

