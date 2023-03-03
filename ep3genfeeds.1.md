---
title: "ep3genfeeds (1) user manual"
author: "R. S. Doiel"
pubDate: 2022-11-28
---

# NAME

ep3genfeeds

# SYNOPSIS

ep3genfeeds [OPTION] JSON_SETTINGS_FILENAME

# DESCRIPTION

ep3genfeeds is a command line program that renders the EPrint harvested
metadata and aggregation tables to JSON documents, non-templated
Markdown documents and the necessary directory structures needed for
representing EPrints repositories as a static site.

The configuration needs to be previously created using the 
ep3harvester tool.

# OPTIONS

-help
: display help

-license
: display license

-version
: display version

-groups
: render groups feeds

-people
: render people feeds

-verbose
: use verbose logging

# EXAMPLES

Harvesting repositories for week month of May, 2022.

~~~
    ep3genfeeds harvester-settings.json
~~~

ep3genfeeds 1.2.4


