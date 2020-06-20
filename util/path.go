package util

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

var currWd string

func init() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "get working directory path error: %v", err)
	}

	currWd = wd
}

// FromRootPath returns absolute path by appending givenPath to the root of this project
func FromRootPath(givenPath string) string {
	return path.Join(currWd, givenPath)
}

// RemoveContents removes all files from given directory
func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
