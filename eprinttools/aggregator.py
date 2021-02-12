import os
import sys
import json

from datetime import date, timedelta

from .normalize import slugify, get_date_year, get_eprint_id, get_object_type, has_creator_ids, has_editor_ids, has_contributor_ids, make_label, get_sort_name, get_sort_year, get_sort_subject, get_sort_publication, get_sort_collection, get_sort_event, get_lastmod_date, get_sort_lastmod, get_sort_issn, get_sort_corp_creator, get_sort_place_of_pub


class Aggregator:
    """This class models the various Eprint aggregations used across Caltech Library repositories"""
    def __init__(self, c_name, objs):
        self.c_name = c_name
        self.objs = objs

    def aggregate_creator(self):
        # now build our people list and create a people, eprint_id, title list
        people = {}
        for obj in self.objs:
            if has_creator_ids(obj):
                # For each author add a reference to object
                for creator in obj['creators']:
                    creator_id = ''
                    if 'id' in creator:
                        creator_id = creator['id']
                    if 'creator_id' in creator:
                        creator_id = creator['creator_id']
                    creator_name = creator['display_name']
                    if creator_id != '':
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
        return people_list

    def aggregate_editor(self):
        # now build our people list based on editors.items
        people = {}
        for obj in self.objs:
            if has_editor_ids(obj):
                # For each author add a reference to object
                for editor in obj['editors']:
                    editor_id = ''
                    if 'id' in editor:
                        editor_id = editor['id']
                    if 'editor_id' in editor:
                        editor_id = editor['editor_id']
                    editor_name = editor['display_name']
                    if editor_id != '':
                        if not editor_id in people:
                            people[editor_id] = { 
                                'key': editor_id,
                                'label': editor_name,
                                'count' : 0,
                                'people_id': editor_id,
                                'sort_name': editor_name,
                                'objects' : []
                            }
                        people[editor_id]['count'] += 1
                        people[editor_id]['objects'].append(obj)
        # Now that we have a people list we need to sort it by name
        people_list = []
        for key in people:
            people_list.append(people[key])
        people_list.sort(key = get_sort_name)
        return people_list

    def aggregate_contributor(self):
        # now build our people list based on contributors.items
        people = {}
        for obj in self.objs:
            if has_contributor_ids(obj):
                # For each author add a reference to object
                for contributor in obj['contributors']:
                    contributor_id = ''
                    if 'id' in contributor:
                        contributor_id = contributor['id']
                    if 'contributor_id' in contributor:
                        contributor_id = contributor['contributor_id']
                    contributor_name = contributor['display_name']
                    if contributor_id != '':
                        if not contributor_id in people:
                            people[contributor_id] = { 
                                'key': contributor_id,
                                'label': contributor_name,
                                'count' : 0,
                                'people_id': contributor_id,
                                'sort_name': contributor_name,
                                'objects' : []
                            }
                        people[contributor_id]['count'] += 1
                        people[contributor_id]['objects'].append(obj)
        # Now that we have a people list we need to sort it by name
        people_list = []
        for key in people:
            people_list.append(people[key])
        people_list.sort(key = get_sort_name)
        return people_list

    def aggregate_by_view_name(self, name, subject_map):
        if name == 'person-az':
            return self.aggregate_person_az()
        elif name == 'person':
            return self.aggregate_person()
        elif name == 'author':
            return self.aggregate_author()
        elif name == 'editor':
            return self.aggregate_editor()
        elif name == 'contributor':
            return self.aggregate_contributor()
        elif name == 'year':
            return self.aggregate_year()
        elif name == 'publication':
            return self.aggregate_publication()
        elif name == 'place_of_pub':
            return self.aggregate_place_of_pub()
        elif name == 'corp_creators':
            return self.aggregate_corp_creators()
        elif name == 'issuing_body':
            return self.aggregate_issuing_body()
        elif name == 'issn':
            return self.aggregate_issn()
        elif name == 'collection':
            return self.aggregate_collection()
        elif name == 'event':
            return self.aggregate_event()
        elif name == 'subjects':
            return self.aggregate_subjects(subject_map)
        elif name == 'ids':
            return self.aggregate_ids()
        elif name == 'types':
            return self.aggregate_types()
        elif name == 'latest':
            return self.aggregate_latest()
        else:
            print(f'WARNING: {name} is unknown aggregation type')
            return None

    def aggregate_person_az(self):
        return self.aggregate_creator()
    
    def aggregate_person(self):
        return self.aggregate_creator()
    
    def aggregate_author(self):
        return self.aggregate_creator()
    
    def aggregate_year(self):
        years = {}
        year = ''
        for obj in self.objs:
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
        return year_list
    
    def aggregate_publication(self):
        publications = {}
        publication = ''
        for obj in self.objs:
            eprint_id = get_eprint_id(obj)
            year = get_date_year(obj)
            if ('publication' in obj):
                publication = obj['publication']
                key = slugify(publication)
                if not publication in publications:
                    publications[publication] = { 
                        'key': key,
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
        return publication_list

    def aggregate_corp_creators(self):
        corp_creators = {}
        for obj in self.objs:
            year = get_date_year(obj)
            if ('corp_creators' in obj) and ('items' in obj['corp_creators']):
                for item in obj['corp_creators']['items']:
                    corp_creator = item['name'].strip()
                    if 'id' in item:
                        key = str(item['id'])
                    else:
                        key = slugify(corp_creator)
                    if not key in corp_creators:
                        corp_creators[key] = { 
                            'key': key,
                            'label': corp_creator,
                            'count': 0,
                            'year': year, 
                            'objects': [] 
                        }
                    corp_creators[key]['count'] += 1
                    corp_creators[key]['objects'].append(obj)
        corp_creator_list = []
        for key in corp_creators:
            corp_creator_list.append(corp_creators[key])
        corp_creator_list.sort(key = get_sort_corp_creator)
        return corp_creator_list

    def aggregate_issuing_body(self):
        return self.aggregate_corp_creators()

    def aggregate_issn(self):
        issns = {}
        issn = ''
        for obj in self.objs:
            eprint_id = get_eprint_id(obj)
            year = get_date_year(obj)
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
        return issn_list

    def aggregate_place_of_pub(self):
        place_of_pubs = {}
        place_of_pub = ''
        for obj in self.objs:
            eprint_id = get_eprint_id(obj)
            year = get_date_year(obj)
            if ('place_of_pub' in obj):
                place_of_pub = obj['place_of_pub'].strip()
                key = slugify(place_of_pub)
                if not place_of_pub in place_of_pubs:
                    place_of_pubs[place_of_pub] = { 
                        'key': key,
                        'label': place_of_pub,
                        'count': 0,
                        'year': year, 
                        'objects': [] 
                    }
                place_of_pubs[place_of_pub]['count'] += 1
                place_of_pubs[place_of_pub]['objects'].append(obj)
        place_of_pub_list = []
        for key in place_of_pubs:
            place_of_pub_list.append(place_of_pubs[key])
        place_of_pub_list.sort(key = get_sort_place_of_pub)
        return place_of_pub_list


    def aggregate_collection(self):
        collections = {}
        collection = ''
        for obj in self.objs:
            eprint_id = get_eprint_id(obj)
            year = get_date_year(obj)
            if ('collection' in obj):
                collection = obj['collection']
                key = slugify(collection)
                if not collection in collections:
                    collections[collection] = { 
                        'key': key,
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
        return collection_list

    def aggregate_event(self):
        events = {}
        event_title = ''
        for obj in self.objs:
            eprint_id = get_eprint_id(obj)
            year = get_date_year(obj)
            event_title = ''
            event_location = ''
            event_dates = ''
            if ('event_title' in obj):
                event_title = obj['event_title']
            if ('event_location' in obj):
                event_location = obj['event_location']
            if ('event_dates' in obj):
                event_dates = obj['event_dates']
            if event_title != '':
                if not event_title in events:
                    key = slugify(event_title)
                    events[event_title] = { 
                        'key': key,
                        'label': event_title,
                        'count': 0,
                        'year': year, 
                        'objects': [] 
                    }
                events[event_title]['count'] += 1
                events[event_title]['objects'].append(obj)
        event_list = []
        for key in events:
            event_list.append(events[key])
        event_list.sort(key = get_sort_event)
        return event_list

    def aggregate_subjects(self, subject_map):
        subjects = {}
        subject = ''
        for obj in self.objs:
            eprint_id = get_eprint_id(obj)
            year = get_date_year(obj)
    
            if ('subjects' in obj):
                for subj in obj['subjects']['items']:
                    subject_name = subject_map.get_subject(subj)
                    if subject_name != None:
                        if not subj in subjects:
                            subjects[subj] = { 
                                'key': subj,
                                'label': subject_name,
                                'count': 0,
                                'subject_id': subj, 
                                'subject_name': subject_name,
                                'objects': [] 
                            }
                        subjects[subj]['count'] += 1
                        subjects[subj]['objects'].append(obj)
        subject_list= []
        for key in subjects:
            subject_list.append(subjects[key])
        subject_list.sort(key = get_sort_subject)
        return subject_list

    def aggregate_ids(self):
        ids = {}
        for obj in self.objs:
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
        return ids_list
    
    def aggregate_types(self):
        types = {}
        for obj in self.objs:
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
        return type_list

    def aggregate_latest(self):
        latest = {}
        today = date.today()
        td = timedelta(days = -7)
        seven_days_ago = (today - td).isoformat()
        for obj in self.objs:
            lastmod = get_lastmod_date(obj)
            if (lastmod != '') and (lastmod >= seven_days_ago):
                key = get_sort_lastmod(obj)
                year = get_date_year(obj)
                if not key in latest:
                    lastest[lastmod] = {
                        'key': key,
                        'label': lastmod,
                        'year': year,
                        'count': 0,
                        'objects': []
                    }
                latest[lastmod]['count'] += 1
                latest[lastmod]['objects'].append(obj)
        latest_list = []
        for key in latest:
            latest_list.append(latest[key])
        latest_list.sort(key = lambda x: x['key'], reverse = True)
        return latest_list
