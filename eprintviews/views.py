import os
import sys
import json

class Views:
    """Model a collection of views to replicate from an EPrints repository"""
    def __init__(self):
        self.views = {}

    def load_views(self, f_name):
        with open(f_name) as f:
            src = f.read()
            self.views = json.loads(src)

    def has_view(self, key):
        if key in self.views:
            return True
        return False

    def get_view(self, key):
        if key in self.views:
            return self.views[key]
        return None

    def get_keys(self):
        keys = []
        for key in self.views:
            keys.append(key)
        return keys
