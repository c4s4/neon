# Neon Quick Start

Neon is a build tool, which means that its purpose is to automate a build process. Its design goals are:

- **Clean and simple build file syntax**: While XML is too verbose and a programming language too complicated, [YAML](http://www.yaml.org/) is perfect with its simple yet powerful syntax.
- **Speed**: Slow startup are irritating while working with a build tool whole day long. Neon happens to be as fast as Make, according to my tests on Make builds ported to Neon.
- **System and language independent**: Neon is written in Go, which run on most operating systems and hardwares, and is not limited to building Go projects. Furthermore, using Neon tasks you can write platform independent builds.
- **Scriptable**: when a task is complicated, it is very handy to script it with a real programming language. Neon build files can embed [Anko](http://github.com/mattn/anko) scripts, which have a syntax very close to Go.

To demonstrate these features, let's see what *Hello World!* looks like using Neon:

```yaml
targets:

  hello:
    steps:
    - print: "Hello World!"
```

This is that simple!

## Build properties

A build file can embed properties that are like variables. For instance, we can put user name in a property named *USER_NAME* like this:

```yaml
properties:
  USER_NAME: "World"

targets:

  hello:
    steps:
    - print: "Hello ={USER_NAME}!"
```

We access the value of this property with `={PROP_NAME}` syntax. Thus, to get value for property `USER_NAME`, we must write `={USER_NAME}`.

These properties live as variables defined in the context of the embedded [Anko VM](http://github.com/mattn/anko). Thus the content of `={...}` is evaluated as an Anko expression, and you can write *Anko* expressions, such as `={uppercase(USER_NAME)}`.

Properties may also have integer, list or hash values. To do so, you should use YAML syntax to write them:

```yaml
properties:
  STRING:  "string"
  INTEGER: 1
  LIST:    [1, "two"]
  HASH:    {one: 1, two: 2}

targets:

  hello:
    steps:
    - print: "string:  ={STRING}"
    - print: "integer: ={INTEGER}"
    - print: "list:    ={LIST}"
    - print: "hash:    ={HASH}"
```

Which will output:

```
$ neon hello
---------------------------------------------------------------------- hello --
string:  string
integer: 1
list:    [1, two]
hash:    [one: 1, two: 2]
OK
```

You can pass properties on command line with `-props` option. For instance, to define property *foo* with value *"bar"*, you would add on command line: `-props='{foo: "bar"}'`. These properties will overwrite those defined in the build file.

## Targets

If we see build files as programs, we could see targets as functions that you can call to achieve a given goal. A build file can define more than one target that may depend on each other. For instance:

```yaml
properties:
  USER_NAME: "World"

targets:

  hello:
    depends: upper
    steps:
    - print: "Hello ={USER_NAME}!"

  upper:
    steps:
    - 'USER_NAME = uppercase(USER_NAME)'
```

Which produces following output:

```
$ neon hello
---------------------------------------------------------------------- upper --
---------------------------------------------------------------------- hello --
Hello WORLD!
OK
```

Target *hello* now depends on target *upper* which puts `USER_NAME` in upper case. Thus, Neon runs first target *upper*, then target *hello*.

When you don't pass any target on command line, Neon selects default one that you can set with `default` field in the build file:

```yaml
default: hello

targets:

  hello:
    steps:
    - print: "Hello World!"
```

You can run more than one default target by settings `default` to a list of targets to run:

```yaml
default: [foo, bar]
```

## Tasks

If targets are functions, tasks are instructions. They can be *shell* scripts, *Anko* code or *Neon* tasks:

#### Shell scripts

To run a shell script, you use task named *$*. Thus, to print the user's name, we could write:

```yaml
- $: 'echo "Hello ${USER}!"'
```

Note that we can surround YAML strings with simple or double quotes. We choose simple ones here so that we can use double for the Shell string. We could also escape double quotes inside YAML string as follows:

```yaml
- $: "echo \"Hello $USER!\""
```

If return value of the script is not *0*, which denotes an error during its execution, the build is interrupted and an error message is printed on the console. For instance, this build file:

```yaml
targets:

  broken:
    steps:
    - $: 'command-that-doesnt-exist'
```

Will produce this output on the console:

```
$ neon broken
--------------------------------------------------------------------- broken --
sh: 1: command-that-doesnt-exist: not found
ERROR running target 'broken': in step 1: executing command: exit status 127
```

A multi line shell scripts can be written using pipe character `|`:

```yaml
targets:

  shell:
    steps:
    - $: |
         echo "This is a long"
         echo "Shell script on"
         echo "More than one line"
```

You can also write commands as a list:

```yaml
targets:

  shell:
    steps:
    - $: ['java', '-jar', 'echo.jar', 'Hello World!']
```

This is useful to write build files that will run on Unix and Windows systems because Windows command interpreter is quite broken regarding options parsing with spaces. In lists, command line options are clearly separated and this
avoids issues.

Note that options in commands as lists won't be evaluated in a shell, thus an option like `${USER}` won't be evaluated. NeON will always evaluate NeON properties, thus an option such as `={USER}` would be evaluated.

### Neon tasks

You could perform most of common build tasks with shell commands, but this would bind your build to a given platform. To develop platform independent builds, you should use Neon tasks. For instance, to copy all XML files to *build* directory, you could call `cp` system command, but you should instead use `copy` Neon task as follows:

```yaml
targets:

  copy:
    steps:
    - copy:  "**/*.xml"
      todir: "build"
```

This will run on all platforms (Unices and Windows) provided you use slashes as path separator instead of platform dependent ones. Furthermore, in file related tasks, you can use extended globs where `**` replaces any number of directories, see [zglob documentation](http://github.com/mattn/zglob) for more information.

To list all available Neon tasks, type `neon -tasks` and to get help on a given one, type `neon -task copy` for instance.

There are special Neon tasks that make it possible to control the execution flow of your build:

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
      - print: "Validating ={file}..."
      - $: ['xmllint', '--noout', '--valid', 'data/={file}']
```

There are also tasks to manage errors:

- **throw**: to interrupt the build with an error.
- **try/catch**: to prevent build failure even on step execution error.

For instance, if you don't want to interrupt validation on error, you could write:

```yaml
targets:

  validate:
    steps:
    - for: file
      in:  find("data", "*.xml")
      do:
      - print: "Validating ={file}..."
      - try:
        - $: ['xmllint', '--noout', '--valid', 'data/={file}']
        catch:
        - print: 'ERROR!'
```

### Scripts

A script is just a simple string, as follows:

```yaml
targets:

  script:
    steps:
    - 'println("Hello World!")'
```

You can also write a long script with the pipe notation:

```yaml
targets:

  script:
    steps:
    - |
      for i in range(10) {
        println(i)
      }
```

In addition to Anko builtin functions, Neon adds handy builtins to perform common build tasks. For instance, function `exists(file)` tells if given file exists.

To list all Neon buitins, type `neon -builtins` and to get help on a given function, type `neon -builtin exists` for instance.

## Parent build files

Your build file can extend existing ones. Parent build files may live in same directory or in a Neon repository, located in *~/.neon/*.

For instance, you can install *c4s4/build* parent build files with command `neon -install c4s4/build`. This will clone *Github* repository *git://github.com/c4s4/build.git* in your Neon repository *~/.neon/c4s4/build*.

You can list parent build files located in your repository with command:

```bash
$ neon -parents
c4s4/build/archive.yml
c4s4/build/buildir.yml
c4s4/build/changelog.yml
c4s4/build/django.yml
c4s4/build/dotenv.yml
c4s4/build/git.yml
c4s4/build/golang.yml
c4s4/build/java.yml
c4s4/build/python.yml
c4s4/build/slides.yml
c4s4/build/xslt.yml
```

You can then extend *golang* parent build file adding in your build file:

```yaml
extends: c4s4/build/golang.yml
```

This provides properties and targets you may list with `neon -info`.

Parent build file projects may also contain templates. You can list templates in your Neon repository with:

```bash
$ neon -templates
c4s4/build/build.tpl
c4s4/build/django.tpl
c4s4/build/flask.tpl
c4s4/build/golang.tpl
c4s4/build/golangws.tpl
c4s4/build/java.tpl
c4s4/build/python.tpl
c4s4/build/slides.tpl
c4s4/build/xslt.tpl
```

You can then generate a *Golang* template project with `neon -template c4s4/build/golang.yml`.

You can play with example parent build file project at *http://github.com/c4s4/build* and customize it to match your need. Best way to create your parent build files is to take existing Neon build files and put them in a *Git* repository, and enjoy!

## Getting help

You can print information on build file typing `neon -info`. You will get help on build properties, environment variables and targets. This will be useful if you have documented your build file. To document a target, you could write:

```yaml
targets:

  validate:
    doc: Validate all XML files in data directory
    steps:
    - for: file
      in:  find("data", "*.xml")
      do:
      - print: 'Validating ={file}...'
      - $: ['xmllint', '--noout', '--valid', 'data/={file}']
```

You can also document the whole build putting a `doc` entry at the root of the build file.

To get help about neon itself, you can:

- List all tasks with `neon -tasks`.
- Get help on a given task with `neon -task copy`.
- List all builtins with `neon -builtins`.
- Get help on a given builtin with `neon -builtin find`.
- List all parent build files with `neon -parents`.
- List all templates with `neon -templates`.

## Go further

Neon as much more to offer, see [User Manual](usermanual.md) for in-depth documentation and [Reference](reference.md) for information about all tasks and builtins functions.

*Enjoy!*
