# pokenv [![Build status](https://ci.appveyor.com/api/projects/status/hrtwo6hrx10d7i88?svg=true)](https://ci.appveyor.com/project/tischda/pokenv)

Windows utility written in [Go](https://www.golang.org) to poke
environment variables into the registry.

### Compile

Tested with GO 1.4.2. There are no dependencies.

~~~
go build
~~~

### Usage

~~~
Usage: pokenv [-checkpaths] [-hkcu|-hklm] infile
  infile: the input file
  -checkpaths=false: values are paths, check that they are valid on this system
  -hkcu="REQUIRED": process input file into HKEY_CURRENT_USER environment
  -hklm="REQUIRED": process input file into HKEY_LOCAL_MACHINE environment
  -version=false: print version and exit
~~~

### Other setters

* http://sourceforge.net/projects/pathmanager/files/?source=navbar
* http://www.rapidee.com/en/about
* http://p-nand-q.com/download/gtools/pathed.html
* http://ss64.com/nt/setx.html
