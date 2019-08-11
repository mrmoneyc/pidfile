package pidfile

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNewAndRemove(t *testing.T) {
	d, err := ioutil.TempDir(os.TempDir(), "pidfile-testdir")
	if err != nil {
		t.Fatal("Could not create test temp folder")
	}

	p := filepath.Join(d, "test.pid")
	f, err := New(p)
	if err != nil {
		t.Fatalf("Could not create test pidfile: %s", err.Error())
	}

	if err = f.Remove(); err != nil {
		t.Fatal("Could not remove test pidfile")
	}
}

func TestRemoveInvalidPath(t *testing.T) {
	p := filepath.Join("foo", "bar", "pidfile-test", "test.pid")
	f := PIDFile{p}

	if err := f.Remove(); err == nil {
		t.Fatalf("Non-existing file %s should cause error", p)
	}
}
