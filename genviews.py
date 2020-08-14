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

from eprints3x import load_subjects, normalize_subject, load_users, get_user, has_user, load_views, has_view

#
# Common EPrint Site Layouts look like:
#
# Home -> index.md
# About -> information.md
# Browse -> browserviews.md
#   Person -> people_list.json -> view/person-az/index.md
#   Year -> year_list.json -> view/year/index.md
#   Subject -> subject_list.json -> view/subject/index.md
#
# Unlinked types include view/ids/ and view/types/
#
# Simple Search and Advanced Search -> Elasticsearch services
# Contact Us -> redirects to https://www.library.caltech.edu/contact
#
doc_tree = {}

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

def get_object_type(obj):
    if 'type' in obj:
        return f'{obj["type"]}'
    return ''

def has_creator_ids(obj):
    for creator in obj['creators']:
        if 'id' in creator:
            return True
    return False

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
    ok = dataset.frame_create(c_name, frame_name, keys, ['.eprint_id', '.title', '.date' ], [ 'eprint_id', 'title', 'date' ])
    if not ok:
        err = dataset.error_message()
        return f'{frame_name} in {c_name} not created, {err}'
    objs = dataset.frame_objects(c_name, frame_name) 
    if len(objs) == 0:
        return f'{frame_name} in {c_name}, missing objects'
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
    ok = dataset.frame_create(c_name, frame_name, keys, [ '.eprint_id', '.title', '.date', '.creators', '.subjects', '.type', '.official_url', '.userid' ], [ 'eprint_id', 'title', 'date', 'creators', 'subjects', 'type', 'official_url', 'userid' ])
    if not ok:
        err = dataset.error_message()
        return f'{frame_name} in {c_name}, not created, {err}'
    return ''

def get_sort_name(o):
    if 'sort_name' in o:
        return o['sort_name']
    return ''

def get_sort_year(o):
    if 'year' in o:
        return o['year']
    return ''

def get_sort_subject(o):
    if 'subject_name' in o:
        return o['subject_name']
    return ''

def normalize_object(obj):
    title = obj['title'].strip()
    year = get_date(obj)
    eprint_id = get_eprint_id(obj)
    creator_list = make_creator_list(obj['creators']['items'])
    if 'userid' in obj:
        key = f'{obj["userid"]}'
        if has_user(key):
            user = get_user(key)
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

def normalize_objects(objs):
    for obj in objs:
        obj = normalize_object(obj)
    return objs

def aggregate_people(c_name, objs):
    # now build our people list and create a people, eprint_id, title list
    people = {}
    for obj in objs:
        if has_creator_ids(obj):
            # For each author add a reference to object
            for creator in obj['creators']:
                creator_id = creator['id']
                creator_name = creator['display_name']
                if not creator_id in people:
                    people[creator_id] = { 
                        'key': creator_id,
                        'label': creator_name,
                        'count' : 0,
                        'people_id': creator_id,
                        'sort_name': creator_name,
                        'objects' : []
                    }
                people[creator_id]['count'] += 1
                people[creator_id]['objects'].append(obj)
    # Now that we have a people list we need to sort it by name
    people_list = []
    for key in people:
        people_list.append(people[key])
    people_list.sort(key = get_sort_name)
    if len(people_list) == 0:
        print(f'Failed to aggregate any people from {c_name}')
        sys.exit(1)
    return people_list

def aggregate_year(c_name, objs):
    years = {}
    year = ''
    for obj in objs:
        if ('date' in obj):
            year = obj['date'][0:4].strip()
            if not year in years:
                years[year] = { 
                    'key': str(year),
                    'label': str(year),
                    'count': 0,
                    'year': year, 
                    'objects': [] 
                }
            years[year]['count'] += 1
            years[year]['objects'].append(obj)
    year_list = []
    for key in years:
        year_list.append(years[key])
    year_list.sort(key = get_sort_year, reverse = True)
    if len(year_list) == 0:
        print(f'Failed to aggregate any years from {c_name}')
        sys.exit(1)
    return year_list

def aggregate_subjects(c_name, objs):
    subjects = {}
    subject = ''
    for obj in objs:
        eprint_id = get_eprint_id(obj)
        year = get_date(obj)

        if ('subjects' in obj):
            for subj in obj['subjects']['items']:
                subject_name = normalize_subject(subj)
                if subject_name != '':
                    if not subj in subjects:
                        subjects[subj] = { 
                            'key': subj,
                            'label': subject_name,
                            'count': 0,
                            'subject_id': subj, 
                            'subject_name': subject_name,
                            'objects': [] }
                    subjects[subj]['count'] += 1
                    subjects[subj]['objects'].append(obj)
    subject_list= []
    for key in subjects:
        subject_list.append(subjects[key])
    subject_list.sort(key = get_sort_subject)
    if len(subject_list) == 0:
        print(f'Failed to aggregate any subjects from {c_name}')
        sys.exit(1)
    return subject_list

