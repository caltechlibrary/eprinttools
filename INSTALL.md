Installation for development of **eprinttools**
===========================================

**eprinttools** Command line tools, services, Golang package and Python module for working with the EPrints 3.x REST API

Quick install with curl or irm
------------------------------

There is an experimental installer.sh script that can be run with the following command to install latest table release. This may work for macOS, Linux and if you’re using Windows with the Unix subsystem. This would be run from your shell (e.g. Terminal on macOS).

~~~shell
curl https://caltechlibrary.github.io/eprinttools/installer.sh | sh
~~~

This will install the programs included in eprinttools in your `$HOME/bin` directory.

If you are running Windows 10 or 11 use the Powershell command below.

~~~ps1
irm https://caltechlibrary.github.io//installer.ps1 | iex
~~~

Installing from source
----------------------

### Required software

- Go &gt;&#x3D; 1.24.2
- GNU Make &gt; 3
- Pandoc &gt; 3
- CMTools &gt;&#x3D; 0.0.23

### Steps

1. git clone https://github.com/caltechlibrary/eprinttools
2. Change directory into the `eprinttools` directory
3. Make to build, test and install

~~~shell
git clone https://github.com/caltechlibrary/eprinttools
cd eprinttools
make
make test
make install
~~~

