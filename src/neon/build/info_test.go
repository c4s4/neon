package build

import (
	"testing"
)

func TestInfoDoc(t *testing.T) {
	build := Build{
		Doc: "Test documentation",
	}
	if build.infoDoc() != "doc: Test documentation\n" {
		t.Errorf("Bad build doc: %s", build.infoDoc())
	}
}

func TestInfoDefault(t *testing.T) {
	build := Build{
		Default: []string{"default"},
	}
	if build.infoDefault() != "default: [default]\n" {
		t.Errorf("Bad build default: %s", build.infoDefault())
	}
}

func TestInfoRepository(t *testing.T) {
	build := Build{
		Repository: "repository",
	}
	if build.infoRepository() != "repository: repository\n" {
		t.Errorf("Bad build repository: %s", build.infoRepository())
	}
}

func TestInfoSingleton(t *testing.T) {
	build := &Build{
		Singleton: "12345",
	}
	context := NewContext(build)
	if build.infoSingleton(context) != "singleton: 12345\n" {
		t.Errorf("Bad build singleton: %s", build.infoSingleton(context))
	}
}

func TestInfoExtends(t *testing.T) {
	build := &Build{
		Extends: []string{"foo", "bar"},
	}
	if build.infoExtends() != "extends:\n- foo\n- bar\n" {
		t.Errorf("Bad build extends: %s", build.infoExtends())
	}
}

func TestInfoConfiguration(t *testing.T) {
	build := &Build{
		Config: []string{"foo", "bar"},
	}
	if build.infoConfiguration() != "configuration:\n- foo\n- bar\n" {
		t.Errorf("Bad build config: %s", build.infoConfiguration())
	}
}

func TestInfoContext(t *testing.T) {
	build := &Build{
		Scripts: []string{"foo", "bar"},
	}
	if build.infoContext() != "context:\n- foo\n- bar\n" {
		t.Errorf("Bad build context: %s", build.infoContext())
	}
}

func TestInfoBuiltins(t *testing.T) {
	BuiltinMap = make(map[string]BuiltinDesc)
	AddBuiltin(BuiltinDesc{
		Name: "test",
		Func: TestInfoBuiltins,
		Help: `Test documentation.`,
	})
	builtins := InfoBuiltins()
	if builtins != "test" {
		t.Errorf("Bad builtins: %s", builtins)
	}
}

func TestInfoBuiltin(t *testing.T) {
	BuiltinMap = make(map[string]BuiltinDesc)
	AddBuiltin(BuiltinDesc{
		Name: "test",
		Func: TestInfoBuiltins,
		Help: `Test documentation.`,
	})
	info := InfoBuiltin("test")
	if info != "Test documentation." {
		t.Errorf("Bad builtin info: %s", info)
	}
}

func TestInfoThemes(t *testing.T) {
	themes := InfoThemes()
	if themes != "bee blue bold cyan fire green magenta marine nature red reverse rgb yellow" {
		t.Errorf("Bad themes")
	}
}
