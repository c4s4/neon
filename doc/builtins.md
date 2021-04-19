# Builtins Reference

## absolute

Return absolute value of a given path.

Arguments:

- The path to get absolute value.

Returns:

- The absolute value of the path.

Examples:

    # get absolute value of path "foo/../bar/spam.txt"
    path = absolute("foo/../bar/spam.txt")
    # returns: "/home/user/build/bar/spam.txt"

## appendpath

Append root directory to paths.

Arguments:

- The root directory.
- The paths to append.

Returns:

- Appended paths as a list.

Examples:

    # append root "foo" to paths "spam" and "eggs"
    appendpath("foo", "spam", "eggs")
	# returns: ["foo/spam", "foo/eggs"] on Linux and
	# ["foo\spam", "foo\eggs"] on Windows

## changelog

Return changelog information from file.

Arguments:

- changelog: the name of the changelog file (look for changelog in current
  directory if empty string).

Note:

- The function returns a Changelog that is a list of Releases struct with
  fields Version, Date and Summary.

Examples:

    # get version of last release:
    - 'VERSION = changelog("")[0].Version'

## contains

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

## directory

Return directory of a given path.

Arguments:

- The path to get directory for as a string.

Returns:

- The directory of the path as a string.

Examples:

    # get directory of path "/foo/bar/spam.txt"
    dir = directory("/foo/bar/spam.txt")
    # returns: "/foo/bar"

## env

Get environment variable.

Arguments:

- The name of the environment variable to get value for.

Returns:

- The value of this environment variable.

Examples:

    # get PATH environment variable
    env("PATH")
    # returns: value of the environment variable PATH

## escapeurl

Escape given URL.

Arguments:

- The URL to escape.

Returns:

- The escaped URL.

Examples:

    # escape given URL
    escapeurl("/foo bar")
    # returns: "/foo%20bar"

## exists

Tells if a given path exists.

Arguments:

- The path to test as a string.

Returns:

- A boolean telling if path exists.

Examples:

    # test if given path exists
    exists("/foo/bar")
    # returns: true if file "/foo/bar" exists

## expand

Expand file name replacing ~/ with home directory.

Arguments:

- The path to expand as a string.

Returns:

- The expanded path as a string.

Examples:

    # expand path ~/.profile
    profile = expand("~/.profile")
    # returns: "/home/casa/.profile" on my machine

## filename

Return filename of a given path.

Arguments:

- The path to get filename for as a string.

Returns:

- The filename of the path as a string.

Examples:

    # get filename of path "/foo/bar/spam.txt"
    filename("/foo/bar/spam.txt")
    # returns: "spam.txt"

## filter

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
    # subdirectories, except those in "github.com/c4s4/neon/neon/build" directory

Notes:

- Works great with find() builtin.

## find

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

## findinpath

Find executables in PATH.

Arguments:

- The executable to find.

Returns:

- A list of absolute paths to the executable, in the order of the PATH.

Examples:

    # find python in path
    findinpath("python")
    # returns: ["/opt/python/current/bin/python", /usr/bin/python"]

## followlink

Follow symbolic link.

Arguments:

- The symbolic link to follow.

Returns:

- The path with symbolic links followed.

Examples:

    # follow symbolic link 'foo'
    followlink("foo")
    # returns: 'bar'

## greater

Check that NeON version is greater that given version.

Arguments:

- The version to check against.

Returns:

- A boolean telling if NeON version is greater than given version.

Examples:

    # check that NeON version is greater than 0.12.0
    greater("0.12.0")
    # return true if version is greater than 0.12.0, false otherwise

## greaterorequal

Check that NeON version is greater or equal that given version.

Arguments:

- The version to check against.

Returns:

- A boolean telling if NeON version is greater or equal than given version.

Examples:

    # check that NeON version is greater or equal than 0.12.0
    greaterorequal("0.12.0")
    # return true if version is greater or equal than 0.12.0, false otherwise

## haskey

Tells if a map contains given key.

Arguments:

- The map to test.
- The key to test.

Returns:

- A boolean telling if the map contains given key.

Examples:

    # Tell if map "map" contains key "key"
    haskey(map, "key")
    # returns: true or false

## join

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

## joinpath

Join file paths.

Arguments:

- The paths to join as a list of strings.

Returns:

- Joined path as a string.

Examples:

    # join paths "foo", "bar" and "spam.txt"
    joinpath("foo", "bar", "spam.txt")
    # returns: "foo/bar/spam.txt" on a Linux box and "foo\bar\spam.txt" on
    # Windows

## jsondecode

Decode given string in Json format.

Arguments:

- The string in Json format to decode.

Returns:

- Decoded string.

Examples:

    # decode given list
    jsondecode("['foo', 'bar']")
    # returns string slice: ["foo", "bar"]

## jsonencode

Encode given variable in Json format.

Arguments:

- The variable to encode in Json format.

Returns:

- Json encoded string.

Examples:

    # encode given list
    jsonencode(["foo", "bar"])
    # returns: "['foo', 'bar']"

## keys

Return keys of gien map.

Arguments:

- The map to get keys for.

Returns:

- A list of keys.

Examples:

    # get keys of a map
    keys({"foo": 1, "bar": 2})
    # returns: ["foo", "bar"]

