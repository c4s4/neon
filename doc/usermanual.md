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
