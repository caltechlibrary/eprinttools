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

from eprinttools import Configuration, Views

#
# mkpage wraps the mkpage command from mkpage using the
# pandoc setup.
#
def mkpage(o_file, template = '', data = [], From = '', To = ''):
    cmd = ['mkpage']
    if From != '':
        cmd.append('-f')
        cmd.append(From)
    if To != '':
        cmd.append('-t')
        cmd.append('html')
    cmd.append('-o')
    cmd.append(o_file)
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
def assemble(cfg, html_filename, template_name, data):
    header = os.path.join(cfg.static, 'header.md')
    footer = os.path.join(cfg.static, 'footer.md')
    nav = os.path.join(cfg.static, 'nav.md')
    announcement = os.path.join(cfg.static, 'announcement.md')
    if cfg.base_url != '':
        data.append(kv('base_url', 'text', cfg.base_url))
    if cfg.base_path != '':
        data.append(kv('base_path', 'text', cfg.base_path))
    if cfg.control_item != '':
        data.append(kv('control_item', 'text', cfg.control_item))
    if os.path.exists(nav):
        data.append(kv('nav', '', nav))
    if os.path.exists(header):
        data.append(kv('header', '', header))
    if os.path.exists(footer):
        data.append(kv('footer', '', footer))
    if os.path.exists(announcement):
        data.append(kv('announcement', '', announcement))
    mkpage(os.path.join(cfg.htdocs, html_filename), os.path.join(cfg.templates, template_name), data)

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
def load_objects(cfg, json_name):
    f_name = os.path.join(cfg.htdocs, json_name)
    with open(f_name) as f:
        src = f.read()
        return json.loads(src)
    return None

def load_keys(cfg, key_name):
    f_name = os.path.join(cfg.htdocs, key_name)
    with open(f_name) as f:
        src = f.read()
        return src.strip().split('\n')
    return None

def build_top_level_pages(cfg):
    site_title, organization, site_welcome = cfg.site_title, cfg.organization, cfg.site_welcome
    print(f'build top level pages')
    css_folder = os.path.join(cfg.static, 'css')
    assets_folder = os.path.join(cfg.static, 'assets')
    js_folder = os.path.join(cfg.static, 'js')
    favicon = os.path.join(cfg.static, 'favicon.ico')
    copy_assets(cfg.htdocs, css_folder, f'{cfg.static}/')
    copy_assets(cfg.htdocs, assets_folder, f'{cfg.static}/')
    copy_assets(cfg.htdocs, js_folder, f'{cfg.static}/')
    copy_assets(cfg.htdocs, favicon, f'{cfg.static}/')

    # Render Homepage (index.html)
    page_title = f'{site_title}'
    title = site_welcome
    md_filename = os.path.join(cfg.static, 'index.md')
    html_filename = 'index.html'
    template_name = 'index-html.tmpl'
    assemble(cfg, html_filename, template_name, [
        kv('organization', 'text', organization),
        kv('site_title', 'text', site_title),
        kv('page_title', 'text', page_title),
        kv('title', 'text', title),
        kv('content', '', md_filename)
    ])

    # Render about page (information.html)
    page_title = f'{site_title}: About'
    title = 'About this Repository'
    md_filename = os.path.join(cfg.static, 'information.md')
    html_filename = 'information.html'
    template_name = 'index-html.tmpl'
    assemble(cfg, html_filename, template_name, [
        kv('organization', 'text', organization),
        kv('site_title', 'text', site_title),
        kv('page_title', 'text', page_title),
        kv('title', 'text', title),
        kv('content', '', md_filename)
    ])

    # Render browseviews page
    page_title = f'{site_title}: Browse'
    title = 'Browse this Repository'
    md_filename = os.path.join(cfg.static, 'browseviews.md')
    html_filename = 'browseviews.html'
    template_name = 'index-html.tmpl'
    assemble(cfg, html_filename, template_name, [
        kv('organization', 'text', organization),
        kv('site_title', 'text', site_title),
        kv('page_title', 'text', page_title),
        kv('title', 'text', title),
        kv('content', '', md_filename)
    ])

    # Render policies page
    page_title = f'{site_title}: Policies'
    title = 'Repository Policies'
    md_filename = os.path.join(cfg.static, 'policies.md')
    if os.path.exists(md_filename):
        html_filename = 'policies.html'
        template_name = 'index-html.tmpl'
        assemble(cfg, html_filename, template_name, [
            kv('organization', 'text', organization),
            kv('site_title', 'text', site_title),
            kv('page_title', 'text', page_title),
            kv('title', 'text', title),
            kv('content', '', md_filename)
        ])

    # Render contact page
    page_title = f'{site_title}: Contact'
    title = 'Contact Informaiton'
    md_filename = os.path.join(cfg.static, 'contact.md')
    if os.path.exists(md_filename):
        html_filename = 'contact.html'
        template_name = 'index-html.tmpl'
        assemble(cfg, html_filename, template_name, [
            kv('organization', 'text', organization),
            kv('site_title', 'text', site_title),
            kv('page_title', 'text', page_title),
            kv('title', 'text', title),
            kv('content', '', md_filename)
        ])


