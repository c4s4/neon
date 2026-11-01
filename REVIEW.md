# NeON code review — 2026-11-01

Review of the `develop` branch (commit `15fb03c`). `go vet`, `golangci-lint` and `go test -race ./...` all pass: the issues below are not covered by the current test suite.

Items marked **[reproduced]** were verified with a binary built from source and a minimal build file (given each time). The others come from reading the code.

Each item is tagged with its severity: Critical, Major, Minor or Quality (quality/maintenance).

## 1. Security

### 1.1 [Critical] Zip Slip / Tar Slip in `unzip` and `untar`

`neon/task/unzip.go:64` and `neon/task/untar.go:86` build the destination path with `filepath.Join(dir, header.Name)` without checking that the result stays inside `dir`. An archive containing an entry `../../.bashrc` writes outside the target directory.
**Fix:** reject entries that are absolute or whose `filepath.Rel(dir, target)` starts with `..` (or use `os.Root`, available since Go 1.24).

### 1.2 [Major] Update and version check over plain HTTP

`neon/build/update.go:18-20`: `versionURL` uses `http://`, and the install command suggested to the user is `sh -c "$(curl http://sweetohm.net/neon/install)"`. An attacker on the network (MITM) can get arbitrary code executed.
**Fix:** switch to HTTPS everywhere, ideally with checksum or signature verification.

### 1.3 [Minor] Singleton port open on all interfaces

`neon/build/build.go:503`: `net.Listen("tcp", ":port")` listens on `0.0.0.0`. A simple local lock should listen on `127.0.0.1`. Port `0` is also accepted: it yields a random port, so no real lock.

### 1.4 [Minor] HTTP calls without timeout

`neon/task/request.go:93` (`&http.Client{}`) and `neon/build/update.go:32` (`http.Get`) have no timeout. A server that doesn't respond blocks the build forever.

### 1.5 [Minor] CI: Codecov uploader via `bash <(curl …)`

`.travis.yml` runs Codecov's bash script, deprecated since the 2021 supply-chain incident. Use the official verified uploader (or the GitHub action).

## 2. Functional bugs

### 2.1 [Critical] False positive in cycle detection with diamond inheritance **[reproduced]**

`neon/build/build.go:60-77`: the `visited` map is never cleaned up when a branch finishes, so it acts as a global "already seen" set instead of a recursion stack. If A extends B and C, and B and C both extend D, the build is rejected with `dependency cycle detected`.
```yaml
# p.yml: extends: [./b.yml, ./c.yml]   b.yml and c.yml: extends: ./d.yml
```
This is a regression introduced by the "Cycle detection" fix in `doc/code-review.md`. Shared parent files (for instance `c4s4/build/common.yml`) are a common case.
**Fix:** `delete(visited, absPath)` once parents are resolved (DFS stack). The error message also lists the path in random order (map iteration): keep an ordered slice.
**Test to add:** a diamond case in `test/bugs/`.

### 2.2 [Critical] Panics on invalid build files **[reproduced]**

Several user inputs make neon panic instead of producing an error message:

| Input | Panic | Location |
|---|---|---|
| `- try: "foo"` (steps field given as a string, integer or boolean) | `reflect.Value.IsNil on string Value` | `neon/build/step.go:115` |
| `- 1: foo` (non-string task key) | `interface {} is int, not string` | `neon/build/step.go:89` |
| `- mkdir: [a, ~]` (`null` element in a list) | nil pointer dereference | `neon/build/task.go:204` (`reflect.TypeOf(nil).Kind()`) |
| property `L: [a, ~]` | `reflect.Value.Set on zero Value` | `neon/build/context.go:300` |
| empty map for a map-typed field (`env: {}`, `headers: {}`) | index out of range | `neon/build/task.go:214` (`MapKeys()[0]`) |
| `threads: -1` | `makechan: size out of range` | `neon/task/threads.go:77` |

Also note: `IsValueOfType` only checks the *first* entry of a map (`task.go:214-216`), and a `null` value in a property map is silently dropped (`SetMapIndex` with a zero value, `context.go:320`).

### 2.3 [Major] Target stack corrupted after an error caught by `try` **[reproduced]**

