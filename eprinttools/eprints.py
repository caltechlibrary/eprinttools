#!/usr/bin/env python3
import os
import sys
import json
from urllib.parse import urlparse
from subprocess import run, Popen, PIPE
from datetime import datetime, timedelta

import progressbar

from py_dataset import dataset

#from .mysql_access import get_recently_modified_keys

#
# This python module provides basic functionality for working
# with an EPrints server via the REST API.
#
skip_and_prune = [ "deletion", "buffer", "inbox" ]

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
        errors.append('missing hostname and port')
    if len(errors) > 0:
        return ', '.join(errors)
    base_url = connection_string
    return ''
    
def eputil(eprint_url, as_json = True, get_document = False, as_text = False):
    cmd = ['eputil']
    if as_text == True:
        cmd.append('-raw')
    if as_json == True:
        cmd.append('-json')
    if get_document == True:
        cmd.append('-document')
    cmd.append(eprint_url)
    src, err = '', ''
    with Popen(cmd, stdout = PIPE, stderr = PIPE, encoding = 'utf-8') as proc:
        err = proc.stderr.read().strip()
        if err != '':
            print(f"\n{' '.join(cmd)}\nerror: {err}")
        src = str(proc.stdout.read())
    return src, err

#
# Get a complete list of keys in the repository
#
def harvest_keys():
    global base_url
    src, err = eputil(f'{base_url}/rest/eprint/', as_json = True)
    if err != '':
        print(f'WARNING: {err}', type(err), file = sys.stderr)
        return []
    if not isinstance(src, bytes):
        src = src.encode('utf-8')
    keys = json.loads(src)
    return keys


#
# harvest_record fetches a single EPrints record including
# it's EPrintsXML as well as related objects. Store them
# in a dataset collection as attachments.
#
def harvest_record(key, verbose = False):
    global base_url, c_name
    obj = {}
    src, err = eputil(f'{base_url}/rest/eprint/{key}.xml')
    if err != '':
        return None, err
    eprint_xml_object = {}
    if src != '':
        eprint_xml_object = json.loads(src)
    else:
        return None,'No data'
    if 'eprint' in eprint_xml_object:
        if len(eprint_xml_object) > 0:
            obj = eprint_xml_object['eprint'][0]
        else:
            return None, "Can't find contents of eprint element in EPrintXML"
    else:
        return None, "Can't find eprint element in EPrintXML"
    key = str(obj['eprint_id'])
    err = ''
    status = ''
    if 'eprint_status' in obj:
        status = obj['eprint_status']
    if dataset.has_key(c_name, key):
        if status in skip_and_prune:
            if verbose:
                print(f'''
WARNING: Removing {key} from {c_name}, reason {status}''')
            ok = dataset.delete(c_name, key)
            if not ok:
                return None, dataset.error_message()
            return None, ''
        ok = dataset.update(c_name, key, obj)
        if not ok:
            return None, dataset.error_message()
        return obj, ''
    if status in skip_and_prune:
        if verbose:
            print(f'''
WARNING: Skipping {key} from {c_name}, reason {status}''')
        return None, ''
    ok = dataset.create(c_name, key, obj)
    if not ok:
        return None, dataset.error_message()
    return obj, ''


#
# harvest_eprintxml retrieves and attaches an EPrintXML document
# for the requested record.
#
def harvest_eprintxml(key):
    global base_url, c_name
    key = str(key)
    # Fetch the EPrintXML document source
    src, err = eputil(f'{base_url}/rest/eprint/{key}.xml', as_json = False)
    if err != '':
        return err
    # Write to temp file
    f_name = f'{key}.xml'
    with open(f_name, 'w') as f:
        f.write(src)
    # Attach file to record in dataset collection
    if dataset.has_key(c_name, key) == False:
        ok = dataset.create(c_nane, key, {})
        if not ok:
            return dataset.error_message()
    ok = dataset.attach(c_name, key, [ f_name ])
    if not ok:
        return dataset.error_message()
    os.remove(f_name)
    return ''

