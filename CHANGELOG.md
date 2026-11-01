# NeON

## 2026-11-01: 1.18.0

- Bug fixes
- Build files, configuration files and `yamldecode` builtin are parsed as YAML 1.2: `yes`, `no`, `on` and `off` are strings, not booleans anymore (use `true` and `false`), and unquoted dates are strings
- Error messages are printed on standard error instead of standard output
- Changing to another directory while threads are running (with `chdir` or `neon` tasks) raises an error
- `threads` task rejects less than 1 thread
- A step that names two tasks raises an error
- `request` task times out after 60 seconds by default (see new `timeout` option)
- `unzip` and `untar` tasks fail on archive entries that would be written outside the destination directory
- `repository` field of the build file is now used (it was ignored), and the `neon` task runs called builds with their own repository
- `neon` task resolves relative build file paths from the build directory instead of current directory
- Relative `dotenv` files declared in parent build files are relative to the parent build file
- `version` field of parent build files is checked too, and version comparison counts missing fields as zeros (`1.0` equals `1.0.0`) and orders suffixes (SNAPSHOT < ALPHA < BETA < RC < release)
- Environment variables set to an empty string are unset, undefined `${VAR}` references are left unchanged and errors in expressions are reported
- `copy` of a file onto itself fails instead of emptying the file

## 2026-09-05: 1.17.0

- Added `cafe` option to `$` task to prevent the system from going to sleep while a command is running (Linux and MacOS)

## 2026-08-10: 1.16.1

- Updated dependencies to fix Anko crash

## 2026-05-05: 1.16.0

- `NEON_CONFIG_PATH` sets configuration file path (which defaults to *~/.neon/settings.yml*)
- Updated Anko to 0.1.12
- Fixed build and archive targets
- Cyclic inheritance of build files is detected and reported as an error (instead of a stack overflow)
- Fixed goroutine leak of singleton port listener

## 2025-07-17: 1.15.0

- Added `pause` task

## 2025-05-06: 1.14.1

- Fix dotenv loading

## 2025-05-06: 1.14.0

- Added `properties` key to `neon` task
- Added `write()` builtin to write text into a file

## 2025-04-14: 1.13.1

- Fixed code warnings
- Fix start task

## 2025-03-13: 1.13.0

- Added `jsonindent()` builtin

## 2025-03-06: 1.12.0

- Added `env` field to `shell` task to set environment

## 2025-02-18: 1.11.1

- Added `-batch` option to force plugin update with `-update` option

## 2025-02-13: 1.11.0

- Environment variables defined in build file are defined everywhere in it, not only in executed commands
- Added `dotenv` key to load environment file in build file
