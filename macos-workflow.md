Notes on using doi2eprintxml on Mac OS X
========================================

NOTE: In recent (version greater than 13.3) versions of macOS, Apple has enhanced their security policies such that it is more difficult to install open source software like eprinttools.  To allow software like eprinttools on your Mac you need to author installation of unsigned software when using the Mac Terminal application.  Go to your system settings, click on "Privacy & Secutiry", then click on "Developer Tools". You should see the Terminal listed. Set the setting so you are allowed to install unsigned software in the Terminal.  You can switch this back off after run the software installed installed.

Installation and setup
----------------------

1. With your web browser goto https://github.com/caltechlibrary/eprinttools/releases/
2. Download the file called "eprinttools-v1.3.0-macos-amd64.zip" (for older Intel based Macs) or "eprinttools-v1.3.0-macos-arm64.zip" (for newer M1 based Macs).
    + if there is a newer version (e.g. v1.3.1, v1.4.0) download that instead
3. Open the Terminal application and change to directory to the Downloads folder
4. Unzip the downloaded zip file
5. Create a home "bin" folder if you haven't in the past
6. Move "bin/doi2eprintxml" to your local **bin** folder 
    (e.g. /Users/rsdoiel/bin)
7. Make sure your **bin** directory is in your path
    (e.g. `export PATH="$HOME/bin:$PATH"` in your **.bashrc** file if your shell is Bash or **.zshrc** if your shell is zsh)
8. Close the Terminl application and then restart it.
9. Run `doi2printxml -h` to see the help page

Here's an example of the commands I after downloading the zip file to the "Downloads" folder on my M1 Mac Mini (steps 4 through 6).

~~~
cd
cd Downloads
unzip eprinttools-v1.3.0.macos-arm64.zip
mkdir -p $HOME/bin
mv bin/doi2eprintxml $HOME/bin/
~~~

If you've not previously add your home "bin" directory to your path
and you're using the new default shell on macOS (i.e. zsh). Enter the
following command.

~~~
echo 'export PATH="$HOME/bin:$PATH"' >>$HOME/.zshrc
~~~

If you're using Bash the command is.

~~~
echo 'export PATH="$HOME/bin:$PATH"' >>$HOME/.bashrc
~~~

Close the Terminal application and restart it then run step 9.



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
using a text editor such as [Atom](https://atom.io) with one DOI 
per line. The `doi2eprintxml` program can then read the list and 
generate an EPrint XML document for importing into EPrints. Follow 
the following steps to generate a file called **eprint.xml** based 
on a list of of DOI in a file called **doi.txt**.

#### General recipe

1. Open your text editor (e.g Atom)
2. With your text editor create a list of DOI on per line 
3. Save this list to your desktop as **doi.txt**
    NOTE: **doi.txt** is the name of your file you saved on your desktop.
4. Open the Terminal app and change directory to your desktop
    `cd ~/Desktop`
5. Generate **eprint.xml** with the following command
    `doi2eprintxml -i doi.txt > eprint.xml`
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

