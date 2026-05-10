package system

import (
	"testing"
)

func TestSpecWalkthrough(t *testing.T) {
	sh := NewShell()
	if sh.Pwd() != "/" {
		t.Fatalf("pwd: got %q", sh.Pwd())
	}

	sh.Mkdir("/a/b/c")
	if sh.Pwd() != "/" {
		t.Fatalf("after mkdir, pwd: got %q", sh.Pwd())
	}

	if err := sh.Cd("a/b"); err != nil {
		t.Fatal(err)
	}
	if sh.Pwd() != "/a/b" {
		t.Fatalf("cd a/b: got %q", sh.Pwd())
	}

	if err := sh.Cd("*"); err != nil {
		t.Fatal(err)
	}
	if sh.Pwd() != "/a/b/c" {
		t.Fatalf("cd *: got %q", sh.Pwd())
	}

	if err := sh.Cd("../*"); err != nil {
		t.Fatal(err)
	}
	if sh.Pwd() != "/a/b/c" {
		t.Fatalf("cd ../*: got %q", sh.Pwd())
	}

	if err := sh.Cd("/"); err != nil {
		t.Fatal(err)
	}
	if err := sh.Cd("/*"); err != nil {
		t.Fatal(err)
	}
	if sh.Pwd() != "/a" {
		t.Fatalf("cd /*: got %q", sh.Pwd())
	}
}

func TestCdStarEmptyRoot(t *testing.T) {
	sh := NewShell()
	if err := sh.Cd("*"); err != nil {
		t.Fatal(err)
	}
	if sh.Pwd() != "/" {
		t.Fatalf("cd * with no children: got %q", sh.Pwd())
	}
}

func TestCdFailsPreservesCwd(t *testing.T) {
	sh := NewShell()
	sh.Mkdir("/a/b")
	if err := sh.Cd("/a/b"); err != nil {
		t.Fatal(err)
	}
	before := sh.Pwd()
	if err := sh.Cd("/nope/*/x"); err == nil {
		t.Fatal("expected error")
	}
	if sh.Pwd() != before {
		t.Fatalf("cwd changed on error: got %q want %q", sh.Pwd(), before)
	}
}

func TestMkdirIdempotentAndSlashes(t *testing.T) {
	sh := NewShell()
	sh.Mkdir("//a///b//")
	if err := sh.Cd("/a/b"); err != nil {
		t.Fatal(err)
	}
	sh.Mkdir("/a/b")
	if sh.Pwd() != "/a/b" {
		t.Fatal(sh.Pwd())
	}
}

func TestMkdirIgnoresFinalDotDot(t *testing.T) {
	sh := NewShell()
	sh.Mkdir("/a/b")
	sh.Mkdir("/a/b/..") // invalid final ..
	if err := sh.Cd("/a"); err != nil {
		t.Fatal(err)
	}
}

func TestWildcardLexSmallest(t *testing.T) {
	sh := NewShell()
	sh.Mkdir("/z")
	sh.Mkdir("/a")
	if err := sh.Cd("/"); err != nil {
		t.Fatal(err)
	}
	if err := sh.Cd("*"); err != nil {
		t.Fatal(err)
	}
	if sh.Pwd() != "/a" {
		t.Fatalf("want /a, got %s", sh.Pwd())
	}
}
