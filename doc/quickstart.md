Neon Quick Start
================

Neon is a build tool.

Installation
------------

Download latest binary archive at <https://github.com/c4s4/neon/releases>. Unzip the archive, put the binary of your platform somewhere in your *PATH* and rename it *neon*.

Usage
-----

To run a build, type on command line:

```bash
$ neon
```

This will launch default target for this build. To run target *foo*, you should type:

```bash
$ neon foo
```

You may pass more than one target on the command line.

To get help on current build file, you can type:

```bash
$ neon -build
Build file to build neon

Properties:
arc_dir   "build/neon-0.1.0" 
build_dir "build" 
name      "neon" 
version   "0.1.0" 

Targets:
archive Generate distribution archive [clean]
build   Build neon binary 
clean   Clean generated files 
deps    Install libraries 
release Perform a release [clean, test, archive]
test    Run unit tests
```

To get help on *neon* usage, you can type:

```bash
$ neon -help
Usage of neon:
  -build
    	Print build help
  -debug
    	Output debugging information
  -file string
    	Build file to run (default "build.yml")
```

*Enjoy!*