`neon/build/target.go:140-162`: the target is pushed (`Stack.Push`), then if a dependency fails or the `chdir` fails, the function returns without `Pop`. Same problem in `RunParentTarget` (`build.go:441-443`). `Stack.Push` also pushes the target when it detects a loop. When the error is caught by `try`, the next call of the same target reports a false infinite loop:
```yaml
targets:
  t: {steps: [{try: [{call: a}]}, {call: a}]}
  a: {depends: b}
  b: {steps: [{throw: boom}]}
# => ERROR ... infinite loop: t -> a -> a
```
**Fix:** `defer context.Stack.Pop()` right after a successful `Push`.

### 2.4 [Major] Inconsistent precedence of inherited targets **[reproduced]**

With `extends: [b, c]` where b and c both define target `x`:
- `GetTarget` (execution) picks **c**'s (last parent), `build.go:342`;
- `GetTargets` (used by `-info`, `-targets` and automatic target selection) keeps **b**'s (first parent), `build.go:293-298`;
- `GetProperties` gives precedence to the last parent.

So `neon -info` shows `b`'s doc while `neon x` runs `c`. **Fix:** iterate in the same direction as `GetTarget`.

### 2.5 [Major] The `version:` field is unusable with binaries built by `neon build` **[reproduced]**

`NeonVersion` is injected with the git hash (`VERSION` = `git rev-parse HEAD[0:8]` in `c4s4/build/golang.yml`), or is `UNKNOWN` without ldflags. The `greater`, `lower`, etc. builtins call `NewVersion(NeonVersion)`, which fails: any build file declaring `version: 'greater("0.1")'` crashes with `Neon version could not be parsed`. For the same reason, `neon -update` always announces a new version, since it compares a hash with a version number (`update.go:47`).
Other version issues: `Compare` considers `1.0` < `1.0.0` and ignores suffixes (`1.0-SNAPSHOT` == `1.0`). Version builtins report errors with `panic` (`neon/builtin/greater.go:32`, etc.).

### 2.6 [Major] `3>` redirection (stdout + stderr to a file) is corrupted **[reproduced]**

`neon/task/shell.go:108-129`: when `3>` is used, the same file is opened **twice** with `os.Create`, giving two descriptors with independent offsets. Writes overwrite each other:
```
echo OUT-LINE-1; echo ERR >&2; echo OUT-LINE-2   →  file: "ERR\nLINE-1\nOUT-LINE-2"
```
**Fix:** open the file once and share the writer. With `3=` (output to a property), the same `strings.Builder` is written concurrently by both `exec.Cmd` copy goroutines (through two different `MultiWriter`s), which is a **data race**. A mutex-protected writer is needed.

### 2.7 [Major] `1>>` / `2>>` / `3>>` fail if the file doesn't exist **[reproduced]**

`shell.go:134,146`: `os.OpenFile(..., O_WRONLY|O_APPEND)` without `O_CREATE`. The expected behaviour (that of `>>` in a shell) is to create the file.

### 2.8 [Major] `untar` doesn't truncate existing files **[reproduced]**

`untar.go:94`: `O_CREATE|O_RDWR` without `O_TRUNC`. Extracting `x` (5 bytes, "short") over a longer existing `x` leaves the tail of the old content: `shorth longer existing content`. Empty directories, symbolic links and mtimes are also ignored.

### 2.9 [Major] `zip` doesn't compress **[reproduced]**

`zip.go:104`: `zip.FileInfoHeader` returns `Method = Store`. The registered `BestCompression` compressor is therefore never used. A 200 KB file of zeros produces a 200 KB archive (`Stored 0%`). **Fix:** `header.Method = zip.Deflate`.

### 2.10 [Major] Ignored close errors: truncated archives with no error

`tar.go:81,86` and `zip.go:75`: `writer.Close()` / `gzip.Close()` / `zipper.Close()` write the archive trailers, but their errors are ignored (`_ =`). A full disk produces a corrupted archive and the build ends with `OK`. Same for `dest.Close()` in `unzip.go:78`.

### 2.11 [Major] Infinite loop if a `dotenv` file is a directory **[reproduced]**

`neon/build/context.go:423-427`: `reader.ReadLine()` only handles `io.EOF`. Any other error (for instance `EISDIR`) loops forever. `isPrefix` is also ignored: lines longer than 4 KB are split and the rest is parsed as a new line. Quotes (`FOO="bar"`) and `export FOO=bar` aren't handled, unlike the usual dotenv format.

### 2.12 [Major] Environment variables: undefined `${VAR}` wrongly replaced **[reproduced]**

