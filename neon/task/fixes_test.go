package task

import (
	t "archive/tar"
	z "archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_build "github.com/c4s4/neon/neon/build"
)

// newTestContext makes a context for a build in given directory
func newTestContext(dir string) *_build.Context {
	build := &_build.Build{
		Dir:   dir,
		Shell: map[string][]string{"default": {"sh", "-c"}},
	}
	return _build.NewContext(build)
}

func TestShellRedirectBothToFile(test *testing.T) {
	dir := test.TempDir()
	file := filepath.Join(dir, "both.txt")
	context := newTestContext(dir)
	args := shellArgs{
		Shell: []string{"echo OUT-1; sleep 0.1; echo ERR >&2; sleep 0.1; echo OUT-2"},
		Red3:  file,
		Del3:  true,
	}
	if err := shell(context, args); err != nil {
		test.Fatalf("running command: %v", err)
	}
	content, err := os.ReadFile(file)
	if err != nil {
		test.Fatal(err)
	}
	if string(content) != "OUT-1\nERR\nOUT-2\n" {
		test.Errorf("bad file content: %q", content)
	}
}

func TestShellAppendCreatesFile(test *testing.T) {
	dir := test.TempDir()
	file := filepath.Join(dir, "new.txt")
	context := newTestContext(dir)
	for i := 0; i < 2; i++ {
		if err := shell(context, shellArgs{Shell: []string{"echo hi"}, App1: file, Del1: true}); err != nil {
			test.Fatalf("running command: %v", err)
		}
	}
	content, err := os.ReadFile(file)
	if err != nil {
		test.Fatal(err)
	}
	if string(content) != "hi\nhi\n" {
		test.Errorf("bad file content: %q", content)
	}
}

func TestShellBothToProperty(test *testing.T) {
	// run with -race to check concurrent writes
	context := newTestContext(test.TempDir())
	args := shellArgs{
		Shell: []string{"for i in 1 2 3 4 5; do echo out; echo err >&2; done"},
		Var3:  "output",
		Del3:  true,
	}
	if err := shell(context, args); err != nil {
		test.Fatalf("running command: %v", err)
	}
	output, err := context.GetProperty("output")
	if err != nil {
		test.Fatal(err)
	}
	if strings.Count(output.(string), "out") != 5 || strings.Count(output.(string), "err") != 5 {
		test.Errorf("bad output: %q", output)
	}
}

func TestShellListWithEnvironmentPath(test *testing.T) {
	dir := test.TempDir()
	script := filepath.Join(dir, "neon-test-command")
	if err := os.WriteFile(script, []byte("#!/bin/sh\necho found\n"), 0755); err != nil {
		test.Fatal(err)
	}
	context := newTestContext(dir)
	context.Build.Environment = map[string]string{"PATH": dir + string(os.PathListSeparator) + "${PATH}"}
	path := os.Getenv("PATH")
	args := shellArgs{Shell: []string{"neon-test-command", "arg"}, Var1: "output", Del1: true}
	if err := shell(context, args); err != nil {
		test.Fatalf("running command: %v", err)
	}
	if os.Getenv("PATH") != path {
		test.Errorf("process PATH should not be changed")
	}
	if output, _ := context.GetProperty("output"); output != "found" {
		test.Errorf("bad output: %v", output)
	}
}

func TestUntarOverwritesAndCreatesDirectories(test *testing.T) {
	tmp := test.TempDir()
	archive := filepath.Join(tmp, "test.tar")
	file, err := os.Create(archive)
	if err != nil {
		test.Fatal(err)
	}
	writer := t.NewWriter(file)
	if err := writer.WriteHeader(&t.Header{Name: "empty/", Mode: 0755, Typeflag: t.TypeDir}); err != nil {
		test.Fatal(err)
	}
	if err := writer.WriteHeader(&t.Header{Name: "x", Mode: 0644, Size: 5, Typeflag: t.TypeReg}); err != nil {
		test.Fatal(err)
	}
	if _, err := writer.Write([]byte("short")); err != nil {
		test.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		test.Fatal(err)
	}
	if err := file.Close(); err != nil {
		test.Fatal(err)
	}
	dir := filepath.Join(tmp, "out")
	if err := os.MkdirAll(dir, 0755); err != nil {
		test.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "x"), []byte("a much longer existing content"), 0644); err != nil {
		test.Fatal(err)
	}
	if err := untarFile(archive, dir); err != nil {
		test.Fatalf("untar: %v", err)
	}
	content, err := os.ReadFile(filepath.Join(dir, "x"))
	if err != nil || string(content) != "short" {
		test.Errorf("existing file should be overwritten: %q (%v)", content, err)
	}
	if stat, err := os.Stat(filepath.Join(dir, "empty")); err != nil || !stat.IsDir() {
		test.Errorf("empty directory should be created")
	}
}

