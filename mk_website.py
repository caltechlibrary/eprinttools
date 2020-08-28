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

from eprintviews import Views

#
# mkpage wraps the mkpage command from mkpage using the
# pandoc setup.
#
def mkpage(o_file, template = '', data = []):
    cmd = ['mkpage', '-t', 'markdown_strict', '-o', o_file]
    for item in data:
        cmd.append(item)
    if template != '':
        cmd.append(template)
    with Popen(cmd, stdout = PIPE, stderr = PIPE) as proc:
        err = proc.stderr.read().strip().decode('utf-8')
        if err != '':
            print(f"\n{' '.join(cmd)}\nerror: {err}")
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
    # Note we're just grabbing the ids from ids_list, we'll index.json
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

# Load Objects a "*.json" file from the htdocs tree for processing.
def load_objects(json_name):
    f_name = os.path.join('htdocs', json_name)
    with open(f_name) as f:
        src = f.read()
        return json.loads(src)
    return None

def load_keys(key_name):
    f_name = os.path.join('htdocs', key_name)
    with open(f_name) as f:
        src = f.read()
        return src.strip().split('\n')
    return None

def build_top_level_pages(htdocs, site_title, organization, site_welcome):
    print(f'build top level pages')
    css_folder = os.path.join('static', 'css')
    assets_folder = os.path.join('static', 'assets')
    js_folder = os.path.join('static', 'js')
    favicon = os.path.join('static', 'favicon.ico')
    copy_assets(htdocs, css_folder, 'static/')
    copy_assets(htdocs, assets_folder, 'static/')
    copy_assets(htdocs, js_folder, 'static/')
    copy_assets(htdocs, favicon, 'static/')

    # Render Homepage (index.html)
    page_title = f'{site_title}'
    title = site_welcome
    md_filename = os.path.join('static', 'index.md')
    html_filename = 'index.html'
    template_name = 'index-html.tmpl'
    assemble(html_filename, template_name, [
        kv('organization', 'text', organization),
        kv('site_title', 'text', organization),
        kv('page_title', 'text', page_title),
        kv('title', 'text', title),
        kv('content', '', md_filename)
    ])

    # Render about page (information.html)
    page_title = f'{site_title}: About'
    title = 'About this Repository'
    md_filename = os.path.join('static', 'information.md')
    html_filename = 'information.html'
    template_name = 'index-html.tmpl'
    assemble(html_filename, template_name, [
        kv('organization', 'text', organization),
        kv('site_title', 'text', organization),
        kv('page_title', 'text', page_title),
        kv('title', 'text', title),
        kv('content', '', md_filename)
    ])

    # Render browseviews page
    page_title = f'{site_title}: Browse'
    title = 'Browse this Repository'
    md_filename = os.path.join('static', 'browseviews.md')
    html_filename = 'browseviews.html'
    template_name = 'index-html.tmpl'
    assemble(html_filename, template_name, [
        kv('organization', 'text', organization),
        kv('site_title', 'text', organization),
        kv('page_title', 'text', page_title),
        kv('title', 'text', title),
        kv('content', '', md_filename)
    ])

    # Render policies page
    page_title = f'{site_title}: Policies'
    title = 'Repository Policies'
    md_filename = os.path.join('static', 'policies.md')
    if os.path.exists(md_filename):
        html_filename = 'policies.html'
        template_name = 'index-html.tmpl'
        assemble(html_filename, template_name, [
            kv('organization', 'text', organization),
            kv('site_title', 'text', organization),
            kv('page_title', 'text', page_title),
            kv('title', 'text', title),
            kv('content', '', md_filename)
        ])

    # Render contact page
    page_title = f'{site_title}: Contact'
    title = 'Contact Informaiton'
    md_filename = os.path.join('static', 'contact.md')
    if os.path.exists(md_filename):
        html_filename = 'contact.html'
        template_name = 'index-html.tmpl'
        assemble(html_filename, template_name, [
            kv('organization', 'text', organization),
            kv('site_title', 'text', organization),
            kv('page_title', 'text', page_title),
            kv('title', 'text', title),
            kv('content', '', md_filename)
        ])


def make_view(view, label, site_title, organization):
    # Render view, e.g. /view/ids/
    page_title = f'{site_title}: Browse by {label}'
    title = f'Browse by {label}'
    list_data = f'view/{view}/{view}_list.json'
    html_filename = f'view/{view}/index.html'
    template_name = 'listing-html.tmpl'
    content = 'Please select a value to browse from the list below.'
    objs = load_objects(list_data)
    tot = len(objs)
    if tot == 0:
        content = 'Nothing available.'
    assemble(html_filename, template_name, [
        kv('organization', 'text', organization),
        kv('site_title', 'text', site_title),
        kv('page_title', 'text', page_title),
        kv('title', 'text', title),
        kv('content', 'text', content),
        kv('listing', '', os.path.join('htdocs', list_data))
    ])
    bar = progressbar.ProgressBar(
          max_value = tot,
            widgets = [
                progressbar.Percentage(), ' ',
                progressbar.Counter(), f'/{tot} ',
                progressbar.AdaptiveETA(),
                f' rendering views/{view}/*'
            ], redirect_stdout=False)
    for i, obj in enumerate(objs):
        key = obj['key']
        count = obj['count']
        object_list = os.path.join('htdocs', 'view', view, f'{key}.json')
        html_filename = f'view/{view}/{key}.html'
        template_name = 'object-list-html.tmpl'
        assemble(html_filename, template_name, [
            kv('organization', 'text', organization),
            kv('site_title', 'text', site_title),
            kv('key', 'text', key),
            kv('page_title', 'text', page_title),
            kv('title', 'text', title),
            kv('count', 'text', count),
            kv('object_list', '', object_list)
        ])
        bar.update(i)
    bar.finish()


