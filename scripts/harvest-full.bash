#!/bin/bash

# Setup working environment
WEEKDAY=$(date +%A)
cd "$HOME/src/github.com/caltechlibrary/ep"
if [ -f "etc/production.bash" ]; then
	. etc/production.bash
fi
if [ ! -d logs ]; then
	mkdir logs
fi

# Export our EPrint data, build site, rebuild index, start webserver
if [ -f "logs/harvest.${WEEKDAY}.log" ]; then
	/bin/rm "logs/harvest.${WEEKDAY}.log"
fi
{
	./bin/ep -export -1
	./bin/ep-genpages
} >>"logs/harvest.${WEEKDAY}.log"
echo "Site rebuilt"
