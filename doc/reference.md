Tasks Reference
===============

$
-

Execute a command and return output and value.

Arguments:

- $: command to run (string or list of strings).
- +: options to pass on command line after command (strings, optional).
- n=: write command output into named property. Values for n are: 1 for stdout,
  2 for stderr and 3 for stdout and stderr.
- n>: write command output in named file. Values for n are: 1 for stdout,
  2 for stderr and 3 for stdout and stderr.
- n>>: append command output to named file. Values for n are: 1 for stdout,
  2 for stderr and 3 for stdout and stderr.
- nx: disable command output. Values for n are: 1 for stdout, 2 for stderr and
  3 for stdout and stderr.

Examples:

    # execute ls command and get result in 'files' variable
    - $:  'ls -al'
      1=: 'files'
    # execute command as a list of strings and output on console
	- $: ['ls', '-al']
	# run pylint on all python files except those in venv
	- $: 'pylint'
	  +: '=filter(find(".", "**/*.py"), "venv/**/*.py")'

Notes:

- Commands defined as a string run in the shell defined by shell field at the
  root of the build file (or 'sh -c' on Unix and 'cmd /c' on Windows by
  default).
- Defining a command as a list of strings is useful on Windows. Default shell on
  Windows is 'cmd' which can't properly manage arguments with spaces.
- Argument of a command defined as a list won't be expanded by shell. Thus
  $USER won't be expanded for instance.

assert
------

Make an assertion and fail if assertion is false.

Arguments:

- assert: the assertion to perform (boolean, expression).

Examples:

    # assert that foo == "bar", and fail otherwise
    - assert: 'foo == "bar"'

call
----

Call a build target.

Arguments:

- call: the name of the target(s) to call (strings, wrap).

Examples:

    # call target 'foo'
    - call: 'foo'

cat
---

Print the content of a given file on the console.

Arguments:

- cat: the name of the file to print on console (string, file).

Examples:

    # print content of LICENSE file on the console
    - cat: "LICENSE"

chdir
-----

Change current working directory.

Arguments:

- chdir: the directory to change to (string, file).

Examples:

    # go to build directory
    - chdir: "build"

Notes:

- Working directory is set to the build file directory before each target.

chmod
-----

Change mode of files.

Arguments:

- chmod: list of globs of files to change mode (strings, file, wrap).
- mode: mode to change to (integer).
- dir: the root directory for globs, defaults to '.' (string, optional, file).
- exclude: globs of files to exclude (strings, optional, file, wrap).

Examples:

    # make foo.sh executable for all users
    - chmod: "foo.sh"
      mode:  0755
    # make all sh files in foo directory executable, except for bar.sh
    - chmod:   "**/*.sh"
      mode:    0755
      exclude: "**/bar.sh"

Notes:

- The mode is an integer, thus must not be surrounded with quotes, or it would
  be a string and parsing of the task would fail.
- We usually set mode with octal integers, starting with '0'. If you don't put
  starting '0', this is decimal integer and you won't probably have expected
  result.

classpath
---------

Build a Java classpath.

Arguments:

- classpath: the property to set with classpath (string).
- classes: class directories to add in classpath (strings, optional, file,
  wrap).
- jars: globs of jar files to add to classpath (strings, optional, file, wrap).
- dependencies: dependency files to add to classpath (strings, optional, file,
  wrap).
- scopes: classpath scope (strings, optional, wrap). If set, will take
  dependencies without scope and listed scopes, if not set, will only take
  dependencies without scope).
- repositories: repository URLs to get dependencies from, defaults to
  'http://repo1.maven.org/maven2' (strings, optional, wrap).
- todir: directory to copy jar files into (string, optional, file).

Examples:

    # build classpath with classes in build/classes directory
    - classpath: 'classpath'
      classes:   'build/classes'
    # build classpath with jar files in lib directory
    - classpath: 'classpath'
      jars:      'lib/*.jar'
    # build classpath with a dependencies file
    - classpath:    'classpath'
      dependencies: 'dependencies.yml'
    # copy classpath's jar files to 'build/lib' directory
    - classpath:    _
      dependencies: 'dependencies.yml'
      todir:        'build/lib'

