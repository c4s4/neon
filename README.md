NeON
====

- Project :   <https://github.com/c4s4/changelog>.
- Downloads : <https://github.com/c4s4/changelog/releases>.

NeON is a build tool the way it should be.

Installation
------------

Download latest binary archive at <https://github.com/c4s4/neon/releases>. Unzip
the archive, put the binary of your platform somewhere in your *PATH* and rename
it *neon*.

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
