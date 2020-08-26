#!/usr/bin/env python3

#
# mk_website.py crawls the htdocs tree and turns all the markdown
# files except nav.md into an HTML page and HTML include page.
#
import os
import sys
import json
from urllib.parse import urlparse
import shutil
from subprocess import run, Popen, PIPE

import progressbar

#
# mkpage wraps the mkpage command from mkpage using the
# pandoc setup.
#
def mkpage(o_file, template = '', data = []):
    cmd = ['mkpage', '-o', o_file]
    for item in data:
        cmd.append(item)
    if template != '':
        cmd.append(template)
    with Popen(cmd, stdout = PIPE, stderr = PIPE) as proc:
        err = proc.stderr.read().strip().decode('utf-8')
        if err != '':
            print(f"{' '.join(cmd[0:3])} error: {err}")
        out = proc.stdout.read().strip().decode('utf-8')
        if out != "":
            print(f"{out}");

def mkpage_version_no():
    cmd = ['mkpage', '-version']
    p = Popen(cmd, stdout = PIPE, stderr = PIPE)
    (version, err) = p.communicate()
    if err.decode('utf-8') != '':
        print(f"ERROR: mkpage -version, {err.decode('utf-8')}")
        sys.exit(1)
    return version.decode('utf-8')

mkpage_version = (mkpage_version_no()).strip()

# Make a key/value pair to pass to mkpage in the data array
def kv(key, protocol = '', value = ''):
    if protocol == '':
        return f'{key}={value}'
    else:
        return f'{key}={protocol}:{value}'

# Assemble page
def assemble(html_filename, template_name, data):
    nav = os.path.join('static', 'nav.md')
    if os.path.exists(nav):
        data.append(kv('nav', '', nav))
    announcement = os.path.join('static', 'announcement.md')
    if os.path.exists(announcement):
        data.append(kv('announcement', '', announcement))
    mkpage(os.path.join('htdocs', html_filename), os.path.join('templates', template_name), data)

# Copy file and folder assets
def copy_assets(htdocs, source, trim_prefix):
    rel_path = source[len(trim_prefix):]
    dst = os.path.join(htdocs, rel_path)
    if os.path.exists(source):
        if os.path.isdir(source):
            print(f'    copying folder {source} to {htdocs}')
            if os.path.exists(dst):
                shutil.rmtree(dst)
            shutil.copytree(source, dst)
        elif os.path.isfile(source):
            print(f'    copying file {source} to {htdocs}')
            if os.path.exists(dst):
                os.remove(dst)
            shutil.copy2(source, dst, follow_symlinks = False)
    else:
        print(f'WARNING: {source} does not exist')

# Load Object a "*.json" file from the htdocs tree for processing.
def load_object(json_name):
    f_name = os.path.join('htdocs', json_name)
    with open(f_name) as f:
        src = f.read()
        return json.loads(src)
    return None

def build_top_level_pages(htdocs, site_title, site_welcome):
    print(f'build top level pages')
    css_folder = os.path.join('static', 'css')
    assets_folder = os.path.join('static', 'assets')
    favicon = os.path.join('static', 'favicon.ico')
    copy_assets(htdocs, css_folder, 'static/')
    copy_assets(htdocs, assets_folder, 'static/')
    copy_assets(htdocs, favicon, 'static/')

    # Render Homepage (index.html)
    page_title = f'{site_title}'
    title = site_welcome
    md_filename = os.path.join('static', 'index.md')
    html_filename = 'index.html'
    template_name = 'index-html.tmpl'
    assemble(html_filename, template_name, [
        kv('page_title', 'text', page_title),
        kv('title', 'text', title),
        kv('content', '', md_filename)
    ])

    # Render about page (information.html)
    page_title = f'{site_title}: About'
    title = 'About the Repository'
    md_filename = os.path.join('static', 'information.md')
    html_filename = 'information.html'
    template_name = 'index-html.tmpl'
    assemble(html_filename, template_name, [
        kv('page_title', 'text', page_title),
        kv('title', 'text', title),
        kv('content', '', md_filename)
    ])

    # Render browseviews page
    page_title = f'{site_title}: Browse'
    title = 'Browse the Repository'
    md_filename = os.path.join('static', 'browseviews.md')
    html_filename = 'browseviews.html'
    template_name = 'index-html.tmpl'
    assemble(html_filename, template_name, [
        kv('page_title', 'text', page_title),
        kv('title', 'text', title),
        kv('content', '', md_filename)
    ])