def make_view(cfg, view, label, site_title, organization):
    # Render view, e.g. /view/ids/
    page_title = f'{site_title}: Browse by {label}'
    title = f'Browse by {label}'
    list_data = f'view/{view}/{view}_list.json'
    html_filename = f'view/{view}/index.html'
    template_name = 'listing-html.tmpl'
    content = 'Please select a value to browse from the list below.'
    objs = load_objects(cfg, list_data)
    tot = len(objs)
    if tot == 0:
        content = 'Nothing available.'
    assemble(cfg, html_filename, template_name, [
        kv('organization', 'text', organization),
        kv('site_title', 'text', site_title),
        kv('page_title', 'text', page_title),
        kv('title', 'text', title),
        kv('content', 'text', content),
        kv('listing', '', os.path.join(cfg.htdocs, list_data))
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
        object_list = os.path.join(cfg.htdocs, 'view', view, f'{key}.json')
        html_filename = f'view/{view}/{key}.html'
        template_name = 'object-list-html.tmpl'
        assemble(cfg, html_filename, template_name, [
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


def build_view_pages(cfg, views):
    htdocs, site_title, organization = cfg.htdocs, cfg.site_title, cfg.organization
    print(f'building view pages')
    # Render view/
    page_title = f'{site_title}: Browse'
    title = 'Browse Items'
    md_filename = os.path.join(cfg.static, 'view_index.md')
    html_filename = 'view/index.html'
    template_name = 'index-html.tmpl'
    assemble(cfg, html_filename, template_name, [
        kv('organization', 'text', organization),
        kv('site_title', 'text', site_title),
        kv('page_title', 'text', page_title),
        kv('title', 'text', title),
        kv('content', '', md_filename)
    ])

    keys = views.get_keys()
    for key in keys:
        label = views.get_view(key)
        make_view(cfg, key, label, site_title, organization)
    print(f'build view pages completed')


def build_landing_pages(cfg):
    htdocs, site_title, organization = cfg.htdocs, cfg.site_title, cfg.organization
    print(f'build landing pages')
    objs = load_objects(cfg, 'view/ids/ids_list.json')
    keys = load_keys(cfg, 'index.keys')

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
        obj = load_objects(cfg, f'{key}/index.json')
        title =  obj['title']
        page_title = f'{title} - {site_title}'
        object_path = os.path.join(cfg.htdocs, f'{key}', f'index.json')
        html_filename = f'{key}/index.html'
        template_name = 'landing-page-html.tmpl'
        assemble(cfg, html_filename, template_name, [
            kv('organization', 'text', organization),
            kv('site_title', 'text', site_title),
            kv('page_title', 'text', page_title),
            kv('title', 'text', title),
            kv('object', '', object_path)
        ])
        bar.update(i)
    bar.finish()
    print(f'build landing pages completed')


def build_search_page(cfg):
    htdocs, site_title, organization = cfg.htdocs, cfg.site_title, cfg.organization
    print(f'build search page')
    data = []
    page_title = f'{site_title}: Search'
    title = f'Search'
    nav = os.path.join(cfg.static, 'nav.md')
    content = os.path.join(cfg.static, 'search_form.md')
    if site_title != '':
        data.append(kv('site_title', 'text', site_title))
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
    mkpage(os.path.join(cfg.htdocs, 'search.html'), os.path.join(cfg.templates, 'search-html.tmpl'), data)


def build_website(cfg):
    views = Views()
    views.load_views(cfg.views)
    build_top_level_pages(cfg)
    build_search_page(cfg)
    build_view_pages(cfg, views)
    build_landing_pages(cfg)
    

if __name__ == "__main__":
    f_name = ''
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
        
    cfg = Configuration()
    if cfg.load_config(f_name) and cfg.required(['htdocs', 'site_title', 'organization', 'control_item', 'views', 'site_welcome' ]):
        f_views, htdocs, site_title, organization, site_welcome = cfg.views, cfg.htdocs, cfg.site_title, cfg.organization, cfg.site_welcome
        if len(args) > 0:
            views = Views()
            views.load_views(f_views)
            for arg in args:
                if arg == 'top':
                    build_top_level_pages(cfg)
                    build_search_page(cfg)
                if arg == 'view':
                    build_view_pages(cfg, views)
                if arg == 'landing':
                    build_landing_pages(cfg)
        else:
            build_website(cfg)
    else:
        sys.exit(1)
