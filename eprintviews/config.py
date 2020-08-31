
import os
import sys
import json

class Configuration:
    '''This is a configuration object that validates a JSON configuration and provides access to the site settings'''

    def __init__(self):
        '''initialize our configuration object'''
        self.htdocs = 'htdocs'
        self.static = 'static'
        self.templates = 'templates'
        self.eprint_url = ''
        self.dataset = ''
        self.number_of_days = 0
        self.control_item = ''
        self.views = ''
        self.subjects = ''
        self.users = ''
        self.organization = ''
        self.site_title = ''
        self.site_welcome = ''
        self.bucket = ''
        self.distribution_id = ''
        self.config_name = ''

    def load_config(self, f_name):
        '''this reads a JSON configuration from disc and configures self'''
        self.config_name = f_name
        ok = True
        if os.path.exists(f_name):
            with open(f_name) as f:
                src = f.read()
                try:
                    data = json.loads(src)
                except Exception as err:
                    print(f'ERROR reading {f_name}, {err}')
                    return False
                if 'htdocs' in data:
                    self.htdocs = data['htdocs']
                if 'static' in data:
                    self.static = data['static']
                if 'templates' in data:
                    self.templates = data['templates']
                if 'eprint_url' in data:
                    self.eprint_url = data['eprint_url']
                if 'dataset' in data:
                    self.dataset = data['dataset']
                if 'number_of_days' in data:
                    self.number_of_days = data['number_of_days']
                if 'control_item' in data:
                    self.control_item = data['control_item']
                if 'subjects' in data:
                    self.subjects = data['subjects']
                if 'views' in data:
                    self.views = data['views']
                    if not os.path.exists(self.views):
                        print(f'''Can't find view {self.views} listed in {f_name}''')
                        ok = False
                if 'users' in data:
                    self.users = data['users']
                    if not os.path.exists(self.users):
                        print(f'''Can't find users {self.users} listed in {f_name}''')
                        ok = False
                if 'organization' in data:
                    self.organization = data['organization']
                if 'site_welcome' in data:
                    self.site_welcome = data['site_welcome']
                if 'site_title' in data:
                    self.site_title = data['site_title']
        else:
            ok = False
        return ok                

    def required(self, settings):
        '''This checks if the list of configuration fields provided have been set'''
        f_name = self.config_name
        ok = True
        if ('htdocs' in settings) and (self.htdocs == ''):
            print(f'htdocs not set in {f_name}')
            ok = False
        if ('static' in settings) and (self.static == ''):
            print(f'static not set in {f_name}')
            ok = False
        if ('templates' in settings) and (self.templates == ''):
            print(f'templates not set in {f_name}')
            ok = False
        if ('eprint_url' in settings) and (self.eprint_url == ''):
            print(f'eprint_url not set in {f_name}')
            ok = False
        if ('dataset' in settings) and (self.dataset == ''):
            print(f'dataset not set in {f_name}')
            ok = False
        if ('number_of_days' in settings) and (self.number_of_days == 0):
            print(f'number_of_days not set in {f_name}')
            ok = False
        if ('control_item' in settings) and (self.control_item == ''):
            print(f'control_item not set in {f_name}')
            ok = False
        if ('users' in settings) and (self.users == ''):
            print(f'users not set in {f_name}')
            ok = False
        if ('subjects' in settings) and (self.subjects == ''):
            print(f'subjects not set in {f_name}')
            ok = False
        if ('views' in settings) and (self.views == ''):
            print(f'views not set in {f_name}')
            ok = False
        if ('organization' in settings) and (self.organization == ''):
            print(f'organization not set in {f_name}')
            ok = False
        if ('site_welcome' in settings) and (self.site_welcome == ''):
            print(f'site_welcome not set in {f_name}')
            ok = False
        if ('site_title' in settings) and (self.site_title == ''):
            print(f'site_title not set in {f_name}')
            ok = False
        return ok                


