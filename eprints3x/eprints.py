#!/usr/bin/env python3
import os
import sys
import json
from urllib.parse import urlparse
from subprocess import run, Popen, PIPE
from datetime import datetime, timedelta

import progressbar

from py_dataset import dataset

from .mysql_access import get_recently_modified_keys

#
# This python module provides basic functionality for working
# with an EPrints server via the REST API.
#

# Internal data for making calls to EPrints.
c_name = ''
scheme = ''
username = ''
password = ''
netloc = ''
base_url = ''

#
# harvest_init setups the connection information for harvesting EPrints content.
#
def harvest_init(collection_name, connection_string):
    global c_name, scheme, username, password, netloc, base_url
    c_name = collection_name
    o = urlparse(connection_string)
    scheme = o.scheme
    username = o.username
    password = o.password
    netloc = o.netloc
    errors = []
    if username == '':
        errors.append('missing username')
    if password == '':
        errors.append('missing password')
    if scheme == '':
        errors.append('missing url scheme')
    if netloc == '':
        errors.eppend('missing hostname and port')
    if len(errors) > 0:
        return errors.join(', ')
    base_url = connection_string
    return ''
    
def eputil(eprint_url, as_json = True):
    cmd = ['eputil']
    if as_json == True:
        cmd.append('-json')
    cmd.append(eprint_url)
    src, err = '', ''
    with Popen(cmd, stdout = PIPE, stderr = PIPE, encoding = 'utf-8') as proc:
        src = str(proc.stdout.read())
        exit_code = proc.returncode
        if exit_code:
            print(f'DEBUG return code => {exit_code}')
            err = str(proc.stderr.read())
            print(f'DEBUG ({proc.returncode}) err -> {err}')
    return src, err

#
# Get a complete list of keys in the repository
#
def harvest_keys():
    global base_url
    src, err = eputil(f'{base_url}/rest/eprint/')
    if err != '':
        print(f'WARNING: {err}', type(err), file = sys.stderr)
        return []
    keys = json.loads(src)
    return keys


#
# harvest_record fetches a single EPrints record including
# it's EPrintsXML as well as related objects. Store them
# in a dataset collection as attachments.
#
def harvest_record(key):
    global base_url, c_name
    src, err = eputil(f'{base_url}/rest/eprint/{key}.xml')
    if err != '':
        return err
    eprint_xml_object = {}
    obj = {}
    if src != '':
        eprint_xml_object = json.loads(src)
    else:
        return 'No data'
    if 'eprint' in eprint_xml_object:
        if len(eprint_xml_object) > 0:
            obj = eprint_xml_object['eprint'][0]
        else:
            return "Can't find contents of eprint element in EPrintXML"
    else:
        return "Can't find eprint element in EPrintXML"
    key = str(obj['eprint_id'])
    err = ''
    if dataset.has_key(c_name, key):
        ok = dataset.update(c_name, key, obj)
        if not ok:
            return dataset.error_message()
    else:
        ok = dataset.create(c_name, key, obj)
        if not ok:
            return dataset.error_message()
    return ''


#
# harvest_eprintxml retrieves and attaches an EPrintXML document
# for the requested record.
#
def harvest_eprintxml(key):
    global base_url, c_name
    return 'harvest_eprintxml({key}) not implemented'

#
# harvest_documents retrieves and attaches the documents
# continue in the EPrint record and attaches them.
#
def harvest_documents(key):
    global base_url, c_name
    return 'harvest_documents({key}) not implemented.'

#
# harvest takes a number of options and replicates functionality
# from the old `ep` golang program used in the feeds project.
# No parameters are provided then a full harvest of metadata will
# be run. Otherwise the harvest is modified according to the parameters.
#
# The following optional params are supported.
#
#  keys - is a list of numeric eprint ids to be harvested, 
#         an empty means use all keys.
#  start_id - harvest the keys start with given ids (ascending)
#  (number_of_days with db_connection) - are used to harvest records
#             the last number of days by doing a SQL query to generate
#             the list of keys. db_connection holds the db connection
#             string in the form of mysql://USER:PASSWORD@DB_HOST/DB_NAME
#  include_documents - will include harvesting the included EPrint 
#  records' documents
#
def harvest(keys = [], start_id = 0, save_exported_keys = '', number_of_days = None, db_connection = None, include_documents = False):
    global base_url, c_name
    repo_name, _ = os.path.splitext(c_name)
    exported_keys = []
    if len(keys) == 0:
        keys = harvest_keys()
    keys.sort(key=int)
    if start_id > 0:
        new_keys = []
        for key in keys:
            if key >= start_id:
                new_keys.append(key)
        keys = new_keys
    if number_of_days:
        days_since = (datetime.now()+timedelta(days=number_of_days)).strftime('%Y-%m-%d')
        if db_connection:
            keys = get_recently_modified_keys(db_connection, no_of_days)

    tot = len(keys)
    e_cnt = 0
    n = 0
    bar = progressbar.ProgressBar(
            max_value = tot,
            widgets = [
                progressbar.Percentage(), ' ',
                progressbar.Counter(), f'/{tot} ',
                progressbar.AdaptiveETA(),
                f' from {repo_name}'
            ], redirect_stdout=False)
    print(f'harvesting {tot} records from {repo_name}')
    bar.start()
    for i, key in enumerate(keys):
        err = harvest_record(key)
        if err != '':
            print(f'WARNING: harvest record {key}, {err}', file = sys.stderr)
            e_cnt += 1
            continue
        err = harvest_eprintxml(key)
        if err != '':
            print(f'WARNING harvest eprint xml {key}, {err}', file = sys.stderr)
            e_cnt += 1
            continue
        if include_documents:
            err = harvest_documents(key)
            if err != '':
                print(f'WARNING harvest documents {key}, {err}', file = sys.stderr)
                e_cnt += 1
                continue
        exported_keys.append(str(key))
        n += 1
        bar.update(i)
    bar.finish()
    print(f'harvested {n}/{tot} harvested from {repo_name}, {e_cnt} warnings')
    if save_exported_keys != '':
        print(f'saving exported keys to {save_exported_keys}')
        with open(save_exported_keys, 'w') as f:
            src = '\n'.join(exported_keys)
            f.write(src)
    return ''
