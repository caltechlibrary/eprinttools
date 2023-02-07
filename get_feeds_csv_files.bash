#!/bin/bash
#

APP_NAME=$(basename "$0")
VERSION="0.0.1"

function usage() {
	cat <<EOT
---
title: "${APP_NAME} (1) user manual"
author: "R. S. Doiel"
pubDate: 2023-02-06
---

# NAME

${APP_NAME}

# SYNOPSIS

${APP_NAME} [OPTIONS]

# DESCRIPTION

${APP_NAME} will fetch the people.csv and groups.csv file from
feeds.library.caltech.edu.

# OPTIONS

-h, -help, --help
: display help

-license
: display license

-version
: display version

# EXAMPLES

~~~
    ${APP_NAME}
~~~

EOT
}

function license() {
	cat <<EOT

Copyright (c) 2022, Caltech
All rights not granted herein are expressly reserved by Caltech.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

EOT

}

function version() {
	echo "${APP_NAME} ${VERSION}"
}

function retrieve_csv_files() {
	curl -L -O https://feeds.library.caltech.edu/people/people.csv
	curl -L -O https://feeds.library.caltech.edu/groups/groups.csv
}

#
# Main
#

# Handle common options
for $ARG in $@; do
	case "$ARG" in
		-h|-help|--help)
		usage
		exit 0
		;;
		-version)
		version
		exit 0
		;;
		-license)
		license
		exit 0
		;;
	esac
done

retrieve_csv_files