Notes:

- Dependency files should list dependencies with YAML syntax as follows:

    - group:    junit
      artifact: junit
      version:  4.12
      scopes:   [test]

- Scopes are optional. If not set, dependency will always be included. If set,
  dependency will be included for classpath with these scopes.

copy
----

Copy file(s).

Arguments:

- copy: globs of files to copy (strings, file, wrap).
- dir: root directory for globs, defaults to '.' (string, optional, file).
- exclude: globs of files to exclude (strings, optional, file, wrap).
- tofile: file to copy file to (string, optional, file).
- todir: directory to copy files to (string, optional, file).
- flat: tells if files should be flatten in destination directory, defaults to
  false (boolean, optional).

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

Notes:

- Parameter 'tofile' is valid if only one file was selected by globs.
- One and only one of parameters 'tofile' and 'todir' might be set.

delete
------

Delete files or directories (recursively).

Arguments:

- delete: glob of files or directories to delete (strings, file, wrap).
- dir: root directory for globs (string, optional, file).
- exclude: globs of files to exclude (strings, optional, file, wrap).

Examples:

    # delete build directory
    - delete: =BUILD_DIR
    # delete all XML files except 'foo.xml'
    - delete:  "**/*.xml"
      exclude: "**/foo.xml"

Notes:

- Handle with care: if globs select directories, they are deleted recursively!

for
---

For loop.

Arguments:

- for: variable name to set at each loop iteration (string).
- in: values or expression to generate values to iterate on (list or
  expression).
- do: steps to execute at each loop iteration (steps).

Examples:

    # create empty files
    - for: file
      in:  ["foo", "bar"]
      do:
    - touch: =file
    # print first 10 integers
    - for: i
      in: range(10)
      do:
      - print: '={i}'

if
--

If condition.

Arguments:

- if: the condition (boolean, expression).
- then: steps to execute if condition is true (steps).
- else: steps to execute if condition is false (optional, steps).

Examples:

    # print hello if x > 10 else print world
    - if: x > 10
      then:
      - print: "hello"
      else:
      - print: "world"

java
----

Run Java virtual machine.

Arguments:

- javac: main Java class name (string).
- cp: classpath to run main class (string).
- args: command line arguments (strings, optional, wrap).

Examples:

    # run class foo.Bar with arguments foo and bar
    - javac: 'foo.Bar'
      cp:    'build/classes'
      args:  ['foo', 'bar']

javac
-----

Compile Java source files.

Arguments:

- javac: glob of Java source files to compile (strings, file, wrap).
- source: directory of source files (string, file).
- exclude: glob of source files to exclude (strings, optional, file, wrap).
- dest: destination directory for generated classes (string, file).
- cp: classpath for compilation (string, optional).

Examples:

    # compile Java source files in src directory
    - javac:  '**/*.java'
      source: 'src'
      dest:   'build/classes'
    # compile Java source files in src directory with given classpath
    - javac:  '**/*.java'
      source: 'src'
      dest:   'build/classes'
      cp:     =classpath

link
----

Create a symbolic link.

Arguments:

- link: source file (string, file).
- to: destination of the link (string, file).

Examples:

    # create a link from file 'foo' to 'bar'
    - link: 'foo''
      to:   'bar''

mkdir
-----

Make a directory.

Arguments:

- mkdir: directories to create (strings, file, wrap).

Examples:

    # create a directory 'build'
    - mkdir: 'build'

move
----

Move file(s).

Arguments:

- move: globs of files to move (strings, file, wrap)
- dir: root directory for globs (string, optional, file).
- exclude: globs of files to exclude (strings, optional, file, wrap).
- tofile: file to move file to (string, optional, file).
- todir: directory to move file(s) to (string, optional, file).
- flat: tells if files should be flatten in destination directory, defaults to
  false (boolean, optional).

Examples:

    # move file foo to bar
    - move:   'foo'
      tofile: 'bar'
    # move text files in directory 'book' (except 'foo.txt') to directory 'text'
    - move:    '**/*.txt'
      dir:     'book'
      exclude: '**/foo.txt'
      todir:   'text'
    # move all go sources to directory 'src', flattening structure
    - move:  '**/*.go'
      todir: 'src'
      flat:  true

