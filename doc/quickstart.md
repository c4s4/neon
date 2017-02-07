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
syntax to write them:</p>

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