def aggregate_ids(c_name, objs):
    ids = {}
    for obj in objs:
        eprint_id = get_eprint_id(obj)
        if not eprint_id in ids:
            ids[eprint_id] = {
                'key': eprint_id,
                'label': eprint_id,
                'eprint_id': eprint_id,
                'count': 0,
                'objects': []
            }
        ids[eprint_id]['count'] += 1
        ids[eprint_id]['objects'].append(obj)
    ids_list = []
    for key in ids:
        ids_list.append(ids[key])
    ids_list.sort(key = lambda x: int(x['key']))
    if len(ids_list) == 0:
        print(f'Failed to aggregate Eprint ID from {c_name}')
        sys.exit(1)
    return ids_list

def aggregate_types(c_name, objs):
    types = {}
    for obj in objs:
        o_type = get_object_type(obj)
        label = make_label(o_type)
        if not o_type in types:
            types[o_type] = {
                'key': o_type,
                'label': label,
                'type': o_type,
                'count': 0,
                'objects': []
            }
        types[o_type]['count'] += 1
        types[o_type]['objects'].append(obj)
    type_list = []
    for o_type in types:
        type_list.append(types[o_type])
    type_list.sort(key = lambda x: x['key'])
    if len(type_list) == 0:
        print(f'Failed to aggregate types from {c_name}')
        sys.exit(1)
    return type_list

#
# Build our aggregated views
#
def aggregate(c_name):
    frame_name = 'date-title'
    if not dataset.has_frame(c_name, frame_name):
        generate_frames(c_name)
    objs = dataset.frame_objects(c_name, frame_name)
    objs = normalize_objects(objs)
    ids = aggregate_ids(c_name, objs)
    types = aggregate_types(c_name, objs)
    people = aggregate_people(c_name, objs)
    years = aggregate_year(c_name, objs)
    subjects = aggregate_subjects(c_name, objs)
    return ids, people, years, subjects, types


#
# Using view_list populate our doc_type
#
def generate_doctree(view_list):
    global doc_tree
    doc_tree = {}
    for key in view_list:
        doc_tree[key] = os.path.join('htdocs', 'view', key)


#
# Build the directory tree
#
def generate_directories(tree):
    for view in tree:
        d_name = tree[view]
        if not os.path.exists(d_name):
            os.makedirs(d_name, mode = 0o777, exist_ok = True)
    
#
# generate_frames makes the frames used to create the various views
# used to populate the doc tree.
#
def generate_frames(c_name):
    err = make_frame_date_title(c_name)
    if err != '':
        print(f'{err}')

#
# landing_filter is used to transform EPrint Objects into something
# friendly to use with Pandoc.
#
def landing_filter(obj):
    return normalize_object(obj)


#
# generate_landings creates index.json to render index.md,
# also deposits attachments in their relative paths.
#
def generate_landings(c_name):
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
        src = json.dumps(landing_filter(obj))
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
        print(f'WARNING {p_name} does not exist, skipping {view_name}')
        return ''
    f_name = os.path.join(p_name, f'{view_name}_list.json')
    with open(f_name, 'w') as f:
        src = json.dumps(object_list)
        f.write(src)
    for obj in object_list:
        key = obj[key_field]
        f_name = os.path.join(p_name, f'{key}.json')
        with open(f_name, 'w') as f:
            src = json.dumps(obj)
            f.write(src)


def generate_views(doc_type, ids, people, years, subjects, types):
    global doc_tree
    # Linked views
    if 'people' in doc_tree:
        make_view('people', doc_tree['people'], 'people_id', people)
    if 'year' in doc_type:
        make_view('year', doc_tree['years'], 'year', years)
    if 'subject' in doc_tree:
        make_view('subject', doc_tree['subjects'], 'subject_id', subjects)
    # possibily unlinked views
    if 'ids' in doc_tree:
        make_view('ids', doc_tree['ids'], 'eprint_id', ids)
    if 'type' in doc_tree:
        make_view('type', doc_tree['types'], 'type', types)
    

def generate_metadata_structure(c_name, userlist_json, doc_tree):
    load_users(userlist_json)
    generate_frames(c_name)
    generate_directories(doc_tree)
    ids, people, years, subjects, types = aggregate(c_name)
    print(f'Found {len(ids)} ids, {len(people)} people, {len(years)} years, {len(subjects)} subjects, {len(types)} types')
    generate_views(doc_tree, ids, people, years, subjects, types)
    generate_landings(c_name)


if __name__ == "__main__":
    f_name = 'config.json'
    c_name = ''
    userlist_json = ''
    viewlist_json = 'views.json'
    subjectlist_json = 'subjects.json'
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
        if 'user_list' in cfg:
            userlist_json = cfg['user_list']
        if 'view_list' in cfg:
            viewlist_json = 'views.json'
        if 'subject_list' in cfg:
            subjectview_json = 'subjects.json'
    generate_doctree(viewlist_json)
    generate_metadata_structure(
        c_name, 
        userlist_json, 
        viewlist_json,
        subjectlist_json)
    print('OK')