`context.go:383`: `x${UNDEFINED}y` becomes `x{UNDEFINED}y`, the `$` is lost. More generally, in `EvaluateEnvironment`:
- an error in a `={...}` expression is silently ignored (`context.go:388-390`);
- `if replaced != ""` (`context.go:395`) makes it impossible to empty a variable (`FOO: ''`);
- escaping only handles one or two `\` (unlike `EvaluateString`);
- variables are evaluated in alphabetical order, so `A: ${B}` doesn't see `B` defined in the same block.

### 2.13 [Major] The build file's `repository:` field is never used

`neon/main.go:233-241` always passes a non-empty repository (`-repo` option, otherwise config, otherwise `~/.neon`), and `ParseRepository` (`neon/build/parse.go:115`) gives precedence to that parameter. The documented `repository:` field (`doc/usermanual.md:1025`) is therefore ignored.

### 2.14 [Major] The `grey` and `time` options from the configuration file are overwritten

`main.go:244-245`: `_build.Gray = opts.Grey` and `configuration.Time = opts.Time` unconditionally replace the values read from `settings.yml`. With `grey: true` in the configuration, output stays in color. An OR is needed (`opts.Grey || configuration.Grey`).

### 2.15 [Minor] `$` (list form) modifies the process `PATH` variable

`shell.go:272-280` does a global `os.Setenv("PATH", …)`, then restores it. Combined with the `threads` task, this creates a race condition: one thread can resolve its executable with another thread's `PATH`. The current directory is also process-wide (`os.Chdir` in `Target.Run` and in the `neon` task), while `threads` runs steps in parallel. **Fix:** look up the executable manually in the computed `PATH`, and set `Dir` explicitly instead of using `os.Getwd()`.

### 2.16 [Minor] `copy` of a file onto itself empties it

`neon/util/files.go:87`: `os.Create(dest)` truncates the source when `source == dest` (a glob that overlaps the destination is enough). Compare with `os.SameFile` before writing.

### 2.17 [Minor] `replace` resets permissions and the replacement order is random

`neon/task/replace.go:69` rewrites the file with `FileMode`: an executable script loses its `+x` bit. Iteration over the `with` map is random: if keys overlap (`a` and `ab`), the result changes from one run to the next. The error message also mentions "copy task".

### 2.18 [Minor] `threads` task

- `threads: 0` with input: nothing runs and the task returns without error.
- After an error, the other threads keep consuming all inputs, with no cancellation; only the first error is reported.
- The "Thread %d done" message `defer` is registered before `wg.Done`, so it can be printed after the task has finished.
- The `error` variable shadows the predeclared type.

### 2.19 [Minor] `neon` task

- `FIXME` in `neon/task/neon.go:43`: the path is resolved relative to the current directory, not the build directory.
- No recursion detection: a build that calls itself runs until memory is exhausted.
- Passed properties are set *after* `Init()`, so derived properties (`B: ={A}`) are not re-evaluated.

### 2.20 [Minor] Resolution of parent, script and template paths

- `RegexpParentName`, `RegexpTemplateName` and `RegexpScriptName` (`neon/build/repo.go:20-29`) are not anchored and the `.` is not escaped: `../common/base.yml` or `a/b/c/d.yml` are taken for repository paths. A relative path without `./` (`extends: base.yml`) isn't supported and yields the misleading error "parent was not found".
- `path.IsAbs` is used instead of `filepath.IsAbs` (`repo.go:76,175,193,209`): Windows absolute paths (`C:\…`) are not recognized.
- `SetDir` replaces the parents' `Dir` with the root build's, so a `context: ./lib.ank` or a relative `dotenv` declared in a parent is resolved relative to the child's directory.
- Parents' `version` and `singleton` fields are ignored: only the root build is checked in `Run`.
- `RegexpLinkName` is unused.

### 2.21 [Minor] Miscellaneous

- `ExpandUserHome` (`files.go:190`) ignores the error from `user.Current()`, which causes a nil pointer dereference (for instance in a container with a UID that has no entry in `/etc/passwd`). Use `os.UserHomeDir()`. A bare `~` is not expanded.
- `checkUnknownArgs` (`neon/build/task.go:152`) uses `continue` instead of `break`, and its condition `len(fields) != 0 || len(args) != 1` is obscure.
- If a step contains two task names, the task picked depends on `TaskMap` iteration order (`step.go:71`), so the error message is not deterministic.
- `MoveFilesToDir` uses `os.Rename`, which fails across filesystems (`EXDEV`).
- `FindFiles` ignores `zglob.Glob` errors, returns duplicates when globs overlap (creating duplicate entries in tar/zip), and fails on a broken symbolic link.
- `updateRepository` aborts the whole update as soon as a subfolder of the repository is not a git repository. In `printNewRelease`, the body is not closed when the status is not 200 (`update.go:36-41`).
- The `chdir` error message in template mode shows `Build.Dir` instead of `Build.Here` (`target.go:156`).
- Configuration is loaded *before* parsing the command line: an invalid `settings.yml` prevents even `neon -version` or `neon -help` from working.

## 3. User experience

- [Minor] Errors (`PrintError`) are written to **stdout**, not stderr (`neon/build/print.go:83-89`). `update.go` uses `println`, which writes to stderr: outputs are mixed inconsistently.
- [Minor] Some errors hide the cause: `SetCommandLineProperties` (`build.go:208`) and `ParseConfiguration` (`parse.go:194`) drop the YAML error, and `runList` replaces the `LookPath` error with "not found in PATH".
- [Minor] Build duration is always printed beyond 10 s, hard-coded in `main.go`, with no option to disable it.

---

## 4. Code quality and maintenance

- [Quality] Dependency on `gopkg.in/yaml.v2` (`yaml.v3` exists). With yaml.v2, the keys `yes`/`no`/`on`/`off` are read as booleans (YAML 1.1), which feeds the panics in §2.2.
- [Quality] `github.com/mcuadros/go-version` (2019, unmaintained) is only used by `sortversions.go`, while the project already has its own `Version` implementation. Two version comparison logics coexist.
- [Quality] Badly grouped imports (`strings` after third-party imports in `build.go:12`; `task.go`, `parse.go` and `repo.go` mix stdlib and module imports). `neon fmt` should run `goimports`.
- [Quality] `NewSteps` and `IsValueOfType` rely on `reflect` without safeguards. A typed walk (`switch v := x.(type)`) would be safer and more readable.
- [Quality] Dead or leftover code: `// (Empty or removed)` comment at the end of `main.go`, commented-out `TODO` in `context.go:212`, useless map loop in `GetShell` (`build.go:456`), `len(stderr) > 0` instead of `> 1` in `shell.go:179`, wrong `len(env)` capacity for `environ`, redundant `bytes.NewBuffer([]byte(body))` in `request.go`, `run_err` variable not in camelCase (`target.go:163`).
- [Quality] Outdated doc comments: `FindBuildFile` (3 return values and one parameter undocumented), `RunParentTarget` ("boolean: that tells if parent target was found" although no boolean is returned), "singleton port port", "parsig target".
- [Quality] Mutable global state (`Gray`, colors, `TaskMap`, `BuiltinMap`, `NeonVersion`) that makes tests hard to parallelize.
- [Quality] `GetProperties`, `GetEnvironment` and `GetDotEnv` recursively re-merge data already merged in `newBuildInternal`. Not wrong, but redundant, and `GetDotEnv` duplicates entries from shared parents (diamond).

