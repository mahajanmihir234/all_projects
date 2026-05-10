package system

import (
	"errors"
	"slices"
	"sort"
	"strings"
)

// ErrCdFailed is returned when cd cannot resolve the full path; the working directory is unchanged.
var ErrCdFailed = errors.New("cd: path resolution failed")

// Shell is an in-memory Unix-like filesystem with a current working directory.
type Shell struct {
	root *dir
	cwd  *dir
}

type dir struct {
	name   string
	parent *dir
	kids   map[string]*dir
}

func newDir(name string, parent *dir) *dir {
	return &dir{
		name:   name,
		parent: parent,
		kids:   make(map[string]*dir),
	}
}

// NewShell returns a shell with CWD at /.
func NewShell() *Shell {
	r := newDir("", nil)
	return &Shell{root: r, cwd: r}
}

func (d *dir) isRoot() bool {
	return d.parent == nil
}

func parentOf(d *dir) *dir {
	if d.isRoot() {
		return d
	}
	return d.parent
}

// Pwd returns the absolute path of the current directory ("/" for root, no trailing slash).
func (s *Shell) Pwd() string {
	if s.cwd.isRoot() {
		return "/"
	}
	var parts []string
	for cur := s.cwd; !cur.isRoot(); cur = cur.parent {
		parts = append(parts, cur.name)
	}
	// collected from leaf to root
	slices.Reverse(parts)
	return "/" + strings.Join(parts, "/")
}

// Mkdir creates the final directory in path, with parents as needed (mkdir -p).
// Relative paths are anchored at CWD; absolute paths at /. Creating an existing directory is a no-op.
// It is invalid for the final segment to be . or ..; such paths are ignored.
func (s *Shell) Mkdir(path string) {
	segs := splitPath(path)
	if len(segs) == 0 {
		return
	}
	last := segs[len(segs)-1]
	if last == "." || last == ".." {
		return
	}

	anchor := s.cwd
	if strings.HasPrefix(path, "/") {
		anchor = s.root
	}

	cur := walkMkdirAnchor(anchor, segs[:len(segs)-1])
	if cur.kids[last] == nil {
		cur.kids[last] = newDir(last, cur)
	}
}

func walkMkdirAnchor(start *dir, segs []string) *dir {
	cur := start
	for _, seg := range segs {
		switch seg {
		case ".":
			continue
		case "..":
			cur = parentOf(cur)
		default:
			if cur.kids[seg] == nil {
				cur.kids[seg] = newDir(seg, cur)
			}
			cur = cur.kids[seg]
		}
	}
	return cur
}

// Cd changes the current directory to path. On failure, CWD is unchanged.
func (s *Shell) Cd(path string) error {
	segs := splitPath(path)
	start := s.cwd
	if strings.HasPrefix(path, "/") {
		start = s.root
	}
	resolved, err := resolveCd(start, segs)
	if err != nil {
		return ErrCdFailed
	}
	s.cwd = resolved
	return nil
}

func resolveCd(cur *dir, segs []string) (*dir, error) {
	for _, seg := range segs {
		if seg == "*" {
			seg = pickWildcard(cur)
		}
		switch seg {
		case ".":
			continue
		case "..":
			cur = parentOf(cur)
		default:
			child, ok := cur.kids[seg]
			if !ok {
				return nil, ErrCdFailed
			}
			cur = child
		}
	}
	return cur, nil
}

// pickWildcard: lexicographically smallest child directory; if there are none, "." (.. in the
// match set is only needed for literal ".." segments; * never resolves to .. when children are absent).
func pickWildcard(cur *dir) string {
	names := make([]string, 0, len(cur.kids))
	for name := range cur.kids {
		names = append(names, name)
	}
	if len(names) > 0 {
		sort.Strings(names)
		return names[0]
	}
	return "."
}

func splitPath(path string) []string {
	return strings.Split(path, "/")
}
