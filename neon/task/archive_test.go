package task

import (
	t "archive/tar"
	z "archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var unsafeNames = []string{"../evil.txt", "foo/../../evil.txt", "/tmp/evil.txt"}

func TestSafeJoin(test *testing.T) {
	path, err := SafeJoin("dir", "foo/bar.txt")
	if err != nil || path != filepath.Join("dir", "foo", "bar.txt") {
		test.Errorf("bad safe join: '%s' (%v)", path, err)
	}
	for _, name := range unsafeNames {
		if _, err := SafeJoin("dir", name); err == nil {
			test.Errorf("unsafe name '%s' should be rejected", name)
		}
	}
}

func TestUntarRejectsUnsafePath(test *testing.T) {
	for _, name := range unsafeNames {
		tmp := test.TempDir()
		archive := filepath.Join(tmp, "evil.tar")
		file, err := os.Create(archive)
		if err != nil {
			test.Fatal(err)
		}
		writer := t.NewWriter(file)
		content := "evil"
		if err := writer.WriteHeader(&t.Header{Name: name, Mode: 0644, Size: int64(len(content)), Typeflag: t.TypeReg}); err != nil {
			test.Fatal(err)
		}
		if _, err := writer.Write([]byte(content)); err != nil {
			test.Fatal(err)
		}
		if err := writer.Close(); err != nil {
			test.Fatal(err)
		}
		if err := file.Close(); err != nil {
			test.Fatal(err)
		}
		dir := filepath.Join(tmp, "a", "b")
		err = untarFile(archive, dir)
		if err == nil || !strings.Contains(err.Error(), "illegal path") {
			test.Errorf("untar of '%s' should fail with illegal path, got: %v", name, err)
		}
		if _, err := os.Stat(filepath.Join(tmp, "evil.txt")); err == nil {
			test.Errorf("untar of '%s' wrote outside destination directory", name)
		}
	}
}

func TestUnzipRejectsUnsafePath(test *testing.T) {
	for _, name := range unsafeNames {
		tmp := test.TempDir()
		archive := filepath.Join(tmp, "evil.zip")
		file, err := os.Create(archive)
		if err != nil {
			test.Fatal(err)
		}
		writer := z.NewWriter(file)
		entry, err := writer.CreateHeader(&z.FileHeader{Name: name})
		if err != nil {
			test.Fatal(err)
		}
		if _, err := entry.Write([]byte("evil")); err != nil {
			test.Fatal(err)
		}
		if err := writer.Close(); err != nil {
			test.Fatal(err)
		}
		if err := file.Close(); err != nil {
			test.Fatal(err)
		}
		dir := filepath.Join(tmp, "a", "b")
		err = unzipFile(archive, dir)
		if err == nil || !strings.Contains(err.Error(), "illegal path") {
			test.Errorf("unzip of '%s' should fail with illegal path, got: %v", name, err)
		}
		if _, err := os.Stat(filepath.Join(tmp, "evil.txt")); err == nil {
			test.Errorf("unzip of '%s' wrote outside destination directory", name)
		}
	}
}

// writeTestTar writes a tar archive with given entries
func writeTestTar(test *testing.T, archive string, headers []*t.Header) {
	file, err := os.Create(archive)
	if err != nil {
		test.Fatal(err)
	}
	writer := t.NewWriter(file)
	for _, header := range headers {
		if err := writer.WriteHeader(header); err != nil {
			test.Fatal(err)
		}
		if header.Size > 0 {
			if _, err := writer.Write([]byte(strings.Repeat("x", int(header.Size)))); err != nil {
				test.Fatal(err)
			}
		}
	}
	if err := writer.Close(); err != nil {
		test.Fatal(err)
	}
	if err := file.Close(); err != nil {
		test.Fatal(err)
	}
}

func TestUntarLinks(test *testing.T) {
	tmp := test.TempDir()
	archive := filepath.Join(tmp, "links.tar")
	writeTestTar(test, archive, []*t.Header{
		{Name: "file.txt", Mode: 0644, Size: 4, Typeflag: t.TypeReg},
		{Name: "symlink.txt", Linkname: "file.txt", Typeflag: t.TypeSymlink},
		{Name: "hardlink.txt", Linkname: "file.txt", Typeflag: t.TypeLink},
	})
	dir := filepath.Join(tmp, "out")
	if err := untarFile(archive, dir); err != nil {
		test.Fatalf("untar: %v", err)
	}
	if target, err := os.Readlink(filepath.Join(dir, "symlink.txt")); err != nil || target != "file.txt" {
		test.Errorf("symbolic link should be created: %s (%v)", target, err)
	}
	if content, err := os.ReadFile(filepath.Join(dir, "hardlink.txt")); err != nil || string(content) != "xxxx" {
		test.Errorf("hard link should be created: %q (%v)", content, err)
	}
}

func TestUntarDoesntWriteThroughSymlink(test *testing.T) {
	tmp := test.TempDir()
	outside := filepath.Join(tmp, "outside")
	if err := os.MkdirAll(outside, 0755); err != nil {
		test.Fatal(err)
	}
	// symbolic link in archive pointing outside, then a file written through it
	archive := filepath.Join(tmp, "evil.tar")
	writeTestTar(test, archive, []*t.Header{
		{Name: "link", Linkname: outside, Typeflag: t.TypeSymlink},
		{Name: "link/evil.txt", Mode: 0644, Size: 4, Typeflag: t.TypeReg},
	})
	if err := untarFile(archive, filepath.Join(tmp, "out1")); err == nil {
		test.Errorf("writing through a symbolic link going outside should fail")
	}
	// existing symbolic link in destination pointing outside
	dir := filepath.Join(tmp, "out2")
	if err := os.MkdirAll(dir, 0755); err != nil {
		test.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(dir, "link")); err != nil {
		test.Fatal(err)
	}
	writeTestTar(test, archive, []*t.Header{
		{Name: "link/evil.txt", Mode: 0644, Size: 4, Typeflag: t.TypeReg},
	})
	if err := untarFile(archive, dir); err == nil {
		test.Errorf("writing through an existing symbolic link going outside should fail")
	}
	if _, err := os.Stat(filepath.Join(outside, "evil.txt")); err == nil {
		test.Errorf("file was written outside destination directory")
	}
	// hard link to a file outside is rejected
	writeTestTar(test, archive, []*t.Header{
		{Name: "hard", Linkname: "../outside/file", Typeflag: t.TypeLink},
	})
	if err := untarFile(archive, filepath.Join(tmp, "out3")); err == nil || !strings.Contains(err.Error(), "illegal path") {
		test.Errorf("hard link outside should be rejected: %v", err)
	}
}

func TestUnzipSymlink(test *testing.T) {
	tmp := test.TempDir()
	archive := filepath.Join(tmp, "link.zip")
	file, err := os.Create(archive)
	if err != nil {
		test.Fatal(err)
	}
	writer := z.NewWriter(file)
	header := &z.FileHeader{Name: "link.txt"}
	header.SetMode(os.ModeSymlink | 0777)
	entry, err := writer.CreateHeader(header)
	if err != nil {
		test.Fatal(err)
	}
	if _, err := entry.Write([]byte("target.txt")); err != nil {
		test.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		test.Fatal(err)
	}
	if err := file.Close(); err != nil {
		test.Fatal(err)
	}
	dir := filepath.Join(tmp, "out")
	if err := unzipFile(archive, dir); err != nil {
		test.Fatalf("unzip: %v", err)
	}
	if target, err := os.Readlink(filepath.Join(dir, "link.txt")); err != nil || target != "target.txt" {
		test.Errorf("symbolic link should be created: %s (%v)", target, err)
	}
}
