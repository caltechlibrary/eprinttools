#!/bin/bash

# Setup working environment
WEEKDAY=$(date +%A)
cd $HOME/src/github.com/caltechlibrary/epgo
if [ -f "etc/production.bash" ]; then
    . etc/production.bash
fi
if [ ! -d logs ]; then
    mkdir logs
fi

# Export our EPrint data, build site, rebuild index, start webserver
./bin/epgo -export -1 2> logs/export.$WEEKDAY.log
./bin/genpages 2> logs/genpages.$WEEKDAY.log
./bin/sitmapper -exclude "$EPGO_REPOSITORY_PATH:affilications" "$EPGO_HTDOCS" "$EPGO_HTDOCS/sitemap.xml" "$EPGO_SITE_URL" 2> logs/sitemapper.$WEEKDAY.log
./bin/indexpages -r 2> logs/indexpages.$WEEKDAY.log
