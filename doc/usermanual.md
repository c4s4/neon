User Manual
===========

This document is a detailed documentation on NeON. For a quick overview, see
[Quick Start](quickstart.md). If you are looking for a documentation on tasks
and builtins, see [Reference](reference.md) documentation.

The build file format
---------------------

NeON build files are in YAML format. YAML (for *Yet Another Markup Language*)
is a lightweight markup language so that you can write structured text in a
natural syntax, with no strange tags or other syntaxic form of torture. For
instance, to describe a list, you would write:

```yaml
- Bilbo
- Frodo
- Gandalf
```
To write a dictionary, you would type:

```yaml
gandalf: istari
bilbo: hobbit
frodo: hobbit
galadriel: elf
```

And you can combine. For instance, a dictionary of lists:

```yaml
elves:
  - Galadriel
  - Arwen
hobbits:
  - Bilbo
  - Frodo
```

There is a compact notation for lists and maps:

```yaml
list: [Bilbo, Frodo]
map:  {good: Frodo, bad: Sauron}
```

Furthermore, data are typed:

```yaml
integers:
  - 1
  - 123
floats:
  - 1.2
  - 3.14
dates:
  - 2015-10-21
strings:
  - This is text
  - "1"
  - '2015-10-21'
```

So, numbers are numbers, dates are ISO formatted and you can force anything to
be a string surrounding it with quotes (simple or double).

Things to know to avoid troubles writing YAML:

- You **cannot indent with tabulations**, this is a syntax error!
- Strings with colons must be surrounded with quotes or they are considered
  maps. Thus you would write string *"see: this is cool"* with quotes if you
  don't want it to be a map.

This introduction to YAML should be enough for you to write valid build files.
If you want more information about YAML, please visit
[YAML website](http://yaml.org).

Build File Structure
--------------------

A build file is a YAML map. First level fields are the following:

- **doc** is the documentation of the build file. This is a string.
- **default** is for default target, which will run if none is passed on
  command line. This might be a string or a list of strings. If no default
  targets are defined, you must pass a target on command line.
- **extends** is the list of extended build files. See inheritance section for
  more details. This is a string or a list of strings.
- **repository** is the default location for parent build files. Defaults to
  directory *.neon* in the home directory of the user.
- **context** is a list of scripts loaded on startup. This is a string or a
  list of strings.
- **singleton** is a port that is opened on startup to ensure that only a
  single instance of the build is running. This is an integer. This should be
  between *1024* and *65535* (or between *1* and *65535* if the build is
  running as root).
- **properties** is a map of properties of the build file. See section
  *Properties* for more information about build properties.
- **configuration** is a list of YAML files to load as build properties.
  These YAML files must be maps with string keys. This might be a string or a
  list of strings.
- **environment** is a map that defines environment for all executed commands.
- **targets** is a map for targets of the build files. This is a map with
  string keys.

Most build files will define documentation, default target, properties and
targets. Thus a simple build file might look like following:

```yaml
doc: This is a sample build file

default: test

properties:
  NAME:      '#{filename(_BASE)}'
  BUILD_DIR: 'build'

targets:

  test:
    doc: Run Go tests
    steps:
    - $: 'go test'

  run:
    doc: Run Goo tool
    steps:
    - $: 'go run "#{NAME}.go"'

  bin:
    doc: Build Go tool
    steps:
    - mkdir: '#{BUILD_DIR}'
    - $: 'go build -o #{BUILD_DIR}/#{NAME}'

  clean:
    doc: Clean generated files
    steps:
    - delete: '#{BUILD_DIR}'
```

Build file properties
---------------------

Build file properties are some sort of build variables. The properties section
is a map with string keys. For instance:

```yaml
properties:
  STRING:  'This is a string'
  OTHER:   '1'
  INTEGER: 42
  FLOAT:   4.2
  LIST:    [1, 1, 2, 3, 5, 8]
  MAP:
    one:   1
    two:   2
    three: 3
```

A build property might use the value of another one. For instance:

```yaml
properties:
  NAME:      'test'
  BUILD_DIR: 'build'
  ARCHIVE:   '#{BUILD_DIR}/#{NAME}.zip'
```

The *ARCHIVE* property uses values of *BUILD_DIR* and *NAME*. Note that the
order of properties is not important as maps are not ordered.

To avoid errors, you should follow these conventions:

- Uppercase properties are constants, defined in properties field.
- Lowercase properties are local variables. Note that they are defined in the
  whole build file you should alway define their value in the current target.
- Properties starting with underscores (such as *_error*) are internal
  variables, defined by NeON. They should not be modified.

These properties are defined in the virtual machine that runs scripts. This
scripting language is [Anko](http://github.com/mattn/anko), which is a kind of
scripted Go.

Thus you can also define and user build properties in scripts. for instance,
you might write:

```yaml
properties:
  BUILD_DIR: 'build'

targets:

  test:
    steps:
    - 'file = joinpath(BUILD_DIR, "test.txt")'
```

Note that to use property *BUILD_DIR*, you write `BUILD_DIR` and not
`#{BUILD_DIR}`. The expression `#{BUILD_DIR}` is used to insert a property
value in a string, not a script.

Note that some tasks define internal properties. For instance, task *try*
will store raised error in internal build property *_error*.

All YAML types might be used to define build properties. Thus, you can define
string, integers, floats, but also lists and maps. You may iterate on values
of a property in the build file.
