#!/usr/bin/env python3

#
# This program generates the htdocs directory structure, and
# metadata needed to replicate EPrint's various views of repository
# content.
#

import os
import sys
import shutil
import json
import string

from urllib.parse import urlparse

import progressbar

from py_dataset import dataset

from eprintviews import Aggregator, Views, Subjects, Users, normalize_object, get_date_year, get_eprint_id, get_title

#
# CaltechES EPrint Site Layouts look like:
#
# Home -> index.html
# About -> information.html
# Browse -> browserviews.html
#   Year -> year_list.json -> view/year/
#   Item Category -> subjects_list.json > view/subject/
#   Author -> people_list.json -> view/person-az/
#   Latest Additions -> latest_list.json -> cgi/latest/
#
# Unlinked types include view/ids/ and view/types/
#
#   Eprint ID -> ids_list.json -> view/ids/index.md
#   Year -> year_list.json -> view/year/
#   Item Category -> subjects_list.json -> view/subject/
#   Person -> person_list.json -> /view/person/
#   Author -> people_list.json -> view/person-az/
#   Type -> types_list.json -> view/types/index.md
#
# Simple Search and Advanced Search -> Elasticsearch services
# Deposit an Item -> http://calteches.library.caltech.edu/cgi/users/home
# Contact Us (broken in produciton) -> should redirect to https://www.library.caltech.edu/contact


def make_frame_date_title(c_name):
    frame_name = 'date-title-unsorted'
    if dataset.has_frame(c_name, frame_name):
        ok = dataset.delete_frame(c_name, frame_name)
        if not ok:
            err = dataset.error_message()
            return f'{frame_name} in {c_name} not deleted, {err}'
    keys = dataset.keys(c_name)
    for i, key in enumerate(keys):
        keys[i] = f'{key}'
    print(f'creating frame {frame_name} in {c_name}')
    ok = dataset.frame_create(c_name, frame_name, keys, ['.eprint_id', '.title', '.date' ], [ 'eprint_id', 'title', 'date' ])
    if not ok:
        err = dataset.error_message()
        return f'{frame_name} in {c_name} not created, {err}'
    objs = dataset.frame_objects(c_name, frame_name) 
    if len(objs) == 0:
        return f'{frame_name} in {c_name}, missing objects'
    print(f'sorting {len(objs)} in frame {frame_name} in {c_name}')
    # Sort the objects by date then title
    objs.sort(key = get_title) # secondary sort value
    objs.sort(reverse = True, key = get_date_year) # primary sort value
    keys = []
    for obj in objs:
        keys.append(get_eprint_id(obj))
    frame_name = 'date-title'
    # Now save the sorted frame
    if dataset.has_frame(c_name, frame_name):
        ok = dataset.delete_frame(c_name, frame_name)
        if not ok:
            err = dataset.error_message()
            return f'{frame_name} in {c_name}, not deleted, {err}'
    print(f'creating frame {frame_name} in {c_name}')
    # NOTE: this frame needs to have the fields you'll use to
    # generate the repository's views.
    ok = dataset.frame_create(c_name, frame_name, keys, 
            [ '.eprint_id', '.title', '.date', '.creators', '.subjects', 
                '.type', '.official_url', '.userid' , 
                '.collection', '.event_title', '.event_location', '.event_dates' ], 
            [ 'eprint_id', 'title', 'date', 'creators', 'subjects', 
                'type', 'official_url', 'userid', 
                'collection', 'event_title', 'event_location', 'event_dates' ])
    if not ok:
        err = dataset.error_message()
        return f'{frame_name} in {c_name}, not created, {err}'
    return ''

def normalize_objects(objs, users):
    for obj in objs:
        obj = normalize_object(obj, users)
    return objs

#
# Build our this repository's aggregated views
#
def aggregate(c_name, views, users, subjects):
    err = make_frame_date_title(c_name)
    if err != '':
        print(f'{err}')
    frame_name = 'date-title'
    aggregations = {}
    objs = dataset.frame_objects(c_name, frame_name)
    objs = normalize_objects(objs, users)
    aggregator = Aggregator(c_name, objs)
    view_keys = views.get_keys()
    for key in view_keys:
        aggregations[key] = aggregator.aggregate_by_view_name(key, subjects)
    return aggregations


#
# Build the directory tree
#
def generate_directories(view_paths):
    doc_tree = {}
    for key in view_paths:
        doc_tree[key] = os.path.join('htdocs', 'view', key)
    for key in doc_tree:
        d_name = doc_tree[key]
        if not os.path.exists(d_name):
            os.makedirs(d_name, mode = 0o777, exist_ok = True)
    

#
# landing_filter is used to transform EPrint Objects into something
# friendly to use with Pandoc.
#
def landing_filter(obj, users):
    return normalize_object(obj, users)