Notes:

- Parameter 'tofile' is valid if only one file was selected by globs.
- One and only one of parameters 'tofile' and 'todir' might be set.

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

- path: globs of files to build the path (strings, file, wrap).
- to: variable to put path into (string).
- dir: root directory for globs, defaults to '.' (string, optional, file).
- exclude: globs of files to exclude (strings, optional, file, wrap).

Examples:

    # build classpath with jar files in lib directory
    - path: 'lib/*.jar'
      to:   'classpath'

print
-----

Print a message on the console.

Arguments:

- print: text to print (string).

Examples:

    # say hello
    - print: 'Hello World!'

prompt
------

Prompt the user for the value of a given property matching a pattern.

Arguments:

- prompt: message to print at prompt that include a description of expected
  pattern (string).
- to: name of the property to set (string).
- default: default value if user doesn't type anything, written into square
  brackets after prompt message (string, optional).
- pattern: a regular expression for prompted value. If this pattern is not
  matched, this task will prompt again. If no pattern is given, any value is
  accepted (string, optional).
- error: error message to print when pattern is not matched (string, optional).

Examples:

    # prompt for age that is a positive number
    - prompt:  'Enter your age'
      to:      'age'
      default: '18'
      pattern: '^\d+$'
      error:   'Age must be a positive integer'

read
----

Read given file as text and put its content in a variable.

Arguments:

- read: file to read (string, file).
- to: name of the variable to set with its content (string).

Examples:

    # put content of LICENSE file in license variable
    - read: 'LICENSE'
      to:   'license'

replace
-------

Replace text matching patterns in files.

Arguments:

- replace: globs of files to process (strings, file, wrap).
- with: map with replacements (map with string keys and values).
- dir: root directory for globs (string, optional, file).
- exclude: globs of files to exclude (strings, optional, files).

Examples:

    # replace foo with bar in file test.txt
    - replace: 'test.txt'
      with:    {'foo': 'bar'}

request
-------

Perform an HTTP request.

Arguments:

- request: URL to request (string).
- method: request method ('GET', 'POST', etc), defaults to 'GET' (string,
  optional).
- headers: request headers (map with string keys and values, optional).
- body: request body (string, optional).
- file: request body as a file (string, optional, file).
- username: user name for authentication (string, optional).
- password: user password for authentication (string, optional).
- status: expected status code, raise an error if different, defaults to 200
  (int, optional).

Examples:

    # get google.com
    - request: 'google.com'

Notes:

- Response status code is stored in variable _status.
- Response body is stored in variable _body.
- Response headers are stored in variable _headers.

sleep
-----

Sleep given number of seconds.
		
Arguments:

- sleep: duration to sleep in seconds (float).

Examples:

    # sleep for 1.5 seconds
    - sleep: 1.5
    # sleep for 3 seconds (3.0 as a float)
    - sleep: 3.0

super
-----

Call target with same name in parent build file.

Arguments:

- none

Examples:

    # call parent target
    - super:

Notes:

- This will raise en error if parent build files have no target with same name.

tar
---

Create a tar archive.

Arguments:

- tar: globs of files to tar (strings, file, wrap).
- dir: root directory for glob, defaults to '.' (string, optional, file).
- exclude: globs of files to exclude (strings, optional, file, wrap).
- tofile: name of the tar file to create (string, file).
- prefix: prefix directory in the archive (optional).

Examples:

    # tar files in build directory in file named build.tar.gz
    - tar:    'build/**/*'
      tofile: 'build.tar.gz'

Notes:

- If archive filename ends with gz (with names such as 'foo.tar.gz' or
  'foo.tgz') the tar archive is also gzip compressed.

threads
-------

Run steps in threads.

Arguments:

- threads: number of threads to run (integer).
- input: values to pass to threads in _input property (list, optional).
- steps: steps to run in threads (steps).
- verbose: if you want thread information on console, defaults to false
  (boolean, optional).

Examples:

    # compute squares of 10 first integers in threads and put them in _output
    - threads: =_NCPU
      input:   =range(10)
      steps:
      - '_output = _input * _input'
      - print: '#{_input}^2 = #{_output}'
    # print squares on the console
    - print: '#{_output}'

