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

from eprintviews import Views, Aggregator
from eprintsubjects import Subjects
from eprintusers import Users

#
# Common EPrint Site Layouts look like:
#
# Home -> index.md
# About -> information.md
# Browse -> browserviews.md
#   Person -> people_list.json -> view/person-az/index.md
#   Person -> people_list.json -> view/person/index.md
#   Author -> people_list.json -> view/author/index.md
#   Year -> year_list.json -> view/year/index.md
#   Subject -> subject_list.json -> view/subject/index.md
#
# Unlinked types include view/ids/ and view/types/
#
# Simple Search and Advanced Search -> Elasticsearch services
# Contact Us -> redirects to https://www.library.caltech.edu/contact
#

def get_title(obj):
    if 'title' in obj:
        return obj['title']
    return ''

def get_date(obj):
    if 'date' in obj:
        return obj['date'][0:4].strip()
    return ''

def get_eprint_id(obj):
    if 'eprint_id' in obj:
        return f"{obj['eprint_id']}"
    return ''

def get_creator_id(creator):
    if 'id' in creator:
        return creator['id']
    return ''

def get_creator_name(creator):
    family, given = '', ''
    if 'name' in creator:
        if 'family' in creator['name']:
            family = creator['name']['family'].strip()
        if 'given' in creator['name']:
            given = creator['name']['given'].strip()
    if len(family) > 0:
        if len(given) > 0:
            return f'{family}, {given}'
        return family
    return ''

def make_creator_list(creators):
    l = []
    for creator in creators:
        display_name = get_creator_name(creator)
        creator_id = get_creator_id(creator)
        creator['display_name'] = display_name
        creator['creator_id'] = creator_id
        l.append(creator)
    return l        

def make_label(s, sep = '_'):
    l = s.split(sep)
    for i, val in enumerate(l):
        l[i] = val.capitalize()
    return ' '.join(l)

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
    objs.sort(reverse = True, key = get_date) # primary sort value
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
    ok = dataset.frame_create(c_name, frame_name, keys, [ '.eprint_id', '.title', '.date', '.creators', '.subjects', '.type', '.official_url', '.userid' , '.local_group', '.issn', '.collection' ], [ 'eprint_id', 'title', 'date', 'creators', 'subjects', 'type', 'official_url', 'userid', 'groups', 'issn', 'collection' ])
    if not ok:
        err = dataset.error_message()
        return f'{frame_name} in {c_name}, not created, {err}'
    return ''

def normalize_object(obj, users):
    title = obj['title'].strip()
    year = get_date(obj)
    eprint_id = get_eprint_id(obj)
    creator_list = make_creator_list(obj['creators']['items'])
    if 'userid' in obj:
        key = f'{obj["userid"]}'
        if users.has_user(key):
            user = users.get_user(key)
            if 'display_name' in user:
                display_name = user['display_name']
                obj['depositor'] = display_name
    if 'date' in obj:
        obj['year'] = year
    if 'type' in obj:
        obj['type_label'] = make_label(obj['type'])
    obj['title'] = title
    obj['creators'] = creator_list
    obj['eprint_id'] = eprint_id
    obj['year'] = year
    return obj

def normalize_objects(objs, users):
    for obj in objs:
        obj = normalize_object(obj, users)
    return objs

#
# Build our aggregated views
#
# After looking across our repositories we have the following aggregations
# some shared between views (e.g. author, person, person-az are all
# people views).  
#
# Example usage:
#
#   aggregation = aggregate(c_nane, users, views, subjects)
#
# The supported aggregations are:
#
#    ids, people, years, subjects, types, groups,
#    options, committee, publication,
#    issn, collection
#
#
# FIXME: this should really be handing back a dict of aggregations!
#
def aggregate(c_name, views, users, subjects):
    err = make_frame_date_title(c_name)
    if err != '':
        print(f'{err}')
    frame_name = 'date-title'
    objs = dataset.frame_objects(c_name, frame_name)
    objs = normalize_objects(objs, users)
    aggregator = Aggregator(c_name, objs)
    aggregation = {}
    print(f'DEBUG view keys -> {views.get_keys()}')
    if views.has_view('ids'):
        aggregation['ids'] = aggregator.aggregate_ids()
    if views.has_view('person'):
        aggregation['person'] = aggregator.aggregate_person()
    if views.has_view('person-az'):
        aggregation['person-az'] = aggregator.aggregate_person_az()
    if views.has_view('author'):
        aggregation['author'] = aggregator.aggregate_person_az()
    if views.has_view('year'):
        aggregation['year'] = aggregator.aggregate_year()
    if views.has_view('subjects'):
        aggregation['subjects'] = aggregator.aggregate_subjects(subjects)
    if views.has_view('types'):
        aggregation['types'] = aggregator.aggregate_types()
    if views.has_view('group'):
        aggregation['group'] = aggregator.aggregate_group()
    if views.has_view('collection'):
        aggregation['collection'] = aggregator.aggregate_collection()
    if views.has_view('latest'):
        aggregation['latest'] = aggregator.aggregate_latest()
    if views.has_view('issn'):
        aggregation['issn'] = aggregator.aggregate_issn()
    if views.has_view('publication'):
        aggregation['publication'] = aggregator.aggregate_publication()
    return aggregation


