package build

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

// writeBuilds writes build files in a temporary directory and returns it
func writeBuilds(t *testing.T, files map[string]string) string {
	dir := t.TempDir()
	for name, content := range files {
		if _, err := WriteFile(dir, name, content); err != nil {
			t.Fatalf("writing build file: %v", err)
		}
	}
	return dir
}

func TestDiamondInheritance(t *testing.T) {
	dir := writeBuilds(t, map[string]string{
		"build.yml": "extends: [./b.yml, ./c.yml]\n",
		"b.yml":     "extends: ./d.yml\n",
		"c.yml":     "extends: ./d.yml\n",
		"d.yml":     "targets:\n  x:\n    steps: []\n",
	})
	build, err := NewBuild(filepath.Join(dir, "build.yml"), dir, "", false)
	if err != nil {
		t.Fatalf("diamond inheritance should be accepted: %v", err)
	}
	if build.GetTarget("x") == nil {
		t.Errorf("target x should be inherited")
	}
}

func TestInheritanceCycle(t *testing.T) {
	dir := writeBuilds(t, map[string]string{
		"a.yml": "extends: ./b.yml\n",
		"b.yml": "extends: ./c.yml\n",
		"c.yml": "extends: ./a.yml\n",
	})
	_, err := NewBuild(filepath.Join(dir, "a.yml"), dir, "", false)
	if err == nil || !strings.Contains(err.Error(), "dependency cycle detected") {
		t.Fatalf("cycle should be detected, got: %v", err)
	}
	// cycle is printed in order
	expected := strings.Join([]string{
		filepath.Join(dir, "a.yml"), filepath.Join(dir, "b.yml"),
		filepath.Join(dir, "c.yml"), filepath.Join(dir, "a.yml")}, " → ")
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("bad cycle description: %v", err)
	}
}

func TestInheritedTargetsPrecedence(t *testing.T) {
	dir := writeBuilds(t, map[string]string{
		"build.yml": "extends: [./b.yml, ./c.yml]\n",
		"b.yml":     "targets:\n  x:\n    doc: b\n",
		"c.yml":     "targets:\n  x:\n    doc: c\n",
	})
	build, err := NewBuild(filepath.Join(dir, "build.yml"), dir, "", false)
	if err != nil {
		t.Fatalf("loading build: %v", err)
	}
	if build.GetTarget("x").Doc != "c" {
		t.Errorf("GetTarget should return last parent target")
	}
	if build.GetTargets()["x"].Doc != "c" {
		t.Errorf("GetTargets should be consistent with GetTarget")
	}
}

func TestStackAfterCaughtError(t *testing.T) {
	// running targets changes current directory, restore it at the end
	t.Chdir(".")
	dir := writeBuilds(t, map[string]string{
		"build.yml": `targets:
  a:
    depends: b
  b:
    steps:
    - 'throw("boom")'
`,
	})
	build, err := NewBuild(filepath.Join(dir, "build.yml"), dir, "", false)
	if err != nil {
		t.Fatalf("loading build: %v", err)
	}
	context := NewContext(build)
	if err := context.Init(); err != nil {
		t.Fatalf("init context: %v", err)
	}
	if err := build.RunTarget(context, "a"); err == nil {
		t.Fatalf("target should fail")
	}
	if len(context.Stack.Targets) != 0 {
		t.Errorf("stack should be empty after error: %v", context.Stack)
	}
	if err := build.RunParentTarget(context, "a"); err == nil {
		t.Fatalf("missing parent target should fail")
	}
	if len(context.Stack.Targets) != 0 {
		t.Errorf("stack should be empty after error: %v", context.Stack)
	}
}

func TestInvalidStepsDontPanic(t *testing.T) {
	for _, steps := range []interface{}{"foo", 42, true, map[interface{}]interface{}{}} {
		if _, err := NewSteps(steps); err == nil {
			t.Errorf("steps %v should be rejected", steps)
		}
	}
	if steps, err := NewSteps(nil); err != nil || len(steps) != 0 {
		t.Errorf("nil steps should give empty steps")
	}
	if _, err := NewTaskStep(TaskArgs{1: "foo"}); err == nil || !strings.Contains(err.Error(), "unknown task '1'") {
		t.Errorf("non string task key should be rejected: %v", err)
	}
}

