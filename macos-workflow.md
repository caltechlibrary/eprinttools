Notes on using doi2eprintxml on Mac OS X
========================================

NOTE: In recent (version greater than 13.3) versions of macOS, Apple has enhanced their security policies such that it is more difficult to install open source software like eprinttools.  To allow software like eprinttools on your Mac you need to author installation of unsigned software when using the Mac Terminal application.  Go to your system settings, click on "Privacy & Secutiry", then click on "Developer Tools". You should see the Terminal listed. Set the setting so you are allowed to install unsigned software in the Terminal.  You can switch this back off after run the software installed installed.

Installation and setup
----------------------

1. Open the Terminal application 
2. Type the "curl" command below
3. Follow the instructions it provides
4. Run the "open" command
5. In the finder window, then right click on `doi2eprintxml` and click "Open" in the dialog
6. Close the Terminal and reopen it, you should now be ready to use doi2eprintxml

Steps 2 through 4

~~~
curl https://caltechlibrary.github.io/eprinttools/installer.sh
open bin
~~~

Step six close the Terminal window and reopen it. You should now be ready
to use `doi2eprintxml`.


Getting started
---------------

The tool `doi2eprintxml` is a command line program, it is normally
run from the terminal's command prompt. The purpose of the program 
is to take a list of DOI contact CrossRef or DataCite and generate 
what is called an EPrint XML document. The EPrint XML document is 
used to import records into EPrints (e.g. CaltechAUTHORS) via 
the "import" button under "Manage Deposits" page.

### Workflow

The typical workflow is to create a plain text file (UTF-8 encoded)
using a text editor such as [VSCode](https://code.visualstudio.com/) or 
[Zed](https://zed.dev/) with one DOI per line. The `doi2eprintxml` 
program can then read the list and generate an EPrint XML document
for importing into EPrints. Follow the following steps to generate 
a file called **eprint.xml** based on a list of of DOI in a file 
called **doi.txt**. In the example below `code` is the VS code editor,
replace that line with your favorite text editor.

~~~
cd ~/Desktop
code doi.txt
doi2eprintxml -i doi.txt >eprints.xml
~~~

#### Generating eprints XML for a single DOI

In this example the DOI is  "10.1021/acsami.7b15651".

~~~
doi2eprintxml "10.1021/acsami.7b15651" > article.xml
~~~

For details of using `doi2eprintxml` see the [manual page](https://caltechlibrary.github.io/eprinttools/doi2eprintxml.1.html)