#
# generate_landings creates index.json to render index.md,
# also deposits attachments in their relative paths.
#
def generate_landings(c_name, views, users, subjects, include_documents = False):
    repo_name, _ = os.path.splitext(c_name)
    keys = dataset.keys(c_name)
    keys.sort(key=int)
    tot = len(keys)
    e_cnt = 0
    bar = progressbar.ProgressBar(
            max_value = tot,
            widgets = [
                progressbar.Percentage(), ' ',
                progressbar.Counter(), f'/{tot} ',
                progressbar.AdaptiveETA(),
                f' from {repo_name}'
            ], redirect_stdout=False)
    f_name = os.path.join('htdocs', 'index.keys')
    with open(f_name, 'w') as f:
        for key in keys:
            f.write(f'{key}\n')
    print(f'generating {tot} landing pages for {repo_name}')
    for i, key in enumerate(keys):
        obj, err = dataset.read(c_name, key)
        if err != '':
            e_cnt += 1
            print(f'''
WARNING: can't read {key} from {c_name}, {err}''')
            continue
        src = json.dumps(landing_filter(obj, users))
        p_name = os.path.join('htdocs', f'{key}')
        os.makedirs(p_name, mode = 0o777, exist_ok = True)
        f_name = os.path.join(p_name, 'index.json')
        with open(f_name, 'w') as f:
            f.write(src)
        #FIXME: We want to have the option of including attachments
        # for the digital files in our collection OR copying from
        # source location to S3 bucket!
        if include_documents:
            # NOTE: we need to copy the attachments into the correct place
            # in our htdocs tree.
            if 'primary_object' in obj:
                b_name = obj['primary_object']['basename']
                semver = obj['primary_object']['version']
                url = obj['primary_object']['url']
                o = urlparse(url)
                p_name = os.path.join('htdocs', 
                         os.path.dirname(o.path).lstrip('/'))
                if not os.path.exists(p_name):
                    os.makedirs(p_name, mode = 0o777, exist_ok = True)
                f_name = os.path.join(p_name, b_name)
                ok = dataset.detach(c_name, key, [ b_name ], semver)
                if not ok:
                    err = dataset.error_message()
                    print(f'''
WARNING could not detach {b_name} in {key} from {c_name}, {err}')''')
                else:
                    # Do final sanity check before copy.
                    if os.path.exists(b_name):
                        shutil.move(b_name, f_name, copy_function = shutil.copy2)
                    else:
                        print(f'''
WARNING detached file missing {b_name} in {key} from {c_name}
cannot move to {f_name}, skipping''')
        bar.update(i)
    bar.finish()
    print(f'generated {tot} landing pages, {e_cnt} errors from {repo_name}')


def make_view(view, p_name, aggregation):
    if not os.path.exists(p_name):
        os.makedirs(p_name, mode = 0o777, exist_ok = True)
    if not os.path.exists(p_name):
        print(f'WARNING {p_name} does not exist, skipping {view}')
        return ''
    f_name = os.path.join(p_name, f'{view}_list.json')
    print(f'writing "{view}" -> {f_name}')
    with open(f_name, 'w') as f:
        src = json.dumps(aggregation)
        f.write(src)
    for obj in aggregation:
        objects = obj['objects']
        key = obj['key']
        f_name = os.path.join(p_name, f'{key}.json')
        #print(f'writing "{view}" ({key}) -> {f_name}')
        with open(f_name, 'w') as f:
            src = json.dumps(objects)
            f.write(src)

def generate_view(key, aggregations):
    if key in aggregations:
        p_name = os.path.join('htdocs', 'view', key)
        if (aggregations[key] != None) and (len(aggregations[key]) > 0):
            make_view(key, p_name, aggregations[key])
        else:
            make_view(key, p_name, [])
    
# generate_views creates the common views across our EPrints
# repositories.
def generate_views(views, aggregations):
    keys = views.get_keys()
    # /view/ views and subviews, may also be included in browseviews.html
    for key in keys:
        generate_view(key, aggregations)


def generate_metadata_structure(c_name, f_views, f_users, f_subjects, include_documents = False):
    views = Views()
    views.load_views(f_views)
    users = Users()
    users.load_users(f_users)
    subjects = Subjects()
    subjects.load_subjects(f_subjects)
    generate_directories(views.get_keys())
    aggregations = aggregate(c_name, views, users, subjects)
    print(f'Found {len(aggregations)} aggregations: ', end = '\n\t')
    for i, key in enumerate(aggregations):
        if i > 0:
            print(', ', end = '')
        if (key in aggregations) and aggregations[key] != None:
            print(f'{len(aggregations[key])} {key}', end = '')
        else:
            print(f'Nothing to aggregate for {key}')
    print('')
    generate_views(views, aggregations)
    generate_landings(c_name, views, users, subjects, include_documents)


if __name__ == "__main__":
    f_name = ''
    c_name = ''
    if len(sys.argv) > 1:
        f_name = sys.argv[1]
    if not os.path.exists(f_name):
        print(f'Missing JSON configuration file')
        sys.exit(1)
    with open(f_name) as f:
        src = f.read()
        cfg = json.loads(src)
        if 'dataset' in cfg:
            c_name = cfg['dataset']
        if 'users' in cfg:
            f_users = cfg['users']
        if 'views' in cfg:
            f_views = cfg['views']
        if 'subjects' in cfg:
            f_subjects = cfg['subjects']
    generate_metadata_structure(c_name, f_views, f_users, f_subjects, include_documents = False) 
    print('OK')
