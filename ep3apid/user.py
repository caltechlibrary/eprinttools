
import json

class User:
    '''User models the EPrints user table as a Python Object'''
    def __init__(self):
        '''Creates an unpopulated User object'''
        self.userid  = 0       # integer id value
        self.uname = ''        # username
        self.email = ''        # email if hide_email is false
        self.hide_email = True # boolean, include or supress user email
        self.display_name = '' # name_family, name_given
        self.role = ''         # type
        self.created = ''      # joined

    def from_dict(self, m):
        '''Takes a dict and populates User object'''
        if 'userid' in m:
            self.userid = m['userid']
        if 'uname' in m:
            self.uname = m['uname']
        if 'username' in m:
            self.uname = m['username']
        if 'email' in m:
            self.email = m['email']
        if 'hide_email' in m:
            self.hide_email = m['hide_email']
        if 'hideemail' in m:
            self.hide_email = m['hideemail']
        if 'role' in m:
            self.role = m['role']
        if 'type' in m:
            self.role = m['type']
        if 'joined' in m:
            self.created = m['joined']
        if 'created' in m:
            self.created = m['created']

    def to_dict(self):
        '''Takes user object and returns dict version'''
        m = {}
        if self.userid:
            m['userid'] = self.userid
        if self.uname:
            m['uname'] = self.uname
        if self.email:
            m['email'] = self.email
        if self.hide_email:
            m['hide_email'] = self.hide_email
        if self.display_name:
            m['display_name'] = self.display_name
        if self.role:
            m['role'] = self.role
        if self.created:
            m['created'] = self.created
        return m

    def from_JSON(self, src):
        '''Takes JSON source and populates user object'''
        if not isinstance(src, byte):
            src = src.encode('utf-8')
        obj = json.loads(src)
        self.fromDict(obj)

    def to_string(self):
        '''returns user object as JSON string'''
        return json.dumps(self.to_dict())

