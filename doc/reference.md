Tasks Reference
===============

cat
---

Print the content of e given file on the console.

Arguments:

- cat: the file to print on console as a string.

Examples:

    # print content of LICENSE file on the console
    - cat: "LICENSE"

chdir
-----

Change current working directory.

Arguments:

- chdir: the directory to change to (as a string).

Examples:

    # go to build directory
    - chdir: "build"

Notes:

- Working directory is set to the build file directory before each target.

chmod
-----

Changes mode of files.

Arguments:

- chmod: the list of globs of files to change mode (as a string or list of
  strings).
- mode: the mode in octal form (such as '0755') as a string
- dir: the root directory for glob (as a string, optional, defaults to '.').
- exclude: globs of files to exclude (as a string or list of strings,
  optional).

Examples:

    # make foo.sh executable for all users
    - chmod: "foo.sh"
      mod: "0755"
    # make all sh files in foo directory executable, except for bar.sh
    - chmod: "**/*.sh"
      mode: "0755"
      exclude: "**/bar.sh"

copy
----

Copy file(s).

Arguments:

- copy: the list of globs of files to copy (as a string or list of strings).
- dir: the root directory for glob (as a string, optional).
- exclude: globs of files to exclude (as a string or list of strings,
  optional).
- tofile: the file to copy to (as a string, optional, only if glob selects a
  single file).
- todir: directory to copy file(s) to (as a string, optional).
- flat: tells if files should be flatten in destination directory (as a boolean,
  optional, defaults to true).

Examples:

    # copy file foo to bar
    - copy:   "foo"
      tofile: "bar"
    # copy text files in directory 'book' (except 'foo.txt') to directory 'text'
    - copy: "**/*.txt"
      dir: "book"
      exclude: "**/foo.txt"
      todir: "text"
    # copy all go sources to directory 'src', preserving directory structure
    - copy: "**/*.go"
      todir: "src"
      flat: false

delete
------

Delete a directory recursively.

Arguments:

- delete: directory or list of directories to delete.

Examples:

    # delete build directory
    - delete: "#{BUILD_DIR}"

Notes:

- Handle with care, this is recursive!

execute
-------

Execute a command and return output and value.

Arguments:

- execute: command to run.
- output: name of the variable to store trimed output into.

Examples:

    # execute ls command and get result in 'files' variable
    - execute: 'ls'
      output:  'files'

for
---

For loop.

Arguments:

- for: the name of the variable to set at each loop iteration.
- in: the list of values or expression that generates this list.
- do: the block of steps to execute at each loop iteration.

Examples:

    # create empty files
    - for: file
      in:  ["foo", "bar"]
      do:
    - touch: "#{file}"
    # print first 10 integers
    - for: i
      in: range(10)
      do:
      - print: "#{i}"

if
--

If condition.

Arguments:

- if: the condition.
- then: the steps to execute if the condition is true.
- else: the steps to execute if the condition is false.

Examples:

    # print hello if x > 10 else print world
    - if: x > 10
      then:
      - print: "hello"
      else:
      - print: "world"

link
----

Create a symbolic link.

Arguments:

- link: the source file.
- to: the destination of the link.

Examples:

    # create a link from file foo to bar
    - link: "foo"
      to: "bar"

mkdir
-----

Make a directory.

Arguments:

- mkdir: directory or list of directories to create.

Examples:

    # create a directory 'build'
    - mkdir: "build"

move
----

Move file(s).

Arguments:

- move: the list of globs of files to move (as a string or list of strings).
- dir: the root directory for glob (as a string, optional).
- exclude: globs of files to exclude (as a string or list of strings,
  optional).
- tofile: the file to move to (as a string, optional, only if glob selects a
  single file).
- todir: directory to move file(s) to (as a string, optional).
- flat: tells if files should be flatten in destination directory (as a boolean,
  optional, defaults to true).

Examples:

    # move file foo to bar
    - move:   "foo"
      tofile: "bar"
    # move text files in directory 'book' (except 'foo.txt') to directory 'text'
    - move: "**/*.txt"
      dir: "book"
      exclude: "**/foo.txt"
      todir: "text"
    # move all go sources to directory 'src', preserving directory structure
    - move: "**/*.go"
      todir: "src"
      flat: false

pass
----

Does nothing.

Arguments:

- none

Examples:

    # do nothing
    - pass:

