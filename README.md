NeON
====

[![Build Status](https://travis-ci.org/c4s4/neon.svg?branch=master)](https://travis-ci.org/c4s4/neon)
[![Code Quality](https://goreportcard.com/badge/github.com/c4s4/neon)](https://goreportcard.com/report/github.com/c4s4/neon)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Codecov](https://codecov.io/gh/c4s4/neon/branch/master/graph/badge.svg)](https://codecov.io/gh/c4s4/neon)

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

There are four sources of documentation:

- [The quick start guide](doc/quickstart.md).
- [The user manual](doc/usermanual.md)
- [The reference guide](doc/reference.md).
- [Slides in French](http://sweetohm.net/slides/slides-neon)

Build
-----

To build the project without *NeON* already installed, follow these steps:

- Define `GOPATH` environment variable on the directory of the project.
    ```bash
    $ export GOPATH=`pwd`
    ```
- Go to the project directory and type command `go build neon`.

This will product the *neon* binary for your platform and OS in the project
directory.

If neon is already installed, simply type `neon bin`, this will generate neon
binary in *bin* directory.

Contributors
------------

Active contributors are:

- [Michel Casabianca](mailto:casa@sweetohm.net)
- [Alexandre Hu](mailto:a.hu@dalloz.fr)

Please feel free to contribute and send your patches or pull requests, they
will be reviewed and integrated as soon as possible.

*Enjoy!*