## 5. Tooling, CI and documentation

- [Major] **Inconsistent Go versions:** `go.mod` declares `go 1.27`, while `.travis.yml` tests on Go `1.26` and `AGENTS.md` states "Go >= 1.26.5". CI can't build the module without a toolchain download. Align all three.
- [Minor] Travis CI (`travis-ci.org`) is no longer free for open source: CI is probably inactive. Consider migrating to GitHub Actions.
- [Minor] No test covers the cases in §2 (diamond, `try` + failing dependency, `3>`, `1>>`, untar over an existing file, zip compression ratio, `null` values). AGENTS.md recommends adding a file under `test/bugs/` for each bug.
- [Minor] The CHANGELOG doesn't mention inheritance cycle detection, nor the switch of `ParseCommandLine` to a struct (changes listed as done in `doc/code-review.md`).
- [Quality] `doc/code-review.md` is a dated review document living in `doc/`: it may be mistaken for user documentation. Merge it with this file or move it.
- [Quality] The locally installed `staticcheck` is too old for Go 1.27 ("export data version" errors). Not a repository problem, but `neon tools` could install it with a pinned version.

## Suggested priorities

1. §1.1 (Zip/Tar Slip) and §2.1 (diamond regression): security and visible regression.
2. §2.2 (panics) and §2.3 (corrupted stack): engine robustness.
3. §2.6 to §2.10: silent file corruption in the `$`, `untar`, `zip` and `tar` tasks.
4. §1.2, §2.5, §2.13, §2.14 and §5 (Go versions and CI).