def build_view_pages(htdocs, site_title):
    print(f'building view pages')
    # Render view/
    page_title = f'{site_title}: Browse'
    title = 'Browse Items'
    md_filename = os.path.join('static', 'view_index.md')
    html_filename = 'view/index.html'
    template_name = 'index-html.tmpl'
    assemble(html_filename, template_name, [
        kv('page_title', 'text', page_title),
        kv('title', 'text', title),
        kv('content', '', md_filename)
    ])

    # Render view/ids/
    page_title = f'{site_title}: Browse by Eprint ID'
    title = 'Browse by Eprint ID'
    list_data = 'view/ids/ids_list.json'
    html_filename = 'view/ids/index.html'
    template_name = 'listing-html.tmpl'
    assemble(html_filename, template_name, [
        kv('page_title', 'text', page_title),
        kv('title', 'text', title),
        kv('content', 'text', 
            'Please select a value to browse from the list below.'),
        kv('listing', '', os.path.join('htdocs', list_data))
    ])
    # Render view/ids/ sub listings
    objs = load_object('view/ids/ids_list.json')
    tot = len(objs)
    bar = progressbar.ProgressBar(
          max_value = tot,
            widgets = [
                progressbar.Percentage(), ' ',
                progressbar.Counter(), f'/{tot} ',
                progressbar.AdaptiveETA(),
                f' rendering views/ids/*'
            ], redirect_stdout=False)
    for i, obj in enumerate(objs):
        key = obj['key']
        count = obj['count']
        object_list = os.path.join('htdocs', 'view', 'ids', f'{key}.json')
        html_filename = f'view/ids/{key}.html'
        template_name = 'object-list-html.tmpl'
        assemble(html_filename, template_name, [
            kv('page_title', 'text', page_title),
            kv('title', 'text', title),
            kv('count', 'text', count),
            kv('object_list', '', object_list)
        ])
        bar.update(i)
    bar.finish()

    # Render view/year/
    page_title = f'{site_title}: Browse by Year'
    title = 'Browse by Year'
    list_data = 'view/year/year_list.json'
    html_filename = 'view/year/index.html'
    template_name = 'listing-html.tmpl'
    assemble(html_filename, template_name, [
        kv('page_title', 'text', page_title),
        kv('title', 'text', title),
        kv('content', 'text', 
            'Please select a value to browse from the list below.'),
        kv('listing', '', os.path.join('htdocs', list_data))
    ])
    # Render view/year/ sub listings
    objs = load_object('view/year/year_list.json')
    tot = len(objs)
    bar = progressbar.ProgressBar(
          max_value = tot,
            widgets = [
                progressbar.Percentage(), ' ',
                progressbar.Counter(), f'/{tot} ',
                progressbar.AdaptiveETA(),
                f' rendering views/year/*'
            ], redirect_stdout=False)
    for i, obj in enumerate(objs):
        key = obj['key']
        count = obj['count']
        object_list = os.path.join('htdocs', 'view', 'year', f'{key}.json')
        html_filename = f'view/year/{key}.html'
        template_name = 'object-list-html.tmpl'
        assemble(html_filename, template_name, [
            kv('page_title', 'text', page_title),
            kv('title', 'text', title),
            kv('count', 'text', count),
            kv('object_list', '', object_list)
        ])
        bar.update(i)
    bar.finish()
        

    # Render view/subjects/
    page_title = f'{site_title}: Browse by Subject'
    title = 'Browse by Subject'
    list_data = 'view/subjects/subject_list.json'
    html_filename = 'view/subjects/index.html'
    template_name = 'listing-html.tmpl'
    assemble(html_filename, template_name, [
        kv('page_title', 'text', page_title),
        kv('title', 'text', title),
        kv('content', 'text', 
            'Please select a value to browse from the list below.'),
        kv('listing', '', os.path.join('htdocs', list_data))
    ])
    # Render view/subjects/* sub lists
    objs = load_object('view/subjects/subject_list.json')
    tot = len(objs)
    bar = progressbar.ProgressBar(
          max_value = tot,
            widgets = [
                progressbar.Percentage(), ' ',
                progressbar.Counter(), f'/{tot} ',
                progressbar.AdaptiveETA(),
                f' rendering views/subjects/*'
            ], redirect_stdout=False)
    for i, obj in enumerate(objs):
        key = obj['key']
        count = obj['count']
        object_list = os.path.join('htdocs', 'view', 'subjects', f'{key}.json')
        html_filename = f'view/subjects/{key}.html'
        template_name = 'object-list-html.tmpl'
        assemble(html_filename, template_name, [
            kv('page_title', 'text', page_title),
            kv('title', 'text', title),
            kv('count', 'text', count),
            kv('object_list', '', object_list)
        ])
        bar.update(i)
    bar.finish()

    # Render view/types/
    page_title = f'{site_title}: Browse by Type'
    title = 'Browse by Type'
    list_data = 'view/types/type_list.json'
    html_filename = 'view/types/index.html'
    template_name = 'listing-html.tmpl'
    assemble(html_filename, template_name, [
        kv('page_title', 'text', page_title),
        kv('title', 'text', title),
        kv('content', 'text', 
            'Please select a value to browse from the list below.'),
        kv('listing', '', os.path.join('htdocs', list_data))
    ])
    # Render view/types/* sub lists
    objs = load_object('view/types/type_list.json')
    tot = len(objs)
    bar = progressbar.ProgressBar(
          max_value = tot,
            widgets = [
                progressbar.Percentage(), ' ',
                progressbar.Counter(), f'/{tot} ',
                progressbar.AdaptiveETA(),
                f' rendering views/types/*'
            ], redirect_stdout=False)
    for i, obj in enumerate(objs):
        key = obj['key']
        count = obj['count']
        object_list = os.path.join('htdocs', 'view', 'types', f'{key}.json')
        html_filename = f'view/types/{key}.html'
        template_name = 'object-list-html.tmpl'
        assemble(html_filename, template_name, [
            kv('page_title', 'text', page_title),
            kv('title', 'text', title),
            kv('count', 'text', count),
            kv('object_list', '', object_list)
        ])
        bar.update(i)
    bar.finish()

    # Render view/person-az/
    page_title = f'{site_title}: Browse by Person'
    title = 'Browse by Person'
    list_data = 'view/person-az/people_list.json'
    html_filename = 'view/person-az/index.html'
    template_name = 'listing-html.tmpl'
    assemble(html_filename, template_name, [
        kv('page_title', 'text', page_title),
        kv('title', 'text', title),
        kv('content', 'text', 
            'Please select a value to browse from the list below.'),
        kv('listing', '', os.path.join('htdocs', list_data))
    ])
    # Render view/person-az/* sub lists
    objs = load_object('view/person-az/people_list.json')
    tot = len(objs)
    bar = progressbar.ProgressBar(
          max_value = tot,
            widgets = [
                progressbar.Percentage(), ' ',
                progressbar.Counter(), f'/{tot} ',
                progressbar.AdaptiveETA(),
                f' rendering views/person-az/*'
            ], redirect_stdout=False)
    for i, obj in enumerate(objs):
        key = obj['key']
        count = obj['count']
        object_list = os.path.join('htdocs', 'view', 'person-az', f'{key}.json')
        html_filename = f'view/person-az/{key}.html'
        template_name = 'object-list-html.tmpl'
        assemble(html_filename, template_name, [
            kv('page_title', 'text', page_title),
            kv('title', 'text', title),
            kv('count', 'text', count),
            kv('object_list', '', object_list)
        ])
        bar.update(i)
    bar.finish()
    print(f'build view pages completed')


