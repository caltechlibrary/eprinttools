import json

users = {}

# load users from a JSON dump of users from EPrints
# NOTE: this function calls normalize_user on each object loaded
# so that we don't accidentally expose confidential info.
def load_users(f_name):
    global users
    objects = []
    with open(f_name) as f:
        src = f.read()
        objects = json.loads(src)
    for i, obj in enumerate(objects):
        user = {}
        # only include the user if we can derive a name from user id
        if not 'userid' in obj:
            continue
        if not 'name' in obj:
            continue
        key = f"{obj['userid']}"
        obj = normalize_user(obj)
        users[key] = obj


def has_user(user_id):
    global users
    if str(user_id) in users:
        return True
    return False

def get_user(user_id):
    global users
    if str(user_id) in users:
        return users[str(user_id)]
    return None

def user_list():
    l = []
    keys = []
    for key in users:
        keys.append(key)
    keys.sort(key=int)
    for key in keys:
        l.append(users[key])
    return l

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

