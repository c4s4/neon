# NeON code review — 2026-11-01

Review of the `develop` branch (commit `15fb03c`). `go vet`, `golangci-lint` and `go test -race ./...` all pass: the issues below are not covered by the current test suite.

Items marked **[reproduced]** were verified with a binary built from source and a minimal build file (given each time). The others come from reading the code.

Each item is tagged with its severity: Critical, Major, Minor or Quality (quality/maintenance).

## 1. Security

All security issues have been fixed: see the [Fixed](#fixed) section at the end of this document.

## 2. Functional bugs

All functional bugs have been fixed: see the [Fixed](#fixed) section.

## 3. User experience

All user experience issues have been fixed: see the [Fixed](#fixed) section.

## 4. Code quality and maintenance

Code quality issues have been fixed (see the [Fixed](#fixed) section), except one that was left as is on purpose:

- [Quality] Mutable global state (`Gray`, colors, `ErrorOutput`, `UserRepository`, `CommandLineRepository`). **Not changed:** these are settings written once by `main` before any build runs. The other globals are either registries filled by `init()` and read-only afterwards (`TaskMap`, `BuiltinMap`), or required by the build (`NeonVersion`, set with `-ldflags`). Injecting these settings would not make builds testable in parallel, because the current directory is global to the process anyway (see 2.23). It would only add plumbing through every message call.

## 5. Tooling, CI and documentation

All tooling, CI and documentation issues have been fixed: see the [Fixed](#fixed) section.

## Suggested priorities

1. Add the `CODECOV_TOKEN` secret to the GitHub repository and check the first runs of the GitHub Actions workflow (it could not be run locally).
2. Release 2.0.0: YAML 1.2 is a breaking change for build files that use `yes`/`no`/`on`/`off` as booleans.

---

## Fixed

Fixed on 2026-11-01. All tests pass (`go test -race ./...`, `golangci-lint`, `test/bugs/*.yml`). The CHANGELOG 2.0.0 entry lists the breaking changes (the other changes are described in this section), `doc/tasks.md` has been regenerated and `doc/usermanual.md` updated.

New tests: `neon/build/fixes_test.go`, `neon/task/fixes_test.go`, `neon/task/archive_test.go`, `neon/task/request_test.go`, `neon/util/fixes_test.go`, and regression build files `test/bugs/8.yml` (diamond inheritance) and `test/bugs/9.yml` (stack after a caught error). Both build files fail with the previous binary.

### 1.1 [Critical] Zip Slip / Tar Slip in `unzip` and `untar`

`unzip.go` and `untar.go` built the destination path with `filepath.Join(dir, name)` without checking that it stayed inside `dir`: an entry `../../.bashrc` was written outside the target directory.
**Fix:** new `SafeJoin` helper (`neon/task/util.go`) that rejects any entry name that is not local (`filepath.IsLocal`: absolute paths, `..` components, Windows reserved names). Extraction now fails with `illegal path in archive: '<name>'` before writing anything for that entry. Tests in `neon/task/archive_test.go` cover `../x`, `foo/../../x` and `/tmp/x` for both formats.
**Update (2.26):** extraction now writes through `os.Root`, so a symbolic link (from the archive or already in the destination directory) can't redirect a write outside of it either.

### 1.2 [Major] Update and version check over plain HTTP

`update.go` fetched the latest version and suggested `sh -c "$(curl http://sweetohm.net/neon/install)"` over plain HTTP, open to MITM.
**Fix:** HTTPS for the version URL, and the suggested install command now uses `https://sweetohm.net/dist/neon/install` (the URL from the README; the old `/neon/install` returned 404). The response body is now also closed when the status is not 200.
**Still open:** the binary downloaded by the install script is not verified with a checksum or signature (server-side script, outside this repository). **Update:** `-update` now reads `https://sweetohm.net/dist/neon/version`, written by the `deploy` target of `build.yml` at each release (it was 404 before, so `-update` never reported a new release). Versions with a leading `v` (git tags such as `v2.0.0`) are now accepted as release versions, and a new release is reported only if it is greater than the running one.

### 1.3 [Minor] Singleton port open on all interfaces

`ListenPort` listened on `:port` (all interfaces) and accepted port 0 (random port, so no lock).
**Fix:** listens on `127.0.0.1` and accepts only ports 1 to 65535 (`neon/build/build.go`). The user manual says so. Test `TestListenPort` checks rejected ports, the loopback address and that a second listener on the same port fails.

### 1.4 [Minor] HTTP calls without timeout

`request` and the `-update` version check could block forever.
**Fix:** the version check times out after 10 seconds. The `request` task has a new optional `timeout` argument in seconds, defaulting to 60 (`TestRequestTimeout`). This is a behavior change: a request longer than 60 seconds now fails unless `timeout` is raised.

### 1.5 [Minor] CI: Codecov uploader via `bash <(curl …)`

`.travis.yml` piped Codecov's deprecated bash script into bash.
**Fix:** downloads the Codecov CLI, checks the GPG signature of its checksum file against the Codecov key (fingerprint `27034E7FDB850E0BBC2C62FF806BB28AED779869`), then the SHA256 checksum, before running `codecov upload-process`. The signature check was run locally. The full Travis job was not run.

### 2.1 [Critical] False positive in cycle detection with diamond inheritance

The `visited` map was never cleaned up, so a build file extended by two parents was reported as a cycle.
**Fix:** cycle detection uses the inheritance chain of the current branch (a stack) instead of a global set (`newBuildInternal` in `build.go`). The cycle is printed in order. A real cycle is still detected.

### 2.2 [Critical] Panics on invalid build files

**Fix:** every case of the table now gives an error message:
- `NewSteps` checks the kind before `IsNil` (`steps must be a list (string provided)`);
- `NewTaskStep` formats non-string keys (`unknown task '1'`) and rejects steps that name two tasks;
- `IsValueOfType` handles `nil` values and empty maps, and checks every map entry;
- `EvaluateObject` keeps `nil` values in lists and maps (`toValue` / `setValue`) and reports type mismatches instead of panicking;
- `threads` rejects less than 1 thread.

### 2.3 [Major] Target stack corrupted after an error caught by `try`

**Fix:** `Target.Run` and `RunParentTarget` pop the stack with `defer` right after a successful push, and `Stack.Push` no longer appends the target when it detects a loop. `TestStack` was updated accordingly.

### 2.4 [Major] Inconsistent precedence of inherited targets

**Fix:** `GetTargets` gives precedence to the last parent, like `GetTarget` and `GetProperties`: `-info` now shows the target that runs.

### 2.5 [Major] The `version:` field is unusable with binaries built by `neon build`

**Fix:** the version check is skipped when `NeonVersion` is not a release number (`IsReleaseVersion`), and `-update` says it runs a development version instead of announcing a new release. Version expressions of parents are now checked too, and a leading `=` is accepted. Version comparison was fixed afterwards (2.25).

### 2.6 [Major] `3>` redirection corrupted, `3=` data race

**Fix:** redirected files are opened once and shared between stdout and stderr. Output captured in properties goes through a mutex-protected builder (`syncBuilder`). `1=` and `2=` with different property names now capture separately.

### 2.7 [Major] `1>>` / `2>>` / `3>>` fail if the file doesn't exist

**Fix:** append redirections open the file with `O_CREATE`.

### 2.8 [Major] `untar` doesn't truncate existing files

**Fix:** files are opened with `O_TRUNC`, directory entries are created, and modification times are restored. Links are created since 2.26.

### 2.9 [Major] `zip` doesn't compress

**Fix:** `header.Method = zip.Deflate`. The 200 KB test file now compresses to a few hundred bytes.

### 2.10 [Major] Ignored close errors

**Fix:** `writeTar` and `writeZip` close the tar, gzip and zip writers explicitly, in order, and return their errors, as well as the error when closing the file. `unzip` closes each entry right after extraction (no more `defer` in the loop) and checks the error.

### 2.11 [Major] Infinite loop if a `dotenv` file is a directory

**Fix:** a single dotenv parser (`ReadDotEnv` in `context.go`), also used by the `dotenv` task, based on `bufio.Scanner` with lines up to 1 MB, which reports read errors. It accepts the `export` prefix and single or double quotes. The `dotenv` task parser also crashed on empty lines; it now uses the same parser.

### 2.12 [Major] Environment variables

**Fix (`EvaluateEnvironment`):**
- an undefined `${VAR}` is left unchanged, `$` included;
- errors in `={...}` expressions are reported;
- variables set to an empty string are unset, as the user manual says (before, the value of the process environment was kept);
- escaping handles any number of backslashes;
- variables are evaluated after the build variables they reference, whatever their names, and cycles are reported.

### 2.13 [Major] The build file's `repository:` field is never used

**Fix:** the repository is now taken from, in order: the `-repo` option, the `repository:` field of the build file, the repository of the configuration file (`UserRepository`), then `~/.neon`.

### 2.14 [Major] The `grey` and `time` options of the configuration file are overwritten

**Fix:** command line options are combined with configuration (`opts.Grey || configuration.Grey`). The command line is now parsed before the configuration is loaded, so `-version` and `-help` work with a broken configuration file.

### 2.15 [Minor] `$` (list) modifies the process `PATH`

**Fix:** the executable is looked up in the computed `PATH` (`lookPath`) without changing the process environment. The global current directory is handled in 2.23.

### 2.16 [Minor] `copy` of a file onto itself empties it

**Fix:** `CopyFile` fails with `source and destination are the same file` when `os.SameFile` matches.

### 2.17 [Minor] `replace` resets permissions and its order is random

**Fix:** the file keeps its mode, and the longest strings are replaced first (deterministic result). The error message no longer mentions the copy task.

### 2.18 [Minor] `threads` task

**Fix:** less than 1 thread is an error. The first error stops other threads, which no longer start new iterations. The "done" message is printed before the task ends. The variable no longer shadows `error`.

### 2.19 [Minor] `neon` task

**Fix:** relative paths are resolved from the build directory (current directory for templates). Recursive calls are detected through the chain of called build files kept in the context (`Context.Calls`). Passed properties are set before build properties are evaluated (`Context.InitWithProperties`). The task help documents this.

### 2.20 [Minor] Resolution of parent, script and template paths

**Fix:**
- repository regexps are anchored and escape the `.`;
- `../` paths and relative paths without `./` to an existing `.yml` file are relative to the build file;
- `filepath.IsAbs` replaces `path.IsAbs`;
- `./` context scripts of a parent are resolved relative to the parent file (new `FileDir` field);
- versions of parents are checked and `singleton` is inherited from parents (`GetSingleton`);
- unused `RegexpLinkName` removed.

Relative `dotenv` files of parents are handled in 2.24.

### 2.21 [Minor] Miscellaneous

**Fix:**
- `ExpandUserHome` uses `os.UserHomeDir()` (no more nil dereference) and expands `~` alone;
- `checkUnknownArgs` breaks out of its loop and its condition is explicit;
- two tasks in the same step give a deterministic error;
- `move` falls back to copy and remove across file systems (`EXDEV`);
- `FindFiles` removes duplicates, skips broken symbolic links and reports glob errors other than "not found";
- `-update` skips repository directories that are not git repositories;
- the `chdir` error message shows the right directory and the cause;
- the configuration is loaded after the command line is parsed (see 2.14).

### 2.22 [Major] `singleton` task with `wait: true` busy-loops

**Fix:** the wait loop pauses 20 ms between attempts (`singletonRetryDelay`; a 1 second pause first made the timing of `test/task/singleton.yml` fail), and an invalid port fails immediately instead of waiting forever.

### 2.23 [Minor] Current directory is global to the process

`Target.Run`, `call`, `chdir` and `neon` change the process current directory, which is shared by the threads of the `threads` task. A real per-thread directory is not possible without reworking builtins such as `exists()` or `find()`, which have no access to the context.
**Fix (chosen behaviour: explicit error):** directory changes go through `build.Chdir`. While threads are running (`EnterParallel` / `LeaveParallel` in `threads`), changing to a directory other than the current one fails with `can't change directory to '…' while threads are running`. Changing to the current directory is still accepted, so calling targets of the same build in threads works. Running another build with the `neon` task, or a `chdir` to another directory, inside `threads` is now an error. Before, it silently mixed the directories of the threads. The `threads` and `neon` task help documents it.

### 2.24 [Minor] Relative `dotenv` files of parents are resolved in the child directory

**Fix (chosen behaviour: relative to the parent file):** relative `dotenv` paths are resolved when the build file is loaded: relative to the file that declares them for parents, like context scripts, and relative to the build directory for the main build file (unchanged). This is a behaviour change for parents that relied on loading the `.env` file of the project. The user manual documents it.

### 2.25 [Minor] Version comparison details

**Fix:** `Version.Compare` treats missing fields as zeros (`1.0` == `1.0.0`), then orders suffixes (`SNAPSHOT` < `ALPHA` < `BETA` < `RC` < release) and suffix numbers (`1.0-RC-1` < `1.0-RC-2`). `TestCompare`, which specified `0` < `0.0`, was updated. The four version builtins share `CompareNeonVersion`, which reports `NeON development version '…' can't be compared with a version number` instead of an obscure message.
Using `panic` in builtins is not a bug: it is how Anko reports errors raised by Go functions. A function returning `(bool, error)` would be seen as a list of two values in expressions.

### 2.26 [Minor] Other leftovers

**Fix:**
- `untar` creates symbolic and hard links, and `unzip` creates symbolic links. Extraction writes through `os.Root`, which refuses any write outside the destination directory, even through a link from the archive or one already in place. A hard link whose target is outside the directory is rejected.
- the `neon` task no longer passes the repository of the calling build: the called build uses its own `repository:` field, unless one was set with `-repo` on the command line (`CommandLineRepository`).
- `FindFiles` validates include and exclude patterns and reports invalid ones (such as `[bad`).

### 3.1 [Minor] Errors printed on standard output

`PrintError` wrote on stdout, while `update.go` used `println` (stderr) for some messages.
**Fix:** `PrintError` writes on standard error (`ErrorOutput`, replaceable in tests), in color only if standard error is a terminal. Git outputs and errors of `-update` also go to standard error, through the same writer. Normal messages (titles, `OK`, build output) stay on standard output. This is a behaviour change for scripts that read `ERROR` messages from the standard output of neon.

### 3.2 [Minor] Errors hiding their cause

Fixed along with the functional bugs: errors of command line properties and configuration files now include the YAML error, and a command not found in `PATH` includes the cause.

### 3.3 [Minor] Build duration always printed beyond 10 s

**Fix:** the `time` setting of the configuration file has three states. `true` always prints the duration, `false` never prints it, and when not set the duration is printed for builds longer than 10 seconds (`AutoTimeThreshold`), as before. The `-time` option always prints it. The rule lives in `PrintDuration` (tested). The configuration example of the user manual no longer sets `time: false`, which would now disable the duration.

### 4.1 [Quality] Dependency on `gopkg.in/yaml.v2`

**Fix (chosen behaviour: migrate to YAML 1.2):** neon uses `go.yaml.in/yaml/v3`, the maintained module of the YAML organization. `yes`, `no`, `on` and `off` are now strings; only `true` and `false` are booleans. **This is a breaking change** (CHANGELOG entry bumped to 2.0.0, user manual updated). No build file in this repository or in `~/.neon` used these values as booleans.
To keep the rest of the code unchanged, `util.DecodeYAML` walks the YAML node tree and produces the same types as before: mappings as `map[interface{}]interface{}`, sequences as `[]interface{}`, timestamps kept as strings (v3 would decode them as `time.Time`). Anchors, aliases and merge keys (`<<`) are supported. All decoding goes through `util` (`DecodeYAML`, `DecodeYAMLObject`, `UnmarshalYAML`), so there is a single place to change. `gopkg.in/yaml.v2` stays as an indirect dependency, required by `github.com/c4s4/changelog`.

### 4.2 [Quality] `github.com/mcuadros/go-version`

**Fix:** `sortversions` uses NeON's own `Version` (`build.SortVersions`), and the dependency is removed. It ignores a leading `v` (git tags), orders suffixes, and places invalid versions first. The suffix number may also be written `RC1` or `RC.1`. On the git tags of this repository, the result is the same as with the old library.

### 4.3 [Quality] Imports grouping

**Fix:** new `.golangci.yml` enabling the `gofmt` (simplify) and `gci` formatters, with three groups: standard library, third party, `github.com/c4s4/neon`. All files were reformatted with `golangci-lint fmt`, which only changes imports. `neon lint` now reports formatting issues, and the `fmt` target of `build.yml` runs `golangci-lint fmt` (AGENTS.md updated).

### 4.4 [Quality] Heavy use of `reflect`

**Fix:** `EvaluateObject` uses a type switch for YAML values (strings, `[]interface{}`, `map[interface{}]interface{}`), split into small functions. Reflection is kept only for other typed slices and maps returned by expressions. Validating and copying task arguments still relies on reflection: tasks declare their arguments as structures with `neon` tags, so reflection is the mechanism itself there, not an accident.

### 4.5 [Quality] Leftover code

**Fix:** removed the commented-out `TODO` in `EvaluateExpression`, and `request` uses `bytes.NewReader(body)`.

### 4.6 [Quality] Redundant merges of inherited data

**Fix:** `GetProperties`, `GetEnvironment` and `GetDotEnv` use the data of parents, already merged when they were loaded, instead of merging recursively again. `dotenv` files and context scripts inherited several times (diamond inheritance) are now loaded once, at their first position. Before, the scripts were run twice.

### 5.1 [Major] Inconsistent Go versions

**Fix:** `go.mod` (`go 1.27`) is the single source of truth. The CI reads it (`go-version-file: go.mod`), and `AGENTS.md` and `README.md` now say Go 1.27. The code needs at least Go 1.25 (`os.Root` methods), so 1.27 was kept rather than lowered.

### 5.2 [Minor] Travis CI no longer running

**Fix:** `.travis.yml` is replaced by `.github/workflows/ci.yml` (GitHub Actions), on pushes to `master` and `develop` and on pull requests, with read-only permissions. Actions are pinned by commit SHA. Two jobs:
- tests with race detector and coverage (uploaded with the Codecov action), then the `test/bugs/` build files run with the freshly built binary;
- golangci-lint v2.15.0, which also checks formatting.

The README badge points to the workflow. **To do by the maintainer:** add the `CODECOV_TOKEN` secret. The workflow YAML was checked, but the workflow itself could not be run locally.

### 5.3 [Minor] CHANGELOG missing changes of the previous review

**Fix:** the 1.16.0 entry (2026-05-05, date of the "Merge develop-code-review" commit) now mentions inheritance cycle detection and the singleton listener goroutine leak. The other changes of that review (`Options` struct, `PrintError` no longer exiting) are internal and have no place in the CHANGELOG. The 2.0.0 entry mentions the move to GitHub Actions and the pinned tools.

### 5.4 [Quality] `doc/code-review.md` in user documentation

**Fix:** its content is merged in the [Previous review](#previous-review-2026-05-05) section below, and the file is removed from `doc/`.

### 5.5 [Quality] Outdated `staticcheck`

**Fix:** staticcheck is one of the standard linters of golangci-lint v2, now made explicit in `.golangci.yml`, so it doesn't need to be installed separately. While fixing this, a real problem showed up: `neon tools` installed `github.com/golangci/golangci-lint/cmd/golangci-lint@latest`, which is **v1** and can't read the v2 `.golangci.yml`. `build.yml` now overrides `GOTOOLS` with pinned versions (gox v1.0.1, golangci-lint v2.15.0, gobinsec v0.10.2). The golangci-lint version is the same as in the CI. `AGENTS.md` explains that they must be kept in sync. The old `staticcheck` binary in `~/.gotools/bin` can be deleted.
**Also fixed:** `AGENTS.md` said the ldflags example of the README had a stale variable path, but it matches `build.yml`; the sentence was corrected.

---

## Previous review (2026-05-05)

Summary of `doc/code-review.md` (in French, merged here and removed), with the status of each item. Fixes were released in 1.16.0 unless stated otherwise:

- `ParseCommandLine` returned 24 values: replaced with an `Options` struct.
- `PrintError` called `os.Exit(1)`, which prevented deferred cleanups and testing: errors are now returned up to `main`, which exits.
- The singleton `ListenPort` goroutine looped forever after the listener was closed: the loop stops on `Accept` error.
- Cyclic inheritance (`extends`) caused a stack overflow: cycle detection was added. Its implementation reported diamond inheritance as a cycle; that regression is fixed in 2.1 of this review.
- Hard-coded configuration path: `NEON_CONFIG_PATH` environment variable added.
- Migration to `gopkg.in/yaml.v3` suggested: done in 4.1 of this review (with `go.yaml.in/yaml/v3`).
- Interfaces for `Build`, `Target` and `Context` suggested for mocking: rejected by the maintainer as not useful.
- Anko updated to 0.1.12 (since then 0.1.15).

