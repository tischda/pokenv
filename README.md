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
Usage: pokenv [options] infile
  infile: the input file
  -hkcu=false: set HKEY_CURRENT_USER environment
  -hklm=false: set HKEY_LOCAL_MACHINE environment
  -version=false: print version
~~~

### Other setters

* http://sourceforge.net/projects/pathmanager/files/?source=navbar
* http://www.rapidee.com/en/about
* http://p-nand-q.com/download/gtools/pathed.html
* http://ss64.com/nt/setx.html
