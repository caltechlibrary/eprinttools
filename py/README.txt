
# eprinttools

This project contains the _eprinttools_, a go package for working with EPrints 
REST API. 

## The command line programs

+ _eputil_ is a command line utility for interacting (e.g. harvesting) JSON and XML from EPrints' REST API
    + uses minimal configuration because it does less!
    + will supercede _ep_
+ _ep_ is a EPrints harvester that integrates with [dataset](https://github.com/caltechlibrary/dataset).

Both utilities can be configured from the environment or command line options. The environment
settings are overridden by command line options. For details running either command envoke the
tool name with the '-help' option. 

## Python 3.6 module

This repository also contains a Python 3.6 module that wraps the basic Go package giving you
the functionality of the command line tools in a Python 3.6 package.