Notes:

- This implementation is super optimized for speed.

path
----

Build a path from files and put it in a variable.

Arguments:

- path: the list of globs of files to build the path (as a string or list of
  strings).
- to: the variable to put path into.
- dir: the root directory for glob (as a string, optional).
- exclude: globs of files to exclude (as a string or list of strings,
  optional).

Examples:

    # build classpath with jar files in lib directory
    - path: "lib/*.jar"
      to: "classpath"

print
-----

Print a message on the console.

Arguments:

- print: the text to print as a string.

Examples:

    # say hello
    - print: "Hello World!"

read
----

Read given file as text and put its content in a variable.

Arguments:

- read: the file to read as a string.
- to: the name of the variable to set with the content.

Examples:

    # put content of LICENSE file on license variable
    - read: "LICENSE"
      to: license

remove
------

Remove file(s).

Arguments:

- remove: file or list of files to remove.

Examples:

    # remove all pyc files
    - remove: "**/*.pyc"

replace
-------

Replace pattern in text files.

Arguments:

- replace: the list of globs of files to work with (as a string or list of strings).
- pattern: the text to replace.
- with: the replacement text.
- dir: the root directory for glob (as a string, optional).
- exclude: globs of files to exclude (as a string or list of strings,
  optional).

Examples:

    # replace foo with bar in file test.txt
    - replace: "test.txt"
      pattern: "foo"
      with: "bar"

script
------

Run an Anko script.

Arguments:

- script: the source of the script to run.

Examples:

    # build a classpath with all jar files in lib directory
    - script: |
        strings = import("strings")
        jars = find("lib", "*.jar")
        classpath = strings.Join(jars, ":")

Notes:

- The scripting language is Anko, which is a scriptable Go. For more information
  please refer to Anko site at http://github.com/mattn/anko. Thanks Mattn!
- Buitlin functions are functions you can access in scripts. To list them, you
  cas type 'neon -builtins', to get help on a given one, you may type for instance
  'neon -builtin find'.
- Properties can be accessed and set in scripts. Variables you define in scripts
  are readable as properties. In other words, scripts and properties share the
  same context.

sleep
-----

Sleep a given number of seconds.
		
Arguments:

- sleep: the duration to sleep in seconds as an integer.

Examples:

    # sleep for 10 seconds
    - sleep: 10

tar
---

Create a tar archive.

Arguments:

- tar: the list of globs of files to tar (as a string or list of strings).
- dir: the root directory for glob (as a string, optional).
- exclude: globs of files to exclude (as a string or list of strings,
  optional).
- tofile: the name of the tar file to create as a string.
- prefix: prefix directory in the archive.

Examples:

    # tar files in build directory in file named build.tar.gz
    - tar: "build/**/*"
      tofile: "build.tar.gz"

Notes:

- If archive filename ends with gz (with a name such as foo.tar.gz or foo.tgz)
  the tar archive is compressed with gzip.

throw
-----

Throws an error.

Arguments:

- throw: the message of the error.

Examples:

    # stop the build because tests don't run
    - throw: "ERROR: tests don't run"

Notes:

- The error message will be printed on the console as the source of the build
  failure.

time
----

Print duration to run a block of steps.

Arguments:

- time: the steps to measure execution duration.

Examples:

    # measure duration to say hello
    - time:
      - print: "Hello World!"

touch
-----

Touch a file (create it or change its time).

Arguments:

- touch: the file or files to create.

Examples:

    # create file in build directory
    - touch: "#{BUILD_DIR}/foo"

Notes:

- If the file already exists it changes it modification time.
- If the file doesn't exist, it creates an empty file.

try
---

Try/catch/finally construct.

Arguments:

- try: steps to execute.
- catch: executed if an error occurs (optional).
- finally: executed in all cases (optional).

Examples:

    # execute a command and continue even if it fails
    - try:
      - "command-that-doesnt-exist"
	- print: "Continue even if command fails"
	# execute a command and print a message if it fails
	- try:
	  - "command-that-doesnt-exist"
	  catch:
	  - print: "There was an error!"
	# execute a command a print message in all cases
	- try:
	  - "command-that-doesnt-exist"
	  finally:
	  - print: "Print whatever happens"

Notes:

- The error message for the failure is stored in '_error' variable as text.

while
-----

While loop.

Arguments:

