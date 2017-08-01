User Manual
===========

This document is a detailed documentation on NeON. For a quick overview, see
[Quick Start](quickstart.md). If you are looking for a documentation on tasks
and builtins, see [Reference](reference.md) documentation.

**Table of Contents**

- [Build file format](#build-file-format)
- [Build file structure](#build-file-structure)
- [Build properties](#build-properties)
- [Build targets](#build-targets)
  - [NeON tasks](#neon-tasks)
  - [Shell tasks](#shell-tasks)
  - [Script tasks](#script-tasks)

Build file format
-----------------

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

Build file structure
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

Build targets
-------------

Build targets might be compared to functions. They are called on command line.
If you call NeON with `neon clean`, you will run target *clean*.

This target might look like:

```yaml
targets:

  clean:
    doc: Clean generated files
    steps:
    - delete: 'build'
```

A target might define following fields:

- **doc** this is the target documentation.
- **depends** to list targets before running this one.
- **steps** is the list of the tasks to run the target.

Tasks might be one of the following:

### NeON tasks

They are tasks defined in NeON engine. This is a way to write platform
independant build files. These tasks are maps with string keys.

There are tasks to manage files (copy, delete, move and so on), archives
(create ZIP or TAR files), directories (create, change to) or links. For
instance, to delete all *.so* files in *build* directory, you would write:

```yaml
targets:

  delete:
    doc: Delete object files
    steps:
    - delete: '**/*.so'
      dir:    'build'
```

Ther are also logical tasks to perform tests and iterate on values. For
instance, to iterate on a list of files, you could write:

```yaml
targets:

  pdf:
    doc: Generate PDF files
    steps:
    - for: 'file'
      in:  'find("md", "*.md")'
      do:
      - $: 'md2pdf -o "build/#{file}.pdf" "md/#{file}"'
```

To generate a file if the source is newer, you would write:

```yaml
targets:

  pdf:
    doc: Generate PDF file
    steps:
    - if: 'older("file.md", "build/file.pdf")'
      then:
      - $: 'md2pdf -o "build/file.pdf" "file.md"'
```

There are also tasks to manage errors. For instance to run a command and catch
any error (that is when the command returns a value different from *0*) to
write an error message, you could write:

```yaml
targets:

  command:
    doc: Try to run a command
    steps:
    - try:
      - $: 'command that might fail'
      catch:
      - throw: 'There was an error running command'
```

To list all available NeON tasks, type command `neon -tasks`. To get help on a
given command *foo*, type `neon -task foo`:

```
$ n -task time
Record duration to run a block of steps.

Arguments:

- time: the steps to measure execution duration.
- to: the property to store duration in seconds as a float (optional,
  print duration on console if not set).

Examples:

    # print duration to say hello
    - time:
      - print: "Hello World!"
      to: duration
    - print: 'duration: #{duration}s'
```

You can get information on available tasks
[on this reference page](reference.md).

### Shell task

A shell task runs a script. This script will run with *sh* on Unix and
*cmd.exe* on Windows. They are a map with *$* field.

This script might be a simple command such as *ls* or it may be a full shell or
batch script. In this case, you should the appropriate YAML syntax:

```yaml
targets:

  shell:
    steps:
    - $: |
         This is a shell script
         with more that one line
```

A shell task will fail if the script returns a value different from *0*. You
might manage errors with *try* task.

Of course a command might be system dependant, but this is not always the case.
For instance, a command such as `java -jar foo.jar` will probably run the same
on all systems. This is also the case for most Git commands.

### Script task

A script task is a piece of code that will run in the NeON scripting engine.
This is also a way to write platform independant code. But this is a way to
write complex scripts that would be complicated to write with system commands.

A script task is a simple string. For instance:

```yaml
targets:

  script:
    steps:
    - 'file = joinpath(BUILD_DIR, "test.txt")'
```

Your scripts might use builtin funtions, defined by Anko scripting engine (such
as *toString()*) or by NeON. To lis NeON builtins, you can type command 
`neon -builtins`. To get help on given builtin, type `neon -builtin split`.
You can get information on available builtin functions
[on this reference page](reference.md).

You can define your own functions in scripts that you load in the build file
with *context* field.

To get more information about 
[Anko scripting language clic here](http://github.com/mattn/anko).

