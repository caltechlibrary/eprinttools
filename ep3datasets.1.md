% ep3datasets(1) user manual
% R. S. Doiel
% 2022-11-28

# NAME

ep3datasets

# SYNOPSIS

ep3datasets [OPTION] JSON_SETTINGS_FILENAME

# DESCRIPTION

ep3datasets is a command line program renders dataset collections
from previously harvested EPrint repositories based on the
configuration in the JSON_SETTINGS_FILENAME.

# OPTIONS

-h, -help
: display help

-license
: display license

-verbose
: use verbose logging

-version
: display version

# EXAMPLES

Rendering harvested repositories for settings.json.

~~~
    ep3datasets settings.json
~~~

