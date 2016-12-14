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

# NOTE: Cycle through the indexes as we rebuild them.
bleveIndexes=${EPGO_BLEVE/:/ }
echo "bleveIndex: [$bleveIndexes]"
for I in $bleveIndexes; do
    echo "Index $I"
    # Bump from the first index to next, rebuild previous
    pids=$(pgrep servepages)
    if [ "$pids" != "" ]; then
        echo "Sending a request to Swaping indexes"
        kill -s HUP $pids
    fi
    echo "Replacing $I"
    ./bin/indexpages -r $I  >> logs/harvest.$WEEKDAY.log
done
echo "Site and Indexes rebuilt"