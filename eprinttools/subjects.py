import os
import sys
import json

class Subjects:
    """Subjects models the Eprint subjects values into a flat list"""

    def __init__ (self):
        self.subjects = {}

    def load_subjects(self, f_name):
        """Load an Eprint subjects file, e.g. /eprint3-1/archives/REPO/cfg/subjects"""
        with open(f_name) as f:
            lines = f.readlines()
            for line in lines:
                # Trim comment
                if "#" in line:
                    parts = line.split(sep = "#", maxsplit = 2)
                    line = parts[0].strip()
                # split on colon
                if ":" in line:
                    parts = line.split(":")
                    if len(parts) > 1:
                        # get key and label
                        key = parts[0].strip()
                        label = parts[1].strip()
                        # add to self.subjects
                        self.subjects[key] = label

        
    def has_subject(self, key):
        if key in self.subjects:
            return True
        return False
        
    def get_subject(self, key):
        if key in self.subjects:
            return self.subjects[key]
        return ''

    def get_keys(self):
        keys = []
        for key in self.subjects:
            keys.append(key)
        return keys
