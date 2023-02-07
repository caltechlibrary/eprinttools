Notes on using doi2eprintxml.exe on Windows 10
==============================================

Installation and setup
----------------------

1. With your web browser goto https://github.com/caltechlibrary/eprinttools/releases/
2. Download the file called "eprinttools-v1.2.2-windows-amd64.zip".
    + if there is a newer version (e.g. v1.2.2) download that instead
3. Unzip the downloaded zip file
4. Copy "doi2eprintxml.exe" to where you want to work (e.g. Desktop)
5. Open the windows "command prompt" and change directory to your desktop
6. Run `doi2printxml.exe -h` to see the help page

Getting started
---------------

The tool `doi2eprintxml.exe` is a command line program, it is normally
run from the command prompt. The purpose of the program is to take
a list of DOI contact CrossRef or DataCite and generate what is called
an EPrint XML document. The EPrint XML document is used to import
records into EPrints (e.g. CaltechAUTHORS) via the "import" button 
under "Manage Deposits" page.

### Workflow

The typical workflow is to create a plain text file (UTF-8 encoded)
using a text editor such as [VSCode](https://code.visualstudio.com/) with one DOI per 
line. The `doi2eprintxml.exe` file can then read the list and 
generate an EPrint XML document for importing into EPrints. Follow 
the following steps to generate a file called **eprint.xml** based 
on a list of of DOI in a file called **doi.txt**.

### General recipe

1. Open your text editor (e.g VSCode)
2. With your text editor create a list of DOI on per line 
3. Save this list to your desktop as **doi.txt**
    NOTE: **doi.txt** is the name of your file you saved on your desktop.
4. Open the Windows command prompt and change directory to your desktop
5. Generate **eprint.xml** with the following command
    `doi2eprintxml.exe -i doi.txt > eprint.xml`
6. From the file manager open **eprint.xml** and review the XML document (e.g. in your web browser or other text editor)
7. With your web browser go to CaltechAUTHORS and log in
8. Change to the "manage deposits" page and select "EPrints XML" from the pull down menu next to the "import" button
9. Press the "import" button
10. Click on "Upload from file"
11. Click the "browse" button and find **eprint.xml** on the desktop
12. From the pull down list next to the filename select "UTF-8"
13. Click "Test without importing", proceed to the next step if everything is OK
14. Click the "browse" button and find **eprint.xml** on the desktop
15. From the pulldown list next to the filename select "UTF-8"
16. Click on "import items" buttons, this should import the items into your buffer to review