- while: the condition that is evaluated at each loop.
- do: steps that run while condition is true.

Examples:

    # loop until i >= 10
    - while: 'i < 10'
      do:
      - script: 'println(i++)'

write
-----

Write text into a given file.

Arguments:

- write: the file to write into as a string.
- text: the text to write into the file.
- append: tells if we should append content to file (defaults to false).

Examples:

    # write 'Hello World!' in file greetings.txt
    - write: "greetings.txt"
      text: "Hello World!"

zip
---

Create a Zip archive.

Arguments:

- zip: the list of globs of files to zip (as a string or list of strings).
- dir: the root directory for glob (as a string, optional).
- exclude: globs of files to exclude (as a string or list of strings,
  optional).
- tofile: the name of the Zip file to create as a string.
- prefix: prefix directory in the archive.

Examples:

    # zip files in build directory in file named build.zip
    - zip: "build/**/*"
      tofile: "build.zip"


Builtins Reference
==================

directory
---------

Return directory of a given path.

Arguments:

- The path to get directory for as a string.

Returns:

- The directory of the path as a string.

Examples:

    // get directory of path "/foo/bar/spam.txt"
    dir = directory("/foo/bar/spam.txt")

exists
------

Tells if a given path exists.

Arguments:

- The path to test as a string.

Returns:

- A boolean telling if path exists.

Examples:

    // test if given path exists
    if exists("/foo/bar") { ...

expand
------

Exapand file name by replace ~/ with home directory.

Arguments:

- The path to expand as a string.

Returns:

- The expanded path as a string.

Examples:

    // expand path ~/.profile
    profile = expand("~/.profile")

filename
--------

Return filename of a given path.

Arguments:

- The path to get filename for as a string.

Returns:

- The filename of the path as a string.

Examples:

    // get filename of path "/foo/bar/spam.txt"
    file = filename("/foo/bar/spam.txt")

filter
------

Filter a list of files with excludes.

Arguments:

- includes: the list of files to filter.
- excludes: a list of patterns for files to exclude.

Returns:

- The list if filtered files as a list of strings.

Examples:

    // filter text files removing those in build directory
    filter(find("**.txt"), "build/**/*")

Notes:

- Works great with find() builtin.

find
----

Find files.

Arguments:

- The directory of files to find.
- The list of pattern for files to find.

Returns:

- Files as a list of strings.

Examples:

    // find all text files in book directory
    find("book", "**/*.txt")
    // find all xml and yml files in src directory
    find("src", "**/*.xml", "**/*.yml")

Notes:

- Files may be filtered with filter() builtin.

join
----

Join strings.

Arguments:

- The strings to join as a list of strings.
- The separator as a string.

Returns:

- Joined strings as a string.

Examples:

    // join "foo" and "bar" with a space
    join(["foo", "bar"], " ")

joinpath
--------

Join file paths.

Arguments:

- The paths to join as a list of strings.

Returns:

- Joined path as a string.

Examples:

    // join paths "/foo", "bar" and "spam.txt"
    joinpath("foo", "bar", "spam.txt")

lowercase
---------

Put a string in lower case.

Arguments:

- The string to put in lower case.

Returns:

- The string in lower case.

Examples:

    // set greetings in lower case
    greetings = lowercase(greetings)

now
---

Return current date and time in ISO format.

Arguments:

- none

Returns:

- ISO date and time as a string.

Examples:

    // put current date and time in dt variable
    dt = now()

read
----

Read given file and return its content as a string.

Arguments:

- The file name to read.

Returns:

- The file content as a string.

Examples:

    // read VERSION file and set variable version with ots content
    version = read("VERSION")

run
---

Run given command and return output.

Arguments:

- The command to run.
- The arguments of the command as strings.

Returns:

- The standard and error output of the command as a string.
- If the command fails, this will cause the script failure.

Examples:

    // zip files of foo directory in bar.zip file
    files = run("zip", "-r", "bar.zip", "foo")

throw
-----

Throw an error that will cause script failure.

Arguments:

- The error message of the failure.

Returns:

- Nothing, but sets the variable 'error' with the error message.

Examples:

    // stop the script with an error message
    throw("Some tests failed")

uppercase
---------

Put a string in upper case.

Arguments:

- The string to put in upper case.

Returns:

- The string in uppercase.

Examples:

    // set greetings in upper case
    greetings = uppercase(greetings)

