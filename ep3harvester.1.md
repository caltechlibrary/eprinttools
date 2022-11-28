% ep3harvester(1) user manual
% R. S. Doiel
% 2022-11-28

# NAME

ep3harvester

# SYNOPSIS

ep3harvester [OPTION] JSON_SETTINGS_FILENAME \
          [START_TIMESTAMP] [END_TIMESTAMP]

# DESCRIPTION

ep3harvester is a command line program for metadata harvesting
of EPrints repositories.

ep3harvester takes a JSON settings file and harvests
all the EPrint repositories defined in the settings file
into a JSON store implemented in MySQL 8. One repository per
MySQL 8 table.

Each MySQL 8 table has several columns id, src (holding the JSON
document as a JSON column) and an updated (holding the timestamp
of when the metadata was harvested).

## CONFIGURING YOUR JSON STORE

ep3harvester uses a MySQL 8 database for a JSON document store.
It will generate one table for EPrint repository. You can
generate a SQL program for creating the MySQL database and
tables from your settings JSON file using the "-sql-schema"
option. Using the option will require a JSON settings filename
parameter. E.g.

~~~
    ep3harvester -sql-schema settings.json
~~~

# OPTIONS

-h, help
: display help

-license
: display license

-sql-schema
: display SQL schema for installing MySQL jsonstore DB

-version
: display version

# EXAMPLES

Harvesting repositories for week month of May, 2022.

~~~
    ep3harvester settings.json "2022-05-01 00:00:00" "2022-05-31 59:59:59"
~~~

