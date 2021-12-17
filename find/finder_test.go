package find

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestNewFinder(t *testing.T) {
	root := "."
	want := finder{
		root:     root,
		filters:  []func(string) bool{},
		maxDepth: 100_000,
	}
	finder := NewFinder(root)

	if finder.root != want.root || len(finder.filters) != len(want.filters) ||
		finder.maxDepth != want.maxDepth {
		t.Errorf("got %+v; want %+v", finder, want)
		return
	}
}

func TestFinder_Filter(t *testing.T) {
	f := NewFinder(".")
	f.Filter(func(s string) bool {
		return true
	})
	if len(f.filters) != 1 {
		t.Errorf("expected finder to have 1 filter, found %d", len(f.filters))
	}
}

func TestFinder_Find(t *testing.T) {
	parentDir := getParentFolder()
	f := NewFinder(parentDir)
	f.Filter(func(s string) bool {
		base := filepath.Base(s)
		return strings.HasPrefix(base, "finder")
	})
	files, _ := f.Find()
	if len(files) != 2 {
		t.Errorf("expected 2 files, got %d\n", len(files))
	}
}

func TestFinder_Find_withMaxDepth(t *testing.T) {
	parentDir := getParentFolder()
	matches, _ := NewFinder(parentDir).
		MaxDepth(1).
		Find()
	if len(matches) != 2 {
		t.Errorf("expected to find 2 files, got %d", len(matches))
	}
}

func TestGetDirContents(t *testing.T) {
	parentFolder := getParentFolder()
	contents, _ := getDirContents(parentFolder)
	for _, file := range contents {
		if filepath.Base(file) == "finder_test.go" {
			return
		}
	}
	t.Error("Expected to find this file (finder_test.go) but didn't")
}

func TestIsDir(t *testing.T) {
	if !isDir(getParentFolder()) {
		t.Error("the dir of this file should be considered a dir")
	}
	notDir := "this better not be a dir"
	if isDir(notDir) {
		t.Error("got true, expected false")
	}
}

func TestFinder_StartsWith(t *testing.T) {
	matches, _ := NewFinder(getParentFolder()).
		StartsWith("finder_").
		Find()
	if filepath.Base(matches[0]) != "finder_test.go" {
		t.Error("expected 1 match and for it to be finder_test.go")
	}
}

func TestFinder_EndsWith(t *testing.T) {
	matches, _ := NewFinder(getParentFolder()).
		EndsWith("_test.go").
		Find()
	for _, file := range matches {
		if !strings.HasSuffix(file, "_test.go") {
			t.Errorf("file %s doesn't have expected suffix", file)
		}
	}
}

func TestFinder_Contains(t *testing.T) {
	matches, _ := NewFinder(getParentFolder()).
		Contains("finder").
		Find()
	if len(matches) != 2 {
		t.Errorf("expected 2 matches, got %d", len(matches))
	}
	matches, _ = NewFinder(getParentFolder()).
		Contains("fiNDer").
		Find()
	if len(matches) != 0 {
		t.Errorf("expected 0 matches, got %d", len(matches))
	}
}

func TestFinder_ContainsInsensitive(t *testing.T) {
	matches, _ := NewFinder(getParentFolder()).
		ContainsInsensitive("fiNDer").
		Find()
	if len(matches) != 2 {
		t.Errorf("expected 2 matches, got %d", len(matches))
	}
}

func TestFinder_Regex(t *testing.T) {
	matches, _ := NewFinder(getParentFolder()).
		Regex("der.{1}[go]{2}").
		Find()
	if len(matches) != 1 {
		t.Errorf("expected 1 matches, got %d", len(matches))
	}
}

func TestFinder_SizeAtLeast(t *testing.T) {
	matches, _ := NewFinder(getParentFolder()).
		SizeAtLeast(10).
		EndsWith(".go").
		Find()
	// We expect 3 files in total and for them all to be > 10 bytes.
	if len(matches) != 2 {
		t.Errorf("expected 2 matches, got %d", len(matches))
	}
}

func TestFinder_SizeAtMost(t *testing.T) {
	matches, _ := NewFinder(getParentFolder()).
		SizeAtMost(10).
		Find()
	// We expect no files to be less than 10 bytes in this dir.
	if len(matches) != 0 {
		t.Errorf("expected 0 matches, got %d", len(matches))
	}
}

// Helpers

func getParentFolder() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filename)
}
