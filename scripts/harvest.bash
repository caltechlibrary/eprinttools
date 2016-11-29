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
if [ -f "logs/harvest.$WEEKDAY.log" ]; then
    /bin/rm "logs/harvest.$WEEKDAY.log"
fi
./bin/epgo -export -1 >> logs/harvest.$WEEKDAY.log
./bin/genpages >> logs/harvest.$WEEKDAY.log
./bin/sitemapper -exclude "$EPGO_REPOSITORY_PATH:affilications" "$EPGO_HTDOCS" "$EPGO_HTDOCS/sitemap.xml" "$EPGO_SITE_URL" >> logs/harvest.$WEEKDAY.log
#FIXME: Need to pick the correct Index Name rather then the default first Bleve DB Name in the list
./bin/indexpages -r >> logs/harvest.$WEEKDAY.log
