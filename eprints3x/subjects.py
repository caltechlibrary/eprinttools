
import json

subject_names = { }

def load_subjects(f_name):
    global subject_names
    subject_names = {}
    with open(f_name) as f:
        src = f.read()
        subject_names = json.loads(src)


def normalize_subject(val):
    global subject_names
    if val in subject_names:
        return subject_names[val]
    return ''
    