func TestIsValueOfTypeNilAndMaps(t *testing.T) {
	stringSlice := []string{}
	stringMap := map[string]string{}
	if IsValueOfType([]interface{}{"a", nil}, reflect.TypeOf(stringSlice)) {
		t.Errorf("nil element should not be a string")
	}
	if !IsValueOfType(map[interface{}]interface{}{}, reflect.TypeOf(stringMap)) {
		t.Errorf("empty map should be a map of strings")
	}
	if IsValueOfType(map[interface{}]interface{}{"a": "b", "c": 1}, reflect.TypeOf(stringMap)) {
		t.Errorf("all map values should be checked")
	}
	if !IsValueOfType(map[interface{}]interface{}{"a": "b", "c": "d"}, reflect.TypeOf(stringMap)) {
		t.Errorf("map of strings should be accepted")
	}
}

func TestEvaluateObjectWithNil(t *testing.T) {
	context := NewContext(&Build{})
	value, err := context.EvaluateObject([]interface{}{"a", nil})
	if err != nil {
		t.Fatalf("evaluating list with nil: %v", err)
	}
	list := value.([]interface{})
	if len(list) != 2 || list[0] != "a" || list[1] != nil {
		t.Errorf("bad evaluated list: %v", list)
	}
	value, err = context.EvaluateObject(map[interface{}]interface{}{"a": nil})
	if err != nil {
		t.Fatalf("evaluating map with nil: %v", err)
	}
	m := value.(map[interface{}]interface{})
	if v, ok := m["a"]; !ok || v != nil {
		t.Errorf("nil value should be kept in map: %v", m)
	}
}

func TestEvaluateEnvironment(t *testing.T) {
	t.Setenv("NEON_TEST_EXISTING", "existing")
	t.Setenv("EMPTY", "should be unset")
	build := &Build{
		Environment: map[string]string{
			"A_REF":     "a${Z_DEF}",
			"EMPTY":     "",
			"UNDEFINED": "x${NEON_TEST_UNDEFINED}y",
			"ESCAPED":   `\${NEON_TEST_EXISTING}`,
			"EXISTING":  "${NEON_TEST_EXISTING}",
			"Z_DEF":     "z",
		},
	}
	context := NewContext(build)
	environment, err := context.EvaluateEnvironment()
	if err != nil {
		t.Fatalf("evaluating environment: %v", err)
	}
	expected := map[string]string{
		"A_REF":     "az",
		"UNDEFINED": "x${NEON_TEST_UNDEFINED}y",
		"ESCAPED":   "${NEON_TEST_EXISTING}",
		"EXISTING":  "existing",
	}
	for name, value := range expected {
		if actual, ok := environment[name]; !ok || actual != value {
			t.Errorf("bad value for %s: '%s' (expected '%s')", name, actual, value)
		}
	}
	if _, ok := environment["EMPTY"]; ok {
		t.Errorf("variable set to empty string should be unset")
	}
	// errors in expressions are reported
	build.Environment = map[string]string{"BAD": "={undefined_function()}"}
	if _, err := context.EvaluateEnvironment(); err == nil {
		t.Errorf("error in expression should be reported")
	}
	// cycles are detected
	build.Environment = map[string]string{"A": "${B}", "B": "${A}"}
	if _, err := context.EvaluateEnvironment(); err == nil || !strings.Contains(err.Error(), "cycle") {
		t.Errorf("cycle should be detected: %v", err)
	}
}

func TestReadDotEnv(t *testing.T) {
	long := strings.Repeat("x", 10000)
	content := "\n# comment\nFOO=bar\nexport EXPORTED=yes\nDOUBLE=\"double quoted\"\nSINGLE='single quoted'\n\nLONG=" + long + "\n"
	environment, err := ReadDotEnv(strings.NewReader(content))
	if err != nil {
		t.Fatalf("reading dotenv: %v", err)
	}
	expected := map[string]string{
		"FOO":      "bar",
		"EXPORTED": "yes",
		"DOUBLE":   "double quoted",
		"SINGLE":   "single quoted",
		"LONG":     long,
	}
	Assert(environment, expected, t)
}

func TestLoadDotEnvDirectory(t *testing.T) {
	dir := t.TempDir()
	if _, err := LoadDotEnv(dir, ""); err == nil {
		t.Errorf("loading a directory should fail")
	}
}

