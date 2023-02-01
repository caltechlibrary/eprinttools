---
title: "ep3datasets (1) user manual"
author: "R. S. Doiel"
pubDate: 2022-11-28
---

# NAME

ep3datasets

# SYNOPSIS

ep3datasets [OPTION] JSON_SETTINGS_FILENAME

# DESCRIPTION

ep3datasets is a command line program renders dataset collections
from previously harvested EPrint repositories based on the
configuration in the JSON_SETTINGS_FILENAME.

# OPTIONS

-help
: display help

-license
: display license

-version
: display version

-verbose
: use verbose logging

# EXAMPLES

Rendering harvested repositories for settings.json.

~~~
    ep3datasets settings.json
~~~

ep3datasets 1.2.1

