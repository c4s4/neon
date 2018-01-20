User Manual
===========

This document is a detailed documentation on NeON. For a quick overview, see
[Quick Start](quickstart.md). If you are looking for a documentation on tasks
and builtins, see [Reference](reference.md) documentation.

**Table of Contents**

- [Build file format](#build-file-format)
- [Build file structure](#build-file-structure)
- [Build properties](#build-properties)
  - [Referencing build properties](#referencing-build-properties)
  - [Predefined build properties](#predefined-build-properties)
  - [Build properties on command line](#build-properties-on-command-line)
  - [Configuration](#configuration)
  - [Properties hierarchy](#properties-hierarchy)
- [Build targets](#build-targets)
  - [NeON task](#neon-task)
    - [File tasks](#file-tasks)
  - [Shell task](#shell-task)
  - [Script task](#script-task)
- [Command line options](#command-line-options)
- [Build inheritance](#build-inheritance)
- [NeON repository](#neon-repository)
- [Project templates](#project-templates)
  - [Creating templates](#creating-templates)

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

[Back to top](#user-manual)

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
  list of strings. This is useful to define your own builtin functions.
- **singleton** is a port that is opened on startup to ensure that only a
  single instance of the build is running. This is an integer. This should be
  between *1024* and *65535* (or between *1* and *65535* if the build is
  running as root).
- **shell** defines the command to execute to run a shell task. This command is
  defined as a list (such as `["sh", "-c"]` or `["cmd", "/c"]` for instance).
  The command to run in the shell will be added as the last argument.
  If you define a single command, this will run for all operating systems. You
  may instead define the shell as a map of command per operating system. For
  instance:
	```yaml
	shell:
	  windows: ['cmd', '/c']
	  default: ['sh', '-c']
	```
  This will define a shell for windows and for other environments. Thus, if
  you deine a single command (say `['sh', '-c']`), this is equivalent to:
	```yaml
	shell:
	  default: ['sh', '-c']
	```
  Note that commands defined as lists will not run with these shells, thus
  options won't be evaluated and `${USER}` will stay as is.
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
  NAME:      '=filename(_BASE)'
  BUILD_DIR: 'build'

targets:

  test:
    doc: Run Go tests
    steps:
    - $: ['go', 'test']

  run:
    doc: Run Goo tool
    steps:
    - $: ['go', 'run', '={NAME}.go']

  bin:
    doc: Build Go tool
    steps:
    - mkdir: '=BUILD_DIR'
    - $: ['go', 'build', '-o', '={BUILD_DIR}/={NAME}']

  clean:
    doc: Clean generated files
    steps:
    - delete: '=BUILD_DIR'
```

[Back to top](#user-manual)

Build properties
----------------

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
  ARCHIVE:   '={BUILD_DIR}/={NAME}.zip'
```

The *ARCHIVE* property uses values of *BUILD_DIR* and *NAME*. Note that the
order of properties is not important as maps are not ordered.

To avoid errors, you should follow these conventions:

- Uppercase properties are constants, defined in properties field.
- Lowercase properties are local variables. Note that they are defined in the
  whole build file, but you should always define their value in current target.
- Properties starting with underscores (such as *_error*) are internal
  variables, defined by NeON. They should not be modified unless you know what
  you are doing.

These properties are defined in the virtual machine that runs scripts. This
scripting language is [Anko](http://github.com/mattn/anko), which is a kind of
scripted Go.

### Referencing build properties

You can reference these properties in a string with the expression
`={PROP_NAME}`. This might be done in the expression of other properties but
also in task fields. For instance:

```yaml
- print: 'Hello ={USER}!'
```

You can also get directly the value of a property with the expression
`=PROP_NAME` (without curly braces). Thus you might write:

```yaml
- print: =USER
```

In this case, if property *USER* is a string, this will not be so different
from preceding example. But this is quite different if the property is not a
string.

For instance, let's see this build file:

```yaml
properties:
  FILES: ['foo.txt', 'bar.txt']

targets:

    test:
      steps:
      - for: file
        in:  =FILES
        do:
        - print: =file
```

In this case expression `=FILES` returns a list that we iterate in *test*
target.

You can also define and use build properties in scripts. for instance, you
might write:

```yaml
properties:
  BUILD_DIR: 'build'

targets:

  test:
    steps:
    - 'file = joinpath(BUILD_DIR, "test.txt")'
```

This would define a build property named *file* that will be a string with
value *"build/test.txt"*.

Note that to use property *BUILD_DIR*, you write `BUILD_DIR` and not
`={BUILD_DIR}` or `=BUILD_DIR`. The expression `={BUILD_DIR}` is used to insert
a property value in a string and `=BUILD_DIR` to get the property value in a
task field.

Note that some tasks define internal properties. For instance, task *try*
will store raised error in internal build property *_error*.

All YAML types might be used to define build properties. Thus, you can define
string, integers, floats, but also lists and maps. You may iterate on values
of a property in the build file.

### Predefined build properties

There are pre-defined build properties:

- **_BASE** is the main build file directory.
- **_HERE** is the current directory when build starts.
- **_OS** is the name of the operating system, such as *linux*.
- **_ARCH**  is the hardware architecture, such as *amd64*.
- **_NCPU** is the number of cores in the processor.

Thus, following build file:

```yaml
default: test

targets:

  test:
    steps:
    - print: 'BASE: ={_BASE}'
    - print: 'HERE: ={_HERE}'
    - print: 'OS:   ={_OS}'
    - print: 'ARCH: ={_ARCH}'
    - print: 'NCPU: ={_NCPU}'
```

Will output on my machine:

```
$ n
----------------------------------------------------------------------- test --
BASE: /home/casa/dsk
HERE: /home/casa/dsk
OS:   linux
ARCH: amd64
NCPU: 2
OK
```

### Build properties on command line

You can pass build properties on command line with `-props` option and a YAML
map to set properties. For instance, calling this build file:

```yaml
default: test

properties:
  FOO: 'foo'

targets:

  test:
    steps:
    - print: 'FOO: ={FOO}'
    - print: 'BAR: ={BAR}'
```

With command line defining properties will print:

```
$ neon -props '{FOO: FOO, BAR: bar}'
----------------------------------------------------------------------- test --
FOO: FOO
BAR: bar
OK
```

Thus:

- Property *FOO* was set to *foo* in build file but was overwritten on command
line with value *FOO*.
- Property *BAR* was not defined in build file but was set on command line.

### Configuration

Sometimes, you don't want to write properties in a build file:

- Some vary depending on the developer's environment.
- Some are confidential and should not be made public in a build file.

You should write these properties in a configuration file that will be loaded
by the build file on startup and overwrite properties of the build file.

Let's say you have following build file:

```yaml
default: test

configuration: configuration.yml

properties:
    TOOL_HOME: ~
    PASSWORD:  ~

targets:

    test:
      steps:
      - $: ['={TOOL_HOME}/bin/tool', 'command', 'line', 'options']
      - $: ['service-that-needs-password', =PASSWORD]
```

You could write configuration in *configuration.yml* file as follows:

```yaml
TOOL_HOME: '/opt/misc/mytool'
PASSWORD:  'fazelirflnazrfl'
```

You should probably exclude *configuration.yml* from your version management
system. Thus, using Git, you would add following line in your *.gitignore*
file:

```
/configuration.yml
```

You should also probably document properties that must be defined in a separate
configuration file. This is a good idea to provide a commented template
configuration file in the project.

### Properties hierarchy

You can define properties in the build file, in a configuration file and on
command line. The hierarchy for properties is the following:

- Properties defined in configuration overwrite those defined in build file
  and previous configuration files (in the order of the list of the
  *configuration* field).
- Properties defined on command line overwrite all other properties.

Thus, with following build file:

```yaml
default: test

configuration: 'configuration.yml'

properties:
  FOO: 'foo'

targets:

  test:
    steps:
    - print: 'FOO: ={FOO}'
```

And configuration file *configuration.yml*:

```yaml
FOO: 'conf'
```

Running the build would produce:

```
$ neon
----------------------------------------------------------------------- test --
FOO: conf
OK
```

And running it redefining property on command line:

```
$ neon -props '{FOO: cmd}'
------------------------------------------------------------------------ test --
FOO: cmd
OK
```

[Back to top](#user-manual)

Build targets
-------------

Build targets might be compared to functions. They are called on command line.
If you call NeON with `neon clean`, you will run target *clean*.

This target might look like:

```yaml
properties:
  BUILD_DIR: 'build'

targets:

  clean:
    doc: Clean generated files
    steps:
    - delete: =BUILD_DIR
```

A target might define following fields:

- **doc** this is the target documentation.
- **depends** to list targets to run before running this one.
- **steps** is the list of tasks to run the target.

Tasks might be one of the following:

[Back to top](#user-manual)

### NeON task

They are tasks defined in NeON engine. This is a way to write platform
independent build files. These tasks are maps with string keys.

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

There are also logical tasks to perform tests and iterate on values. For
instance, to iterate on a list of files, you could write:

```yaml
targets:

  pdf:
    doc: Generate PDF files
    steps:
    - for: 'file'
      in:  'find("md", "*.md")'
      do:
      - $: ['md2pdf', '-o', 'build/={file}.pdf', 'md/={file}']
```

To generate a file if the source is newer, you could write:

```yaml
targets:

  pdf:
    doc: Generate PDF file
    steps:
    - if: 'older("file.md", "build/file.pdf")'
      then:
      - $: ['md2pdf', '-o', 'build/file.pdf', 'file.md']
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
$ neon -task time
Record duration to run a block of steps.

Arguments:

- time: steps we want to measure execution duration (steps).
- to: property to store duration in seconds as a float, if not set, duration is
  printed on the console (string, optional).

Examples:

    # print duration to say hello
    - time:
      - print: 'Hello World!'
      to: duration
    - print: 'duration: ={duration}s'
```

You can get information on available tasks
[on this reference page](reference.md).

### File tasks

Many NeON tasks manage files and have common fields. For instance, *copy* task
as follows:

```yaml
- copy:  ['**/*.txt', '**/*.md']
  dir:   'txt'
  todir: 'dst'
  flat:  true
```

The **task field**, *copy* in this case, defines globs for files to select.
These globs are like these used on command line:

- **\*** to select any number of characters. Thus *\*.txt* will select all
files with *txt* extension.
- **?** to select a character. Thus *?.txt* would select *1.txt* but not
  *12.txt*.
- **\*\*** to select any number of directories. Thus *\*\*/\*.txt* would select
  *foo.txt*, *foo/bar.txt* and any file with *txt* extension in a subdirectory.

For more information on these globs, see 
[zglob documentation](http://github.com/mattn/go-zglob).

The **dir** field is the root directory for globs. Note that selected files are
relative to this directory. If *dir* field is not set, it defaults to current
working directory.

Field **exclude** is a list of globs of files to exclude from selection. This
is optional.

Fields **todir** or **tofile** are destination. One and only one of these fields
might be set, and *tofile* can be set only if globs have selected only one file.

Field **flat** tells if copy will be flat, that is at the root of destination
directory. This is a boolean, thus you should not set it with a string such as
*"false"*, but with a true boolean, such as *false*. This defauts to *true*.

[Back to top](#user-manual)

### Shell task

A shell task runs a script. This script will run with *sh* on Unix and
*cmd.exe* on Windows by default. You can define which shell to use with a
*shell* field at the root of the build file (see section
[Build file structure](build-file-structure) for more information).

This script might be a simple command such as *ls*:

```yaml
targets:

  shell:
    steps:
    - $: 'ls'
```

This may also be a full shell or batch script:

```yaml
targets:

  shell:
    steps:
    - $: |
         set -e
         echo "Renaming branch '={old}' to '={new}'"
         git branch -m ={old} ={new}
         git push origin :={old}
         git push --set-upstream origin ={new}
```

Nevertheless, you might prefer not to rely on shell to run a command. This is
the case on Windows where command interpreter doesn't manage well arguments
with spaces. To do so, you should write commands as lists:

```yaml
targets:

  shell:
    steps:
    - $: ['java', '-jar', 'my.jar', 'Hello World!']
```

This way, you will write portable build files that will run seamlessly on Unix
boxes or Windows machines.

A shell task will fail if the command returns a value different from *0*. You
might manage errors with *try/catch/finally* task.

[Back to top](#user-manual)

### Script task

A script task is a piece of code that will run in the NeON scripting engine.
This is also a way to write platform independent code. But this is a way to
write complex scripts that would be complicated to write with system commands.

A script task is a simple string. For instance:

```yaml
targets:

  script:
    steps:
    - 'file = joinpath(BUILD_DIR, "test.txt")'
```

Your scripts might use builtin functions, defined by Anko scripting engine (such
as *toString()*) or by NeON. To list NeON builtins, you can type command 
`neon -builtins`. To get help on given builtin, type `neon -builtin split`.
You can get information on available builtin functions
[on this reference page](reference.md).

You can define your own functions in scripts that you load in the build file
with *context* field. For instance, to define your own *double* function to use
it in you build files, you could define *context* as following:

```yaml
context: 'myscript.ank'
```

The content of this script might be:

```go
func double(i) {
    return 2*i
}
```

Then, in you build file, you could write:

```yaml
targets:

  script:
    steps:
    - 'd = double(21)'
```

A this will call the *double()* function you defined in your context. This is a
way to write your build utility functions in a separate file and thus keep a
clean build file.

To get more information about 
[Anko scripting language clic here](http://github.com/mattn/anko).

[Back to top](#user-manual)

Command line options
--------------------

To get help on command line options, you can type:

```
$ neon -help
Usage of neon:
  -builtin string
    	Print help on given builtin
  -builtins
    	Print builtins list
  -file string
    	Build file to run (default "build.yml")
  -grey
    	Print on terminal without colors
  -info
    	Print build information
  -install string
    	Install given plugin
  -parents
    	List available parent build files in repository
  -props string
    	Build properties
  -refs
    	Print tasks and builtins reference
  -repo string
    	Neon plugin repository for installation (default "~/.neon")
  -targets
    	Print targets list
  -task string
    	Print help on given task
  -tasks
    	Print tasks list
  -template string
    	Run given template
  -templates
    	List available templates in repository
  -time
    	Print build duration
  -version
    	Print neon version
```

In most cases, you will call NeON passing build targets to invoke. Thus to call
target foo, you would type `neon foo`. You can call more than one target on
command line, with `neon foo bar`. Note that second target will be called even
if it already ran calling *foo*.

Called build file will default to *build.yml* in current directory. If this
file is not found in current directory, it will be searched recursively in
parent directories. You can force build file name with the `-file` option. Thus
to run build file *foo.yml*, you would type `neon -file foo.yml`. Execution
times are always written on console when greater than *10 s*. You can force to
print build execution time with `-time` option. 

You can get information on build file with `-info` option. This will print the
build documentation (written in *doc* field at the root of the build file),
default target(s), repository, extended build files, properties (with their own
help) and targets (with their help). Using this option is a good way to have an
idea of what can perform a build file. You can get targets list with `-targets`
option.

You can define properties on command line with `-props` options and a YAML map
with properties. For instance, to define property *foo* with value *bar*, you
would invoke NeON with command line `neon -props '{foo: bar}'`.

You can set the path to your repository (where live parent build files and
templates) with `-repo` option. This defaults to *~/.neon* but you can set it
anywhere with this option. This option affects builds, but also where are
installed plugin with `-install` option.

The `-install` option will install given plugin in repository. Thus typing
`neon -install foo/bar` will try to clone propject *bar* of user *foo* on
Github into your repository.

You can list parent build build files in your repository with `-parents` option
and templates with `-templates`. These options are affected by `-repo` option.
You can run a template with `-template` option. Thus to run template
*foo/bar/spam.tpl*, you would type `neon -template foo/bar/spam.tpl`.

To list all available builtins, you have option `-builtins`. To get help on a
given builtin, you would type `neon -builtin foo`. To list all available tasks,
you have option `-tasks` and to get help on a given task, you would type
`neon -task foo`. Option `-refs` will output on console help for all builtins
and task in Markdown format (this is the way reference documentation is
generated).

By default, build output is colored on Unix systems for dark terminals (that is
white letters on black background). You can disable colorization with `-grey`
option.

Option `-version` will print NeON version.

Build inheritance
-----------------

A build file can extend another parent build file with the *extends* field.
For instance, with this parent build file called *buildir.yml*:

```yaml
doc: Parent build file to manage build directory

properties:
  BUILD_DIR: 'build'

targets:

  clean:
    doc: Clean generated files
    steps:
    - delete: '=BUILD_DIR'
```

You may reuse this build file in another one:

```yaml
doc: Build file
extends: ./buildir.yml

targets:

  compile:
    doc: Compile
    depends: clean
    steps:
    - $: 'compile sources'
```

The main build file that extends *buildir.yml* will inherit its properties and
targets.

Note that a build file can redefine inherited properties. For instance, you may
decide to set *BUILD_DIR* to *target* with following build file:

```yaml
doc: Build file
extends: ./buildir.yml
properties:
  BUILD_DIR: target

targets:

  compile:
    doc: Compile
    depends: clean
    steps:
    - $: 'compile sources'
```

This will redefine the build directory as expected.

You may also redefine targets. For instance, let's say you want to warn before
deleting build directory. You would write:

```yaml
doc: Build file
extends: ./buildir.yml

targets:

  clean:
    doc: Clea generated files
    steps:
    - print: 'Deleting build directory!!!'
    - super:
```

The path to extended build files is important:

- If this path is **absolute** or starts with **./** (that is in current
  directory), this works as you would expect.
- If this path is relative without starting with *./*, this build file is in
  the NeON repository. See bellow for more explanations.

The *super* task will run steps of parent target.

You can list all parent build files in your repository with following command:

```
$ neon -parents
c4s4/build/buildir.yml
c4s4/build/github.yml
c4s4/build/golang.yml
c4s4/build/java.yml
c4s4/build/release.yml
c4s4/build/slides.yml
c4s4/build/xslt.yml
```

[Back to top](#user-manual)

Neon repository
---------------

NeON repository is the place where live installed parent build files. By
default, NeON repository is in *~/.neon* directory. This may be changed with
*repository* field in the build file.

When you extend a build file with following statement:

```yaml
extends: foo/bar/spam.yml
```

NeON will look in directory *foo/bar* of the NeON repository for file
*spam.yml*. You might install manually your parent build files in your NeON
repository, but you can install them automagically with NeON *install* command:

```bash
$ neon -install foo/bar
```

This will clone Git repository *bar* of user *foo* on *Github*. Thus this will
run command:

```bash
$ git clone git://github.com/foo/bar.git <neon-repo>/foo/bar
```

Thus, if your parent build files are public, you should put them on *Github* so
that they can be easily shared in your team. I personally share my parent build
files in repository <http://github.com/c4s4/build>.

By default this will clone *master* branch. You can change this running
following command in created Git repository:

```bash
$ git checkout develop
```

This will change branch to *develop*. You might get a particular version with:

```bash
$ git checkout 1.2.3
```

In your parent project repository, simply put you parent build files at the root.
You might also put them in any subdirectory. If you put a build file *spam.yml*
in subdirectory *eggs*, you would extend it with:

```yaml
extends:
- foo/bar/eggs/spam.yml
```

Project templates
-----------------

NeON can generate template projects, with the *-template* option. For instance,
to generate template Golang project, you would:

- Install *c4s4/build* plugin, typing `neon -install c4s4/build`.
- Run Golang template with command `neon -template c4s4/build/golang.tpl`.

This will ask you the project name and generate the project:

```
$ neon -template c4s4/build/golang.tpl
------------------------------------------------------------------- template --
Name of this project: test
Making directory '/home/casa/dsk/test'
Copying 6 file(s)
Moving 1 file(s)
Moving 1 file(s)
Replacing text in file '/home/casa/dsk/test/build.yml'
Project generated in 'test' directory
OK
```

This creates a *test* directory with generated project:

```
$ ls test
build.yml  CHANGELOG.yml  LICENSE.txt  README.md  test.go  test_test.go
```

To have an idea of targets in newly created project, go in generated directory
and type:

```
$ neon -info
repository: ~/.neon

extends:
- c4s4/build/golang.yml

properties:
  BUILD_DIR: "build" 
  LIBRARIES: ["github.com/mitchellh/gox"] 
  VERSION:   "1.0.0" 
  ARCHIVE:   "build/test-1.0.0.tar.gz" 
  NAME:      "test" 

targets:
  archive: Build distribution archive 
  bin:     Make binary 
  clean:   Clean build directory 
  fmt:     Format Go code 
  libs:    Install libraries 
  run:     Run Go tool 
  test:    Run Go tests
```

You can list all available templates in you repository typing:

```
$ neon -templates
c4s4/build/golang.tpl
c4s4/build/java.tpl
c4s4/build/slides.tpl
```

### Creating templates

To create your own templates, you can have a look at following Github project;
<http://github.com/c4s4/build> which contains example *Golang* template
project.

The template is made of a build file, *golang.tpl*:

```yaml
# Neon template file (http://github.com/c4s4/neon)

default: template

targets:

  template:
    doc: Generate Golang project
    steps:
    - prompt:  'Name of this project'
      to:      'name'
      pattern: '^\w+$'
      error:   'Project name must be made of letters, numbers, - and _'
    - if: 'exists(joinpath(_HERE, name))'
      then:
      - throw: 'Project directory already exists'
    - mkdir: '={_HERE}/={name}'
    - copy:  '*'
      dir:   '={_BASE}/golang'
      todir: '={_HERE}/={name}'
    - move:   '={_HERE}/={name}/main.go'
      tofile: '={_HERE}/={name}/={name}.go'
    - move:   '={_HERE}/={name}/main_test.go'
      tofile: '={_HERE}/={name}/={name}_test.go'
    - replace: '={_HERE}/={name}/build.yml'
      with:    {'main': =name}
    - print: "Project generated in '={name}' directory"
```

This is a NeON build file that generates project in current working directory.
Note that this is done with *_HERE* property that is current directory, while
*_BASE* is the directory of the build file, that lives in the NeON repository.

This build file prompts the user for the project name and then copies project
files, in the *golang* directory, to the project directory. Then it renames
files and performs some replacements.

Template build files are named with *tpl* extension so that they are identified
as templates, but otherwise they are plain old build files.

*Enjoy!*
