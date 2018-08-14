
# eprinttools

This project contains the _eprinttools_, a go package for working with EPrints 
REST API. 

## The command line programs

+ [eputil](docs/eputil.html) is a command line utility for interacting (e.g. harvesting) JSON and XML from EPrints' REST API
    + uses minimal configuration because it does less!
    + will supercede _ep_
+ [ep](docs/ep.html) is a EPrints harvester that integrates with [dataset](https://github.com/caltechlibrary/dataset)
+ [doi2eprintxml](docs/doi2eprintxml.html) is a command line program for turning metadata harvested from CrossRef and DataCite into an EPrint XML document based on one or more supplied DOI
+ [eprintxml2json](docs/eprintxml2json.html) is a command line program for taking EPrint XML and turning it into JSON 

The first two utilities can be configured from the environment or 
command line options. The environment settings are overridden by command 
line options. For details running either command envoke the
tool name with the '-help' option. 

## Python 3.6 module

This repository also contains a Python 3.6 module that wraps the 
basic Go package giving you the functionality of the command line tools 
in a Python 3.6 package.


