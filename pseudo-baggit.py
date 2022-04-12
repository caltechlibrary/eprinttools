#!/usr/bin/env python3

#
# NOTE: This is an example script of how you could use eputil from Python 3.
#
# pseudo-bagger.py will retrieve metadata for specific keys in an
# eprint repository, then use scp to find and retrieve the materials
# listed in the metadata putting them into a directory structure suitable
# for "bagging".  
#
import os
import sys
import json

from subprocess import Popen, PIPE
from urllib.parse import urlparse

import progressbar

USAGE = '''
USAGE:

    pseudo-bagger.py JSON_CONFIG_FILE JSON_KEY_LIST_FILE

This program takes a JSON configuration and a JSON list of keys
and will retrieve metadata via EPrints REST API and then retrieve
documents via scp to the eprint host system. It stores the results
in a eprint id and pos tree under the name of the repository id.

JSON_KEY_LIST_FILE

The key list file is a JSON list of eprint id to be pseudo bagged.

JSON_CONFIG_FILE

Below is an example JSON document. The strings in all caps should
be replaced with appropriate ones for the system you're doing the
pseudo bagging.

    {
        "user": "EPRINT_USER_ID_GOES_HERE",
        "password": "EPRINT_PASSWORD_GOES_HERE",
        "base_url": "https://EPRINT_REPO_HOSTNAME_GOES_WHERE",
        "repo_id": "REPO_ID_GOES_HERE",
        "install_path": "/coda/eprints-3.3"
    }
    
REQUIREMENTS

This application relies on the *eputil* command line program
that comes with eprinttools version 1.1.6 or better.

    https://github.com/caltechlibrary/eprinttools/releases

'''

# Environment variable map mapping our configuration file to support
# eprinttools cli like eputil.
eprinttools_env_map = {
    "user": "EPRINT_USER",
    "password": "EPRINT_PASSWD",
    "repo_id": "EPRINT_REPO_ID",
    "base_url": "EPRINT_BASE_URL",
    "host": "EPRINT_HOST",
    "install_path": "EPRINT_INSTALL_PATH"
}

def load_json(fname):
    '''read and decode a JSON file and make sure it is a list of dict'''
    src = ''
    try:
        with open(fname, encoding = 'utf-8') as f:
            src = f.read()
    except Exception as err:
            return None, err
    try:
        obj = json.loads(src)
    except Exception as err:
        return None, err
    if isinstance(obj, dict):
        if not obj:
            return obj, f'No data read from {fname}'
    if isinstance(obj, list):
        if len(obj) == 0:
            return obj, f'No data read from {fname}'
    if not isinstance(obj, dict) and not isinstance(obj, list):
        return None, f'{fname} should hold a JSON object or array'
    return obj, None

def run(cmd):
    '''run uses Popen to run a command and capture it's stdout and stderr
       returns two values out holds the what was sent to stdout, err holds
       what was sent to stderr. Note some cli use stderr by default'''
    out, err = None, None
    with Popen(cmd, stdout = PIPE, stderr = PIPE) as proc:
        err = proc.stderr.read().strip().decode('utf-8')
        if err == '':
            err = None
        else:
            print(err)
            sys.stdout.flush()
        out = proc.stdout.read().strip().decode('utf-8')
        if out == '':
            out = None
    return out, err


def retrieve_metadata(config, key):
    '''retrieve_metadata retrieves the EPrint metadata for given eprint id returning a touple of
    metadata object and error value.'''
    # setup the the environment for the child processes running
    repo_id = config["repo_id"]
    base_url = config["base_url"]
    username = config["user"]
    password = config["password"]
    out_name = f'{repo_id}/{repo_id}-{key}/data/{key}.json'
    if not config["base_url"].endswith("/rest/eprint"):
        base_url += "/rest/eprint"
    o = urlparse(base_url)
    u = f'{o.scheme}://{username}:{password}@{o.netloc}/rest/eprint/{key}.xml'
    obj, src = {}, ''
    cmd = [ "eputil", "-json", "-o", out_name, u ]
    _, err = run(cmd)
    if err:
        return None, err
    try:
        with open(out_name) as f:
            src = f.read()
    except Exception as err:
        return None, err
    try:
        obj = json.loads(src)
    except Exception as err:
        return None, err
    cmd = [ "curl", "--silent", "-L", "-o", f'{repo_id}/{repo_id}-{key}/data/{repo_id}-{key}.xml', u ]
    _, err = run(cmd)
    if err:
        return None, err
    return obj['eprint'][0], ""

