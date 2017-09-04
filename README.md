NeON
====

[![Build Status](https://travis-ci.org/c4s4/neon.svg?branch=master)](https://travis-ci.org/c4s4/neon)
[![Code Quality](https://goreportcard.com/badge/github.com/c4s4/neon)](https://goreportcard.com/report/github.com/c4s4/neon)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
<!--
[![Coverage Report](https://coveralls.io/repos/github/c4s4/neon/badge.svg?branch=master)](https://coveralls.io/github/c4s4/neon?branch=master)
-->

- Project :   <https://github.com/c4s4/neon>.
- Downloads : <https://github.com/c4s4/neon/releases>.

NeON is a build tool the way it should be.

Installation
------------

Download latest binary archive at <https://github.com/c4s4/neon/releases>. Unzip
the archive, put the binary of your platform somewhere in your *PATH* and rename
it *neon*.

Unix users can add Bash completion putting file *bash/neon* in
*/etc/bash_completion.d/* directory and adding line
`. /etc/bash_completion.d/neon` in your *~/.bashrc* file. This will enable
following completions:

- Typing `neon ` and hitting TAB will complete on build targets.
- Typing `neon -task ` and hitting TAB will complete on tasks.
- Typing `neon -builtin ` and hitting TAB will complete on builtins.

Documentation
-------------

There are three sources of documentation:

- [The quick start guide](doc/quickstart.md).
- [The user manual](doc/usermanual.md)
- [The reference guide](doc/reference.md).

Build
-----

To build the project without *NeON* already installed, follow these steps:

- Define `GOPATH` environment variable on the directory of the project.
- Write *src/neon/version.go* file with following content :

```go
package main

var VERSION = "UNKNOWN"
```

- Go to the project directory and type command `go build neon`.

This will product the *neon* binary for your platform and OS in the project
directory.

*Enjoy!*
