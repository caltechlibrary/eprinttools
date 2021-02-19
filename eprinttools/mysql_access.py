
#
# Get recently modified keys takes a URL MySQL connection string
# e.g. mysql://db_user:db_password@db_hostname/db_name
#
def get_recently_modified_keys(db_connection, no_of_days = -4):
    o = urlparse(db_connction)
    db_scheme = o.scheme
    db_host = o.hostname
    db_name = o.path.lstrip('/')
    db_user = o.username
    db_password = o.password
    # NOTE: SQL expects a positive number for interval value...
    if no_of_days < 0:
        no_of_days = no_of_days * -1
    sql_src = f'''SELECT eprintid AS eprint_id FROM eprint WHERE CONCAT_WS('-', LPAD(lastmod_year, 4, '0'), LPAD(lastmod_month, 2, '0'), LPAD(lastmod_day, 2, '0')) >= DATE_SUB(NOW(), INTERVAL {no_of_days} DAY) ORDER BY eprintid DESC'''
    #log.print(f'DEBUG SQL: {sql_src}')

    log.print('Updating my.cnf')
    with open('my.cnf', 'w') as f:
        f.write('[client]\n')
        f.write(f'host = {db_host}\n')
        f.write(f'user = {db_user}\n')
        f.write(f'password = {db_password}\n')
    os.chmod('my.cnf', stat.S_IWUSR | stat.S_IRUSR)

    # Verify that SQL file exists
    # Build up command and execute it
    cmd = [ "mysql" ]
    if os.path.exists('my.cnf'):
        my_cnf = os.path.abspath('my.cnf')
        cmd.append(f'--defaults-file={my_cnf}')
    else:
        if db_host != '':
            cmd.append('--host')
            cmd.append(db_host)
        if db_user != '':
            cmd.append('--user')
            cmd.append(db_user)
        if db_password != '':
            cmd.append(f'--password={db_password}')
    cmd.append('--default-character-set=utf8')
    cmd.append(f'--execute={sql_src}')
    cmd.append('--batch')

    if db_name:
        cmd.append(db_name)
    else:
        printf("Missing repository id db name")
        return []

    # NOTE: Now assemble and run the MySQL command to read in the file.
    eprint_ids = []
    with Popen(cmd, stdout = PIPE, encoding = 'utf-8') as proc:
        src = proc.stdout.read()
        for line in src.split("\n"):
            if line.isdigit():
                eprint_ids.append(line.strip(" \n\r"))
    return eprint_ids