Notes:

- You might set number of threads to '_NCPU' which is the number of cores in
  the CPU of the machine.
- Property _thread is set with the thread number (starting with 0)
- Property _input is set with the input for each thread.
- Property _output is set with the output of the threads.
- Each thread should write its output in property _output.

Context of the build is cloned in each thread so that you can read and write
properties, they won't affect other threads. But all properties will be lost
when thread is done, except for _output that will be appended to other in
resulting _output property.

Don't change current directory in threads as it would affect other threads as
well.

throw
-----

Throws an error.

Arguments:

- throw: the message of the error (string).

Examples:

    # stop the build because tests failed
    - throw: "ERROR: tests failed"

Notes:

- You can catch raised errors with try/catch/finally task.
- Property _error is set with the error message.
- If the error was not catch, the error message will be printed on the console
  as the cause of the build failure.

time
----

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

touch
-----

Touch a file (create it or change its time).

Arguments:

- touch: files to touch (strings, file, wrap).

Examples:

    # create file in build directory
    - touch: ['#{BUILD_DIR}/foo', '#{BUILD_DIR}/bar']

Notes:

- If the file already exists it changes it modification time.
- If the file doesn't exist, it creates an empty file.

try
---

Try/catch/finally construct.

Arguments:

- try: steps to execute (steps).
- catch: executed if an error occurs (steps, optional).
- finally: executed in any case (steps, optional).

Examples:

    # execute a command and continue even if it fails
    - try:
      - 'command-that-doesnt-exist'
	- print: 'Continue even if command fails'
	# execute a command and print a message if it fails
	- try:
	  - 'command-that-doesnt-exist'
	  catch:
	  - print: 'There was an error!'
	# execute a command a print message in all cases
	- try:
	  - 'command-that-doesnt-exist'
	  finally:
	  - print: 'Print whatever happens'

Notes:

- The error message for the failure is stored in '_error' variable as text.

untar
-----

Expand a tar file in a directory.

Arguments:

- untar: the tar file to expand (string, file).
- todir: the destination directory (string, file).

Examples:

    # untar foo.tar to build directory
    - untar: 'foo.tar'
      todir: 'build'

Notes:

- If archive filename ends with .gz (with a name such as foo.tar.gz or foo.tgz)
  the tar archive is uncompressed with gzip.

unzip
-----

Expand a zip file in a directory.

Arguments:

- unzip: the zip file to expand (string, file).
- todir: the destination directory (string, file).

Examples:

    # unzip foo.zip to build directory
    - unzip: 'foo.zip'
      todir: 'build'

while
-----

While loop.

Arguments:

- while: condition evaluated at each iteration (string).
- do: steps that run while condition is true (steps).

Examples:

    # loop until i >= 10
    - while: 'i < 10'
      do:
      - script: 'println(i++)'

write
-----

Write text into given file.

Arguments:

- write: file to write into (string, file).
- text: text to write into the file (string, optional).
- append: tells if we should append content to file, default to false (boolean,
  optional).

Examples:

    # write 'Hello World!' in file greetings.txt
    - write: 'greetings.txt'
      text:  'Hello World!'

zip
---

Create a Zip archive.

Arguments:

- zip: globs of files to zip (strings, file, wrap).
- dir: root directory for globs, defaults to '.' (string, optional, file).
- exclude: globs of files to exclude (strings, optional, file, wrap).
- tofile: name of the Zip file to create (string, file).
- prefix: prefix directory in the archive (string, optional).

Examples:

    # zip files of build directory in file named build.zip
    - zip:    'build/**/*'
      tofile: 'build.zip'


Builtins Reference
==================

absolute
--------

Return absolute value of a given path.

Arguments:

- The path to get absolute value.

Returns:

- The absolute value of the path.

Examples:

    # get absolute value of path "foo/../bar/spam.txt"
    path = absolute("foo/../bar/spam.txt")
    # returns: "/home/user/build/bar/spam.txt"

contains
--------

Contains strings.

Arguments:

- List of strings to search into.
- Searched string.

Returns:

- A boolean telling if the string is contained in the list.