#
# Build the directory tree
#
def generate_directories(views):
    doc_tree = {}
    for key in views.get_keys():
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
def generate_landings(c_name, users):
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
                shutil.move(b_name, f_name, copy_function = shutil.copy2)
        bar.update(i)
    bar.finish()
    print(f'generated {tot} landing pages, {e_cnt} errors from {repo_name}')


def make_view(view_name, p_name, key_field, object_list):
    print(f'generating {view_name} view')
    if not os.path.exists(p_name):
        os.makedirs(p_name, mode = 0o777, exist_ok = True)
    if not os.path.exists(p_name):
        print(f'WARNING {p_name} does not exist, skipping {view_name}')
        return ''
    f_name = os.path.join(p_name, f'{view_name}_list.json')
    with open(f_name, 'w') as f:
        src = json.dumps(object_list)
        f.write(src)
    for obj in object_list:
        if not key_field in obj:
            print(f'DEBUG missing key field {key_field} in {obj}')
        else:
            key = obj[key_field]
            f_name = os.path.join(p_name, f'{key}.json')
            with open(f_name, 'w') as f:
                src = json.dumps(obj)
                f.write(src)

# generate_views creates the common views across our EPrints
# repositories.
def generate_views(views, aggregations):
    # map v_name to field id in EPrint XML
    field_ids = {
        'ids': 'eprint_id',
        'people': 'people_id',
        'person-az': 'people_id',
        'person': 'people_id',
        'year': 'year',
        'subjects': 'subject_id',
        'types': 'type',
        'group': 'local_group',
        'option': 'option',
        'advisor': 'thesis_advisor_id',
        'committee': 'thesis_committee_id',
        'collection': 'collection',
        'issn': 'issn',
        'publication': 'publication',
        'event': 'event_id'
    }
    # Commonly Linked views
    for v_name in views.get_keys():
        print(f'DEBUG v_name -> has_view("{v_name}") -> ', end = '')
        print(views.has_view(v_name)) # DEBUG
        if (v_name in [ 'year' ]) and ('year' in aggregations):
            for year in aggregations[v_name]:
                object_list = year['objects']
                label = year['label']
                field_id = 'year'
                p_name = os.path.join('htdocs', 'view', 'year')
                make_view(label, p_name, field_id, object_list)
        elif (v_name in [ 'collection' ]) and ('collection' in aggregations):
            for collection in aggregations[v_name]:
                object_list = collection['objects']
                label = collection['label']
                field_id = 'collection'
                p_name = os.path.join('htdocs', 'view', 'collection')
                make_view(label, p_name, field_id, object_list)
        elif (v_name in [ 'event' ]) and ('event' in aggregations):
            for event in aggregations[v_name]:
                object_list = event['objects']
                label = event['label']
                field_id = 'event_title'
                p_name = os.path.join('htdocs', 'view', 'event')
                make_view(label, p_name, field_id, object_list)
        elif (v_name in aggregations):
            for view in aggregations[v_name]:
                object_list = view['objects']
                label = view['label']
                field_id = field_ids[v_name]
                p_name = os.path.join('htdocs', 'view', v_name)
                make_view(label, p_name, field_id, object_list)
        else:
            print(f'''WARNING: missing aggregation for "{v_name}"''')


def generate_metadata_structure(c_name, views_json, users_json, subjects_json):
    views = Views()
    views.load_views(views_json)
    users = Users()
    users.load_users(users_json)
    subjects = Subjects()
    subjects.load_subjects(subjects_json)
    generate_directories(views)
    # NOTE: To keep the code simple we're only supporting
    # A default list of views found in common with Caltech Library's
    # repositories. If this gets used by other libraries we can
    # look at generalizing this approach.
    aggregations = aggregate(c_name, views, users, subjects)
    print(f'Found {len(aggregations)} aggregations:')
    for key in aggregations:
        print(f'  {len(aggregations[key])} {key}')
    print('')
    generate_views(views, aggregations)
    generate_landings(c_name, users)


if __name__ == "__main__":
    f_name = 'config.json'
    c_name = ''
    views_json = ''
    users_json = ''
    subjects_json = ''
    if len(sys.argv) > 1:
        f_name = sys.argv[1]
    if not os.path.exists(f_name):
        print(f'Missing {f_name} configuration file')
        sys.exit(1)
    with open(f_name) as f:
        src = f.read()
        cfg = json.loads(src)
        if 'dataset' in cfg:
            c_name = cfg['dataset']
        if 'views' in cfg:
            views_json = cfg['views']
        if 'users' in cfg:
            users_json = cfg['users']
        if 'subjects' in cfg:
            subjects_json = cfg['subjects']
    generate_metadata_structure(c_name, views_json, users_json, subjects_json)
    print('OK')
