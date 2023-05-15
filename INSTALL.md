Installation
============

*eputil*, *epfmt*, *doi2eprintxml* are command line programs run from a shell like Bash. *ep2apid* is a web service that can also be run from the command line or setup as a "daemon" and run from systemd. They allow you to 
harvest, work with EPrint repository content, and import content from 
CrossRef and DataCite.

Quick install with curl
-----------------------

There is an experimental installer.sh script that can be run with the
following command to install lastest table release. This may work for
macOS, Linux and if you're using Windows with the Unix subsystem.

~~~
curl https://caltechlibrary.github.io/eprinttools/installer.sh | sh
~~~

Below are generalized instructions for installation of a release.

Compiled version
----------------

Compiled versions are available for macOS (Intel and M1 processors as macos-amd64 or macos-arm64), Linux (amd64 process, linux-amd64), Windows (amd64 and arm64 processor, windows-amd64 and windows-arm64) and Rapsberry Pi (arm7 processor, raspbian-arm7)

VERSION\_NUMBER is a [symantic version number](http://semver.org/) (e.g.
`v1.3.0`)

For all the released version go to the project page on Github and click latest release

> <https://github.com/caltechlibrary/eprinttools/releases/latest>

| Platform    | Zip Filename                                 |
|-------------|----------------------------------------------|
| Windows     | eprinttools-VERSION_NUMBER-windows-amd64.zip |
| Windows     | eprinttools-VERSION_NUMBER-windows-arm64.zip |
| macOS       | eprinttools-VERSION_NUMBER-macos-amd64.zip  |
| macOS       | eprinttools-VERSION_NUMBER-macos-arm64.zip  |
| Linux/Intel | eprinttools-VERSION_NUMBER-linux-amd64.zip   |
| Raspbery Pi | eprinttools-VERSION_NUMBER-raspbian-os-arm7.zip |

The basic recipe
----------------

- Find the Zip file listed matching the architecture you're running
  and download it
      - (e.g. if you're on a Windows 10 laptop/Surface with a amd64
        style CPU you'd choose the Zip file with "windows-amd64" in the
        name).
- Download the zip file and unzip the file.
- Copy the contents of the folder named "bin" to a folder that is in
  your path
      - (e.g. "\$HOME/bin" is common).
- Adjust your PATH if needed
      - (e.g. `export PATH="\$HOME/bin:\$PATH"`)
- Test

### Mac OS X

1.  Download the zip file
2.  Unzip the zip file
3.  Copy the executables to $HOME/bin (or a folder in your path)
4.  Make sure the new location in in our path
5.  Test

Here's an example of the commands run in the Terminal App after downloading the zip file.

#### Intel Hardware

``` shell
    cd Downloads/
    unzip eprinttools-*-macos-amd64.zip
    mkdir -p $HOME/bin
    cp -v bin/* $HOME/bin/
    export PATH=$HOME/bin:$PATH
    eputil -version
```

#### M1 (ARM64) Hardware

``` shell
    cd Downloads/
    unzip eprinttools-*-macos-arm64.zip
    mkdir -p $HOME/bin
    cp -v bin/* $HOME/bin/
    export PATH=$HOME/bin:$PATH
    eputil -version
```


### Windows

1.  Download the zip file
2.  Unzip the zip file
3.  Copy the executables to $HOME/bin (or a folder in your path)
4.  Test

Here's an example of the commands run in from the Bash shell on Windows 10 after downloading the zip file.

#### Intel (amd64) Hardware

``` shell
    cd Downloads/
    unzip eprinttools-*-windows-amd64.zip
    mkdir -p $HOME/bin
    cp -v bin/* $HOME/bin/
    export PATH=$HOME/bin:$PATH
    eputil -version
```


#### ARM64 Hardware

``` shell
    cd Downloads/
    unzip eprinttools-*-windows-arm64.zip
    mkdir -p $HOME/bin
    cp -v bin/* $HOME/bin/
    export PATH=$HOME/bin:$PATH
    eputil -version
```


### Linux

1.  Download the zip file
2.  Unzip the zip file
3.  Copy the executables to $HOME/bin (or a folder in your path)
4.  Test

Here's an example of the commands run in from the Bash shell after downloading the zip file.

``` shell
    cd Downloads/
    unzip eprinttools-*-linux-amd64.zip
    mkdir -p $HOME/bin
    cp -v bin/* $HOME/bin/
    export PATH=$HOME/bin:$PATH
    eputil -version
```

### Raspberry Pi

Released version is for a Raspberry Pi 2 or later use (i.e. requires ARM
7 support).

1.  Download the zip file
2.  Unzip the zip file
3.  Copy the executables to $HOME/bin (or a folder in your path)
4.  Test

Here's an example of the commands run in from the Bash shell after downloading the zip file.

``` shell
    cd Downloads/
    unzip eprinttools-*-raspbian-arm7.zip
    mkdir -p $HOME/bin
    cp -v bin/* $HOME/bin/
    export PATH=$HOME/bin:$PATH
    eputil -version
```


Compiling from source
---------------------

*eprinttools* is "go gettable". Use the "go get" command to download the dependant packages as well as *eprinttools*'s source code.

``` shell
    go get -u github.com/caltechlibrary/eprinttools/...
```

Or clone the repstory and then compile

``` shell
    cd
    git clone https://github.com/caltechlibrary/eprinttools \
        src/github.com/caltechlibrary/eprinttools
    cd src/github.com/caltechlibrary/eprinttools
    make
    make test
    make install
```
