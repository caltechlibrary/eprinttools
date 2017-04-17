
# epgo

This project contains the _epgo_ go package for working with EPrints 
REST API. It also includes a set of command line utilities for 
harvesting content,  building feeds and website pages that run 
indepenantly of EPrints itself (e.g. You could have a website/feeds 
generated from an EPrints repository running on a different system).

## The command line utilities

+ _epgo_ is a command line utility utilizing EPrints' REST API to 
  produce alternative feeds and formats. Currently it supports 
  generating a feed of repository items based on publication dates.
+ _epgo-genpages_ is a command line utility that builds JSON documents 
  based on the content harvested with _epgo_.
    + For HTML, HTML include, RSS 2 documents [mkpage](https://caltechlibrary.github.io/mkpage)

The utilities can be configured from the environment.  The environment
can be overridden by command line options. For details run the individual command wiht the '-help'
option.  E.g. `./bin/epgo -help`