func TestZipCompresses(test *testing.T) {
	dir := test.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "zeros.txt"), make([]byte, 200000), 0644); err != nil {
		test.Fatal(err)
	}
	archive := filepath.Join(dir, "test.zip")
	if err := writeZip(dir, []string{"zeros.txt"}, "", archive); err != nil {
		test.Fatalf("zip: %v", err)
	}
	reader, err := z.OpenReader(archive)
	if err != nil {
		test.Fatal(err)
	}
	defer func() {
		_ = reader.Close()
	}()
	entry := reader.File[0]
	if entry.Method != z.Deflate || entry.CompressedSize64 >= entry.UncompressedSize64/10 {
		test.Errorf("file should be compressed (method %d, %d -> %d bytes)",
			entry.Method, entry.UncompressedSize64, entry.CompressedSize64)
	}
}

func TestArchiveWriteErrors(test *testing.T) {
	dir := test.TempDir()
	missing := filepath.Join(dir, "missing", "archive")
	if err := writeZip(dir, []string{}, "", missing+".zip"); err == nil {
		test.Errorf("zip in missing directory should fail")
	}
	if err := writeTar(dir, []string{}, "", missing+".tar.gz"); err == nil {
		test.Errorf("tar in missing directory should fail")
	}
	// tar and gzip trailers are written
	archive := filepath.Join(dir, "archive.tar.gz")
	if err := os.WriteFile(filepath.Join(dir, "file.txt"), []byte("content"), 0644); err != nil {
		test.Fatal(err)
	}
	if err := writeTar(dir, []string{"file.txt"}, "", archive); err != nil {
		test.Fatalf("tar: %v", err)
	}
	out := filepath.Join(dir, "out")
	if err := untarFile(archive, out); err != nil {
		test.Fatalf("archive should be valid: %v", err)
	}
}

func TestThreadsArguments(test *testing.T) {
	context := newTestContext(test.TempDir())
	for _, count := range []int{0, -1} {
		err := threads(context, threadsArgs{Threads: count, Input: []interface{}{1}})
		if err == nil || !strings.Contains(err.Error(), "at least 1") {
			test.Errorf("threads %d should be rejected: %v", count, err)
		}
	}
}

func TestThreadsStopOnError(test *testing.T) {
	context := newTestContext(test.TempDir())
	steps, err := _build.NewSteps([]interface{}{
		"if _input == 1 { throw(\"boom\") }",
		"_output = _input",
	})
	if err != nil {
		test.Fatal(err)
	}
	var input []interface{}
	for i := 1; i <= 100; i++ {
		input = append(input, i)
	}
	err = threads(context, threadsArgs{Threads: 1, Input: input, Steps: steps})
	if err == nil {
		test.Fatalf("threads should fail")
	}
	output, _ := context.GetProperty(propertyOutput)
	if list, ok := output.([]interface{}); ok && len(list) > 0 {
		test.Errorf("thread should stop on first error, got output %v", list)
	}
}

func TestReplaceKeepsModeAndIsDeterministic(test *testing.T) {
	dir := test.TempDir()
	file := filepath.Join(dir, "script.sh")
	if err := os.WriteFile(file, []byte("ab"), 0755); err != nil {
		test.Fatal(err)
	}
	context := newTestContext(dir)
	with := map[interface{}]interface{}{"a": "X", "ab": "Y"}
	if err := replace(context, replaceArgs{Replace: []string{"script.sh"}, With: with, Dir: dir}); err != nil {
		test.Fatalf("replace: %v", err)
	}
	content, _ := os.ReadFile(file)
	if string(content) != "Y" {
		test.Errorf("longest string should be replaced first: %q", content)
	}
	if stat, err := os.Stat(file); err != nil || stat.Mode().Perm() != 0755 {
		test.Errorf("file mode should be kept")
	}
}