def build_view_pages(views, htdocs, site_title, organization):
    print(f'building view pages')
    # Render view/
    page_title = f'{site_title}: Browse'
    title = 'Browse Items'
    md_filename = os.path.join('static', 'view_index.md')
    html_filename = 'view/index.html'
    template_name = 'index-html.tmpl'
    assemble(html_filename, template_name, [
        kv('organization', 'text', organization),
        kv('site_title', 'text', site_title),
        kv('page_title', 'text', page_title),
        kv('title', 'text', title),
        kv('content', '', md_filename)
    ])

    keys = views.get_keys()
    for key in keys:
        label = views.get_view(key)
        make_view(key, label, site_title, organization)
    print(f'build view pages completed')


def build_landing_pages(htdocs, site_title, organization):
    print(f'build landing pages')
    objs = load_objects('view/ids/ids_list.json')
    keys = load_keys('index.keys')

    tot = len(keys)
    bar = progressbar.ProgressBar(
          max_value = tot,
            widgets = [
                progressbar.Percentage(), ' ',
                progressbar.Counter(), f'/{tot} ',
                progressbar.AdaptiveETA(),
                f' rendering landing pages'
            ], redirect_stdout=False)
    for i, key in enumerate(keys):
        obj = load_objects(f'{key}/index.json')
        title =  obj['title']
        page_title = f'{title} - {site_title}'
        object_path = os.path.join('htdocs', f'{key}', f'index.json')
        html_filename = f'{key}/index.html'
        template_name = 'landing-page-html.tmpl'
        assemble(html_filename, template_name, [
            kv('organization', 'text', organization),
            kv('site_title', 'text', site_title),
            kv('page_title', 'text', page_title),
            kv('title', 'text', title),
            kv('object', '', object_path)
        ])
        bar.update(i)
    bar.finish()
    print(f'build landing pages completed')


def build_search_page(htdocs, site_title, organization):
    print(f'build search page')
    data = []
    page_title = f'{site_title}: Search'
    title = f'Search'
    nav = os.path.join('static', 'nav.md')
    content = os.path.join('static', 'search_form.md')
    if page_title != '':
        data.append(kv('page_title', 'text', page_title))
    if organization != '':
        data.append(kv('organization', 'text', organization))
    if title != '':
        data.append(kv('title', 'text', title))
    if os.path.exists(nav):
        data.append(kv('nav', '', nav))
    if os.path.exists(content):
        data.append(kv('content', '', content))
    else:
        print(f'''Can't find {content}''')
        sys.exit(1)
    mkpage(os.path.join('htdocs', 'search.html'), os.path.join('templates', 'search-html.tmpl'), data)


def build_website(f_views, htdocs, site_title, site_welcome, organization):
    views = Views()
    views.load_views(f_views)
    build_top_level_pages(htdocs, site_title, organization, site_welcome)
    build_view_pages(views, htdocs, site_title, organization)
    build_landing_pages(htdocs, site_title, organization)
    build_search_page(htdocs, site_title, organization)
    

if __name__ == "__main__":
    f_name = ''
    f_views = ''
    htdocs = 'htdocs'
    organization = 'Example Library EDU'
    site_title = 'EPrints Repository'
    site_welcome = 'EPrints Repository public website'
    args = [] 
    if len(sys.argv) > 1:
        f_name = sys.argv[1]
    if len(sys.argv) > 2:
        args = sys.argv[2:]
    if f_name == '':
        print(f'Missing configuraiton filename.')
        sys.exit(1)
    if not os.path.exists(f_name):
        print(f'Missing {f_name} configuration file.')
        sys.exit(1)
        
    with open(f_name) as f:
        src = f.read()
        cfg = json.loads(src)
        if 'organization' in cfg:
            organization = cfg['organization']
        if 'site_welcome' in cfg:
            site_welcome = cfg['site_welcome']
        if 'site_title' in cfg:
            site_title = cfg['site_title']
        if 'views' in cfg:
            f_views = cfg['views']

    if not os.path.exists(htdocs):
        print(f'''Cannot find the htdocs directory''')
        sys.exit(1)
    if f_views == '':
        print(f'Missing views in {f_name} file.')
        sys.exixt(1)
    if not os.path.exists(f_views):
        print(f'Missing views configuration file.')
        sys.exit(1)
    if len(args) > 0:
        views = Views()
        views.load_views(f_views)
        for arg in args:
            if arg == 'top':
                build_top_level_pages(htdocs, site_title, organization, site_welcome)
                build_search_page(htdocs, site_title, organization)
            if arg == 'view':
                build_view_pages(views, htdocs, site_title, organization)
            if arg == 'landing':
                build_landing_pages(htdocs, site_title, organization)
    else:
        build_website(f_views, htdocs, site_title, organization, site_welcome)
