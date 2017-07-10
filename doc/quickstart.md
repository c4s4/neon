Neon Quick Start
================

Neon is a build tool. This means that its purpose is to automate a build
process. Its design goals are:

- **Clean and simple build file syntax**: While XML is too verbose and
  a programming language too complicated, [YAML](http://www.yaml.org/) is
  perfect with its simple yet powerful syntax.
- **Speed**: Slow startup are irritating while working with a build tool the
  whole day. Neon happens to be as fast as Make, according to my tests on Make
  builds ported to Neon.
- **System and language independent**: Neon is written in Go, which run on most
  platforms, and is not limited to building Go projects. Furthermore, using
  Neon tasks you can write platform independent builds.
- **Scriptable**: when a task is complicated, it is very handy to script it
  with a real programming language. Neon build files can embed Anko scripts.
  [Anko](http://github.com/mattn/anko) is an interpreted Go.

To demonstrate these features, let's see what *Hello World!* looks like using
Neon:

```yaml
targets:

  hello:
    steps:
    - print: "Hello World!"
```

This is that simple!

Build properties
----------------

A build file can embed properties that are like variables. For instance, we can
put user name in a property named *USER_NAME* like this:

```yaml
properties:
  USER_NAME: "World"

targets:

  hello:
    steps:
    - print: "Hello #{USER_NAME}!"
```

We access the value of this property with `#{PROP_NAME}` syntax. Thus, to get
value for property `USER_NAME`, we must write `#{USER_NAME}`.

These properties live as variables defined in the context of the embedded Anko
VM. Thus the content of `#{...}` is evaluated as an Anko expression, and you
can write Anko expressions, such as `#{uppercase(USER_NAME)}`.

Properties may also have integer, list or hash values. You must then use a YAML
syntax to write them:

```yaml
properties:
  STRING:  "string"
  INTEGER: 1
  LIST:    [1, "two"]
  HASH:    {one: 1, two: 2}

targets:

  hello:
    steps:
    - print: "string:  #{STRING}"
    - print: "integer: #{INTEGER}"
    - print: "list:    #{LIST}"
    - print: "hash:    #{HASH}"
```

Which will output:

```
$ n hello
Running target hello
string:  string
integer: 1
list:    [1, two]
ahash:   [one: 1, two: 2]
OK
```

You can pass properties on command line with `-props` option. For instance, to
define property *foo* with value *"bar"*, you would add on command line:
`-props='foo: "bar"'`. These properties with overwrite those defined in the
build file.

Targets
-------

If we see build files as programs, we could see targets as functions that you
can call to achieve a given goal. A build file can define more that one target
that may depend on each other. For instance:

```yaml
properties:
  USER_NAME: "World"

targets:

  hello:
    depends: upper
    steps:
    - print: "Hello #{USER_NAME}!"

  upper:
    steps:
    - script: 'USER_NAME = uppercase(USER_NAME)'
```

Which produces following output:

```
$ neon hello
Running target upper
Running target hello
Hello WORLD!
OK
```

Target *hello* now depends on target *upper* which puts `USER_NAME` in upper
case. Thus, Neon runs first target *upper*, then target *hello*.

When you don't pass any target on command line, Neon selects default one that
you can set with `default` field in the build file:

```yaml
default: hello

targets:

  hello:
    steps:
    - print: "Hello World!"
```

You can run more than one default target by settings `default` to a list of
targets to run:

```yaml
default: [foo, bar]
```

Tasks
-----

If targets are functions, tasks are instructions. They can be shell scripts,
Anko code or Neon tasks:

#### Shell scripts

To run a shell script, you just have to put it in a string. Thus, to print the
user's name, we could write:

```yaml
- 'echo "Hello $USER!"'
```

Note that we can surround YAML strings with simple or double quotes. We choose
simple ones here so that we can use double for the Shell string. We could also
escape double quotes inside YAML string as follows:

```yaml
- "echo \"Hello $USER!\""
```

If return value of the script is not *0*, which denotes an error running the
script, the build is interrupted and an error message is printed on the
console. For instance, this script:

```yaml
targets:

  broken:
    steps:
    - 'command-that-doesnt-exist'
```

Will produce this output on the console:

```
$ n broken
Running target broken
sh: 1: command-that-doesnt-exist: not found
ERROR running target 'broken': in step 1: exit status 127
```

A multi line shell scripts can be written using pipe character `|`:

```yaml
targets:

  shell:
    steps:
    - |
      echo "This is a long"
      echo "Shell script on"
      echo "More than one line"
```

### Neon tasks

You could perform most of common build tasks with shell ones, but this would
bind your build to a given platform. To develop platform independant builds,
you should use Neon tasks. For instance, to copy all XML files to *build*
directory, you could call `cp` system command, but you should instead use
copy Neon task as follows:

```yaml
targets:

  copy:
    steps:
    - copy: "**/*.xml"
      todir: "build"
```

This will run on all platforms (Unices and Windows) provided you use slashes
as path separator instead of platform dependent ones. Furthermore, in most
file related tasks, you can use extended globs where `**` replaces any
number of directories, see [zglob documentation](http://github.com/mattn/zglob)
for more information.

To list all available Neon tasks, type `neon -tasks` and to get help on a given
one, type `neon -task copy` for instance.

There are special Neon tasks that make it possible to control the execution
flow of your build:

- **for/in/do**: to make a loop on a list and store value in a variable.
- **if/then/else**: to control execution flow depending on a test.
- **while/do**: too loop while a given condition is met.

For instance, to validate all XML files in *data* directory, you could write:

```yaml
targets:

  validate:
    steps:
    - for: file
      in:  find("data", "*.xml")
      do:
      - print: "Validating #{file}..."
      - 'xmllint --noout --valid data/#{file}'
```

There are also tasks to manage errors:

- **throw**: to interrupt the build with an error.
- **try/catch**: to prevent build failure even on step execution error.

For instance, if you don't want to interrupt validation on error, you could
write:

```yaml
targets:

  validate:
    steps:
    - for: file
      in:  find("data", "*.xml")
      do:
      - print: "Validating #{file}..."
      - try:
        - 'xmllint --noout --valid data/#{file}'
        catch:
        - print: "ERROR!"
```

### Anko scripts

An Anko script is indicated with `script` instruction as follows:

```yaml
targets:

  script:
    steps:
    - script: 'println("Hello World!")'
```

You can also write a long Anko script with the pipe notation:

```yaml
targets:

  script:
    steps:
    - script: |
        for i in range(10) {
          println(i)
        }
```

In addition to Anko builtin functions, Neon adds handy builtins to perform
common build tasks. For instance, function `exists(file)` tells if given file
exists.

To list all Neon buitins, type `neon -builtins` and to get help on a given
function, type `neon -builtin function`.

Example
-------

Here is an example of what you can do with neon. This task validates generated
files against a schematron and produces an HTML report:

```yaml
  schematron:
    doc: Validate geenrated files with schematron
    steps:
    - for: 'encyclo'
      in:  'find(DEST_DIR, "*")'
      do:
      - print: "Validating #{encyclo}:"
      - script: 'errors = {}'
      - for: 'file'
        in:  'find(joinpath(DEST_DIR, encyclo), "**/*.xml")'
        do:
        - try:
          - print: "- #{file}"
          - execute: 'jing -C "#{CATALOG}" "src/sch/Ouvrages_v1.sch" "#{DEST_DIR}/#{encyclo}/#{file}"'
            output:  'output'
          catch:
          - script: 'errors[file] = output'
      - if: 'len(keys(errors)) > 0'
        then:
        - script: |
            sort = import("sort")
            report = "# Validation Schematron des Encyclopedies\n\n"
            files = sort.Strings(keys(errors))
            for file in files {
              report += "### " + file + "\n\n```\n"
              report += errors[file]
              report += "\n```\n\n"
            }
        - write: '#{BUILD_DIR}/#{encyclo}.md'
          from:  'report'
        - 'pandoc -f markdown -t html "#{BUILD_DIR}/#{encyclo}.md" > "#{BUILD_DIR}/#{encyclo}.html"'

```

This task:

- Iterates books in directories.
- It validates each file with *jing*.
- When validation fails, it fills a map with errors for a given file.
- When a book is done, it generates a *markdown* report with errors.
- Then it converts this report to HTML with *pandoc*.

This illustrates how to combine *for* tasks to crawl directories, shell commands
to validate files and *anko* script to generate a report from tool output.

Getting help
------------

You can print help on build file typing `neon -build`. You will get help on
build properties, environment variables and targets. This will be useful if
you have documented your build file. To document a given property, you would
write:

```yaml
targets:

  validate:
    doc: Validate all XML files in data directory
    steps:
    - for: file
      in:  find("data", "*.xml")
      do:
      - print: "Validating #{file}..."
      - 'xmllint --noout --valid data/#{file}'
```

You can also document the whole build putting a `doc` entry at the root of the
build file.

To get help about neon itself, you can:

- List all tasks with `neon -tasks`.
- Get help on a given task with `neon -task copy`.
- List all builtins with `neon -builtins`.
- Get help on a given builtin with `neon -builtin find`.

Go further
----------

Neon as much more to offer, see [User Manual](usermanual.md) for in depth
documentation and [Reference](reference.md) for information about all tasks and
builtins functions.

*Enjoy!*
