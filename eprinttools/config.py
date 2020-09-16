
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
        self.include_documents = False
        self.include_documents_set = False
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
        self.base_url = ''
        self.base_path = ''
        self.elastic_documents_max_no = 25
        self.elastic_documents = ''
        self.elastic_base_endpoint = ''
        self.elastic_api_key = ''
        self.elastic_engine_name = ''

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
                if 'elastic_documents_max_no' in data:
                    self.elastic_documents_max_no = data['elastic_documents_max_no']
                if 'elastic_documents' in data:
                    self.elastic_documents = data['elastic_documents']
                if 'elastic_base_endpoint' in data:
                    self.elastic_base_endpoint = data['elastic_base_endpoint']
                if 'elastic_api_key' in data:
                    self.elastic_api_key = data['elastic_api_key']
                if 'elastic_engine_name' in data:
                    self.elastic_engine_name = data['elastic_engine_name']
                if 'include_documents' in data:
                    self.include_documents = data['include_documents']
                    self.include_documents_set = True
                if 'base_path' in data:
                    self.base_path = data['base_path']
                if 'base_url' in data:
                    self.base_url = data['base_url']
                if 'htdocs' in data:
                    self.htdocs = data['htdocs']
                    if not os.path.exists(self.htdocs):
                        print(f'ERROR: htdocs setting, "{self.htdocs}" from {f_name} does not exist.')
                        ok = False
                if 'static' in data:
                    self.static = data['static']
                    if not os.path.exists(self.static):
                        print(f'ERROR: htdocs setting, "{self.static}" from {f_name} does not exist.')
                        ok = False
                if 'templates' in data:
                    self.templates = data['templates']
                    if not os.path.exists(self.templates):
                        print(f'ERROR: htdocs setting, "{self.templates}" from {f_name} does not exist.')
                        ok = False
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
                if 'bucket' in data:
                    self.bucket = data['bucket']
                if 'distribution_id' in data:
                    self.distribution_id = data['distribution_id']
        else:
            ok = False
        return ok                

    def required(self, settings):
        '''This checks if the list of configuration fields provided have been set'''
        f_name = self.config_name
        ok = True
        if ('elastic_documents_max_no' in settings):
            if (self.elastic_documents_max_no == 0):
                print(f'elastic_documents_max_no not set in {f_name}')
                ok = False
        if ('elastic_documents' in settings):
            if (self.elastic_documents == ''):
                print(f'elastic_documents not set in {f_name}')
                ok = False
        if ('elastic_base_endpoint' in settings):
            if (self.elastic_base_endpoint == ''):
                print(f'elastic_base_endpoint not set in {f_name}')
                ok = False
        if ('elastic_api_key' in settings):
            if (self.elastic_api_key == ''):
                print(f'elastic_api_key not set in {f_name}')
                ok = False
        if ('elastic_engine_name' in settings):
            if (self.elastic_engine_name == ''):
                print(f'elastic_engine_name not set in {f_name}')
                ok = False
        if ('include_documents' in settings):
            if self.include_documents_set == False:
                print(f'include_documents not set in {f_name}')
                ok = False
        if ('control_item' in settings):
            if (self.control_item == ''):
                print(f'control_item not set in {f_name}')
                ok = False
        if ('base_path' in settings):
            if (self.base_path == ''):
                print(f'base_path not set in {f_name}')
                ok = False
        if ('base_url' in settings):
            if (self.base_url == ''):
                print(f'base_url not set in {f_name}')
                ok = False
        if ('htdocs' in settings):
            if (self.htdocs == ''):
                print(f'htdocs not set in {f_name}')
                ok = False
            elif not os.path.exists(self.htdocs):
                print(f'htdocs {self.htdocs} does not exist.')
                ok = False
        if ('static' in settings):
            if (self.static == ''):
                print(f'static not set in {f_name}')
                ok = False
            elif not os.path.exists(self.static):
                print(f'static {self.static} does not exist.')
                ok = False
        if ('templates' in settings):
            if (self.templates == ''):
                print(f'templates not set in {f_name}')
                ok = False
            elif not os.path.exists(self.templates):
                print(f'templates {self.templates} does not exist.')
                ok = False
        if ('eprint_url' in settings) and (self.eprint_url == ''):
            print(f'eprint_url not set in {f_name}')
            ok = False
        if ('dataset' in settings):
            if (self.dataset == ''):
                print(f'dataset not set in {f_name}')
                ok = False
            elif not os.path.exists(self.dataset):
                print(f'dataset {self.dataset} does not exist.')
                ok = False
        if ('number_of_days' in settings) and (self.number_of_days == 0):
            print(f'number_of_days not set in {f_name}')
            ok = False
        if ('control_item' in settings) and (self.control_item == ''):
            print(f'control_item not set in {f_name}')
            ok = False
        if ('users' in settings):
            if (self.users == ''):
                print(f'users not set in {f_name}')
                ok = False
            elif not os.path.exists(self.users):
                print(f'users {self.users} does not exist.')
                ok = False
        if ('subjects' in settings):
            if (self.subjects == ''):
                print(f'subjects not set in {f_name}')
                ok = False
            elif not os.path.exists(self.subjects):
                print(f'subjects {self.subjects} does not exist.')
                ok = False
        if ('views' in settings):
            if (self.views == ''):
                print(f'views not set in {f_name}')
                ok = False
            elif not os.path.exists(self.views):
                print(f'views {self.views} does not exist.')
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
        if ('bucket' in settings) and (self.bucket == ''):
            print(f'bucket not set in {f_name}')
            ok = False
        if ('distribution_id' in settings) and (self.distribution_id == ''):
            print(f'distribution_id not set in {f_name}')
            ok = False
        return ok                

    def toJSON(self):
        o = {}
        if self.elastic_documents_max_no > 0:
            o['elastic_documents_max_no'] = self.elastic_documents_max_no
        if self.elastic_documents != '':
            o['elastic_documents'] = self.elastic_documents
        if self.elastic_base_endpoint != '':
            o['elastic_base_endpoint'] = self.elastic_base_endpoint
            o['elastic_api_key'] = self.elastic_api_key
        if self.include_documents_set == True:
            o['include_documents'] = self.include_documents
        if self.control_item != '':
            o['control_item'] = self.control_item
        if self.base_path != '':
            o['base_path'] = self.base_path
        if self.base_url != '':
            o['base_url'] = self.base_url
        if self.htdocs != '':
            o['htdocs'] = self.htdocs
        if self.static != '':
            o['static'] = self.static
        if self.templates != '':
            o['templates'] = self.templates
        if self.eprint_url != '':
            o['eprint_url'] = self.eprint_url
        if self.dataset != '':
            o['dataset'] = self.dataset
        if self.number_of_days != 0:
            o['number_of_days'] = self.number_of_days
        if self.control_item != '':
            o['control_item'] = self.control_item
        if self.users != '':
            o['users'] = self.users
        if self.subjects != '':
            o['subjects'] = self.subjects
        if self.views != '':
            o['views'] = self.views
        if self.organization != '':
            o['organization'] = self.organization
        if self.site_welcome != '':
            o['site_welcome'] = self.site_welcome
        if self.site_title != '':
            o['site_title'] = self.site_title
        if self.bucket != '':
            o['bucket'] = self.bucket
        if self.distribution_id != '':
            o['distribution_id'] = self.distribution_id
        s = json.dumps(o)
        return s
