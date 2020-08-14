#
# this provides a means of configuring views that will be supported in the
# website.
#

import os
import json

# This is the default view is non is configured by load_views.
# 
# key is view id, value is human readable version
#  "ids": "Eprint ID",
#  "types": "Document Type",
#  "subjects": "Subject",
#  "publication": "Publication Title",
#  "issn": "ISSN",
#  "year": "Year",
#  "person": "Person",
#  "person-az": "People A to Z"
#

# Override the default views with JSON object of view keys and names
def load_views(f_name):
    views = {}
    if os.path.exists(f_name):
        with open(f_name) as f:
            src = f.read()
            views = json.loads(src)
    return views