#
# harvest_documents retrieves and attaches the documents
# continue in the EPrint record and attaches them.
#
def harvest_documents(key, obj):
    global base_url, c_name
    key = str(key)
    if 'primary_object' in obj:
        doc_url = obj['primary_object']['url']
        f_name = obj['primary_object']['basename']
        semver = obj['primary_object']['version']
        src, err = eputil(doc_url, get_document = True)
        if err != '':
            return err
        ok = dataset.attach(c_name, key, [ f_name ], semver)
        if not ok:
            return dataset.error_message()
        os.remove(f_name)
    if 'related_objects' in obj:
        for o in obj['related_objects']:
            doc_url = o['url']
            f_name = o['basename']
            semver = o['version']
            src, err = eputil(doc_url, get_document = True)
            if err != '':
                return err
            ok = dataset.attach(c_name, key, [ f_name ], semver)
            if not ok:
                return dataset.error_message()
            os.remove(f_name)
    return ''

#
# filter_recently_modified_keys takes a list of keys then
# checks the value of lastmod.txt (e.g. rest/eprint/EPRINT_ID/lastmod.txt)
# and filters the key list it returns.
#
def filter_recently_modified_keys(keys, date_since):
    global base_url, c_name
    repo_name, _ = os.path.splitext(c_name)
    filter_keys = []
    tot = len(keys)
    bar = progressbar.ProgressBar(
            max_value = tot,
            widgets = [
                progressbar.Percentage(), ' ',
                progressbar.Counter(), f'/{tot} ',
                progressbar.AdaptiveETA(),
                f' from {repo_name}'
            ], redirect_stdout=False)
    print(f'filtering {tot} modified since {date_since} from {repo_name}')
    bar.start()
    for i, key in enumerate(keys):
        lastmod_url = f'{base_url}/rest/eprint/{key}/lastmod.txt'
        src, err = eputil(lastmod_url, as_text = True)
        if err != '':
            print(f'WARNING: could not get lastmod date for {key}, {err}')
            continue
        if src != '':
            src = src.strip()
            day, hour = src.split(" ")
            if day >= date_since:
                filter_keys.append(key)
        bar.update(i)
    bar.finish()
    print(f'filtered {len(filter_keys)} of {tot} modified since {date_since} from {repo_name}')
    return filter_keys
    
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
def harvest(keys = [], start_id = 0, save_exported_keys = '', number_of_days = None, db_connection = None, include_documents = False, verbose = True):
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
        # NOTE: we want an integer value for the last number of days.
        number_of_days = int(number_of_days)
        if number_of_days > 0:
            number_of_days = number_of_days * -1
        date_since = (datetime.now()+timedelta(days=number_of_days)).strftime('%Y-%m-%d')
        keys = filter_recently_modified_keys(keys, date_since)

    tot = len(keys)
    e_cnt = 0
    pruned = 0
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
        obj, err = harvest_record(key, verbose)
        if err != '':
            print(f'''
WARNING: harvest record {key}, {err}''', file = sys.stderr)
            e_cnt += 1
            continue
        # NOTE: If object is None then the record was be skipped or 
        # pruned. This is not an error but reflects EPrints behavior.
        if obj == None:
            pruned += 1
            bar.update(i)
            continue
        err = harvest_eprintxml(key)
        if err != '':
            print(f'''
WARNING harvest eprint xml {key}, {err}''', file = sys.stderr)
            e_cnt += 1
            bar.update(i)
            continue
        if include_documents:
             err = harvest_documents(key, obj)
             if err != '':
                 print(f'''
WARNING harvest documents {key}, {err}''', file = sys.stderr)
                 e_cnt += 1
                 bar.update(i)
                 continue
        if save_exported_keys:
            exported_keys.append(str(key))
        n += 1
        bar.update(i)
    bar.finish()
    print(f'harvested {n}/{tot}, skipped/pruned {pruned} from {repo_name}, {e_cnt} warnings')
    if save_exported_keys != '':
        print(f'saving exported keys to {save_exported_keys}')
        with open(save_exported_keys, 'w') as f:
            src = '\n'.join(exported_keys)
            f.write(src)
    return ''
