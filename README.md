
# eprinttools

This project contains the _eprinttools_, a go package for working with EPrints 
REST API. It also includes _ep_ and command line utility for 
harvesting content into a [dataset](https://github.com/caltechlibrary/dataset)
collection and rendering JSON documents for web feeds.

## The command line 

+ _eputil_ is a command line utility for interacting (e.g. harvesting) JSON and XML from EPrints' REST API
    + uses minimal configuration because it does less!
    + will supercede _ep_
+ _ep_ is an older version of a EPrints harvester. It features an integration with 
       [dataset](https://github.com/caltechlibrary/dataset) and can be used to produce 
       alternative feeds and formats. 

The _ep_ utility is configured from the environment or command line options. The environment
settings are overridden by command line options. For details running either command envoke the
tool name with the '-help' option. 


