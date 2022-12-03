% ep3feeds(1) user manual
% R. S. Doiel
% 2022-11-28

# NAME

ep3feeds

# SYNOPSIS

ep3feeds [OPTION] JSON_SETTINGS_FILENAME

# DESCRIPTION

ep3feeds is a command line program that renders the EPrint harvested
metadata and aggregation tables to JSON documents, non-templated
Markdown documents and the necessary directory structures needed for
representing EPrints repositories as a static site.

The configuration needs to be previously created using the 
ep3harvester tool.

# OPTIONS

  -h	display help
  -help
    	display help
  -license
    	display license
  -verbose
    	use verbose logging
  -version
    	display version

# EXAMPLES

Harvesting repositories for week month of May, 2022.

~~~
    ep3feeds harvester-settings.json
~~~

