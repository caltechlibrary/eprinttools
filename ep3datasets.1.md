---
title: "ep3datasets (1) user manual"
author: "R. S. Doiel"
pubDate: 2023-02-07
---

# NAME

ep3datasets

# SYNOPSIS

ep3datasets [OPTION] JSON_SETTINGS_FILE

# DESCRIPTION

ep3datasets is a command line program renders dataset collections
from previously harvested EPrint repositories based on the
configuration in the JSON_SETTINGS_FILE.

# OPTIONS

-help
: display help

-license
: display license

-version
: display version

-verbose
: use verbose logging

-repo
: write out the dataset for a specific repo in JSON_SETTINGS_FILE

# EXAMPLES

Rendering harvested repositories for settings.json.

~~~
    ep3datasets settings.json
~~~

Render the harvested repository caltechauthors based on settings.json.

~~~
	ep3datasets -repo caltechauthors settings.json
~~~

ep3datasets 1.2.3