func TestParentPathRelative(t *testing.T) {
	dir := writeBuilds(t, map[string]string{
		"base.yml": "",
	})
	build := &Build{Dir: dir, Repository: "/tmp/repo"}
	path, err := build.ParentPath("base.yml")
	if err != nil || path != filepath.Join(dir, "base.yml") {
		t.Errorf("relative parent without ./ should be found: %s (%v)", path, err)
	}
	path, err = build.ParentPath("../common/base.yml")
	if err != nil || path != filepath.Join(dir, "../common/base.yml") {
		t.Errorf("../ path should be relative to build: %s (%v)", path, err)
	}
	path, err = build.ParentPath("c4s4/build/golang.yml")
	if err != nil || path != "/tmp/repo/c4s4/build/golang.yml" {
		t.Errorf("repository path: %s (%v)", path, err)
	}
	for _, name := range []string{"a/b/c/d.yml", "a/b/cyml"} {
		if RegexpParentName.MatchString(name) {
			t.Errorf("'%s' should not match parent name regexp", name)
		}
	}
}

func TestParentScriptPath(t *testing.T) {
	dir := t.TempDir()
	parentDir := filepath.Join(dir, "parent")
	if _, err := WriteFile(parentDir, "parent.yml", "context: ./lib.ank\n"); err != nil {
		t.Fatal(err)
	}
	if _, err := WriteFile(parentDir, "lib.ank", "func libfunc() { return 42 }\n"); err != nil {
		t.Fatal(err)
	}
	childDir := filepath.Join(dir, "child")
	child, err := WriteFile(childDir, "build.yml", "extends: ../parent/parent.yml\n")
	if err != nil {
		t.Fatal(err)
	}
	build, err := NewBuild(child, childDir, "", false)
	if err != nil {
		t.Fatalf("loading build: %v", err)
	}
	scripts, err := build.GetScripts()
	if err != nil {
		t.Fatalf("getting scripts: %v", err)
	}
	Assert(scripts, []string{filepath.Join(parentDir, "lib.ank")}, t)
}

func TestCheckVersion(t *testing.T) {
	saved := NeonVersion
	defer func() {
		NeonVersion = saved
	}()
	dir := writeBuilds(t, map[string]string{
		"build.yml":  "extends: ./parent.yml\n",
		"parent.yml": "version: '1 == 2'\n",
	})
	build, err := NewBuild(filepath.Join(dir, "build.yml"), dir, "", false)
	if err != nil {
		t.Fatalf("loading build: %v", err)
	}
	context := NewContext(build)
	// development version: check is skipped
	NeonVersion = "15fb03c"
	if err := build.CheckVersion(context); err != nil {
		t.Errorf("version check should be skipped for development version: %v", err)
	}
	// release version: parent version is checked
	NeonVersion = "1.2.3"
	if err := build.CheckVersion(context); err == nil {
		t.Errorf("parent version should be checked")
	}
}

func TestRepositoryField(t *testing.T) {
	saved := UserRepository
	defer func() {
		UserRepository = saved
	}()
	UserRepository = "/user/repo"
	dir := writeBuilds(t, map[string]string{
		"field.yml":    "repository: /build/repo\n",
		"nofield.yml":  "doc: test\n",
		"override.yml": "repository: /build/repo\n",
	})
	build, err := NewBuild(filepath.Join(dir, "field.yml"), dir, "", false)
	if err != nil || build.Repository != "/build/repo" {
		t.Errorf("repository field should be used: %v", err)
	}
	build, err = NewBuild(filepath.Join(dir, "nofield.yml"), dir, "", false)
	if err != nil || build.Repository != "/user/repo" {
		t.Errorf("user repository should be used: %v", err)
	}
	build, err = NewBuild(filepath.Join(dir, "override.yml"), dir, "/cli/repo", false)
	if err != nil || build.Repository != "/cli/repo" {
		t.Errorf("command line repository should override: %v", err)
	}
}

func TestSingletonInherited(t *testing.T) {
	build := &Build{Parents: []*Build{{Singleton: "1234"}}}
	if build.GetSingleton() != "1234" {
		t.Errorf("singleton should be inherited")
	}
}

