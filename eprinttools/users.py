import json
import sys

from .normalize import normalize_user

class Users:
    """Model the Eprint user classes so we can resolve the user id mapping to a human readable name"""
    def __init__(self):
        self.users = {}
    
    # load users from a JSON dump of users from EPrints
    # NOTE: this function calls normalize_user on each object loaded
    # so that we don't accidentally expose confidential info.
    def load_users(self, f_name):
        objects = []
        with open(f_name) as f:
            src = f.read()
            try:
                objects = json.loads(src)
            except Exception as err:
                print(f'''Failed to parse JSON file {f_name}, {err}''')
                sys.exit(1)
        for i, obj in enumerate(objects):
            user = {}
            # only include the user if we can derive a name from user id
            if not 'userid' in obj:
                continue
            if not 'name' in obj:
                continue
            key = f"{obj['userid']}"
            obj = normalize_user(obj)
            self.users[key] = obj
    
    
    def has_user(self, user_id):
        if str(user_id) in self.users:
            return True
        return False
    
    def get_user(self, user_id):
        if str(user_id) in self.users:
            return self.users[str(user_id)]
        return None
    
    def user_list(self):
        l = []
        keys = []
        for key in self.users:
            keys.append(key)
        keys.sort(key=int)
        for key in keys:
            l.append(self.users[key])
        return l
    
