---
title: "ep3harvester (1) user manual"
author: "R. S. Doiel"
pubDate: 2022-11-28
---

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

ep3harvester can generate an example settings JSON document. You
can then edit it with any plain text editor (e.g. nano). Then
you'll need to setup a MySQL 8 database and tables to store
havested data in.

ep3harvester uses a MySQL 8 database for a JSON document store.
It will generate one table for EPrint repository. You can
generate a SQL program for creating the MySQL database and
tables from your settings JSON file using the "-sql-schema"
option. Using the option will require a JSON settings filename
parameter. E.g.

~~~
    ep3harvester -init harvester-settings.json
    nano harvester-settings.json
    ep3harvester -sql-schema harvester-settings.json >collections.sql
~~~

# OPTIONS

-help
: display help

-version
: display version

-license
: display license

-groups
: Harvest groups from CSV files included configuration

-init
: generate a settings JSON file

-people
: Harvest people from CSV files included configuration

-people-groups
: Harvest people and groups from CSV files included configuration

-repo string
: Harvest a specific repository id defined in configuration

-sql-schema
: display SQL schema for installing MySQL jsonstore DB

-verbose
: use verbose logging

# EXAMPLES

Harvesting repositories for week month of May, 2022.

~~~
    ep3harvester harvester-settings.json \
        "2022-05-01 00:00:00" "2022-05-31 59:59:59"
~~~

ep3harvester 1.2.1


