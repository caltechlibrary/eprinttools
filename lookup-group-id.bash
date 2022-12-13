#!/bin/bash

function usage() {
	APP_NAME=$(basename "$0")
	cat <<EOT
% ${APP_NAME}(1) user manual
% R. S. Doiel
% 2022-12-13

# NAME

${APP_NAME}

# SYSNOPSIS

${APP_NAME} GROUP_NAME

# DESCRIPTION

${APP_NAME} takes a group name and looks them up in the _groups table
checking both the official name and alternative names.

# EXAMPLES

~~~
${APP_NAME} "Mary A. Earl McKinney Prize in Literature - Poetry"
~~~

The value printed out should be the group id or an indication it was not found.

EOT

}

#
# Main
# 
if [ "$1" = "" ]; then
	usage
	exit 1
fi

GROUP_NAME="$1"
mysql collections \
	--execute \
"SELECT group_id FROM _groups WHERE (name LIKE '${GROUP_NAME}') OR (LOCATE('${GROUP_NAME}', alternative) > 0)"
