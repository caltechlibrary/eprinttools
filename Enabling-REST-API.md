Enabling REST API for EPrints
=============================

These are just my quick notes for enabling the REST API for EPrints. 

The REST API doesn't "automatically" become enabled even for Admin role users. You can alter this behavior by updating the roles in "archives/REPOSITORY_NAME/cfg/cfg.d/user_roles.pl" (where REPOSITORY_NAME is the name of the respository you setup with _epadmin creeate_) in your eprints directory.

Below is I added "rest" role to the admin role.

```perl
    $c->{user_roles}->{admin} = [qw{
	    rest
	    general
	    edit-own-record
	    saved-searches
	    set-password
	    deposit
	    change-email
	    editor
	    view-status
	    staff-view
	    admin
	    edit-config
    }];
```

## eputil

__eputil__ supports POST and PUT into EPrint's REST API. Content sent by the POST or PUT is assumed to be encoded before it is read from a file or standard input. In the example below the base we are "putting" the value (TRUE) into the lemurprints.local/authors EPrint collection for record 1's referreed field.

```shell
    echo -n "TRUE" | eputil -u "$EP_USER" -p "$EP_PASSWORD" \
        -put http://lemurprints.local/authors/rest/eprint/1/refereed.txt
```


Reference links
---------------

+ [REST API Feature Announcement](http://wiki.eprints.org/w/New_Features_in_EPrints_3.2)
+ [EPrints XML Configuration](https://wiki.eprints.org/w/EPScript) - need to enable REST API access based on role
+ [API:EPrints/Apache/CRUD](http://wiki.eprints.org/w/API:EPrints/Apache/CRUD)
+ [user_roles.pl](https://wiki.eprints.org/w/User_roles.pl)

+ EPrints Tech list archives mention REST API
    + [2012 October](http://mailman.ecs.soton.ac.uk/pipermail/eprints-tech/2012-October/001176.html) - first mention in the archives (Re: AJAX end point)
    + [2013 April](http://mailman.ecs.soton.ac.uk/pipermail/eprints-tech/2013-April/001809.html) - another mention "Edit, Update and delter report from third party tool"
    + [2013 January](http://mailman.ecs.soton.ac.uk/pipermail/eprints-tech/2013-January/)
        + [EPrints REST API documentation](http://mailman.ecs.soton.ac.uk/pipermail/eprints-tech/2013-January/001462.html) -- problem of setup
        + [Re: EPrints REST API documentation](http://mailman.ecs.soton.ac.uk/pipermail/eprints-tech/2013-January/001465.html) -- first helpful response indicating roles need to be enabled
    + [2017 March](http://mailman.ecs.soton.ac.uk/pipermail/eprints-tech/2017-March/)
        + [Question about REST API, review items under EPrints 3.3](http://mailman.ecs.soton.ac.uk/pipermail/eprints-tech/2017-March/006346.html) - Caltech Libray EPrints question
    + [2017 May](http://mailman.ecs.soton.ac.uk/pipermail/eprints-tech/2017-May/)
        + [Bulk updating questions](http://mailman.ecs.soton.ac.uk/pipermail/eprints-tech/2017-May/006516.html)