func TestDotEnvPaths(t *testing.T) {
	dir := t.TempDir()
	parentDir := filepath.Join(dir, "parent")
	if _, err := WriteFile(parentDir, "parent.yml", "dotenv: parent.env\n"); err != nil {
		t.Fatal(err)
	}
	childDir := filepath.Join(dir, "child")
	child, err := WriteFile(childDir, "build.yml", "extends: ../parent/parent.yml\ndotenv: child.env\n")
	if err != nil {
		t.Fatal(err)
	}
	build, err := NewBuild(child, childDir, "", false)
	if err != nil {
		t.Fatalf("loading build: %v", err)
	}
	expected := []string{filepath.Join(parentDir, "parent.env"), filepath.Join(childDir, "child.env")}
	Assert(build.DotEnv, expected, t)
}

func TestCompareNeonVersion(t *testing.T) {
	saved := NeonVersion
	defer func() {
		NeonVersion = saved
	}()
	NeonVersion = "1.2.3"
	if CompareNeonVersion("1.2") <= 0 || CompareNeonVersion("1.2.3") != 0 || CompareNeonVersion("1.3") >= 0 {
		t.Errorf("bad comparison with NeON version")
	}
	NeonVersion = "15fb03c"
	defer func() {
		if r := recover(); r == nil || !strings.Contains(r.(string), "development version") {
			t.Errorf("development version should raise a clear error: %v", r)
		}
	}()
	CompareNeonVersion("1.0")
}

func TestDiamondInheritanceNoDuplicates(t *testing.T) {
	dir := writeBuilds(t, map[string]string{
		"build.yml": "extends: [./b.yml, ./c.yml]\n",
		"b.yml":     "extends: ./d.yml\nproperties:\n  B: b\n",
		"c.yml":     "extends: ./d.yml\n",
		"d.yml":     "dotenv: d.env\ncontext: ./d.ank\nproperties:\n  D: d\n",
		"d.ank":     "",
	})
	build, err := NewBuild(filepath.Join(dir, "build.yml"), dir, "", false)
	if err != nil {
		t.Fatalf("loading build: %v", err)
	}
	Assert(build.DotEnv, []string{filepath.Join(dir, "d.env")}, t)
	scripts, err := build.GetScripts()
	if err != nil {
		t.Fatal(err)
	}
	Assert(scripts, []string{filepath.Join(dir, "d.ank")}, t)
	if build.Properties["B"] != "b" || build.Properties["D"] != "d" {
		t.Errorf("properties should be inherited: %v", build.Properties)
	}
}

func TestScalarsToStrings(t *testing.T) {
	Assert(ScalarsToStrings(1), "1", t)
	Assert(ScalarsToStrings(int64(-7)), "-7", t)
	Assert(ScalarsToStrings(2.5), "2.5", t)
	Assert(ScalarsToStrings(1.0), "1", t)
	Assert(ScalarsToStrings(1e6), "1000000", t)
	Assert(ScalarsToStrings("foo"), "foo", t)
	Assert(ScalarsToStrings(true), "true", t)
	Assert(ScalarsToStrings(time.Date(2026, 11, 1, 12, 30, 0, 0, time.UTC)), "2026-11-01T12:30:00Z", t)
	Assert(ScalarsToStrings(nil), nil, t)
	Assert(ScalarsToStrings([]byte("bytes")), "bytes", t)
	Assert(ScalarsToStrings([]interface{}{[]byte("b")}), []interface{}{"b"}, t)
	Assert(FlattenLists([]interface{}{"a", []interface{}{"b", []string{"c"}}}), []interface{}{"a", "b", "c"}, t)
	Assert(FlattenLists([]string{"a"}), []string{"a"}, t)
	Assert(FlattenLists([]byte("ab")), []byte("ab"), t)
	Assert(ScalarsToStrings([]interface{}{"a", 1, 2.5, false}), []interface{}{"a", "1", "2.5", "false"}, t)
	Assert(ScalarsToStrings([]string{"a"}), []string{"a"}, t)
	Assert(ScalarsToStrings(map[interface{}]interface{}{"A": 1, "B": "b"}),
		map[interface{}]interface{}{"A": "1", "B": "b"}, t)
}

