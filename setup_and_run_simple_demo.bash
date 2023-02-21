#!/bin/bash

if [[ ! -f simple-settings.json ]]; then
	APP_NAME=$(basename $0)
	cat <<EOT
For ${APP_NAME} to work you need to have a 
settings.json file holding the congifuration of the
EPrints repositories you want to demo in feeds form.
You will also need to server out the demo/htdocs directory
on localhost to test with.

Requirements

- Go 1.20 or better
	- compiler for eprinttools
- GNU Make
	- Runs build process
- Bash
	- Runs various scripts for like this one
- MySQL 8 holding  (mysql server and mysql client)
	- It needs to contain EPrints databases
	- It gets used for harvesting and staging the demo content

EOT

fi

if make; then
	./bin/ep3harvester -sql-schema simple-settings.json | mysql
	./bin/ep3harvester -verbose -people-groups simple-settings.json
	./bin/ep3harvester -simple -verbose simple-settings.json
	./bin/ep3datasets -verbose simple-settings.json
	#./bin/ep3genfeeds -verbose settings.json
fi

