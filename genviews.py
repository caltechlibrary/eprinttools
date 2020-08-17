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

from eprints3x import load_subjects, normalize_subject, load_users, get_user, has_user, load_views

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

def get_sort_publication(o):
    if ('publication' in o) and ('item' in publication['publication']):
        return o['publication']['item']
    return ''

def get_sort_collection(o):
    if ('collection' in o):
        return o['collection']
    return ''

def get_sort_event(o):
    if ('event_title' in o):
        return o['event_title']
    return ''

def get_sort_issn(o):
    if ('issn' in o):
        return o['issn']
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

def aggregate_publication(c_name, objs):
    publications = {}
    publication = ''
    for obj in objs:
        eprint_id = get_eprint_id(obj)
        year = get_date(obj)
        if ('publication' in obj):
            publication = obj['publication']
            if not publication in publications:
                publications[publication] = { 
                    'key': str(publication),
                    'label': str(publication),
                    'count': 0,
                    'year': year, 
                    'objects': [] 
                }
            publications[publication]['count'] += 1
            publications[publication]['objects'].append(obj)
    publication_list = []
    for key in publications:
        publication_list.append(publications[key])
    publication_list.sort(key = get_sort_publication)
    if len(publication_list) == 0:
        print(f'Failed to aggregate any publications from {c_name}')
        sys.exit(1)
    return publication_list

def aggregate_issn(c_name, objs):
    issns = {}
    issn = ''
    for obj in objs:
        eprint_id = get_eprint_id(obj)
        year = get_date(obj)
        if ('issn' in obj):
            issn = obj['issn']
            if not issn in issns:
                issns[issn] = { 
                    'key': str(issn),
                    'label': str(issn),
                    'count': 0,
                    'year': year, 
                    'objects': [] 
                }
            issns[issn]['count'] += 1
            issns[issn]['objects'].append(obj)
    issn_list = []
    for key in issns:
        issn_list.append(issns[key])
    issn_list.sort(key = get_sort_issn)
    if len(issn_list) == 0:
        print(f'Failed to aggregate any ISSNs from {c_name}')
        sys.exit(1)
    return issn_list

def aggregate_collection(c_name, objs):
    collections = {}
    collection = ''
    for obj in objs:
        eprint_id = get_eprint_id(obj)
        year = get_date(obj)
        if ('collection' in obj):
            collection = obj['collection']
            if not collection in collections:
                collections[collection] = { 
                    'key': collection,
                    'label': collection,
                    'count': 0,
                    'year': year, 
                    'objects': [] 
                }
            collections[collection]['count'] += 1
            collections[collection]['objects'].append(obj)
    collection_list = []
    for key in collections:
        collection_list.append(collections[key])
    collection_list.sort(key = get_sort_collection)
    if len(collection_list) == 0:
        print(f'Failed to aggregate any collections from {c_name}')
        sys.exit(1)
    return collection_list

def aggregate_event(c_name, objs):
    events = {}
    event_title = ''
    for obj in objs:
        eprint_id = get_eprint_id(obj)
        year = get_date(obj)
        event_title = ''
        event_location = ''
        event_dates = ''
        if ('event_title' in obj):
            event_title = obj['event_title']
        if ('event_location' in obj):
            event_location = obj['event_location']
        if ('event_dates' in obj):
            event_dates = obj['event_dates']
        if not event_title in events:
            events[event_title] = { 
                'key': event_title,
                'label': event_title,
                'count': 0,
                'year': year, 
                'objects': [] 
            }
        events[event_title]['count'] += 1
        events[event_title]['objects'].append(obj)
    event_list = []
    for key in events:
        event_list.append(event_list[key])
    event_list.sort(key = get_sort_event)
    if len(event_list) == 0:
        print(f'Failed to aggregate any events from {c_name}')
        sys.exit(1)
    return event_list

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
def aggregate(c_name, views):
    err = make_frame_date_title(c_name)
    if err != '':
        print(f'{err}')
    frame_name = 'date-title'
    aggregation = {}
    objs = dataset.frame_objects(c_name, frame_name)
    objs = normalize_objects(objs)
    if 'ids' in views:
        aggregation['ids'] = aggregate_ids(c_name, objs)
    if ('person' in views) or ('person-az' in views) or ('author' in views):
        aggregation['person'] = aggregate_people(c_name, objs)
    if 'year' in views:
        aggregation['year'] = aggregate_year(c_name, objs)
    if 'subjects' in views:
        aggregation['subjects'] = aggregate_subjects(c_name, objs)
    if 'types' in views:
        aggregation['types'] = aggregate_types(c_name, objs)
    if 'group' in views:
        aggregation['group'] = aggregate_group(c_name, objs)
    if 'publication' in views:
        aggregation['publication'] = aggregate_publication(c_name, objs)
    if 'issn' in views:
        aggregation['issn'] = aggregate_issn(c_name, objs)
    if 'collection' in views:
        aggregation['collection'] = aggregate_collection(c_name, objs)
    return aggregation


#
# Build the directory tree
#
def generate_directories(views):
    doc_tree = {}
    for key in views:
        doc_tree[key] = os.path.join('htdocs', 'view', key)
    for key in doc_tree:
        d_name = doc_tree[key]
        if not os.path.exists(d_name):
            os.makedirs(d_name, mode = 0o777, exist_ok = True)
    

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
    for v_name in views:
        if (v_name in [ 'people', 'person', 'person-az' ]) and ('people' in aggregations):
            object_list = aggregations['people']['objects']
            label = views[v_name]
            field_id = 'people_id'
            p_name = os.path.join('htdocs', 'view', v_name)
            make_view(label, p_name, field_id, object_list)
        elif (v_name in [ 'year' ]) and ('year' in aggregations):
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
        elif v_name in aggregations:
            object_list = aggregations[v_name]['objects']
            label = views[v_name]
            field_id = field_ids[v_name]
            p_name = os.path.join('htdocs', 'view', v_name)
            make_view(label, p_name, field_id, object_list)
        else:
            print(f'''WARNING: don't know how to create view {v_name}''')


def generate_metadata_structure(c_name, viewlist_json):
    views = load_views(viewlist_json)
    generate_directories(views)
    # NOTE: To keep the code simple we're only supporting
    # A default list of views found in common with Caltech Library's
    # repositories. If this gets used by other libraries we can
    # look at generalizing this approach.
    aggregation = aggregate(c_name, views)
    print(f'Found {len(aggregation)} aggregations:')
    for key in aggregation:
        print(f'  {len(aggregation[key])} {key}')
    print('')
    generate_views(views, aggregation)
    generate_landings(c_name)


if __name__ == "__main__":
    f_name = 'config.json'
    c_name = ''
    viewlist_json = ''
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
        if 'view_list' in cfg:
            viewlist_json = cfg['view_list']
    generate_metadata_structure(c_name, viewlist_json)
    print('OK')
