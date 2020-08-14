
#
# This is a hand translation of the subject listing in 
# /coda/eprints-3.1/archives/caltechoh/cfg/subjects. It also reflects
# Old examples at http://web.archive.org/web/20190302055622/http://oralhistories.library.caltech.edu/view/subjects
#
subject_names = {
    # top level
    "subjects": "Sets", 
    # second level
    "name": "All Records", 
    "sub": "Subjects", 
    # third level, subjects
    "adm": "Administration", 
    "sem": "Alumni Seminar Day",
    "ast": "Astronomy",
    "bio": "Biology",
    "caltech_womens_club": "Caltech Women's Club",
    "chem": "Chemistry",
    "eng": "Engineering",
    "geo": "Geology",
    "hum": "Humanities",
    "jpl": "Jet Propulsion Laboratory",
    "Keck-Observator": "Keck",
    "ligo": "LIGO",
    "math": "Mathematics",
    "phy": "Physics",
    "policy": "Policy Documents",
    "soc": "Social Sciences",
}

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
    
