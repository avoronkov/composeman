package proc

import (
	"os"
	"reflect"
	"testing"
)

func TestCreateMountedRwDirs(t *testing.T) {
	createdDirs := []string{}
	mkdirMock := func(dir string, perm os.FileMode) error {
		createdDirs = append(createdDirs, dir)
		return nil
	}
	oldMkdir := mkdir
	mkdir = mkdirMock
	defer func() {
		mkdir = oldMkdir
	}()

	volumes := []string{
		"./ro-dir:/ro-dir:ro",
		"./rw-dir1:/rw-dir1",
		"./rw-dir2:/rw-dir2:rw",
	}

	if err := createMountedRwDirs(volumes); err != nil {
		t.Fatal(err)
	}

	exp := []string{"./rw-dir1", "./rw-dir2"}
	if !reflect.DeepEqual(createdDirs, exp) {
		t.Errorf("Required dirs are not created: want %v, got %v", exp, createdDirs)
	}
}