## length

Return length of given string.

Arguments:

- The string to get length for.

Returns:

- Length of the given string.

Examples:

    # get length of the string
    l = length("Hello World!")

## list

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

## lower

Check that NeON version is lower that given version.

Arguments:

- The version to check against.

Returns:

- A boolean telling if NeON version is lower than given version.

Examples:

    # check that NeON version is lower than 0.12.0
    greater("0.12.0")
    # return true if version is lower than 0.12.0, false otherwise

## lowercase

Put a string in lower case.

Arguments:

- The string to put in lower case.

Returns:

- The string in lower case.

Examples:

    # set string in lower case
    lowercase("FooBAR")
    # returns: "foobar"

## lowerorequal

Check that NeON version is lower or equal that given version.

Arguments:

- The version to check against.

Returns:

- A boolean telling if NeON version is lower or equal than given version.

Examples:

    # check that NeON version is lower or equal than 0.12.0
    lowerorequal("0.12.0")
    # return true if version is lower or equal than 0.12.0, false otherwise

## match

Tell if given string matches a regular expression.

Arguments:

- The regular expression.
- The string to test.

Returns:

- A boolean telling string matches regular expression.

Examples:

    # tell if string "neon" marchs "n..n" regular expression:
    match("n..n", "neon")
    # return true

## newer

Tells if source files are newer than result ones.

Arguments:

- sources: source file(s) (may not exist).
- results: result file(s) (may not exist).

Returns:

- A boolean that tells if source files are newer than result ones.
  If source files don't exist, returns false.
  If result files don't exist, returns true.

Examples:

    # generate PDF if source Markdown changed
    if newer("source.md", "result.pdf") {
    	compile("source.md")
    }
	# generate binary if source files are newer than generated binary
    if newer(find(".", "**/*.go"), "bin/binary") {
    	generateBinary()
    }

## now

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

## ospath

Convert path to running OS.

Arguments:

- The path to convert.

Returns:

- The converted path.

Examples:

    # convert path foo/bar to OS format
    path = ospath("foo/bar")
    # will return foo/bar on Unix and foo\bar on Windows

## read

Read given file and return its content as a string.

Arguments:

- The file name to read.

Returns:

- The file content as a string.

Examples:

    # read VERSION file and set variable version with ots content
    read("VERSION")
    # returns: the contents of "VERSION" file

## replace

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

## run

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

## setenv

Set environment variable.

Arguments:

- The variable name.
- The variable value.

Examples:

    # set foo to value bar
    setenv("foo", "bar")

## sortversions

Sort a list of versions.

Arguments:

- The list of versions to sort.

Returns:

- nothing but slice of versions is sorted

Examples:

    # sort version ["1.10", "1.1", "1.2"]
    sortversions(["1.10", "1.1", "1.2"])
    # returns nothing but slice of versions is sorted

## split

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

## termwidth

Return terminal width.

Arguments:

- None

Returns:

- Terminal width in characters.

Examples:

	# get terminal width
	width = termwidth()

## throw

Throw an error that will cause script failure.

Arguments:

- The error message of the failure.

Returns:

- Nothing, but sets the variable 'error' with the error message.

Examples:

    # stop the script with an error message
    throw("Some tests failed")
    # returns: nothing, the script is interrupted on error

## toint

Converts int64 value to int.

Arguments:

- The int64 value to convert.

Returns:

- Value converted to int.

Examples:

    # convert len([1, 2, 3]) to int
    toint(len([1, 2, 3]))
    # returns: 3

## trim

Trim spaces from given string.

Arguments:

- The string to trim.

Returns:

- Trimed string.

Examples:

    # trim string "\tfoo bar\n   "
    trim("\tfoo bar\n  ")
    # returns: "foo bar"

## unescapeurl

Unescape given URL.

Arguments:

- The URL to unescape.

Returns:

- The unescaped URL.

Examples:

    # unescape given URL
    escapeurl("foo%20bar")
    # returns: "foo bar"

## unixpath

Convert a path to Unix format.

Arguments:

- The path to convert.

Returns:

- The converted path.

Examples:

    # convert path to unix
    uppercase("c:\foo\bar")
    # returns: "/c/foo/bar"

## uppercase

Put a string in upper case.

Arguments:

- The string to put in upper case.

Returns:

- The string in uppercase.

Examples:

    # set string in upper case
    uppercase("FooBAR")
    # returns: "FOOBAR"

## windowspath

Convert a path to Windows format.

Arguments:

- The path to convert.

Returns:

- The converted path.

Examples:

    # convert path to windows
    uppercase("/c/foo/bar")
    # returns: "c:\foo\bar"

## winexe

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

## yamldecode

Decode given string in YAML format.

Arguments:

- The string in YAML format to decode.

Returns:

- Decoded string.

Examples:

    # decode given list
    yamldecode("['foo', 'bar']")
    # returns string slice: ["foo", "bar"]

## yamlencode

Encode given variable in YAML format.

Arguments:

- The variable to encode in YAML format.

Returns:

- Json encoded string.

Examples:

    # encode given list
    yamlencode(["foo", "bar"])
    # returns: "['foo', 'bar']"
