#!/bin/bash

# Source our configuration
if [ -f etc/orcid-api.bash ]; then
    . etc/orcid-api.bash
else
        echo "Can't find etc/orcid-api.bash holding configuration"
        exit 1
fi
bin/apiexplorer https://pub.orcid.org/v2.0/0000-0003-0900-6903/person