Examples:

    # Tell if the list contains "bar"
    contains(["foo", "bar"], "bar")
    # returns: true

directory
---------

Return directory of a given path.

Arguments:

- The path to get directory for as a string.

Returns:

- The directory of the path as a string.

Examples:

    # get directory of path "/foo/bar/spam.txt"
    dir = directory("/foo/bar/spam.txt")
    # returns: "/foo/bar"

env
---

Get environment variable.

Arguments:

- The name of the environment variable to get value for.

Returns:

- The value of this environment variable.

Examples:

    # get PATH environment variable
    env("PATH")
    # returns: value of the environment variable PATH

escapeurl
---------

Escape given URL.

Arguments:

- The URL to escape.

Returns:

- The escaped URL.

Examples:

    # escape given URL
    escapeurl("/foo bar")
    # returns: "/foo%20bar"

exists
------

Tells if a given path exists.

Arguments:

- The path to test as a string.

Returns:

- A boolean telling if path exists.

Examples:

    # test if given path exists
    exists("/foo/bar")
    # returns: true if file "/foo/bar" exists

expand
------

Expand file name replacing ~/ with home directory.

Arguments:

- The path to expand as a string.

Returns:

- The expanded path as a string.

Examples:

    # expand path ~/.profile
    profile = expand("~/.profile")
    # returns: "/home/casa/.profile" on my machine

filename
--------

Return filename of a given path.

Arguments:

- The path to get filename for as a string.

Returns:

- The filename of the path as a string.

Examples:

    # get filename of path "/foo/bar/spam.txt"
    filename("/foo/bar/spam.txt")
    # returns: "spam.txt"

filter
------

Filter a list of files with excludes.

Arguments:

- includes: the list of files to filter.
- excludes: a list of patterns for files to exclude.

Returns:

- The list if filtered files as a list of strings.

Examples:

    # filter text files removing those in build directory
    filter(find(".", "**/*.txt"), "build/**/*")
    # returns: files with extension "txt" in current directory and
    # subdirectories, except those in "build" directory

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

    # find all text files in book directory
    find("book", "**/*.txt")
    # returns: list of files with extension "txt"
    # find all xml and yml files in src directory
    find("src", "**/*.xml", "**/*.yml")
    # returns: list of "xml" and "yml" files

Notes:

- Files may be filtered with filter() builtin.

findinpath
----------

Find executables in PATH.

Arguments:

- The executable to find.

Returns:

- A list of absolute paths to the executable, in the order of the PATH.

Examples:

    # find python in path
    findinpath("python")
    # returns: ["/opt/python/current/bin/python", /usr/bin/python"]

followlink
----------

Follow symbolic link.

Arguments:

- The symbolic link to follow.

Returns:

- The path with symbolic links followed.

Examples:

    # follow symbolic link 'foo'
    followlink("foo")
    # returns: 'bar'

join
----

Join strings.

Arguments:

- The strings to join as a list of strings.
- The separator as a string.

Returns:

- Joined strings as a string.

Examples:

    # join "foo" and "bar" with a space
    join(["foo", "bar"], " ")
    # returns: "foo bar"

joinpath
--------

Join file paths.

Arguments:

- The paths to join as a list of strings.

Returns:

- Joined path as a string.

Examples:

    # join paths "/foo", "bar" and "spam.txt"
    joinpath("foo", "bar", "spam.txt")
    # returns: "foo/bar/spam.txt" on a Linux box and "foo\bar\spam.txt" on
    # Windows

jsondecode
----------

Decode given string in Json format.

Arguments:

- The string in Json format to decode.

Returns:

- Decoded string.

Examples:

    # decode given list
    jsondecode("['foo', 'bar']")
    # returns string slice: ["foo", "bar"]

jsonencode
----------

Encode given variable in Json format.

Arguments:

- The variable to encode in Json format.

Returns:

- Json encoded string.

Examples:

    # encode given list
    jsonencode(["foo", "bar"])
    # returns: "['foo', 'bar']"

keys
----

Return keys of gien map.

Arguments:

- The map to get keys for.

Returns:

- A list of keys.

Examples:

    # get keys of a map
    keys({"foo": 1, "bar": 2})
    # returns: ["foo", "bar"]