def pairtree(src):
    if isinstance(src, str):
        s = src
    else:
        s = f'{src}'
    # Pairtree's chunks are sized 2
    parts = [ s[i:i+2] for i in range(0, len(s), 2) ]
    return '/'.join(parts)

def retrieve_files(config, key, meta, bar, cnt):
    '''retrieve_files retrieves and saves files indicated via metadata and configuration
       using scp and writing them to a eprint id and pos tree situated in a directory named for the
       repository id.'''
    repo_id = config["repo_id"]
    hostname = config["host"]
    install_path = config["install_path"]
    id_pairtree = pairtree(f'00000000{key}'[-8:])
    errors = []
    # Setup progress bar
    tot = len(meta["documents"])
    pid = os.getpid()
    for i, doc in enumerate(meta["documents"]):
        pos, filename = doc["pos"], doc["main"].strip()
        bar.update(cnt, ky = key, pos = pos, fn = filename)
        dest_filename = filename
        if filename != "":
            pos_pairtree = pairtree(f'00{pos}'[-2:])
            if " " in filename:
                filename = filename.replace(" ", '?')
            if "(" in filename:
                filename = filename.replace("(", '\(')
            if ")" in filename:
                filename = filename.replace(")", '\)')
            scp_source = f'{hostname}:{install_path}/archives/{repo_id}/documents/disk0/{id_pairtree}/{pos_pairtree}/{filename}'
            scp_dest = f'{repo_id}/{repo_id}-{key}/data/{dest_filename}'
            cmd = [ "scp", "-p", scp_source, scp_dest ]
            _, err = run(cmd)
            if err:
                sys.stdout.flush()
                errors.append(err)
    return ("\n".join(errors)).strip()


def pseudo_bag(config, keys):
    '''pseudo_bag takes a configuation dict and key list then attempts
       to pseudo bag the target data'''
    # the eputil command form eprinttools.
    repo_id = config["repo_id"]
    # Setup Progressbar
    tot = len(keys)
    pid = os.getpid()
    bar = progressbar.ProgressBar(
        max_value = tot,
        widgets = [
            progressbar.Percentage(),
            ' ', progressbar.Counter(), f'/{tot}',
            ' ', progressbar.AdaptiveETA(),
            f' {repo_id} pid:{pid}',
            ' ', progressbar.Variable('ky', width = 1),
            ' ', progressbar.Variable('pos', width = 1),
            ' ', progressbar.Variable('fn', width = 1),
        ], redirect_stdout = True)
    bar.start()
    ok = True
    for i, key in enumerate(keys):
        # Update progress bar
        bar.update(i, ky = key, pos = 0, fn = '')
        os.makedirs(f'{repo_id}/{repo_id}-{key}/data', 0o777, True)
        meta, err = retrieve_metadata(config, key)
        if err:
            print(err)
            ok = False
        else:
            err = retrieve_files(config, key, meta, bar, i)
            if err:
                print(err)
                ok = False
    # Cleanup from progressbar
    bar.finish()
    return ok


if __name__ == '__main__':
    if len(sys.argv) != 3:
        print(USAGE)
        sys.exit(1)
    config_file, key_file = sys.argv[1],sys.argv[2]
    config, err = load_json(config_file)
    if err != None:
        print(err)
        sys.exit(1)
    keys, err = load_json(key_file)
    if err != None:
        print(err)
        sys.exit(1)
    if not pseudo_bag(config, keys):
        sys.exit(1)
    sys.exit(0)
