
import json

from py_dataset import dataset

from .config import Configuration
from .normalize import get_title, get_date_year, get_eprint_id

def make_frame_date_title(cfg):
    c_name = cfg.dataset
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
    dot_paths = [ 
        '.eprint_id', 
        '.title', 
        '.date', 
        '.creators', 
        '.editors', 
        '.contributors', 
        '.corp_creators', 
        '.subjects', 
        '.type', 
        '.official_url', 
        '.collection', 
        '.event_title', 
        '.event_location', 
        '.event_dates', 
        '.publication', 
        '.place_of_pub', 
        '.number', 
        '.volume', 
        '.series' 
        '.pagerange', 
        '.userid' , 
        '.lastmod'
    ]
    labels = []
    for label in dot_paths:
        labels.append(label[1:])
    ok = dataset.frame_create(c_name, frame_name, keys, dot_paths, labels)
    if not ok:
        err = dataset.error_message()
        return f'{frame_name} in {c_name}, not created, {err}'
    return ''

