package find

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type finder struct {
	root         string
	filters      []func(string) bool
	maxDepth     uint
	printMatches bool
}

// NewFinder returns a new finder instance that can be used to search for files starting
// at the given root directory. Returns an error if the root dir does not exist.
func NewFinder(root string) *finder {
	return &finder{
		root:         root,
		maxDepth:     100_000, // initialized arbitrarily large
		printMatches: false,
	}
}

// Filter applies a filter function to this finder such that only files that return true
// when passed to this function are kept. Note the finder is not executed until a
// terminal operator is called. i.e. find().
func (f *finder) Filter(filter func(string) bool) *finder {
	f.filters = append(f.filters, filter)
	return f
}

func (f *finder) MaxDepth(n uint) *finder {
	f.maxDepth = n
	return f
}

func (f *finder) StartsWith(prefix string) *finder {
	f.Filter(func(filename string) bool {
		return strings.HasPrefix(filepath.Base(filename), prefix)
	})
	return f
}

func (f *finder) EndsWith(suffix string) *finder {
	f.Filter(func(filename string) bool {
		return strings.HasSuffix(filepath.Base(filename), suffix)
	})
	return f
}

func (f *finder) Contains(s string) *finder {
	f.Filter(func(filename string) bool {
		return strings.Contains(filepath.Base(filename), s)
	})
	return f
}

func (f *finder) DoesntContain(s string) *finder {
	f.Filter(func(file string) bool {
		return !strings.Contains(filepath.Base(file), s)
	})
	return f
}

func (f *finder) ContainsInsensitive(s string) *finder {
	f.Filter(func(file string) bool {
		return strings.Contains(strings.ToLower(file), strings.ToLower(s))
	})
	return f
}

func (f *finder) DoesntContainsInsensitive(s string) *finder {
	f.Filter(func(file string) bool {
		return !strings.Contains(strings.ToLower(file), strings.ToLower(s))
	})
	return f
}

func (f *finder) Regex(pattern string) *finder {
	re := regexp.MustCompile(pattern)
	f.Filter(func(file string) bool {
		base := filepath.Base(file)
		return re.Match([]byte(base))
	})
	return f
}

func (f *finder) SizeAtLeast(n int) *finder {
	closure := func(file string) bool {
		info, err := os.Stat(file)
		if err != nil {
			return false
		}
		return int(info.Size()) >= n
	}
	f.Filter(closure)
	return f
}

func (f *finder) SizeAtMost(n int) *finder {
	f.Filter(func(file string) bool {
		info, err := os.Stat(file)
		if err != nil {
			return false
		}
		return int(info.Size()) <= n
	})
	return f
}

// TestFile returns true if all the filters in the finder return true for the given file.
func (f *finder) TestFile(file string) bool {
	for _, filter := range f.filters {
		if !filter(file) {
			return false
		}
	}
	return true
}

func (f *finder) PrintMatches() *finder {
	f.printMatches = true
	return f
}

func (f *finder) Find() ([]string, error) {
	if !isDir(f.root) {
		return nil, errors.New(fmt.Sprintf("invalid search root dir %s", f.root))
	}

	// Perform DFS
	q := stringQueue{}
	q.add(f.root)
	var depth uint = 0
	var ret []string

	for len(q) > 0 && depth <= f.maxDepth {

		layerKids := len(q) // nodes in next BFS depth layer
		for i := 0; i < layerKids; i++ {
			item, _ := q.poll()
			if isDir(item) {
				children, _ := getDirContents(item)
				for _, child := range children {
					q.add(child)
				}
			} else {
				if f.TestFile(item) {
					ret = append(ret, item)
					if f.printMatches {
						fmt.Println(item)
					}
				}
			}
		}

		depth++
	}
	return ret, nil
}

func getDirContents(dir string) ([]string, error) {
	var ret []string
	f, err := os.Open(dir)
	if err != nil {
		return ret, err
	}
	names, err := f.Readdirnames(0)
	for i, name := range names {
		names[i] = filepath.Join(dir, name)
	}
	return names, err
}

func isDir(file string) bool {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

// stringQueue is a helpful queue implementation for strings.
type stringQueue []string

func (q *stringQueue) add(s string) {
	*q = append(*q, s)
}

func (q *stringQueue) poll() (string, bool) {
	if len(*q) == 0 {
		return "", false
	}
	item := (*q)[0]
	*q = (*q)[1:]
	return item, true
}
