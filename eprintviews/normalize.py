#
# utility methods
#

def slugify(s):
    return s.replace(' ', '_').replace('/','_')

def get_value(obj, key):
    if key in obj:
        return obj[key]
    return None

def get_date_year(obj):
    if 'date' in obj:
        return obj['date'][0:4].strip()
    return ''

def get_eprint_id(obj):
    if 'eprint_id' in obj:
        return f"{obj['eprint_id']}"
    return ''

def get_title(obj):
    if 'title' in obj:
        return obj['title']
    return ''

def has_groups(obj):
    if ('local_group' in obj) and (len(obj['local_group']) > 0):
        return True
    return False

def get_object_type(obj):
    if 'type' in obj:
        return f'{obj["type"]}'
    return ''

def has_creator_ids(obj):
    for creator in obj['creators']:
        if 'id' in creator:
            return True
    return False

def make_label(s, sep = '_'):
    l = s.split(sep)
    for i, val in enumerate(l):
        l[i] = val.capitalize()
    return ' '.join(l)

def get_sort_name(o):
    if 'sort_name' in o:
        return o['sort_name']
    return ''

def get_sort_year(o):
    if 'year' in o:
        return o['year']
    return ''

def get_sort_subject(o):
    if 'subject_name' in o:
        return o['subject_name']
    return ''

def get_sort_publication(o):
    if ('publication' in o) and ('item' in publication['publication']):
        return o['publication']['item']
    return ''

def get_sort_collection(o):
    if ('collection' in o):
        return o['collection']
    return ''

def get_sort_event(o):
    if ('event_title' in o):
        return o['event_title'].strip()
    return ''

def get_lastmod_date(o):
    if ('lastmod' in o):
        return o['lastmod'][0:10]
    return ''

def get_sort_lastmod(o):
    if ('lastmod' in o):
        return o['lastmod']
    return ''

def get_sort_issn(o):
    if ('issn' in o):
        return o['issn']
    return ''


#
# normalize_user() is used in Users when loading a JSON exported
# list of users from EPrints. It transform a problematic
# representation of users containing confidential material to one
# which is useful in rendering public webpages.
# This function is not exported.
#
def normalize_user(obj):
    user = {}
    if 'name' in obj:
        name = obj['name']
        sort_name = []
        display_name = []
        if ('honourific' in name) and name['honourific']:
            display_name.append(f'{name["honourific"]}')
        if ('given' in name) and name['given']:
            display_name.append(name["given"])
        if ('family' in name) and name['family']:
            display_name.append(name["family"])
        if ('lineage' in name) and name['lineage']:
            display_name.append(f'{name["lineage"]}')
        user['display_name'] = ' '.join(display_name).strip()
        if ('family' in name) and name['family']:
            sort_name.append(name["family"])
        if ('given' in name) and name['given']:
            sort_name.append(name["given"])
        if ('lineage' in name) and name['lineage']:
            sort_name.append(f'{name["lineage"]}')
        if ('honourific' in name) and name['honourific']:
            sort_name.append(f'{name["honourific"]}')
        user['sort_name'] = ', '.join(sort_name).strip()
    # filter the remaining object fields for 
    # we want to expose in the templates.
    for field in [ 'dept', 'usertype', 'org', 'name', 'joined' ]:
        if field in obj:
            user[field] = obj[field]
    return user
    
#
# NOTE: get_creator_id, get_creator_name and make_creator_list
# are used in normalize_object. They do not need to 
# be exported.
#
def get_creator_id(creator):
    if 'id' in creator:
        return creator['id']
    return ''

def get_creator_name(creator):
    family, given = '', ''
    if 'name' in creator:
        if 'family' in creator['name']:
            family = creator['name']['family'].strip()
        if 'given' in creator['name']:
            given = creator['name']['given'].strip()
    if len(family) > 0:
        if len(given) > 0:
            return f'{family}, {given}'
        return family
    return ''

def make_creator_list(creators):
    l = []
    for creator in creators:
        display_name = get_creator_name(creator)
        creator_id = get_creator_id(creator)
        creator['display_name'] = display_name
        creator['creator_id'] = creator_id
        l.append(creator)
    return l

#
# Normalize object normalizes an JSON representation of an
# eprints xml object.
#
def normalize_object(obj, users):
    title = obj['title'].strip()
    year = get_date_year(obj)
    eprint_id = get_eprint_id(obj)
    creator_list = make_creator_list(obj['creators']['items'])
    if ('abstract' in obj):
        abstract = obj['abstract'].strip()
        obj['abstract'] = abstract
    if ('event_title' in obj):
        event_id = slugify(obj['event_title'])
        obj['event_id'] = event_id
    if 'userid' in obj:
        key = f'{obj["userid"]}'
        if users.has_user(key):
            user = users.get_user(key)
            if 'display_name' in user:
                display_name = user['display_name']
                obj['depositor'] = display_name
    if 'date' in obj:
        obj['year'] = year
    if 'type' in obj:
        obj['type_label'] = make_label(obj['type'])
    obj['title'] = title
    obj['creators'] = creator_list
    obj['eprint_id'] = eprint_id
    obj['year'] = year
    return obj


