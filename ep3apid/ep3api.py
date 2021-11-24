#!/usr/bin/env python3

#
# ep3apid.py is responsible interfacing with EPrints for via the
# ep3apid web service (part of EPrinttools).
#
import os
import sys

import json
from urllib.request import Request, urlopen
from urllib.error import URLError, HTTPError
from urllib.parse import quote_plus

def dquote(s):
    return '"' + s + '"'

def squote(s):
    return "'" + s + "'"

#
# Handle HTTP Requests
#
def http_get(u, headers = None):
    '''http_get takes a URL and performs a GET. It returns a touple of payload and error'''
    req = Request(u, headers = headers)
    try:
        res = urlopen(req)
    except HTTPError as e:
        return '', f'{e.code}, {e.reason}'
    except URLError as e:
        return '', f'{e.reason}'
    else:
        return res.read(), None

def http_post(u, content_type, data):
    '''http_post takes a URL, content type and data and returns any payload and err'''
    headers = {'Content-Type': content_type}
    req = Request(u, data, headers)
    try:
        res = urlopen(req)
    except HTTPError as e:
        return '', f'{e.code}, {e.reason}'
    except URLError as e:
        return '', f'{e.reason}'
    else:
        return res.read(), None

def get_json_data(u):
    headers = {'Content-Type': 'application/json'}
    src, err = http_get(u, headers = headers)
    if err != None:
        return None, err
    if not isinstance(src, bytes):
        src = src.encode('utf-8')
    return json.loads(src), None

def post_xml(u, xml_src):
    src, err = http_post(u, 'application/xml', xml_src)
    if err != None:
        return None, err
    if not isinstance(src, bytes):
        src = src.encode('utf-8')
    return json.loads(src), None
    
