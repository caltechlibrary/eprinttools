#
# user.py models users as implemented in EPrint 3.3.16 and exposed via ep3apid.
#
import json

class Name:
    '''Name models a person name as implemented in EPrints 3.3.16'''
    def __init__(self, m = None):
        self.honourific = ''
        self.family = ''
        self.given = ''
        self.lineage = ''
        if m != None:
            self.from_dict(m)
    
    def from_dict(self, m):
        '''Take a "name" in dict form and populate a Name object.'''
        if 'family' in m:
            self.family = m['family']
        if 'given' in m:
            self.given = m['given']
        if 'lineage' in m:
            self.lineage = m['lineage']
        if 'honourific' in m:
            self.honourific = m['honourific']

    def to_dict(self):
        '''Return a dict version of Name object'''
        m = {}
        if self.honourific:
            m['honourific'] = self.honourific
        if self.family:
            m['family'] = self.family
        if self.given:
            m['given'] = self.given
        if self.lineage:
            m['lineage'] = self.lineage
        return m
    
    def to_string(self):
        m = self.to_dict()
        return json.dumps(m)

    def display_name(self):
        '''format a display name from Name object.'''
        parts = []
        if self.family != '':
            parts.append(self.family)
        if self.given != '':
            parts.append(self.given)
        #if self.lineage:
        #    parts.append(self.lineage)
        #if self.honourific:
        #    parts.append(self.honourific)
        return ', '.join(parts)

class User:
    '''User models the EPrints user table as a Python Object'''
    def __init__(self, m = None):
        '''Creates an unpopulated User object or a populated one from a dict'''
        self.userid  = 0       # integer id value
        self.uname = ''        # username
        self.name = Name()     # Name object
        self.email = ''        # email if hide_email is false
        self.hide_email = True # boolean, include or supress user email
        self.display_name = '' # name_family, name_given
        self.role = ''         # type
        self.created = ''      # joined
        if m != None:
            self.from_dict(m)

    def from_dict(self, m):
        '''Takes a dict and populates User object'''
        if 'userid' in m:
            self.userid = m['userid']
        if 'uname' in m:
            self.uname = m['uname']
        if 'username' in m:
            self.uname = m['username']
        if 'name' in m:
            self.name = Name(m['name'])
            self.display_name = self.name.display_name()
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
        if self.name:
            m['name'] = self.name.to_dict()
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

    def has_role(self, required = None):
        '''Check if User.role to matche required.  required can be a string or list of string'''
        if not self or not self.role:
            return False
        if isinstance(required, list):
            return (self.role in required)
        else:
            return (self.role == required)