func TestNumbersInStringFields(t *testing.T) {
	type args struct {
		List   []string          `neon:"wrap,strings"`
		Text   string            `neon:"optional,strings"`
		Env    map[string]string `neon:"optional,strings"`
		Strict []string          `neon:"optional,wrap"`
	}
	typ := reflect.TypeOf(args{})
	valid := TaskArgs{"list": []interface{}{"echo", 1, 2.5}, "text": 12, "env": map[interface{}]interface{}{"N": 42}}
	if err := ValidateTaskArgs(valid, typ); err != nil {
		t.Fatalf("numbers should be accepted in string fields: %v", err)
	}
	context := NewContext(&Build{})
	result, err := EvaluateTaskArgs(valid, typ, context)
	if err != nil {
		t.Fatalf("evaluating args: %v", err)
	}
	evaluated := result.(args)
	Assert(evaluated.List, []string{"echo", "1", "2.5"}, t)
	Assert(evaluated.Text, "12", t)
	Assert(evaluated.Env, map[string]string{"N": "42"}, t)
	// wrapped single number
	result, err = EvaluateTaskArgs(TaskArgs{"list": 7}, typ, context)
	if err != nil {
		t.Fatalf("evaluating args: %v", err)
	}
	Assert(result.(args).List, []string{"7"}, t)
	// booleans are converted too
	result, err = EvaluateTaskArgs(TaskArgs{"list": []interface{}{"echo", true}}, typ, context)
	if err != nil {
		t.Fatalf("evaluating args: %v", err)
	}
	Assert(result.(args).List, []string{"echo", "true"}, t)
	// nested lists are flattened
	nested := TaskArgs{"list": []interface{}{"echo", []interface{}{1, []interface{}{"a", true}}, "=[2, [3]]"}}
	if err := ValidateTaskArgs(nested, typ); err != nil {
		t.Fatalf("nested lists should be accepted in list fields: %v", err)
	}
	result, err = EvaluateTaskArgs(nested, typ, context)
	if err != nil {
		t.Fatalf("evaluating args: %v", err)
	}
	Assert(result.(args).List, []string{"echo", "1", "a", "true", "2", "3"}, t)
	// list expression for the whole field is flattened too
	result, err = EvaluateTaskArgs(TaskArgs{"list": "=[1, [2, 3]]"}, typ, context)
	if err != nil {
		t.Fatalf("evaluating args: %v", err)
	}
	Assert(result.(args).List, []string{"1", "2", "3"}, t)
	// maps in lists are rejected with a clear error
	if err := ValidateTaskArgs(TaskArgs{"list": []interface{}{"echo", map[interface{}]interface{}{"a": 1}}}, typ); err == nil {
		t.Errorf("maps should not be accepted in list of strings")
	}
	_, err = EvaluateTaskArgs(TaskArgs{"list": []interface{}{"echo", "= {\"a\": 1}"}}, typ, context)
	if err == nil || !strings.Contains(err.Error(), "must be of type '[]string'") {
		t.Errorf("map returned by expression should give a clear error: %v", err)
	}
	// lists are not flattened in string fields
	if err := ValidateTaskArgs(TaskArgs{"list": "echo", "text": []interface{}{"a"}}, typ); err == nil {
		t.Errorf("list should not be accepted in string field")
	}
	// fields without strings quality still reject numbers
	if err := ValidateTaskArgs(TaskArgs{"list": "echo", "strict": []interface{}{1}}, typ); err == nil {
		t.Errorf("numbers should be rejected in fields without strings quality")
	}
}

func TestPrintNewRelease(t *testing.T) {
	savedURL, savedVersion := versionURL, NeonVersion
	defer func() {
		versionURL, NeonVersion = savedURL, savedVersion
	}()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("v2.0.0\n"))
	}))
	defer server.Close()
	versionURL = server.URL
	cases := map[string]string{
		"v1.17.0": "Neon v2.0.0 is available",
		"2.0.0":   "Your version of neon is the latest one",
		"v2.0.0":  "Your version of neon is the latest one",
		"v2.1.0":  "Your version of neon is the latest one",
		"15fb03c": "development version of neon (15fb03c), latest release is v2.0.0",
	}
	for current, expected := range cases {
		NeonVersion = current
		stdout := os.Stdout
		read, write, _ := os.Pipe()
		os.Stdout = write
		printNewRelease()
		os.Stdout = stdout
		_ = write.Close()
		out, _ := io.ReadAll(read)
		if !strings.Contains(string(out), expected) {
			t.Errorf("with version %s, expected '%s' in output: %q", current, expected, out)
		}
	}
}