class Ep3API:
    '''Ep3API provides data access to the ep3apid web service for a specific repository'''
    url = 'http://localhost:8484'
    repo_id = None

    def __init__(self, url = 'http://localhost:8484', repo_id = None):
        '''Initialize with the URL to ep3apid and repository ID'''
        self.url = url
        self.repo_id = repo_id

    def use(self, repo_id):
        '''Set the repository name Ep3API uses. Returns True if repository found False if not'''
        repositories, err = self.repositories()
        if (not err) and (repo_id in repositories):
            self.repo_id = repo_id
            return True
        return False

    def repositories(self):
        '''Return a list of repositories available'''
        return get_json_data(f'{self.url}/repositories')

    def repository(self):
        '''Return a list of tables in repository'''
        return get_json_data(f'{self.url}/repository/{self.repo_id}')


    #
    # The following methods returns list of eprint ids and error tuples
    #
    def keys(self):
        '''Return a list of eprint ids'''
        return get_json_data(f'{self.url}/{self.repo_id}/keys')


    def doi(self, doi = None):
        if doi == None:
            return get_json_data(f'{self.url}/{self.repo_id}/doi')
        else:
            return get_json_data(f'{self.url}/{self.repo_id}/doi/{doi}')

    def pmid(self, pmid = None):
        if pmid == None:
            return get_json_data(f'{self.url}/{self.repo_id}/pmid')
        else:
            return get_json_data(f'{self.url}/{self.repo_id}/pmid/{pmid}')

    def pmcid(self, pmcid = None):
        if pmcid == None:
            return get_json_data(f'{self.url}/{self.repo_id}/pmcid')
        else:
            return get_json_data(f'{self.url}/{self.repo_id}/pmcid/{pmcid}')


    def creator_id(self, creator_id = None):
        if creator_id == None:
            return get_json_data(f'{self.url}/{self.repo_id}/creator-id')
        else:
            return get_json_data(f'{self.url}/{self.repo_id}/creator-id/{creator_id}')

    def creator_name(self, creator_name = None):
        if creator_name == None:
            return get_json_data(f'{self.url}/{self.repo_name}/creator-name')
        else:
            return get_json_data(f'{self.url}/{self.repo_name}/creator-name/{creator_name}')
        
    def creator_orcid(self, creator_orcid = None):
        if creator_orcid == None:
            return get_json_data(f'{self.url}/{self.repo_orcid}/creator-orcid')
        else:
            return get_json_data(f'{self.url}/{self.repo_orcid}/creator-orcid/{creator_orcid}')
        

    def editor_id(self, editor_id = None):
        if editor_id == None:
            return get_json_data(f'{self.url}/{self.repo_id}/editor-id')
        else:
            return get_json_data(f'{self.url}/{self.repo_id}/editor-id/{editor_id}')
    
    def editor_name(self, editor_name = None):
        if editor_name == None:
            return get_json_data(f'{self.url}/{self.repo_name}/editor-name')
        else:
            return get_json_data(f'{self.url}/{self.repo_name}/editor-name/{editor_name}')
        

    def contributor_id(self, contributor_id = None):
        if contributor_id == None:
            return get_json_data(f'{self.url}/{self.repo_id}/contributor-id')
        else:
            return get_json_data(f'{self.url}/{self.repo_id}/contributor-id/{contributor_id}')
    
    def contributor_name(self, contributor_name = None):
        if contributor_name == None:
            return get_json_data(f'{self.url}/{self.repo_name}/contributor-name')
        else:
            return get_json_data(f'{self.url}/{self.repo_name}/contributor-name/{contributor_name}')
        
    def advisor_id(self, advisor_id = None):
        if advisor_id == None:
            return get_json_data(f'{self.url}/{self.repo_id}/advisor-id')
        else:
            return get_json_data(f'{self.url}/{self.repo_id}/advisor-id/{advisor_id}')
    
    def advisor_name(self, advisor_name = None):
        if advisor_name == None:
            return get_json_data(f'{self.url}/{self.repo_name}/advisor-name')
        else:
            return get_json_data(f'{self.url}/{self.repo_name}/advisor-name/{advisor_name}')
        
    def committee_id(self, committee_id = None):
        if committee_id == None:
            return get_json_data(f'{self.url}/{self.repo_id}/committee-id')
        else:
            return get_json_data(f'{self.url}/{self.repo_id}/committee-id/{committee_id}')
    
    def committee_name(self, committee_name = None):
        if committee_name == None:
            return get_json_data(f'{self.url}/{self.repo_name}/committee-name')
        else:
            return get_json_data(f'{self.url}/{self.repo_name}/committee-name/{committee_name}')
        
    def corp_creator_id(self, corp_creator_id = None):
        if corp_creator_id == None:
            return get_json_data(f'{self.url}/{self.repo_id}/corp-creator-id')
        else:
            return get_json_data(f'{self.url}/{self.repo_id}/corp-creator-id/{corp_creator_id}')
    
    def corp_creator_name(self, corp_creator_name = None):
        if corp_creator_name == None:
            return get_json_data(f'{self.url}/{self.repo_name}/corp-creator-name')
        else:
            return get_json_data(f'{self.url}/{self.repo_name}/corp-creator-name/{corp_creator_name}')
        

    def group_id(self, group_id = None):
        if group_id == None:
            return get_json_data(f'{self.url}/{self.repo_id}/group-id')
        else:
            return get_json_data(f'{self.url}/{self.repo_id}/group-id/{grou_id}')
        
    def funder_id(self, funder_id = None):
        if funder_id == None:
            return get_json_data(f'{self.url}/{self.repo_id}/funder-id')
        else:
            return get_json_data(f'{self.url}/{self.repo_id}/funder-id/{founder_id}')
    
    def grant_number(self, grant_no = None):
        if grant_no == None:
            return get_json_data(f'{self.url}/{self.repo_id}/grant-number')
        else:
            return get_json_data(f'{self.url}/{self.repo_id}/grant-number/{grant_no}')
    
    def issn(self, issn = None):
        if issn == None:
            return get_json_data(f'{self.url}/{self.repo_id}/issn')
        else:
            return get_json_data(f'{self.url}/{self.repo_id}/issn/{issn}')

    def isbn(self, isbn = None):
        if isbn == None:
            return get_json_data(f'{self.url}/{self.repo_id}/isbn')
        else:
            return get_json_data(f'{self.url}/{self.repo_id}/isbn/{isbn}')

    def patent_number(self, patent_number = None):
        if patent_number == None:
            return get_json_data(f'{self.url}/{self.repo_id}/patent-number')
        else:
            return get_json_data(f'{self.url}/{self.repo_id}/patent-number/{patent_number}')

    def patent_applicant(self, patent_applicant = None):
        if patent_applicant == None:
            return get_json_data(f'{self.url}/{self.repo_id}/patent-applicant')
        else:
            return get_json_data(f'{self.url}/{self.repo_id}/patent-applicant/{patent_applicant}')

    def patent_classification(self, patent_classification = None):
        if patent_classification == None:
            return get_json_data(f'{self.url}/{self.repo_id}/patent-classification')
        else:
            return get_json_data(f'{self.url}/{self.repo_id}/patent-classification/{patent_classification}')

    def patent_assignee(self, patent_assignee = None):
        if patent_assignee == None:
            return get_json_data(f'{self.url}/{self.repo_id}/patent-assignee')
        else:
            return get_json_data(f'{self.url}/{self.repo_id}/patent-assignee/{patent_assignee}')

    def year(self, year = None):
        if year == None:
            return get_json_data(f'{self.url}/{self.repo_id}/year')
        else:
            return get_json_data(f'{self.url}/{self.repo_id}/year/{year}')

    def eprint(self, eprint_id):
        return get_json_data(f'{self.url}/{self.repo_id}/eprint/{eprint_id}')

    def eprint_import(self, userid, eprint_xml = None):
        if eprint_xml == None:
            return [], 'missing eprint xml'
        return post_xml(f'{self.url}/{self.repo_id}/eprint-import/{user_id}', eprint_xml)

    def user(self, username_or_id):
        s = quote_plus(username_or_id)
        return get_json_data(f'{self.url}/{self.repo_id}/user/{s}')

    def usernames(self):
        return get_json_data(f'{self.url}/{self.repo_id}/usernames')

    def lookup_userid(self, username):
        return get_json_data(f'{self.url}/{self.repo_id}/lookup-userid/{username}')


