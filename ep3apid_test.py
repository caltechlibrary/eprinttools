#!/usr/bin/env python3
'''
ep3api_test.py tests the ep3api.py python library used in acacia object.
'''
import sys
from cltests import TestSet, T, IsSuccessful
from ep3apid import Ep3API, User

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
    for i, name in enumerate(usernames):
        m, err = api.user(name)
        if err:
            t.Expected(True, False, f'Did not expect error ({i}, {name}) {err}')
        else:
            u = User(m)
            print(f'DEBUG u: {u.to_string()}')
            t.Expected(True, u.uname != '', f'Expected u.uname to be populated ({i}, {name})')
            t.Expected(True, u.userid != 0, f'Expected non-zero u.userid ({i}, {name})')
            t.Expected(True, u.name != None, f'Expected u.name to be populated ({i}, {name})')
            t.Expected(True, u.display_name != '', f'Expected u.display_name to be populated ({i}, {name})')
            #t.Expected(True, u.role != '', f'Expected u.role, got empty string ({i}, {name})')
            #required = []
            #t.Expected(True, u.has_role(required), f'Expected u.has_role({required}) to return True')
    return t.Results()

def test_caltechauthors():
    t = T()
    old_repo_id = api.repo_id
    repo_id = 'caltechauthors'
    print(f'Switching repositories from {old_repo_id} to {repo_id}')
    t.Expected(True, api.use(repo_id), f'Expected API repo change to {repo_id}')
    usernames, err = api.usernames()
    for i, name in enumerate(usernames):
        try:
            m, err = api.user(name)
        except Exception as e:
            print(f'{i}, "{name}", {err}, expection:', e)
        u = User(m)
        required = ['admin', 'editor', 'minuser', 'user' ]
        if u.role:
            t.Expected(True, u.role != '', f'Expected u.role, got empty string ({i}, {name}, {repo_id})')
            t.Expected(True, u.has_role(required), f'Expected u.has_role({required}) to return True, u.role "{u.role}" ({i}, {name}, {repo_id})')
        if err:
            t.Expected(False, True, f'Expected no error in api.user({name}), {err} ({i}, {name}, {repo_id}')
    t.Expected(True, api.use(old_repo_id), f'Expected API repo changed back to {old_repo_id}')
    return t.Results()

def test_demo_setup():
    old_repo_id = api.repo_id
    repo_id = 'caltechauthors'
    t = T()
    t.Expected(True, api.use(repo_id), f'Expected to switch to {repo_id}')
    usernames, err = api.usernames()
    t.Expected(None, err, f"Did not expect an error for usernames, {err}")
    if isinstance(usernames, list):
        t.Expected(True, len(usernames) > 0, 'Expected non-zero usernames')
    userids, err = api.lookup_userid('rsdoiel')
    t.Expected(True, isinstance(userids, list), 'Expected a user id list')
    if isinstance(userids, list):
        t.Expected(userids[0], 5487, f'Expected user id 5487, got {userids[0]}')
    dois, err = api.doi()
    t.Expected(None, err, f"Did not expect an error for doi, {err}")
    if isinstance(dois, list):
        #for doi in dois:
        #    print(f'DEBUG doi {doi}')
        t.Expected(True, len(dois) > 0, 'Expected non-zero dois')
        ids, err = api.doi(dois[0])
        t.Expected(None, err, "Did not expect an error for ids")
        t.Expected(True, len(ids) == 1, f'Expected a single id for DOI, {len(ids)}')
        eprint, err = api.eprint(ids[0])
        t.Expected(True, not err, "Did not expect an eprint record, {err}")
        t.Expected(True, isinstance(eprint, dict), "Expected a dict for eprint")
    t.Expected(True, api.use(old_repo_id), f'Expected to switch back to {old_repo_id}')
    return t.Results()


if __name__ == '__main__':
    ts = TestSet("ep3api_test")
    ts.add(test_eprint_lists)
    ts.add(test_user)
    repositories, err = api.repositories()
    repo_id = 'caltechauthors'
    if repo_id in repositories:
        ts.add(test_demo_setup)
    else:
        print(f'Skipping test_demo_setup')
    if repo_id in repositories:
        ts.add(test_caltechauthors)
    else:
        print(f'Skipping test_{repo_id}')
    IsSuccessful(ts.run())