def build_landing_pages(htdocs, site_title):
    print(f'build landing pages')
    objs = load_object('view/ids/ids_list.json')
    tot = len(objs)
    bar = progressbar.ProgressBar(
          max_value = tot,
            widgets = [
                progressbar.Percentage(), ' ',
                progressbar.Counter(), f'/{tot} ',
                progressbar.AdaptiveETA(),
                f' rendering landing pages'
            ], redirect_stdout=False)
    for i, obj in enumerate(objs):
        key = obj['key']
        title =  obj['label']
        page_title = f'{title} - {site_title}'
        object_path = os.path.join('htdocs', f'{key}', f'index.json')
        html_filename = f'{key}/index.html'
        template_name = 'landing-page-html.tmpl'
        assemble(html_filename, template_name, [
            kv('page_title', 'text', page_title),
            kv('title', 'text', title),
            kv('object', '', object_path)
        ])
        bar.update(i)
    bar.finish()
    print(f'build landing pages completed')


def build_search_page(htdocs, site_title):
    print(f'build search page')
    data = []
    page_title = f'{site_title}: Search'
    title = f'Search'
    nav = os.path.join('static', 'nav.md')
    content = os.path.join('static', 'search_form.md')
    if page_title != '':
        data.append(kv('page_title', 'text', page_title))
    if title != '':
        data.append(kv('title', 'text', title))
    if os.path.exists(nav):
        data.append(kv('nav', '', nav))
    if os.path.exists(content):
        data.append(kv('content', '', content))
    else:
        print(f'''Can't find {content}''')
        sys.exit(1)
    mkpage(os.path.join('htdocs', 'search.html'), os.path.join('templates', 'index-html.tmpl'), data)


def build_website(htdocs, site_title, site_welcome):
    build_top_level_pages(htdocs, site_title, site_welcome)
    build_view_pages(htdocs, site_title)
    build_landing_pages(htdocs, site_title)
    build_search_page(htdocs, site_title)
    

if __name__ == "__main__":
    f_name = ''
    htdocs = 'htdocs'
    site_title = 'EPrints Repository'
    site_welcome = 'EPrints Repository public website'
    if len(sys.argv) > 1:
        f_name = sys.argv[1]
    if f_name == '':
        print(f'Missing configuration filename.')
        sys.exit(1)
    if not os.path.exists(f_name):
        print(f'Missing {f_name} configuration file.')
        sys.exit(1)
        
    with open(f_name) as f:
        src = f.read()
        cfg = json.loads(src)
        if 'site_welcome' in cfg:
            site_welcome = cfg['site_welcome']
        if 'site_title' in cfg:
            site_title = cfg['site_title']
                

    if not os.path.exists(htdocs):
        print(f'''Cannot find the htdocs directory''')
        sys.exit(1)
    build_website(htdocs, site_title, site_welcome)
