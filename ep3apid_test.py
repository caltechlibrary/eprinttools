#!/usr/bin/env python3
'''
ep3api_test.py tests the ep3api.py python library used in acacia object.
'''
from cltests import TestSet, T, IsSuccessful
from ep3apid import Ep3API

api = Ep3API('http://localhost:8484', 'lemurprints')

def test_eprint_lists():
    t = T()
    keys, err = api.keys()
    t.Expected(None, err, f"Did not expect an error for keys, {err}")
    t.Expected(True, keys != None, "Should have a list of keys")
    return t.Results()

def test_user():
    t = T()
    usernames, err = api.usernames()
    t.Expected(None, err, f"Did not expect an error for usernames, {err}")
    t.Expected(True, usernames != None, "Should have some usernames")
    t.Expected(True, isinstance(usernames, list), "Expected a list of usernames")
    return t.Results()


def test_demo_setup():
    demo_api = Ep3API('http://localhost:8484', 'caltechauthors')
    t = T()
    usernames, err = demo_api.usernames()
    t.Expected(None, err, f"Did not expect an error for usernames, {err}")
    if isinstance(usernames, list):
        t.Expected(True, len(usernames) > 0, 'Expected non-zero usernames')
    userids, err = demo_api.lookup_userid('rsdoiel')
    t.Expected(True, isinstance(userids, list), 'Expected a user id list')
    if isinstance(userids, list):
        t.Expected(userids[0], 5487, f'Expected user id 5487, got {userids[0]}')
    dois, err = demo_api.doi()
    t.Expected(None, err, f"Did not expect an error for doi, {err}")
    if isinstance(dois, list):
        #for doi in dois:
        #    print(f'DEBUG doi {doi}')
        t.Expected(True, len(dois) > 0, 'Expected non-zero dois')
        ids, err = demo_api.doi(dois[0])
        t.Expected(None, err, "Did not expect an error for ids")
        t.Expected(True, len(ids) == 1, f'Expected a single id for DOI, {len(ids)}')
        eprint, err = demo_api.eprint(ids[0])
        t.Expected(None, err, "Did not expect an eprint record")
        t.Expected(True, isinstance(eprint, dict), "Expected a dict for eprint")
    return t.Results()


if __name__ == '__main__':
    ts = TestSet("ep3api_test")
    ts.add(test_eprint_lists)
    ts.add(test_user)
    ts.add(test_demo_setup)
    IsSuccessful(ts.run())