list
----

Return a list:
- If the object is already a list, return the object.
- If the object is not a list, wrap it into a list.

Arguments:

- The object to turn into a list.

Returns:

- The list.

Examples:

    # get a list of foo
    list(foo)
	# return foo if already a list or [foo] otherwise

lowercase
---------

Put a string in lower case.

Arguments:

- The string to put in lower case.

Returns:

- The string in lower case.

Examples:

    # set string in lower case
    lowercase("FooBAR")
    # returns: "foobar"

newer
-----

Tells if source is newer than result file (if any).

Arguments:

- source: source file that must exist.
- result: result file (may not exist).

Returns:

- A boolean that tells if source is newer than result. If result file doesn't
  exists, this returns true.

Examples:

    # generate PDF if source Markdown changed
    if newer("source.md", "result.pdf") {
    	compile("source.md")
    }

now
---

Return current date and time in ISO format.

Arguments:

- none

Returns:

- ISO date and time as a string.

Examples:

    # put current date and time in dt variable
    now()
    # returns: "2006-01-02 15:04:05"
    # to get date in ISO format
    now()[0:10]
    # returns: "2006-01-02"

ospath
------

Convert path to running OS.

Arguments:

- The path to convert.

Returns:

- The converted path.

Examples:

    # convert path foo/bar to OS format
    path = ospath("foo/bar")
    # will return foo/bar on Unix and foo\bar on Windows

read
----

Read given file and return its content as a string.

Arguments:

- The file name to read.

Returns:

- The file content as a string.

Examples:

    # read VERSION file and set variable version with ots content
    read("VERSION")
    # returns: the contents of "VERSION" file

replace
-------

Replace string with another.

Arguments:

- The strings where take place replacements.
- The substring to replace.
- The replacement substring.

Returns:

- Replaced string.

Examples:

    # replace "foo" with "bar" in string "spam foo eggs"
    replace("spam foo eggs", "foo", "bar")
    # returns: "spam bar eggs"

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

    # zip files of foo directory in bar.zip file
    run("zip", "-r", "bar.zip", "foo")
    # returns: the trimed output of the command

split
-----

Split strings.

Arguments:

- The strings to split.
- The separator for splitting.

Returns:

- A list of strings.

Examples:

    # split "foo bar" with space
    split("foo bar", " ")
    # returns: ["foo"," "bar"]

throw
-----

Throw an error that will cause script failure.

Arguments:

- The error message of the failure.

Returns:

- Nothing, but sets the variable 'error' with the error message.

Examples:

    # stop the script with an error message
    throw("Some tests failed")
    # returns: nothing, the script is interrupted on error

unescapeurl
-----------

Unescape given URL.

Arguments:

- The URL to unescape.

Returns:

- The unescaped URL.

Examples:

    # unescape given URL
    escapeurl("foo%20bar")
    # returns: "foo bar"

unixpath
--------

Convert a path to Unix format.

Arguments:

- The path to convert.

Returns:

- The converted path.

Examples:

    # convert path to unix
    uppercase("c:\foo\bar")
    # returns: "/c/foo/bar"

uppercase
---------

Put a string in upper case.

Arguments:

- The string to put in upper case.

Returns:

- The string in uppercase.

Examples:

    # set string in upper case
    uppercase("FooBAR")
    # returns: "FOOBAR"

windowspath
-----------

Convert a path to Windows format.

Arguments:

- The path to convert.

Returns:

- The converted path.

Examples:

    # convert path to windows
    uppercase("/c/foo/bar")
    # returns: "c:\foo\bar"

winexe
------

Add '.exe' or '.bat' extensions depending on platform:
- command will stay command on Unix and will become command.exe on Windows.
- script.sh will stay script.sh on Unix and will become script.bat on Windows.
It will also replace / with \ in the executable path.

Arguments:

- The command to process.

Returns:

- Command adapted to host system.

Examples:

    # run command foo on unix and windows
    run(winexe("bin/foo"))
    # will run bin/foo on unix and bin\foo.exe on windows
    # run script script.sh unix and windows
    run(winexe("script.sh"))
    # will run script.sh on unix and script.bat on windows

