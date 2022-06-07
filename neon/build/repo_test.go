package build

import (
	"github.com/c4s4/neon/neon/util"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestPluginPath(t *testing.T) {
	build := &Build{
		Repository: "~/.neon",
	}
	path, err := build.ParentPath("foo/bar/spam.yml")
	Assert(err, nil, t)
	Assert(path, util.ExpandUserHome("~/.neon/foo/bar/spam.yml"), t)
}

func TestFindParents(t *testing.T) {
	repo := "/tmp/neon"
	if _, err := WriteFile(repo+"/foo/bar", "parent1.yml", ""); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if _, err := WriteFile(repo+"/foo/bar", "parent2.yml", ""); err != nil {
		t.Fatalf("write file: %v", err)
	}
	defer os.RemoveAll(repo)
	parents, err := FindParents(repo)
	if err != nil {
		t.Errorf("Error finding parents: %v", err)
	}
	if !reflect.DeepEqual(parents, []string{"foo/bar/parent1.yml", "foo/bar/parent2.yml"}) {
		t.Errorf("Bad parents: %v", parents)
	}
}

func TestFindParent(t *testing.T) {
	repo := "/tmp/neon"
	if _, err := WriteFile(repo+"/foo/bar", "parent1.yml", ""); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if _, err := WriteFile(repo+"/foo/bar", "parent2.yml", ""); err != nil {
		t.Fatalf("write file: %v", err)
	}
	defer os.RemoveAll(repo)
	parents, err := FindParent("parent1", repo)
	if err != nil {
		t.Errorf("Error finding parent: %v", err)
	}
	if !reflect.DeepEqual(parents, []string{"foo/bar/parent1.yml"}) {
		t.Errorf("Error find parent: %v", parents)
	}
}

func TestParentPath(t *testing.T) {
	repo := "/tmp/neon"
	if _, err := WriteFile(repo+"/foo/bar", "parent1.yml", ""); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if _, err := WriteFile(repo+"/foo/bar", "parent2.yml", ""); err != nil {
		t.Fatalf("write file: %v", err)
	}
	defer os.RemoveAll(repo)
	build := &Build{Repository: repo}
	path, err := build.ParentPath("parent1")
	if err != nil {
		t.Errorf("Error finding parent path: %v", err)
	}
	if path != "/tmp/neon/foo/bar/parent1.yml" {
		t.Errorf("Error finding parent path: %v", path)
	}
}

func TestInstallPlugin(t *testing.T) {
	repo := "/tmp/neon"
	if err := os.MkdirAll(repo, util.DirFileMode); err != nil {
		t.Fatalf("making directory: %v", err)
	}
	defer os.RemoveAll(repo)
	err := InstallPlugin("c4s4/build", repo)
	if err != nil {
		t.Errorf("Error installing pluging: %v", err)
	}
	path := filepath.Join(repo, "c4s4/build")
	if !util.DirExists(path) {
		t.Errorf("Plugin path not found")
	}
}

func TestFindTemplates(t *testing.T) {
	repo := "/tmp/neon"
	if _, err := WriteFile(repo+"/foo/bar", "template1.tpl", ""); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if _, err := WriteFile(repo+"/foo/bar", "template2.tpl", ""); err != nil {
		t.Fatalf("write file: %v", err)
	}
	defer os.RemoveAll(repo)
	templates, err := FindTemplates(repo)
	if err != nil {
		t.Errorf("Error finding templates: %v", err)
	}
	if !reflect.DeepEqual(templates, []string{"foo/bar/template1.tpl", "foo/bar/template2.tpl"}) {
		t.Errorf("Templates not found: %v", templates)
	}
}

func TestFindTemplate(t *testing.T) {
	repo := "/tmp/neon"
	if _, err := WriteFile(repo+"/foo/bar", "template1.tpl", ""); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if _, err := WriteFile(repo+"/foo/bar", "template2.tpl", ""); err != nil {
		t.Fatalf("write file: %v", err)
	}
	defer os.RemoveAll(repo)
	templates, err := FindTemplate("template1", repo)
	if err != nil {
		t.Errorf("Error finding template: %v", err)
	}
	if !reflect.DeepEqual(templates, []string{"foo/bar/template1.tpl"}) {
		t.Errorf("Error find template: %v", templates)
	}
}

func TestTemplatePath(t *testing.T) {
	repo := "/tmp/neon"
	if _, err := WriteFile(repo+"/foo/bar", "template1.tpl", ""); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if _, err := WriteFile(repo+"/foo/bar", "template2.tpl", ""); err != nil {
		t.Fatalf("write file: %v", err)
	}
	defer os.RemoveAll(repo)
	path, err := TemplatePath("template1", repo)
	if err != nil {
		t.Errorf("Error finding template path: %v", err)
	}
	if path != "/tmp/neon/foo/bar/template1.tpl" {
		t.Errorf("Error finding template path: %v", path)
	}
}

func TestLinkPath(t *testing.T) {
	if LinkPath("/foo", "/tmp/neon") != "/foo" {
		t.Errorf("Bad linkpath")
	}
	if LinkPath("./foo", "/tmp/neon") != "/tmp/neon/foo" {
		t.Errorf("Bad linkpath")
	}
	if LinkPath("foo", "/tmp/neon") != "/tmp/neon/foo" {
		t.Errorf("Bad linkpath")
	}
}

func TestScriptPath(t *testing.T) {
	build := &Build{
		Repository: "/tmp/repo",
		Dir:        "dir",
	}
	path, err := build.ScriptPath("/foo")
	if err != nil || path != "/foo" {
		t.Errorf("Bad script path: %s", path)
	}
	path, err = build.ScriptPath("./foo")
	if err != nil || path != "dir/foo" {
		t.Errorf("Bad script path: %s", path)
	}
	path, err = build.ScriptPath("/foo/bar/spam")
	if err != nil || path != "/foo/bar/spam" {
		t.Errorf("Bad script path: %s", path)
	}
}