func TestSingletonInvalidPortDoesntWait(test *testing.T) {
	context := newTestContext(test.TempDir())
	if err := singleton(context, singletonArgs{Singleton: 0, Wait: true}); err == nil {
		test.Errorf("invalid port should fail immediately")
	}
}

func TestNeonRecursion(test *testing.T) {
	dir := test.TempDir()
	content := "default: t\ntargets:\n  t:\n    steps:\n    - neon: build.yml\n"
	file := filepath.Join(dir, "build.yml")
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		test.Fatal(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		test.Fatal(err)
	}
	defer func() {
		_ = os.Chdir(cwd)
	}()
	build, err := _build.NewBuild(file, dir, "", false)
	if err != nil {
		test.Fatal(err)
	}
	context := _build.NewContext(build)
	if err := context.Init(); err != nil {
		test.Fatal(err)
	}
	err = build.Run(context, nil)
	if err == nil || !strings.Contains(err.Error(), "recursive call") {
		test.Errorf("recursive call should be detected: %v", err)
	}
}

func TestNeonPropertiesOverride(test *testing.T) {
	dir := test.TempDir()
	content := `properties:
  A: a
  B: '={A}-b'
targets:
  t:
    steps:
    - if: 'B != "x-b"'
      then:
      - throw: 'bad B: ={B}'
`
	if err := os.WriteFile(filepath.Join(dir, "child.yml"), []byte(content), 0644); err != nil {
		test.Fatal(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		test.Fatal(err)
	}
	defer func() {
		_ = os.Chdir(cwd)
	}()
	context := newTestContext(dir)
	args := neonArgs{
		Neon:       "child.yml",
		Targets:    []string{"t"},
		Properties: map[interface{}]interface{}{"A": "x"},
	}
	if err := neon(context, args); err != nil {
		test.Errorf("derived properties should use overridden values: %v", err)
	}
}

func TestThreadsChdir(test *testing.T) {
	test.Chdir(test.TempDir())
	dir, err := os.Getwd()
	if err != nil {
		test.Fatal(err)
	}
	other := test.TempDir()
	context := newTestContext(dir)
	// changing to another directory in threads fails
	steps, err := _build.NewSteps([]interface{}{map[interface{}]interface{}{"chdir": other}})
	if err != nil {
		test.Fatal(err)
	}
	err = threads(context, threadsArgs{Threads: 2, Input: []interface{}{1, 2}, Steps: steps})
	if err == nil || !strings.Contains(err.Error(), "while threads are running") {
		test.Errorf("changing directory in threads should fail: %v", err)
	}
	// changing to current directory is accepted
	steps, err = _build.NewSteps([]interface{}{map[interface{}]interface{}{"chdir": dir}})
	if err != nil {
		test.Fatal(err)
	}
	if err := threads(context, threadsArgs{Threads: 2, Input: []interface{}{1, 2}, Steps: steps}); err != nil {
		test.Errorf("changing to current directory in threads should work: %v", err)
	}
	// after threads, changing directory works again
	if err := _build.Chdir(other); err != nil {
		test.Errorf("changing directory after threads should work: %v", err)
	}
}

func TestNeonUsesOwnRepository(test *testing.T) {
	dir := test.TempDir()
	content := "repository: /child/repo\ntargets:\n  t:\n    steps:\n    - if: '_REPO != \"/child/repo\"'\n      then:\n      - throw: 'bad repository: ={_REPO}'\n"
	if err := os.WriteFile(filepath.Join(dir, "child.yml"), []byte(content), 0644); err != nil {
		test.Fatal(err)
	}
	test.Chdir(dir)
	context := newTestContext(dir)
	context.Build.Repository = "/parent/repo"
	if err := neon(context, neonArgs{Neon: "child.yml", Targets: []string{"t"}}); err != nil {
		test.Errorf("called build should use its own repository: %v", err)
	}
}
