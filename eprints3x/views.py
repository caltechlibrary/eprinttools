#
# this provides a means of configuring views that will be supported in the
# website.
#

# This is the default view is non is configured by load_views.
views = {
    # key is view id, value is human readable version
    "ids": "Eprint ID",
    "types": "Document Type",
    "subjects": "Subject",
    "publication": "Publication Title",
    "issn": "ISSN",
    "year": "Year",
    "person": "Person",
    "person-az": "People A to Z"
}

# Override the default views with JSON object of view keys and names
def load_views(f_name):
    global views
    if os.path.exists(f_name):
        views = {}
        with open(f_name) as f:
            src = f.read()
            views = json.loads(src)

def supported_views():
    global views
    keys = []
    for key in views:
        keys.append(key)
    return keys

def has_view(key):
    return key in views

def normalize_view(key):
    if key in views:
        return views[key]
    return key 
